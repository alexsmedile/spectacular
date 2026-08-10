package workspace

import (
	"fmt"
	"strconv"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

func String(document *Document, name string, required bool) (string, error) {
	node := document.Unknown[name]
	if node == nil {
		if required {
			return "", domain.NewRefusal(domain.RefusalMissingRequiredField, name, "required property is absent", nil)
		}
		return "", nil
	}
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!str" || node.Value == "" {
		return "", domain.NewRefusal(domain.RefusalInvalidKnownField, name, "must be a non-empty string", nil)
	}
	return node.Value, nil
}

func Bool(document *Document, name string, required bool) (bool, error) {
	node := document.Unknown[name]
	if node == nil {
		if required {
			return false, domain.NewRefusal(domain.RefusalMissingRequiredField, name, "required property is absent", nil)
		}
		return false, nil
	}
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!bool" {
		return false, domain.NewRefusal(domain.RefusalInvalidKnownField, name, "must be a boolean", nil)
	}
	return node.Value == "true", nil
}

func Int(document *Document, name string, required bool) (int, error) {
	node := document.Unknown[name]
	if node == nil {
		if required {
			return 0, domain.NewRefusal(domain.RefusalMissingRequiredField, name, "required property is absent", nil)
		}
		return 0, nil
	}
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!int" {
		return 0, domain.NewRefusal(domain.RefusalInvalidKnownField, name, "must be an integer", nil)
	}
	value, err := strconv.Atoi(node.Value)
	if err != nil {
		return 0, domain.NewRefusal(domain.RefusalInvalidKnownField, name, "must be an integer", err)
	}
	return value, nil
}

func Strings(document *Document, name string, required bool) ([]string, error) {
	node := document.Unknown[name]
	if node == nil {
		if required {
			return nil, domain.NewRefusal(domain.RefusalMissingRequiredField, name, "required property is absent", nil)
		}
		return nil, nil
	}
	if node.Kind != yaml.SequenceNode {
		return nil, domain.NewRefusal(domain.RefusalInvalidKnownField, name, "must be a sequence of strings", nil)
	}
	values := make([]string, len(node.Content))
	for i, child := range node.Content {
		if child.Kind != yaml.ScalarNode || child.ShortTag() != "!!str" || child.Value == "" {
			return nil, domain.NewRefusal(domain.RefusalInvalidKnownField, name, fmt.Sprintf("item %d must be a non-empty string", i), nil)
		}
		values[i] = child.Value
	}
	return values, nil
}

// SetString and SetStrings construct canonical unknown-field values for
// governed operations. Unknown fields remain type-specific data rather than
// gaining universal Record authority.
func SetString(document *Document, name, value string) {
	ensureUnknown(document)
	document.Unknown[name] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func SetBool(document *Document, name string, value bool) {
	ensureUnknown(document)
	text := "false"
	if value {
		text = "true"
	}
	document.Unknown[name] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: text}
}

func SetInt(document *Document, name string, value int) {
	ensureUnknown(document)
	document.Unknown[name] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(value)}
}

func SetStrings(document *Document, name string, values []string) {
	ensureUnknown(document)
	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value})
	}
	document.Unknown[name] = node
}

func Delete(document *Document, name string) {
	if document != nil && document.Unknown != nil {
		delete(document.Unknown, name)
	}
}

func ensureUnknown(document *Document) {
	if document.Unknown == nil {
		document.Unknown = map[string]*yaml.Node{}
	}
}
