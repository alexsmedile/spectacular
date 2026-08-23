package missionbundle

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestRecordEvidence_AtomicAndIdempotent(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)

	// Create workspace and a valid active mission
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	s := Service{Workspace: ws}

	planContent := `---
type: MissionPlan
title: Test Mission For Evidence
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: Test atomic clustered evidence.
review: independent
completion:
  - claim: test-claim
    pass_boundary: pass
    proof_requirement: proof
objectives:
  - outcome: First objective
    claims: [test-claim]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [internal/]
  semantic: [testing]
repair_budget: 1
dependencies: []
gaps: []
stops: [data-loss]
---
# Plan
`
	planFile := filepath.Join(root, "plan.md")
	if err := os.WriteFile(planFile, []byte(planContent), 0o644); err != nil {
		t.Fatal(err)
	}
	plan, raw, err := ReadPlan(planFile, nil)
	if err != nil {
		t.Fatalf("failed to read plan: %v", err)
	}
	res, err := s.Start(plan, raw)
	if err != nil {
		t.Fatalf("failed to start mission: %v", err)
	}

	commit, tree := currentGitCoordinates(t, root)

	evidenceDraft := `---
type: EvidenceDraft
title: Attributable test proof for O1
actor: Alex
commit: ` + commit + `
tree: ` + tree + `
objectives: [O1]
claims: [test-claim]
checks:
  - name: go-test
    result: pass
limitations: []
---
# Evidence Details
Test suite execution succeeded with 0 errors.
`
	evResult, err := s.RecordEvidence(res.Ref, "-", []byte(evidenceDraft))
	if err != nil {
		t.Fatalf("failed to record evidence: %v", err)
	}

	if evResult.Operation != "evidence.record" || !strings.HasSuffix(evResult.Ref, "/E1") {
		t.Fatalf("unexpected evidence result: %+v", evResult)
	}

	// Idempotent retry convergence
	retryResult, err := s.RecordEvidence(res.Ref, "-", []byte(evidenceDraft))
	if err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	if retryResult.Ref != evResult.Ref || retryResult.Path != evResult.Path {
		t.Fatalf("expected idempotent retry to match, got %+v vs %+v", retryResult, evResult)
	}

	// Reload bundle and check evidence pointer
	ws, err = discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Load(ws, res.Ref)
	if err != nil {
		t.Fatalf("failed to load bundle: %v", err)
	}
	if len(b.Evidence) != 1 {
		t.Fatalf("expected 1 evidence record, got %d", len(b.Evidence))
	}
	if b.Evidence[0].Document == nil || b.Evidence[0].Document.Actor != "Alex" {
		t.Fatalf("evidence document not resolved or actor mismatch: %+v", b.Evidence[0])
	}
}

func TestRecordEvidence_ValidationFailures(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)
	ws := &discovery.Workspace{Root: root, Config: discovery.DefaultConfig()}
	s := Service{Workspace: ws}

	// Missing actor
	invalidDraft := `---
type: EvidenceDraft
title: Missing actor
commit: 0123456789abcdef0123456789abcdef01234567
tree: 0123456789abcdef0123456789abcdef01234567
---
# Body
`
	_, err := s.RecordEvidence("M1", "-", []byte(invalidDraft))
	if err == nil {
		t.Fatal("expected refusal for missing actor/mission, got none")
	}
}

func initGitRepo(t *testing.T, root string) {
	t.Helper()
	runCmd(t, root, "git", "init")
	runCmd(t, root, "git", "config", "user.name", "Test User")
	runCmd(t, root, "git", "config", "user.email", "test@example.com")
	f := filepath.Join(root, "README.md")
	_ = os.WriteFile(f, []byte("# Test\n"), 0o644)
	specDir := filepath.Join(root, ".spectacular")
	_ = os.MkdirAll(specDir, 0o755)
	_ = os.WriteFile(filepath.Join(specDir, "workspace.yaml"), []byte("schema_version: spectacular.workspace.v1\nrecord_roots: [.]\nproject_anchor: PROJECT.md\n"), 0o644)
	_ = os.WriteFile(filepath.Join(specDir, "PROJECT.md"), []byte("---\ntype: Anchor\nid: 0198a1a0-0000-7000-8000-000000000003\n---\n"), 0o644)
	contractsDir := filepath.Join(specDir, "contracts")
	_ = os.MkdirAll(contractsDir, 0o755)
	const contract = `---
type: Contract
id: 019fe381-5d61-7223-b362-03a5f99a7b10
title: Test Contract
owner: Alex
contract_version: "1"
---
# Contract
`
	_ = os.WriteFile(filepath.Join(contractsDir, "contract.md"), []byte(contract), 0o644)
	runCmd(t, root, "git", "add", ".")
	runCmd(t, root, "git", "commit", "-m", "initial commit")
}

func currentGitCoordinates(t *testing.T, root string) (string, string) {
	t.Helper()
	commit := strings.TrimSpace(runCmd(t, root, "git", "rev-parse", "HEAD"))
	tree := strings.TrimSpace(runCmd(t, root, "git", "rev-parse", "HEAD^{tree}"))
	return commit, tree
}

func runCmd(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %s %v failed in %s: %v\nOutput: %s", name, args, dir, err, string(out))
	}
	return string(out)
}
