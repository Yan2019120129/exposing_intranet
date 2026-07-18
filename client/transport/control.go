package transport

import (
	"errors"
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

var ErrClientDeleted = errors.New("client deleted")

type ConnOptions = toolsPenetrate.ConnOptions

// Client is the client-side control connection and mapping coordinator.
type Client struct {
	addr            string
	OnKey           func(string)
	GetKey          func() string
	conn            *toolsPenetrate.Conn
	connOptions     toolsPenetrate.ConnOptions
	lastControlSeen atomic.Int64
}

func NewClient(addr string) *Client { return &Client{addr: addr} }

func (c *Client) SetConnOptions(options toolsPenetrate.ConnOptions) *Client {
	c.connOptions = options
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
	go func() { results <- copyHalfClose(targetConn, conn) }()
	go func() { results <- copyHalfClose(conn, targetConn) }()
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

func copyHalfClose(dst, src net.Conn) error {
	_, err := io.Copy(dst, src)
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		_ = tcpConn.CloseWrite()
	}
	return err
}

func normalCopyEnd(err error) bool {
	return err == nil || errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed)
}
