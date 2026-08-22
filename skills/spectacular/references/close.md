# Review and complete

Use this when: Primary operator or reviewer assessing claims, collecting Evidence, coordinating review, or completing a Mission.

## Assess each claim

Check every completion claim against its frozen pass boundary using the FROST
framework (Frozen fit, Risk, Operability, System integrity, Truth of proof)
defined in [audit.md](audit.md). For each one:

- does the Evidence meet the stated proof requirement
- is that Evidence attributable to a method, not just asserted
- is there contrary evidence
- is it fresh enough to still be true
- has the required review happened

When something fails, repair with the narrowest action justified by a new
hypothesis, new Evidence, or a materially smaller fix. Every repair consumes the
repair budget.

## Persist Evidence when it is earned

Write an Evidence file when the claim:

- enters assessment as a material claim
- depends on an external or provider observation
- is disputed
- must survive closure

Routine local progress does not need an Evidence file for every passing check.

## Create a review record only when required

Earned only when the Mission requires clustered or independent review.

```bash
spectacular review record <mission-ref> <review.md|-> --json
```

The review input is a Markdown file with `ReviewDraft` frontmatter (get a starter skeleton via `spectacular review record --help`):

```yaml
---
type: ReviewDraft
title: Independent review of M<n>
status: passed
reviewed:
    commit: <40-char-git-commit>
    # tree: optional (auto-derived from commit if omitted)
    activation_fingerprint: <sha256-mission-activation-fingerprint>
reviewer:
    actor: <identity-different-from-operator>
    operator: <operator-name>
    relation_to_operator: independent
    implemented_reviewed_scope: false
    independence_basis: <attributable statement of independence>
    evidence:
        - <attributable command or observation>
claims:
    - claim: <exact-claim-name-from-mission>
      verdict: pass
findings: []
limitations: []
---

# Review body
```

**Independent review needs a reviewer who did not implement the reviewed scope.**
A fresh agent is not automatically an independent one. Freshness is not
independence.

What earns it: the reviewer did not author the work, inspects primary evidence
rather than the executor's summary of it, and where the consequence is high uses
a different method, a qualified human, a specialist tool, or an independent
observation. The reviewer reports verdicts and findings; it never declares the
Mission complete.

### Dual-Path Independent Review Workflow

When an independent review is required, offer the user two execution paths:

#### Path A: In-Harness Subagent Dispatch
If the host runtime supports subagents (e.g. Antigravity `invoke_subagent` or Claude Code subagents):
1. Spawn a dedicated child subagent in a pristine context using the `strict-verifier` or `reasoning` profile.
2. Provide the subagent with the exact reviewed Git commit SHA, tree SHA, frozen completion claims, and the FROST framework instructions ([audit.md](audit.md)).
3. The subagent inspects the Git diff and primary evidence directly, evaluates each claim, and writes the `ReviewDraft` to a temporary file outside `.spectacular/`.
4. Record the review via `spectacular review record <mission-ref> <temporary-review-file> --json`; only the CLI-generated Review belongs in `reviews/`.

#### Path B: External Model / Human Handoff (Clipboard & File Prompt)
When running in a single-agent harness or utilizing a distinct external reasoning model (e.g. OpenAI o3, DeepSeek-R1, an external Claude session, or a human peer):
1. The Skill generates a self-contained, copy-pasteable review prompt in `scratch/` or a temporary directory outside `.spectacular/` (and prints it in chat).
2. The prompt includes:
   - Git baseline, reviewed commit SHA, and tree SHA
   - Exact frozen completion claims with `pass_boundary` and `proof_requirement`
   - FROST inspection protocol (Frozen fit, Risk, Operability, System integrity, Truth of proof)
   - The required `ReviewDraft` YAML schema
3. The external reviewer inspects the work and returns the structured `ReviewDraft`.
4. Pipe the result directly into the CLI:
   ```bash
   spectacular review record <mission-ref> - --json
   ```

## Complete in one flow

1. Run focused checks, then one full tree-bound gate (e.g. the host project's test suite like `npm test`/`pytest`, a clean build/run verification if no test suite exists, or in this repo `bash test/verify.sh all`).
2. Verify every frozen claim, and every required Evidence or review:
   `spectacular mission check <ref> --json`.
3. Confirm the applicable product and specification edits are in the same worktree.
4. Present one owner gate, and stop. Do not proceed on your own reading of intent.
5. After attributable owner confirmation, close it:

   ```bash
   spectacular objective finish <ref>/<objective> --json
   spectacular mission complete <ref> --by <owner> --json
   ```

   `--by` is what makes the confirmation attributable — it names the person who
   accepted. Never pass a name the owner did not actually give you.

There is no Contract reconciliation ritual at completion. Completion checks the
Gaps the Mission declared it would close and refuses while any is still open; it
never writes the Contract itself.

## Close a declared Gap

A Mission that declared `resolves_gaps:` closes each one when the work resolving
it lands, not at completion:

```bash
spectacular contract amend <contract-ref> --gap <gap-ref> --by <owner> --dry-run
spectacular contract amend <contract-ref> --gap <gap-ref> --by <owner>
```

`--dry-run` prints the resolution text, both fingerprints, and every Mission that
would be re-pointed, and writes nothing. The text is the wording frozen in the
Mission's declaration, so the owner approved it at activation.

The amendment reaches the `gaps:` block and editorial fields only. A semantic
change is a new `contract_version:`, not an amendment.

It refuses while a bound Mission that did **not** declare this Gap is live —
that Mission still has the Contract constraining work in flight. The Mission that
declared the Gap is the exception, and closes it while live; an owner
`--resolution` override is never exempt, because its wording was typed at a
prompt rather than approved at an activation gate.

Only the live Mission is re-pointed to the new fingerprint. A completed Mission
keeps the binding it agreed to: `mission check` reports the difference as a
contract-drift notice, and `git log -S <fingerprint>` recovers the exact Contract
text that was in force. Rewriting it would replace a historical fact with today's
answer, so a completed Mission is never re-pointed to silence the notice
(`D10-repoint`).

A Gap is never closed by deleting it.

## Refuse completion

Refuse when any of these is true:

- the frozen criteria changed
- proof is stale, missing, or conflicting
- required independent review is incomplete
- a dependency or Gap still blocks, or a declared resolved Gap is still open
- the repair budget is exhausted
- the baseline drifted
- owner confirmation is missing

Completion and archival are different things. Archival proves nothing.
