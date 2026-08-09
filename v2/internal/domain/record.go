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
	for name, value := range map[string]*string{
		"title":      r.Title,
		"status":     r.Status,
		"created_by": r.CreatedBy,
	} {
		if value != nil && strings.TrimSpace(*value) == "" {
			return NewRefusal(RefusalInvalidKnownField, name, "must not be empty", nil)
		}
	}
	for name, value := range map[string]*string{
		"created": r.Created,
		"updated": r.Updated,
	} {
		if value == nil {
			continue
		}
		if _, err := time.Parse(time.RFC3339, *value); err != nil {
			return NewRefusal(RefusalInvalidKnownField, name, "must be RFC3339", err)
		}
	}
	if recordType == Proposal && r.Source != nil {
		return NewRefusal(RefusalInvalidKnownField, "source", "Proposal grammar has no source relationship", nil)
	}
	if recordType == Mission && r.Source != nil && r.Source.Type != Proposal {
		return NewRefusal(
			RefusalInvalidReference,
			"source",
			fmt.Sprintf("Mission source must target Proposal, got %s", r.Source.Type),
			nil,
		)
	}
	return nil
}

// KnownFieldOrder is the canonical top-level ordering for semantic fields.
func KnownFieldOrder(recordType RecordType) []string {
	fields := []string{
		"type",
		"id",
		"title",
		"description",
		"status",
		"created_by",
		"created",
		"updated",
	}
	if recordType == Mission {
		fields = append(fields, "source")
	}
	return fields
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
