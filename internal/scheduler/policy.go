// internal/scheduler/policy.go
package scheduler

import "context"

// Policy emits execution ticks.
// It never executes work itself.
type Policy interface {
	// Run blocks and emits a signal each time work should execute.
	// The channel is closed when scheduling is complete.
	Run(ctx context.Context) <-chan struct{}
}
