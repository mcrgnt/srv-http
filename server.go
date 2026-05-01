package srvhttp

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/slok/go-http-metrics/metrics"
	"go.uber.org/atomic"
)

type srv[T http.Handler] struct {
	initFn func(t *srv[T])
	http.Server
	handler  T `deps:""`
	listener net.Listener
	err      atomic.Error
	recorder metrics.Recorder
}

func (r *srv[T]) Deps() []any {
	return []any{
		(*T)(nil),
		(*metrics.Recorder)(nil),
	}
}

func (r *srv[T]) Inject(args []any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case T:
			r.Handler = v
		case metrics.Recorder:
			r.recorder = v
		}
	}
}

func (t *srv[T]) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		t.initFn(t)

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

func (t *srv[T]) Close(ctx context.Context) error {
	if err := t.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (t *srv[T]) HealthCheck(_ context.Context) error {
	return t.err.Load()
}
