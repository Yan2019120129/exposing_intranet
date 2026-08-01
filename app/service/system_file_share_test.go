package service

import (
	"errors"
	"testing"
	"time"

	"my-base/app/models"
	"my-base/app/service/dto"
	codeService "my-base/code/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSystemFileShareCreateDownloadAndRevoke 验证 DTO 调用下的创建、下载和撤销行为。
func TestSystemFileShareCreateDownloadAndRevoke(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&models.SystemFile{}, &models.SystemFileShare{}); err != nil {
		t.Fatalf("迁移系统文件分享测试表失败: %v", err)
	}

	file := models.SystemFile{
		OriginalName: "example.txt",
		FileName:     "example.txt",
		StoragePath:  "default/example.txt",
		Status:       1,
		IsPublic:     true,
	}
	if err := db.Create(&file).Error; err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	fileShare := SystemFileShare{Service: codeService.Service{Orm: db}}
	share, err := fileShare.Create(&dto.SystemFileShareCreateInput{
		FileID:        file.Id,
		DurationHours: systemFileShareOneDayHours,
		CreatedBy:     7,
		IsSuperAdmin:  false,
	}, &models.SystemFile{})
	if err != nil {
		t.Fatalf("创建分享失败: %v", err)
	}
	if share.Token == "" || share.SystemFileID != file.Id || share.CreatedBy != 7 {
		t.Fatalf("分享记录不正确: %+v", share)
	}

	downloadFile := models.SystemFile{}
	if err := fileShare.GetDownloadFile(&dto.SystemFileShareDownloadInput{Token: share.Token}, &downloadFile); err != nil {
		t.Fatalf("通过分享令牌获取文件失败: %v", err)
	}
	if downloadFile.Id != file.Id {
		t.Fatalf("分享令牌返回错误文件，期望 %d，实际 %d", file.Id, downloadFile.Id)
	}

	if err := fileShare.UpdateExpiresAt(&dto.SystemFileShareUpdateExpiresAtInput{
		ShareID:      share.Id,
		ExpiresAt:    time.Now().Add(48 * time.Hour),
		IsSuperAdmin: false,
	}); err != nil {
		t.Fatalf("更新分享到期时间失败: %v", err)
	}
	if err := fileShare.Revoke(&dto.SystemFileShareRevokeInput{FileID: file.Id, ShareID: share.Id}); err != nil {
		t.Fatalf("撤销分享失败: %v", err)
	}
	if err := fileShare.GetDownloadFile(&dto.SystemFileShareDownloadInput{Token: share.Token}, &models.SystemFile{}); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("已撤销链接应不可下载，实际错误: %v", err)
	}
}

// TestSystemFileShareRejectsDeletedFile 验证已删除文件不能创建分享链接。
func TestSystemFileShareRejectsDeletedFile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&models.SystemFile{}, &models.SystemFileShare{}); err != nil {
		t.Fatalf("迁移系统文件分享测试表失败: %v", err)
	}

	deletedAt := time.Now()
	file := models.SystemFile{
		BaseModel: models.BaseModel{
			DeletedAt: &deletedAt,
		},
		OriginalName: "deleted.txt",
		FileName:     "deleted.txt",
		StoragePath:  "default/deleted.txt",
		Status:       1,
		IsPublic:     true,
	}
	if err := db.Create(&file).Error; err != nil {
		t.Fatalf("创建已删除测试文件失败: %v", err)
	}

	fileShare := SystemFileShare{Service: codeService.Service{Orm: db}}
	_, err = fileShare.Create(&dto.SystemFileShareCreateInput{
		FileID:        file.Id,
		DurationHours: systemFileShareOneDayHours,
		CreatedBy:     7,
		IsSuperAdmin:  true,
	}, &models.SystemFile{})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("已删除文件不应创建分享链接，实际错误: %v", err)
	}
}
