package service

import (
	"context"
	"errors"
	"time"

	"my-base/client/configs"
	"my-base/client/transport"
)

var ErrClientDeleted = transport.ErrClientDeleted

// Runtime owns the client control connection and reconnect policy.
type Runtime struct {
	Config *configs.ClientConfig
	Client *transport.Client
}

func NewRuntime(cfg *configs.Config) *Runtime {
	clientCfg := cfg.GetClientConfig()
	client := transport.NewClient(clientCfg.GetPublicJobAddr())
	client.SetConnOptions(transport.ConnOptions{
		ReadTimeout:     time.Duration(clientCfg.GetPingInterval()+clientCfg.GetPongTimeout()) * time.Second,
		WriteTimeout:    time.Duration(clientCfg.GetPongTimeout()) * time.Second,
		KeepAlive:       true,
		KeepAlivePeriod: 30 * time.Second,
	})
	client.SetOnKeyFunc(cfg.SetSymbol).SetKeyFunc(cfg.GetSymbol)
	return &Runtime{Config: clientCfg, Client: client}
}

func (r *Runtime) Run(ctx context.Context) error {
	delay := time.Duration(r.Config.GetReconnectBaseDelay()) * time.Second
	maxDelay := time.Duration(r.Config.GetReconnectMaxDelay()) * time.Second
	for {
		started := time.Now()
		err := r.Client.StartWithHeartbeat(
			time.Duration(r.Config.GetPingInterval())*time.Second,
			time.Duration(r.Config.GetPongTimeout())*time.Second,
			r.Config.GetMaxPingFailures(),
		)
		if err == nil || errors.Is(err, transport.ErrClientDeleted) {
			if errors.Is(err, transport.ErrClientDeleted) {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if time.Since(started) >= time.Duration(r.Config.GetPingInterval())*time.Second {
			delay = time.Duration(r.Config.GetReconnectBaseDelay()) * time.Second
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil
		}
		if delay < maxDelay {
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}
}

func (r *Runtime) Close() error { return r.Client.Close() }
