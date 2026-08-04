package file

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"my-base/app/models"
	"my-base/configs"

	adminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestTusServerResumeAndPrepareSystemFile 验证断点恢复、归属校验和文件认领流程。
func TestTusServerResumeAndPrepareSystemFile(t *testing.T) {
	stagingRoot := t.TempDir()
	finalRoot := t.TempDir()
	config := &configs.TusUploadConfig{
		Enabled:     true,
		Endpoint:    "/admin/files/tus/",
		StagingPath: stagingRoot,
		Retention:   "24h",
	}
	server, err := NewTusServer(config)
	if err != nil {
		t.Fatalf("创建 tus 服务失败: %v", err)
	}

	content := []byte("resumable upload content")
	location := createTusUpload(t, server, int64(len(content)), "../测试文件.txt", 7)
	uploadID := uploadIDFromLocation(t, location)
	if !strings.HasSuffix(uploadID, ".txt") || !isSafeTusUploadID(uploadID) {
		t.Fatalf("上传标识不安全: %q", uploadID)
	}

	patchTusUpload(t, server, location, content[:8], 0, 7)
	head := tusRequest(t, server, http.MethodHead, location, nil, 7)
	if head.Code != http.StatusOK || head.Header().Get("Upload-Offset") != "8" {
		t.Fatalf("HEAD 偏移不正确: code=%d offset=%q", head.Code, head.Header().Get("Upload-Offset"))
	}
	forbidden := tusRequest(t, server, http.MethodHead, location, nil, 8)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("其他用户应被拒绝，实际状态码: %d", forbidden.Code)
	}

	service := newTestTusUploadService(stagingRoot, finalRoot)
	_, err = service.PrepareSystemFile(context.Background(), &TusSystemFileInput{
		UploadID:   uploadID,
		Category:   "test",
		IsPublic:   true,
		Status:     1,
		UploaderID: 7,
	})
	if err == nil || !strings.Contains(err.Error(), "not complete") {
		t.Fatalf("未完成上传应拒绝认领，实际错误: %v", err)
	}

	patchTusUpload(t, server, location, content[8:], 8, 7)
	prepared, err := service.PrepareSystemFile(context.Background(), &TusSystemFileInput{
		UploadID:   uploadID,
		Category:   "test",
		IsPublic:   true,
		Status:     1,
		UploaderID: 7,
	})
	if err != nil {
		t.Fatalf("认领完整上传失败: %v", err)
	}
	if prepared.File.OriginalName != "测试文件.txt" || prepared.File.FileSize != int64(len(content)) || prepared.File.UploaderID != 7 {
		t.Fatalf("系统文件元数据不正确: %+v", prepared.File)
	}
	finalContent, err := os.ReadFile(prepared.finalPath)
	if err != nil || !bytes.Equal(finalContent, content) {
		t.Fatalf("正式文件内容不正确: err=%v content=%q", err, finalContent)
	}

	if err := prepared.Rollback(); err != nil {
		t.Fatalf("回滚认领失败: %v", err)
	}
	if _, err := os.Stat(filepath.Join(stagingRoot, uploadID)); err != nil {
		t.Fatalf("回滚后暂存文件不存在: %v", err)
	}

	prepared, err = service.PrepareSystemFile(context.Background(), &TusSystemFileInput{
		UploadID: uploadID, Category: "test", IsPublic: true, Status: 1, UploaderID: 7,
	})
	if err != nil {
		t.Fatalf("回滚后重新认领失败: %v", err)
	}
	if err := prepared.Commit(); err != nil {
		t.Fatalf("提交认领失败: %v", err)
	}
	if _, err := os.Stat(filepath.Join(stagingRoot, uploadID+".info")); !os.IsNotExist(err) {
		t.Fatalf("提交后 tus 元数据仍然存在: %v", err)
	}
}

// TestTusCleanupServiceRemovesOnlyExpiredOrphans 验证清理任务保留数据库引用和活跃文件。
func TestTusCleanupServiceRemovesOnlyExpiredOrphans(t *testing.T) {
	stagingRoot := t.TempDir()
	storeRoot := t.TempDir()
	finalRoot := filepath.Join(storeRoot, "tus")
	oldTime := time.Now().Add(-48 * time.Hour)
	referencedRelative := filepath.FromSlash("tus/2026/08/01/referenced.txt")
	orphanRelative := filepath.FromSlash("tus/2026/08/01/orphan.txt")
	referencedPath := filepath.Join(storeRoot, referencedRelative)
	orphanPath := filepath.Join(storeRoot, orphanRelative)
	writeOldTestFile(t, referencedPath, oldTime)
	writeOldTestFile(t, orphanPath, oldTime)

	staleID := "stale-upload"
	writeOldTestFile(t, filepath.Join(stagingRoot, staleID), oldTime)
	writeOldTestFile(t, filepath.Join(stagingRoot, staleID+".info"), oldTime)
	recentPath := filepath.Join(stagingRoot, "recent-upload")
	if err := os.WriteFile(recentPath, []byte("recent"), 0o664); err != nil {
		t.Fatalf("写入活跃暂存文件失败: %v", err)
	}

	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := database.AutoMigrate(&models.SystemFile{}); err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}
	item := models.SystemFile{OriginalName: "referenced.txt", FileName: "referenced.txt", StoragePath: filepath.ToSlash(referencedRelative)}
	if err := database.Create(&item).Error; err != nil {
		t.Fatalf("创建引用记录失败: %v", err)
	}

	service := &TusCleanupService{
		stagingRoot: stagingRoot,
		finalRoot:   finalRoot,
		retention:   24 * time.Hour,
		now:         time.Now,
	}
	if err := service.Cleanup(context.Background(), database); err != nil {
		t.Fatalf("清理 tus 文件失败: %v", err)
	}
	if _, err := os.Stat(referencedPath); err != nil {
		t.Fatalf("数据库引用文件被错误删除: %v", err)
	}
	if _, err := os.Stat(orphanPath); !os.IsNotExist(err) {
		t.Fatalf("过期孤儿文件未删除: %v", err)
	}
	if _, err := os.Stat(filepath.Join(stagingRoot, staleID)); !os.IsNotExist(err) {
		t.Fatalf("过期暂存文件未删除: %v", err)
	}
	if _, err := os.Stat(recentPath); err != nil {
		t.Fatalf("活跃暂存文件被错误删除: %v", err)
	}
}

// newTestTusUploadService 创建使用临时目录的上传认领服务。
func newTestTusUploadService(stagingRoot, finalRoot string) *TusUploadService {
	return &TusUploadService{
		stagingRoot: stagingRoot,
		finalRoot:   finalRoot,
		store:       filestore.New(stagingRoot),
		locker:      filelocker.New(stagingRoot),
		local: NewLocalService(&adminConfig.Config{Store: adminConfig.Store{
			Path: finalRoot, Prefix: "uploads",
		}}),
		now: func() time.Time { return time.Date(2026, 8, 3, 12, 0, 0, 0, time.Local) },
	}
}

// createTusUpload 创建 tus 上传资源并返回 Location。
func createTusUpload(t *testing.T, server *TusServer, size int64, filename string, userID int64) string {
	t.Helper()
	metadata := "filename " + base64.StdEncoding.EncodeToString([]byte(filename)) + ",filetype " + base64.StdEncoding.EncodeToString([]byte("text/plain"))
	request := httptest.NewRequest(http.MethodPost, "/admin/files/tus/", nil)
	request.Header.Set("Tus-Resumable", "1.0.0")
	request.Header.Set("Upload-Length", strconv.FormatInt(size, 10))
	request.Header.Set("Upload-Metadata", metadata)
	request = request.WithContext(WithTusUser(request.Context(), TusUser{ID: userID}))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("创建 tus 上传失败: code=%d body=%q", response.Code, response.Body.String())
	}
	location := response.Header().Get("Location")
	if location == "" {
		t.Fatal("创建 tus 上传未返回 Location")
	}
	return location
}

// patchTusUpload 向 tus 上传资源追加内容。
func patchTusUpload(t *testing.T, server *TusServer, location string, content []byte, offset int64, userID int64) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPatch, location, bytes.NewReader(content))
	request.Header.Set("Tus-Resumable", "1.0.0")
	request.Header.Set("Upload-Offset", strconv.FormatInt(offset, 10))
	request.Header.Set("Content-Type", "application/offset+octet-stream")
	request = request.WithContext(WithTusUser(request.Context(), TusUser{ID: userID}))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		body, _ := io.ReadAll(response.Result().Body)
		t.Fatalf("追加 tus 内容失败: code=%d body=%q", response.Code, body)
	}
}

// tusRequest 发送带用户上下文的 tus 请求。
func tusRequest(t *testing.T, server *TusServer, method, location string, body io.Reader, userID int64) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, location, body)
	request.Header.Set("Tus-Resumable", "1.0.0")
	request = request.WithContext(WithTusUser(request.Context(), TusUser{ID: userID}))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

// uploadIDFromLocation 从 Location 中提取上传标识。
func uploadIDFromLocation(t *testing.T, location string) string {
	t.Helper()
	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("解析 Location 失败: %v", err)
	}
	return filepath.Base(strings.TrimSuffix(parsed.Path, "/"))
}

// writeOldTestFile 写入并回拨测试文件修改时间。
func writeOldTestFile(t *testing.T, filePath string, modifiedAt time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o775); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	if err := os.WriteFile(filePath, []byte("old"), 0o664); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}
	if err := os.Chtimes(filePath, modifiedAt, modifiedAt); err != nil {
		t.Fatalf("修改测试文件时间失败: %v", err)
	}
}
