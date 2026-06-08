package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
)

// CopyMap 【浅拷贝】会复制一个 map（值类型是可直接赋值的）
func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// DeepCopyMap 深拷贝
func DeepCopyMap[K comparable, V any](src map[K]V, clone func(V) V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = clone(v) // 调用外部提供的复制方法
	}
	return dst
}

// 定义一个 Map 类型
type Map map[string]any

func (m Map) GetInt(key string) int {
	return cast.ToInt(m[key])
}

func (m Map) GetInt64(key string) int64 {
	return cast.ToInt64(m[key])
}

func (m Map) GetString(key string) string {
	return cast.ToString(m[key])
}

func (m Map) GetBool(key string) bool {
	return cast.ToBool(m[key])
}

func (m Map) GetFloat64(key string) float64 {
	return cast.ToFloat64(m[key])
}

func (m Map) GetStrings(key string) []string {
	return cast.ToStringSlice(m[key])
}
func (m *Map) Scan(value interface{}) error {
	if value == nil {
		*m = Map{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Map from value: %v", value)
	}
	return json.Unmarshal(b, m)
}

func (m Map) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
