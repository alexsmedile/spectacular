---
status: pending
updated: 2026-05-23
related:
  - PLAN.md
  - TASKS.md
score: 3/6
---

# Verification — Convention Pack Fabricator

Per [[verification]] 2-of-6 rule: this request scored 3/6 (user-visible new verbs + multi-surface flow + external contract change). Scenarios below are the ones automation cannot cover — interactive grill walkthroughs, source-ingestion judgment calls, dogfood validation.

## Manual scenarios

### S1 — Cold grill walkthrough (`spectacular pack new throwaway-test`)

- [ ] Engine resolves target path `~/.spectacular/packs/throwaway-test/` and creates folder
- [ ] Engine scaffolds `pack.md` from `templates/packs/minimal/pack.md`
- [ ] All 7 slots surface in order (Name & scope → Naming → Taxonomy → Root files & README → Gitignore → File placement → Project types)
- [ ] Reserved pack-id check rejects `blank` / `minimal` / `default` / `pack` / `packs` during slot 1
- [ ] Mini-refine pattern fires on Slot 5 when user puts `.scrapekit/` in `always-add` (rejects; moves to `never-auto-add`)
- [ ] Mini-refine pattern fires on Slot 7 when user declares `project-types.<unknown-type>`
- [ ] Review gate runs checks 4-12 from pack-overrides.md
- [ ] On pass, pack.md frontmatter contains valid YAML for all declared slots
- [ ] `rm -rf ~/.spectacular/packs/throwaway-test/` after test

### S2 — Resume grill (`spectacular pack grill <existing-name>`)

- [ ] Engine detects existing pack.md and resumes at first empty / placeholder slot
- [ ] Already-filled slots are not re-prompted unless user explicitly requests
- [ ] Snapshot rule (snapshot-on-edit: false from registry) confirmed — no `pack@v*.md` files created during grill

### S3 — Source-ingestion mode (`--from`)

- [ ] `spectacular pack new test-src --from ~/code/NAMING_RULES.md` pre-fills Slot 2 (Naming)
- [ ] Each pre-filled slot prompts `Pre-filled from <path>: <answer>. Keep? [Y/n]`
- [ ] Answering `n` falls back to normal slot prompt
- [ ] `--from` with a `.gitignore` file pre-fills Slot 5 `always-add` list
- [ ] Ambiguous source content (e.g. README mentioning "we usually use kebab-case") does NOT pre-fill — slot asked normally
- [ ] Multiple `--from` paths processed in order; later sources can refine earlier pre-fills

### S4 — Review gate negative tests (`spectacular pack review <broken-pack>`)

- [ ] Pack with all 6 rule blocks empty → fails check 4 ("At least one rule category populated")
- [ ] Pack with `templates: [foo.md]` but no `templates/foo.md` file → fails check 6
- [ ] Pack with `version: "not-semver"` → fails check 9
- [ ] Pack with same item in both `gitignore.always-add` and `gitignore.never-auto-add` → fails check 10
- [ ] Pack with `applies-to: [made-up-type]` → produces warning, not error (per check 8 spec)

### S5 — Dogfood validation: `packs/alex-default/`

- [x] Pack.md exists at `packs/alex-default/pack.md`
- [x] Frontmatter parses through awk + Python YAML
- [x] All 6 rule categories present and non-empty
- [x] All 8 project types declared (cli, library, webapp, skill, plugin, content, research, vault-project)
- [x] Templates declared (.gitignore + README.md) exist in `templates/`
- [x] References declared (why-alex-default.md) exists in `references/`
- [x] All 10 archived conventions from `archive/repo-conventions/PLAN.md` expressible in the manifest (per § Schema coverage check in packs-contract.md)
- [ ] `spectacular pack review alex-default` passes the gate end-to-end (requires running the skill grill engine, not mechanical)

### S6 — Schema parity

- [ ] Every field referenced by the grill (pack-overrides.md slots) is documented in the schema (packs-contract.md). No grill prompt produces a field the schema doesn't define.
- [ ] Conversely: no schema field exists that the grill cannot produce (no orphan schema surface).

## Sign-off

- [ ] All scenarios complete or explicitly waived with reason
- [ ] No `❌` items remain; `⚠️` items have a note
- [ ] Verifier: alex
- [ ] Date verified: ____________
