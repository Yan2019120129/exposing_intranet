package view

import "github.com/gin-gonic/gin"

type TypeView string

const (
	TypeViewTable TypeView = "table"
	TypeViewForm  TypeView = "form"
)

type View struct {
	Name string   `json:"name"` // 视图名称
	Desc string   `json:"desc"` // 视图描述
	Type TypeView `json:"type"` // 视图类型
	Body Body     `json:"body"` // 主体
}

// NewView 创建视图
func NewView(name string) *View {
	return &View{
		Name: name,
	}
}

// SetDesc 视图描述
func (v *View) SetDesc(desc string) *View {
	v.Desc = desc
	return v
}

// BodyTable 创建Table主体
func (v *View) BodyTable(model any) *Table {
	table := NewTable(model)
	v.Type = TypeViewTable
	v.Body = table
	return table
}

// Register 注册视图
func (v *View) Register(e *gin.Engine) {
	operates := v.Body.GetOperates()
	for _, operate := range operates {
		switch operate.AskType {
		case AskTypePost:
			e.POST(operate.Url, operate.fun)
		case AskTypeGet:
			e.GET(operate.Url, operate.fun)
		case AskTypePut:
			e.PUT(operate.Url, operate.fun)
		case AskTypeDelete:
			e.DELETE(operate.Url, operate.fun)
		case AskTypePatch:
			e.PATCH(operate.Url, operate.fun)
		case AskTypeHead:
			e.HEAD(operate.Url, operate.fun)
		case AskTypeOptions:
			e.OPTIONS(operate.Url, operate.fun)
		}
	}
}
