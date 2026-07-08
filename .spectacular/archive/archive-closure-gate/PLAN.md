---
status: archived
priority: high
owner: alex
updated: 2026-07-07
build: b22
related:
  - ../../PRD.md
  - ../../specs/index.md
summary: "Archive closure gate + delta-based spec-sync: archive blocks (with recorded override) on unticked TASKS, unwalked VERIFY, or missing spec delta; spec impact is declared as ADDED/MODIFIED/REMOVED bullets and merged mechanically."
archived: 2026-07-07
---

# Plan — archive-closure-gate

## 1. Goal

Close the lifecycle tail: `spectacular archive` refuses (block-with-recorded-override) unless TASKS/VERIFY are closed out and the request declares its spec impact as an explicit delta — serving PRD's "the workspace reflects what is actually built" promise. (Fable review 2026-07-06: anti-drift was the corpus's weakest metric, and it is a lifecycle-*tail* failure — see `docs/reviews/fable-spec-quality-review.md` §1.)

## 2. Constraints

- CLI-is-the-mutator contract untouched: the gate lives in `cmd_archive` (mechanical checks) + `archive.md` (judgment flow); the skill proposes, the CLI enforces.
- Every block is overridable **once, explicitly**, and the override is recorded in the archive note — never a silent `--force` bypass, never an unrecoverable hard wall.
- Existing archived requests are grandfathered — no retroactive validation; `doctor` may *report* legacy gaps but `--fix` never rewrites archive/.
- Delta merge must be mechanical (no LLM in the CLI path): ADDED appends, MODIFIED replaces a quoted bullet, REMOVED deletes a quoted bullet.
- Backwards compatible: a workspace with no `specs/` and a stub SPEC.md still archives (delta `NONE — <why>` is valid).

## Understanding

### How it works now

`cmd_archive` (cli/spectacular:~4766) gates only on lifecycle status (`verified|review`, else `--force`) and emits an *advisory* warning if VERIFY-LOG.md is missing (cli:~4835, verify.md § Lifecycle tie-in — "Advisory, not blocking"). It never inspects TASKS.md completion or VERIFY.md checkbox state. Spec-sync (`spec-sync.md`) is a free-form skill proposal ("here's what I'd update") with no structure `doctor specs` can validate — drift is caught only by a date heuristic. Corpus evidence of the gap: wasabi `m1-two-webviews` archived with 4 open TASKS boxes and a fully unchecked VERIFY.md; harbor SPEC.md frozen at "nothing built yet" against ~30 shipped requests.

### What changes

- `cmd_archive`: three closure checks before the move — (1) every TASKS box `[x]` or `[~]` with a reason, (2) if VERIFY.md exists, VERIFY-LOG.md has ≥1 walk entry, (3) a spec-delta block exists in the archive proposal (or explicit `NONE — <why>`). Each failing check blocks; `--override <check> --reason "<text>"` bypasses once and stamps the reason into the archived PLAN frontmatter (`archive_overrides:`).
- `spec-sync.md`: proposal format becomes structured deltas — `### ADDED` / `### MODIFIED` (quotes current bullet + replacement) / `### REMOVED` (quotes bullet) — merged mechanically into SPEC.md after human confirmation.
- `archive.md`: sequence gains the closure-gate step before the CLI verb; documents the override etiquette.
- `doctor specs`: validates delta integrity structurally (MODIFIED/REMOVED must quote a bullet that exists; ADDED must not duplicate one) — replaces date-only drift detection as the primary signal.

### What stays the same

- The 4-state lifecycle vocabulary and `advance` semantics.
- Spec-sync stays propose-then-confirm — the human still approves every SPEC edit; only the *format* and the *gate* are new.
- VERIFY.md remains opt-in to create; the walk requirement only applies when it exists.
- `archive/` remains append-only history; nothing revalidates or rewrites old archives.

## Decisions

- Chose **block-with-recorded-override** over hard-block and over advisory-warn — advisory is what failed in the corpus (the existing VERIFY-LOG warning was ignored), and a hard wall would force `--force` habits; a recorded override keeps friction *and* an audit trail.
- Chose **delta blocks (ADDED/MODIFIED/REMOVED)** over free-form proposals — mechanical mergeability is what lets `doctor` validate structurally instead of by date heuristic (OpenSpec-inspired, fable review idea #2).
- Chose **`[~]`-with-reason as a valid TASKS closure state** over requiring all-`[x]` — deliberate deferral is legitimate (wasabi `sleep-policy` exemplar); only *unexplained* open boxes block.

## Formats

The four on-disk contracts the CLI gate (M2) and doctor (M3) read. Locking them here so the skill side (M1) writes exactly what the CLI side reads — this is the M1↔M2 handoff.

**1. Spec-delta file — `.spectacular/requests/<slug>/SPEC-DELTA.md`.** The skill writes it during the archive walk (M1); the CLI reads it (M2 check 3) and merges it (M3 doctor validates). A proposal in chat is not enough — the CLI can't see chat, so the delta must land in this file before `archive` runs. Shape:

```md
### ADDED
- specs/billing.md :: team-billing — seats are Stripe-backed; a duplicate webhook never double-charges

### MODIFIED
- SPEC.md :: "billing — single-seat only" -> "billing — multi-seat, Stripe-backed"

### REMOVED
- SPEC.md :: "auth — password-only login"
```

Rules: each line is `<file> :: <payload>`. ADDED payload is the new bullet text. MODIFIED payload is `"<exact current bullet>" -> "<replacement>"` (arrow-separated, current side quoted verbatim so the merge can find it). REMOVED payload is `"<exact current bullet>"`. Empty impact is the single line `NONE — <why>` (no section headings) — a valid, passing delta. The merge is mechanical: ADDED appends under the file's capabilities list, MODIFIED does a quoted-string replace, REMOVED deletes the quoted line.

**2. VERIFY-LOG walk-entry marker.** M2 check 2 ("≥1 walk entry") greps VERIFY-LOG.md for a **result row** — a table line containing `✅` (the pass glyph the walk already writes, see `verify.md` § 3). Presence of ≥1 `✅` row means the walk ran and something passed. A VERIFY-LOG.md that exists but has zero `✅` rows counts as *not walked* (blocks). This reuses the existing recording format — no new marker invented.

**3. `[~]` deferred-task syntax — `- [~] <task> — <reason>`.** An em-dash-or-hyphen-separated reason after the task text. M2 check 1: `- [x]` passes; `- [~]` passes **only if** a ` — ` (space-dash-space) reason follows on the same line; bare `- [ ]` or reasonless `- [~]` blocks, naming the box. The check greps TASKS.md task lines only (under `### M<N>` headings), ignoring the `## v2 (deferred)` section.

**4. `--force` vs `--override` split.** `--force` keeps its **current, single meaning**: bypass the *status gate* only (line 4830) — untouched by this request. The three *new* closure checks are bypassed only by `--override <check> --reason "<text>"`, where `<check>` ∈ `{tasks, verify, spec}`. They do not overlap: `--force` never silences a closure check, `--override` never silences the status gate. Each `--override` stamps `archive_overrides:` into the archived PLAN frontmatter as a list of `{check, reason, date}`. Passing `--override` for a check that would have passed anyway is a no-op (not an error).

## 3. Milestones

- M1 — Skill-side flow: `archive.md` closure-gate step + `spec-sync.md` delta proposal format shipped; a dry-run archive walk on a fixture request writes a `SPEC-DELTA.md` in the Formats §1 shape.
- M2 — CLI gate: `cmd_archive` runs the three closure checks; blocks with actionable messages; `--override <check> --reason` records into archived PLAN frontmatter.
- M3 — Doctor + tests: `doctor specs` validates delta integrity; `tests/` covers gate-blocks, override-recording, `NONE` delta, and grandfathered legacy archives.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None hard. Complements [[spec-audit-mode]] (b11) — this request prevents *new* drift at the convergence point; spec-audit-mode detects *existing* drift. Ship in either order.

## 6. Validation

- M1 — observable: archiving a fixture request in a scratch workspace yields a proposal containing `### ADDED`/`### MODIFIED`/`### REMOVED` (or `NONE — <why>`) sections; free-form prose proposals no longer appear in `spec-sync.md`.
- M2 — run: a test archive with one unticked TASKS box exits non-zero naming the box; re-run with `--override tasks --reason "deferred to b23"` exits 0 and `grep archive_overrides archive/<slug>/PLAN.md` finds the reason.
- M3 — run: `tests/` suite green including the four new cases; `spectacular doctor specs` on a fixture with a MODIFIED delta quoting a nonexistent bullet exits non-zero.

## 7. Deliverables

- Updated `cli/spectacular` (`cmd_archive` closure checks + `--override`), `references/archive.md`, `references/spec-sync.md`, `references/doctor-areas.md` (specs area).
- New test file under `tests/` covering the gate matrix.
- CHANGELOG entry.
