---
pack: alex-default
version: 1.0
description: |
  Opinionated defaults for solo-dev and small-team mono-collections.
  Encodes the maintainer's working conventions: kebab-case naming with
  role suffixes, mono-collection detection for ~/code-style roots,
  AGENTS.md pattern, README contract, full gitignore baseline, and
  scaffold templates for cli / library / skill / plugin / content /
  research / vault-project. Install when you want strong, consistent
  structure across many projects; stick with minimal if you prefer
  per-project freeform.
applies-to:
  - any
rules:
  naming:
    folder-case: kebab-case
    file-case: kebab-case
    pattern: "{anchor}-{descriptor}[-{role}][-{qualifier}]"
    max-words: 3
    role-suffixes:
      - ctrl
      - manager
      - svc
      - orch
      - worker
      - dash
      - viz
      - exp
    forbidden-words:
      - v2
      - new
      - project
      - temp
      - misc
      - util
      - utils
      - common
      - shared
    forbidden-prefixes:
      - app-
      - svc-
    date-formats:
      sandbox: "-YYYYMM"
      archive: "-YYYY-MM"
    language-exceptions:
      python: snake_case   # python package dirs are imported, must be snake
      go: lowercase        # go convention
  taxonomy:
    required: []
    opt-in:
      - src/
      - scripts/
      - tests/
      - test/
      - docs/
      - examples/
      - assets/
      - bin/
      - cli/
      - _research/
      - _archive/
      - _backups/
      - _tmp/
      - scratch/
    mono-collection-detect:
      when-children-have:
        - .git/
        - package.json
        - pyproject.toml
        - SKILL.md
      threshold: 2
    mono-collection-folders:
      - apps/
      - libs/
      - tools/
      - design/
      - dash/
      - sandbox/
      - templates/
      - archive/
      - infra/
  root-files:
    required:
      - README.md
      - .gitignore
    conditional:
      - file: AGENTS.md
        when: agentic
      - file: LICENSE
        when: oss
      - file: CHANGELOG.md
        when: post-v1
    optional:
      - CLAUDE.md
      - STACK.md
      - PRD.md
      - PLAN.md
      - TASKS.md
    readme-contract:
      must-contain-header:
        - Type
        - Stack
        - Run
      must-contain-sections:
        - What it does
        - Setup
        - Usage
  gitignore:
    always-add:
      - .spectacular.local/
      - _archive/
      - _archived/
      - _backup/
      - _backups/
      - _tmp/
      - scratch/
      - .env.local
      - .env.*.local
    opt-in:
      - .cache/
    never-auto-add:
      - .scrapekit/
      - .playwright-mcp/
      - .smart-env/
      - .obsidian/
    language-specific:
      python:
        - __pycache__/
        - "*.pyc"
        - .venv/
        - .pytest_cache/
        - .ruff_cache/
        - dist/
        - build/
        - "*.egg-info/"
      node:
        - node_modules/
        - dist/
        - build/
        - .next/
        - .turbo/
        - coverage/
      go:
        - vendor/
        - "*.exe"
  file-placement:
    helper-script: scripts/<name>.sh
    architecture-doc: docs/<name>.md
    skill-reference: references/<name>.md
    research-artifact: _research/<topic>/
    backup: _backups/<timestamp>/
    generated-cache: .cache/<name>
    sensitive-data: .env.local
    temp-work: scratch/<name>
    request-artifacts: .spectacular/requests/<slug>/artifacts/{kind}/
    large-file-threshold: 5MB
  project-types:
    cli:
      adds:
        - cli/
        - scripts/
        - install.sh
        - README.md
        - LICENSE
      template-dir: templates/repo/cli/
    library:
      adds:
        - src/
        - tests/
        - examples/
        - docs/
        - README.md
        - LICENSE
      template-dir: templates/repo/library/
    webapp:
      adds:
        - src/
        - public/
        - tests/
        - .env.example
        - README.md
      template-dir: templates/repo/webapp/
    skill:
      adds:
        - SKILL.md
        - references/
        - templates/
        - scripts/
        - README.md
      template-dir: templates/repo/skill/
    plugin:
      adds:
        - .claude-plugin/
        - skills/
        - agents/
        - commands/
        - README.md
      template-dir: templates/repo/plugin/
    content:
      adds:
        - articles/
        - _research/
        - assets/
        - drafts/
        - README.md
      template-dir: templates/repo/content/
    research:
      adds:
        - _research/
        - notebooks/
        - data/
        - reports/
        - README.md
      template-dir: templates/repo/research/
    vault-project:
      adds:
        - assets/
        - inbox/
        - brand/
        - business/
        - marketing/
        - offer/
        - projects/
        - tasks/
        - README.md
      template-dir: templates/repo/vault-project/
templates:
  - .gitignore
  - README.md
references:
  - why-alex-default.md
---

# alex-default pack

The maintainer's opinionated defaults, distilled from years of `~/code/` and `~/vault/` practice. Install when you want a single consistent shape across many projects.

See `references/why-alex-default.md` for rationale on each rule, including the AGENTS.md "most-specific wins" pattern (which is runtime behavior, not a schema rule).
