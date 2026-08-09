package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v4"
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
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(document.Unknown[name]); err != nil {
			return nil, domain.NewRefusal(
				domain.RefusalInvalidFrontmatter,
				name,
				"encode unknown property",
				err,
			)
		}
		canonicalizeNode(valueNode)
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

func canonicalizeNode(node *yaml.Node) {
	node.Style = 0
	node.HeadComment = ""
	node.LineComment = ""
	node.FootComment = ""
	node.Anchor = ""

	for _, child := range node.Content {
		canonicalizeNode(child)
	}
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

func nodeSortKey(node *yaml.Node) string {
	encoded, err := yaml.Marshal(node)
	if err != nil {
		return node.Tag + "\x00" + node.Value
	}
	return node.Tag + "\x00" + string(encoded)
}
