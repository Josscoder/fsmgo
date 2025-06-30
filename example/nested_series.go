package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

func main() {
	subSeries1 := state.NewStateSeries([]state.State{
		states.NewPrintState("State sub 1-1"),
		states.NewPrintState("State sub 1-2"),
	})

	subSeries2 := state.NewStateSeries([]state.State{
		states.NewPrintState("State sub 2-1"),
		states.NewPrintState("State sub 2-2"),
	})

	mainSeries := state.NewStateSeries([]state.State{subSeries1, subSeries2})

	mainSeries.Start()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mainSeries.Update()

			if mainSeries.HasEnded() {
				log.Println("Nested state series ended")
				return
			}
		}
	}
}
