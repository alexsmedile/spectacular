// Package command owns the complete public v2 mechanical command registry.
package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	contextcompiler "github.com/alexsmedile/spectacular/v2/internal/context"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/projection"
	spectacularruntime "github.com/alexsmedile/spectacular/v2/internal/runtime"
)

type Effect string

const ReadOnly Effect = "read-only"
const Mutating Effect = "mutating"

type Spec struct {
	Words         []string
	Arguments     string
	ArgumentShape ArgumentShape
	JSONSchema    string
	Effect        Effect
	Operation     Operation
}
type Operation uint8
type ArgumentShape uint8

const (
	argumentsNone ArgumentShape = iota
	argumentsOne
	argumentsScope
	argumentsInput
	argumentsTransition
	argumentsReconcile
	argumentsArchive
	argumentsContext
)

const (
	opAnchorShowProject Operation = iota + 1
	opAnchorShowProduct
	opAnchorShowArchitecture
	opAnchorShowStack
	opMissionList
	opMissionShow
	opGapList
	opGapShow
	opRunShow
	opCheckpointShow
	opEvidenceShow
	opDecisionShow
	opWorkspaceValidate
	opProposalShow
	opObjectiveShow
	opAssessmentShow
	opProposalCheckBase
	opProposalCreate
	opMissionCreate
	opMissionTransition
	opHandoffShow
	opHandoffValidate
	opHandoffCreate
	opHandoffReturn
	opEvidenceCreate
	opDecisionCreate
	opAssessmentRecord
	opContractShow
	opContractReconcile
	opContractReconcileSet
	opMissionArchive
	opWorkspaceContext
	opMissionPrepare
	opMissionAutopilot
)

var Registry = []Spec{
	{[]string{"anchor", "show", "project"}, "[--json]", argumentsNone, "spectacular.anchor.show.v1", ReadOnly, opAnchorShowProject},
	{[]string{"anchor", "show", "product"}, "[--json]", argumentsNone, "spectacular.anchor.show.v1", ReadOnly, opAnchorShowProduct},
	{[]string{"anchor", "show", "architecture"}, "[--json]", argumentsNone, "spectacular.anchor.show.v1", ReadOnly, opAnchorShowArchitecture},
	{[]string{"anchor", "show", "stack"}, "[--json]", argumentsNone, "spectacular.anchor.show.v1", ReadOnly, opAnchorShowStack},
	{[]string{"mission", "list"}, "[--json]", argumentsNone, "spectacular.mission.list.v1", ReadOnly, opMissionList},
	{[]string{"mission", "show"}, "<ref> [--json]", argumentsOne, "spectacular.mission.show.v1", ReadOnly, opMissionShow},
	{[]string{"gap", "list"}, "--scope <ref> [--json]", argumentsScope, "spectacular.gap.list.v1", ReadOnly, opGapList},
	{[]string{"gap", "show"}, "<ref> [--json]", argumentsOne, "spectacular.gap.show.v1", ReadOnly, opGapShow},
	{[]string{"run", "show"}, "<ref> [--json]", argumentsOne, "spectacular.run.show.v1", ReadOnly, opRunShow},
	{[]string{"checkpoint", "show"}, "<ref> [--json]", argumentsOne, "spectacular.checkpoint.show.v1", ReadOnly, opCheckpointShow},
	{[]string{"evidence", "show"}, "<ref> [--json]", argumentsOne, "spectacular.evidence.show.v1", ReadOnly, opEvidenceShow},
	{[]string{"decision", "show"}, "<ref> [--json]", argumentsOne, "spectacular.decision.show.v1", ReadOnly, opDecisionShow},
	{[]string{"workspace", "validate"}, "<scope> [--json]", argumentsOne, "spectacular.workspace.validate.v1", ReadOnly, opWorkspaceValidate},
	{[]string{"workspace", "context"}, "<project|mission-ref> --event <@Event> [--selector <$domain.verb>] [--json]", argumentsContext, "spectacular.workspace.context.v1", ReadOnly, opWorkspaceContext},
	{[]string{"proposal", "show"}, "<ref> [--json]", argumentsOne, "spectacular.proposal.show.v1", ReadOnly, opProposalShow},
	{[]string{"objective", "show"}, "<ref> [--json]", argumentsOne, "spectacular.objective.show.v1", ReadOnly, opObjectiveShow},
	{[]string{"assessment", "show"}, "<ref> [--json]", argumentsOne, "spectacular.assessment.show.v1", ReadOnly, opAssessmentShow},
	{[]string{"proposal", "check-base"}, "<ref> [--json]", argumentsOne, "spectacular.proposal.check-base.v1", ReadOnly, opProposalCheckBase},
	{[]string{"proposal", "create"}, "--input <json-file> [--json]", argumentsInput, "spectacular.proposal.create.v1", Mutating, opProposalCreate},
	{[]string{"mission", "prepare"}, "--input <json-file> [--json]", argumentsInput, "spectacular.mission.prepare.v1", ReadOnly, opMissionPrepare},
	{[]string{"mission", "create"}, "--input <json-file> [--json]", argumentsInput, "spectacular.mission.create.v1", Mutating, opMissionCreate},
	{[]string{"mission", "transition"}, "<ref> --to <state> --authorization <decision-ref> --expected-fingerprint <sha> --idempotency-key <key> [--assessment <ref>] [--reconciliation <ref>] [--disposition <value>] [--terminal-next-action <text>] [--satisfied-objectives <ref,ref>] [--json]", argumentsTransition, "spectacular.mission.transition.v1", Mutating, opMissionTransition},
	{[]string{"mission", "autopilot"}, "--input <json-file> [--json]", argumentsInput, "spectacular.mission.autopilot.v1", ReadOnly, opMissionAutopilot},
	{[]string{"handoff", "show"}, "<ref> [--json]", argumentsOne, "spectacular.handoff.show.v1", ReadOnly, opHandoffShow},
	{[]string{"handoff", "validate"}, "<ref> [--json]", argumentsOne, "spectacular.handoff.validate.v1", ReadOnly, opHandoffValidate},
	{[]string{"handoff", "create"}, "--input <json-file> [--json]", argumentsInput, "spectacular.handoff.create.v1", Mutating, opHandoffCreate},
	{[]string{"handoff", "return"}, "--input <json-file> [--json]", argumentsInput, "spectacular.handoff.return.v1", Mutating, opHandoffReturn},
	{[]string{"evidence", "create"}, "--input <json-file> [--json]", argumentsInput, "spectacular.evidence.create.v1", Mutating, opEvidenceCreate},
	{[]string{"decision", "create"}, "--input <json-file> [--json]", argumentsInput, "spectacular.decision.create.v1", Mutating, opDecisionCreate},
	{[]string{"assessment", "record"}, "--input <json-file> [--json]", argumentsInput, "spectacular.assessment.record.v1", Mutating, opAssessmentRecord},
	{[]string{"contract", "show"}, "<ref> [--json]", argumentsOne, "spectacular.contract.show.v1", ReadOnly, opContractShow},
	{[]string{"contract", "reconcile"}, "<ref> --proposal <ref> --authorization <decision-ref> --expected-fingerprint <sha|absent> --idempotency-key <key> [--json]", argumentsReconcile, "spectacular.contract.reconcile.v1", Mutating, opContractReconcile},
	{[]string{"contract", "reconcile-set"}, "--input <json-file> [--json]", argumentsInput, "spectacular.contract.reconcile-set.v1", Mutating, opContractReconcileSet},
	{[]string{"mission", "archive"}, "<ref> --authorization <decision-ref> --expected-fingerprint <sha> --idempotency-key <key> --terminal-packet <mission-ref> [--json]", argumentsArchive, "spectacular.mission.archive.v1", Mutating, opMissionArchive},
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
	if detail := spec.validateArguments(rest); detail != "" {
		return usage(detail)
	}
	workspace, err := discovery.Open(r.Cwd)
	if err != nil {
		return refuse(err)
	}
	if spec.Effect == Mutating {
		if err := governance.RecoverTransactions(workspace.Root); err != nil {
			return refuse(err)
		}
		workspace, err = discovery.Open(r.Cwd)
		if err != nil {
			return refuse(err)
		}
	}
	b := projection.Builder{Workspace: workspace, Now: r.Now, ShowCommand: registeredShowCommand}
	g := governance.Service{Workspace: workspace, Now: r.Now}
	var value any
	switch spec.Operation {
	case opAnchorShowProject:
		value, err = b.Project()
	case opAnchorShowProduct:
		value, err = b.Detail("PRODUCT", domain.Anchor)
	case opAnchorShowArchitecture:
		value, err = b.Detail("ARCHITECTURE", domain.Anchor)
	case opAnchorShowStack:
		value, err = b.Detail("STACK", domain.Anchor)
	case opMissionList:
		value, err = b.MissionList()
	case opMissionShow:
		value, err = b.Mission(rest[0])
	case opGapList:
		value, err = b.Gaps(rest[1])
	case opGapShow:
		value, err = r.detail(b, rest, domain.Gap)
	case opRunShow:
		value, err = r.detail(b, rest, domain.Run)
	case opCheckpointShow:
		value, err = r.detail(b, rest, domain.Checkpoint)
	case opEvidenceShow:
		value, err = r.detail(b, rest, domain.Evidence)
	case opDecisionShow:
		value, err = r.detail(b, rest, domain.Decision)
	case opWorkspaceValidate:
		value, err = b.Validate(rest[0])
	case opProposalShow:
		value, err = g.ProposalView(rest[0])
	case opObjectiveShow:
		value, err = r.detail(b, rest, domain.Objective)
	case opAssessmentShow:
		value, err = r.detail(b, rest, domain.Assessment)
	case opProposalCheckBase:
		value, err = g.CheckProposalBase(rest[0])
	case opProposalCreate:
		var input governance.ProposalInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CreateProposal(input)
		}
	case opMissionCreate:
		var input governance.MissionInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CreateMission(input)
		}
	case opMissionTransition:
		var input governance.TransitionInput
		input, err = transitionInput(rest)
		if err == nil {
			value, err = g.TransitionMission(input)
		}
	case opHandoffShow:
		value, err = b.Detail(rest[0], domain.Handoff)
	case opHandoffValidate:
		value, err = g.ValidateHandoff(rest[0])
	case opHandoffCreate:
		var input governance.HandoffInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CreateHandoff(input)
		}
	case opHandoffReturn:
		var input governance.HandoffReturnInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.ReturnHandoff(input)
		}
	case opEvidenceCreate:
		var input governance.EvidenceInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CreateEvidence(input)
		}
	case opDecisionCreate:
		var input governance.DecisionInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CreateDecision(input)
		}
	case opAssessmentRecord:
		var input governance.AssessmentInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.RecordAssessment(input)
		}
	case opContractShow:
		value, err = g.ContractView(rest[0])
	case opContractReconcile:
		var input governance.ReconcileInput
		input, err = reconcileInput(rest)
		if err == nil {
			value, err = g.Reconcile(input)
		}
	case opContractReconcileSet:
		var input governance.ReconcileSetInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.ReconcileMany(input.Items)
		}
	case opMissionArchive:
		var input governance.ArchiveInput
		input, err = archiveInput(rest)
		if err == nil {
			value, err = g.ArchiveMission(input)
		}
	case opWorkspaceContext:
		var config contextcompiler.Config
		config, err = contextInput(rest)
		if err == nil {
			value, err = (contextcompiler.Compiler{Workspace: workspace, Now: r.Now}).Compile(config)
		}
	case opMissionPrepare:
		var input spectacularruntime.PreparationInput
		err = readInput(rest[1], &input)
		if err == nil {
			err = validatePreparationSources(workspace, input)
		}
		if err == nil {
			value, err = spectacularruntime.CompilePreparation(input, currentTime(r.Now))
		}
	case opMissionAutopilot:
		var input spectacularruntime.AutopilotInput
		err = readInput(rest[1], &input)
		if err == nil {
			value, err = g.CompileAutopilot(input)
		}
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

func registeredShowCommand(noun domain.RecordType, ref string) (string, bool) {
	if noun == domain.Anchor {
		want := strings.ToLower(ref)
		for _, spec := range Registry {
			if len(spec.Words) == 3 && spec.Words[0] == "anchor" && spec.Words[1] == "show" && spec.Words[2] == want {
				return "spectacular " + strings.Join(spec.Words, " "), true
			}
		}
		return "", false
	}
	want := strings.ToLower(string(noun))
	for _, spec := range Registry {
		if len(spec.Words) < 2 || spec.Words[0] != want || spec.Words[1] != "show" {
			continue
		}
		command := "spectacular " + strings.Join(spec.Words, " ")
		command += " " + ref
		return command, true
	}
	return "", false
}

func (s Spec) validateArguments(args []string) string {
	switch s.ArgumentShape {
	case argumentsNone:
		if len(args) != 0 {
			return strings.Join(s.Words, " ") + " takes no arguments"
		}
	case argumentsOne:
		if len(args) != 1 {
			return strings.Join(s.Words, " ") + " requires exactly one argument"
		}
	case argumentsScope:
		if len(args) != 2 || args[0] != "--scope" {
			return strings.Join(s.Words, " ") + " requires --scope <ref>"
		}
	case argumentsInput:
		if len(args) != 2 || args[0] != "--input" || args[1] == "" {
			return strings.Join(s.Words, " ") + " requires --input <json-file>"
		}
	case argumentsTransition:
		if _, err := transitionInput(args); err != nil {
			return err.Error()
		}
	case argumentsReconcile:
		if _, err := reconcileInput(args); err != nil {
			return err.Error()
		}
	case argumentsArchive:
		if _, err := archiveInput(args); err != nil {
			return err.Error()
		}
	case argumentsContext:
		if _, err := contextInput(args); err != nil {
			return err.Error()
		}
	default:
		return "command registry has an invalid argument shape"
	}
	return ""
}

func contextInput(args []string) (contextcompiler.Config, error) {
	scope, values, err := optionMap(args, true)
	if err != nil {
		return contextcompiler.Config{}, err
	}
	if err := requireOptions(values, "--event"); err != nil {
		return contextcompiler.Config{}, err
	}
	if err := rejectUnknownOptions(values, "--event", "--selector"); err != nil {
		return contextcompiler.Config{}, err
	}
	return contextcompiler.Config{Scope: scope, Event: values["--event"], Selector: values["--selector"]}, nil
}

func currentTime(now func() time.Time) time.Time {
	if now == nil {
		return time.Now().UTC()
	}
	return now().UTC()
}

func validatePreparationSources(workspace *discovery.Workspace, input spectacularruntime.PreparationInput) error {
	if err := validateBoundSource(workspace, input.Proposal); err != nil {
		return err
	}
	for _, source := range input.DirectionSources {
		if err := validateBoundSource(workspace, source); err != nil {
			return err
		}
	}
	return nil
}

func validateBoundSource(workspace *discovery.Workspace, source spectacularruntime.BoundSource) error {
	ref, err := domain.ParseReference(source.Ref)
	if err != nil {
		return err
	}
	entry, err := workspace.Lookup(source.Ref, ref.Type)
	if err != nil {
		return err
	}
	if entry.Fingerprint != source.Fingerprint {
		return domain.NewRefusal(domain.RefusalStaleFingerprint, "source", "bound source fingerprint is stale", nil)
	}
	return nil
}

func readInput(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return domain.NewRefusal(domain.RefusalRecordNotFound, "input", "read confirmed input", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "decode confirmed JSON input", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "input must contain exactly one JSON value", err)
	}
	return nil
}

func optionMap(args []string, firstRef bool) (string, map[string]string, error) {
	ref := ""
	if firstRef {
		if len(args) == 0 || strings.HasPrefix(args[0], "--") {
			return "", nil, fmt.Errorf("requires a record reference")
		}
		ref, args = args[0], args[1:]
	}
	if len(args)%2 != 0 {
		return "", nil, fmt.Errorf("options require values")
	}
	values := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		if !strings.HasPrefix(args[i], "--") || args[i+1] == "" {
			return "", nil, fmt.Errorf("invalid option pair")
		}
		if _, exists := values[args[i]]; exists {
			return "", nil, fmt.Errorf("%s may be supplied at most once", args[i])
		}
		values[args[i]] = args[i+1]
	}
	return ref, values, nil
}

func requireOptions(values map[string]string, names ...string) error {
	for _, name := range names {
		if values[name] == "" {
			return fmt.Errorf("requires %s", name)
		}
	}
	return nil
}

func rejectUnknownOptions(values map[string]string, allowed ...string) error {
	known := map[string]bool{}
	for _, name := range allowed {
		known[name] = true
	}
	for name := range values {
		if !known[name] {
			return fmt.Errorf("unknown option %s", name)
		}
	}
	return nil
}

func transitionInput(args []string) (governance.TransitionInput, error) {
	ref, values, err := optionMap(args, true)
	if err != nil {
		return governance.TransitionInput{}, err
	}
	if err := requireOptions(values, "--to", "--authorization", "--expected-fingerprint", "--idempotency-key"); err != nil {
		return governance.TransitionInput{}, err
	}
	if err := rejectUnknownOptions(values, "--to", "--authorization", "--expected-fingerprint", "--idempotency-key", "--disposition", "--assessment", "--reconciliation", "--terminal-next-action", "--satisfied-objectives"); err != nil {
		return governance.TransitionInput{}, err
	}
	var objectives []string
	if raw := values["--satisfied-objectives"]; raw != "" {
		objectives = strings.Split(raw, ",")
	}
	return governance.TransitionInput{Mission: ref, To: values["--to"], Authorization: values["--authorization"], ExpectedFingerprint: values["--expected-fingerprint"], IdempotencyKey: values["--idempotency-key"], Disposition: values["--disposition"], Assessment: values["--assessment"], Reconciliation: values["--reconciliation"], TerminalNextAction: values["--terminal-next-action"], SatisfiedObjectives: objectives}, nil
}

func reconcileInput(args []string) (governance.ReconcileInput, error) {
	ref, values, err := optionMap(args, true)
	if err != nil {
		return governance.ReconcileInput{}, err
	}
	if err := requireOptions(values, "--proposal", "--authorization", "--expected-fingerprint", "--idempotency-key"); err != nil {
		return governance.ReconcileInput{}, err
	}
	if err := rejectUnknownOptions(values, "--proposal", "--authorization", "--expected-fingerprint", "--idempotency-key"); err != nil {
		return governance.ReconcileInput{}, err
	}
	return governance.ReconcileInput{Contract: ref, Proposal: values["--proposal"], Authorization: values["--authorization"], ExpectedFingerprint: values["--expected-fingerprint"], IdempotencyKey: values["--idempotency-key"]}, nil
}

func archiveInput(args []string) (governance.ArchiveInput, error) {
	ref, values, err := optionMap(args, true)
	if err != nil {
		return governance.ArchiveInput{}, err
	}
	if err := requireOptions(values, "--authorization", "--expected-fingerprint", "--idempotency-key", "--terminal-packet"); err != nil {
		return governance.ArchiveInput{}, err
	}
	if err := rejectUnknownOptions(values, "--authorization", "--expected-fingerprint", "--idempotency-key", "--terminal-packet"); err != nil {
		return governance.ArchiveInput{}, err
	}
	return governance.ArchiveInput{Mission: ref, Authorization: values["--authorization"], ExpectedFingerprint: values["--expected-fingerprint"], IdempotencyKey: values["--idempotency-key"], TerminalPacket: values["--terminal-packet"]}, nil
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
	Expected       string `json:"expected,omitempty"`
	Actual         string `json:"actual,omitempty"`
	Mutation       string `json:"mutation"`
	Recovery       string `json:"recovery"`
}

func (r Runner) refuse(jsonMode bool, invoked string, err error) int {
	code := "internal_error"
	field := ""
	detail := err.Error()
	expected := ""
	actual := ""
	recovery := "correct the refused field or obtain explicit owner authorization, then retry"
	var refusal *domain.Refusal
	if errors.As(err, &refusal) {
		code = string(refusal.Code)
		field = refusal.Field
		detail = refusal.Detail
		expected = refusal.Expected
		actual = refusal.Actual
		if refusal.Recovery != "" {
			recovery = refusal.Recovery
		}
	}
	if jsonMode {
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v1", invoked, 3, code, field, detail, expected, actual, "none", recovery})
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
		_ = writeJSON(r.Stdout, refusalEnvelope{"spectacular.refusal.v1", invoked, 2, "usage", "", detail, "", "", "none", "correct the command invocation using registry-derived help"})
	} else {
		fmt.Fprintln(r.Stderr, "usage:")
		fmt.Fprintf(r.Stderr, "  %s %s\n", VersionInspection.Command, VersionInspection.Arguments)
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
	switch v := envelope.Data.(type) {
	case projection.ProjectView:
		fmt.Fprintf(w, "Spectacular project orientation\n\n%s\n", v.Authoritative.Direction)
		for _, item := range v.Authoritative.Boundaries {
			fmt.Fprintf(w, "Boundary: %s\n", item)
		}
		for _, item := range v.Authoritative.Constraints {
			fmt.Fprintf(w, "Constraint: %s\n", item)
		}
		for _, p := range v.Authoritative.CurrentTruth {
			fmt.Fprintf(w, "Current truth: %-14s %s\n", displayRef(p), p.ShowCommand)
		}
		for _, p := range v.Projection.Missions {
			fmt.Fprintf(w, "Mission:       %-14s %s\n", displayRef(p), p.ShowCommand)
		}
		for _, p := range v.Projection.Gaps {
			fmt.Fprintf(w, "Blocking Gap:  %-14s %s\n", displayRef(p), p.ShowCommand)
		}
		if v.Projection.Continuation != nil {
			fmt.Fprintf(w, "\nNext: %s %s\nAuthorized by: %s\n", v.Projection.Continuation.Operation, displayRef(v.Projection.Continuation.Target), displayRef(v.Projection.Continuation.AuthorizedBy))
		}
		if v.Projection.OwnerGate != nil {
			fmt.Fprintf(w, "\nOwner gate: %s — %s\nInspect: %s\n", v.Projection.OwnerGate.Code, v.Projection.OwnerGate.Detail, v.Projection.OwnerGate.Source.ShowCommand)
		}
		fmt.Fprintf(w, "\nSource: %s\nFreshness: %s\n", v.Source.Path, v.Freshness.State)
	case projection.Card:
		renderCard(w, v)
	case projection.List:
		for _, card := range v.Items {
			renderCard(w, card)
		}
	case projection.Validation:
		fmt.Fprintf(w, "Scope: %s\nValid: %t\nRecords: %d\n", v.Scope, v.Valid, v.Records)
	default:
		data, _ := json.MarshalIndent(v, "", "  ")
		fmt.Fprintln(w, string(data))
	}
	fmt.Fprintf(w, "\nGenerated: %s\nBasis: %s\n", envelope.GeneratedAt, envelope.GenerationBasis)
}
func renderCard(w io.Writer, c projection.Card) {
	ref := c.Ref
	if ref == "" {
		ref = c.ID
	}
	fmt.Fprintf(w, "%s %s — %s\n", c.Noun, ref, c.Title)
	if c.Outcome != "" {
		fmt.Fprintf(w, "Outcome: %s\n", c.Outcome)
	}
	if c.State != "" {
		fmt.Fprintf(w, "State: %s\n", c.State)
	}
	if c.CurrentObjective != nil {
		fmt.Fprintf(w, "Objective: %s  %s\n", displayRef(*c.CurrentObjective), c.CurrentObjective.ShowCommand)
	}
	if c.CurrentRun != nil {
		fmt.Fprintf(w, "Run: %s  %s\n", displayRef(*c.CurrentRun), c.CurrentRun.ShowCommand)
	}
	if c.LatestCheckpoint != nil {
		fmt.Fprintf(w, "Checkpoint: %s  %s\n", displayRef(*c.LatestCheckpoint), c.LatestCheckpoint.ShowCommand)
	}
	for _, p := range c.Pointers {
		if (c.CurrentObjective != nil && p.Ref == c.CurrentObjective.Ref) || (c.CurrentRun != nil && p.Ref == c.CurrentRun.Ref) || (c.LatestCheckpoint != nil && p.Ref == c.LatestCheckpoint.Ref) {
			continue
		}
		fmt.Fprintf(w, "Related: %s %-14s %s\n", p.Noun, displayRef(p), p.ShowCommand)
	}
	for _, p := range c.Gaps {
		fmt.Fprintf(w, "Gap: %s  %s\n", displayRef(p), p.ShowCommand)
	}
	if c.Continuation != nil {
		fmt.Fprintf(w, "Continuation: %s %s (authorized by %s)\n", c.Continuation.Operation, displayRef(c.Continuation.Target), displayRef(c.Continuation.AuthorizedBy))
	}
	if c.OwnerGate != nil {
		fmt.Fprintf(w, "Owner gate: %s — %s\nInspect: %s\n", c.OwnerGate.Code, c.OwnerGate.Detail, c.OwnerGate.Source.ShowCommand)
	}
	fmt.Fprintf(w, "Source: %s\nFreshness: %s\n", c.Source.Path, c.Freshness.State)
}

func displayRef(pointer projection.Pointer) string {
	if pointer.HumanRef != "" {
		return pointer.HumanRef
	}
	return pointer.Ref
}
