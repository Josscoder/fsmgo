package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

// Group runs all child states concurrently and ends when every one finishes.
func main() {
	group := state.NewStateGroup([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})

	group.Start()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var last time.Time
	for t := range ticker.C {
		delta := time.Second
		if !last.IsZero() {
			delta = t.Sub(last)
		}
		last = t

		group.Update(delta)

		if group.HasEnded() {
			log.Println("All parallel states finished.")
			return
		}
	}
}
