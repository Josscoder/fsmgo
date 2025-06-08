package state

type Series struct {
	Holder
	skipping bool
}

func NewStateSeries(states []State) *Series {
	s := &Series{}
	s.states = states
	s.BaseState.Init(s)
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

func (s *Series) OnUpdate() {
	curr := s.Current()
	if curr == nil {
		return
	}

	curr.Update()

	readyToEnd := curr.IsReadyToEnd()
	isFrozen := curr.IsPaused()

	if (readyToEnd && !isFrozen) || s.skipping {
		if s.skipping {
			s.skipping = false
		}
		curr.End()
		s.Next()

		if s.Key() >= len(s.states) {
			s.End()
			return
		}
		if newCurr := s.Current(); newCurr != nil {
			newCurr.Start()
		}
	}
}

func (s *Series) OnEnd() {
	if s.Key() < len(s.states) {
		if curr := s.Current(); curr != nil {
			curr.End()
		}
	}
}

func (s *Series) IsReadyToEnd() bool {
	curr := s.Current()
	if curr == nil {
		return true
	}
	return curr.IsReadyToEnd()
}

func (s *Series) GetDuration() int {
	total := 0
	for _, st := range s.states {
		total += st.GetDuration()
	}
	return total
}

func (s *Series) Skip() {
	s.skipping = true
}

func (s *Series) AddNext(state State) {
	idx := s.Key() + 1
	if idx >= len(s.states) {
		s.states = append(s.states, state)
	} else {
		s.states = append(s.states[:idx+1], s.states[idx:]...)
		s.states[idx] = state
	}
}

func (s *Series) AddNextList(states []State) {
	for _, st := range states {
		s.AddNext(st)
	}
}
