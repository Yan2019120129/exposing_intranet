package models

const (
	ProductSkuStatusEnabled  int8 = 10 // 上架
	ProductSkuStatusDisabled int8 = -1 // 下架
)

// ProductSku 产品SKU
type ProductSku struct {
	BaseModel
	AdminID   uint    `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	ProductID uint    `gorm:"type:int unsigned not null;index;comment:产品ID" json:"productId"`
	Name      string  `gorm:"type:varchar(512) not null;comment:SKU名称" json:"name"`
	Image     string  `gorm:"type:varchar(255) not null;comment:商品图片" json:"image"`
	Stock     uint    `gorm:"type:int unsigned not null;default:1000;comment:库存" json:"stock"`
	LockStock uint    `gorm:"type:int unsigned not null;default:0;comment:锁定库存" json:"lockStock"`
	Sales     uint    `gorm:"type:int unsigned not null;default:0;comment:销量" json:"sales"`
	Money     float64 `gorm:"type:decimal(12, 2) not null;default:100;comment:原价" json:"money"`
	Discount  float64 `gorm:"type:decimal(12, 2);default:0;comment:折扣价" json:"discount"`
	Status    int8    `gorm:"type:tinyint not null;default:10;index;comment:状态(10:上架,-1:下架)" json:"status"`
	Desc      string  `gorm:"type:text;comment:描述" json:"desc"`
}

// ProductSkuAttribute 产品SKU属性
type ProductSkuAttribute struct {
	ProductSku
	SkuAttribute []SkuAttribute `gorm:"-" json:"skuAttribute"`
}

// SkuInfo 购买信息
type SkuInfo struct {
	SkuID  uint `json:"skuId"`  // skuId
	Nums   uint `json:"nums"`   // 购买数量
	CartID uint `json:"cartId"` // 购物车ID
}

// SkuNumsInfo Sku数量
type SkuNumsInfo struct {
	SkuInfo *ProductSku // sku
	Nums    uint        // 购买数量
}

func (cp *ProductSkuAttribute) TableName() string {
	return "product_sku"
}

// GetTotalPrice 获取sku总价,原价*数量
func (p *ProductSku) GetTotalPrice(nums float64) float64 {
	return p.Money * nums
}

// GetFinalPrice 获取最终价,原价*数量-折扣价*数量
func (p *ProductSku) GetFinalPrice(nums float64) float64 {
	// discount是100时不打折
	if p.Discount == 100 {
		return p.GetTotalPrice(nums)
	}
	return p.GetTotalPrice(nums) - p.GetTotalPrice(nums)*(p.Discount/100)
}

// GetEarning 获取收益 (最终价-成本价)*数量
func (p *ProductSku) GetEarning(nums, increase float64) float64 {
	return p.GetFinalPrice(nums) * (increase / 100)
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&ProductSku{}); err != nil {
//		panic("Failed to auto migrate ProductSku table: " + err.Error())
//	}
//}
