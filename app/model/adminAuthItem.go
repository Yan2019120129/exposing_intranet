package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

// AuthItemType 定义管理权限Items类型
const (
	AuthItemTypeRole       int8 = 1 // 角色
	AuthItemTypePermission int8 = 2 // 权限
	AuthItemTypeRoute      int8 = 3 // 路由

	// AdminRoleNamePrefix 管理角色名称前缀
	AdminRoleNamePrefix = "%s(%v)"
	// AdminRolesPrefix 管理角色缓存Key
	AdminRolesPrefix = "adminRoles_%v"

	AuthItemRoleSuperAdmin    = "超级管理员(1)" //	超级管理员
	AuthItemRoleMerchantAdmin = "商户管理员(1)" //	商户管理员
	AuthItemRoleAgentAdmin    = "代理管理员(1)" //	代理管理员
	AuthItemRoleFinanceAdmin  = "财务管理员(1)" //	财务管理员
	AuthItemRoleBackupAdmin   = "备用管理员(1)" //	备用管理员
)

// DefaultPermissions 默认权限
var DefaultPermissions = []string{"控制台信息", "初始化数据", "管理信息", "更新管理信息", "更新管理密码", "上传文件", "菜单配置", "更新语言翻译", "语言标签", "获取选项", "通知详情", "提示音",
	"用户列表", "新增用户", "更新用户", "用户余额", "用户归属", "用户配置", "删除用户",
	"用户认证列表", "新增用户认证", "更新用户认证", "更新用户认证状态", "删除用户认证",
	"用户等级列表", "新增用户等级", "更新用户等级", "删除用户等级",
	"提现账户列表", "新增提现账户", "更新提现账户", "删除提现账户",
	"充值订单列表", "新增充值订单", "更新钱包订单", "删除钱包订单", "更新钱包订单状态",
	"账单列表", "新增账单", "更新账单", "删除账单",
	"提现订单列表", "新增提现订单", "转账订单列表", "新增转账订单", "更新转账订单", "删除转账订单",
	"商品列表", "新增商品", "更新商品", "删除商品", "新增商品评论", "上架商品", "下架商品", "产品扩展设置",
	"默认订单列表", "新增默认订单", "更新默认订单", "删除默认订单",
	"店铺入驻列表", "新增店铺入驻", "更新店铺入驻", "删除店铺入驻", "更新店铺入驻状态",
	"店铺列表", "更新店铺", "删除店铺",
	"收货地址列表", "新增收货地址", "更新收货地址", "删除收货地址",
	"关注收藏列表", "新增关注收藏", "更新关注收藏", "删除关注收藏",
	"评论列表", "更新评论", "删除评论",
	"售后列表", "更新售后", "删除售后", "售后审核",
	"购物车列表", "更新购物车", "删除购物车",
	"店铺订单列表", "删除店铺订单", "店铺订单发货", "店铺订单收货", "店铺订单创建", "店铺扩展设置"}

// AuthItem 管理权限Items
type AuthItem struct {
	BaseModel
	Name        string       `gorm:"type:varchar(50) not null;uniqueIndex:idx_name_type;comment:名称" json:"name"`
	Type        int8         `gorm:"type:tinyint not null;uniqueIndex:idx_name_type;comment:类型(1:角色,2:权限,3:路由)" json:"type"`
	Description string       `gorm:"type:varchar(255) not null;comment:详情" json:"description"`
	Rule        string       `gorm:"type:varchar(255) not null;comment:规则" json:"rule"`
	Data        AuthItemData `gorm:"type:json;comment:数据" json:"data"`
}

// AuthItemData represents custom data for AuthItem
type AuthItemData struct {
	Permissions       []string `json:"permissions"`       //	接受权限  * 接受所有
	FilterPermissions []string `json:"filterPermissions"` //	过滤权限 如果有设置, 那么不会执行接受权限
}

// Value implements the driver.Valuer interface
func (d AuthItemData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *AuthItemData) Scan(value interface{}) error {
	if value == nil {
		// 如果数据库中的值为 NULL，初始化为空的 AuthItemData
		*d = AuthItemData{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// 如果数据为空，初始化为空的 AuthItemData
	if len(bytes) == 0 {
		*d = AuthItemData{}
		return nil
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AuthItem{}); err != nil {
//		panic("Failed to auto migrate AuthItem table: " + err.Error())
//	}
//
//	// Initialize auth items
//	if err := InitAuthItems(db); err != nil {
//		panic("Failed to initialize auth items: " + err.Error())
//	}
//}

// InitAuthItems initializes the default auth items
func InitAuthItems(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AuthItem{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		authItems := []AuthItem{
			{
				Name: AuthItemRoleSuperAdmin,
				Type: AuthItemTypeRole,
				Data: AuthItemData{Permissions: []string{"*"}},
			},
			{
				Name: AuthItemRoleMerchantAdmin,
				Type: AuthItemTypeRole,
				Data: AuthItemData{FilterPermissions: []string{""}},
			},
			{
				Name: AuthItemRoleAgentAdmin,
				Type: AuthItemTypeRole,
				Data: AuthItemData{Permissions: DefaultPermissions},
			},
			{
				Name: AuthItemRoleFinanceAdmin,
				Type: AuthItemTypeRole,
				Data: AuthItemData{Permissions: DefaultPermissions},
			},
			{
				Name: AuthItemRoleBackupAdmin,
				Type: AuthItemTypeRole,
				Data: AuthItemData{Permissions: DefaultPermissions},
			},
		}

		if err := db.Create(&authItems).Error; err != nil {
			return err
		}
	}

	return nil
}
