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
	// server 是包级服务端实例。
	server = &Server{
		addr:   ":1060",
		lock:   sync.RWMutex{},
		client: make(map[string]*Client),
		status: make(map[string]ClientStatus),
	}
)

// ClientStatus 表示供 Web 层展示的客户端在线状态和网络延迟。
type ClientStatus struct {
	// Symbol 是客户端唯一标识。
	Symbol string `json:"symbol"`
	// LastSeenMs 是最后收到客户端消息的 Unix 毫秒时间戳。
	LastSeenMs int64 `json:"lastSeenMs"`
	// RttMs 是最近一次心跳往返时延，未知时为 -1。
	RttMs int64 `json:"rttMs"`
}

// Server 管理控制连接、端口映射监听器及客户端状态。
type Server struct {
	// addr 是控制服务的监听地址。
	addr string
	// lock 保护服务端共享状态。
	lock sync.RWMutex
	// listen 是控制服务的网络监听器。
	listen net.Listener
	// client 按客户端标识保存当前在线客户端。
	client map[string]*Client
	// status 按客户端标识保存最近状态。
	status map[string]ClientStatus
	// db 是客户端与端口映射的持久化存储。
	db *gorm.DB
}

// GetServer 返回包级服务端实例。
func GetServer() *Server {
	return server
}

// NewServer 创建并设置包级服务端实例。
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

// Start 启动控制服务并持续接收客户端连接。
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

// dealWith 处理单个控制连接的首条消息。
func (s *Server) dealWith(conn net.Conn) {
	handConn := transport.NewConnWithOptions(conn, s.controlConnOptions())
	defer handConn.Close()
	// 兜底：单个连接处理 goroutine 的 panic 不应打挂整个服务端
	defer func() {
		if r := recover(); r != nil {
			logger.Error("dealWith panic:", r)
		}
	}()
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

// Status 记录客户端的心跳响应、最后在线时间和往返时延。
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

// GetClientStatus 返回指定客户端的最近状态。
func (s *Server) GetClientStatus(symbol string) (ClientStatus, bool) {
	s.lock.RLock()
	st, ok := s.status[symbol]
	s.lock.RUnlock()
	return st, ok
}

// ListClientStatus 返回所有已记录的客户端状态。
func (s *Server) ListClientStatus() []ClientStatus {
	s.lock.RLock()
	out := make([]ClientStatus, 0, len(s.status))
	for _, v := range s.status {
		out = append(out, v)
	}
	s.lock.RUnlock()
	return out
}

// connect 将客户端数据连接与等待中的公网连接建立双向转发。
func (s *Server) connect(msg message.Message, conn *transport.Conn) {
	logger.Info("server received link request")
	// 每一步都可能在取值前因客户端并发断开而被移除，必须逐级判 nil
	client := s.GetClient(msg.Symbol)
	if client == nil {
		logger.Info("link failed: client or target port not found")
		return
	}
	listen := client.GetListen(msg.TargetSymbol)
	if listen == nil {
		logger.Info("target port does not exist")
		return
	}
	targetConn := listen.ActivateConn(msg.Name)
	if targetConn == nil {
		logger.Info("target port does not exist")
		return
	}
	key := getKey()
	listen.AddConn(key, conn.SetSymbol(key))
	_ = CommunicationWithTimeout(targetConn, conn, s.dataTimeout())
	listen.DelConn(targetConn.GetSymbol(), conn.GetSymbol())
}

// Register 验证并注册客户端，然后维护控制连接及心跳。
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

// AddClient 将客户端写入服务端客户端表。
func (s *Server) AddClient(key string, client *Client) *Server {
	s.lock.Lock()
	s.client[key] = client
	s.lock.Unlock()
	return s
}

// AddClientIfAbsent 仅在客户端标识尚未注册时添加客户端。
func (s *Server) AddClientIfAbsent(key string, client *Client) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, exists := s.client[key]; exists {
		return false
	}
	s.client[key] = client
	return true
}

// GetClient 按客户端标识获取当前在线客户端。
func (s *Server) GetClient(key string) *Client {
	s.lock.RLock()
	client := s.client[key]
	s.lock.RUnlock()
	return client
}

// DelClient 删除一个或多个客户端记录。
func (s *Server) DelClient(key string, keys ...string) {
	s.lock.Lock()
	delete(s.client, key)
	for _, k := range keys {
		delete(s.client, k)
	}
	s.lock.Unlock()
}

// DelClientIf 仅当当前记录仍为指定客户端时删除该记录。
func (s *Server) DelClientIf(key string, client *Client) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if current, ok := s.client[key]; ok && current == client {
		delete(s.client, key)
	}
}

// IsExist 返回指定客户端是否已注册。
func (s *Server) IsExist(key string) bool {
	s.lock.RLock()
	_, ok := s.client[key]
	s.lock.RUnlock()
	return ok
}

// IsListenExist 返回指定客户端是否存在指定的公网端口监听器。
func (s *Server) IsListenExist(clientSymbol, serverPort string) bool {
	if !s.IsExist(clientSymbol) {
		return false
	}
	client := s.GetClient(clientSymbol)
	return client != nil && client.IsExist(serverPort)
}

// NewListen 为指定客户端创建并启动公网端口监听器。
func (s *Server) NewListen(clientSymbol, localPort, serverPort string) error {
	client := s.GetClient(clientSymbol)
	if client == nil {
		return errors.New("client is not exist! ")
	}
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
	// 先登记后启动，避免 Serve 或 Client.Close 在间隙结束后留下失效监听器。
	client.AddListen(serverPort, listen)
	go func() {
		defer client.DelListenIf(serverPort, listen)
		err := listen.Serve()
		if err != nil {
			logger.Info("listen start fail:", err)
		}
	}()
	return nil
}

// Close 关闭控制监听器及所有客户端资源。
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

// DelClientF 通知指定客户端其身份已被永久删除。
func (s *Server) DelClientF(symbols ...string) error {
	for _, v := range symbols {
		// 客户端可能在检查与取值之间断开，GetClient 可能返回 nil
		if client := s.GetClient(v); client != nil {
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

// CloseListen 关闭指定客户端的公网端口监听器。
func (s *Server) CloseListen(clientSymbol, localPort, serverPort string) error {
	client := s.GetClient(clientSymbol)
	if client == nil {
		return errors.New("client is not exist! ")
	}
	listen := client.GetListen(serverPort)
	if listen == nil {
		return errors.New("client of listen is not exist! ")
	}
	client.DelListenIf(serverPort, listen)
	return listen.Stop()
}

// GetAddr 返回控制服务的监听地址。
func (s *Server) GetAddr() string {
	return s.addr
}

// SetAddr 设置控制服务的监听地址。
func (s *Server) SetAddr(addr string) *Server {
	s.addr = addr
	return s
}

// controlConnOptions 返回控制连接使用的超时和 TCP 保活配置。
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

// dataTimeout 返回数据转发连接的读写空闲超时。
func (s *Server) dataTimeout() time.Duration {
	if cfg := configs.GetConnect(); cfg != nil {
		return time.Duration(cfg.GetReadWriteTimeout()) * time.Second
	}
	return 0
}

// pingSentAt 从心跳消息中提取发送时间的 Unix 毫秒时间戳。
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
