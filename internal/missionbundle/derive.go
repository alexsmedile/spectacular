package missionbundle

import (
	"sort"
	"strconv"
	"strings"
)

// Readiness is what a reader wants to know about an Objective that has not
// finished: whether it can be picked up now, or is waiting on something.
type Readiness string

const (
	// ReadyStartable means every Objective this one declares in `after:` has
	// been implemented, so work can begin.
	ReadyStartable Readiness = "startable"
	// ReadyBlocked means at least one declared predecessor is unfinished.
	ReadyBlocked Readiness = "blocked"
	// ReadyActive means the Run currently names this Objective.
	ReadyActive Readiness = "active"
	// ReadyDone means the Objective is implemented.
	ReadyDone Readiness = "done"
)

// ObjectiveState is the derived view of one Objective. Nothing here is stored;
// every field is computed from the bundle on read.
type ObjectiveState struct {
	Ref       string    `json:"ref"`
	Outcome   string    `json:"outcome,omitempty"`
	Status    string    `json:"status,omitempty"`
	Readiness Readiness `json:"readiness"`
	// BlockedBy lists the declared predecessors that are not yet implemented.
	// It is empty unless Readiness is ReadyBlocked.
	BlockedBy      []string `json:"blocked_by,omitempty"`
	After          []string `json:"after,omitempty"`
	AfterInterface []string `json:"after_interface,omitempty"`
	// Level is the longest dependency distance from a root Objective. Roots are
	// level 0. It orders the level-set view and has no meaning beyond ordering.
	Level int `json:"level"`
}

// State is the derived answer to "where is this Mission and what happens next".
// It is computed from bundle fields alone and is never written back.
type State struct {
	Ref          string           `json:"ref"`
	Title        string           `json:"title"`
	Status       string           `json:"status"`
	Run          string           `json:"run,omitempty"`
	RunStatus    string           `json:"run_status,omitempty"`
	Repairs      int              `json:"repairs"`
	Budget       int              `json:"repair_budget"`
	Objectives   []ObjectiveState `json:"objectives,omitempty"`
	Fallbacks    []Fallback       `json:"fallbacks,omitempty"`
	AfterMission []string         `json:"after_mission,omitempty"`
	// Startable, Blocked, and Done count Objectives by derived readiness.
	Startable int `json:"startable"`
	Blocked   int `json:"blocked"`
	Done      int `json:"done"`
	// Next states the next action and Holder states who may take it. A reader
	// should be able to act on these two fields without opening the file.
	Next   string `json:"next"`
	Holder string `json:"holder"`
}

// Derive computes the Mission state. It reads the bundle and writes nothing.
//
// Every field traces to a bundle field: readiness comes from Objective status
// and `after:`, repairs from the live Run, and the next action from the
// combination of Mission status, Run status, and readiness counts.
func (b *Bundle) Derive() State {
	state := State{
		Ref:          b.Ref,
		Title:        b.Title,
		Status:       b.Status,
		Budget:       b.RepairBudget,
		Holder:       holderFor(b),
		Startable:    0,
		Fallbacks:    b.Fallbacks,
		AfterMission: b.AfterMission,
	}
	if b.Run != nil {
		state.Run = b.Run.Ref
		state.RunStatus = b.Run.Status
		state.Repairs = b.Run.Repairs
	}

	implemented := make(map[string]bool, len(b.Objectives))
	for _, objective := range b.Objectives {
		if objective.Status == "implemented" {
			implemented[objective.Ref] = true
		}
	}
	current := ""
	if b.Run != nil {
		current = b.Run.CurrentObjective
	}

	levels := objectiveLevels(b.Objectives)
	for _, objective := range b.Objectives {
		derived := ObjectiveState{
			Ref:            objective.Ref,
			Outcome:        objective.Outcome,
			Status:         objective.Status,
			After:          objective.After,
			AfterInterface: objective.AfterInterface,
			Level:          levels[objective.Ref],
		}
		switch {
		case objective.Status == "implemented":
			derived.Readiness = ReadyDone
			state.Done++
		default:
			for _, predecessor := range objective.After {
				if !implemented[predecessor] {
					derived.BlockedBy = append(derived.BlockedBy, predecessor)
				}
			}
			switch {
			case len(derived.BlockedBy) > 0:
				derived.Readiness = ReadyBlocked
				state.Blocked++
			case objective.Ref == current && current != "":
				derived.Readiness = ReadyActive
				state.Startable++
			default:
				derived.Readiness = ReadyStartable
				state.Startable++
			}
		}
		state.Objectives = append(state.Objectives, derived)
	}

	state.Next = nextAction(b, state)
	return state
}

// Decision is what the authority block answers for one verb.
type Decision string

const (
	// DecisionOperator means the Mission declares the verb for the operator.
	DecisionOperator Decision = "operator"
	// DecisionOwner means the verb is declared but gated on the owner.
	DecisionOwner Decision = "requires-owner"
	// DecisionUndeclared means the Mission does not declare the verb. It is a
	// refusal rather than an implicit owner gate: defaulting would answer a
	// question the record does not answer and turn a typo into a confident
	// wrong result.
	DecisionUndeclared Decision = "undeclared"
)

// AuthorityAnswer is the decision for one verb, with the vocabulary that
// produced it available for a refusal message.
type AuthorityAnswer struct {
	Verb     string   `json:"verb"`
	Decision Decision `json:"decision"`
}

// Authority answers a single authority question by lookup against the declared
// block, rather than asking the reader to recall two arrays.
func (b *Bundle) Authorize(verb string) AuthorityAnswer {
	for _, declared := range b.Authority.Operator {
		if declared == verb {
			return AuthorityAnswer{Verb: verb, Decision: DecisionOperator}
		}
	}
	for _, declared := range b.Authority.RequiresOwner {
		if declared == verb {
			return AuthorityAnswer{Verb: verb, Decision: DecisionOwner}
		}
	}
	return AuthorityAnswer{Verb: verb, Decision: DecisionUndeclared}
}

// AuthorityTable answers every verb in the supported vocabularies at once. It is
// what `mission check` renders: the decision for each verb the system knows
// about, including the ones this Mission declines to declare.
func (b *Bundle) AuthorityTable() []AuthorityAnswer {
	known := append(append([]string{}, SupportedOperatorVerbs...), SupportedOwnerVerbs...)
	answers := make([]AuthorityAnswer, 0, len(known))
	for _, verb := range known {
		answers = append(answers, b.Authorize(verb))
	}
	return answers
}

// nextAction states the single next thing to do, in the order the lifecycle
// forces. Earlier cases are genuinely blocking for later ones.
func nextAction(b *Bundle, state State) string {
	switch {
	case b.Status == "completed":
		return "nothing; the Mission is complete"
	case b.Status == "awaiting-review":
		return "record a review covering every frozen completion claim"
	case b.Run == nil:
		return "start a Run"
	case b.Run.Status == "awaiting-review":
		return "record a review covering every frozen completion claim"
	case state.Budget > 0 && state.Repairs >= state.Budget:
		return "repair budget is exhausted; the owner decides whether to continue"
	case state.Done == len(b.Objectives) && len(b.Objectives) > 0:
		return "every Objective is implemented; record a review"
	case state.Startable == 1:
		return "work " + firstWithReadiness(state, ReadyStartable, ReadyActive)
	case state.Startable > 1:
		return "work any of " + strings.Join(refsWithReadiness(state, ReadyStartable, ReadyActive), ", ")
	case state.Blocked > 0:
		return "nothing is startable; " + strconv.Itoa(state.Blocked) + " Objectives are blocked"
	default:
		return "no Objectives are declared"
	}
}

// holderFor names who may take the next action. Activation and completion are
// owner-gated; ordinary Objective work is not.
func holderFor(b *Bundle) string {
	switch b.Status {
	case "completed":
		return "no one"
	case "awaiting-review":
		return "owner"
	case "planned", "pending":
		return "owner"
	default:
		return "operator"
	}
}

func refsWithReadiness(state State, wanted ...Readiness) []string {
	var refs []string
	for _, objective := range state.Objectives {
		for _, want := range wanted {
			if objective.Readiness == want {
				refs = append(refs, objective.Ref)
				break
			}
		}
	}
	return refs
}

func firstWithReadiness(state State, wanted ...Readiness) string {
	refs := refsWithReadiness(state, wanted...)
	if len(refs) == 0 {
		return ""
	}
	return refs[0]
}

// objectiveLevels computes the longest path from a root for every Objective.
// The graph is already validated acyclic by objective-dependency-dag, so this
// terminates; an unresolvable ref is treated as level 0 rather than failing,
// because reference integrity is a validator's answer to give, not a
// projection's.
func objectiveLevels(objectives []Objective) map[string]int {
	byRef := make(map[string]Objective, len(objectives))
	for _, objective := range objectives {
		byRef[objective.Ref] = objective
	}
	levels := make(map[string]int, len(objectives))
	var resolve func(ref string, seen map[string]bool) int
	resolve = func(ref string, seen map[string]bool) int {
		if level, done := levels[ref]; done {
			return level
		}
		objective, known := byRef[ref]
		if !known || seen[ref] {
			return 0
		}
		seen[ref] = true
		level := 0
		for _, predecessor := range append(append([]string(nil), objective.After...), objective.AfterInterface...) {
			if candidate := resolve(predecessor, seen) + 1; candidate > level {
				level = candidate
			}
		}
		delete(seen, ref)
		levels[ref] = level
		return level
	}
	refs := make([]string, 0, len(objectives))
	for _, objective := range objectives {
		refs = append(refs, objective.Ref)
	}
	sort.Strings(refs)
	for _, ref := range refs {
		resolve(ref, map[string]bool{})
	}
	return levels
}
