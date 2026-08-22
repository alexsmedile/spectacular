package missionbundle

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

var validRunTransitions = map[string][]string{
	"active":          {"paused", "blocked", "awaiting-review", "completed", "stopped"},
	"paused":          {"active", "blocked", "stopped"},
	"blocked":         {"active", "stopped"},
	"awaiting-review": {"active", "completed", "stopped"},
	"completed":       {}, // terminal
	"stopped":         {}, // terminal
}

type TransitionResult struct {
	Operation string   `json:"operation"`
	Ref       string   `json:"ref"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	By        string   `json:"by"`
	Reason    string   `json:"reason"`
	Path      string   `json:"path"`
	Changed   []string `json:"changed"`
}

func ValidateTransition(from, to string) error {
	allowed, ok := validRunTransitions[from]
	if !ok {
		return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("unknown origin run state %q", from), nil)
	}
	for _, a := range allowed {
		if a == to {
			return nil
		}
	}
	if len(allowed) == 0 {
		return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("state %q is terminal and cannot transition to %q", from, to), nil)
	}
	return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("illegal transition from %q to %q (allowed: %s)", from, to, strings.Join(allowed, ", ")), nil)
}

func (s Service) TransitionRun(targetRef, toState, actor, reason, nextAction string) (TransitionResult, error) {
	if strings.TrimSpace(actor) == "" {
		return TransitionResult{}, domain.NewRefusal(domain.RefusalMissingRequiredField, "by", "actor identity is required for run transition", nil)
	}
	if strings.TrimSpace(reason) == "" {
		return TransitionResult{}, domain.NewRefusal(domain.RefusalMissingRequiredField, "reason", "transition reason is required", nil)
	}
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return TransitionResult{}, err
	}
	defer unlock()
	return locked.transitionRun(targetRef, toState, actor, reason, nextAction)
}

func (s Service) transitionRun(targetRef, toState, actor, reason, nextAction string) (TransitionResult, error) {
	parts := strings.Split(targetRef, "/")
	missionRef := parts[0]
	runRef := ""
	if len(parts) == 2 {
		runRef = parts[1]
	} else if len(parts) == 3 {
		runRef = parts[2]
	} else {
		return TransitionResult{}, domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected <mission-ref>/<run-ref> or <mission-ref>/<objective-ref>/<run-ref>", nil)
	}

	bundle, err := Load(s.Workspace, missionRef)
	if err != nil {
		return TransitionResult{}, err
	}

	runs := allRuns(bundle)
	var targetRun *Run
	for i := range runs {
		if runs[i].Ref == runRef || runs[i].ID == runRef {
			targetRun = &runs[i]
			break
		}
	}
	if targetRun == nil {
		return TransitionResult{}, domain.NewRefusal(domain.RefusalRecordNotFound, runRef, fmt.Sprintf("run not found in mission %s", missionRef), nil)
	}

	if err := ValidateTransition(targetRun.Status, toState); err != nil {
		return TransitionResult{}, err
	}

	fromState := targetRun.Status
	now := s.now()
	hist := TransitionHistory{
		At:         now,
		From:       fromState,
		To:         toState,
		By:         actor,
		Reason:     reason,
		NextAction: nextAction,
	}

	targetRun.Status = toState
	targetRun.History = append(targetRun.History, hist)

	var docs []*workspace.Document
	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}

	if targetRun.File == "" {
		// Inline run: update mission bundle document
		workspace.SetValue(bundle.document, "run", targetRun)
		bundle.document.Record.Updated = stringPtr(now)
		docs = append(docs, bundle.document)
	} else {
		// Standalone run file: update run document and mission bundle document
		absolute, pathErr := containedFile(filepath.Dir(bundle.entry.Absolute), targetRun.File)
		if pathErr != nil {
			return TransitionResult{}, pathErr
		}
		doc, readErr := workspace.ReadFile(absolute)
		if readErr != nil {
			return TransitionResult{}, readErr
		}
		doc.Record.Status = stringPtr(toState)
		doc.Record.Updated = stringPtr(now)
		workspace.SetValue(doc, "history", targetRun.History)
		docs = append(docs, doc)
		runRelPath := filepath.ToSlash(filepath.Join(filepath.Dir(bundle.Path), targetRun.File))
		paths[doc.Record.ID] = runRelPath

		// Update pointer in bundle
		var updatedRuns []Run
		for _, r := range runs {
			if r.Ref == targetRun.Ref {
				r.Status = toState
			}
			updatedRuns = append(updatedRuns, Run{Ref: r.Ref, ID: r.ID, File: r.File, Status: r.Status})
		}
		workspace.SetValue(bundle.document, "runs", updatedRuns)
		bundle.document.Record.Updated = stringPtr(now)
		docs = append(docs, bundle.document)
	}

	res, err := s.apply("run.transition:"+bundle.ID+":"+targetRun.ID+":"+toState, docs, paths, "run.transition", bundle.Ref+"/"+targetRun.Ref, bundle.Path)
	if err != nil {
		return TransitionResult{}, err
	}

	return TransitionResult{
		Operation: "run.transition",
		Ref:       bundle.Ref + "/" + targetRun.Ref,
		From:      fromState,
		To:        toState,
		By:        actor,
		Reason:    reason,
		Path:      bundle.Path,
		Changed:   res.Changed,
	}, nil
}
