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
