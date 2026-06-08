package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Convert[T any](v any) (T, error) {
	var zero T

	switch any(zero).(type) {

	// ===== int =====
	case int:
		val, err := toInt(v)
		return any(val).(T), err

	// ===== float64 =====
	case float64:
		val, err := toFloat64(v)
		return any(val).(T), err

	// ===== string =====
	case string:
		val, err := toString(v)
		return any(val).(T), err

	// ===== bool =====
	case bool:
		val, err := toBool(v)
		return any(val).(T), err

	default:
		return zero, fmt.Errorf("不支持的转换类型: %v", reflect.TypeOf(zero))
	}
}

func toInt(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("string转int失败: %v", err)
		}
		return i, nil
	default:
		return 0, errors.New("无法转换为int")
	}
}

func toFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("string转float失败: %v", err)
		}
		return f, nil
	default:
		return 0, errors.New("无法转换为float64")
	}
}

func toString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil
	case float32, float64:
		return fmt.Sprintf("%f", val), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		return "", errors.New("无法转换为string")
	}
}
func toBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("string转bool失败: %v", err)
		}
		return b, nil
	default:
		return false, errors.New("无法转换为bool")
	}
}
