package models

import "gorm.io/gorm"

// AuthAssignment 表示权限分配表
type AuthAssignment struct {
	BaseModel
	AdminId uint   `gorm:"type:int unsigned not null;index:idx_admin_name;comment:管理ID;foreignKey:ID;references:AdminUser" json:"adminId"`
	Name    string `gorm:"type:varchar(50) not null;index:idx_admin_name;comment:权限名称;foreignKey:Name;references:AuthItem" json:"name"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AuthAssignment{}); err != nil {
//		panic("Failed to auto migrate AuthAssignment table: " + err.Error())
//	}
//
//	// Initialize auth assignments
//	if err := InitAuthAssignments(db); err != nil {
//		panic("Failed to initialize auth assignments: " + err.Error())
//	}
//}

// InitAuthAssignments initializes the default auth assignments
func InitAuthAssignments(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AuthAssignment{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		assignments := []AuthAssignment{
			{
				AdminId: SuperAdminID,
				Name:    AuthItemRoleSuperAdmin,
			},
			{
				AdminId: 2,
				Name:    AuthItemRoleMerchantAdmin,
			},
			{
				AdminId: 3,
				Name:    AuthItemRoleAgentAdmin,
			},
		}

		if err := db.Create(&assignments).Error; err != nil {
			return err
		}
	}

	return nil
}
