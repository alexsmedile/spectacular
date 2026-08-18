---
type: Mission
id: 01a01533-1f38-74f4-86f5-c6a84d1e19da
title: Modernize CI/CD pipeline, security scanning, and release automation
status: completed
created: "2026-08-18T14:35:39Z"
updated: "2026-08-18T14:38:28Z"
activation:
    at: "2026-08-18T14:35:39Z"
    by: Alex
    fingerprint: sha256:41ae520071b16ecdc68d816aa6867906a3bfed58ee77de98115ca9ce9f351e21
authority:
    operator:
        - inspect
        - edit-in-scope
        - choose-reversible-implementation
        - run-checks
        - generate-derived-files
        - bounded-repair
        - commit-local
    requires_owner:
        - activate-mission
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
        - secret-change
baseline:
    branch: m13-agent-ready-ci-cd
    commit: eb288796d8e6e163d42accfdd50ddd964f54e835
completion:
    - claim: ci-workflow-is-staged-cached-and-concurrent
      pass_boundary: .github/workflows/verify.yml has concurrency controls with cancel-in-progress, enables Go build/module caching, and decomposes verification into four visible parallel/dependent stages (static analysis, unit & race tests, acceptance fixtures, and Scenario R/S distribution proofs) instead of a monolithic single step.
      proof_requirement: Workflow file validates against GitHub Actions schema; push/PR triggers work with concurrency and caching; full verification suite completes within 60 seconds.
    - claim: manifest-sync-and-security-guardrails-are-enforced
      pass_boundary: test/verify.sh static checks assert that root plugin.json, vendor manifests, and VERSION are synchronized, preventing release drift (P4). Secret scanning with gitleaks is integrated when present.
      proof_requirement: A version mismatch between plugin.json and VERSION fails static checks. Clean state passes.
    - claim: coverage-reporting-is-published-in-ci
      pass_boundary: CI runs unit tests with -coverprofile and uploads an atomic coverage report via actions/upload-artifact, providing background agents and reviewers coverage deltas.
      proof_requirement: CI test step executes with -coverprofile=coverage.out and uploads coverage-report artifact.
    - claim: continuous-delivery-release-workflow-is-automated
      pass_boundary: A dedicated .github/workflows/release.yml triggers on git tags (refs/tags/v*), executes test/verify.sh all, compiles the 4-platform archives via cmd/assemble-release, and publishes the GitHub Release with attached assets and SHA256SUMS.
      proof_requirement: .github/workflows/release.yml triggers on tags v*, builds with cmd/assemble-release, and attaches the 4 archives and SHA256SUMS.
completion_record:
    at: "2026-08-18T14:38:28Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: b41ea052a2bfe484838054ff1c8d7d6fd83eb8c6
contract:
    fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a
    ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
dependencies: []
gaps: []
objectives:
    - claims:
        - ci-workflow-is-staged-cached-and-concurrent
        - coverage-reporting-is-published-in-ci
      id: 01a01533-1f38-7332-83a7-9c1ab2574856
      outcome: Decompose verify.yml into staged jobs with concurrency, caching, and coverage upload.
      ref: O1
      status: implemented
    - claims:
        - manifest-sync-and-security-guardrails-are-enforced
      id: 01a01533-1f38-7a35-993a-5ebb26ab9cb2
      outcome: Add manifest version sync check and host-aware gitleaks to test/verify.sh.
      ref: O2
      status: implemented
    - claims:
        - continuous-delivery-release-workflow-is-automated
      id: 01a01533-1f38-7412-97a1-cd655e4bc77b
      outcome: Create release.yml workflow for automated tag-based GitHub Releases.
      ref: O3
      status: implemented
outcome: CI/CD workflow is broken down into staged parallel jobs with concurrency and caching, secret scanning and manifest synchronization are verified in static checks, and release distribution on git tags is automated.
owner: Alex
ref: M13
repair_budget: 3
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-m13-ci-cd-pipeline-and-release-automation.md
      id: 01a0154c-a278-77b1-a9d0-17f4cbc18f2e
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a01533-1f38-7455-a4bd-c02460d9f52d
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-18T14:35:39Z"
    status: completed
scope:
    mechanical:
        - .github/workflows/
        - test/
        - .spectacular/
    semantic:
        - CI workflow modernization
        - manifest synchronization check
        - release automation.
start_key: sha256:11b8064e1b9adba219c9c4693444d1b9033881667fcfff205fc0afdc86cf68fd
stops:
    - scope-drift
    - public-command-count-change
validation:
    mode: cli
    schema: mission.v2
---
# Modernize CI/CD pipeline, security scanning, and release automation

Modernize Spectacular v2 CI/CD automation and release distribution.
