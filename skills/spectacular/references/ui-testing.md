# UI testing and design refinement

Use this when: auditing visual quality, responsive behavior, action hierarchy,
progressive disclosure, or redesigning an existing interface and its user journey.
Use the owner's scope and authorization; this route introduces no CLI, Mission,
mandatory new document family, or extra approval gate.

## Evidence before critique, hierarchy before polish

Follow this order: inventory → screenshots/material → review and critique →
principles and proposed solution → low-fidelity prototype → implementation →
new screenshots and visual review. Do not substitute a list of suggestions for
requested evidence, or start aesthetic edits before understanding the information
architecture. A reversible implementation already authorized by the owner does
not require another approval merely because a prototype was produced.

Read relevant anchors, constraints, existing UI and available test tooling.
Select screens by frequent tasks, value and failure risk. Include principal
pages and dialogs, drawers, menus and popovers. Honor requested counts; if fewer
real pages exist, name distinct states or sections rather than inventing pages.
Record desktop, tablet and mobile dimensions in CSS pixels; add breakpoint edges,
narrow width, short height, zoom, dark mode and virtual keyboard where material.

Use representative synthetic data: empty and populated lists, long names and
filenames, large amounts, many items, errors, loading, success and disabled states.
Use a disposable or explicitly authorized test environment. Do not edit real
databases directly, expose credentials, publish screenshots, or send external
messages merely to produce UI evidence. Record the runtime/build exercised.

Save screenshots with screen, state, viewport and attempt identifiers. Preserve
baseline and later attempts separately. Verify the destination and visible state
before capture; finish known animations and loading. A screenshot of a closed
dialog or wrong page is missing evidence, even when its filename is correct.
Capture viewport composition and relevant scrolled/footer states; full-page
images alone cannot prove actions are reachable in a normal viewport.

## Review both pixels and the task

Inspect images visually, using contact sheets for comparison and full-resolution
views for findings. DOM geometry, accessibility checks and automated screenshot
diffs supplement visual judgment; they cannot certify elegance or usability.
Name the evidence and distinguish observed defects, source-derived risks,
design judgments and untested hypotheses. Failed and blocked attempts stay visible.

Challenge these principles explicitly:

| Principle | Question / proposed remedy |
|---|---|
| Action hierarchy | Is one contextual primary clear? Reduce equal-weight buttons |
| Progressive disclosure | Are optional/rare controls revealed when needed? |
| Discoverability | Does every hidden capability retain a named touch/keyboard route? |
| Task-first organization | Are frequent user goals ahead of internal technical structure? |
| Read versus edit | Can a person inspect information without entering a form? |
| Feedback and trust | Is success local, truthful and proportionate, without duplicate toasts? |
| Error prevention | Are actions appropriate to lifecycle state and consequences? |
| Continuity | Does opening a related form preserve the work in progress? |
| Responsive meaning | Do key content, amounts, status and actions survive narrowing? |
| Accessibility | Are focus, labels, keyboard and touch alternatives usable? |
| Visual rhythm | Do typography, spacing, surfaces and semantic color establish hierarchy? |
| Language | Are labels concise, localized and consistent with the user's vocabulary? |

Do not confuse neutral design with making every control identical. Do not hide
critical errors, required inputs or task completion behind vague menus. Advanced
settings, organization metadata and diagnostics should not dominate routine work.
Inspirations inform qualities and trade-offs; they do not justify copying a brand
or asserting current product behavior without inspection.

## Prototype the organization at low fidelity

Use ASCII, wireframes or a disposable interactive fragment before visual polish.
Lock known constraints, identify the open decision, show the normal path and
recovery path, and compare concrete alternatives only where uncertainty matters.
For a real multi-option exploration use the rapid-prototyping companion skill;
for a bounded authorized redesign, choose a justified direction within scope.
Show primary/secondary actions, disclosure entry points, mobile adaptation and
completion feedback. Explain how the proposal addresses each priority finding.
Do not introduce new approval requirements; respect owner-requested review gates.

## Implement and verify the user outcome

Fix shared foundations first (dialogs, overflow, focus, feedback), then hierarchy
and flow. Preserve working capabilities and unrelated dirty state. Reuse existing
components where suitable; avoid blanket CSS that erases state semantics.
Exercise meaningful journeys through visible controls, and recapture equivalent
screens/states at the same dimensions. Inspect after images, not just test output.
Measure changed action counts, legibility and completion effort when useful;
do not invent usability scores or claim user validation from a mockup.

The existing plan owns scope, selected direction, prioritized changes and progress.
Use reusable Scenarios for independent journeys and Audits for attempt findings
as described in [scenarios.md](scenarios.md). Keep sketches in the existing raw role.
Keep screenshot locations, capture recipe and evidence status linked from that owner.
Retain a reusable capture/check recipe when ongoing UI work warrants it, using the
host's supported browser mechanism. A test named “UI integration” that only checks
HTTP responses must remain classified as HTTP evidence.
Report implemented, visually verified, failed, blocked and unverified separately.
Completion requires the requested coverage and no unresolved blocking findings;
release and product acceptance retain their own authorization and proof requirements.
