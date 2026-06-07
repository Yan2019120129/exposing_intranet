package configs

import (
	"fmt"
	"my-base/configs/gorm"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"gopkg.in/yaml.v3"
)

// Config 服务端配置实例
var (
	c *Config
)

// Config 服务配置
type Config struct {
	Port  string         `yaml:"port"`  // 服务端口
	Admin *config.Config `yaml:"admin"` // 管理配置信息
	Gorm  *gorm.Config   `yaml:"gorm"`  // gorm 配置
}

// init 初始化服务端配置
func init() {
	path, err := findConfigPath()
	if err != nil {
		fmt.Println("finding config.yaml error:", err)
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("reading config.yaml error:", err)
		return
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		fmt.Println("unmarshal config.yaml error:", err)
		return
	}
}

func findConfigPath() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		candidate := filepath.Join(current, "website", "configs", "config.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", os.ErrNotExist
}

// GetConfig 获取服务信息
func GetConfig() *Config {
	return c
}

// GetGorm 获取gorm 配置
func GetGorm() *gorm.Config {
	return c.Gorm
}

// GetAdmin 获取admin 配置
func GetAdmin() *config.Config {
	return c.Admin
}

// GetAdmin 获取服务端管理信息
func (s *Config) GetAdmin() *config.Config {
	return s.Admin
}

// SetAdmin 设置服务端管理信息
func (s *Config) SetAdmin(admin *config.Config) *Config {
	s.Admin = admin
	return s
}
