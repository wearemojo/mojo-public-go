package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	*mongo.Database
}

func (db Database) Collection(name string, opts ...options.Lister[options.CollectionOptions]) *Collection {
	return &Collection{db.Database.Collection(name, opts...)}
}

func (db Database) SetupSchemas(ctx context.Context, schemaFS fs.FS, collectionNames []string) error {
	// Keep a map of all used collection names to make sure every single collection has a schema defined
	usedCollectionNames := make(map[string]struct{})
	// walk through the filesystem to make sure every single schema defined is being used
	if err := fs.WalkDir(schemaFS, ".", func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		schemaFileName := entry.Name()

		colName, ok := slicefn.Find(collectionNames, func(colName string) bool {
			return fmt.Sprintf("%s.json", colName) == schemaFileName
		})
		if !ok {
			return merr.New(ctx, "schema_defined_but_not_used", merr.M{"schema": schemaFileName}, nil)
		}

		file, err := schemaFS.Open(schemaFileName)
		if err != nil {
			return err
		}

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		var schema any
		if err := json.Unmarshal(data, &schema); err != nil {
			return err
		}

		if err := db.RunCommand(ctx, bson.D{
			{"collMod", colName},
			{"validationLevel", "strict"},
			{"validationAction", "error"},
			{"validator", bson.M{
				"$jsonSchema": schema,
			}},
		}).Err(); err != nil {
			return err
		}

		// mark collection as used
		usedCollectionNames[colName] = struct{}{}
		return nil
	}); err != nil {
		return err
	}

	for _, colName := range collectionNames {
		if _, ok := usedCollectionNames[colName]; !ok {
			return merr.New(ctx, "collection_defined_but_no_matching_schema_defined", merr.M{"collection": colName}, nil)
		}
	}

	return nil
}

func (db *Database) DoTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.DoTxWithOptions(ctx, options.Session(), fn)
}

func (db *Database) DoTxWithOptions(ctx context.Context, opts *options.SessionOptionsBuilder, fn func(ctx context.Context) error) error {
	return db.Client().UseSessionWithOptions(ctx, opts, func(ctx context.Context) error {
		sess := mongo.SessionFromContext(ctx)
		_, err := sess.WithTransaction(ctx, func(ctx context.Context) (any, error) {
			return nil, fn(ctx)
		})
		return err
	})
}
