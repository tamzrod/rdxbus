// internal/scheduler/interval.go
package scheduler

import (
	"context"
	"time"
)

type Interval struct {
	Every time.Duration
}

func (p *Interval) Run(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		ticker := time.NewTicker(p.Every)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ch <- struct{}{}
			}
		}
	}()

	return ch
}
