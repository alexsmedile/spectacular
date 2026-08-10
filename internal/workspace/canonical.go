package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

// Canonical renders normalized semantic content. The result is UTF-8 with LF
// endings, deterministic property ordering, and exactly one terminal newline.
func Canonical(document *Document) ([]byte, error) {
	if document == nil {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"document is nil",
			nil,
		)
	}
	if err := document.Record.Validate(); err != nil {
		return nil, err
	}

	mapping := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	appendString := func(name, value string) {
		mapping.Content = append(
			mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: name},
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value},
		)
	}
	appendString("type", string(document.Record.Type))
	appendString("id", document.Record.ID.String())
	appendOptionalString(mapping, "title", document.Record.Title)
	appendOptionalString(mapping, "description", document.Record.Description)
	appendOptionalString(mapping, "status", document.Record.Status)
	appendOptionalString(mapping, "created_by", document.Record.CreatedBy)
	appendOptionalString(mapping, "created", document.Record.Created)
	appendOptionalString(mapping, "updated", document.Record.Updated)
	if document.Record.Source != nil {
		appendString("source", document.Record.Source.String())
	}

	unknownNames := make([]string, 0, len(document.Unknown))
	for name := range document.Unknown {
		if domain.IsReservedField(name) {
			return nil, domain.NewRefusal(
				domain.RefusalInvalidKnownField,
				name,
				"reserved property cannot be stored as unknown",
				nil,
			)
		}
		unknownNames = append(unknownNames, name)
	}
	sort.Strings(unknownNames)
	unknownValues := make([]*yaml.Node, 0, len(unknownNames))
	for _, name := range unknownNames {
		unknownValues = append(unknownValues, document.Unknown[name])
	}
	if err := validateYAMLTree(unknownValues...); err != nil {
		return nil, err
	}
	for _, name := range unknownNames {
		valueNode, err := normalizeYAMLNode(document.Unknown[name])
		if err != nil {
			return nil, domain.NewRefusal(
				domain.RefusalInvalidFrontmatter,
				name,
				"normalize unknown property",
				err,
			)
		}
		mapping.Content = append(
			mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: name},
			valueNode,
		)
	}

	root := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{mapping}}
	frontmatter, err := yaml.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("encode canonical frontmatter: %w", err)
	}
	frontmatterText := normalizeLineEndings(string(frontmatter))
	frontmatterText = strings.TrimRight(frontmatterText, "\n") + "\n"

	body := normalizeLineEndings(document.Body)
	body = strings.TrimRight(body, "\n")
	result := "---\n" + frontmatterText + "---\n"
	if body != "" {
		result += body + "\n"
	}
	if !utf8.ValidString(result) {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidUTF8,
			"",
			"canonical record must be valid UTF-8",
			nil,
		)
	}
	return []byte(result), nil
}

// Fingerprint returns the lowercase SHA-256 digest of canonical semantic
// content. The digest is computed and never inserted into frontmatter.
func Fingerprint(document *Document) (string, error) {
	canonical, err := Canonical(document)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(canonical)
	return hex.EncodeToString(digest[:]), nil
}

func appendOptionalString(mapping *yaml.Node, name string, value *string) {
	if value == nil {
		return
	}
	mapping.Content = append(
		mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: name},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: *value},
	)
}

func normalizeYAMLNode(root *yaml.Node) (*yaml.Node, error) {
	if root == nil {
		return nil, fmt.Errorf("nil YAML node")
	}
	if root.Kind == yaml.DocumentNode {
		return nil, fmt.Errorf("unknown property must be a YAML value, not a document")
	}

	normalized := *root
	normalized.Style = 0
	normalized.HeadComment = ""
	normalized.LineComment = ""
	normalized.FootComment = ""
	normalized.Anchor = ""
	normalized.Alias = nil
	normalized.Content = nil
	for _, child := range root.Content {
		normalizedChild, err := normalizeYAMLNode(child)
		if err != nil {
			return nil, err
		}
		normalized.Content = append(normalized.Content, normalizedChild)
	}
	if normalized.Kind != yaml.MappingNode || len(normalized.Content)%2 != 0 {
		return &normalized, nil
	}

	type pair struct {
		key   *yaml.Node
		value *yaml.Node
	}
	pairs := make([]pair, 0, len(normalized.Content)/2)
	for index := 0; index < len(normalized.Content); index += 2 {
		pairs = append(pairs, pair{key: normalized.Content[index], value: normalized.Content[index+1]})
	}
	sort.SliceStable(pairs, func(left, right int) bool {
		return nodeSortKey(pairs[left].key) < nodeSortKey(pairs[right].key)
	})
	normalized.Content = normalized.Content[:0]
	for _, item := range pairs {
		normalized.Content = append(normalized.Content, item.key, item.value)
	}
	return &normalized, nil
}

const unsupportedYAMLGraphDetail = "YAML anchors, aliases, shared graphs, and cyclic graphs are unsupported"

func validateYAMLTree(roots ...*yaml.Node) error {
	visited := make(map[*yaml.Node]bool)
	visiting := make(map[*yaml.Node]bool)
	var unsupported func(*yaml.Node) bool
	unsupported = func(node *yaml.Node) bool {
		if node == nil {
			return false
		}
		if node.Anchor != "" || node.Kind == yaml.AliasNode || node.Alias != nil || visiting[node] || visited[node] {
			return true
		}
		visiting[node] = true
		for _, child := range node.Content {
			if unsupported(child) {
				return true
			}
		}
		delete(visiting, node)
		visited[node] = true
		return false
	}
	for _, root := range roots {
		if unsupported(root) {
			return domain.NewRefusal(
				domain.RefusalUnsupportedYAMLGraph,
				"",
				unsupportedYAMLGraphDetail,
				nil,
			)
		}
	}
	return nil
}

func nodeSortKey(node *yaml.Node) string {
	encoded, err := yaml.Marshal(node)
	if err != nil {
		return node.Tag + "\x00" + node.Value
	}
	return node.Tag + "\x00" + string(encoded)
}
