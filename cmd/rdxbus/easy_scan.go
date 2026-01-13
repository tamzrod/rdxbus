// cmd/rdxbus/easy_scan.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/tamzrod/rdxbus/internal/engine"
	"github.com/tamzrod/rdxbus/internal/scan"
)

func runEasyScan(reader *bufio.Reader) {
	fmt.Println("\nScan helpers")
	fmt.Println("------------")
	fmt.Println("  1) Find Unit ID")
	fmt.Println("  2) Scan address range")

	switch promptInt(reader, "Selection", 1) {
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

	eng := &engine.ModbusEngine{TargetAddr: target}
	req := engine.Request{
		UnitID:       1, // valid baseline
		FunctionCode: 3,
		Address:      0,
		Quantity:     1,
		Timeout:      2 * time.Second,
	}

	strat := scan.NewUnitIDScan(req, 1, 247, 50)

	fmt.Println("\nScanning for Unit ID...")
	(&scan.Runner{Engine: eng}).Run(context.Background(), strat)

	fmt.Println("Unit ID scan complete")
}

func easyScanAddress(reader *bufio.Reader) {
	target := prompt(reader, "Target address", "127.0.0.1:502")
	unitID := promptInt(reader, "Unit ID", 1)

	eng := &engine.ModbusEngine{TargetAddr: target}
	req := engine.Request{
		UnitID:       uint8(unitID),
		FunctionCode: 3,
		Quantity:     1,
		Timeout:      2 * time.Second,
	}

	strat := scan.NewAddressScan(req, 0, 1000, 10)

	fmt.Println("\nScanning addresses...")
	(&scan.Runner{Engine: eng}).Run(context.Background(), strat)

	fmt.Println("Address scan complete")
}
