package audit

import (
	"context"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/actor"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type Record struct {
	Timestamp time.Time   `bson:"timestamp" json:"timestamp"`
	Actor     actor.Actor `bson:"actor"     json:"actor"`
	Reason    Reason      `bson:"reason"    json:"reason"`
}

type Reason struct {
	Code string         `bson:"code" json:"code"`
	Meta map[string]any `bson:"meta" json:"meta"`
}

func New(ctx context.Context, code string, meta map[string]any) (*Record, error) {
	actor := actor.GetActor(ctx)
	if actor == nil {
		return nil, merr.New(ctx, "actor_not_found", nil)
	}

	return &Record{
		Timestamp: time.Now(),
		Actor:     *actor,
		Reason: Reason{
			Code: code,
			Meta: meta,
		},
	}, nil
}
