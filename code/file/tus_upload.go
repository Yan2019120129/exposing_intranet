package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"my-base/app/models"
	"my-base/configs"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
)

// TusSystemFileInput 定义将已完成 tus 上传转换为系统文件的输入。
type TusSystemFileInput struct {
	UploadID   string
	Category   string
	IsPublic   bool
	Status     int
	UploaderID int64
}

// TusUploadService 提供 tus 上传校验和正式文件认领能力。
type TusUploadService struct {
	stagingRoot string
	finalRoot   string
	store       filestore.FileStore
	locker      filelocker.FileLocker
	local       *LocalService
	now         func() time.Time
}

// PreparedTusUpload 表示已经移动到正式目录、等待数据库结果的上传文件。
type PreparedTusUpload struct {
	File        models.SystemFile
	stagingPath string
	finalPath   string
	infoPath    string
	lock        handler.Lock
}

// NewTusUploadService 创建本地 tus 上传认领服务。
func NewTusUploadService(tusConfig *configs.TusUploadConfig) (*TusUploadService, error) {
	if tusConfig == nil || strings.TrimSpace(tusConfig.StagingPath) == "" {
		return nil, errors.New("tus staging path is required")
	}
	adminConfig := configs.GetAdmin()
	if adminConfig == nil || strings.TrimSpace(adminConfig.Store.Path) == "" {
		return nil, errors.New("admin file store path is required")
	}
	stagingRoot, err := filepath.Abs(tusConfig.StagingPath)
	if err != nil {
		return nil, fmt.Errorf("resolve tus staging path: %w", err)
	}
	finalRoot, err := filepath.Abs(adminConfig.Store.Path)
	if err != nil {
		return nil, fmt.Errorf("resolve final file path: %w", err)
	}
	return &TusUploadService{
		stagingRoot: stagingRoot,
		finalRoot:   finalRoot,
		store:       filestore.New(stagingRoot),
		locker:      filelocker.New(stagingRoot),
		local:       NewLocalService(adminConfig),
		now:         time.Now,
	}, nil
}

// PrepareSystemFile 校验上传归属和完整性，并把文件移动到正式目录。
func (s *TusUploadService) PrepareSystemFile(ctx context.Context, input *TusSystemFileInput) (*PreparedTusUpload, error) {
	if input == nil || input.UploaderID <= 0 || !isSafeTusUploadID(strings.TrimSpace(input.UploadID)) {
		return nil, errors.New("invalid tus upload input")
	}
	uploadID := strings.TrimSpace(input.UploadID)
	lock, err := s.locker.NewLock(uploadID)
	if err != nil {
		return nil, fmt.Errorf("create tus upload lock: %w", err)
	}
	if err := lock.Lock(ctx, func() {}); err != nil {
		return nil, fmt.Errorf("lock tus upload: %w", err)
	}

	prepared, err := s.prepareLocked(ctx, input, uploadID, lock)
	if err != nil {
		_ = lock.Unlock()
		return nil, err
	}
	return prepared, nil
}

// prepareLocked 在持有上传锁时校验并移动文件。
func (s *TusUploadService) prepareLocked(ctx context.Context, input *TusSystemFileInput, uploadID string, lock handler.Lock) (*PreparedTusUpload, error) {
	upload, err := s.store.GetUpload(ctx, uploadID)
	if err != nil {
		return nil, fmt.Errorf("get tus upload: %w", err)
	}
	info, err := upload.GetInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tus upload info: %w", err)
	}
	if info.SizeIsDeferred || info.Size < 0 || info.Offset != info.Size {
		return nil, errors.New("tus upload is not complete")
	}
	ownerID, err := strconv.ParseInt(info.MetaData["uploader_id"], 10, 64)
	if err != nil || ownerID != input.UploaderID {
		return nil, errors.New("tus upload owner mismatch")
	}
	originalName := sanitizeTusOriginalName(info.MetaData["filename"])
	if originalName == "" {
		return nil, errors.New("tus upload filename is missing")
	}

	stagingPath := info.Storage[filestore.StorageKeyPath]
	infoPath := info.Storage[filestore.StorageKeyInfoPath]
	if err := ensurePathWithinRoot(s.stagingRoot, stagingPath); err != nil {
		return nil, err
	}
	if err := ensurePathWithinRoot(s.stagingRoot, infoPath); err != nil {
		return nil, err
	}

	now := s.now()
	relativePath := path.Join("tus", now.Format("2006/01/02"), uploadID)
	finalPath := filepath.Join(s.finalRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(finalPath), 0o775); err != nil {
		return nil, fmt.Errorf("create final upload directory: %w", err)
	}
	if err := moveTusFile(stagingPath, finalPath); err != nil {
		return nil, fmt.Errorf("move tus upload to final directory: %w", err)
	}

	prepared := &PreparedTusUpload{
		stagingPath: stagingPath,
		finalPath:   finalPath,
		infoPath:    infoPath,
		lock:        lock,
	}
	fileModel, err := s.local.BuildSystemFile(originalName, relativePath, input.Category, input.IsPublic, input.Status)
	if err != nil {
		rollbackErr := prepared.rollbackFile()
		if rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}
	fileModel.UploaderID = input.UploaderID
	prepared.File = fileModel
	return prepared, nil
}

// Commit 确认数据库已保存，并清理 tus 元数据和文件锁。
func (p *PreparedTusUpload) Commit() error {
	if p == nil || p.lock == nil {
		return errors.New("prepared tus upload is invalid")
	}
	removeErr := os.Remove(p.infoPath)
	if os.IsNotExist(removeErr) {
		removeErr = nil
	}
	unlockErr := p.lock.Unlock()
	p.lock = nil
	return errors.Join(removeErr, unlockErr)
}

// Rollback 在数据库保存失败时把正式文件移回 tus 暂存目录。
func (p *PreparedTusUpload) Rollback() error {
	if p == nil || p.lock == nil {
		return nil
	}
	moveErr := p.rollbackFile()
	unlockErr := p.lock.Unlock()
	p.lock = nil
	return errors.Join(moveErr, unlockErr)
}

// rollbackFile 将已移动的文件恢复到暂存位置。
func (p *PreparedTusUpload) rollbackFile() error {
	if _, err := os.Stat(p.finalPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := moveTusFile(p.finalPath, p.stagingPath); err != nil {
		return fmt.Errorf("restore tus staging file: %w", err)
	}
	return nil
}

// ensurePathWithinRoot 校验文件路径位于指定根目录内。
func ensurePathWithinRoot(root, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("tus storage path is missing")
	}
	absolute, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf("resolve tus storage path: %w", err)
	}
	relative, err := filepath.Rel(root, absolute)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) || filepath.IsAbs(relative) {
		return errors.New("tus storage path is outside staging directory")
	}
	return nil
}

// moveTusFile 优先原子移动文件，跨文件系统时安全复制后删除源文件。
func moveTusFile(source, destination string) error {
	if err := os.Rename(source, destination); err == nil {
		return nil
	} else if !errors.Is(err, syscall.EXDEV) {
		return err
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o664)
	if err != nil {
		return err
	}
	copySucceeded := false
	defer func() {
		_ = destinationFile.Close()
		if !copySucceeded {
			_ = os.Remove(destination)
		}
	}()
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}
	if err := destinationFile.Sync(); err != nil {
		return err
	}
	if err := destinationFile.Close(); err != nil {
		return err
	}
	if err := os.Remove(source); err != nil {
		return err
	}
	copySucceeded = true
	return nil
}
