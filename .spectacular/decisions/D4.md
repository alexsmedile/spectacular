# D4 — PRD@v1.1.md snapshot gap

**Decision:** Accept the v1.1 gap in the PRD snapshot sequence (v1.0 → v1.2 → v1.3) rather than fabricate a backfill.

**Why:** Surfaced by `spectacular doctor` during M2 dogfood. Git history shows no commit ever referenced PRD@v1.1 — the version bump was simply skipped during the canonical-docs-rework work. Reconstructing content for a version that never existed would be fabrication.

**Tradeoffs:** Snapshot continuity has a documented gap. Doctor will continue to flag it as a known-acknowledged warning. Future doctor v2 could read DECISIONS.md to suppress flagged-and-documented gaps.
