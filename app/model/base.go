package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"my-base/core/gorm/orm"
	"time"
)

var (
	db = orm.DB
)

// BaseModel 包含所有数据库表的公共字段
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"` //	主键ID
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`    //	创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`    //	更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`   //	删除对象
}

type Strings []string

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据到 Go 类型
func (s *Strings) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("unsupported Scan, storing driver.Value type into type GormStringSlice")
	}
}

// Value 实现 driver.Valuer 接口，用于将 Go 类型写入数据库
func (s Strings) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

type SliceMap []Map

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据到 Go 类型
func (s *SliceMap) Scan(v any) error {
	if v == nil {
		*s = nil
		return nil
	}
	switch value := v.(type) {
	case []byte:
		return json.Unmarshal(value, s)
	case string:
		return json.Unmarshal([]byte(value), s)
	default:
		return errors.New("unsupported Scan, storing driver.Value type into type GormStringSlice")
	}
}

// Value 实现 driver.Valuer 接口，用于将 Go 类型写入数据库
func (s SliceMap) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Add 设置值
func (s SliceMap) Add(v Map) SliceMap {
	s = append(s, v)
	return s
}

type Map map[any]any

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据到 Go 类型
func (m *Map) Scan(v any) error {
	if v == nil {
		*m = nil
		return nil
	}
	switch value := v.(type) {
	case []byte:
		return json.Unmarshal(value, m)
	case string:
		return json.Unmarshal([]byte(value), m)
	default:
		return errors.New("unsupported Scan, storing driver.Value type into type GormStringSlice")
	}
}

// Value 实现 driver.Valuer 接口，用于将 Go 类型写入数据库
func (m Map) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	return json.Marshal(m)
}

// Add 设置值
func (m Map) Add(k, v any) Map {
	m[k] = v
	return m
}

// Get 获取值
func (m Map) Get(k, v any) (any, bool) {
	v, ok := m[k]
	return v, ok
}
