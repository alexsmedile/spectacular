package governance

import (
	"os"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// ProgressMission records implementation progress without claiming Evidence
// sufficiency, owner acceptance, or lifecycle completion.
func (s Service) ProgressMission(input MissionProgressInput) (OperationResult, error) {
	if input.Mission == "" || input.Objective == "" || input.To != "implemented" || input.Actor == "" || input.ExpectedMissionFingerprint == "" || input.ExpectedObjectiveFingerprint == "" || input.IdempotencyKey == "" {
		return OperationResult{}, missing("mission_progress", "Mission, Objective, implemented result, actor, fingerprints, and idempotency are required")
	}
	mission, err := s.Workspace.Lookup(input.Mission, domain.Mission)
	if err != nil {
		return OperationResult{}, err
	}
	objective, err := s.Workspace.Lookup(input.Objective, domain.Objective)
	if err != nil {
		return OperationResult{}, err
	}
	lastKey, err := workspace.String(objective.Document, "progress_idempotency_key", false)
	if err != nil {
		return OperationResult{}, err
	}
	if value(objective.Document.Record.Status) == "implemented" && lastKey == input.IdempotencyKey {
		return result("mission.progress", objective, true, []string{refOf(mission)}), nil
	}
	if mission.Fingerprint != input.ExpectedMissionFingerprint {
		return OperationResult{}, stale("expected_mission_fingerprint", input.ExpectedMissionFingerprint, mission.Fingerprint)
	}
	if objective.Fingerprint != input.ExpectedObjectiveFingerprint {
		return OperationResult{}, stale("expected_objective_fingerprint", input.ExpectedObjectiveFingerprint, objective.Fingerprint)
	}
	if value(mission.Document.Record.Status) != "active" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "mission", "progress requires an active Mission", nil)
	}
	if err := s.validateFresh(mission, "mission"); err != nil {
		return OperationResult{}, err
	}
	if err := s.validateFresh(objective, "objective"); err != nil {
		return OperationResult{}, err
	}
	missionExpiry, err := workspace.String(mission.Document, "expires_at", true)
	if err != nil {
		return OperationResult{}, err
	}
	if _, err := parseFuture(missionExpiry, s.now(), "mission.expires_at"); err != nil {
		return OperationResult{}, err
	}
	allowed, err := workspace.Strings(mission.Document, "allowed_actions", true)
	if err != nil {
		return OperationResult{}, err
	}
	if !contains(allowed, "update-mission-progress") {
		return OperationResult{}, domain.NewRefusal(domain.RefusalUnauthorized, "allowed_actions", "Mission does not allow update-mission-progress", nil)
	}
	objectiveMission, err := workspace.String(objective.Document, "mission", true)
	if err != nil || objectiveMission != refOf(mission) {
		return OperationResult{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "objective", "Objective does not belong to the Mission", err)
	}
	if value(objective.Document.Record.Status) != "pending" {
		return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "objective", "only a pending Objective can become implemented", nil)
	}
	dependencies, err := workspace.Strings(objective.Document, "dependencies", false)
	if err != nil {
		return OperationResult{}, err
	}
	for _, ref := range dependencies {
		dependency, lookupErr := s.Workspace.Lookup(ref, domain.Objective)
		if lookupErr != nil {
			return OperationResult{}, lookupErr
		}
		status := value(dependency.Document.Record.Status)
		if status != "implemented" && status != "satisfied" {
			return OperationResult{}, domain.NewRefusal(domain.RefusalInvalidTransition, "dependencies", "Objective dependency is not implemented: "+ref, nil)
		}
	}
	objective.Document.Record.Status = stringPtr("implemented")
	objective.Document.Record.Updated = stringPtr(s.now().Format(time.RFC3339))
	workspace.SetString(objective.Document, "progress_actor", input.Actor)
	workspace.SetString(objective.Document, "progress_idempotency_key", input.IdempotencyKey)
	data, err := workspace.Canonical(objective.Document)
	if err != nil {
		return OperationResult{}, err
	}
	if err := ApplyTransaction(s.Workspace.Root, input.IdempotencyKey, []FileChange{{Path: objective.Path, Data: data, Mode: os.FileMode(0o644)}}); err != nil {
		return OperationResult{}, err
	}
	fingerprint, err := workspace.Fingerprint(objective.Document)
	if err != nil {
		return OperationResult{}, err
	}
	return OperationResult{Operation: "mission.progress", Ref: refOf(objective), Path: objective.Path, Fingerprint: fingerprint, Sources: []string{refOf(mission)}}, nil
}
