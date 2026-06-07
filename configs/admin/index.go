package admin

// Config 管理配置信息
type Config struct {
	name string `yaml:"name"` // 服务名称
	psw  string `yaml:"psw"`  // 服务密码
}

// SetName 设置管理名
func (a *Config) SetName(name string) *Config {
	a.name = name
	return a
}

// GetName 获取管理名
func (a *Config) GetName() string {
	return a.name
}

// GetPsw 获取管理密码
func (a *Config) GetPsw() string {
	return a.psw
}

// SetPsw 设置管理密码
func (a *Config) SetPsw(psw string) *Config {
	a.psw = psw
	return a
}
