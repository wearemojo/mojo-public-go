package audit

import (
	"context"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/actor"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

type Record struct {
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	Actor     actor.Actor `json:"actor" bson:"actor"`
	Reason    Reason      `json:"reason" bson:"reason"`
}

type Reason struct {
	Code string         `json:"code" bson:"code"`
	Meta map[string]any `json:"meta" bson:"meta"`
}

func New(ctx context.Context, code string, meta map[string]any) (*Record, error) {
	actor := actor.GetActor(ctx)
	if actor == nil {
		return nil, merr.New("actor_not_found", nil)
	}

	return &Record{
		CreatedAt: time.Now(),
		Actor:     *actor,
		Reason: Reason{
			Code: code,
			Meta: meta,
		},
	}, nil
}
