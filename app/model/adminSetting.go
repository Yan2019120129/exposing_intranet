package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	AdminSettingGroupDefault  uint = 1  // 默认分组
	AdminSettingGroupBasic    uint = 10 // 基础分组
	AdminSettingGroupHome     uint = 11 // 首页分组
	AdminSettingGroupWallet   uint = 12 // 钱包分组
	AdminSettingGroupTemplate uint = 13 // 模版分组
	AdminSettingGroupStore    uint = 14 // 商城配置

	AdminSettingTypeInput int8 = 1 // Input类型
)

// AdminSettingGroupOptions 管理设置分组Options
var AdminSettingGroupOptions = []map[string]any{
	{"label": "默认配置", "Value": AdminSettingGroupDefault},
	{"label": "基础配置", "Value": AdminSettingGroupBasic},
	{"label": "首页配置", "Value": AdminSettingGroupHome},
	{"label": "钱包配置", "Value": AdminSettingGroupWallet},
	{"label": "模版配置", "Value": AdminSettingGroupTemplate},
	{"label": "商城配置", "Value": AdminSettingGroupStore},
}

// AdminSetting 管理设置表
type AdminSetting struct {
	BaseModel
	AdminID uint   `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_field;comment:管理ID" json:"adminId"`
	GroupID uint   `gorm:"type:int unsigned not null;index;comment:分组ID" json:"groupId"`
	Name    string `gorm:"type:varchar(60) not null;comment:名称" json:"name"`
	Type    int8   `gorm:"type:tinyint not null;default:1;comment:类型(1:Input类型)" json:"type"`
	Field   string `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_field;comment:键名" json:"field"`
	Value   string `gorm:"type:text;comment:键值" json:"value"`
	//Data    AdminSettingData `gorm:"type:json;comment:配置" json:"data,omitempty"`
}

// AdminSettingData 管理设置配置数据
type AdminSettingData struct {
	Input        Map            `json:"input"`
	ChildrenForm map[string]Map `json:"childrenForm"`
}

// Value implements the driver.Valuer interface
func (d AdminSettingData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *AdminSettingData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&AdminSetting{}); err != nil {
//		panic("Failed to auto migrate AdminSetting table: " + err.Error())
//	}
//
//	// Initialize admin settings
//	if err := InitAdminSettings(db); err != nil {
//		panic("Failed to initialize admin settings: " + err.Error())
//	}
//}

// InitAdminSettings 初始化管理配置
func InitAdminSettings(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AdminSetting{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		settings := []AdminSetting{
			// AdminSettingGroupDefault 默认分组
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupDefault,
				Name:    "前端模版[商户设置]",
				Type:    AdminSettingTypeInput,
				Field:   "merchantTemplate",
				Value:   "default",
				//Data: AdminSettingData{Input: views.Input{Label: "前端模版", Field: "valueInterface", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//	{Label: "默认模版", Value: "default"},
				//}}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupDefault,
				Name:    "代理数量[商户设置]",
				Type:    AdminSettingTypeInput,
				Field:   "merchantAgentNums",
				Value:   "10",
				//Data:    AdminSettingData{Input: views.Input{Label: "代理数量", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupDefault,
				Name:    "虚拟用户数量[商户设置]",
				Type:    AdminSettingTypeInput,
				Field:   "virtualUserNums",
				Value:   "100",
				//Data:    AdminSettingData{Input: views.Input{Label: "虚拟用户数量", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupDefault,
				Name:    "设备数量[商户设置]",
				Type:    AdminSettingTypeInput,
				Field:   "merchantDeviceNums",
				Value:   "5",
				//Data:    AdminSettingData{Input: views.Input{Label: "设备数量", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupDefault,
				Name:    "过期时间[商户设置]",
				Type:    AdminSettingTypeInput,
				Field:   "merchantTokenExpire",
				//Value:   utils.ObjToString(&views.FormatTimePicker{Type: views.TimeTypeMonth, Value: 1}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "过期时间", Field: "valueInterface", Type: views.InputTypeStruct},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {Inputs: [][]*views.Input{
				//			{
				//				{Label: "类型", Field: "type", Type: views.InputTypeSelect, Options: views.FormatTimePickerOptions},
				//				{Label: "时间", Field: "value", Type: views.InputTypeNumber},
				//			},
				//		}},
				//	},
				//},
			},

			// AdminSettingGroupBasic 基础分组
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "站点Logo",
				Type:    AdminSettingTypeInput,
				Field:   "siteLogo",
				Value:   "/logo.png",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点Logo", Field: "valueInterface", Type: views.InputTypeImage}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "站点图标",
				Type:    AdminSettingTypeInput,
				Field:   "siteIcon",
				Value:   "/logo.png",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点图标", Field: "valueInterface", Type: views.InputTypeImage}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "站点名称",
				Type:    AdminSettingTypeInput,
				Field:   "siteName",
				Value:   "Shoppe",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点名称", Field: "valueInterface", Type: views.InputTypeText}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "默认头像",
				Type:    AdminSettingTypeInput,
				Field:   "userAvatar",
				Value:   "/userAvatar.png",
				//Data:    AdminSettingData{Input: views.Input{Label: "默认头像", Field: "valueInterface", Type: views.InputTypeImage}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "时区",
				Type:    AdminSettingTypeInput,
				Field:   "siteTimezone",
				Value:   "Asia/Shanghai",
				//	Data: AdminSettingData{Input: views.Input{Label: "时区", Field: "valueInterface", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//		{Label: "上海", Value: "Asia/Shanghai"},
				//		{Label: "纽约", Value: "America/New_York"},
				//		{Label: "伦敦", Value: "Europe/London"},
				//		{Label: "东京", Value: "Asia/Tokyo"},
				//		{Label: "巴黎", Value: "Europe/Paris"},
				//		{Label: "悉尼", Value: "Australia/Sydney"},
				//		{Label: "莫斯科", Value: "Europe/Moscow"},
				//		{Label: "迪拜", Value: "Asia/Dubai"},
				//		{Label: "洛杉矶", Value: "America/Los_Angeles"},
				//	}}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "最低信用分",
				Type:    AdminSettingTypeInput,
				Field:   "siteMinCreditScore",
				Value:   "60",
				//Data:    AdminSettingData{Input: views.Input{Label: "最低信用分", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "提示音配置",
				Type:    AdminSettingTypeInput,
				Field:   "siteSound",
				//Value: utils.StructToString([]*views.CheckboxOption{
				//	{Label: "余额充值", Value: "balanceDeposit", Checked: true},
				//	{Label: "资产充值", Value: "assetsDeposit", Checked: true},
				//	{Label: "余额提现", Value: "balanceWithdraw", Checked: true},
				//	{Label: "资产提现", Value: "assetsWithdraw", Checked: true},
				//	{Label: "用户认证", Value: "userAuth", Checked: true},
				//}),
				//Data: AdminSettingData{Input: views.Input{Label: "提示音", Field: "valueInterface", Type: views.InputTypeCheckbox, Options: []*views.CheckboxOption{
				//	{Label: "余额充值", Value: "balanceDeposit", Checked: false}, {Label: "资产充值", Value: "assetsDeposit", Checked: false},
				//	{Label: "余额提现", Value: "balanceWithdraw", Checked: false}, {Label: "资产提现", Value: "assetsWithdraw", Checked: false},
				//	{Label: "用户认证", Value: "userAuth", Checked: false},
				//}}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "用户状态配置",
				Type:    AdminSettingTypeInput,
				Field:   "siteUserStatus",
				//Value: utils.StructToString([]*views.CheckboxOption{
				//	{Label: "冻结禁止提现", Value: "freezeWithdraw", Checked: true},
				//	{Label: "冻结禁止订单", Value: "freezeOrder", Checked: true},
				//	{Label: "信用分禁止提现", Value: "creditWithdraw", Checked: true},
				//	{Label: "信用分禁止订单", Value: "creditOrder", Checked: true},
				//	{Label: "冻结禁止转账", Value: "freezeTransfer", Checked: true},
				//	{Label: "信用分禁止转账", Value: "creditTransfer", Checked: true},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "用户状态", Field: "valueInterface", Type: views.InputTypeCheckbox, Options: []*views.CheckboxOption{
				//		{Label: "冻结禁止提现", Value: "freezeWithdraw", Checked: false}, {Label: "冻结禁止订单", Value: "freezeOrder", Checked: false},
				//		{Label: "信用分禁止提现", Value: "creditWithdraw", Checked: false}, {Label: "信用分禁止订单", Value: "creditOrder", Checked: false},
				//		{Label: "冻结禁止转账", Value: "freezeTransfer", Checked: false}, {Label: "信用分禁止转账", Value: "creditTransfer", Checked: false},
				//	}},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "下单状态配置",
				Type:    AdminSettingTypeInput,
				Field:   "siteOrderStatus",
				//Value: utils.StructToString([]*views.CheckboxOption{
				//	{Label: "禁止用户下单", Value: "prohibitUserOrder", Checked: false},
				//	{Label: "禁止商家下单", Value: "prohibitStoreOrder", Checked: false},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "下单状态", Field: "valueInterface", Type: views.InputTypeCheckbox, Options: []*views.CheckboxOption{
				//		{Label: "禁止用户下单", Value: "prohibitUserOrder", Checked: false},
				//		{Label: "禁止商家下单", Value: "prohibitStoreOrder", Checked: false},
				//	}},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "APP文件",
				Type:    AdminSettingTypeInput,
				Field:   "siteAppFile",
				Value:   `{"ios":"/app/ios.apk","android":"/app/android.apk"}`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "APP文件", Field: "valueInterface", Type: views.InputTypeStruct},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{{Label: "IOS", Field: "ios", Type: views.InputTypeFile}},
				//				{{Label: "Android", Field: "android", Type: views.InputTypeFile}},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "底部[社区]",
				Type:    AdminSettingTypeInput,
				Field:   "siteFooterCommunity",
				Value:   `[{"icon":"/social/discord.png","link":"","target":"_self"},{"icon":"/social/telegram.png","link":"","target":"_self"},{"icon":"/social/tiktok.png","link":"","target":"_self"},{"icon":"/social/facebook.png","link":"","target":"_self"},{"icon":"/social/twitter.png","link":"","target":"_self"},{"icon":"/social/reddit.png","link":"","target":"_self"},{"icon":"/social/instagram.png","link":"","target":"_self"},{"icon":"/social/youtube.png","link":"","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "底部[社区]", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "图标", Field: "icon", Type: views.InputTypeImage},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "底部[关于我们]",
				Type:    AdminSettingTypeInput,
				Field:   "siteFooterAboutUs",
				Value:   `[{"name":"articleTitle_aboutUs","link":"/article/aboutUs","target":"_self"},{"name":"articleTitle_contactUs","link":"/article/contactUs","target":"_self"},{"name":"articleTitle_privacyPolicy","link":"/article/privacyPolicy","target":"_self"},{"name":"articleTitle_termsOfService","link":"/article/termsOfService","target":"_self"},{"name":"articleTitle_missionAndVision","link":"/article/missionAndVision","target":"_self"},{"name":"articleTitle_teamCulture","link":"/article/teamCulture","target":"_self"},{"name":"articleTitle_joinUs","link":"/article/joinUs","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "底部[关于我们]", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "名称", Field: "name", Type: views.InputTypeTranslate},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "底部[服务]",
				Type:    AdminSettingTypeInput,
				Field:   "siteFooterService",
				Value:   `[{"name":"articleTitle_registerTerms","link":"/article/registerTerms","target":"_self"},{"name":"articleTitle_userGuidelines","link":"/article/userGuidelines","target":"_self"},{"name":"articleTitle_accountSecurity","link":"/article/accountSecurity","target":"_self"},{"name":"articleTitle_paymentMethods","link":"/article/paymentMethods","target":"_self"},{"name":"articleTitle_legalLiability","link":"/article/legalLiability","target":"_self"},{"name":"articleTitle_disclaimer","link":"/article/disclaimer","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "底部[服务]", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "名称", Field: "name", Type: views.InputTypeTranslate},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "底部[支持]",
				Type:    AdminSettingTypeInput,
				Field:   "siteFooterSupport",
				Value:   `[{"name":"articleTitle_faq","link":"/article/faq","target":"_self"},{"name":"articleTitle_techSupport","link":"/article/techSupport","target":"_self"},{"name":"articleTitle_tutorial","link":"/article/tutorial","target":"_self"},{"name":"articleTitle_feedback","link":"/article/feedback","target":"_self"},{"name":"articleTitle_download","link":"/article/download","target":"_self"},{"name":"articleTitle_afterSale","link":"/article/afterSale","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "底部[支持]", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "名称", Field: "name", Type: views.InputTypeTranslate},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "底部[产品]",
				Type:    AdminSettingTypeInput,
				Field:   "siteFooterProduct",
				Value:   `[{"name":"articleTitle_latestProducts","link":"/article/latestProducts","target":"_self"},{"name":"articleTitle_evaluation","link":"/article/evaluation","target":"_self"},{"name":"articleTitle_productCatalog","link":"/article/productCatalog","target":"_self"},{"name":"articleTitle_crossBorderPayment","link":"/article/crossBorderPayment","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "底部[产品]", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "名称", Field: "name", Type: views.InputTypeTranslate},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "帮助中心",
				Type:    AdminSettingTypeInput,
				Field:   "siteHelpCenter",
				Value:   `[{"name":"articleTitle_helperAccount", "content": "articleContent_helperAccount", "link":"/article/helperAccount","target":"_self"},{"name":"articleTitle_helperPayment", "content": "articleContent_helperPayment", "link":"/article/helperPayment","target":"_self"},{"name":"articleTitle_helperSecurity", "content": "articleContent_helperSecurity", "link":"/article/helperSecurity","target":"_self"},{"name":"articleTitle_helperService", "content": "articleContent_helperService", "link":"/article/helperService","target":"_self"},{"name":"articleTitle_helperFaq", "content": "articleContent_helperFaq", "link":"/article/helperFaq","target":"_self"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "帮助中心", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "名称", Field: "name", Type: views.InputTypeTranslate},
				//					{Label: "内容", Field: "content", Type: views.InputTypeTranslate},
				//					{Label: "链接", Field: "link", Type: views.InputTypeText},
				//					{Label: "方式", Field: "target", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "内部", Value: "_self"}, {Label: "外部", Value: "_blank"},
				//					}},
				//				},
				//			},
				//		},
				//	},
				//},
			},

			// AdminSettingGroupHome 首页分组
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupHome,
				Name:    "站点公告",
				Type:    AdminSettingTypeInput,
				Field:   "siteNotice",
				Value:   "settingNotice",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点公告", Field: "valueInterface", Type: views.InputTypeTranslate}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupHome,
				Name:    "站点信息",
				Type:    AdminSettingTypeInput,
				Field:   "siteIntroduce",
				Value:   "settingIntroduce",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点信息", Field: "valueInterface", Type: views.InputTypeTranslate}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupBasic,
				Name:    "站点版权",
				Type:    AdminSettingTypeInput,
				Field:   "siteCopyright",
				Value:   "Copyright © 2024 BaJie. All rights reserved.",
				//Data:    AdminSettingData{Input: views.Input{Label: "站点版权", Field: "valueInterface", Type: views.InputTypeTextarea}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupHome,
				Name:    "首页轮播图",
				Type:    AdminSettingTypeInput,
				Field:   "siteBanner",
				Value:   `[{"title":"banner1Title","image":"/images/banner1.jpg","desc":"banner1Desc"},{"title":"banner2Title","image":"/images/banner2.jpg","desc":"banner2Desc"}]`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "首页轮播图", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{{Label: "标题", Field: "title", Type: views.InputTypeTranslate},
				//					{Label: "描述", Field: "desc", Type: views.InputTypeTranslate}},
				//				{{Label: "图片", Field: "image", Type: views.InputTypeImage}},
				//			},
				//		},
				//	},
				//},
			},

			// AdminSettingGroupWallet 钱包分组
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "钱包账单",
				Type:    AdminSettingTypeInput,
				Field:   "walletBillOptions",
				//Value:   utils.StructToString(GetWalletBillCheckboxOptions()),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "钱包账单", Field: "valueInterface", Type: views.InputTypeCheckbox, Options: GetWalletBillCheckboxOptions()},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "闪兑手续费",
				Type:    AdminSettingTypeInput,
				Field:   "walletExchangeFee",
				Value:   "0.01",
				//Data:    AdminSettingData{Input: views.Input{Label: "闪兑手续费", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "推荐奖励配置",
				Type:    AdminSettingTypeInput,
				Field:   "walletRecommendReward",
				//Value:   utils.StructToString(&AdminSettingInviteRegisterReward{Register: 10, Invite: 10}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "推荐奖励", Field: "valueInterface", Type: views.InputTypeStruct},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{{Label: "注册奖励", Field: "register", Type: views.InputTypeNumber}},
				//				{{Label: "邀请奖励", Field: "invite", Type: views.InputTypeNumber}},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "分销奖励配置",
				Type:    AdminSettingTypeInput,
				Field:   "walletDistributionReward",
				//Value: utils.StructToString([]*AdminSettingDistributionReward{
				//	{Level: 1, Type: BillTypeProductPurchase, Rate: 30},
				//	{Level: 1, Type: BillTypeProductEarnings, Rate: 30},
				//	{Level: 2, Type: BillTypeProductPurchase, Rate: 20},
				//	{Level: 2, Type: BillTypeProductEarnings, Rate: 20},
				//	{Level: 3, Type: BillTypeProductPurchase, Rate: 10},
				//	{Level: 3, Type: BillTypeProductEarnings, Rate: 10},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "分销奖励", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "分销级数", Field: "level", Type: views.InputTypeNumber},
				//					{Label: "账单类型", Field: "type", Type: views.InputTypeSelect, Options: []*views.SelectOption{
				//						{Label: "购买产品", Value: BillTypeProductPurchase}, {Label: "产品收益", Value: BillTypeProductEarnings},
				//					}},
				//					{Label: "收益比例(%)", Field: "rate", Type: views.InputTypeNumber},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "提现账户数量",
				Type:    AdminSettingTypeInput,
				Field:   "walletAccountNums",
				Value:   "1",
				//Data:    AdminSettingData{Input: views.Input{Label: "钱包余额", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "钱包汇率",
				Type:    AdminSettingTypeInput,
				Field:   "walletRate",
				Value:   "1",
				//Data:    AdminSettingData{Input: views.Input{Label: "钱包汇率", Field: "valueInterface", Type: views.InputTypeNumber}},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupWallet,
				Name:    "提现次数",
				Type:    AdminSettingTypeInput,
				Field:   "walletWithdrawNums",
				Value:   `{"day": 1, "nums": 1}`,
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "提现次数", Field: "valueInterface", Type: views.InputTypeStruct},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{{Label: "每天", Field: "day", Type: views.InputTypeNumber}},
				//				{{Label: "次数", Field: "nums", Type: views.InputTypeNumber}},
				//			},
				//		},
				//	},
				//},
			},

			// AdminSettingGroupTemplate 模版分组
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "基础模版配置",
				Type:    AdminSettingTypeInput,
				Field:   "basicSettings",
				//Value: utils.StructToString([]*views.CheckboxOption{
				//	{Label: "显示首页客服", Value: "showChats", Checked: true},
				//	{Label: "显示资产", Value: "showAssets", Checked: false},
				//	{Label: "显示站内转账", Value: "showTransfer", Checked: false},
				//	{Label: "显示等级", Value: "showLevel", Checked: true},
				//	{Label: "显示积分", Value: "showScore", Checked: true},
				//	{Label: "显示认证", Value: "showAuth", Checked: true},
				//	{Label: "显示认证查看", Value: "showAuthView", Checked: true},
				//	{Label: "显示提现账户编辑", Value: "showWithdrawAccountEdit", Checked: true},
				//	{Label: "显示提现账户删除", Value: "showWithdrawAccountDelete", Checked: true},
				//	{Label: "显示提现账户号码", Value: "showWithdrawAccountNumber", Checked: true},
				//	{Label: "显示修改密码", Value: "showModifyPassword", Checked: true},
				//	{Label: "显示修改提现密码", Value: "showModifyWithdrawPassword", Checked: true},
				//	{Label: "显示绑定邮箱", Value: "showBindEmail", Checked: true},
				//	{Label: "显示绑定手机", Value: "showBindTelephone", Checked: true},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "基础模版配置", Field: "valueInterface", Type: views.InputTypeCheckbox, Options: []*views.CheckboxOption{
				//		{Label: "显示首页客服", Value: "showChats", Checked: false},
				//		{Label: "显示资产", Value: "showAssets", Checked: false},
				//		{Label: "显示站内转账", Value: "showTransfer", Checked: false},
				//		{Label: "显示等级", Value: "showLevel", Checked: false},
				//		{Label: "显示积分", Value: "showScore", Checked: false},
				//		{Label: "显示认证", Value: "showAuth", Checked: false},
				//		{Label: "显示认证查看", Value: "showAuthView", Checked: false},
				//		{Label: "显示提现账户编辑", Value: "showWithdrawAccountEdit", Checked: false},
				//		{Label: "显示提现账户删除", Value: "showWithdrawAccountDelete", Checked: false},
				//		{Label: "显示提现账户号码", Value: "showWithdrawAccountNumber", Checked: false},
				//		{Label: "显示修改密码", Value: "showModifyPassword", Checked: false},
				//		{Label: "显示修改提现密码", Value: "showModifyWithdrawPassword", Checked: false},
				//		{Label: "显示绑定邮箱", Value: "showBindEmail", Checked: false},
				//		{Label: "显示绑定手机", Value: "showBindTelephone", Checked: false},
				//	}},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "邮箱登录模版",
				Type:    AdminSettingTypeInput,
				Field:   "emailLoginTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "email", Field: "email", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "password", Field: "password", Type: vues.InputTypePassword, Style: "width: 100%;"}},
				//	{{Label: "captcha", Field: "captcha", Type: vues.InputTypeCaptcha, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "邮箱登录模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "手机登录模版",
				Type:    AdminSettingTypeInput,
				Field:   "phoneLoginTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "telephone", Field: "telephone", Type: vues.InputTypeTelephone, Style: "width: 100%;"}},
				//	{{Label: "password", Field: "password", Type: vues.InputTypePassword, Style: "width: 100%;"}},
				//	{{Label: "captcha", Field: "captcha", Type: vues.InputTypeCaptcha, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "手机登录模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "邮箱注册模版",
				Type:    AdminSettingTypeInput,
				Field:   "emailRegisterTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "email", Field: "email", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "nickname", Field: "nickname", Type: vues.InputTypeText, Display: true, Style: "width: 100%;"}},
				//	{{Label: "sex", Field: "sex", Type: vues.InputTypeSelect, Options: []*views.SelectOption{{Label: "sexMale", Value: SexMale}, {Label: "sexFemale", Value: SexFemale}}, Display: true, Style: "width: 100%;"}},
				//	{{Label: "password", Field: "password", Type: vues.InputTypePassword, Style: "width: 100%;"}},
				//	{{Label: "confirmPassword", Field: "confirmPassword", Type: vues.InputTypePassword, Display: true, Rules: "(val) => {return val == formParams.value.password || $i18n.t('passwordNotMatch')}", Style: "width: 100%;"}},
				//	{{Label: "securityKey", Field: "securityKey", Type: vues.InputTypePassword, Display: true, Style: "width: 100%;"}},
				//	{{Label: "confirmSecurityKey", Field: "confirmSecurityKey", Type: vues.InputTypePassword, Display: true, Rules: "(val) => {return val == formParams.value.securityKey || $i18n.t('securityKeyNotMatch')}", Style: "width: 100%;"}},
				//	{{Label: "inviteCode", Field: "inviteCode", Type: vues.InputTypeText, Display: true, Style: "width: 100%;"}},
				//	{{Label: "captcha", Field: "captcha", Type: vues.InputTypeCaptcha, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "邮箱注册模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "手机注册模版",
				Type:    AdminSettingTypeInput,
				Field:   "phoneRegisterTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "telephone", Field: "telephone", Type: vues.InputTypeTelephone, Style: "width: 100%;"}},
				//	{{Label: "nickname", Field: "nickname", Type: vues.InputTypeText, Display: true, Style: "width: 100%;"}},
				//	{{Label: "sex", Field: "sex", Type: vues.InputTypeSelect, Options: []*views.SelectOption{{Label: "sexMale", Value: SexMale}, {Label: "sexFemale", Value: SexFemale}}, Display: true, Style: "width: 100%;"}},
				//	{{Label: "password", Field: "password", Type: vues.InputTypePassword, Style: "width: 100%;"}},
				//	{{Label: "confirmPassword", Field: "confirmPassword", Type: vues.InputTypePassword, Display: true, Rules: "(val) => {return val == formParams.value.password || $i18n.t('passwordNotMatch')}", Style: "width: 100%;"}},
				//	{{Label: "securityKey", Field: "securityKey", Type: vues.InputTypePassword, Display: true, Style: "width: 100%;"}},
				//	{{Label: "confirmSecurityKey", Field: "confirmSecurityKey", Type: vues.InputTypePassword, Display: true, Rules: "(val) => {return val == formParams.value.securityKey || $i18n.t('securityKeyNotMatch')}", Style: "width: 100%;"}},
				//	{{Label: "inviteCode", Field: "inviteCode", Type: vues.InputTypeText, Display: true, Style: "width: 100%;"}},
				//	{{Label: "captcha", Field: "captcha", Type: vues.InputTypeCaptcha, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "手机注册模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "账户设置模版",
				Type:    AdminSettingTypeInput,
				Field:   "userSettingTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Icon: "/inputs/avatar.png", Label: "avatar", Small: "avatarSmall", Field: "avatar", Type: vues.InputTypeAvatar, Style: "width: 100%;"}},
				//	{{Icon: "/inputs/nickname.png", Label: "nickname", Small: "nicknameSmall", Field: "nickname", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Icon: "/inputs/sex.png", Label: "sex", Small: "sexSmall", Field: "sex", Type: vues.InputTypeSelect, Options: []*views.SelectOption{
				//		{Label: "sexMale", Value: SexMale},
				//		{Label: "sexFemale", Value: SexFemale},
				//		{Label: "sexUnknown", Value: SexUnknown},
				//	}, Style: "width: 100%;"}},
				//	{{Icon: "/inputs/birthday.png", Label: "birthday", Small: "birthdaySmall", Field: "birthday", Type: vues.InputTypeDatePicker, Style: "width: 100%;"}},
				//	{{Icon: "/inputs/introduction.png", Label: "introduction", Small: "introductionSmall", Field: "desc", Type: vues.InputTypeTextarea, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "账户设置模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "图标", Field: "icon", Type: views.InputTypeIcon},
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "身份认证模版",
				Type:    AdminSettingTypeInput,
				Field:   "userAuthTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "idName", Field: "realName", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "idNumber", Field: "idNumber", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "idAddress", Field: "address", Type: vues.InputTypeTextarea, Style: "width: 100%;"}},
				//	{
				//		{Label: "idPhoto1", Field: "photo1", Type: vues.InputTypeImage},
				//		{Label: "idPhoto2", Field: "photo2", Type: vues.InputTypeImage},
				//		{Label: "idPhoto3", Field: "photo3", Type: vues.InputTypeImage},
				//	},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "身份认证模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "提现账户银行卡模版",
				Type:    AdminSettingTypeInput,
				Field:   "withdrawAccountBankTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "accountBankName", Field: "bankName", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountBankCardNo", Field: "bankCardNo", Type: vues.InputTypeText, Style: "width: 100%;", Mask: "NNNNNNNNNNNNNNNN"}},
				//	{{Label: "accountRealName", Field: "realName", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountBankAddress", Field: "bankAddress", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountBankCode", Field: "bankCode", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "提现账户银行卡模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupTemplate,
				Name:    "提现账户数字货币模版",
				Type:    AdminSettingTypeInput,
				Field:   "withdrawAccountAssetTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "accountTokenName", Field: "bankName", Readonly: true, Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountTokenAddress", Field: "bankCardNo", Type: vues.InputTypeText, Style: "width: 100%;", Mask: "NNNNNNNNNNNNNNNN"}},
				//	{{Label: "accountToken", Field: "realName", Readonly: true, Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountTokenDesc", Field: "bankAddress", Readonly: true, Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "accountTokenSymbol", Field: "bankCode", Readonly: true, Type: vues.InputTypeText, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "提现账户数字货币模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupStore,
				Name:    "超时设置",
				Type:    AdminSettingTypeInput,
				Field:   "storeTimeout",
				//Value:   utils.StructToString(AdminSettingStoreOrderTimeout{Payment: 20, Receipt: 8, Comment: 3, Delivery: 3, DeliveryNums: 10}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "超时设置", Field: "valueInterface", Type: views.InputTypeStruct, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "支付超时(分钟)", Field: "payment", Type: views.InputTypeNumber},
				//					{Label: "收货超时(天数)", Field: "receipt", Type: views.InputTypeNumber},
				//					{Label: "评论超时(天数)", Field: "comment", Type: views.InputTypeNumber},
				//					{Label: "发货超时(天数)", Field: "delivery", Type: views.InputTypeNumber},
				//					{Label: "超时信用分", Field: "deliveryNums", Type: views.InputTypeNumber},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupStore,
				Name:    "商家入驻模版",
				Type:    AdminSettingTypeInput,
				Field:   "storeSettledTemplate",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "storeLoge", Field: "loge", Type: vues.InputTypeImage, Display: true}},
				//	{{Label: "storeName", Field: "name", Type: vues.InputTypeText, Style: "width: 100%;", Display: true}},
				//	{{Label: "storeNumber", Field: "number", Type: vues.InputTypeText, Style: "width: 100%;", Display: true}},
				//	{{Label: "storeAddress", Field: "address", Type: vues.InputTypeText, Style: "width: 100%;", Display: true}},
				//	{
				//		{Label: "storePhoto1", Field: "photo1", Type: vues.InputTypeImage, Display: true},
				//		{Label: "storePhoto2", Field: "photo2", Type: vues.InputTypeImage, Display: true},
				//		{Label: "storePhoto3", Field: "photo3", Type: vues.InputTypeImage, Display: true},
				//	},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "商家入驻模版", Field: "valueInterface", Type: views.InputTypeSlice, Readonly: true},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupStore,
				Name:    "用户地址",
				Type:    AdminSettingTypeInput,
				Field:   "userAddress",
				//Value: utils.StructToString([][]*vues.Input{
				//	{{Label: "recipientName", Field: "name", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "contactDetails", Field: "contact", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "city", Field: "city", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//	{{Label: "detailedAddress", Field: "address", Type: vues.InputTypeText, Style: "width: 100%;"}},
				//}),
				//Data: AdminSettingData{
				//	Input: views.Input{Label: "用户地址", Field: "valueInterface", Type: views.InputTypeSlice},
				//	ChildrenForm: map[string]*views.Form{
				//		"valueInterface": {
				//			Inputs: [][]*views.Input{
				//				{
				//					{Label: "标题", Field: "label", Type: views.InputTypeTranslate},
				//					{Label: "只读", Field: "readonly", Type: views.InputTypeToggle},
				//					{Label: "隐藏", Field: "display", Type: views.InputTypeToggle},
				//					{Label: "验证", Field: "rules", Type: views.InputTypeText},
				//					{Label: "格式", Field: "mask", Type: views.InputTypeText},
				//				},
				//			},
				//		},
				//	},
				//},
			},
			{
				AdminID: SuperAdminID,
				GroupID: AdminSettingGroupStore,
				Name:    "入驻须知",
				Type:    AdminSettingTypeInput,
				Field:   "settledNotice",
				Value:   "settledNoticeValue",
				//Data:    AdminSettingData{Input: views.Input{Label: "入驻须知", Field: "valueInterface", Type: views.InputTypeTranslate}},
			},
		}

		return db.CreateInBatches(settings, len(settings)).Error
	}

	return nil
}
