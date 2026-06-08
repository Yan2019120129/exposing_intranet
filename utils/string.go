package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// EncryptPassword 使用双重MD5加密用户密码
func EncryptPassword(password string) string {
	if password == "" {
		return ""
	}

	// 第一次MD5加密
	firstHash := md5.Sum([]byte(password))
	firstHashStr := hex.EncodeToString(firstHash[:])

	// 第二次MD5加密
	secondHash := md5.Sum([]byte(firstHashStr))
	return hex.EncodeToString(secondHash[:])
}

// ToHump 转换驼峰方法
func ToHump(s string) string {
	sl := strings.Split(s, "_")

	ns := ""
	for i := 0; i < len(sl); i++ {
		ns += strings.ToUpper(sl[i][:1]) + sl[i][1:]
	}
	return ns
}

// ToUnderlinedWords 转换下划线单词
func ToUnderlinedWords(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}

			output = append(output, unicode.ToLower(r))
		}
	}
	return string(output)
}

// GenerateRandomString 随机生成指定长度字符串
func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CamelToSnake 将驼峰命名转换为下划线命名
func CamelToSnake(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// StringToIntList 字符串转int数组
func StringToIntList(s string) []int {
	parts := strings.FieldsFunc(s, func(v rune) bool {
		return v < 48 || v > 57
	})
	intList := make([]int, 0)
	for _, part := range parts {
		if part != "" {
			parseInt, err := strconv.ParseFloat(part, 64)
			if err == nil {
				intList = append(intList, int(parseInt))
			}
		}
	}
	return intList
}

// StringToFloat64List 字符串转int数组
func StringToFloat64List(s string) []float64 {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r != '.' && (r < 48 || r > 57)
	})

	floatList := make([]float64, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			parseFloat, err := strconv.ParseFloat(part, 64)
			if err == nil {
				floatList = append(floatList, parseFloat)
			}
		}
	}
	return floatList
}

// StringIntArrayToIntArray 字符串数组转int数组
func StringIntArrayToIntArray(vales []string) (data []int) {
	for _, vale := range vales {
		v, err := strconv.ParseInt(vale, 10, 64)
		if err != nil {
			return
		}
		data = append(data, int(v))
	}
	return
}
