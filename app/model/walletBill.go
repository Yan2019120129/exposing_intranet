package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"time"
)

// WalletBillTypePrefix 钱包账单类型前缀
const WalletBillTypePrefix = "walletBillType"

const (
	// 信用（传入）账单类型
	BillTypeDeposit            int8 = 1  // 充值
	BillTypeSystemReward       int8 = 2  // 系统奖励
	BillTypeRegisterReward     int8 = 3  // 注册奖励
	BillTypeInviteReward       int8 = 4  // 邀请奖励
	BillTypeDistributionReward int8 = 5  // 分销奖励
	BillTypeBalanceUnfreeze    int8 = 6  // 余额解冻
	BillTypeSystemAddition     int8 = 7  // 系统加款
	BillTypeWithdrawalReject   int8 = 8  // 提现拒绝
	BillTypeProductRefund      int8 = 9  // 产品退款
	BillTypeProductEarnings    int8 = 10 // 产品收益
	BillTypeTransferReceive    int8 = 11 // 转账接收
	BillTypeSwapsReceive       int8 = 12 // 闪兑接收

	// 借记（支出）账单类型
	BillTypeWithdrawal         int8 = -1 // 提现
	BillTypeSystemDeduction    int8 = -2 // 系统扣款
	BillTypeProductPurchase    int8 = -3 // 购买产品
	BillTypeMembershipPurchase int8 = -4 // 购买会员
	BillTypeBalanceFreeze      int8 = -5 // 余额冻结
	BillTypeTransferSend       int8 = -6 // 转账发送
	BillTypeSwapsSend          int8 = -7 // 闪兑发送
)

// BillTypeNames 账单类型名称
var BillTypeNames = map[int8]Map{
	BillTypeDeposit:            {"label": "充值", "value": BillTypeDeposit},
	BillTypeSystemReward:       {"label": "系统奖励", "value": BillTypeSystemReward},
	BillTypeRegisterReward:     {"label": "注册奖励", "value": BillTypeRegisterReward},
	BillTypeInviteReward:       {"label": "邀请奖励", "value": BillTypeInviteReward},
	BillTypeDistributionReward: {"label": "分销奖励", "value": BillTypeDistributionReward},
	BillTypeBalanceUnfreeze:    {"label": "余额解冻", "value": BillTypeBalanceUnfreeze},
	BillTypeSystemAddition:     {"label": "系统加款", "value": BillTypeSystemAddition},
	BillTypeWithdrawalReject:   {"label": "提现拒绝", "value": BillTypeWithdrawalReject},
	BillTypeProductRefund:      {"label": "产品退款", "value": BillTypeProductRefund},
	BillTypeProductEarnings:    {"label": "产品收益", "value": BillTypeProductEarnings},
	BillTypeTransferReceive:    {"label": "转账接收", "value": BillTypeTransferReceive},
	BillTypeSwapsReceive:       {"label": "闪兑接收", "value": BillTypeSwapsReceive},
	BillTypeWithdrawal:         {"label": "提现", "value": BillTypeWithdrawal},
	BillTypeSystemDeduction:    {"label": "系统扣款", "value": BillTypeSystemDeduction},
	BillTypeProductPurchase:    {"label": "购买产品", "value": BillTypeProductPurchase},
	BillTypeMembershipPurchase: {"label": "购买会员", "value": BillTypeMembershipPurchase},
	BillTypeBalanceFreeze:      {"label": "余额冻结", "value": BillTypeBalanceFreeze},
	BillTypeTransferSend:       {"label": "转账发送", "value": BillTypeTransferSend},
	BillTypeSwapsSend:          {"label": "闪兑发送", "value": BillTypeSwapsSend},
}

// GetWalletBillCheckboxOptions 获取钱包账单类型选项
func GetWalletBillCheckboxOptions() SliceMap {
	checkboxes := make(SliceMap, 0)
	for _, v := range BillTypeNames {
		v.Add("label", v["label"])
		v.Add("value", v["value"])
		v.Add("checked", v["checked"])
		checkboxes.Add(v)
	}
	return checkboxes
}

// WalletBill 钱包账单
type WalletBill struct {
	BaseModel
	AdminID  uint           `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	UserID   uint           `gorm:"type:int unsigned not null;index;comment:用户ID" json:"userId"`
	AssetsID uint           `gorm:"type:int unsigned not null;index;comment:资产ID" json:"assetsId"`
	SourceID uint           `gorm:"type:int unsigned not null;index;comment:来源ID" json:"sourceId"`
	Type     int8           `gorm:"type:tinyint not null;index;comment:类型" json:"type"`
	Name     string         `gorm:"type:varchar(60) not null;index;comment:名称" json:"name"`
	Money    float64        `gorm:"type:decimal(16,4) not null;comment:金额" json:"money"`
	Balance  float64        `gorm:"type:decimal(16,4) not null;comment:余额" json:"balance"`
	Desc     string         `gorm:"type:varchar(255);comment:描述" json:"desc"`
	Data     WalletBillData `gorm:"type:json;comment:数据" json:"data"`
}

// WalletBillDailyStats 钱包账单每日统计
type WalletBillDailyStats struct {
	Date    time.Time `json:"date"`
	Money   float64   `json:"money"`   // 已经乘以汇率的金额
	Balance float64   `json:"balance"` // 已经乘以汇率的余额
}

// WalletBillData 钱包账单数据
type WalletBillData struct {
}

// Value implements the driver.Valuer interface
func (d WalletBillData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *WalletBillData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&WalletBill{}); err != nil {
//		panic("Failed to auto migrate WalletBill table: " + err.Error())
//	}
//}
