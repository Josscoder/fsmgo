package main

import (
	"fmt"
	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
	"time"
)

func main() {
	series := state.NewStateSeries([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})
	series.Start()

	for !series.HasEnded() {
		series.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Sequential game state series ended")
}
