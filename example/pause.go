package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
)

func main() {
	ps := states.NewPausablePrintState("Test State")

	ps.Start()

	for i := 0; i < 3; i++ {
		ps.Update()
		time.Sleep(1 * time.Second)
	}

	log.Println("Pausing...")
	ps.Pause()

	time.Sleep(2 * time.Second)

	log.Println("Resuming...")
	ps.Resume()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ps.Update()
			if ps.HasEnded() {
				log.Println("State ended")
				return
			}
		}
	}
}
