package config

import (
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gjson"
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

// The duration fields previously carried a `format:nano` tag, which json/v2
// rejects as unsupported. That failed the whole struct, so every other field
// became unreadable alongside them - including when the duration fields were
// absent from the payload entirely.
func TestMongoDBJSONRoundTrip(t *testing.T) {
	is := is.New(t)

	maxConnIdleTime := gjson.Duration(5 * time.Minute)
	cfg := MongoDB{
		URI:             "mongodb://localhost:27017/test",
		ConnectTimeout:  gjson.Duration(10 * time.Second),
		MaxConnIdleTime: &maxConnIdleTime,
	}

	data, err := json.Marshal(cfg)
	is.NoErr(err)

	out, err := gjson.Unmarshal[MongoDB](data)
	is.NoErr(err)
	is.Equal(cfg.URI, out.URI)
	is.Equal(10*time.Second, out.ConnectTimeout.Duration())
	is.Equal(5*time.Minute, out.MaxConnIdleTime.Duration())

	// the durations being absent must not break the fields alongside them
	partial, err := gjson.Unmarshal[MongoDB]([]byte(`{"uri":"mongodb://localhost:27017/test"}`))
	is.NoErr(err)
	is.Equal("mongodb://localhost:27017/test", partial.URI)
}
