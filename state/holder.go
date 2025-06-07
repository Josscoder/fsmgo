package state

type Holder struct {
	BaseState
	current int
	states  []State
}

func NewStateHolder(states []State) *Holder {
	holder := &Holder{
		current: 0,
		states:  states,
	}
	holder.BaseState.Init(holder)
	return holder
}

func (s *Holder) OnStart() {
	if curr := s.Current(); curr != nil {
		curr.Start()
	}
}

func (s *Holder) OnUpdate() {
	if curr := s.Current(); curr != nil {
		curr.Update()
	}
}

func (s *Holder) OnEnd() {}

func (s *Holder) GetDuration() int {
	if curr := s.Current(); curr != nil {
		return curr.GetDuration()
	}
	return 0
}

func (s *Holder) Current() State {
	if s.Valid() {
		return s.states[s.current]
	}
	return nil
}

func (s *Holder) Key() int {
	return s.current
}

func (s *Holder) Previous() {
	if curr := s.Current(); curr != nil {
		curr.Cleanup()
	}
	s.current--
	if s.current < 0 {
		s.current = 0
	}
	if curr := s.Current(); curr != nil {
		curr.Cleanup()
		curr.Start()
	}
}

func (s *Holder) Next() {
	s.current++
}

func (s *Holder) Valid() bool {
	return s.current >= 0 && s.current < len(s.states)
}

func (s *Holder) Rewind() {
	s.current = 0
	for _, st := range s.states {
		st.Cleanup()
	}
	if curr := s.Current(); curr != nil {
		curr.Start()
	}
}

func (s *Holder) Add(state State) {
	s.states = append(s.states, state)
}

func (s *Holder) AddAll(states []State) {
	s.states = append(s.states, states...)
}

func (s *Holder) SetPaused(frozen bool) {
	for _, st := range s.states {
		st.SetPaused(frozen)
	}
	s.BaseState.SetPaused(frozen)
}
