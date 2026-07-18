package tables

import (
	"strconv"

	"my-base/app/models"
	"my-base/app/service"
	"my-base/code/penetrate"
	"my-base/code/utils"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"gorm.io/gorm"
)

func GetClientTable(ctx *context.Context) table.Table {
	client := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := client.GetInfo().
		HideFilterArea().
		SetTable("client").
		SetTitle("Client").
		SetDescription("Client")

	info.AddField("Id", "id", db.Bigint).FieldFilterable()
	info.AddField("Name", "name", db.Text)
	info.AddField("Symbol", "symbol", db.Text)
	info.AddField("Status", "status", db.Int).
		FieldDisplay(clientStatusDisplay).
		FieldDot(clientStatusDot())
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.SetDeleteFnWithDB(func(db *gorm.DB, ids []string) error {
		return service.NewClientService(db, penetrate.GetServer()).DeleteByIDs(utils.StringIntArrayToIntArray(ids))
	})

	formList := client.GetForm().
		SetTable("client").
		SetTitle("Client").
		SetDescription("Client")
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Name", "name", db.Text, form.Text)
	formList.AddField("Symbol", "symbol", db.Text, form.Text)
	formList.AddField("Status", "status", db.Int, form.SelectSingle).
		FieldOptions(clientStatusOptions())
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)

	return client
}

func clientStatusOptions() types.FieldOptions {
	return types.FieldOptions{
		{Text: "断开", Value: strconv.Itoa(models.StatusOn)},
		{Text: "活跃", Value: strconv.Itoa(models.StatusActive)},
		{Text: "禁用", Value: strconv.Itoa(models.StatusDisable)},
		{Text: "禁用", Value: strconv.Itoa(models.LegacyStatusDisable)},
	}
}

func clientStatusDisplay(value types.FieldModel) interface{} {
	switch value.Value {
	case strconv.Itoa(models.StatusOn):
		return "断开"
	case strconv.Itoa(models.StatusActive):
		return "活跃"
	case strconv.Itoa(models.StatusDisable), strconv.Itoa(models.LegacyStatusDisable):
		return "禁用"
	default:
		return value.Value
	}
}

func clientStatusDot() (map[string]types.FieldDotColor, types.FieldDotColor) {
	return map[string]types.FieldDotColor{
		"断开":   types.FieldDotColorDanger,
		"活跃":   types.FieldDotColorSuccess,
		"无法链接": types.FieldDotColorPrimary,
		"禁用":   types.FieldDotColorInfo,
	}, types.FieldDotColorDanger
}
