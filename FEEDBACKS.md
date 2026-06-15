# FEEDBACKS

Field reports to review and fix. Each entry is a real-world issue hit while using spectacular. Triage, then either cut a request or fix in place.

---

## 2026-06-15 — DECISIONS.md bloats context at scale (Octopus project) — 📌 ROUTED → `decisions-index` (b9, v1.17.0)

Source: Octopus project hit 2028-line DECISIONS.md (D1–D110) after ~6 months of active design. Reading the file in full costs significant context; agents loading it before every planning session pay the price every time.

**Request: implement `index` mode for DECISIONS.md** (the v1.6.x soft-folder shape mentioned in `decisions-rules.md`).

**Desired layout:**

```
.spectacular/
├── DECISIONS.md          ← index only: one line per decision (number, title, one-sentence summary)
└── decisions/
    ├── D1.md             ← full ADR prose for D1
    ├── D2.md
    └── ...
```

**Agent read pattern:**
- Always load `DECISIONS.md` (index) — cheap, ~1 line per decision
- Load `decisions/Dxxx.md` on demand when the decision is directly relevant to current work

**CLI verbs needed:**
- `spectacular decide "<text>"` → appends one-liner to DECISIONS.md index + writes `decisions/D<N>.md` with full ADR prose
- `spectacular decision <slug>` → loads and returns `decisions/<slug>.md`
- `spectacular decisions migrate` → one-shot: reads flat `DECISIONS.md`, splits into `decisions/D*.md` files, rewrites `DECISIONS.md` as index

**Index line format (suggestion):**
```
- **D42** — Reject field-mode storage for v1 — storage is folders-only until v2 ships
```

**Backward compat:** flat DECISIONS.md should still be valid (detected by absence of `decisions/` folder). The `mode: append` frontmatter in `decisions-rules.md` could become `mode: index` when migrated.

**Priority:** medium — not blocking, but the longer a project runs the worse the context tax gets. Octopus will manually split its DECISIONS.md now and adopt the convention ahead of the CLI support.

---

## 2026-05-29 — v1.8.3 dogfooding (external project use)

Source: agent using spectacular to manage a separate iOS/macOS build project.

### 🔴 Bug 1 — `spectacular remember` is broken: "memory entry template not found" — ✅ SHIPPED v1.8.4

- **Repro:** `spectacular remember "<text>" --tag a,b` fails *every time* with `Error: remember: memory entry template not found`.
- **Confirmed reproducible** on current HEAD (v1.8.3) using the installed binary against a fresh `--no-skill` workspace. Exit 1.
- **NOT a 1.8.3 regression** — the broken resolution path has existed since v1.0.0; the symlink-unaware `SCRIPT_DIR` since the binary's inception. Any prior version hit the identical bug.

**Root cause** — `_resolve_template()` (cli/spectacular ~L888) has 3 fallback paths, all fail for the canonical symlinked global install:

| Path | Resolves to | Why it fails |
|---|---|---|
| 1. project override | `.spectacular/templates/memory/entry.md` | absent in normal workspaces |
| 2. `agents_skill_dir` | `$(pwd)/.agents/skills/spectacular/...` (or `~/.agents/` only if scope==global) | empty workspace has no local `.agents/` skill; and scope is never `global` at `remember` time → the `~/.agents` branch never fires |
| 3. `${SCRIPT_DIR}/../skills/...` | `~/.local/bin/../skills/...` | **smoking gun:** `SCRIPT_DIR` (L6) uses `dirname "${BASH_SOURCE[0]}"` which does NOT resolve the symlink. `~/.local/bin/spectacular` → `~/.local/bin`, so Path 3 looks in `~/.local/skills/` which never exists |

**Fix (two parts):**
1. Resolve the symlink when computing `SCRIPT_DIR` (L6) so Path 3 points at the real `cli/` dir:
   ```bash
   SOURCE="${BASH_SOURCE[0]}"
   while [[ -L "$SOURCE" ]]; do
     DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
     SOURCE="$(readlink "$SOURCE")"; [[ "$SOURCE" != /* ]] && SOURCE="$DIR/$SOURCE"
   done
   SCRIPT_DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
   ```
2. Path 2 should also add a scope-independent `~/.agents/skills/spectacular/templates/` fallback so an installed-skill workspace resolves regardless of `OPT_SKILL_SCOPE`.

**Applied** (both parts) in `cli/spectacular`, verified: `remember` via symlinked global binary writes entry + regenerates `MEMORY.md` (exit 0); `--dry-run` previews + writes nothing; direct non-symlink invocation still resolves (no regression). Shipped in v1.8.4 (commit `a968186`).

- **Impact:** had to write the memory file by hand and manually update `MEMORY.md` — defeats the "deterministic mutator" principle. Affects *every* template-backed verb, not just `remember`.
- **Severity:** HIGH — core mutator verb non-functional.

### 🟡 Bug 2 — `spectacular decide` writes empty sections — ✅ SHIPPED v1.8.4

- The verb captured the full decision text into `**Decision:**` but left the other sections blank.
- **Root cause was template drift, now reconciled:** CLI inline entry + `templates/decisions/entry.md` + embedded `doc_decisions` scaffold all used `Decision / Why / Tradeoffs`, while `decisions-rules.md` documented `Context / Decision / Consequences`. Canonical schema chosen = **Context / Decision / Consequences** (Nygard ADR). All three emit sites aligned to it.
- **Fix:** added `--context "..."` and `--consequences "..."` flags to `decide` (positional arg fills `**Decision:**`). Omitted sections emit as empty headers (fillable later by hand or skill `grill`/`refine`) — the verb never invents them. `decisions-rules.md` updated with the flag usage.
- **Verified:** flags populate sections; no-flag path emits empty headers (not dropped); `--dry-run` writes nothing; embedded scaffold uses new schema; `--help` updated. Test suite: 3 pre-existing environmental failures unchanged (identical counts with/without the change → zero regressions). Shipped in v1.8.4 (commit `eb05401`).
- **Severity:** MEDIUM — verb worked but output was thin / misleading.

### 🟡 Bug 3 — broken mutator never surfaces its manual-recovery path — ✅ SHIPPED v1.8.4

- When `remember` errored, nothing told the agent the fallback (write to `.spectacular/memory/<slug>.md` + update index by hand).
- **Fix:** added a shared `die_recover` helper (error message + a "→ Manual recovery:" hint block, all to stderr). Wired into the two template-resolution death sites: `remember` (memory) and `session start` (session). The `remember` hint includes the actual derived slug + the exact frontmatter shape to write.
- **Verified:** forced a real template-resolution failure through the installed binary — `remember` now prints the manual-recovery path with the derived slug instead of a dead-end error; success path unchanged; zero test regressions. Shipped in v1.8.4 (commits `a76d0f0` + `92d01e7`).
- **Note:** now that Bug 1 is fixed, the template-not-found death can't occur in the common case — this is hardening for the residual failure modes (corrupted install, hidden template). Future template-backed verbs should use `die_recover` too.
- **Severity:** MEDIUM → LOW (de-risked by Bug 1 fix) — resilience / DX.

### 🟢 Minor — stale memory not auto-flagged — 📌 ROADMAPPED v1.12.0

- A memory (`m1-build-is-blocked-on-full`, an Xcode/macOS-26 blocker) became obsolete after the project shipped via Command Line Tools (no Xcode needed). Nothing in spectacular noticed the memory went stale when the underlying build path changed.
- Not a bug, but the "be proactive, surface stale state" promise didn't trigger. Worth considering a staleness heuristic (e.g. memory referencing a blocker that a later session/decision contradicts → flag in `doctor memory`).
- Left in place rather than deleting unprompted.
- **Severity:** LOW — enhancement.
- **Disposition (2026-05-29):** Split into two. (a) *Naive age-check* `doctor memory` staleness flag — mirrors the existing `sessions`/`feedback`/`ideas` convention; routed as a side-rider on M3 of the [[cross-request-links]] request (v1.12.0). (b) *Contradiction-check* (memory overturned by a later session/decision) — cross-doc semantic reasoning, likely skill-side; deferred to v2 in that request's TASKS.

**Net (updated 2026-05-29):** All three bugs shipped in **v1.8.4** (template-resolution + `decide` schema + `die_recover`); v1.9.0 followed with the versioning convention. The lone remaining item — the 🟢 stale-memory minor — is now roadmapped to v1.12.0. **No open fixes.**
