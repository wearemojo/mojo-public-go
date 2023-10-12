package errgroup

import (
	"sync"
)

type SettledGroup struct {
	wg sync.WaitGroup

	errs   []error
	errsMu sync.Mutex
}

func (g *SettledGroup) Go(fn func() (err error)) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := fn(); err != nil {
			g.errsMu.Lock()
			defer g.errsMu.Unlock()

			g.errs = append(g.errs, err)
		}
	}()
}

func (g *SettledGroup) Wait() []error {
	g.wg.Wait()

	return g.errs
}
