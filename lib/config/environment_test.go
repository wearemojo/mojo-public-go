package config

import (
	"testing"

	"github.com/matryer/is"
)

type envGetter map[string]string

func (e envGetter) Get(k string) string {
	return e[k]
}

func TestFromEnvironment(t *testing.T) {
	is := is.New(t)

	env := envGetter{ConfigEnvironmentVariable: `{"foo": "bar"}`}

	dest := struct {
		Foo string `json:"foo"`
	}{}

	err := FromEnvironment(env.Get, &dest)
	is.NoErr(err)
	is.Equal("bar", dest.Foo)
}

func TestEnvironmentName(t *testing.T) {
	is := is.New(t)

	env := envGetter{ConfigEnvironmentVariable: `{"env": "prod"}`}

	is.Equal("prod", EnvironmentName(env.Get))
}

func TestEnvironmentNameDev(t *testing.T) {
	is := is.New(t)

	env := envGetter{}

	is.Equal("local", EnvironmentName(env.Get))
}
