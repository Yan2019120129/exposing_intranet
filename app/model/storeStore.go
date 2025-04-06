package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreStatusDisabled int8 = -1 // 禁用
	StoreStatusActivate int8 = 10 // 启用

	StoreTypeDefault int8 = 1 // 普通店铺
)

// Store 店铺
type Store struct {
	BaseModel
	AdminID  uint      `gorm:"type:int unsigned not null;comment:管理员ID" json:"adminId"`
	UserID   uint      `gorm:"type:int unsigned not null;comment:用户ID" json:"userId"`
	Logo     string    `gorm:"type:varchar(255);comment:店铺Logo" json:"logo"`
	Name     string    `gorm:"type:varchar(255) not null;comment:店铺名称" json:"name"`
	Contact  string    `gorm:"type:varchar(50);comment:联系方式" json:"contact"`
	Type     int8      `gorm:"type:tinyint not null;default:1;comment:店铺类型(1:普通店铺)" json:"type"`
	Keywords string    `gorm:"type:varchar(255);index;comment:关键词" json:"keywords"`
	Address  string    `gorm:"type:varchar(2048);comment:店铺地址" json:"address"`
	Desc     string    `gorm:"type:varchar(255);comment:店铺描述" json:"desc"`
	Rating   float64   `gorm:"type:decimal(3,2) not null;default:5;comment:评分" json:"rating"`
	Score    int       `gorm:"type:tinyint not null;default:100;comment:信用分" json:"score"`
	Sales    uint      `gorm:"type:int unsigned not null;default:0;comment:销售量" json:"sales"`
	Visits   uint      `gorm:"type:int unsigned not null;default:0;comment:访问量" json:"visits"`
	Status   int8      `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:激活)" json:"status"`
	Data     StoreData `gorm:"type:json;comment:数据" json:"data"`
}

type StoreData struct{}

// Value implements the driver.Valuer interface
func (d StoreData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *StoreData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Store{}); err != nil {
//		panic("Failed to auto migrate Store table: " + err.Error())
//	}
//}
