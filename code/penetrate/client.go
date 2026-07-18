package penetrate

import (
	"fmt"
	"my-base/code/message"
	transport "my-base/code/transport"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/logger"
)

// Client 客户端链接信息
type Client struct {
	lock              sync.RWMutex
	symbol            string
	conn              *transport.Conn
	Listens           map[string]*Listen
	pingSeq           atomic.Uint64
	heartbeatMu       sync.Mutex
	pendingPings      map[uint64]time.Time
	heartbeatFailures int
}

func NewClient(conn *transport.Conn) *Client {
	return &Client{
		conn:         conn,
		Listens:      make(map[string]*Listen),
		pendingPings: make(map[uint64]time.Time),
	}
}

// NewListen 新建监听
func (c *Client) NewListen(addr, port string) {
	listen := NewListen(addr, port).SetNotify(c.TakeOver)
	c.AddListen(addr, listen)
	defer c.DelListen(addr)
	defer listen.Stop()
	err := listen.Start()
	if err != nil {
		logger.Error("client listen close ", err, addr, port)
	}
}

// GetListen 获取客户端监听端口
func (c *Client) GetListen(addr string) *Listen {
	c.lock.RLock()
	listen := c.Listens[addr]
	c.lock.RUnlock()
	return listen
}

// AddListen 获取客户端监听端口
func (c *Client) AddListen(key string, listen *Listen) *Client {
	c.lock.Lock()
	c.Listens[key] = listen
	c.lock.Unlock()
	return c
}

// IsExist 是否存在
func (c *Client) IsExist(key string) bool {
	c.lock.RLock()
	_, ok := c.Listens[key]
	c.lock.RUnlock()
	return ok
}

// DelListen 删除监听端口
func (c *Client) DelListen(key string, keys ...string) *Client {
	c.lock.Lock()
	delete(c.Listens, key)
	for _, k := range keys {
		delete(c.Listens, k)
	}
	c.lock.Unlock()
	return c
}

func (c *Client) DelListenIf(key string, listen *Listen) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	current, ok := c.Listens[key]
	if !ok || current != listen {
		return false
	}
	delete(c.Listens, key)
	return true
}

// TakeOver 监听器发送信息方法
func (c *Client) TakeOver(msg message.Message) error {
	switch {
	case msg.EqLink():
		msg.Symbol = c.symbol
		if c.conn == nil || c.conn.IsClose() {
			return net.ErrClosed
		}
		return c.conn.Send(msg)
	case msg.EqClose():
		symbol := msg.GetSymbol()
		c.DelListen(symbol)
	}
	return nil
}

// SetConn 设置主链接
func (c *Client) SetConn(conn *transport.Conn) *Client {
	c.conn = conn
	return c
}

// GetSymbol 获取标识
func (c *Client) GetSymbol() string {
	return c.symbol
}

// SetSymbol 设置标识
func (c *Client) SetSymbol(symbol string) *Client {
	c.symbol = symbol
	return c
}

// Send 发送信息
func (c *Client) Send(msg message.Message) error {
	return c.conn.Send(msg)
}

// Ping 心跳检测
func (c *Client) Ping(duration time.Duration) error {
	ch := time.NewTicker(duration)
	for {
		seq := c.pingSeq.Add(1)
		if err := c.conn.Send(message.Message{
			Type:   message.MsgTypePing,
			Symbol: c.symbol,
			Msg: message.PingPayload{
				Seq:      seq,
				SentAtMs: time.Now().UnixMilli(),
			},
		}); err != nil {
			return err
		}
		<-ch.C
	}
}

// PingLoop 带停止信号的心跳检测
func (c *Client) PingLoop(duration time.Duration, stop <-chan struct{}) error {
	return c.PingLoopWithConfig(duration, 3*duration, 3, stop)
}

// PingLoopWithConfig sends heartbeat messages and fails the control
// connection after the configured number of unanswered ping timeouts.
func (c *Client) PingLoopWithConfig(interval, pongTimeout time.Duration, maxFailures int, stop <-chan struct{}) error {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	if pongTimeout <= 0 {
		pongTimeout = 3 * interval
	}
	if maxFailures <= 0 {
		maxFailures = 3
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		seq := c.pingSeq.Add(1)
		sentAt := time.Now()
		c.heartbeatMu.Lock()
		c.pendingPings[seq] = sentAt
		c.heartbeatMu.Unlock()
		if err := c.conn.Send(message.Message{
			Type:   message.MsgTypePing,
			Symbol: c.symbol,
			Msg: message.PingPayload{
				Seq:      seq,
				SentAtMs: sentAt.UnixMilli(),
			},
		}); err != nil {
			return err
		}

		select {
		case <-ticker.C:
			if c.expiredPings(time.Now(), pongTimeout) >= maxFailures {
				return fmt.Errorf("heartbeat timeout: %d failures", maxFailures)
			}
		case <-stop:
			return nil
		}
	}
}

func (c *Client) expiredPings(now time.Time, timeout time.Duration) int {
	c.heartbeatMu.Lock()
	defer c.heartbeatMu.Unlock()
	count := 0
	for seq, sentAt := range c.pendingPings {
		if now.Sub(sentAt) >= timeout {
			count++
			delete(c.pendingPings, seq)
		}
	}
	if count > 0 {
		c.heartbeatFailures += count
	}
	return c.heartbeatFailures
}

// RecordPong removes a pending heartbeat and returns the measured RTT.
func (c *Client) RecordPong(seq uint64) time.Duration {
	c.heartbeatMu.Lock()
	defer c.heartbeatMu.Unlock()
	sentAt, ok := c.pendingPings[seq]
	if !ok {
		return -1
	}
	delete(c.pendingPings, seq)
	c.heartbeatFailures = 0
	return time.Since(sentAt)
}

// Close 关闭资源
func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}

	c.lock.Lock()
	listens := make([]*Listen, 0, len(c.Listens))
	for _, listen := range c.Listens {
		if listen != nil {
			listens = append(listens, listen)
		}
	}
	c.Listens = make(map[string]*Listen)
	c.lock.Unlock()
	for _, listen := range listens {
		if err := listen.Stop(); err != nil {
			return err
		}
	}
	return nil
}
