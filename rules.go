//nolint:unused // ruleguard rules aren't exported
package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func mustWrapMerrCodes(m dsl.Matcher) {
	m.Import("github.com/wearemojo/mojo-public-go/lib/merr")

	m.Match(`return $*x, $y`).
		Where(m["y"].Type.Is("merr.Code")).
		Report(`should wrap errors with merr.New`).
		Suggest(`return $x, merr.New(ctx, $y, nil)`)
}

func merrNoWrappingNil(m dsl.Matcher) {
	m.Import("github.com/wearemojo/mojo-public-go/lib/merr")

	m.Match(`merr.New($a, $b, $c, nil)`).
		Report("merr.New panics with `nil` reasons").
		Suggest(`merr.New($a, $b, $c)`)
}

func ksuidResourcePattern(m dsl.Matcher) {
	m.Import("github.com/wearemojo/mojo-public-go/lib/ksuid")

	m.Match(`ksuid.Generate($x, $y)`).
		Where(m["y"].Text.Matches(`_`)).
		Report(`ksuid resource name must not contain underscores`)

	m.Match(`$z.Generate($x, $y)`).
		Where(
			m["z"].Type.Is("*ksuid.Node") &&
				m["y"].Text.Matches(`_`),
		).
		Report(`ksuid resource name must not contain underscores`)
}

func slicePointers(m dsl.Matcher) {
	m.Match(`[]*$x`).
		Report(`use []T instead of []*T for slices of structs`).
		Suggest(`[]$x`)
}

func ttlcacheNewKeyedAny(m dsl.Matcher) {
	m.Import("github.com/wearemojo/mojo-public-go/lib/ttlcache")

	m.Match(`ttlcache.NewKeyed[$T, any]($ttl)`).
		Report(`use ttlcache.NewSingular instead`)
}

func forbidOmitEmpty(m dsl.Matcher) {
	m.Match(`struct { $*fields }`).
		Where(
			m["fields"].Text.Matches(`,omitempty`),
		).
		Report(`avoid using ,omitempty in struct field tags`)
}
