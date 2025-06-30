package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

func main() {
	series := state.NewStateSeries([]state.State{
		states.NewPrintState("State 1"),
		states.NewPrintState("State 2"),
	})

	series.Start()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			series.Update()

			if series.HasEnded() {
				log.Println("Sequential state series ended")
				return
			}
		}
	}
}
