package missionbundle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestArchivedMissionLayoutRequiresTerminalStatusAndCanonicalPath(t *testing.T) {
	root := t.TempDir()
	path := ".spectacular/archive/missions/M1-delivered/M1-delivered.md"
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0755); err != nil {
		t.Fatal(err)
	}
	// This validator owns placement only; no record/frontmatter template is needed.
	if err := os.WriteFile(filepath.Join(root, path), nil, 0644); err != nil {
		t.Fatal(err)
	}
	ws := &discovery.Workspace{Root: root}
	for _, status := range []string{"completed", "resolved", "superseded", "withdrawn"} {
		if err := validateLayout(ws, &Bundle{Ref: "M1", Path: path, Status: status}); err != nil {
			t.Fatalf("terminal %s refused: %v", status, err)
		}
	}
	for _, status := range []string{"draft", "defined", "active", "invented"} {
		if err := validateLayout(ws, &Bundle{Ref: "M1", Path: path, Status: status}); err == nil {
			t.Fatalf("non-terminal %s accepted in archive", status)
		}
	}
	for _, invalidPath := range []string{
		".spectacular/archive/missions/M2-delivered/M1-delivered.md",
		".spectacular/archive/missions/M1-delivered/../M1-delivered/M1-delivered.md",
		".spectacular/archive/missions/M1-delivered/nested/M1-delivered.md",
	} {
		if err := validateLayout(ws, &Bundle{Ref: "M1", Path: invalidPath, Status: "completed"}); err == nil {
			t.Fatalf("non-canonical archive path accepted: %s", invalidPath)
		}
	}
}
