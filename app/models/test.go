package models

func init() {
	ModelManage.SetModel("test", &Test{}, "示例")
}

type Test struct {
	BaseModel
	Name string `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
}
