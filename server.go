package srvhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	commonv1 "github.com/mcrgnt/proto/gen/go/common/v1"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/atomic"
)

type Config struct {
	Label commonv1.Label
	Host  commonv1.Host
	Port  commonv1.Port
}

func (cfg *Config) Build() (any, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Value, cfg.Port.Value))
	if err != nil {
		return nil, err
	}
	return &srv{
		config:   cfg,
		Server:   http.Server{},
		listener: listener,
	}, nil
}

type srv struct {
	config *Config
	http.Server
	handler  http.Handler `deps:""`
	listener net.Listener
	err      atomic.Error
}

func (t *srv) Label() string {
	return t.config.Label.String()
}

func (t *srv) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if t.Handler != nil {
			t.Handler = otelhttp.NewHandler(t.Handler, t.config.Label.String())
		}
		t.BaseContext = func(net.Listener) context.Context {
			return ctx
		}

		go func() {
			if err := t.Serve(t.listener); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					t.err.Store(err)
				}
			}
		}()
		return nil
	}
}

func (t *srv) Close(ctx context.Context) error {
	if err := t.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (t *srv) HealthCheck(context.Context) error {
	var err = t.err.Load()
	if err != nil {
		return err
	}
	return nil
}
