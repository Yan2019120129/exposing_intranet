package models

const (
	ExtendSettingGroupCategory = 1 // 分类
	ExtendSettingGroupProduct  = 2 // 产品
	ExtendSettingGroupStore    = 3 // 店铺
)

// ExtendSetting 扩展设置表
type ExtendSetting struct {
	BaseModel
	AdminID  uint   `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_source_field;comment:管理ID" json:"adminId"`
	SourceID uint   `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_source_field;comment:来源ID" json:"sourceId"`
	Name     string `gorm:"type:varchar(60) not null;comment:名称" json:"name"`
	GroupID  uint   `gorm:"type:int unsigned not null;index;comment:分组ID" json:"groupId"`
	Type     string `gorm:"type:varchar(50) not null;comment:input类型" json:"type"`
	Field    string `gorm:"type:varchar(50) not null;uniqueIndex:idx_admin_source_field;comment:字段名" json:"field"`
	Value    string `gorm:"type:text;comment:值" json:"value"`
}

// init 初始化
//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&ExtendSetting{}); err != nil {
//		panic("Failed to auto migrate Setting table: " + err.Error())
//	}
//}
