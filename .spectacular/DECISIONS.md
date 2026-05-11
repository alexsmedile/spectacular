---
version: 1.0
updated: 2026-05-11
summary: "Architectural and product decisions for the Spectacular project"
---

# Decisions

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
