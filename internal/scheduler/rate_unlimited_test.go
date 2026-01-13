// internal/scheduler/rate_unlimited_test.go
package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestRatePolicy_Unlimited(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	p := &Rate{
		PerSecond: 0, // unlimited
	}

	ch := p.Run(ctx)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Fatalf("expected ticks in unlimited mode")
	}
}
