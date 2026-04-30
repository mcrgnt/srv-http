package srvhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	commonv1 "github.com/mcrgnt/proto/gen/go/common/v1"
	"go.uber.org/atomic"
)

type Config2[T http.Handler] struct {
	Label commonv1.Label
	Host  commonv1.Host
	Port  commonv1.Port
}

func (cfg *Config2[T]) Build() (any, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Value, cfg.Port.Value))
	if err != nil {
		return nil, err
	}
	fmt.Println(listener.Addr().String())
	return &srv2[T]{
		config:   cfg,
		Server:   http.Server{},
		listener: listener,
	}, nil
}

type srv2[T http.Handler] struct {
	config *Config2[T]
	http.Server
	handler  T `deps:""`
	listener net.Listener
	err      atomic.Error
}

func (t *srv2[T]) Label() string {
	return t.config.Label.String()
}

func (t *srv2[T]) Start(ctx context.Context) error {
	fmt.Printf("%T\n", t.Handler)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// if t.handler != nil {
		// 	t.handler = otelhttp.NewHandler(t.handler, t.config.Label.String())
		// }
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

func (t *srv2[T]) Close(ctx context.Context) error {
	if err := t.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (t *srv2[T]) HealthCheck(context.Context) error {
	var err = t.err.Load()
	if err != nil {
		return err
	}
	return nil
}
