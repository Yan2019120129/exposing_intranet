package models

const (
	ChatsMessagesStatusUnread = 1 // 未读
	ChatsMessagesStatusRead   = 2 // 已读

	ChatsMessagesTypeText  = 1 // 文本
	ChatsMessagesTypeImage = 2 // 图片
	ChatsMessagesTypeFile  = 3 // 文件

	ChatsMessagesSenderTypeUser   = 1 // 用户
	ChatsMessagesSenderTypeAdmin  = 2 // 管理
	ChatsMessagesSenderTypeSystem = 3 // 系统
	ChatsMessagesSenderTypeStore  = 4 // 商家
)

// ChatsMessages 聊天消息
type ChatsMessages struct {
	BaseModel
	AdminID    uint   `gorm:"type:int unsigned;index;comment:管理ID" json:"adminId"`
	SessionID  string `gorm:"type:varchar(255);index;comment:会话ID" json:"sessionId"`
	SenderID   uint   `gorm:"type:int unsigned;index;comment:发送者ID" json:"senderId"`
	ReceiverID uint   `gorm:"type:int unsigned;index;comment:接收者ID" json:"receiverId"`
	SenderType int8   `gorm:"type:tinyint;default:1;comment:发送者类型(1:用户,2:管理,3:系统,4:商家)" json:"senderType"`
	Message    string `gorm:"type:text;comment:消息" json:"message"`
	Type       int8   `gorm:"type:tinyint;default:1;comment:类型(1:文本,2:图片,3:文件)" json:"type"`
	Status     int8   `gorm:"type:tinyint;default:1;comment:状态(1:未读,2:已读)" json:"status"`
}

// ToSessionData 转换为会话数据
func (cm *ChatsMessages) ToSessionData() ChatsSessionsData {
	return ChatsSessionsData{
		ID:         cm.ID,
		SessionID:  cm.SessionID,
		SenderID:   cm.SenderID,
		SenderType: cm.SenderType,
		ReceiverID: cm.ReceiverID,
		Type:       cm.Type,
		Message:    cm.Message,
		Status:     cm.Status,
		CreatedAt:  cm.CreatedAt,
	}
}

// ToChatsConversationData 转换为会话数据
func (cm *ChatsMessages) ToChatsConversationData() ChatsConversationData {
	return ChatsConversationData{
		ID:         cm.ID,
		SessionID:  cm.SessionID,
		SenderID:   cm.SenderID,
		SenderType: cm.SenderType,
		ReceiverID: cm.ReceiverID,
		Type:       cm.Type,
		Message:    cm.Message,
		Status:     cm.Status,
		CreatedAt:  cm.CreatedAt,
	}
}

//func init() {
//	if err := db.AutoMigrate(&ChatsMessages{}); err != nil {
//		panic("Failed to auto migrate ChatsMessages table: " + err.Error())
//	}
//}
