package runtime

import (
	"testing"
	"time"
)

const (
	proposal = "Proposal:0199b000-0000-7000-8000-000000000001"
	mission  = "Mission:0199b000-0000-7000-8000-000000000002"
	decision = "Decision:0199b000-0000-7000-8000-000000000003"
	fp       = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func fixedNow() time.Time { return time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC) }

func TestPreparationKeepsVerdictsSeparateAndSupportsReslicing(t *testing.T) {
	input := PreparationInput{
		Proposal: BoundSource{Ref: proposal, Fingerprint: fp}, Baseline: "git:abc", DirectionSources: []BoundSource{{Ref: proposal, Fingerprint: fp}},
		Candidates: []CandidateSlice{{Name: "vertical", Outcome: "observable workflow", Evidence: []string{"scenario"}, CancellationState: "usable compiler", Reversibility: "local", StandaloneCoherence: "complete", IntegrationPath: "registry", LearningValue: "runtime proof"}},
		Selected:   "vertical", DesignSufficiency: "needs-evidence", DesignRationale: "false assumption", SliceQuality: "dependency-bound", SliceRationale: "shared interface unresolved", BlockingGaps: []string{"Gap:0199b000-0000-7000-8000-000000000004"},
		CompletionBoundary: []string{"cold resume"}, StopConditions: []string{"baseline drift"}, EvidenceClaims: []string{"claim:resume"}, FreshUntil: "2026-08-11T10:00:00Z",
	}
	receipt, err := CompilePreparation(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Ready || receipt.Next != "resume-preparation" || receipt.DesignSufficiency == receipt.SliceQuality || receipt.Fingerprint == "" {
		t.Fatalf("unexpected blocked receipt %#v", receipt)
	}
	input.DesignSufficiency, input.SliceQuality, input.BlockingGaps = "sufficient", "coherent", nil
	input.DesignRationale, input.SliceRationale = "evidence added", "resliced around dependency"
	ready, err := CompilePreparation(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if !ready.Ready || ready.Next != "owner-activation" || ready.Fingerprint == receipt.Fingerprint {
		t.Fatalf("reslicing did not produce a new ready receipt %#v", ready)
	}
}

func TestAutopilotIsExplicitBoundedAndRuntimeNeutral(t *testing.T) {
	input := AutopilotInput{
		Enabled: true, Mission: BoundSource{Ref: mission, Fingerprint: fp}, Authorization: BoundSource{Ref: decision, Fingerprint: fp}, Outcome: "complete local proof", NonGoals: []string{"release"}, AuthoritativeSources: []BoundSource{{Ref: proposal, Fingerprint: fp}},
		DelegatedDecisionDomain: []string{"reversible implementation"}, AllowedActions: []string{"edit", "test"}, ForbiddenEffects: CanonicalForbiddenEffects(), BudgetUnits: 8, RepairBudget: 2, RequiredChecks: []string{"go test ./..."}, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"authority drift"}, RecoveryPoint: "git:abc", ReturnDestination: "Mission owner",
	}
	charter, err := CompileAutopilot(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if !charter.RuntimeNeutral || charter.AutomaticProviderEffects || charter.Handoff.HostPointerRole == "" || charter.Fingerprint == "" {
		t.Fatalf("unsafe charter %#v", charter)
	}
	input.ForbiddenEffects = input.ForbiddenEffects[1:]
	if _, err := CompileAutopilot(input, fixedNow()); err == nil {
		t.Fatal("accepted charter below the initial effect ceiling")
	}
	input.Enabled = false
	if _, err := CompileAutopilot(input, fixedNow()); err == nil {
		t.Fatal("accepted ambient Autopilot")
	}
}
