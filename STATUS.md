# STATUS — Issue portfolio conclusion

**Updated:** 2026-08-07
**Repo:** `alexsmedile/spectacular`
**Verified against:** `main` @ `e6ec0f6` (`Clarify intent routing precedence and adversarial boundaries`)

## Current objective

Turn the remaining, current GitHub issue portfolio into independently dispatchable
work without duplicating solved work or reopening closed decisions.

## Decisions and completed issue actions

- Closed **#4** as superseded by #11; its useful editorial safeguards were added
  to [#11's discussion](https://github.com/alexsmedile/spectacular/issues/11#issuecomment-5181816035).
- Closed **#5** as consolidated into #29, the single owner for the AFK commit
  layer and read-only commit-plan proposal.
- Closed **#9** as mitigated: the intent receipt is implemented; validate it in
  ordinary future use rather than recovering the historical Unwire prompt.
- Closed **#10** as superseded by the GitHub CLI path. `gh pr create` files PRs,
  `gh pr review --approve` submits a reviewer approval, and `gh pr merge` merges
  when repository policy permits. GitHub still disallows self-approval by the PR
  author.
- Closed **#25** as already covered by AFK's durable goal, policy, and explicit
  apply/confirmation gates; #29 owns the missing commit behavior.
- Approved the capture-routing model:

  ```text
  capture → route to a deliberate destination
                   ├─ issue
                   ├─ shared
                   ├─ roadmap
                   ├─ vision
                   └─ request
  ```

  Local IDEAs should become a minimal compatibility/private-offline fallback,
  never the default source of truth. No automatic mirroring is allowed.

## Important constraints

- `shared` is assumed to be a provider-configured external destination, not a
  new committed Spectacular collection, unless the maintainer chooses otherwise.
- Remote issue creation must remain explicit and dry-run/confirmation gated.
- PR authors may merge through GitHub CLI when policy permits, but cannot approve
  their own PR. A separate reviewer identity is required where approval is gated.
- Do not implement #13 before #11/#12 have measured deterministic retrieval.
- Keep Bash 3.2 compatibility and preserve lifecycle, approval, policy, and
  verification safeguards.

## Remaining work, in order

1. Review and explicitly approve or correct draft
   `.spectacular/specs/SPC-007-capture-routing.md`; do not create implementation
   work until that approval.
2. Run the fresh #11 assessment/benchmark. #28 then selects the one compiled
   briefing command/schema before #11/#12 implementation.
3. Implement #17 and #32 on one branch: durable request/PR intent plus a
   plain-language PR opening line.
4. Implement #29: a read-only AFK commit-plan checkpoint. It must not stage,
   commit, amend, push, merge, reset, or stash.
5. Implement #30: a read-only portfolio issue review workflow that checks
   currency, overlap, order, blockers, and emits handoff prompts.
6. Grill #31 with #20 before designing a mission graph; #24 depends on their
   node/execution model.

## Other open issue disposition

- **#18, #19, #33, #42:** presentation/workflow UX; review after the core
  routing/retrieval work. #42's premise needs correction: approved Vision derives
  a draft SPC, not directly a PLAN.
- **#35, #36, #38, #39, #40, #41:** independent, lower-priority hygiene or
  convention work; dispatch only after their small design choices are confirmed.
- **#26:** belongs to the external issue-filing agent's format, unless that
  format is moved into this repository.
- **#13:** parked pending #11/#12 measurements.

## Verification and working tree

- GitHub open-issue state was refreshed on 2026-08-07.
- Closed Issue #6 and its expanded agent-signal-capture comment were read again
  on 2026-08-07. Draft `SPC-007-capture-routing.md` is the successor and has
  passed its document review; it is intentionally unapproved and unimplemented.
- `gh auth status` confirmed the active `alexsmedile` GitHub account with `repo`
  scope; no PR or code changes were created in this conclusion session.
- Working tree was clean at conclusion (`git status --short --branch` produced
  only `## main...origin/main`).

## Resume instruction

Start with item 1: review `SPC-007` against the approved model, correct it if
needed, then wait for explicit approval before creating implementation work.
