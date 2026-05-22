---
status: pending
updated: 2026-05-23
related:
  - PLAN.md
  - TASKS.md
score: 4/6
---

# Verification — Convention Pack Application

Per [[verification]] 2-of-6 rule: this request scored 4/6 (user-visible CLI + multi-surface flow + external contract: config.yaml schema + CLI subcommand). Comprehensive scenarios below — many are covered automatically by `tests/cli/pack.test.sh`; live items require interactive shell sessions.

## Automated (covered by tests/cli/pack.test.sh)

- [x] S1: `pack list` shows bundled + app-store + (when present) user + project-local scopes
- [x] S2: `pack install <name>` copies bundled pack → `~/.spectacular/packs/<name>/` (full tree: pack.md + templates/ + references/)
- [x] S3: re-installing same pack fails with "already installed"
- [x] S4: `pack show <name>` prints scope + path + frontmatter
- [x] S5: `pack remove <name>` deletes user-scope pack
- [x] S6: `pack remove <bundled-pack>` refuses without `--force`
- [x] S7: init with `convention_pack.mode: scaffold` appends pack's gitignore.always-add entries (idempotent)
- [x] S8: `doctor conventions` reports skip-info when no pack declared
- [x] S9: `doctor conventions` in enforce mode flags missing gitignore as errors (exit 2)
- [x] S10: `doctor conventions --fix` mechanically repairs pack-driven gitignore drift
- [x] S11: `pack install <name> --from <local-path>` installs from arbitrary folder
- [x] S12: `pack --help` shows usage with all verbs + scope precedence table

**Automated assert count: 44 passed, 0 failed** (run via `tests/run.sh`).

## Manual scenarios

### S13 — Three modes proved (live)

- [ ] **suggest mode** — declare pack with `mode: suggest`; run `doctor conventions` — confirms pack is active + info-level "drift checks disabled" message; no warnings, no errors
- [ ] **scaffold mode** — declare pack with `mode: scaffold`; re-run init; confirm pack gitignore entries appended; `doctor conventions` reports warnings (not errors) for any drift introduced manually
- [ ] **enforce mode** — declare pack with `mode: enforce`; introduce drift (delete a gitignore entry); `doctor conventions` exits 2 with error severity per missing entry

### S14 — Scope precedence resolution (live)

- [ ] With `minimal` in bundled AND user scope, `pack show minimal` reports `scope: user` (user wins over bundled)
- [ ] With same pack in user AND project-local, `pack show <name>` reports `scope: project-local` (most-specific wins)
- [ ] `pack install <name>` from a name that exists in app-store and bundled: source resolves to project-local first (per precedence), then user, then app-store, then bundled

### S15 — Cross-machine portability (live)

- [ ] Pack installed on machine A (`~/.spectacular/packs/<name>/`) can be tar-copied to machine B's `~/.spectacular/packs/<name>/` and `pack show` / `init` work identically
- [ ] Pack referenced in `config.yaml` (`convention_pack.source: <name>`) but not installed on the target machine produces a clear error from doctor (instructing user to install) — verified via S9 result

### S16 — Interactive init flow (live)

- [ ] Run `spectacular init -i` on a fresh dir
- [ ] Confirm: no pack-related prompts in the interactive flow (v1 limitation — pack selection is post-init via config.yaml; documented in init-workflow.md)
- [ ] Document any UX friction discovered for a v0.5.0 interactive-pack-selection follow-up

### S17 — alex-default end-to-end (live)

- [ ] `spectacular pack install alex-default` from app-store
- [ ] Add `convention_pack: { source: alex-default, mode: enforce }` to config.yaml
- [ ] Re-run init; confirm all 8 alex-default gitignore entries appended
- [ ] Introduce drift (delete `_archive/` line); `doctor conventions` exits 2
- [ ] `doctor conventions --fix` repairs and exits 0

### S18 — Config drift edge cases (live)

- [ ] `convention_pack.source:` typo (e.g. `alex-defualt`) — doctor reports "not found in any scope" with judgment fix suggestion
- [ ] `convention_pack.mode:` unrecognized value (e.g. `strict`) — falls back to `suggest` (default); doctor adds info note about unknown mode
- [ ] config.yaml with `convention_pack:` block but no `source:` field — parser returns empty source; treated as "no pack declared"

## Sign-off

- [ ] All scenarios complete or explicitly waived with reason
- [ ] No `❌` items remain; `⚠️` items have a note
- [ ] CHANGELOG.md updated with v0.4.0 entry
- [ ] README.md packs section + config.yaml example present
- [ ] Verifier: alex
- [ ] Date verified: ____________
