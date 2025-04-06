package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	TransferTypeInternal int8 = 1 // 内部转账

	TransferStatusPending   int8 = 10 // 待处理
	TransferStatusCompleted int8 = 20 // 已完成
	TransferStatusRejected  int8 = -1 // 已拒绝
)

// Transfer 站内转账
type Transfer struct {
	BaseModel
	AdminID    uint         `gorm:"type:int unsigned;index;comment:管理ID" json:"adminId"`
	SenderID   uint         `gorm:"type:int unsigned not null;index;comment:发送者ID" json:"senderId"`
	ReceiverID uint         `gorm:"type:int unsigned not null;index;comment:接收者ID" json:"receiverId"`
	AssetsID   uint         `gorm:"type:int unsigned not null;index;comment:资产ID" json:"assetsId"`
	Amount     float64      `gorm:"type:decimal(16,4) not null;comment:转账金额" json:"amount"`
	Type       int8         `gorm:"type:tinyint not null;default:1;index;comment:转账类型(1:内部转账)" json:"type"`
	Status     int8         `gorm:"type:tinyint not null;default:20;index;comment:状态(10:待处理,20:已完成,-1:已拒绝)" json:"status"`
	Remark     string       `gorm:"type:varchar(255);comment:备注" json:"remark"`
	Data       TransferData `gorm:"type:json;comment:数据" json:"data"`
}

// TransferData 转账数据
type TransferData struct {
}

// Value implements the driver.Valuer interface
func (d TransferData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *TransferData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Transfer{}); err != nil {
//		panic("Failed to auto migrate Transfer table: " + err.Error())
//	}
//}
