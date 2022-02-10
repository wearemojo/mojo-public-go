package secret

import (
	"context"
)

type Wrapper struct {
	secretID string
}

func New(ctx context.Context, secretID string) (w *Wrapper, err error) {
	w = &Wrapper{secretID: secretID}
	_, err = w.Get(ctx)
	return
}

func (w *Wrapper) Get(ctx context.Context) (string, error) {
	return Get(ctx, w.secretID)
}
