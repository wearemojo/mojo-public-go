package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

func (db Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{db.Database.Collection(name, opts...)}
}

func (db Database) SetupSchemas(ctx context.Context, fs fs.FS, collectionNames []string) error {
	for _, colName := range collectionNames {
		file, err := fs.Open(fmt.Sprintf("%s.json", colName))
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
	}

	return nil
}

func (db *Database) DoTx(ctx context.Context, fn func(ctx mongo.SessionContext) error) error {
	return db.DoTxWithOptions(ctx, options.Session(), fn)
}

func (db *Database) DoTxWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(ctx mongo.SessionContext) error) error {
	return db.Client().UseSessionWithOptions(ctx, opts, func(ctx mongo.SessionContext) error {
		_, err := ctx.WithTransaction(ctx, func(ctx mongo.SessionContext) (any, error) {
			return nil, fn(ctx)
		})
		return err
	})
}
