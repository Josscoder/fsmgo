package main

import (
	"fmt"
	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
	"time"
)

func main() {
	wait := states.NewWaitingState("1")
	play := states.NewPlayingState("1")

	series := state.NewStateSeries([]state.State{wait, play})
	series.Start()

	for !series.HasEnded() {
		series.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Game state series ended")
}
