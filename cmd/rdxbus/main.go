package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"modbus-stress/internal/client"
	"modbus-stress/internal/config"
	"modbus-stress/internal/scheduler"
	"modbus-stress/internal/stats"
	"modbus-stress/internal/worker"
)

func runOnce(cfg *config.Config, rate int, duration time.Duration) {
	fmt.Printf("\n=== RAMP STEP: %d req/s ===\n", rate)

	stop := make(chan struct{})
	sched := scheduler.NewScheduler(rate)
	defer sched.Stop()

	results := make(chan worker.Result, cfg.Workers*64)
	counters := &stats.Counters{}
	hist := stats.NewHistogram()

	for i := 0; i < cfg.Workers; i++ {
		conn, err := client.Dial(cfg.TargetAddr, cfg.Timeout)
		if err != nil {
			fmt.Println("dial failed:", err)
			return
		}

		w := worker.NewWorker(
			conn,
			sched,
			cfg.UnitID,
			cfg.FunctionCode,
			cfg.Address,
			cfg.Quantity,
			cfg.Strict,
		)

		go w.Run(stop, results)
	}

	end := time.Now().Add(duration)

	for time.Now().Before(end) {
		select {
		case r := <-results:
			counters.IncRequests()
			if r.Err == nil {
				counters.IncOK()
			} else if _, ok := client.IsModbusException(r.Err); ok {
				counters.IncExceptions()
			} else {
				counters.IncOtherErrs()
			}
			if r.Latency > 0 {
				hist.Record(r.Latency)
			}
		default:
			time.Sleep(200 * time.Microsecond)
		}
	}

	close(stop)
	time.Sleep(300 * time.Millisecond)

	rep := stats.BuildReport(duration, counters, hist.Snapshot())
	fmt.Println(rep.String())
}

func main() {
	cfg := config.Parse()

	// Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(0)
	}()

	if len(cfg.RampRates) > 0 {
		for _, r := range cfg.RampRates {
			runOnce(cfg, r, cfg.StepDuration)
		}
		return
	}

	// Single run (original behavior)
	runOnce(cfg, cfg.Rate, cfg.Duration)
}
