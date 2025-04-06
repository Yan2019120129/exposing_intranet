package models

const (
	MigrationStatusDisabled = -1
	MigrationStatusEnabled  = 10
)

// 数据库的变动记录
type Migrate struct {
	BaseModel
	Name        string `gorm:"type:varchar(255) not null;comment:名称" json:"name"`                            // 名称
	Version     string `gorm:"type:varchar(255) not null;comment:版本" json:"version"`                         // 版本
	Status      int8   `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"` // 状态
	Description string `gorm:"type:text;comment:描述" json:"description"`                                      // 描述
}

//func init() {
//	if err := db.AutoMigrate(&Migrate{}); err != nil {
//		panic("Failed to auto migrate Migrate table: " + err.Error())
//	}
//}
