package missionbundle

import (
	"encoding/json"
	"testing"
)

// equivalencePlan is a Mission with enough shape for every projection surface to
// have something to say: two roots, a fork, and a join.
func equivalencePlan() (Plan, []byte) {
	plan, raw := stressPlan()
	plan.Title = "Representation equivalence fixture"
	plan.Objectives = []Objective{
		{Outcome: "Extract the shared derivation layer", Status: "active", Claims: []string{"atomic"}},
		{Outcome: "Normalize reference resolution", Status: "pending", Claims: []string{"atomic"}},
		{Outcome: "Render the compact state line", Status: "pending", After: []string{"O1"}, Claims: []string{"atomic"}},
		{Outcome: "Compute per-claim drift flags", Status: "pending", After: []string{"O1"}, Claims: []string{"atomic"}},
		{Outcome: "Render the multi-Mission timeline", Status: "pending", After: []string{"O3", "O4"}, Claims: []string{"atomic"}},
	}
	return plan, raw
}

// A promoted Objective lives in its own file and resolves through a different
// decode path than an inline one. Every projection must reach byte-identical
// output across both, or a reader's conclusion depends on where the record
// happens to be stored.
//
// This is the property most likely to rot silently as projection surfaces grow,
// because nothing else fails when it breaks.
func TestInlineAndPromotedRepresentationsRenderIdentically(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)
	plan, raw := equivalencePlan()
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatal(err)
	}

	inline := projectionSnapshot(t, root, started.Ref)

	// Promote a middle Objective and a leaf, so the bundle mixes inline and
	// promoted records rather than converting wholesale.
	for _, ref := range []string{"/O1", "/O5"} {
		service = openMissionService(t, root)
		if _, err := service.PromoteObjective(started.Ref + ref); err != nil {
			t.Fatalf("promote %s: %v", ref, err)
		}
	}

	promoted := projectionSnapshot(t, root, started.Ref)

	for surface, before := range inline {
		after, rendered := promoted[surface]
		if !rendered {
			t.Fatalf("%s is missing after promotion", surface)
		}
		if before != after {
			t.Fatalf("%s differs between representations:\n--- inline ---\n%s\n--- promoted ---\n%s", surface, before, after)
		}
	}
}

// projectionSnapshot renders every surface a reader can reach for one Mission.
// Adding a projection without adding it here would leave the new surface
// unguarded, so the map is deliberately exhaustive.
func projectionSnapshot(t *testing.T, root, ref string) map[string]string {
	t.Helper()
	service := openMissionService(t, root)
	bundle, err := service.Show(ref)
	if err != nil {
		t.Fatal(err)
	}
	check, err := service.Check(ref)
	if err != nil {
		t.Fatal(err)
	}

	// The bundle carries its own storage layout in Source and File, which are
	// expected to differ; the derived conclusions are what must not.
	state, err := json.Marshal(bundle.State)
	if err != nil {
		t.Fatal(err)
	}
	drift, err := json.Marshal(check.Drift)
	if err != nil {
		t.Fatal(err)
	}
	authority, err := json.Marshal(check.Authority)
	if err != nil {
		t.Fatal(err)
	}

	return map[string]string{
		"graph":      bundle.Graph(DefaultGraphWidth),
		"levelSets":  bundle.Graph(1),
		"stateLine":  string(state),
		"drift":      string(drift),
		"authority":  string(authority),
		"notices":    joinLines(check.Notices),
		"checkNames": joinLines(check.Checks),
	}
}

func joinLines(values []string) string {
	out := ""
	for _, value := range values {
		out += value + "\n"
	}
	return out
}
