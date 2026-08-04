package file

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"my-base/configs"

	"github.com/google/uuid"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
)

type tusUserContextKey struct{}

// TusUser 表示已通过后台认证的 tus 上传用户。
type TusUser struct {
	ID         int64
	SuperAdmin bool
}

// TusServer 提供带用户归属校验的 tus HTTP 服务。
type TusServer struct {
	endpoint string
	store    filestore.FileStore
	handler  http.Handler
}

// WithTusUser 将可信的后台用户写入 tus 请求上下文。
func WithTusUser(ctx context.Context, user TusUser) context.Context {
	return context.WithValue(ctx, tusUserContextKey{}, user)
}

// NewTusServer 创建使用本地暂存目录和文件锁的 tus 服务。
func NewTusServer(cfg *configs.TusUploadConfig) (*TusServer, error) {
	if cfg == nil {
		return nil, errors.New("tus upload config is required")
	}
	endpoint := normalizeTusEndpoint(cfg.Endpoint)
	stagingPath := strings.TrimSpace(cfg.StagingPath)
	if stagingPath == "" {
		return nil, errors.New("tus staging path is required")
	}
	if err := os.MkdirAll(stagingPath, 0o775); err != nil {
		return nil, fmt.Errorf("create tus staging directory: %w", err)
	}

	store := filestore.New(stagingPath)
	locker := filelocker.New(stagingPath)
	composer := handler.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	tusServer := &TusServer{endpoint: endpoint, store: store}
	tusHandler, err := handler.NewHandler(handler.Config{
		StoreComposer:           composer,
		BasePath:                endpoint,
		MaxSize:                 cfg.MaxSize,
		DisableDownload:         true,
		DisableConcatenation:    true,
		Cors:                    &handler.CorsConfig{Disable: true},
		PreUploadCreateCallback: tusServer.prepareUpload,
	})
	if err != nil {
		return nil, fmt.Errorf("create tus handler: %w", err)
	}
	tusServer.handler = http.StripPrefix(strings.TrimSuffix(endpoint, "/"), tusHandler)
	return tusServer, nil
}

// ServeHTTP 校验上传归属后把请求交给 tusd 协议处理器。
func (s *TusServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	user, ok := tusUserFromContext(request.Context())
	if !ok || user.ID <= 0 {
		http.Error(response, "unauthorized", http.StatusUnauthorized)
		return
	}

	uploadID, err := s.uploadID(request.URL.Path)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	if uploadID != "" && !user.SuperAdmin {
		if err := s.checkOwner(request.Context(), uploadID, user.ID); err != nil {
			status := http.StatusForbidden
			if errors.Is(err, handler.ErrNotFound) {
				status = http.StatusNotFound
			}
			http.Error(response, http.StatusText(status), status)
			return
		}
	}
	s.handler.ServeHTTP(response, request)
}

// prepareUpload 清理客户端元数据并生成不可伪造的上传标识。
func (s *TusServer) prepareUpload(event handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
	user, ok := tusUserFromContext(event.Context)
	if !ok || user.ID <= 0 {
		return handler.HTTPResponse{}, handler.FileInfoChanges{}, handler.NewError("ERR_UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
	}
	originalName := sanitizeTusOriginalName(event.Upload.MetaData["filename"])
	if originalName == "" {
		return handler.HTTPResponse{}, handler.FileInfoChanges{}, handler.NewError("ERR_INVALID_FILENAME", "filename is required", http.StatusBadRequest)
	}
	metadata := handler.MetaData{
		"filename":    originalName,
		"filetype":    sanitizeTusFileType(event.Upload.MetaData["filetype"]),
		"uploader_id": strconv.FormatInt(user.ID, 10),
	}
	changes := handler.FileInfoChanges{
		ID:       uuid.NewString() + safeTusExtension(originalName),
		MetaData: metadata,
	}
	return handler.HTTPResponse{}, changes, nil
}

// checkOwner 验证上传资源属于指定后台用户。
func (s *TusServer) checkOwner(ctx context.Context, uploadID string, userID int64) error {
	upload, err := s.store.GetUpload(ctx, uploadID)
	if err != nil {
		return err
	}
	info, err := upload.GetInfo(ctx)
	if err != nil {
		return err
	}
	ownerID, err := strconv.ParseInt(info.MetaData["uploader_id"], 10, 64)
	if err != nil || ownerID != userID {
		return errors.New("upload owner mismatch")
	}
	return nil
}

// uploadID 从 tus 访问路径中解析上传标识。
func (s *TusServer) uploadID(requestPath string) (string, error) {
	prefix := strings.TrimSuffix(s.endpoint, "/")
	if !strings.HasPrefix(requestPath, prefix) {
		return "", errors.New("invalid tus upload path")
	}
	value := strings.Trim(strings.TrimPrefix(requestPath, prefix), "/")
	if value == "" {
		return "", nil
	}
	decoded, err := url.PathUnescape(value)
	if err != nil || !isSafeTusUploadID(decoded) {
		return "", errors.New("invalid tus upload id")
	}
	return decoded, nil
}

// tusUserFromContext 获取请求中的可信后台用户。
func tusUserFromContext(ctx context.Context) (TusUser, bool) {
	user, ok := ctx.Value(tusUserContextKey{}).(TusUser)
	return user, ok
}

// normalizeTusEndpoint 规范化 tus 服务访问路径。
func normalizeTusEndpoint(endpoint string) string {
	endpoint = "/" + strings.Trim(strings.TrimSpace(endpoint), "/") + "/"
	if endpoint == "//" {
		return "/admin/files/tus/"
	}
	return endpoint
}

// sanitizeTusOriginalName 清理上传文件的展示名称。
func sanitizeTusOriginalName(value string) string {
	value = filepath.Base(strings.ReplaceAll(strings.TrimSpace(value), "\\", "/"))
	value = strings.Map(func(char rune) rune {
		if unicode.IsControl(char) || char == 0 {
			return -1
		}
		return char
	}, value)
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > 255 {
		value = string(runes[:255])
	}
	return value
}

// sanitizeTusFileType 清理客户端提供的 MIME 提示值。
func sanitizeTusFileType(value string) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > 128 {
		value = string(runes[:128])
	}
	for _, char := range value {
		if unicode.IsControl(char) {
			return ""
		}
	}
	return value
}

// safeTusExtension 返回只包含字母和数字的安全扩展名。
func safeTusExtension(filename string) string {
	extension := strings.ToLower(filepath.Ext(filename))
	if len(extension) <= 1 || len(extension) > 17 {
		return ""
	}
	for _, char := range extension[1:] {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') {
			return ""
		}
	}
	return extension
}

// isSafeTusUploadID 校验上传标识不能形成目录穿越。
func isSafeTusUploadID(value string) bool {
	if value == "" || len(value) > 80 || filepath.Base(value) != value {
		return false
	}
	base := strings.TrimSuffix(value, filepath.Ext(value))
	_, err := uuid.Parse(base)
	return err == nil
}
