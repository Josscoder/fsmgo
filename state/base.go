package state

import (
	"log"
	"sync"
	"time"
)

type Lifecycle interface {
	OnStart()
	OnUpdate()
	OnEnd()
	GetDuration() time.Duration
}

type PauseAware interface {
	OnPause()
	OnResume()
}

type State interface {
	Start()
	Update()
	End()
	Pause()
	Resume()
	HasStarted() bool
	HasEnded() bool
	IsPaused() bool
	GetRemainingTime() time.Duration
	SetRemainingTime(time.Duration)
	IsReadyToEnd() bool
	Cleanup()
}

type BaseState struct {
	lifecycle  Lifecycle
	pauseAware PauseAware

	remaining time.Duration
	started   bool
	ended     bool
	paused    bool

	mu sync.Mutex
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

func (s *BaseState) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.remaining = s.lifecycle.GetDuration()
	s.started = false
	s.ended = false
	s.paused = false
}

func (s *BaseState) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started || s.ended {
		return
	}
	s.started = true
	defer s.recoverPanic("Start")
	s.lifecycle.OnStart()
}

func (s *BaseState) Update() {
	s.mu.Lock()
	if !s.started || s.ended || s.paused {
		s.mu.Unlock()
		return
	}
	if s.remaining <= 0 {
		s.mu.Unlock()
		s.End()
		return
	}
	s.remaining -= time.Second
	s.mu.Unlock()

	defer s.recoverPanic("Update")
	s.lifecycle.OnUpdate()
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

	defer s.recoverPanic("End")
	s.lifecycle.OnEnd()
}

func (s *BaseState) Pause() {
	s.setPaused(true)
}

func (s *BaseState) Resume() {
	s.setPaused(false)
}

func (s *BaseState) setPaused(paused bool) {
	s.mu.Lock()
	if s.paused == paused {
		s.mu.Unlock()
		return
	}
	s.paused = paused
	s.mu.Unlock()

	if s.pauseAware != nil {
		defer s.recoverPanic("Pause/Resume")
		if paused {
			s.pauseAware.OnPause()
		} else {
			s.pauseAware.OnResume()
		}
	}
}

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
	defer s.mu.Unlock()
	s.remaining = d
}

func (s *BaseState) IsReadyToEnd() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ended || s.remaining <= 0
}

func (s *BaseState) recoverPanic(context string) {
	if r := recover(); r != nil {
		log.Printf("[State] Recovered from panic in %s: %v", context, r)
	}
}
