// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package redis

import (
	"errors"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	redigo "github.com/gomodule/redigo/redis"
)

var (
	pool    *redigo.Pool
	cfg     config.Redis
	enabled bool
)

// Init initializes the global redis pool when redis is enabled in config.
func Init(redisCfg config.Redis) {
	if !redisCfg.Enable {
		enabled = false
		pool = nil
		cfg = redisCfg
		return
	}

	redisCfg = redisCfg.SetDefault()
	p := &redigo.Pool{
		MaxIdle:     redisCfg.MaxIdle,
		MaxActive:   redisCfg.MaxActive,
		IdleTimeout: redisCfg.IdleTimeout,
		Wait:        redisCfg.Wait,
		Dial: func() (redigo.Conn, error) {
			return dial(redisCfg)
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	conn := p.Get()
	if err := conn.Err(); err != nil {
		_ = conn.Close()
		logger.Panicf("%s: %v", language.Get("initialize redis connection"), err)
	}
	if _, err := conn.Do("PING"); err != nil {
		_ = conn.Close()
		logger.Panicf("%s: %v", language.Get("initialize redis connection"), err)
	}
	_ = conn.Close()

	cfg = redisCfg
	pool = p
	enabled = true
}

func dial(redisCfg config.Redis) (redigo.Conn, error) {
	options := []redigo.DialOption{
		redigo.DialConnectTimeout(redisCfg.ConnectTimeout),
		redigo.DialReadTimeout(redisCfg.ReadTimeout),
		redigo.DialWriteTimeout(redisCfg.WriteTimeout),
	}
	if redisCfg.Password != "" {
		options = append(options, redigo.DialPassword(redisCfg.Password))
	}
	if redisCfg.DB != 0 {
		options = append(options, redigo.DialDatabase(redisCfg.DB))
	}
	if redisCfg.TLS {
		options = append(options, redigo.DialUseTLS(true))
	}
	if redisCfg.Dsn != "" {
		return redigo.DialURL(redisCfg.Dsn, options...)
	}
	return redigo.Dial(redisCfg.Network, redisCfg.Addr, options...)
}

// Enabled reports whether redis is enabled and initialized.
func Enabled() bool {
	return enabled
}

// Config returns the initialized redis config.
func Config() config.Redis {
	return cfg
}

// Pool returns the global redis pool.
func Pool() (*redigo.Pool, error) {
	if pool == nil {
		return nil, errors.New("redis is not initialized")
	}
	return pool, nil
}

// MustPool returns the global redis pool and panics when it is not initialized.
func MustPool() *redigo.Pool {
	p, err := Pool()
	if err != nil {
		panic(err)
	}
	return p
}

// Conn returns a redis connection from the global pool.
func Conn() (redigo.Conn, error) {
	p, err := Pool()
	if err != nil {
		return nil, err
	}
	conn := p.Get()
	if err := conn.Err(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

// Do executes a redis command with a connection from the global pool.
func Do(commandName string, args ...interface{}) (interface{}, error) {
	conn, err := Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// Close closes the global redis pool.
func Close() error {
	if pool == nil {
		return nil
	}
	err := pool.Close()
	pool = nil
	enabled = false
	return err
}
