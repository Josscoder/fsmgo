package states

import (
	"fmt"
	"github.com/josscoder/fsmgo/state"
)

type PlayingState struct {
	state.BaseState
	Id string
}

func NewPlayingState(id string) *PlayingState {
	ps := &PlayingState{}
	ps.BaseState.Init(ps)
	ps.Id = id
	return ps
}

func (p *PlayingState) OnStart() {
	fmt.Printf("PlayingState %s: started\n", p.Id)
}

func (p *PlayingState) OnUpdate() {
	fmt.Printf("PlayingState %s: updating, remaining: %d\n", p.Id, p.GetRemainingDuration())
}

func (p *PlayingState) OnEnd() {
	fmt.Printf("PlayingState %s: ended\n", p.Id)
}

func (p *PlayingState) GetDuration() int {
	return 5
}
