# Test Pyramid Architecture & Latency Budgets

A production test suite divides verification into distinct operational tiers. Lower tiers run frequently on every keystroke or save; higher tiers execute in CI before merging.

---

## 1. The 4 Verification Tiers

```text
       ▲
      / \        Tier 3: Acceptance & Smoke (CLI fixtures, end-to-end user journeys)
     /---\       [Default budget: ≤ 60s]
    /     \      Tier 2: Hardened & Integration (isolated temp storage, race detector, fault injection)
   /-------\     [Default budget: ≤ 20s]
  /         \    Tier 1: Pure Unit Tests (in-memory logic, zero I/O, domain invariants)
 /-----------\   [Default budget: ≤ 5s]
/             \  Tier 0: Preflight Static Sanity (linting, secret scan, git untracked check)
---------------  [Measured project budget; 1s is an initial target]
```

### Tier 0: Preflight Static Sanity (Measured Project Budget)
- **Goal**: Reject known syntax, policy, or contract failures before expensive downstream verification. Measure the actual runner and use the repository's budget.
- **Scope**: Syntax checking, formatting, linters (fast pass), secret leak scan (`gitleaks`), untracked working tree sanity.
- **Rule**: A failed required preflight blocks heavier tiers. Dirty state alone is not a failure: preserve user changes and classify new files according to repository policy.

### Tier 1: Pure In-Memory Unit Tests (Default: $\le 5$ Seconds)
- **Goal**: Instant feedback loop during development.
- **Scope**: In-memory domain objects, pure functions, state machines, algorithmic invariants.
- **Constraints**: **Zero disk I/O, zero network sockets, zero child processes**. If a test writes to disk or binds a port, it belongs in Tier 2.

### Tier 2: Hardened Integration & Concurrency (Default: $\le 20$ Seconds)
- **Goal**: Verify stateful persistence and concurrency safety.
- **Scope**: Real temporary filesystems (`t.TempDir()`), atomic rollback verification, database transactions, concurrency stress with race detectors enabled (`go test -race`, ThreadSanitizer).

### Tier 3: Acceptance & Subprocess Smoke (Default: $\le 60$ Seconds)
- **Goal**: Verify assembled executables and external contracts.
- **Scope**: Subprocess execution of compiled binaries, CLI exit codes, HTTP server round-trips, cross-platform path validation.

---

## 2. Latency Budget Scaling & Monorepo Overrides

> [!IMPORTANT]
> **Hard vs Configurable Gates**:
> - **Tier 0 success** is a prerequisite when the repository requires it. Its latency target is measured and configurable; do not weaken checks to force a sub-second result.
> - **Tiers 1–3 latency numbers are defaults for modular services**. In large monolithic repositories (e.g. 5,000+ unit tests), enforce budgets per-package or per-module rather than artificially dropping test coverage to meet a global timer.
> - Use existing repository configuration or the verification driver for budgets. Do not create a configuration file solely to satisfy this example.

---

## 3. Cross-Language Tooling & Spectacular Mapping

| Tier | Go Ecosystem | TypeScript / Node | Python | Rust | Spectacular `test/verify.sh` |
|---|---|---|---|---|---|
| **Tier 0** (Preflight) | `golangci-lint run --fast`, `gitleaks` | `eslint --max-warnings 0`, `biome check` | `ruff check`, `ruff format --check` | `cargo check`, `cargo clippy -- -D warnings` | `verify.sh preflight` |
| **Tier 1** (Unit) | `go test ./internal/... -short` | `vitest run --testPathIgnorePatterns integration` | `pytest -m "not integration" -q` | `cargo test --lib` | `verify.sh quick` |
| **Tier 2** (Hardened) | `go test -race ./...` | `vitest run --testPathPattern integration` | `pytest -m integration` | `cargo test --tests` | Project-selected hardened/race gate |
| **Tier 3** (Acceptance) | End-to-end fixture tests against built binary | Playwright smoke, CLI subprocess tests | `pytest -m acceptance` | `cargo test --test cli_e2e` | `verify.sh acceptance` |
| **Tier 4** (Matrix / CI) | Multi-OS cross-compilation & checksum verification | Multi-Node matrix (`18`, `20`, `22`) | Multi-Python matrix (`3.10`–`3.12`) | Cross-target matrix (`x86_64`, `aarch64`) | `verify.sh release` |
