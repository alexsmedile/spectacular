# D9 — Undo v1: single-level breadcrumb (.last-mutation), staleness via timestamp-vs-mtime, idea-promote undo prompts before re

**Context:**
lifecycle-undo (b12) grill resolved 3 open questions before build.

**Decision:**
Undo v1: single-level breadcrumb (.last-mutation), staleness via timestamp-vs-mtime, idea-promote undo prompts before removing the scaffolded request (default leave)

**Consequences:**
v1 scope stays tight: no breadcrumb stack, no git-reflog coupling, no destructive auto-delete. Each can graduate to v2 if real use demands. undo refuses on stale breadcrumb rather than mis-reverting.
