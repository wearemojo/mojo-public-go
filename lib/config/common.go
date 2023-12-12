package config

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wearemojo/mojo-public-go/lib/db/mongodb"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Redis configures a connection to a Redis database.
type Redis struct {
	URI          string        `json:"uri"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// Options returns a configured redis.Options structure.
func (r Redis) Options() (*redis.Options, error) {
	opts, err := redis.ParseURL(r.URI)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = r.DialTimeout
	opts.ReadTimeout = r.ReadTimeout
	opts.WriteTimeout = r.WriteTimeout

	return opts, nil
}

// Connect returns a connected redis.Client instance.
func (r Redis) Connect(ctx context.Context) (*redis.Client, error) {
	opts, err := r.Options()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return client, err
	}

	return client, nil
}

// MongoDB configures a connection to a Mongo database.
type MongoDB struct {
	URI             string         `json:"uri"`
	ConnectTimeout  time.Duration  `json:"connect_timeout"`
	MaxConnIdleTime *time.Duration `json:"max_conn_idle_time"`
	MaxConnecting   *uint64        `json:"max_connecting"`
	MaxPoolSize     *uint64        `json:"max_pool_size"`
	MinPoolSize     *uint64        `json:"min_pool_size"`
}

// Options returns the MongoDB client options and database name.
func (m MongoDB) Options(ctx context.Context) (opts *options.ClientOptions, dbName string, err error) {
	opts = options.Client().ApplyURI(m.URI)
	opts.MaxConnIdleTime = m.MaxConnIdleTime
	opts.MaxConnecting = m.MaxConnecting
	opts.MaxPoolSize = m.MaxPoolSize
	opts.MinPoolSize = m.MinPoolSize

	err = opts.Validate()
	if err != nil {
		return
	}

	// all Go services use majority reads/writes, and this is unlikely to change
	// if it does change, switch to accepting as an argument
	opts.SetReadConcern(readconcern.Majority())
	opts.SetWriteConcern(writeconcern.Majority())

	cs, err := connstring.Parse(m.URI)
	if err != nil {
		return
	}

	dbName = cs.Database
	if dbName == "" {
		err = merr.New(ctx, "mongo_db_name_missing", nil)
	}

	return
}

// Connect returns a connected mongo.Database instance.
func (m MongoDB) Connect(ctx context.Context) (*mongodb.Database, error) {
	opts, dbName, err := m.Options(ctx)
	if err != nil {
		return nil, err
	}

	if m.ConnectTimeout == 0 {
		m.ConnectTimeout = 10 * time.Second
	}

	// this package can only be used for service config
	// so can only happen at init-time - no need to accept context input
	ctx, cancel := context.WithTimeout(ctx, m.ConnectTimeout)
	defer cancel()

	return mongodb.Connect(ctx, opts, dbName)
}

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

// ListenAndServe configures a HTTP server and begins listening for clients.
func (cfg *Server) ListenAndServe(srv *http.Server) error {
	// only set listen address if none is already configured
	if srv.Addr == "" {
		srv.Addr = cfg.Addr
	}

	if cfg.Graceful == 0 {
		cfg.Graceful = DefaultGraceful
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	errs := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		return err

	case <-stop:
		if cfg.Graceful > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Graceful)*time.Second)
			defer cancel()

			return srv.Shutdown(ctx)
		}

		return nil
	}
}

func ContextWithCancelOnSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-stop:
		case <-ctx.Done():
		}
	}()

	return ctx
}
