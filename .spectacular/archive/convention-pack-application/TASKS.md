---
status: review
updated: 2026-05-23
related:
  - PLAN.md
  - VERIFY.md
  - ../convention-pack-schema/PLAN.md
  - ../convention-pack-fabricator/PLAN.md
---

# Tasks — Convention Pack Application

## v1

### M1 — `pack install` / `pack list` / `pack remove`
- [x] Add `pack` subcommand to `cli/spectacular` (sibling of `init` + `doctor`)
- [x] `pack install <name>` — copy pack folder from resolved source (bundled / app-store / --from path) to `~/.spectacular/packs/<name>/`
- [x] Reused install logic (local cp -R for v1; GitHub tarball fetch deferred to v2 once app-store is published)
- [x] `pack list` — shows installed packs grouped by scope (bundled / app-store / user / project-local)
- [x] `pack remove <name>` — rm -rf user-scope pack; refuses bundled / app-store / project-local without `--force`
- [x] `pack show <name>` — bonus verb showing scope + path + frontmatter
- [x] Help text + error handling (unknown pack, already-installed, missing source)

### M2 — config.yaml schema
- [x] Document `convention_pack:` block (commented-out template in init's config.yaml output)
- [x] Fields documented: `source` (pack-id), `mode` (suggest|scaffold|enforce), `overrides` (forward-declared, unused in v1)
- [x] CLI helpers: `config_pack_source()` and `config_pack_mode()` awk parsers
- [x] No separate `spectacular pack apply` verb — activation = editing config.yaml (documented in init-workflow.md)

### M3 — Init wiring
- [x] `scaffold()` consults active pack after doc-set scaffolding completes
- [x] If pack declared + mode=scaffold|enforce: append pack's `gitignore.always-add` entries (deduplicated)
- [x] Always-set wins on conflicts — pack scaffold never overwrites existing lines
- [x] Updated `references/init-workflow.md` — § "Convention packs (v0.4.0+)" with 3-mode matrix + precedence table
- [x] Graceful fallback: missing pack source logged as info, init continues

### M4 — new-request wiring
- [x] Updated `references/new-request.md` — `artifacts/` directory consults active pack's `file-placement.request-artifacts:` rule
- [x] Fallback documented: default to `artifacts/<kind>/` when no pack or no rule
- [ ] **Deferred to v2:** actively consuming `file-placement` rules in the CLI new-request flow (CLI doesn't have a `new-request` subcommand yet — skill-driven, so rule consumption lives in skill)

### M5 — Doctor `conventions` area
- [x] Updated `references/doctor.md` § Check areas — new `conventions` area
- [x] Implemented `check_conventions()` in CLI — gitignore drift detection per pack rules
- [x] Severity: warning in scaffold mode, error in enforce mode, info in suggest mode (no drift checks)
- [x] Updated doctor scoped-areas allowlist to include `conventions`
- [x] Updated `DOC_AREAS` constant + `doctor_usage` help text

### M6 — Three modes proved
- [x] **suggest** — `doctor conventions` reports active pack + info-note "drift checks are disabled" (no warnings/errors)
- [x] **scaffold** — init appends pack gitignore entries; `doctor conventions` flags drift as warnings (exit 1)
- [x] **enforce** — same scaffold behavior; doctor flags drift as errors (exit 2); `--fix` mechanically repairs
- [x] Three-mode behavior documented in `references/init-workflow.md` table + `references/doctor.md`

### M7 — Tests + VERIFY.md
- [x] Created `tests/cli/pack.test.sh` — 12 scenarios, 44 asserts
- [x] Scenarios cover: list, install (bundled + --from), remove (user + bundled-refusal), show, init wiring, doctor (no-pack / enforce / --fix), help, error paths
- [x] All tests pass: `tests/run.sh` reports 83/83 (39 init + 44 pack)
- [x] Created `requests/convention-pack-application/VERIFY.md` — 4-of-6 score, comprehensive (18 scenarios across automated + manual)
- [ ] **Live (VERIFY S13-S18):** three-mode end-to-end walkthrough, scope-precedence resolution, cross-machine portability, interactive init flow, alex-default e2e, config drift edge cases

### Release
- [ ] CHANGELOG.md v0.4.0 entry — pack system as feature release
- [ ] Update README.md — packs section + install/list/remove + config.yaml example
- [ ] Update docs/commands.md — pack subcommands documented
- [ ] Update docs/configuration.md — `convention_pack:` schema + mode semantics

## v2 (deferred — handled by other requests)

- [ ] Pack composition / modular packs → `convention-pack-modules`
- [ ] Pack auto-update from GitHub
- [ ] `spectacular pack publish` — upload via PR
- [ ] Pack signing/verification
- [ ] `spectacular pack migrate` — retroactively apply pack to existing repo (doctor on steroids)
- [ ] Per-pack telemetry (anonymous usage stats — opt-in only)
- [ ] GitHub tarball install (current `install` uses local cp; `--from` accepts arbitrary path; remote URL deferred)
- [ ] Interactive pack selection during `spectacular init -i` (v1 requires post-init config.yaml edit)
- [ ] CLI `new-request` subcommand consuming `file-placement` rules directly (currently skill-side only)
