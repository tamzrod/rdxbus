// internal/scan/runner.go
package scan

import (
	"context"

	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/worker"
)

type Runner struct {
	Engine engine.Engine
}

func (r *Runner) Run(ctx context.Context, strat Strategy) {
	for {
		req, ok := strat.Next()
		if !ok {
			return
		}

		res := worker.Execute(ctx, r.Engine, req)

		if strat.Observe(res.EngineResult) == Stop {
			return
		}
	}
}
