package mongodb

import (
	"context"
	"fmt"
	"io/fs"
	"reflect"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/config"
	"github.com/wearemojo/mojo-public-go/lib/db/mongodb"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"github.com/wearemojo/mojo-public-go/lib/secret"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

const (
	connectTimeout    = 45 * time.Second
	connectAttempts   = 3
	connectRetryDelay = 5 * time.Second
)

type MongoDB struct {
	uriSecretID string
}

func New(uriSecretID string) *MongoDB {
	return &MongoDB{uriSecretID: uriSecretID}
}

func (m *MongoDB) Connect(ctx context.Context, dbStruct any, schemaFS fs.FS) (*mongodb.Database, error) {
	// dbStruct looks like: struct{ Blah *BlahCollection `mongocol:"blah" }{}
	collectionNames := extractCollectionNames(dbStruct)

	uri, err := secret.Get(ctx, m.uriSecretID)
	if err != nil {
		return nil, err
	}

	// TODO: handle reconnection in some way?
	// in case the credentials change since the initial connection
	db, err := connectWithRetries(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = setupCollections(ctx, db, collectionNames)
	if err != nil {
		return nil, err
	}

	err = db.SetupSchemas(ctx, schemaFS, collectionNames)
	if err != nil {
		return nil, merr.New(ctx, "schema_setup_failed", merr.M{"collection_names": collectionNames}, err)
	}

	return db, nil
}

func extractCollectionNames(dbStruct any) []string {
	val := reflect.ValueOf(dbStruct).Type()
	if val.Kind() != reflect.Struct {
		panic("dbStruct must be a struct")
	}

	numFields := val.NumField()
	collectionNames := make([]string, 0, numFields)

	for field := range val.Fields() {
		if !field.IsExported() {
			continue
		}

		name, ok := field.Tag.Lookup("mongocol")
		if !ok {
			panic(fmt.Sprintf("missing mongocol tag on field %s", field.Name))
		}

		collectionNames = append(collectionNames, name)
	}

	return collectionNames
}

// connectWithRetries exists because services connect at instance startup,
// where a failure crashes the whole instance. When many instances boot at
// once (e.g. a deploy rolling every service), MongoDB can be slow to complete
// handshakes, so a single attempt with a short timeout would fail deploys on
// transient slowness. Cloud Run allows several minutes of startup time, so
// retrying is strictly better than exiting on the first slow connection.
func connectWithRetries(ctx context.Context, uri string) (*mongodb.Database, error) {
	var db *mongodb.Database
	var err error

	for attempt := 1; attempt <= connectAttempts; attempt++ {
		db, err = connect(ctx, uri)
		if err == nil {
			return db, nil
		}

		if attempt == connectAttempts {
			break
		}

		mlog.Warn(ctx, merr.New(ctx, "mongo_connect_retrying", merr.M{"attempt": attempt}, err))

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(connectRetryDelay):
		}
	}

	return nil, err
}

func connect(ctx context.Context, uri string) (*mongodb.Database, error) {
	opts, dbName, err := config.MongoDB{URI: uri}.Options(ctx)
	if err != nil {
		return nil, err
	}

	opts.Monitor = otelmongo.NewMonitor()

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	return mongodb.Connect(ctx, opts, dbName)
}

func setupCollections(ctx context.Context, db *mongodb.Database, names []string) error {
	existingNames, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}

outer:
	for _, name := range names {
		for _, existingName := range existingNames {
			if name == existingName {
				continue outer
			}
		}

		err = db.CreateCollection(ctx, name)
		if err != nil {
			return err
		}
	}

	return nil
}
