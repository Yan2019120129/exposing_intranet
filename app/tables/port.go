package tables

import (
	"errors"
	"strconv"

	"my-base/app/service"
	"my-base/code/penetrate"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"gorm.io/gorm"
)

func GetPortTable(ctx *context.Context) table.Table {
	port := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := port.GetInfo().HideFilterArea()
	info.AddField("Id", "id", db.Bigint).FieldFilterable()
	info.AddField("Client_id", "client_id", db.Bigint)
	info.AddField("Client_name", "name", db.Varchar).
		FieldJoin(types.Join{Table: "client", Field: "client_id", JoinField: "id"})
	info.AddField("Server", "server", db.Text)
	info.AddField("Local", "local", db.Text)
	info.AddField("Comment", "comment", db.Text)
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.SetTable("port").SetTitle("Port").SetDescription("Port")
	info.SetDeleteFnWithDB(func(db *gorm.DB, ids []string) error {
		return service.NewPortService(db, penetrate.GetServer()).DeleteByIDs(ids)
	})

	formList := port.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Client_id", "client_id", db.Bigint, form.SelectSingle).
		FieldOptions(clientNameOptions()).FieldMust()
	formList.AddField("Server", "server", db.Text, form.Text).FieldMust().FieldTrimSpace()
	formList.AddField("Local", "local", db.Text, form.Text).FieldMust().FieldTrimSpace()
	formList.AddField("Comment", "comment", db.Text, form.Text).FieldTrimSpace()
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.SetInsertFnWithDB(func(db *gorm.DB, values form2.Values) error {
		clientID, err := strconv.Atoi(values.Get("client_id"))
		if err != nil {
			return err
		}
		err = service.NewPortService(db, penetrate.GetServer()).CreateFromAdmin(
			clientID,
			values.Get("server"),
			values.Get("local"),
			values.Get("comment"),
		)
		if errors.Is(err, service.ErrPortConflict) {
			return errors.New("server port or local port repeat ！！！")
		}
		if errors.Is(err, service.ErrClientNotFound) {
			return errors.New("client not exist ！！！")
		}
		return err
	})

	return port
}

func clientNameOptions() types.FieldOptions {
	items, err := service.NewClientService(nil).NameOptions()
	if err != nil {
		return nil
	}
	options := make(types.FieldOptions, 0, len(items))
	for _, item := range items {
		options = append(options, types.FieldOption{Text: item.Name, Value: strconv.Itoa(item.ID)})
	}
	return options
}
