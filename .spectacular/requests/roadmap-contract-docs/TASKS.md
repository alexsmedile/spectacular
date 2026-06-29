---
status: verified
updated: 2026-06-29
related:
  - PLAN.md
---

# Tasks — roadmap-contract-docs

## v1

### M1 — Spec the ledger
- [x] `spectacular snapshot .spectacular/specs/roadmap/SPEC.md` → @v1 (canonical doc, snapshot first)
- [x] Add Ledger section to specs/roadmap/SPEC.md (build ids, table schema, target-version single-source, status-vs-lifecycle, tbd); frontmatter published/1.1; fixed stale roadmap-overrides ref
- [x] SPEC.md structured-ROADMAP bullet names the ledger + links specs/roadmap/SPEC + ARCHITECTURE + roadmap-rules (fixed stale roadmap-overrides link)
- [x] M-question 1 resolved: ARCHITECTURE is canonical schema home; spec summarizes + points to it (no forked copy)
- [x] `doctor specs` green

### M2 — Define `tbd` + fix contradicting rule
- [x] ARCHITECTURE.md target-version row documents `tbd` sentinel (vs `<TBD>` placeholder)
- [x] roadmap-rules.md: ledger section + `tbd` behavioral rule; placeholder check scoped to "prose slot" + explicit "ledger tbd is not a placeholder" note
- [x] Check: ledger `tbd` passes the gate (own ROADMAP uses tbd, doctor green); prose `<TBD>` still rejected (rule text intact)

### M3 — User docs + tutorial
- [x] docs/commands.md `spectacular roadmap` section + build-id↔version model
- [x] docs/configuration.md documents `last_build:`
- [x] Tutorial walkthrough = versioning.md § "The roadmap ledger" (M-question 2: folded into versioning.md, the natural home, vs a new page); linked from commands.md + configuration.md
- [x] VERIFY-LOG
