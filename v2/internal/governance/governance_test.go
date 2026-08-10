package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/projection"
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
		Baseline: contract.Fingerprint, Budget: "two repairs", ExpiresAt: "2026-08-11T10:00:00Z", Stops: []string{"authority-drift"}, RecoveryPoint: "git-head", ReturnDestination: "central", Authorization: missionDecision, ExpectedProposalFingerprint: proposal.Fingerprint, IdempotencyKey: "mission-create-1",
	}
	_, err = svc.CreateMission(missionInput)
	must(t, err)

	mission := lookup(t, openService(t, root), missionRef, domain.Mission)
	activateDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000013", "mission.transition.active", missionRef, mission.Fingerprint, nil, "activate")
	svc = openService(t, root)
	active, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activateDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-active-1"})
	must(t, err)
	if replay, err := openService(t, root).TransitionMission(TransitionInput{Mission: missionRef, To: "active", Authorization: activateDecision, ExpectedFingerprint: mission.Fingerprint, IdempotencyKey: "mission-active-1"}); err != nil || !replay.IdempotentReplay {
		t.Fatalf("transition replay=%#v err=%v", replay, err)
	}

	handoffDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000014", "handoff.create", handoffRef, "absent", nil, "dispatch")
	svc = openService(t, root)
	handoffInput := HandoffInput{
		ID: strings.TrimPrefix(handoffRef, "Handoff:"), Title: "Bounded fake-provider work", Mission: missionRef, Objective: objectiveRef, Run: runRef, Sender: "owner", Actor: "executor", Destination: "replacement-runtime", HostPointer: "host-task:disposable",
		Scope: []string{"v2"}, AllowedActions: []string{"test"}, ForbiddenEffects: []string{"provider-mutation"}, EvidenceClaims: []string{"claim:closure"}, Budget: "one attempt", ExpiresAt: "2026-08-11T10:00:00Z", Stops: []string{"authority-drift"}, RecoveryPoint: "git-head", ReturnDestination: "central", Authorization: handoffDecision, ExpectedMissionFingerprint: active.Fingerprint, IdempotencyKey: "handoff-create-1",
	}
	outOfScope := handoffInput
	outOfScope.Scope = []string{"v1"}
	if _, err := svc.CreateHandoff(outOfScope); refusalCode(err) != domain.RefusalUnauthorized {
		t.Fatalf("out-of-envelope Handoff err=%v", err)
	}
	handoff, err := svc.CreateHandoff(handoffInput)
	must(t, err)
	svc = openService(t, root)
	validated, err := svc.ValidateHandoff(handoffRef)
	must(t, err)
	if validated["valid"] != true || len(validated["does_not_prove"].([]string)) == 0 {
		t.Fatalf("handoff validation overclaims: %#v", validated)
	}
	_, err = svc.ReturnHandoff(HandoffReturnInput{
		ID: strings.TrimPrefix(returnRef, "Handoff:"), Title: "Returned bounded work", Dispatch: handoffRef, Status: "succeeded", Actor: "replacement-runtime", FinalBaseline: contract.Fingerprint, Result: "fake receipt only", Actions: []string{"test"}, ProviderReceipts: []string{"fake-provider:ok"}, Evidence: []string{evidenceRef}, BudgetUsed: "one attempt", RecoveryPoint: "git-head", NextAction: "record evidence", ExpectedDispatchFingerprint: handoff.Fingerprint, IdempotencyKey: "handoff-return-1",
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
	_, err = svc.CreateEvidence(EvidenceInput{ID: strings.TrimPrefix(badEvidenceRef, "Evidence:"), Title: "Unreviewed executor observation", Mission: missionRef, Objective: objectiveRef, Claim: "claim:closure", Classification: "observation", Scope: []string{"v2"}, Method: "executor report", Actor: "executor", Target: proposalRef, Environment: "disposable", ObservedAt: "2026-08-10T10:00:00Z", FreshnessValidUntil: "2026-08-11T10:00:00Z", RequiredChecks: []string{"go-test"}, CheckResults: []string{"pass:go-test"}, ReviewState: "unreviewed", ExecutorAuthored: true, Authorization: badEvidenceDecision, IdempotencyKey: "evidence-bad-1"})
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

	resolveDecision := createDecision(t, root, "0199b000-0000-7000-8000-000000000019", "mission.transition.resolved", missionRef, awaiting.Fingerprint, []string{assessmentRef, reconciled.Receipt}, "completed")
	svc = openService(t, root)
	resolved, err := svc.TransitionMission(TransitionInput{Mission: missionRef, To: "resolved", Authorization: resolveDecision, ExpectedFingerprint: awaiting.Fingerprint, IdempotencyKey: "mission-resolve-1", Disposition: "completed", Assessment: assessmentRef, Reconciliation: reconciled.Receipt, TerminalNextAction: "review next governed request"})
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
		missionID, _ := domain.ParseID(strings.TrimPrefix(item.mission, "Mission:"))
		mission := svc.document(domain.Mission, missionID, "Atomic Mission", "owner", "awaiting-assessment")
		proposalTyped, _ := domain.ParseReference(item.proposal)
		mission.Record.Source = &proposalTyped
		assessmentID, _ := domain.ParseID(strings.TrimPrefix(item.assessment, "Assessment:"))
		assessment := svc.document(domain.Assessment, assessmentID, "Ready assessment", "reviewer", "recorded")
		workspace.SetString(assessment, "mission", item.mission)
		workspace.SetString(assessment, "verdict", "ready-for-owner")
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
			assessmentID, _ := domain.ParseID(strings.TrimPrefix(assessmentRef, "Assessment:"))
			assessment := svc.document(domain.Assessment, assessmentID, "Ready assessment", "reviewer", "recorded")
			workspace.SetString(assessment, "mission", missionRef)
			workspace.SetString(assessment, "verdict", "ready-for-owner")
			var seed []FileChange
			for _, doc := range []*workspace.Document{proposal, mission, assessment} {
				data, err := workspace.Canonical(doc)
				must(t, err)
				seed = append(seed, FileChange{Path: recordPath(doc.Record.Type, doc.Record.ID), Data: data, Mode: 0o644})
			}
			must(t, ApplyTransaction(root, "seed-zero-delta-"+disposition, seed))
			missionEntry := lookup(t, openService(t, root), missionRef, domain.Mission)
			resolveDecision := createDecision(t, root, "0199b000-0000-7000-8000-0000000000"+string(rune('0'+(suffix+3)/10))+string(rune('0'+(suffix+3)%10)), "mission.transition.resolved", missionRef, missionEntry.Fingerprint, []string{assessmentRef}, "no-contract-delta")
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

func createDecision(t *testing.T, root, id, operation, target, expected string, evidence []string, disposition string) string {
	t.Helper()
	svc := openService(t, root)
	ref := "Decision:" + id
	_, err := svc.CreateDecision(DecisionInput{ID: id, Title: operation, Actor: "Alex", ActorRole: "owner", AuthorityBasis: "accepted B+C contract", Question: operation, Scope: []string{"v2"}, Disposition: disposition, Rationale: "Explicit owner authorization.", Targets: []string{target}, ExpectedFingerprints: []string{expected}, Operation: operation, AuthorizedEffects: []string{operation}, ExpiresAt: "2026-08-11T10:00:00Z", Evidence: evidence, IdempotencyKey: "decision-" + id})
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
