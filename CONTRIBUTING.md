# Contributing to Spectacular

Thanks for your interest. Spectacular is a small, opinionated project — contributions are welcome but the bar is *fit*, not feature count. Read this first.

## Before you open a PR

1. **Open an issue first** for anything non-trivial. A PR landing cold is much harder to land than one that started as an "is this the right shape?" conversation.
2. **Read `.spectacular/PRINCIPLES.md`** — Spectacular has explicit operating principles, and PRs that fight them won't merge.
3. **Read the relevant capability bullet in `.spectacular/SPEC.md`** — this is what's built right now. If your PR contradicts a SPEC bullet, surface that in the PR description.
4. **Use Spectacular's own workflow** to track the work — open a request in `.spectacular/requests/<slug>/` with PLAN.md + TASKS.md. Yes, we dogfood.

## Local setup

```bash
git clone https://github.com/alexsmedile/spectacular
cd spectacular
bash tests/run.sh        # full test suite (CLI + skill behaviors)
```

Tests live in `tests/cli/` — one file per major area (init, doctor, mutator verbs, migrations, packs, specs, docs, archive). New behavior needs a test scenario.

## What's in scope

- **Bug fixes** with a failing test scenario reproducing the bug
- **New convention packs** — add to `packs/<name>/` following the `alex-default` shape
- **New doc-types** for the doc-writing engine — see `skills/spectacular/references/doc-registry.md`
- **CLI verb additions** that fit the **mutation principle** (CLI mutates state, skill orchestrates judgment)
- **Reference doc clarifications** in `skills/spectacular/references/`
- **Migration entries** for workspace-schema upgrades — see `skills/spectacular/references/migrations/`

## What's out of scope

- **Public-facing documentation work** — that lives in the sibling [pageworks](https://github.com/alexsmedile/pageworks) skill. Spectacular hands off `docs/` to pageworks and won't accept renderer/site-generator features.
- **Project management features** — issues, sprints, kanban, time tracking. Spectacular tracks AI-agent work, not team capacity.
- **GUI / web UI** — spectacular is a directory convention + CLI + slash command. Not changing.
- **Backwards-incompatible changes to the workspace schema** without a migration entry under `skills/spectacular/references/migrations/`.

## Commit + PR conventions

- Conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`
- One logical change per PR — don't bundle a feature with a refactor
- Tests must pass: `bash tests/run.sh`
- If your PR touches a canonical `.spectacular/` doc, **snapshot first**: `spectacular snapshot .spectacular/<FILE>.md`
- Bump `CHANGELOG.md` under `## [Unreleased]` — maintainers handle the version bump at release time

## Releases

Releases are tagged by the maintainer. The flow:

1. `/release <version>` (uses git-guard's bump-manifests for plugin.json, marketplace.json, README badge)
2. Tag pushed → GitHub Release created
3. Plugin marketplace updated separately via `/plugin marketplace update spectacular`

## Code of conduct

Be kind. Be specific. Assume good faith.

## Questions

Open a [Discussion](https://github.com/alexsmedile/spectacular/discussions) for design questions, an [Issue](https://github.com/alexsmedile/spectacular/issues) for bugs or concrete proposals.
