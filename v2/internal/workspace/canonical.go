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
	normalizer := yamlGraphNormalizer{
		clones:       make(map[*yaml.Node]*yaml.Node),
		aliasTargets: make(map[*yaml.Node]bool),
	}
	normalized, err := normalizer.clone(root)
	if err != nil {
		return nil, err
	}
	normalizer.sortMappings(normalized, make(map[*yaml.Node]bool))
	if err := normalizer.canonicalizeAnchors(normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}

type yamlGraphNormalizer struct {
	clones       map[*yaml.Node]*yaml.Node
	aliasTargets map[*yaml.Node]bool
}

func (normalizer *yamlGraphNormalizer) clone(node *yaml.Node) (*yaml.Node, error) {
	if node == nil {
		return nil, fmt.Errorf("nil YAML node")
	}
	if existing, found := normalizer.clones[node]; found {
		return existing, nil
	}
	if node.Kind == yaml.DocumentNode {
		return nil, fmt.Errorf("unknown property must be a YAML value, not a document")
	}
	if node.Kind == yaml.AliasNode && node.Alias == nil {
		return nil, fmt.Errorf("YAML alias has no target")
	}

	cloned := *node
	cloned.Style = 0
	cloned.HeadComment = ""
	cloned.LineComment = ""
	cloned.FootComment = ""
	cloned.Content = nil
	cloned.Alias = nil
	if node.Anchor != "" {
		cloned.Anchor = "pending"
	}
	if node.Kind == yaml.AliasNode {
		cloned.Value = ""
	}
	normalizer.clones[node] = &cloned

	for _, child := range node.Content {
		clonedChild, err := normalizer.clone(child)
		if err != nil {
			return nil, err
		}
		cloned.Content = append(cloned.Content, clonedChild)
	}
	if node.Alias != nil {
		clonedTarget, err := normalizer.clone(node.Alias)
		if err != nil {
			return nil, err
		}
		cloned.Alias = clonedTarget
		normalizer.aliasTargets[clonedTarget] = true
	}
	return &cloned, nil
}

func (normalizer *yamlGraphNormalizer) sortMappings(node *yaml.Node, visited map[*yaml.Node]bool) {
	if node == nil || visited[node] {
		return
	}
	visited[node] = true
	for _, child := range node.Content {
		normalizer.sortMappings(child, visited)
	}
	normalizer.sortMappings(node.Alias, visited)
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}

	type pair struct {
		key   *yaml.Node
		value *yaml.Node
	}
	pairs := make([]pair, 0, len(node.Content)/2)
	for index := 0; index < len(node.Content); index += 2 {
		pairs = append(pairs, pair{key: node.Content[index], value: node.Content[index+1]})
	}
	sort.SliceStable(pairs, func(left, right int) bool {
		return nodeSortKey(pairs[left].key) < nodeSortKey(pairs[right].key)
	})
	node.Content = node.Content[:0]
	for _, item := range pairs {
		node.Content = append(node.Content, item.key, item.value)
	}
}

func (normalizer *yamlGraphNormalizer) canonicalizeAnchors(root *yaml.Node) error {
	visited := make(map[*yaml.Node]bool)
	anchorNames := make(map[*yaml.Node]string)
	var assign func(*yaml.Node)
	assign = func(node *yaml.Node) {
		if node == nil || visited[node] {
			return
		}
		visited[node] = true
		if node.Kind != yaml.AliasNode && (node.Anchor != "" || normalizer.aliasTargets[node]) {
			name := fmt.Sprintf("a%d", len(anchorNames)+1)
			anchorNames[node] = name
			node.Anchor = name
		}
		for _, child := range node.Content {
			assign(child)
		}
		assign(node.Alias)
	}
	assign(root)

	visited = make(map[*yaml.Node]bool)
	var updateAliases func(*yaml.Node) error
	updateAliases = func(node *yaml.Node) error {
		if node == nil || visited[node] {
			return nil
		}
		visited[node] = true
		if node.Kind == yaml.AliasNode {
			name, found := anchorNames[node.Alias]
			if !found {
				return fmt.Errorf("YAML alias target has no canonical anchor")
			}
			node.Value = name
		}
		for _, child := range node.Content {
			if err := updateAliases(child); err != nil {
				return err
			}
		}
		return updateAliases(node.Alias)
	}
	return updateAliases(root)
}

func nodeSortKey(root *yaml.Node) string {
	var key strings.Builder
	seen := make(map[*yaml.Node]int)
	var appendNode func(*yaml.Node)
	appendNode = func(node *yaml.Node) {
		if node == nil {
			key.WriteString("nil;")
			return
		}
		if index, found := seen[node]; found {
			fmt.Fprintf(&key, "ref:%d;", index)
			return
		}
		index := len(seen)
		seen[node] = index
		fmt.Fprintf(&key, "node:%d:%d:%d:%s:%d:%s:%t:[", index, node.Kind, len(node.Tag), node.Tag, len(node.Value), node.Value, node.Anchor != "")
		for _, child := range node.Content {
			appendNode(child)
		}
		key.WriteString("]alias:")
		appendNode(node.Alias)
	}
	appendNode(root)
	return key.String()
}
