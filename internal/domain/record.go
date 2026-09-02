package domain

import (
	"fmt"
	"strings"
	"time"
)

// Record contains the authoritative values from the two M1 grammars. Optional
// pointers preserve the distinction between an absent property and a present
// empty value, which validation may refuse.
type Record struct {
	Type        RecordType
	ID          ID
	Title       *string
	Description *string
	Status      *string
	CreatedBy   *string
	Created     *string
	Updated     *string
	Source      *Reference
}

// Validate applies identity, grammar, and static known-field checks. It does
// not perform lifecycle transitions or authorize any effect.
func (r Record) Validate() error {
	recordType, err := ParseRecordType(string(r.Type))
	if err != nil {
		return err
	}
	if _, err := ParseID(r.ID.String()); err != nil {
		return err
	}
	for _, field := range []struct {
		name  string
		value *string
	}{
		{name: "title", value: r.Title},
		{name: "created_by", value: r.CreatedBy},
	} {
		if field.value != nil && strings.TrimSpace(*field.value) == "" {
			return NewRefusal(RefusalInvalidKnownField, field.name, "must not be empty", nil)
		}
	}
	for _, field := range []struct {
		name  string
		value *string
	}{
		{name: "created", value: r.Created},
		{name: "updated", value: r.Updated},
	} {
		if field.value == nil {
			continue
		}
		if _, err := time.Parse(time.RFC3339, *field.value); err != nil {
			return NewRefusal(RefusalInvalidKnownField, field.name, "must be RFC3339", err)
		}
	}
	if r.Status != nil && strings.TrimSpace(*r.Status) == "" {
		return NewRefusal(RefusalInvalidKnownField, "status", "must not be empty", nil)
	}
	if r.Status != nil && (recordType == Proposal || recordType == Mission) && !validStatus(recordType, *r.Status) {
		return NewRefusal(
			RefusalInvalidKnownField,
			"status",
			fmt.Sprintf("%s status must be one of: %s", recordType, allowedStatuses(recordType)),
			nil,
		)
	}
	if recordType == Proposal && r.Source != nil {
		return NewRefusal(RefusalInvalidKnownField, "source", "Proposal grammar has no source relationship", nil)
	}
	if recordType == Mission && r.Source != nil {
		validated, err := ParseReference(r.Source.String())
		if err != nil {
			return err
		}
		if validated.Type != Proposal {
			return NewRefusal(
				RefusalInvalidReference,
				"source",
				fmt.Sprintf("Mission source must target Proposal, got %s", validated.Type),
				nil,
			)
		}
	}
	if recordType != Mission && recordType != Proposal && r.Source != nil {
		return NewRefusal(RefusalInvalidKnownField, "source", "recovery records use explicit typed relationship fields", nil)
	}
	return nil
}

func validStatus(recordType RecordType, status string) bool {
	switch recordType {
	case Proposal:
		switch status {
		case "draft", "submitted", "accepted", "rejected", "withdrawn":
			return true
		}
	case Mission:
		switch status {
		case "draft", "defined", "active", "awaiting-assessment", "resolved", "completed", "superseded", "withdrawn":
			return true
		}
	}
	return false
}

func allowedStatuses(recordType RecordType) string {
	if recordType == Proposal {
		return "draft, submitted, accepted, rejected, withdrawn"
	}
	return "draft, defined, active, awaiting-assessment, resolved, completed, superseded, withdrawn"
}

// IsReservedField identifies names whose interpretation belongs to an M1
// grammar. A field cannot evade validation by appearing on the wrong type.
func IsReservedField(name string) bool {
	switch name {
	case "type", "id", "title", "description", "status", "created_by", "created", "updated", "source":
		return true
	default:
		return false
	}
}
