package crawling

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 全局 Redis 客户端
var rdb *redis.Client

// Storage colly 存储实例
type Storage struct {
	prefix  string
	u       *url.URL
	Expires time.Duration
	mu      sync.RWMutex
}

// InitRedis 初始化 Redis 客户端（推荐在 main() 或初始化逻辑中调用）
func InitRedis(addr, password string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	return nil
}

// NewStorage 新建存储实例
func NewStorage(prefix string) *Storage {
	return &Storage{
		Expires: 5 * time.Minute, // 过期时间 5 分钟
		prefix:  prefix,
		mu:      sync.RWMutex{},
	}
}

// Init 实现初始化接口
func (s *Storage) Init() error {
	if rdb == nil {
		return errors.New("Redis 未初始化，请先调用 InitRedis()")
	}
	return nil
}

// Visited 标记请求已访问
func (s *Storage) Visited(requestID uint64) error {
	ctx := context.Background()
	key := s.getIDStr(requestID)
	err := rdb.Set(ctx, key, "1", s.Expires).Err()
	if err != nil {
		zap.L().Error("Visited Redis 写入失败", zap.Error(err))
	}
	return err
}

// IsVisited 判断请求是否访问过
func (s *Storage) IsVisited(requestID uint64) (bool, error) {
	ctx := context.Background()
	key := s.getIDStr(requestID)
	val, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		zap.L().Error("IsVisited Redis 读取失败", zap.Error(err))
		return false, err
	}
	return val == "1", nil
}

// SetCookies 设置 cookie
func (s *Storage) SetCookies(u *url.URL, cookies string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.u = u
	ctx := context.Background()
	key := s.getCookieID(u.Host)
	err := rdb.Set(ctx, key, cookies, 10*time.Minute).Err()
	if err != nil {
		zap.L().Error("SetCookies Redis 写入失败", zap.Error(err))
	}
}

// Cookies 获取 cookie
func (s *Storage) Cookies(u *url.URL) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.u = u
	ctx := context.Background()
	key := s.getCookieID(u.Host)
	cookiesStr, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		zap.L().Error("Cookies Redis 获取失败", zap.Error(err))
	}
	return cookiesStr
}

// AddRequest 添加请求信息
func (s *Storage) AddRequest(r []byte) error {
	ctx := context.Background()
	key := s.getQueueID()
	err := rdb.RPush(ctx, key, r).Err()
	if err != nil {
		zap.L().Error("AddRequest Redis 写入失败", zap.Error(err))
	}
	return err
}

// GetRequest 获取请求信息
func (s *Storage) GetRequest() ([]byte, error) {
	ctx := context.Background()
	key := s.getQueueID()
	val, err := rdb.LPop(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Error("GetRequest Redis 读取失败", zap.Error(err))
	}
	return val, err
}

// QueueSize 获取队列长度
func (s *Storage) QueueSize() (int, error) {
	ctx := context.Background()
	key := s.getQueueID()
	size, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		zap.L().Error("QueueSize Redis 查询失败", zap.Error(err))
	}
	return int(size), err
}

// 内部辅助方法
func (s *Storage) getIDStr(ID uint64) string {
	return fmt.Sprintf("%s:request:%d", s.getPrefix(), ID)
}

func (s *Storage) getCookieID(c string) string {
	return fmt.Sprintf("%s:cookie:%s", s.getPrefix(), c)
}

func (s *Storage) getQueueID() string {
	return fmt.Sprintf("%s:queue", s.getPrefix())
}

func (s *Storage) getPrefix() string {
	return s.prefix
}
