package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	srvhttp "github.com/mcrgnt/srv-http"
)

func main() {
	// var cfg = &srvhttp.Config{
	// 	Port: commonv1.Port{Value: 8080},
	// }
	// var srv, err = cfg.Build()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(">")

	// if err := srv.(interface {
	// 	Start(context.Context) error
	// }).Start(context.Background()); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(">>")
	// time.Sleep(time.Second * 5)
	// fmt.Println(">>>")

	// if err := srv.(interface{ Close(context.Context) error }).Close(context.Background()); err != nil {
	// 	panic(err)
	// }
	some()
}

type RootConfig struct {
	srvHTTP *srvhttp.Config2[somer]
}

func some() {
	var cfg = new(RootConfig)
	cfg.srvHTTP = new(srvhttp.Config2[somer])
	srv, err := cfg.srvHTTP.Build()
	if err != nil {
		panic(err)
	}

	_ = srv.(interface{ Deps() []any }).Deps()
	srv.(interface{ Inject([]any) }).Inject([]any{newRouter()})

	if err := srv.(interface{ Start(context.Context) error }).Start(context.Background()); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 50)
	if err := srv.(interface {
		Close(ctx context.Context) error
	}).Close(context.Background()); err != nil {
		panic(err)
	}
}

type somer interface {
	http.Handler
	Some()
}

type router struct {
	Mux *chi.Mux
}

func newRouter() *router {
	fmt.Println("called")
	var mux = chi.NewRouter()
	mux.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", *r)
		w.Write([]byte("ok"))
	}))
	return &router{
		Mux: mux,
	}
}

func (t *router) Some() {}

func (t *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%#v\n", *r)
	w.Write([]byte("ok"))
}
