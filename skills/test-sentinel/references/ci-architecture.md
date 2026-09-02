# Senior CI/CD Architecture & Pipeline Engineering

Production CI/CD pipelines are critical production software: they must be fast, deterministic, secure, and observable. This guide codifies the non-negotiable architectural standards for enterprise delivery engines.

---

## 1. Build and Artifact Management

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ 📦 BUILD ONCE, PROMOTE EVERYWHERE                                           │
│ Source Code ──> [Build & Test] ──> Immutable Artifact (Image/Binary)        │
│                                           │                                 │
│                     ┌─────────────────────┴─────────────────────┐           │
│                     ▼                                           ▼           │
│               Staging Deploy                            Production Deploy   │
│         (Inject Staging Config)                     (Inject Production Config)│
└─────────────────────────────────────────────────────────────────────────────┘
```

- **Build Once, Promote Everywhere**: Compile binaries or build container images exactly once early in the pipeline. Never rebuild between staging and production—promote the identical, immutable artifact (by digest or git SHA). Rebuilding breaks determinism and invalidates staging tests.
- **Runtime Configuration Injection**: Artifacts must be environment-agnostic. Never bake secrets, base URLs, or environment flags into build images. Inject them dynamically at runtime via HashiCorp Vault, AWS Secrets Manager, or Kubernetes ConfigMaps.
- **Deterministic Layer & Dependency Caching**:
  - Key cache storage strictly by lockfile checksum (`go.sum`, `package-lock.json`, `poetry.lock`, `Cargo.lock`).
  - Order Dockerfiles from least-frequently changed to most-frequently changed: `OS dependencies → Package manifests → Install dependencies → App source`.
  - Use BuildKit remote inline cache backends to avoid rebuilding unchanged layers.

---

## 2. Speed and Pipeline Efficiency (The Sub-10-Minute SLA)

Developer PR throughput plummets when CI exceeds 10–12 minutes. Target median runs of **5–8 minutes**.

| Optimization | Technique | Impact |
|---|---|---|
| **Path Filtering** | Add `paths-ignore: ['docs/**', '**.md', '.spectacular/**']` | Skips CI entirely on documentation-only edits. |
| **Shallow Clones** | Set `fetch-depth: 1` in `actions/checkout` | Reduces clone time from 30s+ to sub-second on large repos. |
| **Parallel Sharding** | Shard test suites across matrix nodes (`--shard=1/4`, `pytest-xdist`) | Cuts 40m end-to-end suites down to 10m. |
| **RAM Disks (`tmpfs`)** | Mount integration DBs (Postgres, Redis) to `tmpfs` | Eliminates disk I/O bottlenecks during heavy teardown/setup. |
| **Fail-Fast Gating** | Run Tier 0 preflight first; fail immediately on lint/syntax error | Cancels expensive runner matrices before consuming build minutes. |

---

## 3. Branch Protection and Gatekeeping

Protecting the trunk preserves git bisectability, prevents release deadlocks, and eliminates ghost merge breakages.

- **Strict Linear History**: Enforce squash-merging or rebase-merging on the default branch. Banish messy merge bubbles; every commit on `main` must represent an atomic, revertible change.
- **Require Branches to Be Up to Date (With Merge Queue)**: Enforcing *"Require branches to be up to date before merging"* with `cancel-in-progress: true` on busy trunks creates a rebase-requeue treadmill: every merge invalidates all other open PRs.
  - **Mitigation**: Pair up-to-date branch protection with **GitHub Merge Queue** (`merge_group` trigger, `max_entries_to_build: 5`). The queue batches and validates speculative merge commits without bouncing developers back into manual rebase loops.
- **Block Direct Pushes to Main**: Disallow direct pushes to default branches. Require pull requests with at least one approval from designated `CODEOWNERS`.
- **Lean Mandatory Status Checks**: Define a minimal set of blocking checks (Tier 0 preflight, core unit tests, secret scanning). Slow acceptance or multi-OS matrix jobs should run asynchronously or on nightly schedules so they don't bottleneck daily engineering throughput.

### Canonical GitHub Ruleset Configuration
```json
{
  "name": "Trunk Protection",
  "target": "branch",
  "enforcement": "active",
  "conditions": { "ref_name": { "include": ["refs/heads/main"] } },
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" },
    { "type": "required_linear_history" },
    { "type": "required_status_checks", "parameters": {
      "strict_required_status_checks_policy": true,
      "required_status_checks": [
        { "context": "Tier 0 Preflight" },
        { "context": "Tier 1 & 2 Test Suite" }
      ]
    }},
    { "type": "pull_request", "parameters": {
      "required_approving_review_count": 1,
      "dismiss_stale_reviews_on_push": true,
      "require_code_owner_review": true,
      "required_review_thread_resolution": true
    }}
  ]
}
```

---

## 4. Pull Request Ergonomics & Blast Radius

- **Small, Single-Responsibility PRs (Soft Guideline)**: Target PR diffs of **$\le 300\text{--}400$ lines** of changed code.
  - *Enforcement*: Treat line caps as a **soft review guideline**, not a blocking automated status check. Hard blocking checks incentivize stacked PRs that game the metric or artificial splits that break atomicity. If a PR exceeds 400 lines, require explicit owner justification.
- **Automated PR Templates**: Require authors to state context, acceptance criteria, reproduction steps, and links to tracking records (Issues or Missions).
- **Bot Feedback & Guardrails**: Automate code hygiene via bots (e.g. DangerJS, Codecov) to post actionable feedback—such as test coverage drops, missing changelog entries, or bundle-size deltas—directly into PR comments without human reviewer fatigue.

---

## 5. Security & Isolation for Pull Requests

- **Zero Secrets on Fork PRs**: Never expose read/write repository secrets or deployment credentials to pull requests originating from forks.
- **Two-Tier Trigger Model**:
  - `pull_request`: Runs unprivileged, read-only tests with zero secrets for all incoming branches and forks.
  - `workflow_run` or push to main: Triggers downstream deployment workflows only after code is merged into the protected trunk.
- **Quarantined Workflow Permissions**: Set organizational and repository default `GITHUB_TOKEN` permissions to **read-only**. Grant explicit write permissions (e.g. `contents: read`, `pull-requests: write`, `id-token: write`) on a per-job basis.
- **Push Protection Secret Scanning**: Block commits containing API keys, private certificates, or tokens using pre-receive hooks or GitHub push protection before secrets ever land in remote git history.

---

## 6. Continuous Verification & Preview Environments

- **Ephemeral Preview Environments**: Spin up lightweight, isolated preview environments per PR (via Vercel, Cloudflare Pages, or ArgoCD ApplicationSets). Automatically tear them down upon PR close or merge.
- **Policy-as-Code & IaC Drift Verification**: For infrastructure-as-code (Terraform, OpenTofu), run automated `plan` verifications in PRs and post plan diffs to comments. Use Open Policy Agent (OPA) or Conftest to reject non-compliant cloud security groups and IAM roles before merge.

---

## 7. Operational Profiles: Enterprise vs. Open-Source

| Dimension | Enterprise Internal Profile | Open-Source Public Profile |
|---|---|---|
| **Auth** | Keyless OIDC to private AWS/GCP | Zero cloud auth in PRs |
| **Forks** | Same-repo feature branches; secrets accessible | Fork PRs strictly sandboxed; 0 secrets |
| **Gatekeeping** | `CODEOWNERS` + 1 required approval | Maintainer approval + fork approval gate |
| **Previews** | Ephemeral Kubernetes/Argo preview envs | Static preview deployments (Vercel/CF) |
| **Merge Strategy** | Squash & merge (linear history) | Rebase or Squash & merge |
