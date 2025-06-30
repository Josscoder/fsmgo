package states

import (
	"fmt"
	"time"

	"github.com/josscoder/fsmgo/state"
)

type PrintState struct {
	*state.BaseState
	Text string
}

func NewPrintState(text string) *PrintState {
	ps := &PrintState{
		Text: text,
	}
	ps.BaseState = state.NewBaseState(ps)
	return ps
}

func (ps *PrintState) OnStart() {
	fmt.Println("Started:", ps.Text)
}

func (ps *PrintState) OnUpdate() {
	fmt.Printf("Updating: %s, remaining %v\n", ps.Text, ps.GetRemainingTime())
}

func (ps *PrintState) OnEnd() {
	fmt.Println("Ended:", ps.Text)
}

func (ps *PrintState) GetDuration() time.Duration {
	return 5 * time.Second
}
