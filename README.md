<p align="center">
  <img src="/images/logo.png" alt="fsmgo logo" width="300"/>
</p>

<h2 align="center">Finite State Machine for Go</h2>

**fsmgo** is a clean, modular, and reusable Go library for building finite state machines, inspired by [FSMgasm](https://github.com/Minikloon/FSMgasm) but adapted with Go idioms for concurrent and complex systems.

## Installation

```sh
go get github.com/josscoder/fsmgo
```

Then import:

```go
import "github.com/josscoder/fsmgo/state"
```

## Concepts

A **fsmgo** state has a clear lifecycle:

- `OnStart()` – called when the state begins.
- `OnUpdate()` – called every tick or cycle.
- `OnEnd()` – called when the state ends.

Optional pause lifecycle:

- `OnPause()`
- `OnResume()`

States can implement these by using the optional `PauseAware` interface:

```go
type PauseAware interface {
    OnPause()
    OnResume()
}
```

States can define a duration (`time.Duration`) or custom conditions to determine when they complete.

## Examples
See practical examples in the [example](https://github.com/josscoder/fsmgo/tree/master/example) folder.

## License
**fsmgo** is licensed under the [MIT License](./LICENSE). Feel free to use, modify, and distribute it in your projects.