# Active-Loop Transitions

Use this when: an active matrix loop receives A, B, C, M, R, S, or F; receives G only after the
singular merged-result gate; receives an explicit
plain-language request to combine named proposal parts; or receives a plain-language request to grow
the uniquely selected current result. Treat a plain-language combine request as M.

Apply the decision ledger invariant in `SKILL.md` before taking a transition.

For named proposals created outside this loop, inspect each proposal or mark it unavailable with the
attempted access method and exact missing prerequisite. If any required proposal is unavailable,
stop with the unchanged ledger. Otherwise set `Locked` to their shared dimensions, `Open` to the
merge axis, `Level` to their lowest common fidelity, and `Lineage: root`; derive two or three
observable success signals from the merge request and shared constraints; then apply M.

For any plain-language grow request, use the uniquely selected current result. If multiple options
remain and none is named, request A, B, or C and preserve the ledger. If the current result is
singular, use the G/plain-language-grow row without appending an option key.

For M, take the rematrix branch only when the named merge leaves one specific observable design
axis unresolved; record that axis as `Open`. Otherwise clear `Open` and take the singular-result
branch. If the user omitted one result-changing merge part, request only that part before choosing.

| Key | Current level | Transition | Required next step and state |
|---|---:|---|---|
| A / B / C + explicit higher target level | 1–4 | Lock the selected move and set the requested level | For levels 2–4, return to step 1 at that level and run steps 2–4; for level 5, use the next row |
| A / B / C | 1–3 | Lock the selected move and advance one level | Return to step 1 at the new level, then run steps 2–4; reuse A/B/C as input keys and show lineage such as `root → B → A` |
| A / B / C | 4 | Lock the selected screen or flow, set `Open: cleared`, and enter level 5 | Return to step 1 at level 5 and satisfy its full completion criterion; then integrate only that lineage, run step 3, and deliver one selected integrated result with the ledger, changed artifacts, evidence statuses, and uncovered risk |
| M | 1–4, one specific observable axis remains | Merge only explicitly named parts and record that axis as `Open` | Remain at the current level, return to step 1, then run steps 2–4 |
| M | 1–4, no consequential decision remains | Make the named merge `Locked`, set `Open: cleared`, preserve `Lineage`, and produce one result | Satisfy the terminal-M completion criterion below |
| G or plain-language grow | 1–3, any singular current result | Preserve lineage and advance one level | Return to step 1 at the new level, then run steps 2–4 |
| G or plain-language grow | 4, any singular current result | Preserve lineage, set `Open: cleared`, and enter level 5 | Satisfy step 1 at level 5, integrate only the singular current result, run step 3, then deliver under the loop-completion criterion |
| R | 1–4 | Replace the open axis using the user's line | Remain at the current level, return to step 1, then run steps 2–4 |
| S | 2–4 | Move down one level | Return to step 1 at the smaller level, then run steps 2–4 |
| S | 1 | Set `Open` to a named strict sub-decision of the previous axis and state which parent portion is excluded | Remain at level 1, return to step 1, then run steps 2–4 |
| F | 1–4, multiple current options and no named option | Preserve lineage and the current state | Request `F: A`, `F: B`, or `F: C` and end the turn |
| F | 1–4, text-only singular/named output with no artifact, link, externally verifiable claim, or executable claim | If F names A, B, or C, append that key to `Lineage`; if the current result is already singular, preserve `Lineage` | Deliver the singular or named selection alone |
| F | 1–4, every other singular/named output | If F names A, B, or C, append that key to `Lineage`; if the current result is already singular, preserve `Lineage` | Run step 3, then report the ledger, changed files when any, evidence statuses, and uncovered risk |

When M needs clarification, complete the turn by requesting the single result-changing missing part
and leaving the current ledger unchanged.

For terminal M, complete only when one singular result contains the named merge parts, preserves
every locked dimension and `Lineage`, sets `Open: cleared`, exposes the current level's typical
evidence and every success signal, uses a disposable artifact whenever spatial, state, viewport, or
interaction behavior must be demonstrated, reports every applicable step-3 check, and makes these
exact lines the final content of the turn:

```text
[G] Grow this merged result one fidelity level
[F] Finalize this merged result at current scope
```

For A/B/C targeting levels 1–4 and for M/R/S paths that rematrix, satisfy step 4's completion
criterion. Conversational F is complete when only the singular/named selection is delivered.
For A/B/C at level 4 and artifact-backed F, use the loop-completion criterion below.

Complete the loop when the requested artifact exists, relevant checks have defensible statuses,
the decision ledger invariant holds, and every choice required for delivery is recorded as locked.
