package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

const (
	UserAssetStatusDisabled int8 = -1 // 禁用
	UserAssetStatusEnabled  int8 = 10 // 启用
)

// UserAssets 用户资产
type UserAssets struct {
	BaseModel
	AdminID         uint           `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	UserID          uint           `gorm:"type:int unsigned not null;uniqueIndex:idx_user_asset;comment:用户ID" json:"userId"`
	AssetsID        uint           `gorm:"type:int unsigned not null;uniqueIndex:idx_user_asset;comment:资产ID" json:"assetsId"`
	FrozenAmount    float64        `gorm:"type:decimal(16,4) not null;default:0;comment:冻结金额" json:"frozenAmount"`
	AvailableAmount float64        `gorm:"type:decimal(16,4) not null;default:0;comment:可用金额" json:"availableAmount"`
	Status          int8           `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Data            UserAssetsData `gorm:"type:json;comment:数据" json:"data"`
}

// UserAssetsData 用户资产数据
type UserAssetsData struct {
}

// Value implements the driver.Valuer interface
func (d UserAssetsData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *UserAssetsData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&UserAssets{}); err != nil {
//		panic("Failed to auto migrate UserAssets table: " + err.Error())
//	}
//}
