package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

// Series plays states one after another, advancing automatically when each one ends.
func main() {
	series := state.NewStateSeries([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})

	series.Start()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var last time.Time
	for t := range ticker.C {
		delta := time.Second
		if !last.IsZero() {
			delta = t.Sub(last)
		}
		last = t

		series.Update(delta)

		if series.HasEnded() {
			log.Println("Sequential state series ended.")
			return
		}
	}
}
