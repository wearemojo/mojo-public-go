package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func (c Collection) SetupIndexes(models []mongo.IndexModel) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.Indexes().CreateMany(ctx, models)
	return
}

func (c Collection) FindAll(ctx context.Context, filter, results any, opts ...options.Lister[options.FindOptions]) (err error) {
	cur, err := c.Find(ctx, filter, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}
	return
}
