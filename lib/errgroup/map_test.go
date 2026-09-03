package errgroup_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wearemojo/mojo-public-go/lib/errgroup"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

func TestMapAndWaitKeepsInputOrder(t *testing.T) {
	ctx := context.Background()

	// the first input sleeps longest, so completion order is the reverse of
	// input order - results must still come back in input order
	res, err := errgroup.Map(ctx, []int{3, 2, 1}, func(_ context.Context, in int) (int, error) {
		time.Sleep(time.Duration(in) * 10 * time.Millisecond)
		return in * 2, nil
	})

	require.NoError(t, err)
	require.Equal(t, []int{6, 4, 2}, res)
}

func TestMapAndWaitReturnsFirstError(t *testing.T) {
	ctx := context.Background()

	res, err := errgroup.Map(ctx, []int{1, 2}, func(_ context.Context, in int) (int, error) {
		if in == 2 {
			return 0, merr.New(ctx, "boom", nil)
		}
		return in, nil
	})

	mErr, ok := errors.AsType[merr.E](err)
	require.True(t, ok)
	require.Equal(t, merr.Code("boom"), mErr.Code)
	require.Len(t, res, 2)
}

func TestMapAndWaitOnAGroupWithALimit(t *testing.T) {
	g := errgroup.WithContext(context.Background())
	g.SetLimit(1)

	res, err := g.MapAndWait([]string{"a", "b", "c"}, func(_ context.Context, in string) (string, error) {
		return in + "!", nil
	})

	require.NoError(t, err)
	require.Equal(t, []string{"a!", "b!", "c!"}, res)
}
