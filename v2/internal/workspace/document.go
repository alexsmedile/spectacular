package workspace

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

// Document is one canonical Spectacular Markdown record. Unknown contains
// tag-aware YAML nodes that are preserved semantically but have no authority
// in the Proposal or Mission grammar.
type Document struct {
	Record  domain.Record
	Unknown map[string]*yaml.Node
	Body    string
}

// Parse reads flat record frontmatter and its opaque Markdown body.
func Parse(data []byte) (*Document, error) {
	if !utf8.Valid(data) {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidUTF8,
			"",
			"record must be valid UTF-8",
			nil,
		)
	}
	normalized := normalizeLineEndings(string(data))
	frontmatter, body, err := splitFrontmatter(normalized)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal([]byte(frontmatter), &root); err != nil {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"decode YAML mapping",
			err,
		)
	}
	mapping, err := mappingRoot(&root)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]*yaml.Node, len(mapping.Content)/2)
	for index := 0; index < len(mapping.Content); index += 2 {
		keyNode := mapping.Content[index]
		valueNode := mapping.Content[index+1]
		if keyNode.Kind != yaml.ScalarNode || keyNode.ShortTag() != "!!str" || keyNode.Value == "" {
			return nil, domain.NewRefusal(
				domain.RefusalInvalidFrontmatter,
				"",
				"top-level property names must be non-empty strings",
				nil,
			)
		}
		if _, exists := properties[keyNode.Value]; exists {
			return nil, domain.NewRefusal(
				domain.RefusalInvalidFrontmatter,
				keyNode.Value,
				"duplicate top-level property",
				nil,
			)
		}
		properties[keyNode.Value] = valueNode
	}

	typeText, err := requiredString(properties, "type", false)
	if err != nil {
		return nil, err
	}
	recordType, err := domain.ParseRecordType(typeText)
	if err != nil {
		return nil, err
	}
	idText, err := requiredString(properties, "id", false)
	if err != nil {
		return nil, err
	}
	id, err := domain.ParseID(idText)
	if err != nil {
		return nil, err
	}

	record := domain.Record{Type: recordType, ID: id}
	if record.Title, err = optionalString(properties, "title", false); err != nil {
		return nil, err
	}
	if record.Description, err = optionalString(properties, "description", false); err != nil {
		return nil, err
	}
	if record.Status, err = optionalString(properties, "status", false); err != nil {
		return nil, err
	}
	if record.CreatedBy, err = optionalString(properties, "created_by", false); err != nil {
		return nil, err
	}
	if record.Created, err = optionalString(properties, "created", true); err != nil {
		return nil, err
	}
	if record.Updated, err = optionalString(properties, "updated", true); err != nil {
		return nil, err
	}
	if sourceText, sourceErr := optionalString(properties, "source", false); sourceErr != nil {
		return nil, sourceErr
	} else if sourceText != nil {
		reference, parseErr := domain.ParseReference(*sourceText)
		if parseErr != nil {
			return nil, parseErr
		}
		record.Source = &reference
	}
	if err := record.Validate(); err != nil {
		return nil, err
	}

	unknown := make(map[string]*yaml.Node)
	for name, node := range properties {
		if domain.IsReservedField(name) {
			continue
		}
		unknown[name] = cloneYAMLNode(node)
	}

	return &Document{Record: record, Unknown: unknown, Body: body}, nil
}

// ReadFile parses a Markdown record at path.
func ReadFile(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workspace record %q: %w", path, err)
	}
	document, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse workspace record %q: %w", path, err)
	}
	return document, nil
}

func splitFrontmatter(content string) (string, string, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", "", domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"record must start with a --- delimiter",
			nil,
		)
	}
	closing := -1
	for index := 1; index < len(lines); index++ {
		if lines[index] == "---" {
			closing = index
			break
		}
	}
	if closing < 0 {
		return "", "", domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"record has no closing --- delimiter",
			nil,
		)
	}
	frontmatter := strings.Join(lines[1:closing], "\n") + "\n"
	body := strings.Join(lines[closing+1:], "\n")
	return frontmatter, body, nil
}

func mappingRoot(root *yaml.Node) (*yaml.Node, error) {
	if root.Kind != yaml.DocumentNode || len(root.Content) != 1 {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"frontmatter must contain exactly one YAML document",
			nil,
		)
	}
	mapping := root.Content[0]
	if mapping.Kind != yaml.MappingNode || len(mapping.Content)%2 != 0 {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidFrontmatter,
			"",
			"frontmatter root must be a mapping",
			nil,
		)
	}
	return mapping, nil
}

func requiredString(properties map[string]*yaml.Node, name string, allowTimestamp bool) (string, error) {
	value, err := optionalString(properties, name, allowTimestamp)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", domain.NewRefusal(
			domain.RefusalMissingRequiredField,
			name,
			"required property is absent",
			nil,
		)
	}
	return *value, nil
}

func optionalString(properties map[string]*yaml.Node, name string, allowTimestamp bool) (*string, error) {
	node, exists := properties[name]
	if !exists {
		return nil, nil
	}
	if node.Kind != yaml.ScalarNode {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidKnownField,
			name,
			"must be a scalar string",
			nil,
		)
	}
	tag := node.ShortTag()
	if tag != "!!str" && !(allowTimestamp && tag == "!!timestamp") {
		return nil, domain.NewRefusal(
			domain.RefusalInvalidKnownField,
			name,
			"must be a string",
			nil,
		)
	}
	value := node.Value
	return &value, nil
}

func normalizeLineEndings(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	return strings.ReplaceAll(content, "\r", "\n")
}

func cloneYAMLNode(root *yaml.Node) *yaml.Node {
	clones := make(map[*yaml.Node]*yaml.Node)
	var clone func(*yaml.Node) *yaml.Node
	clone = func(node *yaml.Node) *yaml.Node {
		if node == nil {
			return nil
		}
		if existing, found := clones[node]; found {
			return existing
		}
		copied := *node
		copied.Content = nil
		copied.Alias = nil
		clones[node] = &copied
		for _, child := range node.Content {
			copied.Content = append(copied.Content, clone(child))
		}
		copied.Alias = clone(node.Alias)
		return &copied
	}
	return clone(root)
}
