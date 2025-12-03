package gs

// ValueKind represents the type of a value.
type ValueKind int

const (
	ValueNumber ValueKind = iota
	ValueString
	ValueIdentifier
	ValueBool
	ValueFunction // Function reference (@function_name)
)

// String returns the string representation of ValueKind.
func (v ValueKind) String() string {
	switch v {
	case ValueNumber:
		return "number"
	case ValueString:
		return "string"
	case ValueIdentifier:
		return "identifier"
	case ValueBool:
		return "bool"
	case ValueFunction:
		return "function"
	default:
		return "unknown"
	}
}

// Value represents a value in the AST (number, string, identifier, etc.).
type Value struct {
	Kind ValueKind
	Num  float64
	Str  string
	Bool bool
}

// Arg represents a named argument to a call.
type Arg struct {
	Name  string
	Value Value
}

// Call represents a single function call with named arguments.
type Call struct {
	Name string
	Args []Arg
}

// CallChain represents a chain of calls connected by dots (method chaining).
type CallChain struct {
	Calls []Call
}
