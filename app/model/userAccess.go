package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	// Access types
	AccessTypeRefresh int8 = 1  // 刷新
	AccessTypeLogin   int8 = 2  // 登录
	AccessTypeLogout  int8 = 3  // 登出
	AccessTypeVisit   int8 = 10 // 访问

	AccessMethodPost string = "POST" // HTTP方法 - POST
)

// Access 前台访问记录表
type Access struct {
	BaseModel
	AdminID   uint       `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID    uint       `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	Name      string     `gorm:"type:varchar(120) not null;index;comment:名称" json:"name"`
	IP        string     `gorm:"type:varchar(45) not null;index;comment:IP地址" json:"ip"`
	Type      int8       `gorm:"type:tinyint not null;default:1;index;comment:类型(1:刷新,2:登录,3:登出,10:访问)" json:"type"`
	Method    string     `gorm:"type:varchar(10);comment:HTTP方法" json:"method"`
	Route     string     `gorm:"type:varchar(255) not null;index;comment:路由" json:"route"`
	UserAgent string     `gorm:"type:varchar(255);comment:用户代理" json:"userAgent"`
	Referer   string     `gorm:"type:varchar(255);comment:来源页面" json:"referer"`
	Params    string     `gorm:"type:text;comment:查询参数" json:"params"`
	Data      AccessData `gorm:"type:json;comment:数据" json:"data"`
}

// AccessData 访问记录数据
type AccessData struct {
	Headers string `json:"headers"` // 头信息
}

// Value implements the driver.Valuer interface
func (d AccessData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *AccessData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Access{}); err != nil {
//		panic("Failed to auto migrate Access table: " + err.Error())
//	}
//}
