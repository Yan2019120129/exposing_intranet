package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	LangStatusEnabled  int8 = 10 // 启用
	LangStatusDisabled int8 = -1 // 禁用
)

// Lang 系统语言
type Lang struct {
	BaseModel
	AdminID uint     `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_symbol;comment:管理ID" json:"adminId"`
	Name    string   `gorm:"type:varchar(60) not null;index;comment:名称" json:"name"`
	Alias   string   `gorm:"type:varchar(60);not null;comment:语言别名" json:"alias"`
	Symbol  string   `gorm:"type:varchar(60) not null;uniqueIndex:idx_admin_symbol;comment:标识" json:"symbol"`
	Icon    string   `gorm:"type:varchar(255) not null;comment:图标" json:"icon"`
	Sort    int8     `gorm:"type:tinyint unsigned not null;default:99;index;comment:排序" json:"sort"`
	Status  int8     `gorm:"type:tinyint not null;default:10;index;comment:状态(10:启用,-1:禁用)" json:"status"`
	Data    LangData `gorm:"type:json;comment:数据" json:"data"`
}

// LangDisplayData 语言显示数据
type LangDisplayData struct {
	ID     uint   `json:"id"`     //	语言ID
	Alias  string `json:"alias"`  //	语言名称
	Icon   string `json:"icon"`   //	语言图标
	Symbol string `json:"symbol"` //	语言别名
}

// LangData 语言数据
type LangData struct {
}

// Value implements the driver.Valuer interface
func (d LangData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *LangData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Lang{}); err != nil {
//		panic("Failed to auto migrate Lang table: " + err.Error())
//	}
//
//	// Initialize system languages
//	if err := InitSystemLangs(db); err != nil {
//		panic("Failed to initialize system languages: " + err.Error())
//	}
//}

// InitSystemLangs initializes system languages
func InitSystemLangs(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Lang{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		langs := []Lang{
			{
				AdminID: SuperAdminID,
				Name:    "中文",
				Alias:   "简体中文",
				Symbol:  "zh-CN",
				Icon:    "/country/china.png",
				Sort:    1,
				Status:  LangStatusEnabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "English",
				Alias:   "English",
				Symbol:  "en-US",
				Icon:    "/country/usa.png",
				Sort:    2,
				Status:  LangStatusEnabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "繁体中文",
				Alias:   "繁體中文",
				Symbol:  "zh-TW",
				Icon:    "/country/taiwan.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "阿拉伯语",
				Alias:   "عربي",
				Symbol:  "ar-AE",
				Icon:    "/country/united_arab_emirates.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "白俄罗斯语",
				Alias:   "беларускі",
				Symbol:  "be-BY",
				Icon:    "/country/belarus.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "保加利亚语",
				Alias:   "български",
				Symbol:  "bg-BG",
				Icon:    "/country/bulgaria.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "捷克语",
				Alias:   "čeština",
				Symbol:  "cs-CZ",
				Icon:    "/country/czech.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "丹麦语",
				Alias:   "dansk",
				Symbol:  "da-DK",
				Icon:    "/country/denmark.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "德语",
				Alias:   "Deutsch",
				Symbol:  "de-DE",
				Icon:    "/country/germany.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "希腊语",
				Alias:   "Ελληνικά",
				Symbol:  "el-GR",
				Icon:    "/country/greece.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "西班牙语",
				Alias:   "español",
				Symbol:  "es-ES",
				Icon:    "/country/spain.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "爱沙尼亚语",
				Alias:   "eestikeel",
				Symbol:  "et-EE",
				Icon:    "/country/estonia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "芬兰语",
				Alias:   "Suomalainen",
				Symbol:  "fi-FI",
				Icon:    "/country/finland.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "法语",
				Alias:   "Français",
				Symbol:  "fr-FR",
				Icon:    "/country/france.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "克罗地亚语",
				Alias:   "Hrvatski",
				Symbol:  "hr-HR",
				Icon:    "/country/croatia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "匈牙利语",
				Alias:   "Magyar",
				Symbol:  "hu-HU",
				Icon:    "/country/hungary.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "冰岛语",
				Alias:   "íslenskur",
				Symbol:  "is-IS",
				Icon:    "/country/iceland.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "意大利语",
				Alias:   "italiano",
				Symbol:  "it-IT",
				Icon:    "/country/italy.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "日语",
				Alias:   "日本",
				Symbol:  "ja-JP",
				Icon:    "/country/japan.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "马来语",
				Alias:   "Melayu",
				Symbol:  "ms-MY",
				Icon:    "/country/malaysia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "越南语",
				Alias:   "TiếngViệt",
				Symbol:  "vi-VN",
				Icon:    "/country/vietnam.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "韩语",
				Alias:   "한국인",
				Symbol:  "ko-KR",
				Icon:    "/country/north_korea.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "立陶宛语",
				Alias:   "lietuvių",
				Symbol:  "lt-LT",
				Icon:    "/country/lithuania.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "马其顿语",
				Alias:   "македонски",
				Symbol:  "mk-MK",
				Icon:    "/country/macedonia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "荷兰语",
				Alias:   "Nederlands",
				Symbol:  "nl-NL",
				Icon:    "/country/netherlands.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "挪威语",
				Alias:   "norsk",
				Symbol:  "no-NO",
				Icon:    "/country/norway.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "波兰语",
				Alias:   "Polski",
				Symbol:  "pl-PL",
				Icon:    "/country/poland.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "葡萄牙语",
				Alias:   "Português",
				Symbol:  "pt-PT",
				Icon:    "/country/portugal.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "罗马尼亚语",
				Alias:   "Română",
				Symbol:  "ro-RO",
				Icon:    "/country/romania.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "俄语",
				Alias:   "Русский",
				Symbol:  "ru-RU",
				Icon:    "/country/russia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "斯洛伐克语",
				Alias:   "slovenský",
				Symbol:  "sk-SK",
				Icon:    "/country/slovakia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "斯洛文尼亚语",
				Alias:   "Slovenščina",
				Symbol:  "sl-SI",
				Icon:    "/country/slovenia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "阿尔巴尼亚语",
				Alias:   "shqiptare",
				Symbol:  "sq-AL",
				Icon:    "/country/albania.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "瑞典语",
				Alias:   "svenska",
				Symbol:  "sv-SE",
				Icon:    "/country/sweden.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "泰语",
				Alias:   "แบบไทย",
				Symbol:  "th-TH",
				Icon:    "/country/thailand.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "土耳其语",
				Alias:   "Türkçe",
				Symbol:  "tr-TR",
				Icon:    "/country/turkey.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "乌克兰语",
				Alias:   "українська",
				Symbol:  "uk-UA",
				Icon:    "/country/ukraine.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "塞尔维亚语",
				Alias:   "Српски",
				Symbol:  "sr-YU",
				Icon:    "/country/serbia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "希伯来语",
				Alias:   "עִברִית",
				Symbol:  "iw-IL",
				Icon:    "/country/israel.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "印地语",
				Alias:   "हिंदी",
				Symbol:  "hi-IN",
				Icon:    "/country/india.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
			{
				AdminID: SuperAdminID,
				Name:    "印尼语",
				Alias:   "Indonesia",
				Symbol:  "id-ID",
				Icon:    "/country/indonesia.png",
				Sort:    99,
				Status:  LangStatusDisabled,
			},
		}

		return db.CreateInBatches(langs, len(langs)).Error
	}

	return nil
}
