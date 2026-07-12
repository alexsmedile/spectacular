---
status: verified
priority: high
owner: alex
updated: 2026-07-12
build: b29
summary: "CLI-side fixes for the b28 dogfood findings 1–3: policy gate injects each policy's new one-sentence directive (warn=directive+title, block=directive+full, --full restores), advance planned→active scaffolds SESSION.md, doctor repeats non-pass findings after the summary."
related:
  - PRD.md
---

# Plan — cli-gate-ergonomics

## Goal

Remove the three CLI frictions the b28 dogfood run surfaced: the policy gate's full-paragraph tax at every phase entry, the advance verb's open create→warn→repair loop around SESSION.md, and doctor findings that drown mid-report.

## Constraints

- All three changes live in `cli/spectacular` (Bash) — no skill-reference edits except where output examples are quoted.
- Human output only changes shape, never meaning; `--json` forms stay byte-compatible (machine consumers use them).
- The skill's policy-injection loop stays authoritative — the CLI surfaces, it never blocks beyond today's mechanical understand-before-change check.
- Match existing CLI idioms: `doc_add`/emitter patterns for doctor, heredoc templates for scaffolds, `printf` column style for policy.
- Every behavior change lands with a `tests/cli/*.test.sh` case in the existing harness style (`seed_ws`, `assert_output_contains`).

## Understanding

### How it works now

`cmd_policy` hook-form (cli/spectacular ~7788–7812) prints `→ P<n> <full principle line>` for every policy — ~90-word paragraphs injected at every phase gate; `_policy_consult_transition` (~5646) embeds that same output into `advance`. `advance` planned→active writes PLAN/TASKS frontmatter but never creates SESSION.md; `check_lifecycle` (~9321) later flags the absence as a judgment warning. `doctor_emit_text` (~11261) prints pass and non-pass rows interleaved, counts-only at the end, so extracting findings needs greps.

### What changes

(1) POLICY.md's anatomy gains an optional `- directive:` field — one imperative sentence written for injection (the body prose is *designed* as "the instruction injected into context" but the hook form never prints it; today the agent gets the theory layer instead of the practice layer). Hook-form rows become: warn → directive + `P<n> — <title>`; block → directive + full principle line; fallback when a policy has no directive yet → principle title (incremental migration). New `--full` flag restores today's output everywhere; `--json` gains a `directive` key. All ~13 existing policies get a directive authored. (2) `advance` planned→active writes SESSION.md from an embedded heredoc (frontmatter + Current state / Active task / Blockers / Next actions) when absent, printing `✓ scaffolded: SESSION.md`. (3) `doctor_emit_text` appends a `── findings ──` block after the count line repeating every non-pass finding compactly (icon area · file — msg → fix).

### What stays the same

`policy --json`, `policy <id>` (already the single-policy full view), `policy --principle N`; doctor's per-area report body and JSON emitter; the advance state machine and its policy consult; the skill-side SESSION.md template in active-request.md (the CLI heredoc duplicates it deliberately — ~10 stable lines).

## Decisions

- Block=full / warn=title (grilled 2026-07-12): a block is the moment the reasoning must be in-context; warns need only the mnemonic. Rejected title-only-for-all (agents argue with unexplained refusals) and ids-only (loses policy→principle traceability).
- Amended same day: inject the policy's own `directive:` (new one-sentence field) as the primary row text, principle title/line as the trailer — because POLICY.md bodies are the designed injection layer yet the hook form printed only the principle (theory instead of practice). Chose an explicit field over first-sentence extraction (bodies contain tables/laws — extraction is fragile) with title-fallback so migration is per-policy incremental.
- CLI scaffolds SESSION.md, no `--no-session` flag (grilled): mechanical file creation is the CLI's job per the mutation principle; no known active-without-SESSION case exists — doctor treats it as a defect.
- Repeat-after-summary, no `--quiet` flag (grilled): `tail` always works, no new flag to teach; pass rows stay for humans. Rejected hiding pass rows (behavior change users would feel).

## Milestones

- M1 — Policy gate tiered output + directive field: warn=directive+title, block=directive+full, title-fallback, `--full` flag; all existing policies get authored directives; advance gate inherits via `_policy_consult_transition`
- M2 — `advance` planned→active scaffolds SESSION.md when absent
- M3 — Doctor `── findings ──` block after the summary line

## Tasks

See `TASKS.md`.

## Dependencies

- [[self-healing-optimization]] (b28) — source of the dogfood findings; its PLAN § Dogfood review is the evidence base

## Validation

- M1 — run: `tests/cli/policy-output.test.sh` exits 0 (asserts: warn row with directive shows `— <directive>` + `P<n> — <title>`; warn row without directive falls back to title; block row shows directive + full paragraph; `--full` restores paragraphs; `--json` carries `directive`; `advance`'s embedded gate matches); assertable: every enabled policy in POLICY.md greps a `- directive:` line
- M2 — run: seeded workspace, `spectacular advance <slug>` from planned → SESSION.md exists with the four H2s, output contains `✓ scaffolded: SESSION.md`; `spectacular doctor lifecycle` reports 0 warnings; a second advance never overwrites an existing SESSION.md
- M3 — run: seeded workspace with 1 error + 1 warning, `spectacular doctor | tail -6` contains the `── findings ──` block with both rows; clean workspace prints no findings block

## Deliverables

- `cli/spectacular` — the three changes above
- `tests/cli/policy-output.test.sh` (new) + cases in `tests/cli/mutator.test.sh` (advance) and `tests/cli/doctor.test.sh` (findings block)
- CHANGELOG entry; policy-injection.md / doctor.md refreshed only where they quote output shapes
