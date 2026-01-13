package stats

import "sync/atomic"

type Counters struct {
	Requests   uint64
	OK         uint64
	Exceptions uint64
	OtherErrs  uint64
}

func (c *Counters) IncRequests()   { atomic.AddUint64(&c.Requests, 1) }
func (c *Counters) IncOK()         { atomic.AddUint64(&c.OK, 1) }
func (c *Counters) IncExceptions() { atomic.AddUint64(&c.Exceptions, 1) }
func (c *Counters) IncOtherErrs()  { atomic.AddUint64(&c.OtherErrs, 1) }

func (c *Counters) Snapshot() (req, ok, ex, other uint64) {
	req = atomic.LoadUint64(&c.Requests)
	ok = atomic.LoadUint64(&c.OK)
	ex = atomic.LoadUint64(&c.Exceptions)
	other = atomic.LoadUint64(&c.OtherErrs)
	return
}
