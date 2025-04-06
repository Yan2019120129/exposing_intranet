package utils

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

type Reflect struct {
	modelType reflect.Type
	model     interface{}
	field     []string
	value     map[string]interface{}
	client    []*Reflect
}

// NewReflectModel 新建反射工具
func NewReflectModel(model interface{}) *Reflect {
	r := &Reflect{value: make(map[string]interface{}), client: make([]*Reflect, 0), model: model}
	modelType := reflect.TypeOf(model)
	value := reflect.ValueOf(model)
	switch modelType.Kind() {
	case reflect.Struct:
		r.modelType = modelType
	case reflect.Ptr:
		r.modelType = modelType.Elem()
		value = value.Elem()
	case reflect.Slice:
		r.modelType = modelType
		for i := 0; i < value.Len(); i++ {
			element := value.Index(i).Interface()
			r.client = append(r.client, NewReflectModel(element))
		}
	default:
		log.Println(errors.New("it's not struct"))
		return nil
	}
	r.setValue(value)

	return r
}

// SetValue 设置值
func (r *Reflect) setValue(value reflect.Value) {
	if r.modelType.Kind() == reflect.Slice {
		return
	}
	sum := r.modelType.NumField()
	for i := 0; i < sum; i++ {
		structField := r.modelType.Field(i)
		r.field = append(r.field, structField.Name)
		field := value.FieldByName(structField.Name)
		if field.IsValid() && field.CanInterface() {
			r.value[structField.Name] = field.Interface()
		}
	}
}

// GetModelValue 获取模型值
func (r *Reflect) GetModelValue() interface{} {
	return r.model
}

// GetClient 获取子集
func (r *Reflect) GetClient() []*Reflect {
	return r.client
}

// GetClientLen 获取子集大小
func (r *Reflect) GetClientLen() int {
	return len(r.client)
}

// GetName 获取模型名
func (r *Reflect) GetName() string {
	if r.modelType.Kind() == reflect.Slice {
		return ""
	}
	return r.modelType.Name()
}

// GetFields 获取模型字段名
func (r *Reflect) GetFields() (names []string) {
	sum := r.modelType.NumField()
	for i := sum - 1; 0 <= i; i-- {
		names = append(names, r.modelType.Field(i).Name)
	}
	return
}

// GetFieldsInfo 获取模型字段信息
func (r *Reflect) GetFieldsInfo() (fieldsInfo []reflect.StructField) {
	sum := r.modelType.NumField()
	for i := sum - 1; 0 <= i; i-- {
		fieldsInfo = append(fieldsInfo, r.modelType.Field(i))
	}
	return
}

// GetFieldInfo 获取字段信息
func (r *Reflect) GetFieldInfo(name string) reflect.StructField {
	lowerName := strings.ToLower(name)
	for i := r.modelType.NumField() - 1; i >= 0; i-- {
		field := r.modelType.Field(i)
		// 检查字段名是否匹配
		if strings.ToLower(field.Name) != lowerName {
			continue
		}

		return field
	}
	return reflect.StructField{}
}

// GetFieldsDesc 获取模型字段注释
func (r *Reflect) GetFieldsDesc(name, tag, desc string) string {
	// 提前将 name 转换为小写，避免重复操作
	lowerName := strings.ToLower(name)

	// 遍历字段
	for i := r.modelType.NumField() - 1; i >= 0; i-- {
		field := r.modelType.Field(i)

		// 检查字段名是否匹配
		if strings.ToLower(field.Name) != lowerName {
			continue
		}

		// 获取字段的 tag 值
		tagVal := field.Tag.Get(tag)
		if tagVal == "" {
			return ""
		}

		// 提取 desc 对应的值
		return extractDescValue(tagVal, desc)
	}

	return ""
}

// extractDescValue 从 tagVal 中提取 desc 对应的值
func extractDescValue(tagVal, desc string) string {
	// 查找 desc 的起始位置
	descStart := strings.Index(tagVal, desc)
	if descStart == -1 {
		return ""
	}

	// 查找 desc 值的结束位置（分号或字符串末尾）
	descEnd := strings.Index(tagVal[descStart:], ";")
	if descEnd == -1 {
		return tagVal[descStart:]
	}

	// 返回 desc 对应的值
	return tagVal[descStart : descStart+descEnd]
}

// GetValue 获取模型字段值
func (r *Reflect) GetValue(name string) interface{} {
	return r.value[name]
}
