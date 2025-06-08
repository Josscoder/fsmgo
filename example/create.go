package main

import (
	"github.com/josscoder/fsmgo/example/states"
	"time"
)

func main() {
	ps := states.NewPrintState("Hello World")
	ps.Start()

	for !ps.HasEnded() {
		ps.Update()
		time.Sleep(1 * time.Second)
	}

	ps.End()
}
