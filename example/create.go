package main

import (
	"log"
	"time"

	"github.com/josscoder/fsmgo/example/states"
)

func main() {
	ps := states.NewPrintState("Hello World")
	ps.Start()

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
			log.Println("State cycle completed.")
			return
		}
	}
}
