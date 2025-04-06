package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	WalletAssetsTypeDigitalCurrency int8 = 1 // 数字货币

	WalletAssetsStatusDisabled int8 = -1 // 禁用
	WalletAssetsStatusEnabled  int8 = 10 // 启用
)

// WalletAssets 钱包资产管理
type WalletAssets struct {
	BaseModel
	AdminID  uint             `gorm:"type:int unsigned not null;index;comment:管理ID" json:"adminId"`
	Name     string           `gorm:"type:varchar(60) not null;index;comment:名称" json:"name"`
	Currency string           `gorm:"type:varchar(10) not null;comment:货币符号" json:"currency"`
	Icon     string           `gorm:"type:varchar(255) not null;comment:图标" json:"icon"`
	Type     int8             `gorm:"type:tinyint not null;default:1;index;comment:类型(1:数字货币)" json:"type"`
	Rate     float64          `gorm:"type:decimal(16,4) not null;default:1;comment:汇率" json:"rate"`
	Status   int8             `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Decimals uint8            `gorm:"type:tinyint unsigned not null;default:2;comment:小数位数" json:"decimals"`
	Network  string           `gorm:"type:varchar(50);comment:网络" json:"network"`
	Desc     string           `gorm:"type:text;comment:描述" json:"desc"`
	Data     WalletAssetsData `gorm:"type:json;comment:数据" json:"data"`
}

// WalletAssetsDisplayData 钱包资产展示数据
type WalletAssetsDisplayData struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Currency string  `json:"currency"`
	Icon     string  `json:"icon"`
	Rate     float64 `json:"rate"`
	Decimals uint8   `json:"decimals"`
	Network  string  `json:"network"`
}

// WalletAssetsData 钱包资产数据
type WalletAssetsData struct {
}

// Value implements the driver.Valuer interface
func (d WalletAssetsData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *WalletAssetsData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&WalletAssets{}); err != nil {
//		panic("Failed to auto migrate WalletAssets table: " + err.Error())
//	}
//
//	// Initialize wallet assets
//	if err := InitWalletAssets(db); err != nil {
//		panic("Failed to initialize wallet assets: " + err.Error())
//	}
//}

// InitWalletAssets initializes the default wallet assets
func InitWalletAssets(db *gorm.DB) error {
	var count int64
	if err := db.Model(&WalletAssets{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		assets := []WalletAssets{
			{
				BaseModel: BaseModel{
					ID: 1,
				},
				AdminID:  SuperAdminID,
				Name:     "Bitcoin",
				Currency: "BTC",
				Icon:     "/icons/btc.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 4,
				Network:  "Bitcoin",
				Desc:     "Bitcoin is a decentralized digital currency.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 2,
				},
				AdminID:  SuperAdminID,
				Name:     "Ethereum",
				Currency: "ETH",
				Icon:     "/icons/eth.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "Ethereum",
				Desc:     "Ethereum is a decentralized, open-source blockchain with smart contract functionality.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 3,
				},
				AdminID:  SuperAdminID,
				Name:     "Tether",
				Currency: "USDT",
				Icon:     "/icons/usdt.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "Ethereum",
				Desc:     "Tether is a stablecoin pegged to the US dollar.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 4,
				},
				AdminID:  SuperAdminID,
				Name:     "USD Coin",
				Currency: "USDC",
				Icon:     "/icons/usdc.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "Ethereum",
				Desc:     "USD Coin is a stablecoin pegged to the US dollar.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 5,
				},
				AdminID:  SuperAdminID,
				Name:     "Binance Coin",
				Currency: "BNB",
				Icon:     "/icons/bnb.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "Binance Smart Chain",
				Desc:     "Binance Coin is the native cryptocurrency of the Binance ecosystem.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 6,
				},
				AdminID:  SuperAdminID,
				Name:     "Solana",
				Currency: "SOL",
				Icon:     "/icons/sol.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "Solana",
				Desc:     "Solana is a high-performance blockchain supporting smart contracts and decentralized applications.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 7,
				},
				AdminID:  SuperAdminID,
				Name:     "TRON",
				Currency: "TRX",
				Icon:     "/icons/trx.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "TRON",
				Desc:     "TRON is a blockchain-based decentralized platform that aims to build a free, global digital content entertainment system.",
				Data:     WalletAssetsData{},
			},
			{
				BaseModel: BaseModel{
					ID: 8,
				},
				AdminID:  SuperAdminID,
				Name:     "Ripple",
				Currency: "XRP",
				Icon:     "/icons/xrp.png",
				Type:     WalletAssetsTypeDigitalCurrency,
				Rate:     1,
				Status:   WalletAssetsStatusEnabled,
				Decimals: 2,
				Network:  "XRP Ledger",
				Desc:     "XRP is the native cryptocurrency of the XRP Ledger, designed for fast and cost-effective digital payments.",
				Data:     WalletAssetsData{},
			},
		}

		return db.CreateInBatches(assets, len(assets)).Error
	}

	return nil
}
