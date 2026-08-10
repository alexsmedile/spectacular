package governance

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func (s Service) CreateHandoff(input HandoffInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	if mission.Document.Record.Status == nil || *mission.Document.Record.Status != "active" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "mission", "Handoff requires an active Mission", nil)
	}
	if mission.Fingerprint != input.ExpectedMissionFingerprint {
		return OperationResult{}, stale("expected_mission_fingerprint", input.ExpectedMissionFingerprint, mission.Fingerprint)
	}
	objective, err := s.Workspace.Lookup(input.Objective, domain.Objective)
	if err != nil {
		return OperationResult{}, err
	}
	run, err := s.Workspace.Lookup(input.Run, domain.Run)
	if err != nil {
		return OperationResult{}, err
	}
	if mustString(objective.Document, "mission") != input.Mission || mustString(run.Document, "mission") != input.Mission {
		return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "mission", "Handoff Objective and Run must belong to the same Mission", nil)
	}
	if input.Sender == "" || input.Actor == "" || input.Destination == "" || len(input.Scope) == 0 || len(input.EvidenceClaims) == 0 || input.Budget == "" || input.RecoveryPoint == "" || input.ReturnDestination == "" || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("handoff", "actor, destination, scope, authority, evidence, budget, recovery, return, and idempotency are required")
	}
	if _, err := parseFuture(input.ExpiresAt, s.now(), "expires_at"); err != nil {
		return OperationResult{}, err
	}
	if !subset(input.Scope, mustStrings(mission.Document, "scope")) || !subset(input.AllowedActions, mustStrings(mission.Document, "allowed_actions")) || !subset(mustStrings(mission.Document, "forbidden_effects"), input.ForbiddenEffects) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "envelope", "Handoff must be a subset of Mission scope/actions and preserve every forbidden effect", nil)
	}
	ref := string(domain.Handoff) + ":" + id.String()
	if err := s.authorize(input.Authorization, "handoff.create", ref, "absent"); err != nil {
		return OperationResult{}, err
	}
	if input.Supersedes != "" {
		old, lookupErr := s.Workspace.Lookup(input.Supersedes, domain.Handoff)
		if lookupErr != nil || mustString(old.Document, "kind") != "dispatch" || mustString(old.Document, "mission") != input.Mission {
			return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "supersedes", "superseded Handoff must be an existing dispatch in the same Mission", lookupErr)
		}
	}
	doc := s.document(domain.Handoff, id, input.Title, input.Sender, "")
	workspace.SetString(doc, "kind", "dispatch")
	workspace.SetString(doc, "mission", input.Mission)
	workspace.SetString(doc, "objective", input.Objective)
	workspace.SetString(doc, "run", input.Run)
	workspace.SetString(doc, "sender", input.Sender)
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetString(doc, "destination", input.Destination)
	workspace.SetString(doc, "host_pointer", input.HostPointer)
	workspace.SetStrings(doc, "scope", input.Scope)
	workspace.SetStrings(doc, "inputs", input.Inputs)
	workspace.SetStrings(doc, "allowed_actions", input.AllowedActions)
	workspace.SetStrings(doc, "forbidden_effects", input.ForbiddenEffects)
	workspace.SetStrings(doc, "evidence_claims", input.EvidenceClaims)
	workspace.SetString(doc, "budget", input.Budget)
	workspace.SetString(doc, "expires_at", input.ExpiresAt)
	workspace.SetStrings(doc, "stops", input.Stops)
	workspace.SetString(doc, "recovery_point", input.RecoveryPoint)
	workspace.SetString(doc, "return_destination", input.ReturnDestination)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "expected_mission_fingerprint", input.ExpectedMissionFingerprint)
	workspace.SetString(doc, "supersedes", input.Supersedes)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("handoff.create", doc, input.IdempotencyKey, []string{input.Mission, input.Objective, input.Run, input.Authorization})
}

func (s Service) ValidateHandoff(ref string) (map[string]any, error) {
	handoff, err := s.Workspace.Lookup(ref, domain.Handoff)
	if err != nil {
		return nil, err
	}
	if mustString(handoff.Document, "kind") != "dispatch" {
		return nil, invalid("kind", "handoff validate requires a dispatch")
	}
	for _, candidate := range s.Workspace.OfType(domain.Handoff) {
		if mustString(candidate.Document, "supersedes") == refOf(handoff) {
			return nil, domain.NewRefusal(domain.RefusalConflictingAuthority, "supersedes", "Handoff was superseded by "+refOf(candidate), nil)
		}
	}
	missionRef := mustString(handoff.Document, "mission")
	mission, err := s.Workspace.Lookup(missionRef, domain.Mission)
	if err != nil {
		return nil, err
	}
	objective, err := s.Workspace.Lookup(mustString(handoff.Document, "objective"), domain.Objective)
	if err != nil {
		return nil, err
	}
	run, err := s.Workspace.Lookup(mustString(handoff.Document, "run"), domain.Run)
	if err != nil {
		return nil, err
	}
	if mustString(objective.Document, "mission") != missionRef || mustString(run.Document, "mission") != missionRef {
		return nil, domain.NewRefusal(domain.RefusalConflictingAuthority, "mission", "Handoff containment no longer holds", nil)
	}
	if mustString(handoff.Document, "expected_mission_fingerprint") != mission.Fingerprint {
		return nil, stale("expected_mission_fingerprint", mustString(handoff.Document, "expected_mission_fingerprint"), mission.Fingerprint)
	}
	if !subset(mustStrings(handoff.Document, "scope"), mustStrings(mission.Document, "scope")) || !subset(mustStrings(handoff.Document, "allowed_actions"), mustStrings(mission.Document, "allowed_actions")) || !subset(mustStrings(mission.Document, "forbidden_effects"), mustStrings(handoff.Document, "forbidden_effects")) {
		return nil, domain.NewRefusal(domain.RefusalUnauthorized, "envelope", "Handoff no longer fits current Mission envelope", nil)
	}
	if _, err := parseFuture(mustString(handoff.Document, "expires_at"), s.now(), "expires_at"); err != nil {
		return nil, err
	}
	if err := s.authorize(mustString(handoff.Document, "authorization"), "handoff.create", refOf(handoff), "absent"); err != nil {
		return nil, err
	}
	return map[string]any{"ref": refOf(handoff), "fingerprint": handoff.Fingerprint, "valid": true, "proves": []string{"structure", "same_mission_containment", "envelope_subset", "current_authority", "current_baseline", "expiry", "return_shape"}, "does_not_prove": []string{"actor_identity", "actor_competence", "provider_permissions_or_effects", "evidence_truth_or_sufficiency", "mission_success"}}, nil
}

func (s Service) ReturnHandoff(input HandoffReturnInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	dispatch, err := s.Workspace.Lookup(input.Dispatch, domain.Handoff)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := s.ValidateHandoff(input.Dispatch); err != nil {
		return OperationResult{}, err
	}
	if dispatch.Fingerprint != input.ExpectedDispatchFingerprint {
		return OperationResult{}, stale("expected_dispatch_fingerprint", input.ExpectedDispatchFingerprint, dispatch.Fingerprint)
	}
	if input.Status != "succeeded" && input.Status != "blocked" && input.Status != "failed" {
		return OperationResult{}, invalid("status", "return status must be succeeded, blocked, or failed")
	}
	if input.Actor == "" || input.FinalBaseline == "" || input.Result == "" || input.BudgetUsed == "" || input.RecoveryPoint == "" || input.IdempotencyKey == "" || (input.NextAction == "") == (input.OwnerGate == "") {
		return OperationResult{}, missing("return", "actor, result, baseline, budget, recovery, idempotency, and exactly one next_action or owner_gate are required")
	}
	doc := s.document(domain.Handoff, id, input.Title, input.Actor, "")
	workspace.SetString(doc, "kind", "return")
	workspace.SetString(doc, "dispatch", input.Dispatch)
	workspace.SetString(doc, "mission", mustString(dispatch.Document, "mission"))
	workspace.SetString(doc, "objective", mustString(dispatch.Document, "objective"))
	workspace.SetString(doc, "run", mustString(dispatch.Document, "run"))
	workspace.SetString(doc, "return_status", input.Status)
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetString(doc, "final_baseline", input.FinalBaseline)
	workspace.SetString(doc, "result", input.Result)
	workspace.SetStrings(doc, "actions", input.Actions)
	workspace.SetStrings(doc, "provider_receipts", input.ProviderReceipts)
	workspace.SetStrings(doc, "evidence", input.Evidence)
	workspace.SetStrings(doc, "remaining_gaps", input.RemainingGaps)
	workspace.SetString(doc, "budget_used", input.BudgetUsed)
	workspace.SetString(doc, "recovery_point", input.RecoveryPoint)
	workspace.SetString(doc, "next_action", input.NextAction)
	workspace.SetString(doc, "owner_gate", input.OwnerGate)
	workspace.SetString(doc, "expected_dispatch_fingerprint", input.ExpectedDispatchFingerprint)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("handoff.return", doc, input.IdempotencyKey, []string{input.Dispatch})
}

func (s Service) CreateEvidence(input EvidenceInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	if !contains([]string{"direct", "observation", "proxy", "judgment", "unknown"}, input.Classification) {
		return OperationResult{}, invalid("classification", "must be direct, observation, proxy, judgment, or unknown")
	}
	if input.Title == "" || input.Mission == "" || input.Claim == "" || len(input.Scope) == 0 || input.Method == "" || input.Actor == "" || input.Target == "" || input.ObservedAt == "" || input.FreshnessValidUntil == "" || input.ReviewState == "" || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("evidence", "claim attribution, scope, method, target, freshness, review, authorization, and idempotency are required")
	}
	if _, err := s.Workspace.Lookup(input.Mission, domain.Mission); err != nil {
		return OperationResult{}, err
	}
	if _, err := time.Parse(time.RFC3339, input.ObservedAt); err != nil {
		return OperationResult{}, invalid("observed_at", "must be RFC3339")
	}
	if _, err := parseFuture(input.FreshnessValidUntil, s.now(), "freshness_valid_until"); err != nil {
		return OperationResult{}, err
	}
	ref := string(domain.Evidence) + ":" + id.String()
	if err := s.authorize(input.Authorization, "evidence.create", ref, "absent"); err != nil {
		return OperationResult{}, err
	}
	doc := s.document(domain.Evidence, id, input.Title, input.Actor, "")
	workspace.SetString(doc, "mission", input.Mission)
	workspace.SetString(doc, "objective", input.Objective)
	workspace.SetString(doc, "checkpoint", input.Checkpoint)
	workspace.SetString(doc, "claim", input.Claim)
	workspace.SetString(doc, "classification", input.Classification)
	workspace.SetStrings(doc, "scope", input.Scope)
	workspace.SetString(doc, "method", input.Method)
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetString(doc, "target", input.Target)
	workspace.SetString(doc, "environment", input.Environment)
	workspace.SetString(doc, "observed_at", input.ObservedAt)
	workspace.SetStrings(doc, "limitations", input.Limitations)
	workspace.SetStrings(doc, "contrary_evidence", input.ContraryEvidence)
	workspace.SetStrings(doc, "required_checks", input.RequiredChecks)
	workspace.SetStrings(doc, "check_results", input.CheckResults)
	workspace.SetString(doc, "review_state", input.ReviewState)
	workspace.SetBool(doc, "executor_authored", input.ExecutorAuthored)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	s.addFreshness(doc, input.FreshnessValidUntil)
	return s.createOne("evidence.create", doc, input.IdempotencyKey, []string{input.Mission, input.Authorization})
}

func (s Service) RecordAssessment(input AssessmentInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	if !contains([]string{"ready-for-owner", "repair-required", "escalated"}, input.Verdict) {
		return OperationResult{}, invalid("verdict", "must be ready-for-owner, repair-required, or escalated")
	}
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	if input.Actor == "" || len(input.Claims) == 0 || len(input.Evidence) == 0 || input.RecoveryPoint == "" || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("assessment", "actor, claims, Evidence, recovery, authorization, and idempotency are required")
	}
	if input.Verdict == "ready-for-owner" {
		if len(input.BlockingFindings) > 0 {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "blocking_findings", "owner-ready assessment cannot retain blocking findings", nil)
		}
		declared := mustStrings(mission.Document, "evidence_claims")
		for _, claim := range declared {
			if !contains(input.Claims, claim) {
				return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "claims", "material claim is not assessed: "+claim, nil)
			}
			if err := s.sufficientClaim(input.Mission, claim, input.Evidence); err != nil {
				return OperationResult{}, err
			}
		}
	}
	ref := string(domain.Assessment) + ":" + id.String()
	if err := s.authorize(input.Authorization, "assessment.record", ref, "absent"); err != nil {
		return OperationResult{}, err
	}
	doc := s.document(domain.Assessment, id, input.Title, input.Actor, "")
	workspace.SetString(doc, "mission", input.Mission)
	workspace.SetString(doc, "verdict", input.Verdict)
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetStrings(doc, "claims", input.Claims)
	workspace.SetStrings(doc, "evidence", input.Evidence)
	workspace.SetStrings(doc, "blocking_findings", input.BlockingFindings)
	workspace.SetStrings(doc, "limitations", input.Limitations)
	workspace.SetStrings(doc, "repair_attempts", input.RepairAttempts)
	workspace.SetString(doc, "recovery_point", input.RecoveryPoint)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("assessment.record", doc, input.IdempotencyKey, append([]string{input.Mission, input.Authorization}, input.Evidence...))
}

func (s Service) sufficientClaim(missionRef, claim string, refs []string) error {
	var supporting bool
	for _, ref := range refs {
		entry, err := s.Workspace.Lookup(ref, domain.Evidence)
		if err != nil {
			return err
		}
		if mustString(entry.Document, "mission") != missionRef || mustString(entry.Document, "claim") != claim {
			continue
		}
		classification := mustString(entry.Document, "classification")
		if classification == "unknown" || len(mustStrings(entry.Document, "contrary_evidence")) > 0 {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "evidence", "claim has unknown or conflicting evidence: "+claim, nil)
		}
		checks := mustStrings(entry.Document, "required_checks")
		results := mustStrings(entry.Document, "check_results")
		if len(checks) != len(results) {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "check_results", "required checks are incomplete for claim: "+claim, nil)
		}
		for _, result := range results {
			if !strings.HasPrefix(result, "pass:") {
				return domain.NewRefusal(domain.RefusalInsufficientEvidence, "check_results", "required check did not pass for claim: "+claim, nil)
			}
		}
		executorAuthored, _ := workspace.Bool(entry.Document, "executor_authored", false)
		if executorAuthored && mustString(entry.Document, "review_state") != "independent-accepted" {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "review_state", "executor-authored evidence requires independent acceptance for claim: "+claim, nil)
		}
		supporting = true
	}
	if !supporting {
		return domain.NewRefusal(domain.RefusalInsufficientEvidence, "evidence", "no attributable evidence maps to claim: "+claim, nil)
	}
	return nil
}

func (s Service) TransitionMission(input TransitionInput) (OperationResult, error) {
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	current := value(mission.Document.Record.Status)
	if current == input.To && mustString(mission.Document, "last_idempotency_key") == input.IdempotencyKey {
		return result("mission.transition", mission, true, []string{input.Authorization}), nil
	}
	if mission.Fingerprint != input.ExpectedFingerprint {
		return OperationResult{}, stale("expected_fingerprint", input.ExpectedFingerprint, mission.Fingerprint)
	}
	if !legalTransition(current, input.To) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "status", current+" cannot transition to "+input.To, nil)
	}
	if err := s.authorize(input.Authorization, "mission.transition."+input.To, input.Mission, input.ExpectedFingerprint); err != nil {
		return OperationResult{}, err
	}
	if input.To == "awaiting-assessment" {
		if !s.hasReturnAndEvidence(input.Mission) {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "mission", "awaiting-assessment requires an immutable Handoff return and returned Evidence", nil)
		}
	}
	if input.To == "resolved" {
		assessment, lookupErr := s.Workspace.Lookup(input.Assessment, domain.Assessment)
		if lookupErr != nil || mustString(assessment.Document, "mission") != input.Mission || mustString(assessment.Document, "verdict") != "ready-for-owner" {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "assessment", "resolution requires a ready-for-owner Assessment for this Mission", lookupErr)
		}
		if !contains([]string{"completed", "abandoned", "superseded"}, input.Disposition) {
			return OperationResult{}, invalid("disposition", "resolved disposition must be completed, abandoned, or superseded")
		}
		decision, _ := s.Workspace.Lookup(input.Authorization, domain.Decision)
		decisionDisposition := mustString(decision.Document, "disposition")
		if input.Reconciliation == "" && decisionDisposition != "no-contract-delta" && decisionDisposition != "reconciliation-not-required" {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "resolution requires reconciliation or explicit zero-delta disposition", nil)
		}
		if input.TerminalNextAction == "" {
			return OperationResult{}, missing("terminal_next_action", "resolution requires exactly one terminal continuation or owner gate")
		}
		workspace.SetString(mission.Document, "assessment", input.Assessment)
		workspace.SetString(mission.Document, "disposition", input.Disposition)
		workspace.SetString(mission.Document, "reconciliation", input.Reconciliation)
		workspace.SetString(mission.Document, "terminal_next_action", input.TerminalNextAction)
	}
	mission.Document.Record.Status = stringPtr(input.To)
	mission.Document.Record.Updated = stringPtr(s.now().Format(time.RFC3339))
	workspace.SetString(mission.Document, "last_authorization", input.Authorization)
	workspace.SetString(mission.Document, "last_transition_input_fingerprint", input.ExpectedFingerprint)
	workspace.SetString(mission.Document, "last_idempotency_key", input.IdempotencyKey)
	canonical, err := workspace.Canonical(mission.Document)
	if err != nil {
		return OperationResult{}, err
	}
	if err := ApplyTransaction(s.Workspace.Root, input.IdempotencyKey, []FileChange{{Path: mission.Path, Data: canonical, Mode: 0o644}}); err != nil {
		return OperationResult{}, err
	}
	fp, _ := workspace.Fingerprint(mission.Document)
	return OperationResult{Operation: "mission.transition", Ref: input.Mission, Path: mission.Path, Fingerprint: fp, Sources: []string{input.Authorization}}, nil
}

func (s Service) hasReturnAndEvidence(mission string) bool {
	var returned, evidenced bool
	for _, handoff := range s.Workspace.OfType(domain.Handoff) {
		if mustString(handoff.Document, "mission") == mission && mustString(handoff.Document, "kind") == "return" {
			returned = true
		}
	}
	for _, evidence := range s.Workspace.OfType(domain.Evidence) {
		if mustString(evidence.Document, "mission") == mission {
			evidenced = true
		}
	}
	return returned && evidenced
}

func (s Service) Reconcile(input ReconcileInput) (OperationResult, error) {
	results, err := s.ReconcileMany([]ReconcileInput{input})
	if err != nil {
		return OperationResult{}, err
	}
	return results[0], nil
}

func (s Service) ReconcileMany(inputs []ReconcileInput) ([]OperationResult, error) {
	if len(inputs) == 0 {
		return nil, missing("reconciliation", "at least one Contract is required")
	}
	var changes []FileChange
	results := make([]OperationResult, 0, len(inputs))
	transactionKey := inputs[0].IdempotencyKey
	receiptID := receiptID(transactionKey)
	receiptRef := string(domain.Evidence) + ":" + receiptID.String()
	if receipt, found := s.entryByID(receiptID); found {
		if receipt.Document.Record.Type != domain.Evidence || mustString(receipt.Document, "kind") != "reconciliation-receipt" || mustString(receipt.Document, "idempotency_key") != transactionKey {
			return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "transaction receipt identity has different content", nil)
		}
		for _, input := range inputs {
			contract, err := s.Workspace.Lookup(input.Contract, domain.Contract)
			if err != nil || mustString(contract.Document, "reconciliation_idempotency_key") != transactionKey {
				return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "reconciliation replay does not match current Contract set", err)
			}
			result := result("contract.reconcile", contract, true, []string{input.Proposal, input.Authorization})
			result.Receipt = receiptRef
			results = append(results, result)
		}
		return results, nil
	}
	var receiptSources []string
	for _, input := range inputs {
		if input.IdempotencyKey == "" || input.IdempotencyKey != transactionKey {
			return nil, invalid("idempotency_key", "all Contract items must share one transaction key")
		}
		proposal, err := s.Workspace.Lookup(input.Proposal, domain.Proposal)
		if err != nil {
			return nil, err
		}
		if value(proposal.Document.Record.Status) != "accepted" || mustString(proposal.Document, "target_contract") != input.Contract {
			return nil, domain.NewRefusal(domain.RefusalUnauthorized, "proposal", "accepted Proposal must target the exact Contract", nil)
		}
		newCapability, err := workspace.Bool(proposal.Document, "new_capability", true)
		if err != nil {
			return nil, err
		}
		var existing discovery.Entry
		if newCapability {
			if input.ExpectedFingerprint != "absent" {
				return nil, stale("expected_fingerprint", input.ExpectedFingerprint, "absent")
			}
			if _, lookupErr := s.Workspace.Lookup(input.Contract, domain.Contract); lookupErr == nil {
				return nil, domain.NewRefusal(domain.RefusalCollision, "contract", "new capability target exists", nil)
			}
		} else {
			existing, err = s.Workspace.Lookup(input.Contract, domain.Contract)
			if err != nil {
				return nil, err
			}
			if existing.Fingerprint != input.ExpectedFingerprint || mustString(proposal.Document, "base_fingerprint") != input.ExpectedFingerprint {
				return nil, stale("expected_fingerprint", input.ExpectedFingerprint, existing.Fingerprint)
			}
		}
		if err := s.authorize(input.Authorization, "contract.reconcile", input.Contract, input.ExpectedFingerprint); err != nil {
			return nil, err
		}
		if err := s.requireReadyAssessment(input.Authorization, input.Proposal); err != nil {
			return nil, err
		}
		candidate, err := s.candidate(proposal)
		if err != nil {
			return nil, err
		}
		contractRef, _ := domain.ParseReference(input.Contract)
		version := "1"
		if !newCapability {
			current, parseErr := strconv.Atoi(mustString(existing.Document, "contract_version"))
			if parseErr != nil || current < 1 {
				return nil, invalid("contract_version", "must be a positive integer")
			}
			version = strconv.Itoa(current + 1)
		}
		contract := s.document(domain.Contract, contractRef.ID, "Capability Contract", "owner", "current")
		setContract(contract, candidate, version)
		workspace.SetString(contract, "accepted_proposal", input.Proposal)
		workspace.SetString(contract, "authorization", input.Authorization)
		workspace.SetString(contract, "reconciliation_idempotency_key", input.IdempotencyKey)
		canonical, err := workspace.Canonical(contract)
		if err != nil {
			return nil, err
		}
		path := recordPath(domain.Contract, contractRef.ID)
		if !newCapability {
			path = existing.Path
			old, _ := workspace.Canonical(existing.Document)
			history := filepath.ToSlash(filepath.Join(".spectacular", "history", "contracts", contractRef.ID.String()+"@"+mustString(existing.Document, "contract_version")+".md"))
			changes = append(changes, FileChange{Path: history, Data: old, Mode: 0o644})
		}
		changes = append(changes, FileChange{Path: path, Data: canonical, Mode: 0o644})
		fp, _ := workspace.Fingerprint(contract)
		results = append(results, OperationResult{Operation: "contract.reconcile", Ref: input.Contract, Path: path, Fingerprint: fp, Sources: []string{input.Proposal, input.Authorization}})
		receiptSources = append(receiptSources, input.Contract, input.Proposal, input.Authorization)
	}
	receipt := s.document(domain.Evidence, receiptID, "Atomic Contract reconciliation receipt", "spectacular-cli", "recorded")
	workspace.SetString(receipt, "kind", "reconciliation-receipt")
	workspace.SetString(receipt, "classification", "observation")
	workspace.SetString(receipt, "claim", "authorized Contract set was reconciled atomically")
	workspace.SetStrings(receipt, "sources", receiptSources)
	workspace.SetString(receipt, "transaction_key", transactionKey)
	workspace.SetString(receipt, "idempotency_key", transactionKey)
	receiptCanonical, err := workspace.Canonical(receipt)
	if err != nil {
		return nil, err
	}
	changes = append(changes, FileChange{Path: recordPath(domain.Evidence, receiptID), Data: receiptCanonical, Mode: 0o644})
	if err := ApplyTransaction(s.Workspace.Root, transactionKey, changes); err != nil {
		return nil, err
	}
	for i := range results {
		results[i].Receipt = receiptRef
	}
	return results, nil
}

func (s Service) requireReadyAssessment(decisionRef, proposalRef string) error {
	decision, err := s.Workspace.Lookup(decisionRef, domain.Decision)
	if err != nil {
		return err
	}
	evidence, err := workspace.Strings(decision.Document, "evidence", true)
	if err != nil {
		return err
	}
	for _, ref := range evidence {
		typed, parseErr := domain.ParseReference(ref)
		if parseErr != nil || typed.Type != domain.Assessment {
			continue
		}
		assessment, lookupErr := s.Workspace.Lookup(ref, domain.Assessment)
		if lookupErr != nil || mustString(assessment.Document, "verdict") != "ready-for-owner" {
			continue
		}
		mission, lookupErr := s.Workspace.Lookup(mustString(assessment.Document, "mission"), domain.Mission)
		if lookupErr == nil && mission.Document.Record.Source != nil && mission.Document.Record.Source.String() == proposalRef {
			return nil
		}
	}
	return domain.NewRefusal(domain.RefusalInsufficientEvidence, "decision.evidence", "reconciliation Decision must cite a ready-for-owner Assessment for the Proposal Mission", nil)
}

func receiptID(key string) domain.ID {
	digest := sha256.Sum256([]byte("spectacular.reconciliation.receipt.v1\x00" + key))
	digest[6] = (digest[6] & 0x0f) | 0x70
	digest[8] = (digest[8] & 0x3f) | 0x80
	text := fmt.Sprintf("%x-%x-%x-%x-%x", digest[0:4], digest[4:6], digest[6:8], digest[8:10], digest[10:16])
	id, _ := domain.ParseID(text)
	return id
}

func (s Service) ContractView(ref string) (map[string]any, error) {
	contract, err := s.Workspace.Lookup(ref, domain.Contract)
	if err != nil {
		return nil, err
	}
	return map[string]any{"ref": refOf(contract), "path": contract.Path, "fingerprint": contract.Fingerprint, "status": value(contract.Document.Record.Status), "version": mustString(contract.Document, "contract_version"), "contract": candidateFromDocument(contract.Document), "accepted_proposal": mustString(contract.Document, "accepted_proposal")}, nil
}

func (s Service) ArchiveMission(input ArchiveInput) (OperationResult, error) {
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	archived, _ := workspace.Bool(mission.Document, "archived", false)
	if archived && mustString(mission.Document, "last_idempotency_key") == input.IdempotencyKey {
		return result("mission.archive", mission, true, []string{input.Authorization, input.TerminalPacket}), nil
	}
	if mission.Fingerprint != input.ExpectedFingerprint {
		return OperationResult{}, stale("expected_fingerprint", input.ExpectedFingerprint, mission.Fingerprint)
	}
	if value(mission.Document.Record.Status) != "resolved" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "status", "only a resolved Mission may be archived", nil)
	}
	if input.TerminalPacket != input.Mission || mustString(mission.Document, "terminal_next_action") == "" || mustString(mission.Document, "assessment") == "" {
		return OperationResult{}, missing("terminal_packet", "resolved Mission must itself contain the terminal continuity packet")
	}
	if mustString(mission.Document, "reconciliation") == "" {
		decision, lookupErr := s.Workspace.Lookup(input.Authorization, domain.Decision)
		if lookupErr != nil || !contains([]string{"no-contract-delta", "reconciliation-not-required"}, mustString(decision.Document, "disposition")) {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "archive requires settled reconciliation", lookupErr)
		}
	}
	if err := s.authorize(input.Authorization, "mission.archive", input.Mission, input.ExpectedFingerprint); err != nil {
		return OperationResult{}, err
	}
	workspace.SetBool(mission.Document, "archived", true)
	workspace.SetString(mission.Document, "archived_at", s.now().Format(time.RFC3339))
	workspace.SetString(mission.Document, "archive_authorization", input.Authorization)
	workspace.SetString(mission.Document, "archive_input_fingerprint", input.ExpectedFingerprint)
	workspace.SetString(mission.Document, "last_idempotency_key", input.IdempotencyKey)
	mission.Document.Record.Updated = stringPtr(s.now().Format(time.RFC3339))
	canonical, err := workspace.Canonical(mission.Document)
	if err != nil {
		return OperationResult{}, err
	}
	anchor := s.Workspace.ProjectAnchor()
	truth := mustStrings(anchor.Document, "current_truth")
	filtered := make([]string, 0, len(truth)+1)
	for _, ref := range truth {
		if ref != input.Mission {
			filtered = append(filtered, ref)
		}
	}
	if mission.Document.Record.Source == nil || mission.Document.Record.Source.Type != domain.Proposal {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "source", "archived Mission must retain its accepted Proposal source", nil)
	}
	proposalRef := mission.Document.Record.Source.String()
	proposal, err := s.Workspace.Lookup(proposalRef, domain.Proposal)
	if err != nil {
		return OperationResult{}, err
	}
	contractRef := mustString(proposal.Document, "target_contract")
	if contractRef != "" && !contains(filtered, contractRef) {
		filtered = append(filtered, contractRef)
	}
	workspace.SetStrings(anchor.Document, "current_truth", filtered)
	workspace.SetString(anchor.Document, "last_closed_mission", input.Mission)
	anchor.Document.Record.Updated = stringPtr(s.now().Format(time.RFC3339))
	anchorCanonical, err := workspace.Canonical(anchor.Document)
	if err != nil {
		return OperationResult{}, err
	}
	changes := []FileChange{{Path: mission.Path, Data: canonical, Mode: 0o644}, {Path: anchor.Path, Data: anchorCanonical, Mode: 0o644}}
	if err := ApplyTransaction(s.Workspace.Root, input.IdempotencyKey, changes); err != nil {
		return OperationResult{}, err
	}
	fp, _ := workspace.Fingerprint(mission.Document)
	return OperationResult{Operation: "mission.archive", Ref: input.Mission, Path: mission.Path, Fingerprint: fp, Sources: []string{input.Authorization, input.TerminalPacket}}, nil
}

func legalTransition(from, to string) bool {
	return (from == "defined" && to == "active") || (from == "active" && to == "awaiting-assessment") || (from == "awaiting-assessment" && to == "resolved")
}

func subset(values, allowed []string) bool {
	for _, value := range values {
		if !contains(allowed, value) {
			return false
		}
	}
	return true
}
