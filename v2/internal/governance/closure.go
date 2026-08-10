package governance

import (
	"crypto/sha256"
	"encoding/json"
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
	missionAuthority, err := parseMissionEnvelope(mission)
	if err != nil {
		return OperationResult{}, err
	}
	objective, err := s.Workspace.Lookup(input.Objective, domain.Objective)
	if err != nil {
		return OperationResult{}, err
	}
	run, err := s.Workspace.Lookup(input.Run, domain.Run)
	if err != nil {
		return OperationResult{}, err
	}
	objectiveMission, err := workspace.String(objective.Document, "mission", true)
	if err != nil {
		return OperationResult{}, err
	}
	runMission, err := workspace.String(run.Document, "mission", true)
	if err != nil {
		return OperationResult{}, err
	}
	if objectiveMission != input.Mission || runMission != input.Mission {
		return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "mission", "Handoff Objective and Run must belong to the same Mission", nil)
	}
	if input.Sender == "" || input.Actor == "" || input.Destination == "" || len(input.Scope) == 0 || len(input.EvidenceClaims) == 0 || input.BudgetUnits < 1 || input.RecoveryPoint == "" || input.ReturnDestination == "" || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("handoff", "actor, destination, scope, authority, evidence, budget, recovery, return, and idempotency are required")
	}
	if _, err := parseFuture(missionAuthority.ExpiresAt, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(input.ExpiresAt, s.now(), "expires_at"); err != nil {
		return OperationResult{}, err
	}
	if input.BudgetUnits > missionAuthority.BudgetUnits || !subset(input.Scope, missionAuthority.Scope) || !subset(input.AllowedActions, missionAuthority.AllowedActions) || !subset(input.EvidenceClaims, missionAuthority.EvidenceClaims) || !subset(input.Stops, missionAuthority.Stops) || !subset(missionAuthority.ForbiddenEffects, input.ForbiddenEffects) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "envelope", "Handoff must be a subset of Mission scope/actions and preserve every forbidden effect", nil)
	}
	if err := s.validateBoundInputs(input.Inputs); err != nil {
		return OperationResult{}, err
	}
	ref := string(domain.Handoff) + ":" + id.String()
	if err := s.authorize(input.Authorization, "handoff.create", ref, "absent", input.Scope, []string{"handoff.create"}); err != nil {
		return OperationResult{}, err
	}
	if input.Supersedes != "" {
		old, lookupErr := s.Workspace.Lookup(input.Supersedes, domain.Handoff)
		if lookupErr != nil {
			return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "supersedes", "superseded Handoff must be an existing dispatch in the same Mission", lookupErr)
		}
		oldKind, parseErr := workspace.String(old.Document, "kind", true)
		if parseErr != nil {
			return OperationResult{}, parseErr
		}
		oldMission, parseErr := workspace.String(old.Document, "mission", true)
		if parseErr != nil || oldKind != "dispatch" || oldMission != input.Mission {
			return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "supersedes", "superseded Handoff must be an existing dispatch in the same Mission", parseErr)
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
	workspace.SetInt(doc, "budget_units", input.BudgetUnits)
	workspace.SetString(doc, "expires_at", input.ExpiresAt)
	workspace.SetStrings(doc, "stops", input.Stops)
	workspace.SetString(doc, "recovery_point", input.RecoveryPoint)
	workspace.SetString(doc, "return_destination", input.ReturnDestination)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "expected_mission_fingerprint", input.ExpectedMissionFingerprint)
	if input.Supersedes != "" {
		workspace.SetString(doc, "supersedes", input.Supersedes)
	}
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("handoff.create", doc, input.IdempotencyKey, []string{input.Mission, input.Objective, input.Run, input.Authorization})
}

func (s Service) ValidateHandoff(ref string) (map[string]any, error) {
	handoff, err := s.Workspace.Lookup(ref, domain.Handoff)
	if err != nil {
		return nil, err
	}
	kind, err := workspace.String(handoff.Document, "kind", true)
	if err != nil {
		return nil, err
	}
	if kind != "dispatch" {
		return nil, invalid("kind", "handoff validate requires a dispatch")
	}
	for _, candidate := range s.Workspace.OfType(domain.Handoff) {
		supersedes, parseErr := workspace.String(candidate.Document, "supersedes", false)
		if parseErr != nil {
			return nil, parseErr
		}
		if supersedes == refOf(handoff) {
			return nil, domain.NewRefusal(domain.RefusalConflictingAuthority, "supersedes", "Handoff was superseded by "+refOf(candidate), nil)
		}
	}
	handoffAuthority, err := parseHandoffEnvelope(handoff)
	if err != nil {
		return nil, err
	}
	missionRef := handoffAuthority.Mission
	mission, err := s.Workspace.Lookup(missionRef, domain.Mission)
	if err != nil {
		return nil, err
	}
	missionAuthority, err := parseMissionEnvelope(mission)
	if err != nil {
		return nil, err
	}
	objective, err := s.Workspace.Lookup(handoffAuthority.Objective, domain.Objective)
	if err != nil {
		return nil, err
	}
	run, err := s.Workspace.Lookup(handoffAuthority.Run, domain.Run)
	if err != nil {
		return nil, err
	}
	objectiveMission, err := workspace.String(objective.Document, "mission", true)
	if err != nil {
		return nil, err
	}
	runMission, err := workspace.String(run.Document, "mission", true)
	if err != nil {
		return nil, err
	}
	if objectiveMission != missionRef || runMission != missionRef {
		return nil, domain.NewRefusal(domain.RefusalConflictingAuthority, "mission", "Handoff containment no longer holds", nil)
	}
	if _, err := parseFuture(missionAuthority.ExpiresAt, s.now(), "mission.expires_at"); err != nil {
		return nil, err
	}
	if handoffAuthority.ExpectedMissionFingerprint != mission.Fingerprint {
		return nil, stale("expected_mission_fingerprint", handoffAuthority.ExpectedMissionFingerprint, mission.Fingerprint)
	}
	if handoffAuthority.BudgetUnits > missionAuthority.BudgetUnits || !subset(handoffAuthority.Scope, missionAuthority.Scope) || !subset(handoffAuthority.AllowedActions, missionAuthority.AllowedActions) || !subset(handoffAuthority.EvidenceClaims, missionAuthority.EvidenceClaims) || !subset(handoffAuthority.Stops, missionAuthority.Stops) || !subset(missionAuthority.ForbiddenEffects, handoffAuthority.ForbiddenEffects) {
		return nil, domain.NewRefusal(domain.RefusalUnauthorized, "envelope", "Handoff no longer fits current Mission envelope", nil)
	}
	if err := s.validateBoundInputs(handoffAuthority.Inputs); err != nil {
		return nil, err
	}
	if _, err := parseFuture(handoffAuthority.ExpiresAt, s.now(), "expires_at"); err != nil {
		return nil, err
	}
	if err := s.authorize(handoffAuthority.Authorization, "handoff.create", refOf(handoff), "absent", handoffAuthority.Scope, []string{"handoff.create"}); err != nil {
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
	dispatchAuthority, err := parseHandoffEnvelope(dispatch)
	if err != nil {
		return OperationResult{}, err
	}
	if input.Status != "succeeded" && input.Status != "blocked" && input.Status != "failed" {
		return OperationResult{}, invalid("status", "return status must be succeeded, blocked, or failed")
	}
	if input.Actor == "" || input.FinalBaseline == "" || input.Result == "" || input.BudgetUsed < 0 || input.RecoveryPoint == "" || input.IdempotencyKey == "" || (input.NextAction == "") == (input.OwnerGate == "") {
		return OperationResult{}, missing("return", "actor, result, baseline, budget, recovery, idempotency, and exactly one next_action or owner_gate are required")
	}
	if !subset(input.Actions, dispatchAuthority.AllowedActions) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "actions", "Handoff return actions exceed dispatch authority", nil)
	}
	if input.BudgetUsed > dispatchAuthority.BudgetUnits {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "budget_used", "Handoff return exceeded dispatch budget", nil)
	}
	doc := s.document(domain.Handoff, id, input.Title, input.Actor, "")
	workspace.SetString(doc, "kind", "return")
	workspace.SetString(doc, "dispatch", input.Dispatch)
	workspace.SetString(doc, "mission", dispatchAuthority.Mission)
	workspace.SetString(doc, "objective", dispatchAuthority.Objective)
	workspace.SetString(doc, "run", dispatchAuthority.Run)
	workspace.SetString(doc, "return_status", input.Status)
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetString(doc, "final_baseline", input.FinalBaseline)
	workspace.SetString(doc, "result", input.Result)
	workspace.SetStrings(doc, "actions", input.Actions)
	workspace.SetStrings(doc, "provider_receipts", input.ProviderReceipts)
	workspace.SetStrings(doc, "evidence", input.Evidence)
	workspace.SetStrings(doc, "remaining_gaps", input.RemainingGaps)
	workspace.SetInt(doc, "budget_used", input.BudgetUsed)
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
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := time.Parse(time.RFC3339, input.ObservedAt); err != nil {
		return OperationResult{}, invalid("observed_at", "must be RFC3339")
	}
	if _, err := parseFuture(input.FreshnessValidUntil, s.now(), "freshness_valid_until"); err != nil {
		return OperationResult{}, err
	}
	ref := string(domain.Evidence) + ":" + id.String()
	missionScope, err := workspace.Strings(mission.Document, "scope", true)
	if err != nil {
		return OperationResult{}, err
	}
	missionExpiry, err := workspace.String(mission.Document, "expires_at", true)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(missionExpiry, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	if err := s.authorize(input.Authorization, "evidence.create", ref, "absent", missionScope, []string{"evidence.create"}); err != nil {
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
	missionExpiry, err := workspace.String(mission.Document, "expires_at", true)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(missionExpiry, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	if input.Actor == "" || len(input.Claims) == 0 || len(input.Evidence) == 0 || input.RecoveryPoint == "" || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("assessment", "actor, claims, Evidence, recovery, authorization, and idempotency are required")
	}
	if input.Verdict == "ready-for-owner" {
		if len(input.BlockingFindings) > 0 {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "blocking_findings", "owner-ready assessment cannot retain blocking findings", nil)
		}
		declared, err := workspace.Strings(mission.Document, "evidence_claims", true)
		if err != nil {
			return OperationResult{}, err
		}
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
	if err := s.validateRepairAttempts(mission, input.RepairAttempts, ""); err != nil {
		return OperationResult{}, err
	}
	missionScope, err := workspace.Strings(mission.Document, "scope", true)
	if err != nil {
		return OperationResult{}, err
	}
	if err := s.authorize(input.Authorization, "assessment.record", ref, "absent", missionScope, []string{"assessment.record"}); err != nil {
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
	repairs, err := json.Marshal(input.RepairAttempts)
	if err != nil {
		return OperationResult{}, err
	}
	workspace.SetString(doc, "repair_attempts_json", string(repairs))
	workspace.SetString(doc, "recovery_point", input.RecoveryPoint)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("assessment.record", doc, input.IdempotencyKey, append([]string{input.Mission, input.Authorization}, input.Evidence...))
}

func (s Service) sufficientClaim(missionRef, claim string, refs []string) error {
	var supporting bool
	referenced := map[string]bool{}
	for _, ref := range refs {
		referenced[ref] = true
	}
	for _, entry := range s.Workspace.OfType(domain.Evidence) {
		kind, err := workspace.String(entry.Document, "kind", false)
		if err != nil {
			return err
		}
		if kind == "reconciliation-receipt" {
			continue
		}
		entryMission, err := workspace.String(entry.Document, "mission", true)
		if err != nil {
			return err
		}
		entryClaim, err := workspace.String(entry.Document, "claim", true)
		if err != nil {
			return err
		}
		if entryMission != missionRef || entryClaim != claim {
			continue
		}
		parsed, err := parseEvidence(entry)
		if err != nil {
			return err
		}
		if err := s.validateFresh(entry, "evidence"); err != nil {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "evidence", "claim has stale canonical Evidence: "+claim, err)
		}
		if parsed.Classification == "unknown" || len(parsed.Contrary) > 0 {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "evidence", "claim has unknown or conflicting evidence: "+claim, nil)
		}
		checks := parsed.RequiredChecks
		results := parsed.CheckResults
		if len(checks) == 0 || len(checks) != len(results) {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "check_results", "required checks are incomplete for claim: "+claim, nil)
		}
		for _, result := range results {
			if !strings.HasPrefix(result, "pass:") {
				return domain.NewRefusal(domain.RefusalInsufficientEvidence, "check_results", "required check did not pass for claim: "+claim, nil)
			}
		}
		if parsed.ExecutorAuthored && parsed.ReviewState != "independent-accepted" {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "review_state", "executor-authored evidence requires independent acceptance for claim: "+claim, nil)
		}
		if referenced[refOf(entry)] {
			supporting = true
		}
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
	lastKey, err := workspace.String(mission.Document, "last_idempotency_key", false)
	if err != nil {
		return OperationResult{}, err
	}
	if current == input.To && lastKey == input.IdempotencyKey {
		return result("mission.transition", mission, true, []string{input.Authorization}), nil
	}
	if mission.Fingerprint != input.ExpectedFingerprint {
		return OperationResult{}, stale("expected_fingerprint", input.ExpectedFingerprint, mission.Fingerprint)
	}
	if !legalTransition(current, input.To) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "status", current+" cannot transition to "+input.To, nil)
	}
	missionExpiry, err := workspace.String(mission.Document, "expires_at", true)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(missionExpiry, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	missionScope, err := workspace.Strings(mission.Document, "scope", true)
	if err != nil {
		return OperationResult{}, err
	}
	operation := "mission.transition." + input.To
	if err := s.authorize(input.Authorization, operation, input.Mission, input.ExpectedFingerprint, missionScope, []string{operation}); err != nil {
		return OperationResult{}, err
	}
	decision, err := s.Workspace.Lookup(input.Authorization, domain.Decision)
	if err != nil {
		return OperationResult{}, err
	}
	authority, err := s.authority(decision)
	if err != nil {
		return OperationResult{}, err
	}
	if input.To == "awaiting-assessment" {
		hasProof, proofErr := s.hasReturnAndEvidence(input.Mission)
		if proofErr != nil {
			return OperationResult{}, proofErr
		}
		if !hasProof {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "mission", "awaiting-assessment requires an immutable Handoff return and returned Evidence", nil)
		}
	}
	var objectiveChanges []FileChange
	if input.To == "resolved" {
		assessment, lookupErr := s.Workspace.Lookup(input.Assessment, domain.Assessment)
		if lookupErr != nil {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "assessment", "resolution requires a ready-for-owner Assessment for this Mission", lookupErr)
		}
		if err := s.validateReadyAssessment(assessment, mission); err != nil {
			return OperationResult{}, err
		}
		if !contains([]string{"completed", "abandoned", "superseded"}, input.Disposition) {
			return OperationResult{}, invalid("disposition", "resolved disposition must be completed, abandoned, or superseded")
		}
		if authority.Disposition != input.Disposition {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "disposition", "resolution disposition must exactly match the owner Decision", nil)
		}
		if !contains(authority.AuthorizedEffects, "terminal-next-action:"+input.TerminalNextAction) {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "terminal_next_action", "terminal continuation is not mechanically authorized by the resolution Decision", nil)
		}
		if input.Reconciliation == "" && !contains(authority.AuthorizedEffects, "reconciliation.not-required") {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "resolution requires reconciliation or explicit zero-delta disposition", nil)
		}
		if input.Reconciliation != "" {
			if err := s.validateReconciliationReceipt(input.Reconciliation, mission, decision); err != nil {
				return OperationResult{}, err
			}
		}
		if input.TerminalNextAction == "" {
			return OperationResult{}, missing("terminal_next_action", "resolution requires exactly one terminal continuation or owner gate")
		}
		objectives, err := workspace.Strings(mission.Document, "objectives", false)
		if err != nil {
			return OperationResult{}, err
		}
		if !subset(input.SatisfiedObjectives, objectives) {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "satisfied_objectives", "resolution names an Objective outside the Mission", nil)
		}
		for _, ref := range objectives {
			objective, lookupErr := s.Workspace.Lookup(ref, domain.Objective)
			if lookupErr != nil {
				return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "objectives", "Mission Objective containment is invalid", lookupErr)
			}
			objectiveMission, parseErr := workspace.String(objective.Document, "mission", true)
			if parseErr != nil || objectiveMission != input.Mission {
				return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "objectives", "Mission Objective containment is invalid", parseErr)
			}
			status := value(objective.Document.Record.Status)
			if input.Disposition == "completed" && status != "satisfied" {
				if !contains(input.SatisfiedObjectives, ref) || !contains(authority.AuthorizedEffects, "objective.satisfy:"+ref) {
					return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "objectives", "completed resolution requires every Objective satisfied by explicit owner authority", nil)
				}
				objective.Document.Record.Status = stringPtr("satisfied")
				objective.Document.Record.Updated = stringPtr(s.now().Format(time.RFC3339))
				workspace.SetString(objective.Document, "satisfaction_decision", input.Authorization)
				data, canonicalErr := workspace.Canonical(objective.Document)
				if canonicalErr != nil {
					return OperationResult{}, canonicalErr
				}
				objectiveChanges = append(objectiveChanges, FileChange{Path: objective.Path, Data: data, Mode: 0o644})
			}
		}
		workspace.SetString(mission.Document, "assessment", input.Assessment)
		workspace.SetString(mission.Document, "disposition", input.Disposition)
		if input.Reconciliation != "" {
			workspace.SetString(mission.Document, "reconciliation", input.Reconciliation)
		}
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
	changes := append([]FileChange{{Path: mission.Path, Data: canonical, Mode: 0o644}}, objectiveChanges...)
	if err := ApplyTransaction(s.Workspace.Root, input.IdempotencyKey, changes); err != nil {
		return OperationResult{}, err
	}
	fp, _ := workspace.Fingerprint(mission.Document)
	return OperationResult{Operation: "mission.transition", Ref: input.Mission, Path: mission.Path, Fingerprint: fp, Sources: []string{input.Authorization}}, nil
}

func (s Service) hasReturnAndEvidence(mission string) (bool, error) {
	var returned, evidenced bool
	for _, handoff := range s.Workspace.OfType(domain.Handoff) {
		handoffMission, err := workspace.String(handoff.Document, "mission", true)
		if err != nil {
			return false, err
		}
		kind, err := workspace.String(handoff.Document, "kind", true)
		if err != nil {
			return false, err
		}
		if handoffMission == mission && kind == "return" {
			returned = true
		}
	}
	for _, evidence := range s.Workspace.OfType(domain.Evidence) {
		kind, err := workspace.String(evidence.Document, "kind", false)
		if err != nil {
			return false, err
		}
		if kind == "reconciliation-receipt" {
			continue
		}
		evidenceMission, err := workspace.String(evidence.Document, "mission", true)
		if err != nil {
			return false, err
		}
		if evidenceMission == mission {
			evidenced = true
		}
	}
	return returned && evidenced, nil
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
	canonicalInputs, operationDigest, err := canonicalReconcileInputs(inputs)
	if err != nil {
		return nil, err
	}
	inputs = canonicalInputs
	var changes []FileChange
	results := make([]OperationResult, 0, len(inputs))
	transactionKey := inputs[0].IdempotencyKey
	receiptID := receiptID(transactionKey)
	receiptRef := string(domain.Evidence) + ":" + receiptID.String()
	if receipt, found := s.entryByID(receiptID); found {
		parsedReceipt, parseErr := parseReconciliationReceipt(receipt)
		storedKey, keyErr := workspace.String(receipt.Document, "idempotency_key", true)
		if receipt.Document.Record.Type != domain.Evidence || parseErr != nil || keyErr != nil || parsedReceipt.Kind != "reconciliation-receipt" || storedKey != transactionKey {
			return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "transaction receipt identity has different content", nil)
		}
		if parsedReceipt.OperationDigest != operationDigest {
			return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "reconciliation replay changed the canonical Contract set or authority inputs", nil)
		}
		if err := s.validateFresh(receipt, "reconciliation_receipt"); err != nil {
			return nil, err
		}
		for _, input := range inputs {
			contract, err := s.Workspace.Lookup(input.Contract, domain.Contract)
			if err != nil {
				return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "reconciliation replay does not match current Contract set", err)
			}
			storedKey, parseErr := workspace.String(contract.Document, "reconciliation_idempotency_key", true)
			if parseErr != nil || storedKey != transactionKey {
				return nil, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "reconciliation replay does not match current Contract set", parseErr)
			}
			result := result("contract.reconcile", contract, true, []string{input.Proposal, input.Authorization})
			result.Receipt = receiptRef
			results = append(results, result)
		}
		return results, nil
	}
	var receiptSources []string
	var receiptContracts, receiptProposals, receiptDecisions, receiptExpected, receiptMissions []string
	for _, input := range inputs {
		if input.IdempotencyKey == "" || input.IdempotencyKey != transactionKey {
			return nil, invalid("idempotency_key", "all Contract items must share one transaction key")
		}
		contractRef, parseErr := domain.ParseReference(input.Contract)
		if parseErr != nil || contractRef.Type != domain.Contract {
			return nil, domain.NewRefusal(domain.RefusalInvalidReference, "contract", "reconciliation requires an exact Contract reference", parseErr)
		}
		proposal, err := s.Workspace.Lookup(input.Proposal, domain.Proposal)
		if err != nil {
			return nil, err
		}
		targetContract, err := workspace.String(proposal.Document, "target_contract", true)
		if err != nil {
			return nil, err
		}
		if value(proposal.Document.Record.Status) != "accepted" || targetContract != input.Contract {
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
			} else if !domain.RefusalHasCode(lookupErr, domain.RefusalRecordNotFound) {
				return nil, lookupErr
			}
		} else {
			existing, err = s.Workspace.Lookup(input.Contract, domain.Contract)
			if err != nil {
				return nil, err
			}
			baseFingerprint, parseErr := workspace.String(proposal.Document, "base_fingerprint", true)
			if parseErr != nil {
				return nil, parseErr
			}
			if existing.Fingerprint != input.ExpectedFingerprint || baseFingerprint != input.ExpectedFingerprint {
				return nil, stale("expected_fingerprint", input.ExpectedFingerprint, existing.Fingerprint)
			}
		}
		proposalScope, scopeErr := workspace.Strings(proposal.Document, "scope", true)
		if scopeErr != nil {
			return nil, scopeErr
		}
		if err := s.authorize(input.Authorization, "contract.reconcile", input.Contract, input.ExpectedFingerprint, proposalScope, []string{"contract.reconcile"}); err != nil {
			return nil, err
		}
		missionRef, err := s.requireReadyAssessment(input.Authorization, input.Proposal)
		if err != nil {
			return nil, err
		}
		candidate, err := s.candidate(proposal)
		if err != nil {
			return nil, err
		}
		version := "1"
		priorVersion := ""
		if !newCapability {
			currentText, valueErr := workspace.String(existing.Document, "contract_version", true)
			if valueErr != nil {
				return nil, valueErr
			}
			current, parseErr := strconv.Atoi(currentText)
			if parseErr != nil || current < 1 {
				return nil, invalid("contract_version", "must be a positive integer")
			}
			priorVersion = currentText
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
			old, canonicalErr := workspace.Canonical(existing.Document)
			if canonicalErr != nil {
				return nil, canonicalErr
			}
			history := filepath.ToSlash(filepath.Join(".spectacular", "history", "contracts", contractRef.ID.String()+"@"+priorVersion+".md"))
			changes = append(changes, FileChange{Path: history, Data: old, Mode: 0o644})
		}
		changes = append(changes, FileChange{Path: path, Data: canonical, Mode: 0o644})
		fp, fingerprintErr := workspace.Fingerprint(contract)
		if fingerprintErr != nil {
			return nil, fingerprintErr
		}
		results = append(results, OperationResult{Operation: "contract.reconcile", Ref: input.Contract, Path: path, Fingerprint: fp, Sources: []string{input.Proposal, input.Authorization}})
		receiptSources = append(receiptSources, input.Contract, input.Proposal, input.Authorization)
		receiptContracts = append(receiptContracts, input.Contract)
		receiptProposals = append(receiptProposals, input.Proposal)
		receiptDecisions = append(receiptDecisions, input.Authorization)
		receiptExpected = append(receiptExpected, input.ExpectedFingerprint)
		receiptMissions = append(receiptMissions, missionRef)
	}
	receipt := s.document(domain.Evidence, receiptID, "Atomic Contract reconciliation receipt", "spectacular-cli", "recorded")
	workspace.SetString(receipt, "kind", "reconciliation-receipt")
	workspace.SetString(receipt, "classification", "observation")
	workspace.SetString(receipt, "claim", "authorized Contract set was reconciled atomically")
	workspace.SetString(receipt, "operation", "contract.reconcile-set")
	workspace.SetStrings(receipt, "sources", receiptSources)
	workspace.SetStrings(receipt, "contracts", receiptContracts)
	workspace.SetStrings(receipt, "proposals", receiptProposals)
	workspace.SetStrings(receipt, "decisions", receiptDecisions)
	workspace.SetStrings(receipt, "expected_fingerprints", receiptExpected)
	workspace.SetStrings(receipt, "missions", sortedUnique(receiptMissions))
	workspace.SetString(receipt, "operation_digest", operationDigest)
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

func (s Service) requireReadyAssessment(decisionRef, proposalRef string) (string, error) {
	decision, err := s.Workspace.Lookup(decisionRef, domain.Decision)
	if err != nil {
		return "", err
	}
	evidence, err := workspace.Strings(decision.Document, "evidence", true)
	if err != nil {
		return "", err
	}
	for _, ref := range evidence {
		typed, parseErr := domain.ParseReference(ref)
		if parseErr != nil || typed.Type != domain.Assessment {
			continue
		}
		assessment, lookupErr := s.Workspace.Lookup(ref, domain.Assessment)
		if lookupErr != nil {
			continue
		}
		verdict, parseErr := workspace.String(assessment.Document, "verdict", true)
		if parseErr != nil {
			return "", parseErr
		}
		if verdict != "ready-for-owner" {
			continue
		}
		assessmentMission, parseErr := workspace.String(assessment.Document, "mission", true)
		if parseErr != nil {
			return "", parseErr
		}
		mission, lookupErr := s.Workspace.Lookup(assessmentMission, domain.Mission)
		if lookupErr == nil && mission.Document.Record.Source != nil && mission.Document.Record.Source.String() == proposalRef {
			missionExpiry, parseErr := workspace.String(mission.Document, "expires_at", true)
			if parseErr != nil {
				return "", parseErr
			}
			if _, parseErr := parseFuture(missionExpiry, s.now(), "mission.expires_at"); parseErr != nil {
				return "", parseErr
			}
			if validationErr := s.validateReadyAssessment(assessment, mission); validationErr != nil {
				return "", validationErr
			}
			return refOf(mission), nil
		}
	}
	return "", domain.NewRefusal(domain.RefusalInsufficientEvidence, "decision.evidence", "reconciliation Decision must cite a ready-for-owner Assessment for the Proposal Mission", nil)
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
	archived, err := workspace.Bool(mission.Document, "archived", false)
	if err != nil {
		return OperationResult{}, err
	}
	lastKey, err := workspace.String(mission.Document, "last_idempotency_key", false)
	if err != nil {
		return OperationResult{}, err
	}
	if archived && lastKey == input.IdempotencyKey {
		return result("mission.archive", mission, true, []string{input.Authorization, input.TerminalPacket}), nil
	}
	if mission.Fingerprint != input.ExpectedFingerprint {
		return OperationResult{}, stale("expected_fingerprint", input.ExpectedFingerprint, mission.Fingerprint)
	}
	if value(mission.Document.Record.Status) != "resolved" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "status", "only a resolved Mission may be archived", nil)
	}
	if input.TerminalPacket != input.Mission {
		return OperationResult{}, missing("terminal_packet", "resolved Mission must itself contain the terminal continuity packet")
	}
	packet, err := parseTerminalPacket(mission)
	if err != nil {
		return OperationResult{}, err
	}
	missionExpiry, err := workspace.String(mission.Document, "expires_at", true)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(missionExpiry, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	assessment, err := s.Workspace.Lookup(packet.Assessment, domain.Assessment)
	if err != nil {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInsufficientEvidence, "assessment", "archive requires a still-sufficient ready Assessment", err)
	}
	if err := s.validateReadyAssessment(assessment, mission); err != nil {
		return OperationResult{}, err
	}
	resolutionDecision, err := s.Workspace.Lookup(packet.ResolutionDecision, domain.Decision)
	if err != nil {
		return OperationResult{}, err
	}
	resolutionAuthority, err := s.authority(resolutionDecision)
	if err != nil {
		return OperationResult{}, err
	}
	terminal := packet.TerminalNextAction
	if !contains(resolutionAuthority.AuthorizedEffects, "terminal-next-action:"+terminal) || resolutionAuthority.Disposition != packet.Disposition {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "terminal_packet", "terminal packet is not derived from the resolution Decision", nil)
	}
	reconciliation := packet.Reconciliation
	if reconciliation == "" {
		if !contains(resolutionAuthority.AuthorizedEffects, "reconciliation.not-required") {
			return OperationResult{}, domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "archive requires settled reconciliation", nil)
		}
	} else if err := s.validateReconciliationReceipt(reconciliation, mission, resolutionDecision); err != nil {
		return OperationResult{}, err
	}
	if packet.Disposition == "completed" {
		for _, ref := range packet.Objectives {
			objective, lookupErr := s.Workspace.Lookup(ref, domain.Objective)
			if lookupErr != nil || value(objective.Document.Record.Status) != "satisfied" {
				return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "objectives", "completed Mission cannot archive with pending Objectives", lookupErr)
			}
		}
	}
	if err := s.authorize(input.Authorization, "mission.archive", input.Mission, input.ExpectedFingerprint, packet.Scope, []string{"mission.archive"}); err != nil {
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
	truth, err := workspace.Strings(anchor.Document, "current_truth", true)
	if err != nil {
		return OperationResult{}, err
	}
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
	contractRef, err := workspace.String(proposal.Document, "target_contract", true)
	if err != nil {
		return OperationResult{}, err
	}
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
