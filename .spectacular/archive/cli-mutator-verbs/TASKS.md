---
status: verified
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — CLI mutator verbs

## M1 — Shared frontmatter helpers
- [ ] `fm_get <file> <field>` — reads top-level scalar field; returns empty if absent (uses existing awk pattern from `config_workspace_schema`)
- [ ] `fm_set <file> <field> <value>` — writes/replaces top-level scalar; preserves surrounding YAML; adds field if missing (insert before first non-frontmatter line)
- [ ] `fm_add_to_list <file> <field> <item>` — appends item to a YAML list field (creates the list if absent); skips if item already present
- [ ] `fm_touch <file>` — sets `updated:` to today via `fm_set`
- [ ] Helper: `is_canonical_doc <file>` — returns 0 if file is in {PRD, SPEC, ARCHITECTURE, PRINCIPLES, ROADMAP, STACK, DECISIONS, AGENTS, specs/<cap>/SPEC.md}
- [ ] Unit-ish smoke: each helper exercised in `tests/cli/mutator.test.sh` Scenario 1

## M2 — `spectacular touch <file>`
- [ ] New `cmd_touch` function
- [ ] Dispatch added to main() (replaces no current behavior — net-new verb)
- [ ] Validates: file exists, has frontmatter delimiter
- [ ] Idempotent: skip write if `updated:` already today (print "already current"); exit 0
- [ ] Test: scenarios for fresh write, idempotent re-run, non-frontmatter file refused

## M3 — `spectacular new <slug>`
- [ ] Replace existing `skill_verb_message "new" ...` stub with `cmd_new`
- [ ] Args: `<slug>` (required); `--summary "..."`; `--status planned|active|review`; `--priority low|medium|high`
- [ ] Slug validation: kebab-case (or per active convention pack's naming rules); max 64 chars; refuses reserved-id values
- [ ] Duplicate check: refuses if `.spectacular/requests/<slug>/` OR `.spectacular/archive/<slug>/` already exists
- [ ] Scaffolds `PLAN.md` from `templates/plan/base.md` with frontmatter prefilled
- [ ] Scaffolds `TASKS.md` from `templates/tasks/base.md` with frontmatter prefilled
- [ ] Test: scenarios for happy path, dupe rejection, slug validation rejection, --status override, --priority override

## M4 — `spectacular promote <slug> [--to <state>]`
- [ ] Replace existing `skill_verb_message "promote" ...` stub with `cmd_promote`
- [ ] State machine: planned → active → review → verified (linear)
- [ ] Default `--to`: next state in chain; `--to <explicit>` for jumps
- [ ] Refuses backward transitions (e.g. verified → active) unless `--force`
- [ ] Refuses unknown states; refuses if slug not found
- [ ] Mutates: PLAN.md `status:` + `updated:` (via fm_set)
- [ ] Optional: `--archive` flag — if promoting to verified, also runs `cmd_archive` after
- [ ] Test: scenarios for each forward transition, illegal backward, --to jumping, --archive combo, unknown slug

## M5 — `spectacular snapshot <file>`
- [ ] Replace existing `skill_verb_message "snapshot" ...` stub with `cmd_snapshot`
- [ ] Args: `<file>` (path to canonical doc)
- [ ] Validates via `is_canonical_doc`; refuses with explanatory error if not canonical
- [ ] Determines next snapshot number: scan for `<base>@v*.md`; pick highest + 1; default v1 if none
- [ ] Reads current `version:` field; computes new version (default: minor bump from "X.Y" → "X.(Y+1)"; supports `--major` flag for "(X+1).0")
- [ ] Idempotence: if content-hash matches latest snapshot, exit 0 with "no changes since last snapshot"
- [ ] Writes: `<base>@v<N>.md` (copy of current state pre-snapshot), then bumps `version:` + `updated:` in the unversioned file
- [ ] Test: scenarios for first snapshot, subsequent snapshot, idempotent re-run, non-canonical refused, --major bump

## M6 — `spectacular archive <slug>`
- [ ] Replace existing `skill_verb_message "archive" ...` stub with `cmd_archive`
- [ ] Args: `<slug>` (required); `--force` (skip lifecycle gate)
- [ ] Refuses if PLAN.md `status:` is not in {`verified`, `review`} without `--force`
- [ ] Mutates PLAN.md frontmatter: `status: archived` + `archived: <today>` + `updated: <today>`
- [ ] Moves: `.spectacular/requests/<slug>/` → `.spectacular/archive/<slug>/` (git mv if in git repo, else mv)
- [ ] Link rewriter: scans all `.spectacular/requests/*/PLAN.md`, `TASKS.md`, `VERIFY.md` for inbound `related:` paths matching the archived slug; rewrites the relative path to point at `../../archive/<slug>/...`
- [ ] Refuses if slug not found in requests/ (suggest checking archive/)
- [ ] Test: scenarios for verified archive, review archive, non-verified refused without --force, link rewrite (no broken links after), slug-not-found error

## M7 — Skill instruction sync
- [ ] SKILL.md routing table: add 5 new rows pointing each verb to itself (no skill ref needed; CLI does the work)
- [ ] SKILL.md: add top-level principle "Lifecycle mutations go through CLI verbs; only fall back to manual file edits when no CLI verb covers the case." in the Operating principles section
- [ ] `references/new-request.md`: replace manual-scaffold instructions with `spectacular new <slug>` invocation
- [ ] `references/archive.md`: replace manual git-mv + link-rewrite with `spectacular archive <slug>` invocation
- [ ] `references/lifecycle.md`: replace status-edit instructions with `spectacular promote <slug>` invocation
- [ ] `references/versioning.md`: replace snapshot-then-bump instructions with `spectacular snapshot <file>` invocation
- [ ] Spot-check other refs for stale "edit X manually" instructions; update or leave with a "see also: CLI verb" pointer

## M8 — Tests + docs + v0.7.0 release
- [ ] `tests/cli/mutator.test.sh` with full coverage of 5 verbs + edge cases (above scenarios)
- [ ] CHANGELOG.md entry for v0.7.0
- [ ] `.claude-plugin/plugin.json` version bump to 0.7.0
- [ ] SPEC.md capability bullet: "CLI mutator verbs (v0.7.0+)"
- [ ] CLAUDE.md Active Requests update; Archived list update
- [ ] Live doctor: clean
- [ ] Close task #54 (absorbed by M6)

## Verification (per 2-of-6 rule)

Hits: user-visible change (5 new verbs), multi-surface flow (CLI + skill refs both touched), external contract change (verbs become public). 3 of 6 → VERIFY.md required when reaching review.

Carry to VERIFY.md when ready:
- [ ] Each verb invoked end-to-end in a temp workspace; doctor passes after
- [ ] Skill follows new instructions (interactive test: ask skill to archive a request; verify it runs `spectacular archive` not manual mv)
- [ ] No regressions in 7 existing test files (181+ asserts)
- [ ] Live workspace dogfood: archive a real verified request via new verb (cli-mutator-verbs itself, on close-out)
