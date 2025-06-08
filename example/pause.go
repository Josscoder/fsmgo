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

	tick := 0
	for !series.HasEnded() {
		tick++
		fmt.Printf("Tick #%d\n", tick)

		series.Update()

		curr := series.Current()
		if curr != nil {
			fmt.Printf("Progress: %.2f%%\n", curr.GetProgress()*100)
		}

		if tick == 2 {
			fmt.Println("===> Pausing")
			series.Pause()
		}

		if tick == 4 {
			fmt.Println("===> Resuming")
			series.Resume()
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Game state series ended")
}
