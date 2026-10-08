---
name: test-sentinel
description: >-
  Design deterministic test suites, concurrency stress tests, regression anchors, and pinned GitHub Actions CI pipelines.
  Triggers on "test suite", "stress test", "race condition", "flaky test", "regression test", "ci cd", "github actions", "architecture tests", or "dependency boundary checks".
  Do not invoke for generic code styling, non-test refactors, or writing documentation reports.
version: 0.4.1
category: devtools
status: draft
tags: [testing, ci-cd, determinism, regression, github-actions, stress-testing]
---

# Test Sentinel

Deterministic test suites, adversarial stress testing, regression shields, and SHA-pinned GitHub Actions pipelines.

## 1. Route Matrix

| Route | Trigger / Need | Core Action | Complete When |
|---|---|---|---|
| **Preflight** | Fast static sanity, secret check | Lint + secret scan (`gitleaks`) | Exits 0 within the project's measured budget. Read [test pyramid](references/test-pyramid.md). |
| **Unit** | Domain invariant verification | In-memory tests; zero disk/network I/O | All tests pass within the project's measured budget. Read [test pyramid](references/test-pyramid.md). |
| **Hardened** | Concurrency, race, leak verification | ThreadSanitizer (`-race`), ephemeral tempdirs | 0 races, 0 deadlocks under load. Read [determinism](references/determinism-matrix.md). |
| **Regression** | Bug intake or incident fix | TDD reproduction test first $\to$ fix $\to$ anchor | Fails on trunk, passes on fix. Read [regression shield](references/regression-shield.md). |
| **Pipeline** | Automated PR & matrix gate | Update the existing SHA-pinned workflow | Existing workflow passes preflight before its matrix. Read `templates/`. |
| **Benchmark** | Model trials, token sweeps, evals | User-invokable A/B evaluations & parallel harness trials | Evaluates context economy & regressions. Read the project's existing benchmark runner and cost policy. |
| **Module boundaries** | Import rules, isolation, compatibility, atomicity | Execute positive and negative checks against declared boundaries | Forbidden dependencies and mutations fail; allowed behavior passes. |
| **Architecture** | CI/CD audit, gatekeeping & delivery | Audit immutable builds, OIDC, branch protection | Repository delivery policy and measured budgets pass. Read [CI architecture](references/ci-architecture.md). |

---

## 2. Standard Directory & Benchmark Hierarchy

Follow the repository's existing layout and verification driver. The following
layout is an optional starting point; do not relocate tests or add modes solely
to match it. Timing figures are planning targets, measured on the actual runner,
not universal acceptance gates:

```text
test/
├── unit/             (or inline *_test.go in internal/ / src/ — sub-second Tier 1 CI gate)
├── acceptance/       (deterministic integration & CLI fixtures — < 5s Tier 2 CI gate)
├── benchmarks/       (user-invokable model trials, token sweeps & evals — Tier 3 gate)
│   ├── cmd/          (flat benchmark CLI runner: go run ./test/benchmarks/cmd run)
│   ├── adapters/     (harness execution scripts: opencode.sh, claude.sh, codex.sh, agy.sh)
│   ├── cases.json    (evaluation cases, ground-truth prompts, and schemas)
│   └── reports/      (generated golden & A/B comparison markdown reports)
└── verify.sh         (single-entry driver: preflight, quick, acceptance, bench, all)
```

- **Tests vs Benchmarks Separation**:
  - `acceptance/` and `unit/` run on every commit in **GitHub Actions CI** (no model calls; measure runner cost and determinism).
  - `benchmarks/` is **user-invokable** or triggered as an advisory regression gate (supports `--parallel <N>` multi-harness concurrency; protected from accidental token/credit spend).

---

## 3. Direct Negative Constraints (DO NOT)

- **DO NOT write wall-clock sleeps**: Banned: `time.Sleep()`, `setTimeout()`, `time.sleep()`, `thread::sleep()`. Use channel barriers, condition sync, or synthetic clocks (`testing/synctest`, fake timers).
- **DO NOT bind static ports or shared paths**: Banned: `:8080`, `/tmp/test.db`. Use `:0` dynamic port allocation and OS-assigned tempdirs (`t.TempDir()`, `tmp_path`).
- **DO NOT mask flakes with CI retries**: Banned: `retry: 3`, `pytest-rerunfailures`. Quarantine flaky tests immediately; mainline CI must be 100% deterministic.
- **DO NOT allow unexpiring quarantine**: Quarantined tests must have a hard 14-day expiry deadline—either fix the root cause or delete the test.
- **DO NOT rebuild artifacts across environments**: Build binaries and container images once; promote the identical immutable artifact across staging and production.
- **DO NOT bake configuration or secrets into builds**: Inject secrets, base URLs, and environment variables dynamically at runtime (Vault, K8s ConfigMaps).
- **DO NOT store long-lived cloud credentials in CI secrets**: Use OpenID Connect (OIDC) for short-lived, keyless cloud authentication.
- **DO NOT use `pull_request_target` with write tokens**: Never expose deployment secrets to untrusted, unreviewed fork code.
- **DO NOT bypass branch protection or merge stale branches**: Follow the repository's merge and history policy; do not impose a new topology.
- **DO NOT author oversized PRs**: Aim for reviewable, cohesive diffs; 300–400 changed lines is a soft signal. Separate mechanical moves from behavioral changes in the review; follow repository policy without inventing an approval gate.
- **DO NOT create governance records**: `test-sentinel` is read-only on `.spectacular/`. Spectacular owns claims and contracts; `test-sentinel` owns executable test proof.
- **DO NOT generate markdown report sprawl**: Banned: `TEST_PLAN.md`, `COVERAGE.md` in `docs/`. Tests and machine receipts are the only deliverables.
- **DO NOT use floating GitHub Action tags**: Banned: `uses: actions/checkout@v4`. Use verified full commit SHAs (`actions/checkout@<sha> # v4.2.2`).

---

## Module proof and fixture isolation

Consume the owner/API/dependency/transaction facts from `system-architecture`.
Test forbidden imports and unclassified packages, detached read-model aliasing,
consumer ports with substitute adapters, compatible outputs, and rollback of
coordinated writes. Prefer behavior and negative cases over folder/file counts.
Verify the checker itself refuses violations; include platform-specific sources
and connect checks to the existing CI gates. Report unenforced boundaries.

Fixture isolation includes environment and tool configuration, not only tempdirs.
For nested Git repositories, remove or scope inherited `GIT_INDEX_FILE`,
`GIT_DIR`, `GIT_WORK_TREE`, and related overrides. Isolate credentials, working
directory, configuration, and PATH overrides as appropriate. A temporary index
for workspace hygiene must never become the index of a fixture repository.
Preserve the real staged entries; distinguish harness failures from product bugs.

## Claim coverage at milestone gates

For delivery or release proof, use the applicable dimensions in
[Spectacular readiness](../spectacular/references/milestone-readiness.md). Tie
results to the artifact/version, environment, dependency conditions, and material
limitations. Distinguish mocks, integration, real dependencies, and observed user
behavior. Report unverified claims; passing technical checks does not establish
MVP value or publication authority. Recheck affected evidence after material changes.

## 3. Consolidated Command Palette

```bash
# Tier 0 & 1: Fast local checks (≤ 1-5s)
go test -short ./...                              # Go in-memory
npm test -- --testPathIgnorePatterns integration  # Node in-memory
pytest -m "not integration" -q                    # Python in-memory
cargo test --lib                                  # Rust in-memory

# Tier 2: Hardened Concurrency & Stress (≤ 20s default)
go test -v -race -timeout 10m ./...               # Go ThreadSanitizer
npm test -- --ci --runInBand                      # Node isolated
pytest -m integration -v                          # Python integration
cargo test --all-targets --locked                 # Rust locked

# Pipeline templates (examples only; inspect existing files before adapting)
# Select the relevant templates/ci-*.yml and merge into the existing workflow.
# Do not overwrite a workflow with cp or add a duplicate pipeline.
```

---

## 4. The Regression Shield Protocol

When fixing defects: `Failing Repro (Red) → Implement Fix (Green) → Permanent Anchor`.

- For an explicitly selected Mission: use an existing project convention such as **`TestM<N>_<slug>`** (e.g. `TestM14_TokenRefreshRace`).
- For ordinary work, including Spectacular workspaces: use the existing convention or name anchor **`TestRegression_<slug>`** (e.g. `TestRegression_TokenRefreshRace`).

---

## 5. Machine Receipt Standard (`test-sentinel.receipt.v1`)

```json
{
  "schema_version": "test-sentinel.receipt.v1",
  "status": "pass",
  "tier": "tier1-unit",
  "command": "go test -race ./...",
  "duration_ms": 1250,
  "commit": "3a4d2c2",
  "failures": []
}
```
*Read [receipt schema](references/receipt-schema.md) for full specification.*

---

## 6. Expansion Handoffs

| Out-of-Scope Need | Responsible System / Skill |
|---|---|
| Contract drafting, failable claims, mission gates | `spectacular` (`.spectacular/missions/`) |
| Git commit, branch creation, worktrees, PRs | `git-ops` / `gh` CLI |
| Architectural decisions and options comparison | `system-architecture`; governed recording only on explicit owner selection |
| Database schema migration and ER modeling | `data-modeling` |
