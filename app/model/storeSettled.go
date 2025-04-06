package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreSettledStatusRefuse  int8 = -1 // 拒绝
	StoreSettledStatusPending int8 = 10 // 审核中
	StoreSettledStatusPass    int8 = 20 // 通过

	StoreSettledTypeIdCard  int8 = 1 // 身份证
	StoreSettledTypeLicense int8 = 2 // 营业执照
)

// StoreSettled 店铺入驻
type StoreSettled struct {
	BaseModel
	AdminID uint             `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID  uint             `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	Type    int8             `gorm:"type:tinyint not null;default:2;comment:类型(1:身份证,2:营业执照)" json:"type"`
	Name    string           `gorm:"type:varchar(64) not null;comment:店铺名称" json:"name"`
	Address string           `gorm:"type:varchar(2048);comment:店铺地址" json:"address"`
	Logo    string           `gorm:"type:varchar(255);comment:店铺Logo" json:"logo"`
	Photo1  string           `gorm:"type:varchar(120) not null;comment:证件照1" json:"photo1"`
	Photo2  string           `gorm:"type:varchar(120) not null;comment:证件照2" json:"photo2"`
	Photo3  string           `gorm:"type:varchar(120) not null;comment:证件照3" json:"photo3"`
	Number  string           `gorm:"type:varchar(255);comment:证件号" json:"number"`
	Reason  string           `gorm:"type:varchar(255);comment:拒绝原因" json:"reason"`
	Status  int8             `gorm:"type:tinyint not null;default:10;comment:状态(-1:拒绝,10:审核中,20:通过)" json:"status"`
	Data    StoreSettledData `gorm:"type:json;comment:数据" json:"data"`
}

type StoreSettledData struct{}

// Value implements the driver.Valuer interface
func (d StoreSettledData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *StoreSettledData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreSettled{}); err != nil {
//		panic("Failed to auto migrate StoreSettled table: " + err.Error())
//	}
//}
