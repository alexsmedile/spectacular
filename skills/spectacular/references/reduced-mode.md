# Reduced Mode and Fallback Guidance

Use this when: Spectacular CLI is unavailable or incompatible, and the session must operate in read/draft fallback mode.

## Core Principle

Spectacular Skill and CLI are distributed separately:
- The **Skill** travels with the plugin / harness context.
- The **CLI** binary provides typed validation, deterministic SHA-256 fingerprints, UUIDv7 identity allocation, and atomic transitions.

When the CLI is not on `PATH` (or incompatible), the workspace is in **reduced mode**. This is a valid, supported operating tier—not a broken state.

```text
┌─────────────────────────────────────────────────────────────┐
│ REDUCED MODE RULES                                         │
│ • Read canonical Markdown records                          │
│ • Orient on project direction and mission status            │
│ • Explain method and draft Mission plans for activation     │
│ ─────────────────────────────────────────────────────────── │
│ ✕ Do not fabricate command-owned records or fingerprints   │
│ ✕ Do not claim an edit was atomic when it was a plain write │
│ ✕ Do not simulate CLI verification or certification        │
└─────────────────────────────────────────────────────────────┘
```

Without the CLI, you cannot run `mission start`, `objective promote`, `objective finish`, `mission complete`, `contract amend`, `review record`, or `handoff record` because those perform transactional writes and cryptographic binding.

---

## Installing the CLI

To restore full governed execution:

1. Download the platform release archive and `SHA256SUMS` from `https://github.com/alexsmedile/spectacular/releases/latest`.
2. Verify checksum: `shasum -a 256 --check SHA256SUMS --ignore-missing`.
3. Install from the directory holding the archive:

```bash
install/install.sh install \
  --prefix "$HOME/.local" \
  --source "$PWD" \
  --runtime claude \
  --version <VERSION>
```

4. Confirm installation:

```bash
spectacular --version
```

See [docs/installation.md](../../../docs/installation.md) for full installation and platform troubleshooting.

---

## Bundled Read-Only Fallbacks

The Skill bundles standalone zero-dependency scripts in `scripts/` that operate without the CLI:

### Shell Tier (Zero toolchain required)

```bash
sh scripts/doctor.sh          # inspect mode and report available capabilities
sh scripts/orient.sh          # workspace orientation, active Missions, and status
sh scripts/where.sh <ref>     # resolve human ref (e.g. M1, M1/O2) to record path
```

### Node.js Tier (Rich frontmatter & structural parsing)

Where Node.js is available:

```bash
node scripts/show.mjs <ref>    # state, outcome, objectives, dependency edges, gaps
node scripts/check.mjs [<ref>] # structural check; checks all records when ref omitted
```

### Fallback Guarantee & Limits

All bundled fallback scripts **read and report only**. None writes files, calculates fingerprints, or verifies cryptographic bindings. They provide a safe baseline floor for reading and drafting, not a mechanical equivalent to the binary.

---

## Declared Manual-Bootstrap Exception

If the project explicitly requires authoring or bootstrapping records by hand (e.g. self-development when the CLI cannot represent a new schema), follow [bootstrap.md](bootstrap.md). That exception requires explicit owner authorization and must never cite the legacy CLI as proof.
