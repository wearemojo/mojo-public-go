package config

import (
	"context"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/db/mongodb"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

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
		return opts, dbName, err
	}

	// all Go services use majority reads/writes, and this is unlikely to change
	// if it does change, switch to accepting as an argument
	opts.SetReadConcern(readconcern.Majority())
	opts.SetWriteConcern(writeconcern.Majority())

	cs, err := connstring.Parse(m.URI)
	if err != nil {
		return opts, dbName, err
	}

	dbName = cs.Database
	if dbName == "" {
		err = merr.New(ctx, "mongo_db_name_missing", nil)
	}

	return opts, dbName, err
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
