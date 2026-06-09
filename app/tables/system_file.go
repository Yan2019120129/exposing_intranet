package tables

import (
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"my-base/app/models"
	fileService "my-base/module/file"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form1 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"gorm.io/gorm"
)

func GetSystemFileTable(ctx *context.Context) table.Table {
	config := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := config.GetInfo().HideFilterArea()
	info.AddField("Id", "id", db.Bigint).FieldFilterable()
	info.AddField("预览", "storage_path", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return renderSystemFilePreview(fmt.Sprint(value.Row["public_url"]), value.Row["mime_type"])
	})
	info.AddField("原始文件名", "original_name", db.Varchar).FieldCopyable()
	info.AddField("分类", "category", db.Varchar).FieldFilterable()
	info.AddField("大小", "file_size", db.Bigint).FieldDisplay(func(value types.FieldModel) interface{} {
		size, _ := strconv.ParseInt(value.Value, 10, 64)
		return fmt.Sprintf("%s (%d B)", humanFileSize(size), size)
	})
	info.AddField("MIME", "mime_type", db.Varchar)
	info.AddField("URL", "public_url", db.Varchar).FieldCopyable()
	info.AddField("下载", "file_hash", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return renderSystemFileDownloadButton(fmt.Sprint(value.Row["id"]))
	})
	info.AddField("状态", "status", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
		switch value.Value {
		case "1":
			return template.HTML(`<span class="label label-success">正常</span>`)
		case "0":
			return template.HTML(`<span class="label label-warning">禁用</span>`)
		default:
			return template.HTML(`<span class="label label-default">删除</span>`)
		}
	})
	info.AddField("Created_at", "created_at", db.Datetime).
		FieldDisplay(func(value types.FieldModel) interface{} {
			t, err := time.Parse(time.RFC3339Nano, value.Value)
			if err != nil {
				return value.Value
			}
			return t.Format(time.DateTime)
		})
	info.AddField("Updated_at", "updated_at", db.Datetime).
		FieldDisplay(func(value types.FieldModel) interface{} {
			t, err := time.Parse(time.RFC3339Nano, value.Value)
			if err != nil {
				return value.Value
			}
			return t.Format(time.DateTime)
		})
	info.SetTable("system_files").SetTitle("系统文件").SetDescription("系统文件").
		SetDeleteFnWithDB(func(gormDB *gorm.DB, ids []string) error {
			return deleteSystemFilesByDB(gormDB, ids)
		})
	info.Where("status", "!=", -1)
	user := auth.Auth(ctx)
	if !user.IsSuperAdmin() {
		info.Where("is_public", "=", 1)
	}
	formList := config.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("上传文件", "storage_path", db.Varchar, form.File).
		FieldOptionExt(map[string]interface{}{
			"allowedFileTypes": []string{},
			//"allowedFileExtensions": []string{"jpg", "jpeg", "png", "webp", "gif", "pdf", "csv", "json", "txt", "zip", "xlsx", "xls", "docx", "doc"},
		}).
		FieldHelpMsg(template.HTML("新增时必传；编辑时不重新选择文件则保留原文件"))
	formList.AddField("原始文件名", "original_name", db.Varchar, form.Text)
	formList.AddField("文件名", "file_name", db.Varchar, form.Text).FieldHide()
	formList.AddField("扩展名", "file_ext", db.Varchar, form.Text).FieldHide()
	formList.AddField("MIME", "mime_type", db.Varchar, form.Text).FieldHide()
	formList.AddField("大小", "file_size", db.Bigint, form.Number).FieldHide()
	formList.AddField("Hash", "file_hash", db.Varchar, form.Text).FieldHide()
	formList.AddField("存储驱动", "storage_driver", db.Varchar, form.Text).FieldHide()
	formList.AddField("访问URL", "public_url", db.Varchar, form.Text).FieldHide()
	formList.AddField("分类", "category", db.Varchar, form.Text).FieldDefault("default")
	formList.AddField("公开访问", "is_public", db.Tinyint, form.Switch).
		FieldOptions(types.FieldOptions{{Text: "否", Value: "0"}, {Text: "是", Value: "1"}}).
		FieldDefault("1")
	formList.AddField("状态", "status", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{{Text: "正常", Value: "1"}, {Text: "禁用", Value: "0"}}).
		FieldDefault("1")
	formList.AddField("Uploader", "uploader_id", db.Bigint, form.Number).FieldHide()
	formList.SetTable("system_files").SetTitle("系统文件").SetDescription("系统文件")
	formList.SetInsertFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return insertSystemFileByDB(gormDB, values)
	})
	formList.SetUpdateFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return updateSystemFileByDB(gormDB, values)
	})

	return config
}

func insertSystemFileByDB(gormDB *gorm.DB, values form1.Values) error {
	fileModel, err := systemFileFromUploadValues(values, true)
	if err != nil {
		return err
	}
	return gormDB.Create(&fileModel).Error
}

func updateSystemFileByDB(gormDB *gorm.DB, values form1.Values) error {
	id := values.GetIntDefault("id", 0)
	if id <= 0 {
		return errors.New("invalid id")
	}

	if strings.TrimSpace(values.Get("storage_path")) != "" {
		fileModel, err := systemFileFromUploadValues(values, false)
		if err != nil {
			return err
		}
		return gormDB.Model(&models.SystemFile{}).Where("id = ?", id).Updates(&fileModel).Error
	}

	updates := map[string]interface{}{
		"original_name": values.Get("original_name"),
		"category":      systemFileCategory(values),
		"is_public":     values.Get("is_public") != "0",
		"status":        systemFileStatus(values),
		"updated_at":    time.Now(),
	}
	return gormDB.Model(&models.SystemFile{}).Where("id = ?", id).Updates(updates).Error
}

func deleteSystemFilesByDB(gormDB *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return errors.New("invalid id")
	}
	return gormDB.Model(&models.SystemFile{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{"status": -1, "deleted_at": time.Now()}).Error
}

func systemFileFromUploadValues(values form1.Values, isCreate bool) (models.SystemFile, error) {
	storagePath := strings.TrimSpace(values.Get("storage_path"))
	if storagePath == "" && isCreate {
		return models.SystemFile{}, errors.New("file is required")
	}

	service := fileService.NewLocalService(nil)
	fileModel, err := service.BuildSystemFile(values.Get("storage_path_original_name"), storagePath, systemFileCategory(values), values.Get("is_public") != "0", systemFileStatus(values))
	if err != nil {
		return models.SystemFile{}, err
	}
	fileModel.UploaderID = int64(values.GetIntDefault("uploader_id", 0))
	return fileModel, nil
}

func systemFileCategory(values form1.Values) string {
	category := strings.TrimSpace(values.Get("category"))
	if category == "" {
		return "default"
	}
	return category
}

func systemFileStatus(values form1.Values) int {
	status := values.GetIntDefault("status", 1)
	if status != 0 {
		return 1
	}
	return 0
}

func renderSystemFilePreview(publicURL string, mimeType interface{}) template.HTML {
	url := template.HTMLEscapeString(publicURL)
	mime := fmt.Sprint(mimeType)
	if strings.HasPrefix(mime, "image/") {
		return template.HTML(`<a href="` + url + `" target="_blank"><img src="` + url + `" style="max-width:80px;max-height:80px;object-fit:contain;border:1px solid #eee" /></a>`)
	}
	return template.HTML(`<a href="` + url + `" target="_blank">打开</a>`)
}

func renderSystemFileDownloadButton(id string) template.HTML {
	id = template.HTMLEscapeString(id)
	return template.HTML(`<a class="btn btn-xs btn-primary" href="/admin/files/download/` + id + `" target="_blank">下载</a>`)
}

func humanFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size)
	idx := -1
	for value >= 1024 && idx < len(units)-1 {
		value /= 1024
		idx++
	}
	if idx < 0 {
		return fmt.Sprintf("%d B", size)
	}
	return fmt.Sprintf("%.2f %s", value, units[idx])
}
