package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// ID is a canonical lowercase UUIDv7 string.
type ID string

// NewID creates a new canonical UUIDv7 identity.
func NewID() (ID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return "", NewRefusal(RefusalInvalidID, "id", "generate UUIDv7", err)
	}
	return ID(value.String()), nil
}

// ParseID validates both the UUID version and the canonical textual form.
func ParseID(raw string) (ID, error) {
	value, err := uuid.Parse(raw)
	if err != nil {
		return "", NewRefusal(RefusalInvalidID, "id", "parse canonical UUIDv7", err)
	}
	if value.Version() != 7 {
		return "", NewRefusal(
			RefusalInvalidID,
			"id",
			fmt.Sprintf("UUID version %d is not version 7", value.Version()),
			nil,
		)
	}
	if value.String() != raw {
		return "", NewRefusal(RefusalInvalidID, "id", "UUIDv7 is not in canonical lowercase form", nil)
	}
	return ID(raw), nil
}

func (id ID) String() string {
	return string(id)
}
