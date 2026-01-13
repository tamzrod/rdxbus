// internal/engine/modbus_engine.go
package engine

import (
	"context"
	"time"

	"github.com/tamzrod/rdxbus/internal/client"
)

// ModbusEngine is the single execution throat for Modbus TCP.
// It executes exactly one request per call.
type ModbusEngine struct {
	TargetAddr string
	Strict     bool
}

func (e *ModbusEngine) Execute(ctx context.Context, req Request) Result {
	start := time.Now()

	// Dial per call for now (simple + correct).
	conn, err := client.Dial(e.TargetAddr, req.Timeout)
	if err != nil {
		return Result{
			UnitID:       req.UnitID,
			FunctionCode: req.FunctionCode,
			Address:      req.Address,
			Quantity:     req.Quantity,
			Duration:     time.Since(start),
			Err:          err,
		}
	}
	defer conn.Close()

	// Build request frame
	r := client.NewRequest()
	reqBuf := make([]byte, 12) // request.go expects >=12
	frame := r.BuildReadRequest(reqBuf, req.UnitID, req.FunctionCode, req.Address, req.Quantity)
	expectedTxID := r.TxID()

	if err := conn.Write(frame); err != nil {
		return Result{
			UnitID:       req.UnitID,
			FunctionCode: req.FunctionCode,
			Address:      req.Address,
			Quantity:     req.Quantity,
			Duration:     time.Since(start),
			Err:          err,
		}
	}

	// Read MBAP header (7 bytes)
	hdr := make([]byte, 7)
	if err := conn.ReadFull(hdr); err != nil {
		return Result{
			UnitID:       req.UnitID,
			FunctionCode: req.FunctionCode,
			Address:      req.Address,
			Quantity:     req.Quantity,
			Duration:     time.Since(start),
			Err:          err,
		}
	}

	// Parse response into PDU buffer
	parser := client.NewResponseParser(e.Strict)

	// Oversize buffer is fine; higher layers can decode from the prefix.
	pduBuf := make([]byte, 512)

	if err := parser.Parse(conn, expectedTxID, req.FunctionCode, hdr, pduBuf); err != nil {
		return Result{
			UnitID:       req.UnitID,
			FunctionCode: req.FunctionCode,
			Address:      req.Address,
			Quantity:     req.Quantity,
			Duration:     time.Since(start),
			Err:          err,
		}
	}

	// Return protocol-level bytes only (no decoding here).
	raw := make([]byte, len(pduBuf))
	copy(raw, pduBuf)

	return Result{
		UnitID:       req.UnitID,
		FunctionCode: req.FunctionCode,
		Address:      req.Address,
		Quantity:     req.Quantity,
		Raw:          raw,
		Duration:     time.Since(start),
		Err:          nil,
	}
}
