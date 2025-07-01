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
	SetPaused(bool)
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

	mutex sync.Mutex
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.remaining = s.lifecycle.GetDuration()
	s.started = false
	s.ended = false
	s.paused = false
	s.updating = false
}

func (s *BaseState) Start() {
	s.mutex.Lock()
	if s.started || s.ended {
		s.mutex.Unlock()
		return
	}
	s.started = true
	s.mutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered during Start(): %v", r)
		}
	}()
	s.lifecycle.OnStart()
}

func (s *BaseState) Update() {
	s.mutex.Lock()
	if !s.started || s.ended || s.updating {
		s.mutex.Unlock()
		return
	}
	s.updating = true
	s.mutex.Unlock()

	if s.IsReadyToEnd() && !s.paused {
		s.End()
		s.mutex.Lock()
		s.updating = false
		s.mutex.Unlock()
		return
	}

	s.mutex.Lock()
	if !s.paused {
		s.remaining -= time.Second
	}
	s.mutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered during Update(): %v", r)
		}
		s.mutex.Lock()
		s.updating = false
		s.mutex.Unlock()
	}()

	s.lifecycle.OnUpdate()
}

func (s *BaseState) End() {
	s.mutex.Lock()
	if !s.started || s.ended {
		s.mutex.Unlock()
		return
	}
	s.ended = true
	s.remaining = 0
	s.mutex.Unlock()

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
	s.mutex.Lock()
	if s.paused == paused {
		s.mutex.Unlock()
		return
	}
	s.paused = paused
	s.mutex.Unlock()

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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.started
}

func (s *BaseState) HasEnded() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.ended
}

func (s *BaseState) IsPaused() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.paused
}

func (s *BaseState) GetRemainingTime() time.Duration {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.remaining
}

func (s *BaseState) SetRemainingTime(d time.Duration) {
	s.mutex.Lock()
	s.remaining = d
	s.mutex.Unlock()
}

func (s *BaseState) IsReadyToEnd() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.ended || s.remaining <= 0
}
