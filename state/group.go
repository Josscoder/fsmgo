package state

type Group struct {
	Holder
}

func NewStateGroup(states []State) *Group {
	group := &Group{}
	group.states = states
	group.BaseState.Init(group)
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

func (g *Group) GetDuration() int {
	maxDur := 0
	for _, st := range g.states {
		d := st.GetDuration()
		if d > maxDur {
			maxDur = d
		}
	}
	return maxDur
}
