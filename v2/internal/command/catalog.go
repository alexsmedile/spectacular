package command

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type CatalogEntry struct {
	Command   string `json:"command"`
	Arguments string `json:"arguments"`
	Schema    string `json:"schema"`
	Effect    Effect `json:"effect"`
}

func Catalog() []CatalogEntry {
	out := make([]CatalogEntry, 0, len(Registry))
	for _, spec := range Registry {
		out = append(out, CatalogEntry{Command: "spectacular " + strings.Join(spec.Words, " "), Arguments: spec.Arguments, Schema: spec.JSONSchema, Effect: spec.Effect})
	}
	return out
}

func WriteCatalogJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		SchemaVersion string         `json:"schema_version"`
		Commands      []CatalogEntry `json:"commands"`
	}{SchemaVersion: "spectacular.command-catalog.v1", Commands: Catalog()})
}

func WriteCatalogMarkdown(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "# Mechanical interface\n\nGenerated from `internal/command.Registry`; do not edit by hand.\n\n| Command | Arguments | Schema | Effect |\n|---|---|---|---|"); err != nil {
		return err
	}
	for _, entry := range Catalog() {
		if _, err := fmt.Fprintf(w, "| `%s` | `%s` | `%s` | `%s` |\n", entry.Command, entry.Arguments, entry.Schema, entry.Effect); err != nil {
			return err
		}
	}
	return nil
}
