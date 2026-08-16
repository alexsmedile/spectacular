package missionbundle

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func criterion(claim string) Criterion {
	return Criterion{Claim: claim, PassBoundary: claim + " boundary", ProofRequirement: claim + " proof"}
}

// reviewPointer builds a pointer with an attached document, which is how a
// resolved bundle carries review verdicts.
func reviewPointer(fingerprint string, verdicts map[string]string) ReviewPointer {
	document := &Review{}
	document.Reviewed.ActivationFingerprint = fingerprint
	for claim, verdict := range verdicts {
		document.Claims = append(document.Claims, ClaimVerdict{Claim: claim, Verdict: verdict})
	}
	return ReviewPointer{Ref: "RV1", Document: document}
}

func flagsFor(drift []ClaimDrift, claim string) []Flag {
	for _, item := range drift {
		if item.Claim == claim {
			return item.Flags
		}
	}
	return nil
}

func TestDriftFlagsAreNamedPerObservedCause(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		bundle *Bundle
		claim  string
		want   []Flag
	}{
		{
			name: "an unreviewed claim with no implemented Objective is unproven and unimplemented",
			bundle: &Bundle{
				Completion: []Criterion{criterion("alpha")},
				Objectives: []Objective{{Ref: "O1", Status: "pending", Claims: []string{"alpha"}}},
			},
			claim: "alpha", want: []Flag{FlagUnproven, FlagUnimplemented},
		},
		{
			name: "implementing the Objective clears unimplemented but not unproven",
			bundle: &Bundle{
				Completion: []Criterion{criterion("alpha")},
				Objectives: []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
			},
			claim: "alpha", want: []Flag{FlagUnproven},
		},
		{
			name: "a claim no Objective claims is unclaimed rather than unimplemented",
			bundle: &Bundle{
				Completion: []Criterion{criterion("orphan")},
				Objectives: []Objective{{Ref: "O1", Status: "pending", Claims: []string{"alpha"}}},
			},
			claim: "orphan", want: []Flag{FlagUnproven, FlagUnclaimed},
		},
		{
			name: "a passing review on the current fingerprint leaves no flags",
			bundle: &Bundle{
				Completion: []Criterion{criterion("alpha")},
				Objectives: []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
				Activation: &Activation{Fingerprint: "sha256:current"},
				Reviews:    []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "pass"})},
			},
			claim: "alpha", want: nil,
		},
		{
			name: "a non-pass verdict is failed rather than unproven",
			bundle: &Bundle{
				Completion: []Criterion{criterion("alpha")},
				Objectives: []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
				Activation: &Activation{Fingerprint: "sha256:current"},
				Reviews:    []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "fail"})},
			},
			claim: "alpha", want: []Flag{FlagFailed},
		},
		{
			name: "a review bound to a superseded fingerprint is stale",
			bundle: &Bundle{
				Completion: []Criterion{criterion("alpha")},
				Objectives: []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
				Activation: &Activation{Fingerprint: "sha256:current"},
				Reviews:    []ReviewPointer{reviewPointer("sha256:old", map[string]string{"alpha": "pass"})},
			},
			claim: "alpha", want: []Flag{FlagStaleReview},
		},
		{
			name: "consumed repair budget flags every claim, because a Run does not record which claim a repair served",
			bundle: &Bundle{
				Completion:   []Criterion{criterion("alpha")},
				Objectives:   []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
				Activation:   &Activation{Fingerprint: "sha256:current"},
				Reviews:      []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "pass"})},
				RepairBudget: 3,
				Run:          &Run{Ref: "R1", Repairs: 1},
			},
			claim: "alpha", want: []Flag{FlagRepaired},
		},
		{
			name: "reaching the budget adds exhaustion on top of repaired",
			bundle: &Bundle{
				Completion:   []Criterion{criterion("alpha")},
				Objectives:   []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
				Activation:   &Activation{Fingerprint: "sha256:current"},
				Reviews:      []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "pass"})},
				RepairBudget: 2,
				Run:          &Run{Ref: "R1", Repairs: 2},
			},
			claim: "alpha", want: []Flag{FlagRepaired, FlagBudgetExhausted},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if got := flagsFor(testCase.bundle.Drift(), testCase.claim); !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("flags %v, want %v", got, testCase.want)
			}
		})
	}
}

// The default audit target is the most-flagged claim, and ties break on plan
// order so the target does not move between runs for no reason.
func TestAuditTargetPicksTheMostFlaggedClaimAndBreaksTiesStably(t *testing.T) {
	bundle := &Bundle{
		Completion: []Criterion{criterion("alpha"), criterion("beta"), criterion("gamma")},
		Objectives: []Objective{
			{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}},
			{Ref: "O2", Status: "pending", Claims: []string{"beta"}},
			{Ref: "O3", Status: "implemented", Claims: []string{"gamma"}},
		},
		Activation: &Activation{Fingerprint: "sha256:current"},
		Reviews: []ReviewPointer{reviewPointer("sha256:current", map[string]string{
			"alpha": "pass",
			"gamma": "pass",
		})},
	}
	// beta is unproven and unimplemented; alpha and gamma are clean.
	target, ok := bundle.AuditTarget()
	if !ok {
		t.Fatal("a Mission with a flagged claim must yield an audit target")
	}
	if target.Claim != "beta" {
		t.Fatalf("audit target %q, want the most-flagged claim beta", target.Claim)
	}
	if len(target.Flags) != 2 {
		t.Fatalf("target carries %v, want the two flags that selected it", target.Flags)
	}

	// Every claim clean means no target: the caller asks rather than inventing one.
	clean := &Bundle{
		Completion: []Criterion{criterion("alpha")},
		Objectives: []Objective{{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}}},
		Activation: &Activation{Fingerprint: "sha256:current"},
		Reviews:    []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "pass"})},
	}
	if _, ok := clean.AuditTarget(); ok {
		t.Fatal("a Mission with no flagged claim must not name an audit target")
	}
}

// Ordering is by flag count, and equally-flagged claims keep plan order.
func TestDriftOrdersMostFlaggedFirstAndPreservesPlanOrderOnTies(t *testing.T) {
	// Declared deliberately in ascending flag order, so a missing sort leaves
	// the list exactly reversed rather than coincidentally correct.
	bundle := &Bundle{
		Completion: []Criterion{criterion("alpha"), criterion("gamma"), criterion("beta")},
		Objectives: []Objective{
			{Ref: "O1", Status: "implemented", Claims: []string{"alpha"}},
			{Ref: "O3", Status: "implemented", Claims: []string{"gamma"}},
		},
		Activation: &Activation{Fingerprint: "sha256:current"},
		Reviews:    []ReviewPointer{reviewPointer("sha256:current", map[string]string{"alpha": "pass"})},
	}
	// beta: unproven + unclaimed. gamma: unproven. alpha: clean.
	drift := bundle.Drift()
	if drift[0].Claim != "beta" {
		t.Fatalf("most-flagged first: got %q", drift[0].Claim)
	}
	if drift[1].Claim != "gamma" || drift[2].Claim != "alpha" {
		t.Fatalf("order %q, %q, %q; want beta, gamma, alpha", drift[0].Claim, drift[1].Claim, drift[2].Claim)
	}
	if !reflect.DeepEqual(drift[0].Objectives, []string(nil)) {
		t.Fatalf("beta is claimed by no Objective, got %v", drift[0].Objectives)
	}
	if !reflect.DeepEqual(drift[2].Objectives, []string{"O1"}) {
		t.Fatalf("alpha is claimed by O1, got %v", drift[2].Objectives)
	}

	// Equally-flagged claims keep plan order, so the default audit target does
	// not move between runs on an arbitrary tie.
	tied := &Bundle{
		Completion: []Criterion{criterion("delta"), criterion("epsilon"), criterion("zeta")},
	}
	order := tied.Drift()
	if order[0].Claim != "delta" || order[1].Claim != "epsilon" || order[2].Claim != "zeta" {
		t.Fatalf("tied claims must keep plan order, got %q, %q, %q", order[0].Claim, order[1].Claim, order[2].Claim)
	}
}

// Drift is computed from the bundle on read and reaches mission check without a
// new command. Checked against the Missions in this repository.
func TestDriftAgreesWithSelfHostedMissions(t *testing.T) {
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
			drift := bundle.Drift()
			if len(drift) != len(bundle.Completion) {
				t.Fatalf("drift covers %d claims for %d frozen criteria", len(drift), len(bundle.Completion))
			}
			for index := 1; index < len(drift); index++ {
				if len(drift[index-1].Flags) < len(drift[index].Flags) {
					t.Fatalf("drift is not ordered most-flagged first at %d", index)
				}
			}
			// A Mission that consumed repair budget must say so on every claim,
			// which is the signal that was computed and never read before.
			if bundle.Run != nil && bundle.Run.Repairs > 0 {
				for _, item := range drift {
					if !containsFlag(item.Flags, FlagRepaired) {
						t.Fatalf("%s consumed %d repairs but %q is not flagged repaired", ref, bundle.Run.Repairs, item.Claim)
					}
				}
			}
		})
	}
}

func containsFlag(flags []Flag, want Flag) bool {
	for _, flag := range flags {
		if flag == want {
			return true
		}
	}
	return false
}

// The rename from human_ref to ref stopped halfway. A reader must accept both
// spellings while a checker reports the drift, because rewriting a frozen
// record's frontmatter to finish a rename changes fingerprints for a cosmetic
// reason. M8 validates Mission order against refs and needs one spelling.
func TestRefSpellingDriftIsReportedNotRefused(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	service := Service{Workspace: ws}
	for ref, wantDrift := range map[string]bool{
		"M2": true, "M3": true, "M4": true,
		"M5": false, "M6": false, "M7": false,
	} {
		t.Run(ref, func(t *testing.T) {
			check, checkErr := service.Check(ref)
			if checkErr != nil {
				t.Fatalf("legacy spelling must not refuse: %v", checkErr)
			}
			if !check.Valid {
				t.Fatal("a legacy ref spelling is reported, never a validation failure")
			}
			drifted := false
			for _, notice := range check.Notices {
				if strings.Contains(notice, "ref-spelling-drift") {
					drifted = true
				}
			}
			if drifted != wantDrift {
				t.Fatalf("ref-spelling-drift reported=%t, want %t", drifted, wantDrift)
			}
			// Every Mission resolves to its ref through the one decoder,
			// whichever spelling the record uses.
			bundle, loadErr := Load(ws, ref)
			if loadErr != nil {
				t.Fatal(loadErr)
			}
			if bundle.Ref != ref {
				t.Fatalf("resolved ref %q, want %q regardless of spelling", bundle.Ref, ref)
			}
		})
	}
}
