package models

import "gorm.io/gorm"

// AuthChildType 定义权限项之间的关系类型
const (
	AuthChildTypeRoleRole        int8 = 1 // 角色-角色关系
	AuthChildTypePermissionRoute int8 = 2 // 权限-路由关系
	AuthChildTypeRolePermission  int8 = 3 // 角色-权限关系
)

// AuthChild 表示权限项之间的层级关系
type AuthChild struct {
	BaseModel
	Parent string `gorm:"type:varchar(50) not null;index:idx_parent_child;comment:父级权限项;foreignKey:Name;references:AuthItem" json:"parent"`
	Child  string `gorm:"type:varchar(50) not null;index:idx_parent_child;comment:子级权限项;foreignKey:Name;references:AuthItem" json:"child"`
	Type   int8   `gorm:"type:tinyint not null;comment:关系类型 (1:角色-角色,2:权限-路由,3:角色-权限)" json:"type"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AuthChild{}); err != nil {
//		panic("Failed to auto migrate AuthChild table: " + err.Error())
//	}
//
//	// Initialize auth child relationships
//	if err := InitAuthChildRelationships(db); err != nil {
//		panic("Failed to initialize auth child relationships: " + err.Error())
//	}
//}

// InitAuthChildRelationships initializes the default auth child relationships
func InitAuthChildRelationships(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AuthChild{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		relationships := []AuthChild{
			{Parent: AuthItemRoleSuperAdmin, Child: AuthItemRoleMerchantAdmin, Type: AuthChildTypeRoleRole},
			{Parent: AuthItemRoleMerchantAdmin, Child: AuthItemRoleAgentAdmin, Type: AuthChildTypeRoleRole},
			{Parent: AuthItemRoleMerchantAdmin, Child: AuthItemRoleFinanceAdmin, Type: AuthChildTypeRoleRole},
			{Parent: AuthItemRoleMerchantAdmin, Child: AuthItemRoleBackupAdmin, Type: AuthChildTypeRoleRole},
		}

		if err := db.Create(&relationships).Error; err != nil {
			return err
		}
	}

	return nil
}
