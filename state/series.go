package state

import (
	"sync"
	"time"
)

// Series plays states sequentially, advancing to the next one automatically
// when the current state is ready to end.
type Series struct {
	*Holder
	skipMu   sync.Mutex
	skipping bool
}

func NewStateSeries(states []State) *Series {
	s := &Series{
		Holder: NewStateHolder(states),
	}
	s.BaseState = NewBaseState(s)
	return s
}

func (s *Series) OnStart() {
	if len(s.states) == 0 {
		s.End()
		return
	}
	if curr := s.Current(); curr != nil {
		curr.Start()
	}
}

func (s *Series) OnUpdate(delta time.Duration) {
	curr := s.Current()
	if curr == nil {
		return
	}

	curr.Update(delta)

	s.skipMu.Lock()
	skip := s.skipping
	if skip {
		s.skipping = false
	}
	s.skipMu.Unlock()

	if (curr.IsReadyToEnd() && !curr.IsPaused()) || skip {
		curr.End()
		s.Next()

		if !s.Valid() {
			s.End()
			return
		}
		if next := s.Current(); next != nil {
			next.Start()
		}
	}
}

func (s *Series) OnEnd() {
	if curr := s.Current(); curr != nil {
		curr.End()
	}
}

func (s *Series) IsReadyToEnd() bool {
	curr := s.Current()
	if curr == nil {
		return true
	}
	return curr.IsReadyToEnd()
}

// GetDuration returns the sum of all child durations.
func (s *Series) GetDuration() time.Duration {
	s.Holder.mu.Lock()
	states := s.states
	s.Holder.mu.Unlock()

	var total time.Duration
	for _, st := range states {
		if lc, ok := st.(Lifecycle); ok {
			total += lc.GetDuration()
		}
	}
	return total
}

// Skip requests that the current state be skipped on the next Update tick.
// Safe to call from any goroutine.
func (s *Series) Skip() {
	s.skipMu.Lock()
	s.skipping = true
	s.skipMu.Unlock()
}

// AddNext inserts a state immediately after the current position.
func (s *Series) AddNext(state State) {
	s.Holder.mu.Lock()
	defer s.Holder.mu.Unlock()

	idx := s.current + 1
	if idx >= len(s.states) {
		s.states = append(s.states, state)
		return
	}
	// Correct insertion: grow the slice then shift elements right.
	s.states = append(s.states, nil)
	copy(s.states[idx+1:], s.states[idx:])
	s.states[idx] = state
}

// AddNextList inserts multiple states after the current position, in order.
func (s *Series) AddNextList(states []State) {
	for _, st := range states {
		s.AddNext(st)
	}
}
