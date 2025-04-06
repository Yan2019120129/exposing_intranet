package gen

import (
	"gorm.io/gen"
	"my-base/config"
)

var GenDb *gen.Generator

// init 初始化gen
func init() {
	GenDb = gen.NewGenerator(gen.Config{
		//  设置输出路径
		OutPath: config.GetGen().OutPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // 选择生成模式
	})
}
