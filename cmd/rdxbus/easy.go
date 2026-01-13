// cmd/rdxbus/easy.go
package main

import (
	"bufio"
	"fmt"
	"os"
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

	switch promptInt(reader, "Selection", 1) {
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
