package form

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	PostTypeKey           = "__go_admin_post_type"
	PostResultKey         = "__go_admin_post_result"
	PostIsSingleUpdateKey = "__go_admin_is_single_update"

	PreviousKey = "__go_admin_previous_"
	TokenKey    = "__go_admin_t_"
	MethodKey   = "__go_admin_method_"

	NoAnimationKey = "__go_admin_no_animation_"
)

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Values map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (f Values) Get(key string) string {
	if len(f[key]) > 0 {
		return f[key][0]
	}
	return ""
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (f Values) Add(key string, value string) {
	f[key] = []string{value}
}

func (f Values) Set(key string, value interface{}) {
	f[key] = []string{fmt.Sprint(value)}
}

func (f Values) SetBytes(key string, value []byte) {
	f[key] = []string{string(value)}
}

func (f Values) AddInt(key string, value int) {
	f.Add(key, strconv.Itoa(value))
}

func (f Values) AddInt64(key string, value int64) {
	f.Add(key, strconv.FormatInt(value, 10))
}

func (f Values) AddUint(key string, value uint) {
	f.Add(key, strconv.FormatUint(uint64(value), 10))
}

func (f Values) AddUint64(key string, value uint64) {
	f.Add(key, strconv.FormatUint(value, 10))
}

func (f Values) AddFloat(key string, value float64) {
	f.Add(key, strconv.FormatFloat(value, 'f', -1, 64))
}

func (f Values) AddFloat32(key string, value float32) {
	f.Add(key, strconv.FormatFloat(float64(value), 'f', -1, 32))
}

func (f Values) AddBool(key string, value bool) {
	f.Add(key, strconv.FormatBool(value))
}

func (f Values) AddBytes(key string, value []byte) {
	f.Add(key, string(value))
}

func (f Values) GetInt(key string) (int, error) {
	return strconv.Atoi(f.Get(key))
}

func (f Values) GetIntDefault(key string, def int) int {
	value, err := f.GetInt(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetInt64(key string) (int64, error) {
	return strconv.ParseInt(f.Get(key), 10, 64)
}

func (f Values) GetInt64Default(key string, def int64) int64 {
	value, err := f.GetInt64(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetUint(key string) (uint, error) {
	value, err := strconv.ParseUint(f.Get(key), 10, 0)
	return uint(value), err
}

func (f Values) GetUintDefault(key string, def uint) uint {
	value, err := f.GetUint(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetUint64(key string) (uint64, error) {
	return strconv.ParseUint(f.Get(key), 10, 64)
}

func (f Values) GetUint64Default(key string, def uint64) uint64 {
	value, err := f.GetUint64(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetFloat(key string) (float64, error) {
	return strconv.ParseFloat(f.Get(key), 64)
}

func (f Values) GetFloatDefault(key string, def float64) float64 {
	value, err := f.GetFloat(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetFloat32(key string) (float32, error) {
	value, err := strconv.ParseFloat(f.Get(key), 32)
	return float32(value), err
}

func (f Values) GetFloat32Default(key string, def float32) float32 {
	value, err := f.GetFloat32(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetBool(key string) (bool, error) {
	return strconv.ParseBool(f.Get(key))
}

func (f Values) GetBoolDefault(key string, def bool) bool {
	value, err := f.GetBool(key)
	if err != nil {
		return def
	}
	return value
}

func (f Values) GetBytes(key string) []byte {
	return []byte(f.Get(key))
}

// IsEmpty check the key is empty or not.
func (f Values) IsEmpty(key ...string) bool {
	for _, k := range key {
		if f.Get(k) == "" {
			return true
		}
	}
	return false
}

// Has check the key exists or not.
func (f Values) Has(key ...string) bool {
	for _, k := range key {
		if f.Get(k) != "" {
			return true
		}
	}
	return false
}

// Delete deletes the values associated with key.
func (f Values) Delete(key string) {
	delete(f, key)
}

// ToMap turn the values to a map[string]string type.
func (f Values) ToMap() map[string]string {
	var m = make(map[string]string)
	for key, v := range f {
		if len(v) > 0 {
			m[key] = v[0]
		}
	}
	return m
}

// IsUpdatePost check the param if is from an update post request type or not.
func (f Values) IsUpdatePost() bool {
	return f.Get(PostTypeKey) == "0"
}

// IsInsertPost check the param if is from an insert post request type or not.
func (f Values) IsInsertPost() bool {
	return f.Get(PostTypeKey) == "1"
}

// PostError get the post result.
func (f Values) PostError() error {
	msg := f.Get(PostResultKey)
	if msg == "" {
		return nil
	}
	return errors.New(msg)
}

// IsSingleUpdatePost check the param if from an single update post request type or not.
func (f Values) IsSingleUpdatePost() bool {
	return f.Get(PostIsSingleUpdateKey) == "1"
}

// RemoveRemark removes the PostType and IsSingleUpdate flag parameters.
func (f Values) RemoveRemark() Values {
	f.Delete(PostTypeKey)
	f.Delete(PostIsSingleUpdateKey)
	return f
}

// RemoveSysRemark removes all framework post flag parameters.
func (f Values) RemoveSysRemark() Values {
	f.Delete(PostTypeKey)
	f.Delete(PostIsSingleUpdateKey)
	f.Delete(PostResultKey)
	f.Delete(PreviousKey)
	f.Delete(TokenKey)
	f.Delete(MethodKey)
	f.Delete(NoAnimationKey)
	return f
}
