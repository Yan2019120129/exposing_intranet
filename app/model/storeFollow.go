package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreFollowTypeStore   = 1 // 关注店铺
	StoreFollowTypeProduct = 2 // 收藏商品

	StoreFollowStatusCancels = -1 // 取消
	StoreFollowStatusConcern = 10 // 关注
)

// StoreFollow 店铺关注表
type StoreFollow struct {
	BaseModel
	AdminID   uint            `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID    uint            `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	StoreID   uint            `gorm:"type:int unsigned not null;index;comment:店铺ID" json:"storeId"`
	ProductID uint            `gorm:"type:int unsigned not null;index;comment:商品ID" json:"productId"`
	Type      int8            `gorm:"type:tinyint not null;default:1;comment:类型(1:关注店铺,2:收藏商品)" json:"type"`
	Status    int8            `gorm:"type:tinyint not null;default:10;comment:状态(-1:取消,10:关注)" json:"status"`
	Data      StoreFollowData `gorm:"type:json;comment:数据" json:"data"`
}

// StoreFollowData 店铺关注数据
type StoreFollowData struct {
}

// Value implements the driver.Valuer interface
func (d StoreFollowData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *StoreFollowData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreFollow{}); err != nil {
//		panic("Failed to auto migrate StoreOrder table: " + err.Error())
//	}
//}
