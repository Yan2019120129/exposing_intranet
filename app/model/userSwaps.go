package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	SwapsTypeAssets int8 = 1 // 资产交换

	SwapsStatusRejected  int8 = -1 // 拒绝
	SwapsStatusPending   int8 = 10 // 处理中
	SwapsStatusCompleted int8 = 20 // 完成
)

// Swaps 资产交换
type Swaps struct {
	BaseModel
	AdminID         uint      `gorm:"type:int unsigned;index;comment:管理ID" json:"adminId"`
	UserID          uint      `gorm:"type:int unsigned;index;comment:用户ID" json:"userId"`
	SendAssetsID    uint      `gorm:"type:int unsigned;index;comment:发送资产ID" json:"sendAssetsId"`
	ReceiveAssetsID uint      `gorm:"type:int unsigned;index;comment:接收资产ID" json:"receiveAssetsId"`
	SendAmount      float64   `gorm:"type:decimal(16,4);comment:发送金额" json:"sendAmount"`
	ReceiveAmount   float64   `gorm:"type:decimal(16,4);comment:接收金额" json:"receiveAmount"`
	Rate            float64   `gorm:"type:decimal(16,4);comment:汇率" json:"rate"`
	Fee             float64   `gorm:"type:decimal(16,4) not null;default:0;comment:手续费" json:"fee"`
	Type            int8      `gorm:"type:tinyint;default:1;index;comment:类型(1:资产)" json:"type"`
	Status          int8      `gorm:"type:tinyint;default:20;index;comment:状态(-1:已拒绝,10:待处理,20:已完成)" json:"status"`
	Data            SwapsData `gorm:"type:json;comment:数据" json:"data"`
}

// SwapsData 资产交换数据
type SwapsData struct {
}

// Value implements the driver.Valuer interface
func (d SwapsData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *SwapsData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Swaps{}); err != nil {
//		panic("Failed to auto migrate Swaps table: " + err.Error())
//	}
//}
