package missionbundle

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strconv"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

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
