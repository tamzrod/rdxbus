// internal/scheduler/ramp.go
package scheduler

import (
	"context"
	"time"
)

type Ramp struct {
	Rates        []int
	StepDuration time.Duration
}

func (p *Ramp) Run(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for _, r := range p.Rates {
			if r <= 0 {
				continue
			}

			interval := time.Second / time.Duration(r)
			ticker := time.NewTicker(interval)

			stepTimer := time.NewTimer(p.StepDuration)

			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					stepTimer.Stop()
					return

				case <-stepTimer.C:
					ticker.Stop()
					goto nextRate

				case <-ticker.C:
					ch <- struct{}{}
				}
			}
		nextRate:
		}
	}()

	return ch
}
