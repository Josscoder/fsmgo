package states

import (
	"fmt"
	"github.com/josscoder/fsmgo/state"
)

type WaitingState struct {
	state.BaseState
	Id string
}

func NewWaitingState(id string) *WaitingState {
	ws := &WaitingState{}
	ws.BaseState.Init(ws)
	ws.Id = id
	return ws
}

func (w *WaitingState) OnStart() {
	fmt.Printf("WaitingState %s: started\n", w.Id)
}

func (w *WaitingState) OnUpdate() {
	fmt.Printf("WaitingState %s: updating, remaining: %d\n", w.Id, w.GetRemainingDuration())
}

func (w *WaitingState) OnEnd() {
	fmt.Printf("WaitingState %s: ended\n", w.Id)
}

func (w *WaitingState) GetDuration() int {
	return 3
}
