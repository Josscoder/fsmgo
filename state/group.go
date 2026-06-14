package state

import "time"

// Group runs all child states concurrently and ends when every one of them
// has finished. It no longer embeds Holder because it has no use for sequential
// navigation (current index, Next, Previous, etc.).
type Group struct {
	*BaseState
	states []State
}

func NewStateGroup(states []State) *Group {
	g := &Group{states: states}
	g.BaseState = NewBaseState(g)
	return g
}

func (g *Group) OnStart() {
	for _, st := range g.states {
		st.Start()
	}
}

// OnUpdate ticks every child, then ends the group when all children are done.
func (g *Group) OnUpdate(delta time.Duration) {
	for _, st := range g.states {
		st.Update(delta)
	}
	if g.IsReadyToEnd() {
		g.End()
	}
}

func (g *Group) OnEnd() {
	for _, st := range g.states {
		st.End()
	}
}

// IsReadyToEnd returns true only when every child has finished.
func (g *Group) IsReadyToEnd() bool {
	for _, st := range g.states {
		if !st.IsReadyToEnd() {
			return false
		}
	}
	return true
}

// GetDuration returns the longest child duration (groups run in parallel).
func (g *Group) GetDuration() time.Duration {
	var maxDur time.Duration
	for _, st := range g.states {
		if lc, ok := st.(Lifecycle); ok {
			if d := lc.GetDuration(); d > maxDur {
				maxDur = d
			}
		}
	}
	return maxDur
}

// SetPaused propagates pause/resume to every child.
func (g *Group) SetPaused(paused bool) {
	for _, st := range g.states {
		st.SetPaused(paused)
	}
	g.BaseState.SetPaused(paused)
}

// Add appends a state to the group (before Start is called).
func (g *Group) Add(st State) { g.states = append(g.states, st) }
