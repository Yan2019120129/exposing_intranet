package models

// Setting 用户设置
type Setting struct {
	BaseModel
	AdminID uint   `gorm:"type:int unsigned not null;comment:管理ID" json:"adminId"`
	UserID  uint   `gorm:"type:int unsigned not null;uniqueIndex:idx_user_field;index;comment:用户ID" json:"userId"`
	Name    string `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
	Type    string `gorm:"type:varchar(50) not null;comment:类型" json:"type"`
	Field   string `gorm:"type:varchar(50) not null;uniqueIndex:idx_user_field;comment:字段名" json:"field"`
	Value   string `gorm:"type:text;comment:字段值" json:"value"`
}

// SettingData 用户设置数据
type SettingData struct {
	Field string           `json:"field"`
	Type  string           `json:"type"`
	Data  AdminSettingData `json:"data"`
}

//// AddChildrenForm 添加子表单
//func (s *SettingData) AddChildrenForm(field string, form *views.Form) {
//	s.Data.ChildrenForm[field] = form
//}
//
//// NewSettingInput 创建用户设置输入
//func NewSettingInput(settingField string, input views.Input, childrenForm *views.Form) *SettingData {
//	input.Field = "settingValue"
//	input.ScanField = "settingValues." + settingField
//	return &SettingData{
//		Field: settingField,
//		Type:  input.Type,
//		Data: AdminSettingData{
//			Input: input,
//			ChildrenForm: map[string]*views.Form{
//				input.Field: childrenForm,
//			},
//		},
//	}
//}

//// init 初始化
//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Setting{}); err != nil {
//		panic("Failed to auto migrate Setting table: " + err.Error())
//	}
//}
