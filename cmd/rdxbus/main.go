// cmd/rdxbus/main.go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tamzrod/rdxbus/internal/config"
	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/format"
	"github.com/tamzrod/rdxbus/internal/worker"
)

func main() {
	// Friendly (easy) mode
	if len(os.Args) > 1 && os.Args[1] == "easy" {
		runEasy()
		return
	}

	// Expert CLI (single read, legacy behavior)
	cfg := config.Parse()

	eng := &engine.ModbusEngine{
		TargetAddr: cfg.TargetAddr,
		Strict:     cfg.Strict,
	}

	req := engine.Request{
		UnitID:       cfg.UnitID,
		FunctionCode: cfg.FunctionCode,
		Address:      cfg.Address,
		Quantity:     cfg.Quantity,
		Timeout:      cfg.Timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := worker.Execute(ctx, eng, req)

	if res.EngineResult.Err != nil {
		fmt.Fprintln(os.Stderr, "read error:", res.EngineResult.Err)
		os.Exit(1)
	}

	values, err := format.DecodeReadValues(
		res.EngineResult.Raw,
		req.FunctionCode,
		req.Quantity,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode error:", err)
		os.Exit(1)
	}

	fmt.Println("read successful")
	fmt.Println("latency:", res.EngineResult.Duration)
	fmt.Println("values:", values)
}
