package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"time"
)

const (
	// Order Types
	ProductOrderTypeDefault int8 = 1 // 商家订单

	ProductOrderSideBuy  int8 = 1 // 买入
	ProductOrderSideSell int8 = 2 // 卖出

	ProductOrderModeLimit  int8 = 1 // 限价
	ProductOrderModeMarket int8 = 2 // 市价

	// Order Status
	ProductOrderStatusCancelled       int8 = -1 // 取消
	ProductOrderStatusPending         int8 = 10 //待付款
	ProductOrderStatusShipping        int8 = 12 //待发货
	ProductOrderStatusShippingOutTime int8 = 13 //发货超时
	ProductOrderStatusProgress        int8 = 14 //待收货
	ProductOrderStatusCompleted       int8 = 20 //订单完成
)

// Order 店铺规格订单
type Order struct {
	BaseModel
	AdminID      uint      `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID       uint      `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	ProductID    uint      `gorm:"type:int unsigned not null;index;comment:产品ID" json:"productId"`
	StoreID      uint      `gorm:"type:int unsigned not null;comment:店铺ID" json:"storeId"`
	StoreOrderID uint      `json:"storeOrderId" gorm:"type:int unsigned not null;common:店铺订单ID"`
	OrderSN      string    `gorm:"type:varchar(64);not null;uniqueIndex;comment:订单编号" json:"orderSn"`
	Money        float64   `gorm:"type:decimal(20,6) not null;comment:金额" json:"money"`
	Fee          float64   `gorm:"type:decimal(20,6) not null;comment:手续费" json:"fee"`
	Nums         uint      `gorm:"type:int not null;default:1;comment:数量" json:"nums"`
	FinalMoney   float64   `json:"finalMoney" gorm:"type:decimal(12,2) not null;default:0;comment:实际价格"`
	Earnings     float64   `json:"earnings" gorm:"type:decimal(12,2) not null;default:0;comment:订单收益"`
	Type         int8      `gorm:"type:tinyint unsigned not null;default:1;index;comment:类型(1:用户订单)" json:"type"`
	Side         int8      `gorm:"type:tinyint unsigned not null;default:1;index;comment:方向(1:买入,2:卖出)" json:"side"`
	Mode         int8      `gorm:"type:tinyint unsigned not null;default:1;index;comment:模式(1:限价,2:市价)" json:"mode"`
	Status       int8      `gorm:"type:tinyint not null;default:10;index;comment:订单状态(-1:订单取消,10:待付款,12:待发货,13:发货超时,14:待收货,20:订单完成)" json:"status"`
	Data         OrderData `gorm:"type:json;comment:数据" json:"data"`
	ExpiredAt    time.Time `gorm:"type:datetime(3);index;comment:过期时间" json:"expiredAt"`
}

type OrderData struct {
	ProductSkuAttribute ProductSkuAttribute
	LevelData           LevelData
}

// Value implements the driver.Valuer interface
func (d OrderData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *OrderData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Order{}); err != nil {
//		panic("Failed to auto migrate Order table: " + err.Error())
//	}
//}
