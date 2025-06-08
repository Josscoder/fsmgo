package main

import (
	"fmt"
	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
	"time"
)

func main() {
	series := state.NewStateSeries([]state.State{
		state.NewStateSeries([]state.State{
			states.NewPrintState("State sub 1-1"),
			states.NewPrintState("State sub 1-2"),
		}),
		state.NewStateSeries([]state.State{
			states.NewPrintState("State sub 2-1"),
			states.NewPrintState("State sub 2-2"),
		}),
	})
	series.Start()

	for !series.HasEnded() {
		series.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Nested game state series ended")
}
