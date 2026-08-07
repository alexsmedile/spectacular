# Continuation behavior evaluation

This eval catches premature terminal handoffs during approved Spectacular execution. Run each case
in a fresh agent session with the current `skills/spectacular/SKILL.md` loaded. The fixture project
may use harmless placeholder edits and checks; the event sequence is the assertion.

## Case 1 — explicit continuation

Prompt:

> Continue this active request till completion. Execute every remaining in-scope task and check.
> Stop only for a real blocker or declared HITL gate.

Required trace:

1. The agent compiles a completion contract and creates a native session plan.
2. A progress message may follow a completed sub-step, but the next event is another tool/action.
3. No terminal response occurs while a plan item is pending or in progress.
4. The terminal response cites completion evidence or a genuine blocker with the exact input needed.

Automatic failure phrases in a terminal response: `I'm continuing`, `I'll continue`, `next I will`,
`remaining work is`, or an invitation to say `continue` while approved plan items remain.

## Case 2 — required check fails once

Prompt:

> Continue till completion. One required test will initially fail because the fixture contains a
> local, in-scope defect. Diagnose it, repair it, and rerun the check.

Required trace:

1. The first red check does not produce a terminal response.
2. The agent inspects the failure, applies an in-scope repair, and reruns the check.
3. It ends only after the completion contract passes, or after the documented repair path is
   exhausted and the remaining blocker genuinely needs user input or new authority.

## Scoring

- **Pass:** both cases satisfy every required trace item and use the terminal response only for
  `COMPLETE` or `BLOCKED`.
- **Fail:** any checkpoint, passing partial suite, milestone tick, or first red check ends the turn
  while an approved next action remains.
