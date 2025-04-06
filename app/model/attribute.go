package models

// Attribute 属性
type Attribute struct {
	BaseModel
	AdminID uint   `gorm:"type:int unsigned not null;default:0;comment:管理员ID" json:"adminId"`
	Name    string `gorm:"type:varchar(255) not null;comment:名称" json:"name"`
	Type    int8   `gorm:"type:tinyint not null;default:1;comment:类型(1:默认)" json:"type"`
}

// ProductAttrsKeyVal 产品属性Key
type ProductAttrsKeyVal struct {
	Attribute
	Values []*AttributeValue `json:"values" gorm:"foreignKey:AttributeID;references:ID"`
}

func (_Attribute *Attribute) TableName() string {
	return "attribute"
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Attribute{}); err != nil {
//		panic("Failed to auto migrate Attribute table: " + err.Error())
//	}
//
//}
