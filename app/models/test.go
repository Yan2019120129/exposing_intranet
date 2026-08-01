package models

import "gorm.io/gorm"

func init() {
	ModelManage.SetModel("test", &Test{}, "示例")
}

// Test 表示用于演示接口流程的测试记录。
type Test struct {
	gorm.Model
	Name string `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
}

// TableName 返回测试记录对应的数据表名称。
func (Test) TableName() string {
	return "test"
}
