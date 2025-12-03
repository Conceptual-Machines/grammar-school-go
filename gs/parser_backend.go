package gs

// Parser is a pluggable interface for parsing DSL code into CallChain AST.
// This allows different parser backends (participle, pigeon, custom PEG/EBNF).
type Parser interface {
	Parse(input string) (*CallChain, error)
}
