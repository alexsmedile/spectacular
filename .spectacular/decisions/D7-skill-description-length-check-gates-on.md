# D7 — Skill description length check gates on `description` alone at 1024 chars (error), 1000 (warning) — not the description+w

**Context:**
Codex skipped loading spectacular at v1.17.1 because description was 1146 chars (over its 1024 cap; Claude Code's 1536 cap masked it). Need a guard so it can't regress silently.

**Decision:**
Skill description length check gates on `description` alone at 1024 chars (error), 1000 (warning) — not the description+when_to_use concatenation

**Consequences:**
Codex measures description ALONE: the v1.17.2 patch took description 1146→986 while description+when_to_use stayed 1253 (>1024) and the error cleared — proving the concatenation is not what Codex caps. check_skill() and the pre-commit guard both measure description alone.
