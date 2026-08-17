package missionbundle

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// repointWorkspace builds a Contract with one open Gap and one live Mission bound
// to it, whose body carries `quoted` verbatim. The Mission's binding fingerprint is
// computed from the Contract on disk, so the fixture re-points for real rather than
// against an invented value.
//
// The Mission is live and declares the Gap because only a live Mission re-points:
// a completed Mission keeps its historical binding, and a fixture using one would
// pass this test vacuously by never reaching the code it exists to guard.
func repointWorkspace(t *testing.T, quoted func(fingerprint string) string) (*discovery.Workspace, string, string) {
	t.Helper()
	root := t.TempDir()
	metadata := filepath.Join(root, ".spectacular")
	for _, directory := range []string{"contracts", filepath.Join("missions", "M1-bound")} {
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
	write("PROJECT.md", "---\ntype: Anchor\nid: 019fe381-5d61-7223-b362-03a5f99a7b15\ntitle: Repoint fixture\n---\n\nAnchor.\n")

	const contract = `---
type: Contract
id: 01a00aae-8921-7b27-96a9-1a4c175e7dc7
human_ref: CC-repoint
title: Repoint fixture contract
purpose: exercise re-pointing
updated: "2026-08-16T00:00:00Z"
gaps:
    - ref: open-gap
      problem: something is unresolved
      blocked_on: an owner decision
---

Contract body.
`
	write("contracts/CC-repoint.md", contract)
	fingerprint := digest([]byte(contract))

	mission := `---
type: Mission
id: 01a010a6-01b0-7320-acc2-5c695bec2900
ref: M1
title: Bound mission
status: active
owner: Alex
created: "2026-08-16T00:00:00Z"
updated: "2026-08-16T00:00:00Z"
contract:
    fingerprint: ` + fingerprint + `
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc7
outcome: a bound outcome
review: automatic
completion:
    - claim: bound
      pass_boundary: it is bound
      proof_requirement: tests
objectives:
    - claims:
        - bound
      id: 01a010a6-01b0-7320-acc2-5c695bec2901
      outcome: an objective
      ref: O1
      status: implemented
baseline:
    branch: main
    commit: 0000000000000000000000000000000000000000
activation:
    at: "2026-08-16T00:00:00Z"
    by: Alex
    fingerprint: sha256:0000000000000000000000000000000000000000000000000000000000000000
run:
    id: 01a010a6-01b0-7320-acc2-5c695bec2902
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-16T00:00:00Z"
    status: completed
validation:
    mode: cli
    schema: mission.v2
authority:
    operator: []
    requires_owner:
        - amend-contract
resolves_gaps:
    - gap: open-gap
      resolution: the frozen wording approved at activation
scope:
    mechanical: []
    semantic: []
repair_budget: 0
dependencies: []
gaps: []
stops:
    - a stop
---

` + quoted(fingerprint) + `
`
	write("missions/M1-bound/MISSION.md", mission)

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return ws, "Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc7", "open-gap"
}

func TestRepointingRefusesAnAmbiguousFingerprint(t *testing.T) {
	tests := []struct {
		name    string
		body    func(string) string
		refuses bool
	}{
		{
			// The binding appears once, in the contract: block. This is the ordinary
			// case and must keep working exactly as it does today.
			name:    "exactly one occurrence re-points",
			body:    func(string) string { return "This Mission quotes no fingerprint." },
			refuses: false,
		},
		{
			// The Mission quotes its own binding in prose. Rewriting the first
			// occurrence would rewrite the prose and leave the real binding stale.
			name:    "two occurrences refuse",
			body:    func(fingerprint string) string { return "The bound Contract is " + fingerprint + "." },
			refuses: true,
		},
		{
			// The M9-shaped case: a body quoting several fingerprints, none of which
			// is this Mission's binding. Unambiguous, so it re-points.
			name: "a body quoting other fingerprints re-points",
			body: func(string) string {
				return "Earlier revisions were sha256:" + strings.Repeat("a", 64) +
					" and sha256:" + strings.Repeat("b", 64) +
					" and sha256:" + strings.Repeat("c", 64) + "."
			},
			refuses: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ws, contractRef, gap := repointWorkspace(t, test.body)
			before := treeDigest(t, filepath.Join(ws.Root, ".spectacular"))
			service := Service{Workspace: ws}
			result, err := service.AmendContract(contractRef, gap, "Alex", "", false)

			if !test.refuses {
				if err != nil {
					t.Fatalf("an unambiguous Mission must re-point, got: %v", err)
				}
				if len(result.Repointed) != 1 || result.Repointed[0] != "M1" {
					t.Fatalf("re-pointed %v, want the bound Mission M1", result.Repointed)
				}
				return
			}

			var refusal *domain.Refusal
			if !errors.As(err, &refusal) {
				t.Fatalf("expected a typed refusal, got %T: %v", err, err)
			}
			if refusal.Code != domain.RefusalInvalidKnownField || refusal.Field != "contract.fingerprint" {
				t.Fatalf("refusal code=%s field=%s", refusal.Code, refusal.Field)
			}
			// The refusal names the Mission, the file, and every occurrence, so a
			// reader can tell the binding from the quotation.
			for _, want := range []string{"M1", "MISSION.md"} {
				if !strings.Contains(refusal.Error(), want) && !strings.Contains(refusal.Actual, want) {
					t.Fatalf("refusal does not name %q: %+v", want, refusal)
				}
			}
			if !strings.Contains(refusal.Actual, ",") {
				t.Fatalf("refusal names only one occurrence: %q", refusal.Actual)
			}
			if !strings.Contains(refusal.Recovery, "wrote nothing") {
				t.Fatalf("refusal does not state that nothing was written: %q", refusal.Recovery)
			}
			// The whole amendment aborts: the Contract, its log, and every Mission
			// are left byte-identical.
			if after := treeDigest(t, filepath.Join(ws.Root, ".spectacular")); after != before {
				t.Fatal("a refused amendment modified the workspace")
			}
		})
	}
}

// --dry-run reports the ambiguity as a would-refuse. The preview is the only
// thing the owner reads before authorizing, so it must not promise a re-point
// that the real run would refuse.
func TestDryRunReportsTheAmbiguousFingerprint(t *testing.T) {
	ws, contractRef, gap := repointWorkspace(t, func(fingerprint string) string {
		return "The bound Contract is " + fingerprint + "."
	})
	before := treeDigest(t, filepath.Join(ws.Root, ".spectacular"))
	service := Service{Workspace: ws}
	if _, err := service.AmendContract(contractRef, gap, "Alex", "", true); err == nil {
		t.Fatal("a dry run must report the ambiguity as a would-refuse")
	}
	if after := treeDigest(t, filepath.Join(ws.Root, ".spectacular")); after != before {
		t.Fatal("a dry run wrote to the workspace")
	}
}
