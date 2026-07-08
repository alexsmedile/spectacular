---
status: archived
priority: medium
owner: alex
updated: 2026-07-06
build: b18
summary: "Stop ROADMAP.md bloating agent context: enforce the already-stated 'shipped history lives in CHANGELOG' principle by pruning shipped prose blocks down to their ledger row (the index), and/or add a decisions-index-style roadmap mode (cheap ledger + per-version files). ~49% of ROADMAP.md is currently past-tense duplication."
related:
  - ../../PRD.md
  - ../../roadmaps/index.md
  - ../../specs/roadmap.md
  - ../../ARCHITECTURE.md
depends-on: roadmap-contract-docs
archived: 2026-07-06
---

# Plan — roadmap-pruning

> **Origin (2026-06-28):** ROADMAP.md is **525 lines, ~49% past-tense** — 12 shipped
> per-version prose blocks (~230 lines) plus a "Recently shipped" section (~26 lines) that
> **duplicates CHANGELOG.md outright**. An agent loading ROADMAP.md to make a planning
> decision pays for all of it. The kicker: both `specs/roadmap.md:16` and
> `ARCHITECTURE.md:257` already state *"shipped history lives in CHANGELOG"* — the live file
> just doesn't obey it. There's no pruning/archiving mechanism, and no roadmap analog to the
> already-shipped **decisions-index** pattern (cheap root index + per-entry files for large
> collections), even though the ledger table is *already* a compact index — one row per build.

## 1. Goal

Keep ROADMAP.md's agent-context cost bounded as history grows, **without losing the
information needed to make planning choices**:

- The **ledger table** stays the always-present compact index (one row per build — past and future).
- **Forward-looking prose** (planned/active version blocks + vision + Icebox) stays inline — that's the part planning decisions actually read.
- **Shipped prose blocks** stop accumulating: either pruned to their ledger row (CHANGELOG holds the detail, per the stated principle) or split into per-version files behind a cheap index (decisions-index pattern). Decide which.

## 2. Constraints

- **Depends on `roadmap-contract-docs` (b17).** The ledger must be properly specced first, since this request makes the ledger the load-bearing index and prunes the prose around it. Build b17 → then b18.
- **Bash CLI only; deterministic mutator → CLI-owned.** Pruning/archiving is a mutation; the skill may suggest, the CLI performs. Doctor detects the bloat; `--fix` (or a dedicated verb) prunes.
- **Never lose information.** Shipped detail must survive somewhere authoritative — CHANGELOG.md (already the stated home) or per-version `roadmap/` files. Prune is a move/dedup, never a delete of unrecoverable content. Snapshot ROADMAP.md before restructuring.
- **Don't break the prose-block spec.** `specs/roadmap.md` specs per-version blocks; whichever approach wins must update that spec, not silently violate it.
- **Reversible / dry-run first.** Show what would be pruned before doing it (mirror the snapshot-retention + archive patterns).

## Understanding

### How it works now (from the audit)

- 525 lines. ~230 in 12 shipped prose blocks (v1.9–v1.22), ~26 in a "Recently shipped" CHANGELOG mirror (ROADMAP.md:490-515, even links to CHANGELOG). Combined ~256 lines ≈ 49% past-tense.
- Stated-but-violated principle: `specs/roadmap.md:16` + `ARCHITECTURE.md:257` say shipped history belongs in CHANGELOG. The live file keeps full shipped prose *and* the mirror anyway.
- No pruning mechanism anywhere. Doctor's roadmap checks only flag legacy shape (pre-v0.7.x), not bloat. The one soft guard — review-gate check 17 (roadmap-rules.md:366) — warns on too many *full-tier* blocks, but shipped blocks are mostly `themed`, so it doesn't catch this.
- **Decisions-index precedent exists and is fully specced** (`decisions-rules.md:55-89`): flat mode default; index mode for large projects = cheap one-line root index + per-entry files. Shipped in v1.17. No roadmap equivalent.

### Two candidate approaches (pick in M1)

**A. Prune-to-ledger (lean: simplest, honors stated principle).** Shipped version blocks get deleted from ROADMAP.md; the ledger row + CHANGELOG entry are the record. The "Recently shipped" mirror is removed (it's pure duplication). Result: ROADMAP.md = ledger + Next-up + forward prose + Icebox only. A `doctor roadmap` check flags shipped blocks older than the last N as prunable; `--fix` removes them after confirming the CHANGELOG entry exists.

**B. Roadmap-index mode (decisions-index parallel).** Shipped blocks move to `roadmap/v<X.Y.Z>.md` per-version files; ROADMAP.md keeps the ledger + a cheap shipped index + forward prose. More structure, keeps prose discoverable per version, costs a folder + a `spectacular roadmap migrate`.

Trade: A is laziest and matches the already-documented "history → CHANGELOG" rule; B preserves the prose narrative per version (some of which is richer than CHANGELOG). Lean A unless the shipped prose carries planning-relevant rationale CHANGELOG doesn't.

### What stays the same

The ledger table (grows one row per build forever — cheap), forward-looking prose, the Icebox, the review gate. Only the *shipped past-tense prose* is restructured.

## 3. Milestones

### M1 — Decide approach + spec it
- Choose A (prune-to-ledger) vs B (roadmap-index mode); resolve M-questions.
- Update `specs/roadmap.md` + ARCHITECTURE.md to describe the chosen retention/pruning model (so the stated principle is finally enforced, not just asserted).

### M2 — Detection (doctor)
- `doctor roadmap` (or extend the roadmap area): flag shipped prose blocks beyond the keep-window / the "Recently shipped" duplicate mirror as prunable. Info/warning; relayed by `status`.

### M3 — Prune mechanism
- CLI: prune verb or `doctor --fix roadmap` — snapshot ROADMAP.md, then remove/move shipped blocks (A: delete after CHANGELOG-presence check; B: move to `roadmap/v*.md` + write index). Dry-run default.
- Dogfood: prune THIS repo's ROADMAP.md (12 shipped blocks + the mirror) — confirm CHANGELOG covers each first.
- Tests + VERIFY-LOG.

## Decisions (resolved 2026-06-29)

1. **Approach B — roadmap-index mode.** Shipped prose blocks move to per-version files `.spectacular/roadmap/v<X.Y.Z>.md`; ROADMAP.md keeps the ledger + Next-up + a cheap **shipped index** (one line per shipped version → its file) + forward-looking prose + Icebox. Mirrors the shipped decisions-index pattern (`cmd_decisions_migrate` is the template). Chosen over A (prune-to-ledger) to keep the per-version prose narrative discoverable, not just the CHANGELOG facts.
2. **Keep-window = last 2-3 shipped inline.** The most-recent shipped blocks stay inline as "what just landed" context; older shipped blocks move to `roadmap/`. Default keep **3** (mirrors snapshot-retention's recent tier). Configurable later if needed; hardcode-with-named-constant for v1.
3. **Remove the "Recently shipped" mirror** (ROADMAP.md ~490-515) — pure CHANGELOG duplication. The shipped index replaces its navigational purpose.
4. **Reconciliation block-quotes:** prune obviously-stale ones during the dogfood migration; leave current ones. Light touch — not a separate milestone.

## Design (B, mirroring decisions-index)

- **Detection of mode:** `roadmap/` folder exists → index mode (like `decisions/`). Absent → flat mode (all blocks inline; checks skipped).
- **Split unit:** `## v<X.Y.Z> — Title` blocks whose Status is `shipped`. Planned/active/vision blocks NEVER move (they're the forward surface). Only shipped, and only those beyond the keep-window.
- **The ledger is already the index** — but it has no file pointers. The new **shipped index** section in ROADMAP.md adds `- v1.9.0 → roadmap/v1.9.0.md` lines for migrated versions. (Ledger stays as-is; shipped index is the file-pointer layer.)
- **Verb:** `spectacular roadmap migrate [--dry-run] [--keep N]` — snapshot ROADMAP.md, write `roadmap/v*.md` per shipped-beyond-keep block (before touching ROADMAP.md → no data loss on partial run), then rewrite ROADMAP.md with those blocks replaced by index lines. Idempotent (re-run moves only newly-agable blocks). `--keep` overrides the default 3.
- **Doctor `roadmap` area:** when `roadmap/` exists — flag orphan index lines (index → no file), stale files (file → no index line), and shipped blocks still inline beyond keep-window ("N shipped blocks could migrate"). Flat mode: info only.

## Validation

Recorded in `VERIFY-LOG.md`. Success criteria:

- `run: bash tests/cli/roadmap-migrate.test.sh` → all pass (dry-run, migrate, idempotence, doctor: clean/orphan/stale/flat-nudge).
- {assert} `spectacular roadmap migrate` moves shipped-beyond-keep to `roadmap/v*.md`, keeps newest 3 inline, writes `## Shipped` index, leaves planned/active/vision blocks inline.
- {assert} idempotent (re-run = "nothing to migrate"); per-version files written before ROADMAP rewrite (no data loss).
- {assert} `doctor roadmap` detects orphan index lines, stale files, and the flat-mode prune nudge.
- {judge} dogfood: this repo's ROADMAP.md migrated (528 → ~410 lines), `doctor roadmap` clean, no information lost (CHANGELOG + roadmap/v*.md cover the moved prose).
- `run: bash tests/run.sh` → all areas green; `./cli/spectacular doctor` → 0 errors.
