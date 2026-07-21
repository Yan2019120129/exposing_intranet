package transport

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/user"
	"strings"
	"sync"
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
	addr        string
	OnKey       func(string)
	GetKey      func() string
	connMu      sync.RWMutex
	conn        *toolsPenetrate.Conn
	connOptions toolsPenetrate.ConnOptions
	dataTimeout time.Duration
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

func (c *Client) link() (*toolsPenetrate.Conn, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}
	control := toolsPenetrate.NewConnWithOptions(conn, c.connOptions)
	key := ""
	if c.GetKey != nil {
		key = c.GetKey()
	}
	if err := control.Send(message.Message{
		Symbol: key,
		Name:   c.GetName(),
		Type:   message.MsgTypeRegister,
	}); err != nil {
		_ = control.Close()
		return nil, err
	}

	c.connMu.Lock()
	previous := c.conn
	c.conn = control
	c.connMu.Unlock()
	if previous != nil && previous != control {
		_ = previous.Close()
	}
	return control, nil
}

func (c *Client) Link() error {
	_, err := c.link()
	return err
}

func (c *Client) Start() error { return c.StartWithHeartbeat(0, 0, 0) }

func (c *Client) StartWithHeartbeat(interval, pongTimeout time.Duration, maxFailures int) error {
	conn, err := c.link()
	if err != nil {
		return err
	}
	log.Println("connection succeeded")
	lastControlSeen := &atomic.Int64{}
	lastControlSeen.Store(time.Now().UnixNano())
	stopWatch := make(chan struct{})
	watchDone := make(chan struct{})
	if interval > 0 {
		go func() {
			defer close(watchDone)
			c.watchHeartbeat(conn, lastControlSeen, interval, pongTimeout, maxFailures, stopWatch)
		}()
	} else {
		close(watchDone)
	}
	defer func() {
		close(stopWatch)
		<-watchDone
		_ = conn.Close()
		c.clearConnIf(conn)
	}()

	for {
		var msg message.Message
		if err := conn.ParseMsg(&msg); err != nil {
			return err
		}
		lastControlSeen.Store(time.Now().UnixNano())
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
			if err := conn.Send(msg); err != nil {
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

func heartbeatTimeout(interval, pongTimeout time.Duration, maxFailures int) time.Duration {
	if pongTimeout <= 0 {
		pongTimeout = 3 * interval
	}
	if maxFailures <= 0 {
		maxFailures = 3
	}
	return pongTimeout + time.Duration(maxFailures-1)*interval
}

func (c *Client) watchHeartbeat(conn *toolsPenetrate.Conn, lastControlSeen *atomic.Int64, interval, pongTimeout time.Duration, maxFailures int, stop <-chan struct{}) {
	timeout := heartbeatTimeout(interval, pongTimeout, maxFailures)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			last := time.Unix(0, lastControlSeen.Load())
			if time.Since(last) >= timeout {
				_ = conn.Close()
				return
			}
		case <-stop:
			return
		}
	}
}

func (c *Client) Send(msg message.Message) error {
	conn := c.currentConn()
	if conn == nil {
		return net.ErrClosed
	}
	return conn.Send(msg)
}

func (c *Client) Close() error {
	conn := c.currentConn()
	if conn == nil {
		return nil
	}
	return conn.Close()
}

func (c *Client) currentConn() *toolsPenetrate.Conn {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.conn
}

func (c *Client) clearConnIf(conn *toolsPenetrate.Conn) {
	c.connMu.Lock()
	if c.conn == conn {
		c.conn = nil
	}
	c.connMu.Unlock()
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
