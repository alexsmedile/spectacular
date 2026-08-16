package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDefaultsToVersionAndCanonicalOutputs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "VERSION"), []byte("2.1.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(nil, root); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{
		"skills/spectacular/generated/mechanical-interface.json",
		"skills/spectacular/generated/mechanical-interface.md",
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "2.1.0") {
			t.Fatalf("%s does not use VERSION: %s", relative, data)
		}
	}
}
