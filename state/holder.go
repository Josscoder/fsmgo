package state

import (
	"sync"
	"time"
)

// Holder manages a slot that points to one active State at a time.
// Navigation (Next, Previous, Rewind) is safe to call from any goroutine.
type Holder struct {
	*BaseState
	mu      sync.Mutex
	current int
	states  []State
}

func NewStateHolder(states []State) *Holder {
	h := &Holder{
		current: 0,
		states:  states,
	}
	h.BaseState = NewBaseState(h)
	return h
}

func (h *Holder) OnStart() {
	if curr := h.Current(); curr != nil {
		curr.Start()
	}
}

func (h *Holder) OnUpdate(delta time.Duration) {
	if curr := h.Current(); curr != nil {
		curr.Update(delta)
	}
}

func (h *Holder) OnEnd() {}

func (h *Holder) GetDuration() time.Duration {
	if curr := h.Current(); curr != nil {
		if lc, ok := curr.(Lifecycle); ok {
			return lc.GetDuration()
		}
	}
	return 0
}

// Current returns the active state, or nil if the index is out of bounds.
func (h *Holder) Current() State {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.validLocked() {
		return h.states[h.current]
	}
	return nil
}

// Key returns the current index.
func (h *Holder) Key() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.current
}

// Next advances the index by one (no bounds check — caller drives transitions).
func (h *Holder) Next() {
	h.mu.Lock()
	h.current++
	h.mu.Unlock()
}

// Previous steps back one, cleans up both states, and starts the prior one.
// It clamps at index 0.
func (h *Holder) Previous() {
	h.mu.Lock()
	curr := h.currentLocked()
	h.mu.Unlock()

	if curr != nil {
		curr.Cleanup()
	}

	h.mu.Lock()
	if h.current > 0 {
		h.current--
	}
	prev := h.currentLocked()
	h.mu.Unlock()

	if prev != nil {
		prev.Cleanup()
		prev.Start()
	}
}

// Valid reports whether the current index points to an existing state.
func (h *Holder) Valid() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.validLocked()
}

// Rewind resets the cursor to 0, cleans up all states, and starts index 0.
func (h *Holder) Rewind() {
	h.mu.Lock()
	h.current = 0
	states := h.states
	h.mu.Unlock()

	for _, st := range states {
		st.Cleanup()
	}
	if curr := h.Current(); curr != nil {
		curr.Start()
	}
}

func (h *Holder) Add(state State) {
	h.mu.Lock()
	h.states = append(h.states, state)
	h.mu.Unlock()
}

func (h *Holder) AddAll(states []State) {
	h.mu.Lock()
	h.states = append(h.states, states...)
	h.mu.Unlock()
}

// SetPaused propagates pause to every managed state and to BaseState.
func (h *Holder) SetPaused(paused bool) {
	h.mu.Lock()
	states := h.states
	h.mu.Unlock()

	for _, st := range states {
		st.SetPaused(paused)
	}
	h.BaseState.SetPaused(paused)
}

// ── internal helpers (must be called with mu held) ──────────────────────────

func (h *Holder) validLocked() bool {
	return h.current >= 0 && h.current < len(h.states)
}

func (h *Holder) currentLocked() State {
	if h.validLocked() {
		return h.states[h.current]
	}
	return nil
}
