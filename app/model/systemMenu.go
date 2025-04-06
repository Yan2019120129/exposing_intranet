package models

import (
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	MenuIconTypeDark   = "dark"   //	暗色图标
	MenuIconTypeLight  = "light"  //	亮色图标
	MenuIconTypeActive = "active" //	激活图标

	// MenuNameTranslatePrefix 菜单名称翻译前缀
	MenuNameTranslatePrefix = "menu_%s"

	// Menu types
	MenuTypeMobileNavigationBar  int8 = 1  // 移动菜单导航
	MenuTypeCommonUserMenu       int8 = 11 // 用户更多菜单
	MenuTypeCommonUserMoreMenu   int8 = 12 // 用户更多菜单
	MenuTypeCommonMoreMenu       int8 = 13 // 公用更多菜单
	MenuTypeCommonWalletMenu     int8 = 14 // 公用钱包菜单
	MenuTypeDesktopNavigationBar int8 = 21 // 电脑菜单导航

	MenuTargetSelf  = "_self"  //	当前窗口
	MenuTargetBlank = "_blank" //	新窗口

	// Menu statuses
	MenuStatusDisabled int8 = -1 // 禁用
	MenuStatusEnabled  int8 = 10 // 启用
)

// Menu 前台菜单
type Menu struct {
	BaseModel
	AdminID  uint     `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	ParentID uint     `gorm:"type:int unsigned not null;index;comment:父级ID" json:"parentId"`
	Name     string   `gorm:"type:varchar(50) not null;comment:名称" json:"name"`
	Route    string   `gorm:"type:varchar(100) not null;comment:路由" json:"route"`
	Sort     uint8    `gorm:"type:tinyint unsigned not null;default:99;index;comment:排序" json:"sort"`
	Type     int8     `gorm:"type:tinyint not null;default:1;index;comment:类型(1:导航菜单,11:用户菜单,21:快捷菜单)" json:"type"`
	Status   int8     `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Data     MenuData `gorm:"type:json;comment:数据" json:"data,omitempty"`
}

type MenuDisplayData struct {
	ID       uint               `json:"id"`                //	ID
	Name     string             `json:"name"`              //	名称
	Route    string             `json:"route"`             //	路由
	Data     MenuData           `json:"data"`              //	数据
	Children []*MenuDisplayData `json:"children" gorm:"-"` //	下级
}

// MenuData 前台菜单数据
type MenuData struct {
	Target     string `json:"target" views:"label:打开方式;type:select"`   // 打开方式 (_self, _blank)
	Label      string `json:"label" views:"label:翻译;type:translate"`   //	翻译
	Class      string `json:"class" views:"label:自定义CSS类"`             // 自定义CSS类
	DarkIcon   string `json:"darkIcon" views:"label:暗色图标;type:icon"`   // 暗色图标
	LightIcon  string `json:"lightIcon" views:"label:亮色图标;type:icon"`  // 亮色图标
	ActiveIcon string `json:"activeIcon" views:"label:激活图标;type:icon"` // 激活图标
}

// Value implements the driver.Valuer interface
func (d MenuData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *MenuData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Menu{}); err != nil {
//		panic("Failed to auto migrate Menu table: " + err.Error())
//	}
//
//	// Initialize system menus
//	if err := InitSystemMenus(db); err != nil {
//		panic("Failed to initialize system menus: " + err.Error())
//	}
//}

// InitSystemMenus initializes the default system menus
func InitSystemMenus(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Menu{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		menus := initSystemsMenus()
		return db.CreateInBatches(menus, len(menus)).Error
	}

	return nil
}

func initSystemsMenus() []Menu {
	menus := []Menu{
		// 手机端导航菜单
		{AdminID: SuperAdminID, ParentID: 0, Name: "首页", Route: "/", Sort: 1, Type: MenuTypeMobileNavigationBar, Status: MenuStatusEnabled, Data: MenuData{Label: "menuHome", LightIcon: "/menus/light/home.png", ActiveIcon: "/menus/active/home.png", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "分类", Route: "/category", Sort: 1, Type: MenuTypeMobileNavigationBar, Status: MenuStatusEnabled, Data: MenuData{Label: "menuCategory", LightIcon: "/menus/light/category.png", ActiveIcon: "/menus/active/category.png", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "购物车", Route: "/shopping/cart", Sort: 1, Type: MenuTypeMobileNavigationBar, Status: MenuStatusEnabled, Data: MenuData{Label: "menuCard", LightIcon: "/menus/light/shoppingCard.png", ActiveIcon: "/menus/active/shoppingCard.png", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "消息", Route: "/message", Sort: 1, Type: MenuTypeMobileNavigationBar, Status: MenuStatusEnabled, Data: MenuData{Label: "menuMessage", LightIcon: "/menus/light/message.png", ActiveIcon: "/menus/active/message.png", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "我的", Route: "/users", Sort: 2, Type: MenuTypeMobileNavigationBar, Status: MenuStatusEnabled, Data: MenuData{Label: "menuUser", LightIcon: "/menus/light/user.png", ActiveIcon: "/menus/active/user.png", Target: "_self", Class: ""}},

		// 我的菜单 - 总览｜用户信息
		{BaseModel: BaseModel{ID: 99}, AdminID: SuperAdminID, ParentID: 0, Name: "总览", Route: "/users", Sort: 1, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuOverview", DarkIcon: "/menus/dark/overview.svg", LightIcon: "/menus/light/overview.svg", Target: "_self", Class: ""}},

		// 我的菜单 - 钱包
		{BaseModel: BaseModel{ID: 100}, AdminID: SuperAdminID, ParentID: 0, Name: "钱包", Route: "/wallets", Sort: 1, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWallet", DarkIcon: "/menus/dark/wallet.svg", LightIcon: "/menus/light/wallet.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "钱包总览", Route: "/wallets", Sort: 1, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWalletInfo", DarkIcon: "/menus/dark/wallet.svg", LightIcon: "/menus/light/wallet.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "账单明细", Route: "/wallets/bill", Sort: 2, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWalletBill", DarkIcon: "/menus/dark/bill.svg", LightIcon: "/menus/light/bill.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "站内转账", Route: "/wallets/transfer", Sort: 3, Type: MenuTypeCommonUserMenu, Status: MenuStatusDisabled, Data: MenuData{Label: "menuWalletTransfer", DarkIcon: "/menus/dark/transfer.svg", LightIcon: "/menus/light/transfer.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "资产兑换", Route: "/wallets/swaps", Sort: 4, Type: MenuTypeCommonUserMenu, Status: MenuStatusDisabled, Data: MenuData{Label: "menuAssetsSwaps", DarkIcon: "/menus/dark/swap.svg", LightIcon: "/menus/light/swap.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "提现账户", Route: "/wallets/account", Sort: 5, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWalletAccount", DarkIcon: "/menus/dark/card.svg", LightIcon: "/menus/light/card.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "钱包充值", Route: "/wallets/deposit", Sort: 1, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "walletDeposit", DarkIcon: "/menus/dark/deposit.svg", LightIcon: "/menus/light/deposit.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 100, Name: "钱包提现", Route: "/wallets/withdraw", Sort: 2, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "walletWithdraw", DarkIcon: "/menus/dark/withdraw-account.svg", LightIcon: "/menus/light/withdraw-account.svg", Target: "_self", Class: ""}},

		// 我的菜单 - 店铺
		{BaseModel: BaseModel{ID: 110}, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "商城服务", Route: "/merchant/settled", Sort: 2, Status: MenuStatusEnabled, Data: MenuData{Label: "menuShop", DarkIcon: "/menus/dark/services.svg", LightIcon: "/menus/light/services.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "商家入驻", Route: "/merchant/settled", Sort: 1, Status: MenuStatusEnabled, Data: MenuData{Label: "menuMerchantsSettleIn", DarkIcon: "/menus/dark/settled.svg", LightIcon: "/menus/light/settled.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "商家服务", Route: "/merchant", Sort: 2, Status: MenuStatusEnabled, Data: MenuData{Label: "menuStoreService", DarkIcon: "/menus/dark/settled.svg", LightIcon: "/menus/light/settled.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "商品收藏", Route: "/collect/2", Sort: 4, Status: MenuStatusEnabled, Data: MenuData{Label: "menuFollowProduct", DarkIcon: "/menus/dark/follow.svg", LightIcon: "/menus/light/follow.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "店铺关注", Route: "/collect/1", Sort: 5, Status: MenuStatusEnabled, Data: MenuData{Label: "menuFollowStore", DarkIcon: "/menus/dark/follow.svg", LightIcon: "/menus/light/follow.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "浏览记录", Route: "/browsing", Sort: 6, Status: MenuStatusEnabled, Data: MenuData{Label: "menuBrowsingProduct", DarkIcon: "/menus/dark/browsing.svg", LightIcon: "/menus/light/browsing.svg", Target: "_self", Class: ""}},
		{ParentID: 110, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "收货地址", Route: "/address", Sort: 6, Status: MenuStatusEnabled, Data: MenuData{Label: "menuShippingAddress", DarkIcon: "/menus/dark/address.svg", LightIcon: "/menus/light/address.svg", Target: "_self", Class: ""}},

		// 我的菜单 - 账户
		{BaseModel: BaseModel{ID: 130}, AdminID: SuperAdminID, ParentID: 0, Name: "账户", Route: "/users/setting", Sort: 3, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuAccount", DarkIcon: "/menus/dark/account.svg", LightIcon: "/menus/light/account.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "会员权益", Route: "/users/level", Sort: 1, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuLevel", DarkIcon: "/menus/dark/level.svg", LightIcon: "/menus/light/level.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "实名认证", Route: "/users/auth", Sort: 2, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuAuth", DarkIcon: "/menus/dark/auth.svg", LightIcon: "/menus/light/auth.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "邀请奖励", Route: "/users/invite", Sort: 3, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuInvite", DarkIcon: "/menus/dark/invite.svg", LightIcon: "/menus/light/invite.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "我的团队", Route: "/users/teams", Sort: 4, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuTeam", DarkIcon: "/menus/dark/team.svg", LightIcon: "/menus/light/team.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "账户安全", Route: "/users/safety", Sort: 5, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuSafety", DarkIcon: "/menus/dark/safety.svg", LightIcon: "/menus/light/safety.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "账户设置", Route: "/users/setting", Sort: 6, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuSetting", DarkIcon: "/menus/dark/setting.svg", LightIcon: "/menus/light/setting.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 130, Name: "消息通知", Route: "/users/notice", Sort: 7, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuNotice", DarkIcon: "/menus/dark/notice.svg", LightIcon: "/menus/light/notice.svg", Target: "_self", Class: ""}},
		// {AdminID: SuperAdminID, ParentID: 130, Name: "在线客服", Route: "/chats", Sort: 8, Type: MenuTypeCommonUserMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuOnlineChats", DarkIcon: "/menus/dark/chat.svg", LightIcon: "/menus/light/chat.svg", Target: "_self", Class: ""}},

		// 我的菜单 -  快捷菜单
		//{BaseModel: BaseModel{ID: 120}, AdminID: SuperAdminID, Type: MenuTypeCommonUserMenu, Name: "我的订单", Route: "/product/order", Sort: 2, Status: MenuStatusEnabled, Data: MenuData{Label: "menuMyOrder", DarkIcon: "/menus/dark/order.svg", LightIcon: "/menus/light/order.svg", Target: "_self", Class: ""}},
		{ParentID: 0, AdminID: SuperAdminID, Type: MenuTypeCommonUserMoreMenu, Name: "待付款", Route: "/product/order/10", Sort: 1, Status: MenuStatusEnabled, Data: MenuData{Label: "menuPendingPayment", DarkIcon: "/menus/dark/obligation.svg", LightIcon: "/menus/light/obligation.svg", Target: "_self", Class: ""}},
		{ParentID: 0, AdminID: SuperAdminID, Type: MenuTypeCommonUserMoreMenu, Name: "待发货", Route: "/product/order/12", Sort: 2, Status: MenuStatusEnabled, Data: MenuData{Label: "menuAwaitingShipment", DarkIcon: "/menus/dark/pendingShipment.svg", LightIcon: "/menus/light/pendingShipment.svg", Target: "_self", Class: ""}},
		{ParentID: 0, AdminID: SuperAdminID, Type: MenuTypeCommonUserMoreMenu, Name: "待收货", Route: "/product/order/14", Sort: 3, Status: MenuStatusEnabled, Data: MenuData{Label: "menuAwaitingReceipt", DarkIcon: "/menus/dark/received.svg", LightIcon: "/menus/light/received.svg", Target: "_self", Class: ""}},
		{ParentID: 0, AdminID: SuperAdminID, Type: MenuTypeCommonUserMoreMenu, Name: "待评价", Route: "/comment", Sort: 4, Status: MenuStatusEnabled, Data: MenuData{Label: "menuPendingReview", DarkIcon: "/menus/dark/comment.svg", LightIcon: "/menus/light/comment.svg", Target: "_self", Class: ""}},
		{ParentID: 0, AdminID: SuperAdminID, Type: MenuTypeCommonUserMoreMenu, Name: "售后", Route: "/refund", Sort: 5, Status: MenuStatusEnabled, Data: MenuData{Label: "menuPendingRefund", DarkIcon: "/menus/dark/refund.svg", LightIcon: "/menus/light/refund.svg", Target: "_self", Class: ""}},

		//{AdminID: SuperAdminID, ParentID: 0, Name: "转账", Route: "/wallets/transfer", Sort: 3, Type: MenuTypeCommonUserMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuTransfer", DarkIcon: "/menus/dark/transfer.svg", LightIcon: "/menus/light/transfer.svg", Target: "_self", Class: ""}},

		// 公用钱包菜单 - 充值｜提现｜转账｜兑换
		{AdminID: SuperAdminID, ParentID: 0, Name: "充值", Route: "/wallets/deposit", Sort: 1, Type: MenuTypeCommonWalletMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuDeposit", DarkIcon: "/menus/dark/deposit.svg", LightIcon: "/menus/light/deposit.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "提现", Route: "/wallets/withdraw", Sort: 2, Type: MenuTypeCommonWalletMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWithdraw", DarkIcon: "/menus/dark/withdraw-account.svg", LightIcon: "/menus/light/withdraw-account.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "转账", Route: "/wallets/transfer", Sort: 3, Type: MenuTypeCommonWalletMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuTransfer", DarkIcon: "/menus/dark/transfer.svg", LightIcon: "/menus/light/transfer.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "兑换", Route: "/wallets/swaps", Sort: 4, Type: MenuTypeCommonWalletMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuSwaps", DarkIcon: "/menus/dark/swap.svg", LightIcon: "/menus/light/swap.svg", Target: "_self", Class: ""}},

		// 更多菜单
		{AdminID: SuperAdminID, ParentID: 0, Name: "充值", Route: "/wallets/deposit", Sort: 1, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuDeposit", DarkIcon: "/menus/dark/deposit.svg", LightIcon: "/menus/light/deposit.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "提现", Route: "/wallets/withdraw", Sort: 2, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuWithdraw", DarkIcon: "/menus/dark/withdraw.svg", LightIcon: "/menus/light/withdraw.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "实名认证", Route: "/users/auth", Sort: 3, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuAuth", DarkIcon: "/menus/dark/auth.svg", LightIcon: "/menus/light/auth.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "邀请奖励", Route: "/users/invite", Sort: 4, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuInvite", DarkIcon: "/menus/dark/invite.svg", LightIcon: "/menus/light/invite.svg", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "帮助中心", Route: "/helpers", Sort: 5, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuHelps", DarkIcon: "/menus/dark/helper.png", LightIcon: "/menus/light/helper.png", Target: "_self", Class: ""}},
		{AdminID: SuperAdminID, ParentID: 0, Name: "APP下载", Route: "/download", Sort: 6, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuDownload", DarkIcon: "/menus/dark/download.svg", LightIcon: "/menus/light/download.svg", Target: "_self", Class: ""}},
		// {AdminID: SuperAdminID, ParentID: 0, Name: "在线客服", Route: "/chats", Sort: 7, Type: MenuTypeCommonMoreMenu, Status: MenuStatusEnabled, Data: MenuData{Label: "menuChats", DarkIcon: "/menus/dark/chat.svg", LightIcon: "/menus/light/chat.svg", Target: "_self", Class: ""}},
	}

	return menus
}
