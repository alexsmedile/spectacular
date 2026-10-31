# Test credibility

Use this when: reviewing tests directly or assessing whether tests support a code finding.

For each disputed test, identify the claim, production entrypoint, real code executed,
substituted boundaries, oracle and assertion. Ask: which plausible broken behavior would
still make this test green? Demonstrate that counterexample when feasible.

| Signal | Discriminating check |
|---|---|
| Expected value derived from the implementation | Use an independent requirement, hand-calculated fixture or trusted reference; shared helpers can reproduce the same bug. |
| Production function replaced by a mock | Check whether the test asserts only the configured return; keep domain logic real and mock actual I/O boundaries. |
| Broad mocks | Introduce adapter shape/status errors; verify contract and serialization with a real local boundary fixture. |
| Existence / substring / snapshot assertions | State the property they really prove; a file may exist while its contents or behavior are wrong. These assertions are valid for artifact/discovery contracts. |
| Catch-all / conditional assertions | Verify the intended assertion executes and unexpected errors fail; distinguish deliberate optional-dependency skips. |
| Happy-path-only fixtures | Add the smallest negative case exposing a demonstrated bug, including unchanged-state assertions on rejection. |
| Coverage without sensitivity | In a disposable copy, introduce a plausible defect and run the same focused test; a surviving mutation exposes only that test's blind spot. |
| E2E label with substituted pipeline | Trace which boundaries actually execute; label the proven integration scope precisely. |
| Green suite from partial discovery | Check runner patterns, languages, skipped tests and dependency guards; report executed versus omitted suites. |

Run the unmodified focused test first, then the counterexample or isolated mutation, then
confirm the trial did not alter the checkout. Mutation is optional when static reasoning
already demonstrates the gap; never edit production code merely to test test quality.
A failing trial needs diagnosis: distinguish product regression, environment failure and
incorrect test assumptions. Do not weaken assertions to obtain green output.

Suggested behavioral trials: a registered handler returning no artifact; an invalid explicit
config silently treated as absent; ASR over B-roll labeled visible speech; a stale cache
with a reused ID; retry after a simulated committed side effect. Choose only relevant ones.

A test that accurately verifies an incorrect domain assumption is a wrong oracle, not a
no-op test. Passing parse/type/schema checks cannot certify rendered output, accessibility,
audio/visual quality or owner acceptance. Preserve useful low-level tests and identify the
specific higher-level evidence missing.
