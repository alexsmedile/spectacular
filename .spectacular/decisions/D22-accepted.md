---
type: Decision
id: 01a02bb5-c364-74ba-82f4-0c76f7496605
title: Run transition state machine and charter tokenizer specification
created_by: Alex
created: "2026-08-22T23:02:07Z"
updated: "2026-08-22T23:02:07Z"
actor: Alex
actor_role: ""
disposition: accepted
question: ""
rationale: |
    Formalizes the run transition state matrix and token budget threshold evaluation:
    1. Run State Machine: A Run starts active. Active runs may move to paused, blocked, awaiting-review, completed, or stopped. Paused runs may move to active, blocked, or stopped. Blocked runs may move to active or stopped. Awaiting-review runs may move to active, completed, or stopped. Completed and stopped are terminal states. Every transition requires explicit actor (--by) and reason (--reason) attribution.
    2. Charter Tokenizer: Uses the spectacular-charter-tokenizer.v1 specification with o200k_base vocabulary rules and frozen numeric gates: 1,200 max target (pass), 1,400 warning envelope (warn with safe compaction), 1,440 ceiling (split recommended), and >1,440 hard refusal.
ref: D22-accepted
---
# Context

This decision formally binds the lifecycle state graph for execution runs and the prompt budget accounting rules for Context Sandwich compiler envelopes.
