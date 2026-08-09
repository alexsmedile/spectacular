package workspace

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
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
			if !reflect.DeepEqual(original.Unknown, roundTripped.Unknown) {
				t.Fatalf("unknown properties changed:\noriginal: %#v\nround trip: %#v", original.Unknown, roundTripped.Unknown)
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
	if got := document.Unknown["priority"]; got != "high" {
		t.Fatalf("priority = %#v, want high", got)
	}
	review, ok := document.Unknown["review"].(map[string]any)
	if !ok {
		t.Fatalf("review type = %T, want map[string]any", document.Unknown["review"])
	}
	if review["required"] != true || review["reviewer"] != "H28" {
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

	document, err := Parse([]byte("---\nzeta: last\nstatus: approved\nid: 018f2d8e-7b12-7cc3-8a45-123456789abc\ntype: Proposal\nalpha: first\n---\nbody\n"))
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
