package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
)

const (
	UserLevelStatusDisabled int8 = -1 // 禁用
	UserLevelStatusExpired  int8 = 10 // 过期
	UserLevelStatusEnabled  int8 = 20 // 激活
)

// UserLevel 会员表
type UserLevel struct {
	BaseModel
	AdminID   uint          `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID    uint          `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	LevelID   uint          `gorm:"type:int unsigned not null;index;comment:等级ID" json:"levelId"`
	Status    int8          `gorm:"type:tinyint not null;default:20;index;comment:状态(-1:禁用,10:过期,20:激活)" json:"status"`
	Data      UserLevelData `gorm:"type:json;comment:等级信息" json:"data"`
	ExpiredAt time.Time     `gorm:"type:datetime(3);index;comment:过期时间" json:"expiredAt"`
}

// UserLevelData 用户等级数据
type UserLevelData struct {
	Level
}

// Scan implements sql.Scanner interface
func (u *UserLevelData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &u)
}

// Value implements driver.Valuer interface
func (u UserLevelData) Value() (driver.Value, error) {
	return json.Marshal(u)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&UserLevel{}); err != nil {
//		panic("Failed to auto migrate UserLevel table: " + err.Error())
//	}
//}
