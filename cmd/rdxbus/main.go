// cmd/rdxbus/main.go
package main

import (
	"fmt"
	"os"

	"github.com/tamzrod/rdxbus/internal/config"
	"github.com/tamzrod/rdxbus/internal/worker"
)

func main() {
	// Parse CLI/test configuration (temporary)
	cfg := config.Parse()

	// Convert to engine-safe config
	engineCfg := cfg.ToEngineRead()

	// Execute exactly ONE Modbus read via the engine
	result := worker.ExecuteRead(engineCfg)

	if result.Err != nil {
		fmt.Fprintln(os.Stderr, "read error:", result.Err)
		os.Exit(1)
	}

	fmt.Println("read successful")
	fmt.Println("latency:", result.Latency)
	fmt.Println("values:", result.Values)
}
