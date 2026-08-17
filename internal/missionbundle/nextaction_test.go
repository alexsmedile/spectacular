package missionbundle

import (
	"strings"
	"testing"
)

// reviewedBundle builds a Mission whose Objectives are all implemented, so the
// next action turns entirely on what the reviews say.
func reviewedBundle(fingerprint string, pointers []ReviewPointer) *Bundle {
	return &Bundle{
		Ref:    "M_REVIEW",
		Status: "active",
		Objectives: []Objective{
			{Ref: "O1", Status: "implemented"},
			{Ref: "O2", Status: "implemented"},
		},
		Run:          &Run{Ref: "R1", Status: "active"},
		Activation:   &Activation{Fingerprint: fingerprint},
		Reviews:      pointers,
		RepairBudget: 3,
	}
}

func passingPointerBoundTo(fingerprint string) ReviewPointer {
	return ReviewPointer{
		Ref:     "RV1",
		Verdict: "pass",
		Document: &Review{
			Ref:      "RV1",
			Reviewed: Reviewed{ActivationFingerprint: fingerprint},
		},
	}
}

// TestNextActionReadsRecordedReviews covers the four states the next action must
// tell apart. Before this, the "record a review" branch fired on implemented
// Objectives alone and never read b.Reviews, so a recorded review could not
// retire the instruction and the line disagreed with the record it summarized.
func TestNextActionReadsRecordedReviews(t *testing.T) {
	const current = "sha256:current"
	const stale = "sha256:stale"

	failing := passingPointerBoundTo(current)
	failing.Verdict = "fail"
	failing.Document.Claims = []ClaimVerdict{{Claim: "c", Verdict: "fail"}}

	unresolved := ReviewPointer{Ref: "RV1", Verdict: "pass"}

	for _, test := range []struct {
		name       string
		bundle     *Bundle
		wantNext   string
		wantHolder string
	}{
		{
			name:       "no review asks for one",
			bundle:     reviewedBundle(current, nil),
			wantNext:   "record a review",
			wantHolder: "operator",
		},
		{
			name:       "passing review on the current fingerprint hands off to the owner",
			bundle:     reviewedBundle(current, []ReviewPointer{passingPointerBoundTo(current)}),
			wantNext:   "the owner completes the Mission",
			wantHolder: "owner",
		},
		{
			name:       "passing review on a stale fingerprint keeps asking",
			bundle:     reviewedBundle(current, []ReviewPointer{passingPointerBoundTo(stale)}),
			wantNext:   "record a review",
			wantHolder: "operator",
		},
		{
			name:       "failing review keeps asking",
			bundle:     reviewedBundle(current, []ReviewPointer{failing}),
			wantNext:   "record a review",
			wantHolder: "operator",
		},
		{
			name:       "unresolved pointer cannot prove its binding",
			bundle:     reviewedBundle(current, []ReviewPointer{unresolved}),
			wantNext:   "record a review",
			wantHolder: "operator",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			state := test.bundle.Derive()
			if !strings.Contains(state.Next, test.wantNext) {
				t.Fatalf("Next = %q, want it to contain %q", state.Next, test.wantNext)
			}
			if state.Holder != test.wantHolder {
				t.Fatalf("Holder = %q, want %q", state.Holder, test.wantHolder)
			}
		})
	}
}

// TestNextActionPrefersEarlierLifecycleGates asserts a recorded review does not
// let the next action skip a blocking condition ahead of it. Repair exhaustion
// is the owner's decision and outranks completion.
func TestNextActionPrefersEarlierLifecycleGates(t *testing.T) {
	bundle := reviewedBundle("sha256:current", []ReviewPointer{passingPointerBoundTo("sha256:current")})
	bundle.Run.Repairs = 3

	state := bundle.Derive()
	if !strings.Contains(state.Next, "repair budget is exhausted") {
		t.Fatalf("Next = %q, want the exhaustion gate to win", state.Next)
	}
	if state.Holder != "operator" {
		t.Fatalf("Holder = %q, want operator; exhaustion is not the completion handoff", state.Holder)
	}
}

// TestNextActionWithoutActivationJudgesOnVerdictAlone covers a Mission with no
// activation: there is no boundary for a review to be stale against.
func TestNextActionWithoutActivationJudgesOnVerdictAlone(t *testing.T) {
	bundle := reviewedBundle("", []ReviewPointer{passingPointerBoundTo("sha256:anything")})
	bundle.Activation = nil

	state := bundle.Derive()
	if !strings.Contains(state.Next, "the owner completes the Mission") {
		t.Fatalf("Next = %q, want the completion handoff", state.Next)
	}
}
