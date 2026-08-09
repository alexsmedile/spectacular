package domain

import "strings"

// RecordType is the semantic grammar discriminator for an M1 record.
type RecordType string

const (
	Proposal RecordType = "Proposal"
	Mission  RecordType = "Mission"
)

func ParseRecordType(raw string) (RecordType, error) {
	switch RecordType(raw) {
	case Proposal, Mission:
		return RecordType(raw), nil
	default:
		return "", NewRefusal(RefusalInvalidType, "type", "expected Proposal or Mission", nil)
	}
}

// Reference is an exact typed identity relationship.
type Reference struct {
	Type RecordType
	ID   ID
}

func ParseReference(raw string) (Reference, error) {
	typeText, idText, found := strings.Cut(raw, ":")
	if !found || typeText == "" || idText == "" || strings.Contains(idText, ":") {
		return Reference{}, NewRefusal(
			RefusalInvalidReference,
			"source",
			"expected <type>:<canonical-UUIDv7>",
			nil,
		)
	}
	recordType, err := ParseRecordType(typeText)
	if err != nil {
		return Reference{}, NewRefusal(RefusalInvalidReference, "source", "unknown target type", err)
	}
	id, err := ParseID(idText)
	if err != nil {
		return Reference{}, NewRefusal(RefusalInvalidReference, "source", "invalid target identity", err)
	}
	return Reference{Type: recordType, ID: id}, nil
}

func (r Reference) String() string {
	return string(r.Type) + ":" + r.ID.String()
}
