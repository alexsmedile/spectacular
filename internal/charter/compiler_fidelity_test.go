package charter

import (
	"reflect"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

// This is a deterministic M14-style fidelity fixture: a charter must preserve
// the exact frozen claims and execution perimeter of the selected Objective.
func TestCompilePreservesFrozenMissionTruth(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := missionbundle.Load(ws, "M18")
	if err != nil {
		t.Fatal(err)
	}
	charter, err := Compile(ws, "M18", "O1", nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(charter.Layer3.WritesPaths, bundle.Scope.Mechanical) ||
		!reflect.DeepEqual(charter.Layer3.AllowedActions, bundle.Authority.Operator) ||
		!reflect.DeepEqual(charter.Layer3.RequiresOwner, bundle.Authority.RequiresOwner) ||
		!reflect.DeepEqual(charter.Layer3.Stops, bundle.Stops) {
		t.Fatal("charter execution perimeter drifted from the frozen Mission")
	}
	objectiveClaims := map[string]bool{}
	for _, objective := range bundle.Objectives {
		if objective.Ref != "O1" {
			continue
		}
		for _, claim := range objective.Claims {
			objectiveClaims[claim] = true
		}
	}
	for _, criterion := range bundle.Completion {
		if !objectiveClaims[criterion.Claim] {
			continue
		}
		found := false
		for _, claim := range charter.Layer1.Claims {
			if claim.Claim == criterion.Claim && claim.PassBoundary == criterion.PassBoundary && claim.ProofRequirement == criterion.ProofRequirement {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("charter dropped or rewrote frozen claim %q", criterion.Claim)
		}
	}
}
