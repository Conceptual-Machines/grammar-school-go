package gs

import (
	"strconv"
	"strings"
)

// LarkParser is a generic parser that parses method chaining DSL code
// into a CallChain. It works with any grammar that follows the pattern:
// method(param=value).method2(param2=value2)
type LarkParser struct{}

// NewLarkParser creates a new generic Lark parser
func NewLarkParser() *LarkParser {
	return &LarkParser{}
}

// Parse parses DSL code into a CallChain
func (p *LarkParser) Parse(input string) (*CallChain, error) {
	chain := &CallChain{Calls: []Call{}}

	parts := splitMethodCalls(input)

	for _, part := range parts {
		call := parseMethodCall(part)
		if call != nil {
			chain.Calls = append(chain.Calls, *call)
		}
	}

	return chain, nil
}

// splitMethodCalls splits statements into separate calls.
// Handles both method chaining (track(...).method(...)) and concatenated statements (track(...)track(...))
func splitMethodCalls(input string) []string {
	var parts []string
	var current strings.Builder
	depth := 0
	inString := false
	escapeNext := false

	// Statement starters that indicate a new statement (not a method chain)
	statementStarters := []string{"track(", "filter(", "map(", "for_each(", "undo("}

	for i, r := range input {
		char := string(r)

		if escapeNext {
			current.WriteRune(r)
			escapeNext = false
			continue
		}

		if char == "\\" {
			escapeNext = true
			current.WriteRune(r)
			continue
		}

		if char == "\"" {
			inString = !inString
			current.WriteRune(r)
		} else if char == "(" {
			depth++
			current.WriteRune(r)
		} else if char == ")" {
			depth--
			current.WriteRune(r)
			
			// Check if this closing paren ends a statement and next token is a statement starter
			if depth == 0 && !inString {
				// Look ahead to see if next token is a statement starter
				remaining := input[i+1:]
				remaining = strings.TrimSpace(remaining)
				
				// Check if remaining starts with a statement starter
				for _, starter := range statementStarters {
					if strings.HasPrefix(remaining, starter) {
						// This is a statement boundary - split here
						if current.Len() > 0 {
							parts = append(parts, strings.TrimSpace(current.String()))
							current.Reset()
						}
						break
					}
				}
			}
		} else if char == "." && depth == 0 && !inString {
			// Method chaining separator
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// parseMethodCall parses "method(param=value)" into a Call
func parseMethodCall(input string) *Call {
	input = strings.TrimSpace(input)

	parenIndex := strings.Index(input, "(")
	if parenIndex == -1 {
		methodName := strings.TrimSpace(input)
		methodName = capitalizeMethodName(methodName)
		return &Call{
			Name: methodName,
			Args: []Arg{},
		}
	}

	methodName := strings.TrimSpace(input[:parenIndex])
	methodName = capitalizeMethodName(methodName)
	paramsStr := strings.TrimSpace(input[parenIndex+1:])
	paramsStr = strings.TrimSuffix(paramsStr, ")")
	args := parseArgs(paramsStr)

	return &Call{
		Name: methodName,
		Args: args,
	}
}

// parseArgs parses "param1=value1, param2=value2" into []Arg
func parseArgs(paramsStr string) []Arg {
	if paramsStr == "" {
		return []Arg{}
	}

	var args []Arg
	var current strings.Builder
	depth := 0
	inString := false

	for _, r := range paramsStr {
		char := string(r)

		if char == "\"" {
			inString = !inString
			current.WriteRune(r)
		} else if char == "(" {
			depth++
			current.WriteRune(r)
		} else if char == ")" {
			depth--
			current.WriteRune(r)
		} else if char == "," && depth == 0 && !inString {
			argStr := strings.TrimSpace(current.String())
			if argStr != "" {
				args = append(args, parseArg(argStr))
			}
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}

	argStr := strings.TrimSpace(current.String())
	if argStr != "" {
		args = append(args, parseArg(argStr))
	}

	return args
}

// parseArg parses "name=value" into Arg
func parseArg(argStr string) Arg {
	parts := strings.SplitN(argStr, "=", 2)
	if len(parts) != 2 {
		return Arg{
			Name:  "",
			Value: Value{Kind: ValueString, Str: argStr},
		}
	}

	name := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])
	value := parseValue(valueStr)

	return Arg{
		Name:  name,
		Value: value,
	}
}

// capitalizeMethodName converts snake_case to CamelCase (track -> Track, set_selected -> SetSelected)
func capitalizeMethodName(name string) string {
	if name == "" {
		return name
	}

	parts := strings.Split(name, "_")
	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
		}
	}

	return result.String()
}

// parseValue parses a value string into Value
func parseValue(valueStr string) Value {
	valueStr = strings.TrimSpace(valueStr)

	if strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"") {
		return Value{
			Kind: ValueString,
			Str:  valueStr[1 : len(valueStr)-1],
		}
	}

	if valueStr == "true" {
		return Value{Kind: ValueBool, Bool: true}
	}
	if valueStr == "false" {
		return Value{Kind: ValueBool, Bool: false}
	}

	if num, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return Value{Kind: ValueNumber, Num: num}
	}

	return Value{Kind: ValueString, Str: valueStr}
}

