package transport

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/user"
	"strings"
	"sync/atomic"
	"time"

	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
)

var (
	ErrClientDeleted              = errors.New("client deleted")
	ErrClientRegistrationRejected = errors.New("client registration rejected")
)

type ConnOptions = toolsPenetrate.ConnOptions

// Client is the client-side control connection and mapping coordinator.
type Client struct {
	addr            string
	OnKey           func(string)
	GetKey          func() string
	conn            *toolsPenetrate.Conn
	connOptions     toolsPenetrate.ConnOptions
	dataTimeout     time.Duration
	lastControlSeen atomic.Int64
}

const defaultDataTimeout = 300 * time.Second

func NewClient(addr string) *Client { return &Client{addr: addr, dataTimeout: defaultDataTimeout} }

func (c *Client) SetConnOptions(options toolsPenetrate.ConnOptions) *Client {
	c.connOptions = options
	return c
}

// SetDataTimeout sets the idle read/write timeout for data mappings.
func (c *Client) SetDataTimeout(timeout time.Duration) *Client {
	if timeout > 0 {
		c.dataTimeout = timeout
	}
	return c
}

func (c *Client) Link() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = toolsPenetrate.NewConnWithOptions(conn, c.connOptions)
	key := ""
	if c.GetKey != nil {
		key = c.GetKey()
	}
	if err := c.conn.Send(message.Message{
		Symbol: key,
		Name:   c.GetName(),
		Type:   message.MsgTypeRegister,
	}); err != nil {
		_ = c.conn.Close()
		return err
	}
	return nil
}

func (c *Client) Start() error { return c.StartWithHeartbeat(0, 0, 0) }

func (c *Client) StartWithHeartbeat(interval, pongTimeout time.Duration, maxFailures int) error {
	if err := c.Link(); err != nil {
		return err
	}
	defer c.Close()
	log.Println("connection succeeded")
	c.lastControlSeen.Store(time.Now().UnixNano())
	stopWatch := make(chan struct{})
	defer close(stopWatch)
	if interval > 0 {
		go c.watchHeartbeat(interval, pongTimeout, maxFailures, stopWatch)
	}

	for {
		var msg message.Message
		if err := c.conn.ParseMsg(&msg); err != nil {
			return err
		}
		c.lastControlSeen.Store(time.Now().UnixNano())
		switch {
		case msg.EqRegister():
			if c.OnKey != nil {
				c.OnKey(msg.Symbol)
			}
		case msg.EqPing():
			var ping message.PingPayload
			if value, ok := msg.Msg.(map[string]any); ok {
				if seq, ok := value["seq"].(float64); ok {
					ping.Seq = uint64(seq)
				}
				if sentAt, ok := value["sentAtMs"].(float64); ok {
					ping.SentAtMs = int64(sentAt)
				}
			}
			msg.Type = message.MsgTypePong
			msg.Msg = message.PongPayload{
				Seq: ping.Seq, SentAtMs: ping.SentAtMs, ClientAtMs: time.Now().UnixMilli(),
			}
			if err := c.conn.Send(msg); err != nil {
				return err
			}
		case msg.EqLink():
			go c.NewMapping(msg)
		case msg.EqDel():
			if c.OnKey != nil {
				c.OnKey("")
			}
			return ErrClientDeleted
		case msg.EqClose():
			return fmt.Errorf("%w: %v", ErrClientRegistrationRejected, msg.Msg)
		default:
			log.Println(msg.Msg)
		}
	}
}

func (c *Client) watchHeartbeat(interval, pongTimeout time.Duration, maxFailures int, stop <-chan struct{}) {
	if pongTimeout <= 0 {
		pongTimeout = 3 * interval
	}
	if maxFailures <= 0 {
		maxFailures = 3
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			last := time.Unix(0, c.lastControlSeen.Load())
			if time.Since(last) >= time.Duration(maxFailures)*(interval+pongTimeout) {
				_ = c.Close()
				return
			}
		case <-stop:
			return
		}
	}
}

func (c *Client) Send(msg message.Message) error {
	if c.conn == nil {
		return net.ErrClosed
	}
	return c.conn.Send(msg)
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) GetIp() string {
	host, _, err := net.SplitHostPort(c.addr)
	if err == nil {
		return host
	}
	return strings.TrimSuffix(c.addr, ":")
}

func (c *Client) GetName() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	return currentUser.Username
}

func (c *Client) SetOnKeyFunc(fn func(string)) *Client {
	c.OnKey = fn
	return c
}

func (c *Client) SetKeyFunc(fn func() string) *Client {
	c.GetKey = fn
	return c
}

func (c *Client) Copy(targetConn, conn net.Conn) {
	defer conn.Close()
	defer targetConn.Close()
	results := make(chan error, 2)
	go func() { results <- copyHalfClose(targetConn, conn, c.dataTimeout) }()
	go func() { results <- copyHalfClose(conn, targetConn, c.dataTimeout) }()
	first := <-results
	if !normalCopyEnd(first) {
		_ = conn.Close()
		_ = targetConn.Close()
	}
	second := <-results
	if !normalCopyEnd(first) {
		log.Println("mapping copy error:", first)
	} else if !normalCopyEnd(second) {
		log.Println("mapping copy error:", second)
	}
}

func copyHalfClose(dst, src net.Conn, timeout time.Duration) error {
	var reader io.Reader = src
	var writer io.Writer = dst
	if timeout > 0 {
		reader = deadlineReader{conn: src, timeout: timeout}
		writer = deadlineWriter{conn: dst, timeout: timeout}
	}
	_, err := io.Copy(writer, reader)
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		_ = tcpConn.CloseWrite()
	}
	return err
}

type deadlineReader struct {
	conn    net.Conn
	timeout time.Duration
}

func (r deadlineReader) Read(p []byte) (int, error) {
	if err := r.conn.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
		return 0, err
	}
	return r.conn.Read(p)
}

type deadlineWriter struct {
	conn    net.Conn
	timeout time.Duration
}

func (w deadlineWriter) Write(p []byte) (int, error) {
	if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, err
	}
	return w.conn.Write(p)
}

func normalCopyEnd(err error) bool {
	return err == nil || errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed)
}
