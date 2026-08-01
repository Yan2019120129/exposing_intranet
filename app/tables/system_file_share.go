package tables

import (
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"my-base/app/models"
	"my-base/app/service"
	"my-base/app/service/dto"
	codeService "my-base/code/service"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form1 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"gorm.io/gorm"
)

const (
	shareDurationOneDay     = 24
	shareDurationSevenDays  = 24 * 7
	shareDurationThirtyDays = 24 * 30
)

// GetSystemFileShareTable 返回系统文件分享记录的后台管理表。
func GetSystemFileShareTable(ctx *context.Context) table.Table {
	config := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))
	user := auth.Auth(ctx)

	info := config.GetInfo().HideFilterArea()
	info.AddField("ID", "id", db.Bigint).FieldFilterable()
	//info.AddField("文件 ID", "system_file_id", db.Bigint).FieldFilterable()
	info.AddField("分享文件", "original_name", db.Varchar).
		FieldJoin(types.Join{Field: "system_file_id", JoinField: "id", Table: "system_files"}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return fmt.Sprint(value.Row["system_files_goadmin_join_original_name"])
		})
	info.AddField("分享链接", "token", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return renderSystemFileShareDownloadLink(value.Value)
	})
	info.AddField("到期时间", "expires_at", db.Datetime).FieldFilterable().
		FieldDisplay(formatSystemFileShareTime)
	info.AddField("创建人", "created_by", db.Bigint).FieldFilterable()
	info.AddField("撤销时间", "revoked_at", db.Datetime).
		FieldDisplay(formatSystemFileShareTime)
	info.AddColumn("状态", func(value types.FieldModel) interface{} {
		return renderSystemFileShareStatus(value.Row)
	})
	info.AddField("创建时间", "created_at", db.Datetime).FieldFilterable().
		FieldDisplay(formatSystemFileShareTime)
	info.AddField("更新时间", "updated_at", db.Datetime).
		FieldDisplay(formatSystemFileShareTime)
	info.SetTable("system_file_shares").SetTitle("系统文件分享").SetDescription("系统文件分享").
		SetDeleteFnWithDB(func(gormDB *gorm.DB, ids []string) error {
			return revokeSystemFileSharesByDB(gormDB, ids, user.IsSuperAdmin())
		})
	info.WhereRaw("system_file_id IN (SELECT id FROM system_files WHERE status = 1" + systemFileShareVisibilityCondition(user.IsSuperAdmin()) + ")")

	formList := config.GetForm()
	formList.AddField("ID", "id", db.Bigint, form.Default).FieldDisableWhenCreate()
	formList.AddField("系统文件", "system_file_id", db.Bigint, form.SelectSingle).
		FieldOptionsFromTable("system_files", "original_name", "id").
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("分享时长", "duration_hours", db.Int, form.SelectSingle).
		FieldOptions(systemFileShareDurationOptions()).
		FieldDisableWhenUpdate()
	formList.AddField("到期时间", "expires_at", db.Datetime, form.Datetime).
		FieldDisableWhenCreate()
	formList.SetTable("system_file_shares").SetTitle("系统文件分享").SetDescription("系统文件分享")
	formList.SetInsertFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return insertSystemFileShareByDB(gormDB, values, user.Id, user.IsSuperAdmin())
	})
	formList.SetUpdateFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return updateSystemFileShareByDB(gormDB, values, user.IsSuperAdmin())
	})

	return config
}

// insertSystemFileShareByDB 按选定时长创建新的分享记录。
func insertSystemFileShareByDB(gormDB *gorm.DB, values form1.Values, createdBy int64, isSuperAdmin bool) error {
	fileID := values.GetIntDefault("system_file_id", 0)
	durationHours := values.GetIntDefault("duration_hours", 0)
	fileShare := newSystemFileShareService(gormDB)
	_, err := fileShare.Create(&dto.SystemFileShareCreateInput{
		FileID:        fileID,
		DurationHours: durationHours,
		CreatedBy:     createdBy,
		IsSuperAdmin:  isSuperAdmin,
	}, &models.SystemFile{})
	return err
}

// updateSystemFileShareByDB 修改未撤销分享记录的到期时间。
func updateSystemFileShareByDB(gormDB *gorm.DB, values form1.Values, isSuperAdmin bool) error {
	shareID := values.GetIntDefault("id", 0)
	expiresAt, err := parseSystemFileShareTime(values.Get("expires_at"))
	if err != nil {
		return err
	}
	fileShare := newSystemFileShareService(gormDB)
	return fileShare.UpdateExpiresAt(&dto.SystemFileShareUpdateExpiresAtInput{
		ShareID:      shareID,
		ExpiresAt:    expiresAt,
		IsSuperAdmin: isSuperAdmin,
	})
}

// revokeSystemFileSharesByDB 批量逻辑撤销分享记录，任一失败时回滚全部操作。
func revokeSystemFileSharesByDB(gormDB *gorm.DB, ids []string, isSuperAdmin bool) error {
	if len(ids) == 0 {
		return errors.New("分享记录标识不能为空")
	}
	return gormDB.Transaction(func(tx *gorm.DB) error {
		fileShare := newSystemFileShareService(tx)
		for _, id := range ids {
			shareID, err := strconv.Atoi(id)
			if err != nil || shareID <= 0 {
				return errors.New("分享记录标识无效")
			}
			if err := fileShare.RevokeByID(&dto.SystemFileShareRevokeByIDInput{
				ShareID:      shareID,
				IsSuperAdmin: isSuperAdmin,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// newSystemFileShareService 使用当前数据库连接创建分享服务对象。
func newSystemFileShareService(gormDB *gorm.DB) service.SystemFileShare {
	return service.SystemFileShare{Service: codeService.Service{Orm: gormDB}}
}

// systemFileShareDurationOptions 返回后台表单可选的固定分享时长。
func systemFileShareDurationOptions() types.FieldOptions {
	return types.FieldOptions{
		{Text: "1 天", Value: strconv.Itoa(shareDurationOneDay)},
		{Text: "7 天", Value: strconv.Itoa(shareDurationSevenDays)},
		{Text: "30 天", Value: strconv.Itoa(shareDurationThirtyDays)},
	}
}

// systemFileShareVisibilityCondition 返回普通管理员的文件可见范围条件。
func systemFileShareVisibilityCondition(isSuperAdmin bool) string {
	if isSuperAdmin {
		return ""
	}
	return " AND is_public = 1"
}

// formatSystemFileShareTime 格式化后台表格中的时间值。
func formatSystemFileShareTime(value types.FieldModel) interface{} {
	if strings.TrimSpace(value.Value) == "" {
		return "-"
	}
	t, err := time.Parse(time.RFC3339Nano, value.Value)
	if err != nil {
		return value.Value
	}
	return t.Format(time.DateTime)
}

// parseSystemFileShareTime 解析后台表单提交的到期时间。
func parseSystemFileShareTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, errors.New("到期时间不能为空")
	}
	for _, layout := range []string{time.DateTime, time.RFC3339, time.RFC3339Nano} {
		if result, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return result, nil
		}
	}
	return time.Time{}, errors.New("到期时间格式不正确")
}

// renderSystemFileShareDownloadLink 渲染分享链接的下载入口。
func renderSystemFileShareDownloadLink(token string) template.HTML {
	token = template.HTMLEscapeString(token)
	return template.HTML(`<a class="btn btn-xs btn-primary" href="/shares/` + token + `/download" target="_blank">下载</a>`)
}

// renderSystemFileShareStatus 根据撤销和到期状态渲染状态标签。
func renderSystemFileShareStatus(row map[string]interface{}) template.HTML {
	if row["revoked_at"] != nil && fmt.Sprint(row["revoked_at"]) != "" {
		return template.HTML(`<span class="label label-default">已撤销</span>`)
	}
	expiresAt, err := parseSystemFileShareTime(fmt.Sprint(row["expires_at"]))
	if err != nil || !expiresAt.After(time.Now()) {
		return template.HTML(`<span class="label label-warning">已过期</span>`)
	}
	return template.HTML(`<span class="label label-success">有效</span>`)
}
