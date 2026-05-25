---
status: planned
updated: 2026-05-24
related:
  - PLAN.md
  - ../../../_research/agent-memory/REPORT-v2.md
---

# Ideas — memory tool selection by corpus type

Captured from a side-question during memory-protocols scoping (2026-05-24): "what is the best mem tool for a second brain Obsidian database with lots of ideas in an inbox folder; or for a business dir with lots of SOPs and know-how docs?"

The relevance to `memory-protocols`: Spectacular's substrate is one of three distinct corpus types the same research touches. Naming the others sharpens what Spectacular's memory layer is and isn't. Worth folding into the grill when scoping I6 (cross-repo memory inheritance) and the broader "is Spectacular a memory tool or a workspace tool" framing.

---

## Core insight

The corpus shape determines the tool:

| Corpus | Write pattern | Query pattern | Durability |
|---|---|---|---|
| Spectacular workspace | write-heavy, fast, agent-curated | "what did we decide, when, why" — provenance + temporal | git-committed, team-visible |
| Obsidian vault (inbox + ideas) | append-heavy, slow, human-curated | "what's similar to this thought" — semantic + associative | personal, evolving |
| Business SOP / know-how | rare writes, heavy curation, multi-author | "what depends on this policy" — relational + versioned | compliance, audit-trail |

Same research, three different tool stacks. The portable patterns (work for all three): **bi-temporal frontmatter**, **path/glob activation**, **markdown + frontmatter as schema**. Everything else is corpus-specific.

---

## For the Obsidian vault (inbox + ideas — second brain)

**Best fit: MemPalace + Smart Connections (already installed).**

- **Smart Connections** is the right primary tool. It's the only thing in the research actually built for Obsidian. Vector embeddings on idea-heavy text is where vectors genuinely win — semantic similarity matters more than provenance for ideation. The vault's `inbox/ → spaces/ → projects/` flow is exactly the corpus shape Smart Connections targets. Current state: vector index at `.smart-env/`, 3,389 files (per CLAUDE.md).
- **MemPalace as a secondary layer for capture** (not search). Its `mine` flow can ingest Claude Code / chat transcripts into `inbox/` automatically — same MemPalace primitive that fits Spectacular, applied to the vault's intake side. Wing = vault folder, Drawer = verbatim note, Hall = wikilink. Mapping is clean.
- **Skip:** mem0 (cloud-first, opinionated dedupe — kills idea diversity), Letta (agent runtime, overkill), Zep/Graphiti (temporal KG doesn't help inbox triage), GraphRAG triples (great for code, friction for prose).

**Concrete stack:** Smart Connections for search/retrieval + MemPalace MCP for transcript ingestion + existing `core/tools/` for normalization. Add `_research/`-style frontmatter discipline to the inbox so future MemPalace mining has signal.

---

## For the business SOP / know-how directory

**Best fit: GraphRAG triples + sqlite-fts5 (Tier 2 + Tier 3 from REPORT-v2).**

SOPs are **structured, relational, and durable** — the opposite of inbox ideation. The questions asked are "what depends on this policy?", "which SOPs reference the refund process?", "what changed when we updated the onboarding doc?" — those are graph queries, not semantic queries.

- **GraphRAG triples in YAML frontmatter.** Example: `triples: [{head: "REFUND-SOP", relation: "REQUIRES", tail: "VERIFICATION-POLICY"}]`. One `yq` query traces dependencies. Bi-temporal frontmatter (`valid_from`, `superseded_by:`) handles policy revisions — critical for compliance.
- **sqlite-fts5 sidecar** when grep gets slow on long-form SOPs. Local, ephemeral, rebuildable.
- **Anthropic Memory Tool API** if Claude agents work with the SOP dir — its file-based memory pattern aligns natively with markdown SOPs.
- **Skip:** MemPalace (verbatim-storage philosophy clashes with the heavy curation SOPs need), mem0 (dedupe is wrong for canonical docs — you WANT versions), vectors as primary (semantic similarity isn't the question being asked).

**Concrete stack:** triples-in-frontmatter as the relational layer + sqlite-fts5 for full-text + bi-temporal versioning + git for audit trail. Optionally OpenGraph.io to snapshot any external-link references at SOP-publish time (link-rot is real in compliance contexts).

---

## What this means for the Spectacular `memory-protocols` request

1. **Spectacular's stack converges with the SOP stack more than the Obsidian stack.** Both want triples + bi-temporal + git audit. This validates patterns #1 and #8 from REPORT-v2 as high-priority.

2. **Smart Connections shouldn't be a Spectacular pattern.** It's the right answer for the vault but wrong for Spectacular (workspace context is too volatile for stable embeddings; questions are provenance-heavy not similarity-heavy). Vectors stay in Tier 4 — defer until tag/triple retrieval demonstrably fails.

3. **MemPalace mining (#6) generalizes across all three corpora.** It's the universal capture primitive — verbatim transcripts → markdown drawers. This is a strong signal it belongs in `memory-protocols` core, not as an optional add-on. Worth promoting in the M0 grill.

4. **Open question for grill:** is Spectacular a workspace tool that happens to have memory, or a memory tool that happens to scaffold workspaces? The Obsidian/SOP split surfaces this — different framings would justify different scope cuts for v1.6.0.

5. **Bi-product idea:** if Spectacular's substrate is corpus-agnostic enough to work for SOP dirs, that's a parallel adoption surface — `spectacular init --kit sops` or `--kit second-brain` as future kits. Probably not v1.6.0 scope, but worth tracking as a follow-on hypothesis.

---

## Follow-up question raised but not answered

> "Want me to draft separate `memory-protocols` requests for the Obsidian vault and the business dir? Those would live in their own repos with their own `.spectacular/` workspaces. Or are they out-of-scope for this skill and just exploratory?"

**Status:** unanswered, deferred. Treat as exploratory unless user confirms otherwise. If pursued, those become separate request slugs in separate workspaces, not part of `memory-protocols` here.
