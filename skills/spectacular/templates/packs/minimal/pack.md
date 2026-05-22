---
pack: minimal
version: 1.0
description: |
  The default pack shipped with Spectacular. Establishes only the essentials:
  a README contract so projects are triageable in 10 seconds, and a gitignore
  baseline so common temp/archive/sensitive paths never get committed.
  Opinions stop there — folder taxonomy, naming rules, and project-type
  scaffolds are all left empty so projects can grow organically without the
  skill forcing structure. Install a stronger pack (e.g. alex-default from
  the app store) when you want more opinions.
applies-to:
  - any
rules:
  naming: {}
  taxonomy: {}
  root-files:
    required:
      - README.md
      - .gitignore
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
      - _archive/
      - _archived/
      - _backup/
      - _backups/
      - _tmp/
      - scratch/
      - .env.local
      - .env.*.local
      - .spectacular.local/
    never-auto-add:
      - .scrapekit/
      - .playwright-mcp/
      - .smart-env/
      - .obsidian/
  file-placement: {}
  project-types: {}
templates:
  - .gitignore
  - README.md
references:
  - why-minimal.md
---

# Minimal pack

The bundled default. Two opinions only:

1. **A `.gitignore` exists** with safe defaults (archive/backup/temp/sensitive paths).
2. **A `README.md` exists** with a Type/Stack/Run header an agent or human can triage in 10 seconds.

Everything else is yours to shape. Pick a heavier pack from `<repo>/packs/` if you want opinions on folder layout, naming, or project-type scaffolds.
