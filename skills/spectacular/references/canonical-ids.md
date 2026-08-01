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
| `FND-001` | finding (reserved) | `fnd1`; `f1` unavailable while legacy fixes use `F<N>` | `FND1` |
| `FIX-001` | fix/remediation (reserved successor) | `fx1`, `fix1`; current fix records use `f1` → `F1` | `F1` |
| `BUG-001` | bug/defect (reserved) | `bug1`; `b1` is unavailable because it is a roadmap build | `BUG1` |
| `SEC-001` | security vulnerability (reserved) | `sec1` | `SEC1` |
| `BMK-001` | benchmark result (reserved) | `bm1`, `bmk1` | `BMK1` |

Rules:

1. Explicit prefixes win.
2. Naked numbers require an entity context; non-interactive ambiguity refuses rather than guesses.
3. IDs use at least three digits. IDs above 999 render with four digits. A future migration may widen earlier IDs when real scale requires it; creation must not silently bulk-rewrite a workspace.
4. Cross-references store canonical IDs, even when the user supplied an alias.
5. `IDEA` is canonical; `IDE` is accepted only for compatibility.
6. Renames change the descriptive slug, never the ID.
7. Reserved advanced prefixes define future identity and storage only. They are not allocated or resolved until their collection workflow ships.
8. Existing fixes use `F<N>` and own the short alias `f1`. Findings therefore use `fnd1`. Activating `FIX-NNN` requires an explicit preview-first, archive-first migration; no resolver may guess across that collision.
9. Roadmap build IDs own `b<N>`. Bugs therefore use `bug1`; context-free `b1` is never a bug alias.
10. `RES` is the only research prefix. `RCH`, `RSC`, `SER`, and `SRC` are not aliases; use conversational `R1` when brevity matters.
11. `PRT` stays reserved. Prototype artifacts are owned by a request, vision, feedback entry, or technical `SPK` until a standalone lifecycle is proven necessary.
12. “Artifact,” “tracer bullet,” and “technical debt” are concepts, not canonical entity prefixes. Do not invent `ART`, `TRC`, or `DEB` records.

Use `spectacular id resolve <ref> [--context <entity>]` for deterministic normalization.
