package missionbundle

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultGraphWidth is the terminal width assumed when none is supplied. The
// Objective graph falls back to level sets beyond it, because a wrapped ASCII
// graph is worse than no graph.
const DefaultGraphWidth = 100

// Glyph returns the one-character status marker shared by every view. Colour is
// never load-bearing: the glyph survives monochrome terminals and copy-paste
// into plain text.
func Glyph(readiness Readiness) string {
	switch readiness {
	case ReadyDone:
		return "✓"
	case ReadyActive:
		return "◐"
	case ReadyStartable:
		return "▶"
	default:
		return "·"
	}
}

// Graph renders the Objective structure of a Mission, choosing notation by
// shape rather than by a flag: the graph is the default, and level sets take
// over when the graph would exceed width.
func (b *Bundle) Graph(width int) string {
	if width <= 0 {
		width = DefaultGraphWidth
	}
	state := b.Derive()
	if len(state.Objectives) == 0 {
		return b.graphHeading() + "\n  no Objectives are declared\n"
	}
	graph := b.objectiveGraph(state)
	if graphWidth(graph) <= width {
		return graph
	}
	return b.levelSets(state)
}

func (b *Bundle) graphHeading() string {
	return b.Ref + " · " + b.Title
}

// objectiveGraph draws each dependency chain on its own row, so an edge in the
// drawing is an edge in the plan. A node with several dependents branches with
// `┬` and `└`; a node with none stands alone.
func (b *Bundle) objectiveGraph(state State) string {
	nodes := map[string]ObjectiveState{}
	order := make([]string, 0, len(state.Objectives))
	for _, objective := range state.Objectives {
		nodes[objective.Ref] = objective
		order = append(order, objective.Ref)
	}
	successors := map[string][]string{}
	hasPredecessor := map[string]bool{}
	for _, objective := range b.Objectives {
		for _, predecessor := range objective.After {
			if _, known := nodes[predecessor]; !known {
				continue
			}
			successors[predecessor] = append(successors[predecessor], objective.Ref)
			hasPredecessor[objective.Ref] = true
		}
	}

	// Each row is one root-to-leaf path, so every edge appears exactly once and
	// a reader can follow a chain left to right without cross-referencing.
	var rows [][]string
	var walk func(ref string, path []string)
	walk = func(ref string, path []string) {
		path = append(path, ref)
		next := successors[ref]
		if len(next) == 0 {
			rows = append(rows, append([]string(nil), path...))
			return
		}
		sort.Strings(next)
		for _, successor := range next {
			walk(successor, path)
		}
	}
	for _, ref := range order {
		if !hasPredecessor[ref] {
			walk(ref, nil)
		}
	}

	cell := func(ref string) string { return ref + Glyph(nodes[ref].Readiness) }
	lines := make([]string, 0, len(rows))
	var previous []string
	for _, row := range rows {
		// Blank the prefix a row shares with the one above, so a branch reads as
		// a branch rather than as a repeated chain.
		shared := 0
		for shared < len(previous) && shared < len(row) && previous[shared] == row[shared] {
			shared++
		}
		var line strings.Builder
		line.WriteString("  ")
		for index, ref := range row {
			if index < shared {
				line.WriteString(strings.Repeat(" ", len([]rune(cell(ref)))))
			} else {
				line.WriteString(cell(ref))
			}
			if index == len(row)-1 {
				break
			}
			switch {
			// The connector after the last shared node is where this row
			// branches away from the one above.
			case index == shared-1:
				line.WriteString(" └─")
			case index < shared:
				line.WriteString("   ")
			// A node with several dependents forks here.
			case len(successors[ref]) > 1:
				line.WriteString(" ┬─")
			default:
				line.WriteString(" ──")
			}
		}
		lines = append(lines, strings.TrimRight(line.String(), " "))
		previous = row
	}
	return b.graphHeading() + "\n\n" + strings.Join(lines, "\n") + "\n"
}

// levelSets lists one line per dependency level, naming each Objective so a
// reader is not forced to look up what a ref refers to.
func (b *Bundle) levelSets(state State) string {
	byLevel := map[int][]ObjectiveState{}
	maximum := 0
	for _, objective := range state.Objectives {
		byLevel[objective.Level] = append(byLevel[objective.Level], objective)
		if objective.Level > maximum {
			maximum = objective.Level
		}
	}
	var lines []string
	for level := 0; level <= maximum; level++ {
		group := byLevel[level]
		if len(group) == 0 {
			continue
		}
		cells := make([]string, 0, len(group))
		for _, objective := range group {
			cells = append(cells, fmt.Sprintf("%s%s %s", objective.Ref, Glyph(objective.Readiness), ShortName(objective.Outcome)))
		}
		line := fmt.Sprintf("  L%d  %s", level, strings.Join(cells, "  ∥  "))
		if waits := b.declaredPredecessors(group); waits != "" {
			line += "    (after " + waits + ")"
		}
		lines = append(lines, line)
	}
	return b.graphHeading() + "\n\n" + strings.Join(lines, "\n") + "\n\n  ∥ = independent, can run in parallel\n"
}

// declaredPredecessors names the distinct Objectives a level depends on, so a
// reader sees why the level sits where it does. It reports the declared
// dependency rather than the unmet one: the level structure is a property of
// the plan, and does not change as Objectives are implemented.
func (b *Bundle) declaredPredecessors(group []ObjectiveState) string {
	inGroup := make(map[string]bool, len(group))
	for _, objective := range group {
		inGroup[objective.Ref] = true
	}
	seen := map[string]bool{}
	var refs []string
	for _, objective := range b.Objectives {
		if !inGroup[objective.Ref] {
			continue
		}
		for _, predecessor := range objective.After {
			if !seen[predecessor] {
				seen[predecessor] = true
				refs = append(refs, predecessor)
			}
		}
	}
	sort.Strings(refs)
	return strings.Join(refs, ", ")
}

// ShortName shortens an Objective outcome to a few words for the level-set view.
// It is a shortening of text the bundle already carries, never a second stored
// field, so it cannot drift from the outcome it names.
func ShortName(outcome string) string {
	words := strings.Fields(outcome)
	if len(words) == 0 {
		return ""
	}
	// Drop a leading verb phrase that carries no distinguishing information;
	// every Objective outcome starts with one.
	skip := map[string]bool{
		"extract": true, "render": true, "compute": true, "answer": true,
		"validate": true, "normalize": true, "prove": true, "add": true,
		"implement": true, "make": true, "provide": true, "surface": true,
	}
	if skip[strings.ToLower(strings.Trim(words[0], ".,"))] && len(words) > 1 {
		words = words[1:]
	}
	filler := map[string]bool{"a": true, "an": true, "the": true, "that": true}
	kept := make([]string, 0, 3)
	for _, word := range words {
		if filler[strings.ToLower(word)] {
			continue
		}
		kept = append(kept, strings.Trim(word, ".,"))
		if len(kept) == 3 {
			break
		}
	}
	return strings.Join(kept, " ")
}

// graphWidth measures the widest rendered line in runes, since the glyphs are
// multi-byte and a byte count would overstate the width.
func graphWidth(rendered string) int {
	widest := 0
	for _, line := range strings.Split(rendered, "\n") {
		if count := len([]rune(line)); count > widest {
			widest = count
		}
	}
	return widest
}
