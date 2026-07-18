package message

import (
	"encoding/json"
)

const (
	// MsgTypeDefault 默认数据,服务与服务之间的通信数据
	MsgTypeDefault = iota

	// MsgTypePing 服务状态检测类型
	MsgTypePing

	// MsgTypePong 服务状态响应类型
	MsgTypePong

	// MsgTypeRegister 服务注册消息
	MsgTypeRegister

	// MsgTypeLink 链接指定服务消息
	MsgTypeLink

	// MsgTypeClose 关闭类型
	MsgTypeClose

	// MsgTypeDel 删除类型
	MsgTypeDel
)

// PingPayload 心跳请求载荷（用于计算网络 RTT）
// 注意：SentAtMs 由服务端填充，客户端收到后原样回传
type PingPayload struct {
	Seq      uint64 `json:"seq"`
	SentAtMs int64  `json:"sentAtMs"`
}

// PongPayload 心跳响应载荷
type PongPayload struct {
	Seq        uint64 `json:"seq"`
	SentAtMs   int64  `json:"sentAtMs"`
	ClientAtMs int64  `json:"clientAtMs"`
}

// Message 接收的消息
type Message struct {
	Type         int    `json:"msgType"`      // 消息类型
	Name         string `json:"name"`         // 客户端名
	Symbol       string `json:"symbol"`       // 链接标识
	TargetSymbol string `json:"targetSymbol"` // 目标标识
	Msg          any    `json:"msg"`          // 具体消息
	Err          error  `json:"err"`          // 返回的错误
}

func NewMessage(t int, msg any) *Message {
	return &Message{
		Type: t,
		Msg:  msg,
	}
}

// Parse 解析消息
func Parse(msg []byte) (Message, error) {
	tmpMsg := Message{}
	err := json.Unmarshal(msg, &tmpMsg)
	if err != nil {
		tmpMsg.Type = MsgTypeDefault
	}
	return tmpMsg, err
}

// EqRegister 判断是否等于注册类型
func (m *Message) EqRegister() bool {
	return m.Type == MsgTypeRegister
}

// EqLink 判断是否等与连接类型
func (m *Message) EqLink() bool {
	return m.Type == MsgTypeLink
}

// EqDel 判断是否等与删除类型
func (m *Message) EqDel() bool {
	return m.Type == MsgTypeDel
}

// EqPing 判断是否连接检测
func (m *Message) EqPing() bool {
	return m.Type == MsgTypePing
}

// EqPong 判断是否 pong 响应
func (m *Message) EqPong() bool {
	return m.Type == MsgTypePong
}

// EqClose 判断是否等与关闭类型类型
func (m *Message) EqClose() bool {
	return m.Type == MsgTypeClose
}

// EqDefault 判断是否普通链接
func (m *Message) EqDefault() bool {
	return m.Type == MsgTypeDefault
}

// GetType 获取消息类型
func (m *Message) GetType() int {
	return m.Type
}

// SetType 设置消息类型
func (m *Message) SetType(t int) *Message {
	m.Type = t
	return m
}

// GetMsg 获取消息
func (m *Message) GetMsg() any {
	return m.Msg
}

// SetMsg 放置消息
func (m *Message) SetMsg(msg any) *Message {
	m.Msg = msg
	return m
}

// GetSymbol 获取标识
func (m *Message) GetSymbol() string {
	return m.Symbol
}

// SetSymbol 放置标识
func (m *Message) SetSymbol(symbol string) *Message {
	m.Symbol = symbol
	return m
}

// GetTargetSymbol 获取目标标识
func (m *Message) GetTargetSymbol() string {
	return m.TargetSymbol
}

// SetTargetSymbol 放置目标标识
func (m *Message) SetTargetSymbol(symbol string) *Message {
	m.TargetSymbol = symbol
	return m
}

// SetTypeDefault 放置目标标识
func (m *Message) SetTypeDefault() *Message {
	return m.SetType(MsgTypeDefault)
}

// SetTypePing 放置目标标识
func (m *Message) SetTypePing() *Message {
	return m.SetType(MsgTypePing)
}

// SetTypeRegister 放置目标标识
func (m *Message) SetTypeRegister() *Message {
	return m.SetType(MsgTypeRegister)
}

// SetTypeLink 放置目标标识
func (m *Message) SetTypeLink() *Message {
	return m.SetType(MsgTypeLink)
}

// ToString 转换为 string 类型
func (m *Message) ToString() string {
	return string(m.ToBytes())
}

// ToBytes 转换为 []byte 类型
func (m *Message) ToBytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
