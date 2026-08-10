package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/command"
)

func main() {
	jsonPath := flag.String("json", "", "output path for the JSON command catalog")
	markdownPath := flag.String("markdown", "", "output path for the Markdown command catalog")
	flag.Parse()
	if *jsonPath == "" || *markdownPath == "" {
		fmt.Fprintln(os.Stderr, "both --json and --markdown are required")
		os.Exit(2)
	}
	if err := write(*jsonPath, command.WriteCatalogJSON); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := write(*markdownPath, command.WriteCatalogMarkdown); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func write(path string, render func(io.Writer) error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temporary := path + ".tmp"
	file, err := os.OpenFile(temporary, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if err := render(file); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}
