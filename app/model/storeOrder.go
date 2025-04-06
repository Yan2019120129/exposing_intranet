package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreOrderTypeStore int8 = 1 // 商家订单
	StoreTypeDataIndex  int8 = 2 // 数据店铺列表

	// 订单状态常量
	StoreOrderStatusDisable         int8 = -1 // 订单取消
	StoreOrderStatusPending         int8 = 10 // 待支付
	StoreOrderStatusShipping        int8 = 12 // 待发货
	StoreOrderStatusShippingOutTime int8 = 13 // 发货超时
	StoreOrderStatusProgress        int8 = 14 // 待收货
	StoreOrderStatusComplete        int8 = 20 // 订单完成
)

// StoreOrder 店铺订单
type StoreOrder struct {
	BaseModel
	AdminID    uint            `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	AssetsID   uint            `gorm:"type:int unsigned not null;index;comment:资产ID" json:"assetsId"`
	StoreID    uint            `gorm:"type:int unsigned not null;index;comment:店铺ID" json:"storeId"`
	UserID     uint            `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	OrderSN    string          `gorm:"type:varchar(64);not null;index;comment:订单编号" json:"orderSn"`
	Type       int8            `gorm:"type:tinyint unsigned not null;default:1;comment:类型(1:商家订单)" json:"type"`
	Money      float64         `gorm:"type:decimal(20,6) not null;default:0;comment:购买总价" json:"money"`
	FinalMoney float64         `gorm:"type:decimal(10,2) not null;default:0;comment:实际总价" json:"finalMoney"`
	Earnings   float64         `gorm:"type:decimal(10,2) not null;default:0;comment:收益总额" json:"earnings"`
	Status     int8            `gorm:"type:tinyint not null;default:10;comment:状态(-1:取消,10:待支付,12:待发货,14:待收货,20:完成)" json:"status"`
	Data       *StoreOrderData `gorm:"type:json;comment:数据" json:"data"`
}

type StoreOrderData struct {
	StoreShippingAddress StoreShippingAddress
	LevelData            LevelData
}

// Value implements the driver.Valuer interface
func (d *StoreOrderData) Value() (driver.Value, error) {
	v, err := json.Marshal(d)
	return v, err
}

// Scan implements the sql.Scanner interface
func (d *StoreOrderData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreOrder{}); err != nil {
//		panic("Failed to auto migrate StoreOrder table: " + err.Error())
//	}
//}
