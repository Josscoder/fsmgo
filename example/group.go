package main

import (
	"fmt"
	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
	"time"
)

func main() {
	wait1 := states.NewWaitingState("1")
	wait2 := states.NewWaitingState("2")

	group := state.NewStateGroup([]state.State{wait1, wait2})
	group.Start()

	for !group.HasEnded() {
		group.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All parallel states finished")
}
