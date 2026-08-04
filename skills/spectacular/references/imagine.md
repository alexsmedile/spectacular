---
description: Generative pre-request direction discovery — Understand, Imagine, Probe, React, Confirm, then derive a draft SPC from approved Vision.
when_to_use: spectacular imagine <slug>, or when material product, interaction, workflow, or system-shape uncertainty is easier to resolve through concrete proposals than prose requirements.
---

# Imagine — the generative Vision engine

`imagine` is an operation; Vision is its durable result. The human supplies
aspiration, taste, reaction, and approval. The agent grounds the opportunity,
generates concrete alternatives, and makes them experienceable. “Dream” is useful
natural language for the human aspiration, not a second command or lifecycle.

## When to suggest it

Suggest Imagine only when interaction, UX, workflow, product direction, or system
shape is materially uncertain and human reaction would reduce specification risk.
Do not suggest it merely because work has a UI. Skip it for clearly specified
backend, maintenance, mechanical migration, and direct work.

## Preconditions

1. A Spectacular workspace exists. A request, PLAN, and SPC are not required.
2. Scaffold `.spectacular/visions/<slug>/` with `spectacular imagine <slug>` if it
   does not exist.
3. Load project intent plus only the code, behavior, personas, research, or spike
   evidence needed to understand the opportunity.

## The loop: Understand → Imagine → Probe → React → Confirm → Derive

### 1. Understand

Ground `VISION.md § Understanding` before proposing solutions:

- current reality — what exists now, observed rather than assumed;
- users and needs — whose outcome or experience changes;
- constraints — product, technical, policy, and operational boundaries;
- material uncertainties — questions whose answers could change the direction.

This is problem/direction understanding. It does not replace the later
code-grounded PLAN Understanding gate (how the touched system works / what changes /
what stays the same).

### 2. Imagine

Lead with proposals, not an empty questionnaire. Draft the Intent, North star,
Experience signature, meaningful strategies, and a recommended Chosen direction.
Create only the fragments needed to make the uncertainty discussable:

| Kind | Useful when |
|---|---|
| `strategy` | Alternatives and trade-offs need comparison |
| `story` | A user outcome or scenario needs sharpening |
| `flow` | Sequence, states, or handoffs are uncertain |
| `ui` | Screen, CLI output, or interaction needs a vibe check |
| `arch` | System shape or boundary needs comparison |
| `prototype` | A runnable/showable artifact will produce better reaction than prose |

One UI fragment may be enough. A non-visual workflow may need only a story and
flow. There is no one-of-each rule.

### 3. Probe

Use the cheapest sufficient evidence route:

- inspect current code or ask the human directly;
- research a bounded fact in RES;
- run a human-authorized SPK when technical feasibility must be observed;
- build a prototype fragment when human interaction is the evidence.

Research and spikes keep their own truth contracts. Link them from Vision; do not
copy their bodies. A spike can establish that an approach works without choosing
it. Present a showable conclusion, then let the human choose the direction.

### 4. React on parts

Present fragments individually and ask the human to approve, revise, reject, or
supersede them. Record reactions mechanically:

```text
spectacular vision react <slug> <fragment> --approve [--note <why>]
spectacular vision react <slug> <fragment> --revise  [--note <change>]
spectacular vision react <slug> <fragment> --reject  [--note <why>]
```

Revise only the redirected fragment. Keep rejected/superseded fragments as decision
history. Only `reaction: approved` fragments are load-bearing downstream.

### 5. Confirm the whole Vision

Fragments are parts; approval is a coherent whole-direction decision.

1. Resolve every fragment reaction and ensure at least one is approved.
2. Make Chosen direction, Boundaries, evidence, and rationale explicit.
3. Run `spectacular vision propose <slug>`.
4. Present the whole Vision as a short vibe-check card: North star, Experience
   signature, Chosen direction, approved/rejected fragments, material conditions.
5. Only after explicit human approval run:
   `spectacular vision approve <slug> --approved-by <human> [--note <text>]`.

Approval means “derive the contract from this direction.” It does not mean “write
production code.” Rejection records actor and reason; it is not deletion.

### 6. Derive a draft specification

On approval, offer `/spectacular vision derive <slug>`. Create or update one draft
SPC using:

- Intent + North star → specification intent;
- Chosen direction + approved fragments → requirements and scenarios;
- Boundaries → non-goals/constraints;
- linked RES/SPK evidence → evidence and decisions;
- material uncertainties → open questions, never guessed requirements.

Record the Vision path in the SPC's `related:`/evidence and set `derived_to` on
VISION.md after the draft exists. The SPC remains draft until separately approved.
Only the approved SPC may seed a request and PLAN.

## Boundaries

- Pre-request and opt-in. No request or PLAN precondition.
- Vision approval is human-only; agents may propose but never infer approval.
- Prototypes are artifacts, not an entity or mandatory phase.
- `experiment` is ambiguous natural language, not a feedback alias; route by the
  question being answered.
- Legacy `requests/<slug>/vision/` folders remain readable/diagnosable but do not
  gain the new whole-Vision lifecycle automatically.

**Related:** [[vision-rules]], [[discovery-protocol]], [[spike-rules]],
[[spec-lifecycle]], [[plan-rules]], and [[feedback-loop]].
