package errgroup

import (
	"context"

	realerrgroup "golang.org/x/sync/errgroup"
)

type Group struct {
	g *realerrgroup.Group

	//nolint:containedctx // outweighed by the benefit of reducing mistakes
	gctx context.Context
}

func WithContext(ctx context.Context) *Group {
	g, gctx := realerrgroup.WithContext(ctx)

	return &Group{
		g:    g,
		gctx: gctx,
	}
}

func (g *Group) Wait() error {
	return g.g.Wait()
}

func (g *Group) SetLimit(n int) {
	g.g.SetLimit(n)
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.g.Go(func() error {
		return f(g.gctx)
	})
}

func (g *Group) TryGo(f func(ctx context.Context) error) bool {
	return g.g.TryGo(func() error {
		return f(g.gctx)
	})
}
