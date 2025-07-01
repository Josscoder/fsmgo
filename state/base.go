package state

import (
	"log"
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
	updating  bool
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
	s.remaining = s.lifecycle.GetDuration()
	s.started = false
	s.ended = false
	s.paused = false
	s.updating = false
}

func (s *BaseState) Start() {
	if s.started || s.ended {
		return
	}
	s.started = true

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered during Start(): %v", r)
		}
	}()
	s.lifecycle.OnStart()
}

func (s *BaseState) Update() {
	if !s.started || s.ended || s.updating {
		return
	}
	s.updating = true

	if s.IsReadyToEnd() && !s.paused {
		s.End()

		return
	}

	if !s.paused {
		s.remaining -= time.Second
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered during Update(): %v", r)
		}
		s.updating = false
	}()
	s.lifecycle.OnUpdate()

	s.updating = false
}

func (s *BaseState) End() {
	if !s.started || s.ended {
		return
	}
	s.ended = true
	s.remaining = 0

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered during End(): %v", r)
		}
	}()
	s.lifecycle.OnEnd()
}

func (s *BaseState) Pause() {
	s.SetPaused(true)
}

func (s *BaseState) Resume() {
	s.SetPaused(false)
}

func (s *BaseState) SetPaused(paused bool) {
	if s.paused == paused {
		return
	}
	s.paused = paused

	if s.pauseAware != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered during pause state change: %v", r)
			}
		}()
		if paused {
			s.pauseAware.OnPause()
		} else {
			s.pauseAware.OnResume()
		}
	}
}

func (s *BaseState) HasStarted() bool {
	return s.started
}

func (s *BaseState) HasEnded() bool {
	return s.ended
}

func (s *BaseState) IsPaused() bool {
	return s.paused
}

func (s *BaseState) GetRemainingTime() time.Duration {
	return s.remaining
}

func (s *BaseState) SetRemainingTime(d time.Duration) {
	s.remaining = d
}

func (s *BaseState) IsReadyToEnd() bool {
	return s.ended || s.remaining <= 0
}
