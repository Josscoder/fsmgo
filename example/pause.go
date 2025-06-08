package main

import (
	"fmt"
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

	fmt.Println("Pausing...")
	ps.Pause()

	time.Sleep(2 * time.Second)

	fmt.Println("Resuming...")
	ps.Resume()

	for !ps.HasEnded() {
		ps.Update()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("State ended")
}
