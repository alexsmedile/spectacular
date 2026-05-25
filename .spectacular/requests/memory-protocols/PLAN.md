---
status: planned
priority: high
owner: alex
updated: 2026-05-24
target_version: v1.6.0
summary: "Memory protocols layer on top of v1.5.x soft-DB substrate — semantic recall, large-repo onboarding, principle-context persistence, handoff awareness, drift detection. Scope intentionally open; needs grill before locking milestones."
related:
  - ../../ROADMAP.md
  - ../../SPEC.md
  - ../../requests/soft-db-substrate/PLAN.md
  - ../../../_research/agent-memory/REPORT.md
  - ../../../_research/agent-memory/REPORT-v2.md
---

# Plan — memory-protocols

## 1. Goal

v1.5.x shipped the substrate: soft-folder DBs for memory + sessions, CLI mutators, auto-session-linkage, doctor areas. This request layers **protocols** on top — patterns that turn the substrate into a working memory system for real agentic workflows.

**This is a thinking request, not yet a build request.** Scope is intentionally broad; milestones are unranked drafts. A grill session must lock priorities before milestone breakdown.

Research foundation: `_research/agent-memory/REPORT-v2.md` synthesizes 47 NotebookLM sources + scrapekit subagent output across mem0, MemGPT/Letta, Zep/Graphiti, Cognee, MemPalace, OpenGraph.io, GraphRAG, Cline, Cursor, Anthropic Memory Tool. 8 candidate patterns surfaced.

## 2. Constraints

- **File-first invariant.** Markdown stays source of truth. Any derived index (sqlite, vector, graph DB) is a rebuildable sidecar, never the canonical store.
- **Zero new processes by default.** No Postgres, no Redis, no daemon. Optional external services (OpenGraph.io) are opt-in via env vars and never block local workflows.
- **Composes with v1.5.x.** All new behaviors layer on `DECISIONS.md / MEMORY.md / SESSIONS.md` and their entry folders. No substrate breaking changes.
- **Agentic/mechanical split (v1.4.0).** Skill orchestrates; CLI mutates. New verbs respect this — mining/sweeping/unfurling are mechanical; recall/refine/sync are agentic.
- **Scope must shrink before build.** Five user-named usecase families × eight candidate patterns = too much for one version. Grill needs to cut to 3-4 milestones for v1.6.0; rest go to v1.6.x/v1.7.x.

## 3. Usecase families (user-prioritized 2026-05-24)

These five drive the scope. Phrased as user-spoken pain → what the protocol fixes.

### F1 — Session + checkpoint recall
> "I want to recall sessions and requests/checkpoints/milestones."

**Pain today:** session history exists as files but is not addressable. `spectacular recall` is a v1.6.x roadmap stub; no implementation, no semantic match, no temporal filter.

**Candidate patterns:** bi-temporal frontmatter (#1), `recall` verb with tag+date scan (Tier 0), MemPalace mining for verbatim history (#6), two-phase extract→dedupe on close (#4).

### F2 — Large codebase surfing
> "Memory useful for surfing larger codebases."

**Pain today:** SPEC.md is human-maintained and decoupled from code. Agent has no map of "what symbol relates to what spec/decision/memory" without grep loops.

**Candidate patterns:** triples in frontmatter (#8), Aider-style PageRank repo map (REPORT v1 §3), spec sync hook on archive.

### F3 — PRINCIPLES.md context persistence
> "Keep context for files like PRINCIPLES.md throughout the whole session."

**Pain today:** PRINCIPLES.md loaded into context at session start but drifts out of attention as session grows. No re-anchoring mechanism.

**Candidate patterns:** path/glob activation on memory entries (#2), procedural memory refinement loop (#5), "always-pinned" frontmatter on canonical docs (Cursor `alwaysApply` pattern).

### F4 — Larger-scale handoffs
> "Larger scale handoffs."

**Pain today:** handoff = "read PRD + PLAN + scroll backlog." No condensed briefing, no "what's currently loaded," no cross-agent provenance.

**Candidate patterns:** `activeContext.md` pattern (Cline Memory Bank, REPORT v1 §3), MemPalace mining of recent transcripts (#6), Letta-style memory blocks (#3), `agent:` provenance field (deferred from v1.5.0 D8).

### F5 — Onboarding on large repos
> "Onboarding on a larger repo + AGENTS.md + README.md. MEMORY.md will become instructions on how to load the system or relevant files."

**Pain today:** AGENTS.md is static; doesn't reflect what the live system looks like. MEMORY.md (v1.5.0) is a flat table — not a load-order guide.

**Candidate patterns:** MEMORY.md as **load manifest** (new framing — semantic upgrade beyond v1.5.x index), nested AGENTS.md per `.spectacular/<dir>/` (Cursor nested-rules), path-activation (#2), `spectacular tour` walkthrough verb.

### F6 — SPEC drift detection
> "Detecting drifts from SPECS.md."

**Pain today:** SPEC.md describes what's built. Code changes silently invalidate bullets. No alarm bell.

**Candidate patterns:** Aider PageRank + symbol-extraction (REPORT v1 §3), triples linking spec bullets to code paths (#8), `spectacular spec sync --from-code` verb, doctor drift area.

## 4. Cross-cutting patterns from REPORT-v2

The 8 patterns the research surfaced. Numbering matches `_research/agent-memory/REPORT-v2.md`.

1. **Bi-temporal frontmatter on decisions** — `valid_from`, `superseded_by:`. Graphiti pattern. Zero infra. → F1, F6.
2. **Path/glob activation on entries** — `paths:`, `when_to_use:`. Cursor Rules pattern. → F3, F5.
3. **Self-edit verbs** — `memory append|supersede|link`. Letta pattern. Constrains how agents mutate memory. → F3, F4.
4. **Two-phase extract→dedupe `remember`** — skill-side LLM dedupe loop. mem0 pattern. → F1, F4.
5. **Procedural-memory refinement loop** — `principles refine --from sessions/`. LangMem pattern. → F3, F6.
6. **`spectacular mine <dir> --mode convos`** — MemPalace pattern. Backfill memory from raw transcripts. → F1, F4, F5.
7. **`spectacular unfurl <url>`** — OpenGraph.io enrichment. Link metadata + screenshot in frontmatter. → F2 (sort of), F6 (when linking to external docs).
8. **`triples:` frontmatter + `spectacular query --depends-on X`** — GraphRAG + OGB pattern. Multi-hop reasoning via grep+yq. → F2, F6.

## 5. New protocol ideas (not in REPORT — emerged from usecases)

These are **user-suggested ideas** that don't map cleanly to a single REPORT pattern. Each needs grill before becoming a milestone.

### I1 — MEMORY.md as load-manifest, not just index
Reframe MEMORY.md from "table of entries" → "instructions for what to load when." Frontmatter would change from `type: index` to something like `type: load-manifest`. Body becomes a series of `## When working on <X>, load <Y, Z>` blocks generated from entry frontmatter's `paths:` / `when_to_use:` fields.

**Open question:** does this replace the v1.5.0 index format, or add a second view (`MEMORY.md` index + `LOADER.md` manifest)?

### I2 — `spectacular onboard` verb
A walkthrough for cold-start onboarding. Reads AGENTS.md + README.md + SPEC.md + recent sessions, produces a structured briefing tailored to who's onboarding (human / Claude / Codex). Could ride on `spectacular tour` from F5.

**Open question:** is this a separate verb or does it fold into the existing `status` skill flow?

### I3 — Always-pinned canonical docs
Frontmatter flag like `pinned: true` on PRINCIPLES.md / ARCHITECTURE.md / STACK.md that the skill router treats as "re-inject every N turns" rather than "load once." Cursor's `alwaysApply: true` is the precedent.

**Open question:** which docs default to pinned? Is pinning per-doc-type or per-file?

### I4 — Spec-drift watchdog
A `doctor specs --drift` area that compares SPEC.md bullets against code symbols (tree-sitter or simple grep) and flags bullets with no corresponding code reference (and vice versa).

**Open question:** does this require tree-sitter as a dependency, or can we get 80% from path-presence + symbol-name fuzzy match?

### I5 — Checkpoint primitives
Sessions are time-bounded; **checkpoints** would be milestone-bounded — "v1.5.0 ships" is a checkpoint that bundles N decisions + M memories + ~K sessions. A `spectacular checkpoint <name>` verb would snapshot a slice of substrate for retro/handoff purposes.

**Open question:** is a checkpoint a real new doc-type, or just a tag on existing entries that the index queries surface?

### I6 — Cross-repo memory inheritance
For F5 (large-repo onboarding): user-scope memories in `~/.spectacular/memory/` could carry `applies-to: [project-pattern]` and auto-attach when working on matching projects. Pattern from MemPalace's "Halls" connecting Wings.

**Open question:** opt-in per-project? Auto-include based on glob match? Permission prompt on first load?

## 6. Open questions for the grill

Before milestone breakdown:

1. **Scope cut** — REPORT-v2 Option A (patterns 1, 2, 3, 6) vs Option B (+ 7, 8) vs custom mix anchored to F1-F6?
2. **MemPalace stance** — hard dependency / soft adapter / absorb-natively-in-bash? MIT-licensed, ~50 lines for the `sweep` equivalent.
3. **OpenGraph.io stance** — ship `unfurl` as core or as opt-in plugin? Env-var gate?
4. **MEMORY.md re-framing (I1)** — does it stay an index or become a load-manifest?
5. **Triples adoption (#8)** — frontmatter schema OK to invent now, or wait for a real query verb to drive the schema?
6. **Pinned docs (I3)** — which docs default to pinned? PRD, PRINCIPLES, ARCHITECTURE, STACK all candidates.
7. **Drift detection scope (I4 / F6)** — what's "good enough"? Path-existence vs symbol-fuzzy vs tree-sitter?
8. **Checkpoint primitive (I5)** — new doc-type or tag?
9. **Cross-repo memory (I6)** — user-scope `applies-to` matcher pattern?
10. **Recall verb shape** — `spectacular recall "query"` (semantic) vs `spectacular recall --tag x --since 7d` (faceted) vs both?
11. **Provenance** — add `agent:` field now (was deferred from v1.5.0 D8) or wait until multi-agent advisory work in v1.7.x?
12. **Mining target** — `~/.claude/projects/` only, or also Codex/Cursor transcripts? Tool-agnostic by default?

## 7. Non-goals (locked)

- **Vector embeddings as default.** Tier 4 (sqlite-vec / shodh-memory) is deferred until tag-based retrieval demonstrably fails on this repo.
- **Knowledge graph DB.** Bi-temporal frontmatter + triples-in-YAML captures 80% without Neo4j/Postgres.
- **Multi-agent coordination.** v1.7.x advisory work, not here. This request stays single-agent.
- **DECISIONS.md → folder migration.** Deferred to a separate v1.6.x request alongside query verbs.
- **Cloud sync.** Memory is git-committed and team-visible; that's the sync mechanism.

## 8. Dependencies

- v1.5.x soft-db-substrate (in review) — MUST ship and be dogfooded for ≥1 week before this request opens. The substrate needs real usage to surface which patterns matter most.
- v1.6.x query verbs (`spectacular recall`, `spectacular decisions --since`, `spectacular sessions`) — currently in ROADMAP but unscoped. Some of this request's patterns (F1) may absorb that scope.

## 9. Tasks

Deliberately empty. Open via grill once usecase priorities are confirmed.

See `TASKS.md` for the placeholder milestone shell.

## 10. References

- `_research/agent-memory/REPORT.md` — v1 synthesis (top 5 patterns)
- `_research/agent-memory/REPORT-v2.md` — v2 synthesis (top 8 patterns + 5-tier ladder + MemPalace deep-dive)
- `_research/agent-memory/q1-v2-architectures.md` — restated architecture ranking
- `_research/agent-memory/q5-deepdive-mempalace-opengraph.md` — MemPalace + OpenGraph + triples deep-dive
- `_research/agent-memory/sources.txt` — 47 NotebookLM sources
- NotebookLM notebook: `memresearch` (alias for `7de5c309-3eec-4428-a562-64559150e84d`)
- NotebookLM conversation: `9fa79a8a-7fe3-4715-a648-3cc0f51d259b` (warm for follow-ups)
