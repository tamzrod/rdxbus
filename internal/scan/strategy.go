// internal/scan/strategy.go
package scan

import "github.com/tamzrod/rdxbus/internal/engine"

// Decision tells the runner what to do next.
type Decision int

const (
	Continue Decision = iota
	Stop
)

// Strategy defines how scanning progresses.
// It is stateful by design.
type Strategy interface {
	// Next returns the next engine request to try.
	// ok=false means scan is finished.
	Next() (req engine.Request, ok bool)

	// Observe receives the result of the last execution
	// and decides whether to continue or stop.
	Observe(result engine.Result) Decision
}
