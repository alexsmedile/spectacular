package domain

import "strings"

// RecordType is the semantic grammar discriminator for an M1 record.
type RecordType string

const (
	Proposal   RecordType = "Proposal"
	Mission    RecordType = "Mission"
	Anchor     RecordType = "Anchor"
	Gap        RecordType = "Gap"
	Run        RecordType = "Run"
	Checkpoint RecordType = "Checkpoint"
	Evidence   RecordType = "Evidence"
	Decision   RecordType = "Decision"
)

func ParseRecordType(raw string) (RecordType, error) {
	switch RecordType(raw) {
	case Proposal, Mission, Anchor, Gap, Run, Checkpoint, Evidence, Decision:
		return RecordType(raw), nil
	default:
		return "", NewRefusal(RefusalInvalidType, "type", "unknown record noun", nil)
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
