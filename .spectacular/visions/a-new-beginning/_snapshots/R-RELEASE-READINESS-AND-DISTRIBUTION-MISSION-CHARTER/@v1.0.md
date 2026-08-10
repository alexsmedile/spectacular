---
type: mission-charter
mission: release-readiness-and-distribution
version: 1.0
status: active
prepared_at: 2026-08-10
prepared_by: scenario-r-primary
accepted_by: owner
accepted_at: 2026-08-10
activated_at: 2026-08-10
activation_baseline_commit: 65dbe02ce07e645511f8ac81a5d7170168a56a58
activation_baseline_tree: 25f6966fcf4714730094126413e9ee2e6d871a48
branch: codex/feat/v2-r-release-readiness
source_thread: 019febfa-844e-71e2-b6e6-ffbadcb6ba69
design_sufficiency: sufficient
slice_quality: coherent
repair_budget: 2
repair_rounds_used: 0
upstream:
  - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@2.0
  - S-SKILL-AND-RUNTIME-PREREQUISITES-MISSION-CHARTER.md@1.1
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
  - RESPONSIBILITY-PLACEMENT-CONTRACT.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - SHARED-SCAFFOLD-CONTRACT.md@1.0
---

# Scenario R Mission Charter — Release readiness and distribution

## Outcome

Make the accepted v2 CLI and guided Skill locally release-ready as one aligned version: reproducible
checksummed native artifacts for supported macOS and Linux targets, a fail-safe installer/update and
recovery path requiring no user Go runtime, disposable Codex and Claude discovery, complete v1-runtime
exclusion, and a clean user-facing governed-work smoke flow.

## Preparation verdict

- Design Sufficiency: `sufficient`. Accepted architecture fixes native Go distribution, checksum
  verification, v2-only runtime packaging, public grammar, provider boundaries, and evidence rules.
- Slice Quality: `coherent`. Build, package, install, recovery, and smoke evidence share one release
  integrity boundary and leave a usable local release candidate.
- Decision delta: no publication, signing, platform, provider, or irreversible Type-1 choice remains.
  Exact archive encoding, installer layout, and version injection are reversible Type-2 choices.

## Objectives

1. Build byte-reproducible `darwin/arm64`, `darwin/amd64`, `linux/arm64`, and `linux/amd64`
   archives with stable metadata and a deterministically ordered SHA-256 manifest.
2. Bind runtime version inspection, both v2 plugin manifests, the canonical v2 Skill, and the
   generated mechanical catalog to one release version without adding semantic command vocabulary.
3. Select and verify a locally sourced artifact before mutation; install/update only v2 surfaces;
   fail safely; and provide explicit rollback, recovery, and uninstall behavior without Go.
4. Prove fresh Codex and Claude installation/discovery, a guided-orient to governed closure/resume
   smoke flow, and mechanical exclusion of v1 runtime, aliases, readers, tests, Skill references,
   and hidden dependencies from release assembly.

## Completion boundary

- All four supported artifacts rebuild byte-for-byte and their checksum manifest is stable.
- An installed binary reports the same version carried by both v2 manifests and generated catalog.
- Local-source install and update verify checksums before changing the prefix and preserve a usable
  rollback point; mismatch and unsupported-platform cases leave the prefix unchanged.
- Fresh disposable Codex and Claude roots contain only the native v2 binary, canonical v2 Skill,
  generated interface, and the selected runtime manifest.
- The installed runtime completes orientation, governed Mission work, Evidence/assessment,
  reconciliation/closure, and cold resume through accepted v2 surfaces.
- Release evidence states `locally release-ready` and makes no publication, signing, notarization,
  provider-observation, or global-install claim.

## Owned paths

```text
v2/VERSION
v2/cmd/
v2/internal/buildinfo/
v2/internal/command/
v2/skills/spectacular/generated/
v2/.codex-plugin/
v2/.claude-plugin/
v2/install/
v2/release/
v2/testdata/
.spectacular/visions/a-new-beginning/R-RELEASE-READINESS-AND-DISTRIBUTION-MISSION-CHARTER.md
.spectacular/visions/a-new-beginning/R-RELEASE-READINESS-AND-DISTRIBUTION-MISSION-CHARTER.md.sha256
.spectacular/visions/a-new-beginning/evidence/scenario-r-*.md
.spectacular/visions/a-new-beginning/evidence/scenario-r-*.md.sha256
```

Supporting v2 files may receive bounded changes only when required by this outcome and no semantic
invariant, dependency, or public vocabulary changes.

## Authority and prohibited effects

Authorized effects are inspection, local edits, checks, deterministic local artifact assembly,
disposable install roots, coherent local commits, and one fresh exact-head independent review.
Prohibited effects are every v1-path modification; local-Go fallback in the installer; network
fetches in proof; tag, push, PR, upload, publication, deployment, signing, notarization, credential
use, provider mutation, global/user-home installation, v1 migration, or real-project cutover.

The process-level `--version` inspection flag is the narrow public-contract delta required to prove
artifact identity. It reports build metadata only, creates no semantic noun or operation, and is
included in generated release metadata. No `version show` command is introduced.

## Proof and repair

Proof includes full v2 format/vet/test/race/build checks; double assembly comparison; exact archive
inventory and dependency checks; checksum-mismatch and unsupported-platform refusal; install,
same-version update, rollback, uninstall, and recovery; disposable Codex/Claude discovery; generated
interface/version parity; and installed-binary smoke/cold-resume behavior. The v1 suite is not an
acceptance gate and no v1 path is modified.

Two hypothesis-changing repair rounds are authorized. A fresh reviewer who did not author the work
must bind findings to the exact final commit/tree and reproduce release-integrity checks. Stop only
for a new Type-1 decision, authority/safety/irreversible/provider conflict, missing publication or
signing authority, exhausted repair budget, or an unresolved required check.
