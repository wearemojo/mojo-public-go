package actor

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/cuvva/cuvva-public-go/lib/version"
	"github.com/wearemojo/mojo-public-go/lib/stacktrace"
)

type Actorer interface {
	Actor(context.Context) *Actor
}

type Type string

var (
	TypeUnknown  Type = "unknown"
	TypeInternal Type = "internal"
	TypeService  Type = "service"
	TypeUser     Type = "user"
)

type Actor struct {
	Type   Type              `json:"type" bson:"type"`
	Params map[string]string `json:"params" bson:"params"`
}

func NewUnknown(params map[string]string) Actor {
	if params == nil {
		params = map[string]string{}
	}

	return Actor{
		Type:   TypeUnknown,
		Params: params,
	}
}

// NewInternal represents a decision made within your service. Where as service
// is used to mark a decision made from another service calling in.
//
// The argument skip is the number of stack frames to skip before identifying
// the frame to use, with 0 identifying the frame for NewInternal itself and 1
// identifying the caller of NewInternal.
func NewInternal(skip int, codePathFallback string) Actor {
	return Actor{
		Type: TypeInternal,
		Params: map[string]string{
			"revision":  version.Revision,
			"code_path": stacktrace.GetCallerCodePath(skip+1, codePathFallback),
		},
	}
}

func NewService(env, service string) Actor {
	return Actor{
		Type: TypeService,
		Params: map[string]string{
			"env":     env,
			"service": service,
		},
	}
}

func NewUser(userID ksuid.ID) Actor {
	return Actor{
		Type: TypeUser,
		Params: map[string]string{
			"user_id": userID.String(),
		},
	}
}

func NewUserWithSession(sessionID, userID ksuid.ID) Actor {
	return Actor{
		Type: TypeUser,
		Params: map[string]string{
			"session_id": sessionID.String(),
			"user_id":    userID.String(),
		},
	}
}
