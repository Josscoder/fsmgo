<div style="text-align: center;">
    <h1>
        <img alt="fsmgo logo" src="/images/logo.png" height="350" /><br />
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
**fsmgo** revolves around the definition of states and their composition. A state has a lifecycle with three phases:

- `OnStart()`


- `OnUpdate()`


- `OnEnd()`

Additionally, each state can have a duration or a custom end condition.

## Examples
Explore how to use **fsmgo** in the [example](https://github.com/Josscoder/fsmgo/tree/master/example) directory.

## License