package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	// UserAuthMode constants
	UserAuthModeIdentity int8 = 1 // 身份认证
	UserAuthModeAddress  int8 = 2 // 地址认证

	// UserAuthType constants
	UserAuthTypeIDCard   int8 = 1 // 身份证
	UserAuthTypePassport int8 = 2 // 护照
	UserAuthTypeDriver   int8 = 3 // 驾驶证
	// UserAuthStatus constants
	UserAuthStatusRejected  int8 = -1 // 拒绝
	UserAuthStatusPending   int8 = 10 // 待审核
	UserAuthStatusCompleted int8 = 20 // 已通过
)

// UserAuth 实名认证
type UserAuth struct {
	BaseModel
	AdminID  uint         `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID   uint         `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	Mode     int8         `gorm:"type:tinyint not null;default:1;comment:模式(1:身份认证 2:地址认证)" json:"mode"`
	RealName string       `gorm:"type:varchar(50) not null;comment:证件姓名" json:"realName"`
	IDNumber string       `gorm:"type:varchar(50) not null;uniqueIndex;comment:证件号码" json:"idNumber"`
	Photo1   string       `gorm:"type:varchar(255) not null;comment:证件照正面" json:"photo1"`
	Photo2   string       `gorm:"type:varchar(255) not null;comment:证件照反面" json:"photo2"`
	Photo3   string       `gorm:"type:varchar(255);comment:手持证件照" json:"photo3"`
	Address  string       `gorm:"type:varchar(255);comment:详细地址" json:"address"`
	Type     int8         `gorm:"type:tinyint not null;default:1;comment:类型(1:身份证 2:护照 3:驾驶证)" json:"type"`
	Status   int8         `gorm:"type:tinyint not null;default:10;comment:状态(-1:拒绝 10:待审核 20:已通过)" json:"status"`
	Reason   string       `gorm:"type:varchar(255);comment:拒绝原因" json:"reason"`
	Data     UserAuthData `gorm:"type:json;comment:数据" json:"data"`
}

// UserAuthData 实名认证数据
type UserAuthData struct {
}

// Value implements the driver.Valuer interface
func (d UserAuthData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *UserAuthData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&UserAuth{}); err != nil {
//		panic("Failed to auto migrate UserAuth table: " + err.Error())
//	}
//}
