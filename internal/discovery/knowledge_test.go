package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

func knowledgeDocument(t *testing.T, fields map[string]any) string {
	t.Helper()
	data, err := yaml.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return "---\n" + string(data) + "---\n\n# Knowledge\n"
}

func TestOpenSeparatesSoftKnowledgeFromGovernedRecords(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	for _, dir := range []string{"plans", "decisions", "sketch", "scratchpad"} {
		if err := os.MkdirAll(filepath.Join(meta, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, filepath.Join(meta, "workspace.yaml"), defaultTestManifest)
	write(t, filepath.Join(meta, "PROJECT.md"), knowledgeDocument(t, map[string]any{
		"type": "Anchor", "id": "0198a1a0-0000-7000-8000-000000000003",
	}))
	write(t, filepath.Join(meta, "plans", "unfinished.md"), "unfinished raw draft")
	write(t, filepath.Join(meta, "sketch", "free.md"), "no obligations")
	write(t, filepath.Join(meta, "scratchpad", "free.md"), "no obligations")
	write(t, filepath.Join(meta, "INDEX.md"), "manual navigation")
	write(t, filepath.Join(meta, "ONTOLOGY.md"), knowledgeDocument(t, map[string]any{
		"type": "Anchor", "version": "0.1",
	}))
	write(t, filepath.Join(meta, "decisions", "D2-context.md"), knowledgeDocument(t, map[string]any{
		"type": "Decision", "governance": "context", "version": "0.1",
	}))
	write(t, filepath.Join(meta, "knowledge.md"), knowledgeDocument(t, map[string]any{
		"type": "Custom Knowledge", "version": "0.1",
	}))
	write(t, filepath.Join(meta, "decisions", "storage-choice.md"), knowledgeDocument(t, map[string]any{
		"type": "Decision", "version": "0.1", "created": "2026-10-07T17:50:32Z",
	}))
	write(t, filepath.Join(meta, "decisions", "D1-storage.md"), knowledgeDocument(t, map[string]any{
		"type": "Decision", "id": "0198a1a0-0000-7000-8000-000000000004",
	}))
	ws, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ws.Entries) != 2 {
		t.Fatalf("governed entries=%d; want anchor and governed Decision", len(ws.Entries))
	}
	if _, err := ws.Lookup(".spectacular/decisions/storage-choice.md", domain.Decision); err == nil {
		t.Fatal("soft Decision resolved as a governed record")
	}
}

const defaultTestManifest = "schema_version: spectacular.workspace.v1\nrecord_roots: [.]\nproject_anchor: PROJECT.md\n"

func TestOpenDoesNotDowngradeDamagedGovernedRecords(t *testing.T) {
	for _, tc := range []struct {
		name   string
		path   string
		fields map[string]any
	}{
		{"canonical name missing id", "decisions/D1-choice.md", map[string]any{"type": "Decision"}},
		{"invalid id", "decisions/choice.md", map[string]any{"type": "Decision", "id": "broken"}},
		{"unknown type with identity", "note.md", map[string]any{"type": "Custom", "id": "broken"}},
		{"governed ref missing id", "decisions/choice.md", map[string]any{"type": "Decision", "ref": "D1"}},
		{"schema claim", "note.md", map[string]any{"type": "Custom", "schema": "spectacular.unknown"}},
		{"context cannot override identity", "decisions/D2-choice.md", map[string]any{"type": "Decision", "governance": "context", "id": "broken"}},
		{"strict collection", "missions/choice.md", map[string]any{"type": "Custom"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			meta := filepath.Join(root, ".spectacular")
			if err := os.MkdirAll(filepath.Dir(filepath.Join(meta, tc.path)), 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, filepath.Join(meta, "workspace.yaml"), defaultTestManifest)
			write(t, filepath.Join(meta, "PROJECT.md"), knowledgeDocument(t, map[string]any{
				"type": "Anchor", "id": "0198a1a0-0000-7000-8000-000000000003",
			}))
			write(t, filepath.Join(meta, tc.path), knowledgeDocument(t, tc.fields))
			if _, err := Open(root); err == nil {
				t.Fatal("damaged governed record silently ignored")
			}
		})
	}
}
