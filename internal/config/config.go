package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	TargetAddr string

	Workers int
	Rate    int
	Duration time.Duration

	RampRates    []int
	StepDuration time.Duration

	UnitID       uint8
	FunctionCode uint8
	Address      uint16
	Quantity     uint16

	Timeout time.Duration
	Strict  bool
	Quiet   bool
}

func Parse() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.TargetAddr, "target", "127.0.0.1:502", "Modbus TCP target address")

	flag.IntVar(&cfg.Workers, "workers", 10, "Number of concurrent workers")
	flag.IntVar(&cfg.Rate, "rate", 0, "Requests per second (0 = unlimited)")
	flag.DurationVar(&cfg.Duration, "duration", 10*time.Second, "Test duration")

	ramp := flag.String("ramp", "", "Ramp rates, e.g. 100,500,1000")
	flag.DurationVar(&cfg.StepDuration, "step-duration", 5*time.Second, "Duration per ramp step")

	unit := flag.Int("unit", 1, "Modbus Unit ID")
	fc := flag.Int("fc", 3, "Modbus function code")
	addr := flag.Int("address", 0, "Starting register address")
	qty := flag.Int("quantity", 10, "Number of registers")

	flag.DurationVar(&cfg.Timeout, "timeout", 100*time.Millisecond, "Socket timeout")
	flag.BoolVar(&cfg.Strict, "strict", false, "Strict Modbus TCP framing")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "Minimal output")

	flag.Parse()

	cfg.UnitID = uint8(*unit)
	cfg.FunctionCode = uint8(*fc)
	cfg.Address = uint16(*addr)
	cfg.Quantity = uint16(*qty)

	if *ramp != "" {
		parts := strings.Split(*ramp, ",")
		for _, p := range parts {
			var v int
			fmt.Sscanf(p, "%d", &v)
			if v > 0 {
				cfg.RampRates = append(cfg.RampRates, v)
			}
		}
	}

	if err := cfg.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	return cfg
}

func (c *Config) validate() error {
	if c.TargetAddr == "" {
		return fmt.Errorf("target required")
	}
	if c.Workers <= 0 {
		return fmt.Errorf("workers must be > 0")
	}
	if c.Quantity == 0 {
		return fmt.Errorf("quantity must be > 0")
	}
	if len(c.RampRates) > 0 && c.StepDuration <= 0 {
		return fmt.Errorf("step-duration must be > 0")
	}
	return nil
}

// EngineReadConfig is consumed by the read engine.
// It must remain pure data (no flags, no os.Exit, no I/O).
type EngineReadConfig struct {
	TargetAddr   string
	UnitID       uint8
	FunctionCode uint8
	Address      uint16
	Quantity     uint16
	Timeout      time.Duration
	Strict       bool

}

// ToEngineRead converts CLI/test configuration into engine-safe config.
func (c *Config) ToEngineRead() EngineReadConfig {
	return EngineReadConfig{
		TargetAddr:   c.TargetAddr,
		UnitID:       c.UnitID,
		FunctionCode: c.FunctionCode,
		Address:      c.Address,
		Quantity:     c.Quantity,
		Timeout:      c.Timeout,
		Strict:       c.Strict,

	}
}
