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
		CompletionCriteria: []CompletionCriterion{{Claim: "claim:resume", PassBoundary: "cold resume succeeds", ProofRequirement: "real-process scenario", ReviewLevel: ReviewIndependent}}, StopConditions: []string{"baseline drift"}, EvidenceClaims: []string{"claim:resume"}, FreshUntil: "2026-08-11T10:00:00Z",
	}
	receipt, err := CompilePreparation(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Ready || receipt.Next != "adaptive-grill" || len(receipt.UnmetRequirements) != 3 || receipt.DesignSufficiency == receipt.SliceQuality || receipt.Fingerprint == "" {
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

func TestPreparationReportsOnlyUnmetCompletionRequirements(t *testing.T) {
	input := PreparationInput{
		Proposal: BoundSource{Ref: proposal, Fingerprint: fp}, Baseline: "git:abc", DirectionSources: []BoundSource{{Ref: proposal, Fingerprint: fp}},
		Candidates: []CandidateSlice{{Name: "vertical", Outcome: "observable workflow", Evidence: []string{"scenario"}, CancellationState: "usable compiler", Reversibility: "local", StandaloneCoherence: "complete", IntegrationPath: "registry", LearningValue: "runtime proof"}},
		Selected:   "vertical", DesignSufficiency: "sufficient", DesignRationale: "clear", SliceQuality: "coherent", SliceRationale: "bounded",
		CompletionCriteria: []CompletionCriterion{{Claim: "claim:one", PassBoundary: "", ProofRequirement: "test", ReviewLevel: "harsh"}},
		StopConditions:     []string{"authority drift"}, EvidenceClaims: []string{"claim:one", "claim:two"}, FreshUntil: "2026-08-11T10:00:00Z",
	}
	receipt, err := CompilePreparation(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Ready || receipt.Next != "adaptive-grill" || len(receipt.UnmetRequirements) != 3 {
		t.Fatalf("unexpected readiness diagnostics %#v", receipt.UnmetRequirements)
	}
	input.CompletionCriteria = []CompletionCriterion{
		{Claim: "claim:two", PassBoundary: "two passes", ProofRequirement: "test two", ReviewLevel: ReviewClustered},
		{Claim: "claim:one", PassBoundary: "one passes", ProofRequirement: "test one", ReviewLevel: ReviewAutomatic},
	}
	ready, err := CompilePreparation(input, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if !ready.Ready || len(ready.UnmetRequirements) != 0 || ready.CompletionCriteria[0].Claim != "claim:one" {
		t.Fatalf("criteria were not ready and canonicalized %#v", ready)
	}
	reversed := []CompletionCriterion{ready.CompletionCriteria[1], ready.CompletionCriteria[0]}
	if CompletionFingerprint(reversed) != CompletionFingerprint(ready.CompletionCriteria) {
		t.Fatal("completion fingerprint depends on input order")
	}
}

func TestAutopilotIsExplicitBoundedAndRuntimeNeutral(t *testing.T) {
	input := AutopilotInput{
		Enabled: true, Mission: BoundSource{Ref: mission, Fingerprint: fp}, Authorization: BoundSource{Ref: decision, Fingerprint: fp}, Outcome: "complete local proof", NonGoals: []string{"release"}, AuthoritativeSources: []BoundSource{{Ref: proposal, Fingerprint: fp}},
		DelegatedDecisionDomain: []string{"reversible implementation"}, AllowedActions: []string{"edit", "test"}, ForbiddenEffects: CanonicalForbiddenEffects(), BudgetUnits: 8, RepairBudget: 2, RequiredChecks: []string{"go test ./..."}, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"authority drift"}, RecoveryPoint: "git:abc", ReturnDestination: "Mission owner",
	}
	limits := AutopilotLimits{
		Mission: input.Mission, Authorization: input.Authorization, Outcome: input.Outcome,
		AllowedActions: []string{"edit", "test"}, ForbiddenEffects: CanonicalForbiddenEffects(),
		BudgetUnits: 8, RepairBudget: 2, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"authority drift"},
	}
	charter, err := CompileAutopilot(input, limits, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if !charter.RuntimeNeutral || charter.AutomaticProviderEffects || charter.Handoff.HostPointerRole == "" || charter.Fingerprint == "" {
		t.Fatalf("unsafe charter %#v", charter)
	}
	input.ForbiddenEffects = input.ForbiddenEffects[1:]
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("accepted charter below the initial effect ceiling")
	}
	input.ForbiddenEffects = CanonicalForbiddenEffects()
	input.AllowedActions = append(input.AllowedActions, "merge")
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("accepted action outside Mission authority and inside forbidden effects")
	}
	input.AllowedActions = []string{"edit", "test"}
	input.BudgetUnits = 9
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("accepted budget above Mission authority")
	}
	input.BudgetUnits = 8
	input.Enabled = false
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("accepted ambient Autopilot")
	}
}

func TestAutopilotResourceLimitsStateTruthfulEnforcement(t *testing.T) {
	input := AutopilotInput{
		Enabled: true, Mission: BoundSource{Ref: mission, Fingerprint: fp}, Authorization: BoundSource{Ref: decision, Fingerprint: fp}, Outcome: "complete local proof", NonGoals: []string{"release"}, AuthoritativeSources: []BoundSource{{Ref: proposal, Fingerprint: fp}},
		DelegatedDecisionDomain: []string{"reversible implementation"}, AllowedActions: []string{"test"}, ForbiddenEffects: CanonicalForbiddenEffects(), BudgetUnits: 8, RepairBudget: 2, RequiredChecks: []string{"go test ./..."}, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"authority drift"}, RecoveryPoint: "git:abc", ReturnDestination: "Mission owner",
		ResourceLimits: []ResourceLimit{
			{Resource: "wall-time-seconds", Maximum: 3600, Enforcement: EnforcementHard, MeasureCapability: "host.wall.measure", CancelCapability: "host.wall.cancel"},
			{Resource: "tokens", Maximum: 100000, Enforcement: EnforcementObserved, MeasureCapability: "host.tokens.measure"},
			{Resource: "spend-cents", Maximum: 5000, Enforcement: EnforcementUnsupported},
		},
	}
	limits := AutopilotLimits{Mission: input.Mission, Authorization: input.Authorization, Outcome: input.Outcome, AllowedActions: []string{"test"}, ForbiddenEffects: CanonicalForbiddenEffects(), BudgetUnits: 8, RepairBudget: 2, HostCapabilities: []string{"host.wall.measure", "host.wall.cancel", "host.tokens.measure"}, ExpiresAt: input.ExpiresAt, StopConditions: input.StopConditions}
	charter, err := CompileAutopilot(input, limits, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(charter.ResourceLimits) != 3 || charter.ResourceLimits[0].Resource != "spend-cents" {
		t.Fatalf("resource limits were not canonicalized: %#v", charter.ResourceLimits)
	}
	input.ResourceLimits[0].CancelCapability = "host.wall.missing"
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("accepted a false hard limit without host cancellation")
	}
	input.ResourceLimits[0].CancelCapability = "host.wall.cancel"
	limits.HostCapabilities = nil
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("request was allowed to attest to its own host capabilities")
	}
	limits.HostCapabilities = []string{"host.wall.measure", "host.wall.cancel", "host.tokens.measure"}
	input.ResourceLimits[1].CancelCapability = "host.tokens.cancel"
	if _, err := CompileAutopilot(input, limits, fixedNow()); err == nil {
		t.Fatal("observed limit claimed cancellation")
	}
}
