---
status: active
updated: 2026-06-28
related:
  - PLAN.md
  - TASKS.md
---

# Verify log — roadmap-contract-docs

Docs/spec-only request (zero behavior change), so verification is read-and-confirm
of the written artifacts + doctor integrity, not executable tests.

## M1 — Spec the ledger

- **{assert}** `.spectacular/specs/roadmap/SPEC.md` snapshotted before edit → `snapshots/specs/roadmap/SPEC/@v1.md` exists; live spec frontmatter is now `status: published`, `version: 1.1`, `updated: 2026-06-28`. ✅
- **{assert}** Spec now carries a "The ledger (build → version)" section covering build ids, target-version single-source, `tbd`, and status-vs-lifecycle; opens with the two-layer (ledger + prose) framing. ✅
- **{assert}** `.spectacular/SPEC.md` structured-ROADMAP bullet names the ledger + `tbd` and links `[[specs/roadmap/SPEC]]`, `[[ARCHITECTURE]]`, `[[roadmap-rules]]`. Stale `[[roadmap-overrides]]` link removed (in SPEC.md bullet + spec body). ✅
- **{judge}** M-question 1 resolved: ARCHITECTURE.md remains the **canonical** schema home; the spec *summarizes + points to it* (explicit "does not fork a second authoritative copy" note) — no drift risk. ✅
- **`run: ./cli/spectacular doctor specs`** → specs area green (roadmap + doc-engine capability specs present, SPEC.md parses). ✅

## M2 — Define `tbd` + fix contradicting rule

- **{assert}** `ARCHITECTURE.md` `target-version` row now lists `v1.10.0 · tbd` and documents `tbd` as the slotted-but-not-pinned sentinel, explicitly distinct from a `<TBD>` placeholder. ✅
- **{assert}** `roadmap-rules.md` placeholder check rescoped: "in any **prose slot**" + a blockquote "Ledger `tbd` is not a placeholder" + a new "The ledger (build → version)" section with the behavioral `tbd` rule. ✅
- **{judge}** Contradiction resolved without weakening the real check: prose `<TBD>`/`<TODO>`/`???` still rejected (rule text intact, roadmap-rules.md:371); ledger-column `tbd` explicitly exempt. This repo's own ROADMAP ledger uses `tbd` on 7 rows and `doctor` is green — live proof the gate no longer false-positives. ✅

## M3 — User docs + tutorial

- **{assert}** `docs/versioning.md` gains "The roadmap ledger — how builds map to versions" (the tutorial/walkthrough: build-id model, worked ledger example, `tbd`, status-vs-lifecycle, why-build-ids-not-versions). Placed between "Choosing the next version" and "single canonical version source". ✅
- **{assert}** `docs/configuration.md` documents `last_build:` (new section, v1.17.0+) with cross-link to the versioning ledger section. ✅
- **{assert}** `docs/commands.md` gains a `spectacular roadmap` subsection (two-layer model + `tbd` + links to versioning.md and ARCHITECTURE). ✅
- **{judge}** M-question 2 resolved: tutorial folded into `docs/versioning.md` (the natural home — it's how versions get assigned) rather than a new standalone page; commands.md + configuration.md link to it. Single source, no duplication. ✅

## Doctor / integrity

- **`run: ./cli/spectacular doctor specs links docs`** → 0 errors. Fixed a pre-existing stale `related: ../../specs/cli/SPEC.md` link in the snapshot-retention PLAN (that spec doesn't exist yet) surfaced during this verify. ✅
- Remaining doctor warnings are unrelated/transient: (1) this request `active without SESSION.md` (workflow signal, not a docs defect); (2) the known SPEC-drift date heuristic (clears on spec-sync/re-run). ✅

## Result

All three milestones complete; zero behavior change as constrained. Ready for `review` → verify walk.
