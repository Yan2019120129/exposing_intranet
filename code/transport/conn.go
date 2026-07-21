package transport

import (
	"encoding/json"
	"my-base/code/message"
	"net"
	"sync"
	"time"
)

type Status int

const (
	ConnStatusClose  Status = -1   // 关闭
	ConnStatusIdle   Status = iota // 空闲
	ConnStatusActive               // 活跃
	ConnStatusWait                 // 待连接
)

const (
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 10 * time.Second
)

type Conn struct {
	Err             error
	conn            net.Conn
	status          Status
	symbol          string
	decoder         *json.Decoder
	readMu          sync.Mutex
	writeMu         sync.Mutex
	stateMu         sync.RWMutex
	closeOnce       sync.Once
	closeErr        error
	readTimeout     time.Duration
	writeTimeout    time.Duration
	keepAlive       bool
	keepAlivePeriod time.Duration
}

// ConnOptions controls control-plane connection behavior. A zero timeout
// disables the corresponding deadline. TCP options are best-effort and are
// applied only when the underlying connection is a *net.TCPConn.
type ConnOptions struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	KeepAlive       bool
	KeepAlivePeriod time.Duration
}

// NewConn 设置连接
func NewConn(conn net.Conn) *Conn {
	return NewConnWithOptions(conn, ConnOptions{
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	})
}

// NewConnWithOptions creates a connection with explicit timeout and TCP
// keepalive settings while preserving the existing JSON wire format.
func NewConnWithOptions(conn net.Conn, options ConnOptions) *Conn {
	c := &Conn{
		conn:            conn,
		status:          ConnStatusIdle,
		decoder:         json.NewDecoder(conn),
		readTimeout:     options.ReadTimeout,
		writeTimeout:    options.WriteTimeout,
		keepAlive:       options.KeepAlive,
		keepAlivePeriod: options.KeepAlivePeriod,
	}
	c.applyTCPOptions()
	return c
}

func (c *Conn) applyTCPOptions() {
	tcpConn, ok := c.conn.(*net.TCPConn)
	if !ok || !c.keepAlive {
		return
	}
	_ = tcpConn.SetKeepAlive(true)
	if c.keepAlivePeriod > 0 {
		_ = tcpConn.SetKeepAlivePeriod(c.keepAlivePeriod)
	}
}

// GetConn 获取链接
func (c *Conn) GetConn() net.Conn {
	return c.conn
}

// IsIdle 是否空闲
func (c *Conn) IsIdle() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.status == ConnStatusIdle
}

// IsWait 是否待连接状态
func (c *Conn) IsWait() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.status == ConnStatusWait
}

// IsClose 是否关闭
func (c *Conn) IsClose() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.status == ConnStatusClose
}

// Status 连接状态
func (c *Conn) Status() Status {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.status
}

// IsActive 是否活跃
func (c *Conn) IsActive() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.status == ConnStatusActive
}

// SetStatusActive 设置活跃状态
func (c *Conn) SetStatusActive() *Conn {
	c.stateMu.Lock()
	c.status = ConnStatusActive
	c.stateMu.Unlock()
	return c
}

// SetStatusWait 设置待连接状态
func (c *Conn) SetStatusWait() *Conn {
	c.stateMu.Lock()
	c.status = ConnStatusWait
	c.stateMu.Unlock()
	return c
}

// SetStatusClose 设置关闭状态
func (c *Conn) SetStatusClose() *Conn {
	c.stateMu.Lock()
	c.status = ConnStatusClose
	c.stateMu.Unlock()
	return c
}

// ParseMsg 解析数据（带超时）
func (c *Conn) ParseMsg(msg any) error {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	if c.readTimeout > 0 {
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return err
		}
	}
	return c.decoder.Decode(msg)
}

// ParseMsgNoTimeout 无超时解析（用于主连接的消息循环）
func (c *Conn) ParseMsgNoTimeout(msg any) error {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	if err := c.conn.SetReadDeadline(time.Time{}); err != nil {
		return err
	}
	return c.decoder.Decode(msg)
}

// SetSymbol 设置标识
func (c *Conn) SetSymbol(symbol string) *Conn {
	c.symbol = symbol
	return c
}

// GetSymbol 获取标识
func (c *Conn) GetSymbol() string {
	return c.symbol
}

// Send 发送信息（带超时）
func (c *Conn) Send(msg message.Message) error {
	return c.SendWithTimeout(msg, c.writeTimeout)
}

// SendWithTimeout 使用调用方指定的写超时时间发送一条消息。
func (c *Conn) SendWithTimeout(msg message.Message, timeout time.Duration) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	if err := c.conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	encoder := json.NewEncoder(c.conn)
	return encoder.Encode(&msg)
}

// SendNoTimeout 无超时发送
func (c *Conn) SendNoTimeout(msg message.Message) error {
	return c.SendWithTimeout(msg, 0)
}

// Close 关闭连接
func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		c.SetStatusClose()
		c.closeErr = c.conn.Close()
	})
	return c.closeErr
}
