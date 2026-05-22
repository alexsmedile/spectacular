---
version: 1.1
updated: 2026-05-21
summary: "Architectural and product decisions for the Spectacular project"
---

# Decisions

## 2026-05-21

**Decision:** Split `PRD.md` into 4 focused root docs — `PRD.md` (intent), `PRINCIPLES.md` (beliefs), `ARCHITECTURE.md` (structure), `ROADMAP.md` (time). `AGENTS.md` rewritten as in-folder onboarding doc.
**Why:** The original `PRD.md` had grown to 896 lines holding four different documents (PRD + architecture spec + roadmap + skill interaction). This violated principle 3 ("small files over giant documents") and principle 1 (anyone planning had to load the whole thing). Splitting also unblocks future work on agent spec, capabilities, and runtime-enforced principles, none of which could land cleanly while PRD was the catch-all.
**Tradeoffs:** 4 root docs instead of 1 — more cross-links to maintain. Mitigated by the `related:` frontmatter convention and a tail "Related" section on each doc. Old `PRD@v1.3.md` and `AGENTS@v1.0.md` snapshots preserved as the authoritative pre-split content. The split was tracked end-to-end via `requests/canonical-docs-rework/`.

---

## 2026-05-11

**Decision:** Skill architecture uses lean SKILL.md + references/ subdocs
**Why:** Mirrors Spectacular's own philosophy — small files, layered context, routed by trigger. Avoids a monolithic SKILL.md that loads full context on every invocation.
**Tradeoffs:** Slightly more files to maintain; routing table in SKILL.md must stay current.

---

**Decision:** Output format is conversational briefing with minimal embedded table
**Why:** Feels like a collaborator, not a dashboard. Surfaces one priority action rather than dumping all state.
**Tradeoffs:** Less scannable at a glance than a pure table view.

---

**Decision:** Lifecycle state lives in PLAN.md frontmatter only
**Why:** Single source of truth. No state duplication across files.
**Tradeoffs:** Requires reading PLAN.md (not just directory listing) to know request state.

---

**Decision:** .spectacular/memory/ is team-visible, git-committed
**Why:** Memory should survive agent changes, tool changes, and team changes. It's operational learning, not personal context.
**Tradeoffs:** Requires deliberate write discipline (no auto-capture).

---

**Decision:** Versioning via snapshot-before-edit (PRD@v1.0.md naming)
**Why:** Canonical documents should never be silently mutated. Snapshots provide a full audit trail.
**Tradeoffs:** Directory accumulates versioned copies over time.

---

**Decision:** Workflows layer deferred to v2
**Why:** Each project handles release/hotfix/maintenance procedures differently. The design isn't finalized yet.
**Tradeoffs:** Projects can't document procedural sequences in v1.

---

**Decision:** Multi-agent / subagent orchestration deferred to v2
**Why:** v1 focuses on the core convention + skill + CLI. Multi-agent adds significant design complexity.
**Tradeoffs:** Parallel agent workflows not supported in v1.

---

## 2026-05-22 — PRD@v1.1.md snapshot gap

**Decision:** Accept the v1.1 gap in the PRD snapshot sequence (v1.0 → v1.2 → v1.3) rather than fabricate a backfill.

**Why:** Surfaced by `spectacular doctor` during M2 dogfood. Git history shows no commit ever referenced PRD@v1.1 — the version bump was simply skipped during the canonical-docs-rework work. Reconstructing content for a version that never existed would be fabrication.

**Tradeoffs:** Snapshot continuity has a documented gap. Doctor will continue to flag it as a known-acknowledged warning. Future doctor v2 could read DECISIONS.md to suppress flagged-and-documented gaps.
