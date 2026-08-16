package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
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
	"go.yaml.in/yaml/v3"
)

type Service struct {
	Workspace *discovery.Workspace
	Now       func() time.Time
}

type Plan struct {
	Type         string      `yaml:"type"`
	Title        string      `yaml:"title"`
	Owner        string      `yaml:"owner"`
	Contract     Binding     `yaml:"contract"`
	Outcome      string      `yaml:"outcome"`
	Review       string      `yaml:"review"`
	Completion   []Criterion `yaml:"completion"`
	Objectives   []Objective `yaml:"objectives"`
	Authority    Authority   `yaml:"authority"`
	Scope        Scope       `yaml:"scope"`
	RepairBudget int         `yaml:"repair_budget"`
	Dependencies []string    `yaml:"dependencies"`
	Gaps         []string    `yaml:"gaps"`
	Stops        []string    `yaml:"stops"`
	Body         string      `yaml:"-"`
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
	contract, err := resolveContract(s.Workspace, plan.Contract.Ref)
	if err != nil {
		return Result{}, err
	}
	missionID, err := domain.NewID()
	if err != nil {
		return Result{}, err
	}
	missionRef := nextMissionRef(s.Workspace)
	now := s.now()
	objectives := make([]Objective, len(plan.Objectives))
	for i, source := range plan.Objectives {
		id, idErr := domain.NewID()
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
	runID, err := domain.NewID()
	if err != nil {
		return Result{}, err
	}
	commit, branch, err := gitBaseline(s.Workspace.Root)
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
	workspace.SetString(doc, "start_key", startKey)
	temporary := &Bundle{Outcome: plan.Outcome, Review: plan.Review, Completion: plan.Completion, Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget, Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops}
	fingerprint, err := FrozenFingerprint(temporary)
	if err != nil {
		return Result{}, err
	}
	workspace.SetValue(doc, "activation", Activation{By: plan.Owner, At: now, Fingerprint: fingerprint})
	path := filepath.ToSlash(filepath.Join(".spectacular", "missions", missionRef+"-"+humanlayout.Slug(plan.Title), "MISSION.md"))
	candidate := &Bundle{
		ID: missionID.String(), Ref: missionRef, Title: plan.Title, Status: "active", Owner: plan.Owner,
		Contract: contract, Baseline: &Baseline{Commit: commit, Branch: branch}, Outcome: plan.Outcome,
		Review: plan.Review, Completion: plan.Completion, Objectives: objectives, Run: &run,
		Activation: &Activation{By: plan.Owner, At: now, Fingerprint: fingerprint}, Validation: Validation{Schema: Schema, Mode: "cli"},
		Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget,
		Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops, Path: path,
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
	return Load(s.Workspace, ref)
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

func (s Service) StartRun(missionRef, title string) (Result, error) {
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy || bundle.Status != "active" || strings.TrimSpace(title) == "" {
		return Result{}, invalid("run", "new Run requires an active compact Mission and title")
	}
	if _, err := Validate(s.Workspace, bundle); err != nil {
		return Result{}, err
	}
	runs := allRuns(bundle)
	if len(runs) == 0 {
		return Result{}, invalid("run", "Mission has no current Run")
	}
	for _, run := range runs {
		if run.Status == "active" || run.Status == "awaiting-review" {
			// Starting a new execution boundary closes the preceding one.
		}
	}
	docs := []*workspace.Document{}
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	pointers := make([]Run, 0, len(runs)+1)
	for _, run := range runs {
		if run.Status != "completed" {
			run.Status = "completed"
		}
		if run.File == "" {
			doc, path, makeErr := runDocument(bundle, run, runTitle(run))
			if makeErr != nil {
				return Result{}, makeErr
			}
			docs = append(docs, doc)
			paths[doc.Record.ID] = path
			run = Run{Ref: run.Ref, ID: run.ID, File: strings.TrimPrefix(path, filepath.ToSlash(filepath.Dir(bundle.Path))+"/")}
		} else {
			absolute, pathErr := containedFile(filepath.Dir(bundle.entry.Absolute), run.File)
			if pathErr != nil {
				return Result{}, pathErr
			}
			doc, readErr := workspace.ReadFile(absolute)
			if readErr != nil {
				return Result{}, readErr
			}
			doc.Record.Status = stringPtr("completed")
			docs = append(docs, doc)
			paths[doc.Record.ID] = filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), run.File))
			run = Run{Ref: run.Ref, ID: run.ID, File: run.File}
		}
		pointers = append(pointers, run)
	}
	nextID, err := domain.NewID()
	if err != nil {
		return Result{}, err
	}
	current := nextObjective(bundle.Objectives)
	if current == "" {
		current = bundle.Objectives[len(bundle.Objectives)-1].Ref
	}
	next := Run{Ref: "R" + strconv.Itoa(len(runs)+1), ID: nextID.String(), Title: title, Status: "active", Operator: bundle.Owner, StartedAt: s.now(), CurrentObjective: current}
	doc, nextPath, err := runDocument(bundle, next, title)
	if err != nil {
		return Result{}, err
	}
	docs = append(docs, doc)
	paths[doc.Record.ID] = nextPath
	pointers = append(pointers, Run{Ref: next.Ref, ID: next.ID, File: strings.TrimPrefix(nextPath, filepath.ToSlash(filepath.Dir(bundle.Path))+"/")})
	workspace.Delete(bundle.document, "run")
	workspace.SetValue(bundle.document, "runs", pointers)
	bundle.document.Record.Updated = stringPtr(s.now())
	docs = append(docs, bundle.document)
	return s.apply("run.start:"+bundle.ID+":"+next.ID, docs, paths, "run.start", bundle.Ref+"/"+next.Ref, nextPath)
}

func (s Service) RecordReview(missionRef, path string, stdin []byte) (Result, error) {
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
	if draft.Reviewed.ActivationFingerprint != bundle.Activation.Fingerprint || !commitPattern.MatchString(draft.Reviewed.Commit) || !commitPattern.MatchString(draft.Reviewed.Tree) {
		return Result{}, invalid("review.reviewed", "review must bind exact commit, tree, and Mission activation fingerprint")
	}
	id, err := domain.NewID()
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
	relative := filepath.ToSlash(filepath.Join("reviews", ref+"-"+humanlayout.Slug(draft.Title)+".md"))
	bundle.Reviews = append(bundle.Reviews, ReviewPointer{Ref: ref, ID: id.String(), File: relative, Verdict: "pass"})
	workspace.SetValue(bundle.document, "reviews", bundle.Reviews)
	bundle.document.Record.Updated = stringPtr(now)
	reviewPath := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), relative))
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path, id: reviewPath}
	return s.apply("review.record:"+bundle.ID+":"+id.String(), []*workspace.Document{bundle.document, doc}, paths, "review.record", bundle.Ref+"/"+ref, reviewPath)
}

func (s Service) Complete(missionRef, owner string) (Result, error) {
	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return Result{}, err
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
	now := s.now()
	bundle.document.Record.Status = stringPtr("completed")
	bundle.document.Record.Updated = stringPtr(now)
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
		}
		workspace.SetValue(bundle.document, "runs", pointers)
	}
	reviewRef := ""
	reviewedCommit := ""
	if len(bundle.Reviews) > 0 {
		reviewRef = bundle.Reviews[len(bundle.Reviews)-1].Ref
		path, _ := containedFile(filepath.Dir(bundle.entry.Absolute), bundle.Reviews[len(bundle.Reviews)-1].File)
		if doc, readErr := workspace.ReadFile(path); readErr == nil {
			var reviewed struct {
				Commit string `yaml:"commit"`
			}
			if workspace.DecodeValue(doc, "reviewed", &reviewed) == nil {
				reviewedCommit = reviewed.Commit
			}
		}
	}
	record := CompletionRecord{By: owner, At: now, Authorization: "owner supplied --by after schema checks", ReviewedCommit: reviewedCommit, Review: reviewRef}
	workspace.SetValue(bundle.document, "completion_record", record)
	docs = append(docs, bundle.document)
	return s.apply("mission.complete:"+bundle.ID+":"+now, docs, paths, "mission.complete", bundle.Ref, bundle.Path)
}

func (s Service) apply(key string, docs []*workspace.Document, paths map[domain.ID]string, operation, ref, primaryPath string) (Result, error) {
	unlock, err := acquireMutationLock(s.Workspace.Root)
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	for _, doc := range docs {
		for _, entry := range s.Workspace.Entries {
			if entry.Document.Record.ID != doc.Record.ID {
				continue
			}
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
	if err := governance.ApplyTransaction(s.Workspace.Root, key, changes); err != nil {
		return Result{}, err
	}
	changed := make([]string, 0, len(changes))
	for _, change := range changes {
		changed = append(changed, filepath.ToSlash(change.Path))
	}
	return Result{Operation: operation, Ref: ref, Path: primaryPath, Changed: changed}, nil
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
		ref, _ := workspace.String(entry.Document, "ref", false)
		if ref == "" {
			ref, _ = workspace.String(entry.Document, "human_ref", false)
		}
		if missionRefPattern.MatchString(ref) {
			value, _ := strconv.Atoi(strings.TrimPrefix(ref, "M"))
			if value > max {
				max = value
			}
		}
	}
	return "M" + strconv.Itoa(max+1)
}

func gitBaseline(root string) (string, string, error) {
	commitCommand := exec.Command("git", "rev-parse", "HEAD")
	commitCommand.Dir = root
	commit, err := commitCommand.Output()
	if err != nil {
		return "", "", invalidCause("baseline.commit", "read Git HEAD", err)
	}
	branchCommand := exec.Command("git", "branch", "--show-current")
	branchCommand.Dir = root
	branch, err := branchCommand.Output()
	if err != nil || strings.TrimSpace(string(branch)) == "" {
		return "", "", invalidCause("baseline.branch", "read current Git branch", err)
	}
	return strings.TrimSpace(string(commit)), strings.TrimSpace(string(branch)), nil
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
	workspace.SetInt(doc, "repairs", run.Repairs)
	path := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), "runs", run.Ref+"-"+humanlayout.Slug(title)+".md"))
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
