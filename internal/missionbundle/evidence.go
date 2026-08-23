package missionbundle

import (
	"path/filepath"
	"regexp"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

var evidencePattern = regexp.MustCompile(`^E[0-9]+(-[a-z0-9]+)?$`)

// EvidencePointer points from a Mission bundle to one of its recorded Evidence packages.
type EvidencePointer struct {
	Ref      string    `yaml:"ref" json:"ref"`
	ID       string    `yaml:"id" json:"id"`
	File     string    `yaml:"file" json:"file"`
	Document *Evidence `yaml:"-" json:"document,omitempty"`
}

// EvidenceCheck describes one verification check execution and its observed outcome.
type EvidenceCheck struct {
	Name   string `yaml:"name" json:"name"`
	Result string `yaml:"result" json:"result"`
	Output string `yaml:"output,omitempty" json:"output,omitempty"`
}

// Evidence is the typed view of an Evidence package document.
type Evidence struct {
	ID          string          `json:"id"`
	Ref         string          `json:"ref"`
	Title       string          `json:"title"`
	Mission     string          `json:"mission"`
	Actor       string          `json:"actor"`
	ObservedAt  string          `json:"observed_at,omitempty"`
	Commit      string          `json:"commit"`
	Tree        string          `json:"tree"`
	Objectives  []string        `json:"objectives,omitempty"`
	Runs        []string        `json:"runs,omitempty"`
	Claims      []string        `json:"claims,omitempty"`
	Checks      []EvidenceCheck `json:"checks,omitempty"`
	Limitations []string        `json:"limitations,omitempty"`
	Path        string          `json:"path"`
	Body        string          `json:"-"`

	document *workspace.Document
}

// EvidenceDraft is the frontmatter schema authored when submitting an Evidence package.
type EvidenceDraft struct {
	Type        string          `yaml:"type"`
	Title       string          `yaml:"title"`
	Actor       string          `yaml:"actor"`
	Commit      string          `yaml:"commit"`
	Tree        string          `yaml:"tree,omitempty"`
	Objectives  []string        `yaml:"objectives,omitempty"`
	Runs        []string        `yaml:"runs,omitempty"`
	Claims      []string        `yaml:"claims,omitempty"`
	Checks      []EvidenceCheck `yaml:"checks,omitempty"`
	Limitations []string        `yaml:"limitations,omitempty"`
}

// decodeEvidence decodes an Evidence document into typed representation.
func decodeEvidence(doc *workspace.Document, path string) (*Evidence, error) {
	if doc.Record.Type != domain.Evidence {
		return nil, invalid("evidence.file", "pointer must resolve to an Evidence record")
	}
	e := &Evidence{
		ID:         doc.Record.ID.String(),
		Ref:        workspace.RefOrEmpty(doc),
		Title:      value(doc.Record.Title),
		ObservedAt: value(doc.Record.Created),
		Path:       path,
		Body:       doc.Body,

		document: doc,
	}
	e.Mission, _ = workspace.String(doc, "mission", false)
	e.Actor, _ = workspace.String(doc, "actor", false)
	e.Commit, _ = workspace.String(doc, "commit", false)
	e.Tree, _ = workspace.String(doc, "tree", false)
	e.Objectives, _ = workspace.Strings(doc, "objectives", false)
	e.Runs, _ = workspace.Strings(doc, "runs", false)
	e.Claims, _ = workspace.Strings(doc, "claims", false)
	e.Limitations, _ = workspace.Strings(doc, "limitations", false)
	_ = workspace.DecodeValue(doc, "checks", &e.Checks)
	return e, nil
}

// validateEvidenceContent checks schema requirements of an Evidence record against its parent Mission.
func validateEvidenceContent(e *Evidence, b *Bundle) error {
	if e.Mission == "" {
		return invalid("evidence.mission", "an Evidence package must name the Mission it belongs to")
	}
	if e.Mission != b.Ref && e.Mission != b.ID {
		if typed, err := domain.ParseReference(e.Mission); err != nil || typed.ID.String() != b.ID {
			return invalid("evidence.mission", "evidence must name the Mission that carries it")
		}
	}
	if e.Title == "" {
		return invalid("evidence.title", "an Evidence package must have a title")
	}
	if e.Actor == "" {
		return invalid("evidence.actor", "an Evidence package must name the recording actor")
	}
	if !commitPattern.MatchString(e.Commit) || !commitPattern.MatchString(e.Tree) {
		return invalid("evidence.commit", "an Evidence package must bind an exact commit and tree")
	}
	if len(e.Claims) == 0 && len(e.Objectives) == 0 && len(e.Runs) == 0 {
		return invalid("evidence.coverage", "an Evidence package must cover at least one Claim, Objective, or Run")
	}
	return nil
}

// validateEvidence checks every Evidence package pointed to by a Mission bundle.
func validateEvidence(ws *discovery.Workspace, b *Bundle) error {
	if len(b.Evidence) == 0 {
		return nil
	}
	base := filepath.Dir(b.entry.Absolute)
	seenRefs, seenIDs := map[string]bool{}, map[string]bool{}
	for i := range b.Evidence {
		pointer := &b.Evidence[i]
		if _, err := domain.ParseID(pointer.ID); err != nil {
			return invalidCause("evidence.id", "must be canonical UUIDv7", err)
		}
		if !evidencePattern.MatchString(pointer.Ref) || seenRefs[pointer.Ref] || seenIDs[pointer.ID] {
			return invalid("evidence.ref", "evidence refs and identities must be unique E<number> values")
		}
		seenRefs[pointer.Ref], seenIDs[pointer.ID] = true, true

		path, err := containedFile(base, pointer.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		ev, err := decodeEvidence(doc, path)
		if err != nil {
			return err
		}
		if ev.ID != pointer.ID {
			return invalid("evidence.file", "evidence pointer must resolve to the same Evidence identity")
		}
		if err := validateEvidenceContent(ev, b); err != nil {
			return err
		}
		pointer.Document = ev
	}
	return nil
}
