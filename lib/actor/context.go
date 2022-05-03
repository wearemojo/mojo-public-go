package actor

import (
	"context"
	"errors"

	"github.com/wearemojo/mojo-public-go/lib/authparsing"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/ptr"
)

type contextKey string

const contextKeyActor contextKey = "actor"

const ErrActorNotFound = merr.Code("actor_not_found")

func getActor(ctx context.Context) *Actor {
	actor, _ := ctx.Value(contextKeyActor).(*Actor)
	return actor
}

func SetActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, contextKeyActor, &actor)
}

func GetActor(ctx context.Context) (*Actor, error) {
	if ctxActor := getActor(ctx); ctxActor != nil {
		return ctxActor, nil
	}

	authstate := authparsing.GetAuthState(ctx)

	if authstate == nil {
		return nil, merr.New(ErrActorNotFound, nil)
	}

	if a, ok := authstate.(Actorer); ok {
		if actor := a.Actor(ctx); actor != nil {
			return actor, nil
		}
	}

	return nil, merr.New(ErrActorNotFound, nil)
}

func GetActorOrUnknown(ctx context.Context) (*Actor, error) {
	actor, err := GetActor(ctx)
	if !errors.Is(err, ErrActorNotFound) {
		return actor, err
	}

	return ptr.P(NewUnknown()), nil
}
