package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	// Payment Types
	PaymentTypeBankCard        int8 = 1
	PaymentTypeDigitalCurrency int8 = 11
	PaymentTypeThirdParty      int8 = 21

	// Payment Modes
	PaymentModeDeposit         int8 = 1  //	充值
	PaymentModeAssetDeposit    int8 = 2  // 资产充值
	PaymentModeWithdrawal      int8 = 11 // 提现
	PaymentModeAssetWithdrawal int8 = 12 // 资产提现

	// Payment Status
	PaymentStatusDisabled int8 = -1 // 禁用
	PaymentStatusEnabled  int8 = 10 // 激活
)

// WalletPayment 钱包支付管理
type WalletPayment struct {
	BaseModel
	AdminID   uint              `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	AssetsID  uint              `gorm:"type:int unsigned not null;default:0;index;comment:资产ID" json:"assetsId"`
	Name      string            `gorm:"type:varchar(60) not null;index;comment:名称" json:"name"`
	Icon      string            `gorm:"type:varchar(255) not null;comment:图标" json:"icon"`
	Currency  string            `gorm:"type:varchar(60) not null;comment:货币符号" json:"currency"`
	Type      int8              `gorm:"type:tinyint not null;default:1;index;comment:类型(1:银行卡,11:数字货币,21:三方支付)" json:"type"`
	Mode      int8              `gorm:"type:tinyint not null;default:1;index;comment:模式(1:充值,2:资产充值,11:提现,12:资产提现)" json:"mode"`
	MinAmount float64           `gorm:"type:decimal(16,2) not null;default:1;comment:最小金额" json:"minAmount"`
	MaxAmount float64           `gorm:"type:decimal(16,2) not null;default:1000000;comment:最大金额" json:"maxAmount"`
	StartTime string            `gorm:"type:varchar(60) not null;default:'00:00:00';comment:开始时间" json:"startTime"`
	EndTime   string            `gorm:"type:varchar(60) not null;default:'23:59:59';comment:结束时间" json:"endTime"`
	Fee       float64           `gorm:"type:decimal(16,4) not null;default:0;comment:手续费" json:"fee"`
	Level     int8              `gorm:"type:tinyint not null;default:0;comment:等级" json:"level"`
	Status    int8              `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Sort      int16             `gorm:"type:smallint not null;default:99;index;comment:排序" json:"sort"`
	IsProof   int8              `gorm:"type:tinyint not null;default:2;comment:是否需要凭证(1:是,2:否)" json:"isProof"`
	IsChats   int8              `gorm:"type:tinyint not null;default:2;comment:客服(1:真,2:假)" json:"isChats"`
	Data      WalletPaymentData `gorm:"type:json;comment:数据" json:"data"`
	Desc      string            `gorm:"type:varchar(255);comment:描述" json:"desc"`
}

// WalletPaymentData 钱包支付额外数据
type WalletPaymentData struct {
	BankName    string `json:"bankName" views:"label:名称｜公链"`    // 银行名称｜公链名称(Ethereum)
	BankAddress string `json:"bankAddress" views:"label:支付地址"`  // 银行地址｜公链标识(Erc20)
	RealName    string `json:"realName" views:"label:姓名|Token"` // 真实姓名｜公链Token(USDT)
	BankCardNo  string `json:"bankCardNo" views:"label:卡号|地址"`  // 银行卡号｜公链地址(0x1234567890abcdef)
	BankCode    string `json:"bankCode" views:"label:代号|简写"`    // 银行代号｜公链简写(ETH)
}

// Value implements the driver.Valuer interface
func (d WalletPaymentData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *WalletPaymentData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&WalletPayment{}); err != nil {
//		panic("Failed to auto migrate WalletPayment table: " + err.Error())
//	}
//
//	// Initialize wallet payments
//	if err := InitWalletPayments(db); err != nil {
//		panic("Failed to initialize wallet payments: " + err.Error())
//	}
//}

// InitWalletPayments initializes the default wallet payments
func InitWalletPayments(db *gorm.DB) error {
	var count int64
	if err := db.Model(&WalletPayment{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		payments := []WalletPayment{
			{
				AdminID:   SuperAdminID,
				AssetsID:  0,
				Name:      "Bank",
				Icon:      "/icons/bank.png",
				Currency:  "USD",
				Type:      PaymentTypeBankCard,
				Mode:      PaymentModeDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 100,
				MaxAmount: 100000,
				Fee:       0.01,
				Desc:      "Bank transfer payment method",
				Data: WalletPaymentData{
					BankName:    "Example Bank",
					BankCode:    "EXBK",
					RealName:    "John Doe",
					BankCardNo:  "1234567890",
					BankAddress: "123 Example St, City, Country",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  0,
				Name:      "USDT",
				Icon:      "/icons/usdt.png",
				Currency:  "USDT",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 1000000,
				Fee:       0.01, // Adjust fee as needed
				Desc:      "USDT deposit payment method",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDT",
					RealName:    "Tether",
					BankCardNo:  "0x1234567890123456789012345678901234567890",
					BankAddress: "ERC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  0,
				Name:      "Bank",
				Icon:      "/icons/bank.png",
				Currency:  "USD",
				Type:      PaymentTypeBankCard,
				Mode:      PaymentModeWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 100,
				MaxAmount: 50000,
				Fee:       0.01,
				Desc:      "Bank withdrawal payment method",
				Data: WalletPaymentData{
					BankName:    "",
					BankCode:    "",
					RealName:    "",
					BankCardNo:  "",
					BankAddress: "",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  0,
				Name:      "USDT",
				Icon:      "/icons/usdt.png",
				Currency:  "USDT",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 20,
				MaxAmount: 500000,
				Fee:       0.01, // Adjust fee as needed
				Desc:      "USDT withdrawal payment method",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDT",
					RealName:    "Tether",
					BankCardNo:  "",
					BankAddress: "ERC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  1,
				Name:      "BTC",
				Icon:      "/icons/btc.png",
				Currency:  "BTC",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.001,
				MaxAmount: 100,
				Fee:       0.01,
				Desc:      "Bitcoin deposit",
				Data: WalletPaymentData{
					BankName:    "Bitcoin",
					BankCode:    "BTC",
					RealName:    "Bitcoin",
					BankCardNo:  "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
					BankAddress: "Bitcoin",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  1,
				Name:      "BTC",
				Icon:      "/icons/btc.png",
				Currency:  "BTC",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.001,
				MaxAmount: 10,
				Fee:       0.01,
				Desc:      "Bitcoin withdrawal",
				Data: WalletPaymentData{
					BankName:    "Bitcoin",
					BankCode:    "BTC",
					RealName:    "Bitcoin",
					BankCardNo:  "",
					BankAddress: "Bitcoin",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  2,
				Name:      "ETH",
				Icon:      "/icons/eth.png",
				Currency:  "ETH",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.01,
				MaxAmount: 1000,
				Fee:       0.01,
				Desc:      "Ethereum deposit",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "ETH",
					RealName:    "Ethereum",
					BankCardNo:  "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
					BankAddress: "Ethereum",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  2,
				Name:      "ETH",
				Icon:      "/icons/eth.png",
				Currency:  "ETH",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.01,
				MaxAmount: 100,
				Fee:       0.01,
				Desc:      "Ethereum withdrawal",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "ETH",
					RealName:    "Ethereum",
					BankCardNo:  "",
					BankAddress: "Ethereum",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  3,
				Name:      "USDT",
				Icon:      "/icons/usdt.png",
				Currency:  "USDT",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 1000000,
				Fee:       0.01,
				Desc:      "USDT deposit on Ethereum network",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDT",
					RealName:    "Tether",
					BankCardNo:  "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
					BankAddress: "ERC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  3,
				Name:      "USDT",
				Icon:      "/icons/usdt.png",
				Currency:  "USDT",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 100000,
				Fee:       0.01,
				Desc:      "USDT withdrawal on Ethereum network",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDT",
					RealName:    "Tether",
					BankCardNo:  "",
					BankAddress: "ERC20",
				},
			},
			// Add similar structures for USDC, BNB, SOL, TRX, XRP
			// Example for USDC:
			{
				AdminID:   SuperAdminID,
				AssetsID:  4,
				Name:      "USDC",
				Icon:      "/icons/usdc.png",
				Currency:  "USDC",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 1000000,
				Fee:       0.01,
				Desc:      "USDC deposit",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDC",
					RealName:    "USDC Coin",
					BankCardNo:  "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
					BankAddress: "ERC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  4,
				Name:      "USDC",
				Icon:      "/icons/usdc.png",
				Currency:  "USDC",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 100000,
				Fee:       0.01,
				Desc:      "USDC withdrawal",
				Data: WalletPaymentData{
					BankName:    "Ethereum",
					BankCode:    "USDC",
					RealName:    "USDC Coin",
					BankCardNo:  "",
					BankAddress: "ERC20",
				},
			},
			// Continue with similar structures for BNB, SOL, TRX, XRP
			{
				AdminID:   SuperAdminID,
				AssetsID:  5,
				Name:      "BNB",
				Icon:      "/icons/bnb.png",
				Currency:  "BNB",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.1,
				MaxAmount: 10000,
				Fee:       0.01,
				Desc:      "BNB deposit",
				Data: WalletPaymentData{
					BankName:    "Binance",
					BankCode:    "BNB",
					RealName:    "BNB",
					BankCardNo:  "bnb1jxfh2g85q3v0tdq56fnevx6xcxtcnhtsmcu64m",
					BankAddress: "BEP20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  5,
				Name:      "BNB",
				Icon:      "/icons/bnb.png",
				Currency:  "BNB",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 0.1,
				MaxAmount: 1000,
				Fee:       0.01,
				Desc:      "BNB withdrawal",
				Data: WalletPaymentData{
					BankName:    "Binance",
					BankCode:    "BNB",
					RealName:    "BNB",
					BankCardNo:  "",
					BankAddress: "BEP20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  6,
				Name:      "SOL",
				Icon:      "/icons/sol.png",
				Currency:  "SOL",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 1,
				MaxAmount: 100000,
				Fee:       0.01,
				Desc:      "SOL deposit",
				Data: WalletPaymentData{
					BankName:    "Solana",
					BankCode:    "SOL",
					RealName:    "Solana",
					BankCardNo:  "7v91N7iZ9mNicL8WfG6cgSCKyRXydQjLh6UYBWwm6y1Q",
					BankAddress: "Solana",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  6,
				Name:      "SOL",
				Icon:      "/icons/sol.png",
				Currency:  "SOL",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 1,
				MaxAmount: 10000,
				Fee:       0.01,
				Desc:      "SOL withdrawal",
				Data: WalletPaymentData{
					BankName:    "Solana",
					BankCode:    "SOL",
					RealName:    "Solana",
					BankCardNo:  "",
					BankAddress: "Solana",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  7,
				Name:      "TRX",
				Icon:      "/icons/trx.png",
				Currency:  "TRX",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 100,
				MaxAmount: 10000000,
				Fee:       0.01,
				Desc:      "TRX deposit",
				Data: WalletPaymentData{
					BankName:    "TRON",
					BankCode:    "TRX",
					RealName:    "TRX",
					BankCardNo:  "TJYeasTPDdxYbsqd4Q5ci2ETpW5khfvjZj",
					BankAddress: "TRC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  7,
				Name:      "TRX",
				Icon:      "/icons/trx.png",
				Currency:  "TRX",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 100,
				MaxAmount: 1000000,
				Fee:       0.01,
				Desc:      "TRX withdrawal",
				Data: WalletPaymentData{
					BankName:    "TRON",
					BankCode:    "TRX",
					RealName:    "TRX",
					BankCardNo:  "",
					BankAddress: "TRC20",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  8,
				Name:      "XRP",
				Icon:      "/icons/xrp.png",
				Currency:  "XRP",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetDeposit,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 1000000,
				Fee:       0.01,
				Desc:      "XRP deposit",
				Data: WalletPaymentData{
					BankName:    "Ripple",
					BankCode:    "XRP",
					RealName:    "XRP",
					BankCardNo:  "rEb8TK3gBgk5auZkwc6sHnwrGVJH8DuaLh",
					BankAddress: "XRPL",
				},
			},
			{
				AdminID:   SuperAdminID,
				AssetsID:  8,
				Name:      "XRP",
				Icon:      "/icons/xrp.png",
				Currency:  "XRP",
				Type:      PaymentTypeDigitalCurrency,
				Mode:      PaymentModeAssetWithdrawal,
				Status:    PaymentStatusEnabled,
				MinAmount: 10,
				MaxAmount: 100000,
				Fee:       0.01,
				Desc:      "XRP withdrawal",
				Data: WalletPaymentData{
					BankName:    "Ripple",
					BankCode:    "XRP",
					RealName:    "XRP",
					BankCardNo:  "",
					BankAddress: "XRPL",
				},
			},
		}

		return db.CreateInBatches(payments, len(payments)).Error
	}

	return nil
}
