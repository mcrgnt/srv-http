package chttp

import (
	"net/http"

	srvhttp "github.com/mcrgnt/srv-http"
)

type Config struct {
	srvhttp.Config
}

func (t *Config) Build() (any, error) {
	return &srv{}, nil
}

type some interface {
	http.Handler
	Some()
}

type srv struct {
	handler some
}
