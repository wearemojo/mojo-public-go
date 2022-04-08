package kmsjwt

import (
	kms "cloud.google.com/go/kms/apiv1"
	jwtinterface "github.com/wearemojo/mojo-public-go/lib/jwt"
)

var (
	_ jwtinterface.Signer   = (*KMSJWT)(nil)
	_ jwtinterface.Verifier = (*KMSJWT)(nil)
)

type KMSJWT struct {
	*Signer
	*Verifier
}

func New(client *kms.KeyManagementClient, projectID, env, serviceName string) *KMSJWT {
	return &KMSJWT{
		Signer:   NewSigner(client, projectID, env, serviceName),
		Verifier: NewVerifier(client, projectID),
	}
}
