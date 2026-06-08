package models

func init() {
	ModelManage.SetModel("system_files", &SystemFile{}, "系统文件")
}

type SystemFile struct {
	BaseModel
	OriginalName  string `gorm:"type:varchar(255);not null;comment:原始文件名" json:"originalName"`
	FileName      string `gorm:"type:varchar(255);not null;comment:存储文件名" json:"fileName"`
	FileExt       string `gorm:"type:varchar(32);not null;default:'';comment:扩展名" json:"fileExt"`
	MimeType      string `gorm:"type:varchar(128);not null;default:'';comment:MIME类型" json:"mimeType"`
	FileSize      int64  `gorm:"not null;default:0;comment:文件大小" json:"fileSize"`
	FileHash      string `gorm:"type:varchar(64);not null;default:'';index;comment:文件SHA256" json:"fileHash"`
	StorageDriver string `gorm:"type:varchar(32);not null;default:'local';comment:存储驱动" json:"storageDriver"`
	StoragePath   string `gorm:"type:varchar(500);not null;comment:存储路径" json:"storagePath"`
	PublicURL     string `gorm:"type:varchar(500);not null;default:'';comment:访问URL" json:"publicUrl"`
	Category      string `gorm:"type:varchar(64);not null;default:'default';index;comment:分类" json:"category"`
	IsPublic      bool   `gorm:"type:tinyint(1);not null;default:1;comment:是否公开访问" json:"isPublic"`
	Status        int    `gorm:"type:tinyint;not null;default:1;index;comment:状态 1正常 0禁用 -1删除" json:"status"`
	UploaderID    int64  `gorm:"not null;default:0;comment:上传人ID" json:"uploaderId"`
}

func (SystemFile) TableName() string {
	return "system_files"
}
