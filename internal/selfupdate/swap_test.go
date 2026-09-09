package selfupdate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceBinarySwapsAndKeepsTheDisplacedBinary(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "spectacular")
	if err := os.WriteFile(path, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	backup, err := ReplaceBinary(path, []byte("new binary"))
	if err != nil {
		t.Fatalf("ReplaceBinary: %v", err)
	}

	installed, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(installed) != "new binary" {
		t.Errorf("installed binary is %q, want the new one", installed)
	}
	// The displaced binary is what makes a bad update undoable by hand.
	displaced, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("the displaced binary must survive: %v", err)
	}
	if string(displaced) != "old binary" {
		t.Errorf("displaced binary is %q, want the old one", displaced)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Error("the installed binary must be executable")
	}
}

func TestReplaceBinaryLeavesNoStagingFileBehind(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "spectacular")
	if err := os.WriteFile(path, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplaceBinary(path, []byte("new")); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	// Exactly the binary and its backup: a leftover staging file would
	// accumulate one dotfile per update in the user's bin directory.
	if len(entries) != 2 {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Errorf("directory holds %v, want only the binary and its backup", names)
	}
}

func TestReplaceBinaryRefusesEmptyData(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "spectacular")
	if err := os.WriteFile(path, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplaceBinary(path, nil); err == nil {
		t.Fatal("installing an empty binary must refuse")
	}
	// A refusal must not have touched the working install.
	surviving, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the existing binary must survive a refusal: %v", err)
	}
	if string(surviving) != "old binary" {
		t.Errorf("binary is %q, want it untouched", surviving)
	}
}

func TestReplaceBinaryReplacesAnOlderBackup(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "spectacular")
	if err := os.WriteFile(path, []byte("v1"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplaceBinary(path, []byte("v2")); err != nil {
		t.Fatal(err)
	}
	backup, err := ReplaceBinary(path, []byte("v3"))
	if err != nil {
		t.Fatal(err)
	}
	// The backup tracks the immediately previous version, not the first one:
	// a stale backup would roll a user back further than they expect.
	displaced, err := os.ReadFile(backup)
	if err != nil {
		t.Fatal(err)
	}
	if string(displaced) != "v2" {
		t.Errorf("backup holds %q, want the immediately previous binary", displaced)
	}
}
