// internal/scheduler/ramp_test.go
package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestRampPolicy_MultipleRates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	p := &Ramp{
		Rates:        []int{50, 100},
		StepDuration: 20 * time.Millisecond,
	}

	ch := p.Run(ctx)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Fatalf("expected ticks during ramp, got %d", count)
	}
}
