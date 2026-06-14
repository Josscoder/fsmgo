package main

import (
	"context"
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
	"github.com/josscoder/fsmgo/state"
)

// ScheduledStateSeries drives itself internally via a ticker — no manual
// Update loop needed. Pass a context to control its lifetime.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	series := state.NewScheduledStateSeries(ctx, []state.State{
		states.NewPrintState("Auto State 1"),
		states.NewPrintState("Auto State 2"),
	}, time.Second)

	series.Start()

	// Poll until the series ends or the context expires.
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if series.HasEnded() {
			log.Println("Scheduled series ended.")
			return
		}
		if ctx.Err() != nil {
			log.Println("Context expired before series finished.")
			return
		}
	}
}
