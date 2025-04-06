package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
)

const (
	// Notify modes
	NotifyModeBackend  int8 = 1  // 后台
	NotifyModeFrontend int8 = 11 // 前台

	// Notify types
	NotifyTypeSystem int8 = 1 // 系统通知

	// Notify statuses
	NotifyStatusUnread int8 = 10 // 未读
	NotifyStatusRead   int8 = 20 // 已读
)

// Notify 系统通知表
type Notify struct {
	BaseModel
	AdminID uint       `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID  uint       `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	Mode    int8       `gorm:"type:tinyint not null;index;comment:模式(1:后台,11:前台)" json:"mode"`
	Type    int8       `gorm:"type:tinyint not null;default:1;index;comment:类型(1:系统通知)" json:"type"`
	Title   string     `gorm:"type:varchar(1020) not null;comment:标题" json:"title"`
	Content string     `gorm:"type:text;comment:内容" json:"content"`
	Status  int8       `gorm:"type:tinyint not null;default:10;index;comment:状态(10:未读,20:已读)" json:"status"`
	Data    NotifyData `gorm:"type:json;comment:数据" json:"data"`
}

// NotifyDisplayData 通知展示数据
type NotifyDisplayData struct {
	ID        uint      `json:"id"`
	Type      int8      `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetDisplayData 获取通知展示数据
func (n *Notify) GetDisplayData() *NotifyDisplayData {
	return &NotifyDisplayData{
		ID:        n.ID,
		Type:      n.Type,
		Title:     n.Title,
		Content:   n.Content,
		Status:    n.Status,
		CreatedAt: n.CreatedAt,
	}
}

// NotifyData 通知数据
type NotifyData struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Value implements the driver.Valuer interface
func (d NotifyData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *NotifyData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Notify{}); err != nil {
//		panic("Failed to auto migrate Notify table: " + err.Error())
//	}
//}
