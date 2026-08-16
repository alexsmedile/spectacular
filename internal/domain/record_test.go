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
		Status:    pointer("accepted"),
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

func TestRecordValidateExactStatusVocabulary(t *testing.T) {
	t.Parallel()

	id := mustID(t, proposalIDText)
	allowed := map[RecordType][]string{
		Proposal: {"draft", "submitted", "accepted", "rejected", "withdrawn"},
		Mission:  {"defined", "active", "awaiting-assessment", "resolved", "completed"},
	}
	for recordType, statuses := range allowed {
		recordType := recordType
		for _, status := range statuses {
			status := status
			t.Run(string(recordType)+"/"+status, func(t *testing.T) {
				t.Parallel()
				record := Record{Type: recordType, ID: id, Status: pointer(status)}
				if err := record.Validate(); err != nil {
					t.Fatalf("Validate %s status %q: %v", recordType, status, err)
				}
			})
		}
	}

	invalid := []struct {
		recordType RecordType
		status     string
	}{
		{recordType: Proposal, status: "approved"},
		{recordType: Proposal, status: " accepted "},
		{recordType: Mission, status: "definitely-not-a-status"},
	}
	for _, test := range invalid {
		test := test
		t.Run(string(test.recordType)+"/invalid/"+test.status, func(t *testing.T) {
			t.Parallel()
			record := Record{Type: test.recordType, ID: id, Status: pointer(test.status)}
			if err := record.Validate(); !RefusalHasCode(err, RefusalInvalidKnownField) {
				t.Fatalf("Validate %s status %q error = %v, want invalid_known_field", test.recordType, test.status, err)
			}
		})
	}
}

func TestRecordValidateRejectsInvalidProgrammaticReferenceID(t *testing.T) {
	t.Parallel()

	record := Record{
		Type: Mission,
		ID:   mustID(t, missionIDText),
		Source: &Reference{
			Type: Proposal,
			ID:   ID("not-a-uuid"),
		},
	}
	if err := record.Validate(); !RefusalHasCode(err, RefusalInvalidReference) {
		t.Fatalf("Validate error = %v, want invalid_reference", err)
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
