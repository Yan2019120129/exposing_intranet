package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	StoreRefundStatusRefuse   int8 = -1 // 拒绝
	StoreRefundStatusPending  int8 = 10 // 申请
	StoreRefundStatusComplete int8 = 20 // 同意
)

// StoreOrderRefund 订单售后
type StoreOrderRefund struct {
	BaseModel
	AdminID   uint            `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID    uint            `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	OrderID   uint            `gorm:"type:int unsigned not null;index;comment:订单ID" json:"orderId"`
	StoreID   uint            `gorm:"type:int unsigned not null;index;comment:店铺ID" json:"storeId"`
	ProductID uint            `gorm:"type:int unsigned not null;index;comment:产品ID" json:"productId"`
	Reason    string          `gorm:"type:varchar(512) not null;comment:申请理由" json:"name"`
	Images    Strings         `gorm:"type:varchar(4096) not null;comment:凭证图片" json:"images"`
	Amount    float64         `gorm:"type:decimal(12,2) not null;default:0;comment:金额" json:"money"`
	Status    int8            `gorm:"type:tinyint not null;default:10;comment:状态(-1:拒绝,10:申请,20:同意)" json:"status"`
	Reject    string          `gorm:"type:varchar(512);comment:拒绝原因" json:"reject"`
	Data      OrderRefundData `gorm:"type:json;comment:数据" json:"data"`
}

type OrderRefundData struct{}

// Value implements the driver.Valuer interface
func (d OrderRefundData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *OrderRefundData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreOrderRefund{}); err != nil {
//		panic("Failed to auto migrate StoreOrderRefund table: " + err.Error())
//	}
//}
