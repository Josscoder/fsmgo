package states

import (
	"fmt"
	"time"

	"github.com/josscoder/fsmgo/state"
)

type PausablePrintState struct {
	*state.BaseState
	Text string
}

func NewPausablePrintState(text string) *PausablePrintState {
	ps := &PausablePrintState{
		Text: text,
	}
	ps.BaseState = state.NewBaseState(ps)
	return ps
}

func (ps *PausablePrintState) OnStart() {
	fmt.Println("Started:", ps.Text)
}

func (ps *PausablePrintState) OnUpdate() {
	fmt.Printf("Updating: %s, remaining %v\n", ps.Text, ps.GetRemainingTime())
}

func (ps *PausablePrintState) OnEnd() {
	fmt.Println("Ended:", ps.Text)
}

func (ps *PausablePrintState) OnPause() {
	fmt.Println("Paused:", ps.Text)
}

func (ps *PausablePrintState) OnResume() {
	fmt.Println("Resumed:", ps.Text)
}

func (ps *PausablePrintState) GetDuration() time.Duration {
	return 5 * time.Second
}
