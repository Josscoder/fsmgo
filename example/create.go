package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
)

func main() {
	ps := states.NewPrintState("Hello World")

	ps.Start()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ps.Update()

			if ps.HasEnded() {
				log.Println("State cycle completed.")
				return
			}
		}
	}
}
