<p align="center">
  <img src="/images/logo.png" alt="fsmgo logo" width="300"/>
</p>

<h2 align="center">Finite State Machine for Go</h2>

**fsmgo** is a clean, modular, and concurrent Go library for building finite state machines, inspired by [FSMgasm](https://github.com/Minikloon/FSMgasm) and adapted with Go idioms for safe, real-time systems.

## Installation

```sh
go get github.com/josscoder/fsmgo
```

```go
import "github.com/josscoder/fsmgo/state"
```

---

## Concepts

### Lifecycle

Every state has three lifecycle hooks:

| Hook | When it's called |
|------|-----------------|
| `OnStart()` | Once, when the state begins |
| `OnUpdate(delta time.Duration)` | Every tick; receives real elapsed time |
| `OnEnd()` | Once, when the state finishes |

`delta` lets your state logic be frame-rate independent — subtract it from timers or use it to interpolate values rather than assuming a fixed interval.

### PauseAware (optional)

Implement the `PauseAware` interface to receive pause/resume notifications. `BaseState` detects it automatically at construction time — no extra wiring needed.

```go
type PauseAware interface {
    OnPause()
    OnResume()
}
```

---

## Creating a state

```go
package states

import (
    "fmt"
    "time"

    "github.com/josscoder/fsmgo/state"
)

type PrintState struct {
    *state.BaseState
    Text string
}

func NewPrintState(text string) *PrintState {
    ps := &PrintState{Text: text}
    ps.BaseState = state.NewBaseState(ps)
    return ps
}

func (ps *PrintState) OnStart()                       { fmt.Println("Started:", ps.Text) }
func (ps *PrintState) OnUpdate(delta time.Duration)   { fmt.Printf("tick %v — remaining: %v\n", delta, ps.GetRemainingTime()) }
func (ps *PrintState) OnEnd()                         { fmt.Println("Ended:", ps.Text) }
func (ps *PrintState) GetDuration() time.Duration     { return 5 * time.Second }
```

---

## Running a single state

```go
ps := states.NewPrintState("Hello World")
ps.Start()

ticker := time.NewTicker(time.Second)
defer ticker.Stop()

var last time.Time
for t := range ticker.C {
    delta := time.Second
    if !last.IsZero() {
        delta = t.Sub(last) // real elapsed time, absorbs scheduler jitter
    }
    last = t

    ps.Update(delta)

    if ps.HasEnded() {
        return
    }
}
```

---

## Containers

### Series — sequential

States run one after another. When the current state ends, the next one starts automatically.

```go
series := state.NewStateSeries([]state.State{
    states.NewPrintState("Step 1"),
    states.NewPrintState("Step 2"),
})
series.Start()
// ... tick loop calling series.Update(delta)
```

Series can be **nested**: a `Series` is itself a `State`, so it can be placed inside another `Series` or `Group`.

#### Dynamic insertion

```go
series.Skip()                     // skip the current state on the next tick
series.AddNext(newState)          // insert a state right after the current one
series.AddNextList([]state.State{ // insert multiple states after current
    stateA, stateB,
})
```

### Group — parallel

All child states run concurrently. The group ends when every child has finished.

```go
group := state.NewStateGroup([]state.State{
    states.NewPrintState("Worker A"),
    states.NewPrintState("Worker B"),
})
group.Start()
// ... tick loop calling group.Update(delta)
```

### ScheduledStateSeries — self-ticking series

Drives its own internal ticker so you don't need a manual update loop. Pass a `context.Context` to control its lifetime.

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

series := state.NewScheduledStateSeries(ctx, []state.State{
    states.NewPrintState("Auto 1"),
    states.NewPrintState("Auto 2"),
}, time.Second)

series.Start()

// Just wait — no Update calls needed.
for !series.HasEnded() {
    time.Sleep(200 * time.Millisecond)
}
```

---

## Pause / Resume

Any state (or container) can be paused. Pause propagates to all children in a `Group`, `Holder`, or `Series`. If the state also implements `PauseAware`, `OnPause` / `OnResume` are called automatically.

```go
ps := states.NewPausablePrintState("My State")
ps.Start()

// ... run a few ticks ...

ps.Pause()   // countdown freezes; OnPause() called if PauseAware
time.Sleep(2 * time.Second)
ps.Resume()  // countdown resumes; OnResume() called if PauseAware
```

---

## Thread safety

All public methods on `BaseState` and the built-in containers are safe to call from multiple goroutines. `ScheduledStateSeries` uses a `context.Context` for cancellation and `sync.Once` to ensure `OnEnd` fires exactly once.

---

## License

**fsmgo** is licensed under the [MIT License](./LICENSE).