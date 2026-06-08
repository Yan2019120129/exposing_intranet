package file

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"my-base/app/models"
	"my-base/configs"

	adminconfig "github.com/GoAdminGroup/go-admin/modules/config"
)

type LocalService struct {
	root   string
	prefix string
}

func NewLocalService(cfg *adminconfig.Config) *LocalService {
	if cfg == nil {
		cfg = configs.GetAdmin()
	}
	return &LocalService{
		root:   cfg.Store.Path,
		prefix: strings.Trim(cfg.Store.Prefix, "/"),
	}
}

func (s *LocalService) BuildSystemFile(originalName, storagePath, category string, isPublic bool, status int) (models.SystemFile, error) {
	storagePath, err := s.normalizeStoragePath(storagePath)
	if err != nil {
		return models.SystemFile{}, err
	}

	absPath := filepath.Join(s.root, filepath.FromSlash(storagePath))
	stat, err := os.Stat(absPath)
	if err != nil {
		return models.SystemFile{}, err
	}
	if stat.IsDir() {
		return models.SystemFile{}, errors.New("uploaded path is a directory")
	}

	hash, mimeType, err := inspectFile(absPath)
	if err != nil {
		return models.SystemFile{}, err
	}

	if strings.TrimSpace(category) == "" {
		category = "default"
	}

	fileName := path.Base(storagePath)
	return models.SystemFile{
		OriginalName:  originalName,
		FileName:      fileName,
		FileExt:       strings.TrimPrefix(strings.ToLower(path.Ext(fileName)), "."),
		MimeType:      mimeType,
		FileSize:      stat.Size(),
		FileHash:      hash,
		StorageDriver: "local",
		StoragePath:   storagePath,
		PublicURL:     s.publicURL(storagePath),
		Category:      category,
		IsPublic:      isPublic,
		Status:        status,
	}, nil
}

func (s *LocalService) normalizeStoragePath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("file is required")
	}
	value = strings.TrimPrefix(value, "/")
	if s.prefix != "" && strings.HasPrefix(value, s.prefix+"/") {
		value = strings.TrimPrefix(value, s.prefix+"/")
	}
	cleaned := path.Clean(strings.ReplaceAll(value, "\\", "/"))
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return "", errors.New("invalid file path")
	}
	return cleaned, nil
}

func (s *LocalService) publicURL(storagePath string) string {
	if s.prefix == "" {
		return "/" + strings.TrimPrefix(storagePath, "/")
	}
	return "/" + s.prefix + "/" + strings.TrimPrefix(storagePath, "/")
}

func inspectFile(absPath string) (string, string, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	h := sha256.New()
	buf := make([]byte, 512)
	n, readErr := io.ReadFull(f, buf)
	if readErr != nil && !errors.Is(readErr, io.EOF) && !errors.Is(readErr, io.ErrUnexpectedEOF) {
		return "", "", readErr
	}
	mimeType := http.DetectContentType(buf[:n])

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", "", err
	}
	if _, err := io.Copy(h, f); err != nil {
		return "", "", err
	}
	return hex.EncodeToString(h.Sum(nil)), mimeType, nil
}
