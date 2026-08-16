package missionbundle

import (
	"sort"

	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// Flag names one observed reason a claim may deserve attention. Flags are named
// rather than scored: a score defaults cleanly but cannot be argued with, and
// owner disagreement is the control that makes an aimed audit safe to default.
type Flag string

const (
	// FlagUnproven means no review has returned a verdict for the claim.
	FlagUnproven Flag = "unproven"
	// FlagFailed means a review returned a verdict other than pass.
	FlagFailed Flag = "failed"
	// FlagStaleReview means the review that covered the claim was bound to an
	// activation fingerprint that is no longer current, so its verdict
	// describes a Mission that has since been re-frozen.
	FlagStaleReview Flag = "stale-review"
	// FlagRepaired means the Run consumed repair budget. It applies to every
	// claim because the Run does not record which claim a repair served.
	FlagRepaired Flag = "repaired"
	// FlagBudgetExhausted means repairs have reached the frozen budget.
	FlagBudgetExhausted Flag = "budget-exhausted"
	// FlagUnimplemented means no Objective claiming this has been implemented.
	FlagUnimplemented Flag = "unimplemented"
	// FlagUnclaimed means no Objective claims it at all, so nothing in the plan
	// will produce evidence for it.
	FlagUnclaimed Flag = "unclaimed"
)

// ClaimDrift is the derived audit signal for one frozen completion claim.
type ClaimDrift struct {
	Claim string `json:"claim"`
	Flags []Flag `json:"flags,omitempty"`
	// Objectives lists the Objective refs that claim this, in plan order.
	Objectives []string `json:"objectives,omitempty"`
	// Verdict is the most recent recorded verdict, empty when unreviewed.
	Verdict string `json:"verdict,omitempty"`
}

// Drift computes per-claim flags from inputs the bundle already carries:
// repairs consumed against budget, review verdicts, the activation fingerprint
// a review was bound to, and which Objectives claim what.
//
// The result is ordered most-flagged first so the head of the list is the
// default audit target. Ties break on plan order, which is stable.
func (b *Bundle) Drift() []ClaimDrift {
	implementedClaims := map[string]bool{}
	claimObjectives := map[string][]string{}
	for _, objective := range b.Objectives {
		for _, claim := range objective.Claims {
			claimObjectives[claim] = append(claimObjectives[claim], objective.Ref)
			if objective.Status == "implemented" {
				implementedClaims[claim] = true
			}
		}
	}

	verdicts := map[string]string{}
	stale := map[string]bool{}
	current := ""
	if b.Activation != nil {
		current = b.Activation.Fingerprint
	}
	for _, pointer := range b.Reviews {
		if pointer.Document == nil {
			continue
		}
		outdated := current != "" && pointer.Document.Reviewed.ActivationFingerprint != "" &&
			pointer.Document.Reviewed.ActivationFingerprint != current
		for _, claim := range pointer.Document.Claims {
			verdicts[claim.Claim] = claim.Verdict
			if outdated {
				stale[claim.Claim] = true
			}
		}
	}

	repairs, budget := 0, b.RepairBudget
	if b.Run != nil {
		repairs = b.Run.Repairs
	}

	drift := make([]ClaimDrift, 0, len(b.Completion))
	for _, criterion := range b.Completion {
		item := ClaimDrift{
			Claim:      criterion.Claim,
			Objectives: claimObjectives[criterion.Claim],
			Verdict:    verdicts[criterion.Claim],
		}
		switch verdict, reviewed := verdicts[criterion.Claim]; {
		case !reviewed:
			item.Flags = append(item.Flags, FlagUnproven)
		case verdict != "pass":
			item.Flags = append(item.Flags, FlagFailed)
		}
		if stale[criterion.Claim] {
			item.Flags = append(item.Flags, FlagStaleReview)
		}
		if len(claimObjectives[criterion.Claim]) == 0 {
			item.Flags = append(item.Flags, FlagUnclaimed)
		} else if !implementedClaims[criterion.Claim] {
			item.Flags = append(item.Flags, FlagUnimplemented)
		}
		if repairs > 0 {
			item.Flags = append(item.Flags, FlagRepaired)
		}
		if budget > 0 && repairs >= budget {
			item.Flags = append(item.Flags, FlagBudgetExhausted)
		}
		drift = append(drift, item)
	}

	// Stable sort keeps plan order among equally-flagged claims, so the default
	// audit target does not move between runs on an arbitrary tie.
	sort.SliceStable(drift, func(i, j int) bool {
		return len(drift[i].Flags) > len(drift[j].Flags)
	})
	return drift
}

// AuditTarget names the claim an unnamed audit should aim at, with the flags
// that selected it. It returns false when nothing is flagged, so a caller asks
// the owner rather than inventing a target.
func (b *Bundle) AuditTarget() (ClaimDrift, bool) {
	drift := b.Drift()
	if len(drift) == 0 || len(drift[0].Flags) == 0 {
		return ClaimDrift{}, false
	}
	return drift[0], true
}

// Notices reports observations that are worth telling a reader about but must
// not fail validation. The legacy `human_ref:` spelling is the current case:
// the rename to `ref:` stopped halfway, and rewriting a frozen record's
// frontmatter to finish it would change fingerprints for a cosmetic reason.
func (b *Bundle) Notices() []string {
	if b.document == nil {
		return nil
	}
	var notices []string
	if _, legacy, err := workspace.Ref(b.document); err == nil && legacy {
		notices = append(notices, "ref-spelling-drift: uses legacy `human_ref:`; new records declare `ref:`")
	}
	return notices
}
