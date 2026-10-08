package missionbundle

import (
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

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
