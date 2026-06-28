---
status: planned
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — roadmap-contract-docs

## v1

### M1 — Spec the ledger
- [ ] `spectacular snapshot .spectacular/specs/roadmap/SPEC.md` (canonical doc, snapshot first)
- [ ] Add Ledger section to specs/roadmap/SPEC.md (build ids, table schema, target-version single-source, status-vs-lifecycle, tbd); refresh frontmatter (published, 2026-06-28)
- [ ] SPEC.md structured-ROADMAP bullet names the ledger + links specs/roadmap/SPEC
- [ ] Resolve M-question 1 (canonical schema home: ARCHITECTURE vs spec)
- [ ] `doctor specs` green

### M2 — Define `tbd` + fix contradicting rule
- [ ] ARCHITECTURE.md target-version row documents `tbd` sentinel
- [ ] roadmap-rules.md: ledger `tbd` rule + scope placeholder check (prose `<TBD>` still rejected)
- [ ] Check: ledger `tbd` passes the gate; prose `<TBD>` still fails

### M3 — User docs + tutorial
- [ ] docs/commands.md `spectacular roadmap` section + build-id↔version model
- [ ] docs/configuration.md documents `last_build:`
- [ ] Short tutorial walkthrough (location per M-question 2); link from docs/versioning.md
- [ ] VERIFY-LOG
