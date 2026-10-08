package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProductionBoundaries(t *testing.T) {
	if err := check("../.."); err != nil {
		t.Fatal(err)
	}
}

func TestRejectBoundaryViolations(t *testing.T) {
	for _, tc := range []struct{ name, policy, source, want string }{
		{"forbidden sibling", `{"module":"example.test/product","packages":{"a":{"internal":[]},"b":{"internal":[]}}}`, `package a; import _ "example.test/product/internal/b"`, "forbidden dependency"},
		{"pure package IO", `{"module":"example.test/product","packages":{"a":{"internal":[],"external":[]},"b":{"internal":[]}}}`, `package a; import _ "os"`, "forbidden external dependency"},
		{"unclassified package", `{"module":"example.test/product","packages":{"b":{"internal":[]}}}`, `package a`, "has no architecture policy"},
		{"composition inversion", `{"module":"example.test/product","packages":{"a":{"internal":[]},"b":{"internal":[]}}}`, `package a; import _ "example.test/product/cmd/app"`, "must not import product composition"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for name, data := range map[string]string{"scripts/architecture-boundaries.json": tc.policy, "internal/a/a.go": tc.source, "internal/b/b_windows.go": "package b"} {
				path := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(data), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if err := check(root); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want %s", err, tc.want)
			}
		})
	}
}

func TestCyclesAcrossPlatformFiles(t *testing.T) {
	if err := acyclic(map[string][]string{"a": {"b"}, "b": {"a"}}); err == nil {
		t.Fatal("cycle accepted")
	}
	if err := acyclic(map[string][]string{"a": {"b", "c"}, "b": {"c"}}); err != nil {
		t.Fatal(err)
	}
}
