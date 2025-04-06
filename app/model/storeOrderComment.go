package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	CommentTypeProduct int8 = 1 // 商品评论
	CommentTypeStore   int8 = 2 // 店铺评论

	CommentsStatusPending  int8 = 10 // 待评论
	CommentsStatusComplete int8 = 20 //	已评论
)

// StoreOrderComment 店铺评论表
type StoreOrderComment struct {
	BaseModel
	AdminID   uint             `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID    uint             `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	StoreID   uint             `gorm:"type:int unsigned not null;comment:店铺ID" json:"storeId"`
	ProductID uint             `gorm:"type:int unsigned not null;comment:商品ID" json:"productId"`
	OrderID   uint             `gorm:"type:int unsigned not null;comment:订单ID" json:"orderId"`
	Content   string           `gorm:"type:varchar(512) not null;comment:评论内容" json:"content"`
	Score     float64          `gorm:"type:decimal(3,2) not null;default:0;comment:评分" json:"score"`
	Status    int8             `gorm:"type:tinyint not null;default:10;comment:状态(10:待评论,20:已评论)" json:"status"`
	Type      int8             `gorm:"type:tinyint not null;default:1;comment:评论类型(1:商品评论,2:店铺评论)" json:"type"`
	Images    Strings          `gorm:"type:text;comment:图片组" json:"images"`
	Data      OrderCommentData `gorm:"type:json;comment:配置" json:"data"`
}

type OrderCommentData struct{}

// Value implements the driver.Valuer interface
func (d OrderCommentData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *OrderCommentData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&StoreOrderComment{}); err != nil {
//		panic("Failed to auto migrate StoreOrderComment table: " + err.Error())
//	}
//}
