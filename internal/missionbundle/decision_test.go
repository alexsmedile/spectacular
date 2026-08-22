package missionbundle

import (
	"strings"
	"testing"
)

func TestRecordDecision_StdinAndAtomicity(t *testing.T) {
	payload := `---
type: DecisionDraft
title: Test Recording Decision Architecture
actor: Alex
actor_role: owner
question: Should we use atomic decision recording?
disposition: accepted
rationale: Atomic transactions guarantee consistent indexes and zero corruption.
scope: [v2]
---
# Test Recording Decision Architecture

Content body for test decision.
`

	draft, _, err := ReadDecisionDraft("-", []byte(payload))
	if err != nil {
		t.Fatalf("ReadDecisionDraft failed: %v", err)
	}
	if draft.Title != "Test Recording Decision Architecture" {
		t.Errorf("expected title 'Test Recording Decision Architecture', got %s", draft.Title)
	}
	if draft.Disposition != "accepted" {
		t.Errorf("expected disposition 'accepted', got %s", draft.Disposition)
	}
}

func TestReadDecisionDraft_ValidationErrors(t *testing.T) {
	invalidType := `---
type: SomethingElse
title: Test
disposition: accepted
---
`
	_, _, err := ReadDecisionDraft("-", []byte(invalidType))
	if err == nil || !strings.Contains(err.Error(), "must declare type: DecisionDraft") {
		t.Fatalf("expected error for invalid type, got: %v", err)
	}

	missingTitle := `---
type: DecisionDraft
disposition: accepted
---
`
	_, _, err = ReadDecisionDraft("-", []byte(missingTitle))
	if err == nil || !strings.Contains(err.Error(), "title is required") {
		t.Fatalf("expected error for missing title, got: %v", err)
	}

	missingDisposition := `---
type: DecisionDraft
title: Test
---
`
	_, _, err = ReadDecisionDraft("-", []byte(missingDisposition))
	if err == nil || !strings.Contains(err.Error(), "disposition is required") {
		t.Fatalf("expected error for missing disposition, got: %v", err)
	}
}
