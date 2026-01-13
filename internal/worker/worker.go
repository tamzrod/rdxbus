// internal/worker/worker.go
package worker

import (
	"context"

	"github.com/tamzrod/rdxbus/internal/engine"
)

// Result is a direct pass-through of engine execution.
type Result struct {
	EngineResult engine.Result
}

// Execute performs exactly ONE engine execution.
// No protocol logic. No parsing. No retries. No scheduling.
func Execute(
	ctx context.Context,
	eng engine.Engine,
	req engine.Request,
) Result {
	return Result{
		EngineResult: eng.Execute(ctx, req),
	}
}
