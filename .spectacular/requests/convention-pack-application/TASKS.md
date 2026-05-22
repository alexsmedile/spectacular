---
status: planned
updated: 2026-05-23
related:
  - PLAN.md
  - ../convention-pack-schema/PLAN.md
  - ../convention-pack-fabricator/PLAN.md
---

# Tasks — Convention Pack Application

## v1

### M1 — `pack install` / `pack list` / `pack remove`
- [ ] Add `pack` subcommand to `cli/spectacular` (sibling of `init` + `doctor`)
- [ ] `pack install <name>` — fetch tarball of `<repo>/packs/<name>/` from GitHub, extract to `~/.spectacular/packs/<name>/`
- [ ] Reuse skill-install download logic (GitHub release/tarball + sha verification)
- [ ] `pack list` — show installed packs grouped by scope (bundled / user / project-local)
- [ ] `pack remove <name>` — rm -rf user-scope pack; refuse to remove bundled or project-local unless --force
- [ ] Help text + error handling (unknown pack, network failure, etc.)

### M2 — config.yaml schema
- [ ] Document `convention_pack:` block in ARCHITECTURE.md
- [ ] Fields: `source` (path or name), `mode` (suggest|scaffold|enforce), `overrides` (optional list of rule-skips)
- [ ] CLI helper to read pack from config (resolves source: name → ~/.spectacular/packs/<name>/, path → literal path)
- [ ] `spectacular pack apply <name>` — set active pack in current repo's config.yaml (skill verb? CLI? decide)

### M3 — Init wiring
- [ ] Update `cli/spectacular` `resolve_doc_set()` (or sibling) to consult active pack
- [ ] If pack declared + mode=scaffold: scaffold pack's `applies-to` files/folders alongside always-set
- [ ] Always-set wins on conflicts (never overwrites; uses smart-init pre-flight)
- [ ] Update `references/init-workflow.md` to document pack integration
- [ ] Pre-flight: validate pack exists before init proceeds; fall back gracefully if not

### M4 — new-request wiring
- [ ] Update `references/new-request.md` — consult active pack for artifact directory layout
- [ ] If pack has `file-placement` rules for `research`, `screenshots`, `benchmarks` — create those subdirs on request scaffold
- [ ] If no pack or no rule for a kind, fall back to default (`artifacts/<kind>/`)

### M5 — Doctor `conventions` area
- [ ] Update `references/doctor.md` § Check areas — add `conventions` (active only when mode=enforce)
- [ ] Implement `check_conventions()` in CLI — walks pack's naming/taxonomy/root-files/gitignore/file-placement rules, flags violations
- [ ] Severity per violation type (naming = warning; missing required folder = error; gitignore drift = info)
- [ ] Update doctor scoped-areas list to include `conventions`

### M6 — Three modes proved
- [ ] suggest: skill briefing mentions pack opinions when relevant ("This project would normally have scripts/ — want me to create it?")
- [ ] scaffold: init/new-request actively apply pack
- [ ] enforce: doctor flags drift; recommends `pack apply --strict` or `--relaxed` to toggle
- [ ] Document mode semantics in `references/packs-contract.md` and `init-workflow.md`

### M7 — Tests + VERIFY.md
- [ ] Extend `tests/cli/pack.test.sh` with install/list/remove scenarios (mocked GitHub fetch or local source path)
- [ ] Scenario: install bundled minimal → list shows it → init consults it → remove it
- [ ] Scenario: init with config.yaml `convention_pack:` → expected folders scaffolded per pack's applies-to
- [ ] Scenario: enforce mode + workspace drift → doctor `conventions` reports warnings/errors
- [ ] Scenario: project-local pack overrides user-scope same-name pack
- [ ] Create `requests/convention-pack-application/VERIFY.md` (4-of-6 score — comprehensive)
- [ ] Manual scenarios: full interactive flow on a fresh repo + alex-default applied

### Release
- [ ] CHANGELOG.md v0.4.0 entry — pack system as feature release
- [ ] Update README.md — packs section + install/list/remove + config.yaml example
- [ ] Update docs/commands.md — pack subcommands documented
- [ ] Update docs/configuration.md — `convention_pack:` schema + mode semantics

## v2 (deferred)

- [ ] Pack composition (multi-pack stack with precedence)
- [ ] Pack auto-update from GitHub
- [ ] `spectacular pack publish` — upload via PR
- [ ] Pack signing/verification
- [ ] `spectacular pack migrate` — retroactively apply pack to existing repo (doctor on steroids)
- [ ] Per-pack telemetry (anonymous usage stats — opt-in only)
