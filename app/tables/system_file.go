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
	ConfigureFileNameDisplay(info)
	info.AddField("Id", "id", db.Bigint).
		FieldWidth(50).
		FieldFilterable()
	info.AddField("预览", "storage_path", db.Varchar).
		FieldWidth(120).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return renderSystemFilePreview(fmt.Sprint(value.Row["public_url"]), value.Row["mime_type"])
		})
	info.AddField("原始文件名", "original_name", db.Varchar).
		FieldWidth(150).
		FieldDisplay(DisplayFileName)
	info.AddField("分类", "category", db.Varchar).FieldFilterable()
	info.AddField("大小", "file_size", db.Bigint).
		FieldWidth(200).
		FieldDisplay(func(value types.FieldModel) interface{} {
			size, _ := strconv.ParseInt(value.Value, 10, 64)
			return fmt.Sprintf("%s (%d B)", humanFileSize(size), size)
		})
	info.AddField("MIME", "mime_type", db.Varchar).
		FieldWidth(200)
	//info.AddField("URL", "public_url", db.Varchar).FieldWidth(120).FieldDisplay(func(value types.FieldModel) interface{} {
	//	return renderSystemFileQuickShareLink(fmt.Sprint(value.Row["id"]), value.Value)
	//})
	info.AddField("下载", "file_hash", db.Varchar).
		FieldWidth(70).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return renderSystemFileDownloadButton(fmt.Sprint(value.Row["id"]))
		})
	info.AddColumn("分享", func(value types.FieldModel) interface{} {
		return renderSystemFileShareButton(fmt.Sprint(value.Row["id"]))
	}).FieldWidth(70)
	info.AddField("状态", "status", db.Int).
		FieldWidth(65).
		FieldDisplay(func(value types.FieldModel) interface{} {
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
		FieldWidth(180).
		FieldDisplay(func(value types.FieldModel) interface{} {
			t, err := time.Parse(time.RFC3339Nano, value.Value)
			if err != nil {
				return value.Value
			}
			return t.Format(time.DateTime)
		})
	info.AddField("Updated_at", "updated_at", db.Datetime).
		FieldWidth(180).
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
	info.SetFooterHtml(systemFileShareModal())
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

// renderSystemFileQuickShareLink 渲染创建一天有效期分享链接的 URL 入口。
func renderSystemFileQuickShareLink(id, publicURL string) template.HTML {
	id = template.HTMLEscapeString(id)
	publicURL = template.HTMLEscapeString(publicURL)
	if publicURL == "" {
		publicURL = "创建 1 天分享链接"
	}
	return template.HTML(`<button type="button" class="btn btn-link btn-xs system-file-quick-share-link" data-file-id="` + id + `">` + publicURL + `</button>`)
}

// renderSystemFileShareButton 渲染打开文件分享弹窗的按钮。
func renderSystemFileShareButton(id string) template.HTML {
	id = template.HTMLEscapeString(id)
	return template.HTML(`<button type="button" class="btn btn-xs btn-success system-file-share-button" data-file-id="` + id + `">分享</button>`)
}

// systemFileShareModal 渲染分享时长选择弹窗及交互脚本。
func systemFileShareModal() template.HTML {
	return template.HTML(`
<div class="modal fade" id="system-file-share-modal" tabindex="-1" role="dialog" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title">文件分享</h4>
      </div>
      <div class="modal-body">
        <p>选择有效期后将创建新的公开下载链接。</p>
        <div class="system-file-share-duration-buttons" role="group" style="display:flex;flex-wrap:wrap;gap:10px">
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="24">分享 1 天</button>
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="168">分享 7 天</button>
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="720">分享 30 天</button>
        </div>
        <div class="system-file-share-result" style="display:none;margin-top:16px">
          <label for="system-file-share-link">分享链接</label>
          <div class="input-group">
            <input id="system-file-share-link" type="text" class="form-control system-file-share-link" readonly>
            <span class="input-group-btn">
              <button type="button" class="btn btn-default system-file-share-copy">复制链接</button>
            </span>
          </div>
        </div>
        <div class="alert alert-success system-file-share-notice" role="status" style="display:none;margin-top:12px;margin-bottom:0"></div>
      </div>
    </div>
  </div>
</div>
<script>
(function () {
  var modal = $('#system-file-share-modal');
  var noticeTimer;

  function request(path, method, body) {
    return $.ajax({
      url: path,
      type: method,
      contentType: 'application/json',
      dataType: 'json',
      data: body ? JSON.stringify(body) : undefined
    });
  }

  function errorMessage(response) {
    if (response && response.responseJSON && response.responseJSON.msg) {
      return response.responseJSON.msg;
    }
    return '请求失败，请稍后重试';
  }

  function copyLink(link) {
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(link);
    }
    var input = document.createElement('textarea');
    input.value = link;
    input.style.position = 'fixed';
    input.style.opacity = '0';
    document.body.appendChild(input);
    input.select();
    var copied = document.execCommand('copy');
    document.body.removeChild(input);
    return copied ? Promise.resolve() : Promise.reject(new Error('copy failed'));
  }

  function resetShareResult() {
    modal.find('.system-file-share-link').val('');
    modal.find('.system-file-share-result').hide();
  }

  function showShareResult(link) {
    modal.find('.system-file-share-link').val(link);
    modal.find('.system-file-share-result').show();
  }

  function showNotice(message) {
    var notice = modal.find('.system-file-share-notice');
    window.clearTimeout(noticeTimer);
    notice.text(message).stop(true, true).fadeIn(150);
    noticeTimer = window.setTimeout(function () {
      notice.fadeOut(300);
    }, 2000);
  }

  function createShare(fileID, durationHours, button) {
    if (!fileID) {
      alert('文件标识无效');
      return;
    }
    button.prop('disabled', true);
    request('/admin/files/' + fileID + '/shares', 'POST', { durationHours: durationHours }).done(function (response) {
      if (!response || response.status === 'error' || !response.data || !response.data.downloadUrl) {
        alert(response.msg || '创建分享链接失败');
        return;
      }
      showShareResult(response.data.downloadUrl);
    }).fail(function (response) {
      alert(errorMessage(response));
    }).always(function () {
      button.prop('disabled', false);
    });
  }

  $(document).off('click.systemFileShare', '.system-file-share-button').on('click.systemFileShare', '.system-file-share-button', function () {
    modal.data('file-id', $(this).data('file-id'));
    resetShareResult();
    modal.modal('show');
  });

  $(document).off('click.systemFileShareCreate', '.system-file-share-create').on('click.systemFileShareCreate', '.system-file-share-create', function () {
    var button = $(this);
    createShare(modal.data('file-id'), Number(button.data('hours')), button);
  });

  $(document).off('click.systemFileShareCopy', '.system-file-share-copy').on('click.systemFileShareCopy', '.system-file-share-copy', function () {
    var link = modal.find('.system-file-share-link').val();
    if (!link) {
      return;
    }
    copyLink(link).then(function () {
      showNotice('分享链接已复制');
    }).catch(function () {
      window.prompt('请复制分享链接', link);
    });
  });

  $(document).off('click.systemFileQuickShare', '.system-file-quick-share-link').on('click.systemFileQuickShare', '.system-file-quick-share-link', function () {
    var button = $(this);
    modal.data('file-id', button.data('file-id'));
    resetShareResult();
    modal.modal('show');
    createShare(button.data('file-id'), 24, button);
  });
}());
</script>`)
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
