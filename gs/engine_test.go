package gs

import (
	"reflect"
	"strings"
	"testing"
)

//nolint:funlen // Test function is long but readable with table-driven tests
func TestBuildOpenAICFGTool(t *testing.T) {
	tests := []struct {
		name     string
		config   CFGConfig
		validate func(t *testing.T, tool map[string]any)
	}{
		{
			name: "correct structure",
			config: CFGConfig{
				ToolName:    "magda_dsl",
				Description: "Generates MAGDA DSL code",
				Grammar:     "start: track",
				Syntax:      SyntaxLark,
			},
			validate: func(t *testing.T, tool map[string]any) {
				if tool["type"] != "custom" {
					t.Errorf("expected type 'custom', got %v", tool["type"])
				}
				if tool["name"] != "magda_dsl" {
					t.Errorf("expected name 'magda_dsl', got %v", tool["name"])
				}
				if tool["description"] != "Generates MAGDA DSL code" {
					t.Errorf("expected description 'Generates MAGDA DSL code', got %v", tool["description"])
				}

				format, ok := tool["format"].(map[string]any)
				if !ok {
					t.Fatalf("expected format to be map[string]any, got %T", tool["format"])
				}
				if format["type"] != "grammar" {
					t.Errorf("expected format.type 'grammar', got %v", format["type"])
				}
				if format["syntax"] != SyntaxLark {
					t.Errorf("expected format.syntax 'lark', got %v", format["syntax"])
				}
				if _, ok := format["definition"].(string); !ok {
					t.Errorf("expected format.definition to be string, got %T", format["definition"])
				}
			},
		},
		{
			name: "grammar cleaning applied",
			config: CFGConfig{
				ToolName:    "test_tool",
				Description: "Test tool",
				Grammar:     "%import common\nstart: track\ntrack: \"track\"\n",
				Syntax:      SyntaxLark,
			},
			validate: func(t *testing.T, tool map[string]any) {
				format := tool["format"].(map[string]any)
				definition := format["definition"].(string)

				// Verify %import directive was removed
				if strings.Contains(definition, "%import") {
					t.Errorf("expected %%import directive to be removed, but found in: %s", definition)
				}

				// Verify the actual grammar content is still there
				if !strings.Contains(definition, "start: track") {
					t.Errorf("expected 'start: track' in definition, got: %s", definition)
				}
				if !strings.Contains(definition, "track: \"track\"") {
					t.Errorf("expected 'track: \"track\"' in definition, got: %s", definition)
				}

				// Verify it matches what CleanGrammarForCFG would produce
				expectedCleaned := CleanGrammarForCFG("%import common\nstart: track\ntrack: \"track\"\n")
				if definition != expectedCleaned {
					t.Errorf("expected cleaned grammar to match CleanGrammarForCFG output, got: %s", definition)
				}
			},
		},
		{
			name: "default syntax handling",
			config: CFGConfig{
				ToolName:    "test_tool",
				Description: "Test tool",
				Grammar:     "start: test",
				Syntax:      "", // Empty string should default to "lark"
			},
			validate: func(t *testing.T, tool map[string]any) {
				format := tool["format"].(map[string]any)
				if format["syntax"] != SyntaxLark {
					t.Errorf("expected syntax to default to 'lark', got %v", format["syntax"])
				}
			},
		},
		{
			name: "regex syntax preserved",
			config: CFGConfig{
				ToolName:    "test_tool",
				Description: "Test tool",
				Grammar:     "^\\d+$",
				Syntax:      SyntaxRegex,
			},
			validate: func(t *testing.T, tool map[string]any) {
				format := tool["format"].(map[string]any)
				if format["syntax"] != SyntaxRegex {
					t.Errorf("expected syntax 'regex', got %v", format["syntax"])
				}
			},
		},
		{
			name: "all config fields used",
			config: CFGConfig{
				ToolName:    "custom_tool",
				Description: "Custom description with special chars: !@#$",
				Grammar:     "start: custom_rule\ncustom_rule: \"value\"",
				Syntax:      SyntaxLark,
			},
			validate: func(t *testing.T, tool map[string]any) {
				if tool["name"] != "custom_tool" {
					t.Errorf("expected name 'custom_tool', got %v", tool["name"])
				}
				if tool["description"] != "Custom description with special chars: !@#$" {
					t.Errorf("expected description with special chars, got %v", tool["description"])
				}
				format := tool["format"].(map[string]any)
				if format["syntax"] != SyntaxLark {
					t.Errorf("expected syntax 'lark', got %v", format["syntax"])
				}
				definition := format["definition"].(string)
				if !strings.Contains(definition, "custom_rule") {
					t.Errorf("expected 'custom_rule' in definition, got: %s", definition)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := BuildOpenAICFGTool(tt.config)
			tt.validate(t, tool)
		})
	}
}

func TestGetOpenAITextFormatForCFG(t *testing.T) {
	t.Run("text format structure", func(t *testing.T) {
		textFormat := GetOpenAITextFormatForCFG()

		// Verify it's a map
		if textFormat == nil {
			t.Fatal("expected non-nil map")
		}

		// Verify "format" key exists
		format, ok := textFormat["format"].(map[string]any)
		if !ok {
			t.Fatalf("expected format to be map[string]any, got %T", textFormat["format"])
		}

		// Verify "type" key exists and is "text"
		if format["type"] != TextFormatType {
			t.Errorf("expected format.type to be 'text', got %v", format["type"])
		}
	})

	t.Run("text format consistency", func(t *testing.T) {
		format1 := GetOpenAITextFormatForCFG()
		format2 := GetOpenAITextFormatForCFG()

		if !reflect.DeepEqual(format1, format2) {
			t.Errorf("expected consistent output, got different results: %v vs %v", format1, format2)
		}

		// Verify both have the correct structure
		if format1["format"].(map[string]any)["type"] != "text" {
			t.Errorf("expected format1.type to be 'text', got %v", format1["format"].(map[string]any)["type"])
		}
		if format2["format"].(map[string]any)["type"] != "text" {
			t.Errorf("expected format2.type to be 'text', got %v", format2["format"].(map[string]any)["type"])
		}
	})
}
