package missionbundle

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

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
