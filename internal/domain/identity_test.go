package domain

import "testing"

const (
	proposalIDText = "018f2d8e-7b12-7cc3-8a45-123456789abc"
	missionIDText  = "018f2d8e-7b13-7aa1-9b34-acdeffedcba9"
)

func TestParseIDRequiresCanonicalUUIDv7(t *testing.T) {
	t.Parallel()

	valid, err := ParseID(proposalIDText)
	if err != nil {
		t.Fatalf("ParseID(valid): %v", err)
	}
	if valid.String() != proposalIDText {
		t.Fatalf("ParseID(valid) = %q", valid)
	}

	tests := map[string]string{
		"version 4": "550e8400-e29b-41d4-a716-446655440000",
		"uppercase": "018F2D8E-7B12-7CC3-8A45-123456789ABC",
		"compact":   "018f2d8e7b127cc38a45123456789abc",
		"malformed": "not-a-uuid",
	}
	for name, input := range tests {
		input := input
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseID(input)
			if !RefusalHasCode(err, RefusalInvalidID) {
				t.Fatalf("ParseID(%q) error = %v, want invalid_id refusal", input, err)
			}
		})
	}
}

func TestNewIDCreatesCanonicalUUIDv7(t *testing.T) {
	t.Parallel()

	id, err := NewID()
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	parsed, err := ParseID(id.String())
	if err != nil {
		t.Fatalf("generated ID is not canonical UUIDv7: %v", err)
	}
	if parsed != id {
		t.Fatalf("ParseID(NewID()) = %q, want %q", parsed, id)
	}
}
