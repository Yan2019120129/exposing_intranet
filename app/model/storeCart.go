package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

// StoreCart 店铺购物车
type StoreCart struct {
	BaseModel
	AdminID   uint          `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID    uint          `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	StoreID   uint          `gorm:"type:int unsigned not null;index;comment:店铺ID" json:"storeId"`
	ProductID uint          `gorm:"type:int unsigned not null;index;comment:商品ID" json:"productId"`
	SkuID     uint          `gorm:"type:int unsigned not null;index;comment:规格ID" json:"skuId"`
	Nums      uint          `gorm:"type:smallint unsigned not null;default:1;comment:数量" json:"nums"`
	Data      StoreCartData `gorm:"type:json;comment:数据" json:"data"`
}

// StoreCartInfo 购物车店铺信息
type StoreCartInfo struct {
	Store
	CartList []*CartInfo `json:"cartList" gorm:"foreignKey:StoreID"`
}

func (StoreCartInfo) TableName() string {
	return "store"
}

// CartInfo 购物车信息
type CartInfo struct {
	StoreCart
	ProductInfo Product             `json:"productInfo" gorm:"foreignKey:ProductID"`
	SkuInfo     ProductSkuAttribute `json:"skuInfo" gorm:"foreignKey:SkuID"`
}

func (CartInfo) TableName() string {
	return "store_cart"
}

// StoreCartData 店铺购物车数据
type StoreCartData struct{}

// Value implements the driver.Valuer interface
func (d StoreCartData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *StoreCartData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreCart{}); err != nil {
//		panic("Failed to auto migrate StoreCart table: " + err.Error())
//	}
//}
