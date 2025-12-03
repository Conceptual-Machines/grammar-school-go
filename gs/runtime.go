package gs

import (
	"context"
	"fmt"
)

// Runtime is the interface for runtime implementations that execute actions.
//
// The Runtime is responsible for taking Actions (produced by verb handlers) and
// performing the actual side effects - printing, database operations, API calls,
// file I/O, etc.
//
// This separation allows:
// - Verb handlers to be pure (just return Action data structures)
// - Runtime to handle all side effects and state management
// - Same Engine can work with different Runtimes (testing, production, etc.)
type Runtime interface {
	ExecuteAction(ctx context.Context, a Action) error
}

// DefaultRuntime is a default runtime that prints actions to stdout.
// This is used when no runtime is provided to NewEngine, matching Python's behavior.
// Output goes to standard output (stdout) - typically the console/terminal.
//
// For custom output destinations (files, databases, APIs, etc.),
// create a custom Runtime implementation.
type DefaultRuntime struct{}

// ExecuteAction prints the action to stdout (standard output/console).
func (r *DefaultRuntime) ExecuteAction(ctx context.Context, a Action) error {
	fmt.Printf("Action: %s with payload: %v\n", a.Kind, a.Payload)
	return nil
}
