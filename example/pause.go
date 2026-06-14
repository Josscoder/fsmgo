package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
)

// PausablePrintState implements PauseAware, so BaseState will call OnPause/OnResume
// automatically whenever SetPaused changes.
func main() {
	ps := states.NewPausablePrintState("Test State")
	ps.Start()

	// Run for 3 ticks.
	for range 3 {
		ps.Update(time.Second)
		time.Sleep(time.Second)
	}

	log.Println("Pausing...")
	ps.Pause()
	time.Sleep(2 * time.Second)

	log.Println("Resuming...")
	ps.Resume()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var last time.Time
	for t := range ticker.C {
		delta := time.Second
		if !last.IsZero() {
			delta = t.Sub(last)
		}
		last = t

		ps.Update(delta)

		if ps.HasEnded() {
			log.Println("State ended.")
			return
		}
	}
}
