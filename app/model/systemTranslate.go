package models

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	DefaultLang = "zh-CN"

	TranslateEmptyValue = "- -"

	TranslateTypeSystem   int8 = 1 // 系统翻译
	TranslateTypeFrontend int8 = 2 // 前台翻译
	TranslateTypeBackend  int8 = 3 // 后台翻译

	TranslateStatusDisabled int8 = -1 // 禁用
	TranslateStatusEnabled  int8 = 10 // 启用
)

// Translate 系统语言翻译表
type Translate struct {
	BaseModel
	AdminID uint   `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_lang_field;comment:管理ID" json:"adminId"`
	Lang    string `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_lang_field;comment:标识" json:"lang"`
	Name    string `gorm:"type:varchar(60) not null;comment:名称" json:"name"`
	Type    int8   `gorm:"type:tinyint unsigned not null;default:1;comment:类型(1:系统,2:前台,3:后台)" json:"type"`
	Field   string `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_lang_field;comment:键名" json:"field"`
	Value   string `gorm:"type:text;comment:键值" json:"value"`
	Desc    string `gorm:"type:varchar(255);comment:描述" json:"desc"`
	Status  int8   `gorm:"type:tinyint not null;default:10;comment:状态(10:启用,-1:禁用)" json:"status"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Translate{}); err != nil {
//		panic("Failed to auto migrate Translate table: " + err.Error())
//	}
//
//	// Initialize system translates
//	if err := InitSystemTranslates(db); err != nil {
//		panic("Failed to initialize system translates: " + err.Error())
//	}
//}

// InitSystemTranslates 初始化系统翻译数据
func InitSystemTranslates(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Translate{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		translates := []Translate{
			// validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[必需的]", Type: TranslateTypeSystem, Field: "required", Value: "%s 值必需的 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[最小长度]", Type: TranslateTypeSystem, Field: "min", Value: "%s 最小长度 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[最大长度]", Type: TranslateTypeSystem, Field: "max", Value: "%s 最大长度 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[列表中存在]", Type: TranslateTypeSystem, Field: "oneof", Value: "%s 必需在 %s 中存在", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[大于]", Type: TranslateTypeSystem, Field: "gt", Value: "%s 大于 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[大于等级]", Type: TranslateTypeSystem, Field: "gte", Value: "%s 大于等于 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[小于]", Type: TranslateTypeSystem, Field: "lt", Value: "%s 小于 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[小于等级]", Type: TranslateTypeSystem, Field: "lte", Value: "%s 小于等于 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},
			{AdminID: 0, Lang: "zh-CN", Name: "参数验证[邮箱]", Type: TranslateTypeSystem, Field: "email", Value: "%s 邮箱格式不正确 %s", Desc: "validator 验证数据, 第一个 %s errs.Field() 第二个 %s errs.Param()"},

			// 公共翻译数据
			{AdminID: 0, Lang: "zh-CN", Name: "余额不足", Type: TranslateTypeSystem, Field: "insufficientBalance", Value: "余额不足"},
			{AdminID: 0, Lang: "zh-CN", Name: "数据库错误", Type: TranslateTypeSystem, Field: "dataError", Value: "数据库错误 - %s"},

			// 后台翻译数据
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "用户不存在", Type: TranslateTypeSystem, Field: "userNotExist", Value: "用户不存在", Desc: "用于展示在前台的用户不存在提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "资产不存在", Type: TranslateTypeSystem, Field: "assetsNotExist", Value: "资产不存在", Desc: "用于展示在前台的资产不存在提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "账户或密码不正确", Type: TranslateTypeSystem, Field: "accountOrPasswordError", Value: "账户或密码不正确"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "验证码错误", Type: TranslateTypeSystem, Field: "captchaError", Value: "验证码错误"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付密码错误", Type: TranslateTypeSystem, Field: "securityPasswordError", Value: "支付密码错误"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请码错误", Type: TranslateTypeSystem, Field: "inviteCodeError", Value: "邀请码错误"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "账户已存在", Type: TranslateTypeSystem, Field: "usernameExist", Value: "账户已存在, 请重新设置账户～"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "客服第一句问候语", Type: TranslateTypeSystem, Field: "chatsFirstGreeting", Value: "您好，有什么可以帮您的？", Desc: "用于展示在前台的客服第一句问候语"},

			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "不能为空", Type: TranslateTypeSystem, Field: "notEmpty", Value: "不能为空"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "不正确", Type: TranslateTypeSystem, Field: "notCorrect", Value: "不正确"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "格式不正确", Type: TranslateTypeSystem, Field: "formatError", Value: "格式不正确"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "数量限制", Type: TranslateTypeSystem, Field: "numberLimit", Value: "数量限制, 不能超过 <= %v"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "范围限制[数值]", Type: TranslateTypeSystem, Field: "rangeLimit", Value: "范围限制, 值 > %v 并且 < %v"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "范围限制[字符串]", Type: TranslateTypeSystem, Field: "rangeLimitString", Value: "范围限制, 值 > %s 并且 < %s"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "关于我们", Type: TranslateTypeSystem, Field: "aboutUs", Value: "关于我们", Desc: "用于展示在前台的关于信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "社区", Type: TranslateTypeSystem, Field: "community", Value: "社区", Desc: "用于展示在前台的社区信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "产品", Type: TranslateTypeSystem, Field: "product", Value: "产品", Desc: "用于展示在前台的产品信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "服务", Type: TranslateTypeSystem, Field: "service", Value: "服务", Desc: "用于展示在前台的服务信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支持", Type: TranslateTypeSystem, Field: "support", Value: "支持", Desc: "用于展示在前台的支持信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站点公告", Type: TranslateTypeSystem, Field: "settingNotice", Value: "欢迎来到我们的网站！我们正在进行系统升级，可能会影响部分功能。感谢您的理解和支持。", Desc: "用于显示在前台的站点公告内容"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站点信息", Type: TranslateTypeSystem, Field: "settingIntroduce", Value: "我们是一家专注于提供高质量服务的公司。我们的使命是为客户创造价值，为员工提供发展机会，为社会做出贡献。", Desc: "用于展示在前台的站点简介信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请规则内容", Type: TranslateTypeSystem, Field: "inviteRuleContent", Value: `<div><font size="3" style="line-height: 28px;"><div>1.活动对象：本活动面向所有已注册并登录的用户。只要您是平台的有效用户，均可参与此活动并通过邀请好友获取奖励。</div>
<div>2.邀请方式: 用户通过个人专属邀请链接或邀请码邀请好友注册，好友成功注册后，用户将获得奖励。</div>
<div>3.注册要求: 被邀请的好友必须为首次注册且在7天内完成注册, 才能视为有效邀请。</div>
<div>4.奖励内容: 每成功邀请一位符合条件的好友, 您可获得${inviteReward}现金奖励, 好友可获得${signupReward}的现金奖励以及新人大礼包。累计邀请20位好友, 您将免费获得一个月的钻石会员。</div>
<div>5.奖励发放: 奖励将在好友完成任务后的3个工作日内发放至邀请人的账户。</div>
<div>6.奖励限制: 被邀请的好友必须为新用户, 即未在平台注册过账户。邀请奖励仅通过平台提供的有效链接或邀请码进行才有效。如发现恶意行为, 平台有权取消相关奖励。</div></font></div>`, Desc: "用于展示在前台的邀请规则内容"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[我的]", Type: TranslateTypeSystem, Field: "menuUser", Value: "我的", Desc: "用于展示在前台的我的按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[首页]", Type: TranslateTypeSystem, Field: "menuHome", Value: "首页", Desc: "用于展示在前台的首页按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[APP下载]", Type: TranslateTypeSystem, Field: "menuDownload", Value: "APP下载", Desc: "用于展示在前台的APP下载按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[帮助中心]", Type: TranslateTypeSystem, Field: "menuHelps", Value: "帮助中心", Desc: "用于展示在前台的帮助中心按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[在线客服]", Type: TranslateTypeSystem, Field: "menuOnlineChats", Value: "在线客服", Desc: "用于展示在前台的在线客服按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[客服]", Type: TranslateTypeSystem, Field: "menuChats", Value: "客服", Desc: "用于展示在前台的客服按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[充值]", Type: TranslateTypeSystem, Field: "menuDeposit", Value: "充值", Desc: "用于展示在前台的充值按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[提现]", Type: TranslateTypeSystem, Field: "menuWithdraw", Value: "提现", Desc: "用于展示在前台的提现按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[兑换]", Type: TranslateTypeSystem, Field: "menuSwaps", Value: "兑换", Desc: "用于展示在前台的兑换按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[转账]", Type: TranslateTypeSystem, Field: "menuTransfer", Value: "转账", Desc: "用于展示在前台的转账按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[订单]", Type: TranslateTypeSystem, Field: "menuOrder", Value: "订单", Desc: "用于展示在前台的订单按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[账户]", Type: TranslateTypeSystem, Field: "menuAccount", Value: "账户", Desc: "用于展示在前台的账户按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[认证管理]", Type: TranslateTypeSystem, Field: "menuAuth", Value: "认证管理", Desc: "用于展示在前台的认证管理按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[邀请奖励]", Type: TranslateTypeSystem, Field: "menuInvite", Value: "邀请奖励", Desc: "用于展示在前台的邀请奖励按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[账户安全]", Type: TranslateTypeSystem, Field: "menuSafety", Value: "账户安全", Desc: "用于展示在前台的账户安全按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[会员权益]", Type: TranslateTypeSystem, Field: "menuLevel", Value: "会员权益", Desc: "用于展示在前台的会员权益按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[我的团队]", Type: TranslateTypeSystem, Field: "menuTeam", Value: "我的团队", Desc: "用于展示在前台的我的团队按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[钱包]", Type: TranslateTypeSystem, Field: "menuWallet", Value: "钱包", Desc: "用于展示在前台的钱包按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[钱包总览]", Type: TranslateTypeSystem, Field: "menuWalletInfo", Value: "钱包总览", Desc: "用于展示在前台的钱包总览按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[账单明细]", Type: TranslateTypeSystem, Field: "menuWalletBill", Value: "账单明细", Desc: "用于展示在前台的账单明细按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[站内转账]", Type: TranslateTypeSystem, Field: "menuWalletTransfer", Value: "站内转账", Desc: "用于展示在前台的站内转账按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[资产兑换]", Type: TranslateTypeSystem, Field: "menuAssetsSwaps", Value: "资产兑换", Desc: "用于展示在前台的资产兑换按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[提现账户]", Type: TranslateTypeSystem, Field: "menuWalletAccount", Value: "提现账户", Desc: "用于展示在前台的提现账户按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[账户设置]", Type: TranslateTypeSystem, Field: "menuSetting", Value: "账户设置", Desc: "用于展示在前台的账户设置按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[消息通知]", Type: TranslateTypeSystem, Field: "menuNotice", Value: "消息通知", Desc: "用于展示在前台的消息通知按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "菜单[总览]", Type: TranslateTypeSystem, Field: "menuOverview", Value: "总览", Desc: "用于展示在前台的总览按钮"},

			// 前台翻译数据

			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "联系客服", Type: TranslateTypeFrontend, Field: "contactCustomerService", Value: "联系客服", Desc: "用于展示在前台的联系客服"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未设置在线客服", Type: TranslateTypeFrontend, Field: "noOnlineChats", Value: "未设置在线客服", Desc: "用于展示在前台的未设置在线客服"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "你好", Type: TranslateTypeFrontend, Field: "hello", Value: "你好～", Desc: "用于展示在前台的好你好"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "备注", Type: TranslateTypeFrontend, Field: "remark", Value: "备注", Desc: "用于展示在前台的备注"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "活跃", Type: TranslateTypeFrontend, Field: "activated", Value: "活跃", Desc: "用于展示在前台的已激活"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "冻结", Type: TranslateTypeFrontend, Field: "frozen", Value: "冻结", Desc: "用于展示在前台的已冻结"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "立即绑定", Type: TranslateTypeFrontend, Field: "bindNow", Value: "立即绑定", Desc: "用于展示在前台的绑定按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已绑定", Type: TranslateTypeFrontend, Field: "binded", Value: "已绑定", Desc: "用于展示在前台的已绑定状态"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "总资产估值", Type: TranslateTypeFrontend, Field: "totalAssets", Value: "总资产估值", Desc: "用于展示在前台的总资产估值"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "可用余额", Type: TranslateTypeFrontend, Field: "availableBalance", Value: "可用余额", Desc: "用于展示在前台的可用余额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "冻结余额", Type: TranslateTypeFrontend, Field: "frozenBalance", Value: "冻结余额", Desc: "用于展示在前台的冻结余额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包账户", Type: TranslateTypeFrontend, Field: "walletAccount", Value: "账户", Desc: "用于展示在前台的账户余额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包资产", Type: TranslateTypeFrontend, Field: "walletAssets", Value: "资产", Desc: "用于展示在前台的资产"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "银行信息", Type: TranslateTypeFrontend, Field: "bankInfo", Value: "银行信息", Desc: "用于展示在前台的银行信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包信息", Type: TranslateTypeFrontend, Field: "walletInfo", Value: "钱包信息", Desc: "用于展示在前台的钱包信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "今日盈亏", Type: TranslateTypeFrontend, Field: "todayProfit", Value: "今日盈亏", Desc: "用于展示在前台的今日盈亏"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "货币符号", Type: TranslateTypeFrontend, Field: "currencySymbol", Value: "$", Desc: "用于展示在前台的货币符号"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "注册时间", Type: TranslateTypeFrontend, Field: "signupTime", Value: "注册时间", Desc: "用于展示在前台的注册时间"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "登录时间", Type: TranslateTypeFrontend, Field: "lastLoginTime", Value: "登录时间", Desc: "用于展示在前台的登录时间"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "登录IP", Type: TranslateTypeFrontend, Field: "lastLoginIP", Value: "登录IP", Desc: "用于展示在前台的登录IP"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "时间", Type: TranslateTypeFrontend, Field: "time", Value: "时间", Desc: "用于展示在前台的时间"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "金额", Type: TranslateTypeFrontend, Field: "amount", Value: "金额", Desc: "用于展示在前台的金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "余额", Type: TranslateTypeFrontend, Field: "balance", Value: "余额", Desc: "用于展示在前台的余额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "总金额", Type: TranslateTypeFrontend, Field: "totalAmount", Value: "总金额", Desc: "用于展示在前台的总金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "手续费", Type: TranslateTypeFrontend, Field: "fee", Value: "手续费", Desc: "用于展示在前台的手续费"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "汇率", Type: TranslateTypeFrontend, Field: "rate", Value: "汇率", Desc: "用于展示在前台的汇率"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "到期时间", Type: TranslateTypeFrontend, Field: "expiredTime", Value: "到期时间", Desc: "用于展示在前台的到期时间"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "信用分", Type: TranslateTypeFrontend, Field: "creditScore", Value: "信用分", Desc: "用于展示在前台的信用分"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "普通会员", Type: TranslateTypeFrontend, Field: "defaultMember", Value: "普通会员", Desc: "用于展示在前台的普通会员"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "状态", Type: TranslateTypeFrontend, Field: "status", Value: "状态", Desc: "用于展示在前台的状态"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未认证", Type: TranslateTypeFrontend, Field: "unverified", Value: "未认证", Desc: "用于展示在前台的未认证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已认证", Type: TranslateTypeFrontend, Field: "verified", Value: "已认证", Desc: "用于展示在前台的已认证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "认证失败", Type: TranslateTypeFrontend, Field: "verifyFailed", Value: "认证失败", Desc: "用于展示在前台的认证失败提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "安全退出", Type: TranslateTypeFrontend, Field: "logout", Value: "安全退出", Desc: "用于展示在前台的安全退出按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "暂无记录", Type: TranslateTypeFrontend, Field: "noRecord", Value: "暂无记录", Desc: "用于展示在前台的暂无记录"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "启用", Type: TranslateTypeFrontend, Field: "enabled", Value: "启用", Desc: "用于展示在前台的启用按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "编辑", Type: TranslateTypeFrontend, Field: "edit", Value: "编辑", Desc: "用于展示在前台的编辑按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "修改", Type: TranslateTypeFrontend, Field: "update", Value: "更新", Desc: "用于展示在前台的更新按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "删除", Type: TranslateTypeFrontend, Field: "delete", Value: "删除", Desc: "用于展示在前台的删除按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确定删除吗", Type: TranslateTypeFrontend, Field: "confirmDelete", Value: "确定要删除这条记录吗?", Desc: "用于展示在前台的确定删除提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "新增", Type: TranslateTypeFrontend, Field: "add", Value: "新增", Desc: "用于展示在前台的新增按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "取消", Type: TranslateTypeFrontend, Field: "cancel", Value: "取消", Desc: "用于展示在前台的取消按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "关闭", Type: TranslateTypeFrontend, Field: "close", Value: "关闭", Desc: "用于展示在前台的关闭按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确定", Type: TranslateTypeFrontend, Field: "confirm", Value: "确定", Desc: "用于展示在前台的确定按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "保存", Type: TranslateTypeFrontend, Field: "save", Value: "保存", Desc: "用于展示在前台的保存按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提交", Type: TranslateTypeFrontend, Field: "submit", Value: "提交", Desc: "用于展示在前台的提交按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "重新提交", Type: TranslateTypeFrontend, Field: "resubmit", Value: "重新提交", Desc: "用于展示在前台的重新提交按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "继续", Type: TranslateTypeFrontend, Field: "continue", Value: "继续", Desc: "用于展示在前台的继续按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "返回", Type: TranslateTypeFrontend, Field: "back", Value: "返回", Desc: "用于展示在前台的返回按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "设置", Type: TranslateTypeFrontend, Field: "setting", Value: "设置", Desc: "用于展示在前台的设置按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "查看", Type: TranslateTypeFrontend, Field: "view", Value: "查看", Desc: "用于展示在前台的查看按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "查看更多", Type: TranslateTypeFrontend, Field: "viewMore", Value: "查看更多", Desc: "用于展示在前台的查看更多按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "查看详情", Type: TranslateTypeFrontend, Field: "viewDetails", Value: "查看详情", Desc: "用于展示在前台的查看详情按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "查看全部", Type: TranslateTypeFrontend, Field: "viewAll", Value: "查看全部", Desc: "用于展示在前台的查看全部按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已读", Type: TranslateTypeFrontend, Field: "read", Value: "已读", Desc: "用于展示在前台的已读"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未读", Type: TranslateTypeFrontend, Field: "unread", Value: "未读", Desc: "用于展示在前台的未读"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "搜索", Type: TranslateTypeFrontend, Field: "search", Value: "搜索", Desc: "用于展示在前台的搜索按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "刷新", Type: TranslateTypeFrontend, Field: "refresh", Value: "刷新", Desc: "用于展示在前台的刷新按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "预览", Type: TranslateTypeFrontend, Field: "preview", Value: "预览", Desc: "用于展示在前台的预览按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "天", Type: TranslateTypeFrontend, Field: "day", Value: "天", Desc: "用于展示在前台的天"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "等待", Type: TranslateTypeFrontend, Field: "waiting", Value: "等待", Desc: "用于展示在前台的等待"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "运行中", Type: TranslateTypeFrontend, Field: "running", Value: "运行中", Desc: "用于展示在前台的运行中"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "审核", Type: TranslateTypeFrontend, Field: "audit", Value: "审核", Desc: "用于展示在前台的审核"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "完成", Type: TranslateTypeFrontend, Field: "completed", Value: "完成", Desc: "用于展示在前台的完成"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "失败", Type: TranslateTypeFrontend, Field: "failed", Value: "失败", Desc: "用于展示在前台的失败"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "拒绝", Type: TranslateTypeFrontend, Field: "rejected", Value: "拒绝", Desc: "用于展示在前台的拒绝"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "复制", Type: TranslateTypeFrontend, Field: "copy", Value: "复制", Desc: "用于展示在前台的复制按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "类型", Type: TranslateTypeFrontend, Field: "type", Value: "类型", Desc: "用于展示在前台的类型"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "资料", Type: TranslateTypeFrontend, Field: "material", Value: "资料", Desc: "用于展示在前台的资料"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "名称", Type: TranslateTypeFrontend, Field: "name", Value: "名称", Desc: "用于展示在前台的名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "操作", Type: TranslateTypeFrontend, Field: "operation", Value: "操作", Desc: "用于展示在前台的操作"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "数量", Type: TranslateTypeFrontend, Field: "quantity", Value: "数量", Desc: "用于展示在前台的数量"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "全部", Type: TranslateTypeFrontend, Field: "all", Value: "全部", Desc: "用于展示在前台的全部"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "消耗", Type: TranslateTypeFrontend, Field: "consumed", Value: "消耗", Desc: "用于展示在前台的消耗"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "获得", Type: TranslateTypeFrontend, Field: "obtained", Value: "获得", Desc: "用于展示在前台的获得"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请选择", Type: TranslateTypeFrontend, Field: "pleaseSelect", Value: "请选择", Desc: "用于展示在前台的请选择"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "筛选", Type: TranslateTypeFrontend, Field: "filter", Value: "筛选", Desc: "用于展示在前台的筛选"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请输入", Type: TranslateTypeFrontend, Field: "pleaseInput", Value: "请输入", Desc: "用于展示在前台的请输入"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "复制成功", Type: TranslateTypeFrontend, Field: "copySuccess", Value: "复制成功", Desc: "用于展示在前台的复制成功提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "复制失败", Type: TranslateTypeFrontend, Field: "copyFailed", Value: "复制失败", Desc: "用于展示在前台的复制失败提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "操作成功", Type: TranslateTypeFrontend, Field: "operationSuccess", Value: "操作成功", Desc: "用于展示在前台的操作成功提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "操作失败", Type: TranslateTypeFrontend, Field: "operationFailed", Value: "操作失败", Desc: "用于展示在前台的操作失败提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "网络异常", Type: TranslateTypeFrontend, Field: "networkError", Value: "网络异常", Desc: "用于展示在前台的网络异常提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站点SEO[标题]", Type: TranslateTypeFrontend, Field: "metaTitle", Value: "最优购物体验，尽在掌握", Desc: "用于展示在前台的首页标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站点SEO[描述]", Type: TranslateTypeFrontend, Field: "metaDescription", Value: "快速、安全的购物体验，让您轻松购物，尽享生活。加入我们，发现更多惊喜！", Desc: "用于展示在前台的首页描述"},

			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "身份证", Type: TranslateTypeFrontend, Field: "idCard", Value: "身份证", Desc: "用于展示在前台的身份证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "身份证副标题", Type: TranslateTypeFrontend, Field: "idCardSmall", Value: "请准备好您的身份证正反面照片，确保照片清晰完整", Desc: "用于展示在前台的身份证副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "护照", Type: TranslateTypeFrontend, Field: "passport", Value: "护照", Desc: "用于展示在前台的护照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "护照副标题", Type: TranslateTypeFrontend, Field: "passportSmall", Value: "请准备好您的护照照片，确保照片清晰完整", Desc: "用于展示在前台的护照副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "驾驶证", Type: TranslateTypeFrontend, Field: "driverLicense", Value: "驾驶证", Desc: "用于展示在前台的驾驶证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "驾驶证副标题", Type: TranslateTypeFrontend, Field: "driverLicenseSmall", Value: "请准备好您的驾驶证照片，确保照片清晰完整", Desc: "用于展示在前台的驾驶证副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "认证", Type: TranslateTypeFrontend, Field: "auth", Value: "认证", Desc: "用于展示在前台的认证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "选择证件类型", Type: TranslateTypeFrontend, Field: "selectIdType", Value: "选择证件类型", Desc: "用于展示在前台的证件类型选择"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件名称", Type: TranslateTypeFrontend, Field: "idName", Value: "证件名称", Desc: "用于展示在前台的证件名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件号码", Type: TranslateTypeFrontend, Field: "idNumber", Value: "证件号码", Desc: "用于展示在前台的证件号码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件正面照", Type: TranslateTypeFrontend, Field: "idPhoto1", Value: "证件正面照", Desc: "用于展示在前台的证件正面照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件反面照", Type: TranslateTypeFrontend, Field: "idPhoto2", Value: "证件反面照", Desc: "用于展示在前台的证件反面照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件手持照", Type: TranslateTypeFrontend, Field: "idPhoto3", Value: "证件手持照", Desc: "用于展示在前台的证件手持照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件地址", Type: TranslateTypeFrontend, Field: "idAddress", Value: "证件地址", Desc: "用于展示在前台的证件地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "等待审核", Type: TranslateTypeFrontend, Field: "waitingAudit", Value: "等待审核", Desc: "用于展示在前台的等待审核提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "审核通过", Type: TranslateTypeFrontend, Field: "auditPassed", Value: "审核通过", Desc: "用于展示在前台的审核通过提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "用户名", Type: TranslateTypeFrontend, Field: "username", Value: "用户名", Desc: "用于展示在前台的用户名"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "昵称", Type: TranslateTypeFrontend, Field: "nickname", Value: "昵称", Desc: "用于展示在前台的昵称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "昵称副标题", Type: TranslateTypeFrontend, Field: "nicknameSmall", Value: "设置个性化昵称让其他用户更容易记住和识别你,展现你的个性特色", Desc: "用于展示在前台的昵称副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "上传", Type: TranslateTypeFrontend, Field: "upload", Value: "上传", Desc: "用于展示在前台的上传"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "头像", Type: TranslateTypeFrontend, Field: "avatar", Value: "头像", Desc: "用于展示在前台的头像"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "头像副标题", Type: TranslateTypeFrontend, Field: "avatarSmall", Value: "上传个性化头像来展示你的形象,让你的账户更具个人特色和识别度", Desc: "用于展示在前台的头像副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "性别", Type: TranslateTypeFrontend, Field: "sex", Value: "性别", Desc: "用于展示在前台的性别"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "性别副标题", Type: TranslateTypeFrontend, Field: "sexSmall", Value: "选择你的性别信息,帮助我们为你提供更加个性化的服务体验和推荐", Desc: "用于展示在前台的性别副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "男", Type: TranslateTypeFrontend, Field: "sexMale", Value: "男", Desc: "用于展示在前台的男"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "女", Type: TranslateTypeFrontend, Field: "sexFemale", Value: "女", Desc: "用于展示在前台的女"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未知", Type: TranslateTypeFrontend, Field: "sexUnknown", Value: "未知", Desc: "用于展示在前台的未知"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "生日", Type: TranslateTypeFrontend, Field: "birthday", Value: "生日", Desc: "用于展示在前台的生日"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "生日副标题", Type: TranslateTypeFrontend, Field: "birthdaySmall", Value: "填写你的生日信息,我们将在特殊的日子为你送上祝福,并提供专属优惠", Desc: "用于展示在前台的生日副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "个性签名", Type: TranslateTypeFrontend, Field: "introduction", Value: "个性签名", Desc: "用于展示在前台的个性签名"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "个性签名副标题", Type: TranslateTypeFrontend, Field: "introductionSmall", Value: "写下你想说的话,展示你的个性和态度,让其他用户更好地了解真实的你", Desc: "用于展示在前台的个性签名副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邮箱", Type: TranslateTypeFrontend, Field: "email", Value: "邮箱", Desc: "用于展示在前台的邮箱"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "UID", Type: TranslateTypeFrontend, Field: "uid", Value: "UID", Desc: "用于展示在前台的UID"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "国家", Type: TranslateTypeFrontend, Field: "country", Value: "国家", Desc: "用于展示在前台的国家"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "电话号码", Type: TranslateTypeFrontend, Field: "telephone", Value: "电话号码", Desc: "用于展示在前台的电话号码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "密码", Type: TranslateTypeFrontend, Field: "password", Value: "密码", Desc: "用于展示在前台的密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "登录密码", Type: TranslateTypeFrontend, Field: "loginPassword", Value: "登录密码", Desc: "用于展示在前台的登录密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "登录密码副标题", Type: TranslateTypeFrontend, Field: "loginPasswordSmall", Value: "设置登录密码,确保账户安全", Desc: "用于展示在前台的登录密码副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "旧密码", Type: TranslateTypeFrontend, Field: "oldPassword", Value: "旧密码", Desc: "用于展示在前台的旧密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "新密码", Type: TranslateTypeFrontend, Field: "newPassword", Value: "新密码", Desc: "用于展示在前台的新密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认密码", Type: TranslateTypeFrontend, Field: "confirmPassword", Value: "确认密码", Desc: "用于展示在前台的确认密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付密码", Type: TranslateTypeFrontend, Field: "securityKey", Value: "支付密码", Desc: "用于展示在前台的支付密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付密码副标题", Type: TranslateTypeFrontend, Field: "securityKeySmall", Value: "设置支付密码,确保账户资金安全", Desc: "用于展示在前台的支付密码副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认支付密码", Type: TranslateTypeFrontend, Field: "confirmSecurityKey", Value: "确认支付密码", Desc: "用于展示在前台的确认支付密码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "两次密码不一致", Type: TranslateTypeFrontend, Field: "passwordNotMatch", Value: "两次密码不一致", Desc: "用于展示在前台的两次密码不一致提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "两次支付密码不一致", Type: TranslateTypeFrontend, Field: "securityKeyNotMatch", Value: "两次支付密码不一致", Desc: "用于展示在前台的两次支付密码不一致提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请码", Type: TranslateTypeFrontend, Field: "inviteCode", Value: "邀请码", Desc: "用于展示在前台的邀请码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请好友标题", Type: TranslateTypeFrontend, Field: "inviteFriendsTitle", Value: "邀请好友一起加入", Desc: "用于展示在前台的邀请好友标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请好友副标题", Type: TranslateTypeFrontend, Field: "inviteFriendsSmall", Value: "分享给好友, 让更多人体验优质服务, 邀请成功即可获得奖励, 更多奖励等你来拿", Desc: "用于展示在前台的邀请好友副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请链接", Type: TranslateTypeFrontend, Field: "inviteLink", Value: "邀请链接", Desc: "用于展示在前台的邀请链接"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "总人数", Type: TranslateTypeFrontend, Field: "totalNums", Value: "总人数", Desc: "用于展示在前台的总人数"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "总奖励", Type: TranslateTypeFrontend, Field: "totalRewards", Value: "总奖励", Desc: "用于展示在前台的总奖励"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "收益记录", Type: TranslateTypeFrontend, Field: "incomeRecord", Value: "收益记录", Desc: "用于展示在前台的收益记录"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "团队人数", Type: TranslateTypeFrontend, Field: "teamNums", Value: "团队人数", Desc: "用于展示在前台的团队人数"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "团队收益", Type: TranslateTypeFrontend, Field: "teamRewards", Value: "团队收益", Desc: "用于展示在前台的团队收益"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请规则", Type: TranslateTypeFrontend, Field: "inviteRule", Value: "邀请规则", Desc: "用于展示在前台的邀请规则"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请记录", Type: TranslateTypeFrontend, Field: "inviteRecord", Value: "邀请记录", Desc: "用于展示在前台的邀请记录"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "验证码", Type: TranslateTypeFrontend, Field: "captcha", Value: "验证码", Desc: "用于展示在前台的验证码"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "创建账户", Type: TranslateTypeFrontend, Field: "createAccount", Value: "创建账户", Desc: "用于展示在前台的创建账户按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "立即登录", Type: TranslateTypeFrontend, Field: "loginNow", Value: "立即登录", Desc: "用于展示在前台的立即登录按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邮箱登录", Type: TranslateTypeFrontend, Field: "emailLogin", Value: "邮箱登录", Desc: "用于展示在前台的邮箱登录按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "手机登录", Type: TranslateTypeFrontend, Field: "phoneLogin", Value: "手机登录", Desc: "用于展示在前台的手机登录按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邮箱注册", Type: TranslateTypeFrontend, Field: "emailRegister", Value: "邮箱注册", Desc: "用于展示在前台的邮箱注册按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "手机注册", Type: TranslateTypeFrontend, Field: "phoneRegister", Value: "手机注册", Desc: "用于展示在前台的手机注册按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "忘记密码", Type: TranslateTypeFrontend, Field: "forgetPassword", Value: "忘记密码", Desc: "用于展示在前台的忘记密码按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "同意创建协议", Type: TranslateTypeFrontend, Field: "registerAgreement", Value: `点击 "创建账户" 即表示您同意 《注册协议》`, Desc: "用于展示在前台的同意创建协议按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "欢迎来到", Type: TranslateTypeFrontend, Field: "welcome", Value: "欢迎来到", Desc: "用于展示在前台的欢迎信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "欢迎加入", Type: TranslateTypeFrontend, Field: "welcomeJoin", Value: "欢迎加入", Desc: "用于展示在前台的欢迎加入信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "充值", Type: TranslateTypeFrontend, Field: "deposit", Value: "充值", Desc: "用于展示在前台的充值按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "充值金额", Type: TranslateTypeFrontend, Field: "depositAmount", Value: "充值金额", Desc: "用于展示在前台的充值金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "充值凭证", Type: TranslateTypeFrontend, Field: "depositProof", Value: "充值凭证", Desc: "用于展示在前台的充值凭证"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提现", Type: TranslateTypeFrontend, Field: "withdraw", Value: "提现", Desc: "用于展示在前台的提现按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提现金额", Type: TranslateTypeFrontend, Field: "withdrawAmount", Value: "提现金额", Desc: "用于展示在前台的提现金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提现提示", Type: TranslateTypeFrontend, Field: "withdrawTips", Value: "需要绑定对应的提现账户, 才能进行提现", Desc: "用于展示在前台的提现提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "转账", Type: TranslateTypeFrontend, Field: "transfer", Value: "转账", Desc: "用于展示在前台的转账按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "兑换", Type: TranslateTypeFrontend, Field: "swaps", Value: "兑换", Desc: "用于展示在前台的兑换按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "安全验证", Type: TranslateTypeFrontend, Field: "securityVerification", Value: "安全验证", Desc: "用于展示在前台的安全验证标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邮箱验证", Type: TranslateTypeFrontend, Field: "emailVerification", Value: "邮箱验证", Desc: "用于展示在前台的邮箱验证标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邮箱验证提示", Type: TranslateTypeFrontend, Field: "emailVerificationTips", Value: "开启提现和安全设置, 可以增加安全强度级别", Desc: "用于展示在前台的邮箱验证提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "手机验证", Type: TranslateTypeFrontend, Field: "telephoneVerification", Value: "手机验证", Desc: "用于展示在前台的手机验证标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "手机验证提示", Type: TranslateTypeFrontend, Field: "telephoneVerificationTips", Value: "账户手机短信验证为您的账号资金带来第二重保护, 登录、提现等操作时使用。", Desc: "用于展示在前台的手机验证提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "身份认证", Type: TranslateTypeFrontend, Field: "identityVerification", Value: "身份认证", Desc: "用于展示在前台的身份认证标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "身份认证提示", Type: TranslateTypeFrontend, Field: "identityVerificationTips", Value: "身份认证的主要优势在于提升安全性，防止账号被盗用，保护用户隐私。它能有效验证用户身份，避免未经授权的访问，并支持多种验证方式，增强用户信任感，确保线上交易和信息交换的安全性。", Desc: "用于展示在前台的身份认证提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "地址认证", Type: TranslateTypeFrontend, Field: "addressVerification", Value: "地址认证", Desc: "用于展示在前台的地址认证标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "地址认证提示", Type: TranslateTypeFrontend, Field: "addressVerificationTips", Value: "地址认证的主要优势在于提升安全性，防止账号被盗用，保护用户隐私。它能有效验证用户身份，避免未经授权的访问，并支持多种验证方式，增强用户信任感，确保线上交易和信息交换的安全性。", Desc: "用于展示在前台的地址认证提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "立即认证", Type: TranslateTypeFrontend, Field: "verifyNow", Value: "立即认证", Desc: "用于展示在前台的立即认证按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "会员权益", Type: TranslateTypeFrontend, Field: "memberBenefits", Value: "会员权益", Desc: "用于展示在前台的会员权益标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "会员权益提示", Type: TranslateTypeFrontend, Field: "memberBenefitsTips", Value: "会员等级系统根据用户的活跃度、消费金额和使用时长等因素进行划分，用户随着等级提升可享受更多专属权益，如优先服务、折扣优惠、积分奖励等，旨在提升用户体验并激励长期使用和参与。", Desc: "用于展示在前台的会员权益提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "暂未获得", Type: TranslateTypeFrontend, Field: "notObtained", Value: "暂未获得", Desc: "用于展示在前台的暂未获得"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "立即购买", Type: TranslateTypeFrontend, Field: "buyNow", Value: "立即购买", Desc: "用于展示在前台的立即购买按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "购买", Type: TranslateTypeFrontend, Field: "buy", Value: "购买", Desc: "用于展示在前台的购买按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已拥有", Type: TranslateTypeFrontend, Field: "alreadyOwned", Value: "已拥有", Desc: "用于展示在前台的已拥有"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "我的资产", Type: TranslateTypeFrontend, Field: "myAssets", Value: "我的资产", Desc: "用于展示在前台的我的资产按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "充提订单", Type: TranslateTypeFrontend, Field: "depositWithdrawOrder", Value: "充提订单", Desc: "用于展示在前台的充提订单按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付方式", Type: TranslateTypeFrontend, Field: "paymentMethod", Value: "支付方式", Desc: "用于展示在前台的支付方式"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "转账方式", Type: TranslateTypeFrontend, Field: "transferMethod", Value: "转账方式", Desc: "用于展示在前台的转账方式"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "转账金额", Type: TranslateTypeFrontend, Field: "transferAmount", Value: "转账金额", Desc: "用于展示在前台的转账金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]银行名称", Type: TranslateTypeFrontend, Field: "accountBankName", Value: "银行名称", Desc: "用于提现账户的银行名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]公链名称", Type: TranslateTypeFrontend, Field: "accountTokenName", Value: "公链名称", Desc: "用于提现账户的公链名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]银行卡号", Type: TranslateTypeFrontend, Field: "accountBankCardNo", Value: "银行卡号", Desc: "用于提现账户的银行卡号"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]Token地址", Type: TranslateTypeFrontend, Field: "accountTokenAddress", Value: "Token地址", Desc: "用于提现账户的Token地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]真实姓名", Type: TranslateTypeFrontend, Field: "accountRealName", Value: "真实姓名", Desc: "用于提现账户的真实姓名"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]Token", Type: TranslateTypeFrontend, Field: "accountToken", Value: "Token", Desc: "用于提现账户的Token"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]银行代号", Type: TranslateTypeFrontend, Field: "accountBankCode", Value: "银行代号", Desc: "用于提现账户的银行代号"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]公链描述", Type: TranslateTypeFrontend, Field: "accountTokenDesc", Value: "公链描述", Desc: "用于提现账户的公链描述"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]银行地址", Type: TranslateTypeFrontend, Field: "accountBankAddress", Value: "银行地址", Desc: "用于提现账户的银行地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "[提现账户]公链标识", Type: TranslateTypeFrontend, Field: "accountTokenSymbol", Value: "公链标识", Desc: "用于提现账户的公链标识"},

			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "账户设置", Type: TranslateTypeFrontend, Field: "userSetting", Value: "账户设置", Desc: "用于展示在前台的账户设置按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "安全设置", Type: TranslateTypeFrontend, Field: "securitySetting", Value: "安全设置", Desc: "用于展示在前台的安全设置按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "密码管理", Type: TranslateTypeFrontend, Field: "passwordManagement", Value: "密码管理", Desc: "用于展示在前台的密码管理按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "认证管理", Type: TranslateTypeFrontend, Field: "verificationManagement", Value: "认证管理", Desc: "用于展示在前台的认证管理按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "邀请好友", Type: TranslateTypeFrontend, Field: "inviteFriends", Value: "邀请好友", Desc: "用于展示在前台的邀请好友按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "我的团队", Type: TranslateTypeFrontend, Field: "myTeam", Value: "我的团队", Desc: "用于展示在前台的我的团队按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "消息通知", Type: TranslateTypeFrontend, Field: "messageNotice", Value: "消息通知", Desc: "用于展示在前台的消息通知按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包总览", Type: TranslateTypeFrontend, Field: "walletOverview", Value: "钱包总览", Desc: "用于展示在前台的钱包总览按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包账单", Type: TranslateTypeFrontend, Field: "walletBill", Value: "钱包账单", Desc: "用于展示在前台的钱包账单按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站内转账", Type: TranslateTypeFrontend, Field: "walletTransfer", Value: "站内转账", Desc: "用于展示在前台的站内转账按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "站内转账副标题", Type: TranslateTypeFrontend, Field: "walletTransferSmall", Value: "即时到账｜手续费低｜安全性高", Desc: "用于展示在前台的站内转账副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "资产兑换", Type: TranslateTypeFrontend, Field: "walletSwaps", Value: "资产兑换", Desc: "用于展示在前台的资产兑换按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "资产兑换副标题", Type: TranslateTypeFrontend, Field: "walletSwapsSmall", Value: "即时到账｜手续费低｜安全性高", Desc: "用于展示在前台的资产兑换副标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提现账户", Type: TranslateTypeFrontend, Field: "walletWithdrawAccount", Value: "提现账户", Desc: "用于展示在前台的提现账户按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包充值", Type: TranslateTypeFrontend, Field: "walletDeposit", Value: "钱包充值", Desc: "用于展示在前台的钱包充值按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "钱包提现", Type: TranslateTypeFrontend, Field: "walletWithdraw", Value: "钱包提现", Desc: "用于展示在前台的钱包提现按钮"},

			//商城
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺名称", Type: TranslateTypeFrontend, Field: "storeName", Value: "店铺名称", Desc: "用于展示在前台的店铺名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "证件号", Type: TranslateTypeFrontend, Field: "storeNumber", Value: "证件号", Desc: "用于展示在前台的证件号"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺logo", Type: TranslateTypeFrontend, Field: "storeLoge", Value: "店铺logo", Desc: "用于展示在前台的店铺logo"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "营业执照", Type: TranslateTypeFrontend, Field: "storePhoto1", Value: "营业执照", Desc: "用于展示在前台的营业执照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "营业执照副本", Type: TranslateTypeFrontend, Field: "storePhoto2", Value: "营业执照副本", Desc: "用于展示在前台的营业执照副本"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "营业执照手持照", Type: TranslateTypeFrontend, Field: "storePhoto3", Value: "营业执照手持照", Desc: "用于展示在前台的营业执照手持照"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "分类", Type: TranslateTypeFrontend, Field: "menuCategory", Value: "分类", Desc: "用于展示在前台的分类"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "购物车", Type: TranslateTypeFrontend, Field: "menuCard", Value: "购物车", Desc: "用于展示在前台的购物车"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "消息", Type: TranslateTypeFrontend, Field: "menuMessage", Value: "消息", Desc: "用于展示在前台的消息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺", Type: TranslateTypeFrontend, Field: "menuShop", Value: "店铺", Desc: "用于展示在前台的店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商家入驻", Type: TranslateTypeFrontend, Field: "menuMerchantsSettleIn", Value: "商家入驻", Desc: "用于展示在前台的商家入驻"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商家服务", Type: TranslateTypeFrontend, Field: "menuStoreService", Value: "商家服务", Desc: "用于展示在前台的商家服务"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "我的订单", Type: TranslateTypeFrontend, Field: "menuMyOrder", Value: "我的订单", Desc: "用于展示在前台的我的订单"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品收藏", Type: TranslateTypeFrontend, Field: "menuFollowProduct", Value: "商品收藏", Desc: "用于展示在前台的商品收藏"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺关注", Type: TranslateTypeFrontend, Field: "menuFollowStore", Value: "店铺关注", Desc: "用于展示在前台的店铺关注"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "浏览记录", Type: TranslateTypeFrontend, Field: "menuBrowsingProduct", Value: "浏览记录", Desc: "用于展示在前台的浏览记录"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "收货地址", Type: TranslateTypeFrontend, Field: "menuShippingAddress", Value: "收货地址", Desc: "用于展示在前台的收货地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待评价", Type: TranslateTypeFrontend, Field: "menuPendingReview", Value: "待评价", Desc: "用于展示在前台的待评价"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "售后", Type: TranslateTypeFrontend, Field: "menuPendingRefund", Value: "售后", Desc: "用于展示在前台的售后"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待付款", Type: TranslateTypeFrontend, Field: "menuPendingPayment", Value: "待付款", Desc: "用于展示在前台的待付款"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待发货", Type: TranslateTypeFrontend, Field: "menuAwaitingShipment", Value: "待发货", Desc: "用于展示在前台的待发货"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待收货", Type: TranslateTypeFrontend, Field: "menuAwaitingReceipt", Value: "待收货", Desc: "用于展示在前台的待收货"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待支付", Type: TranslateTypeFrontend, Field: "pendingPayment", Value: "待支付", Desc: "用于展示在前台的待支付"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待退货", Type: TranslateTypeFrontend, Field: "toBeReturned", Value: "待退货", Desc: "用于展示在前台的待退货"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待退款", Type: TranslateTypeFrontend, Field: "pendingRefund", Value: "待退款", Desc: "用于展示在前台的待退货"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "客服", Type: TranslateTypeFrontend, Field: "chats", Value: "客服", Desc: "用于展示在前台的客服"},

			// 购物车
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "管理", Type: TranslateTypeFrontend, Field: "manage", Value: "管理", Desc: "用于展示在前台的管理按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "收藏", Type: TranslateTypeFrontend, Field: "collect", Value: "收藏", Desc: "用于展示在前台的收藏"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "猜你喜欢", Type: TranslateTypeFrontend, Field: "lovely", Value: "猜你喜欢", Desc: "用于展示在前台的猜你喜欢"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品", Type: TranslateTypeFrontend, Field: "commodity", Value: "商品", Desc: "用于展示在前台的商品"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "结算", Type: TranslateTypeFrontend, Field: "settlement", Value: "结算", Desc: "用于展示在前台的结算"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已选", Type: TranslateTypeFrontend, Field: "selected", Value: "已选", Desc: "用于展示在前台的已选提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "合计", Type: TranslateTypeFrontend, Field: "total", Value: "合计", Desc: "用于展示在前台的合计提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "件", Type: TranslateTypeFrontend, Field: "items", Value: "件", Desc: "用于展示在前台的合计提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "选择规模", Type: TranslateTypeFrontend, Field: "size", Value: "选择规模", Desc: "用于展示在前台的合计提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请选择商品", Type: TranslateTypeFrontend, Field: "selectOneItem", Value: "请选择商品", Desc: "用于展示在前台的合计提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "删除订单", Type: TranslateTypeFrontend, Field: "deleteOrder", Value: "删除订单", Desc: "用于展示在前台的合计提示"},

			//评价中心
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已评价", Type: TranslateTypeFrontend, Field: "evaluated", Value: "已评价", Desc: "用于展示在前台的合计提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评论", Type: TranslateTypeFrontend, Field: "comment", Value: "评论", Desc: "用于展示在前台的评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评价", Type: TranslateTypeFrontend, Field: "evaluate", Value: "评价", Desc: "用于展示在前台的评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评价中心", Type: TranslateTypeFrontend, Field: "assessmentCenter", Value: "评价中心", Desc: "用于展示在前台的评价中心提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "发表评价", Type: TranslateTypeFrontend, Field: "postEvaluation", Value: "发表评价", Desc: "用于展示在前台的发布评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请输入评价", Type: TranslateTypeFrontend, Field: "enterYourReview", Value: "请输入评价内容", Desc: "用于展示在前台的发布评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评价成功", Type: TranslateTypeFrontend, Field: "commentSussessDesc", Value: "评价成功,您可前往评价中心查看评价", Desc: "用于展示在前台的评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "返回首页", Type: TranslateTypeFrontend, Field: "returnToHomepage", Value: "返回首页", Desc: "用于展示在前台的发布评价提示"},

			//产品详细
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "申请退款", Type: TranslateTypeFrontend, Field: "requestRefund", Value: "申请退款", Desc: "用于展示在前台的申请退款提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "追加评价", Type: TranslateTypeFrontend, Field: "additionalEvaluation", Value: "追加评价", Desc: "用于展示在前台的追加评价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品总价", Type: TranslateTypeFrontend, Field: "totalPrice", Value: "商品总价", Desc: "用于展示在前台的商品总价提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "优惠金额", Type: TranslateTypeFrontend, Field: "discountAmount", Value: "优惠金额", Desc: "用于展示在前台的优惠提示"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "实付款", Type: TranslateTypeFrontend, Field: "disbursements", Value: "实付款", Desc: "用于展示在前台的实付款"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单编号", Type: TranslateTypeFrontend, Field: "orderId", Value: "订单编号", Desc: "用于展示在前台的订单编号"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "创建时间", Type: TranslateTypeFrontend, Field: "creationTime", Value: "创建时间", Desc: "用于展示在前台的创建时间"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单详细", Type: TranslateTypeFrontend, Field: "orderDetails", Value: "订单详细", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "取消订单", Type: TranslateTypeFrontend, Field: "cancelOrder", Value: "取消订单", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认收货", Type: TranslateTypeFrontend, Field: "confirmReceipt", Value: "确认收货", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "交易已关闭", Type: TranslateTypeFrontend, Field: "transactionClosed", Value: "交易已关闭", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "交易关闭", Type: TranslateTypeFrontend, Field: "transactionClose", Value: "交易关闭", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "发货超时", Type: TranslateTypeFrontend, Field: "deliveryTimeout", Value: "发货超时", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "交易完成", Type: TranslateTypeFrontend, Field: "transactionCompletion", Value: "交易完成", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "搜索订单", Type: TranslateTypeFrontend, Field: "searchOrders", Value: "搜索订单", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单超时", Type: TranslateTypeFrontend, Field: "orderTimeout", Value: "订单超时", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "应付款", Type: TranslateTypeFrontend, Field: "accountsPayable", Value: "应付款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单号", Type: TranslateTypeFrontend, Field: "orderNo", Value: "订单号", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已完成", Type: TranslateTypeFrontend, Field: "Finished", Value: "已完成", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "添加收货地址", Type: TranslateTypeFrontend, Field: "addAddress", Value: "添加收货地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "价格明细", Type: TranslateTypeFrontend, Field: "priceBreakdown", Value: "价格明细", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "共", Type: TranslateTypeFrontend, Field: "altogether", Value: "共", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提交订单", Type: TranslateTypeFrontend, Field: "submitOrder", Value: "提交订单", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认收货地址", Type: TranslateTypeFrontend, Field: "confirmAddress", Value: "确认收货地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认订单信息", Type: TranslateTypeFrontend, Field: "confirmInformation", Value: "确认订单信息", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付密码不能为空", Type: TranslateTypeFrontend, Field: "paymentPasswordEmpty", Value: "支付密码不能为空", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退款原因", Type: TranslateTypeFrontend, Field: "refundReason", Value: "退款原因", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "您还可以输入", Type: TranslateTypeFrontend, Field: "canAlsoEnter", Value: "您还可以输入", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "字", Type: TranslateTypeFrontend, Field: "word", Value: "字", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "上传图片", Type: TranslateTypeFrontend, Field: "uploadPictures", Value: "上传图片", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "提交申请", Type: TranslateTypeFrontend, Field: "submitApplication", Value: "提交申请", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请输入退款原因或上传图片", Type: TranslateTypeFrontend, Field: "uploadImage", Value: "请输入退款原因或上传图片", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退款中", Type: TranslateTypeFrontend, Field: "refunding", Value: "退款中", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退款成功", Type: TranslateTypeFrontend, Field: "refundSuccessful", Value: "退款成功", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退款金额", Type: TranslateTypeFrontend, Field: "refundAmount", Value: "退款金额", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "申请时间", Type: TranslateTypeFrontend, Field: "applicationTime", Value: "申请时间", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "取消申请", Type: TranslateTypeFrontend, Field: "cancelledWithdraw", Value: "取消申请", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "售后详情", Type: TranslateTypeFrontend, Field: "afterSalesDetails", Value: "售后详情", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "取消退款", Type: TranslateTypeFrontend, Field: "cancelRefund", Value: "取消退款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已退款", Type: TranslateTypeFrontend, Field: "refunded", Value: "已退款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "售后中心", Type: TranslateTypeFrontend, Field: "afterSalesCenter", Value: "售后中心", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退款申请成功", Type: TranslateTypeFrontend, Field: "refundRequestSuccessful", Value: "退款申请成功,您可前往售后中心查看退款进度", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "申请成功", Type: TranslateTypeFrontend, Field: "applicationSuccessful", Value: "申请成功", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品详情", Type: TranslateTypeFrontend, Field: "shopDetails", Value: "商品详情", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "产品详情", Type: TranslateTypeFrontend, Field: "productDetails", Value: "产品详情", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "销量", Type: TranslateTypeFrontend, Field: "salesVolume", Value: "销量", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "选择规格", Type: TranslateTypeFrontend, Field: "selectSpecifications", Value: "选择规格", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请选择商品规格", Type: TranslateTypeFrontend, Field: "selectProductSize", Value: "请选择商品规格", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "综合", Type: TranslateTypeFrontend, Field: "comprehensive", Value: "综合", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "价格", Type: TranslateTypeFrontend, Field: "price", Value: "价格", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "新品", Type: TranslateTypeFrontend, Field: "newProduct", Value: "新品", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "结算支付成功", Type: TranslateTypeFrontend, Field: "paymentSuccessful", Value: "支付成功,您可前往订单中心查看订单", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "支付成功", Type: TranslateTypeFrontend, Field: "paySuccessful", Value: "支付成功", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单中心", Type: TranslateTypeFrontend, Field: "orderCenter", Value: "订单中心", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评价管理", Type: TranslateTypeFrontend, Field: "evaluationManagement", Value: "评价管理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "搜索订单号", Type: TranslateTypeFrontend, Field: "searchOrderId", Value: "搜索订单号", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "批发价", Type: TranslateTypeFrontend, Field: "wholesalePrice", Value: "批发价", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "利润", Type: TranslateTypeFrontend, Field: "profit", Value: "利润", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品金额", Type: TranslateTypeFrontend, Field: "productAmount", Value: "商品金额", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "买家姓名", Type: TranslateTypeFrontend, Field: "buyerName", Value: "买家姓名", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认发货", Type: TranslateTypeFrontend, Field: "confirmShipment", Value: "确认发货", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已发货", Type: TranslateTypeFrontend, Field: "shipped", Value: "已发货", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "图片", Type: TranslateTypeFrontend, Field: "picture", Value: "图片", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺信息", Type: TranslateTypeFrontend, Field: "storeInformation", Value: "店铺信息", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "库存", Type: TranslateTypeFrontend, Field: "inventory", Value: "库存", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "零售价", Type: TranslateTypeFrontend, Field: "retailPrice", Value: "零售价", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "上架", Type: TranslateTypeFrontend, Field: "listing", Value: "上架", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "全部上架", Type: TranslateTypeFrontend, Field: "allListed", Value: "全部上架", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品名称", Type: TranslateTypeFrontend, Field: "productName", Value: "商品名称", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "下架", Type: TranslateTypeFrontend, Field: "removeProducts", Value: "下架", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已下架", Type: TranslateTypeFrontend, Field: "removed", Value: "已下架", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "全部下架", Type: TranslateTypeFrontend, Field: "removeAll", Value: "全部下架", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "出售中", Type: TranslateTypeFrontend, Field: "onSale", Value: "出售中", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "全部商品", Type: TranslateTypeFrontend, Field: "allProducts", Value: "全部商品", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "财务管理", Type: TranslateTypeFrontend, Field: "financialManagement", Value: "财务管理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品描述", Type: TranslateTypeFrontend, Field: "productDescription", Value: "商品描述", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "拒绝退款", Type: TranslateTypeFrontend, Field: "refuseRefund", Value: "拒绝退款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "拒绝原因", Type: TranslateTypeFrontend, Field: "rejectionReason", Value: "拒绝原因", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "同意", Type: TranslateTypeFrontend, Field: "agree", Value: "同意", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "同意退款", Type: TranslateTypeFrontend, Field: "agreeRefund", Value: "同意退款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "确认退款", Type: TranslateTypeFrontend, Field: "confirmRefund", Value: "确认退款", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "搜索售后", Type: TranslateTypeFrontend, Field: "searchAfterSales", Value: "搜索售后", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已拒绝", Type: TranslateTypeFrontend, Field: "hasBeenRejected", Value: "已拒绝", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "售后管理", Type: TranslateTypeFrontend, Field: "afterSalesManagement", Value: "售后管理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待处理", Type: TranslateTypeFrontend, Field: "toBeProcessed", Value: "待处理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "已处理", Type: TranslateTypeFrontend, Field: "isProcessed", Value: "已处理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "规格", Type: TranslateTypeFrontend, Field: "specifications", Value: "规格", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品价格", Type: TranslateTypeFrontend, Field: "productPrice", Value: "商品价格", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "退出管理", Type: TranslateTypeFrontend, Field: "exitManagement", Value: "退出管理", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "新增地址", Type: TranslateTypeFrontend, Field: "addNewAddress", Value: "新增地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "默认", Type: TranslateTypeFrontend, Field: "default", Value: "默认", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "默认地址", Type: TranslateTypeFrontend, Field: "defaultAddress", Value: "默认地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "城市", Type: TranslateTypeFrontend, Field: "city", Value: "城市", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "地址", Type: TranslateTypeFrontend, Field: "address", Value: "地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "保存地址", Type: TranslateTypeFrontend, Field: "saveAddress", Value: "保存地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "删除地址", Type: TranslateTypeFrontend, Field: "deleteAddress", Value: "删除地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "是否删除该地址", Type: TranslateTypeFrontend, Field: "confirmDeleteAddress", Value: "是否删除该地址", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "访问量", Type: TranslateTypeFrontend, Field: "pageViews", Value: "访问量", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "取消关注", Type: TranslateTypeFrontend, Field: "unfollow", Value: "取消关注", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "关注", Type: TranslateTypeFrontend, Field: "follow", Value: "关注", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "访问店铺", Type: TranslateTypeFrontend, Field: "visitStore", Value: "访问店铺", Desc: "用于展示在前台的提示/按钮"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "加入", Type: TranslateTypeFrontend, Field: "join", Value: "加入", Desc: "用于展示在前台的加入"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "综合评分", Type: TranslateTypeFrontend, Field: "overallScore", Value: "综合评分", Desc: "用于展示在前台的综合评分"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "商品管理", Type: TranslateTypeFrontend, Field: "productManage", Value: "商品管理", Desc: "用于展示在前台的商品管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单管理", Type: TranslateTypeFrontend, Field: "orderManage", Value: "订单管理", Desc: "用于展示在前台的订单管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "售后管理", Type: TranslateTypeFrontend, Field: "refundManage", Value: "售后管理", Desc: "用于展示在前台的售后管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "财务管理", Type: TranslateTypeFrontend, Field: "financialManage", Value: "财务管理", Desc: "用于展示在前台的财务管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺管理", Type: TranslateTypeFrontend, Field: "storeManage", Value: "店铺管理", Desc: "用于展示在前台的店铺管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "评论管理", Type: TranslateTypeFrontend, Field: "commentManage", Value: "评论管理", Desc: "用于展示在前台的评论管理"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "收益", Type: TranslateTypeFrontend, Field: "earnings", Value: "收益", Desc: "用于展示在前台的收益"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单", Type: TranslateTypeFrontend, Field: "order", Value: "订单", Desc: "用于展示在前台的订单"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "总销售额", Type: TranslateTypeFrontend, Field: "totalSales", Value: "总销售额", Desc: "用于展示在前台的总销售额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "今日销售额", Type: TranslateTypeFrontend, Field: "todaySalesMoney", Value: "今日销售额", Desc: "用于展示在前台的今日销售额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "待处理金额", Type: TranslateTypeFrontend, Field: "pendingAmount", Value: "待处理金额", Desc: "用于展示在前台的待处理金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "联系方式", Type: TranslateTypeFrontend, Field: "contactDetails", Value: "联系方式", Desc: "用于展示在前台的联系方式"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺地址", Type: TranslateTypeFrontend, Field: "storeAddress", Value: "店铺地址", Desc: "用于展示在前台的店铺地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "关键字", Type: TranslateTypeFrontend, Field: "storeKeywords", Value: "关键字", Desc: "用于展示在前台的关键字"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺描述", Type: TranslateTypeFrontend, Field: "storesDesc", Value: "店铺描述", Desc: "用于展示在前台的店铺描述"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺logo", Type: TranslateTypeFrontend, Field: "storeLogo", Value: "店铺logo", Desc: "用于展示在前台的店铺logo"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "订单信息", Type: TranslateTypeFrontend, Field: "orderInfo", Value: "订单信息", Desc: "用于展示在前台的订单信息"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "常用应用", Type: TranslateTypeFrontend, Field: "commonApplications", Value: "常用应用", Desc: "用于展示在前台的常用应用"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "店铺数据", Type: TranslateTypeFrontend, Field: "storeData", Value: "店铺数据", Desc: "用于展示在前台的店铺数据"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "请选择收货地址", Type: TranslateTypeFrontend, Field: "pleaseSelectAddress", Value: "请选择收货地址", Desc: "用于展示在前台的请选择收货地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未发货金额", Type: TranslateTypeFrontend, Field: "unshippedAmount", Value: "未发货金额", Desc: "用于展示在前台的未发货金额"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "首页", Type: TranslateTypeFrontend, Field: "home", Value: "首页", Desc: "用于展示在前台的首页"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "免运费", Type: TranslateTypeFrontend, Field: "shippingTitle", Value: "免运费", Desc: "用于展示在前台的免运费"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "所有订单免费送货", Type: TranslateTypeFrontend, Field: "shippingDesc", Value: "所有订单免费送货", Desc: "用于展示在前台的所有订单免费送货"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "资金担保", Type: TranslateTypeFrontend, Field: "guaranteeTitle", Value: "资金担保", Desc: "用于展示在前台的资金担保"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "30天退款", Type: TranslateTypeFrontend, Field: "guaranteeDesc", Value: "30天退款", Desc: "用于展示在前台的30天退款"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "24/7支持", Type: TranslateTypeFrontend, Field: "supportTitle", Value: "24/7支持", Desc: "用于展示在前台的24/7支持"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "友好的24/7支持", Type: TranslateTypeFrontend, Field: "supportDesc", Value: "友好的24/7支持", Desc: "用于展示在前台的友好的24/7支持"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "安全支付", Type: TranslateTypeFrontend, Field: "secureTitle", Value: "安全支付", Desc: "用于展示在前台的安全支付"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "接受所有卡", Type: TranslateTypeFrontend, Field: "secureDesc", Value: "接受所有卡", Desc: "用于展示在前台的接受所有卡"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "畅销产品", Type: TranslateTypeFrontend, Field: "bestProducts", Value: "畅销产品", Desc: "用于展示在前台的畅销产品"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "畅销产品描述", Type: TranslateTypeFrontend, Field: "bestProductsDesc", Value: "精选热销商品，尽享购物乐趣，快来看看大家都在抢购的必备好物！", Desc: "用于展示在前台的畅销产品描述"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "入驻协议", Type: TranslateTypeFrontend, Field: "settledNotice", Value: "入驻协议", Desc: "用于显示在前台的入驻须知"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "入驻协议内容", Type: TranslateTypeFrontend, Field: "settledNoticeValue", Value: "<pre><li><p><strong>协议双方</strong>：</p><ul><li>明确协议的双方是平台（通常是电商平台、服务平台等）与商家，列出商家和平台的基本信息。</li></ul></li><li><p><strong>入驻要求与资格</strong>：</p><ul><li>说明商家入驻平台的资格条件，如注册公司、经营范围、合法资质、商品合法性等。</li><li>商家需要提交的相关材料、审核流程等。</li></ul></li><li><p><strong>平台服务内容</strong>：</p><ul><li>列明平台为商家提供的服务，如线上展示、交易、物流支持、支付结算等。</li><li>描述平台提供的工具、数据分析等服务。</li></ul></li><li><p><strong>商家责任与义务</strong>：</p><ul><li>商家承诺提供真实、合法的商品或服务，不侵犯知识产权。</li><li>商家应按照平台要求进行商品发布、维护店铺、履行售后服务等。</li><li>商家需确保遵守平台的经营规则，包括商品描述规范、客户服务标准等。</li></ul></li><li><p><strong>平台责任与义务</strong>：</p><ul><li>平台承诺为商家提供良好的交易环境，保障商家的合法权益。</li><li>平台需要及时结算商家的销售收入，并提供技术支持。</li></ul></li><li><p><strong>费用及结算</strong>：</p><ul><li>商家应支付的费用（如入驻费用、技术服务费、平台抽佣等）。</li><li>结算方式、周期及支付方式的说明。</li></ul></li><li><p><strong>知识产权</strong>：</p><ul><li>确保商家不侵犯他人的知识产权，并授权平台使用相关商标、品牌、图片等。</li></ul></li><li><p><strong>数据保护与隐私</strong>：</p><ul><li>平台如何处理商家的个人信息、商业数据等，商家如何保护消费者数据等。</li></ul></li><li><p><strong>违约责任与争议解决</strong>：</p><ul><li>违约责任的条款，商家或平台违反协议的后果。</li><li>争议的解决方式，通常包括协商、调解、仲裁或诉讼。</li></ul></li><li><p><strong>协议的有效期与终止</strong>：</p><ul><li>协议的有效期限，终止条件（如商家撤销入驻、平台违规等）。</li></ul></li><li><p><strong>其他条款</strong>：</p><ul><li>适用法律、不可抗力、修改协议的程序等。</li></ul></li></pre>", Desc: "用于显示在前台的入驻须知内容"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "收件人名称", Type: TranslateTypeFrontend, Field: "recipientName", Value: "收件人名称", Desc: "用于展示在前台的收件人名称"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "详细地址", Type: TranslateTypeFrontend, Field: "detailedAddress", Value: "详细地址", Desc: "用于展示在前台的详细地址"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "搜索记录", Type: TranslateTypeFrontend, Field: "searchRecord", Value: "搜索记录", Desc: "用于展示在前台的搜索记录"},

			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到订单", Type: TranslateTypeSystem, Field: "notFoundOrder", Value: "未找到订单", Desc: "用于展示在前台的未找到订单"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到店铺订单", Type: TranslateTypeSystem, Field: "notFoundStoreOrder", Value: "未找到店铺订单", Desc: "用于展示在前台的未找到订单"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到店铺", Type: TranslateTypeSystem, Field: "notFoundStore", Value: "未找到店铺", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到售后订单", Type: TranslateTypeSystem, Field: "notFoundRefund", Value: "未找到售后订单", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "更新失败", Type: TranslateTypeSystem, Field: "updateFail", Value: "更新失败", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到产品", Type: TranslateTypeSystem, Field: "notFoundProduct", Value: "未找到产品", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "没有数据", Type: TranslateTypeSystem, Field: "notData", Value: "没有数据", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未发货不能下架", Type: TranslateTypeSystem, Field: "notShippedRemoved", Value: "未发货不能下架", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "未找到商品规格", Type: TranslateTypeSystem, Field: "notFoundSku", Value: "未找到商品规格", Desc: "用于展示在前台的未找到店铺"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "轮播图1标题", Type: TranslateTypeSystem, Field: "banner1Title", Value: "会员专享，尊享福利", Desc: "用于展示在前台的轮播图1标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "轮播图1描述", Type: TranslateTypeSystem, Field: "banner1Desc", Value: "超值折扣、积分返现、优先抢购新品等你来！立即注册，享受会员专属福利，让购物更省心！", Desc: "用于展示在前台的轮播图1描述"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "轮播图2标题", Type: TranslateTypeSystem, Field: "banner2Title", Value: "新品上架，快来看看！", Desc: "用于展示在前台的轮播图2标题"},
			{AdminID: SuperAdminID, Lang: "zh-CN", Name: "轮播图2描述", Type: TranslateTypeSystem, Field: "banner2Desc", Value: "全新商品上线，潮流新品等你来选！从 时尚服饰 到 智能科技，应有尽有！赶紧抢先体验新品，点击购买，让你引领潮流", Desc: "用于展示在前台的轮播图2描述"},
		}

		// 添加钱包账单类型翻译
		for _, v := range BillTypeNames {
			label, _ := v["label"].(string)
			translates = append(translates, Translate{
				AdminID: SuperAdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("钱包账单类型[%s]", label),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(WalletBillTypePrefix+"%d", v.Value),
				Value:   label,
			})
		}

		//	添加菜单翻译
		menus := initSystemsMenus()
		existingFields := make(map[string]bool)
		for _, menu := range menus {
			field := fmt.Sprintf(MenuNameTranslatePrefix, menu.Data.Label)
			if !existingFields[field] {
				translates = append(translates, Translate{
					AdminID: menu.AdminID,
					Lang:    DefaultLang,
					Name:    fmt.Sprintf("菜单[%s]", menu.Name),
					Type:    TranslateTypeSystem,
					Field:   field,
					Value:   menu.Name,
				})
				existingFields[field] = true
			}
		}

		// 添加等级翻译
		levels := initSystemLevels()
		for _, level := range levels {
			translates = append(translates, Translate{
				AdminID: level.AdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("等级[%s]", level.Name),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(LevelNameTranslatePrefix, level.Symbol),
				Value:   level.Name,
			})

			translates = append(translates, Translate{
				AdminID: level.AdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("等级[%s详情]", level.Name),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(LevelDescTranslatePrefix, level.Symbol),
				Value:   level.Desc,
			})
		}

		// 添加商品分类翻译
		categors := initCategorys()
		for _, v := range categors {
			translates = append(translates, Translate{
				AdminID: v.AdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("商品分类[%s]", v.Name),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(CategoryTranslatePrefix, v.Symbol),
				Value:   v.Name,
			})
		}

		// 添加文章翻译
		articles := initSystemArticles()
		for _, v := range articles {
			translates = append(translates, Translate{
				AdminID: v.AdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("文章[%s]标题", v.Title),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(ArticleTitleTranslatePrefix, v.Symbol),
				Value:   v.Title,
			})

			translates = append(translates, Translate{
				AdminID: v.AdminID,
				Lang:    DefaultLang,
				Name:    fmt.Sprintf("文章[%s]内容", v.Title),
				Type:    TranslateTypeSystem,
				Field:   fmt.Sprintf(ArticleContentTranslatePrefix, v.Symbol),
				Value:   v.Content,
			})
		}

		if err := db.CreateInBatches(translates, len(translates)).Error; err != nil {
			return err
		}
	}

	return nil
}
