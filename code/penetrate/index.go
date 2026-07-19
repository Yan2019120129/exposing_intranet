package penetrate

import (
	"errors"
	"fmt"
	"my-base/app/models"
	"my-base/app/repository"
	"my-base/code/message"
	transport "my-base/code/transport"
	"my-base/configs"
	"net"
	"sync"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/logger"
	"gorm.io/gorm"
)

var (
	server = &Server{
		addr:   ":1060",
		lock:   sync.RWMutex{},
		client: make(map[string]*Client),
		status: make(map[string]ClientStatus),
	}
)

// ClientStatus 供后续 Web 展示的客户端状态（先实现 RTT/在线时间）
type ClientStatus struct {
	Symbol     string `json:"symbol"`
	LastSeenMs int64  `json:"lastSeenMs"`
	RttMs      int64  `json:"rttMs"`
}

type Server struct {
	addr   string
	lock   sync.RWMutex
	listen net.Listener
	client map[string]*Client
	status map[string]ClientStatus
	db     *gorm.DB
}

// GetServer 获取服务端
func GetServer() *Server {
	return server
}

// NewServer 创建服务
func NewServer(addr string, db *gorm.DB) *Server {
	server = &Server{
		addr:   addr,
		lock:   sync.RWMutex{},
		client: make(map[string]*Client),
		status: make(map[string]ClientStatus),
		db:     db,
	}
	return server
}

// Start 启动服务
func (s *Server) Start() error {
	var err error
	s.listen, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := s.listen.Accept()
		if err != nil {
			return err
		}
		go s.dealWith(conn)
	}
}

// dealWith 处理方法
func (s *Server) dealWith(conn net.Conn) {
	handConn := transport.NewConnWithOptions(conn, s.controlConnOptions())
	defer handConn.Close()
	msg := message.Message{}
	err := handConn.ParseMsg(&msg)
	if err != nil {
		logger.Error("read err", err)
		return
	}
	logger.Info("received control message type:", msg.Type)
	switch {
	case msg.EqLink():
		s.connect(msg, handConn)
	case msg.EqRegister():
		s.Register(msg, handConn)
	}
}

// 客户端
func (s *Server) Status(msg message.Message) {
	var pong message.PongPayload
	switch v := msg.Msg.(type) {
	case string:
		// 兼容旧客户端：msg.Msg == "pong"
	case map[string]any:
		if seq, ok := v["seq"].(float64); ok {
			pong.Seq = uint64(seq)
		}
		if sentAt, ok := v["sentAtMs"].(float64); ok {
			pong.SentAtMs = int64(sentAt)
		}
		if clientAt, ok := v["clientAtMs"].(float64); ok {
			pong.ClientAtMs = int64(clientAt)
		}
	}

	now := time.Now().UnixMilli()
	rtt := int64(-1)
	if pong.SentAtMs > 0 {
		rtt = now - pong.SentAtMs
	}
	if client := s.GetClient(msg.Symbol); client != nil {
		if measured := client.RecordPong(pong.Seq); measured >= 0 {
			rtt = measured.Milliseconds()
		}
	}

	s.lock.Lock()
	s.status[msg.Symbol] = ClientStatus{
		Symbol:     msg.Symbol,
		LastSeenMs: now,
		RttMs:      rtt,
	}
	s.lock.Unlock()

	fmt.Println("pong received, seq:", pong.Seq, "rttMs:", rtt)
}

// GetClientStatus 获取单个客户端状态（后续 Web API 可直接调用）
func (s *Server) GetClientStatus(symbol string) (ClientStatus, bool) {
	s.lock.RLock()
	st, ok := s.status[symbol]
	s.lock.RUnlock()
	return st, ok
}

// ListClientStatus 列表（后续 Web API 可直接调用）
func (s *Server) ListClientStatus() []ClientStatus {
	s.lock.RLock()
	out := make([]ClientStatus, 0, len(s.status))
	for _, v := range s.status {
		out = append(out, v)
	}
	s.lock.RUnlock()
	return out
}

// 建立连接
func (s *Server) connect(msg message.Message, conn *transport.Conn) {
	logger.Info("server received link request")
	if s.IsExist(msg.Symbol) {
		client := s.GetClient(msg.Symbol)
		if client.IsExist(msg.TargetSymbol) {
			listen := client.GetListen(msg.TargetSymbol)
			targetConn := listen.ActivateConn(msg.Name)
			if targetConn != nil {
				key := getKey()
				listen.AddConn(key, conn.SetSymbol(key))
				_ = CommunicationWithTimeout(targetConn, conn, s.dataTimeout())
				listen.DelConn(targetConn.GetSymbol(), conn.GetSymbol())
			} else {
				logger.Info("target port does not exist")
			}
			// waitConnList := listen.GetWaitConnList()
			// for _, targetConn := range waitConnList {
			// 	key := getKey()
			// 	listen.AddConn(key, conn.SetSymbol(key))
			// 	_ = Communication(targetConn, conn)
			// 	listen.DelConn(targetConn.GetSymbol(), conn.GetSymbol())
			// 	break
			// }
		} else {
			logger.Info("target port does not exist")
		}
	} else {
		logger.Info("link failed: client or target port not found")
	}
}

// Register 注册客户端
func (s *Server) Register(param message.Message, conn *transport.Conn) {
	// 必须提供 symbol（通过 HTTP 认证获取）
	if param.Symbol == "" {
		param.SetMsg("Client not registered. Please authenticate first via HTTP API.")
		param.SetType(message.MsgTypeClose)
		_ = conn.Send(param)
		_ = conn.Close()
		logger.Error("client register failed: symbol is empty")
		return
	}

	// 验证 symbol 是否在数据库中存在
	clientRepository := repository.NewClientRepository(s.db)
	clientInfo, err := clientRepository.FindBySymbol(param.Symbol)
	if err != nil {
		param.SetMsg("Invalid symbol. Please authenticate first via HTTP API.")
		param.SetType(message.MsgTypeClose)
		_ = conn.Send(param)
		_ = conn.Close()
		logger.Error("client register failed: symbol not found in database")
		return
	}

	// 检查是否已被禁用
	if models.IsDisabledStatus(clientInfo.Status) {
		param.SetMsg("Client is disabled.")
		param.SetType(message.MsgTypeClose)
		_ = conn.Send(param)
		_ = conn.Close()
		logger.Error("client register failed: client is disabled")
		return
	}

	// 如果已经有对应的客户端在线，拒绝连接
	client := NewClient(conn.SetStatusActive()).SetSymbol(param.Symbol)
	if !s.AddClientIfAbsent(param.Symbol, client) {
		param.SetMsg("The client already exists or is not completely closed.")
		param.SetType(message.MsgTypeClose)
		_ = conn.Send(param)
		_ = conn.Close()
		return
	}

	// 创建客户端并保证旧连接退出时不会删除新连接
	stop := make(chan struct{})
	var stopOnce sync.Once
	defer func() {
		stopOnce.Do(func() { close(stop) })
		s.DelClientIf(param.Symbol, client)
		err := client.Close()
		if err != nil {
			logger.Error("end client close", err)
		}
	}()
	// 获取客户端信息
	// 更新客户端状态为活跃
	if clientInfo.Id > 0 {
		// 更新客户端名称（如果提供了新的名称）
		if param.Name != "" {
			_ = clientRepository.UpdateBySymbol(param.Symbol, &models.Client{Name: param.Name, Status: models.StatusActive})
		} else {
			_ = clientRepository.UpdateStatusBySymbol(param.Symbol, models.StatusActive)
		}

		// 恢复端口映射；每条映射都同步 bind，单条失败不影响客户端主连接。
		for _, portAndAddr := range clientInfo.PortList {
			if err := s.NewListen(param.Symbol, portAndAddr.Local, portAndAddr.Server); err != nil {
				logger.Error("restore listen failed:", err)
			}
		}
	}

	// 启动心跳检测
	go func() {
		cfg := configs.GetConnect()
		if cfg == nil {
			return
		}
		if err := client.PingLoopWithConfig(
			time.Duration(cfg.GetPingInterval())*time.Second,
			time.Duration(cfg.GetPongTimeout())*time.Second,
			cfg.GetMaxPingFailures(),
			stop,
		); err != nil {
			logger.Error("ping loop err:", err)
			_ = clientRepository.UpdateStatusBySymbol(param.Symbol, models.StatusOn)
			_ = client.Close()
		}
	}()

	// 主连接读循环：消费客户端 pong（以及未来扩展的状态上报）
	for {
		m := message.Message{}
		if err := conn.ParseMsg(&m); err != nil {
			stopOnce.Do(func() { close(stop) })
			_ = clientRepository.UpdateStatusBySymbol(param.Symbol, models.StatusOn)
			return
		}
		switch {
		case m.EqPong():
			s.Status(m)
		case m.EqPing():
			m.Type = message.MsgTypePong
			m.Msg = message.PongPayload{SentAtMs: pingSentAt(m), ClientAtMs: time.Now().UnixMilli()}
			if err := conn.Send(m); err != nil {
				stopOnce.Do(func() { close(stop) })
				return
			}
		default:
			// 保持兼容：先打印，后续可扩展为客户端状态上报等
			logger.Info("client message type:", m.Type)
		}
	}
}

// AddClient 添加客户端
func (s *Server) AddClient(key string, client *Client) *Server {
	s.lock.Lock()
	s.client[key] = client
	s.lock.Unlock()
	return s
}

func (s *Server) AddClientIfAbsent(key string, client *Client) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, exists := s.client[key]; exists {
		return false
	}
	s.client[key] = client
	return true
}

// GetClient 添加客户端
func (s *Server) GetClient(key string) *Client {
	s.lock.RLock()
	client := s.client[key]
	s.lock.RUnlock()
	return client
}

// DelClient 删除客户端
func (s *Server) DelClient(key string, keys ...string) {
	s.lock.Lock()
	delete(s.client, key)
	for _, k := range keys {
		delete(s.client, k)
	}
	s.lock.Unlock()
}

func (s *Server) DelClientIf(key string, client *Client) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if current, ok := s.client[key]; ok && current == client {
		delete(s.client, key)
	}
}

// IsExist 是否存在
func (s *Server) IsExist(key string) bool {
	s.lock.RLock()
	_, ok := s.client[key]
	s.lock.RUnlock()
	return ok
}

// IsListenExist reports whether a public port listener is active for a
// connected client.
func (s *Server) IsListenExist(clientSymbol, serverPort string) bool {
	if !s.IsExist(clientSymbol) {
		return false
	}
	client := s.GetClient(clientSymbol)
	return client != nil && client.IsExist(serverPort)
}

// NewListen 申请添加客户端，添加端口
func (s *Server) NewListen(clientSymbol, localPort, serverPort string) error {
	if !s.IsExist(clientSymbol) {
		return errors.New("client is not exist! ")
	}
	client := s.GetClient(clientSymbol)
	if client.IsExist(serverPort) {
		return errors.New("client port exist")
	}
	s.lock.RLock()
	for symbol, otherClient := range s.client {
		if symbol != clientSymbol && otherClient.IsExist(serverPort) {
			s.lock.RUnlock()
			return errors.New("server port is already in use")
		}
	}
	s.lock.RUnlock()

	listen := NewListen(serverPort, localPort).SetNotify(client.TakeOver)
	if cfg := configs.GetConnect(); cfg != nil {
		wait := time.Duration(cfg.GetConnectWait()) * time.Second
		cleanup := time.Duration(cfg.GetWaitConnTimeout()) * time.Second
		keepAlive := time.Duration(cfg.GetTCPKeepAlive()) * time.Second
		listen.SetTimeouts(wait, cleanup, keepAlive)
	}
	if err := listen.Bind(); err != nil {
		return err
	}
	go func() {
		defer client.DelListenIf(serverPort, listen)
		err := listen.Serve()
		if err != nil {
			logger.Info("listen start fail:", err)
		}
	}()

	client.AddListen(serverPort, listen)
	return nil
}

// Close 关闭资源
func (s *Server) Close() error {
	s.lock.RLock()
	clients := make([]*Client, 0, len(s.client))
	for _, client := range s.client {
		clients = append(clients, client)
	}
	listener := s.listen
	s.lock.RUnlock()
	if listener != nil {
		_ = listener.Close()
	}
	for _, client := range clients {
		err := client.Close()
		if err != nil {
			return err
		}
		clientRepository := repository.NewClientRepository(s.db)
		_ = clientRepository.UpdateStatusBySymbol(client.symbol, models.StatusOn)
	}
	return nil
}

// DelClientF 永久删除客户端
func (s *Server) DelClientF(symbols ...string) error {
	for _, v := range symbols {
		if s.IsExist(v) {
			client := s.GetClient(v)
			err := client.Send(message.Message{
				Type:   message.MsgTypeDel,
				Symbol: v,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CloseListen 关闭监听
func (s *Server) CloseListen(clientSymbol, localPort, serverPort string) error {
	if s.IsExist(clientSymbol) {
		client := s.GetClient(clientSymbol)
		if client.IsExist(serverPort) {
			listen := client.GetListen(serverPort)
			client.DelListenIf(serverPort, listen)
			err := listen.Stop()
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("client of listen is not exist! ")
	}
	return errors.New("client is not exist! ")
}

// GetAddr 获取服务端
func (s *Server) GetAddr() string {
	return s.addr
}

// SetAddr 设置服务端监听地址
func (s *Server) SetAddr(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) controlConnOptions() transport.ConnOptions {
	options := transport.ConnOptions{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if cfg := configs.GetConnect(); cfg != nil {
		options.ReadTimeout = time.Duration(cfg.GetPongTimeout()+cfg.GetPingInterval()) * time.Second
		options.WriteTimeout = time.Duration(cfg.GetReadWriteTimeout()) * time.Second
		options.KeepAlive = true
		options.KeepAlivePeriod = time.Duration(cfg.GetTCPKeepAlive()) * time.Second
	}
	return options
}

func (s *Server) dataTimeout() time.Duration {
	if cfg := configs.GetConnect(); cfg != nil {
		return time.Duration(cfg.GetReadWriteTimeout()) * time.Second
	}
	return 0
}

func pingSentAt(msg message.Message) int64 {
	switch value := msg.Msg.(type) {
	case message.PingPayload:
		return value.SentAtMs
	case map[string]any:
		if sentAt, ok := value["sentAtMs"].(float64); ok {
			return int64(sentAt)
		}
	}
	return 0
}
