package gs

// Action represents a runtime action produced by the interpreter.
type Action struct {
	Kind    string                 `json:"kind"`
	Payload map[string]interface{} `json:"payload"`
}
