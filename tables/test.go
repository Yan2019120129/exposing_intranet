package tables

import (
	"errors"
	"my-base/utils"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/db/dialect"
	form1 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetTestTable(ctx *context.Context) table.Table {
	test := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := test.GetInfo().HideFilterArea()
	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Name", "name", db.Varchar).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.SetTable("test").SetTitle("Test").SetDescription("Test").
		SetDeleteFnWithDB(func(sql *db.SQL, ids []string) error {
			return deleteTestByDB(sql, ids)
		})

	formList := test.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.SetTable("test").SetTitle("Test").SetDescription("Test")
	formList.SetInsertFnWithDB(func(sql *db.SQL, values form1.Values) error {
		return insertTestByDB(sql, values)
	})
	formList.SetUpdateFnWithDB(func(sql *db.SQL, values form1.Values) error {
		return updateTestByDB(sql, values)
	})

	return test
}

func insertTestByDB(sql *db.SQL, values form1.Values) error {
	_, err := sql.Table("test").Insert(dialect.H{
		"name":       values.Get("name"),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	return err
}

func updateTestByDB(sql *db.SQL, values form1.Values) error {
	_, err := sql.Table("test").Where("id", "=", values.GetIntDefault("id", 0)).Update(dialect.H{
		"name":       values.Get("name"),
		"updated_at": time.Now(),
	})
	return err
}

func deleteTestByDB(sql *db.SQL, ids []string) error {
	if len(ids) == 0 {
		return errors.New("invalid id")
	}

	return sql.Table("test").WhereIn("id", utils.SliceToAny(ids)).Delete()
}
