# D1 — Request slugs stay unnumbered

**Decision:** Requests live at `requests/<slug>/`, not `requests/NNN-<slug>/`. ID = slug.

**Why:** Numbering was raised (TODO.md) to fix folder ordering and give each request a stable ID. Costs outweigh the win:
- Slug *is* already a stable ID — every PLAN.md, cross-link, archive path, and CLAUDE.md table references it. Adding a number prefix means every existing reference becomes a rename or breaks.
- Folder ordering is solved cheaply elsewhere: `CLAUDE.md` "Active Requests" table is the human entry point and is explicitly ordered; `spectacular status` orders by lifecycle state, not directory listing.
- Numbering invites bikeshedding (zero-pad width? gaps when archived? renumber on insert?) — each answer has its own cost.
- Lifecycle already provides a temporal axis via `PLAN.md` frontmatter (`updated`, `status`) — directory order isn't load-bearing.

**Tradeoffs:** `ls requests/` shows alphabetical, not chronological. Acceptable — directory listing is not the intended discovery surface; `spectacular status` and CLAUDE.md are. Revisit if a project hits 30+ active requests and discovery genuinely degrades.
