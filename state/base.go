package state

type State interface {
	OnStart()
	OnUpdate()
	OnEnd()
	GetDuration() int

	Start()
	Update()
	End()
	Cleanup()
	IsReadyToEnd() bool
	GetRemainingTime() int
	SetRemainingTime(int)
	HasStarted() bool
	HasEnded() bool
	IsPaused() bool
	SetPaused(bool)
	Pause()
	Resume()
}

type BaseState struct {
	time     int
	started  bool
	ended    bool
	paused   bool
	updating bool
	self     State
}

func (s *BaseState) Init(self State) {
	s.self = self
	s.Cleanup()
}

func (s *BaseState) Cleanup() {
	s.time = s.self.GetDuration()
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
	s.self.OnStart()
}

func (s *BaseState) Update() {
	if !s.started || s.ended || s.updating {
		return
	}
	s.updating = true

	if s.IsReadyToEnd() && !s.paused {
		s.End()
		s.updating = false
		return
	}

	if !s.paused {
		s.time--
	}

	s.self.OnUpdate()
	s.updating = false
}

func (s *BaseState) IsReadyToEnd() bool {
	return s.ended || s.GetRemainingTime() <= 0
}

func (s *BaseState) GetRemainingTime() int {
	if s.time < 0 {
		return 0
	}
	return s.time
}

func (s *BaseState) SetRemainingTime(remaining int) {
	s.time = remaining
}

func (s *BaseState) End() {
	if !s.started || s.ended {
		return
	}
	s.ended = true
	s.time = 0
	s.self.OnEnd()
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

func (s *BaseState) SetPaused(frozen bool) {
	s.paused = frozen
}

func (s *BaseState) Pause() {
	s.SetPaused(true)
}

func (s *BaseState) Resume() {
	s.SetPaused(false)
}
