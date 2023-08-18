package errgroup

import (
	"context"
)

func Map[TIn, TOut any](ctx context.Context, inputs []TIn, fn func(ctx context.Context, input TIn) (TOut, error)) (res []TOut, err error) {
	g := WithContext(ctx)
	return GroupMapAndWait(g, inputs, fn)
}

func GroupMapAndWait[TIn, TOut any](group *Group, inputs []TIn, fn func(ctx context.Context, input TIn) (TOut, error)) (res []TOut, err error) {
	res = make([]TOut, len(inputs))

	for idx, input := range inputs {
		idx, input := idx, input
		group.Go(func(ctx context.Context) (err error) {
			res[idx], err = fn(ctx, input)
			return
		})
	}

	err = group.Wait()
	return
}
