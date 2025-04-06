package models

const (
	AccountTypeBalance int8 = 1  // 余额
	AccountTypeAsset   int8 = 11 // 资产

	AccountStatusDisabled int8 = -1 // 禁用
	AccountStatusEnabled  int8 = 10 // 启用
)

// WalletAccount 钱包账户管理
type WalletAccount struct {
	BaseModel
	AdminID   uint              `gorm:"type:int unsigned not null;comment:管理员ID" json:"adminId"`
	UserID    uint              `gorm:"type:int unsigned not null;comment:用户ID" json:"userId"`
	PaymentID uint              `gorm:"type:int unsigned not null;comment:支付ID" json:"paymentId"`
	Name      string            `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
	Type      int8              `gorm:"type:tinyint not null;default:1;comment:类型(1:余额,11:资产)" json:"type"`
	Status    int8              `gorm:"type:tinyint not null;default:10;comment:状态(-1:禁用,10:启用)" json:"status"`
	Data      WalletPaymentData `gorm:"type:json;comment:数据" json:"data"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&WalletAccount{}); err != nil {
//		panic("Failed to auto migrate WalletAccount table: " + err.Error())
//	}
//}
