package penetrate

import (
	"my-base/code/message"
	transport "my-base/code/transport"
	"net"
	"sync"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/logger"
)

type Listen struct {
	lock            sync.RWMutex
	addr            string
	mapPort         string
	Notify          func(msg message.Message) error
	listen          net.Listener
	connList        map[string]*transport.Conn
	waitSince       map[string]time.Time
	waitTimeout     time.Duration
	cleanupInterval time.Duration
	keepAlivePeriod time.Duration
	stopOnce        sync.Once
	stopCh          chan struct{}
}

// NewListen 新建监听
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

// SetTimeouts configures waiting-connection cleanup and TCP keepalive.
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

// Bind reserves the public port synchronously. Callers can safely persist a
// mapping only after Bind returns nil.
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

// Start binds and serves the listener. Bind is intentionally synchronous.
func (l *Listen) Start() error {
	if err := l.Bind(); err != nil {
		return err
	}
	return l.Serve()
}

// Serve accepts public connections after Bind has succeeded.
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

// dealWith 处理方法
func (l *Listen) dealWith(conn net.Conn) {
	key := getKey()
	handConn := transport.NewConn(conn).SetSymbol(key).SetStatusWait()
	l.AddConn(key, handConn)
	l.lock.Lock()
	l.waitSince[key] = time.Now()
	l.lock.Unlock()
	if err := l.Notify(message.Message{Name: key, Symbol: key, Type: message.MsgTypeLink, TargetSymbol: l.addr, Msg: l.mapPort}); err != nil {
		l.DelConn(key)
		_ = handConn.Close()
		logger.Error("notify err", err)
		return
	}
}

// Stop 关闭服务
func (l *Listen) Stop() error {
	var closeErr error
	l.stopOnce.Do(func() {
		close(l.stopCh)
		l.lock.Lock()
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

// SetNotify 通知上级的方法
func (l *Listen) SetNotify(fn func(msg message.Message) error) *Listen {
	l.Notify = fn
	return l
}

// GetAddr 获取地址
func (l *Listen) GetAddr() string {
	return l.addr
}

// GetMapPort 获取映射的端口
func (l *Listen) GetMapPort() string {
	return l.mapPort
}

// SetMapPort 设置映射的端口
func (l *Listen) SetMapPort(mapPort string) *Listen {
	l.mapPort = mapPort
	return l
}

// GetConn 获取连接
func (l *Listen) GetConn(key string) *transport.Conn {
	l.lock.RLock()
	conn := l.connList[key]
	l.lock.RUnlock()
	return conn
}

// ActivateConn marks a waiting public connection as matched. The state
// transition and removal from the wait-timeout set are kept under the same
// listener lock so the reaper cannot close the connection between the match
// lookup and activation.
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

// DelConn 删除连接
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

// GetWaitConnList 获取等待的链接
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

// AddConn 添加连接
func (l *Listen) AddConn(key string, conn *transport.Conn) *Listen {
	l.lock.Lock()
	l.connList[key] = conn
	l.lock.Unlock()
	return l
}

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
					// A matched/active connection is no longer waiting even if
					// an older caller left a timestamp behind.
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
