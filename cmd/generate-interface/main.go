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
	version := flag.String("version", "", "release version written into generated catalogs")
	flag.Parse()
	if *jsonPath == "" || *markdownPath == "" || *version == "" {
		fmt.Fprintln(os.Stderr, "--json, --markdown, and --version are required")
		os.Exit(2)
	}
	if err := write(*jsonPath, func(w io.Writer) error { return command.WriteCatalogJSONVersion(w, *version) }); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := write(*markdownPath, func(w io.Writer) error { return command.WriteCatalogMarkdownVersion(w, *version) }); err != nil {
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
