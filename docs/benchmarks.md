# Spectacular Benchmarks & Model Evaluations

Spectacular includes a dedicated benchmark suite designed to measure prompt context economy, decision fidelity, and AI agent behavior across different LLM harnesses.

---

## 1. Directory Structure

Benchmarks live in a flat, predictable directory structure under `test/benchmarks/`:

```text
test/
├── acceptance/           # Deterministic Go integration fixtures (< 5s, runs in CI)
│   └── cli_test.go
│
├── benchmarks/           # Consolidated Model Benchmarks & Token Sweeps (Flat)
│   ├── cmd/              # Flat CLI runner (go run ./test/benchmarks/cmd ...)
│   ├── adapters/         # Multi-harness adapters (opencode, claude, codex, agy)
│   ├── evals.json        # Benchmark test catalog & ground-truth assertions
│   ├── reports/          # Generated A/B comparison and static golden reports
│   └── *_test.go         # Fast Go evaluation unit tests (charter budget, trace decoders)
│
└── verify.sh             # Canonical verification router (preflight, quick, acceptance, bench, all)
```

---

## 2. CI Tests vs. Live Model Benchmarks

Spectacular strictly separates **continuous integration tests** from **model benchmarks**:

| Dimension | CI Tests (`verify.sh quick` / `acceptance`) | Model Benchmarks (`verify.sh bench` / `cmd/bench`) |
|---|---|---|
| **Execution Trigger** | Automatic on every git commit / PR (GitHub Actions). | **User-invokable** on demand OR as a post-mission regression gate. |
| **Duration & Cost** | Sub-second to sub-5-seconds. **Zero token cost**, zero network calls. | 1 to 5 minutes. Sends live model API calls to evaluate agent behavior. |
| **Purpose** | Guarantees code doesn't crash, types match, and contracts hold. | Measures **AI agent behavior**, prompt context savings, and model accuracy. |
| **Concurrency** | Go `-race` parallel execution. | **Multi-Harness Parallel Dispatch**: Runs trials concurrently with `--parallel <N>`. |

---

## 3. How to Run Benchmarks

### A. Local Static & Token Sweeps (Zero API Cost)

To evaluate prompt context size across all mission files or compare static reductions against a baseline without spending LLM credits:

```bash
# 1. Run local token sweep across all mission files
go run ./test/benchmarks/cmd/sweep-tokens

# 2. Run static A/B context comparison against baseline
go run ./test/benchmarks/cmd static --candidate-dir skills/spectacular

# 3. Fast verification gate entry point
bash test/verify.sh bench
```

### B. Live Multi-Harness Trials (Consumes API Credits)

To run behavioral A/B trials against live model endpoints (e.g. Claude Sonnet 5, OpenCode Muse Spark, Codex):

```bash
# 1. First, build the pinned candidate binary
go build -o /tmp/spectacular ./cmd/spectacular

# 2. Run focused smoke trials in parallel (e.g., 4 concurrent workers)
go run ./test/benchmarks/cmd/bench run \
  --spectacular-cli /tmp/spectacular \
  --candidate $(git rev-parse HEAD) \
  --model opencode/muse-spark-1.2-contributor-free \
  --adapter test/benchmarks/adapters/opencode-adapter.sh \
  --tier micro \
  --parallel 4 \
  --out test/benchmarks/reports/live-micro

# 3. Run Multi-Harness Matrix (Codex & Claude in parallel)
go run ./test/benchmarks/cmd/bench matrix \
  --spectacular-cli /tmp/spectacular \
  --candidate $(git rev-parse HEAD) \
  --models "codex:gpt-5.6-terra,claude:claude-opus-5" \
  --tier micro \
  --parallel 2

# 4. Run A/A Noise-Floor Calibration on identical code
go run ./test/benchmarks/cmd/bench calibrate \
  --model gpt-5.6-terra \
  --adapter test/benchmarks/adapters/codex-adapter.sh \
  --spectacular-cli /tmp/spectacular \
  --repeats 2

# 5. Inspect Historical Execution Trends & Economy Metrics
go run ./test/benchmarks/cmd/bench history

# Or trigger via verify.sh
LIVE_BENCHMARK=1 BENCH_TIER=micro BENCH_PARALLEL=4 bash test/verify.sh bench
```

---

## 4. Benchmark Tiers & Scopes

- **`micro`**: 3 essential cases $\times$ 1 repeat $\times$ 2 variants = 6 model calls (~60 seconds in parallel).
- **`smoke`**: 11 representative cases $\times$ 1 repeat $\times$ 2 variants = 22 model calls (~2 minutes in parallel with `--parallel 4`).
- **`full`**: All 26 catalog cases across full repetitions.
- **`held-out`**: Unseen generalization cases for anti-overfitting checks.

---

## 5. What the Benchmarks Measure

1. **Context Economy & Token Efficiency**: Measures input token savings achieved by compact single-file envelopes and scoped charters ($\ge 40\%$ reduction target), reporting `TokensPerSuccess`.
2. **Tool Economy & Waste Mitigation**: Measures tool calls per task success and penalizes `BranchPollution` (unnecessary mission files or lingering worktrees).
3. **Decision Fidelity & Orchestration Boundaries**: Asserts that the model complies with locked architectural decisions (`D<N>`), handles mid-run design gap escalations (`ORCH-02`), and rejects incomplete return receipts (`ORCH-03`).
4. **Scope Containment & Safety**: Detects whether the model escapes its authorized write perimeter or reads security canary files (`CAMPAIGN-CANARY`, `ARCHIVE-CANARY`).
5. **Case-Validity Quarantine**: Isolate fixture or adapter defects into quarantined partitions while evaluating conditional comparative effects on valid active cases.
6. **Harness Parity**: Normalizes tool events, token usage, and semantic observations across Claude Code, OpenCode, Codex, and Antigravity.
