---
name: test-sentinel
description: >-
  Design deterministic test suites, concurrency stress tests, regression anchors, and pinned GitHub Actions CI pipelines.
  Triggers on "test suite", "stress test", "race condition", "flaky test", "regression test", "ci cd", or "github actions".
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
| **Preflight** | Fast static sanity, secret check | Lint + secret scan (`gitleaks`) | Exits 0 in $\le 1$s. Read [test pyramid](references/test-pyramid.md). |
| **Unit** | Domain invariant verification | In-memory tests; zero disk/network I/O | All tests pass in $\le 5$s. Read [test pyramid](references/test-pyramid.md). |
| **Hardened** | Concurrency, race, leak verification | ThreadSanitizer (`-race`), ephemeral tempdirs | 0 races, 0 deadlocks under load. Read [determinism](references/determinism-matrix.md). |
| **Regression** | Bug intake or incident fix | TDD reproduction test first $\to$ fix $\to$ anchor | Fails on trunk, passes on fix. Read [regression shield](references/regression-shield.md). |
| **Pipeline** | Automated PR & matrix gate | Deploy SHA-pinned `.github/workflows/ci.yml` | Workflow passes preflight before matrix. Read `templates/`. |
| **Architecture** | CI/CD audit, gatekeeping & delivery | Audit immutable builds, OIDC, branch protection | Linear history, sub-10m SLA, zero secrets. Read [CI architecture](references/ci-architecture.md). |

---

## 2. Direct Negative Constraints (DO NOT)

- **DO NOT write wall-clock sleeps**: Banned: `time.Sleep()`, `setTimeout()`, `time.sleep()`, `thread::sleep()`. Use channel barriers, condition sync, or synthetic clocks (`testing/synctest`, fake timers).
- **DO NOT bind static ports or shared paths**: Banned: `:8080`, `/tmp/test.db`. Use `:0` dynamic port allocation and OS-assigned tempdirs (`t.TempDir()`, `tmp_path`).
- **DO NOT mask flakes with CI retries**: Banned: `retry: 3`, `pytest-rerunfailures`. Quarantine flaky tests immediately; mainline CI must be 100% deterministic.
- **DO NOT allow unexpiring quarantine**: Quarantined tests must have a hard 14-day expiry deadline—either fix the root cause or delete the test.
- **DO NOT rebuild artifacts across environments**: Build binaries and container images once; promote the identical immutable artifact across staging and production.
- **DO NOT bake configuration or secrets into builds**: Inject secrets, base URLs, and environment variables dynamically at runtime (Vault, K8s ConfigMaps).
- **DO NOT store long-lived cloud credentials in CI secrets**: Use OpenID Connect (OIDC) for short-lived, keyless cloud authentication.
- **DO NOT use `pull_request_target` with write tokens**: Never expose deployment secrets to untrusted, unreviewed fork code.
- **DO NOT bypass branch protection or merge stale branches**: Enforce strict linear history (squash/rebase) and require branches to be up to date (or use GitHub merge queues) before merging.
- **DO NOT author oversized PRs**: Cap PR diffs at $\le 300\text{--}400$ lines (soft review guideline; require explicit owner justification if exceeded).
- **DO NOT create governance records**: `test-sentinel` is read-only on `.spectacular/`. Spectacular owns claims and contracts; `test-sentinel` owns executable test proof.
- **DO NOT generate markdown report sprawl**: Banned: `TEST_PLAN.md`, `COVERAGE.md` in `docs/`. Tests and machine receipts are the only deliverables.
- **DO NOT use floating GitHub Action tags**: Banned: `uses: actions/checkout@v4`. Use verified full commit SHAs (`actions/checkout@<sha> # v4.2.2`).

---

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

# Pipeline Deployment (copy and adapt template)
cp templates/ci-go.yml .github/workflows/ci.yml
cp templates/ci-node.yml .github/workflows/ci.yml
cp templates/ci-python.yml .github/workflows/ci.yml
cp templates/ci-rust.yml .github/workflows/ci.yml
```

---

## 4. The Regression Shield Protocol

When fixing defects: `Failing Repro (Red) → Implement Fix (Green) → Permanent Anchor`.

- In Spectacular workspaces: Name anchor **`TestM<N>_<slug>`** (e.g. `TestM14_TokenRefreshRace`).
- In standalone workspaces: Name anchor **`TestRegression_<slug>`** (e.g. `TestRegression_TokenRefreshRace`).

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
| Architectural decisions and options comparison | `system-architecture` (`spectacular decide`) |
| Database schema migration and ER modeling | `data-modeling` |
