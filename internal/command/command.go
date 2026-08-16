// Package command owns the complete public v2 mechanical command registry.
package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

type Effect string

const ReadOnly Effect = "read-only"
const Mutating Effect = "mutating"

type argumentShape uint8

const (
	one argumentShape = iota
	two
	titleOption
	byOption
)

type operation uint8

const (
	opMissionStart operation = iota + 1
	opMissionShow
	opMissionCheck
	opObjectiveShow
	opObjectivePromote
	opObjectiveFinish
	opRunShow
	opRunStart
	opReviewRecord
	opMissionComplete
)

type Spec struct {
	Words         []string
	Arguments     string
	ArgumentShape argumentShape
	JSONSchema    string
	Effect        Effect
	Operation     operation
}

var Registry = []Spec{
	{[]string{"mission", "start"}, "<plan.md|-> [--json]", one, "spectacular.mission.start.v2", Mutating, opMissionStart},
	{[]string{"mission", "show"}, "<ref> [--json]", one, "spectacular.mission.show.v2", ReadOnly, opMissionShow},
	{[]string{"mission", "check"}, "<ref> [--json]", one, "spectacular.mission.check.v2", ReadOnly, opMissionCheck},
	{[]string{"objective", "show"}, "<mission-ref>/<objective-ref> [--json]", one, "spectacular.objective.show.v2", ReadOnly, opObjectiveShow},
	{[]string{"objective", "promote"}, "<mission-ref>/<objective-ref> [--json]", one, "spectacular.objective.promote.v2", Mutating, opObjectivePromote},
	{[]string{"objective", "finish"}, "<mission-ref>/<objective-ref> [--json]", one, "spectacular.objective.finish.v2", Mutating, opObjectiveFinish},
	{[]string{"run", "show"}, "<mission-ref>/<run-ref> [--json]", one, "spectacular.run.show.v2", ReadOnly, opRunShow},
	{[]string{"run", "start"}, "<mission-ref> --title <title> [--json]", titleOption, "spectacular.run.start.v2", Mutating, opRunStart},
	{[]string{"review", "record"}, "<mission-ref> <review.md|-> [--json]", two, "spectacular.review.record.v2", Mutating, opReviewRecord},
	{[]string{"mission", "complete"}, "<ref> --by <owner> [--json]", byOption, "spectacular.mission.complete.v2", Mutating, opMissionComplete},
}

type Runner struct {
	Cwd    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Now    func() time.Time
}

type envelope struct {
	SchemaVersion string `json:"schema_version"`
	GeneratedAt   string `json:"generated_at"`
	Data          any    `json:"data"`
}

func (r Runner) Run(args []string) int {
	original := append([]string(nil), args...)
	invoked := "spectacular " + strings.Join(original, " ")
	jsonMode, duplicate := removeJSON(&args)
	if duplicate {
		return r.usage(jsonMode, invoked, "--json may be supplied at most once")
	}
	graphMode, duplicateGraph := removeFlag(&args, "--graph")
	if duplicateGraph {
		return r.usage(jsonMode, invoked, "--graph may be supplied at most once")
	}
	spec, rest, ok := match(args)
	if !ok {
		return r.usage(jsonMode, invoked, "unknown or incomplete command")
	}
	if detail := validateArguments(spec, rest); detail != "" {
		return r.commandUsage(jsonMode, invoked, spec, detail)
	}
	if graphMode && spec.Operation != opMissionShow {
		return r.commandUsage(jsonMode, invoked, spec, "--graph applies to mission show")
	}
	ws, err := discovery.Open(r.Cwd)
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	if spec.Effect == Mutating {
		if err := governance.RecoverTransactions(ws.Root); err != nil {
			return r.refuse(jsonMode, invoked, err)
		}
		ws, err = discovery.Open(r.Cwd)
		if err != nil {
			return r.refuse(jsonMode, invoked, err)
		}
	}
	service := missionbundle.Service{Workspace: ws, Now: r.Now}
	var value any
	switch spec.Operation {
	case opMissionStart:
		path := inputPath(r.Cwd, rest[0])
		stdin, readErr := r.stdinIfNeeded(rest[0])
		if readErr != nil {
			err = readErr
			break
		}
		var plan missionbundle.Plan
		var raw []byte
		plan, raw, err = missionbundle.ReadPlan(path, stdin)
		if err == nil {
			value, err = service.Start(plan, raw)
		}
	case opMissionShow:
		value, err = service.Show(rest[0])
	case opMissionCheck:
		value, err = service.Check(rest[0])
	case opObjectiveShow:
		var objective missionbundle.Objective
		objective, _, err = service.Objective(rest[0])
		value = objective
	case opObjectivePromote:
		value, err = service.PromoteObjective(rest[0])
	case opObjectiveFinish:
		value, err = service.FinishObjective(rest[0])
	case opRunShow:
		var run missionbundle.Run
		run, _, err = service.Run(rest[0])
		value = run
	case opRunStart:
		value, err = service.StartRun(rest[0], rest[2])
	case opReviewRecord:
		stdin, readErr := r.stdinIfNeeded(rest[1])
		if readErr != nil {
			err = readErr
			break
		}
		value, err = service.RecordReview(rest[0], inputPath(r.Cwd, rest[1]), stdin)
	case opMissionComplete:
		value, err = service.Complete(rest[0], rest[2])
	}
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	now := r.Now
	if now == nil {
		now = time.Now
	}
	output := envelope{SchemaVersion: spec.JSONSchema, GeneratedAt: now().UTC().Format(time.RFC3339Nano), Data: value}
	if jsonMode {
		if err := writeJSON(r.Stdout, output); err != nil {
			return r.refuse(true, invoked, err)
		}
	} else if graphMode {
		bundle, isBundle := value.(*missionbundle.Bundle)
		if !isBundle {
			return r.usage(jsonMode, invoked, "--graph applies to mission show")
		}
		fmt.Fprint(r.Stdout, bundle.Graph(terminalWidth()))
	} else {
		renderHuman(r.Stdout, value)
	}
	return 0
}

func validateArguments(spec Spec, args []string) string {
	switch spec.ArgumentShape {
	case one:
		if len(args) != 1 || args[0] == "" {
			return "requires exactly one argument"
		}
	case two:
		if len(args) != 2 || args[0] == "" || args[1] == "" {
			return "requires exactly two arguments"
		}
	case titleOption:
		if len(args) != 3 || args[0] == "" || args[1] != "--title" || args[2] == "" {
			return "requires <mission-ref> --title <title>"
		}
	case byOption:
		if len(args) != 3 || args[0] == "" || args[1] != "--by" || args[2] == "" {
			return "requires <ref> --by <owner>"
		}
	default:
		return "command registry has an invalid argument shape"
	}
	return ""
}

func (r Runner) stdinIfNeeded(path string) ([]byte, error) {
	if path != "-" {
		return nil, nil
	}
	reader := r.Stdin
	if reader == nil {
		reader = os.Stdin
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "read stdin", err)
	}
	return data, nil
}

func inputPath(cwd, path string) string {
	if path == "-" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(cwd, path)
}

// terminalWidth reports the width the Objective graph is measured against. The
// environment is consulted rather than the terminal device so the value is
// reproducible in tests and identical for piped output.
func terminalWidth() int {
	if columns := os.Getenv("COLUMNS"); columns != "" {
		if parsed, err := strconv.Atoi(columns); err == nil && parsed > 0 {
			return parsed
		}
	}
	return missionbundle.DefaultGraphWidth
}

func renderHuman(writer io.Writer, value any) {
	switch item := value.(type) {
	case *missionbundle.Bundle:
		state := item.Derive()
		fmt.Fprintf(writer, "%s — %s\n", item.Ref, item.Title)
		fmt.Fprintln(writer, stateLine(item, state))
		fmt.Fprintf(writer, "NEXT: %s (%s)\n", state.Next, state.Holder)
		fmt.Fprintf(writer, "Outcome: %s\nPath: %s\n", item.Outcome, item.Path)
		for _, objective := range state.Objectives {
			suffix := ""
			if len(objective.BlockedBy) > 0 {
				suffix = " — waits " + strings.Join(objective.BlockedBy, ", ")
			}
			fmt.Fprintf(writer, "  %s %s/%s — %s%s\n", readinessGlyph(objective.Readiness), item.Ref, objective.Ref, objective.Outcome, suffix)
		}
	case missionbundle.Check:
		fmt.Fprintf(writer, "%s valid=%t schema=%s checks=%d\n", item.Ref, item.Valid, item.Schema, len(item.Checks))
		for _, notice := range item.Notices {
			fmt.Fprintf(writer, "notice: %s\n", notice)
		}
		if len(item.Drift) > 0 {
			fmt.Fprintln(writer, "CLAIMS")
			for _, claim := range item.Drift {
				flags := "clean"
				if len(claim.Flags) > 0 {
					names := make([]string, 0, len(claim.Flags))
					for _, flag := range claim.Flags {
						names = append(names, string(flag))
					}
					flags = strings.Join(names, ", ")
				}
				fmt.Fprintf(writer, "  %-20s %s\n", claim.Claim, flags)
			}
			if count := len(item.Drift[0].Flags); count > 0 {
				noun := "flags"
				if count == 1 {
					noun = "flag"
				}
				tied := 0
				for _, claim := range item.Drift {
					if len(claim.Flags) == count {
						tied++
					}
				}
				suffix := ""
				if tied > 1 {
					suffix = fmt.Sprintf("; %d claims tied, plan order breaks it", tied)
				}
				fmt.Fprintf(writer, "  audit defaults to %s (%d %s%s)\n", item.Drift[0].Claim, count, noun, suffix)
			}
		}
		if len(item.Authority) > 0 {
			fmt.Fprintln(writer, "AUTHORITY")
			for _, decision := range []missionbundle.Decision{missionbundle.DecisionOperator, missionbundle.DecisionOwner, missionbundle.DecisionUndeclared} {
				var verbs []string
				for _, answer := range item.Authority {
					if answer.Decision == decision {
						verbs = append(verbs, answer.Verb)
					}
				}
				if len(verbs) > 0 {
					fmt.Fprintf(writer, "  %-14s %s\n", decision, strings.Join(verbs, ", "))
				}
			}
			fmt.Fprintln(writer, "  any other verb is undeclared and refused")
		}
	case missionbundle.Objective:
		fmt.Fprintf(writer, "%s — %s (%s)\n", item.Ref, item.Outcome, item.Status)
	case missionbundle.Run:
		fmt.Fprintf(writer, "%s — %s (%s)\n", item.Ref, item.Title, item.Status)
	case missionbundle.Result:
		fmt.Fprintf(writer, "%s %s\nPath: %s\n", item.Operation, item.Ref, item.Path)
	default:
		data, _ := json.MarshalIndent(value, "", "  ")
		fmt.Fprintln(writer, string(data))
	}
}

// stateLine states lifecycle position in one line. Every field is derived; the
// line is never stored, so it cannot disagree with the record it summarizes.
func stateLine(bundle *missionbundle.Bundle, state missionbundle.State) string {
	parts := []string{"State: " + state.Status}
	if state.Run != "" {
		parts = append(parts, fmt.Sprintf("run %s (%s)", state.Run, state.RunStatus))
	}
	if len(state.Objectives) > 0 {
		parts = append(parts, fmt.Sprintf("%d/%d done", state.Done, len(state.Objectives)))
		if state.Startable > 0 {
			parts = append(parts, fmt.Sprintf("%d startable", state.Startable))
		}
		if state.Blocked > 0 {
			parts = append(parts, fmt.Sprintf("%d blocked", state.Blocked))
		}
	}
	if state.Budget > 0 {
		parts = append(parts, fmt.Sprintf("repairs %d/%d", state.Repairs, state.Budget))
	}
	if len(bundle.Completion) > 0 {
		parts = append(parts, fmt.Sprintf("%d claims", len(bundle.Completion)))
	}
	return strings.Join(parts, " · ")
}

func readinessGlyph(readiness missionbundle.Readiness) string {
	switch readiness {
	case missionbundle.ReadyDone:
		return "✓"
	case missionbundle.ReadyActive:
		return "◐"
	case missionbundle.ReadyStartable:
		return "▶"
	default:
		return "·"
	}
}

func removeJSON(args *[]string) (bool, bool) {
	return removeFlag(args, "--json")
}

// removeFlag strips a boolean flag, reporting whether it was present and
// whether it was supplied more than once.
func removeFlag(args *[]string, name string) (bool, bool) {
	found, duplicate := false, false
	out := (*args)[:0]
	for _, arg := range *args {
		if arg == name {
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
		matched := true
		for i := range spec.Words {
			if args[i] != spec.Words[i] {
				matched = false
				break
			}
		}
		if matched {
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
	Problem        string `json:"problem"`
	Expected       string `json:"expected,omitempty"`
	Actual         string `json:"actual,omitempty"`
	Mutation       string `json:"mutation"`
	Correction     string `json:"safe_correction"`
}

func (r Runner) refuse(jsonMode bool, invoked string, err error) int {
	code, field, problem := "internal_error", "", err.Error()
	expected, actual := "", ""
	correction := "correct the named problem and retry; no files were changed"
	var refusal *domain.Refusal
	if errors.As(err, &refusal) {
		code, field, problem = string(refusal.Code), refusal.Field, refusal.Detail
		expected, actual = refusal.Expected, refusal.Actual
		if refusal.Recovery != "" {
			correction = refusal.Recovery
		}
	}
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v2", invoked, 3, code, field, problem, expected, actual, "none", correction})
	} else {
		fmt.Fprintf(r.Stderr, "refused %s", code)
		if field != "" {
			fmt.Fprintf(r.Stderr, " field %s", field)
		}
		fmt.Fprintf(r.Stderr, ": %s\ncorrection: %s\n", problem, correction)
	}
	return 3
}

func (r Runner) commandUsage(jsonMode bool, invoked string, spec Spec, detail string) int {
	corrected := strings.TrimSpace("spectacular " + strings.Join(spec.Words, " ") + " " + spec.Arguments)
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v2", invoked, 2, "usage", "", detail, "", "", "none", corrected})
	} else {
		fmt.Fprintf(r.Stderr, "usage: %s\n%s\n", corrected, detail)
	}
	return 2
}

func (r Runner) usage(jsonMode bool, invoked, detail string) int {
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v2", invoked, 2, "usage", "", detail, "", "", "none", "use a command from the public registry"})
	} else {
		fmt.Fprintln(r.Stderr, "usage:")
		fmt.Fprintf(r.Stderr, "  %s %s\n", VersionInspection.Command, VersionInspection.Arguments)
		for _, spec := range Registry {
			fmt.Fprintf(r.Stderr, "  spectacular %s %s\n", strings.Join(spec.Words, " "), spec.Arguments)
		}
		fmt.Fprintln(r.Stderr, detail)
	}
	return 2
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
