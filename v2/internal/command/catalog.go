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

var VersionInspection = CatalogEntry{
	Command: "spectacular --version", Arguments: "[--json]", Schema: "spectacular.build-info.v1", Effect: ReadOnly,
}

func Catalog() []CatalogEntry {
	out := make([]CatalogEntry, 0, len(Registry))
	for _, spec := range Registry {
		out = append(out, CatalogEntry{Command: "spectacular " + strings.Join(spec.Words, " "), Arguments: spec.Arguments, Schema: spec.JSONSchema, Effect: spec.Effect})
	}
	return out
}

func WriteCatalogJSONVersion(w io.Writer, version string) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		SchemaVersion     string         `json:"schema_version"`
		ReleaseVersion    string         `json:"release_version"`
		ReleaseInspection CatalogEntry   `json:"release_inspection"`
		Commands          []CatalogEntry `json:"commands"`
	}{SchemaVersion: "spectacular.command-catalog.v1", ReleaseVersion: version, ReleaseInspection: VersionInspection, Commands: Catalog()})
}

func WriteCatalogMarkdownVersion(w io.Writer, version string) error {
	if _, err := fmt.Fprintf(w, "# Mechanical interface\n\nGenerated from `internal/command.Registry`; do not edit by hand.\n\nRelease version: `%s`\n\nRelease inspection: `%s %s` (`%s`, `%s`)\n\n| Command | Arguments | Schema | Effect |\n|---|---|---|---|\n", version, VersionInspection.Command, VersionInspection.Arguments, VersionInspection.Schema, VersionInspection.Effect); err != nil {
		return err
	}
	for _, entry := range Catalog() {
		if _, err := fmt.Fprintf(w, "| `%s` | `%s` | `%s` | `%s` |\n", entry.Command, entry.Arguments, entry.Schema, entry.Effect); err != nil {
			return err
		}
	}
	return nil
}
