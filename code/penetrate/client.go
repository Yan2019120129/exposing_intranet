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

// Client 保存一个在线客户端的控制连接、端口监听器和心跳状态。
type Client struct {
	// lock 保护客户端连接、监听器和关闭状态。
	lock sync.RWMutex
	// symbol 是客户端唯一标识。
	symbol string
	// conn 是客户端的控制连接。
	conn *transport.Conn
	// Listens 按公网端口保存客户端拥有的监听器。
	Listens map[string]*Listen
	// closed 表示客户端是否已关闭，受 lock 保护。
	closed bool
	// pingSeq 是递增的心跳序号。
	pingSeq atomic.Uint64
	// heartbeatMu 保护未响应心跳及失败计数。
	heartbeatMu sync.Mutex
	// pendingPings 记录尚未收到响应的心跳发送时间。
	pendingPings map[uint64]time.Time
	// heartbeatFailures 是连续心跳超时次数。
	heartbeatFailures int
}

// NewClient 创建客户端连接状态。
func NewClient(conn *transport.Conn) *Client {
	return &Client{
		conn:         conn,
		Listens:      make(map[string]*Listen),
		pendingPings: make(map[uint64]time.Time),
	}
}

// NewListen 创建并运行客户端拥有的公网端口监听器。
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

// GetListen 按公网端口获取监听器。
func (c *Client) GetListen(addr string) *Listen {
	c.lock.RLock()
	listen := c.Listens[addr]
	c.lock.RUnlock()
	return listen
}

// AddListen 登记监听器；客户端已关闭时立即停止该监听器。
func (c *Client) AddListen(key string, listen *Listen) *Client {
	c.lock.Lock()
	if c.closed {
		c.lock.Unlock()
		if listen != nil {
			_ = listen.Stop()
		}
		return c
	}
	c.Listens[key] = listen
	c.lock.Unlock()
	return c
}

// IsExist 返回指定公网端口是否已登记监听器。
func (c *Client) IsExist(key string) bool {
	c.lock.RLock()
	_, ok := c.Listens[key]
	c.lock.RUnlock()
	return ok
}

// DelListen 删除一个主监听器及可选的关联监听器记录。
func (c *Client) DelListen(key string, keys ...string) *Client {
	c.lock.Lock()
	delete(c.Listens, key)
	for _, k := range keys {
		delete(c.Listens, k)
	}
	c.lock.Unlock()
	return c
}

// DelListenIf 仅当当前端口仍关联指定监听器时删除其记录。
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

// TakeOver 将监听器事件转发到客户端控制连接。
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

// SetConn 设置客户端控制连接。
func (c *Client) SetConn(conn *transport.Conn) *Client {
	c.conn = conn
	return c
}

// GetSymbol 返回客户端唯一标识。
func (c *Client) GetSymbol() string {
	return c.symbol
}

// SetSymbol 设置客户端唯一标识。
func (c *Client) SetSymbol(symbol string) *Client {
	c.symbol = symbol
	return c
}

// Send 通过客户端控制连接发送消息。
func (c *Client) Send(msg message.Message) error {
	return c.conn.Send(msg)
}

// SendWithTimeout 在指定的写超时时间内发送控制消息。
func (c *Client) SendWithTimeout(msg message.Message, timeout time.Duration) error {
	return c.conn.SendWithTimeout(msg, timeout)
}

// Ping 按指定间隔持续发送心跳，直到发送失败。
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

// PingLoop 使用默认超时策略运行可停止的心跳循环。
func (c *Client) PingLoop(duration time.Duration, stop <-chan struct{}) error {
	return c.PingLoopWithConfig(duration, 3*duration, 3, stop)
}

// PingLoopWithConfig 按配置发送心跳，并在连续超时达到阈值后返回错误。
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

// expiredPings 清理超时心跳并返回累计失败次数。
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

// RecordPong 移除已响应心跳并返回测得的往返时延。
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

// Close 关闭控制连接和客户端拥有的全部监听器。
func (c *Client) Close() error {
	var closeErr error
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			closeErr = err
		}
	}

	c.lock.Lock()
	c.closed = true
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
			if closeErr == nil {
				closeErr = err
			}
		}
	}
	return closeErr
}
