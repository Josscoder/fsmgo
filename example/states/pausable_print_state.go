package states

import (
	"fmt"
	"github.com/josscoder/fsmgo/state"
)

type PausablePrintState struct {
	state.BaseState
	Text string
}

func NewPausablePrintState(text string) *PausablePrintState {
	ps := &PausablePrintState{}
	ps.BaseState.Init(ps)
	ps.Text = text
	return ps
}

func (ps *PausablePrintState) OnStart() {
	fmt.Println("Started:", ps.Text)
}

func (ps *PausablePrintState) OnUpdate() {
	fmt.Printf("Updating: %s, remaining %d\n", ps.Text, ps.GetRemainingTime())
}

func (ps *PausablePrintState) OnEnd() {
	fmt.Println("Ended:", ps.Text)
}

// OnPause is optionally called when the state is paused.
// Override this method if you need to handle pause-specific behavior.
func (ps *PausablePrintState) OnPause() {
	fmt.Println("Paused:", ps.Text)
}

// OnResume is optionally called when the state is resumed.
// Override this method if you need to handle resume-specific behavior.
func (ps *PausablePrintState) OnResume() {
	fmt.Println("Resumed:", ps.Text)
}

var _ state.PauseAware = (*PausablePrintState)(nil)

func (ps *PausablePrintState) GetDuration() int {
	return 5
}
