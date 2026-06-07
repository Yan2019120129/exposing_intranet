package models

import (
	"time"

	"gorm.io/gorm"
)

var (
	ModelManage = &Manage{
		table: make(map[string]*modelInfo),
	}
)

type BaseModel struct {
	Id        int        `json:"id" gorm:"primarykey; comment:主键;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"comment:修改时间"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"comment:删除时间" sql:"index"`
}

type modelInfo struct {
	name    string // 表名
	model   any    // 模型
	comment string
}

type Manage struct {
	table map[string]*modelInfo // 模型
}

// SetModel 放置模型
func (m *Manage) SetModel(name string, model any, comment string) *Manage {
	m.table[name] = &modelInfo{
		name:    name,
		model:   model,
		comment: comment,
	}
	return m
}

// Create 创建表
func (m *Manage) Create(db *gorm.DB, names ...string) error {
	tables := make([]*modelInfo, 0)
	for _, name := range names {
		if p, ok := m.table[name]; ok {
			tables = append(tables, p)
		}
	}

	if len(names) == 0 {
		for _, v := range m.table {
			tables = append(tables, v)
		}
	}

	if len(tables) > 0 {
		for _, v := range tables {
			_ = db.Set("gorm:table_options", "COMMENT='"+v.comment+";'").AutoMigrate(v.model)
		}
	}
	return nil
}

// Delete 删除表
func (m *Manage) Delete(db *gorm.DB, names ...string) error {
	var tables []any
	for _, name := range names {
		if p, ok := m.table[name]; ok {
			tables = append(tables, p)
		}
	}

	if len(names) == 0 {
		for _, v := range m.table {
			tables = append(tables, v)
		}
	}

	if len(tables) > 0 {
		if err := db.Migrator().DropTable(tables...); err != nil {
			return err
		}
	}
	return nil
}

// Reset 重置表
func (m *Manage) Reset(db *gorm.DB, names ...string) error {
	err := m.Delete(db, names...)
	if err != nil {
		return err
	}
	err = m.Create(db, names...)
	if err != nil {
		return err
	}
	return nil
}
