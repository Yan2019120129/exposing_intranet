package models

const (
	ShippingAddressStatusDisabled = -1 // 禁用
	ShippingAddressStatusActivate = 10 // 激活

	ShippingAddressIsShowYes = 2 // 默认地址
	ShippingAddressIsShowNo  = 1 // 非默认地址

	ShippingAddressTypeReceiving = 1 // 收货地址
	ShippingAddressTypeShipments = 2 // 发货地址
)

// StoreShippingAddress 购物地址
type StoreShippingAddress struct {
	BaseModel
	AdminID        uint   `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID         uint   `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	Name           string `gorm:"type:varchar(50) not null;comment:收件人名称" json:"name"`
	Contact        string `gorm:"type:varchar(50) not null;comment:联系方式" json:"contact"`
	City           string `gorm:"type:varchar(255) not null;comment:城市" json:"city"`
	Address        string `gorm:"type:varchar(255) not null;comment:详细地址" json:"address"`
	Type           int8   `gorm:"type:tinyint not null;default:1;comment:类型(1:收货地址,2:发货地址)" json:"type"`
	Status         int8   `gorm:"type:tinyint not null;default:10;comment:状态(-1:禁用,10:激活)" json:"status"`
	DefaultAddress int8   `gorm:"type:tinyint not null;default:1;comment:(1:非默认,2:默认)" json:"defaultAddress"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreShippingAddress{}); err != nil {
//		panic("Failed to auto migrate StoreShippingAddress table: " + err.Error())
//	}
//}
