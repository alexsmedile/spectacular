---
status: active
updated: 2026-05-30
related:
  - PLAN.md
---

# Tasks — verify-walk

## Design decisions (resolved 2026-05-30)

- **Checks are TYPED — verification is multi-authority, not one thing.** Five kinds along a spine (deterministic → judgment → human), each with its own authority + walk behavior:
  - `executable` (`` `run: <cmd>` ``) → external command exit code (deterministic).
  - `assertable` (`{assert}`) → agent checks a binary property of files/state (deterministic, no subprocess, no opinion).
  - `judgable` (`{judge}`) → LLM reasons over named artifacts; fuzzy; uncertainty = blocker.
  - `observable` (`{observable}`, untagged default) → human looks/confirms (passive).
  - `manual` (`{manual}`) → human performs an action first, then confirms (active).
  - No overlaps: exec=external-tool vs assert=agent-check; judge=fuzzy vs assert=binary; observe=passive vs manual=active.
- **Syntax — two accepted shapes:** (a) **inline** per-line tag + optional `run:`; (b) **section-grouped** — `## Title {kind}` applies to all checks under it. **Section is absolute** (inline tags inside a tagged section are ignored). `## Title {run}` = executable section where each line IS the command; inline executable still uses per-line `run:`. Untagged line/section = observable. Tags stay literal in the file.
- **Exec safety:** confirm-each `run:` command (show it, y/n/skip); batch-allow option at walk start. Never run an unshown command.
- **All-pass outcome:** configurable, **default propose** (human confirms via `promote`); opt-in `verify.auto_promote` / `--auto` — the seam where [[policy-engine]] severity plugs in later.
- **Blocker handling:** record + **keep walking**, report all at end, stay `review`. No fabricated passes (any kind).
- **Result recording (both):** tick checkboxes inline in VERIFY.md (live state) + append a timestamped entry to **`VERIFY-LOG.md`** (audit trail, append-only, records the `[kind]` per check). VERIFY-LOG is a new per-request artifact this request introduces.

> **Scope note:** typed checks expand M1 beyond the original "ask human, LLM judges" sketch. The VERIFY.md format now carries kind tags + `run:` commands — `scaffold-reference.md`'s VERIFY.md stub + the `verification.md` doc both need updating to document the typed syntax (added to M4).

## v1

### M1 — Walk algorithm
- [x] Write `references/verify.md`: locate VERIFY.md or fall back to PLAN § Validation
- [x] Iterate each check; prompt for evidence; record pass / blocker per item
- [x] Define the gate (all-pass vs any-blocker outcomes)

### M2 — Lifecycle wiring (specified in verify.md; needs CLI/surface work to be live)
- [x] All-pass → flip `review → verified` via `promote` (algorithm defined; default-propose + configurable auto)
- [x] Any blocker → stay `review`, write the blocker list to VERIFY-LOG
- [ ] Confirm `config.yaml` `verify.auto_promote` key is read where the gate checks it

### M3 — Retrospective + archive tie-in
- [x] End-of-walk optional "what surprised you?" prompt → `memory/` entry (defined in verify.md § 5)
- [ ] `spectacular archive` warns when `verified` was reached with no VERIFY-LOG (CLI/archive.md change)

### M4 — Surface + docs
- [ ] SKILL.md routing-table entry for `verify`
- [ ] CLI redirect: `spectacular verify <slug>` → skill-only message
- [ ] `docs/commands.md` agentic-verbs section covers `verify`
- [ ] VERIFY-LOG.md stub added to `scaffold-reference.md`

### M5 — Dogfood + ship
- [ ] Drive 1+ real request through the walk to `verified`
- [ ] CHANGELOG [1.11.0] entry; plugin bump to v1.11.0

## v2 (deferred)

- [ ] Auto-suggest the walk when a request hits `review` (proactive surface)
- [ ] Per-check evidence persistence (store evidence inline in VERIFY.md)
