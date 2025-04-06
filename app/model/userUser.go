package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"time"
)

const (
	// Sex constants
	SexUnknown int8 = -1 // 未知性别
	SexMale    int8 = 1  // 男性
	SexFemale  int8 = 2  // 女性

	// User type constants
	UserTypeVirtual int8 = -1 // 虚拟用户
	UserTypeDefault int8 = 1  // 普通用户

	// User status constants
	UserStatusFrozen int8 = -1 // 冻结状态
	UserStatusActive int8 = 10 // 激活状态
)

// User 用户表
type User struct {
	BaseModel
	AdminID         uint      `gorm:"type:int unsigned not null;default:1;index;comment:管理ID;uniqueIndex:idx_admin_username;uniqueIndex:idx_admin_email;uniqueIndex:idx_admin_telephone" json:"adminId"`
	ChannelID       uint      `gorm:"type:int unsigned not null;index;comment:渠道ID" json:"channelId"`
	ParentID        uint      `gorm:"type:int unsigned not null;index;comment:父级ID" json:"parentId"`
	CountryID       uint      `gorm:"type:int unsigned not null;index;comment:国家ID" json:"countryId"`
	Username        string    `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_username;comment:用户名" json:"username"`
	Nickname        string    `gorm:"type:varchar(60) not null;index;comment:昵称" json:"nickname"`
	Email           string    `gorm:"type:varchar(60);uniqueIndex:idx_admin_email;comment:邮箱;default:null" json:"email"`
	Telephone       string    `gorm:"type:varchar(50);uniqueIndex:idx_admin_telephone;comment:手机号码;default:null" json:"telephone"`
	Lang            string    `gorm:"type:varchar(60) not null;index;comment:语言;" json:"lang"`
	Avatar          string    `gorm:"type:varchar(255) not null;comment:头像" json:"avatar"`
	Score           int       `gorm:"type:int not null;default:100;comment:信用分" json:"score"`
	Sex             int8      `gorm:"type:tinyint not null;default:-1;comment:性别(-1:未知,1:男,2:女)" json:"sex"`
	Birthday        time.Time `gorm:"type:date;comment:生日" json:"birthday"`
	Password        string    `gorm:"type:varchar(255) not null;comment:密码" json:"-"`
	SecurityKey     string    `gorm:"type:varchar(255) not null;comment:密钥" json:"-"`
	FrozenAmount    float64   `gorm:"type:decimal(16,2) not null;default:0;comment:冻结金额" json:"frozenAmount"`
	AvailableAmount float64   `gorm:"type:decimal(16,2) not null;default:0;comment:可用金额" json:"availableAmount"`
	InviteCode      string    `gorm:"type:varchar(20);uniqueIndex;comment:邀请码" json:"inviteCode"`
	Type            int8      `gorm:"type:tinyint not null;default:1;index;comment:类型(-1:虚拟用户,1:普通用户,10:渠道用户)" json:"type"`
	Status          int8      `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:冻结,10:激活)" json:"status"`
	LastLoginAt     time.Time `gorm:"type:timestamp;comment:登录时间" json:"lastLoginAt"`
	LastLoginIP     string    `gorm:"type:varchar(45);comment:登录IP" json:"lastLoginIP"`
	Desc            string    `gorm:"type:text;comment:详情" json:"desc"`
	Data            UserData  `gorm:"type:json;comment:数据" json:"data"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID                uint          `json:"id"`                         // 用户ID
	CountryID         uint          `json:"countryId"`                  // 国家ID
	Username          string        `json:"username"`                   // 用户名
	Nickname          string        `json:"nickname"`                   // 昵称
	Email             string        `json:"email"`                      // 邮箱
	Telephone         string        `json:"telephone"`                  // 电话号码
	Avatar            string        `json:"avatar"`                     // 头像
	Lang              string        `json:"lang"`                       // 语言
	Score             int           `json:"score"`                      // 信用分
	MinScore          int           `json:"minScore" gorm:"-"`          // 最低信用分
	Sex               int8          `json:"sex"`                        // 性别
	Birthday          time.Time     `json:"birthday"`                   // 生日
	FrozenAmount      float64       `json:"frozenAmount"`               // 冻结金额
	AvailableAmount   float64       `json:"availableAmount"`            // 可用金额
	InviteCode        string        `json:"inviteCode"`                 // 邀请码
	Desc              string        `json:"desc"`                       // 详情
	CreatedAt         time.Time     `json:"createdAt"`                  // 注册时间
	LastLoginAt       time.Time     `json:"lastLoginAt"`                // 最后登录时间
	LastLoginIP       string        `json:"lastLoginIP"`                // 最后登录IP
	Type              int8          `json:"type"`                       // 类型
	Status            int8          `json:"status"`                     // 状态
	AuthStatus        int8          `json:"authStatus" gorm:"-"`        // 认证状态
	AuthAddressStatus int8          `json:"authAddressStatus" gorm:"-"` // 地址认证状态
	LevelInfo         UserLevelInfo `json:"levelInfo" gorm:"-"`         // 等级信息
	SettledInfo       StoreSettled  `json:"settledInfo" gorm:"-"`       // 入驻信息
}

// UserLevelInfo 用户等级信息
type UserLevelInfo struct {
	Name      string    `json:"name"`      // 名称
	Icon      string    `json:"icon"`      // 图标
	Symbol    int8      `json:"symbol"`    // 符号
	ExpiredAt time.Time `json:"expiredAt"` // 过期时间
}

// ToUserInfo 转换为用户信息
func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:                u.ID,
		CountryID:         u.CountryID,
		Username:          u.Username,
		Nickname:          u.Nickname,
		Email:             u.Email,
		Telephone:         u.Telephone,
		Avatar:            u.Avatar,
		Lang:              u.Lang,
		Score:             u.Score,
		Sex:               u.Sex,
		Birthday:          u.Birthday,
		FrozenAmount:      u.FrozenAmount,
		AvailableAmount:   u.AvailableAmount,
		InviteCode:        u.InviteCode,
		Desc:              u.Desc,
		AuthStatus:        UserAuthStatusPending,
		AuthAddressStatus: UserAuthStatusPending,
		LevelInfo:         UserLevelInfo{},
		CreatedAt:         u.CreatedAt,
		LastLoginAt:       u.LastLoginAt,
		LastLoginIP:       u.LastLoginIP,
		Type:              u.Type,
		Status:            u.Status,
	}
}

// UserData 用户数据
type UserData struct {
	Address     string `json:"address"`     // 地址
	IDCardNo    string `json:"idCardNo"`    // 身份证号
	Occupation  string `json:"occupation"`  // 职业
	Education   string `json:"education"`   // 教育程度
	Interests   string `json:"interests"`   // 兴趣爱好
	SocialMedia string `json:"socialMedia"` // 社交媒体账号
}

// Value implements the driver.Valuer interface
func (d UserData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *UserData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist and set auto increment start value
//	if err := db.Set("gorm:table_options", "AUTO_INCREMENT=34567891").AutoMigrate(&User{}); err != nil {
//		panic("Failed to auto migrate User table: " + err.Error())
//	}
//}
