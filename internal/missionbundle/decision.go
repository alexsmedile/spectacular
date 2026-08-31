package missionbundle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

type DecisionDraft struct {
	Type              string   `yaml:"type"`
	Title             string   `yaml:"title"`
	Actor             string   `yaml:"actor"`
	ActorRole         string   `yaml:"actor_role"`
	Question          string   `yaml:"question"`
	Disposition       string   `yaml:"disposition"`
	Rationale         string   `yaml:"rationale"`
	Alternatives      []string `yaml:"alternatives,omitempty"`
	AuthorityBasis    string   `yaml:"authority_basis,omitempty"`
	AuthorizedEffects []string `yaml:"authorized_effects,omitempty"`
	Conditions        []string `yaml:"conditions,omitempty"`
	Scope             []string `yaml:"scope,omitempty"`
	Targets           []string `yaml:"targets,omitempty"`
	Supersedes        string   `yaml:"supersedes,omitempty"`
	Body              string   `yaml:"-"`
}

type DecisionResult struct {
	Operation string   `json:"operation"`
	Ref       string   `json:"ref"`
	ID        string   `json:"id"`
	Path      string   `json:"path"`
	Unblocked []string `json:"unblocked,omitempty"`
	Changed   []string `json:"changed"`
}

func ValidateDecisionDraft(draft *DecisionDraft) error {
	if draft.Type == "" {
		draft.Type = "DecisionDraft"
	}
	if draft.Type != "DecisionDraft" && draft.Type != "Decision" {
		return invalid("type", "Decision input must declare type: DecisionDraft or Decision")
	}
	if strings.TrimSpace(draft.Title) == "" {
		return invalid("title", "Decision title is required")
	}
	if strings.TrimSpace(draft.Disposition) == "" {
		return invalid("disposition", "Decision disposition is required")
	}
	validDispositions := map[string]bool{
		"accepted":   true,
		"rejected":   true,
		"deferred":   true,
		"superseded": true,
	}
	if !validDispositions[draft.Disposition] {
		return invalid("disposition", fmt.Sprintf("invalid disposition %q; must be accepted, rejected, deferred, or superseded", draft.Disposition))
	}
	if strings.TrimSpace(draft.Rationale) == "" {
		return invalid("rationale", "Decision rationale is required")
	}
	return nil
}

func ReadDecisionDraft(path string, stdin []byte) (DecisionDraft, string, error) {
	data, err := readInput(path, stdin)
	if err != nil {
		return DecisionDraft{}, "", err
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var draft DecisionDraft
		if err := json.Unmarshal(trimmed, &draft); err != nil {
			return DecisionDraft{}, "", invalidCause("input", "decode Decision draft JSON", err)
		}
		if err := ValidateDecisionDraft(&draft); err != nil {
			return DecisionDraft{}, "", err
		}
		return draft, draft.Body, nil
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return DecisionDraft{}, "", err
	}
	var draft DecisionDraft
	if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
		return DecisionDraft{}, "", invalidCause("input", "decode Decision draft frontmatter", err)
	}
	if err := ValidateDecisionDraft(&draft); err != nil {
		return DecisionDraft{}, "", err
	}
	draft.Body = body
	return draft, string(data), nil
}

func (s Service) RecordDecision(path string, stdin []byte) (DecisionResult, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return DecisionResult{}, err
	}
	defer unlock()
	return locked.recordDecision(path, stdin)
}

func (s Service) RecordDecisionDraft(draft DecisionDraft) (DecisionResult, error) {
	if err := ValidateDecisionDraft(&draft); err != nil {
		return DecisionResult{}, err
	}
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return DecisionResult{}, err
	}
	defer unlock()
	return locked.recordDecisionDraft(draft)
}

var refNumberRegex = regexp.MustCompile(`^D([0-9]+)`)

func (s Service) recordDecision(path string, stdin []byte) (DecisionResult, error) {
	draft, _, err := ReadDecisionDraft(path, stdin)
	if err != nil {
		return DecisionResult{}, err
	}
	return s.recordDecisionDraft(draft)
}

func (s Service) recordDecisionDraft(draft DecisionDraft) (DecisionResult, error) {
	// 1. Validate supersession target if declared
	if draft.Supersedes != "" {
		found := false
		for _, entry := range s.Workspace.Entries {
			base := strings.TrimSuffix(filepath.Base(entry.Path), ".md")
			if base == draft.Supersedes || strings.HasPrefix(base, draft.Supersedes+"-") {
				found = true
				break
			}
			if docRef, _, _ := workspace.Ref(entry.Document); docRef == draft.Supersedes {
				found = true
				break
			}
		}
		if !found {
			return DecisionResult{}, domain.NewRefusal(domain.RefusalRecordNotFound, "supersedes", fmt.Sprintf("superseded decision %q not found in workspace", draft.Supersedes), nil)
		}
	}

	// 2. Retry convergence: if exact identical decision already exists, return existing receipt
	for _, entry := range s.Workspace.Entries {
		if entry.Document != nil && entry.Document.Record.Type == domain.Decision {
			docTitle := ""
			if entry.Document.Record.Title != nil {
				docTitle = *entry.Document.Record.Title
			}
			docDisp, _ := workspace.String(entry.Document, "disposition", false)
			docRat, _ := workspace.String(entry.Document, "rationale", false)
			if docTitle == draft.Title && docDisp == draft.Disposition && strings.TrimSpace(docRat) == strings.TrimSpace(draft.Rationale) {
				ref, _, _ := workspace.Ref(entry.Document)
				return DecisionResult{
					Operation: "decision.record",
					Ref:       ref,
					ID:        entry.Document.Record.ID.String(),
					Path:      entry.Path,
					Changed:   nil,
				}, nil
			}
		}
	}

	maxN := 0
	for _, entry := range s.Workspace.Entries {
		base := filepath.Base(entry.Path)
		if m := refNumberRegex.FindStringSubmatch(base); len(m) > 1 {
			if n, parseErr := strconv.Atoi(m[1]); parseErr == nil && n > maxN {
				maxN = n
			}
		}
	}
	nextN := maxN + 1

	slug := slugifyTitle(draft.Title)
	if slug == "" {
		slug = slugify(draft.Disposition)
	}
	if slug == "" {
		slug = "decision"
	}
	ref := fmt.Sprintf("D%d-%s", nextN, slug)

	domainID, err := domain.NewID()
	if err != nil {
		return DecisionResult{}, fmt.Errorf("generate decision id: %w", err)
	}

	now := s.now()
	body := draft.Body
	if strings.TrimSpace(body) == "" {
		body = "# " + draft.Title + "\n\n## Rationale\n" + draft.Rationale + "\n"
	}

	actor := draft.Actor
	if actor == "" {
		actor = "Alex"
	}

	doc := &workspace.Document{
		Record: domain.Record{
			Type:      domain.Decision,
			ID:        domainID,
			Title:     stringPtr(draft.Title),
			CreatedBy: stringPtr(actor),
			Created:   stringPtr(now),
			Updated:   stringPtr(now),
		},
		Unknown: map[string]*yaml.Node{},
		Body:    body,
	}

	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "actor", actor)
	workspace.SetString(doc, "actor_role", draft.ActorRole)
	workspace.SetString(doc, "question", draft.Question)
	workspace.SetString(doc, "disposition", draft.Disposition)
	workspace.SetString(doc, "rationale", draft.Rationale)
	if len(draft.Alternatives) > 0 {
		workspace.SetStrings(doc, "alternatives", draft.Alternatives)
	}
	if draft.AuthorityBasis != "" {
		workspace.SetString(doc, "authority_basis", draft.AuthorityBasis)
	}
	if len(draft.AuthorizedEffects) > 0 {
		workspace.SetStrings(doc, "authorized_effects", draft.AuthorizedEffects)
	}
	if len(draft.Conditions) > 0 {
		workspace.SetStrings(doc, "conditions", draft.Conditions)
	}
	if len(draft.Scope) > 0 {
		workspace.SetStrings(doc, "scope", draft.Scope)
	}
	if len(draft.Targets) > 0 {
		workspace.SetStrings(doc, "targets", draft.Targets)
	}
	if draft.Supersedes != "" {
		workspace.SetString(doc, "supersedes", draft.Supersedes)
	}

	decisionRelPath := fmt.Sprintf(".spectacular/decisions/%s.md", ref)
	paths := map[domain.ID]string{domainID: decisionRelPath}

	res, err := s.apply("decision.record:"+domainID.String(), []*workspace.Document{doc}, paths, "decision.record", ref, decisionRelPath)
	if err != nil {
		return DecisionResult{}, err
	}

	// Check for unblocked Objectives
	var unblocked []string
	for _, entry := range s.Workspace.Entries {
		if entry.Document.Record.Type == domain.Mission {
			bundle, decodeErr := decode(s.Workspace, entry)
			if decodeErr == nil && bundle.Status == "active" {
				for _, obj := range bundle.Objectives {
					for _, dep := range obj.After {
						if dep == ref || dep == fmt.Sprintf("D%d", nextN) {
							unblocked = append(unblocked, fmt.Sprintf("%s/%s", bundle.Ref, obj.Ref))
						}
					}
				}
			}
		}
	}

	return DecisionResult{
		Operation: "decision.record",
		Ref:       ref,
		ID:        domainID.String(),
		Path:      decisionRelPath,
		Unblocked: unblocked,
		Changed:   res.Changed,
	}, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func slugifyTitle(title string) string {
	title = strings.ToLower(title)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	words := strings.Fields(reg.ReplaceAllString(title, " "))
	var meaningful []string
	for _, w := range words {
		if len(w) > 2 && w != "the" && w != "and" && w != "for" && w != "with" && w != "from" && w != "that" {
			meaningful = append(meaningful, w)
		}
		if len(meaningful) >= 5 {
			break
		}
	}
	if len(meaningful) == 0 {
		meaningful = words
	}
	res := strings.Join(meaningful, "-")
	if len(res) > 40 {
		res = res[:40]
		res = strings.Trim(res, "-")
	}
	return res
}
