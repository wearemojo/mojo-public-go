package config

import (
	"testing"

	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

func TestMongoDBOptions(t *testing.T) {
	is := is.New(t)

	//nolint:gosec // G101 - not a real password
	m := &MongoDB{
		URI: "mongodb://foo:bar@127.0.0.1/demo?authSource=admin",
	}

	opts, dbName, err := m.Options(t.Context())

	is.NoErr(err)
	is.Equal(dbName, "demo")
	is.Equal(opts.Hosts, []string{"127.0.0.1"})
	is.Equal(opts.WriteConcern, writeconcern.Majority())

	is.Equal(opts.Auth, &options.Credential{
		AuthSource:  "admin",
		Username:    "foo",
		Password:    "bar",
		PasswordSet: true,
	})
}
