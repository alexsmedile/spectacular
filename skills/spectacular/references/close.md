# Close & Review Micro-Kernel

## 1. Trigger Context
Primary operator or reviewer verifying claims, collecting receipts, or completing a Mission.

## 2. CLI Palette & Completion
```bash
spectacular mission check <ref> --json               # Read-only claim validation
spectacular mission complete <ref> --by <owner> --json # Final completion gate
spectacular review record <ref> review.md --json     # Record dedicated review (high-risk only)
```

## 3. Zero-Sprawl Verification Policy
- **Routine Tasks**: Passing test suite (`exit 0`) + clean Git commit **is** the primary proof.
- **Separate Evidence (`evidence/`)**: Created only when third-party provider receipts or disputed behaviors must be preserved.
- **Separate Review (`reviews/`)**: Required only for high-risk operations (`mode: control`, payments, auth, DB migrations).

## 4. Reviewer Hygiene (Observe ≠ Act)
- Reviewers inspect diffs and test logs; they **NEVER edit files or fix bugs directly**.
- All defects are returned to the Orchestrator as structured findings (`pass | repair | owner-gate`).

## 5. Negative Constraints (DO NOT)
- **DO NOT** create multi-file review/evidence folders for routine features.
- **DO NOT** complete a Mission with failing post-checks or unaddressed blockers.
- **DO NOT** echo completed mission YAML into chat; return 1-line confirmation.
