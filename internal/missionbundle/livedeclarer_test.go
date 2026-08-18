package missionbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

// liveDeclarerWorkspace builds a Contract with one open Gap and one *live* Mission
// bound to it. declares controls whether that Mission's resolves_gaps names the
// Gap, which is the only thing separating the two sides of the rule.
func liveDeclarerWorkspace(t *testing.T, declares bool) (*discovery.Workspace, string, string) {
	t.Helper()
	root := t.TempDir()
	metadata := filepath.Join(root, ".spectacular")
	for _, directory := range []string{"contracts", filepath.Join("missions", "M1-live")} {
		if err := os.MkdirAll(filepath.Join(metadata, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(relative, body string) {
		if err := os.WriteFile(filepath.Join(metadata, filepath.FromSlash(relative)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("workspace.yaml", "schema_version: spectacular.workspace.v1\nrecord_roots:\n  - .\nproject_anchor: PROJECT.md\n")
	write("PROJECT.md", "---\ntype: Anchor\nid: 019fe381-5d61-7223-b362-03a5f99a7b15\ntitle: Live declarer fixture\n---\n\nAnchor.\n")

	const contract = `---
type: Contract
id: 01a00aae-8921-7b27-96a9-1a4c175e7dc8
human_ref: CC-livedeclarer
title: Live declarer fixture contract
purpose: exercise the declaring-Mission exemption
updated: "2026-08-16T00:00:00Z"
gaps:
    - ref: open-gap
      problem: something is unresolved
      blocked_on: an owner decision
---

Contract body.
`
	write("contracts/CC-livedeclarer.md", contract)
	fingerprint := digest([]byte(contract))

	resolves := ""
	if declares {
		resolves = "resolves_gaps:\n    - gap: open-gap\n      resolution: the frozen wording approved at activation\n"
	}
	requires := "[]"
	if declares {
		requires = "\n        - amend-contract"
	}

	mission := `---
type: Mission
id: 01a010a6-01b0-7320-acc2-5c695bec2910
ref: M1
title: Live mission
status: active
owner: Alex
created: "2026-08-16T00:00:00Z"
updated: "2026-08-16T00:00:00Z"
contract:
    fingerprint: ` + fingerprint + `
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc8
outcome: a bound outcome
review: automatic
completion:
    - claim: bound
      pass_boundary: it is bound
      proof_requirement: tests
objectives:
    - claims:
        - bound
      id: 01a010a6-01b0-7320-acc2-5c695bec2911
      outcome: an objective
      ref: O1
      status: implemented
baseline:
    branch: main
    commit: "0000000000000000000000000000000000000000"
activation:
    at: "2026-08-16T00:00:00Z"
    by: Alex
    fingerprint: sha256:0000000000000000000000000000000000000000000000000000000000000000
run:
    id: 01a010a6-01b0-7320-acc2-5c695bec2912
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-16T00:00:00Z"
    status: active
validation:
    mode: cli
    schema: mission.v2
authority:
    operator: []
    requires_owner: ` + requires + `
scope:
    mechanical: []
    semantic: []
repair_budget: 0
dependencies: []
gaps: []
` + resolves + `stops:
    - a stop
---

Mission body.
`
	write("missions/M1-live/M1-live.md", mission)

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return ws, "Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc8", "open-gap"
}

// Completion refuses while a declared Gap is open and amendment refused while the
// declaring Mission was live, which deadlocked the first Mission ever to declare
// one. The Mission that froze the resolution wording at its own activation gate is
// the Mission authorized to close that Gap, so it is exempt from the live guard.
func TestTheDeclaringMissionMayAmendWhileItIsLive(t *testing.T) {
	ws, contractRef, gap := liveDeclarerWorkspace(t, true)
	service := Service{Workspace: ws}
	result, err := service.AmendContract(contractRef, gap, "Alex", "", true)
	if err != nil {
		t.Fatalf("the declaring Mission must be able to close its own Gap while live, got: %v", err)
	}
	if result.Mission != "M1" || result.Source != "mission-declared" {
		t.Fatalf("amendment mission=%q source=%q, want M1/mission-declared", result.Mission, result.Source)
	}
	if !strings.Contains(result.Resolution, "approved at activation") {
		t.Fatalf("resolution did not come from the Mission's frozen wording: %q", result.Resolution)
	}
}

// The exemption is narrow. A live bound Mission that never declared this Gap still
// has the Contract constraining work in flight, and still blocks.
func TestALiveMissionThatDeclaresNothingStillBlocks(t *testing.T) {
	ws, contractRef, gap := liveDeclarerWorkspace(t, false)
	service := Service{Workspace: ws}
	if _, err := service.AmendContract(contractRef, gap, "Alex", "an owner-supplied resolution", true); err == nil {
		t.Fatal("a live bound Mission that declares nothing must still block the amendment")
	} else if !strings.Contains(err.Error(), "live") {
		t.Fatalf("refusal does not name the live Mission: %v", err)
	}
}

// A completed Mission's binding is the historical fact of what it agreed to.
// Re-pointing it would write today's Contract over that fact, destroying the record
// re-pointing was meant to protect; the stale binding is instead reported as a
// contract-drift notice and stays recoverable through `git log -S <fingerprint>`.
// Only the live Mission, which is still working against the Contract, re-points.
func TestOnlyALiveMissionIsRepointed(t *testing.T) {
	ws, contractRef, gap := liveDeclarerWorkspace(t, true)

	// A second Mission on the same Contract, completed, bound at the same text.
	completed := filepath.Join(ws.Root, ".spectacular", "missions", "M2-done")
	if err := os.MkdirAll(completed, 0o755); err != nil {
		t.Fatal(err)
	}
	live, err := os.ReadFile(filepath.Join(ws.Root, ".spectacular", "missions", "M1-live", "M1-live.md"))
	if err != nil {
		t.Fatal(err)
	}
	done := strings.NewReplacer(
		"ref: M1", "ref: M2",
		"status: active", "status: completed",
		"5c695bec2910", "5c695bec2920",
		"5c695bec2911", "5c695bec2921",
		"5c695bec2912", "5c695bec2922",
	).Replace(string(live))
	if err := os.WriteFile(filepath.Join(completed, "M2-done.md"), []byte(done), 0o644); err != nil {
		t.Fatal(err)
	}
	reopened, err := discovery.Open(ws.Root)
	if err != nil {
		t.Fatal(err)
	}

	before, err := os.ReadFile(filepath.Join(completed, "M2-done.md"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := (Service{Workspace: reopened}).AmendContract(contractRef, gap, "Alex", "", false)
	if err != nil {
		t.Fatalf("the amendment must succeed, got: %v", err)
	}
	if len(result.Repointed) != 1 || result.Repointed[0] != "M1" {
		t.Fatalf("re-pointed %v, want only the live Mission M1", result.Repointed)
	}
	after, err := os.ReadFile(filepath.Join(completed, "M2-done.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("a completed Mission's binding was rewritten; its historical record must survive verbatim")
	}
}

// An override's wording was typed at a prompt, not frozen at an activation gate, so
// it has no declaring Mission to be. The exemption must not become a way to hand-write
// a resolution past the live guard by naming a Gap some live Mission happens to declare.
func TestAnOverrideIsNotExemptFromTheLiveGuard(t *testing.T) {
	ws, contractRef, gap := liveDeclarerWorkspace(t, true)
	service := Service{Workspace: ws}
	if _, err := service.AmendContract(contractRef, gap, "Alex", "wording typed at a prompt", true); err == nil {
		t.Fatal("an owner override must not inherit the declaring Mission's exemption")
	} else if !strings.Contains(err.Error(), "live") {
		t.Fatalf("refusal does not name the live Mission: %v", err)
	}
}
