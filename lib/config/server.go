package config

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// DefaultGraceful is the graceful shutdown timeout applied when no
// configuration value is given.
const DefaultGraceful = 5

// Server configures the binding and security of an HTTP server.
type Server struct {
	Addr string `json:"addr"`

	// Graceful enables graceful shutdown and is the time in seconds to wait
	// for all outstanding requests to terminate before forceably killing the
	// server. When no value is given, DefaultGraceful is used. Graceful
	// shutdown is disabled when less than zero.
	Graceful int `json:"graceful"`
}

// Listen configures a HTTP server and binds a listener without serving it yet.
func (cfg *Server) Listen(ctx context.Context, srv *http.Server) (net.Listener, error) {
	addr := ":http"
	if srv.Addr != "" {
		addr = srv.Addr
	} else if cfg.Addr != "" {
		addr = cfg.Addr
	}

	var lc net.ListenConfig
	return lc.Listen(ctx, "tcp", addr)
}

// Serve begins serving an already-bound listener, runs afterStart once the
// server is accepting requests, and handles graceful shutdown.
func (cfg *Server) Serve(ctx context.Context, srv *http.Server, listener net.Listener, afterStart func(context.Context) error) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stop)

	errs := make(chan error, 1)

	go func() {
		err := srv.Serve(listener)
		if !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		return err
	default:
	}

	if afterStart != nil {
		if err := afterStart(ctx); err != nil {
			if shutdownErr := cfg.shutdown(ctx, srv); shutdownErr != nil {
				return merr.New(ctx, "shutdown_failed_following_after_start_error", nil, shutdownErr, err)
			}
			return err
		}
	}

	select {
	case err := <-errs:
		return err

	case <-stop:
		return cfg.shutdown(ctx, srv)
	}
}

func (cfg *Server) shutdown(ctx context.Context, srv *http.Server) error {
	graceful := DefaultGraceful
	if cfg.Graceful != 0 {
		graceful = cfg.Graceful
	}

	if graceful > 0 {
		tsCtx, cancel := context.WithTimeout(ctx, time.Duration(graceful)*time.Second)
		defer cancel()

		return srv.Shutdown(tsCtx)
	}

	return srv.Close()
}

// ListenAndServe configures a HTTP server and begins listening for clients.
func (cfg *Server) ListenAndServe(ctx context.Context, srv *http.Server, afterStart func(context.Context) error) error {
	listener, err := cfg.Listen(ctx, srv)
	if err != nil {
		return err
	}

	return cfg.Serve(ctx, srv, listener, afterStart)
}
