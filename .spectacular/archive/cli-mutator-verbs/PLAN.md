---
status: archived
priority: high
owner: alex
updated: 2026-05-23
target_version: 0.7.0
summary: "5 CLI verbs (promote, archive, snapshot, new, touch) replace skill-side manual file edits for the most common lifecycle mutations. CLI mutates deterministically; skill orchestrates."
related:
  - ../../archive/workspace-migrations/PLAN.md
  - ../../../skills/spectacular/references/lifecycle.md
  - ../../../skills/spectacular/references/archive.md
  - ../../../skills/spectacular/references/versioning.md
provenance:
  - source: live-session insight (2026-05-23)
    captures: "User flagged that lifecycle mutations done via skill free-form edits are fragile, non-deterministic, and untestable. Should formalize 'skill orchestrates, CLI mutates' as the contract."
archived: 2026-05-23
---

# Plan — CLI mutator verbs

## Goal

Make the canonical lifecycle mutations of Spectacular **first-class CLI verbs**. The skill orchestrates (reads, decides, communicates); the CLI mutates (atomically, deterministically, tested). Manual file edits remain available for edge cases the verbs don't cover but become the exception, not the rule.

## Why

Today most lifecycle mutations happen via the skill writing free-form edits:

- Promoting a request: skill rewrites `status:` in PLAN frontmatter
- Archiving: skill `git mv`s the dir + greps for `related:` paths and sed-rewrites them (this is exactly what task #54 was about)
- Snapshotting a canonical doc: skill copies the file with `@vN` suffix
- Scaffolding a new request: skill writes PLAN.md + TASKS.md from a template

Cost of doing it this way:

1. **Non-determinism** — two agents, two indentation styles, two frontmatter key orderings
2. **Untestable** — there's no way to assert "promoting works"; we trust the agent
3. **Drift** — when conventions change (e.g. `status` field gains values, snapshot naming gains rules), every agent invocation needs to relearn
4. **No safety net** — manual `git mv` doesn't update inbound `related:` links; manual frontmatter rewrites don't snapshot

The fix: 5 CLI verbs absorb the high-frequency mutations. Skill calls them. CLI has tests. Edge cases (e.g. a manual frontmatter edit for an unusual field) still work as fallback.

## Scope

### In scope (v1 of this request)

**Five CLI verbs**, each replacing a stubbed skill-only message:

1. **`spectacular promote <slug> [--to <state>]`**
   - States: `planned | active | review | verified`
   - Default `--to`: advances to next state in the lifecycle (planned → active → review → verified)
   - Rewrites PLAN.md frontmatter `status:` + `updated:` fields
   - Refuses illegal transitions (e.g. verified → active without `--force`)
   - Auto-triggers `archive` when promoting to `verified` AND user passes `--archive`

2. **`spectacular archive <slug>`**
   - Moves `.spectacular/requests/<slug>/` → `.spectacular/archive/<slug>/`
   - Rewrites inbound `related:` paths in any other PLAN.md / TASKS.md / VERIFY.md
   - Bumps PLAN.md frontmatter: adds `archived: <today>`, sets `status: archived`
   - Refuses if not in `verified` or `review` state (use `--force` to override)
   - Absorbs task #54 (archive --check / auto-rewrite related: paths) — same behavior, different verb name

3. **`spectacular snapshot <file>`**
   - Snapshots a canonical doc to `<FILE>@v<N>.md` (next N inferred from existing snapshots)
   - Bumps `version:` field in the unversioned file to the new N (e.g. `1.1` → `1.2`)
   - Bumps `updated:` to today
   - Refuses if file is not a registered canonical doc (PRD, SPEC, ARCHITECTURE, PRINCIPLES, ROADMAP, STACK, DECISIONS, AGENTS, specs/<cap>/SPEC.md)
   - Exits cleanly if no changes since last snapshot (idempotent — same content hash)

4. **`spectacular new <slug> [--summary "..."] [--status planned]`**
   - Scaffolds `.spectacular/requests/<slug>/PLAN.md` + `TASKS.md`
   - Frontmatter prefilled: `status:` (default `planned`), `priority:`, `owner:` (from config), `updated:`, `summary:`, empty `related:`
   - Body has the 7-slot PLAN template + minimal TASKS scaffold
   - Refuses if `<slug>` already exists (in requests/ OR archive/)
   - Slug validation: kebab-case, no dots, max 64 chars (per pack rules if a pack is active)

5. **`spectacular touch <file>`**
   - Updates frontmatter `updated:` field to today's date
   - No-op if already today
   - Works on any markdown file with a frontmatter block
   - Useful when content is edited via tooling that doesn't auto-update timestamps

### Out of scope (v2+)

- `spectacular decisions add "<title>"` — append DECISIONS entry
- `spectacular remember "<lesson>"` — write memory entry (currently a skill stub)
- `spectacular link <from> <to>` — add cross-link in `related:`
- `spectacular task done <slug> <task-pattern>` — toggle TASKS checkboxes (probably too micro)
- `spectacular bump <file>` — version bump without snapshot
- Atomic batch operations (`spectacular pipeline promote-archive <slug>` — combo)

### Explicit anti-patterns

- Verbs that wrap a single sed/grep (no abstraction value)
- Verbs that need an AI agent to make a judgment call inside them (those belong in skill flow)
- Mutating files outside `.spectacular/` (e.g. host project code) — out of scope for Spectacular CLI
- Silently rewriting fields the user didn't ask to change (e.g. promoting doesn't reorder frontmatter keys)

## Skill enforcement

Both surfaces updated per design lock (2026-05-23):

1. **SKILL.md routing table** — entries added for each verb plus a top-level principle: "Lifecycle mutations go through CLI verbs; only fall back to manual file edits when no CLI verb covers the case."
2. **Reference doc rewrites** — `new-request.md`, `archive.md`, `lifecycle.md`, `versioning.md` updated to instruct the skill to call the CLI verb instead of doing the mutation directly.

## Decisions (locked 2026-05-23 via interview)

- **Scope cut**: 5 verbs (promote, archive, snapshot, new, touch) shipping as v0.7.0. Comprehensive 11-verb cut deferred.
- **Skill enforcement**: both SKILL.md routing AND per-reference rewrites. Belt-and-suspenders to resist skill drift.
- **CLI exits non-zero on illegal mutations** — agents handle the error; not silent.
- **Verbs are read-then-write, not transactional** — race conditions between concurrent agents are out of scope (single-user tool).
- **Frontmatter writer is one shared helper** — `fm_set <file> <field> <value>` (and `fm_get`, `fm_add_to_list`); all 5 verbs use it. No 5 different YAML rewriters.
- **Archive absorbs task #54** — same goal, no separate request needed.

## Validation

- `spectacular promote <slug>` advances status one step; `--to <state>` jumps; illegal transitions exit 2 with clear message
- `spectacular archive <verified-slug>` moves the dir AND rewrites every inbound `related:` link; doctor `links` area passes after
- `spectacular snapshot .spectacular/PRD.md` creates `PRD@v<N+1>.md` AND bumps PRD.md `version:` field; second snapshot with no changes exits 0 with "no changes since last snapshot"
- `spectacular new test-throwaway` scaffolds PLAN.md + TASKS.md with correct frontmatter; `spectacular new test-throwaway` again refuses with "slug already exists"
- `spectacular touch .spectacular/PRD.md` updates `updated:` field to today; running again same day is a no-op
- After this ships, the skill stub messages for these verbs (currently "this is an interactive skill flow") are replaced with actual CLI behavior
- Doctor stays clean across all 5 verb invocations
- All 7 existing test files still green

## Milestones

1. **M1 — Shared frontmatter helpers** — `fm_get`, `fm_set`, `fm_add_to_list`, `fm_touch` in `cli/spectacular`. Used by all 5 verbs.
2. **M2 — `spectacular touch`** — simplest verb. Validates frontmatter helpers work end-to-end.
3. **M3 — `spectacular new`** — replaces stub with actual scaffold. Includes slug validation + duplicate detection.
4. **M4 — `spectacular promote`** — state machine + frontmatter mutation. Includes `--to` flag + illegal-transition refusal.
5. **M5 — `spectacular snapshot`** — file copy with @vN naming + version bump. Includes idempotence check.
6. **M6 — `spectacular archive`** — git mv + inbound link rewrite. Absorbs task #54.
7. **M7 — Skill instruction sync** — SKILL.md routing + reference doc rewrites (new-request.md, archive.md, lifecycle.md, versioning.md). Add top-level principle.
8. **M8 — Tests + docs + v0.7.0 release** — `tests/cli/mutator.test.sh` covering all 5 verbs + edge cases; CHANGELOG; plugin bump; SPEC.md capability bullet.

## Risks

- **Frontmatter rewriters break on edge cases** (multi-line strings, comments, unicode). Mitigation: shared helpers tested separately + Python YAML fallback for verification.
- **Slug validation diverges from pack rules** when a convention pack is active. Mitigation: read pack's `naming.folder-case:` if `convention_pack:` is declared; otherwise default kebab-case.
- **Archive link-rewrite over-matches** (e.g. rewrites a related: that just happens to share a path prefix). Mitigation: literal-string match with full-path resolution; test scenarios cover near-miss patterns.
- **Snapshot version bump conflicts with manual version edits** in the doc. Mitigation: read existing `version:` field; bump by 0.1 if minor change implied (no flag), or accept `--major` for breaking changes. Show the proposed bump; refuse if non-numeric.
- **Skill keeps doing manual edits anyway** because the new instruction is one paragraph in a long file. Mitigation: aggressive routing table entries + the affected reference docs ALSO updated (not just SKILL.md). Belt-and-suspenders.

## Open questions (for v2+)

- Should there be a `spectacular pipeline <name>` verb for combo operations (promote-archive, snapshot-bump-touch)? Likely yes once we have 3+ verbs that commonly chain.
- Should the verbs emit structured JSON output (`--format json`) for programmatic consumption? Probably yes for promote + archive (state machine outcomes), no for touch (trivial).
- How does this interact with workspace-migrations? Migration apply-fns might want to use these helpers (`fm_set` for `workspace_schema:` bump). Refactor opportunity once both ship.
