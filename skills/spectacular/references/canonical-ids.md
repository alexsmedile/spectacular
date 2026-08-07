---
description: UUIDv7 identity, slug, alias, and cross-reference contract for durable entities.
when_to_use: Creating, resolving, linking, renaming, or migrating decisions, questions, ideas, research, spikes, prototypes, specs, or reserved tasks.
---

# Canonical IDs

New durable records use UUIDv7 as the stable machine reference. Their slug is the human-facing locator and filename; it can change without changing the ID. Markdown plus YAML frontmatter remains the canonical portable store.

```yaml
id: 019fdce6-72b0-7108-8024-b42a68765a85 # immutable UUIDv7
slug: durable-identity-contract-workflow    # human locator and filename
kind: spec
scope: project
status: draft
summary: "..."
related: []
references: []
```

Legacy numeric IDs are read-only aliases after an explicit `spectacular id migrate --uuidv7 --apply --yes`; newly created records never receive a counter alias. `spectacular id resolve <slug|UUIDv7|legacy-alias> --context <entity>` returns the UUIDv7 to persist in cross-references.

The historical table below applies only while old records have not migrated.

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
| `FND-001` | finding (reserved) | `fnd1` | `FND1` |
| `FIX-001` | fix/remediation (reserved successor) | `fix1`; before migration, `fix1` resolves legacy `F1` | `FIX1` |
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
8. Use `fix1` for fixes and `fnd1` for findings. Context-free `f1` is refused because it hides which entity the user meant. Before an explicit preview-first, archive-first migration activates `FIX-NNN`, `fix1` resolves the existing legacy fix `F1`; explicit canonical `FIX-001` remains reserved.
9. Roadmap build IDs own `b<N>`. Bugs therefore use `bug1`; context-free `b1` is never a bug alias.
10. `RES` is the only research prefix. `RCH`, `RSC`, `SER`, and `SRC` are not aliases; use conversational `R1` when brevity matters.
11. `PRT` stays reserved. Prototype artifacts are owned by a request, vision, feedback entry, or technical `SPK` until a standalone lifecycle is proven necessary.
12. “Artifact,” “tracer bullet,” and “technical debt” are concepts, not canonical entity prefixes. Do not invent `ART`, `TRC`, or `DEB` records.

Use `spectacular id resolve <ref> [--context <entity>]` for deterministic normalization.
