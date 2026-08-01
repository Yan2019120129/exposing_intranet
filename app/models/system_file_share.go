package models

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	ModelManage.SetModel("system_file_shares", &SystemFileShare{}, "系统文件分享")
}

// SystemFileShare 表示可公开下载系统文件的分享链接。
type SystemFileShare struct {
	gorm.Model
	SystemFileID uint       `gorm:"not null;index;comment:系统文件ID" json:"systemFileId"`
	Token        string     `gorm:"type:varchar(64);not null;uniqueIndex;comment:分享令牌" json:"-"`
	ExpiresAt    time.Time  `gorm:"not null;index;comment:过期时间" json:"expiresAt"`
	RevokedAt    *time.Time `gorm:"index;comment:撤销时间" json:"revokedAt"`
	CreatedBy    int64      `gorm:"not null;default:0;comment:创建人ID" json:"createdBy"`
}

// TableName 返回系统文件分享表名。
func (SystemFileShare) TableName() string {
	return "system_file_shares"
}
