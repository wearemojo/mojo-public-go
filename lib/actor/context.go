package actor

import (
	"context"

	"github.com/wearemojo/mojo-public-go/lib/authparsing"
)

type contextKey string

const contextKeyActor contextKey = "actor"

func getActor(ctx context.Context) *Actor {
	if actor, ok := ctx.Value(contextKeyActor).(Actor); ok {
		return &actor
	}

	return nil
}

func SetActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, contextKeyActor, actor)
}

func GetActor(ctx context.Context) *Actor {
	if ctxActor := getActor(ctx); ctxActor != nil {
		return ctxActor
	}

	authstate := authparsing.GetAuthState(ctx)

	if authstate == nil {
		return nil
	}

	if a, ok := authstate.(Actorer); ok {
		if actor := a.Actor(ctx); actor != nil {
			return actor
		}
	}

	return nil
}

func GetActorOrUnknown(ctx context.Context) Actor {
	if actor := GetActor(ctx); actor != nil {
		return *actor
	}

	return NewUnknown()
}
