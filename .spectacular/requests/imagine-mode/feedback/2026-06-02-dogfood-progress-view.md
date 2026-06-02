---
type: feedback
target: dogfood-progress-view
scope: skill
status: resolved
opened: 2026-06-02
resolved: 2026-06-02
proposal_summary: "M6 dogfood — ran `spectacular imagine` end-to-end on the progress-view feature idea; surfaced + fixed 2 real bugs"
next_action: ship-as-is
request: imagine-mode
spawned_request: null
promoted_to: null
related: []
---

# Feedback — dogfood-progress-view

## Target
The `imagine` mode itself (M6 dogfood), exercised by imagining the `progress-view` feature — a workspace-wide `spectacular progress` (no slug) dashboard. Vision built: spine (end-goal/phases/flow) + 4 fragments (2 ui, 1 story, 1 arch), all approved by the human.

## Hypothesis / hunch
The engine doc (`imagine.md`) fully specifies the render→react→derive loop; running it for real should be smooth and mostly confirm the spec.

## Proposal
Ran the full loop live: scaffolded `requests/progress-view/vision/`, rendered ASCII artifacts, presented for per-fragment reaction. Human approved all 4 fragments, chose to fix surfaced bugs, and held PLAN derivation to refine first.

## Question asked
"Which fragments do you approve? What to do with the surfaced bugs? Derive the PLAN now?"

## User response
Approved all 4 fragments. **Fix the caption bug in imagine-mode.** Hold derivation — refine fragments first.

## Insight
**The dogfood found two real bugs the spec missed — which is the entire point of M6:**

1. **`--caption` rejects values starting with `-`.** The shared `require_value` guard (`val != -*`) protects against forgotten flag values but wrongly blocks free-text captions like `"--stalled triage view"`. *Fix:* added the attached `--caption=<text>` form (bypasses the guard cleanly; mirrors existing `--since=`). Engine doc updated to use it.

2. **Manifest drift only detected fragment presence, not approval-state.** Approving a fragment (`approved: pending → true`) left the spine manifest stale and `doctor vision` reported clean — defeating the manifest's whole purpose (mirroring approval). *Fix:* drift check now compares the live manifest body against a freshly-computed one (catches presence, approval, caption changes). Factored out `_vision_manifest_lines` shared by writer + checker. Engine doc now says: after approving, run `doctor vision --fix`.

Both are exactly the kind of gap that only shows up when a human actually drives the loop — not visible from reading the spec. Verified: suite 9/9, doctor clean, manifest now shows 4× approved.

## Decision
Resolved within M6 — both fixes shipped in the same milestone. The dogfood validated M3/M4/M5's loop spec end-to-end *and* hardened it. The imagined `progress-view` feature was **folded into [[visual-layer]]** (v1.15.0) — it's a slice of the already-roadmapped Visual layer (`spectacular progress` rendering), not a separate request. Its `vision/` artifacts moved to `requests/visual-layer/vision/` as render-spec input; the standalone request was retired.
