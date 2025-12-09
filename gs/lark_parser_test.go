package gs

import (
	"reflect"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
	}{
		{
			name:  "simple string",
			input: `"hello"`,
			expected: Value{
				Kind: ValueString,
				Str:  "hello",
			},
		},
		{
			name:  "string with special characters",
			input: `"Sunlit Echoes 7f3a2"`,
			expected: Value{
				Kind: ValueString,
				Str:  "Sunlit Echoes 7f3a2",
			},
		},
		{
			name:  "string with parentheses",
			input: `"track(name="`,
			expected: Value{
				Kind: ValueString,
				Str:  "track(name=",
			},
		},
		{
			name:  "string with semicolon",
			input: `"test;value"`,
			expected: Value{
				Kind: ValueString,
				Str:  "test;value",
			},
		},
		{
			name:  "string with escaped quote",
			input: `"test\"value"`,
			expected: Value{
				Kind: ValueString,
				Str:  "test\"value",
			},
		},
		{
			name:  "boolean true",
			input: "true",
			expected: Value{
				Kind: ValueBool,
				Bool: true,
			},
		},
		{
			name:  "boolean false",
			input: "false",
			expected: Value{
				Kind: ValueBool,
				Bool: false,
			},
		},
		{
			name:  "number",
			input: "42",
			expected: Value{
				Kind: ValueNumber,
				Num:  42,
			},
		},
		{
			name:  "decimal number",
			input: "3.14",
			expected: Value{
				Kind: ValueNumber,
				Num:  3.14,
			},
		},
		{
			name:  "unquoted string",
			input: "identifier",
			expected: Value{
				Kind: ValueString,
				Str:  "identifier",
			},
		},
		{
			name:  "malformed string - no closing quote",
			input: `"unclosed string`,
			expected: Value{
				Kind: ValueString,
				Str:  `"unclosed string`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValue(tt.input)
			if result.Kind != tt.expected.Kind {
				t.Errorf("expected kind %v, got %v", tt.expected.Kind, result.Kind)
			}
			if result.Str != tt.expected.Str {
				t.Errorf("expected str %q, got %q", tt.expected.Str, result.Str)
			}
			if result.Num != tt.expected.Num {
				t.Errorf("expected num %v, got %v", tt.expected.Num, result.Num)
			}
			if result.Bool != tt.expected.Bool {
				t.Errorf("expected bool %v, got %v", tt.expected.Bool, result.Bool)
			}
		})
	}
}

func TestLarkParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *CallChain
		wantErr  bool
	}{
		{
			name:  "single track with name",
			input: `track(name="Test Track")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Test Track"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "track with special characters in name",
			input: `track(name="Sunlit Echoes 7f3a2")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Sunlit Echoes 7f3a2"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "multiple tracks with names",
			input: `track(name="Track 1");track(name="Track 2")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Track 1"},
							},
						},
					},
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Track 2"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "track with name containing problematic characters",
			input: `track(name="Sunlit Echoes 7f3a2);track(name=")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Sunlit Echoes 7f3a2);track(name="},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "method chaining",
			input: `track(name="Test").set_name(name="Renamed")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Test"},
							},
						},
					},
					{
						Name: "SetName",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Renamed"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "track with multiple parameters",
			input: `track(name="Test", instrument="Piano")`,
			expected: &CallChain{
				Calls: []Call{
					{
						Name: "Track",
						Args: []Arg{
							{
								Name:  "name",
								Value: Value{Kind: ValueString, Str: "Test"},
							},
							{
								Name:  "instrument",
								Value: Value{Kind: ValueString, Str: "Piano"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLarkParser()
			result, err := parser.Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result.Calls) != len(tt.expected.Calls) {
					t.Errorf("expected %d calls, got %d", len(tt.expected.Calls), len(result.Calls))
					return
				}

				for i, expectedCall := range tt.expected.Calls {
					actualCall := result.Calls[i]
					if actualCall.Name != expectedCall.Name {
						t.Errorf("call %d: expected name %q, got %q", i, expectedCall.Name, actualCall.Name)
					}

					if len(actualCall.Args) != len(expectedCall.Args) {
						t.Errorf("call %d: expected %d args, got %d", i, len(expectedCall.Args), len(actualCall.Args))
						continue
					}

					for j, expectedArg := range expectedCall.Args {
						actualArg := actualCall.Args[j]
						if actualArg.Name != expectedArg.Name {
							t.Errorf("call %d, arg %d: expected name %q, got %q", i, j, expectedArg.Name, actualArg.Name)
						}
						if !reflect.DeepEqual(actualArg.Value, expectedArg.Value) {
							t.Errorf("call %d, arg %d: expected value %+v, got %+v", i, j, expectedArg.Value, actualArg.Value)
						}
					}
				}
			}
		})
	}
}

func TestSplitMethodCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single statement",
			input:    `track(name="Test")`,
			expected: []string{`track(name="Test")`},
		},
		{
			name:     "two statements with semicolon",
			input:    `track(name="Test 1");track(name="Test 2")`,
			expected: []string{`track(name="Test 1")`, `track(name="Test 2")`},
		},
		{
			name:     "two statements concatenated",
			input:    `track(name="Test 1")track(name="Test 2")`,
			expected: []string{`track(name="Test 1")`, `track(name="Test 2")`},
		},
		{
			name:     "method chaining",
			input:    `track(name="Test").set_name(name="Renamed")`,
			expected: []string{`track(name="Test")`, `set_name(name="Renamed")`},
		},
		{
			name:     "string with semicolon inside",
			input:    `track(name="Test;Value")`,
			expected: []string{`track(name="Test;Value")`},
		},
		{
			name:     "string with parentheses inside",
			input:    `track(name="Test(value)")`,
			expected: []string{`track(name="Test(value)")`},
		},
		{
			name:     "multiple concatenated tracks",
			input:    `track(name="A")track(name="B")track(name="C")`,
			expected: []string{`track(name="A")`, `track(name="B")`, `track(name="C")`},
		},
		{
			name:     "filter with complex predicate",
			input:    `filter(tracks, track.name=="Test")`,
			expected: []string{`filter(tracks, track.name=="Test")`},
		},
		{
			name:     "nested parentheses in string",
			input:    `track(name="Test(value)")`,
			expected: []string{`track(name="Test(value)")`},
		},
		{
			name:     "empty string",
			input:    `track(name="")`,
			expected: []string{`track(name="")`},
		},
		{
			name:     "string with escaped quotes",
			input:    `track(name="Test\"Value")`,
			expected: []string{`track(name="Test\"Value")`},
		},
		{
			name:     "complex chaining",
			input:    `track(name="A").set_volume(volume_db=-3.0).set_pan(pan=0.5)`,
			expected: []string{`track(name="A")`, `set_volume(volume_db=-3.0)`, `set_pan(pan=0.5)`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitMethodCalls(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d parts, got %d", len(tt.expected), len(result))
				t.Logf("expected: %v", tt.expected)
				t.Logf("got: %v", result)
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("part %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Arg
	}{
		{
			name:     "single arg",
			input:    `name="Test"`,
			expected: []Arg{{Name: "name", Value: Value{Kind: ValueString, Str: "Test"}}},
		},
		{
			name:  "multiple args",
			input: `name="Test", instrument="Piano"`,
			expected: []Arg{
				{Name: "name", Value: Value{Kind: ValueString, Str: "Test"}},
				{Name: "instrument", Value: Value{Kind: ValueString, Str: "Piano"}},
			},
		},
		{
			name:  "mixed types",
			input: `name="Test", volume_db=-3.0, selected=true`,
			expected: []Arg{
				{Name: "name", Value: Value{Kind: ValueString, Str: "Test"}},
				{Name: "volume_db", Value: Value{Kind: ValueNumber, Num: -3.0}},
				{Name: "selected", Value: Value{Kind: ValueBool, Bool: true}},
			},
		},
		{
			name:     "string with comma inside",
			input:    `name="Test, Value"`,
			expected: []Arg{{Name: "name", Value: Value{Kind: ValueString, Str: "Test, Value"}}},
		},
		{
			name:     "string with parentheses inside",
			input:    `name="Test(value)"`,
			expected: []Arg{{Name: "name", Value: Value{Kind: ValueString, Str: "Test(value)"}}},
		},
		{
			name:     "string with semicolon inside",
			input:    `name="Test;Value"`,
			expected: []Arg{{Name: "name", Value: Value{Kind: ValueString, Str: "Test;Value"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArgs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d args, got %d", len(tt.expected), len(result))
				t.Logf("expected: %+v", tt.expected)
				t.Logf("got: %+v", result)
				return
			}

			for i, expected := range tt.expected {
				if result[i].Name != expected.Name {
					t.Errorf("arg %d: expected name %q, got %q", i, expected.Name, result[i].Name)
				}
				if !reflect.DeepEqual(result[i].Value, expected.Value) {
					t.Errorf("arg %d: expected value %+v, got %+v", i, expected.Value, result[i].Value)
				}
			}
		})
	}
}

func TestParseArg(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Arg
	}{
		{
			name:     "string value",
			input:    `name="Test"`,
			expected: Arg{Name: "name", Value: Value{Kind: ValueString, Str: "Test"}},
		},
		{
			name:     "number value",
			input:    `volume_db=-3.0`,
			expected: Arg{Name: "volume_db", Value: Value{Kind: ValueNumber, Num: -3.0}},
		},
		{
			name:     "boolean true",
			input:    `selected=true`,
			expected: Arg{Name: "selected", Value: Value{Kind: ValueBool, Bool: true}},
		},
		{
			name:     "boolean false",
			input:    `selected=false`,
			expected: Arg{Name: "selected", Value: Value{Kind: ValueBool, Bool: false}},
		},
		{
			name:     "positional arg (no name)",
			input:    `"Test"`,
			expected: Arg{Name: "", Value: Value{Kind: ValueString, Str: "\"Test\""}},
		},
		{
			name:     "string with special chars",
			input:    `name="Sunlit Echoes 7f3a2"`,
			expected: Arg{Name: "name", Value: Value{Kind: ValueString, Str: "Sunlit Echoes 7f3a2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArg(tt.input)
			if result.Name != tt.expected.Name {
				t.Errorf("expected name %q, got %q", tt.expected.Name, result.Name)
			}
			if !reflect.DeepEqual(result.Value, tt.expected.Value) {
				t.Errorf("expected value %+v, got %+v", tt.expected.Value, result.Value)
			}
		})
	}
}

func TestCapitalizeMethodName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "track", "Track"},
		{"snake_case", "set_selected", "SetSelected"},
		{"multiple_underscores", "set_track_name", "SetTrackName"},
		{"already_capitalized", "Track", "Track"},
		{"empty", "", ""},
		{"single_char", "a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capitalizeMethodName(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLarkParser_RegressionCases(t *testing.T) {
	// Test cases that previously caused regressions
	tests := []struct {
		name        string
		input       string
		description string
		validate    func(t *testing.T, chain *CallChain)
	}{
		{
			name:        "track name with problematic characters",
			input:       `track(name="Sunlit Echoes 7f3a2")`,
			description: "Track name that previously caused parsing to include next statement",
			validate: func(t *testing.T, chain *CallChain) {
				if len(chain.Calls) != 1 {
					t.Fatalf("expected 1 call, got %d", len(chain.Calls))
				}
				call := chain.Calls[0]
				if call.Name != "Track" {
					t.Errorf("expected call name 'Track', got %q", call.Name)
				}
				if len(call.Args) != 1 {
					t.Fatalf("expected 1 arg, got %d", len(call.Args))
				}
				arg := call.Args[0]
				if arg.Name != "name" {
					t.Errorf("expected arg name 'name', got %q", arg.Name)
				}
				if arg.Value.Str != "Sunlit Echoes 7f3a2" {
					t.Errorf("expected name value 'Sunlit Echoes 7f3a2', got %q", arg.Value.Str)
				}
			},
		},
		{
			name:        "multiple tracks with random names",
			input:       `track(name="Track 1");track(name="Track 2");track(name="Track 3")`,
			description: "Multiple tracks that should parse correctly",
			validate: func(t *testing.T, chain *CallChain) {
				if len(chain.Calls) != 3 {
					t.Fatalf("expected 3 calls, got %d", len(chain.Calls))
				}
				for i, call := range chain.Calls {
					if call.Name != "Track" {
						t.Errorf("call %d: expected name 'Track', got %q", i, call.Name)
					}
					if len(call.Args) != 1 {
						t.Errorf("call %d: expected 1 arg, got %d", i, len(call.Args))
					}
					expectedName := "Track " + string(rune('1'+i))
					if call.Args[0].Value.Str != expectedName {
						t.Errorf("call %d: expected name %q, got %q", i, expectedName, call.Args[0].Value.Str)
					}
				}
			},
		},
		{
			name:        "track name with set_name call",
			input:       `track(name="Track 1").set_name(name="Sunlit Echoes 7f3a2")`,
			description: "Track creation followed by name change",
			validate: func(t *testing.T, chain *CallChain) {
				if len(chain.Calls) != 2 {
					t.Fatalf("expected 2 calls, got %d", len(chain.Calls))
				}
				if chain.Calls[0].Name != "Track" {
					t.Errorf("expected first call 'Track', got %q", chain.Calls[0].Name)
				}
				if chain.Calls[1].Name != "SetName" {
					t.Errorf("expected second call 'SetName', got %q", chain.Calls[1].Name)
				}
				if chain.Calls[1].Args[0].Value.Str != "Sunlit Echoes 7f3a2" {
					t.Errorf("expected renamed value 'Sunlit Echoes 7f3a2', got %q", chain.Calls[1].Args[0].Value.Str)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLarkParser()
			chain, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			tt.validate(t, chain)
		})
	}
}

