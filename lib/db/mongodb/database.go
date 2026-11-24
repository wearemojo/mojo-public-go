package mongodb

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"time"

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
			return merr.New(ctx, "schema_defined_but_not_used", merr.M{"schema": schemaFileName})
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
			return merr.New(ctx, "collection_defined_but_no_matching_schema_defined", merr.M{"collection": colName})
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

const transactionAbortTimeout = time.Second

// DoTxWithCommitDeadline uses ctx for the commit because the driver's
// WithTransaction helper detaches its deadline. fn may run twice after a
// transient error. Cancellation after sending a commit can leave its outcome unknown.
func (db *Database) DoTxWithCommitDeadline(ctx context.Context, fn func(ctx context.Context) error) error {
	err := db.doTxWithCommitDeadline(ctx, fn)
	if mongoErr, ok := errors.AsType[mongo.ServerError](err); ok && ctx.Err() == nil {
		if mongoErr.HasErrorLabel("TransientTransactionError") && !mongoErr.HasErrorLabel("UnknownTransactionCommitResult") {
			return db.doTxWithCommitDeadline(ctx, fn)
		}
	}
	return err
}

func (db *Database) doTxWithCommitDeadline(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}

	committing := false
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), transactionAbortTimeout)
		defer cancel()
		if committing {
			// An in-flight commit may have succeeded. Avoid a server-side abort.
			cancel()
		}
		session.EndSession(cleanupCtx)
	}()

	if err = session.StartTransaction(); err != nil {
		return err
	}
	ctx = mongo.NewSessionContext(ctx, session)
	if err = fn(ctx); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	committing = true
	return session.CommitTransaction(ctx)
}
