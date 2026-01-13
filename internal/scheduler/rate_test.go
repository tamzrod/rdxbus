// internal/scheduler/rate_test.go
package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestRatePolicy_LimitedRate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	p := &Rate{
		PerSecond: 100, // ~1 tick per 10ms
	}

	ch := p.Run(ctx)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Fatalf("expected ticks, got %d", count)
	}

	// Upper bound sanity check (should not explode)
	if count > 10 {
		t.Fatalf("unexpectedly high tick count: %d", count)
	}
}
