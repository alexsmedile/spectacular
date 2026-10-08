package missionbundle

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

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
