package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestOpenAscendsToNearestExplicitV2Marker(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	records := filepath.Join(meta, "records")
	if err := os.MkdirAll(records, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(meta, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots: [records]\nproject_anchor: records/project.md\n")
	write(t, filepath.Join(records, "project.md"), "---\ntype: Anchor\nid: 0198a1a0-0000-7000-8000-000000000003\nfreshness: current\nmissions: []\n---\n")
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := Open(nested)
	if err != nil {
		t.Fatal(err)
	}
	if got.Root != root {
		t.Fatalf("root=%q want %q", got.Root, root)
	}
}

func TestOpenRefusesUnsupportedMarkerAndEscapingRoot(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(meta, "workspace.yaml")
	write(t, marker, "schema_version: spectacular.workspace.v0\nrecord_roots: [records]\nproject_anchor: records/project.md\n")
	if _, err := Open(root); !domain.RefusalHasCode(err, domain.RefusalInvalidManifest) {
		t.Fatalf("unsupported marker error=%v", err)
	}
	write(t, marker, "schema_version: spectacular.workspace.v1\nrecord_roots: [../outside]\nproject_anchor: records/project.md\n")
	if _, err := Open(root); !domain.RefusalHasCode(err, domain.RefusalInvalidManifest) {
		t.Fatalf("escaping root error=%v", err)
	}
}

func TestOpenRefusesSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	records := filepath.Join(meta, "records")
	outside := t.TempDir()
	if err := os.MkdirAll(records, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(meta, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots: [records]\nproject_anchor: records/project.md\n")
	write(t, filepath.Join(records, "project.md"), "---\ntype: Anchor\nid: 0198a1a0-0000-7000-8000-000000000003\nfreshness: current\nmissions: []\n---\n")
	write(t, filepath.Join(outside, "escaped.md"), "not authoritative\n")
	if err := os.Symlink(filepath.Join(outside, "escaped.md"), filepath.Join(records, "escaped.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(root); !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
		t.Fatalf("symlink escape error=%v", err)
	}
}

func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}
