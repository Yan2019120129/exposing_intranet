package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// AdminMenuStatus 管理员菜单状态常量
const (
	AdminMenuStatusDisabled int8 = -1 // 禁用
	AdminMenuStatusEnabled  int8 = 10 // 启用
)

// AdminMenu 管理员菜单
type AdminMenu struct {
	BaseModel
	ParentID uint          `gorm:"type:int unsigned not null;comment:父级ID" json:"parentId"`
	Name     string        `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
	Route    string        `gorm:"type:varchar(100);default:null;unique;comment:路径" json:"route"`
	Sort     uint8         `gorm:"type:tinyint unsigned not null;default:99;comment:排序" json:"sort"`
	Status   int8          `gorm:"type:tinyint not null;default:10;comment:状态(10:启用,-1:禁用)" json:"status"`
	Level    int8          `gorm:"type:tinyint not null;comment:层级" json:"level"`
	Data     AdminMenuData `gorm:"type:json;comment:数据" json:"data,omitempty"`
}

// ToDisplayData 将 AdminMenu 转换为 AdminMenuDisplayData
func (m *AdminMenu) ToDisplayData() *AdminMenuDisplayData {
	return &AdminMenuDisplayData{
		ID:       m.ID,
		Name:     m.Name,
		Route:    m.Route,
		Sort:     m.Sort,
		Status:   m.Status,
		Data:     m.Data,
		Children: make([]*AdminMenuDisplayData, 0),
	}
}

// AdminMenuDisplayData 管理菜单显示数据
type AdminMenuDisplayData struct {
	ID       uint                    `json:"id"`
	Name     string                  `json:"name"`
	Route    string                  `json:"route"`
	Sort     uint8                   `json:"sort"`
	Status   int8                    `json:"status"`
	Level    int8                    `json:"level"`
	Data     AdminMenuData           `json:"data"`
	Children []*AdminMenuDisplayData `json:"children"`
}

// AdminMenuData 管理菜单配置
type AdminMenuData struct {
	Icon    string `json:"icon" views:"label:内置图标"`    // 图标
	VueFile string `json:"vueFile" views:"label:视图文件"` // vue文件
}

// Value implements the driver.Valuer interface
func (d AdminMenuData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *AdminMenuData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AdminMenu{}); err != nil {
//		panic("Failed to auto migrate AdminMenu table: " + err.Error())
//	}
//
//	// Initialize admin menus
//	if err := InitAdminMenus(db); err != nil {
//		panic("Failed to initialize admin menus: " + err.Error())
//	}
//}

// InitAdminMenus 初始化管理菜单数据
func InitAdminMenus(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AdminMenu{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		menuData := []*AdminMenuDisplayData{
			{
				Name:   "用户管理",
				Sort:   1,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "group",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "用户列表",
						Route:  "/users/user/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "people",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "用户资产",
						Route:  "/users/assets/index",
						Sort:   2,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "account_balance_wallet",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "认证管理",
						Route:  "/users/auth/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "verified_user",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "会员列表",
						Route:  "/users/level/index",
						Sort:   4,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "star",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "访问记录",
						Route:  "/users/access/index",
						Sort:   5,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "history",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "提现账户",
						Route:  "/users/account/index",
						Sort:   6,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "account_balance",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "商品管理",
				Sort:   2,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "inventory_2",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "商品列表",
						Route:  "/products/product/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "list",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "商品规格",
						Route:  "/products/sku/index",
						Sort:   2,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "storage",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "商品分类",
						Route:  "/products/category/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "category",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "店铺管理",
				Route:  "",
				Sort:   4,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "store",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "店铺入驻",
						Route:  "/stores/settled/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "add_business",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "店铺列表",
						Route:  "/stores/store/index",
						Sort:   2,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "storefront",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "购物车",
						Route:  "/stores/cart/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "shopping_cart",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "店铺订单",
						Route:  "/stores/order/index",
						Sort:   4,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "local_shipping",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "商品订单",
						Route:  "/products/order/index",
						Sort:   5,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "currency_bitcoin",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "规格订单",
						Route:  "/products/order/default/index",
						Sort:   6,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "receipt_long",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "收货地址",
						Route:  "/stores/shippingAddress/index",
						Sort:   7,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "person_pin_circle",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "关注收藏",
						Route:  "/stores/follow/index",
						Sort:   8,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "follow_the_signs",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "商品评论",
						Route:  "/stores/comment/index",
						Sort:   9,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "speaker_notes",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "售后管理",
						Route:  "/stores/refund/index",
						Sort:   10,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "shopping_bag",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "资产管理",
				Route:  "",
				Sort:   4,
				Status: AdminMenuStatusDisabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "account_balance_wallet",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "充值订单",
						Route:  "/wallets/assets/deposit/index",
						Sort:   1,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "add_card",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "提现订单",
						Route:  "/wallets/assets/withdraw/index",
						Sort:   2,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "money_off",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "转账记录",
						Route:  "/wallets/assets/transfer/index",
						Sort:   3,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "swap_horiz",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "闪兑记录",
						Route:  "/wallets/exchange/index",
						Sort:   4,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "currency_exchange",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "账单资产",
						Route:  "/wallets/assets/bill/index",
						Sort:   5,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "receipt",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "资产列表",
						Route:  "/wallets/assets/index",
						Sort:   5,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "list",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "支付管理",
						Route:  "/wallets/assets/payment/index",
						Sort:   6,
						Status: AdminMenuStatusDisabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "payment",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "财务管理",
				Sort:   5,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "payments",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "充值订单",
						Route:  "/wallets/deposit/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "add_card",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "提现订单",
						Route:  "/wallets/withdraw/index",
						Sort:   2,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "money_off",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "转账记录",
						Route:  "/wallets/transfer/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "swap_horiz",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "账单记录",
						Route:  "/wallets/bill/index",
						Sort:   4,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "receipt",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "支付管理",
						Route:  "/wallets/payment/index",
						Sort:   5,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "payment",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "后台管理",
				Sort:   90,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "admin_panel_settings",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "管理列表",
						Route:  "/admins/manage/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "manage_accounts",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "参数配置",
						Route:  "/admins/settings/index",
						Sort:   2,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "tune",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "角色列表",
						Route:  "/admins/role/index",
						Sort:   2,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "security",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "菜单管理",
						Route:  "/admins/menu/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "menu",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "操作日志",
						Route:  "/admins/logs/index",
						Sort:   4,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "assignment",
							VueFile: "table.vue",
						},
					},
				},
			},
			{
				Name:   "系统配置",
				Sort:   91,
				Status: AdminMenuStatusEnabled,
				Level:  1,
				Data: AdminMenuData{
					Icon:    "settings",
					VueFile: "",
				},
				Children: []*AdminMenuDisplayData{
					{
						Name:   "前台菜单",
						Route:  "/systems/menu/index",
						Sort:   1,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "menu",
							VueFile: "table.vue",
						},
					},

					{
						Name:   "翻译配置",
						Route:  "/systems/translate/index",
						Sort:   3,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "translate",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "文章配置",
						Route:  "/systems/article/index",
						Sort:   4,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "article",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "国家配置",
						Route:  "/systems/country/index",
						Sort:   5,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "public",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "语言配置",
						Route:  "/systems/lang/index",
						Sort:   6,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "language",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "渠道管理",
						Route:  "/systems/channel/index",
						Sort:   7,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "device_hub",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "等级配置",
						Route:  "/systems/level/index",
						Sort:   8,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "star_border",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "通知配置",
						Route:  "/systems/notify/index",
						Sort:   9,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "notifications",
							VueFile: "table.vue",
						},
					},
					{
						Name:   "虚拟地址",
						Route:  "/systems/address/index",
						Sort:   10,
						Status: AdminMenuStatusEnabled,
						Level:  2,
						Data: AdminMenuData{
							Icon:    "library_books",
							VueFile: "table.vue",
						},
					},
				},
			},
		}

		menus := make([]*AdminMenu, 0)
		var nextID uint = 1
		var createMenu func(menu *AdminMenuDisplayData, parentID uint) *AdminMenu
		createMenu = func(menu *AdminMenuDisplayData, parentID uint) *AdminMenu {
			adminMenu := &AdminMenu{
				BaseModel: BaseModel{ID: nextID},
				ParentID:  parentID,
				Name:      menu.Name,
				Route:     menu.Route,
				Sort:      menu.Sort,
				Status:    menu.Status,
				Data:      menu.Data,
			}
			menus = append(menus, adminMenu)
			nextID++

			for _, child := range menu.Children {
				createMenu(child, adminMenu.ID)
			}

			return adminMenu
		}

		for _, menu := range menuData {
			createMenu(menu, 0)
		}

		return db.Create(&menus).Error
	}

	return nil
}
