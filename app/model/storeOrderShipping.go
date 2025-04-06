package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	LogisticsInfoTypeDefault = 1 //默认类型

	LogisticsInfoStatusShipping  = 10 // 运输中
	LogisticsInfoStatusDelivered = 20 // 已送达
	LogisticsInfoStatusFailed    = 30 // 配送失败
	LogisticsInfoStatusReturning = 40 // 退回中
)

// StoreOrderShipping 店铺船运
type StoreOrderShipping struct {
	BaseModel
	AdminID      uint              `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID       uint              `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	StoreID      uint              `gorm:"type:int unsigned not null;comment:店铺ID" json:"storeId"`
	OrderID      uint              `gorm:"type:int unsigned not null;common:订单ID" json:"orderId"`
	ProductID    uint              `gorm:"type:int unsigned not null;comment:商品ID" json:"productId"`
	SerialNumber string            `gorm:"type:varchar(64);not null;index;common:物流编号" json:"serialNumber"`
	Form         string            `gorm:"type:varchar(64);not null;index;common:发货地址" json:"form"`
	To           string            `gorm:"type:varchar(64);not null;index;common:目标地址" json:"to"`
	Location     string            `gorm:"type:varchar(64);not null;index;common:当前位置" json:"currentLocation"`
	Latitude     string            `gorm:"type:varchar(64);not null;index;common:经纬度" json:"latitude"`
	Type         int8              `gorm:"type:tinyint not null;default:1;common:类型(1:默认类型)" json:"type"`
	Status       int8              `gorm:"type:tinyint not null;default:10;index;comment:状态(10:运输中,20:已送达,30:配送失败,40:退回中)" json:"status"`
	Data         OrderShippingData `gorm:"type:json;comment:数据" json:"data"`
}

type OrderShippingData struct{}

// Value implements the driver.Valuer interface
func (d OrderShippingData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *OrderShippingData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreOrderShipping{}); err != nil {
//		panic("Failed to auto migrate StoreOrderShipping table: " + err.Error())
//	}
//}
