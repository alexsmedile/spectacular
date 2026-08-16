package missionbundle

import (
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// ProposalCheck is the validation result for one Proposal record. Proposals are
// authored as Markdown and checked; there is no creation command, so this
// validates a shape a human wrote rather than one the system generated.
type ProposalCheck struct {
	Ref     string   `json:"ref"`
	ID      string   `json:"id"`
	Path    string   `json:"path"`
	Valid   bool     `json:"valid"`
	Checks  []string `json:"checks"`
	Notices []string `json:"notices,omitempty"`
}

// proposalFields is the compact schema: the shape P5 and P6 are actually
// authored in. Legacy records carry more, and their extra fields are preserved
// rather than refused.
var proposalRequiredFields = []string{"type", "title", "status", "created_by", "created", "updated", "target_contract"}

// ValidateProposal checks one Proposal against the compact schema.
//
// It reports rather than refuses wherever a frozen record would otherwise have
// to be rewritten: the legacy `human_ref:` spelling, and the superseded
// candidate_* body that P2 through P4 carry. A Proposal proposes; the Contract
// it targets holds the frozen result.
func ValidateProposal(ws *discovery.Workspace, ref string) (ProposalCheck, error) {
	entry, err := ws.Lookup(ref, domain.Proposal)
	if err != nil {
		return ProposalCheck{}, err
	}
	doc := entry.Document
	check := ProposalCheck{Path: entry.Path, Valid: true}

	resolved, legacy, err := workspace.Ref(doc)
	if err != nil {
		return ProposalCheck{}, err
	}
	check.Ref = resolved
	check.ID = doc.Record.ID.String()
	if resolved == "" {
		return ProposalCheck{}, domain.NewRefusal(domain.RefusalMissingRequiredField, workspace.RefField,
			"a Proposal must declare ref (legacy human_ref is accepted)", nil)
	}
	check.Checks = append(check.Checks, "ref-present")
	if legacy {
		check.Notices = append(check.Notices,
			"ref-spelling-drift: uses legacy `human_ref:`; new records declare `ref:`")
	}

	for _, field := range proposalRequiredFields {
		// Most of the compact schema lives on the typed Record rather than in
		// the Unknown map, so absence from Unknown is not absence.
		if typedProposalField(doc, field) {
			continue
		}
		if _, present := doc.Unknown[field]; present {
			continue
		}
		return ProposalCheck{}, domain.NewRefusal(domain.RefusalMissingRequiredField, field,
			"the compact Proposal schema requires this property", nil)
	}
	check.Checks = append(check.Checks, "compact-schema")
	check.Checks = append(check.Checks, "uuidv7-identity", "status-vocabulary")

	if target, targetErr := workspace.String(doc, "target_contract", false); targetErr == nil && target != "" {
		typed, parseErr := domain.ParseReference(target)
		if parseErr != nil || typed.Type != domain.Contract {
			return ProposalCheck{}, domain.NewRefusal(domain.RefusalInvalidKnownField, "target_contract",
				"must be Contract:<UUIDv7>", parseErr)
		}
		check.Checks = append(check.Checks, "contract-reference")
	}

	if _, present := doc.Unknown["candidate_purpose"]; present {
		check.Notices = append(check.Notices,
			"legacy-candidate-body: carries superseded candidate_* fields; a Proposal proposes and the Contract holds the frozen result")
	}
	for _, superseded := range []string{"idempotency_key", "authorization", "base_fingerprint"} {
		if _, present := doc.Unknown[superseded]; present {
			check.Notices = append(check.Notices,
				"legacy-ceremony: carries superseded `"+superseded+"`; the compact schema does not require it")
		}
	}
	return check, nil
}

// typedProposalField reports whether a required property is carried on the
// typed Record rather than in the Unknown map.
func typedProposalField(doc *workspace.Document, field string) bool {
	switch field {
	case "type":
		return doc.Record.Type != ""
	case "title":
		return doc.Record.Title != nil
	case "status":
		return doc.Record.Status != nil
	case "created_by":
		return doc.Record.CreatedBy != nil
	case "created":
		return doc.Record.Created != nil
	case "updated":
		return doc.Record.Updated != nil
	}
	return false
}
