---
type: independent-review
mission: release-readiness-and-distribution
status: accepted
observed_at: 2026-08-10
actor: scenario-r-independent-reviewer
reviewed_commit: 79784772b6bfc7a358852cedb3b7606a242a0030
reviewed_tree: 521cb225582580545f6ec36d891e3a3e86d5f512
baseline_commit: 65dbe02ce07e645511f8ac81a5d7170168a56a58
baseline_tree: 25f6966fcf4714730094126413e9ee2e6d871a48
release_version: 2.0.0
verdict: ACCEPT
publication_state: unpublished-local-artifacts-only
---

# Scenario R independent release-integrity review

## Verdict

`ACCEPT`. No blocking release-integrity, scope, recovery, or installed-runtime finding was found.
The reviewed v2 release is locally release-ready within the accepted Scenario R charter.

This review is bound to commit `79784772b6bfc7a358852cedb3b7606a242a0030`, tree
`521cb225582580545f6ec36d891e3a3e86d5f512`, against baseline commit
`65dbe02ce07e645511f8ac81a5d7170168a56a58`, tree
`25f6966fcf4714730094126413e9ee2e6d871a48`. The release payload embeds implementation commit
`f8102a8a4e978f79165344dec0170de0171c5e8a`; the reviewed head adds only the scoped charter and
primary evidence record after that implementation commit.

## Scope and source review

- Read the repository and workspace agent rules, the accepted Scenario R charter, checksummed
  primary evidence, the complete canonical v2 Skill bundle, and directly relevant release,
  installer, manifest, catalog, build-info, smoke, and test sources.
- `git diff` from the specified baseline changes only `v2/` and Scenario R charter/evidence paths.
  Root `cli/`, `tests/`, `skills/`, manifests, docs, tracking files, and the live legacy workspace
  are unchanged. No root installer, v1 deletion, archival, or cutover is present.
- A fresh source tree was created with `git archive <reviewed-head> v2`, so release assembly had no
  repository parent and could not source root v1 surfaces. Exact archive inventories contain only
  the native binary, release metadata, both v2 manifests, and the canonical v2 Skill bundle.
- The v2 CLI dependency graph contains only the v2 module, `github.com/google/uuid`, and
  `go.yaml.in/yaml/v3`; no root v1 module dependency is present. Skill/package scans found none of
  the prohibited v1 aliases, readers, request paths, or version archives.

## Reproduced checks

| Boundary | Independent observation | Result |
|---|---|---|
| Version and interface parity | `v2/VERSION`, installed `--version`/JSON, both v2 manifests, Skill frontmatter, regenerated JSON catalog, and regenerated Markdown catalog all resolve to `2.0.0` | pass |
| Platform identity | Assembly produced Mach-O `amd64`/`arm64` and ELF `amd64`/`arm64` artifacts; archive `RELEASE.json` identifies the matching OS and architecture | pass |
| Deterministic archives | The standard double-assembly comparison passed; a separate `git archive` source extraction reproduced the primary artifact SHA-256 values exactly | pass |
| Canonical metadata | Unit inspection and independent archive listing confirmed stable order, epoch timestamps, uid/gid `0`, root names, modes `0644`/`0755`, and regular-file-only assembly | pass |
| Checksums and refusal | Ordered `SHA256SUMS` verified all four artifacts; checksum corruption and unsupported platform selection refused without altering sentinel state | pass |
| Installer safety | Absolute-prefix, protected symlink, archive inventory, binary/catalog/Skill/manifest parity, and pre-mutation checksum gates were inspected; Codex proof and an additional Claude symlink escape probe both refused safely | pass |
| Lifecycle and no user Go | Install, same-version update, rollback, recoverable uninstall, and recovery passed for both runtimes; an additional full Claude lifecycle passed with `PATH=/usr/bin:/bin`, demonstrating that the user installer does not invoke Go | pass |
| Runtime packaging | Disposable Codex and Claude roots contain the native binary, canonical v2 Skill/catalog, and only the selected runtime manifest | pass |
| Governed smoke and cold resume | The installed binary completed `@Orient`, Proposal/Mission preparation and activation, Handoff/return, Evidence and assessment, reconciliation, resolution, `@Resume`, archive, and final orientation for both runtimes | pass |
| v2 quality gates | `gofmt`, module verification with `GOPROXY=off`, vet, full tests, race tests, three command builds, Bash syntax, Scenario S staging, Scenario R release proof, generated-catalog comparison, and repository version guard | pass |

The independently assembled artifact digests were:

```text
329094d5d7d304f817db4b4fbc0a16cf80385eace63d94d9c7cc3f28d6bab679  spectacular-v2.0.0-darwin-amd64.tar.gz
4fe7a94ddbd4429ae61fceaeece3fe16cb51f365f2dddcbc472278ff24ed5cbe  spectacular-v2.0.0-darwin-arm64.tar.gz
6fa9f92995d5593963ed66bfa8665e27ebbfa4df9ebfbb0700ae3aa4123198bb  spectacular-v2.0.0-linux-amd64.tar.gz
96717c1f571bf9d2e891851191beb669ba2073e0cc73b796f345912586a09ad1  spectacular-v2.0.0-linux-arm64.tar.gz
```

## Limits and publication state

- Only the host-native `darwin/arm64` binary was executed. Cross-target identity was verified by
  parsing Mach-O/ELF headers; cross-host execution remains a later publication-pipeline concern.
- SHA-256 verifies artifact bytes against the local manifest; it does not authenticate a publisher.
  Signing and notarization were not authorized or attempted.
- No v1 acceptance suite was run. Root v1 surfaces were inspected only through read-only diff and
  version-consistency checks and are outside Scenario R acceptance.
- No tag, push, PR, upload, publication, deployment, provider mutation, credential use, global
  install, real-project migration, root cutover, deletion, or archival occurred.

Publication state remains `unpublished-local-artifacts-only`. Root cutover and remote publication
remain separate owner-authorized Missions/gates.
