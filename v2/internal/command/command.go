// Package command owns the complete public Scenario A command registry.
package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/projection"
)

type Effect string

const ReadOnly Effect = "read-only"

type Spec struct {
	Words      []string
	Arguments  string
	JSONSchema string
	Effect     Effect
	Operation  Operation
}
type Operation uint8

const (
	opAnchorShowProject Operation = iota + 1
	opMissionList
	opMissionShow
	opGapList
	opGapShow
	opRunShow
	opCheckpointShow
	opEvidenceShow
	opDecisionShow
	opWorkspaceValidate
)

var Registry = []Spec{
	{[]string{"anchor", "show", "project"}, "[--json]", "spectacular.anchor.show.v1", ReadOnly, opAnchorShowProject},
	{[]string{"mission", "list"}, "[--json]", "spectacular.mission.list.v1", ReadOnly, opMissionList},
	{[]string{"mission", "show"}, "<ref> [--json]", "spectacular.mission.show.v1", ReadOnly, opMissionShow},
	{[]string{"gap", "list"}, "--scope <ref> [--json]", "spectacular.gap.list.v1", ReadOnly, opGapList},
	{[]string{"gap", "show"}, "<ref> [--json]", "spectacular.gap.show.v1", ReadOnly, opGapShow},
	{[]string{"run", "show"}, "<ref> [--json]", "spectacular.run.show.v1", ReadOnly, opRunShow},
	{[]string{"checkpoint", "show"}, "<ref> [--json]", "spectacular.checkpoint.show.v1", ReadOnly, opCheckpointShow},
	{[]string{"evidence", "show"}, "<ref> [--json]", "spectacular.evidence.show.v1", ReadOnly, opEvidenceShow},
	{[]string{"decision", "show"}, "<ref> [--json]", "spectacular.decision.show.v1", ReadOnly, opDecisionShow},
	{[]string{"workspace", "validate"}, "<scope> [--json]", "spectacular.workspace.validate.v1", ReadOnly, opWorkspaceValidate},
}

type Runner struct {
	Cwd    string
	Stdout io.Writer
	Stderr io.Writer
	Now    func() time.Time
}

func (r Runner) Run(args []string) int {
	original := append([]string(nil), args...)
	invoked := "spectacular " + strings.Join(original, " ")
	jsonMode, duplicateJSON := removeJSON(&args)
	usage := func(detail string) int { return r.usage(jsonMode, invoked, detail) }
	refuse := func(err error) int { return r.refuse(jsonMode, invoked, err) }
	if duplicateJSON {
		return usage("--json may be supplied at most once")
	}
	spec, rest, ok := match(args)
	if !ok {
		return usage("unknown or incomplete command")
	}
	workspace, err := discovery.Open(r.Cwd)
	if err != nil {
		return refuse(err)
	}
	b := projection.Builder{Workspace: workspace, Now: r.Now}
	var value any
	switch spec.Operation {
	case opAnchorShowProject:
		if len(rest) != 0 {
			return usage("anchor show project takes no arguments")
		}
		value, err = b.Project()
	case opMissionList:
		if len(rest) != 0 {
			return usage("mission list takes no arguments")
		}
		value, err = b.MissionList()
	case opMissionShow:
		if len(rest) != 1 {
			return usage("mission show requires exactly one ref")
		}
		value, err = b.Mission(rest[0])
	case opGapList:
		if len(rest) != 2 || rest[0] != "--scope" {
			return usage("gap list requires --scope <ref>")
		}
		value, err = b.Gaps(rest[1])
	case opGapShow:
		if len(rest) != 1 {
			return usage("gap show requires exactly one ref")
		}
		value, err = r.detail(b, rest, domain.Gap)
	case opRunShow:
		if len(rest) != 1 {
			return usage("run show requires exactly one ref")
		}
		value, err = r.detail(b, rest, domain.Run)
	case opCheckpointShow:
		if len(rest) != 1 {
			return usage("checkpoint show requires exactly one ref")
		}
		value, err = r.detail(b, rest, domain.Checkpoint)
	case opEvidenceShow:
		if len(rest) != 1 {
			return usage("evidence show requires exactly one ref")
		}
		value, err = r.detail(b, rest, domain.Evidence)
	case opDecisionShow:
		if len(rest) != 1 {
			return usage("decision show requires exactly one ref")
		}
		value, err = r.detail(b, rest, domain.Decision)
	case opWorkspaceValidate:
		if len(rest) != 1 {
			return usage("workspace validate requires exactly one scope")
		}
		value, err = b.Validate(rest[0])
	}
	if err != nil {
		return refuse(err)
	}
	envelope := b.Envelope(spec.JSONSchema, value)
	if jsonMode {
		if err := writeJSON(r.Stdout, envelope); err != nil {
			return r.refuse(true, invoked, err)
		}
	} else {
		renderHuman(r.Stdout, envelope)
	}
	return 0
}

func (r Runner) detail(b projection.Builder, args []string, noun domain.RecordType) (any, error) {
	if len(args) != 1 {
		return nil, domain.NewRefusal(domain.RefusalInvalidReference, "ref", "show requires exactly one ref", nil)
	}
	return b.Detail(args[0], noun)
}

func removeJSON(args *[]string) (bool, bool) {
	found := false
	duplicate := false
	out := (*args)[:0]
	for _, arg := range *args {
		if arg == "--json" {
			if found {
				duplicate = true
			}
			found = true
			continue
		}
		out = append(out, arg)
	}
	*args = out
	return found, duplicate
}
func match(args []string) (Spec, []string, bool) {
	for _, spec := range Registry {
		if len(args) < len(spec.Words) {
			continue
		}
		ok := true
		for i, w := range spec.Words {
			if args[i] != w {
				ok = false
				break
			}
		}
		if ok {
			return spec, args[len(spec.Words):], true
		}
	}
	return Spec{}, nil, false
}

type refusalEnvelope struct {
	SchemaVersion  string `json:"schema_version"`
	InvokedCommand string `json:"invoked_command"`
	ExitStatus     int    `json:"exit_status"`
	Code           string `json:"code"`
	Field          string `json:"field,omitempty"`
	Detail         string `json:"detail"`
}

func (r Runner) refuse(jsonMode bool, invoked string, err error) int {
	code := "internal_error"
	field := ""
	detail := err.Error()
	var refusal *domain.Refusal
	if errors.As(err, &refusal) {
		code = string(refusal.Code)
		field = refusal.Field
		detail = refusal.Detail
	}
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v1", invoked, 3, code, field, detail})
	} else {
		fmt.Fprintf(r.Stderr, "refused %s", code)
		if field != "" {
			fmt.Fprintf(r.Stderr, " field %s", field)
		}
		if detail != "" {
			fmt.Fprintf(r.Stderr, ": %s", detail)
		}
		fmt.Fprintln(r.Stderr)
		fmt.Fprintf(r.Stderr, "command: %s\nexit_status: 3\n", invoked)
	}
	return 3
}
func (r Runner) usage(jsonMode bool, invoked, detail string) int {
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v1", invoked, 2, "usage", "", detail})
	} else {
		fmt.Fprintln(r.Stderr, "usage:")
		for _, s := range Registry {
			fmt.Fprintf(r.Stderr, "  spectacular %s %s\n", strings.Join(s.Words, " "), s.Arguments)
		}
		fmt.Fprintln(r.Stderr, detail)
		fmt.Fprintf(r.Stderr, "command: %s\nexit_status: 2\n", invoked)
	}
	return 2
}
func writeJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(v)
}
func renderHuman(w io.Writer, envelope projection.Envelope) {
	fmt.Fprintf(w, "%s\nGenerated: %s\nBasis: %s\n", envelope.SchemaVersion, envelope.GeneratedAt, envelope.GenerationBasis)
	switch v := envelope.Data.(type) {
	case projection.ProjectView:
		fmt.Fprintf(w, "Project: %s\nSource: %s %s\nFreshness: %s (%s → %s; %s)\nAuthoritative direction: %s\n", v.Authoritative.Identity.Path, v.Source.Path, v.Source.Fingerprint, v.Freshness.State, v.Freshness.CheckedAt, v.Freshness.ValidUntil, v.Freshness.Source.Path, v.Authoritative.Direction)
		for _, item := range v.Authoritative.Boundaries {
			fmt.Fprintf(w, "Authoritative boundary: %s\n", item)
		}
		for _, item := range v.Authoritative.Constraints {
			fmt.Fprintf(w, "Authoritative constraint: %s\n", item)
		}
		for _, p := range v.Authoritative.CurrentTruth {
			fmt.Fprintf(w, "Authoritative current truth: %s [%s]\n", p.Ref, p.ShowCommand)
		}
		for _, p := range v.Projection.Missions {
			fmt.Fprintf(w, "Projected Mission: %s [%s]\n", p.Ref, p.ShowCommand)
		}
		for _, p := range v.Projection.Gaps {
			fmt.Fprintf(w, "Projected Gap: %s [%s]\n", p.Ref, p.ShowCommand)
		}
		if v.Projection.Continuation != nil {
			fmt.Fprintf(w, "Projected continuation: %s %s (authorized by %s)\n", v.Projection.Continuation.Operation, v.Projection.Continuation.Target.Ref, v.Projection.Continuation.AuthorizedBy.Ref)
		}
		if v.Projection.OwnerGate != nil {
			fmt.Fprintf(w, "Projected owner gate: %s — %s [%s]\n", v.Projection.OwnerGate.Code, v.Projection.OwnerGate.Detail, v.Projection.OwnerGate.Source.ShowCommand)
		}
	case projection.Card:
		renderCard(w, v)
	case projection.List:
		for _, card := range v.Items {
			renderCard(w, card)
		}
	case projection.Validation:
		fmt.Fprintf(w, "Scope: %s\nValid: %t\nRecords: %d\n", v.Scope, v.Valid, v.Records)
	}
}
func renderCard(w io.Writer, c projection.Card) {
	fmt.Fprintf(w, "%s %s\nTitle: %s\nFreshness: %s (%s → %s)\nSource: %s %s\n", c.Noun, c.ID, c.Title, c.Freshness.State, c.Freshness.CheckedAt, c.Freshness.ValidUntil, c.Source.Path, c.Source.Fingerprint)
	for _, p := range c.Pointers {
		fmt.Fprintf(w, "Pointer: %s %s %s %s [%s]\n", p.Noun, p.Ref, p.Path, p.Fingerprint, p.ShowCommand)
	}
	for _, p := range c.Gaps {
		fmt.Fprintf(w, "Gap: %s %s [%s]\n", p.Ref, p.Path, p.ShowCommand)
	}
	if c.Continuation != nil {
		fmt.Fprintf(w, "Continuation: %s %s (authorized by %s)\n", c.Continuation.Operation, c.Continuation.Target.Ref, c.Continuation.AuthorizedBy.Ref)
	}
	if c.OwnerGate != nil {
		fmt.Fprintf(w, "Owner gate: %s — %s [%s]\n", c.OwnerGate.Code, c.OwnerGate.Detail, c.OwnerGate.Source.ShowCommand)
	}
}
