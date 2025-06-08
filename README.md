<div style="text-align: center;">
  <h1>
    <img alt="fsmgo logo" src="/images/logo.png" style="max-width: 70%; height: auto;" />
    <br />
    Finite State Machine for Go
  </h1>
</div>

**fsmgo** is a Go library for building finite state machines in a clear, modular, and reusable way. It's inspired by [FSMgasm](https://github.com/Minikloon/FSMgasm), a Kotlin library originally designed for Minecraft minigames, but fsmgo has been adapted with Go idioms to fit naturally into concurrent architectures and complex systems.


## Installation

```shell
go get github.com/josscoder/fsmgo
```

Import the package into your code:

```go
import "github.com/josscoder/fsmgo/state"
```

## Concepts
**fsmgo** revolves around the definition of states and their composition. A state has a lifecycle with three main phases:

- `OnStart()`

- `OnUpdate()`

- `OnEnd()`

Additionally, states can implement optional lifecycle hooks for pause and resume:

- `OnPause()`

- `OnResume()`

These two methods are not required as part of the main State interface. Instead, they belong to an optional interface named PauseAware:

```go
type PauseAware interface {
    OnPause()
    OnResume()
}
```
When a state implements PauseAware, the FSM system will automatically call `OnPause()` and `OnResume()` when appropriate. States that don't require pause/resume behavior can ignore this interface entirely.

Each state can also have a duration or a custom end condition to determine when it completes.

## Examples
Explore how to use **fsmgo** in the [example](https://github.com/Josscoder/fsmgo/tree/master/example) directory.

## License
**fsmgo** is licensed under the [MIT License](./LICENSE). Feel free to use, modify, and distribute it in your projects.