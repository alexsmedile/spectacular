package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

type authorityRecord struct {
	Operation            string
	Scope                []string
	AuthorizedEffects    []string
	Conditions           []string
	Targets              []string
	ExpectedFingerprints []string
	Evidence             []string
	Disposition          string
}

type evidenceRecord struct {
	Mission          string
	Claim            string
	Classification   string
	Contrary         []string
	RequiredChecks   []string
	CheckResults     []string
	ReviewState      string
	ExecutorAuthored bool
}

type reconciliationReceipt struct {
	Kind                 string
	Operation            string
	Contracts            []string
	Proposals            []string
	Decisions            []string
	ExpectedFingerprints []string
	Missions             []string
	OperationDigest      string
}

type missionEnvelope struct {
	Scope            []string
	AllowedActions   []string
	ForbiddenEffects []string
	EvidenceClaims   []string
	Stops            []string
	BudgetUnits      int
	ExpiresAt        string
}

func parseMissionEnvelope(entry discovery.Entry) (missionEnvelope, error) {
	scope, err := workspace.Strings(entry.Document, "scope", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	allowed, err := workspace.Strings(entry.Document, "allowed_actions", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	forbidden, err := workspace.Strings(entry.Document, "forbidden_effects", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	claims, err := workspace.Strings(entry.Document, "evidence_claims", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	stops, err := workspace.Strings(entry.Document, "stops", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	budget, err := workspace.Int(entry.Document, "budget_units", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	expires, err := workspace.String(entry.Document, "expires_at", true)
	if err != nil {
		return missionEnvelope{}, err
	}
	return missionEnvelope{Scope: scope, AllowedActions: allowed, ForbiddenEffects: forbidden, EvidenceClaims: claims, Stops: stops, BudgetUnits: budget, ExpiresAt: expires}, nil
}

type handoffEnvelope struct {
	Mission                    string
	Objective                  string
	Run                        string
	Scope                      []string
	Inputs                     []string
	AllowedActions             []string
	ForbiddenEffects           []string
	EvidenceClaims             []string
	Stops                      []string
	BudgetUnits                int
	ExpiresAt                  string
	Authorization              string
	ExpectedMissionFingerprint string
}

type terminalPacket struct {
	Assessment         string
	ResolutionDecision string
	TerminalNextAction string
	Disposition        string
	Reconciliation     string
	Scope              []string
	Objectives         []string
}

func parseTerminalPacket(mission discovery.Entry) (terminalPacket, error) {
	assessment, err := workspace.String(mission.Document, "assessment", true)
	if err != nil {
		return terminalPacket{}, err
	}
	decision, err := workspace.String(mission.Document, "last_authorization", true)
	if err != nil {
		return terminalPacket{}, err
	}
	next, err := workspace.String(mission.Document, "terminal_next_action", true)
	if err != nil {
		return terminalPacket{}, err
	}
	disposition, err := workspace.String(mission.Document, "disposition", true)
	if err != nil {
		return terminalPacket{}, err
	}
	reconciliation, err := workspace.String(mission.Document, "reconciliation", false)
	if err != nil {
		return terminalPacket{}, err
	}
	scope, err := workspace.Strings(mission.Document, "scope", true)
	if err != nil {
		return terminalPacket{}, err
	}
	objectives, err := workspace.Strings(mission.Document, "objectives", false)
	if err != nil {
		return terminalPacket{}, err
	}
	return terminalPacket{Assessment: assessment, ResolutionDecision: decision, TerminalNextAction: next, Disposition: disposition, Reconciliation: reconciliation, Scope: scope, Objectives: objectives}, nil
}

func parseHandoffEnvelope(entry discovery.Entry) (handoffEnvelope, error) {
	mission, err := workspace.String(entry.Document, "mission", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	objective, err := workspace.String(entry.Document, "objective", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	run, err := workspace.String(entry.Document, "run", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	scope, err := workspace.Strings(entry.Document, "scope", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	inputs, err := workspace.Strings(entry.Document, "inputs", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	allowed, err := workspace.Strings(entry.Document, "allowed_actions", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	forbidden, err := workspace.Strings(entry.Document, "forbidden_effects", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	claims, err := workspace.Strings(entry.Document, "evidence_claims", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	stops, err := workspace.Strings(entry.Document, "stops", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	budget, err := workspace.Int(entry.Document, "budget_units", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	expires, err := workspace.String(entry.Document, "expires_at", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	authorization, err := workspace.String(entry.Document, "authorization", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	expected, err := workspace.String(entry.Document, "expected_mission_fingerprint", true)
	if err != nil {
		return handoffEnvelope{}, err
	}
	return handoffEnvelope{Mission: mission, Objective: objective, Run: run, Scope: scope, Inputs: inputs, AllowedActions: allowed, ForbiddenEffects: forbidden, EvidenceClaims: claims, Stops: stops, BudgetUnits: budget, ExpiresAt: expires, Authorization: authorization, ExpectedMissionFingerprint: expected}, nil
}

func parseReconciliationReceipt(entry discovery.Entry) (reconciliationReceipt, error) {
	kind, err := workspace.String(entry.Document, "kind", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	operation, err := workspace.String(entry.Document, "operation", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	contracts, err := workspace.Strings(entry.Document, "contracts", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	proposals, err := workspace.Strings(entry.Document, "proposals", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	decisions, err := workspace.Strings(entry.Document, "decisions", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	expected, err := workspace.Strings(entry.Document, "expected_fingerprints", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	missions, err := workspace.Strings(entry.Document, "missions", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	digest, err := workspace.String(entry.Document, "operation_digest", true)
	if err != nil {
		return reconciliationReceipt{}, err
	}
	return reconciliationReceipt{Kind: kind, Operation: operation, Contracts: contracts, Proposals: proposals, Decisions: decisions, ExpectedFingerprints: expected, Missions: missions, OperationDigest: digest}, nil
}

func parseEvidence(entry discovery.Entry) (evidenceRecord, error) {
	mission, err := workspace.String(entry.Document, "mission", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	claim, err := workspace.String(entry.Document, "claim", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	classification, err := workspace.String(entry.Document, "classification", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	contrary, err := workspace.Strings(entry.Document, "contrary_evidence", false)
	if err != nil {
		return evidenceRecord{}, err
	}
	checks, err := workspace.Strings(entry.Document, "required_checks", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	results, err := workspace.Strings(entry.Document, "check_results", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	review, err := workspace.String(entry.Document, "review_state", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	executor, err := workspace.Bool(entry.Document, "executor_authored", true)
	if err != nil {
		return evidenceRecord{}, err
	}
	return evidenceRecord{Mission: mission, Claim: claim, Classification: classification, Contrary: contrary, RequiredChecks: checks, CheckResults: results, ReviewState: review, ExecutorAuthored: executor}, nil
}

func (s Service) authority(entry discovery.Entry) (authorityRecord, error) {
	operation, err := workspace.String(entry.Document, "operation", true)
	if err != nil {
		return authorityRecord{}, err
	}
	scope, err := workspace.Strings(entry.Document, "scope", true)
	if err != nil || len(scope) == 0 {
		return authorityRecord{}, domain.NewRefusal(domain.RefusalUnauthorized, "scope", "Decision requires non-empty mechanical scope", err)
	}
	effects, err := workspace.Strings(entry.Document, "authorized_effects", true)
	if err != nil || len(effects) == 0 {
		return authorityRecord{}, domain.NewRefusal(domain.RefusalUnauthorized, "authorized_effects", "Decision requires explicit authorized effects", err)
	}
	conditions, err := workspace.Strings(entry.Document, "conditions", false)
	if err != nil {
		return authorityRecord{}, err
	}
	targets, err := workspace.Strings(entry.Document, "targets", true)
	if err != nil {
		return authorityRecord{}, err
	}
	fingerprints, err := workspace.Strings(entry.Document, "expected_fingerprints", true)
	if err != nil || len(fingerprints) != len(targets) {
		return authorityRecord{}, domain.NewRefusal(domain.RefusalUnauthorized, "expected_fingerprints", "Decision target/fingerprint cardinality mismatch", err)
	}
	evidence, err := workspace.Strings(entry.Document, "evidence", false)
	if err != nil {
		return authorityRecord{}, err
	}
	disposition, err := workspace.String(entry.Document, "disposition", true)
	if err != nil {
		return authorityRecord{}, err
	}
	return authorityRecord{Operation: operation, Scope: scope, AuthorizedEffects: effects, Conditions: conditions, Targets: targets, ExpectedFingerprints: fingerprints, Evidence: evidence, Disposition: disposition}, nil
}

func (s Service) evaluateConditions(authority authorityRecord, expected string) error {
	for _, condition := range authority.Conditions {
		switch {
		case condition == "no-provider-effects" || condition == "no provider effects":
		case condition == "target-absent":
			if expected != "absent" {
				return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "target-absent condition is not satisfied", nil)
			}
		case condition == "target-current":
			if expected == "absent" {
				return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "target-current condition is not satisfied", nil)
			}
		case strings.HasPrefix(condition, "not-before:"):
			when, err := time.Parse(time.RFC3339, strings.TrimPrefix(condition, "not-before:"))
			if err != nil || s.now().Before(when) {
				return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "not-before condition is not satisfied", err)
			}
		case strings.HasPrefix(condition, "requires-evidence:"):
			ref := strings.TrimPrefix(condition, "requires-evidence:")
			typed, err := domain.ParseReference(ref)
			if err != nil {
				return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "requires-evidence condition has invalid reference", err)
			}
			if _, err := s.Workspace.Lookup(ref, typed.Type); err != nil {
				return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "requires-evidence condition is not satisfied", err)
			}
		default:
			return domain.NewRefusal(domain.RefusalUnauthorized, "conditions", "Decision condition is not mechanically evaluable: "+condition, nil)
		}
	}
	return nil
}

func (s Service) validateFresh(entry discovery.Entry, field string) error {
	checkedText, err := workspace.String(entry.Document, "freshness_checked_at", true)
	if err != nil {
		return err
	}
	validText, err := workspace.String(entry.Document, "freshness_valid_until", true)
	if err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339, checkedText); err != nil {
		return domain.NewRefusal(domain.RefusalInvalidKnownField, "freshness_checked_at", "must be RFC3339", err)
	}
	validUntil, err := time.Parse(time.RFC3339, validText)
	if err != nil || !validUntil.After(s.now()) {
		return domain.NewRefusal(domain.RefusalInsufficientEvidence, field, "canonical Evidence or authority is stale", err)
	}
	source, err := workspace.String(entry.Document, "freshness_source", true)
	if err != nil {
		return err
	}
	expected, err := workspace.String(entry.Document, "freshness_source_fingerprint", true)
	if err != nil {
		return err
	}
	_, actual, err := s.Workspace.Source(source)
	if err != nil {
		return err
	}
	if expected != actual {
		return domain.NewStateRefusal(domain.RefusalInsufficientEvidence, field, "freshness source fingerprint changed", expected, actual, "refresh evidence and repeat independent assessment", nil)
	}
	return nil
}

func (s Service) validateBoundInputs(inputs []string) error {
	if len(inputs) == 0 {
		return missing("inputs", "Handoff requires at least one authoritative ref@fingerprint input")
	}
	for _, binding := range inputs {
		index := strings.LastIndex(binding, "@")
		if index <= 0 || index == len(binding)-1 {
			return invalid("inputs", "authoritative inputs must use ref@fingerprint")
		}
		ref, expected := binding[:index], binding[index+1:]
		_, actual, err := s.Workspace.Source(ref)
		if err != nil {
			return err
		}
		if actual != expected {
			return stale("inputs", expected, actual)
		}
	}
	return nil
}

func canonicalReconcileInputs(inputs []ReconcileInput) ([]ReconcileInput, string, error) {
	canonical := append([]ReconcileInput(nil), inputs...)
	sort.Slice(canonical, func(i, j int) bool { return canonical[i].Contract < canonical[j].Contract })
	for i, item := range canonical {
		if item.Contract == "" || item.Proposal == "" || item.Authorization == "" || item.ExpectedFingerprint == "" || item.IdempotencyKey == "" {
			return nil, "", missing("reconciliation", "every set item requires Contract, Proposal, Decision, expected fingerprint, and idempotency key")
		}
		if i > 0 && item.Contract == canonical[i-1].Contract {
			return nil, "", domain.NewRefusal(domain.RefusalCollision, "contract", "reconciliation set contains duplicate Contract", nil)
		}
	}
	data, err := json.Marshal(canonical)
	if err != nil {
		return nil, "", err
	}
	digest := sha256.Sum256(append([]byte("spectacular.contract.reconcile-set.v1\x00"), data...))
	return canonical, hex.EncodeToString(digest[:]), nil
}

func containsAll(available, required []string) bool {
	for _, value := range required {
		if !contains(available, value) && !contains(available, "*") {
			return false
		}
	}
	return true
}

func (s Service) validateRepairAttempts(mission discovery.Entry, attempts []RepairAttempt, excludeAssessment string) error {
	budget, err := workspace.Int(mission.Document, "repair_budget", false)
	if err != nil {
		return err
	}
	claims, err := workspace.Strings(mission.Document, "evidence_claims", true)
	if err != nil {
		return err
	}
	consumed := 0
	for _, assessment := range s.Workspace.OfType(domain.Assessment) {
		if refOf(assessment) == excludeAssessment {
			continue
		}
		missionRef, parseErr := workspace.String(assessment.Document, "mission", true)
		if parseErr != nil {
			return parseErr
		}
		if missionRef != refOf(mission) {
			continue
		}
		encoded, parseErr := workspace.String(assessment.Document, "repair_attempts_json", false)
		if parseErr != nil {
			return parseErr
		}
		if encoded == "" {
			continue
		}
		var prior []RepairAttempt
		decoder := json.NewDecoder(strings.NewReader(encoded))
		decoder.DisallowUnknownFields()
		if parseErr := decoder.Decode(&prior); parseErr != nil {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts_json", "decode prior typed repair attempts", parseErr)
		}
		if parseErr := decoder.Decode(&struct{}{}); parseErr != io.EOF {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts_json", "prior repair attempts must contain exactly one JSON value", parseErr)
		}
		for _, attempt := range prior {
			if attempt.BudgetConsumed < 1 || len(attempt.AffectedClaims) == 0 || !subset(attempt.AffectedClaims, claims) || (attempt.NewHypothesis == attempt.PreviousHypothesis && len(attempt.NewEvidence) == 0 && attempt.NarrowerAction == "") {
				return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts_json", "prior repair attempt violates typed repair invariants", nil)
			}
			consumed += attempt.BudgetConsumed
		}
	}
	for index, attempt := range attempts {
		if len(attempt.AffectedClaims) == 0 || !subset(attempt.AffectedClaims, claims) || attempt.Actor == "" || len(attempt.Checks) == 0 || attempt.Result == "" || attempt.BudgetConsumed < 1 || attempt.RecoveryPoint == "" {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts", "repair attempt is incomplete or affects undeclared claims", nil)
		}
		changedHypothesis := attempt.NewHypothesis != "" && attempt.NewHypothesis != attempt.PreviousHypothesis
		if !changedHypothesis && len(attempt.NewEvidence) == 0 && attempt.NarrowerAction == "" {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts", "repair attempt must change hypothesis/evidence or narrow the action", nil)
		}
		for _, refs := range [][]string{attempt.BeforeEvidence, attempt.AfterEvidence, attempt.NewEvidence} {
			for _, ref := range refs {
				evidence, lookupErr := s.Workspace.Lookup(ref, domain.Evidence)
				if lookupErr != nil {
					return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts", "repair Evidence must resolve within the same Mission", lookupErr)
				}
				evidenceMission, parseErr := workspace.String(evidence.Document, "mission", true)
				if parseErr != nil || evidenceMission != refOf(mission) {
					return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts", "repair Evidence must resolve within the same Mission", parseErr)
				}
			}
		}
		consumed += attempt.BudgetConsumed
		if consumed > budget {
			return domain.NewRefusal(domain.RefusalUnauthorized, "repair_budget", "repair attempt budget exceeds the Mission declaration at item "+strconv.Itoa(index), nil)
		}
	}
	return nil
}

func (s Service) validateReadyAssessment(assessment, mission discovery.Entry) error {
	assessmentMission, err := workspace.String(assessment.Document, "mission", true)
	if err != nil {
		return err
	}
	verdict, err := workspace.String(assessment.Document, "verdict", true)
	if err != nil {
		return err
	}
	if assessmentMission != refOf(mission) || verdict != "ready-for-owner" {
		return domain.NewRefusal(domain.RefusalInsufficientEvidence, "assessment", "ready-for-owner Assessment must belong to the exact Mission", nil)
	}
	if err := s.validateFresh(assessment, "assessment"); err != nil {
		return err
	}
	claims, err := workspace.Strings(assessment.Document, "claims", true)
	if err != nil {
		return err
	}
	evidence, err := workspace.Strings(assessment.Document, "evidence", true)
	if err != nil {
		return err
	}
	materialClaims, err := workspace.Strings(mission.Document, "evidence_claims", true)
	if err != nil {
		return err
	}
	for _, claim := range materialClaims {
		if !contains(claims, claim) {
			return domain.NewRefusal(domain.RefusalInsufficientEvidence, "claims", "Assessment omitted material claim: "+claim, nil)
		}
		if err := s.sufficientClaim(refOf(mission), claim, evidence); err != nil {
			return err
		}
	}
	repairText, err := workspace.String(assessment.Document, "repair_attempts_json", false)
	if err != nil {
		return err
	}
	if repairText != "" {
		var attempts []RepairAttempt
		decoder := json.NewDecoder(strings.NewReader(repairText))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&attempts); err != nil {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts_json", "decode typed repair attempts", err)
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, "repair_attempts_json", "repair attempts must contain exactly one JSON value", err)
		}
		if err := s.validateRepairAttempts(mission, attempts, refOf(assessment)); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) validateReconciliationReceipt(ref string, mission, resolutionDecision discovery.Entry) error {
	receipt, err := s.Workspace.Lookup(ref, domain.Evidence)
	if err != nil {
		return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "reconciliation receipt does not resolve", err)
	}
	parsed, err := parseReconciliationReceipt(receipt)
	if err != nil {
		return err
	}
	if parsed.Kind != "reconciliation-receipt" || parsed.Operation != "contract.reconcile-set" {
		return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "Evidence is not an atomic reconciliation receipt", nil)
	}
	if err := s.validateFresh(receipt, "reconciliation_receipt"); err != nil {
		return err
	}
	if mission.Document.Record.Source == nil || !contains(parsed.Missions, refOf(mission)) || !contains(parsed.Proposals, mission.Document.Record.Source.String()) {
		return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "receipt is not bound to the exact Mission and Proposal", nil)
	}
	resolutionAuthority, err := s.authority(resolutionDecision)
	if err != nil {
		return err
	}
	if !contains(resolutionAuthority.Evidence, ref) {
		return domain.NewRefusal(domain.RefusalUnauthorized, "decision.evidence", "resolution Decision must cite the exact reconciliation receipt", nil)
	}
	contracts := parsed.Contracts
	proposals := parsed.Proposals
	decisions := parsed.Decisions
	expected := parsed.ExpectedFingerprints
	if len(contracts) == 0 || len(contracts) != len(proposals) || len(contracts) != len(decisions) || len(contracts) != len(expected) {
		return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "receipt Contract/Proposal/Decision/fingerprint cardinality is invalid", nil)
	}
	for index, contractRef := range contracts {
		contract, lookupErr := s.Workspace.Lookup(contractRef, domain.Contract)
		if lookupErr != nil {
			return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "current Contract does not match receipt provenance", lookupErr)
		}
		acceptedProposal, parseErr := workspace.String(contract.Document, "accepted_proposal", true)
		if parseErr != nil {
			return parseErr
		}
		contractDecision, parseErr := workspace.String(contract.Document, "authorization", true)
		if parseErr != nil || acceptedProposal != proposals[index] || contractDecision != decisions[index] {
			return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "current Contract does not match receipt provenance", parseErr)
		}
		decision, lookupErr := s.Workspace.Lookup(decisions[index], domain.Decision)
		if lookupErr != nil {
			return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "receipt reconciliation Decision does not resolve", lookupErr)
		}
		operation, parseErr := workspace.String(decision.Document, "operation", true)
		if parseErr != nil || operation != "contract.reconcile" {
			return domain.NewRefusal(domain.RefusalUnsettledReconcile, "reconciliation", "receipt reconciliation Decision does not resolve", parseErr)
		}
	}
	return nil
}
