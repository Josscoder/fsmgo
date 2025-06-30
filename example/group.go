package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

func main() {
	group := state.NewStateGroup([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})

	group.Start()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			group.Update()

			if group.HasEnded() {
				log.Println("All parallel states finished")
				return
			}
		}
	}
}
