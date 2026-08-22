package charter

import "strings"

// declaredSourceRefs joins every explicit charter source in its contractual
// order and removes only later duplicates. It deliberately does no lookup or
// semantic matching: resolution belongs to Compile.
func declaredSourceRefs(contract string, mission, objective, invocation []string) []string {
	raw := make([]string, 0, 1+len(mission)+len(objective)+len(invocation))
	if contract != "" {
		raw = append(raw, contract)
	}
	raw = append(raw, mission...)
	raw = append(raw, objective...)
	raw = append(raw, invocation...)

	seen := make(map[string]bool, len(raw))
	ordered := make([]string, 0, len(raw))
	for _, ref := range raw {
		ref = strings.TrimSpace(ref)
		if ref == "" || seen[ref] {
			continue
		}
		seen[ref] = true
		ordered = append(ordered, ref)
	}
	return ordered
}
