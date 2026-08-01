package dto

import "time"

// FileShareCreatePayload 定义创建文件分享链接的请求体。
type FileShareCreatePayload struct {
	DurationHours int `json:"durationHours"`
}

// FileShareItem 定义文件分享记录的接口返回内容。
type FileShareItem struct {
	Id          int        `json:"id"`
	DownloadURL string     `json:"downloadUrl"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	RevokedAt   *time.Time `json:"revokedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	Status      string     `json:"status"`
}
