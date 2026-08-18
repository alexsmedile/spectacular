# Manual bootstrap

Load this only when the Mission declares `manual-bootstrap` as its validation mode.

In that mode the CLI cannot mutate or validate the Mission. Treat it as out of
service. Edit the canonical Markdown directly and verify the result yourself.

The file you are editing is `.spectacular/missions/<slug>/<ref>-<slug>.md`. For the
field shape it must keep, see [mission-anatomy.md](mission-anatomy.md); for a live
one, read any existing `.spectacular/missions/*/*.md`.

## Verify by hand

Work down the list. Each item is a separate check.

- **YAML** parses, and every key is one the schema knows.
- **UUIDv7 identity** is present, well-formed, and unchanged from before the edit.
- **Refs** resolve. No pointer to a missing file, Objective, or Run.
- **Fingerprints** match the frozen semantic envelope. Mutable progress is
  excluded from them.
- **Claim coverage** is complete: every completion claim maps to an Objective, and
  every Objective serves at least one claim.
- **Dependency DAG** is acyclic, and each named dependency exists.
- **Run state** is legal for the Mission status, with exactly one current Run.
- **Authority** is recorded: owner, activation time, forbidden-effect ceiling.
- **Scope** is stated on both axes, mechanical and semantic.
- **File layout** matches the canonical shape. The root Mission record is the index;
  `objectives/` and `runs/` exist only where they were earned.

## Rules

- Never cite the old CLI as proof of anything. An incompatible CLI refusal is a
  tooling Gap, not evidence that the Mission is invalid.
- Prefer a focused script over a manual read when the check is exact and you will
  repeat it.
- Do not route the work through a legacy command sequence to make it look
  validated. Create the same canonical shape directly instead.
