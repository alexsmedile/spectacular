---
type: central-acceptance-evidence
mission: release-readiness-and-distribution
status: accepted
accepted_at: 2026-08-10
accepted_by: central-orchestration
feature_commit: b8e0d4d7046d3afde7e61ab752191c41f606025e
feature_tree: 57e5cdccc76524edbd05be3a9f6687a93fabed08
merge_commit: 7837bf3d210985d018c2471eda3c232be4bea643
merge_tree: 57e5cdccc76524edbd05be3a9f6687a93fabed08
independent_review: scenario-r-independent-review.md
disposition: accept
publication_state: unpublished-local-artifacts-only
next_action: prepare-v2-root-cutover-mission
---

# Scenario R central acceptance

Central orchestration accepts and integrates the release-readiness and distribution Mission. The
reviewed implementation supplies deterministic v2.0.0 archives for Darwin/Linux on amd64/arm64,
ordered checksums, aligned CLI/Skill/runtime manifests, a verified no-Go installer, safe
update/rollback/recoverable-uninstall behavior, Codex/Claude staging, and a full governed-work smoke
and cold-resume proof without shipping or depending on v1 runtime surfaces.

## Acceptance evidence

- The feature branch descends from accepted baseline `65dbe02…` and changed only `v2/` plus the
  Scenario R charter, snapshot, and evidence records.
- All Scenario R SHA-256 sidecars verify.
- The independent reviewer inspected commit `7978477…`, reproduced the release artifacts and
  installer lifecycle, and returned `ACCEPT` with no findings.
- Commits after the reviewed head contain only the immutable review, charter reconciliation,
  snapshot, sidecars, and Mission conclusion; no product or release code changed after review.
- Central reproduced module verification, vet, full tests, race tests, command builds, Bash syntax,
  and `v2/release/test.sh` from a clean `git archive` using a disposable Go build cache.
- The resulting merge commit is `7837bf3…`, with tree `57e5cd…`, identical to the final feature tree.

## Limits

This acceptance proves local release readiness only. It does not claim cross-host execution,
publisher authentication, a final v1 tag, root cutover, push, PR, GitHub release, upload,
publication, signing, notarization, provider mutation, global installation, or deployment.

The v2 root cutover is now the only next-ready preparation action. This acceptance does not activate
that Mission or authorize deletion, tag creation, push, publication, or any provider effect.
