// Package projection builds disposable, source-attributed recovery views.
package projection

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

type Pointer struct {
	Noun        string `json:"noun"`
	Ref         string `json:"ref"`
	HumanRef    string `json:"human_ref,omitempty"`
	Path        string `json:"path"`
	Fingerprint string `json:"fingerprint"`
	ShowCommand string `json:"show_command,omitempty"`
}

type Source struct {
	Path        string `json:"path"`
	Fingerprint string `json:"fingerprint"`
}
type Freshness struct {
	State               string   `json:"state"`
	Basis               string   `json:"basis"`
	Source              Source   `json:"source"`
	Pointer             *Pointer `json:"source_pointer,omitempty"`
	CheckedAt           string   `json:"checked_at"`
	ValidUntil          string   `json:"valid_until"`
	ExpectedFingerprint string   `json:"expected_fingerprint"`
	ActualFingerprint   string   `json:"actual_fingerprint"`
}
type OwnerGate struct {
	Code   string  `json:"code"`
	Detail string  `json:"detail"`
	Source Pointer `json:"source"`
}
type Continuation struct {
	Operation    string  `json:"operation"`
	Target       Pointer `json:"target"`
	AuthorizedBy Pointer `json:"authorized_by"`
}
type Card struct {
	Noun             string        `json:"noun"`
	Ref              string        `json:"ref"`
	ID               string        `json:"id"`
	Title            string        `json:"title,omitempty"`
	Outcome          string        `json:"outcome,omitempty"`
	State            string        `json:"state,omitempty"`
	Freshness        Freshness     `json:"freshness"`
	Source           Source        `json:"source"`
	Sources          []Pointer     `json:"sources"`
	Pointers         []Pointer     `json:"pointers"`
	Gaps             []Pointer     `json:"gaps"`
	Conflicts        []string      `json:"conflicts"`
	Omissions        []string      `json:"omissions"`
	Continuation     *Continuation `json:"continuation,omitempty"`
	OwnerGate        *OwnerGate    `json:"owner_gate,omitempty"`
	CurrentObjective *Pointer      `json:"current_objective,omitempty"`
	CurrentRun       *Pointer      `json:"current_run,omitempty"`
	LatestCheckpoint *Pointer      `json:"latest_checkpoint,omitempty"`
}
type List struct {
	Items []Card `json:"items"`
}
type Validation struct {
	Scope   string   `json:"scope"`
	Valid   bool     `json:"valid"`
	Records int      `json:"records"`
	Sources []Source `json:"sources"`
}
type Envelope struct {
	SchemaVersion   string `json:"schema_version"`
	GeneratedAt     string `json:"generated_at"`
	GenerationBasis string `json:"generation_basis"`
	Data            any    `json:"data"`
}

type ProjectAuthority struct {
	Identity     Pointer   `json:"identity"`
	Direction    string    `json:"direction"`
	Boundaries   []string  `json:"boundaries"`
	Constraints  []string  `json:"constraints"`
	CurrentTruth []Pointer `json:"current_truth"`
}
type ProjectProjection struct {
	Missions     []Pointer     `json:"missions"`
	Gaps         []Pointer     `json:"gaps"`
	Conflicts    []string      `json:"conflicts"`
	Omissions    []string      `json:"omissions"`
	Continuation *Continuation `json:"continuation,omitempty"`
	OwnerGate    *OwnerGate    `json:"owner_gate,omitempty"`
}
type ProjectView struct {
	Source        Source            `json:"source"`
	Freshness     Freshness         `json:"freshness"`
	Authoritative ProjectAuthority  `json:"authoritative"`
	Projection    ProjectProjection `json:"projection"`
}

type Builder struct {
	Workspace   *discovery.Workspace
	Now         func() time.Time
	ShowCommand func(domain.RecordType, string) (string, bool)
}

func (b Builder) Envelope(schema string, data any) Envelope {
	now := b.Now
	if now == nil {
		now = time.Now
	}
	return Envelope{SchemaVersion: schema, GeneratedAt: now().UTC().Format(time.RFC3339Nano), GenerationBasis: b.generationBasis(), Data: data}
}

func (b Builder) Project() (ProjectView, error) {
	anchor := b.Workspace.ProjectAnchor()
	card, err := b.base(anchor)
	if err != nil {
		return ProjectView{}, err
	}
	if card.Freshness.State == "stale" {
		return ProjectView{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", anchor.Path+" is stale", nil)
	}
	var anchorGate *OwnerGate
	if card.Freshness.State != "current" {
		anchorGate = &OwnerGate{Code: "confirm_freshness", Detail: "project Anchor freshness is unknown", Source: b.pointer(anchor)}
	}
	truthRefs, err := workspace.Strings(anchor.Document, "current_truth", true)
	if err != nil {
		return ProjectView{}, err
	}
	for _, ref := range truthRefs {
		typed, parseErr := domain.ParseReference(ref)
		if parseErr != nil {
			return ProjectView{}, parseErr
		}
		if typed.Type != domain.Mission {
			continue
		}
		target, err := b.Workspace.Lookup(ref, domain.Mission)
		if err != nil {
			return ProjectView{}, err
		}
		card.Pointers = append(card.Pointers, b.pointer(target))
	}
	if len(card.Pointers) == 0 {
		lastClosed, valueErr := workspace.String(anchor.Document, "last_closed_mission", false)
		if valueErr != nil {
			return ProjectView{}, valueErr
		}
		if lastClosed != "" {
			target, lookupErr := b.Workspace.Lookup(lastClosed, domain.Mission)
			if lookupErr != nil {
				return ProjectView{}, lookupErr
			}
			card.Pointers = append(card.Pointers, b.pointer(target))
		}
	}
	if len(card.Pointers) == 1 {
		mission, missionErr := b.Mission(card.Pointers[0].Ref)
		if missionErr != nil {
			return ProjectView{}, missionErr
		}
		card.Gaps = mission.Gaps
		card.Conflicts = mission.Conflicts
		card.Continuation = mission.Continuation
		card.OwnerGate = mission.OwnerGate
	}
	if anchorGate != nil {
		card.Continuation = nil
		card.OwnerGate = anchorGate
	}
	if len(card.Pointers) != 1 {
		card.OwnerGate = &OwnerGate{Code: "select_mission", Detail: fmt.Sprintf("project has %d Missions; owner must select exactly one", len(card.Pointers)), Source: b.pointer(anchor)}
	}
	direction, err := workspace.String(anchor.Document, "direction", true)
	if err != nil {
		return ProjectView{}, err
	}
	boundaries, err := workspace.Strings(anchor.Document, "boundaries", true)
	if err != nil {
		return ProjectView{}, err
	}
	constraints, err := workspace.Strings(anchor.Document, "constraints", true)
	if err != nil {
		return ProjectView{}, err
	}
	truth := []Pointer{}
	for _, ref := range truthRefs {
		typed, e := domain.ParseReference(ref)
		if e != nil {
			return ProjectView{}, e
		}
		target, e := b.Workspace.Lookup(ref, typed.Type)
		if e != nil {
			return ProjectView{}, e
		}
		truth = append(truth, b.pointer(target))
	}
	return ProjectView{Source: card.Source, Freshness: card.Freshness, Authoritative: ProjectAuthority{Identity: b.pointer(anchor), Direction: direction, Boundaries: boundaries, Constraints: constraints, CurrentTruth: truth}, Projection: ProjectProjection{Missions: card.Pointers, Gaps: card.Gaps, Conflicts: card.Conflicts, Omissions: card.Omissions, Continuation: card.Continuation, OwnerGate: card.OwnerGate}}, nil
}

func (b Builder) MissionList() (List, error) {
	entries := b.Workspace.OfType(domain.Mission)
	list := List{Items: make([]Card, 0, len(entries))}
	for _, entry := range entries {
		archived, _ := workspace.Bool(entry.Document, "archived", false)
		if archived {
			continue
		}
		card, err := b.Mission(entry.Document.Record.ID.String())
		if err != nil {
			return List{}, err
		}
		list.Items = append(list.Items, card)
	}
	return list, nil
}

func (b Builder) Mission(ref string) (Card, error) {
	mission, err := b.Workspace.Lookup(ref, domain.Mission)
	if err != nil {
		return Card{}, err
	}
	card, err := b.base(mission)
	if err != nil {
		return Card{}, err
	}
	card.Outcome, _ = workspace.String(mission.Document, "outcome", false)
	if objectiveRefs, objectiveErr := workspace.Strings(mission.Document, "objectives", false); objectiveErr == nil {
		for _, objectiveRef := range objectiveRefs {
			objective, lookupErr := b.Workspace.Lookup(objectiveRef, domain.Objective)
			if lookupErr != nil {
				return Card{}, lookupErr
			}
			pointer := b.pointer(objective)
			if card.CurrentObjective == nil || value(objective.Document.Record.Status) != "satisfied" {
				card.CurrentObjective = &pointer
			}
			if value(objective.Document.Record.Status) != "satisfied" {
				break
			}
		}
	}
	if mission.Document.Record.Source != nil {
		proposal, lookupErr := b.Workspace.Lookup(mission.Document.Record.Source.String(), domain.Proposal)
		if lookupErr != nil {
			return Card{}, lookupErr
		}
		card.Sources = append(card.Sources, b.pointer(proposal))
	}
	if assessmentRef, _ := workspace.String(mission.Document, "assessment", false); assessmentRef != "" {
		assessment, lookupErr := b.Workspace.Lookup(assessmentRef, domain.Assessment)
		if lookupErr != nil {
			return Card{}, lookupErr
		}
		appendPointer(&card.Pointers, b.pointer(assessment))
		if evidenceRefs, parseErr := workspace.Strings(assessment.Document, "evidence", false); parseErr == nil {
			for _, evidenceRef := range evidenceRefs {
				evidence, evidenceErr := b.Workspace.Lookup(evidenceRef, domain.Evidence)
				if evidenceErr != nil {
					return Card{}, evidenceErr
				}
				appendPointer(&card.Pointers, b.pointer(evidence))
			}
		}
	}
	archived, _ := workspace.Bool(mission.Document, "archived", false)
	if mission.Document.Record.Status != nil && *mission.Document.Record.Status == "resolved" {
		for _, field := range []struct {
			name string
			noun domain.RecordType
		}{{"assessment", domain.Assessment}, {"reconciliation", domain.Evidence}, {"last_authorization", domain.Decision}, {"archive_authorization", domain.Decision}} {
			ref, _ := workspace.String(mission.Document, field.name, false)
			if ref == "" {
				continue
			}
			target, lookupErr := b.Workspace.Lookup(ref, field.noun)
			if lookupErr != nil {
				return Card{}, lookupErr
			}
			card.Pointers = append(card.Pointers, b.pointer(target))
		}
		if archived {
			authorizationRef, _ := workspace.String(mission.Document, "last_authorization", false)
			if authorizationRef != "" {
				authorization, lookupErr := b.Workspace.Lookup(authorizationRef, domain.Decision)
				if lookupErr != nil {
					return Card{}, lookupErr
				}
				next, _ := workspace.String(mission.Document, "terminal_next_action", false)
				card.Continuation = &Continuation{Operation: next, Target: b.pointer(mission), AuthorizedBy: b.pointer(authorization)}
			}
		}
		return card, nil
	}
	runRef, err := workspace.String(mission.Document, "current_run", true)
	if err != nil {
		return Card{}, err
	}
	run, err := b.Workspace.Lookup(runRef, domain.Run)
	if err != nil {
		return Card{}, err
	}
	card.Pointers = append(card.Pointers, b.pointer(run))
	runPointer := b.pointer(run)
	card.CurrentRun = &runPointer
	if checkpointRef, checkpointErr := workspace.String(run.Document, "latest_checkpoint", false); checkpointErr == nil && checkpointRef != "" {
		checkpoint, lookupErr := b.Workspace.Lookup(checkpointRef, domain.Checkpoint)
		if lookupErr != nil {
			return Card{}, lookupErr
		}
		checkpointPointer := b.pointer(checkpoint)
		card.LatestCheckpoint = &checkpointPointer
	}
	expectedRun, err := fingerprint(mission.Document, "expected_run_fingerprint")
	if err != nil {
		return Card{}, err
	}
	if expectedRun != run.Fingerprint {
		return Card{}, domain.NewRefusal(domain.RefusalInvalidFingerprint, "expected_run_fingerprint", "Mission fingerprint does not match current Run", nil)
	}
	for _, gap := range b.Workspace.OfType(domain.Gap) {
		scope, e := workspace.String(gap.Document, "scope", true)
		if e != nil {
			return Card{}, e
		}
		typed, e := domain.ParseReference(scope)
		if e != nil {
			return Card{}, e
		}
		if typed.Type != domain.Mission {
			return Card{}, domain.NewRefusal(domain.RefusalNounMismatch, "scope", "Gap scope must be Mission", nil)
		}
		if typed.ID == mission.Document.Record.ID {
			card.Gaps = append(card.Gaps, b.pointer(gap))
		}
	}
	if card.Freshness.State == "stale" {
		return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", mission.Path+" is stale", nil)
	}
	if card.Freshness.State != "current" {
		card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Mission freshness is unknown", Source: b.pointer(mission)}
		return card, nil
	}
	for _, gapPtr := range card.Gaps {
		gap, _ := b.Workspace.Lookup(gapPtr.Ref, domain.Gap)
		gapFresh, e := b.freshness(gap.Document)
		if e != nil {
			return Card{}, e
		}
		if gapFresh.State == "stale" {
			return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", gap.Path+" is stale", nil)
		}
		if gapFresh.State != "current" {
			card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Gap freshness is unknown", Source: gapPtr}
			return card, nil
		}
		blocking, e := workspace.Bool(gap.Document, "blocking", true)
		if e != nil {
			return Card{}, e
		}
		if blocking {
			card.OwnerGate = &OwnerGate{Code: "resolve_blocking_gap", Detail: "blocking Gap prevents continuation", Source: gapPtr}
			return card, nil
		}
	}
	runFresh, err := b.freshness(run.Document)
	if err != nil {
		return Card{}, err
	}
	if runFresh.State == "stale" {
		return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", run.Path+" is stale", nil)
	}
	if runFresh.State != "current" {
		card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Run freshness is unknown", Source: b.pointer(run)}
		return card, nil
	}
	missionBack, err := workspace.String(run.Document, "mission", true)
	if err != nil {
		return Card{}, err
	}
	mt, err := domain.ParseReference(missionBack)
	if err != nil || mt.Type != domain.Mission || mt.ID != mission.Document.Record.ID {
		return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "mission", "Run does not identify selected Mission", err)
	}
	cpRef, err := workspace.String(run.Document, "latest_checkpoint", false)
	if err != nil {
		return Card{}, err
	}
	if cpRef == "" {
		return b.governedRunContinuation(card, mission, run)
	}
	cp, err := b.Workspace.Lookup(cpRef, domain.Checkpoint)
	if err != nil {
		return Card{}, err
	}
	card.Pointers = append(card.Pointers, b.pointer(cp))
	cpPointer := b.pointer(cp)
	card.LatestCheckpoint = &cpPointer
	expectedCP, err := fingerprint(run.Document, "expected_checkpoint_fingerprint")
	if err != nil {
		return Card{}, err
	}
	if expectedCP != cp.Fingerprint {
		return Card{}, domain.NewRefusal(domain.RefusalInvalidFingerprint, "expected_checkpoint_fingerprint", "Run fingerprint does not match latest Checkpoint", nil)
	}
	cpFresh, err := b.freshness(cp.Document)
	if err != nil {
		return Card{}, err
	}
	if cpFresh.State == "stale" {
		return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", cp.Path+" is stale", nil)
	}
	if cpFresh.State != "current" {
		card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Checkpoint freshness is unknown", Source: b.pointer(cp)}
		return card, nil
	}
	cpRun, err := workspace.String(cp.Document, "run", true)
	if err != nil {
		return Card{}, err
	}
	rt, err := domain.ParseReference(cpRun)
	if err != nil || rt.Type != domain.Run || rt.ID != run.Document.Record.ID {
		return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "run", "Checkpoint does not identify selected Run", err)
	}
	evidenceRefs, err := workspace.Strings(cp.Document, "evidence", false)
	if err != nil {
		return Card{}, err
	}
	for _, evidenceRef := range evidenceRefs {
		evidence, e := b.Workspace.Lookup(evidenceRef, domain.Evidence)
		if e != nil {
			return Card{}, e
		}
		evFresh, e := b.freshness(evidence.Document)
		if e != nil {
			return Card{}, e
		}
		if evFresh.State == "stale" {
			return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", evidence.Path+" is stale", nil)
		}
		if evFresh.State != "current" {
			card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Evidence freshness is unknown", Source: b.pointer(evidence)}
			return card, nil
		}
		checkpointBack, e := workspace.String(evidence.Document, "checkpoint", true)
		if e != nil {
			return Card{}, e
		}
		back, e := domain.ParseReference(checkpointBack)
		if e != nil || back.Type != domain.Checkpoint || back.ID != cp.Document.Record.ID {
			return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "checkpoint", "Evidence does not identify selected Checkpoint", e)
		}
		appendPointer(&card.Pointers, b.pointer(evidence))
	}
	if governed, ok, governedErr := b.tryGovernedRunContinuation(card, mission, run); ok || governedErr != nil {
		return governed, governedErr
	}
	decisions := []discovery.Entry{}
	for _, d := range b.Workspace.OfType(domain.Decision) {
		if state := value(d.Document.Record.Status); state == "superseded" || state == "rejected" {
			continue
		}
		mref, e := workspace.String(d.Document, "mission", false)
		if e != nil {
			return Card{}, e
		}
		// Governance Decisions target records through their targets list. Only
		// legacy Scenario A continuation Decisions carry a direct mission field;
		// once selected, their operation is still validated strictly below.
		if mref == "" {
			continue
		}
		typed, e := domain.ParseReference(mref)
		if e != nil {
			return Card{}, e
		}
		if typed.Type == domain.Mission && typed.ID == mission.Document.Record.ID {
			decisions = append(decisions, d)
		}
	}
	if len(decisions) == 0 {
		card.OwnerGate = &OwnerGate{Code: "authorize_continuation", Detail: "no Decision authorizes the selected Mission chain", Source: b.pointer(mission)}
		return card, nil
	}
	if len(decisions) > 1 {
		return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "decision", "multiple Decisions claim the selected Mission", nil)
	}
	d := decisions[0]
	decidedBy, err := workspace.String(d.Document, "decided_by", true)
	if err != nil {
		return Card{}, err
	}
	if strings.TrimSpace(decidedBy) == "" {
		return Card{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "decided_by", "must identify an attributable decision maker", nil)
	}
	df, err := b.freshness(d.Document)
	if err != nil {
		return Card{}, err
	}
	if df.State == "stale" {
		return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", d.Path+" is stale", nil)
	}
	if df.State != "current" {
		card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: "Decision freshness is unknown", Source: b.pointer(d)}
		return card, nil
	}
	expectedMission, err := fingerprint(d.Document, "expected_mission_fingerprint")
	if err != nil {
		return Card{}, err
	}
	if expectedMission != mission.Fingerprint {
		return Card{}, domain.NewRefusal(domain.RefusalInvalidFingerprint, "expected_mission_fingerprint", "Decision fingerprint does not match current Mission", nil)
	}
	op, err := workspace.String(d.Document, "operation", true)
	if err != nil {
		return Card{}, err
	}
	if op != "resume" {
		return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "operation", "Decision operation must be exactly resume", nil)
	}
	targetRef, err := workspace.String(d.Document, "target", true)
	if err != nil {
		return Card{}, err
	}
	target, err := b.Workspace.Lookup(targetRef, domain.Run)
	if err != nil {
		return Card{}, err
	}
	if target.Document.Record.ID != run.Document.Record.ID {
		return Card{}, domain.NewRefusal(domain.RefusalConflictingAuthority, "target", "Decision target is not selected Run", nil)
	}
	card.Pointers = append(card.Pointers, b.pointer(d))
	card.Continuation = &Continuation{Operation: op, Target: b.pointer(target), AuthorizedBy: b.pointer(d)}
	return card, nil
}

func (b Builder) governedRunContinuation(card Card, mission, run discovery.Entry) (Card, error) {
	governed, ok, err := b.tryGovernedRunContinuation(card, mission, run)
	if err != nil || ok {
		return governed, err
	}
	card.OwnerGate = &OwnerGate{Code: "authorize_continuation", Detail: "active Mission has no owner authorization for its current Run", Source: b.pointer(mission)}
	return card, nil
}

func (b Builder) tryGovernedRunContinuation(card Card, mission, run discovery.Entry) (Card, bool, error) {
	authorizationRef, err := workspace.String(mission.Document, "last_authorization", false)
	if err != nil || authorizationRef == "" {
		return card, false, err
	}
	authorization, err := b.Workspace.Lookup(authorizationRef, domain.Decision)
	if err != nil {
		return Card{}, true, err
	}
	role, err := workspace.String(authorization.Document, "actor_role", true)
	if err != nil || role != "owner" {
		return Card{}, true, domain.NewRefusal(domain.RefusalUnauthorized, "actor_role", "current Run requires an owner Decision", err)
	}
	operation, err := workspace.String(authorization.Document, "operation", true)
	if err != nil || operation != "mission.transition.active" {
		return card, false, nil
	}
	targets, err := workspace.Strings(authorization.Document, "targets", true)
	if err != nil {
		return Card{}, true, err
	}
	missionRef := string(domain.Mission) + ":" + mission.Document.Record.ID.String()
	if !stringSliceContains(targets, missionRef) {
		return Card{}, true, domain.NewRefusal(domain.RefusalUnauthorized, "targets", "activation Decision does not target the selected Mission", nil)
	}
	freshness, err := b.freshness(authorization.Document)
	if err != nil {
		return Card{}, true, err
	}
	if freshness.State != "current" {
		return Card{}, true, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", authorization.Path+" is not current", nil)
	}
	card.Pointers = append(card.Pointers, b.pointer(authorization))
	card.Continuation = &Continuation{Operation: "continue-run", Target: b.pointer(run), AuthorizedBy: b.pointer(authorization)}
	return card, true, nil
}

func (b Builder) Detail(ref string, noun domain.RecordType) (Card, error) {
	entry, err := b.Workspace.Lookup(ref, noun)
	if err != nil {
		return Card{}, err
	}
	card, err := b.base(entry)
	if err != nil {
		return Card{}, err
	}
	if card.Freshness.State == "stale" {
		return Card{}, domain.NewRefusal(domain.RefusalStaleRequired, "freshness", entry.Path+" freshness basis is stale", nil)
	}
	if card.Freshness.State == "unknown" {
		card.OwnerGate = &OwnerGate{Code: "confirm_freshness", Detail: string(noun) + " freshness is unknown", Source: b.pointer(entry)}
	}
	for _, field := range relationshipFields(noun) {
		value, e := workspace.String(entry.Document, field, false)
		if e != nil {
			return Card{}, e
		}
		if value == "" {
			continue
		}
		typed, e := domain.ParseReference(value)
		if e != nil {
			return Card{}, e
		}
		target, e := b.Workspace.Lookup(value, typed.Type)
		if e != nil {
			return Card{}, e
		}
		card.Pointers = append(card.Pointers, b.pointer(target))
	}
	if noun == domain.Checkpoint {
		refs, e := workspace.Strings(entry.Document, "evidence", false)
		if e != nil {
			return Card{}, e
		}
		for _, ref := range refs {
			target, e := b.Workspace.Lookup(ref, domain.Evidence)
			if e != nil {
				return Card{}, e
			}
			card.Pointers = append(card.Pointers, b.pointer(target))
		}
	}
	if noun == domain.Decision {
		op, e := workspace.String(entry.Document, "operation", true)
		if e != nil {
			return Card{}, e
		}
		if op == "resume" {
			if _, e := workspace.String(entry.Document, "decided_by", true); e != nil {
				return Card{}, e
			}
			expected, e := fingerprint(entry.Document, "expected_mission_fingerprint")
			if e != nil {
				return Card{}, e
			}
			missionRef, e := workspace.String(entry.Document, "mission", true)
			if e != nil {
				return Card{}, e
			}
			mission, e := b.Workspace.Lookup(missionRef, domain.Mission)
			if e != nil {
				return Card{}, e
			}
			state := value(entry.Document.Record.Status)
			if state != "superseded" && state != "rejected" && expected != mission.Fingerprint {
				return Card{}, domain.NewRefusal(domain.RefusalInvalidFingerprint, "expected_mission_fingerprint", "Decision fingerprint does not match current Mission", nil)
			}
		} else {
			if role, e := workspace.String(entry.Document, "actor_role", true); e != nil || role != "owner" {
				return Card{}, domain.NewRefusal(domain.RefusalUnauthorized, "actor_role", "governance Decision must identify owner authority", e)
			}
			if _, e := workspace.Strings(entry.Document, "targets", true); e != nil {
				return Card{}, e
			}
			if _, e := workspace.Strings(entry.Document, "expected_fingerprints", true); e != nil {
				return Card{}, e
			}
		}
	}
	return card, nil
}

func (b Builder) Gaps(scope string) (List, error) {
	mission, err := b.Workspace.Lookup(scope, domain.Mission)
	if err != nil {
		return List{}, err
	}
	out := List{Items: []Card{}}
	for _, gap := range b.Workspace.OfType(domain.Gap) {
		ref, e := workspace.String(gap.Document, "scope", true)
		if e != nil {
			return List{}, e
		}
		typed, e := domain.ParseReference(ref)
		if e != nil {
			return List{}, e
		}
		if typed.Type != domain.Mission {
			return List{}, domain.NewRefusal(domain.RefusalNounMismatch, "scope", "Gap scope must be Mission", nil)
		}
		if typed.ID == mission.Document.Record.ID {
			card, e := b.Detail(gap.Document.Record.ID.String(), domain.Gap)
			if e != nil {
				return List{}, e
			}
			out.Items = append(out.Items, card)
		}
	}
	return out, nil
}

func (b Builder) Validate(scope string) (Validation, error) {
	var entries []discovery.Entry
	if scope == "project" {
		entries = b.Workspace.Entries
		if _, err := b.Project(); err != nil {
			return Validation{}, err
		}
		for _, entry := range entries {
			noun := entry.Document.Record.Type
			if noun == domain.Proposal || noun == domain.Mission || noun == domain.Anchor {
				continue
			}
			if _, err := b.Detail(entry.Document.Record.ID.String(), noun); err != nil {
				return Validation{}, err
			}
		}
	} else {
		mission, err := b.Workspace.Lookup(scope, domain.Mission)
		if err != nil {
			return Validation{}, domain.NewStateRefusal(domain.RefusalInvalidScope, "scope", err.Error(), "project or an exact Mission reference", scope, "spectacular workspace validate project", err)
		}
		entries = []discovery.Entry{mission}
		if _, err := b.Mission(scope); err != nil {
			return Validation{}, err
		}
	}
	sources := make([]Source, 0, len(entries))
	for _, e := range entries {
		sources = append(sources, Source{Path: e.Path, Fingerprint: e.Fingerprint})
	}
	return Validation{Scope: scope, Valid: true, Records: len(entries), Sources: sources}, nil
}

func (b Builder) base(entry discovery.Entry) (Card, error) {
	fresh, err := b.freshness(entry.Document)
	if err != nil {
		return Card{}, err
	}
	title := ""
	if entry.Document.Record.Title != nil {
		title = *entry.Document.Record.Title
	}
	ref := humanRef(entry)
	if ref == "" {
		ref = string(entry.Document.Record.Type) + ":" + entry.Document.Record.ID.String()
	}
	state := ""
	if entry.Document.Record.Status != nil {
		state = *entry.Document.Record.Status
	}
	return Card{Noun: string(entry.Document.Record.Type), Ref: ref, ID: entry.Document.Record.ID.String(), Title: title, State: state, Freshness: fresh, Source: Source{Path: entry.Path, Fingerprint: entry.Fingerprint}, Sources: []Pointer{}, Pointers: []Pointer{}, Gaps: []Pointer{}, Conflicts: []string{}, Omissions: []string{}}, nil
}
func (b Builder) freshness(doc *workspace.Document) (Freshness, error) {
	checkedText, err := workspace.String(doc, "freshness_checked_at", true)
	if err != nil {
		return Freshness{}, err
	}
	validText, err := workspace.String(doc, "freshness_valid_until", true)
	if err != nil {
		return Freshness{}, err
	}
	sourceRef, err := workspace.String(doc, "freshness_source", true)
	if err != nil {
		return Freshness{}, err
	}
	expected, err := fingerprint(doc, "freshness_source_fingerprint")
	if err != nil {
		return Freshness{}, err
	}
	checked, err := time.Parse(time.RFC3339, checkedText)
	if err != nil {
		return Freshness{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "freshness_checked_at", "must be RFC3339", err)
	}
	valid, err := time.Parse(time.RFC3339, validText)
	if err != nil {
		return Freshness{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "freshness_valid_until", "must be RFC3339", err)
	}
	if !valid.After(checked) {
		return Freshness{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "freshness_valid_until", "must be after freshness_checked_at", nil)
	}
	entry, fp, err := b.Workspace.Source(sourceRef)
	if err != nil {
		return Freshness{}, err
	}
	source := Source{Path: entry.Path, Fingerprint: fp}
	var sourcePointer *Pointer
	if entry.Document != nil {
		p := b.pointer(entry)
		sourcePointer = &p
	}
	now := b.Now
	if now == nil {
		now = time.Now
	}
	at := now().UTC()
	state := "current"
	if at.Before(checked) {
		state = "unknown"
	} else if at.After(valid) {
		state = "stale"
	}
	if expected != fp {
		state = "stale"
	}
	return Freshness{State: state, Basis: "explicit-validity-window+source-fingerprint", Source: source, Pointer: sourcePointer, CheckedAt: checked.UTC().Format(time.RFC3339), ValidUntil: valid.UTC().Format(time.RFC3339), ExpectedFingerprint: expected, ActualFingerprint: fp}, nil
}
func fingerprint(doc *workspace.Document, field string) (string, error) {
	v, err := workspace.String(doc, field, true)
	if err != nil {
		return "", err
	}
	if len(v) != 64 {
		return "", domain.NewRefusal(domain.RefusalInvalidFingerprint, field, "expected lowercase SHA-256 hex", nil)
	}
	if _, err := hex.DecodeString(v); err != nil || strings.ToLower(v) != v {
		return "", domain.NewRefusal(domain.RefusalInvalidFingerprint, field, "expected lowercase SHA-256 hex", err)
	}
	return v, nil
}
func (b Builder) pointer(e discovery.Entry) Pointer {
	ref := string(e.Document.Record.Type) + ":" + e.Document.Record.ID.String()
	human := humanRef(e)
	if e.Document.Record.Type == domain.Anchor {
		if human == "" || human == "PROJECT" {
			ref = "project"
			human = ""
		}
	}
	noun := strings.ToLower(string(e.Document.Record.Type))
	p := Pointer{Noun: noun, Ref: ref, HumanRef: human, Path: e.Path, Fingerprint: e.Fingerprint}
	if b.ShowCommand != nil {
		showRef := ref
		if human != "" {
			showRef = human
		}
		p.ShowCommand, _ = b.ShowCommand(e.Document.Record.Type, showRef)
	}
	return p
}

func humanRef(e discovery.Entry) string {
	if ref, err := workspace.String(e.Document, "human_ref", false); err == nil && ref != "" {
		return ref
	}
	return ""
}

func appendPointer(pointers *[]Pointer, candidate Pointer) {
	for _, existing := range *pointers {
		if existing.Ref == candidate.Ref {
			return
		}
	}
	*pointers = append(*pointers, candidate)
}

func value(text *string) string {
	if text == nil {
		return ""
	}
	return *text
}

func stringSliceContains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
func relationshipFields(noun domain.RecordType) []string {
	switch noun {
	case domain.Gap:
		return []string{"scope"}
	case domain.Run:
		return []string{"mission", "latest_checkpoint"}
	case domain.Checkpoint:
		return []string{"run"}
	case domain.Evidence:
		return []string{"checkpoint"}
	case domain.Decision:
		return []string{"mission", "target"}
	case domain.Objective:
		return []string{"mission"}
	case domain.Handoff:
		return []string{"mission", "objective", "run", "dispatch", "authorization", "supersedes"}
	case domain.Assessment:
		return []string{"mission", "authorization"}
	}
	return nil
}
func (b Builder) generationBasis() string {
	parts := make([]string, len(b.Workspace.Entries))
	for i, e := range b.Workspace.Entries {
		parts[i] = e.Path + ":" + e.Fingerprint
	}
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}
