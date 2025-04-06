package config

import (
	"gopkg.in/yaml.v3"
	"gorm.io/gen"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	DatabaseTypeMysql      = "mysql"
	DatabaseTypePostgresql = "postgresql"
	// LogPath 日志输出位置
	LogPath = "/Users/taozi/Documents/Golang/gotest/common/file/logs/app.log"
)

// GenMode gorm gen 生成模式
var GenMode = map[string]interface{}{
	"WithoutContext":     gen.WithoutContext,
	"WithDefaultQuery":   gen.WithDefaultQuery,
	"WithQueryInterface": gen.WithQueryInterface,
}

type Config struct {
	Gorm          GormConfig          `yaml:"gorm"`
	Gin           GinConfig           `yaml:"gin"`
	Redis         RedisConfig         `yaml:"redis"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	Logs          LogConfig           `yaml:"logs"`
}

// ElasticsearchConfig 配置参数。
type ElasticsearchConfig struct {
	IpAddress             []string `yaml:"ip-address"`              // IP地址
	MaxIdleConnsPerHost   int      `yaml:"max-idle-conns-per-host"` // 每个主机的最大空闲连接数
	ResponseHeaderTimeout int      `yaml:"response-header-timeout"` // 接收响应头的超时时间
	DialerTimeout         int      `yaml:"dialer-timeout"`          // 建立连接的上下文和超时设置
}

// Database 数据库配置
type Database struct {
	Mysql      DatabaseConfig `yaml:"mysql"`
	Postgresql DatabaseConfig `yaml:"postgresql"`
}

// GormConfig gorm 配置
type GormConfig struct {
	QueryFields   bool      `yaml:"query-fields"`   // 是否开启全字段匹配查询
	SingularTable bool      `yaml:"singular-table"` // 是否关闭数据库表复数s
	UseDatabase   string    `yaml:"use-database"`   // 数据库选择
	GenConfig     GenConfig `yaml:"gen"`            // 数据库配置
	Database      Database  `yaml:"database"`       // 数据库配置
}

// GenConfig gen配置
type GenConfig struct {
	OutPath string   `yaml:"out-path"` // 生成模型接口文件路径
	mode    []string `yaml:"mode"`     // 生成模式
}

// GinConfig gin配置参数
type GinConfig struct {
	Port           string `yaml:"port"`           // gin端口号
	ReadTimeout    int    `yaml:"readTimeout"`    // 读取超时时间
	WriteTimeout   int    `yaml:"writeTimeout"`   // 写入超时时间
	MaxHeaderBytes int    `yaml:"maxHeaderBytes"` // 请求头的最大字节数
}

// DatabaseConfig 数据库基本配置
type DatabaseConfig struct {
	Host   string `yaml:"host"`    // 服务器地址
	DbName string `yaml:"db-name"` // 数据库名
	User   string `yaml:"user"`    // 用户名
	Pass   string `yaml:"pass"`    // 密码
	Port   int    `yaml:"port"`    // 端口
}

// RedisConfig redis配置
type RedisConfig struct {
	UsePub bool            `yaml:"use-pub"` // 是否配置redis订阅功能
	Pool   RedisPoolConfig `yaml:"pool"`    // 连接池配置
}

// RedisPoolConfig redisPool 连接池连接配置
type RedisPoolConfig struct {
	Network         string `yaml:"network"`            // 连接协议
	Server          string `yaml:"server"`             // 服务地址
	Port            int    `yaml:"port"`               // 端口号
	Pass            string `yaml:"pass"`               // 密码
	DbName          int    `yaml:"dbname"`             // 数据库
	ConnectTimeout  int    `yaml:"connect-timeout"`    // 连接超时时间
	ReadTimeout     int    `yaml:"read-timeout"`       // 读取超时时间
	WriteTimeout    int    `yaml:"Write-timeout"`      // 写入超时时间
	MaxOpenConn     int    `yaml:"max-open-conn"`      // 设置最大连接数
	ConnMaxIdleTime int    `yaml:"conn-max-idle-time"` // 空闲连接数
	MaxIdleConn     int    `yaml:"max-idle-conn"`      // 最大空闲连接数
	Wait            bool   `yaml:"wait"`               // 如果超时最大连接数是否等待
}

// LogConfig 日志配置
type LogConfig struct {
	UseLog   string      `yaml:"use-log"`  // 使用哪个日志
	Instance LogInstance `yaml:"instance"` // 日志实例
}

// LogInstance 日志实例
type LogInstance struct {
	Zap ZapConfig `yaml:"zap"` // zap 日志
}

// ZapConfig zap 日志配置
type ZapConfig struct {
	Mode        string   `yaml:"mode"`          // 日志模式，自定义模式custom，开发模式Devel ，生产模式product，
	Level       string   `yaml:"level"`         // 日志级别Debug,Info,WarnLevel,Set,DPanic,Panic,Fatal,
	Encoding    string   `yaml:"encoding"`      // 日志级输出格式# 输出格式json，控制台 console
	FormatTime  string   `yaml:"format-time"`   // 选择格式化日期
	FileLength  string   `yaml:"file-length"`   // 文件地址类型
	OutPath     []string `yaml:"out-path"`      // 日志输出路径
	MaxSize     int64    `yaml:"max-size"`      // 文件分割大小单位字节
	MaxBackups  int      `yaml:"max-backups"`   // 保留旧文件的最大个数
	MaxAge      int      `yaml:"max-age"`       // 保留旧文件的最大天数
	IsLocalTime bool     `yaml:"is-local-time"` // 是否使用本地时间
	IsCompress  bool     `yaml:"is-compress"`   // 是否压缩归档旧文件
}

// once 用于初始化config变量，并保证只初始化一次
var _once sync.Once

// 定义全局变量config，并初始化为nil
var cfg *Config

func init() {
	if cfg == nil {
		_once.Do(
			func() {
				path := GetConfigPath()
				configByte, err := os.ReadFile(path)
				if err != nil {
					log.Print("config init err:", err)
				}
				if err = yaml.Unmarshal(configByte, &cfg); err != nil {
					log.Print("config  unmarshal err：", err)
				}
			},
		)
	} else {
		log.Print("config  already exists！！！")
	}
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Print("get config Caller path err！！！")
	}
	// 获取绝对路径
	absolutePath, err := filepath.Abs(file)
	if err != nil {
		log.Print("get config abs path err:", err)
	}

	// 修改为配置文件路径并返回
	return strings.ReplaceAll(absolutePath, ".go", ".yml")
}

// GetGorm  获取gorm 配置
func GetGorm() *GormConfig {
	return &cfg.Gorm
}

// GetGen  获取gen 配置
func GetGen() *GenConfig {
	return &cfg.Gorm.GenConfig
}

// GetMysql  获取mysql 配置
func GetMysql() *DatabaseConfig {
	return &cfg.Gorm.Database.Mysql
}

// GetPostgres  获取postgres 配置
func GetPostgres() *DatabaseConfig {
	return &cfg.Gorm.Database.Postgresql
}

// GetGin  获取gin 配置
func GetGin() *GinConfig {
	return &cfg.Gin
}

// GetRedis  获取redis 配置
func GetRedis() *RedisConfig {
	return &cfg.Redis
}

// GetElasticsearch  获取elasticsearch 配置
func GetElasticsearch() *ElasticsearchConfig {
	return &cfg.Elasticsearch
}

// GetLog 获取日志配置
func GetLog() *LogConfig {
	return &cfg.Logs
}

// GetZap 获取zap 日志配置
func GetZap() *ZapConfig {
	return &cfg.Logs.Instance.Zap
}
