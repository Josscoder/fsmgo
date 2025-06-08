package states

import (
	"fmt"
	"github.com/josscoder/fsmgo/state"
)

type PrintState struct {
	state.BaseState
	Text string
}

func NewPrintState(text string) *PrintState {
	ps := &PrintState{}
	ps.BaseState.Init(ps)
	ps.Text = text
	return ps
}

func (ps *PrintState) OnStart() {
	fmt.Println("Started:", ps.Text)
}

func (ps *PrintState) OnUpdate() {
	fmt.Printf("Updating: %s, remaining %d\n", ps.Text, ps.GetRemainingTime())
}

func (ps *PrintState) OnEnd() {
	fmt.Println("Ended:", ps.Text)
}

func (ps *PrintState) GetDuration() int {
	return 5
}
