package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/projection"
	spectacularruntime "github.com/alexsmedile/spectacular/v2/internal/runtime"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

const (
	contractRef   = "Contract:0199b000-0000-7000-8000-000000000001"
	proposalRef   = "Proposal:0199b000-0000-7000-8000-000000000002"
	missionRef    = "Mission:0199b000-0000-7000-8000-000000000003"
	objectiveRef  = "Objective:0199b000-0000-7000-8000-000000000004"
	runRef        = "Run:0199b000-0000-7000-8000-000000000005"
	handoffRef    = "Handoff:0199b000-0000-7000-8000-000000000006"
	returnRef     = "Handoff:0199b000-0000-7000-8000-000000000007"
	evidenceRef   = "Evidence:0199b000-0000-7000-8000-000000000008"
	assessmentRef = "Assessment:0199b000-0000-7000-8000-000000000009"
)

var testNow = time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC)

func TestProviderNeutralGovernedLoopAndSecondColdResume(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	contract := lookup(t, svc, contractRef, domain.Contract)

	proposalDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000010", "proposal.create", proposalRef, "absent", nil, "approve")
	svc = openService(t, root)
	proposalInput := ProposalInput{
		ID: strings.TrimPrefix(proposalRef, "Proposal:"), Title: "Add governed closure", Actor: "owner", Status: "accepted",
		TargetContract: contractRef, BaseVersion: "1", BaseFingerprint: contract.Fingerprint,
		Additions: []string{"Archive only after terminal continuity is complete."}, Rationale: "Close the governed loop.", Scope: []string{"v2"},
		Authorization: proposalDecision, IdempotencyKey: "proposal-create-1",
	}
	expiredService := svc
	expiredService.Now = func() time.Time { return testNow.Add(48 * time.Hour) }
	if _, err := expiredService.CreateProposal(proposalInput); refusalCode(err) != domain.RefusalExpiredAuthority {
		t.Fatalf("expired authority err=%v", err)
	}
	proposalResult, err := svc.CreateProposal(proposalInput)
	must(t, err)
	svc = openService(t, root)
	view, err := svc.ProposalView(proposalRef)
	must(t, err)
	if view["candidate_authoritative"] != false || !strings.Contains(strings.Join(view["candidate_contract"].(ContractCandidate).RequiredBehavior, "\n"), "Archive only") {
		t.Fatalf("Proposal candidate/delta not rendered honestly: %#v", view)
	}
	replayService := svc
	replayService.Now = func() time.Time { return testNow.Add(time.Hour) }
	replay, err := replayService.CreateProposal(proposalInput)
	must(t, err)
	if !replay.IdempotentReplay || replay.Fingerprint != proposalResult.Fingerprint {
		t.Fatalf("identity replay = %#v", replay)
	}
	stale := proposalInput
	stale.ID = "0199b000-0000-7000-8000-000000000011"
	stale.BaseFingerprint = strings.Repeat("0", 64)
	if _, err := svc.CreateProposal(stale); refusalCode(err) != domain.RefusalStaleFingerprint {
		t.Fatalf("stale Proposal base err=%v", err)
	}

	missionDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000012", "mission.create", missionRef, "absent", nil, "activate")
	svc = openService(t, root)
	proposal := lookup(t, svc, proposalRef, domain.Proposal)
	missionInput := MissionInput{
		ID: strings.TrimPrefix(missionRef, "Mission:"), Title: "Govern closure", Actor: "owner", Proposal: proposalRef, Outcome: "Reconcile and recover.",
		Objectives: []ObjectiveInput{{ID: strings.TrimPrefix(objectiveRef, "Objective:"), Outcome: "Prove the loop", ExpectedProof: []string{"claim:closure"}}}, InitialRunID: strings.TrimPrefix(runRef, "Run:"),
		DesignSufficiency: "sufficient", SliceQuality: "coherent", EvidenceClaims: []string{"claim:closure"}, Scope: []string{"v2"}, AllowedActions: []string{"test", "write-v2"}, ForbiddenEffects: []string{"provider-mutation"},
		Baseline: contract.Fingerprint, BudgetUnits: 2, RepairBudget: 2, ExpiresAt: "2026-08-11T10:00:00Z", Stops: []string{"authority-drift"}, RecoveryPoint: "git-head", ReturnDestination: "central", Authorization: missionDecision, ExpectedProposalFingerprint: proposal.Fingerprint, IdempotencyKey: "mission-create-1",
	}
	preparation, err := spectacularruntime.CompilePreparation(spectacularruntime.PreparationInput{
		Proposal: spectacularruntime.BoundSource{Ref: proposalRef, Fingerprint: proposal.Fingerprint}, Baseline: contract.Fingerprint,
		DirectionSources: []spectacularruntime.BoundSource{{Ref: contractRef, Fingerprint: contract.Fingerprint}},
		Candidates:       []spectacularruntime.CandidateSlice{{Name: "governed-loop", Outcome: "Reconcile and recover.", Evidence: []string{"claim:closure"}, CancellationState: "Proposal remains inspectable", Reversibility: "local fixture", StandaloneCoherence: "complete B+C loop", IntegrationPath: "governance service", LearningValue: "cold recovery proof"}},
		Selected:         "governed-loop", DesignSufficiency: "sufficient", DesignRationale: "accepted authority and closure contracts", SliceQuality: "coherent", SliceRationale: "serial governed loop",
		CompletionBoundary: []string{"claim:closure"}, StopConditions: []string{"authority-drift"}, EvidenceClaims: []string{"claim:closure"}, FreshUntil: "2026-08-10T11:00:00Z",
	}, testNow)
	must(t, err)
	missionInput.Preparation = &preparation
	tampered := preparation
	tampered.Selected = "tampered"
	missionInput.Preparation = &tampered
	if _, err := svc.CreateMission(missionInput); refusalCode(err) != domain.RefusalInvalidKnownField {
		t.Fatalf("tampered preparation receipt err=%v", err)
	}
	missionInput.Preparation = &preparation
	_, err = svc.CreateMission(missionInput)
	must(t, err)

	mission := lookup(t, openService(t, root), missionRef, domain.Mission)
	activateDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000013", "mission.transition.active", missionRef, mission.Fingerprint, nil, "activate")
	svc = openService(t, root)
	stalePreparationService := svc
	stalePreparationService.Now = func() time.Time { return testNow.Add(2 * time.Hour) }
	if _, err := stalePreparationService.TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activateDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-active-stale"}); refusalCode(err) != domain.RefusalExpiredAuthority {
		t.Fatalf("stale preparation did not block activation: %v", err)
	}
	active, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activateDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-active-1"})
	must(t, err)
	if replay, err := openService(t, root).TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activateDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-active-1"}); err != nil || !replay.IdempotentReplay {
		t.Fatalf("transition replay=%#v err=%v", replay, err)
	}

	handoffDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000014", "handoff.create", handoffRef, "absent", nil, "dispatch")
	svc = openService(t, root)
	handoffInput := HandoffInput{
		ID: strings.TrimPrefix(handoffRef, "Handoff:"), Title: "Bounded fake-provider work", Mission: missionRef, Objective: objectiveRef, Run: runRef, Sender: "owner", Actor: "executor", Destination: "replacement-runtime", HostPointer: "host-task:disposable",
		Scope: []string{"v2"}, Inputs: []string{contractRef + "@" + contract.Fingerprint}, AllowedActions: []string{"test"}, ForbiddenEffects: []string{"provider-mutation"}, EvidenceClaims: []string{"claim:closure"}, BudgetUnits: 1, ExpiresAt: "2026-08-11T10:00:00Z", Stops: []string{"authority-drift"}, RecoveryPoint: "git-head", ReturnDestination: "central", Authorization: handoffDecision, ExpectedMissionFingerprint: active.Fingerprint, IdempotencyKey: "handoff-create-1",
	}
	outOfScope := handoffInput
	outOfScope.Scope = []string{"v1"}
	if _, err := svc.CreateHandoff(outOfScope); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-envelope Handoff err=%v", err)
	}
	outOfClaims := handoffInput
	outOfClaims.EvidenceClaims = []string{"claim:undeclared"}
	if _, err := svc.CreateHandoff(outOfClaims); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-envelope Handoff claims err=%v", err)
	}
	outOfStops := handoffInput
	outOfStops.Stops = []string{"undeclared-stop"}
	if _, err := svc.CreateHandoff(outOfStops); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-envelope Handoff stops err=%v", err)
	}
	overBudgetHandoff := handoffInput
	overBudgetHandoff.BudgetUnits = 3
	if _, err := svc.CreateHandoff(overBudgetHandoff); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-envelope Handoff budget err=%v", err)
	}
	staleInput := handoffInput
	staleInput.Inputs = []string{contractRef + "@" + strings.Repeat("0", 64)}
	if _, err := svc.CreateHandoff(staleInput); refusalCode(err) != domain.RefusalStaleFingerprint {
		t.Fatalf("stale authoritative Handoff input err=%v", err)
	}
	handoff, err := svc.CreateHandoff(handoffInput)
	must(t, err)
	svc = openService(t, root)
	validated, err := svc.ValidateHandoff(handoffRef)
	must(t, err)
	if validated["valid"] != true || len(validated["does_not_prove"].([]string)) == 0 {
		t.Fatalf("handoff validation overclaims: %#v", validated)
	}
	lateValidation := svc
	lateValidation.Now = func() time.Time { return testNow.Add(48 * time.Hour) }
	if _, err := lateValidation.ValidateHandoff(handoffRef); refusalCode(err) != domain.RefusalExpiredAuthority {
		t.Fatalf("expired Mission/Handoff envelope err=%v", err)
	}
	invalidReturn := HandoffReturnInput{ID: strings.TrimPrefix(returnRef, "Handoff:"), Title: "Over-broad return", Dispatch: handoffRef, Status: "succeeded", Actor: "replacement-runtime", FinalBaseline: contract.Fingerprint, Result: "over-broad", Actions: []string{"write-v2"}, BudgetUsed: 1, RecoveryPoint: "git-head", NextAction: "record evidence", ExpectedDispatchFingerprint: handoff.Fingerprint, IdempotencyKey: "handoff-return-1"}
	if _, err := svc.ReturnHandoff(invalidReturn); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("over-broad Handoff return actions err=%v", err)
	}
	_, err = svc.ReturnHandoff(HandoffReturnInput{
		ID: strings.TrimPrefix(returnRef, "Handoff:"), Title: "Returned bounded work", Dispatch: handoffRef, Status: "succeeded", Actor: "replacement-runtime", FinalBaseline: contract.Fingerprint, Result: "fake receipt only", Actions: []string{"test"}, ProviderReceipts: []string{"fake-provider:ok"}, Evidence: []string{evidenceRef}, BudgetUsed: 1, RecoveryPoint: "git-head", NextAction: "record evidence", ExpectedDispatchFingerprint: handoff.Fingerprint, IdempotencyKey: "handoff-return-1",
	})
	must(t, err)
	supersedingRef := "Handoff:0199b000-0000-7000-8000-000000000030"
	supersedingDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000031", "handoff.create", supersedingRef, "absent", nil, "redispatch")
	svc = openService(t, root)
	superseding := handoffInput
	superseding.ID = strings.TrimPrefix(supersedingRef, "Handoff:")
	superseding.Authorization = supersedingDecision
	superseding.IdempotencyKey = "handoff-create-2"
	superseding.Supersedes = handoffRef
	_, err = svc.CreateHandoff(superseding)
	must(t, err)
	if _, err := openService(t, root).ValidateHandoff(handoffRef); refusalCode(err) != domain.RefusalConflictingAuthority {
		t.Fatalf("superseded Handoff remained valid: %v", err)
	}

	evidenceDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000015", "evidence.create", evidenceRef, "absent", nil, "record")
	svc = openService(t, root)
	_, err = svc.CreateEvidence(EvidenceInput{
		ID: strings.TrimPrefix(evidenceRef, "Evidence:"), Title: "Closure proof", Mission: missionRef, Objective: objectiveRef, Claim: "claim:closure", Classification: "direct", Scope: []string{"v2"}, Method: "deterministic tests", Actor: "executor", Target: proposalRef, Environment: "disposable", ObservedAt: "2026-08-10T10:00:00Z", FreshnessValidUntil: "2026-08-11T10:00:00Z", RequiredChecks: []string{"go-test"}, CheckResults: []string{"pass:go-test"}, ReviewState: "independent-accepted", ExecutorAuthored: true, Authorization: evidenceDecision, IdempotencyKey: "evidence-create-1",
	})
	must(t, err)
	badEvidenceRef := "Evidence:0199b000-0000-7000-8000-000000000032"
	badEvidenceDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000033", "evidence.create", badEvidenceRef, "absent", nil, "record")
	svc = openService(t, root)
	_, err = svc.CreateEvidence(EvidenceInput{ID: strings.TrimPrefix(badEvidenceRef, "Evidence:"), Title: "Unreviewed executor observation", Mission: missionRef, Objective: objectiveRef, Claim: "claim:other", Classification: "observation", Scope: []string{"v2"}, Method: "executor report", Actor: "executor", Target: proposalRef, Environment: "disposable", ObservedAt: "2026-08-10T10:00:00Z", FreshnessValidUntil: "2026-08-11T10:00:00Z", RequiredChecks: []string{"go-test"}, CheckResults: []string{"pass:go-test"}, ReviewState: "unreviewed", ExecutorAuthored: true, Authorization: badEvidenceDecision, IdempotencyKey: "evidence-bad-1"})
	must(t, err)

	mission = lookup(t, openService(t, root), missionRef, domain.Mission)
	awaitDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000016", "mission.transition.awaiting-assessment", missionRef, mission.Fingerprint, nil, "await")
	svc = openService(t, root)
	awaiting, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "awaiting-assessment", Authorization: awaitDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-await-1"})
	must(t, err)

	assessmentDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000017", "assessment.record", assessmentRef, "absent", nil, "record")
	svc = openService(t, root)
	badAssessment := AssessmentInput{ID: strings.TrimPrefix(assessmentRef, "Assessment:"), Title: "Premature assessment", Mission: missionRef, Verdict: "ready-for-owner", Actor: "reviewer", Claims: []string{"claim:closure"}, Evidence: []string{badEvidenceRef}, RecoveryPoint: "git-head", Authorization: assessmentDecision, IdempotencyKey: "assessment-1"}
	if _, err := svc.RecordAssessment(badAssessment); refusalCode(err) != domain.RefusalInsufficientEvidence {
		t.Fatalf("executor-dependent Evidence did not block: %v", err)
	}
	prematureResolve := createDecision(t, root, "0199b000-0000-7000-8000-000000000034", "mission.transition.resolved", missionRef, awaiting.Fingerprint, nil, "completed")
	svc = openService(t, root)
	if _, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: prematureResolve, ExpectedFingerprint: awaiting.Fingerprint, IdempotencyKey: "premature-resolve", Disposition: "completed", TerminalNextAction: "continue"}); refusalCode(err) != domain.RefusalInsufficientEvidence {
		t.Fatalf("resolution without Assessment did not refuse: %v", err)
	}
	svc = openService(t, root)
	_, err = svc.RecordAssessment(AssessmentInput{ID: strings.TrimPrefix(assessmentRef, "Assessment:"), Title: "Owner readiness assessment", Mission: missionRef, Verdict: "ready-for-owner", Actor: "reviewer", Claims: []string{"claim:closure"}, Evidence: []string{evidenceRef}, RecoveryPoint: "git-head", Authorization: assessmentDecision, IdempotencyKey: "assessment-1"})
	must(t, err)

	reconcileDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000018", "contract.reconcile", contractRef, contract.Fingerprint, []string{assessmentRef}, "accept")
	svc = openService(t, root)
	reconciled, err := svc.Reconcile(ReconcileInput{Contract: contractRef, Proposal: proposalRef, Authorization: reconcileDecision, ExpectedFingerprint: contract.Fingerprint, IdempotencyKey: "reconcile-1"})
	must(t, err)
	if !strings.HasPrefix(reconciled.Receipt, "Evidence:") {
		t.Fatalf("durable receipt missing: %#v", reconciled)
	}
	svc = openService(t, root)
	contractView, err := svc.ContractView(contractRef)
	must(t, err)
	if contractView["version"] != "2" || contractView["accepted_proposal"] != proposalRef {
		t.Fatalf("reconciled Contract = %#v", contractView)
	}
	if replay, err := svc.Reconcile(ReconcileInput{Contract: contractRef, Proposal: proposalRef, Authorization: reconcileDecision, ExpectedFingerprint: contract.Fingerprint, IdempotencyKey: "reconcile-1"}); err != nil || !replay.IdempotentReplay {
		t.Fatalf("reconcile replay=%#v err=%v", replay, err)
	}
	for name, changedReplay := range map[string]ReconcileInput{
		"nonexistent Proposal":  {Contract: contractRef, Proposal: "Proposal:0199b000-0000-7000-8000-000000000099", Authorization: reconcileDecision, ExpectedFingerprint: contract.Fingerprint, IdempotencyKey: "reconcile-1"},
		"nonexistent Decision":  {Contract: contractRef, Proposal: proposalRef, Authorization: "Decision:0199b000-0000-7000-8000-000000000099", ExpectedFingerprint: contract.Fingerprint, IdempotencyKey: "reconcile-1"},
		"unrelated fingerprint": {Contract: contractRef, Proposal: proposalRef, Authorization: reconcileDecision, ExpectedFingerprint: strings.Repeat("f", 64), IdempotencyKey: "reconcile-1"},
	} {
		if _, err := svc.Reconcile(changedReplay); refusalCode(err) != domain.RefusalIdempotencyConflict {
			t.Fatalf("%s reconciliation replay err=%v", name, err)
		}
	}
	unboundResolution := createDecision(t, root, "0199b000-0000-7000-8000-000000000035", "mission.transition.resolved", missionRef, awaiting.Fingerprint, []string{assessmentRef}, "completed", "objective.satisfy:"+objectiveRef, "terminal-next-action:review next governed request")
	svc = openService(t, root)
	if _, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: unboundResolution, ExpectedFingerprint: awaiting.Fingerprint, IdempotencyKey: "unbound-resolution", Disposition: "completed", Assessment: assessmentRef, Reconciliation: reconciled.Receipt, TerminalNextAction: "review next governed request", SatisfiedObjectives: []string{objectiveRef}}); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("resolution Decision omitted exact receipt but succeeded: %v", err)
	}

	resolveDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000019", "mission.transition.resolved", missionRef, awaiting.Fingerprint, []string{assessmentRef, reconciled.Receipt}, "completed", "objective.satisfy:"+objectiveRef, "terminal-next-action:review next governed request")
	svc = openService(t, root)
	resolved, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: resolveDecision, ExpectedFingerprint: awaiting.Fingerprint, IdempotencyKey: "mission-resolve-1", Disposition: "completed", Assessment: assessmentRef, Reconciliation: reconciled.Receipt, TerminalNextAction: "review next governed request", SatisfiedObjectives: []string{objectiveRef}})
	must(t, err)

	archiveDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000020", "mission.archive", missionRef, resolved.Fingerprint, []string{assessmentRef, reconciled.Receipt}, "archive")
	svc = openService(t, root)
	archiveInput := ArchiveInput{Mission: missionRef, Authorization: archiveDecision, ExpectedFingerprint: resolved.Fingerprint, IdempotencyKey: "mission-archive-1", TerminalPacket: missionRef}
	_, err = svc.ArchiveMission(archiveInput)
	must(t, err)
	if replay, err := openService(t, root).ArchiveMission(archiveInput); err != nil || !replay.IdempotentReplay {
		t.Fatalf("archive replay=%#v err=%v", replay, err)
	}

	// This is the second cold runtime: only cwd/canonical records are reopened.
	cold := openService(t, root)
	builder := projection.Builder{Workspace: cold.Workspace, Now: func() time.Time { return testNow }}
	project, err := builder.Project()
	must(t, err)
	if len(project.Authoritative.CurrentTruth) != 1 || project.Authoritative.CurrentTruth[0].Ref != contractRef {
		t.Fatalf("cold current truth = %#v", project.Authoritative.CurrentTruth)
	}
	if len(project.Projection.Missions) != 1 || project.Projection.Continuation == nil || project.Projection.Continuation.Operation != "review next governed request" {
		t.Fatalf("cold terminal continuity = %#v", project.Projection)
	}
	if project.Projection.Continuation.AuthorizedBy.Ref != resolveDecision {
		t.Fatalf("archive Decision incorrectly authorized terminal continuation: %#v", project.Projection.Continuation)
	}
	card, err := builder.Mission(missionRef)
	must(t, err)
	if len(card.Pointers) < 4 { // Assessment, reconciliation receipt, resolution Decision, archive Decision.
		t.Fatalf("historical closure provenance missing: %#v", card.Pointers)
	}
}

func TestTransactionInterruptionRollsBackEveryTarget(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	must(t, os.WriteFile(first, []byte("old-first"), 0o644))
	must(t, os.WriteFile(second, []byte("old-second"), 0o644))
	err := applyTransaction(root, "interrupt", []FileChange{{Path: "first", Data: []byte("new-first")}, {Path: "second", Data: []byte("new-second")}}, 1)
	if err == nil {
		t.Fatal("injected interruption succeeded")
	}
	for path, want := range map[string]string{first: "old-first", second: "old-second"} {
		got, readErr := os.ReadFile(path)
		must(t, readErr)
		if string(got) != want {
			t.Fatalf("rollback %s = %q, want %q", path, got, want)
		}
	}
}

func TestTransactionRejectsSymlinkedParentEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	must(t, os.Symlink(outside, filepath.Join(root, "linked")))
	err := ApplyTransaction(root, "symlink-parent", []FileChange{{Path: filepath.Join("linked", "escaped"), Data: []byte("escaped"), Mode: 0o644}})
	if refusalCode(err) != domain.RefusalPathEscape {
		t.Fatalf("symlink parent err=%v", err)
	}
	if _, statErr := os.Stat(filepath.Join(outside, "escaped")); !os.IsNotExist(statErr) {
		t.Fatalf("transaction wrote through symlinked parent: %v", statErr)
	}
}

func TestTransactionInstallResistsValidatedParentSwap(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	parent := filepath.Join(root, "records")
	displaced := filepath.Join(root, "records-original")
	must(t, os.Mkdir(parent, 0o755))
	must(t, os.WriteFile(filepath.Join(parent, "contract.md"), []byte("original"), 0o644))

	err := applyTransactionWithInstallHook(root, "parent-swap", []FileChange{{
		Path: filepath.Join("records", "contract.md"),
		Data: []byte("candidate"),
		Mode: 0o644,
	}}, -1, func(index int) {
		if index != 0 {
			t.Fatalf("unexpected install index %d", index)
		}
		must(t, os.Rename(parent, displaced))
		must(t, os.Symlink(outside, parent))
	})
	if err == nil {
		t.Fatal("parent swap unexpectedly installed outside the rooted workspace")
	}
	if refusalCode(err) != domain.RefusalPathEscape {
		t.Fatalf("parent swap refusal=%v", err)
	}
	if _, statErr := os.Stat(filepath.Join(outside, "contract.md")); !os.IsNotExist(statErr) {
		t.Fatalf("outside target changed: %v", statErr)
	}
	data, readErr := os.ReadFile(filepath.Join(displaced, "contract.md"))
	must(t, readErr)
	if string(data) != "original" {
		t.Fatalf("displaced original=%q", data)
	}
}

func TestRecoveryRejectsCraftedJournalOutsideWorkspace(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "victim")
	must(t, os.WriteFile(outside, []byte("original"), 0o644))
	txRoot := filepath.Join(root, ".spectacular", "transactions")
	must(t, os.MkdirAll(txRoot, 0o755))
	backup := filepath.Join(txRoot, "crafted.old")
	temporary := filepath.Join(txRoot, "crafted.new")
	must(t, os.WriteFile(backup, []byte("attacker"), 0o600))
	must(t, os.WriteFile(temporary, []byte("attacker"), 0o600))
	journal := transactionJournal{Schema: "spectacular.transaction.v1", Key: "crafted", Files: []transactionFile{{Target: outside, Temporary: temporary, Backup: backup, HadOriginal: true, Mode: 0o644}}}
	data, err := json.Marshal(journal)
	must(t, err)
	journalHash := sha256.Sum256([]byte("crafted"))
	must(t, os.WriteFile(filepath.Join(txRoot, hex.EncodeToString(journalHash[:16])+".json"), data, 0o600))
	if err := RecoverTransactions(root); refusalCode(err) != domain.RefusalPathEscape {
		t.Fatalf("crafted journal err=%v", err)
	}
	got, err := os.ReadFile(outside)
	must(t, err)
	if string(got) != "original" {
		t.Fatalf("crafted journal overwrote outside target: %q", got)
	}
}

func TestTransactionPreparationFailureCleansArtifacts(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "original")
	must(t, os.WriteFile(target, []byte("original"), 0o644))
	err := ApplyTransaction(root, "prepare-cleanup", []FileChange{
		{Path: "original", Data: []byte("changed"), Mode: 0o644},
		{Path: filepath.Join("..", "escape"), Data: []byte("escape"), Mode: 0o644},
	})
	if refusalCode(err) != domain.RefusalInvalidWorkspacePath {
		t.Fatalf("preparation err=%v", err)
	}
	got, err := os.ReadFile(target)
	must(t, err)
	if string(got) != "original" {
		t.Fatalf("preparation changed original: %q", got)
	}
	entries, err := os.ReadDir(filepath.Join(root, ".spectacular", "transactions"))
	must(t, err)
	if len(entries) != 0 {
		t.Fatalf("pre-journal artifacts remain: %v", entries)
	}
}

func TestMultiContractReconciliationRefusesBeforeAnyPartialWrite(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	first := lookup(t, svc, contractRef, domain.Contract)
	secondRef := "Contract:0199b000-0000-7000-8000-000000000040"
	secondID, _ := domain.ParseID(strings.TrimPrefix(secondRef, "Contract:"))
	second := svc.document(domain.Contract, secondID, "Second Capability Contract", "owner", "current")
	setContract(second, candidateFromDocument(first.Document), "1")
	workspace.SetString(second, "accepted_proposal", "Proposal:0199b000-0000-7000-8000-000000000099")
	secondBytes, err := workspace.Canonical(second)
	must(t, err)
	must(t, ApplyTransaction(root, "seed-second-contract", []FileChange{{Path: recordPath(domain.Contract, secondID), Data: secondBytes, Mode: 0o644}}))

	proposalOne := "Proposal:0199b000-0000-7000-8000-000000000041"
	proposalTwo := "Proposal:0199b000-0000-7000-8000-000000000042"
	missionOne := "Mission:0199b000-0000-7000-8000-000000000043"
	missionTwo := "Mission:0199b000-0000-7000-8000-000000000044"
	assessmentOne := "Assessment:0199b000-0000-7000-8000-000000000045"
	assessmentTwo := "Assessment:0199b000-0000-7000-8000-000000000046"
	svc = openService(t, root)
	secondEntry := lookup(t, svc, secondRef, domain.Contract)
	seed := []struct {
		proposal, mission, assessment, target, fp, addition string
	}{
		{proposalOne, missionOne, assessmentOne, contractRef, first.Fingerprint, "First atomic addition."},
		{proposalTwo, missionTwo, assessmentTwo, secondRef, secondEntry.Fingerprint, "Second atomic addition."},
	}
	var changes []FileChange
	for _, item := range seed {
		proposalID, _ := domain.ParseID(strings.TrimPrefix(item.proposal, "Proposal:"))
		proposal := svc.document(domain.Proposal, proposalID, "Accepted atomic delta", "owner", "accepted")
		workspace.SetString(proposal, "target_contract", item.target)
		workspace.SetString(proposal, "base_version", "1")
		workspace.SetString(proposal, "base_fingerprint", item.fp)
		workspace.SetBool(proposal, "new_capability", false)
		workspace.SetStrings(proposal, "additions", []string{item.addition})
		workspace.SetStrings(proposal, "modification_from", nil)
		workspace.SetStrings(proposal, "modification_to", nil)
		workspace.SetStrings(proposal, "removals", nil)
		workspace.SetStrings(proposal, "scope", []string{"v2"})
		missionID, _ := domain.ParseID(strings.TrimPrefix(item.mission, "Mission:"))
		mission := svc.document(domain.Mission, missionID, "Atomic Mission", "owner", "awaiting-assessment")
		proposalTyped, _ := domain.ParseReference(item.proposal)
		mission.Record.Source = &proposalTyped
		workspace.SetStrings(mission, "evidence_claims", []string{})
		workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
		assessmentID, _ := domain.ParseID(strings.TrimPrefix(item.assessment, "Assessment:"))
		assessment := svc.document(domain.Assessment, assessmentID, "Ready assessment", "reviewer", "recorded")
		workspace.SetString(assessment, "mission", item.mission)
		workspace.SetString(assessment, "verdict", "ready-for-owner")
		workspace.SetStrings(assessment, "claims", []string{})
		workspace.SetStrings(assessment, "evidence", []string{})
		for _, doc := range []*workspace.Document{proposal, mission, assessment} {
			data, canonicalErr := workspace.Canonical(doc)
			must(t, canonicalErr)
			changes = append(changes, FileChange{Path: recordPath(doc.Record.Type, doc.Record.ID), Data: data, Mode: 0o644})
		}
	}
	must(t, ApplyTransaction(root, "seed-atomic-records", changes))

	decisionOne := createDecision(t, root, "0199b000-0000-7000-8000-000000000047", "contract.reconcile", contractRef, first.Fingerprint, []string{assessmentOne}, "accept")
	decisionTwo := createDecision(t, root, "0199b000-0000-7000-8000-000000000048", "contract.reconcile", secondRef, secondEntry.Fingerprint, []string{assessmentTwo}, "accept")
	beforeFirst := fileDigest(t, filepath.Join(root, ".spectacular", "records", "contract-0199b000-0000-7000-8000-000000000001.md"))
	beforeSecond := fileDigest(t, filepath.Join(root, ".spectacular", "records", "contract-0199b000-0000-7000-8000-000000000040.md"))
	svc = openService(t, root)
	bad := []ReconcileInput{
		{Contract: contractRef, Proposal: proposalOne, Authorization: decisionOne, ExpectedFingerprint: first.Fingerprint, IdempotencyKey: "atomic-two"},
		{Contract: secondRef, Proposal: proposalTwo, Authorization: decisionTwo, ExpectedFingerprint: strings.Repeat("0", 64), IdempotencyKey: "atomic-two"},
	}
	if _, err := svc.ReconcileMany(bad); refusalCode(err) != domain.RefusalStaleFingerprint {
		t.Fatalf("stale multi-Contract set err=%v", err)
	}
	if fileDigest(t, filepath.Join(root, ".spectacular", "records", "contract-0199b000-0000-7000-8000-000000000001.md")) != beforeFirst || fileDigest(t, filepath.Join(root, ".spectacular", "records", "contract-0199b000-0000-7000-8000-000000000040.md")) != beforeSecond {
		t.Fatal("multi-Contract validation partially mutated current truth")
	}
	svc = openService(t, root)
	bad[1].ExpectedFingerprint = secondEntry.Fingerprint
	results, err := svc.ReconcileMany(bad)
	must(t, err)
	if len(results) != 2 || results[0].Receipt == "" || results[0].Receipt != results[1].Receipt {
		t.Fatalf("atomic reconciliation results=%#v", results)
	}
	for _, ref := range []string{contractRef, secondRef} {
		view, viewErr := openService(t, root).ContractView(ref)
		must(t, viewErr)
		if view["version"] != "2" {
			t.Fatalf("%s version=%v", ref, view["version"])
		}
	}
}

func TestZeroDeltaAbandonedAndSupersededMissionsResolveThenArchive(t *testing.T) {
	for index, disposition := range []string{"abandoned", "superseded"} {
		t.Run(disposition, func(t *testing.T) {
			root := governedFixture(t)
			svc := openService(t, root)
			suffix := 60 + index*5
			proposalRef := "Proposal:0199b000-0000-7000-8000-0000000000" + string(rune('0'+suffix/10)) + string(rune('0'+suffix%10))
			missionRef := "Mission:0199b000-0000-7000-8000-0000000000" + string(rune('0'+(suffix+1)/10)) + string(rune('0'+(suffix+1)%10))
			assessmentRef := "Assessment:0199b000-0000-7000-8000-0000000000" + string(rune('0'+(suffix+2)/10)) + string(rune('0'+(suffix+2)%10))
			proposalID, _ := domain.ParseID(strings.TrimPrefix(proposalRef, "Proposal:"))
			proposal := svc.document(domain.Proposal, proposalID, "Zero delta Proposal", "owner", "accepted")
			workspace.SetString(proposal, "target_contract", contractRef)
			missionID, _ := domain.ParseID(strings.TrimPrefix(missionRef, "Mission:"))
			mission := svc.document(domain.Mission, missionID, "Terminal Mission", "owner", "awaiting-assessment")
			proposalTyped, _ := domain.ParseReference(proposalRef)
			mission.Record.Source = &proposalTyped
			workspace.SetStrings(mission, "scope", []string{"v2"})
			workspace.SetStrings(mission, "evidence_claims", []string{})
			workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
			assessmentID, _ := domain.ParseID(strings.TrimPrefix(assessmentRef, "Assessment:"))
			assessment := svc.document(domain.Assessment, assessmentID, "Ready assessment", "reviewer", "recorded")
			workspace.SetString(assessment, "mission", missionRef)
			workspace.SetString(assessment, "verdict", "ready-for-owner")
			workspace.SetStrings(assessment, "claims", []string{})
			workspace.SetStrings(assessment, "evidence", []string{})
			var seed []FileChange
			for _, doc := range []*workspace.Document{proposal, mission, assessment} {
				data, err := workspace.Canonical(doc)
				must(t, err)
				seed = append(seed, FileChange{Path: recordPath(doc.Record.Type, doc.Record.ID), Data: data, Mode: 0o644})
			}
			must(t, ApplyTransaction(root, "seed-zero-delta-"+disposition, seed))
			missionEntry := lookup(t, openService(t, root), missionRef, domain.Mission)
			resolveDecision := createDecision(t, root, "0199b000-0000-7000-8000-0000000000"+string(rune('0'+(suffix+3)/10))+string(rune('0'+(suffix+3)%10)), "mission.transition.resolved", missionRef, missionEntry.Fingerprint, []string{assessmentRef}, disposition, "reconciliation.not-required", "terminal-next-action:owner selects next Mission")
			svc = openService(t, root)
			resolved, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: resolveDecision, ExpectedFingerprint: missionEntry.Fingerprint, IdempotencyKey: "resolve-zero-" + disposition, Disposition: disposition, Assessment: assessmentRef, TerminalNextAction: "owner selects next Mission"})
			must(t, err)
			archiveDecision := createDecision(t, root, "0199b000-0000-7000-8000-0000000000"+string(rune('0'+(suffix+4)/10))+string(rune('0'+(suffix+4)%10)), "mission.archive", missionRef, resolved.Fingerprint, []string{assessmentRef}, "no-contract-delta")
			svc = openService(t, root)
			_, err = svc.ArchiveMission(ArchiveInput{Mission: missionRef, Authorization: archiveDecision, ExpectedFingerprint: resolved.Fingerprint, IdempotencyKey: "archive-zero-" + disposition, TerminalPacket: missionRef})
			must(t, err)
			archived := lookup(t, openService(t, root), missionRef, domain.Mission)
			if mustString(archived.Document, "disposition") != disposition {
				t.Fatalf("disposition=%s", mustString(archived.Document, "disposition"))
			}
		})
	}
}

func TestCompletedResolutionRejectsPendingObjectiveAndFakeReceipt(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	proposalRef := "Proposal:0199b000-0000-7000-8000-000000000070"
	missionRef := "Mission:0199b000-0000-7000-8000-000000000071"
	objectiveRef := "Objective:0199b000-0000-7000-8000-000000000072"
	assessmentRef := "Assessment:0199b000-0000-7000-8000-000000000073"
	proposalID, _ := domain.ParseID(strings.TrimPrefix(proposalRef, "Proposal:"))
	proposal := svc.document(domain.Proposal, proposalID, "Accepted delta", "owner", "accepted")
	workspace.SetString(proposal, "target_contract", contractRef)
	missionID, _ := domain.ParseID(strings.TrimPrefix(missionRef, "Mission:"))
	mission := svc.document(domain.Mission, missionID, "Pending Objective Mission", "owner", "awaiting-assessment")
	proposalTyped, _ := domain.ParseReference(proposalRef)
	mission.Record.Source = &proposalTyped
	workspace.SetStrings(mission, "objectives", []string{objectiveRef})
	workspace.SetStrings(mission, "scope", []string{"v2"})
	workspace.SetStrings(mission, "evidence_claims", []string{})
	workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
	objectiveID, _ := domain.ParseID(strings.TrimPrefix(objectiveRef, "Objective:"))
	objective := svc.document(domain.Objective, objectiveID, "Still pending", "owner", "pending")
	workspace.SetString(objective, "mission", missionRef)
	assessmentID, _ := domain.ParseID(strings.TrimPrefix(assessmentRef, "Assessment:"))
	assessment := svc.document(domain.Assessment, assessmentID, "Ready assessment", "reviewer", "recorded")
	workspace.SetString(assessment, "mission", missionRef)
	workspace.SetString(assessment, "verdict", "ready-for-owner")
	workspace.SetStrings(assessment, "claims", []string{})
	workspace.SetStrings(assessment, "evidence", []string{})
	var changes []FileChange
	for _, doc := range []*workspace.Document{proposal, mission, objective, assessment} {
		data, err := workspace.Canonical(doc)
		must(t, err)
		changes = append(changes, FileChange{Path: recordPath(doc.Record.Type, doc.Record.ID), Data: data, Mode: 0o644})
	}
	must(t, ApplyTransaction(root, "seed-pending-objective", changes))
	missionEntry := lookup(t, openService(t, root), missionRef, domain.Mission)
	decisionRef := createDecision(t, root, "0199b000-0000-7000-8000-000000000074", "mission.transition.resolved", missionRef, missionEntry.Fingerprint, []string{assessmentRef}, "completed", "terminal-next-action:arbitrary executor text", "objective.satisfy:"+objectiveRef)
	svc = openService(t, root)
	_, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: decisionRef, ExpectedFingerprint: missionEntry.Fingerprint, IdempotencyKey: "fake-receipt", Disposition: "completed", Assessment: assessmentRef, Reconciliation: "Evidence:0199b000-0000-7000-8000-000000000075", TerminalNextAction: "arbitrary executor text"})
	if refusalCode(err) != domain.RefusalUnsettledReconcile {
		t.Fatalf("pending Objective/fake receipt err=%v", err)
	}
	mission.Record.Status = stringPtr("resolved")
	workspace.SetString(mission, "assessment", assessmentRef)
	workspace.SetString(mission, "disposition", "completed")
	workspace.SetString(mission, "reconciliation", "Evidence:0199b000-0000-7000-8000-000000000075")
	workspace.SetString(mission, "terminal_next_action", "archive-chosen text")
	workspace.SetString(mission, "last_authorization", decisionRef)
	forged, err := workspace.Canonical(mission)
	must(t, err)
	must(t, ApplyTransaction(root, "forge-resolved-counterexample", []FileChange{{Path: recordPath(domain.Mission, missionID), Data: forged, Mode: 0o644}}))
	forgedEntry := lookup(t, openService(t, root), missionRef, domain.Mission)
	archiveDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000076", "mission.archive", missionRef, forgedEntry.Fingerprint, []string{assessmentRef}, "archive")
	_, err = openService(t, root).ArchiveMission(ArchiveInput{Mission: missionRef, Authorization: archiveDecision, ExpectedFingerprint: forgedEntry.Fingerprint, IdempotencyKey: "archive-forged", TerminalPacket: missionRef})
	if refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("archive authority supplied arbitrary terminal continuation: %v", err)
	}
	workspace.SetString(mission, "terminal_next_action", "arbitrary executor text")
	forged, err = workspace.Canonical(mission)
	must(t, err)
	must(t, ApplyTransaction(root, "forge-receipt-counterexample", []FileChange{{Path: recordPath(domain.Mission, missionID), Data: forged, Mode: 0o644}}))
	forgedEntry = lookup(t, openService(t, root), missionRef, domain.Mission)
	archiveDecision = createDecision(t, root, "0199b000-0000-7000-8000-000000000077", "mission.archive", missionRef, forgedEntry.Fingerprint, []string{assessmentRef}, "archive")
	_, err = openService(t, root).ArchiveMission(ArchiveInput{Mission: missionRef, Authorization: archiveDecision, ExpectedFingerprint: forgedEntry.Fingerprint, IdempotencyKey: "archive-fake-receipt", TerminalPacket: missionRef})
	if refusalCode(err) != domain.RefusalUnsettledReconcile {
		t.Fatalf("forged resolved Mission archived with fake receipt/pending Objective: %v", err)
	}
}

func TestAssessmentRejectsStaleSupportAndOmittedContraryEvidence(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	missionRef := "Mission:0199b000-0000-7000-8000-000000000080"
	missionID, _ := domain.ParseID(strings.TrimPrefix(missionRef, "Mission:"))
	mission := svc.document(domain.Mission, missionID, "Evidence Mission", "owner", "awaiting-assessment")
	workspace.SetStrings(mission, "evidence_claims", []string{"claim:material"})
	workspace.SetStrings(mission, "scope", []string{"v2"})
	workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
	var changes []FileChange
	missionBytes, err := workspace.Canonical(mission)
	must(t, err)
	changes = append(changes, FileChange{Path: recordPath(domain.Mission, missionID), Data: missionBytes, Mode: 0o644})
	for number, contrary := range map[int]bool{81: false, 82: true} {
		idText := "0199b000-0000-7000-8000-0000000000" + string(rune('0'+number/10)) + string(rune('0'+number%10))
		id, _ := domain.ParseID(idText)
		evidence := svc.document(domain.Evidence, id, "Claim evidence", "executor", "recorded")
		workspace.SetString(evidence, "mission", missionRef)
		workspace.SetString(evidence, "claim", "claim:material")
		workspace.SetString(evidence, "classification", "direct")
		workspace.SetStrings(evidence, "required_checks", []string{"check"})
		workspace.SetStrings(evidence, "check_results", []string{"pass:check"})
		workspace.SetString(evidence, "review_state", "independent-accepted")
		workspace.SetBool(evidence, "executor_authored", false)
		if contrary {
			workspace.SetStrings(evidence, "contrary_evidence", []string{"contradicts support"})
		} else {
			workspace.SetString(evidence, "freshness_valid_until", "2026-08-10T09:00:00Z")
		}
		data, canonicalErr := workspace.Canonical(evidence)
		must(t, canonicalErr)
		changes = append(changes, FileChange{Path: recordPath(domain.Evidence, id), Data: data, Mode: 0o644})
	}
	must(t, ApplyTransaction(root, "seed-conflicting-evidence", changes))
	assessmentRef := "Assessment:0199b000-0000-7000-8000-000000000083"
	decisionRef := createDecision(t, root, "0199b000-0000-7000-8000-000000000084", "assessment.record", assessmentRef, "absent", nil, "record")
	svc = openService(t, root)
	_, err = svc.RecordAssessment(AssessmentInput{ID: strings.TrimPrefix(assessmentRef, "Assessment:"), Title: "Omitting contrary evidence", Mission: missionRef, Verdict: "ready-for-owner", Actor: "reviewer", Claims: []string{"claim:material"}, Evidence: []string{"Evidence:0199b000-0000-7000-8000-000000000081"}, RecoveryPoint: "git-head", Authorization: decisionRef, IdempotencyKey: "assessment-omission"})
	if refusalCode(err) != domain.RefusalInsufficientEvidence {
		t.Fatalf("stale/omitted contrary Evidence err=%v", err)
	}
}

func TestAuthorityEffectsConditionsAndMissionExpiryAreEnforced(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	contract := lookup(t, svc, contractRef, domain.Contract)
	proposalRef := "Proposal:0199b000-0000-7000-8000-000000000085"
	decisionRef := "Decision:0199b000-0000-7000-8000-000000000086"
	_, err := svc.CreateDecision(DecisionInput{ID: strings.TrimPrefix(decisionRef, "Decision:"), Title: "Wrong effect", Actor: "Alex", ActorRole: "owner", AuthorityBasis: "test", Question: "create", Scope: []string{"v2"}, Disposition: "approve", Rationale: "test", Targets: []string{proposalRef}, ExpectedFingerprints: []string{"absent"}, Operation: "proposal.create", AuthorizedEffects: []string{"mission.create"}, Conditions: []string{"unknown-human-condition"}, ExpiresAt: "2026-08-11T10:00:00Z", IdempotencyKey: "wrong-effect"})
	must(t, err)
	svc = openService(t, root)
	_, err = svc.CreateProposal(ProposalInput{ID: strings.TrimPrefix(proposalRef, "Proposal:"), Title: "Unauthorized delta", Actor: "owner", Status: "accepted", TargetContract: contractRef, BaseVersion: "1", BaseFingerprint: contract.Fingerprint, Additions: []string{"unauthorized"}, Rationale: "test", Scope: []string{"v2"}, Authorization: decisionRef, IdempotencyKey: "unauthorized-proposal"})
	if refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("wrong Decision effect/condition err=%v", err)
	}
	conditionProposal := "Proposal:0199b000-0000-7000-8000-000000000096"
	conditionDecision := "Decision:0199b000-0000-7000-8000-000000000097"
	_, err = openService(t, root).CreateDecision(DecisionInput{ID: strings.TrimPrefix(conditionDecision, "Decision:"), Title: "Unknown condition", Actor: "Alex", ActorRole: "owner", AuthorityBasis: "test", Question: "create", Scope: []string{"v2"}, Disposition: "approve", Rationale: "test", Targets: []string{conditionProposal}, ExpectedFingerprints: []string{"absent"}, Operation: "proposal.create", AuthorizedEffects: []string{"proposal.create"}, Conditions: []string{"unknown-human-condition"}, ExpiresAt: "2026-08-11T10:00:00Z", IdempotencyKey: "unknown-condition"})
	must(t, err)
	conditionInput := ProposalInput{ID: strings.TrimPrefix(conditionProposal, "Proposal:"), Title: "Condition delta", Actor: "owner", Status: "accepted", TargetContract: contractRef, BaseVersion: "1", BaseFingerprint: contract.Fingerprint, Additions: []string{"condition"}, Rationale: "test", Scope: []string{"v2"}, Authorization: conditionDecision, IdempotencyKey: "condition-proposal"}
	if _, err := openService(t, root).CreateProposal(conditionInput); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("unevaluated Decision condition err=%v", err)
	}
	scopeProposal := "Proposal:0199b000-0000-7000-8000-000000000098"
	scopeDecision := "Decision:0199b000-0000-7000-8000-000000000099"
	_, err = openService(t, root).CreateDecision(DecisionInput{ID: strings.TrimPrefix(scopeDecision, "Decision:"), Title: "Narrow scope", Actor: "Alex", ActorRole: "owner", AuthorityBasis: "test", Question: "create", Scope: []string{"v1"}, Disposition: "approve", Rationale: "test", Targets: []string{scopeProposal}, ExpectedFingerprints: []string{"absent"}, Operation: "proposal.create", AuthorizedEffects: []string{"proposal.create"}, ExpiresAt: "2026-08-11T10:00:00Z", IdempotencyKey: "narrow-scope"})
	must(t, err)
	scopeInput := conditionInput
	scopeInput.ID = strings.TrimPrefix(scopeProposal, "Proposal:")
	scopeInput.Authorization = scopeDecision
	scopeInput.IdempotencyKey = "scope-proposal"
	if _, err := openService(t, root).CreateProposal(scopeInput); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-scope Decision err=%v", err)
	}

	missionRef := "Mission:0199b000-0000-7000-8000-000000000087"
	missionID, _ := domain.ParseID(strings.TrimPrefix(missionRef, "Mission:"))
	mission := svc.document(domain.Mission, missionID, "Expired Mission", "owner", "defined")
	workspace.SetStrings(mission, "scope", []string{"v2"})
	workspace.SetString(mission, "expires_at", "2026-08-10T09:00:00Z")
	data, err := workspace.Canonical(mission)
	must(t, err)
	must(t, ApplyTransaction(root, "seed-expired-mission", []FileChange{{Path: recordPath(domain.Mission, missionID), Data: data, Mode: 0o644}}))
	missionEntry := lookup(t, openService(t, root), missionRef, domain.Mission)
	activate := createDecision(t, root, "0199b000-0000-7000-8000-000000000088", "mission.transition.active", missionRef, missionEntry.Fingerprint, nil, "activate")
	_, err = openService(t, root).TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activate, ExpectedFingerprint: missionEntry.Fingerprint, IdempotencyKey: "expired-activate"})
	if refusalCode(err) != domain.RefusalExpiredAuthority {
		t.Fatalf("expired Mission envelope err=%v", err)
	}
}

func TestTypedRepairAttemptsMustChangeAndStayWithinBudget(t *testing.T) {
	root := governedFixture(t)
	svc := openService(t, root)
	missionRef := "Mission:0199b000-0000-7000-8000-000000000089"
	beforeRef := "Evidence:0199b000-0000-7000-8000-000000000090"
	afterRef := "Evidence:0199b000-0000-7000-8000-000000000091"
	missionID, _ := domain.ParseID(strings.TrimPrefix(missionRef, "Mission:"))
	mission := svc.document(domain.Mission, missionID, "Repair Mission", "owner", "awaiting-assessment")
	workspace.SetStrings(mission, "scope", []string{"v2"})
	workspace.SetStrings(mission, "evidence_claims", []string{"claim:repair"})
	workspace.SetInt(mission, "repair_budget", 1)
	workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
	var changes []FileChange
	for _, ref := range []string{beforeRef, afterRef} {
		id, _ := domain.ParseID(strings.TrimPrefix(ref, "Evidence:"))
		evidence := svc.document(domain.Evidence, id, "Repair Evidence", "executor", "recorded")
		workspace.SetString(evidence, "mission", missionRef)
		workspace.SetString(evidence, "claim", "claim:repair")
		data, err := workspace.Canonical(evidence)
		must(t, err)
		changes = append(changes, FileChange{Path: recordPath(domain.Evidence, id), Data: data, Mode: 0o644})
	}
	missionData, err := workspace.Canonical(mission)
	must(t, err)
	changes = append(changes, FileChange{Path: recordPath(domain.Mission, missionID), Data: missionData, Mode: 0o644})
	must(t, ApplyTransaction(root, "seed-repair", changes))
	assessmentRef := "Assessment:0199b000-0000-7000-8000-000000000092"
	decisionRef := createDecision(t, root, "0199b000-0000-7000-8000-000000000093", "assessment.record", assessmentRef, "absent", nil, "repair-required")
	base := AssessmentInput{ID: strings.TrimPrefix(assessmentRef, "Assessment:"), Title: "Repair assessment", Mission: missionRef, Verdict: "repair-required", Actor: "reviewer", Claims: []string{"claim:repair"}, Evidence: []string{afterRef}, RecoveryPoint: "git-head", Authorization: decisionRef, IdempotencyKey: "repair-assessment"}
	unchanged := RepairAttempt{AffectedClaims: []string{"claim:repair"}, PreviousHypothesis: "same", NewHypothesis: "same", Actor: "executor", BeforeEvidence: []string{beforeRef}, AfterEvidence: []string{afterRef}, Checks: []string{"test"}, Result: "failed", BudgetConsumed: 1, RecoveryPoint: "git-head"}
	base.RepairAttempts = []RepairAttempt{unchanged}
	if _, err := openService(t, root).RecordAssessment(base); refusalCode(err) != domain.RefusalInvalidKnownField {
		t.Fatalf("unchanged repair attempt err=%v", err)
	}
	overBudget := unchanged
	overBudget.NewHypothesis = "changed"
	overBudget.BudgetConsumed = 2
	base.RepairAttempts = []RepairAttempt{overBudget}
	if _, err := openService(t, root).RecordAssessment(base); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("over-budget repair attempt err=%v", err)
	}
	valid := overBudget
	valid.BudgetConsumed = 1
	base.RepairAttempts = []RepairAttempt{valid}
	_, err = openService(t, root).RecordAssessment(base)
	must(t, err)
	secondAssessment := "Assessment:0199b000-0000-7000-8000-000000000094"
	secondDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000095", "assessment.record", secondAssessment, "absent", nil, "repair-required")
	base.ID = strings.TrimPrefix(secondAssessment, "Assessment:")
	base.Authorization = secondDecision
	base.IdempotencyKey = "repair-assessment-second"
	if _, err := openService(t, root).RecordAssessment(base); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("cumulative repair budget bypass err=%v", err)
	}
}

func governedFixture(t *testing.T) string {
	t.Helper()
	source := filepath.Join("..", "..", "testdata", "scenario-b-c")
	root := t.TempDir()
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(source, path)
		dest := filepath.Join(root, rel)
		if info.IsDir() {
			return os.MkdirAll(dest, info.Mode())
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		return os.WriteFile(dest, data, info.Mode())
	})
	must(t, err)
	return root
}

func openService(t *testing.T, root string) Service {
	t.Helper()
	ws, err := discovery.Open(root)
	must(t, err)
	return Service{Workspace: ws, Now: func() time.Time { return testNow }}
}

func lookup(t *testing.T, svc Service, ref string, noun domain.RecordType) discovery.Entry {
	t.Helper()
	entry, err := svc.Workspace.Lookup(ref, noun)
	must(t, err)
	return entry
}

func createDecision(t *testing.T, root, id, operation, target, expected string, evidence []string, disposition string, extraEffects ...string) string {
	t.Helper()
	svc := openService(t, root)
	ref := "Decision:" + id
	effects := append([]string{operation}, extraEffects...)
	_, err := svc.CreateDecision(DecisionInput{ID: id, Title: operation, Actor: "Alex", ActorRole: "owner", AuthorityBasis: "accepted B+C contract", Question: operation, Scope: []string{"v2"}, Disposition: disposition, Rationale: "Explicit owner authorization.", Targets: []string{target}, ExpectedFingerprints: []string{expected}, Operation: operation, AuthorizedEffects: effects, ExpiresAt: "2026-08-11T10:00:00Z", Evidence: evidence, IdempotencyKey: "decision-" + id})
	must(t, err)
	return ref
}

func refusalCode(err error) domain.RefusalCode {
	var refusal *domain.Refusal
	if errors.As(err, &refusal) {
		return refusal.Code
	}
	return ""
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func fileDigest(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	must(t, err)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
