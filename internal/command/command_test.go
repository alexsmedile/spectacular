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

	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func TestPublicRegistryIsMinimalAndTyped(t *testing.T) {
	// Twenty-four commands. Owner authorized init and guard.
	want := []string{
		"mission start", "mission list", "mission show", "mission check", "mission amend-scope", "mission close", "objective show", "objective promote",
		"objective finish", "run show", "run start", "run transition", "review record", "handoff record",
		"evidence record", "mission complete", "proposal check", "campaign check", "contract amend", "contract create",
		"charter", "decide", "init", "guard",
	}
	if len(Registry) != len(want) {
		t.Fatalf("registry has %d commands, want %d", len(Registry), len(want))
	}
	if len(want) != 24 {
		t.Fatalf("the public surface is %d commands; owner authorized 24", len(want))
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

// TestDeletedV1SurfaceStaysDeleted asserts the removed v1 packages cannot return
// by import. The context-compiler chain (context, projection, guardrails) was
// unreachable from the main package and removed as a unit; index was the v1
// predecessor of discovery.Workspace.Lookup and removed after. A reintroduced
// import is the first symptom of any of them growing back, and it would compile
// silently without this check.
func TestDeletedV1SurfaceStaysDeleted(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	// Split so this file does not match its own search strings.
	modulePrefix := "github.com/alexsmedile/spectacular/v2/internal/"
	deleted := []string{"con" + "text", "pro" + "jection", "guard" + "rails", "in" + "dex"}
	self, err := filepath.Abs("command_test.go")
	if err != nil {
		t.Fatal(err)
	}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			// Skip trees that are not this module's source. `.claude/worktrees/`
			// may hold agent worktrees: separate checkouts at arbitrary commits,
			// so one sitting before a deletion would report it as live. Empty is
			// normal. `_archive` and `_backups` keep superseded copies on purpose,
			// and `node_modules` is vendored.
			switch entry.Name() {
			case ".git", ".claude", "_archive", "_backups", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || path == self {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		for _, name := range deleted {
			if strings.Contains(string(data), modulePrefix+name) {
				rel, _ := filepath.Rel(root, path)
				t.Errorf("%s references deleted package %s%s", rel, modulePrefix, name)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestWorkingTreeHasNoUnexplainedUntrackedPaths asserts every path is tracked,
// ignored by a stated rule, or absent. An untracked path is invisible to the
// record and silently lost on a fresh clone, so "not committed yet" and
// "deliberately excluded" must not look the same.
//
// It reports rather than fails when git is unavailable: this is a repository
// hygiene check, not a property of the code under test.
func TestWorkingTreeHasNoUnexplainedUntrackedPaths(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "status", "--porcelain", "--untracked-files=all")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("git unavailable in this environment: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if !strings.HasPrefix(line, "??") {
			continue
		}
		path := strings.TrimSpace(strings.TrimPrefix(line, "??"))
		// Records this Mission is actively producing are staged by the commit
		// that closes it, not by an ignore rule.
		if strings.HasPrefix(path, ".spectacular/missions/") {
			continue
		}
		t.Errorf("untracked path %q is neither tracked nor ignored by a stated rule", path)
	}
}

// TestGapsDoNotReferenceDeletedPackages asserts no open Gap points at something
// this repository no longer has. A Gap naming a deleted package is worse than a
// closed one: it reads as live work and sends a reader looking for code that is
// gone.
func TestGapsDoNotReferenceDeletedPackages(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	deleted := []string{"internal/context", "internal/projection", "internal/guardrails"}
	governance := filepath.Join(root, "internal", "governance")
	if _, statErr := os.Stat(governance); statErr != nil {
		t.Fatalf("internal/governance is referenced by an open Gap but is missing: %v", statErr)
	}
	err = filepath.WalkDir(filepath.Join(root, ".spectacular"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		// Only the gaps: block matters; body prose may discuss history freely.
		// The block runs until the next top-level key, so an empty `gaps: []`
		// contributes no lines rather than swallowing the rest of the file.
		var block []string
		inGaps := false
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "gaps:") {
				inGaps = true
				continue
			}
			if inGaps {
				indented := strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
				if strings.TrimSpace(line) != "" && !indented {
					break
				}
				block = append(block, line)
			}
		}
		for _, pkg := range deleted {
			if strings.Contains(strings.Join(block, "\n"), pkg) {
				rel, _ := filepath.Rel(root, path)
				t.Errorf("%s has a Gap referencing deleted package %s", rel, pkg)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
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

func TestCampaignCheckRendersOrderedProjection(t *testing.T) {
	root, _ := fixture(t)
	campaigns := filepath.Join(root, ".spectacular", "campaigns")
	if err := os.MkdirAll(campaigns, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(campaigns, "launch.md")
	write(t, path, `---
type: Campaign
schema: spectacular.campaign.v2
title: Launch readiness
focus: Ship a safe release path.
current: B2
exit_condition: The release path is proven.
blocks:
  - ref: B1
    title: CI foundation
    state: complete
    after: []
    missions: []
  - ref: B2
    title: Release hardening
    state: active
    after: [B1]
    missions: []
---
# Campaign: Launch readiness
`)
	output := run(t, root, nil, 0, "campaign", "check", ".spectacular/campaigns/launch.md", "--json")
	for _, want := range []string{`"schema_version":"spectacular.campaign.check.v2"`, `"current":"B2"`, `"order":["B1","B2"]`, `"mermaid":"flowchart LR`} {
		if !strings.Contains(output, want) {
			t.Fatalf("campaign check lacks %q:\n%s", want, output)
		}
	}
	human := run(t, root, nil, 0, "campaign", "check", ".spectacular/campaigns/launch.md")
	if !strings.Contains(human, "CURRENT CAMPAIGN BLOCK: B2 — Release hardening") || !strings.Contains(human, "ORDER") || !strings.Contains(human, "MERMAID") {
		t.Fatalf("human campaign projection=%s", human)
	}

	asciiOut := run(t, root, nil, 0, "campaign", "check", ".spectacular/campaigns/launch.md", "--ascii")
	if !strings.Contains(asciiOut, "Campaign DAG: Launch") || !strings.Contains(asciiOut, "B1") || !strings.Contains(asciiOut, "B2") {
		t.Fatalf("ascii campaign projection=%s", asciiOut)
	}
	write(t, path, `---
type: Campaign
schema: spectacular.campaign.v2
title: Cycle
focus: Demonstrate a cycle.
current: B1
exit_condition: Never.
blocks:
  - ref: B1
    title: First
    state: planned
    after: [B2]
  - ref: B2
    title: Second
    state: planned
    after: [B1]
---
# Campaign: Cycle
`)
	cycle := run(t, root, nil, 3, "campaign", "check", ".spectacular/campaigns/launch.md", "--json")
	if !strings.Contains(cycle, "contains a cycle") {
		t.Fatalf("cycle refusal=%s", cycle)
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
allow_main: true
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
	inline := run(t, root, nil, 0, "objective", "show", "M1/O1", "--json")
	run(t, root, nil, 0, "objective", "promote", "M1/O1", "--json")
	promoted := run(t, root, nil, 0, "objective", "show", "M1/O1", "--json")
	for _, name := range []string{"id", "ref", "outcome", "status"} {
		if field(t, inline, name) != field(t, promoted, name) {
			t.Fatalf("promotion changed %s", name)
		}
	}
	run(t, root, nil, 0, "objective", "finish", "M1/O1", "--json")
	run(t, root, nil, 0, "run", "start", "M1", "--title", "Second execution boundary", "--json")
	run(t, root, nil, 0, "run", "start", "M1", "--title", "Third execution boundary", "--json")
	if third := run(t, root, nil, 0, "run", "show", "M1/R3", "--json"); !strings.Contains(third, `"status":"active"`) {
		t.Fatalf("third Run=%s", third)
	}
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
	missions, _ := filepath.Glob(filepath.Join(root, ".spectacular", "missions", "M1-*", "M1-*.md"))
	data, err := os.ReadFile(missions[0])
	if err != nil || !strings.Contains(string(data), "Implement the accepted plan.") || !strings.Contains(string(data), "start_key: sha256:") {
		t.Fatalf("round trip lost Markdown or unknown fields: %v\n%s", err, data)
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

// TestHumanRefusalSurfacesTheUnderlyingCause asserts the human output carries the
// wrapped cause, not only the field. A YAML syntax error surfaces as
// invalid_known_field on `input`, which reads as "you named a field that does not
// exist" and sends the reader hunting through field names. The parser's line
// number is the only text that points at the real fault, and the JSON envelope
// already carried it while the human path dropped it.
func TestHumanRefusalSurfacesTheUnderlyingCause(t *testing.T) {
	root, contractRef := fixture(t)
	// A scalar opening with a backtick is a YAML reserved indicator, not a field error.
	plan := "---\ntype: MissionPlan\ntitle: Probe\nowner: Alex\ncontract:\n  ref: " +
		contractRef + "\noutcome: `backtick opens this scalar`\n---\n# Probe\n"

	human := run(t, root, []byte(plan), 3, "mission", "start", "-")
	if !strings.Contains(human, "cause:") || !strings.Contains(human, "line 6") {
		t.Fatalf("human refusal hides the parser cause and its line number: %s", human)
	}

	// The JSON envelope must keep carrying it in the `actual` field.
	structured := run(t, root, []byte(plan), 3, "mission", "start", "-", "--json")
	if !strings.Contains(structured, "line 6") {
		t.Fatalf("JSON refusal lost the parser cause: %s", structured)
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
		if rel == filepath.Join("transactions", ".lock") {
			return nil
		}
		data, _ := os.ReadFile(path)
		hash.Write([]byte(rel))
		hash.Write(data)
		return nil
	})
	return hex.EncodeToString(hash.Sum(nil))
}

// The frozen boundary names `show` first: a Mission with inline Objectives and
// the same Mission with those Objectives promoted must render identically. The
// package-level equivalence test compares derived structures, which leaves the
// human renderer — the surface a reader actually sees — unguarded. This drives
// the real CLI, so a change to renderHuman that leaks storage layout into the
// output fails here.
func TestShowRendersIdenticallyAcrossInlineAndPromotedObjectives(t *testing.T) {
	root, contractRef := fixture(t)
	plan := `---
type: MissionPlan
title: Render equivalence across representations
owner: Alex
contract:
  ref: ` + contractRef + `
outcome: Promotion changes storage, never the rendered conclusion.
review: independent
completion:
  - claim: equivalence
    pass_boundary: Inline and promoted bundles render identically.
    proof_requirement: Real command output is compared byte for byte.
objectives:
  - outcome: Extract the shared derivation layer.
    claims: [equivalence]
  - outcome: Render the compact state line.
    claims: [equivalence]
    after: [O1]
  - outcome: Compute per-claim drift flags.
    claims: [equivalence]
    after: [O1]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [internal/]
  semantic: [Representation equivalence.]
repair_budget: 1
allow_main: true
dependencies: []
gaps: []
stops: [scope-drift]
---
# Mission

Prove the renderer is blind to storage layout.
`
	run(t, root, []byte(plan), 0, "mission", "start", "-", "--json")

	// Every surface a reader can reach without --json. Adding a human-rendered
	// surface without adding it here leaves it unguarded, which is the failure
	// this test exists to prevent.
	surfaces := [][]string{
		{"mission", "show", "M1"},
		{"mission", "show", "M1", "--graph"},
		{"mission", "show", "M1", "--timeline"},
		{"objective", "show", "M1/O1"},
		{"objective", "show", "M1/O2"},
		{"mission", "check", "M1"},
	}
	inline := map[string]string{}
	for _, args := range surfaces {
		inline[strings.Join(args, " ")] = run(t, root, nil, 0, args...)
	}

	// Promote a root and a leaf, so the bundle mixes both representations rather
	// than converting wholesale.
	for _, ref := range []string{"M1/O1", "M1/O3"} {
		run(t, root, nil, 0, "objective", "promote", ref, "--json")
	}
	if matches, _ := filepath.Glob(filepath.Join(root, ".spectacular", "missions", "M1-*", "objectives", "*.md")); len(matches) != 2 {
		t.Fatalf("promotion wrote %d Objective files, want 2", len(matches))
	}

	for _, args := range surfaces {
		name := strings.Join(args, " ")
		promoted := run(t, root, nil, 0, args...)
		if inline[name] != promoted {
			t.Errorf("`spectacular %s` differs between representations:\n--- inline ---\n%s\n--- promoted ---\n%s", name, inline[name], promoted)
		}
	}
}

func TestTimelineCommandValidation(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("timeline on valid mission", func(t *testing.T) {
		out := run(t, root, nil, 0, "mission", "show", "M6", "--timeline")
		if !strings.Contains(out, "Timeline") || !strings.Contains(out, "M6 · ") {
			t.Fatalf("expected timeline output, got:\n%s", out)
		}
	})

	t.Run("cannot combine graph and timeline", func(t *testing.T) {
		var out bytes.Buffer
		runner := Runner{Cwd: root, Stdout: &out, Stderr: &out}
		code := runner.Run([]string{"mission", "show", "M6", "--graph", "--timeline"})
		if code == 0 {
			t.Fatalf("expected failure when combining --graph and --timeline")
		}
	})
}

func TestRepairExhaustionSurfacesAllFallbacksAndNeverSuppressesAlternatives(t *testing.T) {
	// Test directly against renderHuman with a Bundle at repair exhaustion
	bundle := &missionbundle.Bundle{
		Ref:          "M8",
		Title:        "Freeze the schema and record what was asked for",
		Status:       "active",
		RepairBudget: 2,
		Run: &missionbundle.Run{
			Ref:     "R1",
			Status:  "active",
			Repairs: 2,
		},
		Fallbacks: []missionbundle.Fallback{
			{
				Approach:        "Keep two package roots",
				RejectedBecause: "Doubles decode surface",
				InvalidatedIf:   "Single decoder fails",
				Recommendation:  true,
			},
			{
				Approach:        "Generate parallel schema adapters",
				RejectedBecause: "High overhead",
				InvalidatedIf:   "Unmaintainable",
				Recommendation:  false,
			},
		},
	}

	var out bytes.Buffer
	renderHuman(&out, bundle)
	rendered := out.String()

	if !strings.Contains(rendered, "FALLBACKS") {
		t.Fatalf("rendered output lacks FALLBACKS section:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Keep two package roots [recommendation]") {
		t.Fatalf("rendered output lacks recommendation label:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Generate parallel schema adapters") {
		t.Fatalf("rendered output suppressed alternative fallback:\n%s", rendered)
	}
}

func TestHelpAndSchemaFlags(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("global help", func(t *testing.T) {
		outHelp := run(t, root, nil, 0, "--help")
		if !strings.Contains(outHelp, "spectacular — governed execution") || !strings.Contains(outHelp, "spectacular mission start") {
			t.Fatalf("expected global help, got:\n%s", outHelp)
		}

		outH := run(t, root, nil, 0, "-h")
		if !strings.Contains(outH, "spectacular — governed execution") {
			t.Fatalf("expected global help with -h, got:\n%s", outH)
		}

		outJSON := run(t, root, nil, 0, "--help", "--json")
		if !strings.Contains(outJSON, `"schema_version":"spectacular.command-catalog.v1"`) {
			t.Fatalf("expected json envelope for global help, got:\n%s", outJSON)
		}
	})

	t.Run("subcommand help with templates", func(t *testing.T) {
		startHelp := run(t, root, nil, 0, "mission", "start", "--help")
		if !strings.Contains(startHelp, "type: MissionPlan") || !strings.Contains(startHelp, "spectacular mission start") {
			t.Fatalf("expected mission start template, got:\n%s", startHelp)
		}

		startH := run(t, root, nil, 0, "mission", "start", "-h")
		if !strings.Contains(startH, "type: MissionPlan") {
			t.Fatalf("expected mission start template with -h, got:\n%s", startH)
		}

		reviewHelp := run(t, root, nil, 0, "review", "record", "--help")
		if !strings.Contains(reviewHelp, "type: ReviewDraft") {
			t.Fatalf("expected review record template, got:\n%s", reviewHelp)
		}

		handoffHelp := run(t, root, nil, 0, "handoff", "record", "--help")
		if !strings.Contains(handoffHelp, "type: HandoffDraft") {
			t.Fatalf("expected handoff record template, got:\n%s", handoffHelp)
		}

		jsonHelp := run(t, root, nil, 0, "mission", "start", "--help", "--json")
		if !strings.Contains(jsonHelp, `"schema_version":"spectacular.command-help.v1"`) || !strings.Contains(jsonHelp, `"input_type":"MissionPlan"`) {
			t.Fatalf("expected json help envelope, got:\n%s", jsonHelp)
		}
	})

	t.Run("subcommand schema inspection", func(t *testing.T) {
		for _, spec := range Registry {
			cmdArgs := append([]string(nil), spec.Words...)
			cmdArgs = append(cmdArgs, "--schema")
			schemaOut := run(t, root, nil, 0, cmdArgs...)
			if !strings.Contains(schemaOut, `"schema_version":"spectacular.command-schema.v1"`) || !strings.Contains(schemaOut, spec.JSONSchema) {
				t.Fatalf("schema output for %v missing schema %s:\n%s", spec.Words, spec.JSONSchema, schemaOut)
			}
		}
	})
}

func TestCharterAndDecideCLI(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("charter command human and json", func(t *testing.T) {
		outHuman := run(t, root, nil, 0, "charter", "M17/O1")
		if !strings.Contains(outHuman, "FROZEN TRUTH") || !strings.Contains(outHuman, "Tokens:") {
			t.Fatalf("expected charter markdown output, got:\n%s", outHuman)
		}

		outJSON := run(t, root, nil, 0, "charter", "M17/O1", "--json")
		if !strings.Contains(outJSON, `"schema_version":"spectacular.charter.show.v2"`) || !strings.Contains(outJSON, `"token_count"`) {
			t.Fatalf("expected charter json envelope, got:\n%s", outJSON)
		}

		outPrompt := run(t, root, nil, 0, "charter", "M17/O1", "--prompt")
		if !strings.Contains(outPrompt, "# Objective: M17/O1") || !strings.Contains(outPrompt, "## 1. Initial Code Grounding") || !strings.Contains(outPrompt, "## 3. Worker Contract") {
			t.Fatalf("expected charter prompt output, got:\n%s", outPrompt)
		}
	})

	t.Run("decide command help", func(t *testing.T) {
		outHelp := run(t, root, nil, 0, "decide", "--help")
		if !strings.Contains(outHelp, "type: DecisionDraft") || !strings.Contains(outHelp, "spectacular decide") {
			t.Fatalf("expected decide help template, got:\n%s", outHelp)
		}
	})

	t.Run("guard command execution", func(t *testing.T) {
		outPass := run(t, root, nil, 0, "guard", "M17/O1", "--", "echo", "clean")
		if !strings.Contains(outPass, "Perimeter Guard: PASS") {
			t.Fatalf("expected guard pass, got:\n%s", outPass)
		}

		outExec := run(t, root, nil, 0, "guard", "M17/O1", "--exec", "echo")
		if !strings.Contains(outExec, "Perimeter Guard: PASS") || !strings.Contains(outExec, "# Objective: M17/O1") {
			t.Fatalf("expected guard --exec pass with prompt output, got:\n%s", outExec)
		}

		outRefuse := run(t, root, nil, 3, "guard", "M17/O1", "--exec", "echo", "--", "echo")
		if !strings.Contains(outRefuse, "cannot combine --exec") {
			t.Fatalf("expected refusal on combining --exec and --, got:\n%s", outRefuse)
		}
	})

	t.Run("mission check verify mode with token efficiency", func(t *testing.T) {
		outJSON := run(t, root, nil, 0, "mission", "check", "M17", "--json")
		if !strings.Contains(outJSON, `"valid":true`) {
			t.Fatalf("expected mission check valid=true, got:\n%s", outJSON)
		}
	})
}

// TestSchemaTemplatesSatisfyTheirOwnValidator fills every published input template
// with plausible values and feeds it back to the command that emitted it.
//
// This is the round trip that three shipped template bugs slipped through: a
// Campaign template naming `depends_on:` where the parser reads `after:`, the same
// template omitting the required `state:`, and a Proposal template marking the
// required `target_contract` optional. Each was invisible to a test that only
// asserted --schema returns something, because YAML drops unknown keys in silence.
// A template is an authoring interface; an agent that follows it must land on a
// document the validator accepts.
func TestSchemaTemplatesSatisfyTheirOwnValidator(t *testing.T) {
	placeholders := strings.NewReplacer(
		"<title>", "Round trip fixture",
		"<first block title>", "Foundation",
		"<second block title>", "Hardening",
		"<what this campaign is steering toward>", "Prove the template validates.",
		"<observable condition that ends the campaign>", "The template round trips.",
		"<owner>", "Alex",
		"<uuidv7>", "01a02faa-ed22-736d-94a0-e1596184921f",
		"P<N>", "P99",
		"Contract:<uuidv7>", "Contract:019fe381-5d61-7223-b362-03a5f99a7b10",
		"<RFC3339>", "2026-08-23T12:00:00Z",
	)

	for _, testCase := range []struct {
		name     string
		words    []string
		filename string
		check    func(t *testing.T, root, path string)
	}{
		{
			name:     "campaign",
			words:    []string{"campaign", "check"},
			filename: filepath.Join(".spectacular", "campaigns", "round-trip.md"),
			check: func(t *testing.T, root, path string) {
				// Validation passing is not enough. YAML drops an unknown key in
				// silence, so a template naming the wrong field for dependencies
				// still validates while every edge disappears. Assert the parsed
				// projection carries the semantics the template claimed.
				output := run(t, root, nil, 0, "campaign", "check", filepath.ToSlash(path), "--json")
				if !strings.Contains(output, `"after":["B1"]`) {
					t.Fatalf("template dependency edge did not survive parsing; the key it names is not the key the parser reads:\n%s", output)
				}
			},
		},
		{
			name:     "proposal",
			words:    []string{"proposal", "check"},
			filename: filepath.Join(".spectacular", "proposals", "P99-round-trip.md"),
			check: func(t *testing.T, root, path string) {
				run(t, root, nil, 0, "proposal", "check", "P99", "--json")
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root, _ := fixture(t)
			var spec *Spec
			for i := range Registry {
				if strings.Join(Registry[i].Words, " ") == strings.Join(testCase.words, " ") {
					spec = &Registry[i]
					break
				}
			}
			if spec == nil {
				t.Fatalf("no spec for %v", testCase.words)
			}
			if spec.Template == "" {
				t.Fatalf("%v publishes no template", testCase.words)
			}
			path := filepath.Join(root, testCase.filename)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			document := placeholders.Replace(spec.Template)
			// A template that still carries an unresolved placeholder cannot be
			// round tripped, and silently substituting one would hide that.
			if strings.Contains(document, "<") && strings.Contains(document, ">") {
				for _, line := range strings.Split(document, "\n") {
					if strings.Contains(line, "<") && strings.Contains(line, ">") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
						t.Fatalf("template line has no placeholder substitution: %q", line)
					}
				}
			}
			write(t, path, document)
			testCase.check(t, root, testCase.filename)
		})
	}
}

func TestContractCreateCLI(t *testing.T) {
	root, _ := fixture(t)
	out := run(t, root, nil, 0, "contract", "create", "CC-analytics-pipeline", "--title", "Analytics Pipeline", "--json")
	if !strings.Contains(out, `"operation":"contract.create"`) || !strings.Contains(out, `"ref":"CC-analytics-pipeline"`) {
		t.Fatalf("contract create --json output=%s", out)
	}
}

func TestMissionAmendScopeAndCloseCLI(t *testing.T) {
	root, contractRef := fixture(t)
	plan := `---
type: MissionPlan
title: Command execution mission
owner: Alex
contract:
  ref: ` + contractRef + `
outcome: Test amend-scope and close.
review: automatic
completion:
  - claim: smoke
    pass_boundary: Tests pass.
    proof_requirement: Verified.
objectives:
  - outcome: Run smoke test
    claims: [smoke]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [internal/]
  semantic: [Testing amend-scope.]
repair_budget: 1
allow_main: true
dependencies: []
gaps: []
stops: [scope-drift]
---
# Description
`
	startOut := run(t, root, []byte(plan), 0, "mission", "start", "-", "--json")
	if !strings.Contains(startOut, `"operation":"mission.start"`) {
		t.Fatalf("mission start failed: %s", startOut)
	}

	// Amend scope with dry-run
	dryOut := run(t, root, nil, 0, "mission", "amend-scope", "M1", "--add", "Harbor/MenuBar.swift", "--by", "Alex", "--reason", "found menu bar dependency", "--dry-run", "--json")
	if !strings.Contains(dryOut, `"operation":"mission.amend_scope"`) || !strings.Contains(dryOut, `"fingerprint"`) {
		t.Fatalf("mission amend-scope dry-run failed: %s", dryOut)
	}

	// Amend scope real
	amendOut := run(t, root, nil, 0, "mission", "amend-scope", "M1", "--add", "Harbor/MenuBar.swift", "--by", "Alex", "--reason", "found menu bar dependency", "--json")
	if !strings.Contains(amendOut, `"operation":"mission.amend_scope"`) {
		t.Fatalf("mission amend-scope failed: %s", amendOut)
	}

	// Close mission
	closeOut := run(t, root, nil, 0, "mission", "close", "M1", "--by", "Alex", "--json")
	if !strings.Contains(closeOut, `"operation":"mission.complete"`) {
		t.Fatalf("mission close failed: %s", closeOut)
	}
}

func TestMissionListAndBranchAndEvidenceCLI(t *testing.T) {
	root, contractRef := fixture(t)
	plan := `---
type: MissionPlan
title: List and branch test
owner: Alex
contract:
  ref: ` + contractRef + `
outcome: Test mission list and branch guardrails.
review: automatic
completion:
  - claim: smoke
    pass_boundary: Tests pass.
    proof_requirement: Verified.
objectives:
  - outcome: Run smoke test
    claims: [smoke]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [internal/]
  semantic: [Testing branch.]
repair_budget: 1
dependencies: []
gaps: []
stops: [scope-drift]
---
# Description
`
	// Start with --create-branch
	startOut := run(t, root, []byte(plan), 0, "mission", "start", "-", "--create-branch", "--json")
	if !strings.Contains(startOut, `"operation":"mission.start"`) {
		t.Fatalf("mission start --create-branch failed: %s", startOut)
	}

	// Mission list json
	listOut := run(t, root, nil, 0, "mission", "list", "--json")
	if !strings.Contains(listOut, `"schema_version":"spectacular.mission.list.v2"`) || !strings.Contains(listOut, `"ref":"M1"`) {
		t.Fatalf("mission list --json failed: %s", listOut)
	}

	// Mission list human
	humanList := run(t, root, nil, 0, "mission", "list")
	if !strings.Contains(humanList, "REF") || !strings.Contains(humanList, "M1") {
		t.Fatalf("mission list human failed: %s", humanList)
	}

	// Evidence record from test output
	testJSON := filepath.Join(root, "test-output.json")
	write(t, testJSON, `{"Time":"2026-08-30T10:00:00Z","Action":"pass","Test":"TestFeature"}
`)
	evOut := run(t, root, nil, 0, "evidence", "record", "M1", "--from", testJSON, "--json")
	if !strings.Contains(evOut, `"operation":"evidence.record"`) {
		t.Fatalf("evidence record --from failed: %s", evOut)
	}

	// Evidence record from stdin
	commit := git(t, root, "rev-parse", "HEAD")
	tree := git(t, root, "rev-parse", "HEAD^{tree}")
	evStdinDraft := `---
type: EvidenceDraft
title: Test stdin evidence
actor: Alex
commit: ` + commit + `
tree: ` + tree + `
claims: [smoke]
checks:
  - name: test-check
    result: pass
---
# Evidence
Verified via stdin.
`
	evStdinOut := run(t, root, []byte(evStdinDraft), 0, "evidence", "record", "M1", "-", "--json")
	if !strings.Contains(evStdinOut, `"operation":"evidence.record"`) {
		t.Fatalf("evidence record from stdin failed: %s", evStdinOut)
	}
}

func TestCampaignValidateDerivesLiveState(t *testing.T) {
	root, _ := fixture(t)
	campDir := filepath.Join(root, ".spectacular", "campaigns")
	if err := os.MkdirAll(campDir, 0o755); err != nil {
		t.Fatal(err)
	}

	campaignDoc := `---
type: Campaign
schema: spectacular.campaign.v2
title: Test Campaign
focus: Test focus
current: B1
exit_condition: All complete
blocks:
  - ref: B1
    title: Phase 1
    state: planned
---
# Campaign Details
`
	campFile := filepath.Join(campDir, "test-campaign.md")
	write(t, campFile, campaignDoc)

	// Campaign check
	campOut := run(t, root, nil, 0, "campaign", "check", campFile, "--json")
	if !strings.Contains(campOut, `"live_state":"planned"`) || !strings.Contains(campOut, `"title":"Test Campaign"`) {
		t.Fatalf("campaign check failed: %s", campOut)
	}
}

func TestInitFreshWorkspace(t *testing.T) {
	root := t.TempDir()

	// 1. Run init in fresh directory
	out := run(t, root, nil, 0, "init", "--name", "My Cool Project")
	if !strings.Contains(out, "Initialized Spectacular workspace") {
		t.Fatalf("expected initialized message, got: %s", out)
	}
	if !strings.Contains(out, ".spectacular/workspace.yaml") || !strings.Contains(out, ".spectacular/PROJECT.md") {
		t.Fatalf("expected created files listed, got: %s", out)
	}

	// 2. Check workspace.yaml exists and is valid
	wsYAML, err := os.ReadFile(filepath.Join(root, ".spectacular", "workspace.yaml"))
	if err != nil {
		t.Fatalf("read workspace.yaml: %v", err)
	}
	if !strings.Contains(string(wsYAML), "schema_version: spectacular.workspace.v1") {
		t.Fatalf("unexpected workspace.yaml content: %s", string(wsYAML))
	}

	// 3. Check PROJECT.md exists with custom name
	projMD, err := os.ReadFile(filepath.Join(root, ".spectacular", "PROJECT.md"))
	if err != nil {
		t.Fatalf("read PROJECT.md: %v", err)
	}
	if !strings.Contains(string(projMD), "title: My Cool Project") {
		t.Fatalf("unexpected PROJECT.md content: %s", string(projMD))
	}

	// 4. Verify discovery and commands immediately work
	listOut := run(t, root, nil, 0, "mission", "list")
	if !strings.Contains(listOut, "No missions found.") {
		t.Fatalf("mission list failed after init: %s", listOut)
	}
}

func TestInitAvoidsOverwritesOnExistingWorkspace(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create custom workspace.yaml and custom PROJECT.md
	customWS := "schema_version: spectacular.workspace.v1\nrecord_roots: [.]\nproject_anchor: PROJECT.md\nguardrails: CUSTOM_GUARDRAILS.md\n"
	write(t, filepath.Join(meta, "workspace.yaml"), customWS)

	customProj := "---\ntype: Anchor\ntitle: Custom Original Title\ncurrent_truth: [.spectacular/PROJECT.md]\n---\n# Custom Original Content\n"
	write(t, filepath.Join(meta, "PROJECT.md"), customProj)

	// Run init again
	out := run(t, root, nil, 0, "init", "--name", "Different Name", "--json")
	if !strings.Contains(out, `"schema_version":"spectacular.init.v2"`) {
		t.Fatalf("expected spectacular.init.v2 JSON schema, got: %s", out)
	}
	if !strings.Contains(out, `"already_initialized":true`) {
		t.Fatalf("expected already_initialized: true, got: %s", out)
	}
	if !strings.Contains(out, `".spectacular/workspace.yaml"`) || !strings.Contains(out, `".spectacular/PROJECT.md"`) {
		t.Fatalf("expected skipped files listed, got: %s", out)
	}

	// Assert custom content was NOT overwritten
	afterWS, err := os.ReadFile(filepath.Join(meta, "workspace.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(afterWS) != customWS {
		t.Fatalf("workspace.yaml was overwritten! got: %s, want: %s", string(afterWS), customWS)
	}

	afterProj, err := os.ReadFile(filepath.Join(meta, "PROJECT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(afterProj) != customProj {
		t.Fatalf("PROJECT.md was overwritten! got: %s, want: %s", string(afterProj), customProj)
	}
}

func TestDecideDirectFlagsAndListFiltering(t *testing.T) {
	root, contractRef := fixture(t)

	// 1. Test direct decision recording via flags
	decideOut := run(t, root, nil, 0, "decide", "--title", "Direct Flag Decision", "--disposition", "accepted", "--rationale", "Testing direct flags without temp files", "--actor", "Alex", "--json")
	if !strings.Contains(decideOut, `"operation":"decision.record"`) || !strings.Contains(decideOut, `"ref":"D1-direct-flag-decision"`) {
		t.Fatalf("decide with direct flags failed: %s", decideOut)
	}

	// 1b. Test decision recording via stdin JSON payload
	jsonStdin := []byte(`{"title":"Stdin JSON Decision","disposition":"accepted","rationale":"Testing stdin JSON without temp files","actor":"Alex"}`)
	decideStdinOut := run(t, root, jsonStdin, 0, "decide", "-", "--json")
	if !strings.Contains(decideStdinOut, `"operation":"decision.record"`) || !strings.Contains(decideStdinOut, `"ref":"D2-stdin-json-decision"`) {
		t.Fatalf("decide with stdin JSON failed: %s", decideStdinOut)
	}

	// 1c. Test decision recording via inline JSON string
	decideInlineOut := run(t, root, nil, 0, "decide", `{"title":"Inline JSON Decision","disposition":"accepted","rationale":"Testing inline JSON argument","actor":"Alex"}`, "--json")
	if !strings.Contains(decideInlineOut, `"operation":"decision.record"`) || !strings.Contains(decideInlineOut, `"ref":"D3-inline-json-decision"`) {
		t.Fatalf("decide with inline JSON failed: %s", decideInlineOut)
	}

	// Verify the decision file was written inside .spectacular/decisions/
	decPath := filepath.Join(root, ".spectacular", "decisions", "D1-direct-flag-decision.md")
	content, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatalf("failed to read recorded decision file: %v", err)
	}
	if !strings.Contains(string(content), "Direct Flag Decision") || !strings.Contains(string(content), "Testing direct flags without temp files") {
		t.Fatalf("unexpected decision file content: %s", string(content))
	}

	// 2. Test mission list active-first vs --all
	plan := `---
type: MissionPlan
title: Test Mission Active List
owner: Alex
contract:
  ref: ` + contractRef + `
outcome: Test active mission list filtering.
review: automatic
completion:
  - claim: smoke
    pass_boundary: Smoke passes.
    proof_requirement: Verified.
objectives:
  - outcome: Implement smoke
    claims: [smoke]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [internal/]
  semantic: [Testing mission list.]
repair_budget: 1
dependencies: []
gaps: []
stops: [scope-drift]
---
# Description
`
	_ = run(t, root, []byte(plan), 0, "mission", "start", "-", "--create-branch", "--json")

	// Active list should show M1
	activeList := run(t, root, nil, 0, "mission", "list", "--json")
	if !strings.Contains(activeList, `"ref":"M1"`) {
		t.Fatalf("mission list active view should show M1: %s", activeList)
	}

	// Finish objective & complete mission
	_ = run(t, root, nil, 0, "objective", "finish", "M1/O1", "--json")
	_ = run(t, root, nil, 0, "mission", "complete", "M1", "--by", "Alex", "--json")

	// Now active list should NOT show completed M1
	afterComplete := run(t, root, nil, 0, "mission", "list", "--json")
	if strings.Contains(afterComplete, `"ref":"M1"`) {
		t.Fatalf("mission list active view should NOT show completed M1: %s", afterComplete)
	}

	// But mission list --all SHOULD show M1
	allList := run(t, root, nil, 0, "mission", "list", "--all", "--json")
	if !strings.Contains(allList, `"ref":"M1"`) {
		t.Fatalf("mission list --all should show completed M1: %s", allList)
	}
}

func TestUnifiedJSONPayloadProtocol(t *testing.T) {
	root, contractRef := fixture(t)

	// 1. Mission start with JSON payload via --data
	planJSON := `{
		"type": "MissionPlan",
		"title": "JSON Mission Start",
		"owner": "Alex",
		"contract": { "ref": "` + contractRef + `" },
		"outcome": "Verify JSON protocol across commands",
		"review": "automatic",
		"completion": [
			{ "claim": "json-protocol", "pass_boundary": "Passes", "proof_requirement": "Checked" }
		],
		"objectives": [
			{ "outcome": "Test JSON commands", "claims": ["json-protocol"] }
		],
		"authority": {
			"operator": ["inspect", "edit-in-scope", "choose-reversible-implementation", "run-checks", "generate-derived-files", "bounded-repair", "commit-local"],
			"requires_owner": ["activate-mission", "change-outcome-or-completion", "expand-scope", "push", "merge", "release", "irreversible-change", "destructive-data", "secret-change"]
		},
		"scope": {
			"mechanical": ["internal/"],
			"semantic": ["Testing JSON commands"]
		},
		"repair_budget": 1,
		"dependencies": [],
		"gaps": [],
		"stops": ["scope-drift"]
	}`
	startOut := run(t, root, nil, 0, "mission", "start", "--data", planJSON, "--create-branch", "--json")
	if !strings.Contains(startOut, `"operation":"mission.start"`) || !strings.Contains(startOut, `"ref":"M1"`) {
		t.Fatalf("mission start --data failed: %s", startOut)
	}

	// 2. Evidence record via --data
	evJSON := `{
		"type": "EvidenceDraft",
		"title": "JSON Evidence Test",
		"claims": ["json-protocol"]
	}`
	evOut := run(t, root, nil, 0, "evidence", "record", "M1", "--data", evJSON, "--json")
	if !strings.Contains(evOut, `"operation":"evidence.record"`) {
		t.Fatalf("evidence record --data failed: %s", evOut)
	}

	// 3. Handoff record via --data
	hoJSON := `{
		"type": "HandoffDraft",
		"title": "JSON Handoff Test",
		"task": "Perform automated review"
	}`
	hoOut := run(t, root, nil, 0, "handoff", "record", "M1", "--data", hoJSON, "--json")
	if !strings.Contains(hoOut, `"operation":"handoff.record"`) {
		t.Fatalf("handoff record --data failed: %s", hoOut)
	}

	// 4. Review record via --data
	rvJSON := `{
		"type": "ReviewDraft",
		"title": "JSON Review Test",
		"status": "passed",
		"claims": [
			{ "claim": "json-protocol", "verdict": "pass" }
		]
	}`
	rvOut := run(t, root, nil, 0, "review", "record", "M1", "--data", rvJSON, "--json")
	if !strings.Contains(rvOut, `"operation":"review.record"`) {
		t.Fatalf("review record --data failed: %s", rvOut)
	}

	// 5. Run transition via --data
	trJSON := `{
		"target": "M1/R1",
		"to": "paused",
		"by": "Alex",
		"reason": "testing transition json"
	}`
	trOut := run(t, root, nil, 0, "run", "transition", "M1/R1", "--data", trJSON, "--json")
	if !strings.Contains(trOut, `"operation":"run.transition"`) {
		t.Fatalf("run transition --data failed: %s", trOut)
	}

	// Transition back to active via inline JSON
	trActiveJSON := `{"target":"M1/R1","to":"active","by":"Alex","reason":"resume work"}`
	trActiveOut := run(t, root, nil, 0, "run", "transition", "M1/R1", trActiveJSON, "--json")
	if !strings.Contains(trActiveOut, `"operation":"run.transition"`) {
		t.Fatalf("run transition inline JSON failed: %s", trActiveOut)
	}

	// 6. Objective finish
	_ = run(t, root, nil, 0, "objective", "finish", "M1/O1", "--json")

	// 7. Mission complete without --by (auto-resolves default owner)
	completeOut := run(t, root, nil, 0, "mission", "complete", "M1", "--json")
	if !strings.Contains(completeOut, `"operation":"mission.complete"`) {
		t.Fatalf("mission complete without --by failed: %s", completeOut)
	}
}
