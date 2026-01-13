package stats

import (
	"fmt"
	"time"
)

type Report struct {
	Duration time.Duration

	Requests   uint64
	OK         uint64
	Exceptions uint64
	OtherErrs  uint64

	MinNS uint64
	AvgNS uint64
	P95NS uint64
	P99NS uint64
	MaxNS uint64

	RPS float64
}

func BuildReport(duration time.Duration, c *Counters, h HistSnapshot) Report {
	req, ok, ex, other := c.Snapshot()

	rps := 0.0
	if duration > 0 {
		rps = float64(req) / duration.Seconds()
	}

	return Report{
		Duration: duration,

		Requests:   req,
		OK:         ok,
		Exceptions: ex,
		OtherErrs:  other,

		MinNS: h.MinNS,
		AvgNS: h.AvgNS(),
		P95NS: h.QuantileNS(0.95),
		P99NS: h.QuantileNS(0.99),
		MaxNS: h.MaxNS,

		RPS: rps,
	}
}

func (r Report) String() string {
	return fmt.Sprintf(
		"Requests:   %d\nOK:         %d\nExceptions: %d\nOtherErrs:  %d\n\nLatency (ms):\n  min  %.3f\n  avg  %.3f\n  p95  %.3f\n  p99  %.3f\n  max  %.3f\n\nThroughput:\n  %.1f req/s\n",
		r.Requests, r.OK, r.Exceptions, r.OtherErrs,
		nsToMS(r.MinNS),
		nsToMS(r.AvgNS),
		nsToMS(r.P95NS),
		nsToMS(r.P99NS),
		nsToMS(r.MaxNS),
		r.RPS,
	)
}

func nsToMS(ns uint64) float64 {
	return float64(ns) / 1_000_000.0
}
