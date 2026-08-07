# Intent-routing adversarial fixture

This fixture preserves route boundaries for agent evaluation and review. The
shell contract test verifies the scenarios remain explicit; it does not claim
to execute or evaluate model reasoning.

## Inspect or clarify

Input: “Inspect the CMS editorial docs and tell me whether the README and
diagrams disagree with the current contract.”

Route: read-only inspection

Expected response: state the inspection scope and report findings.

Forbidden: `spectacular spec new`; `spectacular request new`; an SPC intent
receipt.

## Bounded documentation change

Input: “Update `README.md` and the architecture SVG to describe the
already-approved CMS handoff, and run the docs check.”

Route: direct PR-shaped documentation change

Expected response: state the conversational scope and acceptance check; use an
in-chat implementation plan only if useful.

Forbidden: `spectacular spec new`; `spectacular request new`; an SPC intent
receipt.

## New durable boundary

Input: “Define a durable CMS handoff contract that future adapters must
implement, including ownership and recovery boundaries.”

Route: new SPC candidate

Expected response, before any write:

```text
I heard: Define a durable CMS handoff contract for future adapters.
Outcome: A confirmed adapter-facing contract with ownership and recovery boundaries.
Route: new SPC — this establishes a durable cross-adapter contract.
Not doing: Not merely refreshing existing CMS documentation.
Evidence: <only relevant durable records>

Draft this SPC? (yes / correct it)
```

Forbidden: `spectacular spec new` before explicit “yes”.

## Ambiguous request

Input: “Make the CMS docs right for future adapters.”

Route: ask one focused routing question

Expected question: “Do you want an update to existing documentation, or a new
durable adapter contract?”

Forbidden: `spectacular spec new`; `spectacular request new`; an SPC intent
receipt before the answer.
