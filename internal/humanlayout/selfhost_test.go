package humanlayout

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestSelfHostedIndexesAreRebuildableCollectionCaches(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := Indexes(opened.Entries, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for path, expected := range generated {
		target := filepath.Join(root, filepath.FromSlash(path))
		if os.Getenv("WRITE_INDEXES") == "1" {
			_ = os.MkdirAll(filepath.Dir(target), 0755)
			if err := os.WriteFile(target, expected, 0644); err != nil {
				t.Fatalf("write generated index %s: %v", path, err)
			}
		}
		actual, readErr := os.ReadFile(target)
		if readErr != nil {
			t.Fatalf("read generated index %s: %v", path, readErr)
		}
		if !bytes.Equal(actual, expected) {
			t.Fatalf("generated collection cache drifted: %s", path)
		}
	}
	for _, mission := range []string{"M5-implement-compact-expandable-missions", "M6-implement-compact-mission-cli"} {
		if _, err := os.Stat(filepath.Join(root, ".spectacular", "missions", mission, "index.md")); !os.IsNotExist(err) {
			t.Fatalf("compact Mission has a local index: %s", mission)
		}
	}
}
