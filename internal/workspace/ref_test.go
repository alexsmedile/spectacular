package workspace

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func documentWith(fields map[string]string) *Document {
	doc := &Document{Unknown: map[string]*yaml.Node{}}
	for name, value := range fields {
		node := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
		doc.Unknown[name] = node
	}
	return doc
}

// One decoder resolves both spellings so every caller compares refs that are
// spelled one way, and reports which spelling supplied the answer.
func TestRefAcceptsBothSpellingsAndReportsTheLegacyOne(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		fields map[string]string
		want   string
		legacy bool
	}{
		{
			name:   "the current spelling supplies the ref and reports no drift",
			fields: map[string]string{"ref": "M7"},
			want:   "M7", legacy: false,
		},
		{
			name:   "the legacy spelling supplies the ref and reports drift",
			fields: map[string]string{"human_ref": "M2"},
			want:   "M2", legacy: true,
		},
		{
			name:   "the current spelling wins when a record carries both",
			fields: map[string]string{"ref": "M5", "human_ref": "M5-old"},
			want:   "M5", legacy: false,
		},
		{
			name:   "a record with neither spelling yields an empty ref, not an error",
			fields: map[string]string{},
			want:   "", legacy: false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ref, legacy, err := Ref(documentWith(testCase.fields))
			if err != nil {
				t.Fatalf("Ref returned %v, want no error", err)
			}
			if ref != testCase.want {
				t.Fatalf("ref %q, want %q", ref, testCase.want)
			}
			if legacy != testCase.legacy {
				t.Fatalf("legacy=%t, want %t", legacy, testCase.legacy)
			}
		})
	}
}

// A malformed ref is an error for callers that refuse, and empty for
// presentation callers that must not fail on an optional field.
func TestRefSurfacesMalformedFieldsWhileRefOrEmptyStaysSilent(t *testing.T) {
	doc := &Document{Unknown: map[string]*yaml.Node{
		"ref": {Kind: yaml.SequenceNode, Tag: "!!seq"},
	}}
	if _, _, err := Ref(doc); err == nil {
		t.Fatal("a ref that is not a string must be an error for callers that refuse")
	}
	if got := RefOrEmpty(doc); got != "" {
		t.Fatalf("RefOrEmpty returned %q, want empty on a malformed field", got)
	}
}
