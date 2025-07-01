package state

import "time"

type Group struct {
	*Holder
}

func NewStateGroup(states []State) *Group {
	group := &Group{
		Holder: NewStateHolder(states),
	}
	group.BaseState = NewBaseState(group)
	return group
}

func (g *Group) OnStart() {
	for _, st := range g.states {
		st.Start()
	}
}

func (g *Group) OnUpdate() {
	for _, st := range g.states {
		st.Update()
	}
	allEnded := true
	for _, st := range g.states {
		if !st.HasEnded() {
			allEnded = false
			break
		}
	}
	if allEnded {
		g.End()
	}
}

func (g *Group) OnEnd() {
	for _, st := range g.states {
		st.End()
	}
}

func (g *Group) IsReadyToEnd() bool {
	for _, st := range g.states {
		if !st.IsReadyToEnd() {
			return false
		}
	}
	return true
}

func (g *Group) GetDuration() time.Duration {
	var maxDur time.Duration
	for _, st := range g.states {
		d := st.(Lifecycle).GetDuration()
		if d > maxDur {
			maxDur = d
		}
	}
	return maxDur
}
