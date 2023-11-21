package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/config"
	"github.com/wearemojo/mojo-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/crpc/example"
	"github.com/wearemojo/mojo-public-go/lib/middleware/request"
)

type ServerConfig struct {
	Logging clog.Config `json:"logging"`

	Server config.Server `json:"server"`
}

type ExampleServer struct{}

func (es *ExampleServer) Ping(ctx context.Context) error {
	return nil
}

func (es *ExampleServer) Greet(ctx context.Context, req *example.GreetRequest) (*example.GreetResponse, error) {
	clog.Get(ctx).Info("just an example")

	return &example.GreetResponse{
		Greeting: fmt.Sprintf("Hello %s!", req.Name),
	}, nil
}

func main() {
	cfg := &ServerConfig{
		Logging: clog.Config{
			Format: "text",
			Debug:  true,
		},

		Server: config.Server{
			Addr: "127.0.0.1:3000",
		},
	}

	log := cfg.Logging.Configure(context.Background())

	var svc example.Service = &ExampleServer{}

	// create a new RPC server
	rpc := crpc.NewServer(unsafeNoAuthentication)

	// add logging middleware
	rpc.Use(crpc.Logger())

	// register Ping and Greet (version 2017-11-08)
	rpc.Register("ping", "2017-11-08", nil, svc.Ping)
	rpc.Register("greet", "2017-11-08", example.GreetRequestSchema, svc.Greet)

	mux := chi.NewRouter()
	mux.Use(request.Logger(log))
	mux.With(request.StripPrefix("/v1")).Handle("/v1/*", rpc)

	httpServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.WithField("addr", cfg.Server.Addr).Info("listening")

	if err := cfg.Server.ListenAndServe(httpServer); err != nil {
		log.WithError(err).Fatal("listen failed")
	}
}

func unsafeNoAuthentication(next crpc.HandlerFunc) crpc.HandlerFunc {
	return func(res http.ResponseWriter, req *crpc.Request) error {
		return next(res, req)
	}
}
