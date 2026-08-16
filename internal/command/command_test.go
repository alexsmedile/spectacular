package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPublicRegistryIsMinimalAndTyped(t *testing.T) {
	want := []string{
		"mission start", "mission show", "mission check", "objective show", "objective promote",
		"objective finish", "run show", "run start", "review record", "mission complete",
	}
	if len(Registry) != len(want) {
		t.Fatalf("registry has %d commands, want %d", len(Registry), len(want))
	}
	for i, spec := range Registry {
		if got := strings.Join(spec.Words, " "); got != want[i] {
			t.Fatalf("registry[%d]=%q want %q", i, got, want[i])
		}
		if !strings.HasSuffix(spec.JSONSchema, ".v2") {
			t.Fatalf("registry command %q lacks v2 schema", want[i])
		}
	}
	for _, forbidden := range []string{"proposal create", "mission prepare", "mission transition", "contract reconcile", "workspace context"} {
		for _, spec := range Registry {
			if strings.Join(spec.Words, " ") == forbidden {
				t.Fatalf("superseded command remains public: %s", forbidden)
			}
		}
	}
}

func TestSelfHostedCompactMissionsUseSharedCheck(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"M5", "M6"} {
		output := run(t, root, nil, 0, "mission", "check", ref, "--json")
		if !strings.Contains(output, `"valid":true`) || !strings.Contains(output, `"schema":"mission.v2"`) || !strings.Contains(output, `"activation-fingerprint"`) {
			t.Fatalf("check %s returned %s", ref, output)
		}
	}
	legacy := run(t, root, nil, 0, "mission", "show", "M2", "--json")
	if !strings.Contains(legacy, `"legacy":true`) {
		t.Fatalf("legacy Mission did not use read-only shared decoder: %s", legacy)
	}
}

func TestCompactMissionLifecycleUsesMarkdownAndAtomicCommands(t *testing.T) {
	root, contractRef := fixture(t)
	plan := `---
type: MissionPlan
title: Ship compact mechanics
owner: Alex
contract:
  ref: ` + contractRef + `
outcome: The compact command loop works.
review: independent
completion:
  - claim: loop
    pass_boundary: The lifecycle completes.
    proof_requirement: Real command process assertions pass.
objectives:
  - outcome: Implement the loop.
    claims: [loop]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [internal/]
  semantic: [Compact lifecycle behavior.]
repair_budget: 1
dependencies: []
gaps: []
stops: [scope-drift]
---
# Mission

Implement the accepted plan.
`
	started := run(t, root, []byte(plan), 0, "mission", "start", "-", "--json")
	if !strings.Contains(started, `"ref":"M1"`) {
		t.Fatalf("start=%s", started)
	}
	replayed := run(t, root, []byte(plan), 0, "mission", "start", "-", "--json")
	if !strings.Contains(replayed, `"ref":"M1"`) {
		t.Fatalf("replay=%s", replayed)
	}
	if matches, _ := filepath.Glob(filepath.Join(root, ".spectacular", "missions", "M*-*")); len(matches) != 1 {
		t.Fatalf("retry created %d Missions", len(matches))
	}
	run(t, root, nil, 0, "objective", "promote", "M1/O1", "--json")
	run(t, root, nil, 0, "objective", "finish", "M1/O1", "--json")
	run(t, root, nil, 0, "run", "start", "M1", "--title", "Second execution boundary", "--json")
	activation := field(t, run(t, root, nil, 0, "mission", "check", "M1", "--json"), "fingerprint")
	commit := git(t, root, "rev-parse", "HEAD")
	tree := git(t, root, "rev-parse", "HEAD^{tree}")
	review := `---
type: ReviewDraft
title: Compact lifecycle review
status: passed
reviewed:
  commit: ` + commit + `
  tree: ` + tree + `
  activation_fingerprint: ` + activation + `
reviewer:
  actor: Fresh reviewer
  operator: Primary operator
  relation_to_operator: independent
  implemented_reviewed_scope: false
  independence_basis: Separate test reviewer bound to the exact tree.
  evidence: [test:independent-review]
claims:
  - claim: loop
    verdict: pass
findings: []
limitations: []
---
# Review

The claim passes.
`
	run(t, root, []byte(review), 0, "review", "record", "M1", "-", "--json")
	run(t, root, nil, 0, "mission", "complete", "M1", "--by", "Alex", "--json")
	checked := run(t, root, nil, 0, "mission", "check", "M1", "--json")
	if !strings.Contains(checked, `"valid":true`) {
		t.Fatalf("completed check=%s", checked)
	}
	shown := run(t, root, nil, 0, "mission", "show", "M1", "--json")
	if !strings.Contains(shown, `"status":"completed"`) || !strings.Contains(shown, `"file":"objectives/O1-implement-the-loop.md"`) {
		t.Fatalf("completed show=%s", shown)
	}
}

func TestRefusalIsTypedAndDoesNotMutate(t *testing.T) {
	root, _ := fixture(t)
	before := digestTree(t, filepath.Join(root, ".spectacular"))
	output := run(t, root, nil, 3, "objective", "finish", "M1/O1", "--json")
	for _, want := range []string{`"schema_version":"spectacular.refusal.v2"`, `"field":"M1"`, `"mutation":"none"`, `"safe_correction"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("refusal lacks %s: %s", want, output)
		}
	}
	if after := digestTree(t, filepath.Join(root, ".spectacular")); after != before {
		t.Fatal("refusal mutated workspace")
	}
}

func run(t *testing.T, root string, stdin []byte, want int, args ...string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	exit := (Runner{Cwd: root, Stdin: bytes.NewReader(stdin), Stdout: &stdout, Stderr: &stderr, Now: func() time.Time {
		return time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	}}).Run(args)
	if exit != want {
		t.Fatalf("spectacular %s exit=%d want=%d\nstdout=%s\nstderr=%s", strings.Join(args, " "), exit, want, stdout.String(), stderr.String())
	}
	return stdout.String() + stderr.String()
}

func fixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	if err := os.MkdirAll(filepath.Join(meta, "contracts"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(meta, "workspace.yaml"), "schema_version: spectacular.workspace.v1\nrecord_roots:\n  - .\nproject_anchor: PROJECT.md\n")
	write(t, filepath.Join(meta, "PROJECT.md"), "---\ntype: Anchor\nid: 0199a000-0000-7000-8000-000000000001\ntitle: Test project\nref: PROJECT\n---\n# Project\n")
	contractID := "0199a000-0000-7000-8000-000000000002"
	write(t, filepath.Join(meta, "contracts", "CC-test.md"), "---\ntype: Contract\nid: "+contractID+"\ntitle: Test Contract\nref: CC-test\nstatus: current\n---\n# Contract\n")
	git(t, root, "init", "-q")
	git(t, root, "config", "user.name", "Test")
	git(t, root, "config", "user.email", "test@example.invalid")
	git(t, root, "add", ".spectacular")
	git(t, root, "commit", "-qm", "fixture")
	return root, "Contract:" + contractID
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func field(t *testing.T, raw, name string) string {
	t.Helper()
	var value struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatal(err)
	}
	result, _ := value.Data[name].(string)
	if result == "" {
		t.Fatalf("missing %s in %s", name, raw)
	}
	return result
}

func digestTree(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		data, _ := os.ReadFile(path)
		hash.Write([]byte(rel))
		hash.Write(data)
		return nil
	})
	return hex.EncodeToString(hash.Sum(nil))
}
