package missionbundle

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

// objective is a terse constructor for the table below. Derive reads only
// exported fields, so a Bundle can be built in-code without a workspace.
func objective(ref, status string, after ...string) Objective {
	return Objective{Ref: ref, Outcome: ref + " outcome", Status: status, After: after, Claims: []string{"c"}}
}

func TestDeriveReadinessNextActionAndHolder(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		bundle    *Bundle
		startable int
		blocked   int
		done      int
		next      string
		holder    string
	}{
		{
			name: "two parallel roots are both startable",
			bundle: &Bundle{
				Ref: "M7", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active"},
				Objectives: []Objective{
					objective("O1", "pending"),
					objective("O2", "pending"),
					objective("O3", "pending", "O1"),
					objective("O4", "pending", "O2"),
				},
			},
			startable: 2, blocked: 2, done: 0,
			next: "work any of O1, O2", holder: "operator",
		},
		{
			name: "a strict chain exposes exactly one startable Objective",
			bundle: &Bundle{
				Ref: "M6", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active"},
				Objectives: []Objective{
					objective("O1", "pending"),
					objective("O2", "pending", "O1"),
					objective("O3", "pending", "O2"),
				},
			},
			startable: 1, blocked: 2, done: 0,
			next: "work O1", holder: "operator",
		},
		{
			name: "implementing a predecessor releases its dependent",
			bundle: &Bundle{
				Ref: "M6", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active"},
				Objectives: []Objective{
					objective("O1", "implemented"),
					objective("O2", "pending", "O1"),
					objective("O3", "pending", "O2"),
				},
			},
			startable: 1, blocked: 1, done: 1,
			next: "work O2", holder: "operator",
		},
		{
			name: "the Run's current Objective is active, not merely startable",
			bundle: &Bundle{
				Ref: "M7", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active", CurrentObjective: "O1"},
				Objectives: []Objective{
					objective("O1", "pending"),
					objective("O2", "pending"),
				},
			},
			startable: 2, blocked: 0, done: 0,
			next: "work any of O1, O2", holder: "operator",
		},
		{
			name: "every Objective implemented asks for a review",
			bundle: &Bundle{
				Ref: "M6", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active"},
				Objectives: []Objective{
					objective("O1", "implemented"),
					objective("O2", "implemented", "O1"),
				},
			},
			startable: 0, blocked: 0, done: 2,
			next: "every Objective is implemented; record a review", holder: "operator",
		},
		{
			name: "a Mission with no Run is told to start one",
			bundle: &Bundle{
				Ref: "M8", Status: "active", RepairBudget: 3,
				Objectives: []Objective{objective("O1", "pending")},
			},
			startable: 1, blocked: 0, done: 0,
			next: "start a Run", holder: "operator",
		},
		{
			name: "exhausted repair budget outranks startable work",
			bundle: &Bundle{
				Ref: "M7", Status: "active", RepairBudget: 3,
				Run: &Run{Ref: "R1", Status: "active", Repairs: 3},
				Objectives: []Objective{
					objective("O1", "pending"),
					objective("O2", "pending"),
				},
			},
			startable: 2, blocked: 0, done: 0,
			next:   "repair budget is exhausted; the owner decides whether to continue",
			holder: "operator",
		},
		{
			name: "a completed Mission has no next action and no holder",
			bundle: &Bundle{
				Ref: "M6", Status: "completed", RepairBudget: 3,
				Run:        &Run{Ref: "R1", Status: "completed", Repairs: 1},
				Objectives: []Objective{objective("O1", "implemented")},
			},
			startable: 0, blocked: 0, done: 1,
			next: "nothing; the Mission is complete", holder: "no one",
		},
		{
			name: "a Mission awaiting review is held by the owner",
			bundle: &Bundle{
				Ref: "M6", Status: "awaiting-review", RepairBudget: 3,
				Run:        &Run{Ref: "R1", Status: "awaiting-review"},
				Objectives: []Objective{objective("O1", "implemented")},
			},
			startable: 0, blocked: 0, done: 1,
			next:   "record a review covering every frozen completion claim",
			holder: "owner",
		},
		{
			name:      "a Mission with no Objectives says so rather than inventing work",
			bundle:    &Bundle{Ref: "M9", Status: "active", RepairBudget: 3, Run: &Run{Ref: "R1", Status: "active"}},
			startable: 0, blocked: 0, done: 0,
			next: "no Objectives are declared", holder: "operator",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			state := testCase.bundle.Derive()
			if state.Startable != testCase.startable || state.Blocked != testCase.blocked || state.Done != testCase.done {
				t.Fatalf("counts: startable=%d blocked=%d done=%d, want %d/%d/%d",
					state.Startable, state.Blocked, state.Done, testCase.startable, testCase.blocked, testCase.done)
			}
			if state.Next != testCase.next {
				t.Fatalf("next: %q, want %q", state.Next, testCase.next)
			}
			if state.Holder != testCase.holder {
				t.Fatalf("holder: %q, want %q", state.Holder, testCase.holder)
			}
			if total := state.Startable + state.Blocked + state.Done; total != len(state.Objectives) {
				t.Fatalf("every Objective must fall in exactly one readiness bucket: %d buckets for %d Objectives", total, len(state.Objectives))
			}
		})
	}
}

// BlockedBy must name only the predecessors that are actually unfinished. Naming
// every declared predecessor would tell a reader to wait on work already done.
func TestDeriveBlockedByNamesOnlyUnfinishedPredecessors(t *testing.T) {
	bundle := &Bundle{
		Ref: "M7", Status: "active", RepairBudget: 3,
		Run: &Run{Ref: "R1", Status: "active"},
		Objectives: []Objective{
			objective("O1", "implemented"),
			objective("O2", "pending"),
			objective("O3", "pending", "O1", "O2"),
		},
	}
	state := bundle.Derive()
	byRef := map[string]ObjectiveState{}
	for _, item := range state.Objectives {
		byRef[item.Ref] = item
	}
	if got := byRef["O3"].BlockedBy; !reflect.DeepEqual(got, []string{"O2"}) {
		t.Fatalf("O3 blocked by %v, want only the unfinished predecessor [O2]", got)
	}
	if byRef["O3"].Readiness != ReadyBlocked {
		t.Fatalf("O3 readiness %q, want %q", byRef["O3"].Readiness, ReadyBlocked)
	}
	if byRef["O1"].BlockedBy != nil {
		t.Fatalf("an implemented Objective must name no blockers, got %v", byRef["O1"].BlockedBy)
	}
}

// Level orders the level-set view. A dependent must never sort above the
// predecessor it waits on, and a diamond must take the longest path.
func TestDeriveLevelsExceedEveryPredecessor(t *testing.T) {
	bundle := &Bundle{
		Ref: "M7", Status: "active",
		Objectives: []Objective{
			objective("O1", "pending"),
			objective("O2", "pending", "O1"),
			objective("O3", "pending", "O1"),
			objective("O4", "pending", "O2", "O3"),
			objective("O5", "pending", "O1", "O4"),
		},
	}
	state := bundle.Derive()
	levels := map[string]int{}
	for _, item := range state.Objectives {
		levels[item.Ref] = item.Level
	}
	for _, item := range bundle.Objectives {
		for _, predecessor := range item.After {
			if levels[item.Ref] <= levels[predecessor] {
				t.Fatalf("%s is level %d but waits on %s at level %d", item.Ref, levels[item.Ref], predecessor, levels[predecessor])
			}
		}
	}
	if levels["O5"] != 3 {
		t.Fatalf("O5 level %d, want 3: the longest path O1->O2->O4->O5 wins over the direct O1 edge", levels["O5"])
	}
}

// An unresolvable `after:` ref is a validator's refusal to give, not a
// projection's. Derive must still return a usable answer.
func TestDeriveToleratesDanglingDependencyWithoutPanicking(t *testing.T) {
	bundle := &Bundle{
		Ref: "M7", Status: "active",
		Run:        &Run{Ref: "R1", Status: "active"},
		Objectives: []Objective{objective("O1", "pending", "O9")},
	}
	state := bundle.Derive()
	if state.Blocked != 1 {
		t.Fatalf("an Objective waiting on an unknown ref is blocked, got %d blocked", state.Blocked)
	}
	if got := state.Objectives[0].BlockedBy; !reflect.DeepEqual(got, []string{"O9"}) {
		t.Fatalf("blocked by %v, want the unresolved ref [O9]", got)
	}
}

// An undeclared verb is refused rather than defaulted to requires-owner.
// Defaulting never permits what it shouldn't, but it answers a question the
// record does not answer and turns a typo into a confident wrong result.
func TestAuthorizeAnswersByLookupAndRefusesUndeclaredVerbs(t *testing.T) {
	bundle := &Bundle{
		Ref: "M7",
		Authority: Authority{
			Operator:      []string{"inspect", "run-checks"},
			RequiresOwner: []string{"push", "release"},
		},
	}
	for verb, want := range map[string]Decision{
		"inspect":    DecisionOperator,
		"run-checks": DecisionOperator,
		"push":       DecisionOwner,
		"release":    DecisionOwner,
		// Supported by the vocabulary but not declared by this Mission.
		"commit-local": DecisionUndeclared,
		"merge":        DecisionUndeclared,
		// Not a verb at all.
		"deploy": DecisionUndeclared,
		"":       DecisionUndeclared,
	} {
		if got := bundle.Authorize(verb).Decision; got != want {
			t.Fatalf("Authorize(%q) = %q, want %q", verb, got, want)
		}
	}
}

// The table renders from the same vocabularies the validator enforces, so a
// reader is never shown a verb the validator would reject, or vice versa.
func TestAuthorityTableCoversTheValidatedVocabularyExactly(t *testing.T) {
	bundle := &Bundle{
		Ref: "M7",
		Authority: Authority{
			Operator:      SupportedOperatorVerbs,
			RequiresOwner: SupportedOwnerVerbs,
		},
	}
	table := bundle.AuthorityTable()
	if len(table) != len(SupportedOperatorVerbs)+len(SupportedOwnerVerbs) {
		t.Fatalf("table has %d entries for %d supported verbs", len(table), len(SupportedOperatorVerbs)+len(SupportedOwnerVerbs))
	}
	for _, answer := range table {
		if answer.Decision == DecisionUndeclared {
			t.Fatalf("%q is undeclared although the Mission declares the full vocabulary", answer.Verb)
		}
	}
	// A Mission declaring nothing must answer undeclared for every verb rather
	// than inventing a permission.
	empty := (&Bundle{Ref: "M7"}).AuthorityTable()
	for _, answer := range empty {
		if answer.Decision != DecisionUndeclared {
			t.Fatalf("%q = %q on a Mission with no authority block, want undeclared", answer.Verb, answer.Decision)
		}
	}
}

// Check must carry the table the authority-vocabulary validator resolves, so
// the answer arrives without a new command or noun.
func TestCheckCarriesTheAuthorityTableWithoutANewCommand(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	check, err := Service{Workspace: ws}.Check("M7")
	if err != nil {
		t.Fatal(err)
	}
	if len(check.Authority) == 0 {
		t.Fatal("mission check must render the authority decision table it already computes")
	}
	declared := 0
	for _, answer := range check.Authority {
		if answer.Decision != DecisionUndeclared {
			declared++
		}
	}
	if declared == 0 {
		t.Fatal("M7 declares an authority block, so the table cannot be entirely undeclared")
	}
}

// Derived state travels in JSON so an agent reading --json reaches the same
// conclusion a human reads from the rendered output, and it must never reach
// the canonical file. Mutations set canonical fields individually, so this
// guards the property rather than the mechanism.
func TestDerivedStateTravelsInJSONButNeverToTheCanonicalFile(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Service{Workspace: ws}.Show("M7")
	if err != nil {
		t.Fatal(err)
	}
	if bundle.State == nil {
		t.Fatal("Show must derive state so a JSON reader sees readiness")
	}
	if bundle.State.Next == "" || len(bundle.State.Objectives) != len(bundle.Objectives) {
		t.Fatalf("derived state must cover every Objective and name a next action: next=%q %d of %d",
			bundle.State.Next, len(bundle.State.Objectives), len(bundle.Objectives))
	}
	encoded, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(encoded, []byte(`"readiness"`)) {
		t.Fatal("JSON output must carry per-Objective readiness")
	}
	source, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(bundle.Path)))
	if err != nil {
		t.Fatal(err)
	}
	for _, derived := range []string{"state:", "readiness:", "blocked_by:", "startable:", "holder:"} {
		if bytes.Contains(source, []byte("\n"+derived)) {
			t.Fatalf("derived field %q reached the canonical file", derived)
		}
	}
}

// The lifecycle transitions in nextAction and holderFor were inferred from
// status strings rather than a declared state model. This checks them against
// the Missions in this repository.
func TestDeriveAgreesWithSelfHostedMissions(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"M5", "M6", "M7"} {
		t.Run(ref, func(t *testing.T) {
			bundle, loadErr := Load(ws, ref)
			if loadErr != nil {
				t.Fatal(loadErr)
			}
			state := bundle.Derive()
			if state.Ref != ref {
				t.Fatalf("derived ref %q, want %q", state.Ref, ref)
			}
			if state.Next == "" || state.Holder == "" {
				t.Fatalf("every Mission must derive a next action and a holder, got next=%q holder=%q", state.Next, state.Holder)
			}
			if total := state.Startable + state.Blocked + state.Done; total != len(state.Objectives) {
				t.Fatalf("readiness buckets total %d for %d Objectives", total, len(state.Objectives))
			}
			if bundle.Status == "completed" {
				if state.Holder != "no one" || state.Done != len(state.Objectives) {
					t.Fatalf("a completed Mission must hold no next action and have every Objective done: holder=%q done=%d/%d",
						state.Holder, state.Done, len(state.Objectives))
				}
			}
			for _, item := range state.Objectives {
				if item.Readiness == ReadyBlocked && len(item.BlockedBy) == 0 {
					t.Fatalf("%s is blocked but names no blocker", item.Ref)
				}
				if item.Readiness != ReadyBlocked && len(item.BlockedBy) > 0 {
					t.Fatalf("%s is %q but names blockers %v", item.Ref, item.Readiness, item.BlockedBy)
				}
			}
		})
	}
}
