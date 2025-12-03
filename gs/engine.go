package gs

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Engine is the main Grammar School engine that orchestrates parsing, interpretation, and execution.
// It corresponds to the Grammar class in Python.
//
// The Engine uses a unified interface:
// - Methods execute directly when called
// - Methods can have side effects and manage state via struct fields
// - No Runtime needed - methods contain their implementation
type Engine struct {
	grammar string
	parser  Parser
	methods map[string]MethodHandler
	dsl     interface{}
}

// NewEngine creates a new Engine with the given grammar, DSL instance, and parser.
// Methods on the DSL struct are automatically discovered and registered.
func NewEngine(grammar string, dsl interface{}, parser Parser) (*Engine, error) {
	engine := &Engine{
		grammar: grammar,
		parser:  parser,
		methods: make(map[string]MethodHandler),
		dsl:     dsl,
	}

	if err := engine.collectMethods(); err != nil {
		return nil, fmt.Errorf("failed to collect methods: %w", err)
	}

	return engine, nil
}

// collectMethods uses reflection to find all methods on the DSL instance
// that match the MethodHandler signature and register them.
// Method signature: func (d *MyDSL) MethodName(args Args) error
func (e *Engine) collectMethods() error {
	dslType := reflect.TypeOf(e.dsl)
	dslValue := reflect.ValueOf(e.dsl)

	// Get methods from the pointer type to access all methods (including pointer receivers)
	var methodsType reflect.Type
	var methodsValue reflect.Value

	if dslType.Kind() == reflect.Ptr {
		// Already a pointer, use it directly
		methodsType = dslType
		methodsValue = dslValue
	} else {
		// Not a pointer, get methods from pointer type
		methodsType = reflect.PtrTo(dslType)
		// Create a pointer to the value
		ptr := reflect.New(dslType)
		ptr.Elem().Set(dslValue)
		methodsValue = ptr
	}

	for i := 0; i < methodsType.NumMethod(); i++ {
		method := methodsType.Method(i)
		methodType := method.Type

		// Method signature: func (receiver) MethodName(args Args) error
		// NumIn: 2 (receiver + args)
		// NumOut: 1 (error)
		if methodType.NumIn() != 2 || methodType.NumOut() != 1 {
			continue
		}

		// Check second parameter is Args
		if methodType.In(1) != reflect.TypeOf(Args{}) {
			continue
		}

		// Check return type is error
		if methodType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		// Register the method
		// Capture method name and value in closure
		methodName := method.Name
		methodValue := methodsValue.Method(i)
		e.methods[methodName] = func(args Args) error {
			results := methodValue.Call([]reflect.Value{
				reflect.ValueOf(args),
			})

			if !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil
		}
	}

	return nil
}

// Execute parses and executes DSL code by calling methods directly.
func (e *Engine) Execute(ctx context.Context, code string) error {
	callChain, err := e.parser.Parse(code)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	return e.interpret(ctx, callChain)
}

// Stream parses DSL code and executes methods as they're called (streaming).
// This allows for memory-efficient processing and real-time execution of large DSL programs.
//
// The channel will be closed when all methods have been executed or an error occurs.
// If an error occurs, it will be sent on the error channel before closing.
//
// Example:
//
//	errors := engine.Stream(context.Background(), `track(name="A").track(name="B")`)
//	for err := range errors {
//	  if err != nil {
//	    log.Fatal(err)
//	  }
//	}
func (e *Engine) Stream(ctx context.Context, code string) <-chan error {
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		callChain, err := e.parser.Parse(code)
		if err != nil {
			errors <- fmt.Errorf("parse error: %w", err)
			return
		}

		if err := e.interpretStream(ctx, callChain); err != nil {
			errors <- err
		}
	}()

	return errors
}

// interpret walks the CallChain and calls methods directly.
func (e *Engine) interpret(ctx context.Context, callChain *CallChain) error {
	for _, call := range callChain.Calls {
		handler, ok := e.methods[call.Name]
		if !ok {
			return fmt.Errorf("unknown method: %s", call.Name)
		}

		args := make(Args)
		for _, arg := range call.Args {
			args[arg.Name] = arg.Value
		}

		if err := handler(args); err != nil {
			return fmt.Errorf("method %s error: %w", call.Name, err)
		}
	}

	return nil
}

// interpretStream walks the CallChain and executes methods as they're called (streaming).
func (e *Engine) interpretStream(ctx context.Context, callChain *CallChain) error {
	for _, call := range callChain.Calls {
		handler, ok := e.methods[call.Name]
		if !ok {
			return fmt.Errorf("unknown method: %s", call.Name)
		}

		args := make(Args)
		for _, arg := range call.Args {
			args[arg.Name] = arg.Value
		}

		if err := handler(args); err != nil {
			return fmt.Errorf("method %s error: %w", call.Name, err)
		}
	}

	return nil
}

// CleanGrammarForCFG cleans a grammar string for use with CFG systems (e.g., GPT-5).
//
// Removes parser-specific directives that aren't supported in standard CFG:
// - Lines starting with % (Lark directives like %import, %ignore)
// - Empty lines for cleaner output
// - Other parser-specific meta-directives
//
// This is useful when exporting a grammar definition to use as a CFG constraint
// for LLM tools like GPT-5, which require standard CFG format without parser-specific directives.
//
// Args:
//   - grammar: Grammar string that may contain parser-specific directives
//
// Returns:
//   - Cleaned grammar string suitable for CFG systems
//
// Example:
//
//	cleaned := gs.CleanGrammarForCFG(grammarString)
//	// Use cleaned grammar with GPT-5 CFG
func CleanGrammarForCFG(grammar string) string {
	lines := strings.Split(grammar, "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Remove lines starting with % (Lark directives)
		// Remove empty lines for cleaner output
		if trimmed != "" && !strings.HasPrefix(trimmed, "%") {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}

const (
	// SyntaxLark is the default syntax for CFG grammars
	SyntaxLark = "lark"
	// SyntaxRegex is the regex syntax for CFG grammars
	SyntaxRegex = "regex"
	// TextFormatType is the text format type for OpenAI CFG requests
	TextFormatType = "text"
)

// CFGConfig contains configuration for building an OpenAI CFG tool.
type CFGConfig struct {
	ToolName    string // Name of the tool that will receive the DSL output
	Description string // Description of what the tool does
	Grammar     string // Lark or regex grammar definition
	Syntax      string // "lark" or "regex" (default: "lark")
}

// BuildOpenAICFGTool builds an OpenAI CFG tool payload from a CFGConfig.
//
// This function:
//   - Cleans the grammar using CleanGrammarForCFG
//   - Returns the properly formatted OpenAI tool structure
//   - Ensures the syntax defaults to "lark" if not specified
//
// Args:
//   - config: CFGConfig containing tool name, description, grammar, and syntax
//
// Returns:
//   - map[string]any: OpenAI tool structure ready to be added to the tools array
//
// Example:
//
//	tool := gs.BuildOpenAICFGTool(gs.CFGConfig{
//		ToolName:    "magda_dsl",
//		Description: "Generates MAGDA DSL code for REAPER automation",
//		Grammar:     grammarString,
//		Syntax:      "lark",
//	})
//	// Add tool to OpenAI request: tools = append(tools, tool)
func BuildOpenAICFGTool(config CFGConfig) map[string]any {
	// Clean the grammar for CFG
	cleanedGrammar := CleanGrammarForCFG(config.Grammar)

	// Default to "lark" if syntax is not specified
	syntax := config.Syntax
	if syntax == "" {
		syntax = SyntaxLark
	}

	// Build the OpenAI CFG tool structure
	return map[string]any{
		"type":        "custom",
		"name":        config.ToolName,
		"description": config.Description,
		"format": map[string]any{
			"type":       "grammar",
			"syntax":     syntax,
			"definition": cleanedGrammar,
		},
	}
}

// GetOpenAITextFormatForCFG returns the text format configuration that should be used
// when making OpenAI requests with CFG tools.
//
// When using CFG, the text format must be set to "text" (not JSON schema) because
// the output is DSL code, not JSON.
//
// Returns:
//   - map[string]any: Text format config: {"format": {"type": "text"}}
//
// Example:
//
//	paramsMap["text"] = gs.GetOpenAITextFormatForCFG()
func GetOpenAITextFormatForCFG() map[string]any {
	return map[string]any{
		"format": map[string]any{
			"type": TextFormatType,
		},
	}
}

// OpenAICFG is a convenient struct for building OpenAI CFG tool configurations.
//
// This struct encapsulates all the functionality needed to create OpenAI CFG tools
// from Grammar School grammars, eliminating the need to import and compose
// multiple utilities.
//
// Example:
//
//	// Use default Grammar School grammar (if available)
//	cfg := &gs.OpenAICFG{
//		ToolName:    "task_dsl",
//		Description: "Executes task management operations using Grammar School DSL.",
//		Grammar:     "", // Empty string will use default if set
//		Syntax:      gs.SyntaxLark,
//	}
//
//	// Build the tool and get text format
//	tool := cfg.BuildTool()
//	textFormat := cfg.GetTextFormat()
type OpenAICFG struct {
	ToolName    string // Name of the tool that will receive the DSL output
	Description string // Description of what the tool does
	Grammar     string // Lark or regex grammar definition (empty string uses default if available)
	Syntax      string // "lark" or "regex" (default: "lark")
}

// BuildTool builds the OpenAI CFG tool payload.
//
// Returns:
//   - map[string]any: OpenAI tool structure ready to be added to the tools array
//
// Example:
//
//	cfg := &gs.OpenAICFG{
//		ToolName:    "my_tool",
//		Description: "My tool",
//		Grammar:     grammarString,
//	}
//	tool := cfg.BuildTool()
//	// Use in OpenAI request: tools = append(tools, tool)
func (c *OpenAICFG) BuildTool() map[string]any {
	return BuildOpenAICFGTool(CFGConfig{
		ToolName:    c.ToolName,
		Description: c.Description,
		Grammar:     c.Grammar,
		Syntax:      c.Syntax,
	})
}

// GetTextFormat returns the text format configuration for OpenAI requests with CFG.
//
// Returns:
//   - map[string]any: Text format config: {"format": {"type": "text"}}
//
// Example:
//
//	cfg := &gs.OpenAICFG{
//		ToolName:    "my_tool",
//		Description: "My tool",
//	}
//	textFormat := cfg.GetTextFormat()
//	// Use in OpenAI request: text = textFormat
func (c *OpenAICFG) GetTextFormat() map[string]any {
	return GetOpenAITextFormatForCFG()
}

// BuildRequestConfig builds a complete request configuration with both tool and text format.
//
// This is a convenience method that returns both the tool and text format
// in a single map structure that can be easily merged into OpenAI request params.
//
// Returns:
//   - map[string]any: Map with "tool" and "text" keys ready for OpenAI request
//
// Example:
//
//	cfg := &gs.OpenAICFG{
//		ToolName:    "my_tool",
//		Description: "My tool",
//		Grammar:     grammarString,
//	}
//	config := cfg.BuildRequestConfig()
//	// Use in OpenAI request:
//	// tools = append(tools, config["tool"].(map[string]any))
//	// text = config["text"].(map[string]any)
func (c *OpenAICFG) BuildRequestConfig() map[string]any {
	return map[string]any{
		"tool": c.BuildTool(),
		"text": c.GetTextFormat(),
	}
}
