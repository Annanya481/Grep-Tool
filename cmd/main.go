package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Annanya481/Grep-Tool/pkg/grep"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // Error code for incorrect usage
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // Read input from stdin, assume we only deal with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(2)
	}

	// Check if the line matches the pattern
	ok, err := grep.MatchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		log.Print("No match!")
		os.Exit(1) // No match
	}

	// Match success, exit with code 0
	log.Print("Match found!")
	os.Exit(0)
}
