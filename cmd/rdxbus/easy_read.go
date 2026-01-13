// cmd/rdxbus/easy_read.go
package main

import (
	"bufio"
	"context"
	"os"
	"time"

	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/format"
	"github.com/tamzrod/rdxbus/internal/output"
	"github.com/tamzrod/rdxbus/internal/render"
	"github.com/tamzrod/rdxbus/internal/scheduler"
	"github.com/tamzrod/rdxbus/internal/worker"
)

func easyReadOnce(reader *bufio.Reader, target string, unitID int) {
	fc := promptInt(reader, "Function code (1–4)", 3)
	addr := promptInt(reader, "Start address", 0)
	qty := promptInt(reader, "Quantity", 10)

	eng := &engine.ModbusEngine{TargetAddr: target}

	req := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: uint8(fc),
		Address:      uint16(addr),
		Quantity:     uint16(qty),
		Timeout:      2 * time.Second,
	}

	res := worker.Execute(context.Background(), eng, req)
	if res.EngineResult.Err != nil {
		render.Render(os.Stdout, output.Output{Error: res.EngineResult.Err.Error()})
		return
	}

	values, err := format.DecodeReadValues(
		res.EngineResult.Raw,
		req.FunctionCode,
		req.Quantity,
	)
	if err != nil {
		render.Render(os.Stdout, output.Output{Error: err.Error()})
		return
	}

	render.Render(os.Stdout, output.Output{
		Meta: output.Meta{
			Mode:     "read",
			Target:   target,
			UnitID:   uint8(unitID),
			Function: uint8(fc),
			Latency:  res.EngineResult.Duration,
		},
		Table: buildTable(buildRows(req.Address, values)),
	})
}

func easyPoll(reader *bufio.Reader, target string, unitID int) {
	fc := promptInt(reader, "Function code (1–4)", 3)
	addr := promptInt(reader, "Start address", 0)
	qty := promptInt(reader, "Quantity", 10)
	intervalMs := promptInt(reader, "Poll interval (ms)", 1000)

	eng := &engine.ModbusEngine{TargetAddr: target}

	req := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: uint8(fc),
		Address:      uint16(addr),
		Quantity:     uint16(qty),
		Timeout:      2 * time.Second,
	}

	policy := &scheduler.Interval{
		Every: time.Duration(intervalMs) * time.Millisecond,
	}

	for range policy.Run(context.Background()) {
		res := worker.Execute(context.Background(), eng, req)
		if res.EngineResult.Err != nil {
			render.Render(os.Stdout, output.Output{Error: res.EngineResult.Err.Error()})
			return
		}

		values, _ := format.DecodeReadValues(
			res.EngineResult.Raw,
			req.FunctionCode,
			req.Quantity,
		)

		render.Render(os.Stdout, output.Output{
			Meta: output.Meta{
				Mode:     "poll",
				Target:   target,
				UnitID:   uint8(unitID),
				Function: uint8(fc),
				Latency:  res.EngineResult.Duration,
			},
			Table: buildTable(buildRows(req.Address, values)),
		})
	}
}
