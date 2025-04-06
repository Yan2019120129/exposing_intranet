package models

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

const (
	// LevelNameTranslatePrefix 等级名称翻译前缀
	LevelNameTranslatePrefix = "levelName%v"
	// LevelDescTranslatePrefix 等级详情翻译前缀
	LevelDescTranslatePrefix = "levelDesc%v"

	LevelStatusDisabled int8 = -1 // 禁用
	LevelStatusEnabled  int8 = 10 // 启用

	LevelTypeMember int8 = 1 // 会员等级
)

// Level 系统等级配置
type Level struct {
	BaseModel
	AdminID uint      `gorm:"type:int unsigned not null;uniqueIndex:idx_admin_symbol;comment:管理ID" json:"adminId"`
	Name    string    `gorm:"type:varchar(60) not null;comment:名称" json:"name"`
	Icon    string    `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Symbol  int8      `gorm:"type:tinyint not null;uniqueIndex:idx_admin_symbol;comment:标识" json:"symbol"`
	Type    int8      `gorm:"type:tinyint not null;default:1;index;comment:类型(1:等级)" json:"type"`
	Money   float64   `gorm:"type:decimal(12,2) not null;comment:金额" json:"money"`
	Days    int       `gorm:"type:smallint not null;comment:天数" json:"days"`
	Status  int8      `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
	Desc    string    `gorm:"type:text;comment:详情" json:"desc"`
	Data    LevelData `gorm:"type:json;comment:数据" json:"data"`
}

// LevelData 系统等级数据
type LevelData struct {
	Discount float64  `json:"discount"` // 折扣率
	Benefits SliceMap `json:"benefits"` // 等级权益
}

// Value implements the driver.Valuer interface
func (d LevelData) Value() (driver.Value, error) {
	return json.Marshal(&d)
}

// Scan implements the sql.Scanner interface
func (d *LevelData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &d)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Level{}); err != nil {
//		panic("Failed to auto migrate Level table: " + err.Error())
//	}
//
//	// Initialize system levels
//	if err := InitSystemLevels(db); err != nil {
//		panic("Failed to initialize system levels: " + err.Error())
//	}
//}

// InitSystemLevels initializes the default system levels
func InitSystemLevels(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Level{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		levels := initSystemLevels()
		for _, v := range levels {
			v.Name = fmt.Sprintf(LevelNameTranslatePrefix, v.Symbol)
			v.Desc = fmt.Sprintf(LevelDescTranslatePrefix, v.Symbol)
		}
		return db.CreateInBatches(levels, len(levels)).Error
	}

	return nil
}

func initSystemLevels() []*Level {
	levels := []*Level{
		{
			AdminID: SuperAdminID,
			Symbol:  1,
			Name:    "临时会员",
			Desc: `<div>
				<ul>
					<li>体验正式会员所有标准权益</li>
					<li>逐步解锁订单：系统推送订单，需完成后方可接受下一个</li>
					<li>正式入驻平台，对接Shopee全球供应链体系</li>
				</ul>
			</div>`,
			Icon:  "/icons/level1.png",
			Money: 0,
			Days:  30,
			Data: LevelData{
				Discount: 5,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  2,
			Name:    "正式会员",
			Desc: `<div>
				<ul>
					<li>首次加入赠送$1,000购物基金</li>
					<li>可同时接受5个仓位订单</li>
					<li>享受95折优惠</li>
					<li>更强大的精准引流系统</li>
					<li>专属会员服务</li>
				</ul>
			</div>`,
			Icon:  "/icons/level2.png",
			Money: 0,
			Days:  365,
			Data: LevelData{
				Discount: 5,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  3,
			Name:    "白银会员",
			Desc: `<div>
				<ul>
					<li>支持7个仓位订单</li>
					<li>享受93折优惠</li>
					<li>更高的流量扶持</li>
					<li>每月定期举办会员专属福利活动</li>
				</ul>
			</div>`,
			Icon:  "/icons/level3.png",
			Money: 0,
			Days:  365,
			Data: LevelData{
				Discount: 7,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  4,
			Name:    "黄金会员",
			Desc: `<div>
				<ul>
					<li>额外奖励$20,000购物额度</li>
					<li>可同时管理10个仓位订单</li>
					<li>尊享92折优惠</li>
					<li>可接受奢侈品类高利润订单</li>
				</ul>
			</div>`,
			Icon:  "/icons/level4.png",
			Money: 0,
			Days:  1095, // 3年
			Data: LevelData{
				Discount: 8,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  5,
			Name:    "白金会员",
			Desc: `<div>
				<ul>
					<li>额外奖励$30,000购物额度</li>
					<li>订单管理无仓位限制</li>
					<li>尊享91折优惠</li>
					<li>支持"先发货，后付款"模式</li>
					<li>专属仓库和VIP物流服务</li>
				</ul>
			</div>`,
			Icon:  "/icons/level5.png",
			Money: 0,
			Days:  1095, // 3年
			Data: LevelData{
				Discount: 9,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  6,
			Name:    "钻石会员",
			Desc: `<div>
				<ul>
					<li>额外奖励$100,000购物折扣额度</li>
					<li>尊享86折采购折扣</li>
					<li>支持先发货后付款</li>
					<li>专享定制物流、VIP仓储</li>
					<li>专属LINE客服、VIP流量扶持</li>
					<li>每年15日新加坡总部豪华旅行</li>
				</ul>
			</div>`,
			Icon:  "/icons/level6.png",
			Money: 0,
			Days:  1825, // 5年
			Data: LevelData{
				Discount: 14,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  7,
			Name:    "合伙人",
			Desc: `<div>
				<ul>
					<li>成为公司合伙人，获得专属身份勋章</li>
					<li>机场贵宾厅、全球酒店VIP通道</li>
					<li>品牌方直面洽谈</li>
					<li>每年3个月全球商旅计划</li>
					<li>定期参加顶级商业培训和研讨会</li>
				</ul>
			</div>`,
			Icon:  "/icons/level7.png",
			Money: 0,
			Days:  -1, // 永久（100年）
			Data: LevelData{
				Discount: 20,
			},
		},
		{
			AdminID: SuperAdminID,
			Symbol:  8,
			Name:    "至尊合伙人",
			Desc: `<div>
				<ul>
					<li>行业最低采购价</li>
					<li>独立品牌运营权</li>
					<li>全球仓储和极速物流</li>
					<li>1个月全球私人游轮旅行</li>
					<li>终身酒店VIP待遇</li>
					<li>每年5次私人飞机出行</li>
					<li>尊享公司股份分红</li>
				</ul>
			</div>`,
			Icon:  "/icons/level8.png",
			Money: 0,
			Days:  -1, // 永久（100年）
			Data: LevelData{
				Discount: 30,
			},
		},
	}

	return levels
}
