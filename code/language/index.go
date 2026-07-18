package language

import adminlanguage "github.com/GoAdminGroup/go-admin/modules/language"

// AddDefault 添加默认翻译
func AddDefault() {
	adminlanguage.AppendTo(adminlanguage.CN, map[string]string{
		"port":     "端口",
		"server":   "服务端",
		"client":   "客户端",
		"computer": "电脑",
	})

	adminlanguage.AppendTo(adminlanguage.EN, map[string]string{
		"port":     "port",
		"server":   "server",
		"client":   "client",
		"computer": "computer",
	})
}
