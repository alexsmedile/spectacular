---
status: planned
priority: medium
owner: alex
updated: 2026-06-28
build: b18
summary: "Stop ROADMAP.md bloating agent context: enforce the already-stated 'shipped history lives in CHANGELOG' principle by pruning shipped prose blocks down to their ledger row (the index), and/or add a decisions-index-style roadmap mode (cheap ledger + per-version files). ~49% of ROADMAP.md is currently past-tense duplication."
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../specs/roadmap/SPEC.md
  - ../../ARCHITECTURE.md
depends-on: roadmap-contract-docs
---

# Plan — roadmap-pruning

> **Origin (2026-06-28):** ROADMAP.md is **525 lines, ~49% past-tense** — 12 shipped
> per-version prose blocks (~230 lines) plus a "Recently shipped" section (~26 lines) that
> **duplicates CHANGELOG.md outright**. An agent loading ROADMAP.md to make a planning
> decision pays for all of it. The kicker: both `specs/roadmap/SPEC.md:16` and
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
- **Don't break the prose-block spec.** `specs/roadmap/SPEC.md` specs per-version blocks; whichever approach wins must update that spec, not silently violate it.
- **Reversible / dry-run first.** Show what would be pruned before doing it (mirror the snapshot-retention + archive patterns).

## Understanding

### How it works now (from the audit)

- 525 lines. ~230 in 12 shipped prose blocks (v1.9–v1.22), ~26 in a "Recently shipped" CHANGELOG mirror (ROADMAP.md:490-515, even links to CHANGELOG). Combined ~256 lines ≈ 49% past-tense.
- Stated-but-violated principle: `specs/roadmap/SPEC.md:16` + `ARCHITECTURE.md:257` say shipped history belongs in CHANGELOG. The live file keeps full shipped prose *and* the mirror anyway.
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
- Update `specs/roadmap/SPEC.md` + ARCHITECTURE.md to describe the chosen retention/pruning model (so the stated principle is finally enforced, not just asserted).

### M2 — Detection (doctor)
- `doctor roadmap` (or extend the roadmap area): flag shipped prose blocks beyond the keep-window / the "Recently shipped" duplicate mirror as prunable. Info/warning; relayed by `status`.

### M3 — Prune mechanism
- CLI: prune verb or `doctor --fix roadmap` — snapshot ROADMAP.md, then remove/move shipped blocks (A: delete after CHANGELOG-presence check; B: move to `roadmap/v*.md` + write index). Dry-run default.
- Dogfood: prune THIS repo's ROADMAP.md (12 shipped blocks + the mirror) — confirm CHANGELOG covers each first.
- Tests + VERIFY-LOG.

## M-questions

1. **A vs B?** Lean A (prune-to-ledger; CHANGELOG is already the documented home, laziest, biggest context win). Pick B only if shipped prose holds planning rationale worth keeping structured.
2. **Keep-window:** prune *all* shipped blocks, or keep the last K shipped inline for recent context? Lean: keep the last ~2-3 shipped inline (recent "what just landed" context), prune older. Mirrors the snapshot-retention recent-tier instinct.
3. **The "Recently shipped" mirror (ROADMAP.md:490-515):** remove entirely (pure CHANGELOG dup) — any objection? Lean: remove.
4. **Reconciliation block-quotes** (5 narrative notes, ROADMAP.md:17-21/87/201/299/420): prune the stale ones too, or leave? Out of strict scope; note for cleanup.
