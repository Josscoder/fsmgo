package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

// Series can be nested: each sub-series is itself a State that runs sequentially
// inside the outer series.
func main() {
	subSeries1 := state.NewStateSeries([]state.State{
		states.NewPrintState("Sub 1-1"),
		states.NewPrintState("Sub 1-2"),
	})

	subSeries2 := state.NewStateSeries([]state.State{
		states.NewPrintState("Sub 2-1"),
		states.NewPrintState("Sub 2-2"),
	})

	mainSeries := state.NewStateSeries([]state.State{subSeries1, subSeries2})
	mainSeries.Start()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var last time.Time
	for t := range ticker.C {
		delta := time.Second
		if !last.IsZero() {
			delta = t.Sub(last)
		}
		last = t

		mainSeries.Update(delta)

		if mainSeries.HasEnded() {
			log.Println("Nested state series ended.")
			return
		}
	}
}
