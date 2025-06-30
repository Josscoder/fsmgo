package state

import "time"

type Holder struct {
	*BaseState
	current int
	states  []State
}

func NewStateHolder(states []State) *Holder {
	holder := &Holder{
		current: 0,
		states:  states,
	}
	holder.BaseState = NewBaseState(holder)
	return holder
}

func (h *Holder) OnStart() {
	if curr := h.Current(); curr != nil {
		curr.Start()
	}
}

func (h *Holder) OnUpdate() {
	if curr := h.Current(); curr != nil {
		curr.Update()
	}
}

func (h *Holder) OnEnd() {}

func (h *Holder) GetDuration() time.Duration {
	if curr := h.Current(); curr != nil {
		return curr.GetRemainingTime()
	}
	return 0
}

func (h *Holder) Current() State {
	if h.Valid() {
		return h.states[h.current]
	}
	return nil
}

func (h *Holder) Key() int {
	return h.current
}

func (h *Holder) Previous() {
	if curr := h.Current(); curr != nil {
		curr.Cleanup()
	}
	h.current--
	if h.current < 0 {
		h.current = 0
	}
	if curr := h.Current(); curr != nil {
		curr.Cleanup()
		curr.Start()
	}
}

func (h *Holder) Next() {
	h.current++
}

func (h *Holder) Valid() bool {
	return h.current >= 0 && h.current < len(h.states)
}

func (h *Holder) Rewind() {
	h.current = 0
	for _, st := range h.states {
		st.Cleanup()
	}
	if curr := h.Current(); curr != nil {
		curr.Start()
	}
}

func (h *Holder) Add(state State) {
	h.states = append(h.states, state)
}

func (h *Holder) AddAll(states []State) {
	h.states = append(h.states, states...)
}

func (h *Holder) SetPaused(paused bool) {
	for _, st := range h.states {
		st.Pause()
	}
	if !paused {
		for _, st := range h.states {
			st.Resume()
		}
	}
	h.BaseState.setPaused(paused)
}
