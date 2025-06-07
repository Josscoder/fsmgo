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
	GetRemainingDuration() int
	SetRemainingDuration(int)
	HasStarted() bool
	HasEnded() bool
	IsFrozen() bool
	SetFrozen(bool)
	Freeze()
	Unfreeze()
}

type BaseState struct {
	time     int
	started  bool
	ended    bool
	frozen   bool
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
	s.frozen = false
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

	if s.IsReadyToEnd() && !s.frozen {
		s.End()
		s.updating = false
		return
	}

	if !s.frozen {
		s.time--
	}

	s.self.OnUpdate()
	s.updating = false
}

func (s *BaseState) IsReadyToEnd() bool {
	return s.ended || s.GetRemainingDuration() <= 0
}

func (s *BaseState) GetRemainingDuration() int {
	if s.time < 0 {
		return 0
	}
	return s.time
}

func (s *BaseState) SetRemainingDuration(remaining int) {
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

func (s *BaseState) IsFrozen() bool {
	return s.frozen
}

func (s *BaseState) SetFrozen(frozen bool) {
	s.frozen = frozen
}

func (s *BaseState) Freeze() {
	s.SetFrozen(true)
}

func (s *BaseState) Unfreeze() {
	s.SetFrozen(false)
}
