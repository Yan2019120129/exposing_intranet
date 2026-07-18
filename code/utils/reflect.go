package utils

import (
	"errors"
	"log"
	"reflect"
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

	}
	return r.modelType.Name()
}

// GetFields 获取模型字段名
func (r *Reflect) GetFields() (names []string) {
	sum := r.modelType.NumField()
	for i := sum - 1; 0 < i; i-- {
		names = append(names, r.modelType.Field(i).Name)
	}
	return
}

// GetValue 获取模型字段名
func (r *Reflect) GetValue(name string) interface{} {
	return r.value[name]
}
