package tables

import (
	"errors"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form1 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"gorm.io/gorm"
)

func GetTestTable(ctx *context.Context) table.Table {
	test := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := test.GetInfo().HideFilterArea()
	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Name", "name", db.Varchar).
		FieldFilterable()
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
	info.SetTable("test").SetTitle("Test").SetDescription("Test").
		SetDeleteFnWithDB(func(gormDB *gorm.DB, ids []string) error {
			return deleteTestByDB(gormDB, ids)
		})

	formList := test.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.SetTable("test").SetTitle("Test").SetDescription("Test")
	formList.SetInsertFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return insertTestByDB(gormDB, values)
	})
	formList.SetUpdateFnWithDB(func(gormDB *gorm.DB, values form1.Values) error {
		return updateTestByDB(gormDB, values)
	})

	return test
}

func insertTestByDB(gormDB *gorm.DB, values form1.Values) error {
	return gormDB.Table("test").Create(map[string]interface{}{
		"name": values.Get("name"),
	}).Error
}

func updateTestByDB(gormDB *gorm.DB, values form1.Values) error {
	id := values.GetIntDefault("id", 0)
	if id <= 0 {
		return errors.New("invalid id")
	}

	return gormDB.Table("test").Where("id = ?", id).Updates(map[string]interface{}{
		"name": values.Get("name"),
	}).Error
}

func deleteTestByDB(gormDB *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return errors.New("invalid id")
	}

	return gormDB.Table("test").Where("id IN ?", ids).Delete(map[string]interface{}{}).Error
}
