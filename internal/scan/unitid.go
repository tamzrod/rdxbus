// internal/scan/unitid.go
package scan

import "github.com/tamzrod/rdxbus/internal/engine"

type UnitIDScan struct {
	Base engine.Request

	Start uint8
	End   uint8
	Step  uint8

	current uint8
	done    bool
}

func NewUnitIDScan(base engine.Request, start, end, step uint8) *UnitIDScan {
	return &UnitIDScan{
		Base:    base,
		Start:   start,
		End:     end,
		Step:    step,
		current: start,
	}
}

func (s *UnitIDScan) Next() (engine.Request, bool) {
	if s.done || s.current > s.End {
		return engine.Request{}, false
	}

	req := s.Base
	req.UnitID = s.current

	s.current += s.Step
	return req, true
}

func (s *UnitIDScan) Observe(result engine.Result) Decision {
	if result.Err == nil {
		s.done = true
		return Stop
	}
	return Continue
}
