package apis

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"my-base/app/models"
	"my-base/app/service"
	"my-base/app/service/dto"
	"my-base/code/api"
	"my-base/configs"

	adminModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type File struct {
	api.Api
}

func (e File) Download(ctx *gin.Context) {
	if err := e.MakeContext(ctx).MakeOrm().Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		e.Error(http.StatusBadRequest, errors.New("invalid id"), "invalid id")
		return
	}

	item := models.SystemFile{}
	if err := e.Orm.Where("id = ? AND status = ?", id, 1).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(http.StatusNotFound, err, "file not found")
			return
		}
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	e.serveDownload(ctx, &item)
}

// ShareList 返回当前后台用户可管理文件的分享记录。
func (e File) ShareList(ctx *gin.Context) {
	s := service.FileShare{}
	if err := e.MakeContext(ctx).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	fileID, ok := e.parseFileID(ctx)
	if !ok {
		return
	}
	user, ok := e.shareUser(ctx)
	if !ok {
		return
	}
	item := models.SystemFile{}
	if err := s.GetManageableFile(fileID, user.IsSuperAdmin(), &item); err != nil {
		e.handleShareFileError(err)
		return
	}

	shares := make([]models.SystemFileShare, 0)
	if err := s.List(fileID, &shares); err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}
	e.OK(e.shareItems(ctx, shares), "查询成功")
}

// CreateShare 创建并返回可公开下载的文件分享链接。
func (e File) CreateShare(ctx *gin.Context) {
	s := service.FileShare{}
	if err := e.MakeContext(ctx).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	fileID, ok := e.parseFileID(ctx)
	if !ok {
		return
	}
	user, ok := e.shareUser(ctx)
	if !ok {
		return
	}
	payload := dto.FileShareCreatePayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e.Error(http.StatusBadRequest, err, "invalid json body")
		return
	}

	item := models.SystemFile{}
	share, err := s.Create(fileID, payload.DurationHours, user.Id, user.IsSuperAdmin(), &item)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(http.StatusNotFound, err, "file not found")
			return
		}
		if errors.Is(err, service.ErrInvalidShareDuration) {
			e.Error(http.StatusBadRequest, err, err.Error())
			return
		}
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}
	e.OK(e.shareItem(ctx, *share), "创建成功")
}

// RevokeShare 撤销指定文件的分享链接。
func (e File) RevokeShare(ctx *gin.Context) {
	s := service.FileShare{}
	if err := e.MakeContext(ctx).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	fileID, ok := e.parseFileID(ctx)
	if !ok {
		return
	}
	user, ok := e.shareUser(ctx)
	if !ok {
		return
	}
	item := models.SystemFile{}
	if err := s.GetManageableFile(fileID, user.IsSuperAdmin(), &item); err != nil {
		e.handleShareFileError(err)
		return
	}
	shareID, err := strconv.Atoi(ctx.Param("shareId"))
	if err != nil || shareID <= 0 {
		e.Error(http.StatusBadRequest, errors.New("invalid share id"), "invalid share id")
		return
	}
	if err := s.Revoke(fileID, shareID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(http.StatusNotFound, err, "share not found")
			return
		}
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}
	e.OK(nil, "撤销成功")
}

// ShareDownload 根据分享令牌下载文件，无需后台登录。
func (e File) ShareDownload(ctx *gin.Context) {
	s := service.FileShare{}
	if err := e.MakeContext(ctx).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	item := models.SystemFile{}
	if err := s.GetDownloadFile(strings.TrimSpace(ctx.Param("token")), &item); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(http.StatusNotFound, err, "file not found")
			return
		}
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}
	e.serveDownload(ctx, &item)
}

// serveDownload 校验本地文件并写入下载响应。
func (e File) serveDownload(ctx *gin.Context, item *models.SystemFile) {
	absPath, err := localSystemFilePath(item.StoragePath)
	if err != nil {
		e.Error(http.StatusBadRequest, err, err.Error())
		return
	}
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			e.Error(http.StatusNotFound, err, "file not found")
			return
		}
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	filename := strings.TrimSpace(item.OriginalName)
	if filename == "" {
		filename = item.FileName
	}
	if filename == "" {
		filename = filepath.Base(absPath)
	}

	contentType := strings.TrimSpace(item.MimeType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", contentDisposition(filename))
	ctx.File(absPath)
}

// parseFileID 解析请求路径中的文件标识。
func (e File) parseFileID(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		e.Error(http.StatusBadRequest, errors.New("invalid id"), "invalid id")
		return 0, false
	}
	return id, true
}

// shareUser 获取已通过文件管理权限校验的后台用户。
func (e File) shareUser(ctx *gin.Context) (adminModels.UserModel, bool) {
	user, ok := ctx.Get("user")
	if !ok {
		e.Error(http.StatusInternalServerError, errors.New("admin user not found"), "admin user not found")
		return adminModels.UserModel{}, false
	}
	result, ok := user.(adminModels.UserModel)
	if !ok {
		e.Error(http.StatusInternalServerError, errors.New("invalid admin user"), "invalid admin user")
		return adminModels.UserModel{}, false
	}
	return result, true
}

// handleShareFileError 返回文件不可管理时的统一响应。
func (e File) handleShareFileError(err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		e.Error(http.StatusNotFound, err, "file not found")
		return
	}
	e.Error(http.StatusInternalServerError, err, err.Error())
}

// shareItems 将分享记录转换为接口返回内容。
func (e File) shareItems(ctx *gin.Context, shares []models.SystemFileShare) []dto.FileShareItem {
	items := make([]dto.FileShareItem, 0, len(shares))
	for _, share := range shares {
		items = append(items, e.shareItem(ctx, share))
	}
	return items
}

// shareItem 将单条分享记录转换为接口返回内容。
func (e File) shareItem(ctx *gin.Context, share models.SystemFileShare) dto.FileShareItem {
	return dto.FileShareItem{
		Id:          share.Id,
		DownloadURL: shareDownloadURL(ctx, share.Token),
		ExpiresAt:   share.ExpiresAt,
		RevokedAt:   share.RevokedAt,
		CreatedAt:   share.CreatedAt,
		Status:      shareStatus(share),
	}
}

// shareStatus 计算分享记录当前状态。
func shareStatus(share models.SystemFileShare) string {
	if share.RevokedAt != nil {
		return "已撤销"
	}
	if !share.ExpiresAt.After(time.Now()) {
		return "已过期"
	}
	return "有效"
}

// shareDownloadURL 生成当前请求站点下的完整分享下载链接。
func shareDownloadURL(ctx *gin.Context, token string) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	if forwarded := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Proto")); forwarded != "" {
		candidate := strings.ToLower(strings.TrimSpace(strings.Split(forwarded, ",")[0]))
		if candidate == "http" || candidate == "https" {
			scheme = candidate
		}
	}
	return (&url.URL{Scheme: scheme, Host: ctx.Request.Host, Path: "/shares/" + token + "/download"}).String()
}

func localSystemFilePath(storagePath string) (string, error) {
	storagePath = strings.TrimSpace(strings.ReplaceAll(storagePath, "\\", "/"))
	if storagePath == "" {
		return "", errors.New("empty storage path")
	}
	storagePath = strings.TrimPrefix(storagePath, "/")
	prefix := strings.Trim(configs.GetAdmin().Store.Prefix, "/")
	if prefix != "" && strings.HasPrefix(storagePath, prefix+"/") {
		storagePath = strings.TrimPrefix(storagePath, prefix+"/")
	}
	cleaned := filepath.Clean(filepath.FromSlash(storagePath))
	if cleaned == "." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) || filepath.IsAbs(cleaned) {
		return "", errors.New("invalid storage path")
	}
	return filepath.Join(configs.GetAdmin().Store.Path, cleaned), nil
}

func contentDisposition(filename string) string {
	asciiName := sanitizeASCIIFileName(filename)
	if asciiName == "" {
		asciiName = "download"
	}
	return fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", asciiName, url.PathEscape(filename))
}

func sanitizeASCIIFileName(filename string) string {
	filename = strings.TrimSpace(filepath.Base(filename))
	var b strings.Builder
	for _, r := range filename {
		if r < 32 || r == 127 || r == '"' || r == '\\' || r == '/' {
			continue
		}
		if r > 126 {
			b.WriteByte('_')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
