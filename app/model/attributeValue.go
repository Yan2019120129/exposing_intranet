package models

// AttributeValue 属性值
type AttributeValue struct {
	BaseModel
	AttributeID uint   `gorm:"type:int unsigned not null;index;comment:属性ID" json:"attributeId"`
	Name        string `gorm:"type:varchar(255) not null;comment:名称" json:"name"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AttributeValue{}); err != nil {
//		panic("Failed to auto migrate attribute value table: " + err.Error())
//	}
//}
