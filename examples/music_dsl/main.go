package main

import (
	"fmt"
	"github.com/Conceptual-Machines/grammar-school-go/gs"
)

// MusicDSL is a simple music DSL for creating tracks and clips.
type MusicDSL struct{}

// Track creates a new track.
func (d *MusicDSL) Track(args gs.Args) error {
	name := args["name"].Str
	color := ""
	if c, ok := args["color"]; ok {
		color = c.Str
	}

	fmt.Printf("Creating track: name=%s, color=%s\n", name, color)
	return nil
}

// AddClip adds a clip to the current track.
func (d *MusicDSL) AddClip(args gs.Args) error {
	start := args["start"].Num
	length := args["length"].Num

	fmt.Printf("Adding clip: start=%.2f, length=%.2f\n", start, length)
	return nil
}

// Mute mutes the current track.
func (d *MusicDSL) Mute(args gs.Args) error {
	fmt.Println("Muting track")
	return nil
}

// MusicParser is a placeholder parser (would need actual implementation).
type MusicParser struct{}

func (p *MusicParser) Parse(input string) (*gs.CallChain, error) {
	// This is a placeholder - in a real implementation, you'd use
	// a parser library like participle or pigeon here.
	// For now, we'll return an error indicating a parser is needed.
	return nil, fmt.Errorf("parser not implemented - use a real parser backend")
}

func main() {
	dsl := &MusicDSL{}
	parser := &MusicParser{}

	// Create engine with new unified API (no runtime needed)
	_, err := gs.NewEngine("", dsl, parser)
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		return
	}

	// Note: This won't work until a real parser is implemented
	code := `track(name="Drums").add_clip(start=0, length=8)`
	fmt.Printf("Code: %s\n", code)
	fmt.Println("\nNote: Parser implementation needed to run this example")
	fmt.Println("See parser_backend.go for the Parser interface")

	// Once parser is implemented, you could do:
	// engine, _ := gs.NewEngine("", dsl, parser)
	// err := engine.Execute(context.Background(), code)
	// if err != nil {
	// 	fmt.Printf("Error executing: %v\n", err)
	// 	return
	// }
}
