package state

import (
	"context"
	"sync"
	"time"
)

// ScheduledStateSeries is a Series that drives itself via an internal ticker,
// so callers don't need an external game-loop to call Update.
//
// Cancellation is handled through context: pass a cancellable context and the
// goroutine will exit cleanly when the context is done.
type ScheduledStateSeries struct {
	*Series
	interval time.Duration
	cancel   context.CancelFunc
	stopOnce sync.Once
}

// NewScheduledStateSeries creates a series that ticks every `interval`.
// Pass context.Background() if you don't need external cancellation.
func NewScheduledStateSeries(ctx context.Context, states []State, interval time.Duration) *ScheduledStateSeries {
	ctx, cancel := context.WithCancel(ctx)
	s := &ScheduledStateSeries{
		Series:   NewStateSeries(states),
		interval: interval,
		cancel:   cancel,
	}
	s.BaseState = NewBaseState(s)

	go s.loop(ctx)
	return s
}

func (s *ScheduledStateSeries) loop(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	var last time.Time
	for {
		select {
		case t := <-ticker.C:
			delta := s.interval
			if !last.IsZero() {
				delta = t.Sub(last) // real elapsed time, absorbs jitter
			}
			last = t
			// Call the public Update so BaseState guards (started, ended,
			// updating) are respected — never bypass them.
			s.Update(delta)

		case <-ctx.Done():
			return
		}
	}
}

// OnStart delegates to Series and does NOT spawn a goroutine — the goroutine
// is started in the constructor so its lifetime is tied to context, not to
// the Start/End cycle (which can be replayed via Cleanup).
func (s *ScheduledStateSeries) OnStart() {
	s.Series.OnStart()
}

// OnEnd stops the ticker goroutine exactly once via context cancellation.
func (s *ScheduledStateSeries) OnEnd() {
	s.stopOnce.Do(s.cancel)
	s.Series.OnEnd()
}
