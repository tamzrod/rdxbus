package stats

import (
	"math/bits"
	"time"
)

type Histogram struct {
	buckets [64]uint64

	minNS uint64
	maxNS uint64
	sumNS uint64
	count uint64
}

func NewHistogram() *Histogram {
	return &Histogram{}
}

func (h *Histogram) Record(d time.Duration) {
	ns := uint64(d.Nanoseconds())
	if ns == 0 {
		ns = 1
	}

	if h.count == 0 {
		h.minNS = ns
		h.maxNS = ns
	} else {
		if ns < h.minNS {
			h.minNS = ns
		}
		if ns > h.maxNS {
			h.maxNS = ns
		}
	}

	h.count++
	h.sumNS += ns

	idx := bucketIndex(ns)
	h.buckets[idx]++
}

func bucketIndex(ns uint64) int {
	return 63 - bits.LeadingZeros64(ns)
}

type HistSnapshot struct {
	Buckets [64]uint64
	MinNS   uint64
	MaxNS   uint64
	SumNS   uint64
	Count   uint64
}

func (h *Histogram) Snapshot() HistSnapshot {
	return HistSnapshot{
		Buckets: h.buckets,
		MinNS:   h.minNS,
		MaxNS:   h.maxNS,
		SumNS:   h.sumNS,
		Count:   h.count,
	}
}

func (s HistSnapshot) AvgNS() uint64 {
	if s.Count == 0 {
		return 0
	}
	return s.SumNS / s.Count
}

func (s HistSnapshot) QuantileNS(q float64) uint64 {
	if s.Count == 0 {
		return 0
	}
	if q <= 0 {
		return s.MinNS
	}
	if q >= 1 {
		return s.MaxNS
	}

	target := uint64(float64(s.Count) * q)
	if target == 0 {
		target = 1
	}

	var seen uint64
	for i := 0; i < len(s.Buckets); i++ {
		seen += s.Buckets[i]
		if seen >= target {
			return 1 << uint(i)
		}
	}
	return s.MaxNS
}
