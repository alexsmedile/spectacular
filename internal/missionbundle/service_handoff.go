package missionbundle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

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

type activeReservation struct {
	missionRef string
	handoffRef string
	path       string
}
