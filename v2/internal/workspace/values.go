package workspace

import (
	"fmt"

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
