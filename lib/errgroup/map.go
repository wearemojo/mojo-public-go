package errgroup

import (
	"context"
)

func Map[TIn, TOut any](ctx context.Context, inputs []TIn, fn func(ctx context.Context, input TIn) (TOut, error)) (res []TOut, err error) {
	return WithContext(ctx).MapAndWait(inputs, fn)
}

// MapAndWait runs fn over every input concurrently, then waits for the group.
//
// Results keep the order of inputs regardless of completion order.
func (g *Group) MapAndWait[TIn, TOut any](inputs []TIn, fn func(ctx context.Context, input TIn) (TOut, error)) (res []TOut, err error) {
	res = make([]TOut, len(inputs))

	for idx, input := range inputs {
		g.Go(func(ctx context.Context) (err error) {
			res[idx], err = fn(ctx, input)
			return err
		})
	}

	err = g.Wait()
	return res, err
}
