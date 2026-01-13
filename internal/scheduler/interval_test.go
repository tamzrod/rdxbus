// internal/scheduler/interval_test.go
package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestIntervalPolicy_EmitsTicks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &Interval{
		Every: 10 * time.Millisecond,
	}

	ch := p.Run(ctx)

	count := 0
	timeout := time.After(50 * time.Millisecond)

	for {
		select {
		case <-ch:
			count++
			if count >= 3 {
				cancel()
			}
		case <-timeout:
			cancel()
			goto done
		}
	}

done:
	if count < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", count)
	}
}
