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

	"github.com/alexsmedile/spectacular/v2/internal/campaign"
	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/guard"
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
	twoByOption
	amendOptions
	atLeastOne
	transitionOptions
	contractCreateOptions
	amendScopeOptions
	optionalStatus
	missionStartOptions
	evidenceRecordOptions
	initOptions
	decideOptions
)

type operation uint8

const (
	opMissionStart operation = iota + 1
	opMissionList
	opMissionShow
	opMissionCheck
	opMissionAmendScope
	opMissionClose
	opObjectiveShow
	opObjectivePromote
	opObjectiveFinish
	opRunShow
	opRunStart
	opRunTransition
	opReviewRecord
	opHandoffRecord
	opEvidenceRecord
	opMissionComplete
	opProposalCheck
	opCampaignCheck
	opContractAmend
	opContractCreate
	opCharter
	opDecide
	opInit
	opGuard
)

type Spec struct {
	Words         []string
	Arguments     string
	ArgumentShape argumentShape
	JSONSchema    string
	Effect        Effect
	Operation     operation
	Description   string
	InputType     string
	OutputType    string
	Template      string
}

var Registry = []Spec{
	{
		Words:         []string{"mission", "start"},
		Arguments:     "[<plan.md|->] [--data <json>] [--allow-main] [--create-branch] [--json]",
		ArgumentShape: missionStartOptions,
		JSONSchema:    "spectacular.mission.start.v2",
		Effect:        Mutating,
		Operation:     opMissionStart,
		Description:   "Activates a new compact Mission from a valid MissionPlan markdown document, stdin, or JSON payload.",
		InputType:     "MissionPlan",
		OutputType:    "Mission",
		Template: `---
type: MissionPlan
title: <title>
owner: <owner>
contract:
  ref: <contract-ref>
outcome: <outcome>
review: independent
completion:
  - claim: <claim-name>
    pass_boundary: <boundary>
    proof_requirement: <proof>
objectives:
  - outcome: <outcome>
    claims: [<claim-name>]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [<paths>]
  semantic: [<scope>]
repair_budget: 1
dependencies: []
gaps: []
stops: [scope-drift]
---
# Mission

<Description and approach>
`,
	},
	{
		Words:         []string{"mission", "list"},
		Arguments:     "[--all] [--status <status>] [--json]",
		ArgumentShape: optionalStatus,
		JSONSchema:    "spectacular.mission.list.v2",
		Effect:        ReadOnly,
		Operation:     opMissionList,
		Description:   "Lists all discovered Missions with current status, holder, and next action.",
		OutputType:    "MissionListResult",
	},
	{
		Words:         []string{"mission", "show"},
		Arguments:     "<ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.mission.show.v2",
		Effect:        ReadOnly,
		Operation:     opMissionShow,
		Description:   "Displays the live state, objectives, handoffs, and next actions for a Mission.",
		OutputType:    "Bundle",
	},
	{
		Words:         []string{"mission", "check"},
		Arguments:     "<ref> [--verify] [--json]",
		ArgumentShape: atLeastOne,
		JSONSchema:    "spectacular.mission.check.v2",
		Effect:        ReadOnly,
		Operation:     opMissionCheck,
		Description:   "Validates structural integrity, contract drift, and claim drift of a Mission.",
		OutputType:    "Check",
	},
	{
		Words:         []string{"mission", "amend-scope"},
		Arguments:     "<ref> --add <paths> --by <owner> [--reason <text>] [--dry-run] [--json]",
		ArgumentShape: amendScopeOptions,
		JSONSchema:    "spectacular.mission.amend_scope.v2",
		Effect:        Mutating,
		Operation:     opMissionAmendScope,
		Description:   "Amends the frozen mechanical scope envelope of an active Mission with owner authorization.",
		OutputType:    "Result",
	},
	{
		Words:         []string{"mission", "close"},
		Arguments:     "<ref> [--by <owner>] [--json]",
		ArgumentShape: byOption,
		JSONSchema:    "spectacular.mission.close.v2",
		Effect:        Mutating,
		Operation:     opMissionClose,
		Description:   "Performs atomic final verification, objective completion, and closeout of a Mission with owner sign-off.",
		OutputType:    "Result",
	},
	{
		Words:         []string{"objective", "show"},
		Arguments:     "<mission-ref>/<objective-ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.objective.show.v2",
		Effect:        ReadOnly,
		Operation:     opObjectiveShow,
		Description:   "Displays details of an Objective.",
		OutputType:    "Objective",
	},
	{
		Words:         []string{"objective", "promote"},
		Arguments:     "<mission-ref>/<objective-ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.objective.promote.v2",
		Effect:        Mutating,
		Operation:     opObjectivePromote,
		Description:   "Promotes an inline Objective into a standalone document.",
		OutputType:    "Objective",
	},
	{
		Words:         []string{"objective", "finish"},
		Arguments:     "<mission-ref>/<objective-ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.objective.finish.v2",
		Effect:        Mutating,
		Operation:     opObjectiveFinish,
		Description:   "Marks an Objective as complete.",
		OutputType:    "Objective",
	},
	{
		Words:         []string{"run", "show"},
		Arguments:     "<mission-ref>/<run-ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.run.show.v2",
		Effect:        ReadOnly,
		Operation:     opRunShow,
		Description:   "Displays details of an execution Run.",
		OutputType:    "Run",
	},
	{
		Words:         []string{"run", "start"},
		Arguments:     "<mission-ref>[/<objective-ref>] --title <title> [--json]",
		ArgumentShape: titleOption,
		JSONSchema:    "spectacular.run.start.v2",
		Effect:        Mutating,
		Operation:     opRunStart,
		Description:   "Starts a new execution Run under an active Mission or Objective.",
		OutputType:    "Run",
	},
	{
		Words:         []string{"run", "transition"},
		Arguments:     "<target-ref> --to <state> --by <actor> --reason <text> [--next-action <action>] [--data <json>] [--json]",
		ArgumentShape: transitionOptions,
		JSONSchema:    "spectacular.run.transition.v2",
		Effect:        Mutating,
		Operation:     opRunTransition,
		Description:   "Transitions an active, paused, or blocked Run to a new state with mandatory attribution.",
		OutputType:    "TransitionResult",
	},
	{
		Words:         []string{"review", "record"},
		Arguments:     "<mission-ref> [<review.md|->] [--data <json>] [--json]",
		ArgumentShape: two,
		JSONSchema:    "spectacular.review.record.v2",
		Effect:        Mutating,
		Operation:     opReviewRecord,
		Description:   "Records a formal review against frozen completion criteria.",
		InputType:     "ReviewDraft",
		OutputType:    "Review",
		Template: `---
type: ReviewDraft
title: <title>
status: passed
reviewed:
  commit: <commit-sha>
  # tree: optional (auto-derived from commit if omitted)
  activation_fingerprint: <activation-fingerprint>
reviewer:
  actor: <reviewer-name>
  operator: <operator-name>
  relation_to_operator: independent
  implemented_reviewed_scope: false
  independence_basis: <basis>
  evidence: [<evidence-refs>]
claims:
  - claim: <claim-name>
    verdict: pass
findings: []
limitations: []
---
# Review

<Review observations and assessment>
`,
	},
	{
		Words:         []string{"handoff", "record"},
		Arguments:     "<mission-ref> [<handoff.md|->] [--by <sender>] [--data <json>] [--json]",
		ArgumentShape: twoByOption,
		JSONSchema:    "spectacular.handoff.record.v2",
		Effect:        Mutating,
		Operation:     opHandoffRecord,
		Description:   "Records an explicit handoff between operators or sessions.",
		InputType:     "HandoffDraft",
		OutputType:    "Handoff",
		Template: `---
type: HandoffDraft
title: <title>
task: <task in receiver's terms>
supersedes: ""
reviewed:
  commit: <commit-sha>
  # tree: optional (auto-derived from commit if omitted)
sender:
  actor: <sender-name>
  relation_to_receiver: <relation>
asserted: []
assumed: []
stops: [<stops>]
returns: [<returns>]
---
# Handoff

<Handoff context and receiver instructions>
`,
	},
	{
		Words:         []string{"evidence", "record"},
		Arguments:     "<mission-ref> [draft.md|-] [--from <test-output>] [--data <json>] [--json]",
		ArgumentShape: evidenceRecordOptions,
		JSONSchema:    "spectacular.evidence.record.v2",
		Effect:        Mutating,
		Operation:     opEvidenceRecord,
		Description:   "Records an attributable Evidence package covering Objectives, Runs, or claims on a Mission.",
		InputType:     "EvidenceDraft",
		OutputType:    "Evidence",
		Template: `---
type: EvidenceDraft
title: <title>
actor: <actor-name>
commit: <commit-sha>
# tree: optional (auto-derived from commit if omitted)
objectives: [<objective-refs>]
runs: [<run-refs>]
claims: [<claim-names>]
checks:
  - name: <check-name>
    result: pass
limitations: []
---
# Evidence

<Attributable test output and proof observations>
`,
	},
	{
		Words:         []string{"mission", "complete"},
		Arguments:     "<ref> [--by <owner>] [--json]",
		ArgumentShape: byOption,
		JSONSchema:    "spectacular.mission.complete.v2",
		Effect:        Mutating,
		Operation:     opMissionComplete,
		Description:   "Completes an active Mission with owner sign-off.",
		OutputType:    "Result",
	},
	{
		Words:         []string{"proposal", "check"},
		Arguments:     "<ref> [--json]",
		ArgumentShape: one,
		JSONSchema:    "spectacular.proposal.check.v2",
		Effect:        ReadOnly,
		Operation:     opProposalCheck,
		Description:   "Validates a Proposal document.",
		InputType:     "Proposal",
		OutputType:    "ProposalCheck",
		Template: `---
type: Proposal
id: <uuidv7>
ref: P<N>
title: <title>
status: draft
created_by: <owner>
created: "<RFC3339>"
updated: "<RFC3339>"
scope: [v2]
target_contract: Contract:<uuidv7>
# atlas: optional ../atlas/<map>.md attachment; recorded, not resolved
---
# <title>

Exploration for a possible Mission. Nothing here is frozen — this Proposal
carries no execution authority and binds only when a Mission plan freezes its
claims.
`,
	},
	{
		Words:         []string{"campaign", "check"},
		Arguments:     "<path> [--ascii] [--json]",
		ArgumentShape: atLeastOne,
		JSONSchema:    "spectacular.campaign.check.v2",
		Effect:        ReadOnly,
		Operation:     opCampaignCheck,
		Description:   "Validates a Campaign block map and renders its ordered Mermaid projection.",
		InputType:     "Campaign",
		OutputType:    "CampaignCheck",
		Template: `---
type: Campaign
schema: spectacular.campaign.v2
title: <title>
focus: <what this campaign is steering toward>
current: B2
exit_condition: <observable condition that ends the campaign>
# atlas: optional ../atlas/<map>.md attachment; recorded, not resolved
blocks:
    - ref: B1
      title: <first block title>
      state: complete
      after: []
      missions: []
    - ref: B2
      title: <second block title>
      state: active
      after: [B1]
      missions: []
---
# <title>

Campaigns are non-governing planning documents. They are skip-listed in
discovery and grant no authority.
`,
	},
	{
		Words:         []string{"contract", "amend"},
		Arguments:     "<contract-ref> --gap <gap-ref> --by <owner> [--resolution <text>] [--dry-run] [--json]",
		ArgumentShape: amendOptions,
		JSONSchema:    "spectacular.contract.amend.v2",
		Effect:        Mutating,
		Operation:     opContractAmend,
		Description:   "Amends a Contract Gap resolution with owner authorization.",
		OutputType:    "Amendment",
	},
	{
		Words:         []string{"contract", "create"},
		Arguments:     "<ref> [--title <title>] [--json]",
		ArgumentShape: contractCreateOptions,
		JSONSchema:    "spectacular.contract.create.v2",
		Effect:        Mutating,
		Operation:     opContractCreate,
		Description:   "Scaffolds a new Contract document with generated UUIDv7 and valid schema.",
		OutputType:    "Result",
	},
	{
		Words:         []string{"charter"},
		Arguments:     "<mission-ref>/<objective-ref> [sources...] [--prompt] [--json]",
		ArgumentShape: atLeastOne,
		JSONSchema:    "spectacular.charter.show.v2",
		Effect:        ReadOnly,
		Operation:     opCharter,
		Description:   "Compiles and displays a 3-layer Context Sandwich charter for an Objective.",
		OutputType:    "Charter",
	},
	{
		Words:         []string{"decide"},
		Arguments:     "[<decision.md|->] [--title <title>] [--disposition <accepted|rejected|deferred|superseded>] [--rationale <rationale>] [--actor <name>] [--supersedes <ref>] [--json]",
		ArgumentShape: decideOptions,
		JSONSchema:    "spectacular.decision.record.v2",
		Effect:        Mutating,
		Operation:     opDecide,
		Description:   "Atomically validates and records an architectural Decision package.",
		InputType:     "DecisionDraft",
		OutputType:    "DecisionResult",
		Template: `---
type: DecisionDraft
title: <title>
actor: <actor>
actor_role: owner
question: <question>
disposition: <disposition>
rationale: <rationale>
alternatives: []
scope: [v2]
targets: []
supersedes: ""
---
# <title>

<Detailed decision context and impact>
`,
	},
	{
		Words:         []string{"init"},
		Arguments:     "[<path>] [--name <name>] [--json]",
		ArgumentShape: initOptions,
		JSONSchema:    "spectacular.init.v2",
		Effect:        Mutating,
		Operation:     opInit,
		Description:   "Initializes a new Spectacular workspace safely without overwriting existing files.",
		OutputType:    "InitResult",
	},
	{
		Words:         []string{"guard"},
		Arguments:     "<mission-ref>/<objective-ref> [--watch] [--exec <command>] [--json] [-- <command...>]",
		ArgumentShape: atLeastOne,
		JSONSchema:    "spectacular.guard.v2",
		Effect:        ReadOnly,
		Operation:     opGuard,
		Description:   "Executes a command under mechanical perimeter supervision with automatic rollback on escapes.",
		OutputType:    "GuardResult",
	},
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
	helpMode, duplicateHelp := removeFlag(&args, "--help")
	hMode, duplicateH := removeFlag(&args, "-h")
	if duplicateHelp || duplicateH {
		return r.usage(jsonMode, invoked, "--help/-h may be supplied at most once")
	}
	help := helpMode || hMode

	schemaMode, duplicateSchema := removeFlag(&args, "--schema")
	if duplicateSchema {
		return r.usage(jsonMode, invoked, "--schema may be supplied at most once")
	}

	graphMode, duplicateGraph := removeFlag(&args, "--graph")
	if duplicateGraph {
		return r.usage(jsonMode, invoked, "--graph may be supplied at most once")
	}
	timelineMode, duplicateTimeline := removeFlag(&args, "--timeline")
	if duplicateTimeline {
		return r.usage(jsonMode, invoked, "--timeline may be supplied at most once")
	}
	if graphMode && timelineMode {
		return r.usage(jsonMode, invoked, "cannot combine --graph and --timeline")
	}
	dryRun, duplicateDryRun := removeFlag(&args, "--dry-run")
	if duplicateDryRun {
		return r.usage(jsonMode, invoked, "--dry-run may be supplied at most once")
	}
	allowMain, duplicateAllowMain := removeFlag(&args, "--allow-main")
	if duplicateAllowMain {
		return r.usage(jsonMode, invoked, "--allow-main may be supplied at most once")
	}
	createBranch, duplicateCreateBranch := removeFlag(&args, "--create-branch")
	if duplicateCreateBranch {
		return r.usage(jsonMode, invoked, "--create-branch may be supplied at most once")
	}
	override, _, badOverride := removeValueFlag(&args, "--resolution")
	if badOverride {
		return r.usage(jsonMode, invoked, "--resolution requires exactly one value")
	}
	allMissions, duplicateAll := removeFlag(&args, "--all")
	if duplicateAll {
		return r.usage(jsonMode, invoked, "--all may be supplied at most once")
	}
	fromFile, _, badFrom := removeValueFlag(&args, "--from")
	if badFrom {
		return r.usage(jsonMode, invoked, "--from requires exactly one value")
	}
	statusFilter, _, badStatus := removeValueFlag(&args, "--status")
	if badStatus {
		return r.usage(jsonMode, invoked, "--status requires exactly one value")
	}
	nameOpt, _, badName := removeValueFlag(&args, "--name")
	if badName {
		return r.usage(jsonMode, invoked, "--name requires exactly one value")
	}
	dataPayload, _, badData := removeValueFlag(&args, "--data")
	if badData {
		return r.usage(jsonMode, invoked, "--data requires exactly one value")
	}
	jsonPayloadOpt, _, badJSONPayload := removeValueFlag(&args, "--json-payload")
	if badJSONPayload {
		return r.usage(jsonMode, invoked, "--json-payload requires exactly one value")
	}
	if dataPayload == "" {
		dataPayload = jsonPayloadOpt
	}

	if len(args) == 0 {
		if help {
			return r.globalHelp(jsonMode)
		}
		return r.usage(jsonMode, invoked, "unknown or incomplete command")
	}

	spec, rest, ok := match(args)
	if !ok {
		if help {
			return r.globalHelp(jsonMode)
		}
		return r.usage(jsonMode, invoked, "unknown or incomplete command")
	}

	if help {
		return r.commandHelp(jsonMode, spec)
	}
	if schemaMode {
		return r.commandSchema(jsonMode, spec)
	}

	if detail := validateArguments(spec, rest); detail != "" {
		return r.commandUsage(jsonMode, invoked, spec, detail)
	}
	if graphMode && spec.Operation != opMissionShow {
		return r.commandUsage(jsonMode, invoked, spec, "--graph applies to mission show")
	}
	if timelineMode && spec.Operation != opMissionShow {
		return r.commandUsage(jsonMode, invoked, spec, "--timeline applies to mission show")
	}
	if dryRun && spec.Operation != opContractAmend && spec.Operation != opMissionAmendScope {
		return r.commandUsage(jsonMode, invoked, spec, "--dry-run applies to contract amend and mission amend-scope")
	}
	if (allowMain || createBranch) && spec.Operation != opMissionStart {
		return r.commandUsage(jsonMode, invoked, spec, "--allow-main and --create-branch apply to mission start")
	}
	if allMissions && spec.Operation != opMissionList {
		return r.commandUsage(jsonMode, invoked, spec, "--all applies to mission list")
	}
	if fromFile != "" && spec.Operation != opEvidenceRecord {
		return r.commandUsage(jsonMode, invoked, spec, "--from applies to evidence record")
	}
	if statusFilter != "" && spec.Operation != opMissionList {
		return r.commandUsage(jsonMode, invoked, spec, "--status applies to mission list")
	}
	if nameOpt != "" && spec.Operation != opInit {
		return r.commandUsage(jsonMode, invoked, spec, "--name applies to init")
	}
	if override != "" && spec.Operation != opContractAmend {
		return r.commandUsage(jsonMode, invoked, spec, "--resolution applies to contract amend")
	}
	if spec.Operation == opInit {
		targetDir := r.Cwd
		if len(rest) == 1 {
			targetDir = inputPath(r.Cwd, rest[0])
		}
		initRes, initErr := InitWorkspace(targetDir, nameOpt)
		if initErr != nil {
			return r.refuse(jsonMode, invoked, initErr)
		}
		now := r.Now
		if now == nil {
			now = time.Now
		}
		output := envelope{SchemaVersion: spec.JSONSchema, GeneratedAt: now().UTC().Format(time.RFC3339Nano), Data: initRes}
		if jsonMode {
			if err := writeJSON(r.Stdout, output); err != nil {
				return r.refuse(true, invoked, err)
			}
		} else {
			renderHuman(r.Stdout, initRes)
		}
		return 0
	}
	ws, err := discovery.Open(r.Cwd)
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	if spec.Effect == Mutating && !dryRun {
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
		var plan missionbundle.Plan
		var raw []byte
		if dataPayload != "" {
			plan, raw, err = missionbundle.ReadPlan("", []byte(dataPayload))
		} else if len(rest) > 0 && strings.HasPrefix(strings.TrimSpace(rest[0]), "{") {
			plan, raw, err = missionbundle.ReadPlan("", []byte(rest[0]))
		} else if len(rest) > 0 {
			path := inputPath(r.Cwd, rest[0])
			stdin, readErr := r.stdinIfNeeded(rest[0])
			if readErr != nil {
				err = readErr
				break
			}
			plan, raw, err = missionbundle.ReadPlan(path, stdin)
		} else {
			err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "requires a plan file, stdin (-), inline JSON, or --data", nil)
			break
		}
		if allowMain {
			plan.AllowMain = true
		}
		if createBranch {
			plan.CreateBranch = true
		}
		if err == nil {
			value, err = service.Start(plan, raw)
		}
	case opMissionList:
		value, err = service.ListMissions(statusFilter, allMissions)
	case opMissionShow:
		value, err = service.Show(rest[0])
	case opMissionCheck:
		targetRef := ""
		verifyMode := false
		for _, arg := range rest {
			if arg == "--verify" {
				verifyMode = true
			} else if targetRef == "" {
				targetRef = arg
			}
		}
		if targetRef == "" {
			err = domain.NewRefusal(domain.RefusalInvalidReference, "", "expected <mission-ref>", nil)
			break
		}
		if verifyMode {
			checkVal, checkErr := service.CheckWithVerify(targetRef)
			if checkErr == nil {
				if bundle, loadErr := missionbundle.Load(ws, targetRef); loadErr == nil && len(bundle.Objectives) > 0 {
					totalTokens := 0
					counted := 0
					for _, obj := range bundle.Objectives {
						if c, compileErr := charter.Compile(ws, bundle.Ref, obj.Ref, nil); compileErr == nil && c != nil {
							if count, countErr := tokenizer.Count(c.RenderPrompt()); countErr == nil {
								totalTokens += count
								counted++
							}
						}
					}
					if counted > 0 {
						avg := totalTokens / counted
						status := "within-sweet-spot"
						if avg > 900 {
							status = "over-budget"
						} else if avg > 600 {
							status = "under-ceiling"
						}
						checkVal.TokenEfficiency = &missionbundle.TokenEfficiency{
							TotalTokens:    totalTokens,
							AverageTokens:  avg,
							ObjectiveCount: counted,
							BudgetStatus:   status,
						}
					}
				}
			}
			value, err = checkVal, checkErr
		} else {
			value, err = service.Check(targetRef)
		}
	case opMissionAmendScope:
		addPaths := strings.Split(rest[2], ",")
		var cleaned []string
		for _, p := range addPaths {
			if t := strings.TrimSpace(p); t != "" {
				cleaned = append(cleaned, t)
			}
		}
		reason := ""
		if len(rest) == 7 && rest[5] == "--reason" {
			reason = rest[6]
		}
		value, err = service.AmendScope(rest[0], cleaned, rest[4], reason, dryRun)
	case opMissionClose:
		by := ""
		if len(rest) == 3 && rest[1] == "--by" {
			by = rest[2]
		}
		if by == "" {
			by = ws.Config.Defaults.Operator
		}
		if by == "" {
			by = "Alex"
		}
		value, err = service.CloseMission(rest[0], by)
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
		missionRef := rest[0]
		if dataPayload != "" {
			var draft missionbundle.ReviewDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.ReviewDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		} else if len(rest) > 1 {
			stdin, readErr := r.stdinIfNeeded(rest[1])
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordReview(missionRef, inputPath(r.Cwd, rest[1]), stdin)
		} else {
			draft := missionbundle.ReviewDraft{
				Type:   "ReviewDraft",
				Title:  "Mission review for " + missionRef,
				Status: "passed",
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		}
	case opHandoffRecord:
		missionRef := rest[0]
		sender := ""
		for i := 1; i < len(rest); i++ {
			if (rest[i] == "--by" || rest[i] == "--sender") && i+1 < len(rest) {
				sender = rest[i+1]
			}
		}
		if dataPayload != "" {
			var draft missionbundle.HandoffDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.HandoffDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		} else if len(rest) > 1 && rest[1] != "--by" && !strings.HasPrefix(rest[1], "--") {
			stdin, readErr := r.stdinIfNeeded(rest[1])
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordHandoff(missionRef, inputPath(r.Cwd, rest[1]), sender, stdin)
		} else {
			draft := missionbundle.HandoffDraft{
				Type:  "HandoffDraft",
				Title: "Handoff for " + missionRef,
				Task:  "Continue mission implementation",
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		}
	case opEvidenceRecord:
		missionRef := rest[0]
		if dataPayload != "" {
			var draft missionbundle.EvidenceDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordEvidenceDraft(missionRef, draft)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.EvidenceDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordEvidenceDraft(missionRef, draft)
		} else {
			targetPath := ""
			if len(rest) > 1 && !strings.HasPrefix(rest[1], "--") {
				targetPath = rest[1]
			}
			stdin, readErr := r.stdinIfNeeded(targetPath)
			if readErr != nil {
				err = readErr
				break
			}
			var absTarget string
			if targetPath != "" {
				absTarget = inputPath(r.Cwd, targetPath)
			}
			var absFrom string
			if fromFile != "" {
				absFrom = inputPath(r.Cwd, fromFile)
			}
			value, err = service.RecordEvidence(missionRef, absTarget, stdin, absFrom)
		}
	case opMissionComplete:
		by := ""
		if len(rest) == 3 && rest[1] == "--by" {
			by = rest[2]
		}
		if by == "" {
			by = ws.Config.Defaults.Operator
		}
		if by == "" {
			by = "Alex"
		}
		value, err = service.Complete(rest[0], by)
	case opProposalCheck:
		value, err = missionbundle.ValidateProposal(ws, rest[0])
	case opCampaignCheck:
		targetPath := ""
		asciiMode := false
		for _, arg := range rest {
			if arg == "--ascii" {
				asciiMode = true
			} else if targetPath == "" {
				targetPath = arg
			}
		}
		if targetPath == "" {
			err = domain.NewRefusal(domain.RefusalInvalidWorkspacePath, "", "expected <campaign-path>", nil)
			break
		}
		var c campaign.Check
		c, err = campaign.Validate(ws, targetPath)
		if err == nil {
			c.ASCIIMode = asciiMode
			value = c
		}
	case opContractAmend:
		value, err = service.AmendContract(rest[0], rest[2], rest[4], override, dryRun)
	case opContractCreate:
		title := ""
		if len(rest) == 3 && rest[1] == "--title" {
			title = rest[2]
		}
		value, err = service.CreateContract(rest[0], title, ws.Config.Defaults.Operator)
	case opCharter:
		targetRef := ""
		promptMode := false
		var extraSources []string
		for _, arg := range rest {
			if arg == "--prompt" {
				promptMode = true
			} else if targetRef == "" {
				targetRef = arg
			} else {
				extraSources = append(extraSources, arg)
			}
		}
		if targetRef == "" {
			err = domain.NewRefusal(domain.RefusalInvalidReference, "", "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
			break
		}
		parts := strings.Split(targetRef, "/")
		if len(parts) != 2 {
			err = domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
			break
		}
		var c *charter.Charter
		c, err = charter.Compile(ws, parts[0], parts[1], extraSources)
		if err == nil && c != nil {
			c.PromptMode = promptMode
			value = c
		}
	case opDecide:
		flags, parseErr := parseDecideArgs(rest)
		if parseErr != nil {
			err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", parseErr.Error(), nil)
			break
		}
		if dataPayload != "" {
			flags.jsonPayload = dataPayload
		}
		if flags.jsonPayload != "" {
			var draft missionbundle.DecisionDraft
			if unmarshalErr := json.Unmarshal([]byte(flags.jsonPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			actor := draft.Actor
			if actor == "" {
				actor = flags.actor
			}
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			draft.Actor = actor
			if draft.ActorRole == "" {
				draft.ActorRole = "owner"
			}
			value, err = service.RecordDecisionDraft(draft)
		} else if flags.file != "" {
			stdin, readErr := r.stdinIfNeeded(flags.file)
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordDecision(inputPath(r.Cwd, flags.file), stdin)
		} else {
			if strings.TrimSpace(flags.title) == "" || strings.TrimSpace(flags.disposition) == "" || strings.TrimSpace(flags.rationale) == "" {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "direct decision recording requires --title, --disposition, and --rationale", nil)
				break
			}
			actor := flags.actor
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			draft := missionbundle.DecisionDraft{
				Type:        "DecisionDraft",
				Title:       flags.title,
				Disposition: flags.disposition,
				Rationale:   flags.rationale,
				Actor:       actor,
				ActorRole:   "owner",
				Question:    flags.question,
				Supersedes:  flags.supersedes,
				Scope:       flags.scope,
				Targets:     flags.targets,
			}
			value, err = service.RecordDecisionDraft(draft)
		}
	case opRunTransition:
		targetRef := rest[0]
		if dataPayload != "" || (len(rest) > 0 && strings.HasPrefix(strings.TrimSpace(rest[len(rest)-1]), "{")) {
			rawJSON := dataPayload
			if rawJSON == "" {
				rawJSON = strings.TrimSpace(rest[len(rest)-1])
			}
			var payload struct {
				Target     string `json:"target"`
				To         string `json:"to"`
				By         string `json:"by"`
				Reason     string `json:"reason"`
				NextAction string `json:"next_action"`
			}
			if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+err.Error(), nil)
				break
			}
			if payload.Target != "" {
				targetRef = payload.Target
			}
			actor := payload.By
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			value, err = service.TransitionRun(targetRef, payload.To, actor, payload.Reason, payload.NextAction)
		} else {
			toState := ""
			actor := ""
			reason := ""
			nextAction := ""
			for i := 1; i < len(rest); i++ {
				switch rest[i] {
				case "--to":
					if i+1 < len(rest) {
						toState = rest[i+1]
						i++
					}
				case "--by":
					if i+1 < len(rest) {
						actor = rest[i+1]
						i++
					}
				case "--reason":
					if i+1 < len(rest) {
						reason = rest[i+1]
						i++
					}
				case "--next-action":
					if i+1 < len(rest) {
						nextAction = rest[i+1]
						i++
					}
				}
			}
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			if toState == "" || reason == "" {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "run transition requires --to <state> and --reason <text>", nil)
				break
			}
			value, err = service.TransitionRun(targetRef, toState, actor, reason, nextAction)
		}
	case opGuard:
		targetRef := ""
		watchMode := false
		execCmd := ""
		var cmdArgs []string
		inCmd := false
		for i := 0; i < len(rest); i++ {
			arg := rest[i]
			if inCmd {
				cmdArgs = append(cmdArgs, arg)
			} else if arg == "--" {
				inCmd = true
			} else if arg == "--watch" {
				watchMode = true
			} else if arg == "--exec" && i+1 < len(rest) {
				execCmd = rest[i+1]
				i++
			} else if targetRef == "" {
				targetRef = arg
			}
		}
		if targetRef == "" || (len(cmdArgs) == 0 && execCmd == "") {
			err = domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected spectacular guard <mission-ref>/<objective-ref> [--watch] [--exec <command>] -- <command...>", nil)
			break
		}
		value, err = guard.Run(ws, targetRef, watchMode, execCmd, cmdArgs)
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
	} else if timelineMode {
		bundle, isBundle := value.(*missionbundle.Bundle)
		if !isBundle {
			return r.usage(jsonMode, invoked, "--timeline applies to mission show")
		}
		fmt.Fprint(r.Stdout, bundle.Timeline(ws, terminalWidth()))
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
		if len(args) < 1 || args[0] == "" {
			return "requires <mission-ref> [<review.md|->]"
		}
		if len(args) > 2 {
			return "takes at most two arguments"
		}
	case titleOption:
		if len(args) != 3 || args[0] == "" || args[1] != "--title" || args[2] == "" {
			return "requires <mission-ref> --title <title>"
		}
	case byOption:
		if len(args) == 1 {
			if args[0] == "" {
				return "requires <ref> [--by <owner>]"
			}
		} else if len(args) == 3 {
			if args[0] == "" || args[1] != "--by" || args[2] == "" {
				return "requires <ref> [--by <owner>]"
			}
		} else {
			return "requires <ref> [--by <owner>]"
		}
	case twoByOption:
		if len(args) < 1 || args[0] == "" {
			return "requires <mission-ref> [<handoff.md|->] [--by <sender>]"
		}
	case amendOptions:
		if len(args) != 5 || args[0] == "" || args[1] != "--gap" || args[2] == "" || args[3] != "--by" || args[4] == "" {
			return "requires <contract-ref> --gap <gap-ref> --by <owner>"
		}
	case atLeastOne:
		if len(args) < 1 || args[0] == "" {
			return "requires at least one argument"
		}
	case transitionOptions:
		if len(args) < 1 || args[0] == "" {
			return "requires <target-ref> --to <state> --by <actor> --reason <text> [--next-action <action>]"
		}
	case contractCreateOptions:
		if len(args) == 1 {
			if args[0] == "" {
				return "requires <ref> [--title <title>]"
			}
		} else if len(args) == 3 {
			if args[0] == "" || args[1] != "--title" || args[2] == "" {
				return "requires <ref> [--title <title>]"
			}
		} else {
			return "requires <ref> [--title <title>]"
		}
	case amendScopeOptions:
		if len(args) != 5 && len(args) != 7 {
			return "requires <ref> --add <paths> --by <owner> [--reason <text>]"
		}
		if args[0] == "" || args[1] != "--add" || args[2] == "" || args[3] != "--by" || args[4] == "" {
			return "requires <ref> --add <paths> --by <owner> [--reason <text>]"
		}
		if len(args) == 7 && (args[5] != "--reason" || args[6] == "") {
			return "requires [--reason <text>]"
		}
	case optionalStatus:
		if len(args) != 0 {
			return "takes no positional arguments"
		}
	case missionStartOptions:
		if len(args) > 1 {
			return "takes at most one plan argument"
		}
	case evidenceRecordOptions:
		if len(args) < 1 || args[0] == "" {
			return "requires <mission-ref> [draft.md|-]"
		}
	case initOptions:
		if len(args) > 1 {
			return "takes at most one path argument"
		}
	case decideOptions:
		if len(args) == 0 {
			return "requires <decision.md|-> or --title, --disposition, and --rationale flags"
		}
	default:
		return "command registry has an invalid argument shape"
	}
	return ""
}

type decideFlags struct {
	file        string
	jsonPayload string
	title       string
	disposition string
	rationale   string
	actor       string
	question    string
	supersedes  string
	scope       []string
	targets     []string
}

func parseDecideArgs(args []string) (decideFlags, error) {
	var flags decideFlags
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--data", "--json-payload":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("%s requires a non-empty value", arg)
			}
			flags.jsonPayload = args[i+1]
			i++
		case "--title":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--title requires a non-empty value")
			}
			flags.title = args[i+1]
			i++
		case "--disposition":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--disposition requires a non-empty value")
			}
			flags.disposition = args[i+1]
			i++
		case "--rationale":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--rationale requires a non-empty value")
			}
			flags.rationale = args[i+1]
			i++
		case "--actor", "--by":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--actor/--by requires a non-empty value")
			}
			flags.actor = args[i+1]
			i++
		case "--question":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--question requires a non-empty value")
			}
			flags.question = args[i+1]
			i++
		case "--supersedes":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--supersedes requires a non-empty value")
			}
			flags.supersedes = args[i+1]
			i++
		case "--scope":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--scope requires a non-empty value")
			}
			flags.scope = strings.Split(args[i+1], ",")
			i++
		case "--targets":
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return flags, fmt.Errorf("--targets requires a non-empty value")
			}
			flags.targets = strings.Split(args[i+1], ",")
			i++
		default:
			if strings.HasPrefix(arg, "--") {
				return flags, fmt.Errorf("unknown flag %s", arg)
			}
			trimmed := strings.TrimSpace(arg)
			if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
				flags.jsonPayload = trimmed
				continue
			}
			if flags.file != "" {
				return flags, fmt.Errorf("multiple input files specified")
			}
			flags.file = arg
		}
	}
	return flags, nil
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
	case missionbundle.MissionListResult:
		if len(item.Missions) == 0 {
			fmt.Fprintln(writer, "No missions found.")
			return
		}
		fmt.Fprintf(writer, "%-6s %-12s %-40s %-10s %s\n", "REF", "STATUS", "TITLE", "HOLDER", "NEXT")
		for _, m := range item.Missions {
			title := m.Title
			if len(title) > 38 {
				title = title[:35] + "..."
			}
			holder := m.Holder
			if holder == "" {
				holder = "-"
			}
			fmt.Fprintf(writer, "%-6s %-12s %-40s %-10s %s\n", m.Ref, m.Status, title, holder, m.Next)
		}
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
		if len(item.Handoffs) > 0 {
			fmt.Fprintln(writer, "HANDOFFS")
			superseded := map[string]string{}
			for _, handoff := range item.Handoffs {
				if handoff.Document != nil && handoff.Document.Supersedes != "" {
					superseded[handoff.Document.Supersedes] = handoff.Ref
				}
			}
			for _, handoff := range item.Handoffs {
				title, sender, binding := "", "", ""
				if handoff.Document != nil {
					title = handoff.Document.Title
					sender = handoff.Document.Sender.Actor
					// A receiver acts on the state the sender bound. Saying whether
					// that tree is still the working tree is the difference between
					// a pointer and a usable instruction.
					binding = "tree moved since it was sent"
					if handoff.Document.TreeCurrent {
						binding = "tree matches"
					}
				}
				fmt.Fprintf(writer, "  %s/%s — %s (from %s; %s)\n", item.Ref, handoff.Ref, title, sender, binding)
				// A reader arriving at a superseded Handoff is pointed forward
				// rather than left to act on a correction that already happened.
				// Through a chain, the pointer names the record that is current
				// rather than the next link, which would leave the reader walking.
				if replacement, ok := superseded[handoff.Ref]; ok {
					if newest := missionbundle.NewestHandoff(item, handoff.Ref); newest != nil && newest.Ref != replacement {
						fmt.Fprintf(writer, "      superseded by %s/%s (current: %s/%s)\n", item.Ref, replacement, item.Ref, newest.Ref)
					} else {
						fmt.Fprintf(writer, "      superseded by %s/%s\n", item.Ref, replacement)
					}
				}
			}
		}
		if len(item.Fallbacks) > 0 && state.Budget > 0 && state.Repairs >= state.Budget {
			fmt.Fprintln(writer, "FALLBACKS")
			for _, fb := range item.Fallbacks {
				rec := ""
				if fb.Recommendation {
					rec = " [recommendation]"
				}
				fmt.Fprintf(writer, "  - %s%s (invalidated if: %s)\n", fb.Approach, rec, fb.InvalidatedIf)
			}
		}
	case missionbundle.ProposalCheck:
		fmt.Fprintf(writer, "%s valid=%t checks=%d\n", item.Ref, item.Valid, len(item.Checks))
		for _, notice := range item.Notices {
			fmt.Fprintf(writer, "notice: %s\n", notice)
		}
	case campaign.Check:
		if item.ASCIIMode {
			fmt.Fprint(writer, item.RenderASCII())
		} else {
			fmt.Fprintf(writer, "Campaign: %s\n", item.Title)
			if item.StrategicGoal != "" {
				fmt.Fprintf(writer, "Outcome: %s\n", item.StrategicGoal)
			}
			fmt.Fprintf(writer, "CURRENT CAMPAIGN BLOCK: %s — %s\n", item.CurrentBlock.Ref, item.CurrentBlock.Title)
			if len(item.CurrentBlock.Missions) > 0 {
				fmt.Fprintf(writer, "LINKED MISSIONS: %s\n", strings.Join(item.CurrentBlock.Missions, ", "))
			}
			if len(item.Next) > 0 {
				next := make([]string, 0, len(item.Next))
				for _, block := range item.Next {
					next = append(next, block.Ref+" — "+block.Title)
				}
				fmt.Fprintf(writer, "NEXT MAP BLOCKS: %s\n", strings.Join(next, ", "))
			}
			if item.ExitCondition != "" {
				fmt.Fprintf(writer, "Exit: %s\n", item.ExitCondition)
			}
			fmt.Fprintln(writer, "ORDER")
			for index, block := range item.Order {
				fmt.Fprintf(writer, "  %d. %s\n", index+1, block)
			}
			fmt.Fprintln(writer, "MERMAID")
			fmt.Fprint(writer, item.Mermaid)
		}
	case missionbundle.Amendment:
		verb := "amended"
		switch {
		case item.DryRun:
			verb = "would amend"
		case item.NoOp:
			verb = "already closed on"
		}
		fmt.Fprintf(writer, "%s %s\n", verb, item.Contract)
		fmt.Fprintf(writer, "  gaps.%s: blocked_on -> resolution (declared by %s)\n", item.Gap, item.Mission)
		if !item.NoOp {
			fmt.Fprintf(writer, "  %s\n", item.Resolution)
			fmt.Fprintf(writer, "  fingerprint %s -> %s\n", short(item.From), short(item.To))
			if len(item.Repointed) > 0 {
				fmt.Fprintf(writer, "  re-points contract.fingerprint on %s\n", strings.Join(item.Repointed, " "))
			}
			fmt.Fprintf(writer, "  logged in %s\n", item.Log)
		}
		if item.DryRun {
			fmt.Fprintln(writer, "no files written")
		}
	case missionbundle.Check:
		fmt.Fprintf(writer, "%s valid=%t schema=%s checks=%d", item.Ref, item.Valid, item.Schema, len(item.Checks))
		if item.ContractVersion > 0 {
			fmt.Fprintf(writer, " contract=v%d", item.ContractVersion)
		}
		fmt.Fprintln(writer)
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
		if item.TokenEfficiency != nil {
			fmt.Fprintf(writer, "TOKEN EFFICIENCY: avg %d tokens/obj (%s, %d objectives total %d tokens)\n",
				item.TokenEfficiency.AverageTokens, item.TokenEfficiency.BudgetStatus, item.TokenEfficiency.ObjectiveCount, item.TokenEfficiency.TotalTokens)
		}
	case missionbundle.Objective:
		fmt.Fprintf(writer, "%s — %s (%s)\n", item.Ref, item.Outcome, item.Status)
	case missionbundle.Run:
		fmt.Fprintf(writer, "%s — %s (%s)\n", item.Ref, item.Title, item.Status)
	case missionbundle.Result:
		fmt.Fprintf(writer, "%s %s\nPath: %s\n", item.Operation, item.Ref, item.Path)
	case *charter.Charter:
		if item.PromptMode {
			fmt.Fprint(writer, item.RenderPrompt())
		} else {
			fmt.Fprint(writer, item.RenderMarkdown())
			fmt.Fprintln(writer, "\n---")
			fmt.Fprintf(writer, "Tokens: %d (%s, %s)\n", item.TokenCount, item.Disposition, tokenizer.Version)
			if item.Compacted {
				fmt.Fprintln(writer, "Safe compaction applied to stay within token budget.")
			}
		}
	case missionbundle.DecisionResult:
		fmt.Fprintf(writer, "Recorded decision %s (%s)\nPath: %s\n", item.Ref, item.ID, item.Path)
		if len(item.Unblocked) > 0 {
			fmt.Fprintf(writer, "Unblocked objectives: %s\n", strings.Join(item.Unblocked, ", "))
		}
		for _, changed := range item.Changed {
			fmt.Fprintf(writer, "  updated: %s\n", changed)
		}
	case missionbundle.TransitionResult:
		fmt.Fprintf(writer, "Transitioned run %s (%s -> %s)\nBy: %s\nReason: %s\nPath: %s\n", item.Ref, item.From, item.To, item.By, item.Reason, item.Path)
		for _, changed := range item.Changed {
			fmt.Fprintf(writer, "  updated: %s\n", changed)
		}
	case InitResult:
		if item.AlreadyInitialized {
			fmt.Fprintf(writer, "Spectacular workspace already initialized in %s\n", item.MetadataDir)
			if len(item.SkippedFiles) > 0 {
				fmt.Fprintf(writer, "Preserved existing files: %s\n", strings.Join(item.SkippedFiles, ", "))
			}
			return
		}
		fmt.Fprintf(writer, "Initialized Spectacular workspace in %s\n", item.MetadataDir)
		if len(item.CreatedFiles) > 0 {
			fmt.Fprintf(writer, "Created: %s\n", strings.Join(item.CreatedFiles, ", "))
		}
		if len(item.SkippedFiles) > 0 {
			fmt.Fprintf(writer, "Preserved existing files: %s\n", strings.Join(item.SkippedFiles, ", "))
		}
	case *guard.GuardResult:
		if item.Status == "pass" {
			fmt.Fprintf(writer, "Perimeter Guard: PASS (mode: %s)\nTarget: %s/%s\nExit code: %d\n", item.Mode, item.MissionRef, item.ObjectiveRef, item.ExitCode)
			if item.Output != "" {
				fmt.Fprint(writer, item.Output)
			}
		} else {
			fmt.Fprintf(writer, "Perimeter Guard: QUARANTINED (%s, mode: %s)\nTarget: %s/%s\n", item.Status, item.Mode, item.MissionRef, item.ObjectiveRef)
			if len(item.PreservedPaths) > 0 {
				fmt.Fprintf(writer, "Preserved Valid Paths: %s\n", strings.Join(item.PreservedPaths, ", "))
			}
			if len(item.EscapedPaths) > 0 {
				fmt.Fprintf(writer, "Escaped Paths (Purged & Rolled Back): %s\n", strings.Join(item.EscapedPaths, ", "))
			}
			fmt.Fprintf(writer, "Allowed Write Paths: %s\n", strings.Join(item.AllowedPaths, ", "))
			if item.FeedbackPrompt != "" {
				fmt.Fprintf(writer, "\nFeedback Turn for Session Continuation:\n%s\n", item.FeedbackPrompt)
			}
		}
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
// removeValueFlag extracts a `--name value` pair. Returns the value, whether it was
// present, and whether it was malformed — supplied twice or with no value after it.
func removeValueFlag(args *[]string, name string) (string, bool, bool) {
	value, found, bad := "", false, false
	out := (*args)[:0]
	for i := 0; i < len(*args); i++ {
		if (*args)[i] != name {
			out = append(out, (*args)[i])
			continue
		}
		if found || i+1 >= len(*args) || (*args)[i+1] == "" {
			bad = true
			if i+1 < len(*args) {
				i++
			}
			continue
		}
		value, found = (*args)[i+1], true
		i++
	}
	*args = out
	return value, found, bad
}

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
		fmt.Fprintf(r.Stderr, ": %s\n", problem)
		// The JSON envelope already carries `actual`; the human path dropped it,
		// so a decode failure named the field and hid the parser's line number.
		// A reader then hunts for a bad field when a mistyped character several
		// lines away is the real fault.
		if detail := strings.TrimSpace(strings.ReplaceAll(actual, "\n", "; ")); detail != "" && detail != problem {
			fmt.Fprintf(r.Stderr, "cause: %s\n", detail)
		}
		fmt.Fprintf(r.Stderr, "correction: %s\n", correction)
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

// short abbreviates a sha256 fingerprint for human output. The full value stays in
// --json and in the amendment log; a reader comparing two fingerprints by eye needs
// the first few characters, not sixty-four.
func short(fingerprint string) string {
	trimmed := strings.TrimPrefix(fingerprint, "sha256:")
	if len(trimmed) > 12 {
		return trimmed[:12] + "…"
	}
	return trimmed
}

type CommandHelpData struct {
	Command     string `json:"command"`
	Arguments   string `json:"arguments"`
	Schema      string `json:"schema"`
	Effect      Effect `json:"effect"`
	Description string `json:"description"`
	InputType   string `json:"input_type,omitempty"`
	OutputType  string `json:"output_type,omitempty"`
	Template    string `json:"template,omitempty"`
}

func (r Runner) globalHelp(jsonMode bool) int {
	now := r.Now
	if now == nil {
		now = time.Now
	}
	if jsonMode {
		output := envelope{
			SchemaVersion: "spectacular.command-catalog.v1",
			GeneratedAt:   now().UTC().Format(time.RFC3339Nano),
			Data: struct {
				ReleaseInspection CatalogEntry   `json:"release_inspection"`
				Commands          []CatalogEntry `json:"commands"`
			}{
				ReleaseInspection: VersionInspection,
				Commands:          Catalog(),
			},
		}
		_ = writeJSON(r.Stdout, output)
		return 0
	}
	fmt.Fprintln(r.Stdout, "spectacular — governed execution for agents and developers")
	fmt.Fprintln(r.Stdout)
	fmt.Fprintln(r.Stdout, "Usage:")
	fmt.Fprintf(r.Stdout, "  %s %s\n", VersionInspection.Command, VersionInspection.Arguments)
	for _, spec := range Registry {
		fmt.Fprintf(r.Stdout, "  spectacular %s %s\n", strings.Join(spec.Words, " "), spec.Arguments)
	}
	fmt.Fprintln(r.Stdout)
	fmt.Fprintln(r.Stdout, "Flags:")
	fmt.Fprintln(r.Stdout, "  -h, --help    Show command help, usage, and starter YAML templates")
	fmt.Fprintln(r.Stdout, "  --schema      Inspect machine-readable schema specification")
	fmt.Fprintln(r.Stdout, "  --json        Emit machine-readable JSON output envelope")
	return 0
}

func (r Runner) commandHelp(jsonMode bool, spec Spec) int {
	now := r.Now
	if now == nil {
		now = time.Now
	}
	data := CommandHelpData{
		Command:     "spectacular " + strings.Join(spec.Words, " "),
		Arguments:   spec.Arguments,
		Schema:      spec.JSONSchema,
		Effect:      spec.Effect,
		Description: spec.Description,
		InputType:   spec.InputType,
		OutputType:  spec.OutputType,
		Template:    spec.Template,
	}
	if jsonMode {
		output := envelope{
			SchemaVersion: "spectacular.command-help.v1",
			GeneratedAt:   now().UTC().Format(time.RFC3339Nano),
			Data:          data,
		}
		_ = writeJSON(r.Stdout, output)
		return 0
	}
	fmt.Fprintf(r.Stdout, "%s — %s (%s)\n", data.Command, spec.Effect, spec.JSONSchema)
	if spec.Description != "" {
		fmt.Fprintln(r.Stdout, spec.Description)
	}
	fmt.Fprintln(r.Stdout)
	fmt.Fprintf(r.Stdout, "Usage:\n  %s %s\n", data.Command, spec.Arguments)
	if spec.Template != "" {
		fmt.Fprintf(r.Stdout, "\nInput Template (%s YAML frontmatter):\n%s", spec.InputType, spec.Template)
	}
	return 0
}

type CommandSchemaData struct {
	Command    string `json:"command"`
	Schema     string `json:"schema"`
	Effect     Effect `json:"effect"`
	Arguments  string `json:"arguments"`
	InputType  string `json:"input_type,omitempty"`
	OutputType string `json:"output_type,omitempty"`
	Template   string `json:"template,omitempty"`
}

func (r Runner) commandSchema(jsonMode bool, spec Spec) int {
	now := r.Now
	if now == nil {
		now = time.Now
	}
	data := CommandSchemaData{
		Command:    "spectacular " + strings.Join(spec.Words, " "),
		Schema:     spec.JSONSchema,
		Effect:     spec.Effect,
		Arguments:  spec.Arguments,
		InputType:  spec.InputType,
		OutputType: spec.OutputType,
		Template:   spec.Template,
	}
	output := envelope{
		SchemaVersion: "spectacular.command-schema.v1",
		GeneratedAt:   now().UTC().Format(time.RFC3339Nano),
		Data:          data,
	}
	_ = writeJSON(r.Stdout, output)
	return 0
}
