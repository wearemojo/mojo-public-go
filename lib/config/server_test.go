package config

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestServerServeReturnsWhenContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	cfg := Server{
		Addr:     "127.0.0.1:0",
		Graceful: 1,
	}
	srv := &http.Server{
		Handler:           http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	listener, err := cfg.Listen(ctx, srv)
	if err != nil {
		t.Fatal(err)
	}

	started := make(chan struct{})
	errs := make(chan error, 1)
	go func() {
		errs <- cfg.Serve(ctx, srv, listener, func(context.Context) error {
			close(started)
			return nil
		})
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("server did not start")
	}

	cancel()

	select {
	case err := <-errs:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not stop after context cancellation")
	}
}
