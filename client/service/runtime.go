package service

import (
	"context"
	"errors"
	"time"

	"my-base/client/configs"
	"my-base/client/transport"
)

var ErrClientDeleted = transport.ErrClientDeleted
var ErrClientRegistrationRejected = transport.ErrClientRegistrationRejected

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
	// StartWithHeartbeat blocks while reading control messages. Closing the
	// connection on cancellation wakes that read so signal-driven shutdown does
	// not have to wait for a network timeout or peer activity.
	stopClose := make(chan struct{})
	defer close(stopClose)
	go func() {
		select {
		case <-ctx.Done():
			_ = r.Client.Close()
		case <-stopClose:
		}
	}()

	delay := time.Duration(r.Config.GetReconnectBaseDelay()) * time.Second
	maxDelay := time.Duration(r.Config.GetReconnectMaxDelay()) * time.Second
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		started := time.Now()
		err := r.Client.StartWithHeartbeat(
			time.Duration(r.Config.GetPingInterval())*time.Second,
			time.Duration(r.Config.GetPongTimeout())*time.Second,
			r.Config.GetMaxPingFailures(),
		)
		if errors.Is(err, transport.ErrClientDeleted) || errors.Is(err, transport.ErrClientRegistrationRejected) {
			return err
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
