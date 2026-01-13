// internal/scan/address_test.go
package scan

import (
	"testing"

	"github.com/tamzrod/rdxbus/internal/engine"
)

func TestAddressScan_StepThenRefine(t *testing.T) {
	base := engine.Request{
		FunctionCode: 3,
		Quantity:     1,
	}

	scan := NewAddressScan(base, 0, 100, 10)

	// Step scan phase
	expect := []uint16{0, 10, 20}
	for _, addr := range expect {
		req, ok := scan.Next()
		if !ok || req.Address != addr {
			t.Fatalf("expected address %d, got %+v", addr, req)
		}
		if scan.Observe(engine.Result{Err: fakeErr()}) != Continue {
			t.Fatalf("expected Continue during step scan")
		}
	}

	// First success at address 30
	req, ok := scan.Next()
	if !ok || req.Address != 30 {
		t.Fatalf("expected address 30, got %+v", req)
	}
	if scan.Observe(engine.Result{Address: 30, Err: nil}) != Continue {
		t.Fatalf("expected Continue after first success")
	}

	// Refinement phase: stop on first success
req, ok = scan.Next()
if !ok || req.Address != 21 {
	t.Fatalf("expected first refine address 21, got %+v", req)
}

if scan.Observe(engine.Result{Address: 21, Err: nil}) != Stop {
	t.Fatalf("expected Stop on first refine success")
}

// Should be done
if _, ok := scan.Next(); ok {
	t.Fatalf("expected scan to finish after refine success")
}


	// Should be done
	if _, ok := scan.Next(); ok {
		t.Fatalf("expected scan to finish")
	}
}
