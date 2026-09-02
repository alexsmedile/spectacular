# Audit & FROST Micro-Kernel

## 1. Trigger Context
Reviewer or bounded auditor conducting retrospective claim or proof challenges.

## 2. CLI Palette & Inspection
```bash
spectacular mission show <ref> --json    # Inspect frozen claims & state
spectacular mission check <ref> --json   # Validate mechanical integrity
git log --oneline <baseline>..HEAD       # Inspect actual diffs
```

## 3. Reviewer 4-Point Rubric & FROST (Observe ≠ Act)
- **FROST Evaluation**: Check `Frozen fit` (claims vs non-goals), Risk, Operability, System integrity, and `Truth of proof` (attributable un-fabricated receipts).
- **4-Point Rubric**:
  1. **Decision Compliance (`D<N>`)**: Does the diff violate any locked Decisions or STACK constraints?
  2. **Atlas State Coverage**: If state machines exist in `atlas/`, does code implement all declared transitions?
  3. **Claim Fidelity**: Do deliverables satisfy the frozen completion checklist in `M<N>.md`?
  4. **Proof Validity**: Are test receipts genuine and reproducible (`exit 0`)? Zero fabricated claims.

## 4. Completed Mission Contract Freeze
- On completed missions, `contract-drift` is an informational notice, not an error.
- Never re-point a completed mission's contract fingerprint.

## 5. Negative Constraints (DO NOT)
- **DO NOT** edit code or fix bugs during an audit. Report findings to Orchestrator.
- **DO NOT** widen audit scope beyond the named claim without explicit instruction.
- **DO NOT** read catalog caches or projections; inspect canonical records via CLI.
