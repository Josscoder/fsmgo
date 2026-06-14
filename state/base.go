package state

import (
	"log/slog"
	"sync"
	"time"
)

// Lifecycle defines the hooks a concrete state must implement.
type Lifecycle interface {
	OnStart()
	OnUpdate(delta time.Duration)
	OnEnd()
	GetDuration() time.Duration
}

// PauseAware is an optional extension for states that react to pause/resume.
type PauseAware interface {
	OnPause()
	OnResume()
}

// State is the public interface consumed by containers (Series, Group, Holder).
type State interface {
	Start()
	Update(delta time.Duration)
	End()
	SetPaused(bool)
	HasStarted() bool
	HasEnded() bool
	IsPaused() bool
	GetRemainingTime() time.Duration
	SetRemainingTime(time.Duration)
	IsReadyToEnd() bool
	Cleanup()
}

// BaseState is the thread-safe skeleton every concrete state embeds.
type BaseState struct {
	lifecycle  Lifecycle
	pauseAware PauseAware

	remaining time.Duration
	started   bool
	ended     bool
	paused    bool
	updating  bool

	mu      sync.Mutex
	endOnce sync.Once // guards OnEnd so it fires exactly once
}

func NewBaseState(l Lifecycle) *BaseState {
	bs := &BaseState{
		lifecycle: l,
		remaining: l.GetDuration(),
	}
	if pa, ok := l.(PauseAware); ok {
		bs.pauseAware = pa
	}
	return bs
}

// Cleanup resets all runtime fields so the state can be replayed.
func (s *BaseState) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.remaining = s.lifecycle.GetDuration()
	s.started = false
	s.ended = false
	s.paused = false
	s.updating = false
	s.endOnce = sync.Once{} // allow OnEnd to fire again after Cleanup
}

func (s *BaseState) Start() {
	s.mu.Lock()
	if s.started || s.ended {
		s.mu.Unlock()
		return
	}
	s.started = true
	s.mu.Unlock()

	safeCall("Start", s.lifecycle.OnStart)
}

// Update advances the state by delta. Callers (e.g. ScheduledStateSeries)
// pass the real elapsed time so the countdown is accurate regardless of
// scheduler jitter.
func (s *BaseState) Update(delta time.Duration) {
	s.mu.Lock()
	if !s.started || s.ended || s.updating {
		s.mu.Unlock()
		return
	}
	s.updating = true
	paused := s.paused
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.updating = false
		s.mu.Unlock()
	}()

	// Tick the countdown before checking readiness so a delta that exhausts
	// remaining time is handled in the same cycle.
	if !paused {
		s.mu.Lock()
		if s.remaining > 0 {
			s.remaining -= delta
			if s.remaining < 0 {
				s.remaining = 0
			}
		}
		s.mu.Unlock()
	}

	if s.isReadyToEndLocked() && !paused {
		s.End()
		return
	}

	safeCall("Update", func() { s.lifecycle.OnUpdate(delta) })
}

func (s *BaseState) End() {
	s.mu.Lock()
	if !s.started || s.ended {
		s.mu.Unlock()
		return
	}
	s.ended = true
	s.remaining = 0
	s.mu.Unlock()

	s.endOnce.Do(func() {
		safeCall("End", s.lifecycle.OnEnd)
	})
}

// SetPaused toggles the paused flag and, if the lifecycle is PauseAware,
// notifies it. Safe to call from any goroutine.
func (s *BaseState) SetPaused(paused bool) {
	s.mu.Lock()
	if s.paused == paused {
		s.mu.Unlock()
		return
	}
	s.paused = paused
	s.mu.Unlock()

	if s.pauseAware == nil {
		return
	}
	if paused {
		safeCall("OnPause", s.pauseAware.OnPause)
	} else {
		safeCall("OnResume", s.pauseAware.OnResume)
	}
}

func (s *BaseState) Pause()  { s.SetPaused(true) }
func (s *BaseState) Resume() { s.SetPaused(false) }

func (s *BaseState) HasStarted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

func (s *BaseState) HasEnded() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ended
}

func (s *BaseState) IsPaused() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paused
}

func (s *BaseState) GetRemainingTime() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remaining
}

func (s *BaseState) SetRemainingTime(d time.Duration) {
	s.mu.Lock()
	s.remaining = d
	s.mu.Unlock()
}

func (s *BaseState) IsReadyToEnd() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isReadyToEndLocked()
}

// isReadyToEndLocked must be called with mu held.
func (s *BaseState) isReadyToEndLocked() bool {
	return s.ended || s.remaining <= 0
}

// safeCall executes fn and recovers from any panic, logging context info.
func safeCall(phase string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recovered", "phase", phase, "recover", r)
		}
	}()
	fn()
}
