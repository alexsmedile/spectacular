---
description: Canonical identity, alias, padding, and cross-reference contract for Wayfinding entities.
when_to_use: Creating, resolving, linking, renaming, or migrating decisions, questions, ideas, research, spikes, prototypes, specs, or reserved tasks.
---

# Canonical IDs

Canonical IDs are stable machine references. Conversational aliases are input conveniences only; frontmatter always persists the canonical form.

| Canonical | Entity | Accepted aliases | Preferred speech |
|---|---|---|---|
| `DEC-001` | decision | `d1`, `dec1`, `DEC-001` | `D1` |
| `QUE-001` | question | `q1`, `que1`, `QUE-001` | `Q1` |
| `IDEA-001` | idea | `i1`, `ide1`, `idea1`, legacy `IDE-001` | `I1` |
| `RES-001` | research | `r1`, `res1` | `R1` |
| `SPK-001` | spike | `spk1` | `SPK1` |
| `PRT-001` | prototype | `prt1` | `PRT1` |
| `SPC-001` | specification | `s1`, `spc1`, legacy `spec1` | `S1` |
| `TSK-001` | task (reserved) | `t1`, `tsk1` | `T1` |

Rules:

1. Explicit prefixes win.
2. Naked numbers require an entity context; non-interactive ambiguity refuses rather than guesses.
3. IDs use at least three digits. IDs above 999 render with four digits. A future migration may widen earlier IDs when real scale requires it; creation must not silently bulk-rewrite a workspace.
4. Cross-references store canonical IDs, even when the user supplied an alias.
5. `IDEA` is canonical; `IDE` is accepted only for compatibility.
6. Renames change the descriptive slug, never the ID.

Use `spectacular id resolve <ref> [--context <entity>]` for deterministic normalization.
