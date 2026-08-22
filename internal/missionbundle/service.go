package missionbundle

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"github.com/google/uuid"
	"go.yaml.in/yaml/v3"
)

type Service struct {
	Workspace        *discovery.Workspace
	Now              func() time.Time
	ApplyTransaction func(root, key string, changes []governance.FileChange) error
}

type Plan struct {
	Type         string      `yaml:"type"`
	Title        string      `yaml:"title"`
	Owner        string      `yaml:"owner"`
	Contract     Binding     `yaml:"contract"`
	Outcome      string      `yaml:"outcome"`
	Request      *Request    `yaml:"request,omitempty"`
	Review       string      `yaml:"review"`
	Completion   []Criterion `yaml:"completion"`
	Objectives   []Objective `yaml:"objectives"`
	Authority    Authority   `yaml:"authority"`
	Scope        Scope       `yaml:"scope"`
	RepairBudget int         `yaml:"repair_budget"`
	Dependencies []string    `yaml:"dependencies"`
	Gaps         []string    `yaml:"gaps"`
	Stops        []string    `yaml:"stops"`
	Fallbacks    []Fallback  `yaml:"fallbacks,omitempty"`
	AfterMission []string    `yaml:"after_mission,omitempty"`
	// ResolvesGaps names Gaps on the bound Contract that this Mission closes at
	// completion, with the resolution text it will write. Frozen at activation so
	// the authority to amend a Contract cannot be acquired afterwards.
	ResolvesGaps []ResolvedGap `yaml:"resolves_gaps,omitempty"`
	Body         string        `yaml:"-"`
}

type ReviewDraft struct {
	Type     string `yaml:"type"`
	Title    string `yaml:"title"`
	Status   string `yaml:"status"`
	Reviewed struct {
		Commit                string `yaml:"commit"`
		Tree                  string `yaml:"tree"`
		ActivationFingerprint string `yaml:"activation_fingerprint"`
	} `yaml:"reviewed"`
	Reviewer Reviewer `yaml:"reviewer"`
	Claims   []struct {
		Claim   string `yaml:"claim"`
		Verdict string `yaml:"verdict"`
	} `yaml:"claims"`
	Findings    []string `yaml:"findings"`
	Limitations []string `yaml:"limitations"`
	Body        string   `yaml:"-"`
}

func ReadPlan(path string, stdin []byte) (Plan, []byte, error) {
	data, err := readInput(path, stdin)
	if err != nil {
		return Plan{}, nil, err
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return Plan{}, nil, err
	}
	var plan Plan
	if err := yaml.Unmarshal(frontmatter, &plan); err != nil {
		return Plan{}, nil, invalidCause("input", "decode Mission plan frontmatter", err)
	}
	if plan.Type != "MissionPlan" {
		return Plan{}, nil, invalid("type", "Mission start input must declare type: MissionPlan")
	}
	plan.Body = body
	return plan, data, nil
}

func (s Service) Start(plan Plan, raw []byte) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.start(plan, raw)
}

func (s Service) start(plan Plan, raw []byte) (Result, error) {
	if s.Workspace == nil {
		return Result{}, invalid("workspace", "workspace is required")
	}
	startDigest := sha256.Sum256(raw)
	startKey := "sha256:" + hex.EncodeToString(startDigest[:])
	for _, entry := range s.Workspace.OfType(domain.Mission) {
		if existing, _ := workspace.String(entry.Document, "start_key", false); existing == startKey {
			ref, _ := compactRef(entry.Document)
			return Result{Operation: "mission.start", Ref: ref, Path: entry.Path, Fingerprint: entry.Fingerprint}, nil
		}
	}
	if err := validatePlan(plan); err != nil {
		return Result{}, err
	}
	plan.Dependencies = presentStrings(plan.Dependencies)
	plan.Gaps = presentStrings(plan.Gaps)
	plan.Stops = presentStrings(plan.Stops)
	contract, err := resolveContract(s.Workspace, plan.Contract.Ref)
	if err != nil {
		return Result{}, err
	}
	commit, branch, baselineAt, err := gitBaseline(s.Workspace.Root)
	if err != nil {
		return Result{}, err
	}
	missionID, err := stableID(baselineAt, startKey+":mission")
	if err != nil {
		return Result{}, err
	}
	missionRef := nextMissionRef(s.Workspace)
	now := s.now()
	objectives := make([]Objective, len(plan.Objectives))
	for i, source := range plan.Objectives {
		id, idErr := stableID(baselineAt, startKey+":objective:"+strconv.Itoa(i+1))
		if idErr != nil {
			return Result{}, idErr
		}
		source.Ref = "O" + strconv.Itoa(i+1)
		source.ID = id.String()
		if source.Status == "" {
			source.Status = "pending"
		}
		objectives[i] = source
	}
	runID, err := stableID(baselineAt, startKey+":run:R1")
	if err != nil {
		return Result{}, err
	}
	doc := &workspace.Document{Record: domain.Record{
		Type: domain.Mission, ID: missionID, Title: stringPtr(plan.Title), Status: stringPtr("active"),
		Created: stringPtr(now), Updated: stringPtr(now),
	}, Unknown: map[string]*yaml.Node{}, Body: plan.Body}
	workspace.SetString(doc, "ref", missionRef)
	workspace.SetString(doc, "owner", plan.Owner)
	workspace.SetValue(doc, "contract", contract)
	workspace.SetValue(doc, "baseline", Baseline{Commit: commit, Branch: branch})
	workspace.SetString(doc, "outcome", plan.Outcome)
	if plan.Request != nil {
		workspace.SetValue(doc, "request", plan.Request)
	}
	workspace.SetString(doc, "review", plan.Review)
	workspace.SetValue(doc, "completion", plan.Completion)
	workspace.SetValue(doc, "objectives", objectives)
	run := Run{Ref: "R1", ID: runID.String(), Status: "active", Operator: plan.Owner, StartedAt: now, CurrentObjective: objectives[0].Ref}
	workspace.SetValue(doc, "run", run)
	workspace.SetValue(doc, "validation", Validation{Schema: Schema, Mode: "cli"})
	workspace.SetValue(doc, "authority", plan.Authority)
	workspace.SetValue(doc, "scope", plan.Scope)
	workspace.SetInt(doc, "repair_budget", plan.RepairBudget)
	workspace.SetStrings(doc, "dependencies", plan.Dependencies)
	workspace.SetStrings(doc, "gaps", plan.Gaps)
	workspace.SetStrings(doc, "stops", plan.Stops)
	if len(plan.Fallbacks) > 0 {
		workspace.SetValue(doc, "fallbacks", plan.Fallbacks)
	}
	if len(plan.AfterMission) > 0 {
		workspace.SetStrings(doc, "after_mission", plan.AfterMission)
	}
	if len(plan.ResolvesGaps) > 0 {
		workspace.SetValue(doc, "resolves_gaps", plan.ResolvesGaps)
	}
	workspace.SetString(doc, "start_key", startKey)
	temporary := &Bundle{Outcome: plan.Outcome, Request: plan.Request, Review: plan.Review, Completion: plan.Completion, Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget, Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops, Fallbacks: plan.Fallbacks, AfterMission: plan.AfterMission, ResolvesGaps: plan.ResolvesGaps}
	fingerprint, err := FrozenFingerprint(temporary)
	if err != nil {
		return Result{}, err
	}
	workspace.SetValue(doc, "activation", Activation{By: plan.Owner, At: now, Fingerprint: fingerprint})
	slug := humanlayout.Slug(plan.Title)
	path := filepath.ToSlash(filepath.Join(".spectacular", "missions", missionRef+"-"+slug, missionRef+"-"+slug+".md"))
	candidate := &Bundle{
		ID: missionID.String(), Ref: missionRef, Title: plan.Title, Status: "active", Owner: plan.Owner,
		Contract: contract, Baseline: &Baseline{Commit: commit, Branch: branch}, Outcome: plan.Outcome,
		Request: plan.Request, Review: plan.Review, Completion: plan.Completion, Objectives: objectives, Run: &run,
		Activation: &Activation{By: plan.Owner, At: now, Fingerprint: fingerprint}, Validation: Validation{Schema: Schema, Mode: "cli"},
		Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget,
		Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops, Fallbacks: plan.Fallbacks, AfterMission: plan.AfterMission,
		ResolvesGaps: plan.ResolvesGaps, Path: path,
		document: doc,
	}
	for _, check := range registry {
		if check.name == "safe-file-layout" {
			continue
		}
		if err := check.run(s.Workspace, candidate); err != nil {
			return Result{}, err
		}
	}
	return s.apply("mission.start:"+startKey, []*workspace.Document{doc}, map[domain.ID]string{missionID: path}, "mission.start", missionRef, path)
}

func (s Service) Show(ref string) (*Bundle, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return nil, err
	}
	state := bundle.Derive()
	bundle.State = &state
	return bundle, nil
}

func (s Service) Check(ref string) (Check, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Check{}, err
	}
	return Validate(s.Workspace, bundle)
}

func (s Service) Objective(ref string) (Objective, *Bundle, error) {
	missionRef, local, err := scopedRef(ref)
	if err != nil {
		return Objective{}, nil, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Objective{}, nil, err
	}
	for _, objective := range bundle.Objectives {
		if objective.Ref == local {
			return objective, bundle, nil
		}
	}
	return Objective{}, nil, domain.NewRefusal(domain.RefusalRecordNotFound, "ref", "Objective does not exist in Mission", nil)
}

func (s Service) PromoteObjective(ref string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.promoteObjective(ref)
}

func (s Service) promoteObjective(ref string) (Result, error) {
	objective, bundle, err := s.Objective(ref)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy {
		return Result{}, invalid("mission", "legacy Mission is read-only")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	if objective.File != "" {
		return Result{Operation: "objective.promote", Ref: ref, Path: filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), objective.File))}, nil
	}
	id, _ := domain.ParseID(objective.ID)
	doc := &workspace.Document{Record: domain.Record{Type: domain.Objective, ID: id, Title: stringPtr(objective.Outcome), Status: stringPtr(objective.Status)}, Unknown: map[string]*yaml.Node{}, Body: "# Objective\n\n" + objective.Outcome + "\n"}
	workspace.SetString(doc, "ref", objective.Ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetString(doc, "outcome", objective.Outcome)
	workspace.SetStrings(doc, "after", objective.After)
	workspace.SetStrings(doc, "claims", objective.Claims)
	relative := filepath.ToSlash(filepath.Join("objectives", objective.Ref+"-"+humanlayout.Slug(objective.Outcome)+".md"))
	for i := range bundle.Objectives {
		if bundle.Objectives[i].Ref == objective.Ref {
			bundle.Objectives[i] = Objective{Ref: objective.Ref, ID: objective.ID, File: relative}
		}
	}
	workspace.SetValue(bundle.document, "objectives", bundle.Objectives)
	bundle.document.Record.Updated = stringPtr(s.now())
	missionID := bundle.document.Record.ID
	path := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), relative))
	paths := map[domain.ID]string{missionID: bundle.Path, id: path}
	return s.apply("objective.promote:"+bundle.ID+":"+objective.ID, []*workspace.Document{bundle.document, doc}, paths, "objective.promote", ref, path)
}

func (s Service) FinishObjective(ref string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.finishObjective(ref)
}

func (s Service) finishObjective(ref string) (Result, error) {
	objective, bundle, err := s.Objective(ref)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy {
		return Result{}, invalid("mission", "legacy Mission is read-only")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	if objective.Status == "implemented" {
		return Result{Operation: "objective.finish", Ref: ref, Path: bundle.Path}, nil
	}
	states := map[string]string{}
	for _, item := range bundle.Objectives {
		states[item.Ref] = item.Status
	}
	for _, dependency := range objective.After {
		if states[dependency] != "implemented" {
			return Result{}, invalid("objectives.after", "finish dependencies before this Objective")
		}
	}
	var docs []*workspace.Document
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	if objective.File == "" {
		for i := range bundle.Objectives {
			if bundle.Objectives[i].Ref == objective.Ref {
				bundle.Objectives[i].Status = "implemented"
			}
		}
		workspace.SetValue(bundle.document, "objectives", bundle.Objectives)
	} else {
		absolute, _ := containedFile(filepath.Dir(bundle.entry.Absolute), objective.File)
		doc, readErr := workspace.ReadFile(absolute)
		if readErr != nil {
			return Result{}, readErr
		}
		doc.Record.Status = stringPtr("implemented")
		docs = append(docs, doc)
		paths[doc.Record.ID] = filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), objective.File))
	}
	bundle.document.Record.Updated = stringPtr(s.now())
	docs = append(docs, bundle.document)
	return s.apply("objective.finish:"+bundle.ID+":"+objective.ID, docs, paths, "objective.finish", ref, bundle.Path)
}

func (s Service) Run(ref string) (Run, *Bundle, error) {
	missionRef, local, err := scopedRef(ref)
	if err != nil {
		return Run{}, nil, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Run{}, nil, err
	}
	for _, run := range allRuns(bundle) {
		if run.Ref == local {
			return run, bundle, nil
		}
	}
	return Run{}, nil, domain.NewRefusal(domain.RefusalRecordNotFound, "ref", "Run does not exist in Mission", nil)
}

func (s Service) StartRun(targetRef, title string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.startRun(targetRef, title)
}

func (s Service) startRun(targetRef, title string) (Result, error) {
	parts := strings.Split(targetRef, "/")
	missionRef := parts[0]
	targetObjective := ""
	if len(parts) > 1 {
		targetObjective = parts[1]
	}

	if err := CheckPassiveGitState(s.Workspace.Root); err != nil {
		return Result{}, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("run", "new Run requires an active compact Mission")
	}
	if strings.TrimSpace(title) == "" {
		if targetObjective != "" {
			title = "Run for " + targetObjective
		} else {
			return Result{}, invalid("run", "new Run requires an active compact Mission and title")
		}
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}

	if targetObjective != "" {
		var currentObj *Objective
		for i := range bundle.Objectives {
			if bundle.Objectives[i].Ref == targetObjective || bundle.Objectives[i].ID == targetObjective {
				currentObj = &bundle.Objectives[i]
				break
			}
		}
		if currentObj == nil {
			return Result{}, invalid("objective", fmt.Sprintf("objective %s not found in mission %s", targetObjective, missionRef))
		}
		for _, r := range allRuns(bundle) {
			if (r.CurrentObjective == targetObjective || r.Objective == targetObjective) && (r.Status == "active" || r.Status == "paused" || r.Status == "blocked" || r.Status == "awaiting-review") {
				return Result{}, domain.NewRefusal(domain.RefusalCollision, "objective", fmt.Sprintf("objective %s already has an active run reserving it (%s in state %s)", targetObjective, r.Ref, r.Status), nil)
			}
		}
		// Upstream dependency locking
		if len(currentObj.After) > 0 {
			for _, depRef := range currentObj.After {
				for _, r := range allRuns(bundle) {
					if (r.CurrentObjective == depRef || r.Objective == depRef) && (r.Status == "blocked" || r.Status == "stopped") {
						return Result{}, domain.NewRefusal(domain.RefusalCollision, "dependency", fmt.Sprintf("cannot start run on %s: upstream dependency %s is in state %s; resolve blocker with owner Decision first", targetObjective, depRef, r.Status), nil)
					}
				}
			}
		}
	}

	runs := allRuns(bundle)
	if len(runs) == 0 && targetObjective == "" {
		return Result{}, invalid("run", "Mission has no current Run")
	}
	if len(runs) > 0 {
		last := runs[len(runs)-1]
		if last.Status == "active" && last.Title == title {
			path := bundle.Path
			if last.File != "" {
				path = filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), last.File))
			}
			return Result{Operation: "run.start", Ref: bundle.Ref + "/" + last.Ref, Path: path}, nil
		}
	}
	docs := []*workspace.Document{}
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	pointers := make([]Run, 0, len(runs)+1)
	for _, run := range runs {
		if run.Status != "completed" && run.Status != "stopped" {
			run.Status = "completed"
		}
		if run.File == "" {
			doc, path, makeErr := runDocument(bundle, run, runTitle(run))
			if makeErr != nil {
				return Result{}, makeErr
			}
			docs = append(docs, doc)
			paths[doc.Record.ID] = path
			run = Run{Ref: run.Ref, ID: run.ID, File: strings.TrimPrefix(path, filepath.ToSlash(filepath.Dir(bundle.Path))+"/"), Status: run.Status}
		} else {
			absolute, pathErr := containedFile(filepath.Dir(bundle.entry.Absolute), run.File)
			if pathErr != nil {
				return Result{}, pathErr
			}
			doc, readErr := workspace.ReadFile(absolute)
			if readErr != nil {
				return Result{}, readErr
			}
			doc.Record.Status = stringPtr(run.Status)
			docs = append(docs, doc)
			paths[doc.Record.ID] = filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), run.File))
			run = Run{Ref: run.Ref, ID: run.ID, File: run.File, Status: run.Status}
		}
		pointers = append(pointers, run)
	}
	nextRef := "R" + strconv.Itoa(len(runs)+1)
	nextID, err := stableID(bundle.Activation.At, "run.start:"+bundle.ID+":"+nextRef+":"+title)
	if err != nil {
		return Result{}, err
	}
	current := targetObjective
	if current == "" {
		current = nextObjective(bundle.Objectives)
	}
	if current == "" && len(bundle.Objectives) > 0 {
		current = bundle.Objectives[len(bundle.Objectives)-1].Ref
	}
	next := Run{Ref: nextRef, ID: nextID.String(), Title: title, Status: "active", Operator: bundle.Owner, StartedAt: s.now(), CurrentObjective: current, Objective: current}
	doc, nextPath, err := runDocument(bundle, next, title)
	if err != nil {
		return Result{}, err
	}
	docs = append(docs, doc)
	paths[doc.Record.ID] = nextPath
	pointers = append(pointers, Run{Ref: next.Ref, ID: next.ID, File: strings.TrimPrefix(nextPath, filepath.ToSlash(filepath.Dir(bundle.Path))+"/"), Status: "active"})
	workspace.Delete(bundle.document, "run")
	workspace.SetValue(bundle.document, "runs", pointers)
	bundle.document.Record.Updated = stringPtr(s.now())
	docs = append(docs, bundle.document)
	return s.apply("run.start:"+bundle.ID+":"+next.ID, docs, paths, "run.start", bundle.Ref+"/"+next.Ref, nextPath)
}

func (s Service) RecordReview(missionRef, path string, stdin []byte) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.recordReview(missionRef, path, stdin)
}

func (s Service) recordReview(missionRef, path string, stdin []byte) (Result, error) {
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "review recording requires an active compact Mission")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	data, err := readInput(path, stdin)
	if err != nil {
		return Result{}, err
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return Result{}, err
	}
	var draft ReviewDraft
	if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
		return Result{}, invalidCause("input", "decode ReviewDraft frontmatter", err)
	}
	if draft.Type != "ReviewDraft" || draft.Title == "" || draft.Status != "passed" || len(draft.Claims) == 0 {
		return Result{}, invalid("review", "requires type ReviewDraft, title, status passed, and claim verdicts")
	}
	known := map[string]bool{}
	for _, criterion := range bundle.Completion {
		known[criterion.Claim] = true
	}
	seen := map[string]bool{}
	for _, claim := range draft.Claims {
		if !known[claim.Claim] || claim.Verdict != "pass" {
			return Result{}, invalid("review.claims", "every frozen claim must have a pass verdict")
		}
		seen[claim.Claim] = true
	}
	if len(seen) != len(known) {
		return Result{}, invalid("review.claims", "review must cover every frozen claim exactly")
	}
	if bundle.Review == "independent" && (draft.Reviewer.RelationToOperator != "independent" || draft.Reviewer.ImplementedReviewedScope || draft.Reviewer.Actor == "" || draft.Reviewer.Operator == "" || draft.Reviewer.Actor == draft.Reviewer.Operator || draft.Reviewer.IndependenceBasis == "" || len(draft.Reviewer.Evidence) == 0) {
		return Result{}, invalid("review.reviewer", "independent review requires distinct reviewer/operator identities, a non-implementation statement, basis, and attributable evidence")
	}
	if draft.Reviewed.ActivationFingerprint != bundle.Activation.Fingerprint || !commitPattern.MatchString(draft.Reviewed.Commit) {
		return Result{}, invalid("review.reviewed", "review must bind exact commit, tree, and Mission activation fingerprint")
	}
	if draft.Reviewed.Tree == "" {
		resolvedTree, err := resolveCommitTree(s.Workspace.Root, draft.Reviewed.Commit)
		if err != nil {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "review.reviewed.commit", "reviewed commit does not exist in this repository", "existing exact commit", draft.Reviewed.Commit, "review the committed tree, then record its exact commit and tree", err)
		}
		draft.Reviewed.Tree = resolvedTree
	} else if !commitPattern.MatchString(draft.Reviewed.Tree) {
		return Result{}, invalid("review.reviewed", "review must bind exact commit, tree, and Mission activation fingerprint")
	}
	if err := verifyReviewedGit(s.Workspace.Root, draft.Reviewed.Commit, draft.Reviewed.Tree, "review"); err != nil {
		return Result{}, err
	}
	for _, existing := range bundle.Reviews {
		if existing.Document != nil && existing.Document.Reviewed.Commit == draft.Reviewed.Commit {
			return Result{
				Operation: "review.record",
				Ref:       bundle.Ref + "/" + existing.Ref,
				Path:      filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), existing.File)),
			}, nil
		}
	}
	id, err := stableID(bundle.Activation.At, "review.record:"+bundle.ID+":"+draft.Reviewed.Commit)
	if err != nil {
		return Result{}, err
	}
	ref := "RV" + strconv.Itoa(len(bundle.Reviews)+1)
	now := s.now()
	doc := &workspace.Document{Record: domain.Record{Type: domain.Review, ID: id, Title: stringPtr(draft.Title), Status: stringPtr("passed"), Created: stringPtr(now)}, Unknown: map[string]*yaml.Node{}, Body: body}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetValue(doc, "reviewed", draft.Reviewed)
	workspace.SetValue(doc, "reviewer", draft.Reviewer)
	workspace.SetValue(doc, "claims", draft.Claims)
	workspace.SetStrings(doc, "findings", draft.Findings)
	workspace.SetStrings(doc, "limitations", draft.Limitations)
	reviewPath, relative, err := s.missionRecordPath(bundle, doc, ref)
	if err != nil {
		return Result{}, err
	}
	bundle.Reviews = append(bundle.Reviews, ReviewPointer{Ref: ref, ID: id.String(), File: relative, Verdict: "pass"})
	workspace.SetValue(bundle.document, "reviews", bundle.Reviews)
	bundle.document.Record.Updated = stringPtr(now)
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path, id: reviewPath}
	return s.apply("review.record:"+bundle.ID+":"+id.String(), []*workspace.Document{bundle.document, doc}, paths, "review.record", bundle.Ref+"/"+ref, reviewPath)
}

func (s Service) RecordHandoff(missionRef, path, sender string, stdin []byte) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.recordHandoff(missionRef, path, sender, stdin)
}

func (s Service) recordHandoff(missionRef, path, sender string, stdin []byte) (Result, error) {
	if err := CheckPassiveGitState(s.Workspace.Root); err != nil {
		return Result{}, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "handoff recording requires an active compact Mission")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	data, err := readInput(path, stdin)
	if err != nil {
		return Result{}, err
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return Result{}, err
	}
	var draft HandoffDraft
	if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
		return Result{}, invalidCause("input", "decode HandoffDraft frontmatter", err)
	}
	if draft.Type != "HandoffDraft" || draft.Title == "" {
		return Result{}, invalid("handoff", "requires type HandoffDraft and a title")
	}
	// The asserted/assumed split is why this record exists, so an absent list is
	// refused at the boundary rather than defaulted to empty. Empty stays legal:
	// a sender who verified nothing is saying something.
	if draft.Asserted == nil {
		return Result{}, invalid("handoff.asserted", "a Handoff must state asserted, even as an empty list")
	}
	if draft.Assumed == nil {
		return Result{}, invalid("handoff.assumed", "a Handoff must state assumed, even as an empty list")
	}
	// --by names the sender of record. A draft may carry the same identity, but
	// the two disagreeing means the caller and the record would attribute the
	// delegation to different people.
	if draft.Sender.Actor != "" && draft.Sender.Actor != sender {
		return Result{}, domain.NewStateRefusal(domain.RefusalUnauthorized, "by",
			"handoff sender must match the identity recording it", draft.Sender.Actor, sender,
			"record the Handoff as its stated sender, or correct the draft", nil)
	}
	draft.Sender.Actor = sender
	if !commitPattern.MatchString(draft.Reviewed.Commit) {
		return Result{}, invalid("handoff.reviewed", "handoff must bind an exact commit and tree")
	}
	if draft.Reviewed.Tree == "" {
		resolvedTree, err := resolveCommitTree(s.Workspace.Root, draft.Reviewed.Commit)
		if err != nil {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "handoff.reviewed.commit", "reviewed commit does not exist in this repository", "existing exact commit", draft.Reviewed.Commit, "review the committed tree, then record its exact commit and tree", err)
		}
		draft.Reviewed.Tree = resolvedTree
	} else if !commitPattern.MatchString(draft.Reviewed.Tree) {
		return Result{}, invalid("handoff.reviewed", "handoff must bind an exact commit and tree")
	}
	// The git binding is verified against the real repository, so a Handoff cannot
	// point at a commit that does not exist or a tree that is not that commit's.
	// A receiver re-verifies from this binding, and a binding that was never
	// checked would send them to a state that never existed.
	if err := verifyReviewedGit(s.Workspace.Root, draft.Reviewed.Commit, draft.Reviewed.Tree, "handoff"); err != nil {
		return Result{}, err
	}
	if draft.Supersedes != "" {
		found := false
		for _, existing := range bundle.Handoffs {
			if existing.Ref == draft.Supersedes {
				found = true
				break
			}
		}
		if !found {
			return Result{}, invalid("handoff.supersedes", "supersedes must name a Handoff recorded on this Mission: "+draft.Supersedes)
		}
	}

	// Recording the same logical Handoff twice converges rather than duplicating.
	// Identity is derived from the Mission, the bound commit, and what is being
	// superseded, so a retry after a crash lands on the record already written.
	key := "handoff.record:" + bundle.ID + ":" + draft.Reviewed.Commit + ":" + draft.Supersedes
	id, err := stableID(bundle.Activation.At, key)
	if err != nil {
		return Result{}, err
	}
	for _, existing := range bundle.Handoffs {
		if existing.ID == id.String() {
			return Result{
				Operation: "handoff.record",
				Ref:       bundle.Ref + "/" + existing.Ref,
				Path:      filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), existing.File)),
			}, nil
		}
	}

	ref := "H" + strconv.Itoa(len(bundle.Handoffs)+1) + "-" + humanlayout.ShortKey(id)
	now := s.now()
	doc := &workspace.Document{
		Record:  domain.Record{Type: domain.Handoff, ID: id, Title: stringPtr(draft.Title), Created: stringPtr(now)},
		Unknown: map[string]*yaml.Node{}, Body: body,
	}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetValue(doc, "reviewed", draft.Reviewed)
	workspace.SetValue(doc, "sender", draft.Sender)
	workspace.SetString(doc, "task", draft.Task)
	workspace.SetStrings(doc, "asserted", *draft.Asserted)
	workspace.SetStrings(doc, "assumed", *draft.Assumed)
	workspace.SetStrings(doc, "stops", draft.Stops)
	workspace.SetStrings(doc, "returns", draft.Returns)
	if err := ValidateWritePaths(draft.Writes); err != nil {
		return Result{}, err
	}
	for _, existingPointer := range bundle.Handoffs {
		if existingPointer.Document != nil && existingPointer.Document.Supersedes == "" && existingPointer.Document.Ref != draft.Supersedes {
			for _, ew := range existingPointer.Document.Writes {
				for _, dw := range draft.Writes {
					if PathsOverlap(dw, ew) {
						return Result{}, domain.NewRefusal(domain.RefusalInvalidScope, "writes", fmt.Sprintf("cannot reserve write path %q: overlaps with active Handoff %s reservation %q", dw, existingPointer.Ref, ew), nil)
					}
				}
			}
		}
	}
	if draft.Supersedes != "" {
		workspace.SetString(doc, "supersedes", draft.Supersedes)
	}
	if len(draft.Writes) > 0 {
		workspace.SetStrings(doc, "writes", draft.Writes)
	}
	handoffPath, relative, err := s.missionRecordPath(bundle, doc, ref)
	if err != nil {
		return Result{}, err
	}

	// The record is validated before it is written, so the schema refuses a bad
	// Handoff at the command rather than leaving one on disk for a later read to
	// reject.
	candidate, err := decodeHandoff(doc, handoffPath)
	if err != nil {
		return Result{}, err
	}
	if err := validateHandoffContent(candidate, bundle); err != nil {
		return Result{}, err
	}

	bundle.Handoffs = append(bundle.Handoffs, HandoffPointer{Ref: ref, ID: id.String(), File: relative})
	workspace.SetValue(bundle.document, "handoffs", bundle.Handoffs)
	bundle.document.Record.Updated = stringPtr(now)
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path, id: handoffPath}
	return s.apply("handoff.record:"+bundle.ID+":"+id.String(), []*workspace.Document{bundle.document, doc}, paths, "handoff.record", bundle.Ref+"/"+ref, handoffPath)
}

func (s Service) Complete(missionRef, owner string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.complete(missionRef, owner)
}

func (s Service) complete(missionRef, owner string) (Result, error) {
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if !bundle.Legacy && bundle.Status == "completed" && owner == bundle.Owner {
		return Result{Operation: "mission.complete", Ref: bundle.Ref, Path: bundle.Path}, nil
	}
	if bundle.Legacy || bundle.Status != "active" || owner != bundle.Owner {
		return Result{}, domain.NewStateRefusal(domain.RefusalUnauthorized, "by", "completion requires the active compact Mission owner", bundle.Owner, owner, "ask the named owner to confirm completion", nil)
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	for _, objective := range bundle.Objectives {
		if objective.Status != "implemented" {
			return Result{}, invalid("objectives.status", "finish every Objective before Mission completion")
		}
	}
	if bundle.Review != "automatic" && len(bundle.Reviews) == 0 {
		return Result{}, invalid("reviews", "record the required review before Mission completion")
	}
	// Completion enforces resolves_gaps rather than executing it. A Mission that
	// declared it would close a Gap has not finished until the Gap is closed, and
	// the write belongs to `contract amend` — the amendment happens when the work
	// resolving the Gap lands, not as a side effect of a lifecycle transition.
	if err := s.assertDeclaredGapsClosed(bundle); err != nil {
		return Result{}, err
	}
	now := s.now()
	bundle.document.Record.Status = stringPtr("completed")
	bundle.document.Record.Updated = stringPtr(now)
	bundle.Status = "completed"
	bundle.Updated = now
	var docs []*workspace.Document
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	if bundle.Run != nil {
		bundle.Run.Status = "completed"
		workspace.SetValue(bundle.document, "run", bundle.Run)
	} else {
		pointers := make([]Run, len(bundle.Runs))
		copy(pointers, bundle.Runs)
		for i := range bundle.Runs {
			if bundle.Runs[i].File == "" {
				pointers[i].Status = "completed"
				bundle.Runs[i].Status = "completed"
				continue
			}
			absolute, _ := containedFile(filepath.Dir(bundle.entry.Absolute), bundle.Runs[i].File)
			doc, readErr := workspace.ReadFile(absolute)
			if readErr != nil {
				return Result{}, readErr
			}
			doc.Record.Status = stringPtr("completed")
			docs = append(docs, doc)
			paths[doc.Record.ID] = filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), bundle.Runs[i].File))
			pointers[i] = Run{Ref: bundle.Runs[i].Ref, ID: bundle.Runs[i].ID, File: bundle.Runs[i].File}
			bundle.Runs[i].Status = "completed"
		}
		workspace.SetValue(bundle.document, "runs", pointers)
	}
	reviewRef := ""
	reviewedCommit := ""
	if len(bundle.Reviews) > 0 {
		reviewRef = bundle.Reviews[len(bundle.Reviews)-1].Ref
		if resolved := bundle.Reviews[len(bundle.Reviews)-1].Document; resolved != nil {
			reviewedCommit = resolved.Reviewed.Commit
		}
	}
	record := CompletionRecord{By: owner, At: now, Authorization: "owner supplied --by after schema checks", ReviewedCommit: reviewedCommit, Review: reviewRef}
	workspace.SetValue(bundle.document, "completion_record", record)
	bundle.CompletionRecord = &record
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	docs = append(docs, bundle.document)
	return s.apply("mission.complete:"+bundle.ID+":"+now, docs, paths, "mission.complete", bundle.Ref, bundle.Path)
}

func (s Service) apply(key string, docs []*workspace.Document, paths map[domain.ID]string, operation, ref, primaryPath string) (Result, error) {
	for _, doc := range docs {
		known := false
		for _, entry := range s.Workspace.Entries {
			if entry.Document.Record.ID != doc.Record.ID {
				continue
			}
			known = true
			current, readErr := workspace.ReadFile(entry.Absolute)
			if readErr != nil {
				return Result{}, readErr
			}
			fingerprint, fingerprintErr := workspace.Fingerprint(current)
			if fingerprintErr != nil {
				return Result{}, fingerprintErr
			}
			if fingerprint != entry.Fingerprint {
				return Result{}, domain.NewStateRefusal(domain.RefusalStaleFingerprint, entry.Path, "canonical source changed after command validation", entry.Fingerprint, fingerprint, "reload the Mission and retry the typed command", nil)
			}
			break
		}
		if !known {
			target := filepath.Clean(filepath.FromSlash(paths[doc.Record.ID]))
			if filepath.IsAbs(target) || target == ".." || strings.HasPrefix(target, ".."+string(filepath.Separator)) {
				return Result{}, domain.NewStateRefusal(domain.RefusalPathEscape, "path", "new record target escapes the workspace", "canonical workspace-relative path", paths[doc.Record.ID], "use the schema-derived path inside .spectacular", nil)
			}
			if _, statErr := os.Lstat(filepath.Join(s.Workspace.Root, target)); statErr == nil {
				return Result{}, domain.NewStateRefusal(domain.RefusalCollision, filepath.ToSlash(target), "new record target already exists", "unused schema-derived path", "occupied", "inspect the existing file and choose a distinct title or resolve the collision", nil)
			} else if !os.IsNotExist(statErr) {
				return Result{}, invalidCause(filepath.ToSlash(target), "inspect new record target", statErr)
			}
		}
	}
	changes := make([]governance.FileChange, 0, len(docs)+4)
	for _, doc := range docs {
		data, err := workspace.Canonical(doc)
		if err != nil {
			return Result{}, err
		}
		changes = append(changes, governance.FileChange{Path: paths[doc.Record.ID], Data: data, Mode: 0o644})
	}
	indexes, err := humanlayout.Indexes(s.Workspace.Entries, docs, paths)
	if err != nil {
		return Result{}, err
	}
	indexPaths := make([]string, 0, len(indexes))
	for path := range indexes {
		indexPaths = append(indexPaths, path)
	}
	sort.Strings(indexPaths)
	for _, path := range indexPaths {
		changes = append(changes, governance.FileChange{Path: path, Data: indexes[path], Mode: 0o644})
	}
	apply := s.ApplyTransaction
	if apply == nil {
		apply = governance.ApplyTransaction
	}
	if err := apply(s.Workspace.Root, key, changes); err != nil {
		return Result{}, err
	}
	changed := make([]string, 0, len(changes))
	for _, change := range changes {
		changed = append(changed, filepath.ToSlash(change.Path))
	}
	return Result{Operation: operation, Ref: ref, Path: primaryPath, Changed: changed}, nil
}

func (s Service) beginMutation() (Service, func(), error) {
	if s.Workspace == nil {
		return Service{}, nil, invalid("workspace", "workspace is required")
	}
	unlock, err := acquireMutationLock(s.Workspace.Root)
	if err != nil {
		return Service{}, nil, err
	}
	fresh, err := discovery.Open(s.Workspace.Root)
	if err != nil {
		unlock()
		return Service{}, nil, err
	}
	s.Workspace = fresh
	return s, unlock, nil
}

func acquireMutationLock(root string) (func(), error) {
	directory := filepath.Join(root, ".spectacular", "transactions")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, invalidCause("transactions", "create mutation lock directory", err)
	}
	file, err := os.OpenFile(filepath.Join(directory, ".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, invalidCause("transactions", "open mutation lock", err)
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		return nil, domain.NewStateRefusal(domain.RefusalCollision, "transactions", "another Mission mutation is in progress", "exclusive mutation lock", "busy", "wait for the active command to finish, reload the Mission, and retry", err)
	}
	return func() {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	}, nil
}

func validatePlan(plan Plan) error {
	if plan.Title == "" || plan.Owner == "" || plan.Outcome == "" || plan.Contract.Ref == "" || len(plan.Completion) == 0 || len(plan.Objectives) == 0 || len(plan.Stops) == 0 {
		return invalid("input", "plan requires title, owner, Contract ref, outcome, completion, Objectives, and stops")
	}
	if plan.Review != "automatic" && plan.Review != "clustered" && plan.Review != "independent" {
		return invalid("review", "must be automatic, clustered, or independent")
	}
	for _, objective := range plan.Objectives {
		if objective.Outcome == "" || len(objective.Claims) == 0 {
			return invalid("objectives", "each plan Objective needs outcome and claims")
		}
	}
	return nil
}

func resolveContract(ws *discovery.Workspace, ref string) (Binding, error) {
	typed, err := domain.ParseReference(ref)
	if err != nil || typed.Type != domain.Contract {
		return Binding{}, invalidCause("contract.ref", "must be Contract:<UUIDv7>", err)
	}
	entry, err := ws.Lookup(ref, domain.Contract)
	if err != nil {
		return Binding{}, err
	}
	data, err := os.ReadFile(entry.Absolute)
	if err != nil {
		return Binding{}, err
	}
	digest := sha256.Sum256(data)
	return Binding{Ref: ref, Fingerprint: "sha256:" + hex.EncodeToString(digest[:])}, nil
}

func nextMissionRef(ws *discovery.Workspace) string {
	max := 0
	for _, entry := range ws.OfType(domain.Mission) {
		ref := workspace.RefOrEmpty(entry.Document)
		if missionRefPattern.MatchString(ref) {
			value, _ := strconv.Atoi(strings.TrimPrefix(ref, "M"))
			if value > max {
				max = value
			}
		}
	}
	return "M" + strconv.Itoa(max+1)
}

func gitBaseline(root string) (string, string, string, error) {
	commitCommand := exec.Command("git", "rev-parse", "HEAD")
	commitCommand.Dir = root
	commit, err := commitCommand.Output()
	if err != nil {
		return "", "", "", invalidCause("baseline.commit", "read Git HEAD", err)
	}
	branchCommand := exec.Command("git", "branch", "--show-current")
	branchCommand.Dir = root
	branch, err := branchCommand.Output()
	if err != nil || strings.TrimSpace(string(branch)) == "" {
		return "", "", "", invalidCause("baseline.branch", "read current Git branch", err)
	}
	timeCommand := exec.Command("git", "show", "-s", "--format=%cI", strings.TrimSpace(string(commit)))
	timeCommand.Dir = root
	baselineAt, err := timeCommand.Output()
	if err != nil {
		return "", "", "", invalidCause("baseline.commit", "read Git commit time", err)
	}
	return strings.TrimSpace(string(commit)), strings.TrimSpace(string(branch)), strings.TrimSpace(string(baselineAt)), nil
}

func stableID(at, key string) (domain.ID, error) {
	instant, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return "", invalidCause("id", "derive retry-stable UUIDv7 timestamp", err)
	}
	digest := sha256.Sum256([]byte(key))
	var raw [16]byte
	milliseconds := uint64(instant.UnixMilli())
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], milliseconds)
	copy(raw[:6], encoded[2:])
	copy(raw[6:], digest[:10])
	raw[6] = (raw[6] & 0x0f) | 0x70
	raw[8] = (raw[8] & 0x3f) | 0x80
	return domain.ParseID(uuid.UUID(raw).String())
}

// missionRecordPath resolves a Mission-scoped record's canonical location through
// the layout system rather than through a join written at the call site, so a
// record created by any path lands where the layout rule says it lands. It
// returns the workspace-relative path and the bundle-relative pointer the
// Mission stores.
func (s Service) missionRecordPath(bundle *Bundle, doc *workspace.Document, ref string) (string, string, error) {
	// The layout system reads the scoped ref from the field workspace.Ref reads,
	// so the Mission scope is written there for the resolution and the record's
	// own ref is restored afterward. A record stores its leaf ref and names its
	// Mission in mission:; the scoped spelling is how the layout rule is asked
	// where that record belongs, not something the file carries.
	unscoped := workspace.RefOrEmpty(doc)
	workspace.SetString(doc, workspace.RefField, bundle.Ref+"/"+ref)
	path, err := humanlayout.PlannedPath(s.Workspace.Entries, doc)
	workspace.SetString(doc, workspace.RefField, unscoped)
	if err != nil {
		return "", "", err
	}
	path = filepath.ToSlash(path)
	bundleDirectory := filepath.ToSlash(filepath.Dir(bundle.Path))
	relative, err := filepath.Rel(bundleDirectory, path)
	if err != nil {
		return "", "", invalidCause("record", "resolve record path relative to its Mission", err)
	}
	return path, filepath.ToSlash(relative), nil
}

func resolveCommitTree(root, commit string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", commit+"^{tree}")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// verifyReviewedGit checks a record's git binding against the real repository.
// The field prefix names the record being verified, so a Handoff's refusal does
// not tell its sender to correct a review.
func verifyReviewedGit(root, commit, tree string, field string) error {
	if !commitPattern.MatchString(commit) || !commitPattern.MatchString(tree) {
		return invalid(field+".reviewed", "a bound record must carry canonical 40-character Git commit and tree IDs")
	}
	commitCommand := exec.Command("git", "rev-parse", "--verify", commit+"^{commit}")
	commitCommand.Dir = root
	resolvedCommit, err := commitCommand.Output()
	if err != nil || strings.TrimSpace(string(resolvedCommit)) != commit {
		return domain.NewStateRefusal(domain.RefusalInvalidKnownField, field+".reviewed.commit", "reviewed commit does not exist in this repository", "existing exact commit", commit, "review the committed tree, then record its exact commit and tree", err)
	}
	treeCommand := exec.Command("git", "rev-parse", "--verify", commit+"^{tree}")
	treeCommand.Dir = root
	resolvedTree, err := treeCommand.Output()
	actual := strings.TrimSpace(string(resolvedTree))
	if err != nil || actual != tree {
		return domain.NewStateRefusal(domain.RefusalStaleFingerprint, field+".reviewed.tree", "reviewed tree does not belong to the reviewed commit", actual, tree, "record the exact tree from git rev-parse <commit>^{tree}", err)
	}
	return nil
}

func runDocument(bundle *Bundle, run Run, title string) (*workspace.Document, string, error) {
	id, err := domain.ParseID(run.ID)
	if err != nil {
		return nil, "", err
	}
	run.Title = title
	doc := &workspace.Document{Record: domain.Record{Type: domain.Run, ID: id, Title: stringPtr(title), Status: stringPtr(run.Status)}, Unknown: map[string]*yaml.Node{}, Body: "# Run\n\n" + title + "\n"}
	workspace.SetString(doc, "ref", run.Ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetString(doc, "operator", run.Operator)
	workspace.SetString(doc, "started_at", run.StartedAt)
	workspace.SetString(doc, "current_objective", run.CurrentObjective)
	slug := humanlayout.Slug(title)
	path := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), "runs", run.Ref+"-"+slug, run.Ref+"-"+slug+".md"))
	return doc, path, nil
}

func runTitle(run Run) string {
	if run.Title != "" {
		return run.Title
	}
	if run.Ref == "R1" {
		return "Initial run"
	}
	return "Run " + run.Ref
}

func nextObjective(objectives []Objective) string {
	for _, objective := range objectives {
		if objective.Status != "implemented" {
			return objective.Ref
		}
	}
	return ""
}

func (s Service) now() string {
	now := s.Now
	if now == nil {
		now = time.Now
	}
	return now().UTC().Truncate(time.Second).Format(time.RFC3339)
}

func presentStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func readInput(path string, stdin []byte) ([]byte, error) {
	if path == "-" {
		if len(stdin) == 0 {
			return nil, invalid("input", "stdin is empty")
		}
		return stdin, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, domain.NewRefusal(domain.RefusalRecordNotFound, "input", "read Markdown input", err)
	}
	return data, nil
}

func splitInput(data []byte) ([]byte, string, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return nil, "", invalid("input", "Markdown input must start with YAML frontmatter")
	}
	rest := strings.TrimPrefix(text, "---\n")
	index := strings.Index(rest, "\n---\n")
	if index < 0 {
		return nil, "", invalid("input", "Markdown input needs a closing --- delimiter")
	}
	return []byte(rest[:index] + "\n"), strings.TrimLeft(rest[index+5:], "\n"), nil
}

func stringPtr(value string) *string { return &value }
