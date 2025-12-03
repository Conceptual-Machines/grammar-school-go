package gs

import "testing"

func TestValue(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name:     "number value",
			value:    Value{Kind: ValueNumber, Num: 42},
			expected: "number",
		},
		{
			name:     "string value",
			value:    Value{Kind: ValueString, Str: "hello"},
			expected: "string",
		},
		{
			name:     "identifier value",
			value:    Value{Kind: ValueIdentifier, Str: "myVar"},
			expected: "identifier",
		},
		{
			name:     "bool value",
			value:    Value{Kind: ValueBool, Bool: true},
			expected: "bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.Kind.String() != tt.expected {
				t.Errorf("expected kind %s, got %s", tt.expected, tt.value.Kind.String())
			}
		})
	}
}

func TestArg(t *testing.T) {
	value := Value{Kind: ValueString, Str: "test"}
	arg := Arg{Name: "name", Value: value}

	if arg.Name != "name" {
		t.Errorf("expected name 'name', got %s", arg.Name)
	}
	if arg.Value.Str != "test" {
		t.Errorf("expected value 'test', got %s", arg.Value.Str)
	}
}

func TestCall(t *testing.T) {
	value1 := Value{Kind: ValueString, Str: "test"}
	value2 := Value{Kind: ValueNumber, Num: 42}

	call := Call{
		Name: "greet",
		Args: []Arg{
			{Name: "name", Value: value1},
			{Name: "count", Value: value2},
		},
	}

	if call.Name != "greet" {
		t.Errorf("expected name 'greet', got %s", call.Name)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

func TestCallChain(t *testing.T) {
	call1 := Call{Name: "track", Args: []Arg{}}
	call2 := Call{Name: "add_clip", Args: []Arg{}}

	chain := CallChain{
		Calls: []Call{call1, call2},
	}

	if len(chain.Calls) != 2 {
		t.Errorf("expected 2 calls, got %d", len(chain.Calls))
	}
	if chain.Calls[0].Name != "track" {
		t.Errorf("expected first call name 'track', got %s", chain.Calls[0].Name)
	}
	if chain.Calls[1].Name != "add_clip" {
		t.Errorf("expected second call name 'add_clip', got %s", chain.Calls[1].Name)
	}
}
