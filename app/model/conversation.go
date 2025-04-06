package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"time"
)

const (
	// ConversationStatusActive 会话状态正常
	ConversationStatusActive = 10

	// ConversationStatusDisable 会话状态结束
	ConversationStatusDisable = -1

	// ConversationTypePrivate 私信
	ConversationTypePrivate = 1

	// ConversationOnline 在线
	ConversationOnline = 1

	// ConversationOffline 离线
	ConversationOffline = 2

	// ConversationModeUser 用户会话模式
	ConversationModeUser = 1

	// ConversationModeStore 商家会话模式
	ConversationModeStore = 4

	// ConversationModeSystem 系统会话
	ConversationModeSystem = 3
)

// ChatsConversation 聊天会话
type ChatsConversation struct {
	BaseModel
	SessionID  string                `gorm:"type:varchar(255) not null;comment:会话Key" json:"sessionId"`
	AdminID    uint                  `gorm:"type:int unsigned;index;comment:管理ID" json:"adminId"`
	UserID     uint                  `gorm:"type:int unsigned not null;comment:用户ID" json:"userId"`
	ReceiverID uint                  `gorm:"type:int unsigned not null;comment:接收ID" json:"receiverId"`
	Name       string                `gorm:"type:varchar(255);comment:名称" json:"name"`
	Type       int8                  `gorm:"type:tinyint;default:1;comment:类型(1:私信)" json:"type"`
	Number     int8                  `gorm:"type:smallint;default:0;comment:未读消息" json:"number"`
	Mode       int8                  `gorm:"type:tinyint;default:1;comment:会话模式(1:用户,4:商家)" json:"mode"`
	Online     int8                  `gorm:"type:tinyint(1);default:1;comment:在线(1:离线,2:在线)" json:"online"`
	Status     int8                  `gorm:"type:smallint;default:10;comment:状态(-1:屏蔽,10:正常)" json:"status"`
	Data       ChatsConversationData `gorm:"type:json;comment:数据" json:"data"`
}

// ChatsConversationData 聊天会话数据
type ChatsConversationData struct {
	ID         uint      `json:"id"`         //	最后消息ID
	SessionID  string    `json:"sessionId"`  //	最后消息会话ID
	SenderID   uint      `json:"senderId"`   //	最后消息发送者ID
	SenderType int8      `json:"senderType"` //	最后消息发送者类型
	ReceiverID uint      `json:"receiverId"` //	最后消息接收者ID
	Type       int8      `json:"type"`       //	最后消息类型
	Message    string    `json:"message"`    //	最后消息内容
	Status     int8      `json:"status"`     //	最后消息状态
	CreatedAt  time.Time `json:"createdAt"`  //	最后消息时间
}

// Value implements the driver.Valuer interface
func (d ChatsConversationData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *ChatsConversationData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&ChatsConversation{}); err != nil {
//		panic("Failed to auto migrate WalletAccount table: " + err.Error())
//	}
//}
