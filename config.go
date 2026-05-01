package srvhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	commonv1 "github.com/mcrgnt/proto/gen/go/common/v1"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Config[T http.Handler] struct {
	Label commonv1.Label
	Host  commonv1.Host
	Port  commonv1.Port
}

func (cfg *Config[T]) Build() (any, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Value, cfg.Port.Value))
	if err != nil {
		return nil, err
	}
	return &srv[T]{
		initFn: func(t *srv[T]) {
			mdlw := middleware.New(middleware.Config{
				Recorder: t.recorder,
				Service:  cfg.Label.String(),
			})
			t.Handler = std.Handler("", mdlw, t.Handler)

			t.Handler = otelhttp.NewHandler(t.Handler, cfg.Label.String())

			t.BaseContext = func(net.Listener) context.Context {
				logger := slog.Default().With("srv", cfg.Label.String())
				return context.WithValue(context.Background(), "srvhttp", logger)
			}
		},
		Server:   http.Server{},
		listener: listener,
	}, nil
}
