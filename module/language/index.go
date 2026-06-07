package language

import "github.com/GoAdminGroup/go-admin/modules/language"

// AddDefault 添加默认翻译
func AddDefault() {
	language.AppendTo(language.CN, map[string]string{
		"port":     "端口",
		"server":   "服务端",
		"client":   "客户端",
		"computer": "电脑",
	})

	language.AppendTo(language.EN, map[string]string{
		"port":     "port",
		"server":   "server",
		"client":   "client",
		"computer": "computer",
	})
}
