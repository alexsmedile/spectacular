package missionbundle

import (
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// contractGap is one entry in a Contract's `gaps:` block. Only the ref and the
// two terminal-state fields are read: a Gap is closed by rewriting `blocked_on:`
// to `resolution:`, never by deleting the entry, so knowing which of the two a
// Gap currently carries is what tells an amendment whether there is anything to
// do. The problem statement is prose for a reader and is not interpreted here.
type contractGap struct {
	Ref        string `yaml:"ref"`
	Problem    string `yaml:"problem"`
	BlockedOn  string `yaml:"blocked_on"`
	Resolution string `yaml:"resolution"`
}

// ContractGaps reads the `gaps:` block of a bound Contract. A Contract with no
// Gaps is normal and returns nothing rather than refusing: most Contracts carry
// none, and an absent block is silence rather than a defect.
func ContractGaps(ws *discovery.Workspace, ref string) ([]contractGap, error) {
	entry, err := ws.Lookup(ref, domain.Contract)
	if err != nil {
		return nil, err
	}
	if entry.Document == nil || entry.Document.Unknown["gaps"] == nil {
		return nil, nil
	}
	var gaps []contractGap
	if err := workspace.DecodeValue(entry.Document, "gaps", &gaps); err != nil {
		return nil, err
	}
	return gaps, nil
}

// validateResolvedGaps holds a `resolves_gaps:` declaration to the Contract it
// names. The refs are checked at plan-freeze rather than at completion because a
// Mission that discovers at the gate that it never had the authority to close a
// Gap has already done the work under a false premise.
func validateResolvedGaps(ws *discovery.Workspace, b *Bundle) error {
	if len(b.ResolvesGaps) == 0 {
		return nil
	}
	gaps, err := ContractGaps(ws, b.Contract.Ref)
	if err != nil {
		return err
	}
	available := make(map[string]bool, len(gaps))
	names := make([]string, 0, len(gaps))
	for _, gap := range gaps {
		if gap.Ref == "" {
			continue
		}
		available[gap.Ref] = true
		names = append(names, gap.Ref)
	}
	sort.Strings(names)
	declared := strings.Join(names, ", ")
	if declared == "" {
		declared = "the bound Contract declares no Gaps"
	}

	seen := make(map[string]bool, len(b.ResolvesGaps))
	for _, resolved := range b.ResolvesGaps {
		ref := strings.TrimSpace(resolved.Gap)
		if ref == "" {
			return invalid("resolves_gaps.gap", "each entry must name a Gap ref on the bound Contract")
		}
		if strings.TrimSpace(resolved.Resolution) == "" {
			return invalid("resolves_gaps.resolution",
				"Gap "+ref+" needs the resolution text this Mission will write; it is frozen so the owner approves the wording")
		}
		if seen[ref] {
			return invalid("resolves_gaps.gap", "Gap declared twice: "+ref)
		}
		seen[ref] = true
		if !available[ref] {
			return domain.NewStateRefusal(domain.RefusalInvalidKnownField, "resolves_gaps.gap",
				"no such Gap on the bound Contract: "+ref, ref, declared,
				"name a Gap the bound Contract declares, or drop the entry", nil)
		}
	}
	return nil
}

// DeclaredNotices are the observation names a Contract may declare in
// mandatory_validation. A notice reports and never refuses, so it is not a
// validator — CC-projsurf declares ref-spelling-drift alongside its validators and
// annotates it "legacy human_ref reported, not refused", which is a notice
// described accurately and filed in the wrong list.
var DeclaredNotices = []string{"ref-spelling-drift"}

// ProposalValidations are check names satisfied by the Proposal validator rather
// than the Mission registry. A Contract that governs Proposal shape declares them,
// and they run through `proposal check`.
var ProposalValidations = []string{"proposal-schema-v2"}

// ResolveDeclaredValidation reports whether a name a Contract declares in
// mandatory_validation corresponds to something that actually runs. A declaration
// nothing verifies is the same defect as a Gap left reading blocked_on after it was
// resolved: the record promises what the system does not enforce, and nothing
// catches the divergence.
func ResolveDeclaredValidation(name string) (kind string, ok bool) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", false
	}
	for _, validator := range registry {
		if validator.name == trimmed {
			return "validator", true
		}
	}
	for _, notice := range DeclaredNotices {
		if notice == trimmed {
			return "notice", true
		}
	}
	for _, proposal := range ProposalValidations {
		if proposal == trimmed {
			return "proposal-validator", true
		}
	}
	return "", false
}
