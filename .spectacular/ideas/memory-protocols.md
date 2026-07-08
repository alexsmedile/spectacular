---
type: idea
status: parked
priority: high
owner: alex
updated: 2026-05-26
origin: requests/memory-protocols (abandoned 2026-05-26 — too broad; a narrower spec request will replace it)
related:
  - ../roadmaps/index.md
  - ../specs/index.md
  - ../../_research/agent-memory/REPORT.md
  - ../../_research/agent-memory/REPORT-v2.md
---

# Idea — memory protocols (parked thinking-doc)

> **Status note (2026-05-26):** the `memory-protocols` request was abandoned because its scope was too broad to ship as a single milestone (5-6 usecase families × 8 candidate patterns × 6 new ideas = no closing surface). This document collects the full thinking so a **narrower, spec-shaped request** can be cut from it without losing the research.
>
> All content below is the merged contents of the abandoned request's `PLAN.md`, `TASKS.md`, and `IDEAS_MEMORY.md`. Treat it as raw material for the next request, not as a plan.

---

## 1. Goal (carried over from PLAN)

v1.5.x shipped the substrate: soft-folder DBs for memory + sessions, CLI mutators, auto-session-linkage, doctor areas. The original `memory-protocols` request set out to layer **protocols** on top — patterns that turn the substrate into a working memory system for real agentic workflows.

Research foundation: `_research/agent-memory/REPORT-v2.md` synthesizes 47 NotebookLM sources + scrapekit subagent output across mem0, MemGPT/Letta, Zep/Graphiti, Cognee, MemPalace, OpenGraph.io, GraphRAG, Cline, Cursor, Anthropic Memory Tool. 8 candidate patterns surfaced.

## 2. Constraints (carried over from PLAN)

- **File-first invariant.** Markdown stays source of truth. Any derived index (sqlite, vector, graph DB) is a rebuildable sidecar, never the canonical store.
- **Zero new processes by default.** No Postgres, no Redis, no daemon. Optional external services (OpenGraph.io) are opt-in via env vars and never block local workflows.
- **Composes with v1.5.x.** All new behaviors layer on `DECISIONS.md / MEMORY.md / SESSIONS.md` and their entry folders. No substrate breaking changes.
- **Agentic/mechanical split (v1.4.0).** Skill orchestrates; CLI mutates. New verbs respect this — mining/sweeping/unfurling are mechanical; recall/refine/sync are agentic.
- **Scope must shrink before build.** This is the reason this got parked — five user-named usecase families × eight candidate patterns is too much for one version.

## 3. Usecase families (user-prioritized 2026-05-24)

Phrased as user-spoken pain → what the protocol would fix.

### F1 — Session + checkpoint recall
> "I want to recall sessions and requests/checkpoints/milestones."

**Pain today:** session history exists as files but is not addressable. `spectacular recall` is a v1.6.x roadmap stub; no implementation, no semantic match, no temporal filter.

**Candidate patterns:** bi-temporal frontmatter (#1), `recall` verb with tag+date scan (Tier 0), MemPalace mining for verbatim history (#6), two-phase extract→dedupe on close (#4).

### F2 — Large codebase surfing
> "Memory useful for surfing larger codebases."

**Pain today:** SPEC.md is human-maintained and decoupled from code. Agent has no map of "what symbol relates to what spec/decision/memory" without grep loops.

**Candidate patterns:** triples in frontmatter (#8), Aider-style PageRank repo map (REPORT v1 §3), spec sync hook on archive, Graphify-style code-graph snapshot (see § Graphify below).

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

**Candidate patterns:** MEMORY.md as **load manifest** (new framing — semantic upgrade beyond v1.5.x index), nested AGENTS.md per `.spectacular/<dir>/` (Cursor nested-rules), path-activation (#2), `spectacular tour` walkthrough verb, Graphify "god nodes" as cold-start reading order.

### F6 — SPEC drift detection
> "Detecting drifts from SPECS.md."

**Pain today:** SPEC.md describes what's built. Code changes silently invalidate bullets. No alarm bell.

**Candidate patterns:** Aider PageRank + symbol-extraction (REPORT v1 §3), triples linking spec bullets to code paths (#8), `spectacular spec sync --from-code` verb, doctor drift area, Graphify diff between successive `graph.json` snapshots.

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

Each needs grill before becoming a milestone.

### I1 — MEMORY.md as load-manifest, not just index
Reframe MEMORY.md from "table of entries" → "instructions for what to load when." Frontmatter would change from `type: index` to something like `type: load-manifest`. Body becomes a series of `## When working on <X>, load <Y, Z>` blocks generated from entry frontmatter's `paths:` / `when_to_use:` fields.

**Refinement (2026-05-26):** borrow Graphify's "god nodes" concept — surface the most-connected/referenced docs at the top of `MEMORY.md` as the default cold-start reading order.

**Open question:** does this replace the v1.5.0 index format, or add a second view (`MEMORY.md` index + `LOADER.md` manifest)?

### I2 — `spectacular onboard` verb
A walkthrough for cold-start onboarding. Reads AGENTS.md + README.md + SPEC.md + recent sessions, produces a structured briefing tailored to who's onboarding (human / Claude / Codex). Could ride on `spectacular tour` from F5.

**Open question:** is this a separate verb or does it fold into the existing `status` skill flow?

### I3 — Always-pinned canonical docs
Frontmatter flag like `pinned: true` on PRINCIPLES.md / ARCHITECTURE.md / STACK.md that the skill router treats as "re-inject every N turns" rather than "load once." Cursor's `alwaysApply: true` is the precedent.

**Open question:** which docs default to pinned? Is pinning per-doc-type or per-file?

### I4 — Spec-drift watchdog
A `doctor specs --drift` area that compares SPEC.md bullets against code symbols (tree-sitter or simple grep) and flags bullets with no corresponding code reference (and vice versa).

**Open question:** does this require tree-sitter as a dependency, or can we get 80% from path-presence + symbol-name fuzzy match? Or — can a Graphify rebuild + diff *be* the drift detector?

### I5 — Checkpoint primitives
Sessions are time-bounded; **checkpoints** would be milestone-bounded — "v1.5.0 ships" is a checkpoint that bundles N decisions + M memories + ~K sessions. A `spectacular checkpoint <name>` verb would snapshot a slice of substrate for retro/handoff purposes.

**Open question:** is a checkpoint a real new doc-type, or just a tag on existing entries that the index queries surface?

### I6 — Cross-repo memory inheritance
For F5 (large-repo onboarding): user-scope memories in `~/.spectacular/memory/` could carry `applies-to: [project-pattern]` and auto-attach when working on matching projects. Pattern from MemPalace's "Halls" connecting Wings.

**Open question:** opt-in per-project? Auto-include based on glob match? Permission prompt on first load?

## 6. Open questions (the 13 that blocked milestone breakdown)

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
13. **Graphify integration** — should code-graph snapshots be a Tier-2 recommended pattern for F2/F5? Trigger threshold (file count? token-burn signal?)? Stale-snapshot policy? Doctor warning when `GRAPH_REPORT.md` is older than the SPEC tag? See § Graphify below for the full pros/cons + compression matrix.

## 7. Non-goals (the things we already locked out)

- **Vector embeddings as default.** Tier 4 (sqlite-vec / shodh-memory) is deferred until tag-based retrieval demonstrably fails on this repo.
- **Knowledge graph DB.** Bi-temporal frontmatter + triples-in-YAML captures 80% without Neo4j/Postgres.
- **Multi-agent coordination.** v1.7.x advisory work, not here.
- **DECISIONS.md → folder migration.** Separate v1.6.x request alongside query verbs.
- **Cloud sync.** Memory is git-committed and team-visible; that's the sync mechanism.

## 8. Dependencies

- v1.5.x soft-db-substrate (shipped) — dogfood ≥1 week before the next request opens.
- v1.6.x query verbs (`spectacular recall`, `spectacular decisions --since`, `spectacular sessions`) — currently in ROADMAP but unscoped. Some F1 patterns may absorb that scope.

---

# Memory tool selection by corpus type (original IDEAS_MEMORY content)

Captured from a side-question during memory-protocols scoping (2026-05-24): "what is the best mem tool for a second brain Obsidian database with lots of ideas in an inbox folder; or for a business dir with lots of SOPs and know-how docs?"

The relevance: Spectacular's substrate is one of three distinct corpus types the same research touches. Naming the others sharpens what Spectacular's memory layer is and isn't.

## Core insight

The corpus shape determines the tool:

| Corpus | Write pattern | Query pattern | Durability |
|---|---|---|---|
| Spectacular workspace | write-heavy, fast, agent-curated | "what did we decide, when, why" — provenance + temporal | git-committed, team-visible |
| Obsidian vault (inbox + ideas) | append-heavy, slow, human-curated | "what's similar to this thought" — semantic + associative | personal, evolving |
| Business SOP / know-how | rare writes, heavy curation, multi-author | "what depends on this policy" — relational + versioned | compliance, audit-trail |

Same research, three different tool stacks. The portable patterns (work for all three): **bi-temporal frontmatter**, **path/glob activation**, **markdown + frontmatter as schema**. Everything else is corpus-specific.

## For the Obsidian vault (inbox + ideas — second brain)

**Best fit: MemPalace + Smart Connections (already installed).**

- **Smart Connections** is the right primary tool. It's the only thing in the research actually built for Obsidian. Vector embeddings on idea-heavy text is where vectors genuinely win — semantic similarity matters more than provenance for ideation. The vault's `inbox/ → spaces/ → projects/` flow is exactly the corpus shape Smart Connections targets. Current state: vector index at `.smart-env/`, 3,389 files (per CLAUDE.md).
- **MemPalace as a secondary layer for capture** (not search). Its `mine` flow can ingest Claude Code / chat transcripts into `inbox/` automatically. Wing = vault folder, Drawer = verbatim note, Hall = wikilink. Mapping is clean.
- **Skip:** mem0 (cloud-first, opinionated dedupe — kills idea diversity), Letta (agent runtime, overkill), Zep/Graphiti (temporal KG doesn't help inbox triage), GraphRAG triples (great for code, friction for prose).

**Concrete stack:** Smart Connections for search/retrieval + MemPalace MCP for transcript ingestion + existing `core/tools/` for normalization. Add `_research/`-style frontmatter discipline to the inbox so future MemPalace mining has signal.

## For the business SOP / know-how directory

**Best fit: GraphRAG triples + sqlite-fts5 (Tier 2 + Tier 3 from REPORT-v2).**

SOPs are **structured, relational, and durable** — the opposite of inbox ideation. The questions asked are "what depends on this policy?", "which SOPs reference the refund process?", "what changed when we updated the onboarding doc?" — those are graph queries, not semantic queries.

- **GraphRAG triples in YAML frontmatter.** Example: `triples: [{head: "REFUND-SOP", relation: "REQUIRES", tail: "VERIFICATION-POLICY"}]`. One `yq` query traces dependencies. Bi-temporal frontmatter (`valid_from`, `superseded_by:`) handles policy revisions — critical for compliance.
- **sqlite-fts5 sidecar** when grep gets slow on long-form SOPs. Local, ephemeral, rebuildable.
- **Anthropic Memory Tool API** if Claude agents work with the SOP dir — its file-based memory pattern aligns natively with markdown SOPs.
- **Skip:** MemPalace (verbatim-storage philosophy clashes with the heavy curation SOPs need), mem0 (dedupe is wrong for canonical docs — you WANT versions), vectors as primary (semantic similarity isn't the question being asked).

**Concrete stack:** triples-in-frontmatter as the relational layer + sqlite-fts5 for full-text + bi-temporal versioning + git for audit trail. Optionally OpenGraph.io to snapshot any external-link references at SOP-publish time (link-rot is real in compliance contexts).

## What this means for Spectacular

1. **Spectacular's stack converges with the SOP stack more than the Obsidian stack.** Both want triples + bi-temporal + git audit. This validates patterns #1 and #8 from REPORT-v2 as high-priority.
2. **Smart Connections shouldn't be a Spectacular pattern.** Right answer for the vault, wrong for Spectacular (workspace context is too volatile for stable embeddings; questions are provenance-heavy not similarity-heavy). Vectors stay in Tier 4 — defer until tag/triple retrieval demonstrably fails.
3. **MemPalace mining (#6) generalizes across all three corpora.** It's the universal capture primitive — verbatim transcripts → markdown drawers. Strong signal it belongs in core memory work, not as an optional add-on.
4. **Open question:** is Spectacular a workspace tool that happens to have memory, or a memory tool that happens to scaffold workspaces? The Obsidian/SOP split surfaces this — different framings would justify different scope cuts.
5. **Bi-product idea:** if Spectacular's substrate is corpus-agnostic enough to work for SOP dirs, that's a parallel adoption surface — `spectacular init --kit sops` or `--kit second-brain` as future kits.

---

# Graphify — candidate memory pattern (added 2026-05-26)

Researched via the NotebookLM `Agent Memory Systems — Spectacular Research` notebook (sources: safishamsi/graphify README, `docs/how-it-works.md` v8, DEV.to Graphify+code-review-graph article, Graphify YouTube intro, GraphRAG-Pureinsights, Zep arXiv 2501.13956).

## What it is

CLI + AI-assistant skill (`/graphify` or `graphify .` in PowerShell). Ingests a folder of mixed inputs — source code (25 languages via tree-sitter), SQL schemas, shell/R scripts, markdown, PDFs, images, video, audio — and emits **three artifacts** into `graphify-out/`:

- `graph.html` — interactive browser-viewable graph
- `GRAPH_REPORT.md` — extracted highlights + surprising connections + suggested questions
- `graph.json` — the full graph, queried by the agent on subsequent turns instead of re-reading raw files

The build runs in three passes:

1. **Code structure (local, free)** — tree-sitter extracts classes, functions, imports, call graphs, comments; SQL gets tables/views/FKs/JOINs deterministically. No LLM. If the corpus is code-only, Pass 3 is skipped entirely.
2. **Video/audio (local, free)** — `faster-whisper` transcription, prompt-seeded with top "god nodes" from Pass 1 to bias toward your domain. Cached.
3. **Docs/papers/images (Claude subagents, costs tokens)** — Claude runs in parallel, each subagent outputs JSON fragments (nodes, edges, group relationships) merged into one graph.

**Confidence-tagged edges**: `EXTRACTED` (1.0), `INFERRED` (discrete 0.55–0.95 rubric), `AMBIGUOUS` (flagged for review). Communities found via the **Leiden algorithm** — no embeddings, no vector DB; semantic-similarity edges Claude extracts (`semantically_similar_to`) directly shape clustering.

**Reported token saving**: "71.5× fewer tokens per query" on a mixed 52-file corpus vs reading raw files. First-run cost is real; amortizes over re-queries.

## Pros

| | |
|---|---|
| **Zero infra** | No DB, no embedding service, no vector store. Three files on disk, queryable from any agent. |
| **Code-first pipeline is free** | Tree-sitter pass needs no LLM. For pure-code repos, the whole graph is built locally at near-zero cost. |
| **Multi-platform skill** | Ships installers for Claude Code, Codex, Cursor, Gemini CLI, Aider, Copilot CLI, VS Code, OpenCode, Kiro, and ~10 others. Project-scoped install (`--project`) writes into `.claude/skills/` or `.agents/skills/` — same install surface Spectacular uses. |
| **Confidence tagging is honest** | Every edge knows whether it's a fact or a guess. Spectacular's `provenance:` and lifecycle states are kindred ideas. |
| **Graph-native, not vector-native** | Aligns with the "questions are provenance-heavy, not similarity-heavy" framing above. The graph **is** the index. |
| **Mature** | 53k stars, 537 commits, v0.8.18, v8 branch active; Neo4j export, MCP stdio server, Mermaid call-flow HTML export — production-shaped. |

## Cons

| | |
|---|---|
| **Stale by default** | The graph is a snapshot. Re-running `/graphify` is manual. No file-watcher, no incremental rebuild called out in the README. For an active workspace where PRD/SPEC change daily, this is a real cost (token-wise + workflow-wise). |
| **First-run cost on doc-heavy corpora** | Pass 3 uses Claude subagents in parallel — on a docs+papers+images corpus this is non-trivial tokens before the savings start. |
| **Token-saving claim is unverified at our scale** | "71.5×" is one mixed corpus of 52 files. Spectacular workspaces are typically <100 markdown files — much smaller corpus, fewer reads per file. The saving could be far less compelling here. |
| **No durable preferences / episodic memory** | Graphify maps **the corpus**, not the conversation. Does not solve session recall, checkpoint restoration, user-preference durability (F1, F3). |
| **No drift detection by itself** | A graph snapshot can't tell you the code drifted from the SPEC unless you re-run and diff. Useful **inside** F6, not a solution to F6. |
| **PyPI quirk** | Official package is `graphifyy` (double-y) — name-squatting on `graphify` exists. Minor but operationally annoying. |
| **Adds another runtime** | Python 3.10+, `uv` or `pipx`, optional extras for PDF/office/video/SQL/neo4j. Each repo needs the install or a one-time global. |
| **Marketing density** | Star count and breadth of platform support smell like growth-hacking. Worth a second-look review of *actual* maintenance vs *apparent* maintenance. |

## Compression against other candidates

| Solution | Storage shape | Stale handling | What it answers | Where it fits in Spectacular |
|---|---|---|---|---|
| **Graphify** | Local JSON+HTML graph, snapshot | Manual re-run | "How does the codebase fit together; what depends on what" | **F2**, **F5**, optionally F6. Layer **on top of** soft-DB. |
| **MEMORY.md / SESSIONS.md** (v1.5.0) | Markdown + frontmatter, git-committed | Continuous via `spectacular remember`, `session start/end` | "What did we decide; what was active when" | Core substrate. Always present. |
| **MemPalace** | Verbatim markdown drawers + wikilink "halls" | Append-only mining from transcripts | "What did the user say verbatim; capture without interpretation" | F1, F4. Complements curated decisions with raw transcript trail. |
| **Mem0** | Cloud (or self-host), opinionated dedupe | Continuous, automatic | "What does the user generally want / prefer" | F3 (durable preferences). Conflicts with git-committed/team-visible philosophy — likely skip. |
| **Zep / Graphiti** | Temporal knowledge graph (Neo4j-backed) | Bi-temporal — valid time vs ingestion time | "What was true when; what changed over time" | F6, F4. Heavier infra than v1.6 wants. |
| **GraphRAG** (Pureinsights framing) | LLM-extracted triples + community summaries | Manual re-run, like Graphify | "Multi-hop reasoning across a corpus" | Same niche as Graphify — Graphify is a packaged version of this. |
| **Smart Connections / vectors** | Embedding index (`.smart-env/`) | Continuous re-index on file save | "What's semantically similar to this" | Wrong tool for Spectacular. Right tool for the Obsidian vault. |
| **CLAUDE.md / AGENTS.md** | Static markdown, loaded into context | Manual edits | "What are the rules; what's the orientation" | Already used. Floor, not ceiling. |

**Key compression**: Graphify and MEMORY.md/SESSIONS.md are **orthogonal**, not competing. Graphify maps the **corpus**; soft-DB tracks the **conversation and decisions**. The interesting question isn't "Graphify *or* soft-DB" — it's "does Graphify earn its install in addition to soft-DB?"

## Fit for Spectacular — recommendation

**Adopt as an optional Tier-2 pattern for F2/F5, do not bundle into core.**

1. **F2 (large codebase surfing)** — Graphify's killer use case. A new agent on a large repo can read `GRAPH_REPORT.md` + query `graph.json` instead of grepping. Real value.
2. **F5 (onboarding cold)** — same logic. `spectacular onboard` (idea I2) could optionally include `graphify .` as a step when the repo is over some file-count threshold.
3. **Do not** include Graphify in the always-set kit. Most Spectacular workspaces are markdown-heavy and small — the snapshot would be stale within hours and the token saving doesn't materialize.
4. **Document the integration shape, don't ship the runtime.** Spectacular's job is to know Graphify exists and reference its outputs (`graphify-out/`) when present; not to install or invoke it. Same posture as `pageworks` post-v1.2.0.
5. **Optional convention pack** — `convention_pack: graphify-aware` could add a `graphify-out/` rule to `.gitignore` (or keep it committed, user choice) and add a `doctor` warning when `GRAPH_REPORT.md` is older than the SPEC tag.

## Open questions for the next spec request

1. **Threshold for recommending Graphify** — file count? Token budget burned on file reads in a recent session? Both?
2. **Stale-snapshot policy** — if Spectacular surfaces Graphify outputs to agents, does the doctor flag staleness? Hours? Days?
3. **Cost transparency** — Pass 3 token cost should be reported and accepted before run. Should `spectacular` proxy this with a dry-run estimate?
4. **Overlap with `spec-drift watchdog` (I4)** — can a Graphify rebuild + diff against the previous `graph.json` *be* the drift detector, or is that overkill?
5. **Cross-repo memory (I6)** — Graphify is single-repo. If we want cross-repo, what's the merge story for multiple `graph.json` files?
6. **What `MEMORY.md` learns from Graphify** — the **load-manifest** idea (I1) could borrow Graphify's "god nodes" concept.
7. **Trust posture** — Graphify ships with `git add` hints for committable artifacts. Do we recommend committing `graph.json`? It's deterministic-ish, but Pass 3 LLM output isn't reproducible. Probably no — `.gitignore` it.

---

# References

- `_research/agent-memory/REPORT.md` — v1 synthesis (top 5 patterns)
- `_research/agent-memory/REPORT-v2.md` — v2 synthesis (top 8 patterns + 5-tier ladder + MemPalace deep-dive)
- `_research/agent-memory/q1-v2-architectures.md` — restated architecture ranking
- `_research/agent-memory/q5-deepdive-mempalace-opengraph.md` — MemPalace + OpenGraph + triples deep-dive
- `_research/agent-memory/sources.txt` — 47 NotebookLM sources
- NotebookLM notebook: `memresearch` (alias for `7de5c309-3eec-4428-a562-64559150e84d`)
- NotebookLM conversation: `9fa79a8a-7fe3-4715-a648-3cc0f51d259b` (warm for follow-ups)
- Graphify — https://github.com/safishamsi/graphify
- NotebookLM report artifact: `8a40a36c-0646-4ef0-ad60-503b0d92cbcc` ("Research Note: Graphify as a Candidate Memory Pattern for AI Coding Agents", generated 2026-05-26)

---

# Hook for the next request

When you cut a narrower spec request from this material, the recommended starting moves:

- **Pick one usecase family** (F1-F6) as the spec's anchor. Multiple families → broad scope → same trap.
- **Pick at most two patterns** from §4 + §5 to implement against that family.
- **Lock the non-goals from §7 as-is** — they were resolved decisions, not open questions.
- **Defer §6 questions that don't bear on the chosen family.** Most are scope-cuts in disguise.
- **Graphify and Smart Connections stay external** — Spectacular references their outputs, never installs them. This is a posture decision (matches the pageworks-migration precedent), not an open question.
