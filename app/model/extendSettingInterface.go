package models

import (
	"time"
)

const (
	// 扩展设置 - 分类设置
	CategorySettings = "categorySettings"

	// 扩展设置 - 产品设置
	ProductSettings = "productSettings"

	// 访问量配置 - 店铺设置
	StoreSettingsCreateAccess = "storeSettingsCreateAccess"

	// 下单配置 - 店铺设置
	StoreSettingsCreateOrder = "storeSettingsCreateOrder"
)

type StoreTrafficAdd struct {
	Time time.Time `json:"time"` // 时间
	Nums int       `json:"nums"` // 访问记录信息
}

// AutomaticOrderPlacement 自动下单信息
type AutomaticOrderPlacement struct {
	UserId         uint              `json:"userId"`         // 用户Id
	Username       string            `json:"username"`       // 用户名称
	Money          float64           `json:"money"`          // 下单总金额
	SkuIds         []uint            `json:"skuIds"`         // skuIds
	ProductSkuList []*ProductSkuInfo `json:"productSkuList"` // 产品Sku信息集合
	Time           time.Time         `json:"time"`           // 时间
}

// ProductSkuInfo 下单产品信息
type ProductSkuInfo struct {
	ProductID uint    `json:"productId"` // 产品ID
	SkuID     uint    `json:"skuId"`     // skuId
	SkuName   string  `json:"skuName"`   //规格名称
	Image     string  `json:"image"`     // 产品图
	Name      string  `json:"name"`      // 产品名
	Money     float64 `json:"money"`     // 产品金额
	Numb      uint    `json:"numb"`      // 下单数量
}
