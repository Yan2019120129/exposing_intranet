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
	"my-base/code/api"
	"my-base/configs"

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
