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
	c         *Config
	configErr error
)

// Config 服务配置
type Config struct {
	Port       string         `yaml:"port"`       // 服务端口
	Admin      *config.Config `yaml:"admin"`      // 管理配置信息
	Gorm       *gorm.Config   `yaml:"gorm"`       // gorm 配置
	ListenPort string         `yaml:"listenPort"` // 监听端口
	Connect    *ConnectConfig `yaml:"connect"`    // 连接配置
}

// ConnectConfig 连接配置
type ConnectConfig struct {
	PingInterval     int `yaml:"ping-interval"`     // 心跳发送间隔（秒）
	PongTimeout      int `yaml:"pong-timeout"`      // pong 响应超时（秒）
	MaxPingFailures  int `yaml:"max-ping-failures"` // 最大 ping 失败次数
	ConnectWait      int `yaml:"connect-wait"`      // 连接匹配等待超时（秒）
	TCPKeepAlive     int `yaml:"tcp-keepalive"`     // TCP Keep-Alive 间隔（秒）
	ReadWriteTimeout int `yaml:"readwrite-timeout"` // 数据转发读写超时（秒）
	WaitConnTimeout  int `yaml:"wait-conn-timeout"` // 等待连接超时清理（秒）
}

const (
	defaultPingInterval     = 30
	defaultPongTimeout      = 90
	defaultMaxPingFailures  = 3
	defaultConnectWait      = 10
	defaultTCPKeepAlive     = 30
	defaultReadWriteTimeout = 300
	defaultWaitConnTimeout  = 15
)

// GetConnect returns connection settings with safe defaults when the config
// file is missing or leaves a value unset.
func GetConnect() *ConnectConfig {
	if c == nil || c.Connect == nil {
		return &ConnectConfig{
			PingInterval: defaultPingInterval, PongTimeout: defaultPongTimeout,
			MaxPingFailures: defaultMaxPingFailures, ConnectWait: defaultConnectWait,
			TCPKeepAlive: defaultTCPKeepAlive, ReadWriteTimeout: defaultReadWriteTimeout,
			WaitConnTimeout: defaultWaitConnTimeout,
		}
	}
	return c.Connect
}

func (c *ConnectConfig) GetPingInterval() int {
	if c.PingInterval <= 0 {
		return defaultPingInterval
	}
	return c.PingInterval
}

func (c *ConnectConfig) GetPongTimeout() int {
	if c.PongTimeout <= 0 {
		return defaultPongTimeout
	}
	return c.PongTimeout
}

func (c *ConnectConfig) GetMaxPingFailures() int {
	if c.MaxPingFailures <= 0 {
		return defaultMaxPingFailures
	}
	return c.MaxPingFailures
}

func (c *ConnectConfig) GetConnectWait() int {
	if c.ConnectWait <= 0 {
		return defaultConnectWait
	}
	return c.ConnectWait
}

func (c *ConnectConfig) GetTCPKeepAlive() int {
	if c.TCPKeepAlive <= 0 {
		return defaultTCPKeepAlive
	}
	return c.TCPKeepAlive
}

func (c *ConnectConfig) GetReadWriteTimeout() int {
	if c.ReadWriteTimeout <= 0 {
		return defaultReadWriteTimeout
	}
	return c.ReadWriteTimeout
}

func (c *ConnectConfig) GetWaitConnTimeout() int {
	if c.WaitConnTimeout <= 0 {
		return defaultWaitConnTimeout
	}
	return c.WaitConnTimeout
}

// init 初始化服务端配置
func init() {
	c, configErr = loadConfig()
}

func loadConfig() (*Config, error) {
	path, err := findConfigPath()
	if err != nil {
		return nil, fmt.Errorf("find config.yaml: %w", err)
	}
	return loadConfigFile(path)
}

func loadConfigFile(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var cfg Config
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &cfg, nil
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

// RequireConfig returns the loaded configuration or its original load error.
// Callers that cannot operate without server configuration should use this at
// their boundary rather than dereferencing a nil Config later.
func RequireConfig() (*Config, error) {
	if configErr != nil {
		return nil, configErr
	}
	if c == nil {
		return nil, fmt.Errorf("server configuration was not loaded")
	}
	return c, nil
}

// GetGorm 获取gorm 配置
func GetGorm() *gorm.Config {
	if c == nil || c.Gorm == nil {
		return &gorm.Config{}
	}
	return c.Gorm
}

// GetAdmin 获取admin 配置
func GetAdmin() *config.Config {
	if c == nil {
		return nil
	}
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
