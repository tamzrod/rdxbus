// internal/scan/unitid_test.go
package scan

import (
	"testing"

	"github.com/tamzrod/rdxbus/internal/engine"
)

func TestUnitIDScan_StopOnFirstSuccess(t *testing.T) {
	base := engine.Request{
		FunctionCode: 3,
		Address:      0,
		Quantity:     1,
	}

	scan := NewUnitIDScan(base, 1, 200, 50)

	// First attempt: UnitID = 1 → error
	req, ok := scan.Next()
	if !ok || req.UnitID != 1 {
		t.Fatalf("expected first UnitID=1, got %+v", req)
	}
	if scan.Observe(engine.Result{Err: fakeErr()}) != Continue {
		t.Fatalf("expected Continue on error")
	}

	// Second attempt: UnitID = 51 → success
	req, ok = scan.Next()
	if !ok || req.UnitID != 51 {
		t.Fatalf("expected UnitID=51, got %+v", req)
	}
	if scan.Observe(engine.Result{Err: nil}) != Stop {
		t.Fatalf("expected Stop on first success")
	}

	// No more requests
	if _, ok := scan.Next(); ok {
		t.Fatalf("expected scan to stop")
	}
}
