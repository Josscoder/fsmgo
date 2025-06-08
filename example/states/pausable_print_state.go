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

// OnPause is optionally called when the state is paused.
// Override this method if you need to handle pause-specific behavior.
func (ps *PrintState) OnPause() {
	fmt.Println("Paused:", ps.Text)
}

// OnResume is optionally called when the state is resumed.
// Override this method if you need to handle resume-specific behavior.
func (ps *PrintState) OnResume() {
	fmt.Println("Resumed:", ps.Text)
}

var _ state.PauseAware = (*PrintState)(nil)

func (ps *PrintState) GetDuration() int {
	return 5
}
