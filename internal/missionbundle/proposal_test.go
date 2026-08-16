package missionbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func noticeContaining(notices []string, want string) bool {
	for _, notice := range notices {
		if strings.Contains(notice, want) {
			return true
		}
	}
	return false
}

// proposalWorkspace writes the minimum workspace a Proposal can be validated
// in: the manifest plus the anchor and guardrails it names.
func proposalWorkspace(t *testing.T, name, body string) *discovery.Workspace {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".spectacular")
	if err := os.MkdirAll(filepath.Join(dir, "proposals"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		"workspace.yaml": "schema_version: spectacular.workspace.v1\nrecord_roots:\n  - .\nproject_anchor: PROJECT.md\nguardrails: GUARDRAILS.md\n",
		"PROJECT.md":     "---\ntype: Anchor\nid: 019fe381-5d61-7223-b362-03a5f99a7b13\ntitle: Test workspace\nupdated: \"2026-08-16T10:31:30Z\"\n---\n# Test workspace\n",
		"GUARDRAILS.md":  "# Guardrails\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "proposals", name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return ws
}

// The compact schema records the shape practice settled on. P5 and P6 must
// validate unchanged, and the older Proposals must validate with their legacy
// fields preserved rather than refused -- a Proposal that was accepted months
// ago is not rewritten to satisfy a schema written afterwards.
func TestProposalsValidateAgainstTheCompactSchema(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, testCase := range []struct {
		ref           string
		wantCandidate bool
	}{
		{"P1", false},
		{"P2", true},
		{"P3", true},
		{"P4", true},
		{"P5", false},
		{"P6", false},
	} {
		t.Run(testCase.ref, func(t *testing.T) {
			check, checkErr := ValidateProposal(ws, testCase.ref)
			if checkErr != nil {
				t.Fatalf("the compact schema must accept an authored Proposal: %v", checkErr)
			}
			if !check.Valid || check.Ref != testCase.ref {
				t.Fatalf("valid=%t ref=%q, want true and %q", check.Valid, check.Ref, testCase.ref)
			}
			if check.ID == "" {
				t.Fatal("a Proposal must resolve a UUIDv7 identity")
			}
			// Every Proposal in this repository predates the ref rename.
			if !noticeContaining(check.Notices, "ref-spelling-drift") {
				t.Fatalf("%s uses human_ref and must report the drift, got %v", testCase.ref, check.Notices)
			}
			if got := noticeContaining(check.Notices, "legacy-candidate-body"); got != testCase.wantCandidate {
				t.Fatalf("candidate body notice=%t, want %t (notices: %v)", got, testCase.wantCandidate, check.Notices)
			}
			// The superseded body is reported, never a validation failure: the
			// record stays exactly as it was accepted.
			if testCase.wantCandidate && !check.Valid {
				t.Fatal("a legacy candidate body is reported, not refused")
			}
		})
	}
}

// A missing required property is a refusal naming the exact field, so an author
// can correct it without guessing. Proposals are hand-written, which is
// precisely why the refusal has to be specific.
func TestProposalRefusalsNameTheMissingField(t *testing.T) {
	ws := proposalWorkspace(t, "P9-missing-target.md", `---
type: Proposal
id: 01a00a93-4757-7547-b64e-e91d2c291ce4
ref: P9
title: A Proposal missing its target
status: draft
created_by: Alex
created: "2026-08-16T12:36:59Z"
updated: "2026-08-16T12:36:59Z"
---
# A Proposal missing its target
`)
	_, err := ValidateProposal(ws, "P9")
	if err == nil {
		t.Fatal("a Proposal without target_contract must be refused")
	}
	if !strings.Contains(err.Error(), "target_contract") {
		t.Fatalf("the refusal must name the missing field, got %v", err)
	}
}

// The current spelling is accepted without a drift notice, which is what a
// newly authored Proposal looks like.
func TestProposalWithCurrentRefSpellingReportsNoDrift(t *testing.T) {
	ws := proposalWorkspace(t, "P9-compact.md", `---
type: Proposal
id: 01a00a93-4757-7547-b64e-e91d2c291ce4
ref: P9
title: A compact Proposal
status: draft
created_by: Alex
created: "2026-08-16T12:36:59Z"
updated: "2026-08-16T12:36:59Z"
scope:
    - v2
target_contract: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
---
# A compact Proposal
`)
	check, err := ValidateProposal(ws, "P9")
	if err != nil {
		t.Fatalf("the compact schema must accept its own reference shape: %v", err)
	}
	if len(check.Notices) != 0 {
		t.Fatalf("a compact Proposal reports nothing, got %v", check.Notices)
	}
	if !check.Valid || check.Ref != "P9" {
		t.Fatalf("valid=%t ref=%q", check.Valid, check.Ref)
	}
}
