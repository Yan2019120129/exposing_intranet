package view

type FieldType string
type Permissions int8

// 只读 只写
const (
	PermissionsNull Permissions = iota // 无权限
	PermissionsX                       // 执行(创建)
	PermissionsW                       // 只写(修改)
	PermissionsXW                      // 执行(创建)-只写(修改)
	PermissionsR                       // 只读
	PermissionsRX                      // 读-执行(创建)
	PermissionsRW                      // 读-只写(修改)
	PermissionsRWX                     // 读-只写(修改)-执行(创建)

	Number FieldType = "number"
	Text   FieldType = "text"
	Bool   FieldType = "bool"
	Slice  FieldType = "slice"
	Map    FieldType = "Map"
	Struct FieldType = "struct"
)

type Field struct {
	Permissions Permissions `json:"permissions"`
	Type        FieldType   `json:"type"` // 字段类型
	Name        string      `json:"name"` // 字段名
	Desc        string      `json:"desc"` // 字段描述
}

// NewField 新建字段
func NewField(name string, t FieldType) *Field {
	fieldInstance := &Field{Type: t, Name: name}
	fieldInstance.Permissions = PermissionsR
	return fieldInstance
}

// SetPermissions 设置权限
func (f *Field) SetPermissions(p Permissions) *Field {
	f.Permissions = p
	return f
}

// PermissionsNull 没有权限
func (f *Field) PermissionsNull() *Field {
	f.Permissions = PermissionsNull
	return f
}

// PermissionsX 执行(创建)
func (f *Field) PermissionsX() *Field {
	f.Permissions = PermissionsX
	return f
}

// PermissionsW 只读
func (f *Field) PermissionsW() *Field {
	f.Permissions = PermissionsW
	return f
}

// PermissionsXW 只写(修改)
func (f *Field) PermissionsXW() *Field {
	f.Permissions = PermissionsXW
	return f
}

// PermissionsR 执行(创建)-只写(修改)
func (f *Field) PermissionsR() *Field {
	f.Permissions = PermissionsR
	return f
}

// PermissionsRX 读-执行(创建)
func (f *Field) PermissionsRX() *Field {
	f.Permissions = PermissionsRX
	return f
}

// PermissionsRW 读-只写(修改)
func (f *Field) PermissionsRW() *Field {
	f.Permissions = PermissionsRW
	return f
}

// PermissionsRWX 读-只写(修改)-执行(创建)
func (f *Field) PermissionsRWX() *Field {
	f.Permissions = PermissionsRWX
	return f
}

// SetDesc 设置字段描述
func (f *Field) SetDesc(desc string) *Field {
	f.Desc = desc
	return f
}
