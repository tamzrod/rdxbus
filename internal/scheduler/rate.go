package scheduler

import (
	"time"
)

// Scheduler controls request rate.
type Scheduler struct {
	tokens <-chan struct{}
	stop   chan struct{}
}

// NewScheduler creates a rate scheduler.
// rate = requests per second (0 = unlimited)
func NewScheduler(rate int) *Scheduler {
	if rate <= 0 {
		// Unlimited mode: always ready
		ch := make(chan struct{})
		close(ch)
		return &Scheduler{
			tokens: ch,
			stop:   nil,
		}
	}

	tokens := make(chan struct{}, rate)
	stop := make(chan struct{})

	interval := time.Second / time.Duration(rate)
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokens <- struct{}{}:
				default:
					// drop token if workers are behind
				}
			case <-stop:
				return
			}
		}
	}()

	return &Scheduler{
		tokens: tokens,
		stop:   stop,
	}
}

// Wait blocks until a token is available or returns false if stopped.
func (s *Scheduler) Wait() bool {
	_, ok := <-s.tokens
	return ok
}

// Stop terminates the scheduler.
func (s *Scheduler) Stop() {
	if s.stop != nil {
		close(s.stop)
	}
}
