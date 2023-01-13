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
	TypeUnknown           Type = "unknown"
	TypeInternal          Type = "internal"
	TypeService           Type = "service" // Mojo service
	TypeUser              Type = "user"    // Mojo user
	TypeSession           Type = "session"
	TypeExternalCloudAuth Type = "external_cloud_auth" // non-Mojo system
	TypeExternalUser      Type = "external_user"       // non-Mojo user
)

type Actor struct {
	Type   Type           `json:"type" bson:"type"`
	Params map[string]any `json:"params" bson:"params"`
}

func NewUnknown(params map[string]any) Actor {
	if params == nil {
		params = map[string]any{}
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
		Params: map[string]any{
			"revision":  version.Revision,
			"code_path": stacktrace.GetCallerCodePath(skip+1, codePathFallback),
		},
	}
}

func NewService(env, service string) Actor {
	return Actor{
		Type: TypeService,
		Params: map[string]any{
			"env":     env,
			"service": service,
		},
	}
}

func NewUser(sessionID, userID ksuid.ID) Actor {
	return Actor{
		Type: TypeUser,
		Params: map[string]any{
			"session_id": sessionID.String(),
			"user_id":    userID.String(),
		},
	}
}

func NewSession(sessionID ksuid.ID) Actor {
	return Actor{
		Type: TypeSession,
		Params: map[string]any{
			"session_id": sessionID.String(),
		},
	}
}

// Deprecated: remove once WP auth is gone
//
// TODO: remove
func NewUserWithoutSession(userID ksuid.ID) Actor {
	return Actor{
		Type: TypeUser,
		Params: map[string]any{
			"session_id": nil,
			"user_id":    userID.String(),
		},
	}
}

func NewExternalCloudAuth(typ, service string) Actor {
	return Actor{
		Type: TypeExternalCloudAuth,
		Params: map[string]any{
			"type":    typ,
			"service": service,
		},
	}
}

func NewExternalUser(typ, id, reference string) Actor {
	return Actor{
		Type: TypeExternalUser,
		Params: map[string]any{
			"type":      typ,
			"id":        id,
			"reference": reference,
		},
	}
}
