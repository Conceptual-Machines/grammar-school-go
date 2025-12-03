package main

import (
	"fmt"
	"github.com/Conceptual-Machines/grammar-school-go/gs"
	"strings"
)

// FunctionalDSL demonstrates functional programming patterns.
// Users implement their own functional methods (map, filter, etc.) as needed.
type FunctionalDSL struct {
}

// Square squares a number.
func (d *FunctionalDSL) Square(args gs.Args) error {
	x := args["x"].Num
	result := x * x
	fmt.Printf("Square(%.2f) = %.2f\n", x, result)
	return nil
}

// Double doubles a number.
func (d *FunctionalDSL) Double(args gs.Args) error {
	x := args["x"].Num
	result := x * 2
	fmt.Printf("Double(%.2f) = %.2f\n", x, result)
	return nil
}

// IsEven checks if a number is even.
func (d *FunctionalDSL) IsEven(args gs.Args) error {
	x := args["x"].Num
	result := int(x)%2 == 0
	fmt.Printf("IsEven(%.2f) = %v\n", x, result)
	return nil
}

// FunctionalParser is a placeholder parser.
type FunctionalParser struct{}

func (p *FunctionalParser) Parse(input string) (*gs.CallChain, error) {
	return nil, fmt.Errorf("parser not implemented - use a real parser backend")
}

func main() {
	dsl := &FunctionalDSL{}
	parser := &FunctionalParser{}

	// Create engine with new unified API (no runtime needed)
	_, err := gs.NewEngine("", dsl, parser)
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		return
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Functional DSL Examples")
	fmt.Println(strings.Repeat("=", 60))

	// Note: These examples won't run without a real parser that supports
	// function references (@function_name syntax)
	fmt.Println("\nFunctional operations available:")
	fmt.Println("  - map(@Square, data)")
	fmt.Println("  - filter(@IsEven, data)")
	fmt.Println("  - reduce(@Add, data, 0)")
	fmt.Println("  - compose(@Square, @Double)")
	fmt.Println("  - pipe(data, @Double, @Square)")

	// Once parser is implemented, you could do:
	// engine, _ := gs.NewEngine("", dsl, parser)
	// err := engine.Execute(context.Background(), "map(@Square, data)")
	// if err != nil {
	// 	fmt.Printf("Error executing: %v\n", err)
	// 	return
	// }
}
