package missionbundle

import (
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/governance"
)

const transactionDecisionDraft = `---
type: DecisionDraft
title: Keep exact Decision payloads atomic
actor: Alex
actor_role: owner
question: Should decision recording be all-or-nothing?
disposition: accepted
rationale: The CLI must never leave a Decision without its indexes.
scope: [v2]
---
# Keep exact Decision payloads atomic

The operator supplied this exact body.
`

func TestRecordDecisionRetriesWithoutDuplicateOrRewrite(t *testing.T) {
	root := missionServiceFixture(t)
	first, err := openMissionService(t, root).RecordDecision("-", []byte(transactionDecisionDraft))
	if err != nil {
		t.Fatal(err)
	}
	before := canonicalTreeDigest(t, root)
	retry, err := openMissionService(t, root).RecordDecision("-", []byte(transactionDecisionDraft))
	if err != nil {
		t.Fatal(err)
	}
	if retry.Ref != first.Ref || retry.ID != first.ID || retry.Path != first.Path || len(retry.Changed) != 0 {
		t.Fatalf("retry did not converge: first=%+v retry=%+v", first, retry)
	}
	if after := canonicalTreeDigest(t, root); after != before {
		t.Fatal("retry rewrote canonical Decision or indexes")
	}
}

func TestRecordDecisionRollsBackEveryWriteBoundary(t *testing.T) {
	template := missionServiceFixture(t)
	probeRoot := t.TempDir()
	copyTree(t, template, probeRoot)
	probe, err := openMissionService(t, probeRoot).RecordDecision("-", []byte(transactionDecisionDraft))
	if err != nil {
		t.Fatal(err)
	}
	if len(probe.Changed) < 2 {
		t.Fatalf("decision transaction changed %d files; expected Decision plus indexes", len(probe.Changed))
	}

	for boundary := 0; boundary < len(probe.Changed); boundary++ {
		root := t.TempDir()
		copyTree(t, template, root)
		before := canonicalTreeDigest(t, root)
		service := openMissionService(t, root)
		service.ApplyTransaction = func(transactionRoot, key string, changes []governance.FileChange) error {
			return governance.ApplyTransactionWithFailure(transactionRoot, key, changes, boundary)
		}
		if _, err := service.RecordDecision("-", []byte(transactionDecisionDraft)); err == nil {
			t.Fatalf("boundary %d: injected failure did not surface", boundary)
		}
		if err := governance.RecoverTransactions(root); err != nil {
			t.Fatalf("boundary %d: recovery failed: %v", boundary, err)
		}
		if after := canonicalTreeDigest(t, root); after != before {
			t.Fatalf("boundary %d: failed Decision transaction changed canonical files", boundary)
		}
	}
}
