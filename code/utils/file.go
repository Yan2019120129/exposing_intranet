package utils

import (
	"os"
	"strings"
)

// IsExistFile 判断文件路径是否存在
func IsExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// MkdirAll 创建全部路径
func MkdirAll(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// MkdirFile 创建目录和文件
func MkdirFile(path string) {
	index := strings.LastIndex(path, "/")
	dir := path[:index]
	if !IsExistFile(path) {
		MkdirAll(dir)
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	file.Close()
}
