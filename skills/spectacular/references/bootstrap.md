# Manual bootstrap

Use this when: Operator with explicitly approved manual-bootstrap authority drafting a future record shape that the current CLI cannot represent.

In that mode the CLI cannot mutate or validate the Mission. Treat it as out of
service. Draft the proposed Markdown outside the governed lifecycle, normally in
`scratch/`, and verify only the structure you can observe.

For the field shape the future record must keep, see
[mission-anatomy.md](mission-anatomy.md); compare one existing
`.spectacular/missions/*/*.md` without modifying it.

## Verify by hand

Work down the list. Each item is a separate check.

- **YAML** parses, and every key is one the schema knows.
- **Identity placeholders** are clearly marked for later CLI allocation; never invent a UUIDv7.
- **Refs** resolve. No pointer to a missing file, Objective, or Run.
- **Fingerprint fields** are absent or clearly marked for later CLI computation; never invent or claim to verify one.
- **Claim coverage** is complete: every completion claim maps to an Objective, and
  every Objective serves at least one claim.
- **Dependency DAG** is acyclic, and each named dependency exists.
- **Run state** is legal for the Mission status, with exactly one current Run.
- **Authority intent** is drafted: proposed owner and forbidden-effect ceiling; no activation time or authority is claimed.
- **Scope** is stated on both axes, mechanical and semantic.
- **File layout** matches the canonical shape. The root Mission record is the index;
  `objectives/` and `runs/` exist only where they were earned.

## Rules

- Manual bootstrap cannot create, activate, transition, or complete a fingerprint-bound governed record. It produces a draft for a later compatible CLI or an explicitly governed implementation Mission.
- Never cite the old CLI as proof of anything. An incompatible CLI refusal is a
  tooling Gap, not evidence that the Mission is invalid.
- Prefer a focused script over a manual read when the check is exact and you will
  repeat it.
- Do not route the work through a legacy command sequence to make it look
  validated. Preserve the draft and record the incompatible CLI as a tooling Gap.
