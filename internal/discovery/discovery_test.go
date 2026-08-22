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

func TestOpenResolvesHumanReferencesAndIgnoresGeneratedNavigation(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	missionDir := filepath.Join(meta, "missions", "M4-readable")
	if err := os.MkdirAll(missionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(meta, "campaigns"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(meta, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots: [.]\nproject_anchor: PROJECT.md\n")
	write(t, filepath.Join(meta, "PROJECT.md"), "---\ntype: Anchor\nid: 0198a1a0-0000-7000-8000-000000000003\nhuman_ref: PROJECT\n---\n")
	write(t, filepath.Join(missionDir, "M4-readable.md"), "---\ntype: Mission\nid: 0198a1a0-0000-7000-8000-000000000002\nhuman_ref: M4\n---\n")
	write(t, filepath.Join(meta, "index.md"), "# generated and intentionally not canonical\n")
	write(t, filepath.Join(meta, "catalog.md"), "# generated and intentionally not canonical\n")
	write(t, filepath.Join(meta, "campaigns", "launch.md"), "# Campaign: launch\n")
	opened, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	entry, err := opened.Lookup("M4", domain.Mission)
	if err != nil || entry.Path != ".spectacular/missions/M4-readable/M4-readable.md" {
		t.Fatalf("human lookup=%#v err=%v", entry, err)
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

func TestOpenRefusesAuthorityMarkerSymlinksBeforeReading(t *testing.T) {
	t.Run("workspace marker", func(t *testing.T) {
		root := t.TempDir()
		meta := filepath.Join(root, ".spectacular")
		outside := t.TempDir()
		if err := os.MkdirAll(meta, 0o755); err != nil {
			t.Fatal(err)
		}
		externalMarker := filepath.Join(outside, "workspace.yaml")
		write(t, externalMarker, "schema_version: spectacular.workspace.v1\nrecord_roots: [records]\nproject_anchor: records/project.md\n")
		if err := os.Symlink(externalMarker, filepath.Join(meta, "workspace.yaml")); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
			t.Fatalf("symlinked marker error=%v", err)
		}
	})

	t.Run("metadata directory", func(t *testing.T) {
		root := t.TempDir()
		outside := t.TempDir()
		write(t, filepath.Join(outside, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots: [records]\nproject_anchor: records/project.md\n")
		if err := os.Symlink(outside, filepath.Join(root, ".spectacular")); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
			t.Fatalf("symlinked metadata directory error=%v", err)
		}
	})
}

func TestOpenLoadsConfigYAML(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(meta, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots: [.]\nproject_anchor: PROJECT.md\n")
	write(t, filepath.Join(meta, "PROJECT.md"), "---\ntype: Anchor\nid: 0198a1a0-0000-7000-8000-000000000003\nhuman_ref: PROJECT\n---\n")

	// 1. Without config.yaml -> should have DefaultConfig
	ws1, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if ws1.Config.TokenBudgets.MissionActive.Target != 900 {
		t.Fatalf("expected default active mission target 900, got %d", ws1.Config.TokenBudgets.MissionActive.Target)
	}

	// 2. With config.yaml -> should overlay custom overrides
	configContent := `schema_version: "spectacular.config.v1"
token_budgets:
  mission_active:
    target: 850
    hard_cap: 1100
defaults:
  operator: "CustomOperator"
`
	write(t, filepath.Join(meta, "config.yaml"), configContent)
	ws2, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if ws2.Config.TokenBudgets.MissionActive.Target != 850 {
		t.Fatalf("expected custom target 850, got %d", ws2.Config.TokenBudgets.MissionActive.Target)
	}
	if ws2.Config.TokenBudgets.MissionActive.HardCap != 1100 {
		t.Fatalf("expected custom hard_cap 1100, got %d", ws2.Config.TokenBudgets.MissionActive.HardCap)
	}
	if ws2.Config.Defaults.Operator != "CustomOperator" {
		t.Fatalf("expected custom operator CustomOperator, got %s", ws2.Config.Defaults.Operator)
	}
	// Non-overridden values should retain defaults
	if ws2.Config.TokenBudgets.Charter.Target != 1200 {
		t.Fatalf("expected default charter target 1200, got %d", ws2.Config.TokenBudgets.Charter.Target)
	}
}

func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}
