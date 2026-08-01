package service

import (
	"errors"
	"testing"

	"my-base/app/models"
	"my-base/app/service/dto"
	codeService "my-base/code/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSystemFileGetDownloadFile 验证系统文件服务只返回正常文件。
func TestSystemFileGetDownloadFile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&models.SystemFile{}); err != nil {
		t.Fatalf("迁移系统文件测试表失败: %v", err)
	}

	file := models.SystemFile{
		OriginalName: "example.txt",
		FileName:     "example.txt",
		StoragePath:  "default/example.txt",
		Status:       1,
	}
	if err := db.Create(&file).Error; err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	fileService := SystemFile{Service: codeService.Service{Orm: db}}
	item := models.SystemFile{}
	if err := fileService.GetDownloadFile(&dto.SystemFileDownloadInput{FileID: file.Id}, &item); err != nil {
		t.Fatalf("查询正常文件失败: %v", err)
	}
	if item.Id != file.Id {
		t.Fatalf("查询到错误文件，期望 %d，实际 %d", file.Id, item.Id)
	}

	if err := db.Model(&file).Update("status", 0).Error; err != nil {
		t.Fatalf("禁用测试文件失败: %v", err)
	}
	if err := fileService.GetDownloadFile(&dto.SystemFileDownloadInput{FileID: file.Id}, &models.SystemFile{}); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("禁用文件不应可下载，实际错误: %v", err)
	}
}
