# Verify — wayfinding-contract

## Contract
- [x] {assert} Canonical aliases resolve to normalized IDs and ambiguous naked numbers require context.
- [x] {assert} Canonical references, padding, reserved prefixes, and compatibility behavior are documented.
- [x] Typed record, spec lifecycle, fog/frontier, migration, and doctor fixtures pass. `run: bash tests/cli/wayfinding-contract.test.sh`
- [x] Bash syntax is valid. `run: bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit`
- [x] The guarded version strings are consistent. `run: scripts/hooks/pre-commit --check`
- [x] The complete regression suite passes. `run: bash tests/run.sh`

## Rollback
- [x] {assert} Migration output records its mapping and archived originals under `.spectacular/archive/id-migrations/`.
- [x] {assert} No automatic repository migration was performed while shipping the contract.
