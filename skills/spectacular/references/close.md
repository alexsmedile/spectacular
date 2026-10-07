# Close & Review Micro-Kernel

Use this when: the primary operator or reviewer is verifying claims, collecting receipts, or completing a Mission.

## 2. CLI Palette & Completion
```bash
spectacular mission check <ref> --json               # Read-only claim validation
spectacular mission complete <ref> --by <owner> --json # Final completion gate
spectacular review record <ref> review.md --json     # Record review required by the Mission
```

## 3. Zero-Sprawl Verification Policy
- **Frozen claims**: Account for every `pass_boundary` and `proof_requirement`
  against inspected, attributable evidence. Passing tests prove only what they
  cover; a clean commit does not prove claim coverage. Unproved claims and
  unresolved blockers prevent Mission completion.
- **Frozen proof and review requirements**: The selected Mission is authoritative.
  Record its specified Evidence and independent Reviews even for routine work or
  local observations. Metadata profiles do not waive those requirements.
- **When the Mission leaves record choices open**: Prefer the smallest useful
  footprint. Separate Evidence preserves attributable receipts or disputed
  observations; separate Reviews serve requested independent evaluation. Risk
  alone does not activate this governed route.

## 4. Reviewer Hygiene (Observe ≠ Act)
- Reviewers inspect diffs and test logs; they **NEVER edit files or fix bugs directly**.
- All defects are returned to the Orchestrator as structured findings (`pass | repair | owner-gate`).

## 5. Negative Constraints (DO NOT)
- **DO NOT** create extra review/evidence folders without a frozen requirement or independently useful proof need.
- **DO NOT** complete a Mission with failing post-checks or unaddressed blockers.
- **DO NOT** echo completed mission YAML into chat; return 1-line confirmation.
