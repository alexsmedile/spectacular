package guard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestGuardPassesOnCleanWrites(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Write within allowed path src/
	res, err := Run(ws, "M1/O1", false, "", []string{"sh", "-c", "echo 'package main' > src/app.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "pass" {
		t.Fatalf("expected pass, got %s (escapes: %v)", res.Status, res.EscapedPaths)
	}
	if len(res.PreservedPaths) == 0 {
		t.Fatalf("expected src/app.go to be recorded in preserved paths")
	}
}

func TestGuardSurgicalQuarantinePreservesValidWork(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Subagent writes BOTH valid code in src/ and rogue .gitignore
	res, err := Run(ws, "M1/O1", false, "", []string{"sh", "-c", "echo 'package main' > src/valid.go && echo '*.db' > .gitignore"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "violation" {
		t.Fatalf("expected violation, got %s", res.Status)
	}
	if !res.AutoHealed {
		t.Fatalf("expected auto_healed to be true")
	}

	// 1. Rogue file .gitignore must be purged
	if _, err := os.Stat(filepath.Join(ws.Root, ".gitignore")); !os.IsNotExist(err) {
		t.Fatalf("expected .gitignore to be quarantined and deleted")
	}

	// 2. ZERO WASTED WORK: Valid file src/valid.go must be preserved!
	if _, err := os.Stat(filepath.Join(ws.Root, "src/valid.go")); err != nil {
		t.Fatalf("expected valid file src/valid.go to be preserved, got error: %v", err)
	}

	// 3. Feedback prompt must mention preserved and quarantined paths
	if !strings.Contains(res.FeedbackPrompt, "src/valid.go") || !strings.Contains(res.FeedbackPrompt, ".gitignore") {
		t.Fatalf("feedback prompt missing key details:\n%s", res.FeedbackPrompt)
	}
}

func TestGuardRealtimeWatcherKillsRogueProcess(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Script that writes rogue file and sleeps briefly
	res, err := Run(ws, "M1/O1", true, "", []string{"sh", "-c", "echo 'evil' > evil.log && sleep 0.3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "killed" {
		t.Fatalf("expected killed status, got %s", res.Status)
	}
	if !res.RolledBack || !res.AutoHealed {
		t.Fatalf("expected rolled_back and auto_healed to be true")
	}

	// Verify evil.log was deleted
	if _, err := os.Stat(filepath.Join(ws.Root, "evil.log")); !os.IsNotExist(err) {
		t.Fatalf("expected evil.log to be rolled back and removed")
	}
}

func TestGuardPreExistingFileModificationIsRestored(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Pre-create README.md outside allowed perimeter
	readmePath := filepath.Join(ws.Root, "README.md")
	if err := os.WriteFile(readmePath, []byte("original text"), 0o644); err != nil {
		t.Fatalf("failed to write original README: %v", err)
	}

	// Subagent overwrites README.md with same-length string and writes valid code in src/
	res, err := Run(ws, "M1/O1", false, "", []string{"sh", "-c", "echo -n 'modified text' > README.md && echo 'package main' > src/app.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "violation" {
		t.Fatalf("expected violation on same-length edit outside perimeter, got %s", res.Status)
	}

	// Original README.md MUST be restored
	content, _ := os.ReadFile(readmePath)
	if string(content) != "original text" {
		t.Fatalf("expected README.md to be restored to 'original text', got '%s'", string(content))
	}

	// Valid work in src/app.go must be preserved
	if _, err := os.Stat(filepath.Join(ws.Root, "src/app.go")); err != nil {
		t.Fatalf("expected valid file src/app.go to be preserved")
	}
}

func TestGuardRealtimeWatcherCatchesInstantEscape(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Script that writes rogue file and exits immediately without waiting for ticker
	res, err := Run(ws, "M1/O1", true, "", []string{"sh", "-c", "echo evil > instant.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "violation" && res.Status != "killed" {
		t.Fatalf("expected violation or killed status, got %s", res.Status)
	}
	if _, err := os.Stat(filepath.Join(ws.Root, "instant.log")); !os.IsNotExist(err) {
		t.Fatalf("expected instant.log to be quarantined and deleted")
	}
}

func TestGuardExecShorthandPipesPrompt(t *testing.T) {
	ws, cleanup := setupTestWorkspace(t)
	defer cleanup()

	res, err := Run(ws, "M1/O1", false, "echo --flag", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "pass" {
		t.Fatalf("expected pass, got %s", res.Status)
	}
	if !strings.Contains(res.Output, "# Objective: M1/O1") {
		t.Fatalf("expected output to receive prompt, got:\n%s", res.Output)
	}
}

func setupTestWorkspace(t *testing.T) (*discovery.Workspace, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spectacular-guard-test-*")
	if err != nil {
		t.Fatal(err)
	}

	meta := filepath.Join(dir, ".spectacular")
	missions := filepath.Join(meta, "missions")
	contracts := filepath.Join(meta, "contracts")
	src := filepath.Join(dir, "src")

	_ = os.MkdirAll(missions, 0o755)
	_ = os.MkdirAll(contracts, 0o755)
	_ = os.MkdirAll(src, 0o755)

	contract := `---
id: 01955f1a-1234-7123-8123-123456789abc
type: Contract
schema: spectacular.contract.v2
ref: CC-test
title: Test Contract
contract_version: 1
---
`
	_ = os.WriteFile(filepath.Join(contracts, "CC-test.md"), []byte(contract), 0o644)

	mission := `---
id: 01955f1a-5678-7123-8123-123456789def
type: Mission
schema: spectacular.mission.v2
ref: M1
title: Test Mission
owner: Alex
status: active
contract:
  ref: CC-test
outcome: Test
review: none
completion:
  - claim: test-claim
    pass_boundary: tests pass
    proof_requirement: tests pass
objectives:
  - ref: O1
    outcome: Build app
    status: active
    claims: [test-claim]
scope:
  mechanical:
    - src/**
  semantic:
    - src/**
authority:
  operator: [read, write]
  requires_owner: [deploy]
repair_budget: 3
dependencies: []
gaps: []
stops: []
validation:
  schema: spectacular.mission.v2
  mode: strict
---
`
	wsYaml := `schema_version: spectacular.workspace.v1
record_roots:
  - .
project_anchor: PROJECT.md
`
	_ = os.WriteFile(filepath.Join(meta, "workspace.yaml"), []byte(wsYaml), 0o644)
	projectMd := `---
type: Anchor
id: 019fe381-5d61-7223-b362-03a5f99a7b13
title: Test Project
---
# Test Project
`
	_ = os.WriteFile(filepath.Join(meta, "PROJECT.md"), []byte(projectMd), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte(projectMd), 0o644)
	_ = os.WriteFile(filepath.Join(missions, "M1.md"), []byte(mission), 0o644)

	ws, err := discovery.Open(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}

	return ws, func() {
		_ = os.RemoveAll(dir)
	}
}
