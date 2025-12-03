# Grammar School - Go Implementation

Go implementation of Grammar School, a lightweight framework for building LLM-friendly DSLs.

## Installation

```bash
go get github.com/Conceptual-Machines/grammar-school-go
```

## Usage

```go
import "github.com/Conceptual-Machines/grammar-school-go/gs"

engine, err := gs.NewEngine(grammar, dslInstance, nil)
err = engine.Execute(ctx, dslCode)
```

## Documentation

For complete documentation, see the main [Grammar School](https://github.com/Conceptual-Machines/grammar-school) repository.

## License

AGPL v3 - See LICENSE file for details.
