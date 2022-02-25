package authparsing

import (
	"github.com/cuvva/cuvva-public-go/lib/ksuid"
)

type AuthStateType string

const (
	// U2S represents user-to-service auth
	U2S AuthStateType = "u2s"

	// S2S represents service-to-service auth
	S2S AuthStateType = "s2s"
)

type AuthState struct {
	Type AuthStateType

	// UserID is only set for AuthStateTypeU2S
	UserID ksuid.ID
}
