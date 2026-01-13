// internal/scheduler/rate.go
package scheduler

import (
	"context"
	"time"
)

type Rate struct {
	PerSecond int
}

func (p *Rate) Run(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		if p.PerSecond <= 0 {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- struct{}{}
				}
			}
		}

		interval := time.Second / time.Duration(p.PerSecond)
		ticker := time.NewTicker(interval)
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
