package indirect

import (
	"context"

	"github.com/wearemojo/mojo-public-go/lib/merr"
)

var (
	Debug func(ctx context.Context, err merr.Merrer)
	Info  func(ctx context.Context, err merr.Merrer)
	Warn  func(ctx context.Context, err merr.Merrer)
	Error func(ctx context.Context, err merr.Merrer)
)
