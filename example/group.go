package main

import (
	"fmt"
	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
	"time"
)

func main() {
	group := state.NewStateGroup([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})
	group.Start()

	for !group.HasEnded() {
		group.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All parallel states finished")
}
