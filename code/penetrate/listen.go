package penetrate

import (
	"my-base/code/message"
	transport "my-base/code/transport"
	"net"
	"sync"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/logger"
)

// Listen 管理一个公网端口监听器及其等待配对的数据连接。
type Listen struct {
	// lock 保护监听器状态和连接索引。
	lock sync.RWMutex
	// addr 是公网监听地址。
	addr string
	// mapPort 是客户端本地映射地址。
	mapPort string
	// Notify 用于将新公网连接通知给所属客户端。
	Notify func(msg message.Message) error
	// listen 是底层网络监听器。
	listen net.Listener
	// connList 按连接标识保存等待或正在转发的数据连接。
	connList map[string]*transport.Conn
	// waitSince 记录等待配对连接的开始时间。
	waitSince map[string]time.Time
	// waitTimeout 是等待配对连接的最长等待时间。
	waitTimeout time.Duration
	// cleanupInterval 是等待连接清理任务的执行间隔。
	cleanupInterval time.Duration
	// keepAlivePeriod 是公网 TCP 连接的保活探测间隔。
	keepAlivePeriod time.Duration
	// stopped 表示监听器是否已停止，受 lock 保护。
	stopped bool
	// stopOnce 确保停止逻辑只执行一次。
	stopOnce sync.Once
	// stopCh 用于通知接收和清理协程退出。
	stopCh chan struct{}
}

// NewListen 创建一个尚未绑定端口的公网监听器。
func NewListen(addr, mapPort string) *Listen {
	return &Listen{
		addr:            addr,
		mapPort:         mapPort,
		connList:        make(map[string]*transport.Conn),
		waitSince:       make(map[string]time.Time),
		waitTimeout:     10 * time.Second,
		cleanupInterval: 1 * time.Second,
		stopCh:          make(chan struct{}),
	}
}

// SetTimeouts 设置等待连接清理和 TCP 保活参数。
func (l *Listen) SetTimeouts(waitTimeout, cleanupInterval, keepAlivePeriod time.Duration) *Listen {
	if waitTimeout > 0 {
		l.waitTimeout = waitTimeout
	}
	if cleanupInterval > 0 {
		l.cleanupInterval = cleanupInterval
	}
	l.keepAlivePeriod = keepAlivePeriod
	return l
}

// Bind 同步绑定公网端口；仅在成功后调用方才可持久化端口映射。
func (l *Listen) Bind() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.listen != nil {
		return nil
	}
	listener, err := net.Listen("tcp", l.addr)
	if err != nil {
		return err
	}
	if tcpListener, ok := listener.(*net.TCPListener); ok {
		_ = tcpListener.SetDeadline(time.Time{})
	}
	l.listen = listener
	return nil
}

// Start 同步绑定端口后开始接收公网连接。
func (l *Listen) Start() error {
	if err := l.Bind(); err != nil {
		return err
	}
	return l.Serve()
}

// Serve 在端口绑定成功后持续接收公网连接。
func (l *Listen) Serve() error {
	l.lock.RLock()
	listener := l.listen
	l.lock.RUnlock()
	if listener == nil {
		return net.ErrClosed
	}
	go l.reapWaiters()
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-l.stopCh:
				return nil
			default:
			}
			// 非预期 Accept 错误时必须释放监听器和全部跟踪连接，避免端口泄漏。
			_ = l.Stop()
			return err
		}
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			_ = tcpConn.SetKeepAlive(true)
			if l.keepAlivePeriod > 0 {
				_ = tcpConn.SetKeepAlivePeriod(l.keepAlivePeriod)
			}
		}
		go l.dealWith(conn)
	}
}

// dealWith 将新公网连接登记为等待状态并通知所属客户端建立映射。
func (l *Listen) dealWith(conn net.Conn) {
	key := getKey()
	handConn := transport.NewConn(conn).SetSymbol(key).SetStatusWait()
	if !l.addWaitingConn(key, handConn) {
		_ = handConn.Close()
		return
	}
	if err := l.Notify(message.Message{Name: key, Symbol: key, Type: message.MsgTypeLink, TargetSymbol: l.addr, Msg: l.mapPort}); err != nil {
		l.DelConn(key)
		_ = handConn.Close()
		logger.Error("notify err", err)
		return
	}
}

// Stop 关闭监听器、已跟踪连接及等待连接清理协程。
func (l *Listen) Stop() error {
	var closeErr error
	l.stopOnce.Do(func() {
		close(l.stopCh)
		l.lock.Lock()
		l.stopped = true
		connections := make([]*transport.Conn, 0, len(l.connList))
		for _, conn := range l.connList {
			connections = append(connections, conn)
		}
		listener := l.listen
		l.connList = make(map[string]*transport.Conn)
		l.waitSince = make(map[string]time.Time)
		l.listen = nil
		l.lock.Unlock()

		for _, conn := range connections {
			if err := conn.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
		if listener != nil {
			if err := listener.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
	})
	return closeErr
}

// SetNotify 设置新公网连接的通知回调。
func (l *Listen) SetNotify(fn func(msg message.Message) error) *Listen {
	l.Notify = fn
	return l
}

// GetAddr 返回公网监听地址。
func (l *Listen) GetAddr() string {
	return l.addr
}

// GetMapPort 返回客户端本地映射地址。
func (l *Listen) GetMapPort() string {
	return l.mapPort
}

// SetMapPort 设置客户端本地映射地址。
func (l *Listen) SetMapPort(mapPort string) *Listen {
	l.mapPort = mapPort
	return l
}

// GetConn 按连接标识获取已跟踪的数据连接。
func (l *Listen) GetConn(key string) *transport.Conn {
	l.lock.RLock()
	conn := l.connList[key]
	l.lock.RUnlock()
	return conn
}

// ActivateConn 将等待中的公网连接标记为已配对，并取消其等待超时清理。
func (l *Listen) ActivateConn(key string) *transport.Conn {
	l.lock.Lock()
	defer l.lock.Unlock()
	conn := l.connList[key]
	if conn == nil || !conn.IsWait() {
		return nil
	}
	conn.SetStatusActive()
	delete(l.waitSince, key)
	return conn
}

// DelConn 删除一个主连接及可选的关联连接记录。
func (l *Listen) DelConn(key string, keys ...string) *Listen {
	l.lock.Lock()
	delete(l.connList, key)
	delete(l.waitSince, key)
	for _, k := range keys {
		delete(l.connList, k)
		delete(l.waitSince, k)
	}
	l.lock.Unlock()
	return l
}

// GetWaitConnList 返回所有等待与客户端配对的连接。
func (l *Listen) GetWaitConnList() (connList []*transport.Conn) {
	l.lock.RLock()
	for _, conn := range l.connList {
		if conn.IsWait() {
			connList = append(connList, conn)
		}
	}
	l.lock.RUnlock()
	return
}

// AddConn 在监听器未停止时登记数据连接。
func (l *Listen) AddConn(key string, conn *transport.Conn) *Listen {
	l.lock.Lock()
	if !l.stopped {
		l.connList[key] = conn
	}
	l.lock.Unlock()
	return l
}

// addWaitingConn 在监听器存活时原子地登记等待连接及其开始时间。
func (l *Listen) addWaitingConn(key string, conn *transport.Conn) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.stopped {
		return false
	}
	l.connList[key] = conn
	l.waitSince[key] = time.Now()
	return true
}

// reapWaiters 定期关闭等待配对超时的公网连接。
func (l *Listen) reapWaiters() {
	interval := l.cleanupInterval
	if interval <= 0 {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			var expired []*transport.Conn
			l.lock.Lock()
			for key, since := range l.waitSince {
				if l.waitTimeout <= 0 || now.Sub(since) < l.waitTimeout {
					continue
				}
				conn := l.connList[key]
				if conn == nil {
					delete(l.connList, key)
					delete(l.waitSince, key)
					continue
				}
				if !conn.IsWait() {
					// 已配对连接不再等待，即使旧调用方遗留了等待时间记录。
					delete(l.waitSince, key)
					continue
				}
				expired = append(expired, conn)
				delete(l.connList, key)
				delete(l.waitSince, key)
			}
			l.lock.Unlock()
			for _, conn := range expired {
				_ = conn.Close()
			}
		case <-l.stopCh:
			return
		}
	}
}
