// cmd/rdxbus/easy.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/format"
	"github.com/tamzrod/rdxbus/internal/output"
	"github.com/tamzrod/rdxbus/internal/render"
	"github.com/tamzrod/rdxbus/internal/scheduler"
	"github.com/tamzrod/rdxbus/internal/scan"
	"github.com/tamzrod/rdxbus/internal/worker"
)

func runEasy() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("RDXBus Easy Mode")
	fmt.Println("----------------")

	target := prompt(reader, "Target address", "127.0.0.1:502")
	unitID := promptInt(reader, "Unit ID", 1)

	fmt.Println("\nChoose operation:")
	fmt.Println("  1) Read once")
	fmt.Println("  2) Poll continuously")
	fmt.Println("  3) Scan helpers")

	mode := promptInt(reader, "Selection", 1)

	switch mode {
	case 1:
		easyReadOnce(reader, target, unitID)
	case 2:
		easyPoll(reader, target, unitID)
	case 3:
		runEasyScan(reader)
	default:
		fmt.Println("Invalid selection")
	}
}

//
// ---- READ ONCE -----------------------------------------------------------
//

func easyReadOnce(reader *bufio.Reader, target string, unitID int) {
	fc := promptInt(reader, "Function code (1–4)", 3)
	addr := promptInt(reader, "Start address", 0)
	qty := promptInt(reader, "Quantity", 10)

	eng := &engine.ModbusEngine{
		TargetAddr: target,
		Strict:     false,
	}

	req := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: uint8(fc),
		Address:      uint16(addr),
		Quantity:     uint16(qty),
		Timeout:      2 * time.Second,
	}

	ctx := context.Background()
	res := worker.Execute(ctx, eng, req)

	if res.EngineResult.Err != nil {
		render.Render(os.Stdout, output.Output{
			Error: res.EngineResult.Err.Error(),
		})
		return
	}

	values, err := format.DecodeReadValues(
		res.EngineResult.Raw,
		req.FunctionCode,
		req.Quantity,
	)
	if err != nil {
		render.Render(os.Stdout, output.Output{
			Error: err.Error(),
		})
		return
	}

	rows := buildRows(req.Address, values)

	out := output.Output{
		Meta: output.Meta{
			Mode:      "read",
			Target:    target,
			UnitID:    uint8(unitID),
			Function:  uint8(fc),
			Latency:   res.EngineResult.Duration,
			Timestamp: time.Now(),
		},
		Table: buildTable(rows),
	}

	render.Render(os.Stdout, out)
}

//
// ---- POLL MODE -----------------------------------------------------------
//

func easyPoll(reader *bufio.Reader, target string, unitID int) {
	fc := promptInt(reader, "Function code (1–4)", 3)
	addr := promptInt(reader, "Start address", 0)
	qty := promptInt(reader, "Quantity", 10)
	intervalMs := promptInt(reader, "Poll interval (ms)", 1000)

	eng := &engine.ModbusEngine{
		TargetAddr: target,
		Strict:     false,
	}

	req := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: uint8(fc),
		Address:      uint16(addr),
		Quantity:     uint16(qty),
		Timeout:      2 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	policy := &scheduler.Interval{
		Every: time.Duration(intervalMs) * time.Millisecond,
	}

	fmt.Printf("\nPolling every %d ms — press Ctrl+C to stop\n\n", intervalMs)

	for range policy.Run(ctx) {
		res := worker.Execute(ctx, eng, req)

		if res.EngineResult.Err != nil {
			render.Render(os.Stdout, output.Output{
				Error: res.EngineResult.Err.Error(),
			})
			return
		}

		values, err := format.DecodeReadValues(
			res.EngineResult.Raw,
			req.FunctionCode,
			req.Quantity,
		)
		if err != nil {
			render.Render(os.Stdout, output.Output{
				Error: err.Error(),
			})
			return
		}

		rows := buildRows(req.Address, values)

		out := output.Output{
			Meta: output.Meta{
				Mode:      "poll",
				Target:    target,
				UnitID:    uint8(unitID),
				Function:  uint8(fc),
				Latency:   res.EngineResult.Duration,
				Timestamp: time.Now(),
			},
			Table: buildTable(rows),
		}

		render.Render(os.Stdout, out)
	}
}

//
// ---- SCAN HELPERS --------------------------------------------------------
//

func runEasyScan(reader *bufio.Reader) {
	fmt.Println("\nScan helpers")
	fmt.Println("------------")

	fmt.Println("  1) Find Unit ID")
	fmt.Println("  2) Scan address range")

	choice := promptInt(reader, "Selection", 1)

	switch choice {
	case 1:
		easyScanUnitID(reader)
	case 2:
		easyScanAddress(reader)
	default:
		fmt.Println("Invalid selection")
	}
}

func easyScanUnitID(reader *bufio.Reader) {
	target := prompt(reader, "Target address", "127.0.0.1:502")
	start := promptInt(reader, "Start Unit ID", 1)
	end := promptInt(reader, "End Unit ID", 247)
	step := promptInt(reader, "Step", 50)

	eng := &engine.ModbusEngine{
		TargetAddr: target,
		Strict:     false,
	}

	baseReq := engine.Request{
		FunctionCode: 3,
		Address:      0,
		Quantity:     1,
		Timeout:      2 * time.Second,
	}

	strat := scan.NewUnitIDScan(
		baseReq,
		uint8(start),
		uint8(end),
		uint8(step),
	)

	runner := scan.Runner{Engine: eng}

	fmt.Println("\nScanning for Unit ID...")
	ctx := context.Background()
	runner.Run(ctx, strat)

	render.Render(os.Stdout, output.Output{
		Message: "Unit ID scan complete",
	})
}

func easyScanAddress(reader *bufio.Reader) {
	target := prompt(reader, "Target address", "127.0.0.1:502")
	unitID := promptInt(reader, "Unit ID", 1)
	start := promptInt(reader, "Start address", 0)
	end := promptInt(reader, "End address", 1000)
	step := promptInt(reader, "Step", 10)

	eng := &engine.ModbusEngine{
		TargetAddr: target,
		Strict:     false,
	}

	baseReq := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: 3,
		Quantity:     1,
		Timeout:      2 * time.Second,
	}

	strat := scan.NewAddressScan(
		baseReq,
		uint16(start),
		uint16(end),
		uint16(step),
	)

	runner := scan.Runner{Engine: eng}

	fmt.Println("\nScanning addresses...")
	ctx := context.Background()
	runner.Run(ctx, strat)

	render.Render(os.Stdout, output.Output{
		Message: "Address scan complete",
	})
}

//
// ---- SHARED HELPERS ------------------------------------------------------
//

func buildRows(start uint16, values []uint16) []output.Row {
	rows := make([]output.Row, 0, len(values))
	for i, v := range values {
		rows = append(rows, output.Row{
			Cells: map[string]any{
				"address": int(start) + i,
				"raw":     v,
			},
		})
	}
	return rows
}

func buildTable(rows []output.Row) *output.Table {
	return &output.Table{
		Columns: []output.Column{
			{Key: "address", Title: "Address"},
			{Key: "raw", Title: "Raw"},
		},
		Rows: rows,
	}
}

func prompt(reader *bufio.Reader, label, def string) string {
	fmt.Printf("%s [%s]: ", label, def)
	in, _ := reader.ReadString('\n')
	in = strings.TrimSpace(in)
	if in == "" {
		return def
	}
	return in
}

func promptInt(reader *bufio.Reader, label string, def int) int {
	for {
		v := prompt(reader, label, strconv.Itoa(def))
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
		fmt.Println("Invalid number")
	}
}
