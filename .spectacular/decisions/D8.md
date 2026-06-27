# D8 — Stub rules-file bodies thin to frontmatter + one pointer; the shared 3-verb default lives once in doc-index.md

**Context:**
6 of 18 <doc>-rules.md were mode:stub with near-identical ~21-line bodies (grill no-op / refine rewrite / review structural). The boilerplate drifts per-file and adds registry surface for zero per-doc info.

**Decision:**
Stub rules-file bodies thin to frontmatter + one pointer; the shared 3-verb default lives once in doc-index.md

**Consequences:**
architecture/principles/stack thin to pointer; agents keeps its top-level-AGENTS.md delta; spec and tasks keep real bodies (spec=index+sync role, tasks=full review/refine spec — mislabeled mode:stub). doc-index.md gains a 'stub default behavior' section. Frontmatter untouched everywhere — engine dispatch intact.
