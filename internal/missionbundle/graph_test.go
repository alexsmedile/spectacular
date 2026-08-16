package missionbundle

import (
	"strings"
	"testing"
)

func graphBundle(objectives ...Objective) *Bundle {
	return &Bundle{
		Ref: "M7", Title: "Render derived state", Status: "active",
		Run:        &Run{Ref: "R1", Status: "active"},
		Objectives: objectives,
	}
}

// Notation is part of the Contract: a view whose glyphs drift between commands
// is worse than no view.
func TestGlyphsAreStableAcrossReadiness(t *testing.T) {
	for readiness, want := range map[Readiness]string{
		ReadyDone:      "✓",
		ReadyActive:    "◐",
		ReadyStartable: "▶",
		ReadyBlocked:   "·",
	} {
		if got := Glyph(readiness); got != want {
			t.Fatalf("Glyph(%q) = %q, want %q", readiness, got, want)
		}
	}
}

// An edge in the drawing is an edge in the plan. A chain draws on one row, and
// a fork draws one row per path.
func TestObjectiveGraphDrawsRealEdges(t *testing.T) {
	chain := graphBundle(
		objective("O1", "implemented"),
		objective("O2", "pending", "O1"),
		objective("O3", "pending", "O2"),
	)
	rendered := chain.Graph(DefaultGraphWidth)
	if !strings.Contains(rendered, "O1✓") || !strings.Contains(rendered, "O3·") {
		t.Fatalf("a chain must draw every Objective:\n%s", rendered)
	}
	if lines := graphBody(rendered); len(lines) != 1 {
		t.Fatalf("a linear chain draws on one row, got %d:\n%s", len(lines), rendered)
	}

	fork := graphBundle(
		objective("O1", "implemented"),
		objective("O2", "pending", "O1"),
		objective("O3", "pending", "O1"),
	)
	rendered = fork.Graph(DefaultGraphWidth)
	lines := graphBody(rendered)
	if len(lines) != 2 {
		t.Fatalf("a fork draws one row per path, got %d:\n%s", len(lines), rendered)
	}
	if !strings.Contains(lines[0], "┬") {
		t.Fatalf("a node with several dependents forks with ┬:\n%s", rendered)
	}
	if !strings.Contains(lines[1], "└") {
		t.Fatalf("a branching row joins with └:\n%s", rendered)
	}
	// The shared prefix is blanked, so a branch reads as a branch rather than
	// as a repeated chain.
	if strings.Contains(lines[1], "O1") {
		t.Fatalf("the shared prefix must not repeat on a branch row:\n%s", rendered)
	}
}

// Two Objectives with no predecessors are two roots, and both must be drawn.
func TestObjectiveGraphDrawsEveryRoot(t *testing.T) {
	bundle := graphBundle(
		objective("O1", "pending"),
		objective("O2", "pending"),
		objective("O3", "pending", "O1"),
	)
	rendered := bundle.Graph(DefaultGraphWidth)
	for _, ref := range []string{"O1", "O2", "O3"} {
		if !strings.Contains(rendered, ref) {
			t.Fatalf("%s is missing from the graph:\n%s", ref, rendered)
		}
	}
	if len(graphBody(rendered)) != 2 {
		t.Fatalf("two roots draw two rows:\n%s", rendered)
	}
}

// Notation is chosen by shape, not by a flag: the graph is the default and level
// sets take over only when the graph would exceed the width.
func TestGraphFallsBackToLevelSetsOnlyWhenTooWide(t *testing.T) {
	bundle := graphBundle(
		objective("O1", "implemented"),
		objective("O2", "pending", "O1"),
		objective("O3", "pending", "O2"),
		objective("O4", "pending", "O3"),
	)
	wide := bundle.Graph(DefaultGraphWidth)
	if strings.Contains(wide, "L0") {
		t.Fatalf("a graph that fits must not fall back to level sets:\n%s", wide)
	}
	narrow := bundle.Graph(12)
	if !strings.Contains(narrow, "L0") || !strings.Contains(narrow, "∥ = independent") {
		t.Fatalf("a graph that overflows must fall back to level sets:\n%s", narrow)
	}
	// The same Mission, both notations: no Objective may be lost in either.
	for _, ref := range []string{"O1", "O2", "O3", "O4"} {
		if !strings.Contains(wide, ref) || !strings.Contains(narrow, ref) {
			t.Fatalf("%s missing from one notation", ref)
		}
	}
}

// Level sets name each Objective, so a reader is not forced to look up a ref,
// and mark which work is independent.
func TestLevelSetsNameObjectivesAndMarkParallelWork(t *testing.T) {
	bundle := graphBundle(
		Objective{Ref: "O1", Outcome: "Extract the shared derivation layer", Status: "implemented", Claims: []string{"c"}},
		Objective{Ref: "O2", Outcome: "Render the compact state line", Status: "pending", After: []string{"O1"}, Claims: []string{"c"}},
		Objective{Ref: "O3", Outcome: "Compute per-claim drift flags", Status: "pending", After: []string{"O1"}, Claims: []string{"c"}},
	)
	rendered := bundle.Graph(1)
	if !strings.Contains(rendered, "shared derivation layer") {
		t.Fatalf("level sets carry a short name per Objective:\n%s", rendered)
	}
	if !strings.Contains(rendered, "∥") {
		t.Fatalf("independent Objectives on one level are marked parallel:\n%s", rendered)
	}
	if !strings.Contains(rendered, "(after O1)") {
		t.Fatalf("a level names what it waits on:\n%s", rendered)
	}
}

// The short name is a shortening of the outcome the bundle already carries, so
// it cannot drift from the text it names.
func TestShortNameDropsLeadingVerbsAndFiller(t *testing.T) {
	for outcome, want := range map[string]string{
		"Extract the shared derivation layer over the typed bundle": "shared derivation layer",
		"Render the compact state line and NEXT line":               "compact state line",
		"Compute per-claim drift flags":                             "per-claim drift flags",
		"":                                                          "",
		"Timeline":                                                  "Timeline",
	} {
		if got := ShortName(outcome); got != want {
			t.Fatalf("ShortName(%q) = %q, want %q", outcome, got, want)
		}
	}
}

// A Mission with no Objectives says so rather than drawing an empty frame.
func TestGraphStatesWhenThereIsNothingToDraw(t *testing.T) {
	rendered := graphBundle().Graph(DefaultGraphWidth)
	if !strings.Contains(rendered, "no Objectives are declared") {
		t.Fatalf("an empty Mission must say so:\n%s", rendered)
	}
}

// Every view is labelled <ref> · <title>: the ref is what a reader types, the
// title is what they recognise.
func TestGraphLabelsTheMissionByRefAndTitle(t *testing.T) {
	bundle := graphBundle(objective("O1", "pending"))
	for _, width := range []int{DefaultGraphWidth, 1} {
		if !strings.Contains(bundle.Graph(width), "M7 · Render derived state") {
			t.Fatalf("width %d lost the Mission label", width)
		}
	}
}

// graphBody returns the drawn rows, without the heading or the legend.
func graphBody(rendered string) []string {
	var rows []string
	for _, line := range strings.Split(rendered, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.Contains(line, " · ") || strings.HasPrefix(trimmed, "∥ =") {
			continue
		}
		rows = append(rows, line)
	}
	return rows
}
