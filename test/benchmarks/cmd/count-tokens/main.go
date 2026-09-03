package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: count-tokens <file-path|->")
		os.Exit(1)
	}

	target := os.Args[1]
	var data []byte
	var err error

	if target == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(target)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", target, err)
		os.Exit(1)
	}

	content := string(data)
	lines := len(strings.Split(content, "\n"))
	words := len(strings.Fields(content))
	chars := len(content)

	tokens, err := tokenizer.Count(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Tokenizer error: %v\n", err)
		os.Exit(1)
	}

	disp, status := tokenizer.EvaluateDisposition(tokens)

	fmt.Printf("Target: %s\n", filepath.ToSlash(target))
	fmt.Printf("  Lines:      %d\n", lines)
	fmt.Printf("  Words:      %d\n", words)
	fmt.Printf("  Characters: %d\n", chars)
	fmt.Printf("  Tokens:     %d (%s, %s)\n", tokens, disp, tokenizer.Version)
	fmt.Printf("  Status:     %s\n", status)
}
