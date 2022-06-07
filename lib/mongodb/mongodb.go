package mongodb

import (
	"context"
	"io/fs"

	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/db/mongodb"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/secret"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDB struct {
	uriSecretID string
}

func New(uriSecretID string) *MongoDB {
	return &MongoDB{uriSecretID: uriSecretID}
}

func (m *MongoDB) Connect(ctx context.Context, schemaFS fs.FS, collectionNames []string) (*mongodb.Database, error) {
	uri, err := secret.Get(ctx, m.uriSecretID)
	if err != nil {
		return nil, err
	}

	// TODO: handle reconnection in some way?
	// in case the credentials change since the initial connection
	db, err := config.MongoDB{URI: uri}.Connect()
	if err != nil {
		return nil, err
	}

	err = setupCollections(ctx, db, collectionNames)
	if err != nil {
		return nil, err
	}

	err = db.SetupSchemas(ctx, schemaFS, collectionNames)
	if err != nil {
		return nil, merr.Wrap(ctx, err, "schema_setup_failed", merr.M{"collection_names": collectionNames})
	}

	return db, nil
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
