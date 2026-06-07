package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"
)

// IsExistFile 判断文件或目录是否存在
func IsExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
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

type File struct {
	path    string `json:"-"`        // 文件路径
	Name    string `json:"name"`     // 文件名
	Size    int64  `json:"size"`     // 文件大小
	IsDir   bool   `json:"is_dir"`   // 是否是目录
	ModTime string `json:"mod_time"` // 修改时间
	err     error  `json:"-"`        // 错误信息
}

// NewFileManager 创建文件管理器
func NewFileManager(path string) (*File, error) {
	info := &File{}
	fileInfo, err := os.Stat(path)
	info.err = err
	if err != nil {
		return nil, err
	}
	info.path = path
	info.Name = fileInfo.Name()
	info.Size = fileInfo.Size()
	info.IsDir = fileInfo.IsDir()
	info.ModTime = fileInfo.ModTime().Format(time.DateTime)
	return info, nil
}

// ToString 转换为字符串
func (f *File) ToString() string {
	data, _ := os.ReadFile(f.path)
	return string(data)
}

// ToBytes 转换为字节码
func (f *File) ToBytes() []byte {
	data, _ := os.ReadFile(f.path)
	return data
}

// JsonToStruct 转换为结构体
func (f *File) JsonToStruct(param any) error {
	data, _ := os.ReadFile(f.path)
	return json.Unmarshal(data, &param)
}

// CeateFile 创建文件
func (f *File) CeateFile() error {
	if f.IsDir {
		MkdirFile(f.path)
	}
	file, err := os.Create(f.path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// 逐行读取文件内容
func (f *File) ReadLine(fu func(val string)) error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fu(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
