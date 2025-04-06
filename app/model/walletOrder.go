package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	// Order Types 订单类型
	WalletOrderTypeDeposit         int8 = 1  // 充值
	WalletOrderTypeWithdrawal      int8 = 2  // 提现
	WalletOrderTypeAssetDeposit    int8 = 11 // 资产充值
	WalletOrderTypeAssetWithdrawal int8 = 12 // 资产提现

	// Order Status 订单状态
	WalletOrderStatusRejected  int8 = -1 // 拒绝
	WalletOrderStatusPending   int8 = 10 // 审核中
	WalletOrderStatusCompleted int8 = 20 // 完成
)

// WalletOrder 钱包订单
type WalletOrder struct {
	BaseModel
	AdminID  uint            `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID   uint            `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	AssetsID uint            `gorm:"type:int unsigned not null;index;comment:资产ID" json:"assetsId"`
	SourceID uint            `gorm:"type:int unsigned not null;index;comment:来源ID" json:"sourceId"`
	Type     int8            `gorm:"type:tinyint not null;default:1;index;comment:类型(1:充值,2:提现,11:资产充值,12:资产提现)" json:"type"`
	OrderSN  string          `gorm:"type:varchar(60) not null;uniqueIndex;comment:订单编号" json:"orderSn"`
	Money    float64         `gorm:"type:decimal(16,4) not null;comment:金额" json:"money"`
	Fee      float64         `gorm:"type:decimal(16,4) not null;comment:手续费" json:"fee"`
	IsHint   int8            `gorm:"type:tinyint not null;default:1;comment:提示音(1:开启,2:关闭)" json:"isHint"`
	Status   int8            `gorm:"type:tinyint not null;default:10;index;comment:订单状态(-1:拒绝,10:审核中,20:完成)" json:"status"`
	Proof    string          `gorm:"type:varchar(255);comment:凭证" json:"proof"`
	Reason   string          `gorm:"type:varchar(255);comment:拒绝原因" json:"reason"`
	Data     WalletOrderData `gorm:"type:json;comment:数据" json:"data"`
}

// WalletOrderData 钱包订单数据
type WalletOrderData struct {
	WalletPayment
}

// Value implements the driver.Valuer interface
func (d WalletOrderData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *WalletOrderData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&WalletOrder{}); err != nil {
//		panic("Failed to auto migrate WalletOrder table: " + err.Error())
//	}
//}
