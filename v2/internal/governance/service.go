package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	spectacularruntime "github.com/alexsmedile/spectacular/v2/internal/runtime"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

type Service struct {
	Workspace *discovery.Workspace
	Now       func() time.Time
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

func (s Service) CreateDecision(input DecisionInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	if input.Actor == "" || input.ActorRole != "owner" || input.AuthorityBasis == "" || input.Question == "" || input.Disposition == "" || input.Rationale == "" || input.Operation == "" || len(input.Scope) == 0 || len(input.AuthorizedEffects) == 0 || input.IdempotencyKey == "" {
		return OperationResult{}, missing("decision", "actor, owner actor_role, authority_basis, question, disposition, rationale, operation, and idempotency_key are required")
	}
	if len(input.Targets) == 0 || len(input.ExpectedFingerprints) != len(input.Targets) {
		return OperationResult{}, invalid("expected_fingerprints", "one expected fingerprint is required per target")
	}
	if input.ExpiresAt != "" {
		if _, err := parseFuture(input.ExpiresAt, s.now(), "expires_at"); err != nil {
			return OperationResult{}, err
		}
	}
	doc := s.document(domain.Decision, id, input.Title, input.Actor, "")
	workspace.SetString(doc, "actor", input.Actor)
	workspace.SetString(doc, "actor_role", input.ActorRole)
	workspace.SetString(doc, "authority_basis", input.AuthorityBasis)
	workspace.SetString(doc, "question", input.Question)
	workspace.SetStrings(doc, "scope", input.Scope)
	workspace.SetString(doc, "disposition", input.Disposition)
	workspace.SetString(doc, "rationale", input.Rationale)
	workspace.SetStrings(doc, "alternatives", input.Alternatives)
	workspace.SetStrings(doc, "targets", input.Targets)
	workspace.SetStrings(doc, "expected_fingerprints", input.ExpectedFingerprints)
	workspace.SetString(doc, "operation", input.Operation)
	workspace.SetStrings(doc, "authorized_effects", input.AuthorizedEffects)
	workspace.SetStrings(doc, "conditions", input.Conditions)
	workspace.SetString(doc, "expires_at", input.ExpiresAt)
	workspace.SetStrings(doc, "evidence", input.Evidence)
	workspace.SetString(doc, "supersedes", input.Supersedes)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	return s.createOne("decision.create", doc, input.IdempotencyKey, []string{"owner-direct:" + input.AuthorityBasis})
}

func (s Service) CreateProposal(input ProposalInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	target, err := domain.ParseReference(input.TargetContract)
	if err != nil || target.Type != domain.Contract {
		return OperationResult{}, invalid("target_contract", "must be an exact Contract reference")
	}
	if input.Title == "" || input.Actor == "" || input.Rationale == "" || len(input.Scope) == 0 || input.Authorization == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("proposal", "title, actor, rationale, scope, authorization, and idempotency_key are required")
	}
	if input.Status != "submitted" && input.Status != "accepted" {
		return OperationResult{}, invalid("status", "confirmed Proposal creation requires submitted or accepted")
	}
	if input.NewCapability {
		if input.BaseVersion != "absent" || input.BaseFingerprint != "absent" {
			return OperationResult{}, invalid("base", "new capability requires explicit absent base")
		}
		if _, lookupErr := s.Workspace.Lookup(input.TargetContract, domain.Contract); lookupErr == nil {
			return OperationResult{}, domain.NewRefusal(domain.RefusalCollision, "target_contract", "new capability target already exists", nil)
		}
		if err := validateCandidate(input.Candidate); err != nil {
			return OperationResult{}, err
		}
	} else {
		contract, lookupErr := s.Workspace.Lookup(input.TargetContract, domain.Contract)
		if lookupErr != nil {
			return OperationResult{}, lookupErr
		}
		version, valueErr := workspace.String(contract.Document, "contract_version", true)
		if valueErr != nil {
			return OperationResult{}, valueErr
		}
		if version != input.BaseVersion || contract.Fingerprint != input.BaseFingerprint {
			return OperationResult{}, stale("base_fingerprint", input.BaseFingerprint, contract.Fingerprint)
		}
		if len(input.Additions)+len(input.Modifications)+len(input.Removals) == 0 {
			return OperationResult{}, invalid("delta", "an existing capability Proposal requires a non-empty delta")
		}
	}
	ref := string(domain.Proposal) + ":" + id.String()
	if err := s.authorize(input.Authorization, "proposal.create", ref, "absent", input.Scope, []string{"proposal.create"}); err != nil {
		return OperationResult{}, err
	}
	doc := s.document(domain.Proposal, id, input.Title, input.Actor, input.Status)
	workspace.SetString(doc, "target_contract", input.TargetContract)
	workspace.SetString(doc, "base_version", input.BaseVersion)
	workspace.SetString(doc, "base_fingerprint", input.BaseFingerprint)
	workspace.SetBool(doc, "new_capability", input.NewCapability)
	workspace.SetStrings(doc, "additions", input.Additions)
	from, to := splitModifications(input.Modifications)
	workspace.SetStrings(doc, "modification_from", from)
	workspace.SetStrings(doc, "modification_to", to)
	workspace.SetStrings(doc, "removals", input.Removals)
	setCandidate(doc, input.Candidate)
	workspace.SetString(doc, "rationale", input.Rationale)
	workspace.SetStrings(doc, "scope", input.Scope)
	workspace.SetStrings(doc, "gaps", input.Gaps)
	workspace.SetString(doc, "authorization", input.Authorization)
	workspace.SetString(doc, "idempotency_key", input.IdempotencyKey)
	s.addFreshness(doc, input.FreshnessValidUntil)
	return s.createOne("proposal.create", doc, input.IdempotencyKey, []string{input.Authorization, input.TargetContract})
}

func (s Service) CheckProposalBase(ref string) (map[string]any, error) {
	proposal, err := s.Workspace.Lookup(ref, domain.Proposal)
	if err != nil {
		return nil, err
	}
	targetRef, err := workspace.String(proposal.Document, "target_contract", true)
	if err != nil {
		return nil, err
	}
	baseFP, err := workspace.String(proposal.Document, "base_fingerprint", true)
	if err != nil {
		return nil, err
	}
	newCapability, err := workspace.Bool(proposal.Document, "new_capability", true)
	if err != nil {
		return nil, err
	}
	if newCapability {
		_, lookupErr := s.Workspace.Lookup(targetRef, domain.Contract)
		if lookupErr == nil {
			return nil, domain.NewRefusal(domain.RefusalCollision, "target_contract", "absence base no longer holds", nil)
		}
		return map[string]any{"proposal": refOf(proposal), "target": targetRef, "base_state": "absent", "matches": true}, nil
	}
	contract, err := s.Workspace.Lookup(targetRef, domain.Contract)
	if err != nil {
		return nil, err
	}
	if contract.Fingerprint != baseFP {
		return nil, stale("base_fingerprint", baseFP, contract.Fingerprint)
	}
	return map[string]any{"proposal": refOf(proposal), "target": targetRef, "base_fingerprint": baseFP, "matches": true}, nil
}

func (s Service) ProposalView(ref string) (map[string]any, error) {
	proposal, err := s.Workspace.Lookup(ref, domain.Proposal)
	if err != nil {
		return nil, err
	}
	candidate, err := s.candidate(proposal)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"ref": refOf(proposal), "path": proposal.Path, "fingerprint": proposal.Fingerprint,
		"status": value(proposal.Document.Record.Status), "target_contract": mustString(proposal.Document, "target_contract"),
		"base_version": mustString(proposal.Document, "base_version"), "base_fingerprint": mustString(proposal.Document, "base_fingerprint"),
		"delta":              map[string]any{"additions": mustStrings(proposal.Document, "additions"), "modification_from": mustStrings(proposal.Document, "modification_from"), "modification_to": mustStrings(proposal.Document, "modification_to"), "removals": mustStrings(proposal.Document, "removals")},
		"candidate_contract": candidate, "candidate_authoritative": false,
	}, nil
}

func (s Service) CreateMission(input MissionInput) (OperationResult, error) {
	id, err := domain.ParseID(input.ID)
	if err != nil {
		return OperationResult{}, err
	}
	proposal, err := s.Workspace.Lookup(input.Proposal, domain.Proposal)
	if err != nil {
		return OperationResult{}, err
	}
	if proposal.Document.Record.Status == nil || *proposal.Document.Record.Status != "accepted" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "proposal", "Mission requires an accepted Proposal", nil)
	}
	if proposal.Fingerprint != input.ExpectedProposalFingerprint {
		return OperationResult{}, stale("expected_proposal_fingerprint", input.ExpectedProposalFingerprint, proposal.Fingerprint)
	}
	if _, err := s.CheckProposalBase(input.Proposal); err != nil {
		return OperationResult{}, err
	}
	if input.Title == "" || input.Actor == "" || input.Outcome == "" || len(input.Objectives) == 0 || input.DesignSufficiency != "sufficient" || input.SliceQuality != "coherent" || len(input.EvidenceClaims) == 0 || len(input.Scope) == 0 || input.Baseline == "" || input.BudgetUnits < 1 || input.RepairBudget < 0 || input.RecoveryPoint == "" || input.ReturnDestination == "" || input.Authorization == "" || input.IdempotencyKey == "" || input.Preparation == nil {
		return OperationResult{}, missing("mission", "bounded outcome, Objectives, preparation verdicts, evidence plan, envelope, recovery, authorization, and idempotency are required")
	}
	if err := spectacularruntime.ValidatePreparationReceipt(*input.Preparation, s.now()); err != nil {
		return OperationResult{}, err
	}
	selectedOutcome := ""
	for _, candidate := range input.Preparation.Candidates {
		if candidate.Name == input.Preparation.Selected {
			selectedOutcome = candidate.Outcome
		}
	}
	if input.Preparation.Proposal.Ref != input.Proposal || input.Preparation.Proposal.Fingerprint != input.ExpectedProposalFingerprint || input.Preparation.Baseline != input.Baseline || input.Preparation.DesignSufficiency != input.DesignSufficiency || input.Preparation.SliceQuality != input.SliceQuality || selectedOutcome != input.Outcome || !sameStrings(input.Preparation.EvidenceClaims, input.EvidenceClaims) || !sameStrings(input.Preparation.StopConditions, input.Stops) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "preparation", "preparation receipt does not bind the Mission input", nil)
	}
	preparationBindings := []string{input.Preparation.Proposal.Ref + "@" + input.Preparation.Proposal.Fingerprint}
	for _, source := range input.Preparation.DirectionSources {
		preparationBindings = append(preparationBindings, source.Ref+"@"+source.Fingerprint)
	}
	if err := s.validateBoundInputs(preparationBindings); err != nil {
		return OperationResult{}, err
	}
	preparationJSON, err := json.Marshal(input.Preparation)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(input.ExpiresAt, s.now(), "expires_at"); err != nil {
		return OperationResult{}, err
	}
	missionRef := string(domain.Mission) + ":" + id.String()
	if err := s.authorize(input.Authorization, "mission.create", missionRef, "absent", input.Scope, []string{"mission.create"}); err != nil {
		return OperationResult{}, err
	}
	mission := s.document(domain.Mission, id, input.Title, input.Actor, "defined")
	source, _ := domain.ParseReference(input.Proposal)
	mission.Record.Source = &source
	workspace.SetString(mission, "outcome", input.Outcome)
	workspace.SetString(mission, "design_sufficiency", input.DesignSufficiency)
	workspace.SetString(mission, "slice_quality", input.SliceQuality)
	workspace.SetStrings(mission, "dependencies", input.Dependencies)
	workspace.SetStrings(mission, "gaps", input.Gaps)
	workspace.SetStrings(mission, "evidence_claims", input.EvidenceClaims)
	workspace.SetStrings(mission, "scope", input.Scope)
	workspace.SetStrings(mission, "allowed_actions", input.AllowedActions)
	workspace.SetStrings(mission, "forbidden_effects", input.ForbiddenEffects)
	workspace.SetString(mission, "baseline", input.Baseline)
	workspace.SetInt(mission, "budget_units", input.BudgetUnits)
	workspace.SetInt(mission, "repair_budget", input.RepairBudget)
	workspace.SetString(mission, "expires_at", input.ExpiresAt)
	workspace.SetStrings(mission, "stops", input.Stops)
	workspace.SetString(mission, "recovery_point", input.RecoveryPoint)
	workspace.SetString(mission, "return_destination", input.ReturnDestination)
	workspace.SetString(mission, "activation_decision", input.Authorization)
	workspace.SetString(mission, "idempotency_key", input.IdempotencyKey)
	workspace.SetString(mission, "preparation_fingerprint", input.Preparation.Fingerprint)
	workspace.SetString(mission, "preparation_valid_until", input.Preparation.FreshUntil)
	workspace.SetString(mission, "preparation_baseline", input.Preparation.Baseline)
	workspace.SetStrings(mission, "preparation_sources", preparationBindings)
	workspace.SetString(mission, "preparation_receipt", string(preparationJSON))
	var documents []*workspace.Document
	documents = append(documents, mission)
	objectiveRefs := make([]string, 0, len(input.Objectives))
	for _, item := range input.Objectives {
		objectiveID, parseErr := domain.ParseID(item.ID)
		if parseErr != nil || item.Outcome == "" || len(item.ExpectedProof) == 0 {
			return OperationResult{}, invalid("objectives", "each Objective needs UUIDv7 identity, outcome, and expected proof")
		}
		objective := s.document(domain.Objective, objectiveID, item.Outcome, input.Actor, "pending")
		workspace.SetString(objective, "mission", missionRef)
		workspace.SetString(objective, "outcome", item.Outcome)
		workspace.SetStrings(objective, "dependencies", item.Dependencies)
		workspace.SetStrings(objective, "expected_proof", item.ExpectedProof)
		documents = append(documents, objective)
		objectiveRefs = append(objectiveRefs, string(domain.Objective)+":"+objectiveID.String())
	}
	workspace.SetStrings(mission, "objectives", objectiveRefs)
	if input.InitialRunID != "" {
		runID, parseErr := domain.ParseID(input.InitialRunID)
		if parseErr != nil {
			return OperationResult{}, parseErr
		}
		run := s.document(domain.Run, runID, "Initial governed Run", input.Actor, "active")
		workspace.SetString(run, "mission", missionRef)
		workspace.SetString(run, "baseline", input.Baseline)
		workspace.SetString(run, "authority", input.Authorization)
		documents = append(documents, run)
		workspace.SetString(mission, "current_run", string(domain.Run)+":"+runID.String())
	}
	return s.createMany("mission.create", mission, documents, input.IdempotencyKey, []string{input.Proposal, input.Authorization})
}

// CompileAutopilot derives hard limits from the exact current Mission and
// owner Decision before the runtime compiler can emit a charter.
func (s Service) CompileAutopilot(input spectacularruntime.AutopilotInput) (spectacularruntime.AutopilotCharter, error) {
	mission, err := s.Workspace.Lookup(input.Mission.Ref, domain.Mission)
	if err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	if mission.Fingerprint != input.Mission.Fingerprint {
		return spectacularruntime.AutopilotCharter{}, stale("mission", input.Mission.Fingerprint, mission.Fingerprint)
	}
	if err := s.validateFresh(mission, "mission"); err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	if value(mission.Document.Record.Status) != "active" {
		return spectacularruntime.AutopilotCharter{}, domain.NewRefusal(domain.RefusalInvalidTransition, "mission", "Autopilot requires an active Mission", nil)
	}
	decision, err := s.Workspace.Lookup(input.Authorization.Ref, domain.Decision)
	if err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	if decision.Fingerprint != input.Authorization.Fingerprint {
		return spectacularruntime.AutopilotCharter{}, stale("authorization", input.Authorization.Fingerprint, decision.Fingerprint)
	}
	if err := s.validateFresh(decision, "authorization"); err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	for _, source := range input.AuthoritativeSources {
		_, actual, lookupErr := s.Workspace.Source(source.Ref)
		if lookupErr != nil {
			return spectacularruntime.AutopilotCharter{}, lookupErr
		}
		if actual != source.Fingerprint {
			return spectacularruntime.AutopilotCharter{}, stale("authoritative_sources", source.Fingerprint, actual)
		}
	}
	envelope, err := parseMissionEnvelope(mission)
	if err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	missionOutcome, err := workspace.String(mission.Document, "outcome", true)
	if err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	repairBudget, err := workspace.Int(mission.Document, "repair_budget", true)
	if err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	requiredEffects := append([]string{"mission.autopilot"}, input.AllowedActions...)
	for _, provider := range input.AllowedProviders {
		requiredEffects = append(requiredEffects, "provider:"+provider)
	}
	if err := s.authorize(input.Authorization.Ref, "mission.autopilot", input.Mission.Ref, input.Mission.Fingerprint, envelope.Scope, requiredEffects); err != nil {
		return spectacularruntime.AutopilotCharter{}, err
	}
	limitExpiry := envelope.ExpiresAt
	if decisionExpiry, _ := workspace.String(decision.Document, "expires_at", false); decisionExpiry != "" {
		missionTime, missionErr := time.Parse(time.RFC3339, limitExpiry)
		decisionTime, decisionErr := time.Parse(time.RFC3339, decisionExpiry)
		if missionErr != nil || decisionErr != nil {
			return spectacularruntime.AutopilotCharter{}, invalid("expires_at", "Mission and Decision expiry must be RFC3339")
		}
		if decisionTime.Before(missionTime) {
			limitExpiry = decisionExpiry
		}
	}
	limits := spectacularruntime.AutopilotLimits{
		Mission: input.Mission, Authorization: input.Authorization, Outcome: missionOutcome,
		AllowedProviders: input.AllowedProviders, AllowedActions: envelope.AllowedActions,
		ForbiddenEffects: envelope.ForbiddenEffects, BudgetUnits: envelope.BudgetUnits,
		RepairBudget: repairBudget, ExpiresAt: limitExpiry, StopConditions: envelope.Stops,
	}
	return spectacularruntime.CompileAutopilot(input, limits, s.now())
}

func (s Service) document(noun domain.RecordType, id domain.ID, title, actor, status string) *workspace.Document {
	now := s.now().Format(time.RFC3339)
	doc := &workspace.Document{Record: domain.Record{Type: noun, ID: id}, Unknown: map[string]*yaml.Node{}, Body: "# " + string(noun) + "\n"}
	if title != "" {
		doc.Record.Title = stringPtr(title)
	}
	if actor != "" {
		doc.Record.CreatedBy = stringPtr(actor)
	}
	if status != "" {
		doc.Record.Status = stringPtr(status)
	}
	doc.Record.Created = stringPtr(now)
	doc.Record.Updated = stringPtr(now)
	s.addFreshness(doc, "")
	return doc
}

func (s Service) addFreshness(doc *workspace.Document, validUntil string) {
	now := s.now()
	if validUntil == "" {
		validUntil = now.Add(24 * time.Hour).Format(time.RFC3339)
	}
	marker := filepath.Join(s.Workspace.Root, ".spectacular", "workspace.yaml")
	data, _ := os.ReadFile(marker)
	sum := sha256.Sum256(data)
	workspace.SetString(doc, "freshness_checked_at", now.Format(time.RFC3339))
	workspace.SetString(doc, "freshness_valid_until", validUntil)
	workspace.SetString(doc, "freshness_source", ".spectacular/workspace.yaml")
	workspace.SetString(doc, "freshness_source_fingerprint", hex.EncodeToString(sum[:]))
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func (s Service) createOne(operation string, doc *workspace.Document, key string, sources []string) (OperationResult, error) {
	return s.createMany(operation, doc, []*workspace.Document{doc}, key, sources)
}

func (s Service) createMany(operation string, primary *workspace.Document, docs []*workspace.Document, key string, sources []string) (OperationResult, error) {
	if key == "" {
		return OperationResult{}, missing("idempotency_key", "required")
	}
	allReplay := true
	changes := make([]FileChange, 0, len(docs))
	for _, doc := range docs {
		canonical, err := workspace.Canonical(doc)
		if err != nil {
			return OperationResult{}, err
		}
		existing, found := s.entryByID(doc.Record.ID)
		if found {
			// System timestamps and freshness checks are creation metadata, not
			// caller content. Normalize them to the existing record so an exact
			// identity replay remains exact even in a later process.
			doc.Record.Created = existing.Document.Record.Created
			doc.Record.Updated = existing.Document.Record.Updated
			for _, field := range []string{"freshness_checked_at", "freshness_valid_until", "freshness_source", "freshness_source_fingerprint"} {
				workspace.SetString(doc, field, mustString(existing.Document, field))
			}
			canonical, err = workspace.Canonical(doc)
			if err != nil {
				return OperationResult{}, err
			}
			existingCanonical, canonicalErr := workspace.Canonical(existing.Document)
			if canonicalErr != nil {
				return OperationResult{}, canonicalErr
			}
			if string(existingCanonical) != string(canonical) {
				return OperationResult{}, domain.NewRefusal(domain.RefusalIdempotencyConflict, doc.Record.ID.String(), "identity exists with different content", nil)
			}
			continue
		}
		path := filepath.ToSlash(recordPath(doc.Record.Type, doc.Record.ID))
		if occupied, pathFound := s.entryByPath(path); pathFound {
			return OperationResult{}, domain.NewRefusal(domain.RefusalCollision, path, "record path is occupied by "+refOf(occupied), nil)
		}
		allReplay = false
		changes = append(changes, FileChange{Path: path, Data: canonical, Mode: 0o644})
	}
	if allReplay {
		entry, _ := s.entryByID(primary.Record.ID)
		return result(operation, entry, true, sources), nil
	}
	if len(changes) != len(docs) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalIdempotencyConflict, "idempotency_key", "partial replay would create an incomplete record set", nil)
	}
	if err := ApplyTransaction(s.Workspace.Root, key, changes); err != nil {
		return OperationResult{}, err
	}
	path := recordPath(primary.Record.Type, primary.Record.ID)
	fp, err := workspace.Fingerprint(primary)
	if err != nil {
		return OperationResult{}, err
	}
	return OperationResult{Operation: operation, Ref: string(primary.Record.Type) + ":" + primary.Record.ID.String(), Path: filepath.ToSlash(path), Fingerprint: fp, Sources: sources}, nil
}

func (s Service) authorize(ref, operation, target, expected string, requiredScope, requiredEffects []string) error {
	decision, err := s.Workspace.Lookup(ref, domain.Decision)
	if err != nil {
		return err
	}
	role, err := workspace.String(decision.Document, "actor_role", true)
	if err != nil || role != "owner" {
		return domain.NewRefusal(domain.RefusalUnauthorized, "actor_role", "owner Decision required", err)
	}
	authority, err := s.authority(decision)
	if err != nil || authority.Operation != operation {
		return domain.NewRefusal(domain.RefusalUnauthorized, "operation", "Decision does not authorize "+operation, err)
	}
	if !containsAll(authority.Scope, requiredScope) {
		return domain.NewRefusal(domain.RefusalUnauthorized, "scope", "Decision scope does not contain the requested operation scope", nil)
	}
	if !containsAll(authority.AuthorizedEffects, requiredEffects) {
		return domain.NewRefusal(domain.RefusalUnauthorized, "authorized_effects", "Decision does not authorize every requested effect", nil)
	}
	if err := s.evaluateConditions(authority, expected); err != nil {
		return err
	}
	if expiry, _ := workspace.String(decision.Document, "expires_at", false); expiry != "" {
		when, parseErr := time.Parse(time.RFC3339, expiry)
		if parseErr != nil || !when.After(s.now()) {
			return domain.NewRefusal(domain.RefusalExpiredAuthority, "expires_at", "Decision authority expired", parseErr)
		}
	}
	targets := authority.Targets
	if !contains(targets, target) {
		return domain.NewRefusal(domain.RefusalUnauthorized, "targets", "Decision does not name target "+target, err)
	}
	fingerprints := authority.ExpectedFingerprints
	for i := range targets {
		if targets[i] == target && fingerprints[i] != expected {
			return stale("authorization.expected_fingerprint", fingerprints[i], expected)
		}
	}
	for _, candidate := range s.Workspace.OfType(domain.Decision) {
		supersedes, _ := workspace.String(candidate.Document, "supersedes", false)
		if supersedes == ref && candidate.Document.Record.ID != decision.Document.Record.ID {
			return domain.NewRefusal(domain.RefusalConflictingAuthority, "authorization", "Decision was superseded by "+refOf(candidate), nil)
		}
	}
	return nil
}

func (s Service) entryByID(id domain.ID) (discovery.Entry, bool) {
	for _, entry := range s.Workspace.Entries {
		if entry.Document.Record.ID == id {
			return entry, true
		}
	}
	return discovery.Entry{}, false
}

func (s Service) entryByPath(path string) (discovery.Entry, bool) {
	for _, entry := range s.Workspace.Entries {
		if entry.Path == path {
			return entry, true
		}
	}
	return discovery.Entry{}, false
}

func (s Service) candidate(proposal discovery.Entry) (ContractCandidate, error) {
	newCapability, err := workspace.Bool(proposal.Document, "new_capability", true)
	if err != nil {
		return ContractCandidate{}, err
	}
	if newCapability {
		return candidateFromDocument(proposal.Document), nil
	}
	target, err := workspace.String(proposal.Document, "target_contract", true)
	if err != nil {
		return ContractCandidate{}, err
	}
	contract, err := s.Workspace.Lookup(target, domain.Contract)
	if err != nil {
		return ContractCandidate{}, err
	}
	candidate := candidateFromDocument(contract.Document)
	additions, _ := workspace.Strings(proposal.Document, "additions", false)
	removals, _ := workspace.Strings(proposal.Document, "removals", false)
	from, _ := workspace.Strings(proposal.Document, "modification_from", false)
	to, _ := workspace.Strings(proposal.Document, "modification_to", false)
	if len(from) != len(to) {
		return ContractCandidate{}, invalid("modifications", "from/to cardinality mismatch")
	}
	behavior := append([]string(nil), candidate.RequiredBehavior...)
	for _, removal := range removals {
		var found bool
		behavior, found = removeExact(behavior, removal)
		if !found {
			return ContractCandidate{}, domain.NewRefusal(domain.RefusalCollision, "removals", "base Contract does not contain exact clause: "+removal, nil)
		}
	}
	for i := range from {
		index := indexOf(behavior, from[i])
		if index < 0 {
			return ContractCandidate{}, domain.NewRefusal(domain.RefusalCollision, "modifications", "base Contract does not contain exact clause: "+from[i], nil)
		}
		behavior[index] = to[i]
	}
	for _, addition := range additions {
		if contains(behavior, addition) {
			return ContractCandidate{}, domain.NewRefusal(domain.RefusalCollision, "additions", "candidate already contains clause: "+addition, nil)
		}
		behavior = append(behavior, addition)
	}
	candidate.RequiredBehavior = behavior
	return candidate, validateCandidate(candidate)
}

func setCandidate(doc *workspace.Document, candidate ContractCandidate) {
	workspace.SetString(doc, "candidate_purpose", candidate.Purpose)
	workspace.SetString(doc, "candidate_outcome", candidate.Outcome)
	workspace.SetStrings(doc, "candidate_applies_when", candidate.AppliesWhen)
	workspace.SetStrings(doc, "candidate_does_not_apply_when", candidate.DoesNotApplyWhen)
	workspace.SetStrings(doc, "candidate_does_not_provide", candidate.DoesNotProvide)
	workspace.SetStrings(doc, "candidate_required_behavior", candidate.RequiredBehavior)
	workspace.SetStrings(doc, "candidate_operating_cases", candidate.OperatingCases)
	workspace.SetStrings(doc, "candidate_persistent_information", candidate.PersistentInformation)
	workspace.SetStrings(doc, "candidate_conformance_checks", candidate.ConformanceChecks)
	workspace.SetString(doc, "candidate_authority_freshness", candidate.AuthorityFreshness)
	workspace.SetStrings(doc, "candidate_related_material", candidate.RelatedMaterial)
}

func candidateFromDocument(doc *workspace.Document) ContractCandidate {
	prefix := ""
	if doc.Record.Type == domain.Proposal {
		prefix = "candidate_"
	}
	return ContractCandidate{
		Purpose: mustString(doc, prefix+"purpose"), Outcome: mustString(doc, prefix+"outcome"),
		AppliesWhen: mustStrings(doc, prefix+"applies_when"), DoesNotApplyWhen: mustStrings(doc, prefix+"does_not_apply_when"),
		DoesNotProvide: mustStrings(doc, prefix+"does_not_provide"), RequiredBehavior: mustStrings(doc, prefix+"required_behavior"),
		OperatingCases: mustStrings(doc, prefix+"operating_cases"), PersistentInformation: mustStrings(doc, prefix+"persistent_information"),
		ConformanceChecks: mustStrings(doc, prefix+"conformance_checks"), AuthorityFreshness: mustString(doc, prefix+"authority_freshness"),
		RelatedMaterial: mustStrings(doc, prefix+"related_material"),
	}
}

func validateCandidate(candidate ContractCandidate) error {
	if candidate.Purpose == "" || candidate.Outcome == "" || len(candidate.AppliesWhen) == 0 || len(candidate.DoesNotApplyWhen) == 0 || len(candidate.DoesNotProvide) == 0 || len(candidate.RequiredBehavior) == 0 || len(candidate.OperatingCases) == 0 || len(candidate.PersistentInformation) == 0 || len(candidate.ConformanceChecks) == 0 || candidate.AuthorityFreshness == "" {
		return missing("candidate", "complete MCC candidate fields are required")
	}
	return nil
}

func setContract(doc *workspace.Document, candidate ContractCandidate, version string) {
	workspace.SetString(doc, "contract_version", version)
	workspace.SetString(doc, "purpose", candidate.Purpose)
	workspace.SetString(doc, "outcome", candidate.Outcome)
	workspace.SetStrings(doc, "applies_when", candidate.AppliesWhen)
	workspace.SetStrings(doc, "does_not_apply_when", candidate.DoesNotApplyWhen)
	workspace.SetStrings(doc, "does_not_provide", candidate.DoesNotProvide)
	workspace.SetStrings(doc, "required_behavior", candidate.RequiredBehavior)
	workspace.SetStrings(doc, "operating_cases", candidate.OperatingCases)
	workspace.SetStrings(doc, "persistent_information", candidate.PersistentInformation)
	workspace.SetStrings(doc, "conformance_checks", candidate.ConformanceChecks)
	workspace.SetString(doc, "authority_freshness", candidate.AuthorityFreshness)
	workspace.SetStrings(doc, "related_material", candidate.RelatedMaterial)
}

func recordPath(noun domain.RecordType, id domain.ID) string {
	return filepath.ToSlash(filepath.Join(".spectacular", "records", strings.ToLower(string(noun))+"-"+id.String()+".md"))
}

func result(operation string, entry discovery.Entry, replay bool, sources []string) OperationResult {
	return OperationResult{Operation: operation, Ref: refOf(entry), Path: entry.Path, Fingerprint: entry.Fingerprint, IdempotentReplay: replay, Sources: sources}
}

func refOf(entry discovery.Entry) string {
	return string(entry.Document.Record.Type) + ":" + entry.Document.Record.ID.String()
}
func stringPtr(value string) *string { return &value }
func value(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func missing(field, detail string) error {
	return domain.NewRefusal(domain.RefusalMissingRequiredField, field, detail, nil)
}
func invalid(field, detail string) error {
	return domain.NewRefusal(domain.RefusalInvalidKnownField, field, detail, nil)
}
func stale(field, expected, actual string) error {
	return domain.NewStateRefusal(domain.RefusalStaleFingerprint, field, "expected "+expected+", actual "+actual, expected, actual, "reload current canonical state and obtain fresh explicit authorization", nil)
}

func parseFuture(raw string, now time.Time, field string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil || !parsed.After(now) {
		return time.Time{}, domain.NewRefusal(domain.RefusalExpiredAuthority, field, "must be a future RFC3339 time", err)
	}
	return parsed, nil
}

func splitModifications(values []Modification) ([]string, []string) {
	from := make([]string, len(values))
	to := make([]string, len(values))
	for i, item := range values {
		from[i], to[i] = item.From, item.To
	}
	return from, to
}
func contains(values []string, want string) bool { return indexOf(values, want) >= 0 }
func indexOf(values []string, want string) int {
	for i, value := range values {
		if value == want {
			return i
		}
	}
	return -1
}
func removeExact(values []string, want string) ([]string, bool) {
	index := indexOf(values, want)
	if index < 0 {
		return values, false
	}
	return append(values[:index:index], values[index+1:]...), true
}

func mustString(doc *workspace.Document, name string) string {
	value, _ := workspace.String(doc, name, false)
	return value
}
func mustStrings(doc *workspace.Document, name string) []string {
	values, _ := workspace.Strings(doc, name, false)
	if values == nil {
		return []string{}
	}
	return values
}

func sortedUnique(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	out := result[:0]
	for _, value := range result {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}
