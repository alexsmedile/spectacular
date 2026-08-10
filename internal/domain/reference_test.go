package domain

import "testing"

func TestParseReference(t *testing.T) {
	t.Parallel()

	reference, err := ParseReference("Proposal:" + proposalIDText)
	if err != nil {
		t.Fatalf("ParseReference: %v", err)
	}
	if reference.Type != Proposal || reference.ID.String() != proposalIDText {
		t.Fatalf("ParseReference returned %#v", reference)
	}
	if got := reference.String(); got != "Proposal:"+proposalIDText {
		t.Fatalf("Reference.String() = %q", got)
	}

	for _, input := range []string{
		proposalIDText,
		"Unknown:" + proposalIDText,
		"Proposal:not-a-uuid",
		"Proposal:" + proposalIDText + ":extra",
	} {
		if _, err := ParseReference(input); !RefusalHasCode(err, RefusalInvalidReference) {
			t.Errorf("ParseReference(%q) error = %v, want invalid_reference", input, err)
		}
	}
}
