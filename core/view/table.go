package view

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-base/utils"
	"net/http"
	"reflect"
	"strings"
)

// Table 创建表单视图
type Table struct {
	Fields    []*Field   `json:"fields"` // 字段信息
	model     any        // 模型
	tableName string     // 模型名
	Operates  []*Operate // 操作
}

// NewTable 创建表视图
func NewTable(model any) *Table {
	table := &Table{}
	reflectInstance := utils.NewReflectModel(model)
	table.tableName = strings.ToLower(reflectInstance.GetName())
	for _, v := range reflectInstance.GetFieldsInfo() {
		table.Fields = append(table.Fields, NewField(utils.CamelToSnake(v.Name), getFieldType(v.Type.Kind())).SetDesc(reflectInstance.GetFieldsDesc(v.Name, "gorm", "comment:")))
	}
	table.model = model

	return table
}

// Field 设置模型字段类型
func (t *Table) Field(field string, fieldType FieldType) *Field {
	fieldInstance := NewField(field, fieldType)
	for i, v := range t.Fields {
		if v.Name == field {
			t.Fields[i] = fieldInstance
			return fieldInstance
		}
	}
	t.Fields = append(t.Fields, fieldInstance)
	return fieldInstance
}

// GetOperates 获取操作方法
func (t *Table) GetOperates() []*Operate {
	// 查找
	t.Operates = append(t.Operates, NewOperate("/"+t.tableName, AskTypeGet, func(c *gin.Context) {
		v, ok := c.Get("db")
		if ok {
			fmt.Println("------", v)
		}
		c.JSON(http.StatusOK, "ok")
	}))

	// 更新
	t.Operates = append(t.Operates, NewOperate("/"+t.tableName, AskTypePut, func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	}))

	// 删除
	t.Operates = append(t.Operates, NewOperate("/"+t.tableName+"/:ids", AskTypeDelete, func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	}))

	// 创建
	t.Operates = append(t.Operates, NewOperate("/"+t.tableName, AskTypePost, func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	}))

	return t.Operates
}

// getFieldType 获取字段类型
func getFieldType(v reflect.Kind) FieldType {
	fieldType := Number
	switch v {
	case reflect.Bool:
		fieldType = Bool
	case reflect.Int:
		fieldType = Number
	case reflect.Int8:
		fieldType = Number
	case reflect.Int16:
		fieldType = Number
	case reflect.Int32:
		fieldType = Number
	case reflect.Int64:
		fieldType = Number
	case reflect.Uint:
		fieldType = Number
	case reflect.Uint8:
		fieldType = Number
	case reflect.Uint16:
		fieldType = Number
	case reflect.Uint32:
		fieldType = Number
	case reflect.Uint64:
		fieldType = Number
	case reflect.Uintptr:
		fieldType = Number
	case reflect.Float32:
		fieldType = Number
	case reflect.Float64:
		fieldType = Number
	case reflect.Complex64:
		fieldType = Number
	case reflect.Complex128:
		fieldType = Number
	case reflect.Array:
		fieldType = Slice
	case reflect.Chan:
		panic("unhandled default case")
	case reflect.Func:
		panic("unhandled default case")
	case reflect.Interface:
		fieldType = Number
	case reflect.Map:
		fieldType = Map
	case reflect.Pointer:
		panic("unhandled default case")
	case reflect.Slice:
		fieldType = Slice
	case reflect.String:
		fieldType = Text
	case reflect.Struct:
		fieldType = Struct
	case reflect.UnsafePointer:
		panic("unhandled default case")
	default:
		panic("unhandled default case")
	}
	return fieldType
}
