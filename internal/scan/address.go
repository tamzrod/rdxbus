// internal/scan/address.go
package scan

import "github.com/tamzrod/rdxbus/internal/engine"

type AddressScan struct {
	Base engine.Request

	Start uint16
	End   uint16
	Step  uint16

	current uint16

	foundAt *uint16
	refine  bool
	done    bool
}

func NewAddressScan(base engine.Request, start, end, step uint16) *AddressScan {
	return &AddressScan{
		Base:    base,
		Start:   start,
		End:     end,
		Step:    step,
		current: start,
	}
}

func (s *AddressScan) Next() (engine.Request, bool) {
	if s.done {
		return engine.Request{}, false
	}

	if s.current > s.End {
		s.done = true
		return engine.Request{}, false
	}

	req := s.Base
	req.Address = s.current

	if s.refine {
		s.current++
	} else {
		s.current += s.Step
	}

	return req, true
}

func (s *AddressScan) Observe(result engine.Result) Decision {
	if result.Err != nil {
		return Continue
	}

	// First success: switch to refine mode
	if !s.refine {
		addr := result.Address
		start := addr
		if addr >= s.Step {
			start = addr - s.Step + 1
		}

		s.current = start
		s.refine = true
		return Continue
	}

	// Refine success → stop
	s.done = true
	return Stop
}
