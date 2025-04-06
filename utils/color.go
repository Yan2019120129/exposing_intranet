package utils

import (
	"fmt"
)

type Color string

const (
	// Bold 加粗
	Bold = ";1m"
	// Reset 重置
	Reset Color = "\033[0m"
	// Purple 紫色
	Purple Color = "\x1b[35m"
	// Green 绿色
	Green Color = "\x1b[32m"
	// Yellow 黄色
	Yellow Color = "\x1b[33m"
	// Red 红色
	Red Color = "\x1b[31m"
	// Blue 蓝色
	Blue Color = "\033[34m"
	// Magenta 洋红/紫红
	Magenta Color = "\033[35m"
	// Cyan 青色
	Cyan Color = "\033[36m"
	// White 白色
	White Color = "\033[37m"
)

// SetColor 设置颜色
func SetColor(color Color, msg interface{}, isBold bool) string {
	if isBold {
		return fmt.Sprintf("%v%v%v%v%v", color, Bold, msg, Reset)
	}
	return fmt.Sprintf("%v%v%v", color, msg, Reset)
}

// SetPurple 设置紫色
func SetPurple(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Purple, msg, Reset)
}

// SetGreen 设置绿色
func SetGreen(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Green, msg, Reset)
}

// SetYellow 设置黄色
func SetYellow(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Yellow, msg, Reset)
}

// SetRed 设置红色
func SetRed(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Red, msg, Reset)
}

// SetBlue 设置紫色
func SetBlue(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Blue, msg, Reset)
}

// SetMagenta 设置绿色
func SetMagenta(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Magenta, msg, Reset)
}

// SetCyan 设置黄色
func SetCyan(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", Cyan, msg, Reset)
}

// SetWhite 设置红色
func SetWhite(msg interface{}) string {
	return fmt.Sprintf("%v%v%v", White, msg, Reset)
}
