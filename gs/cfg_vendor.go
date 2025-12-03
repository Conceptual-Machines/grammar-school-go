package gs

import (
	"context"
)

// CFGProvider is the interface that LLM providers must implement
// to support Context-Free Grammar (CFG) constraints with Grammar School.
//
// This interface allows Grammar School to work with different LLM providers
// (OpenAI, Anthropic, Google, etc.) that implement CFG in their own way.
type CFGProvider interface {
	// BuildTool builds the CFG tool payload for this provider.
	BuildTool(toolName, description, grammar, syntax string) map[string]any

	// GetTextFormat returns the text format configuration for this provider.
	GetTextFormat() map[string]any

	// Generate generates a response from the provider's API.
	Generate(
		ctx context.Context,
		prompt, model string,
		tools []map[string]any,
		textFormat map[string]any,
		client interface{},
		kwargs map[string]any,
	) (interface{}, error)

	// ExtractDSLCode extracts DSL code from the provider's response.
	ExtractDSLCode(response interface{}) (string, error)
}

// OpenAICFGProvider is the OpenAI implementation of the CFG provider interface.
type OpenAICFGProvider struct{}

// Backward compatibility: CFGVendor is an alias for CFGProvider
type CFGVendor = CFGProvider

// Backward compatibility: OpenAICFGVendor is an alias for OpenAICFGProvider
type OpenAICFGVendor = OpenAICFGProvider

// BuildTool builds OpenAI CFG tool payload.
func (v *OpenAICFGProvider) BuildTool(toolName, description, grammar, syntax string) map[string]any {
	return BuildOpenAICFGTool(CFGConfig{
		ToolName:    toolName,
		Description: description,
		Grammar:     grammar,
		Syntax:      syntax,
	})
}

// GetTextFormat returns OpenAI text format configuration.
func (v *OpenAICFGProvider) GetTextFormat() map[string]any {
	return GetOpenAITextFormatForCFG()
}

// Generate generates response from OpenAI API.
// Note: This is a placeholder - actual OpenAI client integration would go here.
// For now, this demonstrates the interface structure.
func (v *OpenAICFGProvider) Generate(
	ctx context.Context,
	prompt, model string,
	tools []map[string]any,
	textFormat map[string]any,
	client interface{},
	kwargs map[string]any,
) (interface{}, error) {
	// This would call the OpenAI SDK
	// For now, return nil to indicate it needs to be implemented
	// or use the OpenAI SDK directly
	return nil, nil
}

// ExtractDSLCode extracts DSL code from OpenAI response.
// Note: This is a placeholder - actual response parsing would go here.
func (v *OpenAICFGProvider) ExtractDSLCode(response interface{}) (string, error) {
	// This would parse the OpenAI response structure
	// For now, return empty string to indicate it needs to be implemented
	return "", nil
}
