package missionbundle

import (
	"reflect"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestReadViewUsesCanonicalDecoder(t *testing.T) {
	ws, err := discovery.Open("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"M3", "M5", "M16", "M18"} {
		t.Run(ref, func(t *testing.T) {
			bundle, err := Load(ws, ref)
			if err != nil {
				t.Fatal(err)
			}
			view, err := ReadView(ws, ref)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(view, projectView(bundle)) {
				t.Fatalf("read view differs from canonical decoder for %s", ref)
			}
		})
	}
	if _, err := ReadView(ws, "M999999"); err == nil {
		t.Fatal("missing Mission accepted")
	}
}

func TestReadViewCannotMutateCanonicalBundle(t *testing.T) {
	b := &Bundle{
		Scope:        Scope{Mechanical: []string{"src/"}},
		Authority:    Authority{Operator: []string{"inspect"}, RequiresOwner: []string{"release"}},
		Stops:        []string{"scope-drift"},
		Completion:   []Criterion{{Claim: "ship"}},
		Objectives:   []Objective{{Ref: "O1", Claims: []string{"ship"}, Sources: []string{"D1"}}},
		ResolvesGaps: []ResolvedGap{{Gap: "G1", Resolution: "resolved"}},
		Handoffs:     []HandoffPointer{{Document: &Handoff{Writes: []string{"src/"}}}, {}},
	}
	before := projectView(b)
	view := projectView(b)
	view.WritesPaths[0] = "escaped/"
	view.AllowedActions[0] = "release"
	view.RequiresOwner[0] = "inspect"
	view.Stops[0] = "none"
	view.Completion[0].Claim = "different"
	view.Objectives[0].Claims[0] = "different"
	view.Objectives[0].Sources[0] = "D2"
	view.ResolvesGaps[0].Resolution = "different"
	view.Handoffs[0].Writes[0] = "escaped/"
	if !reflect.DeepEqual(before, projectView(b)) {
		t.Fatal("read model aliases authoritative bundle memory")
	}
}
