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

	"my-base/app/models"
	"my-base/app/service"
	"my-base/app/service/dto"
	"my-base/code/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SystemFile 提供系统文件下载接口。
type SystemFile struct {
	api.Api
}

// Download 下载指定的正常系统文件。
func (e SystemFile) Download(ctx *gin.Context) {
	s := service.SystemFile{}
	if err := e.MakeContext(ctx).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		e.Error(http.StatusBadRequest, errors.New("invalid id"), "invalid id")
		return
	}

	item := models.SystemFile{}
	if err := s.GetDownloadFile(&dto.SystemFileDownloadInput{FileID: id}, &item); err != nil {
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
func (e SystemFile) serveDownload(ctx *gin.Context, item *models.SystemFile) {
	fileService := service.SystemFile{}
	absPath, err := fileService.ResolveStoragePath(item.StoragePath)
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
