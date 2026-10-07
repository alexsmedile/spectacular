package discovery

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

// Canonical record names remain evidence of governed intent even when their
// frontmatter identity is damaged. Knowledge uses descriptive, unnumbered names.
var governedName = regexp.MustCompile(`^(?:CC-|[A-Z]+[0-9]+(?:[-.]|$))`)

// softKnowledge recognizes only well-formed, unclaimed knowledge. It never
// downgrades a malformed record, a project binding, or a governed collection.
func softKnowledge(meta, anchor, absolute string, data []byte) bool {
	rel, err := filepath.Rel(meta, absolute)
	if err != nil || filepath.ToSlash(rel) == anchor {
		return false
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	for _, part := range parts[:len(parts)-1] {
		switch part {
		case "contracts", "proposals", "missions", "evidence", "gaps", "handoffs", "assessments", "reviews", "archive":
			return false
		}
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return false
	}
	frontmatter, _, found := strings.Cut(text[4:], "\n---\n")
	if !found {
		return false
	}
	var fields map[string]any
	if yaml.Unmarshal([]byte(frontmatter), &fields) != nil {
		return false
	}
	for _, claim := range []string{"id", "ref", "human_ref", "schema", "schema_version"} {
		if _, exists := fields[claim]; exists {
			return false
		}
	}
	// Explicit context metadata resolves naming ambiguity without overriding an
	// identity/schema claim or a strict governed collection. Historical damaged
	// numbered records without that declaration still fail governed validation.
	if fields["governance"] == "context" {
		typeName, ok := fields["type"].(string)
		return ok && strings.TrimSpace(typeName) != ""
	}
	if governedName.MatchString(filepath.Base(rel)) {
		return false
	}
	typeName, ok := fields["type"].(string)
	if !ok || strings.TrimSpace(typeName) == "" {
		return false
	}
	kind, err := domain.ParseRecordType(typeName)
	if err != nil {
		return true // OKF types are extensible; no governed identity was claimed.
	}
	return kind == domain.Anchor || (kind == domain.Decision && len(parts) > 1 && parts[0] == "decisions")
}
