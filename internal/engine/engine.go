// internal/engine/engine.go
package engine

import (
	"context"
	"time"
)

type Request struct {
	UnitID       uint8
	FunctionCode uint8
	Address      uint16
	Quantity     uint16
	Timeout      time.Duration
}

type Result struct {
	UnitID       uint8
	FunctionCode uint8
	Address      uint16
	Quantity     uint16

	// Raw is protocol-level bytes only (PDU-ish buffer as returned by parser).
	// No decoding, scaling, or interpretation happens here.
	Raw []byte

	Duration time.Duration
	Err      error
}

type Engine interface {
	Execute(ctx context.Context, req Request) Result
}
