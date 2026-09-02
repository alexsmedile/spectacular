# Cross-Language Determinism Matrix

Flaky tests are caused by non-deterministic runtime assumptions. This matrix provides the exact anti-patterns and required deterministic replacements across Go, TypeScript/JavaScript, Python, and Rust.

---

## 1. No Wall-Clock Sleeps (The Zero-Sleep Principle)

Arbitrary sleep statements (`sleep(100ms)`) pass on fast developer laptops and fail under CPU contention on distributed CI runners.

| Language | ❌ Banned (Anti-Pattern) | ✅ Required (Deterministic Pattern) |
|---|---|---|
| **Go** | `time.Sleep(100 * time.Millisecond)` | Channel waits `<-done`, condition sync, or synthetic time via Go 1.24+ `testing/synctest` (`synctest.Run(func() { ... })`). |
| **TypeScript / Node** | `await new Promise(r => setTimeout(r, 100))` | Jest/Vitest fake timers (`jest.useFakeTimers()`, `vi.useFakeTimers()`, `vi.advanceTimersByTime(100)`), or event listener promises. |
| **Python** | `time.sleep(0.5)` | Explicit condition polling (`wait_for(condition)` with 1ms intervals), `asyncio.Event`, or `freezegun` / `time-machine`. |
| **Rust** | `std::thread::sleep(Duration::from_millis(50))` | Tokio mock time (`tokio::time::pause()`, `tokio::time::advance()`), or crossbeam/tokio channels `rx.recv().await`. |

---

## 2. Unordered Iteration (Maps, Sets, Dictionaries)

Many languages intentionally randomize hash seed order per process run. Asserting JSON or string outputs containing un-sorted map keys causes intermittent test failures.

| Language | ❌ Flaky Assumption | ✅ Deterministic Pattern |
|---|---|---|
| **Go** | Iterating directly over `map[K]V` into string or serialized output. | Extract keys, sort them (`sort.Strings(keys)` or `slices.Sort(keys)`), then iterate in sorted order. |
| **Python** | Iterating over `set` or relying on dict insertion order across serializations. | Sort collections before assertion: `sorted(my_set)` or `json.dumps(d, sort_keys=True)`. |
| **TypeScript / JS** | Assuming `Object.keys()` maintains insertion order across numeric vs string keys. | `Object.keys(obj).sort().map(...)` or sort array items before equality assertions. |
| **Rust** | Iterating over `std::collections::HashMap` or `HashSet`. | Use `BTreeMap` / `BTreeSet` when order matters, or collect into `Vec` and `.sort()` before assertion. |

---

## 3. Shared External State (Ports, Databases, Tempfiles)

The number one cause of failures when running parallel test runners (`-parallel`, Jest workers, pytest-xdist) is port collisions and dirty filesystems.

| Resource | ❌ Flaky Anti-Pattern | ✅ Deterministic Pattern |
|---|---|---|
| **Network Ports** | Hardcoded `:8080`, `:3000`, `:5432` | Bind to ephemeral port `:0` (e.g. `net.Listen("tcp", ":0")` in Go, `server.listen(0)` in Node), then read assigned port dynamically. |
| **Files & Dirs** | Writing to shared `/tmp/test.json` or local `./temp` | Use framework-provided unique isolated temp directories: `t.TempDir()` (Go), `tmp_path` fixture (Python), `fs.mkdtempSync` (Node), `tempfile::tempdir()` (Rust). Cleaned up automatically. |
| **Environment Variables** | `os.Setenv("FOO", "bar")` without restoration | Always clean up on exit: `t.Setenv("FOO", "bar")` (Go auto-restores), `monkeypatch.setenv()` (Python), or restore original env in `finally`/`defer`. |
| **Databases** | Running tests against one shared persistent database | Run tests inside transactions rolled back at teardown (`BEGIN ... ROLLBACK`), or spin up ephemeral Testcontainers with dynamic port mappings. |

---

## 4. Supply-Chain & Dependency Pinning

Unpinned dependencies produce tests that pass today and break tomorrow without any code changes in the repository.

- **Lockfile Enforcement in CI**:
  - Go: `go test` with check that `go.mod` and `go.sum` are tidy (`git diff --exit-code go.mod go.sum`).
  - Node: `npm ci` (never `npm install`).
  - Python: `pip install --require-hashes` or `poetry install --frozen`.
  - Rust: `cargo test --locked`.
- **GitHub Actions SHA Pinning Rules**:
  - Banned: Floating tags (`uses: actions/checkout@v4`).
  - Required: Full 40-character commit SHA with tag annotation:
    `uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2`
  - **Verification Invariant**: A pin is **never written** until verified via GitHub API:
    ```bash
    gh api repos/<owner>/<repo>/git/ref/tags/<tag>
    ```
- **Continuous Delivery & Keyless Publishing (Out of Scope for Test Tier)**:
  Pure test suites run with `permissions: contents: read`. If extending pipelines to publish packages or deploy:
  - Add `permissions: id-token: write` for keyless OIDC auth (e.g. `aws-actions/configure-aws-credentials`).
  - Use `actions/attest-build-provenance` or Sigstore/Cosign for SLSA cryptographic provenance attestations.
  - **Never** grant `contents: write` globally.

---

## 5. The Flaky Quarantine Protocol (No `retry: 3` & 14-Day Expiry)

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ ⚠️ THE RETRY BAN                                                            │
│ Banned in CI: actions/retry, retry: 3, pytest-rerunfailures                 │
│ Rationale: Retries hide non-determinism, mask memory leaks, and let         │
│ concurrency bugs escape to production.                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

When a test flakes in CI:
1. **Quarantine Immediately**: Move the test into an explicit quarantine suite (e.g. `@pytest.mark.quarantine`, `// +build quarantine`, or tagged group).
2. **Never Mask Mainline**: Mainline PR gates run non-quarantined tests and must remain 100% green without retries.
3. **Run Quarantined Asynchronously**:
   ```yaml
   quarantine:
     name: Quarantined Tests (Non-blocking)
     if: always()
     runs-on: ubuntu-latest
     continue-on-error: true
     steps:
       - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
       - run: go test -tags=quarantine -count=20 ./...
   ```
4. **Hard 14-Day Expiry**: Quarantine without a deadline is a graveyard. Any test quarantined for **> 14 days** must either be fixed with deterministic synchronization or permanently deleted.
