// Package context compiles bounded, source-backed runtime context. Its output
// is a disposable projection and never owns Mission or Contract truth.
package context

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/guardrails"
	"github.com/alexsmedile/spectacular/v2/internal/projection"
)

const SchemaVersion = "spectacular.context.v1"

type Config struct {
	Scope    string `json:"scope"`
	Event    string `json:"event,omitempty"`
	Selector string `json:"selector,omitempty"`
}

type Source struct {
	Role        string `json:"role"`
	Noun        string `json:"noun"`
	Ref         string `json:"ref"`
	Path        string `json:"path"`
	Fingerprint string `json:"fingerprint"`
	Authority   string `json:"authority"`
}

type Next struct {
	Kind         string `json:"kind"`
	Operation    string `json:"operation,omitempty"`
	Target       string `json:"target,omitempty"`
	AuthorizedBy string `json:"authorized_by,omitempty"`
	Code         string `json:"code,omitempty"`
	Detail       string `json:"detail,omitempty"`
	Source       string `json:"source,omitempty"`
}

type Bundle struct {
	SchemaVersion       string               `json:"schema_version"`
	Scope               string               `json:"scope"`
	GeneratedAt         string               `json:"generated_at"`
	GenerationBasis     string               `json:"generation_basis"`
	UniversalInvariants []string             `json:"universal_invariants"`
	Authoritative       []Source             `json:"authoritative"`
	ProjectionSources   []Source             `json:"projection_sources"`
	Guidance            []guardrails.Section `json:"guidance"`
	Gaps                []Source             `json:"gaps"`
	Conflicts           []string             `json:"conflicts"`
	Omissions           []string             `json:"omissions"`
	Next                Next                 `json:"next"`
	LoadedRecords       int                  `json:"loaded_records"`
	AvailableRecords    int                  `json:"available_records"`
}

type Compiler struct {
	Workspace *discovery.Workspace
	Now       func() time.Time
}

func (c Compiler) Compile(config Config) (Bundle, error) {
	if c.Workspace == nil {
		return Bundle{}, domain.NewRefusal(domain.RefusalWorkspaceNotFound, "workspace", "context compiler requires an open workspace", nil)
	}
	if config.Scope == "" {
		return Bundle{}, domain.NewRefusal(domain.RefusalInvalidScope, "scope", "project or exact Mission reference is required", nil)
	}
	now := c.Now
	if now == nil {
		now = time.Now
	}
	bundle := Bundle{
		SchemaVersion: SchemaVersion,
		Scope:         config.Scope,
		GeneratedAt:   now().UTC().Format(time.RFC3339Nano),
		UniversalInvariants: []string{
			"Owner authority is required for Mission and current Contract disposition.",
			"Compiled context and generated views are non-authoritative projections.",
			"Runtime execution and Handoffs cannot broaden the current authority envelope.",
			"Missing, stale, or conflicting required sources stop consequential action.",
		},
		Authoritative:     []Source{},
		ProjectionSources: []Source{},
		Guidance:          []guardrails.Section{},
		Gaps:              []Source{},
		Conflicts:         []string{},
		Omissions:         []string{},
		AvailableRecords:  len(c.Workspace.Entries),
	}
	builder := projection.Builder{Workspace: c.Workspace, Now: now}
	if config.Scope == "project" {
		view, err := builder.Project()
		if err != nil {
			return Bundle{}, err
		}
		bundle.Authoritative = appendUnique(bundle.Authoritative, sourceFromPointer("project-anchor", view.Authoritative.Identity, "authoritative"))
		for _, pointer := range view.Authoritative.CurrentTruth {
			bundle.Authoritative = appendUnique(bundle.Authoritative, sourceFromPointer("current-truth", pointer, "authoritative"))
		}
		for _, pointer := range view.Projection.Missions {
			bundle.ProjectionSources = appendUnique(bundle.ProjectionSources, sourceFromPointer("mission-projection", pointer, "projection"))
		}
		for _, pointer := range view.Projection.Gaps {
			bundle.Gaps = appendUnique(bundle.Gaps, sourceFromPointer("gap", pointer, "authoritative"))
		}
		bundle.Conflicts = append(bundle.Conflicts, view.Projection.Conflicts...)
		bundle.Omissions = append(bundle.Omissions, view.Projection.Omissions...)
		bundle.Next = next(view.Projection.Continuation, view.Projection.OwnerGate)
	} else {
		if _, err := c.Workspace.Lookup(config.Scope, domain.Mission); err != nil {
			return Bundle{}, domain.NewRefusal(domain.RefusalInvalidScope, "scope", "context scope must be project or an exact Mission reference", err)
		}
		card, err := builder.Mission(config.Scope)
		if err != nil {
			return Bundle{}, err
		}
		bundle.Authoritative = appendUnique(bundle.Authoritative, Source{Role: "mission-anchor", Noun: "Mission", Ref: "Mission:" + card.ID, Path: card.Source.Path, Fingerprint: card.Source.Fingerprint, Authority: "authoritative"})
		for _, pointer := range card.Sources {
			bundle.Authoritative = appendUnique(bundle.Authoritative, sourceFromPointer("accepted-source", pointer, "authoritative"))
		}
		for _, pointer := range card.Pointers {
			bundle.Authoritative = appendUnique(bundle.Authoritative, sourceFromPointer("mission-context", pointer, "authoritative"))
		}
		for _, pointer := range card.Gaps {
			bundle.Gaps = appendUnique(bundle.Gaps, sourceFromPointer("gap", pointer, "authoritative"))
		}
		bundle.Conflicts = append(bundle.Conflicts, card.Conflicts...)
		bundle.Omissions = append(bundle.Omissions, card.Omissions...)
		bundle.Next = next(card.Continuation, card.OwnerGate)
	}
	if config.Event != "" {
		if c.Workspace.Manifest.Guardrails == "" {
			bundle.Omissions = append(bundle.Omissions, "owner Guardrails are supported but unused")
		} else {
			data, fingerprint, err := c.Workspace.ReadMetadataFile(c.Workspace.Manifest.Guardrails)
			if err != nil {
				return Bundle{}, err
			}
			document, err := guardrails.Parse(data)
			if err != nil {
				return Bundle{}, domain.NewRefusal(domain.RefusalInvalidManifest, "guardrails", err.Error(), err)
			}
			bundle.Guidance, err = document.Select(config.Event, config.Selector)
			if err != nil {
				return Bundle{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "event", err.Error(), err)
			}
			bundle.ProjectionSources = appendUnique(bundle.ProjectionSources, Source{Role: "owner-guardrails", Noun: "Guardrails", Ref: c.Workspace.Manifest.Guardrails, Path: ".spectacular/" + c.Workspace.Manifest.Guardrails, Fingerprint: fingerprint, Authority: "owner-guidance"})
		}
	} else if config.Selector != "" {
		return Bundle{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "selector", "selector requires an event", nil)
	}
	bundle.LoadedRecords = len(bundle.Authoritative) + len(bundle.Gaps)
	bundle.GenerationBasis = basis(bundle)
	return bundle, nil
}

func sourceFromPointer(role string, pointer projection.Pointer, authority string) Source {
	return Source{Role: role, Noun: pointer.Noun, Ref: pointer.Ref, Path: pointer.Path, Fingerprint: pointer.Fingerprint, Authority: authority}
}

func appendUnique(sources []Source, source Source) []Source {
	for _, existing := range sources {
		if existing.Path == source.Path && existing.Fingerprint == source.Fingerprint {
			return sources
		}
	}
	return append(sources, source)
}

func next(continuation *projection.Continuation, gate *projection.OwnerGate) Next {
	if gate != nil {
		return Next{Kind: "owner-gate", Code: gate.Code, Detail: gate.Detail, Source: gate.Source.Ref}
	}
	if continuation != nil {
		return Next{Kind: "continuation", Operation: continuation.Operation, Target: continuation.Target.Ref, AuthorizedBy: continuation.AuthorizedBy.Ref}
	}
	return Next{Kind: "owner-gate", Code: "no_safe_continuation", Detail: "canonical sources do not justify a continuation"}
}

func basis(bundle Bundle) string {
	var parts []string
	for _, source := range append(append([]Source{}, bundle.Authoritative...), append(bundle.Gaps, bundle.ProjectionSources...)...) {
		parts = append(parts, source.Path+":"+source.Fingerprint+":"+source.Authority)
	}
	for _, section := range bundle.Guidance {
		parts = append(parts, section.Event+":"+section.Selector+":"+section.Prose)
	}
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}
