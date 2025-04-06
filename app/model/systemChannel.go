package models

import (
	"database/sql/driver"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	ChannelTypeDefault  int8 = 1 // 默认渠道
	ChannelTypeApple    int8 = 2 // 苹果渠道
	ChannelTypeGoogle   int8 = 3 // 谷歌渠道
	ChannelTypeTelegram int8 = 4 // Telegram渠道

	ChannelStatusDisabled int8 = -1 // 禁用状态
	ChannelStatusEnabled  int8 = 10 // 激活状态
)

// Channel 渠道表
type Channel struct {
	BaseModel
	AdminID   uint        `gorm:"type:int unsigned not null;default:1;uniqueIndex:idx_admin_symbol;comment:管理ID" json:"adminId"`
	AppID     string      `gorm:"type:varchar(100) not null;uniqueIndex;comment:应用ID" json:"appId"`
	SecretKey string      `gorm:"type:varchar(255) not null;comment:密钥" json:"secretKey"`
	Icon      string      `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Name      string      `gorm:"type:varchar(60) not null;index;comment:名称" json:"name"`
	Symbol    string      `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_symbol;comment:标识" json:"symbol"`
	Type      int8        `gorm:"type:tinyint not null;default:1;comment:类型(1:默认渠道,2:苹果渠道,3:谷歌渠道,4:Telegram渠道)" json:"type"`
	Status    int8        `gorm:"type:tinyint not null;default:10;index;comment:状态" json:"status"`
	Desc      string      `gorm:"type:varchar(255);comment:描述" json:"desc,omitempty"`
	Data      ChannelData `gorm:"type:json;comment:数据" json:"data,omitempty"`
}

// ChannelDisplayData 渠道显示数据
type ChannelDisplayData struct {
	ID     uint   `json:"id"`
	Icon   string `json:"icon"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// ChannelData 渠道数据
type ChannelData struct {
	Gateway     string `json:"gateway,omitempty"`     // 网关
	CallbackURL string `json:"callbackUrl,omitempty"` // 回调URL
	IPWhitelist string `json:"ipWhitelist,omitempty"` // IP白名单
	ReturnURL   string `json:"returnUrl,omitempty"`   // 返回地址
}

// Value implements the driver.Valuer interface
func (d ChannelData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *ChannelData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Channel{}); err != nil {
//		panic("Failed to auto migrate Channel table: " + err.Error())
//	}
//
//	// Initialize system channels
//	if err := InitSystemChannels(db); err != nil {
//		panic("Failed to initialize system channels: " + err.Error())
//	}
//}

// InitSystemChannels initializes the default system channels
func InitSystemChannels(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Channel{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		channels := []*Channel{
			{
				AdminID:   SuperAdminID,
				AppID:     gofakeit.UUID(),
				SecretKey: gofakeit.UUID(),
				Name:      "Apple",
				Icon:      "/social/apple.png",
				Symbol:    "apple",
				Type:      ChannelTypeApple,
				Status:    ChannelStatusEnabled,
				Data:      ChannelData{},
			},
			{
				AdminID:   SuperAdminID,
				AppID:     gofakeit.UUID(),
				SecretKey: gofakeit.UUID(),
				Name:      "Google",
				Icon:      "/social/google.png",
				Symbol:    "google",
				Type:      ChannelTypeGoogle,
				Status:    ChannelStatusEnabled,
				Data:      ChannelData{},
			},
			{
				AdminID:   SuperAdminID,
				AppID:     gofakeit.UUID(),
				SecretKey: gofakeit.UUID(),
				Name:      "Telegram",
				Icon:      "/social/telegram.png",
				Symbol:    "telegram",
				Type:      ChannelTypeTelegram,
				Status:    ChannelStatusEnabled,
				Data:      ChannelData{},
			},
		}

		if err := db.Create(&channels).Error; err != nil {
			return err
		}
	}

	return nil
}
