package missionbundle

import (
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

func ReadDecisionDraft(path string, stdin []byte) (DecisionDraft, string, error) {
	data, err := readInput(path, stdin)
	if err != nil {
		return DecisionDraft{}, "", err
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return DecisionDraft{}, "", err
	}
	var draft DecisionDraft
	if err := yaml.Unmarshal(frontmatter, &draft); err != nil {
		return DecisionDraft{}, "", invalidCause("input", "decode Decision draft frontmatter", err)
	}
	if draft.Type != "DecisionDraft" && draft.Type != "Decision" {
		return DecisionDraft{}, "", invalid("type", "Decision input must declare type: DecisionDraft or Decision")
	}
	if strings.TrimSpace(draft.Title) == "" {
		return DecisionDraft{}, "", invalid("title", "Decision title is required")
	}
	if strings.TrimSpace(draft.Disposition) == "" {
		return DecisionDraft{}, "", invalid("disposition", "Decision disposition is required")
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

var refNumberRegex = regexp.MustCompile(`^D([0-9]+)`)

func (s Service) recordDecision(path string, stdin []byte) (DecisionResult, error) {
	draft, _, err := ReadDecisionDraft(path, stdin)
	if err != nil {
		return DecisionResult{}, err
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

	slug := slugify(draft.Disposition)
	if slug == "" {
		slug = slugify(draft.Title)
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
