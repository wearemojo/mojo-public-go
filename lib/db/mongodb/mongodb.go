package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, opts *options.ClientOptions, dbName string) (db *Database, err error) {
	client, err := mongo.Connect(opts)
	if err != nil {
		return db, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return db, err
	}

	db = &Database{client.Database(dbName)}
	return db, err
}
