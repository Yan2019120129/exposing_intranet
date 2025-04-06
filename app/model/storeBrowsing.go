package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreBrowsingTypeStore   int8 = 1 // 商城
	StoreBrowsingTypeProduct int8 = 2 // 商品
)

// StoreBrowsing 店铺浏览记录
type StoreBrowsing struct {
	BaseModel
	AdminID   uint              `json:"adminId" gorm:"type:int unsigned not null;index;comment:管理员ID"`
	UserID    uint              `json:"userId" gorm:"type:int unsigned not null;index;comment:用户ID"`
	StoreID   uint              `json:"storeId" gorm:"type:int unsigned not null;index;comment:店铺ID"`
	ProductID uint              `json:"productId" gorm:"type:int unsigned not null;index;comment:商品ID"`
	Type      int8              `json:"type" gorm:"type:tinyint not null;default:1;comment:类型(1:商城)"`
	Status    int8              `json:"status" gorm:"type:tinyint not null;default:10;comment:状态(10:默认)"`
	Nums      int64             `json:"nums" gorm:"type:int not null;default:1;comment:次数"`
	Data      StoreBrowsingData `gorm:"type:json;comment:数据" json:"data"`
}

// Value implements the driver.Valuer interface
func (d StoreBrowsingData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *StoreBrowsingData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

// StoreBrowsingData 店铺浏览记录数据
type StoreBrowsingData struct{}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreBrowsing{}); err != nil {
//		panic("Failed to auto migrate StoreBrowsing table: " + err.Error())
//	}
//}
