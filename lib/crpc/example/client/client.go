package main

import (
	"context"
	"net/http"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/crpc"
	"github.com/wearemojo/mojo-public-go/lib/crpc/example"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

type ExampleClient struct {
	*crpc.Client
}

func (ec *ExampleClient) Ping(ctx context.Context) error {
	return ec.Do(ctx, "ping", "2017-11-08", nil, nil)
}

func (ec *ExampleClient) Greet(ctx context.Context, req *example.GreetRequest) (res *example.GreetResponse, err error) {
	return res, ec.Do(ctx, "greet", "2017-11-08", req, &res)
}

func main() {
	var client example.Service = &ExampleClient{
		Client: crpc.NewClient(context.Background(), "http://127.0.0.1:3000/v1", &http.Client{
			Timeout: 5 * time.Second,
		}),
	}

	ctx := context.Background()

	if err := client.Ping(ctx); err != nil {
		mlog.Warn(ctx, merr.New(ctx, "ping_failed", nil, err))
		return
	}

	res, err := client.Greet(ctx, &example.GreetRequest{Name: "James"})
	if err != nil {
		mlog.Warn(ctx, merr.New(ctx, "greet_failed", nil, err))
		return
	}

	mlog.Info(ctx, merr.New(ctx, "greeting_received", merr.M{"greeting": res.Greeting}))
}
