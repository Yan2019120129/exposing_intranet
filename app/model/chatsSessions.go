package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"time"
)

const (
	ChatsSessionsTypeUser = 1 // 用户

	ChatsSessionsStatusPending = 10 // 等待
	ChatsSessionsStatusActive  = 20 // 活跃
	ChatsSessionsStatusClosed  = 30 // 结束
)

// ChatsSessions 聊天会话
type ChatsSessions struct {
	BaseModel
	SessionID string            `gorm:"type:varchar(255);uniqueIndex;comment:会话ID" json:"sessionId"`
	AdminID   uint              `gorm:"type:int unsigned;index;comment:管理ID" json:"adminId"`
	UserID    uint              `gorm:"type:int unsigned;index;comment:用户ID" json:"userId"`
	Name      string            `gorm:"type:varchar(255);comment:名称" json:"name"`
	Type      int8              `gorm:"type:tinyint;default:1;comment:类型(1:用户)" json:"type"`
	Number    int8              `gorm:"type:smallint;default:0;comment:未读消息" json:"number"`
	Status    int8              `gorm:"type:tinyint;default:10;comment:状态(10:等待,20:活跃,30:结束)" json:"status"`
	Data      ChatsSessionsData `gorm:"type:json;comment:数据" json:"data"`
}

// ChatsSessionsData 聊天会话数据
type ChatsSessionsData struct {
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
func (d ChatsSessionsData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *ChatsSessionsData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	if err := db.AutoMigrate(&ChatsSessions{}); err != nil {
//		panic("Failed to auto migrate ChatsSessions table: " + err.Error())
//	}
//}
