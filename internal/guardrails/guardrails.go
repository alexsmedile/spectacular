// Package guardrails selects owner-written workflow guidance without creating
// authority or interpreting the selected prose.
package guardrails

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var selectorPattern = regexp.MustCompile(`^\$[a-z][a-z0-9-]*\.[a-z][a-z0-9-]*$`)

var events = map[string]bool{
	"@Orient": true, "@Prepare": true, "@Start": true, "@Resume": true,
	"@Run": true, "@Assess": true, "@Reconcile": true, "@Resolve": true,
}

type Section struct {
	Event    string `json:"event"`
	Selector string `json:"selector,omitempty"`
	Prose    string `json:"prose"`
}

type Document struct {
	Sections []Section `json:"sections"`
}

func Parse(data []byte) (Document, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	var document Document
	seen := map[string]bool{}
	for i := 0; i < len(lines); {
		line := lines[i]
		if !strings.HasPrefix(line, "## @") {
			i++
			continue
		}
		fields := strings.Fields(strings.TrimPrefix(line, "## "))
		if len(fields) < 1 || len(fields) > 2 || !events[fields[0]] {
			return Document{}, fmt.Errorf("invalid Guardrails heading on line %d", i+1)
		}
		selector := ""
		if len(fields) == 2 {
			selector = fields[1]
			if !selectorPattern.MatchString(selector) {
				return Document{}, fmt.Errorf("invalid Guardrails selector on line %d", i+1)
			}
		}
		key := fields[0] + "\x00" + selector
		if seen[key] {
			return Document{}, fmt.Errorf("duplicate Guardrails section %s %s", fields[0], selector)
		}
		seen[key] = true
		start := i + 1
		i = start
		for i < len(lines) && !strings.HasPrefix(lines[i], "## @") {
			i++
		}
		prose := strings.Trim(strings.Join(lines[start:i], "\n"), "\n")
		if strings.TrimSpace(prose) == "" {
			return Document{}, fmt.Errorf("empty Guardrails section %s %s", fields[0], selector)
		}
		document.Sections = append(document.Sections, Section{Event: fields[0], Selector: selector, Prose: prose})
	}
	return document, nil
}

func (d Document) Select(event, selector string) ([]Section, error) {
	if !events[event] {
		return nil, fmt.Errorf("unknown Guardrails event %q", event)
	}
	if selector != "" && !selectorPattern.MatchString(selector) {
		return nil, fmt.Errorf("invalid Guardrails selector %q", selector)
	}
	var selected []Section
	for _, section := range d.Sections {
		if section.Event == event && (section.Selector == "" || section.Selector == selector) {
			selected = append(selected, section)
		}
	}
	sort.SliceStable(selected, func(i, j int) bool { return selected[i].Selector < selected[j].Selector })
	return selected, nil
}
