package domain

import (
	"errors"
	"testing"
)

func pointer(value string) *string {
	return &value
}

func mustID(t *testing.T, value string) ID {
	t.Helper()
	id, err := ParseID(value)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestRecordValidateProposalAndMissionGrammars(t *testing.T) {
	t.Parallel()

	proposalID := mustID(t, proposalIDText)
	missionID := mustID(t, missionIDText)
	source := Reference{Type: Proposal, ID: proposalID}

	proposal := Record{
		Type:      Proposal,
		ID:        proposalID,
		Title:     pointer("Semantic substrate"),
		Status:    pointer("approved"),
		CreatedBy: pointer("owner"),
		Created:   pointer("2026-08-09T10:00:00Z"),
		Updated:   pointer("2026-08-09T10:30:00Z"),
	}
	if err := proposal.Validate(); err != nil {
		t.Fatalf("proposal.Validate: %v", err)
	}

	mission := Record{Type: Mission, ID: missionID, Source: &source}
	if err := mission.Validate(); err != nil {
		t.Fatalf("mission.Validate: %v", err)
	}
}

func TestRecordValidateRefusesInvalidKnownFields(t *testing.T) {
	t.Parallel()

	proposalID := mustID(t, proposalIDText)
	missionID := mustID(t, missionIDText)
	proposalSource := Reference{Type: Proposal, ID: proposalID}
	missionSource := Reference{Type: Mission, ID: missionID}

	tests := []struct {
		name   string
		record Record
		code   RefusalCode
	}{
		{
			name:   "empty title",
			record: Record{Type: Proposal, ID: proposalID, Title: pointer(" ")},
			code:   RefusalInvalidKnownField,
		},
		{
			name:   "invalid timestamp",
			record: Record{Type: Proposal, ID: proposalID, Created: pointer("yesterday")},
			code:   RefusalInvalidKnownField,
		},
		{
			name:   "proposal source",
			record: Record{Type: Proposal, ID: proposalID, Source: &proposalSource},
			code:   RefusalInvalidKnownField,
		},
		{
			name:   "mission source type",
			record: Record{Type: Mission, ID: missionID, Source: &missionSource},
			code:   RefusalInvalidReference,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := test.record.Validate(); !RefusalHasCode(err, test.code) {
				t.Fatalf("Validate error = %v, want %s", err, test.code)
			}
		})
	}
}

func TestOnlyTypeAndIDAreUniversalRequirements(t *testing.T) {
	t.Parallel()

	record := Record{Type: Proposal, ID: mustID(t, proposalIDText)}
	if err := record.Validate(); err != nil {
		t.Fatalf("minimal record rejected: %v", err)
	}
}

func TestRecordValidationRefusalOrderIsDeterministic(t *testing.T) {
	t.Parallel()

	record := Record{
		Type:      Proposal,
		ID:        mustID(t, proposalIDText),
		Title:     pointer(""),
		Status:    pointer(""),
		CreatedBy: pointer(""),
	}
	for attempt := 0; attempt < 100; attempt++ {
		err := record.Validate()
		var refusal *Refusal
		if !errors.As(err, &refusal) {
			t.Fatalf("Validate error = %v, want refusal", err)
		}
		if refusal.Field != "title" {
			t.Fatalf("Validate refusal field = %q, want title", refusal.Field)
		}
	}
}
