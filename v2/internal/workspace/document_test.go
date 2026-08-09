package workspace

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

func fixturePath(parts ...string) string {
	all := append([]string{"..", "..", "testdata", "m1"}, parts...)
	return filepath.Join(all...)
}

func TestApprovedProposalAndMissionSemanticRoundTrip(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"proposal.md", "mission.md"} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			original, err := ReadFile(fixturePath("positive", name))
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			canonical, err := Canonical(original)
			if err != nil {
				t.Fatalf("Canonical: %v", err)
			}
			roundTripped, err := Parse(canonical)
			if err != nil {
				t.Fatalf("Parse(canonical): %v\n%s", err, canonical)
			}
			canonicalAgain, err := Canonical(roundTripped)
			if err != nil {
				t.Fatalf("Canonical(roundTripped): %v", err)
			}
			if !bytes.Equal(canonical, canonicalAgain) {
				t.Fatalf("canonical round-trip changed:\nfirst:\n%s\nsecond:\n%s", canonical, canonicalAgain)
			}
			if !reflect.DeepEqual(original.Record, roundTripped.Record) {
				t.Fatalf("record changed:\noriginal: %#v\nround trip: %#v", original.Record, roundTripped.Record)
			}
			if original.Body != roundTripped.Body {
				t.Fatalf("Markdown body changed:\noriginal: %q\nround trip: %q", original.Body, roundTripped.Body)
			}
		})
	}
}

func TestUnknownPropertiesAndMarkdownBodyArePreserved(t *testing.T) {
	t.Parallel()

	document, err := ReadFile(fixturePath("positive", "proposal.md"))
	if err != nil {
		t.Fatal(err)
	}
	priority := document.Unknown["priority"]
	if priority == nil || priority.ShortTag() != "!!str" || priority.Value != "high" {
		t.Fatalf("priority = %#v, want tagged string high", priority)
	}
	review := document.Unknown["review"]
	if review == nil || review.Kind != yaml.MappingNode {
		t.Fatalf("review = %#v, want mapping node", review)
	}
	if required := mappingValue(review, "required"); required == nil || required.ShortTag() != "!!bool" || required.Value != "true" {
		t.Fatalf("review.required = %#v", required)
	}
	if reviewer := mappingValue(review, "reviewer"); reviewer == nil || reviewer.ShortTag() != "!!str" || reviewer.Value != "H28" {
		t.Fatalf("review = %#v", review)
	}
	if !strings.Contains(document.Body, "Unknown properties and this Markdown body") {
		t.Fatalf("body not preserved: %q", document.Body)
	}

	canonical, err := Canonical(document)
	if err != nil {
		t.Fatal(err)
	}
	text := string(canonical)
	if strings.Contains(text, "fingerprint:") {
		t.Fatalf("canonical output stored a fingerprint:\n%s", text)
	}
	if !strings.HasSuffix(text, "\n") || strings.HasSuffix(text, "\n\n") {
		t.Fatalf("canonical output must have exactly one terminal newline: %q", text)
	}
}

func TestParseRefusesMissingOrInvalidKnownFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		code domain.RefusalCode
	}{
		{
			name: "missing id",
			data: "---\ntype: Proposal\n---\n",
			code: domain.RefusalMissingRequiredField,
		},
		{
			name: "non-string title",
			data: "---\ntype: Proposal\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\ntitle: [invalid]\n---\n",
			code: domain.RefusalInvalidKnownField,
		},
		{
			name: "source on proposal",
			data: "---\ntype: Proposal\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\nsource: Proposal:018f2d8e-7b12-7cc3-8a45-123456789abc\n---\n",
			code: domain.RefusalInvalidKnownField,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse([]byte(test.data))
			if !domain.RefusalHasCode(err, test.code) {
				t.Fatalf("Parse error = %v, want %s", err, test.code)
			}
		})
	}
}

func TestCanonicalKnownAndUnknownFieldOrdering(t *testing.T) {
	t.Parallel()

	document, err := Parse([]byte("---\nzeta: last\nstatus: accepted\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\ntype: Proposal\nalpha: first\n---\nbody\n"))
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := Canonical(document)
	if err != nil {
		t.Fatal(err)
	}
	wantOrder := []string{"type:", "id:", "status:", "alpha:", "zeta:"}
	last := -1
	for _, marker := range wantOrder {
		position := strings.Index(string(canonical), marker)
		if position <= last {
			t.Fatalf("%q is not in canonical order:\n%s", marker, canonical)
		}
		last = position
	}
}

func TestComplexKeyUnknownValueIsPreserved(t *testing.T) {
	t.Parallel()

	input := "---\n" +
		"type: Proposal\n" +
		"id: 018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
		"opaque:\n" +
		"  ? [alpha, beta]\n" +
		"  : sequence-key\n" +
		"---\n" +
		"Complex keys are valid unknown YAML values.\n"
	document, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse complex-key value: %v", err)
	}
	opaque := document.Unknown["opaque"]
	if opaque == nil || opaque.Kind != yaml.MappingNode || len(opaque.Content) != 2 {
		t.Fatalf("opaque = %#v, want one-entry mapping", opaque)
	}
	if opaque.Content[0].Kind != yaml.SequenceNode {
		t.Fatalf("opaque key kind = %v, want sequence", opaque.Content[0].Kind)
	}
	canonical, err := Canonical(document)
	if err != nil {
		t.Fatalf("Canonical complex-key value: %v", err)
	}
	roundTripped, err := Parse(canonical)
	if err != nil {
		t.Fatalf("Parse canonical complex-key value: %v\n%s", err, canonical)
	}
	roundTrippedOpaque := roundTripped.Unknown["opaque"]
	if roundTrippedOpaque == nil || roundTrippedOpaque.Kind != yaml.MappingNode || len(roundTrippedOpaque.Content) != 2 {
		t.Fatalf("round-tripped opaque = %#v, want one-entry mapping", roundTrippedOpaque)
	}
	if roundTrippedOpaque.Content[0].Kind != yaml.SequenceNode || len(roundTrippedOpaque.Content[0].Content) != 2 {
		t.Fatalf("round-tripped complex key = %#v", roundTrippedOpaque.Content[0])
	}
	if roundTrippedOpaque.Content[0].Content[0].Value != "alpha" || roundTrippedOpaque.Content[0].Content[1].Value != "beta" {
		t.Fatalf("round-tripped complex key values = %#v", roundTrippedOpaque.Content[0].Content)
	}
	canonicalAgain, err := Canonical(roundTripped)
	if err != nil {
		t.Fatalf("Canonical round trip: %v", err)
	}
	if !bytes.Equal(canonical, canonicalAgain) {
		t.Fatalf("complex-key value changed:\nfirst:\n%s\nsecond:\n%s", canonical, canonicalAgain)
	}
}

func TestInvalidUTF8IsRefused(t *testing.T) {
	t.Parallel()

	input := append([]byte("---\ntype: Proposal\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\n---\n"), 0xff)
	if _, err := Parse(input); !domain.RefusalHasCode(err, domain.RefusalInvalidUTF8) {
		t.Fatalf("Parse invalid UTF-8 error = %v, want invalid_utf8", err)
	}

	document, err := Parse([]byte("---\ntype: Proposal\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\n---\n"))
	if err != nil {
		t.Fatal(err)
	}
	document.Body = string([]byte{0xff})
	if _, err := Canonical(document); !domain.RefusalHasCode(err, domain.RefusalInvalidUTF8) {
		t.Fatalf("Canonical invalid UTF-8 error = %v, want invalid_utf8", err)
	}
}

func TestCanonicalRefusesInvalidProgrammaticReference(t *testing.T) {
	t.Parallel()

	missionID, err := domain.ParseID("018f2d8e-7b13-7aa1-9b34-acdeffedcba9")
	if err != nil {
		t.Fatal(err)
	}
	document := &Document{Record: domain.Record{
		Type: domain.Mission,
		ID:   missionID,
		Source: &domain.Reference{
			Type: domain.Proposal,
			ID:   domain.ID("not-a-uuid"),
		},
	}}
	if _, err := Canonical(document); !domain.RefusalHasCode(err, domain.RefusalInvalidReference) {
		t.Fatalf("Canonical error = %v, want invalid_reference", err)
	}
}

func TestReadFileIncludesPathOnParseFailure(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "broken.md")
	if err := os.WriteFile(path, []byte("no frontmatter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadFile(path)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("ReadFile error = %v, want path", err)
	}
}

func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}
