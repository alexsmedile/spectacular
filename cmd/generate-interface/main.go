package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/command"
)

func main() {
	if err := run(os.Args[1:], "."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, root string) error {
	flags := flag.NewFlagSet("generate-interface", flag.ContinueOnError)
	jsonPath := flags.String("json", filepath.Join(root, "skills", "spectacular", "generated", "mechanical-interface.json"), "output path for the JSON command catalog")
	markdownPath := flags.String("markdown", filepath.Join(root, "skills", "spectacular", "generated", "mechanical-interface.md"), "output path for the Markdown command catalog")
	version := flags.String("version", "", "release version written into generated catalogs")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	if *version == "" {
		data, err := os.ReadFile(filepath.Join(root, "VERSION"))
		if err != nil {
			return fmt.Errorf("read VERSION: %w", err)
		}
		*version = strings.TrimSpace(string(data))
	}
	if *version == "" {
		return fmt.Errorf("VERSION is empty")
	}
	if err := write(*jsonPath, func(w io.Writer) error { return command.WriteCatalogJSONVersion(w, *version) }); err != nil {
		return err
	}
	if err := write(*markdownPath, func(w io.Writer) error { return command.WriteCatalogMarkdownVersion(w, *version) }); err != nil {
		return err
	}
	return nil
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
