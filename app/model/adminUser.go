package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"my-base/utils"
	"time"
)

const (
	// AdminUserStatusEnabled 表示管理员用户处于激活状态
	AdminUserStatusEnabled int8 = 10
	// AdminUserStatusDisabled 表示管理员用户处于禁用状态
	AdminUserStatusDisabled int8 = -1

	// SuperAdminID 表示超级管理员的ID
	SuperAdminID uint = 1

	// DefaultMerchantID 表示默认商户的ID
	DefaultMerchantID uint = 2
)

// AdminUser 管理表
type AdminUser struct {
	BaseModel
	ParentID      uint          `gorm:"type:int unsigned;index;comment:上级ID" json:"parentId"`
	Username      string        `gorm:"uniqueIndex;type:varchar(60);comment:用户名" json:"username"`
	Nickname      string        `gorm:"type:varchar(60);comment:昵称" json:"nickname"`
	Email         string        `gorm:"type:varchar(60);uniqueIndex;comment:邮箱;default:null" json:"email"`
	Telephone     string        `gorm:"type:varchar(50);uniqueIndex;comment:号码;default:null" json:"telephone"`
	Avatar        string        `gorm:"type:varchar(255);default:'';comment:头像" json:"avatar"`
	Password      string        `gorm:"type:varchar(255);comment:密码" json:"-"`
	SecurityKey   string        `gorm:"type:varchar(255);comment:密钥" json:"-"`
	Money         float64       `gorm:"type:decimal(12,2);default:0.00;comment:金额" json:"money"`
	Domains       string        `gorm:"type:varchar(2046);comment:域名" json:"domains"`
	ChatURL       string        `gorm:"type:varchar(255);comment:客服链接" json:"chatUrl"`
	Status        int8          `gorm:"type:tinyint;default:10;comment:状态(10:启用,-1:禁用)" json:"status"`
	Data          AdminUserData `gorm:"type:json;comment:数据" json:"data,omitempty"`
	ExpiredAt     time.Time     `gorm:"type:timestamp;comment:过期时间" json:"expiredAt"`
	LastLoginTime time.Time     `gorm:"type:timestamp;autoCreateTime;comment:最后登录时间" json:"lastLoginTime"`
	LastLoginIP   string        `gorm:"type:varchar(45);comment:最后登录IP" json:"lastLoginIP"`
}

// ToDisplayData 将 AdminUser 转换为 AdminUserDisplayData
func (au *AdminUser) ToDisplayData() AdminUserDisplayData {
	return AdminUserDisplayData{
		ID:            au.ID,
		ParentID:      au.ParentID,
		Username:      au.Username,
		Nickname:      au.Nickname,
		Email:         au.Email,
		Telephone:     au.Telephone,
		Avatar:        au.Avatar,
		Money:         au.Money,
		Domains:       au.Domains,
		ChatURL:       au.ChatURL,
		Status:        au.Status,
		Data:          au.Data,
		ExpiredAt:     au.ExpiredAt,
		LastLoginTime: au.LastLoginTime,
		LastLoginIP:   au.LastLoginIP,
	}
}

// AdminUserDisplayData 管理用户显示数据
type AdminUserDisplayData struct {
	ID            uint          `json:"id" column:"label:ID"`
	ParentID      uint          `json:"parentId"`
	Avatar        string        `json:"avatar" column:"label:头像;type:icon"`
	Username      string        `json:"username" column:"label:账户"`
	Nickname      string        `json:"nickname" column:"label:昵称"`
	Email         string        `json:"email"`
	Telephone     string        `json:"telephone"`
	Money         float64       `json:"money" column:"label:金额"`
	Domains       string        `json:"domains" column:"label:前台域名"`
	ChatURL       string        `json:"chatUrl" column:"label:客服链接"`
	Status        int8          `json:"status" column:"label:状态;"`
	Data          AdminUserData `json:"data"`
	ExpiredAt     time.Time     `json:"expiredAt" column:"label:过期时间;type:date"`
	UpdatedAt     time.Time     `json:"updatedAt" column:"label:活跃时间;type:date"`
	LastLoginTime time.Time     `json:"lastLoginTime"`
	LastLoginIP   string        `json:"lastLoginIP"`
}

// AdminUserData 管理数据Json
type AdminUserData struct {
	Auth AdminUserDataAuth `json:"auth"` //	Google验证器配置
}

// AdminUserDataAuth Google验证器配置
type AdminUserDataAuth struct {
	Key    string `json:"key" views:"label:Google验证器Key"`        //	Google验证器Key
	Enable bool   `json:"enable" views:"label:是否开启;type:toggle"` //	是否开启
}

// Value implements the driver.Valuer interface
func (d AdminUserData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *AdminUserData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AdminUser{}); err != nil {
//		panic("Failed to auto migrate AdminUser table: " + err.Error())
//	}
//
//	// Initialize admin users
//	if err := InitAdminUsers(db); err != nil {
//		panic("Failed to initialize admin users: " + err.Error())
//	}
//}

// InitAdminUsers 初始化管理员数据
func InitAdminUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AdminUser{}).Count(&count).Error; err != nil {
		return err
	}

	adminKey, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "八戒科技",
		AccountName: "admin",
	})

	merchantKey, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "八戒科技",
		AccountName: "merchant",
	})

	agentKey, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "八戒科技",
		AccountName: "agent",
	})

	if count == 0 {
		adminUsers := []AdminUser{
			{
				Username:    "admin",
				Nickname:    "Admin",
				Avatar:      "/avatar.png",
				Email:       "admin@example.com",
				Password:    utils.EncryptPassword("taozijun"),
				SecurityKey: utils.EncryptPassword("taozijun"),
				Status:      AdminUserStatusEnabled,
				ExpiredAt:   time.Now().Add(365 * 24 * time.Hour),
				Data: AdminUserData{
					Auth: AdminUserDataAuth{
						Key:    adminKey.Secret(),
						Enable: false,
					},
				},
			},
			{
				ParentID:    SuperAdminID,
				Username:    "merchant",
				Nickname:    "Merchant",
				Avatar:      "/avatar.png",
				Email:       "merchant@example.com",
				Password:    utils.EncryptPassword("abc123"),
				SecurityKey: utils.EncryptPassword("abc123"),
				Status:      AdminUserStatusEnabled,
				Domains:     "localhost:9100,192.168.5.10:9100",
				ExpiredAt:   time.Now().Add(365 * 24 * time.Hour),
				Data: AdminUserData{
					Auth: AdminUserDataAuth{
						Key:    merchantKey.Secret(),
						Enable: false,
					},
				},
			},
			{
				ParentID:    DefaultMerchantID,
				Username:    "agent",
				Nickname:    "Agent",
				Avatar:      "/avatar.png",
				Email:       "agent@example.com",
				Password:    utils.EncryptPassword("abc123"),
				SecurityKey: utils.EncryptPassword("abc123"),
				Status:      AdminUserStatusEnabled,
				ExpiredAt:   time.Now().Add(365 * 24 * time.Hour),
				Data: AdminUserData{
					Auth: AdminUserDataAuth{
						Key:    agentKey.Secret(),
						Enable: false,
					},
				},
			},
		}

		for _, user := range adminUsers {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
