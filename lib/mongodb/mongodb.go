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
	"github.com/wearemojo/mojo-public-go/lib/secret"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
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
	db, err := connect(ctx, uri)
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

	for idx := range numFields {
		field := val.Field(idx)
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

func connect(ctx context.Context, uri string) (*mongodb.Database, error) {
	opts, dbName, err := config.MongoDB{URI: uri}.Options(ctx)
	if err != nil {
		return nil, err
	}

	opts.Monitor = otelmongo.NewMonitor()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
