---
type: implementation-evidence
mission: release-readiness-and-distribution
status: ready-for-independent-review
observed_at: 2026-08-10
actor: scenario-r-primary
baseline_commit: 65dbe02ce07e645511f8ac81a5d7170168a56a58
baseline_tree: 25f6966fcf4714730094126413e9ee2e6d871a48
implementation_commit: f8102a8a4e978f79165344dec0170de0171c5e8a
implementation_tree: a9812c4beabd0703bd4094ca12df8c0b022db937
release_version: 2.0.0
repair_rounds_used: 0
publication_state: unpublished-local-artifacts-only
---

# Scenario R primary implementation evidence

## Claim boundary

This evidence supports `locally release-ready` for the accepted v2 macOS/Linux and Codex/Claude
boundary. It does not claim a tag, push, PR, upload, endpoint, publication, deployment, signing,
notarization, credential use, provider observation, global installation, root cutover, v1 cleanup,
or real-project migration.

## Exact implementation

- Baseline: commit `65dbe02ce07e645511f8ac81a5d7170168a56a58`, tree
  `25f6966fcf4714730094126413e9ee2e6d871a48`.
- Implementation: commit `f8102a8a4e978f79165344dec0170de0171c5e8a`, tree
  `a9812c4beabd0703bd4094ca12df8c0b022db937`.
- Scope inspection found changes only under `v2/` plus this Mission's charter/evidence paths.
  Root v1 CLI, Skill, manifests, tests, docs, tracking files, and live workspace were not modified.

## Claim-mapped results

| Claim | Direct proof | Result |
|---|---|---|
| One aligned v2 version | `v2/VERSION`; runtime `--version[ --json]`; Skill frontmatter; both v2 manifests; generated JSON/Markdown catalog; preflight and manifest tests | pass at `2.0.0` |
| Four correct native targets | release assembler cross-build plus Mach-O/ELF CPU inspection for `darwin/{amd64,arm64}` and `linux/{amd64,arm64}` | pass |
| Complete-archive reproducibility | two independent clean-v2 source copies; `GOPROXY=off`; byte comparison of archives, `VERSION`, and ordered `SHA256SUMS`; unit inspection of tar/gzip order, ownership, modes, and timestamps | pass |
| Verified installation without user Go | local-source selector and installer verify SHA-256 and archive inventory before mutation; installed binary is precompiled; checksum mismatch and unsupported platform refuse without changing sentinel state | pass |
| Safe update and recovery | disposable native install, same-version update, rollback, recoverable uninstall, recovery, and symlink-path refusal for both Codex and Claude | pass |
| Codex/Claude discovery | each disposable prefix contains the common v2 Skill/catalog, native binary, and only its selected v2 runtime manifest; cold `@Orient` context compiles | pass |
| Realistic v2 smoke | installed binary performs `@Orient`, owner Decisions, Proposal preparation, Mission creation/activation, Handoff/return, Evidence, awaiting-assessment, assessment, reconciliation, resolution, `@Resume`, archive, and final project orientation | pass |
| No v1 runtime ships | assembly runs from copied `v2/` roots with no repository parent; exact inventory whitelist excludes CLI/tests/testdata/requests/version archives; textual Skill rejection and Go dependency check reject v1 aliases/readers/references/module dependencies | pass |
| Existing v2 behavior remains sound | format check, module verification, vet, full tests, race tests, builds, Scenario S staging, Bash syntax, and repository version guard | pass |

## Exact artifact checksums

Locally assembled from the exact implementation commit with canonical release metadata:

```text
329094d5d7d304f817db4b4fbc0a16cf80385eace63d94d9c7cc3f28d6bab679  spectacular-v2.0.0-darwin-amd64.tar.gz
4fe7a94ddbd4429ae61fceaeece3fe16cb51f365f2dddcbc472278ff24ed5cbe  spectacular-v2.0.0-darwin-arm64.tar.gz
6fa9f92995d5593963ed66bfa8665e27ebbfa4df9ebfbb0700ae3aa4123198bb  spectacular-v2.0.0-linux-amd64.tar.gz
96717c1f571bf9d2e891851191beb669ba2073e0cc73b796f345912586a09ad1  spectacular-v2.0.0-linux-arm64.tar.gz
```

The disposable artifact directory was `/tmp/spectacular-r-artifacts-f8102a8-exact`. It is evidence
input, not a published or committed distribution location.

## Reproduced checks

```text
gofmt -l v2/**/*.go                                      clean
cd v2 && GOPROXY=off go mod verify                       pass
cd v2 && GOPROXY=off go vet ./...                        pass
cd v2 && GOPROXY=off go test ./...                       pass
cd v2 && GOPROXY=off go test -race ./...                 pass
cd v2 && GOPROXY=off go build ./cmd/{spectacular,assemble-release,release-smoke}  pass
bash -n v2/install/*.sh v2/release/test.sh                pass
bash v2/install/test.sh                                   pass
bash v2/release/test.sh                                   pass
scripts/hooks/pre-commit --check                          pass
```

The legacy 31-file v1 suite was deliberately not used as a v2 acceptance gate. The version guard
was read-only and confirmed that untouched root/v1 release surfaces remain internally consistent.

## Limitations and next gate

- Only the host-native artifact was executed; cross-target binaries were compiled with `CGO_ENABLED=0`
  and their Mach-O/ELF architecture was parsed directly. Cross-host execution remains a publication
  pipeline concern.
- SHA-256 proves local artifact integrity against the assembled manifest, not publisher identity.
  Signing and notarization remain unauthorized.
- Final acceptance requires an independent reviewer bound to the exact evidence head. Actual
  cutover and push remain separate owner/provider gates.
