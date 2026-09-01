package missionbundle

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	AllowMain    bool          `yaml:"allow_main,omitempty"`
	CreateBranch bool          `yaml:"create_branch,omitempty"`
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
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var plan Plan
		if err := json.Unmarshal(trimmed, &plan); err != nil {
			return Plan{}, nil, invalidCause("input", "decode Mission plan JSON", err)
		}
		if plan.Type == "" {
			plan.Type = "MissionPlan"
		}
		if plan.Type != "MissionPlan" {
			return Plan{}, nil, invalid("type", "Mission start input must declare type: MissionPlan")
		}
		return plan, data, nil
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
	missionRef := nextMissionRef(s.Workspace)
	if (branch == "main" || branch == "master") && !plan.AllowMain {
		if plan.CreateBranch {
			featBranch := "feat/" + missionRef + "-" + humanlayout.Slug(plan.Title)
			cmd := exec.Command("git", "checkout", "-b", featBranch)
			cmd.Dir = s.Workspace.Root
			if out, err := cmd.CombinedOutput(); err != nil {
				return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "baseline.branch", "failed to create feature branch: "+string(out), "clean feature branch", branch, "check git status and create branch manually", err)
			}
			branch = featBranch
		} else {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "baseline.branch", "activating directly on "+branch+" is prohibited without a dedicated feature branch", "feature branch (e.g. feat/"+missionRef+"-"+humanlayout.Slug(plan.Title)+")", branch, "checkout a feature branch ('git checkout -b <branch>') or supply --create-branch / --allow-main", nil)
		}
	}
	missionID, err := stableID(baselineAt, startKey+":mission")
	if err != nil {
		return Result{}, err
	}
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

func (s Service) CheckWithVerify(ref string) (Check, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Check{}, err
	}
	check, err := Validate(s.Workspace, bundle)
	if err != nil {
		return check, err
	}

	root := ""
	if s.Workspace != nil {
		root = s.Workspace.Root
	}

	// 1. Run domain verification if script exists
	testCmd := ""
	checkScriptPath := filepath.Join(root, "tests", "check.sh")
	if stat, statErr := os.Stat(checkScriptPath); statErr == nil && !stat.IsDir() {
		testCmd = "sh tests/check.sh"
	}

	if testCmd != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", "-c", testCmd)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			check.Valid = false
			check.Notices = append(check.Notices, fmt.Sprintf("domain verification failed: %s", strings.TrimSpace(string(out))))
			return check, nil
		}
		check.Checks = append(check.Checks, "domain-verification-pass")
	}

	// 2. Replay check if declared
	if bundle.Replay != nil && bundle.Replay.Command != "" {
		for _, p := range bundle.Replay.CachePaths {
			full := filepath.Join(root, p)
			_ = os.RemoveAll(full)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		replayCmd := exec.CommandContext(ctx, "sh", "-c", bundle.Replay.Command)
		replayCmd.Dir = root
		if out, err := replayCmd.CombinedOutput(); err != nil {
			check.Valid = false
			check.Notices = append(check.Notices, fmt.Sprintf("replay command failed: %s", strings.TrimSpace(string(out))))
			return check, nil
		}

		if testCmd != "" {
			ctxPost, cancelPost := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancelPost()
			cmd := exec.CommandContext(ctxPost, "sh", "-c", testCmd)
			cmd.Dir = root
			if out, err := cmd.CombinedOutput(); err != nil {
				check.Valid = false
				check.Notices = append(check.Notices, fmt.Sprintf("post-replay verification failed: %s", strings.TrimSpace(string(out))))
				return check, nil
			}
		}
		check.Checks = append(check.Checks, "replay-reconstruction-pass")
	}

	// 3. Git working tree cleanliness check
	ctxGit, cancelGit := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelGit()
	gitStatus := exec.CommandContext(ctxGit, "git", "status", "--porcelain")
	gitStatus.Dir = root
	if out, err := gitStatus.Output(); err == nil {
		if len(bytes.TrimSpace(out)) == 0 {
			check.Checks = append(check.Checks, "git-working-tree-clean")
		} else {
			check.Notices = append(check.Notices, "git working tree contains untracked or uncommitted changes")
		}
	}

	return check, nil
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

	docs := []*workspace.Document{}
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}

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
						if !isDependencyUnblockedByDecision(s.Workspace, bundle.Ref, depRef, r.Ref, targetObjective) {
							return Result{}, domain.NewRefusal(domain.RefusalCollision, "dependency", fmt.Sprintf("cannot start run on %s: upstream dependency %s is in state %s; resolve blocker with owner Decision first", targetObjective, depRef, r.Status), nil)
						}
					}
				}
			}
		}
		if currentObj.File == "" {
			objID, _ := domain.ParseID(currentObj.ID)
			objDoc := &workspace.Document{
				Record:  domain.Record{Type: domain.Objective, ID: objID, Title: stringPtr(currentObj.Outcome), Status: stringPtr(currentObj.Status)},
				Unknown: map[string]*yaml.Node{},
				Body:    "# Objective\n\n" + currentObj.Outcome + "\n",
			}
			workspace.SetString(objDoc, "ref", currentObj.Ref)
			workspace.SetString(objDoc, "mission", bundle.Ref)
			workspace.SetString(objDoc, "outcome", currentObj.Outcome)
			workspace.SetStrings(objDoc, "after", currentObj.After)
			workspace.SetStrings(objDoc, "claims", currentObj.Claims)
			objRelative := filepath.ToSlash(filepath.Join("objectives", currentObj.Ref+"-"+humanlayout.Slug(currentObj.Outcome)+".md"))
			currentObj.File = objRelative
			for i := range bundle.Objectives {
				if bundle.Objectives[i].Ref == currentObj.Ref {
					bundle.Objectives[i].File = objRelative
				}
			}
			workspace.SetValue(bundle.document, "objectives", bundle.Objectives)
			objPath := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), objRelative))
			docs = append(docs, objDoc)
			paths[objID] = objPath
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

func (s Service) RecordReviewDraft(missionRef string, draft ReviewDraft) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	bundle, err := Load(locked.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "review recording requires an active compact Mission")
	}
	if _, err := Validate(locked.Workspace, bundle); err != nil {
		return Result{}, err
	}
	return locked.recordReviewDraft(bundle, draft, draft.Body)
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
	trimmed := bytes.TrimSpace(data)
	var draft ReviewDraft
	var body string
	if len(trimmed) > 0 && trimmed[0] == '{' {
		if err := json.Unmarshal(trimmed, &draft); err != nil {
			return Result{}, invalidCause("input", "decode ReviewDraft JSON", err)
		}
		if draft.Type == "" {
			draft.Type = "ReviewDraft"
		}
	} else {
		frontmatter, b, err := splitInput(data)
		if err != nil {
			return Result{}, err
		}
		body = b
		if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
			return Result{}, invalidCause("input", "decode ReviewDraft frontmatter", err)
		}
	}
	return s.recordReviewDraft(bundle, draft, body)
}

func (s Service) recordReviewDraft(bundle *Bundle, draft ReviewDraft, body string) (Result, error) {
	if draft.Type == "" {
		draft.Type = "ReviewDraft"
	}
	if draft.Type != "ReviewDraft" || draft.Title == "" || draft.Status != "passed" || len(draft.Claims) == 0 {
		return Result{}, invalid("review", "requires type ReviewDraft, title, status passed, and claim verdicts")
	}
	if draft.Reviewed.ActivationFingerprint == "" && bundle.Activation != nil {
		draft.Reviewed.ActivationFingerprint = bundle.Activation.Fingerprint
	}
	if draft.Reviewed.Commit == "" || draft.Reviewed.Commit == "HEAD" {
		commit, tree, err := currentGitCommitAndTree(s.Workspace.Root)
		if err != nil {
			return Result{}, err
		}
		draft.Reviewed.Commit = commit
		if draft.Reviewed.Tree == "" {
			draft.Reviewed.Tree = tree
		}
	}
	if draft.Reviewer.Actor == "" {
		draft.Reviewer.Actor = s.Workspace.Config.Defaults.Operator
	}
	if draft.Reviewer.Actor == "" {
		draft.Reviewer.Actor = "Alex"
	}
	if draft.Reviewer.Operator == "" && bundle.Run != nil {
		draft.Reviewer.Operator = bundle.Run.Operator
	}
	if draft.Reviewer.Operator == "" {
		draft.Reviewer.Operator = "Alex"
	}
	if draft.Reviewer.RelationToOperator == "" {
		if draft.Reviewer.Actor == draft.Reviewer.Operator {
			draft.Reviewer.RelationToOperator = "same-actor"
			draft.Reviewer.ImplementedReviewedScope = true
		} else {
			draft.Reviewer.RelationToOperator = "independent"
			draft.Reviewer.ImplementedReviewedScope = false
		}
	}
	if draft.Reviewer.IndependenceBasis == "" {
		if draft.Reviewer.RelationToOperator == "independent" {
			draft.Reviewer.IndependenceBasis = "distinct operator verification"
		} else {
			draft.Reviewer.IndependenceBasis = "operator self-verification"
		}
	}
	if len(draft.Reviewer.Evidence) == 0 {
		for _, e := range bundle.Evidence {
			draft.Reviewer.Evidence = append(draft.Reviewer.Evidence, e.Ref)
		}
		if len(draft.Reviewer.Evidence) == 0 {
			draft.Reviewer.Evidence = []string{"git-commit-proof"}
		}
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

func (s Service) RecordHandoffDraft(missionRef string, draft HandoffDraft, sender string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	if err := CheckPassiveGitState(locked.Workspace.Root); err != nil {
		return Result{}, err
	}
	bundle, err := Load(locked.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "handoff recording requires an active compact Mission")
	}
	if _, err := Validate(locked.Workspace, bundle); err != nil {
		return Result{}, err
	}
	return locked.recordHandoffDraft(bundle, draft, sender, "")
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
	trimmed := bytes.TrimSpace(data)
	var draft HandoffDraft
	var body string
	if len(trimmed) > 0 && trimmed[0] == '{' {
		if err := json.Unmarshal(trimmed, &draft); err != nil {
			return Result{}, invalidCause("input", "decode HandoffDraft JSON", err)
		}
		if draft.Type == "" {
			draft.Type = "HandoffDraft"
		}
	} else {
		frontmatter, b, err := splitInput(data)
		if err != nil {
			return Result{}, err
		}
		body = b
		if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
			return Result{}, invalidCause("input", "decode HandoffDraft frontmatter", err)
		}
		if draft.Asserted == nil {
			return Result{}, invalid("handoff.asserted", "a Handoff must state asserted, even as an empty list")
		}
		if draft.Assumed == nil {
			return Result{}, invalid("handoff.assumed", "a Handoff must state assumed, even as an empty list")
		}
	}
	return s.recordHandoffDraft(bundle, draft, sender, body)
}

func (s Service) recordHandoffDraft(bundle *Bundle, draft HandoffDraft, sender string, body string) (Result, error) {
	if draft.Type == "" {
		draft.Type = "HandoffDraft"
	}
	if draft.Type != "HandoffDraft" || draft.Title == "" {
		return Result{}, invalid("handoff", "requires type HandoffDraft and a title")
	}
	if draft.Asserted == nil {
		empty := []string{}
		draft.Asserted = &empty
	}
	if draft.Assumed == nil {
		empty := []string{}
		draft.Assumed = &empty
	}
	if len(draft.Stops) == 0 {
		draft.Stops = []string{"scope-drift"}
	}
	if len(draft.Returns) == 0 {
		draft.Returns = []string{"test-receipt"}
	}
	if sender == "" && draft.Sender.Actor != "" {
		sender = draft.Sender.Actor
	}
	if sender == "" {
		sender = s.Workspace.Config.Defaults.Operator
	}
	if sender == "" {
		sender = "Alex"
	}
	if draft.Sender.Actor == "" {
		draft.Sender.Actor = sender
	}
	if draft.Sender.RelationToReceiver == "" {
		draft.Sender.RelationToReceiver = "delegation"
	}
	if draft.Sender.Actor != sender {
		return Result{}, domain.NewStateRefusal(domain.RefusalUnauthorized, "by",
			"handoff sender must match the identity recording it", draft.Sender.Actor, sender,
			"record the Handoff as its stated sender, or correct the draft", nil)
	}
	if draft.Reviewed.Commit == "" || draft.Reviewed.Commit == "HEAD" {
		commit, tree, err := currentGitCommitAndTree(s.Workspace.Root)
		if err != nil {
			return Result{}, err
		}
		draft.Reviewed.Commit = commit
		if draft.Reviewed.Tree == "" {
			draft.Reviewed.Tree = tree
		}
	}
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
	if draft.RuntimePointer != nil {
		switch draft.RuntimePointer.WorkspaceMode {
		case "", "share", "branch", "inherit":
		default:
			return Result{}, invalid("handoff.runtime_pointer.workspace_mode", fmt.Sprintf("workspace_mode must be one of: share, branch, inherit; got %q", draft.RuntimePointer.WorkspaceMode))
		}
		if draft.RuntimePointer.Harness != "" || draft.RuntimePointer.ThreadID != "" || draft.RuntimePointer.WorkspaceMode != "" {
			workspace.SetValue(doc, "runtime_pointer", draft.RuntimePointer)
		}
	}
	workspace.SetString(doc, "task", draft.Task)
	workspace.SetStrings(doc, "asserted", *draft.Asserted)
	workspace.SetStrings(doc, "assumed", *draft.Assumed)
	workspace.SetStrings(doc, "stops", draft.Stops)
	workspace.SetStrings(doc, "returns", draft.Returns)
	if err := ValidateWritePaths(draft.Writes); err != nil {
		return Result{}, err
	}
	activeRes, err := collectActiveWriteReservations(s.Workspace, bundle, draft.Supersedes)
	if err != nil {
		return Result{}, err
	}
	for _, res := range activeRes {
		for _, dw := range draft.Writes {
			if PathsOverlap(dw, res.path) {
				return Result{}, domain.NewRefusal(domain.RefusalInvalidScope, "writes", fmt.Sprintf("cannot reserve write path %q: overlaps with active Handoff %s/%s reservation %q", dw, res.missionRef, res.handoffRef, res.path), nil)
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

func (s Service) CreateContract(ref, title, owner string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.createContract(ref, title, owner)
}

func (s Service) createContract(ref, title, owner string) (Result, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Result{}, invalid("ref", "contract ref is required (e.g. CC-my-feature)")
	}
	if !strings.HasPrefix(ref, "CC-") {
		ref = "CC-" + ref
	}
	if title == "" {
		title = strings.TrimPrefix(ref, "CC-")
	}
	if owner == "" {
		owner = s.Workspace.Config.Defaults.Operator
		if owner == "" {
			owner = "Alex"
		}
	}
	now := s.now()
	contractID, err := domain.NewID()
	if err != nil {
		return Result{}, err
	}
	slug := humanlayout.Slug(title)
	fileName := ref + "-" + slug + ".md"
	relPath := filepath.ToSlash(filepath.Join(".spectacular", "contracts", fileName))

	doc := &workspace.Document{
		Record: domain.Record{
			Type:    domain.Contract,
			ID:      contractID,
			Title:   stringPtr(title),
			Status:  stringPtr("current"),
			Created: stringPtr(now),
			Updated: stringPtr(now),
		},
		Unknown: map[string]*yaml.Node{},
		Body:    "# " + title + "\n\n" + "## Overview\n\n" + title + " contract specification.\n",
	}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "owner", owner)
	workspace.SetString(doc, "contract_version", "1")
	workspace.SetString(doc, "purpose", title)
	workspace.SetString(doc, "outcome", title)
	workspace.SetStrings(doc, "applies_when", []string{"Work on " + title + " begins."})
	workspace.SetStrings(doc, "does_not_apply_when", []string{})
	workspace.SetStrings(doc, "does_not_provide", []string{})
	workspace.SetStrings(doc, "required_behavior", []string{})
	workspace.SetStrings(doc, "command_surface", []string{})
	workspace.SetStrings(doc, "mandatory_validation", []string{})

	paths := map[domain.ID]string{contractID: relPath}
	return s.apply("contract.create:"+contractID.String(), []*workspace.Document{doc}, paths, "contract.create", ref, relPath)
}

func (s Service) AmendScope(ref string, addPaths []string, owner, reason string, dryRun bool) (Result, error) {
	if !dryRun {
		locked, unlock, err := s.beginMutation()
		if err != nil {
			return Result{}, err
		}
		defer unlock()
		return locked.amendScope(ref, addPaths, owner, reason, false)
	}
	fresh, err := discovery.Open(s.Workspace.Root)
	if err != nil {
		return Result{}, err
	}
	s.Workspace = fresh
	return s.amendScope(ref, addPaths, owner, reason, true)
}

func (s Service) amendScope(ref string, addPaths []string, owner, reason string, dryRun bool) (Result, error) {
	if strings.TrimSpace(owner) == "" {
		return Result{}, invalid("by", "an amendment requires the owner who authorizes it")
	}
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy {
		return Result{}, invalid("mission", "legacy Mission is read-only")
	}
	if len(addPaths) == 0 {
		return Result{}, invalid("add", "at least one path must be specified to add to mechanical scope")
	}
	seen := map[string]bool{}
	for _, p := range bundle.Scope.Mechanical {
		seen[strings.TrimSuffix(p, "/")] = true
	}
	for _, p := range addPaths {
		clean := strings.TrimSuffix(p, "/")
		if clean == "" || filepath.IsAbs(clean) || filepath.ToSlash(filepath.Clean(clean)) != clean || strings.HasPrefix(clean, "../") || strings.Contains(clean, "\\") {
			return Result{}, invalid("scope.mechanical", "paths must be canonical workspace-relative paths: "+p)
		}
		if !seen[clean] {
			bundle.Scope.Mechanical = append(bundle.Scope.Mechanical, p)
			seen[clean] = true
		}
	}
	newFP, err := FrozenFingerprint(bundle)
	if err != nil {
		return Result{}, err
	}
	if dryRun {
		return Result{Operation: "mission.amend_scope", Ref: bundle.Ref, Path: bundle.Path, Fingerprint: newFP}, nil
	}
	if bundle.Activation != nil {
		bundle.Activation.Fingerprint = newFP
		workspace.SetValue(bundle.document, "activation", bundle.Activation)
	}
	workspace.SetValue(bundle.document, "scope", bundle.Scope)
	if strings.TrimSpace(reason) != "" {
		workspace.SetString(bundle.document, "scope_amendment_reason", reason)
	}
	bundle.document.Record.Updated = stringPtr(s.now())

	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	return s.apply("mission.amend_scope:"+bundle.ID+":"+newFP, []*workspace.Document{bundle.document}, paths, "mission.amend_scope", bundle.Ref, bundle.Path)
}

func (s Service) CloseMission(missionRef, owner string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.closeMission(missionRef, owner)
}

func (s Service) closeMission(missionRef, owner string) (Result, error) {
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if !bundle.Legacy && bundle.Status == "completed" && owner == bundle.Owner {
		return Result{Operation: "mission.complete", Ref: bundle.Ref, Path: bundle.Path}, nil
	}
	if bundle.Legacy || bundle.Status != "active" || owner != bundle.Owner {
		return Result{}, domain.NewStateRefusal(domain.RefusalUnauthorized, "by", "close requires the active compact Mission owner", bundle.Owner, owner, "ask the named owner to confirm closeout", nil)
	}
	now := s.now()
	var docs []*workspace.Document
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}

	for i := range bundle.Objectives {
		if bundle.Objectives[i].Status != "implemented" {
			bundle.Objectives[i].Status = "implemented"
			if bundle.Objectives[i].File != "" {
				objDoc, objErr := s.Workspace.Lookup(bundle.Objectives[i].ID, domain.Objective)
				if objErr == nil {
					objDoc.Document.Record.Status = stringPtr("implemented")
					objDoc.Document.Record.Updated = stringPtr(now)
					docs = append(docs, objDoc.Document)
					paths[objDoc.Document.Record.ID] = objDoc.Path
				}
			}
		}
	}
	workspace.SetValue(bundle.document, "objectives", bundle.Objectives)

	if bundle.Review != "automatic" && len(bundle.Reviews) == 0 {
		return Result{}, invalid("reviews", "record the required review before Mission completion")
	}
	if err := s.assertDeclaredGapsClosed(bundle); err != nil {
		return Result{}, err
	}

	bundle.document.Record.Status = stringPtr("completed")
	bundle.document.Record.Updated = stringPtr(now)
	bundle.Status = "completed"
	bundle.Updated = now

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
	return s.apply("mission.close:"+bundle.ID+":"+now, docs, paths, "mission.complete", bundle.Ref, bundle.Path)
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
		existingData, readErr := os.ReadFile(filepath.Join(s.Workspace.Root, filepath.FromSlash(path)))
		if readErr == nil && bytes.Equal(existingData, indexes[path]) {
			continue
		}
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

// acquireMutationLock takes the workspace-wide exclusive lock that makes a
// concurrent mutation refuse instead of interleaving. The lock itself is
// load-bearing; only the kernel call that takes it is platform-specific, so
// lockFile and unlockFile are supplied by service_unix.go and
// service_windows.go.
func acquireMutationLock(root string) (func(), error) {
	directory := filepath.Join(root, ".spectacular", "transactions")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, invalidCause("transactions", "create mutation lock directory", err)
	}
	file, err := os.OpenFile(filepath.Join(directory, ".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, invalidCause("transactions", "open mutation lock", err)
	}
	if err := lockFile(file); err != nil {
		file.Close()
		return nil, domain.NewStateRefusal(domain.RefusalCollision, "transactions", "another Mission mutation is in progress", "exclusive mutation lock", "busy", "wait for the active command to finish, reload the Mission, and retry", err)
	}
	return func() {
		_ = unlockFile(file)
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
	entry, err := ws.Lookup(ref, domain.Contract)
	if err != nil {
		return Binding{}, invalidCause("contract.ref", "must be a valid Contract reference (e.g. CC-* or Contract:<UUIDv7>)", err)
	}
	data, err := os.ReadFile(entry.Absolute)
	if err != nil {
		return Binding{}, err
	}
	digest := sha256.Sum256(data)
	canonicalRef := string(domain.Contract) + ":" + entry.Document.Record.ID.String()
	return Binding{Ref: canonicalRef, Fingerprint: "sha256:" + hex.EncodeToString(digest[:])}, nil
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
	if path == "-" || (path == "" && len(stdin) > 0) {
		if len(stdin) == 0 {
			return nil, invalid("input", "stdin is empty")
		}
		return stdin, nil
	}
	if path == "" {
		return nil, invalid("input", "input path is empty")
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

func (s Service) RecordEvidence(missionRef, path string, stdin []byte, fromPath string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.recordEvidence(missionRef, path, stdin, fromPath)
}

func (s Service) RecordEvidenceDraft(missionRef string, draft EvidenceDraft) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.recordEvidenceDraft(missionRef, draft, "")
}

func (s Service) recordEvidenceDraft(missionRef string, draft EvidenceDraft, body string) (Result, error) {
	if err := CheckPassiveGitState(s.Workspace.Root); err != nil {
		return Result{}, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "evidence recording requires an active compact Mission")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	if draft.Type == "" {
		draft.Type = "EvidenceDraft"
	}
	if draft.Title == "" {
		draft.Title = "Verification evidence"
	}
	if draft.Actor == "" {
		draft.Actor = s.Workspace.Config.Defaults.Operator
	}
	if draft.Actor == "" {
		draft.Actor = "Alex"
	}
	if draft.Commit == "" || draft.Commit == "HEAD" {
		commit, tree, gitErr := currentGitCommitAndTree(s.Workspace.Root)
		if gitErr != nil {
			return Result{}, gitErr
		}
		draft.Commit = commit
		if draft.Tree == "" {
			draft.Tree = tree
		}
	}
	if draft.Tree == "" {
		resolvedTree, err := resolveCommitTree(s.Workspace.Root, draft.Commit)
		if err != nil {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "evidence.commit", "evidence commit does not exist in this repository", "existing exact commit", draft.Commit, "record exact committed tree", err)
		}
		draft.Tree = resolvedTree
	}
	if len(draft.Claims) == 0 {
		for _, c := range bundle.Completion {
			draft.Claims = append(draft.Claims, c.Claim)
		}
	}
	if len(draft.Checks) == 0 {
		draft.Checks = []EvidenceCheck{{Name: "verification", Result: "pass"}}
	}
	if err := verifyReviewedGit(s.Workspace.Root, draft.Commit, draft.Tree, "evidence"); err != nil {
		return Result{}, err
	}

	key := "evidence.record:" + bundle.ID + ":" + draft.Commit + ":" + draft.Title
	id, err := stableID(bundle.Activation.At, key)
	if err != nil {
		return Result{}, err
	}
	for _, existing := range bundle.Evidence {
		if existing.ID == id.String() {
			return Result{
				Operation: "evidence.record",
				Ref:       bundle.Ref + "/" + existing.Ref,
				Path:      filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), existing.File)),
			}, nil
		}
	}

	ref := "E" + strconv.Itoa(len(bundle.Evidence)+1)
	now := s.now()
	doc := &workspace.Document{
		Record:  domain.Record{Type: domain.Evidence, ID: id, Title: stringPtr(draft.Title), Created: stringPtr(now)},
		Unknown: map[string]*yaml.Node{}, Body: body,
	}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetString(doc, "actor", draft.Actor)
	workspace.SetString(doc, "commit", draft.Commit)
	workspace.SetString(doc, "tree", draft.Tree)
	if len(draft.Objectives) > 0 {
		workspace.SetStrings(doc, "objectives", draft.Objectives)
	}
	if len(draft.Runs) > 0 {
		workspace.SetStrings(doc, "runs", draft.Runs)
	}
	if len(draft.Claims) > 0 {
		workspace.SetStrings(doc, "claims", draft.Claims)
	}
	if len(draft.Checks) > 0 {
		workspace.SetValue(doc, "checks", draft.Checks)
	}
	if len(draft.Limitations) > 0 {
		workspace.SetStrings(doc, "limitations", draft.Limitations)
	}
	evidencePath, relative, err := s.missionRecordPath(bundle, doc, ref)
	if err != nil {
		return Result{}, err
	}
	bundle.Evidence = append(bundle.Evidence, EvidencePointer{Ref: ref, ID: id.String(), File: relative})
	workspace.SetValue(bundle.document, "evidence", bundle.Evidence)
	bundle.document.Record.Updated = stringPtr(now)
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path, id: evidencePath}
	return s.apply("evidence.record:"+bundle.ID+":"+id.String(), []*workspace.Document{bundle.document, doc}, paths, "evidence.record", bundle.Ref+"/"+ref, evidencePath)
}

func (s Service) recordEvidence(missionRef, path string, stdin []byte, fromPath string) (Result, error) {
	if err := CheckPassiveGitState(s.Workspace.Root); err != nil {
		return Result{}, err
	}
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" {
		return Result{}, invalid("mission", "evidence recording requires an active compact Mission")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}

	var draft EvidenceDraft
	body := ""

	if fromPath != "" && (path == "" || path == "-") && len(stdin) == 0 {
		fromData, readErr := os.ReadFile(fromPath)
		if readErr != nil {
			return Result{}, domain.NewRefusal(domain.RefusalInvalidWorkspacePath, fromPath, "read test output file", readErr)
		}
		commit, tree, gitErr := currentGitCommitAndTree(s.Workspace.Root)
		if gitErr != nil {
			return Result{}, gitErr
		}
		var claims []string
		for _, c := range bundle.Completion {
			claims = append(claims, c.Claim)
		}
		checks := parseChecksFromTestOutput(fromData)
		draft = EvidenceDraft{
			Type:   "EvidenceDraft",
			Title:  "Automated test verification from " + filepath.Base(fromPath),
			Actor:  bundle.Owner,
			Commit: commit,
			Tree:   tree,
			Claims: claims,
			Checks: checks,
		}
		body = "# Verification\n\nDerived from `" + filepath.Base(fromPath) + "`.\n"
	} else {
		data, readErr := readInput(path, stdin)
		if readErr != nil {
			return Result{}, readErr
		}
		trimmed := bytes.TrimSpace(data)
		if len(trimmed) > 0 && trimmed[0] == '{' {
			if err := json.Unmarshal(trimmed, &draft); err != nil {
				return Result{}, invalidCause("input", "decode EvidenceDraft JSON", err)
			}
			if draft.Type == "" {
				draft.Type = "EvidenceDraft"
			}
		} else {
			frontmatter, parsedBody, splitErr := splitInput(data)
			if splitErr != nil {
				return Result{}, splitErr
			}
			body = parsedBody
			if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
				return Result{}, invalidCause("input", "decode EvidenceDraft frontmatter", err)
			}
		}
		if fromPath != "" {
			fromData, readErr := os.ReadFile(fromPath)
			if readErr == nil {
				parsedChecks := parseChecksFromTestOutput(fromData)
				draft.Checks = append(draft.Checks, parsedChecks...)
			}
		}
	}
	if draft.Type != "EvidenceDraft" || draft.Title == "" || draft.Actor == "" {
		return Result{}, invalid("evidence", "requires type EvidenceDraft, title, and actor")
	}
	if !commitPattern.MatchString(draft.Commit) {
		return Result{}, invalid("evidence.commit", "evidence must bind an exact commit")
	}
	if draft.Tree == "" {
		resolvedTree, err := resolveCommitTree(s.Workspace.Root, draft.Commit)
		if err != nil {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "evidence.commit", "evidence commit does not exist in this repository", "existing exact commit", draft.Commit, "record exact committed tree", err)
		}
		draft.Tree = resolvedTree
	} else if !commitPattern.MatchString(draft.Tree) {
		return Result{}, invalid("evidence.tree", "evidence must bind an exact tree")
	}
	if err := verifyReviewedGit(s.Workspace.Root, draft.Commit, draft.Tree, "evidence"); err != nil {
		return Result{}, err
	}

	key := "evidence.record:" + bundle.ID + ":" + draft.Commit + ":" + draft.Title
	id, err := stableID(bundle.Activation.At, key)
	if err != nil {
		return Result{}, err
	}
	for _, existing := range bundle.Evidence {
		if existing.ID == id.String() {
			return Result{
				Operation: "evidence.record",
				Ref:       bundle.Ref + "/" + existing.Ref,
				Path:      filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), existing.File)),
			}, nil
		}
	}

	ref := "E" + strconv.Itoa(len(bundle.Evidence)+1)
	now := s.now()
	doc := &workspace.Document{
		Record:  domain.Record{Type: domain.Evidence, ID: id, Title: stringPtr(draft.Title), Created: stringPtr(now)},
		Unknown: map[string]*yaml.Node{}, Body: body,
	}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "mission", bundle.Ref)
	workspace.SetString(doc, "actor", draft.Actor)
	workspace.SetString(doc, "commit", draft.Commit)
	workspace.SetString(doc, "tree", draft.Tree)
	if len(draft.Objectives) > 0 {
		workspace.SetStrings(doc, "objectives", draft.Objectives)
	}
	if len(draft.Runs) > 0 {
		workspace.SetStrings(doc, "runs", draft.Runs)
	}
	if len(draft.Claims) > 0 {
		workspace.SetStrings(doc, "claims", draft.Claims)
	}
	if len(draft.Checks) > 0 {
		workspace.SetValue(doc, "checks", draft.Checks)
	}
	if len(draft.Limitations) > 0 {
		workspace.SetStrings(doc, "limitations", draft.Limitations)
	}
	evidencePath, relative, err := s.missionRecordPath(bundle, doc, ref)
	if err != nil {
		return Result{}, err
	}

	candidate, err := decodeEvidence(doc, evidencePath)
	if err != nil {
		return Result{}, err
	}
	if err := validateEvidenceContent(candidate, bundle); err != nil {
		return Result{}, err
	}

	bundle.Evidence = append(bundle.Evidence, EvidencePointer{Ref: ref, ID: id.String(), File: relative})
	workspace.SetValue(bundle.document, "evidence", bundle.Evidence)
	bundle.document.Record.Updated = stringPtr(now)
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path, id: evidencePath}
	return s.apply("evidence.record:"+bundle.ID+":"+id.String(), []*workspace.Document{bundle.document, doc}, paths, "evidence.record", bundle.Ref+"/"+ref, evidencePath)
}

type activeReservation struct {
	missionRef string
	handoffRef string
	path       string
}

func collectActiveWriteReservations(ws *discovery.Workspace, currentBundle *Bundle, draftSupersedes string) ([]activeReservation, error) {
	var reservations []activeReservation
	if ws == nil {
		return reservations, nil
	}
	for _, entry := range ws.OfType(domain.Mission) {
		b, err := Load(ws, entry.Path)
		if err != nil || b == nil || b.Status != "active" {
			continue
		}
		// Track superseded handoffs in this bundle
		superseded := make(map[string]bool)
		if b.Ref == currentBundle.Ref && draftSupersedes != "" {
			superseded[draftSupersedes] = true
		}
		for _, hPtr := range b.Handoffs {
			if hPtr.Document != nil && hPtr.Document.Supersedes != "" {
				superseded[hPtr.Document.Supersedes] = true
			}
		}
		for _, hPtr := range b.Handoffs {
			if superseded[hPtr.Ref] {
				continue
			}
			if hPtr.Document != nil {
				for _, w := range hPtr.Document.Writes {
					reservations = append(reservations, activeReservation{
						missionRef: b.Ref,
						handoffRef: hPtr.Ref,
						path:       w,
					})
				}
			}
		}
	}
	return reservations, nil
}

func isDependencyUnblockedByDecision(ws *discovery.Workspace, missionRef, depRef, depRunRef, targetObjRef string) bool {
	if ws == nil {
		return false
	}
	for _, entry := range ws.OfType(domain.Decision) {
		doc := entry.Document
		if doc == nil {
			continue
		}
		status := ""
		if doc.Record.Status != nil {
			status = *doc.Record.Status
		}
		disp, _ := workspace.String(doc, "disposition", false)
		if status != "accepted" && disp != "accepted" {
			continue
		}
		targets, _ := workspace.Strings(doc, "targets", false)
		unblocked, _ := workspace.Strings(doc, "unblocked", false)
		scope, _ := workspace.Strings(doc, "scope", false)
		allRefs := append(targets, unblocked...)
		allRefs = append(allRefs, scope...)
		fullDepObj := missionRef + "/" + depRef
		fullDepRun := missionRef + "/" + depRunRef
		fullTargetObj := missionRef + "/" + targetObjRef
		for _, ref := range allRefs {
			if ref == depRef || ref == depRunRef || ref == fullDepObj || ref == fullDepRun || ref == targetObjRef || ref == fullTargetObj {
				return true
			}
		}
	}
	return false
}

type MissionSummary struct {
	Ref    string `json:"ref"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Holder string `json:"holder"`
	Next   string `json:"next"`
	Path   string `json:"path"`
}

type MissionListResult struct {
	SchemaVersion string           `json:"schema_version"`
	Missions      []MissionSummary `json:"missions"`
}

func (s Service) ListMissions(statusFilter string, all bool) (MissionListResult, error) {
	if s.Workspace != nil {
		if fresh, err := discovery.Open(s.Workspace.Root); err == nil {
			s.Workspace = fresh
		}
	}
	var results []MissionSummary
	seen := map[string]bool{}
	for _, entry := range s.Workspace.Entries {
		if entry.Document.Record.Type != domain.Mission {
			continue
		}
		human, _, err := workspace.Ref(entry.Document)
		if err != nil || human == "" || seen[human] {
			continue
		}
		seen[human] = true
		bundle, err := Load(s.Workspace, human)
		if err != nil {
			continue
		}
		if statusFilter != "" && statusFilter != "all" {
			if bundle.Status != statusFilter {
				continue
			}
		} else if !all && statusFilter != "all" {
			// Active-first: exclude completed, resolved, superseded, cancelled
			switch bundle.Status {
			case "completed", "resolved", "superseded", "cancelled":
				continue
			}
		}
		derived := bundle.Derive()
		title := bundle.Title
		if title == "" && entry.Document.Record.Title != nil {
			title = *entry.Document.Record.Title
		}
		results = append(results, MissionSummary{
			Ref:    bundle.Ref,
			Title:  title,
			Status: bundle.Status,
			Holder: derived.Holder,
			Next:   derived.Next,
			Path:   bundle.Path,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Ref < results[j].Ref
	})
	return MissionListResult{
		SchemaVersion: "spectacular.mission.list.v2",
		Missions:      results,
	}, nil
}

func parseChecksFromTestOutput(data []byte) []EvidenceCheck {
	var checks []EvidenceCheck
	content := string(data)
	lines := strings.Split(content, "\n")
	hasJSON := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var evt struct {
			Action  string  `json:"Action"`
			Test    string  `json:"Test"`
			Elapsed float64 `json:"Elapsed"`
		}
		if err := json.Unmarshal([]byte(line), &evt); err == nil && evt.Action != "" && evt.Test != "" {
			hasJSON = true
			if evt.Action == "pass" || evt.Action == "fail" {
				checks = append(checks, EvidenceCheck{
					Name:   evt.Test,
					Result: evt.Action,
				})
			}
		}
	}
	if hasJSON && len(checks) > 0 {
		return checks
	}
	result := "pass"
	if strings.Contains(content, "FAIL") || strings.Contains(content, "failed") {
		result = "fail"
	}
	return []EvidenceCheck{
		{Name: "automated-tests", Result: result},
	}
}

func currentGitCommitAndTree(root string) (string, string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return "", "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, "git", "git rev-parse HEAD", err)
	}
	commit := strings.TrimSpace(string(out))
	treeCmd := exec.Command("git", "rev-parse", "HEAD^{tree}")
	treeCmd.Dir = root
	treeOut, err := treeCmd.Output()
	if err != nil {
		return "", "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, "git", "git rev-parse HEAD^{tree}", err)
	}
	tree := strings.TrimSpace(string(treeOut))
	return commit, tree, nil
}
