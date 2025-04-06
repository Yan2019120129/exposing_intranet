package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// AdminLogs 表示管理员日志条目
type AdminLogs struct {
	BaseModel
	AdminID uint       `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	IP      string     `gorm:"type:varchar(45) not null;index;comment:IP4" json:"ip"`
	Action  string     `gorm:"type:varchar(120) not null;index;comment:操作" json:"action"`
	Route   string     `gorm:"type:varchar(255) not null;index;comment:路由" json:"route"`
	Method  string     `gorm:"type:varchar(10) not null;comment:方法" json:"method"`
	Params  string     `gorm:"type:text;comment:参数" json:"params"`
	Headers HeaderInfo `gorm:"type:json;comment:头信息" json:"headers"`
	Data    string     `gorm:"type:text;comment:数据" json:"data"`
}

type HeaderInfo struct {
	UserAgent      string `json:"userAgent"`
	Referer        string `json:"referer"`
	AcceptLanguage string `json:"acceptLanguage"`
}

// Value 实现了 driver.Valuer 接口，用于将 HeaderInfo 转换为数据库可存储的值
func (h HeaderInfo) Value() (driver.Value, error) {
	return json.Marshal(h)
}

// Scan 实现了 sql.Scanner 接口，用于将数据库中的值解析为 HeaderInfo 结构
func (h *HeaderInfo) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}

	return json.Unmarshal(bytes, &h)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AdminLogs{}); err != nil {
//		panic("Failed to auto migrate AdminLogs table: " + err.Error())
//	}
//}
