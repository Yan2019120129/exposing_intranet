package models

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	CountryStatusDisabled int8 = -1 // 禁用
	CountryStatusEnabled  int8 = 10 // 启用
)

// Country 系统国家
type Country struct {
	BaseModel
	AdminID uint        `gorm:"type:int unsigned;uniqueIndex:idx_admin_iso2;not null;comment:管理ID" json:"adminId"`
	Name    string      `gorm:"type:varchar(60);not null;index;comment:国家名称" json:"name"`
	Alias   string      `gorm:"type:varchar(60);not null;comment:国家别名" json:"alias"`
	Icon    string      `gorm:"type:varchar(255);not null;comment:国旗图标URL" json:"icon"`
	ISO2    string      `gorm:"type:char(2);not null;uniqueIndex:idx_admin_iso2;comment:ISO 3166-1 alpha-2代码" json:"iso2"`
	Sort    int8        `gorm:"type:tinyint;not null;default:99;index;comment:排序" json:"sort"`
	Code    string      `gorm:"type:varchar(10);not null;comment:国家区号" json:"code"`
	Status  int8        `gorm:"type:tinyint;not null;default:10;index;comment:状态" json:"status"`
	Data    CountryData `gorm:"type:json;comment:数据" json:"data"`
}

// CountryDisplayData 国家显示数据
type CountryDisplayData struct {
	ID    uint   `json:"id"`    //	国家ID
	Alias string `json:"alias"` //	国家名称
	Icon  string `json:"icon"`  //	国家图标
	ISO2  string `json:"iso2"`  //	国家ISO2
	Code  string `json:"code"`  //	电话区号
}

// CountryData 国家数据
type CountryData struct {
}

// Value implements the driver.Valuer interface
func (d CountryData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *CountryData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Country{}); err != nil {
//		panic("Failed to auto migrate Country table: " + err.Error())
//	}
//
//	// Initialize system countries
//	if err := InitSystemCountries(db); err != nil {
//		panic("Failed to initialize system countries: " + err.Error())
//	}
//}

// InitSystemCountries initializes the default system countries
func InitSystemCountries(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Country{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		countries := []Country{
			{AdminID: SuperAdminID, Name: "中国", Alias: "China", Code: "86", Icon: "/country/china.png", ISO2: "CN", Sort: 1, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "美国", Alias: "United States", Code: "1", Icon: "/country/usa.png", ISO2: "US", Sort: 2, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "日本", Alias: "Japan", Code: "81", Icon: "/country/japan.png", ISO2: "JP", Sort: 3, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "韩国", Alias: "South Korea", Code: "82", Icon: "/country/south_korea.png", ISO2: "KR", Sort: 4, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "英国", Alias: "United Kingdom", Code: "44", Icon: "/country/uk.png", ISO2: "GB", Sort: 5, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "德国", Alias: "Germany", Code: "49", Icon: "/country/germany.png", ISO2: "DE", Sort: 6, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "法国", Alias: "France", Code: "33", Icon: "/country/france.png", ISO2: "FR", Sort: 7, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "加拿大", Alias: "Canada", Code: "1", Icon: "/country/canada.png", ISO2: "CA", Sort: 8, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "澳大利亚", Alias: "Australia", Code: "61", Icon: "/country/australia.png", ISO2: "AU", Sort: 9, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "新加坡", Alias: "Singapore", Code: "65", Icon: "/country/singapore.png", ISO2: "SG", Sort: 10, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "香港", Alias: "Hong kong", Code: "852", Icon: "/country/hongkong.png", ISO2: "HK", Sort: 11, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "台湾", Alias: "Taiwan", Code: "886", Icon: "/country/taiwan.png", ISO2: "TW", Sort: 12, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "澳门", Alias: "Macao", Code: "853", Icon: "/country/macao.png", ISO2: "MO", Sort: 13, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "印度", Alias: "India", Code: "91", Icon: "/country/india.png", ISO2: "IN", Sort: 14, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "俄罗斯", Alias: "Russia", Code: "7", Icon: "/country/russia.png", ISO2: "RU", Sort: 15, Status: CountryStatusEnabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "蒙古", Alias: "Mongolia", Code: "976", Icon: "/country/mongolia.png", ISO2: "MN", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "朝鲜", Alias: "North Korea", Code: "850", Icon: "/country/north_korea.png", ISO2: "KP", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "菲律宾", Alias: "Philippines", Code: "63", Icon: "/country/philippines.png", ISO2: "PH", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "越南", Alias: "Vietnam", Code: "84", Icon: "/country/vietnam.png", ISO2: "VN", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "老挝", Alias: "Laos", Code: "856", Icon: "/country/laos.png", ISO2: "LA", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "柬埔寨", Alias: "Cambodia", Code: "855", Icon: "/country/cambodia.png", ISO2: "KH", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "缅甸", Alias: "Myanmar", Code: "95", Icon: "/country/myanmar.png", ISO2: "MM", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "泰国", Alias: "Thailand", Code: "66", Icon: "/country/thailand.png", ISO2: "TH", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "马来西亚", Alias: "Malaysia", Code: "60", Icon: "/country/malaysia.png", ISO2: "MY", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "文莱", Alias: "Brunei", Code: "673", Icon: "/country/brunei.png", ISO2: "BN", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "印度尼西亚", Alias: "Indonesia", Code: "62", Icon: "/country/indonesia.png", ISO2: "ID", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "东帝汶", Alias: "Timor-Leste", Code: "670", Icon: "/country/east_timor.png", ISO2: "TL", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "尼泊尔", Alias: "Nepal", Code: "977", Icon: "/country/nepal.png", ISO2: "NP", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "不丹", Alias: "Bhutan", Code: "975", Icon: "/country/bhutan.png", ISO2: "BT", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "孟加拉国", Alias: "Bangladesh", Code: "880", Icon: "/country/bangladesh.png", ISO2: "BD", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "巴基斯坦", Alias: "Pakistan", Code: "92", Icon: "/country/pakistan.png", ISO2: "PK", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "斯里兰卡", Alias: "SriLanka", Code: "94", Icon: "/country/sri_lanka.png", ISO2: "LK", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "马尔代夫", Alias: "Maldives", Code: "960", Icon: "/country/maldives.png", ISO2: "MV", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "哈萨克斯坦", Alias: "Kazakhstan", Code: "7", Icon: "/country/kazakhstan.png", ISO2: "KZ", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "吉尔吉斯斯坦", Alias: "Kyrgyzstan", Code: "996", Icon: "/country/kyrgyzstan.png", ISO2: "KG", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "塔吉克斯坦", Alias: "Tajikistan", Code: "992", Icon: "/country/tajikistan.png", ISO2: "TJ", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "乌兹别克斯坦", Alias: "Uzbekistan", Code: "998", Icon: "/country/uzbekistan.png", ISO2: "UZ", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "土库曼斯坦", Alias: "Turkmenistan", Code: "993", Icon: "/country/turkmenistan.png", ISO2: "TM", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "阿富汗", Alias: "Afghanistan", Code: "93", Icon: "/country/afghanistan.png", ISO2: "AF", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "伊拉克", Alias: "Iraq", Code: "964", Icon: "/country/iraq.png", ISO2: "IQ", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "伊朗", Alias: "Iran", Code: "98", Icon: "/country/iran.png", ISO2: "IR", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "叙利亚", Alias: "Syria", Code: "963", Icon: "/country/syria.png", ISO2: "SY", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "约旦", Alias: "Jordan", Code: "962", Icon: "/country/jordan.png", ISO2: "JO", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "黎巴嫩", Alias: "Lebanon", Code: "961", Icon: "/country/lebanon.png", ISO2: "LB", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "以色列", Alias: "Israel", Code: "972", Icon: "/country/israel.png", ISO2: "IL", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "巴勒斯坦", Alias: "Palestine", Code: "970", Icon: "/country/palestine.png", ISO2: "PS", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "沙特阿拉伯", Alias: "SaudiArabia", Code: "966", Icon: "/country/saudi_arabia.png", ISO2: "SA", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "巴林", Alias: "Bahrain", Code: "973", Icon: "/country/bahrain.png", ISO2: "BH", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "卡塔尔", Alias: "Qatar", Code: "974", Icon: "/country/qatar.png", ISO2: "QA", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "科威特", Alias: "Kuwait", Code: "965", Icon: "/country/kuwait.png", ISO2: "KW", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "阿拉伯联合酋长国", Alias: "United Arab Emirates", Code: "971", Icon: "/country/united_arab_emirates.png", ISO2: "AE", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "阿曼", Alias: "Oman", Code: "968", Icon: "/country/oman.png", ISO2: "OM", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "也门", Alias: "Yemen", Code: "967", Icon: "/country/yemen.png", ISO2: "YE", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "格鲁吉亚", Alias: "Georgia", Code: "995", Icon: "/country/georgia.png", ISO2: "GE", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "亚美尼亚", Alias: "Armenia", Code: "374", Icon: "/country/armenia.png", ISO2: "AM", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "阿塞拜疆", Alias: "Azerbaijan", Code: "994", Icon: "/country/azerbaijan.png", ISO2: "AZ", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "土耳其", Alias: "Turkey", Code: "90", Icon: "/country/turkey.png", ISO2: "TR", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
			{AdminID: SuperAdminID, Name: "塞浦路斯", Alias: "Cyprus", Code: "357", Icon: "/country/cyprus.png", ISO2: "CY", Sort: 99, Status: CountryStatusDisabled, Data: CountryData{}},
		}

		if err := db.Create(&countries).Error; err != nil {
			return err
		}
	}

	return nil
}
