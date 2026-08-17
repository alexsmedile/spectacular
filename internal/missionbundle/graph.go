package missionbundle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
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

type edgeTarget struct {
	ref           string
	interfaceOnly bool
}

type pathNode struct {
	ref           string
	interfaceOnly bool
}

// objectiveGraph draws each dependency chain on its own row, so an edge in the
// drawing is an edge in the plan. A node with several dependents branches with
// `┬` and `└`; a node with none stands alone. Artifact dependencies draw solid (──);
// interface dependencies draw dotted (┄─).
func (b *Bundle) objectiveGraph(state State) string {
	nodes := map[string]ObjectiveState{}
	order := make([]string, 0, len(state.Objectives))
	for _, objective := range state.Objectives {
		nodes[objective.Ref] = objective
		order = append(order, objective.Ref)
	}
	successors := map[string][]edgeTarget{}
	hasPredecessor := map[string]bool{}
	for _, objective := range b.Objectives {
		for _, predecessor := range objective.After {
			if _, known := nodes[predecessor]; !known {
				continue
			}
			successors[predecessor] = append(successors[predecessor], edgeTarget{ref: objective.Ref, interfaceOnly: false})
			hasPredecessor[objective.Ref] = true
		}
		for _, predecessor := range objective.AfterInterface {
			if _, known := nodes[predecessor]; !known {
				continue
			}
			successors[predecessor] = append(successors[predecessor], edgeTarget{ref: objective.Ref, interfaceOnly: true})
			hasPredecessor[objective.Ref] = true
		}
	}

	// Each row is one root-to-leaf path, so every edge appears exactly once and
	// a reader can follow a chain left to right without cross-referencing.
	var rows [][]pathNode
	var walk func(current pathNode, path []pathNode)
	walk = func(current pathNode, path []pathNode) {
		path = append(path, current)
		next := successors[current.ref]
		if len(next) == 0 {
			rows = append(rows, append([]pathNode(nil), path...))
			return
		}
		sort.Slice(next, func(i, j int) bool { return next[i].ref < next[j].ref })
		for _, target := range next {
			walk(pathNode{ref: target.ref, interfaceOnly: target.interfaceOnly}, path)
		}
	}
	for _, ref := range order {
		if !hasPredecessor[ref] {
			walk(pathNode{ref: ref, interfaceOnly: false}, nil)
		}
	}

	cell := func(ref string) string { return ref + Glyph(nodes[ref].Readiness) }
	lines := make([]string, 0, len(rows))
	var previous []pathNode
	for _, row := range rows {
		// Blank the prefix a row shares with the one above, so a branch reads as
		// a branch rather than as a repeated chain.
		shared := 0
		for shared < len(previous) && shared < len(row) && previous[shared] == row[shared] {
			shared++
		}
		var line strings.Builder
		line.WriteString("  ")
		for index, node := range row {
			ref := node.ref
			if index < shared {
				line.WriteString(strings.Repeat(" ", len([]rune(cell(ref)))))
			} else {
				line.WriteString(cell(ref))
			}
			if index == len(row)-1 {
				break
			}
			nextIsInterface := row[index+1].interfaceOnly
			switch {
			// The connector after the last shared node is where this row
			// branches away from the one above.
			case index == shared-1:
				if nextIsInterface {
					line.WriteString(" └┄")
				} else {
					line.WriteString(" └─")
				}
			case index < shared:
				line.WriteString("   ")
			// A node with several dependents forks here.
			case len(successors[ref]) > 1:
				if nextIsInterface {
					line.WriteString(" ┬┄")
				} else {
					line.WriteString(" ┬─")
				}
			default:
				if nextIsInterface {
					line.WriteString(" ┄─")
				} else {
					line.WriteString(" ──")
				}
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
			line += "    (" + waits + ")"
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
	seenAfter := map[string]bool{}
	seenInterface := map[string]bool{}
	var afterRefs []string
	var interfaceRefs []string
	for _, objective := range b.Objectives {
		if !inGroup[objective.Ref] {
			continue
		}
		for _, predecessor := range objective.After {
			if !seenAfter[predecessor] {
				seenAfter[predecessor] = true
				afterRefs = append(afterRefs, predecessor)
			}
		}
		for _, predecessor := range objective.AfterInterface {
			if !seenInterface[predecessor] {
				seenInterface[predecessor] = true
				interfaceRefs = append(interfaceRefs, predecessor)
			}
		}
	}
	sort.Strings(afterRefs)
	sort.Strings(interfaceRefs)
	parts := make([]string, 0, 2)
	if len(afterRefs) > 0 {
		parts = append(parts, "after "+strings.Join(afterRefs, ", "))
	}
	if len(interfaceRefs) > 0 {
		parts = append(parts, "after_interface "+strings.Join(interfaceRefs, ", "))
	}
	return strings.Join(parts, ", ")
}

// Timeline renders an ASCII sequence of Missions in the workspace showing
// declared Mission-to-Mission order edges and completion state.
func (b *Bundle) Timeline(ws *discovery.Workspace, width int) string {
	if width <= 0 {
		width = DefaultGraphWidth
	}
	heading := b.graphHeading() + " · Timeline"
	if ws == nil {
		return heading + "\n\n  " + b.Ref + " · " + b.Title + " " + missionGlyph(b.Status) + "\n"
	}

	entries := ws.OfType(domain.Mission)
	if len(entries) == 0 {
		return heading + "\n\n  no Missions discovered\n"
	}

	bundles := make([]*Bundle, 0, len(entries))
	for _, entry := range entries {
		m, err := decode(ws, entry)
		if err == nil && m != nil {
			bundles = append(bundles, m)
		}
	}

	sort.Slice(bundles, func(i, j int) bool {
		if bundles[i].Created != "" && bundles[j].Created != "" && bundles[i].Created != bundles[j].Created {
			return bundles[i].Created < bundles[j].Created
		}
		return bundles[i].Ref < bundles[j].Ref
	})

	var lines []string
	for _, m := range bundles {
		glyph := missionGlyph(m.Status)
		label := fmt.Sprintf("%s · %s %s", m.Ref, m.Title, glyph)
		if len(m.AfterMission) > 0 {
			label += " (after " + strings.Join(m.AfterMission, ", ") + ")"
		}
		lines = append(lines, "  "+label)
	}

	return heading + "\n\n" + strings.Join(lines, "\n") + "\n"
}

func missionGlyph(status string) string {
	switch status {
	case "completed":
		return "✓"
	case "active":
		return "◐"
	case "defined":
		return "▶"
	default:
		return "·"
	}
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
