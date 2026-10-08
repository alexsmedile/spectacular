package missionbundle

import (
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

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
