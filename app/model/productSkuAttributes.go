package models

// ProductSkuAttributes SKU属性关联表
type ProductSkuAttributes struct {
	BaseModel
	AdminID          uint `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	SkuID            uint `gorm:"type:int unsigned not null;index;comment:SKU ID" json:"skuId"`
	AttributeID      uint `gorm:"type:int unsigned not null;index;comment:属性ID" json:"attributeId"`
	AttributeValueID uint `gorm:"type:int unsigned not null;index;comment:属性值ID" json:"attributeValueId"`
}

// SkuAttribute SKU属性
type SkuAttribute struct {
	AttributeName  string
	AttributeValue string
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&ProductSkuAttributes{}); err != nil {
//		panic("Failed to auto migrate ProductSkuAttributes table: " + err.Error())
//	}
//}
