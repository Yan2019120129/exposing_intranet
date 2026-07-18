package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "website/configs/config.client.yaml"
	defaultSymbolPath = "website/configs/client.key"

	defaultPingInterval       = 30
	defaultPongTimeout        = 90
	defaultMaxPingFailures    = 3
	defaultReconnectBaseDelay = 2
	defaultReconnectMaxDelay  = 60
)

// Config is the configuration root for the client executable.
type Config struct {
	Client *ClientConfig `yaml:"client"`
}

// ClientConfig contains settings for the server control connection and the
// local client runtime.
type ClientConfig struct {
	ServerAddr         string `yaml:"server-addr"`
	JobPort            string `yaml:"job-port"`
	Port               string `yaml:"port"`
	PingInterval       int    `yaml:"ping-interval"`
	PongTimeout        int    `yaml:"pong-timeout"`
	MaxPingFailures    int    `yaml:"max-ping-failures"`
	ReconnectBaseDelay int    `yaml:"reconnect-base-delay"`
	ReconnectMaxDelay  int    `yaml:"reconnect-max-delay"`
}

var current = load()

func load() *Config {
	path := resolvePath(os.Getenv("EXPOSING_INTRANET_CLIENT_CONFIG"), defaultConfigPath)

	data, err := os.ReadFile(path)
	if err != nil {
		return withDefaults(&Config{Client: &ClientConfig{}})
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("unmarshal client config: %v\n", err)
		return withDefaults(&Config{Client: &ClientConfig{}})
	}
	if cfg.Client == nil {
		cfg.Client = &ClientConfig{}
	}
	return withDefaults(&cfg)
}

func withDefaults(cfg *Config) *Config {
	if cfg.Client == nil {
		cfg.Client = &ClientConfig{}
	}
	if cfg.Client.PingInterval <= 0 {
		cfg.Client.PingInterval = defaultPingInterval
	}
	if cfg.Client.PongTimeout <= 0 {
		cfg.Client.PongTimeout = defaultPongTimeout
	}
	if cfg.Client.MaxPingFailures <= 0 {
		cfg.Client.MaxPingFailures = defaultMaxPingFailures
	}
	if cfg.Client.ReconnectBaseDelay <= 0 {
		cfg.Client.ReconnectBaseDelay = defaultReconnectBaseDelay
	}
	if cfg.Client.ReconnectMaxDelay <= 0 {
		cfg.Client.ReconnectMaxDelay = defaultReconnectMaxDelay
	}
	return cfg
}

func symbolPath() string {
	if path := os.Getenv("EXPOSING_INTRANET_SYMBOL_PATH"); path != "" {
		return filepath.Clean(path)
	}
	return resolvePath("", defaultSymbolPath)
}

func resolvePath(configured, relative string) string {
	if configured != "" {
		return filepath.Clean(configured)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return filepath.Clean(relative)
	}
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(currentDir, relative)
		if _, err := os.Stat(filepath.Dir(candidate)); err == nil {
			return candidate
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
	return filepath.Clean(relative)
}

// GetConfig returns the process configuration.
func GetConfig() *Config { return current }

// GetClientConfig returns the client settings.
func (c *Config) GetClientConfig() *ClientConfig {
	if c == nil || c.Client == nil {
		return withDefaults(&Config{Client: &ClientConfig{}}).Client
	}
	return c.Client
}

// GetClientConfig returns the process client settings.
func GetClientConfig() *ClientConfig { return current.GetClientConfig() }

// GetSymbol loads the authenticated client symbol.
func (c *Config) GetSymbol() string {
	path := symbolPath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		fmt.Printf("read client symbol: %v\n", err)
		return ""
	}
	return string(data)
}

// SetSymbol persists the authenticated client symbol.
func (c *Config) SetSymbol(symbol string) {
	path := symbolPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Printf("create symbol directory: %v\n", err)
		return
	}
	if err := os.WriteFile(path, []byte(symbol), 0644); err != nil {
		fmt.Printf("write client symbol: %v\n", err)
	}
}

func (c *ClientConfig) GetServerAddr() string { return c.ServerAddr }
func (c *ClientConfig) GetJobPort() string    { return c.JobPort }
func (c *ClientConfig) GetPort() string       { return c.Port }

func (c *ClientConfig) GetPublicAddr() string {
	return c.GetServerAddr() + ":" + c.GetPort()
}

func (c *ClientConfig) GetPublicJobAddr() string {
	return c.GetServerAddr() + ":" + c.GetJobPort()
}

func (c *ClientConfig) GetPingInterval() int {
	if c.PingInterval <= 0 {
		return defaultPingInterval
	}
	return c.PingInterval
}

func (c *ClientConfig) GetPongTimeout() int {
	if c.PongTimeout <= 0 {
		return defaultPongTimeout
	}
	return c.PongTimeout
}

func (c *ClientConfig) GetMaxPingFailures() int {
	if c.MaxPingFailures <= 0 {
		return defaultMaxPingFailures
	}
	return c.MaxPingFailures
}

func (c *ClientConfig) GetReconnectBaseDelay() int {
	if c.ReconnectBaseDelay <= 0 {
		return defaultReconnectBaseDelay
	}
	return c.ReconnectBaseDelay
}

func (c *ClientConfig) GetReconnectMaxDelay() int {
	if c.ReconnectMaxDelay <= 0 {
		return defaultReconnectMaxDelay
	}
	return c.ReconnectMaxDelay
}
