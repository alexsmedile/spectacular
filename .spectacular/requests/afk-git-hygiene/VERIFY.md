# Verify — afk-git-hygiene

## Policy and playgrounds
- [x] AFK policy, branch naming, playground cleanup, and PR handoff fixtures pass. `run: bash tests/cli/afk-git-hygiene.test.sh`
- [x] Existing wayfinding contracts and sequencing remain compatible. `run: bash tests/cli/wayfinding-contract.test.sh && bash tests/cli/wayfinding-sequencer.test.sh`
- [x] {assert} Status and configuration are read-only/dry-run by default, and host branch prefixes compose with every Spectacular branch class.
- [x] {assert} Branch start rejects dirty repositories and records canonical provenance after explicit authorization.

## Cleanup and handoff
- [x] {assert} Cleanup refuses remote deletion, archives disposition/outcome/evidence before local deletion, and documents reflog recovery.
- [x] {assert} PR handoff requires a current spec, verified request evidence, passing tests, and separate breaking-change approval; it never merges.
- [x] Bash syntax and guarded versions remain valid. `run: bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit && scripts/hooks/pre-commit --check`
- [x] Complete regression suite passes. `run: bash tests/run.sh`
