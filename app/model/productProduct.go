package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	// 翻译键名
	ProductTranslateName = "productName"
	// 翻译键名
	ProductTranslateDesc = "productDesc"

	// Product Types
	ProductTypeWholesale int8 = 1 //批发商品类型
	ProductTypeStore     int8 = 2 // 店铺商品类型

	// Product Status
	ProductStatusDisabled int8 = -1 // 禁用
	ProductStatusEnabled  int8 = 10 // 启用
)

// Product 产品表
type Product struct {
	BaseModel
	AdminID     uint        `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	ParentID    uint        `gorm:"type:int unsigned;index;comment:上级ID" json:"parentId"`
	CategoryID  uint        `gorm:"type:int unsigned not null;index;comment:类目ID" json:"categoryId"`
	AssetsID    uint        `gorm:"type:int unsigned not null;index;comment:资产ID" json:"assetsId"`
	StoreID     uint        `gorm:"type:int unsigned not null;comment:店铺ID" json:"storeId"`
	Name        string      `gorm:"type:varchar(500) not null;index;comment:产品名称" json:"name"`
	Symbol      string      `gorm:"type:varchar(64) not null;index;comment:产品标识" json:"symbol"`
	Images      Strings     `gorm:"type:text;comment:产品图片" json:"images"`
	Money       float64     `gorm:"type:decimal(16,2) not null;comment:金额" json:"money"`
	Type        int8        `gorm:"type:tinyint not null;default:1;index;comment:类型(1:批发商品,2:店铺商品)" json:"type"`
	Sort        int16       `gorm:"type:smallint not null;default:99;index;comment:排序" json:"sort"`
	Recommended int8        `gorm:"type:tinyint not null;default:2;index;comment:是否推荐(1:真,2:假)" json:"recommended"`
	IsTranslate int8        `gorm:"type:tinyint not null;default:2;index;comment:是否翻译(1:是,2:否)" json:"isTranslate"`
	Sales       int         `gorm:"type:int unsigned not null;default:0;comment:销售量" json:"sales"`
	Status      int8        `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Desc        string      `gorm:"type:text;comment:产品描述" json:"desc"`
	Data        ProductData `gorm:"type:json;comment:配置" json:"data"`
}

// ProductData 产品数据
type ProductData struct {
}

// Value implements the driver.Valuer interface
func (d ProductData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *ProductData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Product{}); err != nil {
//		panic("Failed to auto migrate Product table: " + err.Error())
//	}
//
//}
