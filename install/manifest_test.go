package install

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type manifest struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Skills      string         `json:"skills"`
	Hooks       any            `json:"hooks"`
	Interface   map[string]any `json:"interface"`
}

func TestRuntimeManifestsPackageOnlyTheCanonicalV2Skill(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	versionBytes, err := os.ReadFile(filepath.Join(root, "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	version := string(bytes.TrimSpace(versionBytes))
	for _, path := range []string{".codex-plugin/plugin.json", ".claude-plugin/plugin.json"} {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		var got manifest
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if got.Name != "spectacular" || got.Version != version || got.Description == "" || got.Hooks != nil {
			t.Fatalf("unsafe or incomplete %s: %#v", path, got)
		}
		if path == ".codex-plugin/plugin.json" && (got.Skills != "./skills/" || got.Interface["displayName"] != "Spectacular") {
			t.Fatalf("Codex Skill discovery metadata is invalid: %#v", got)
		}
	}
	for _, path := range []string{"skills/spectacular/SKILL.md", "skills/spectacular/generated/mechanical-interface.json"} {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Fatalf("packaged v2 surface missing %s: %v", path, err)
		}
	}
}
