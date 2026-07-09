### ADDED
- specs/index.md :: workspace-layout — the `.spectacular/` workspace follows the Open Knowledge Format (OKF) v2.0 layout: all collection directories are plural (`specs/`, `decisions/`, `memories/`, `sessions/`, `roadmaps/`, `feedbacks/`, `audits/`, `fixes/`, `debugs/`, `requests/`, `ideas/`), each collection's index lives inside its folder as `index.md` (not a root `DECISIONS.md`/`MEMORY.md`/`ROADMAP.md`), capability specs are flat files (`specs/<cap>.md`, no `specs/<cap>/SPEC.md` nesting), and decision/memory entries carry sequential prefixes (`D<N>-<slug>.md`, `M<N>-<slug>.md`).
- specs/index.md :: migration — `spectacular migrate` (and `doctor --fix`) upgrades a schema-0.6 workspace to 2.0 mechanically via `migration_apply_v06_to_v20`: renames singular dirs → plural, relocates root index files → `<dir>/index.md`, flattens nested capability specs, sequentially prefixes decision/memory entries with collision-safe de-dup, and rewrites `[[…]]`/`(…)` link targets + YAML `related:`/`depends-on:`/`blocks:` fields depth-correctly (link-only, idempotent, never touches bare prose).

### MODIFIED
- ARCHITECTURE.md :: workspace tree + ID-namespace convention now document the OKF v2.0 plural-dir + `index.md`-in-folder + `D<N>`/`M<N>`-prefix layout as the canonical shape (was: mixed singular/plural dirs, root index files, nested `specs/<cap>/SPEC.md`).

### REMOVED
- (none — no capability was removed; the change is a layout/convention migration plus a new migration path.)

### NOTES
- M1's literal term "project-wide anchors" was not coined as a defined heading; the concept ships (PRD.md is described as an "anchor doc", the ID-namespace table is marked "project-wide"), but the exact phrase is deferred as cosmetic. Recorded in TASKS.md as `[~]`.
- v2 items (compose subcommands, external convention-pack auto-migration) are out of scope and tracked under the `convention-pack-modules` request.
