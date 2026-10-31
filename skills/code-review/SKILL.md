---
name: code-review
description: >-
  Review code changes or audit a codebase for bugs, dead code, disconnected features,
  duplication and tests that give false confidence, including vibe-coded projects.
  Use for evidence-based code review and test credibility audits. Not for routine
  implementation, formatting-only cleanup, Git operations or visual UI critique.
metadata:
  version: "1.0.0"
---

# Code Review

Produce prioritized, reproducible findings about actual behavior. Portable Agent Skills
profile; implicit invocation is enabled by the host default for relevant review requests.
Review is read-only by default; remediation requires an implementation request.

## Choose scope

| Request | Route and required context |
|---|---|
| Review a diff | Read changed code, its callers, contracts and affected tests; compare against the named base. If no base is available, state that limit and inspect current behavior. |
| Study a codebase / vibe coding audit | Map entrypoints and ownership, then read [failure-patterns.md](references/failure-patterns.md); use its probes across the selected subsystems. |
| Review test credibility | Read [test-credibility.md](references/test-credibility.md); trace claimed behavior to assertions and production code. |

For diff or whole-repo review, also read test-credibility when a finding depends on test
adequacy, mocks, coverage or green-suite claims. Use existing repository tools; missing
analyzers are limitations, not a reason to install tools or download dependencies.
Spectacular workspace initialization, CLI, Mission and formal Decision are unnecessary.
If a workspace already exists, consult its relevant current contracts; retain its owners.

## Execute

1. **Establish the baseline.** Read applicable instructions, status/diff and entrypoints.
   Record the inspected revision plus dirty/untracked scope; distinguish user changes from
   this review. Identify intended behavior from contracts and actual consumers, flagging
   contradictory authorities. Inventory subsystem boundaries and generated/vendor/legacy
   surfaces. Done when the selected scope, acceptance claims and excluded surfaces are explicit.
2. **Trace execution.** Follow representative user actions from CLI/UI/API through validators,
   domain logic and persistence to observable output. For each selected subsystem inspect
   its main path plus a failure/retry or boundary path. Check registrations and adapters,
   not merely implementations. For whole-repo audits report which subsystems were deeply
   traced versus only inventoried. Done when each selected path has concrete callers,
   writers and consumers, or an explicit unresolved integration gap.
3. **Challenge candidates.** Read surrounding code and search all relevant call surfaces.
   For suspected dead code check exports, registries, config strings, dynamic imports,
   plugins, subprocesses, scripts and supported external consumers. Distinguish locally
   unused computation, internally unreachable code and public compatibility surfaces.
   For duplication compare semantics, validation and ownership before suggesting sharing;
   generated copies and intentionally independent boundaries may be legitimate.
   Done when each candidate has supporting evidence and plausible counterevidence checked.
4. **Verify consequences.** Prefer a minimal isolated reproduction or focused existing test.
   Before executing tests inspect side effects and external dependencies; keep trials in
   temporary fixtures and avoid production media, network calls and publishing. Use the
   test reference for credibility checks. If tooling is unavailable, retain a static finding
   only when the control/data flow proves it; otherwise mark it unverified.
   Done when every reported defect has a trigger, actual result, expected contract and
   consequence, with evidence type identified. A green suite proves only covered assertions.
5. **Deliver and stop.** Rank by impact and reachability: P0 immediate critical harm,
   P1 serious reachable failure, P2 bounded correctness/reliability gap, P3 maintenance cost.
   Keep confidence separate: reproduced, statically demonstrated, or suspected.
   Deduplicate findings with one root cause; attach related symptoms and test gaps.
   Done when every finding satisfies the output contract, checks and blind spots are
   reported, and no unrequested fixes, deletion, installation or publication occurred.

## Output contract

Lead with consequential findings. For each include ID, priority, confidence, exact file/line,
trigger and caller, expected versus actual behavior, impact, evidence/reproduction,
smallest suitable correction and a behavioral regression check. Evidence can be a cited
control-flow proof when execution is unsuitable; identify what remains untested.
Put unverified candidates and clarification needs separately; do not count them as defects.
Report commands, results, failures/skips and coverage limits. If no defects were established,
say so without claiming the codebase is clean. Smell counts are not a quality score.

Never manufacture findings to meet a quota. Tests proving schema, discovery,
artifact integrity or help behavior are valid at those layers; name the missing claim.
When saving was requested, use the project's existing audit owner and link the report.
