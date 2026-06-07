package crawling

type ProductAttr struct {
	OriginalPrice float64             // 原价
	CurrentPrice  float64             // 现价
	Images        []string            // 图片
	Style         map[string][]string // 样式
	Title         string              // 标题
	Describe      string              // 描述
}
