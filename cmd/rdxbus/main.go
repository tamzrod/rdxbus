// cmd/rdxbus/main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tamzrod/rdxbus/internal/config"
	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/format"
	"github.com/tamzrod/rdxbus/internal/worker"
)

func main() {
	// Parse CLI configuration (CLI concern only)
	cfg := config.Parse()

	// Build engine (single execution throat)
	eng := &engine.ModbusEngine{
		TargetAddr: cfg.TargetAddr,
		Strict:     cfg.Strict,
	}

	// Build engine request (pure data)
	req := engine.Request{
		UnitID:       cfg.UnitID,
		FunctionCode: cfg.FunctionCode,
		Address:      cfg.Address,
		Quantity:     cfg.Quantity,
		Timeout:      cfg.Timeout,
	}

	// Execute exactly ONE request
	ctx := context.Background()
	wr := worker.Execute(ctx, eng, req)

	if wr.EngineResult.Err != nil {
		fmt.Fprintln(os.Stderr, "read error:", wr.EngineResult.Err)
		os.Exit(1)
	}

	// Decode values at the CLI/UI layer (optional, not engine responsibility)
	values, err := format.DecodeReadValues(
		wr.EngineResult.Raw,
		req.FunctionCode,
		req.Quantity,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode error:", err)
		os.Exit(1)
	}

	// Output
	fmt.Println("read successful")
	fmt.Println("latency:", wr.EngineResult.Duration)
	fmt.Println("values:", values)
}
