# Template — implementation Mission handoff

Do not use this template until S12 has produced an approved specification and Mission boundary.

```text
You are the implementation-Mission orchestrator for <MISSION-ID / slug>.

Approved outcome:
<one independently reviewable contract delta>

Authority and baseline:
- Approved specification: <stable ref and approval commit>
- Exact base commit: <SHA>
- Integration target: <branch>
- Traffic result: parallel | conditional | serialized
- Branch/worktree creation: allowed | forbidden
- Required branch: codex/feat/v2-<mission-slug>
- File mutation: allowed only in <paths>
- Local commits: allowed | forbidden
- Push: allowed | forbidden
- Draft PR: allowed | forbidden
- PR ready, merge, deployment, destructive cleanup, remote deletion: forbidden unless separately
  authorized

Goal:
<observable result>

Constraints and non-goals:
- <accepted product/architecture constraints>
- <explicit exclusions>

Shared interface/join:
<exact schema, signature, registry key, ordering, or none>

Milestones:
1. <closed milestone and expected output>
2. <closed milestone and expected output>

Required evidence:
- <check and expected result>
- <compatibility/migration/rollback evidence>

Agent-routing rules:
- Inspect unfamiliar code read-only before planning.
- You are the only lifecycle/checkpoint mutator.
- Build inline unless there are at least three independent, closed, disjoint-file units.
- If dispatching, give each builder Goal / Constraints / Approach / Expected output / Success
  criteria plus the identical shared join.
- Parallel builders may touch only disjoint files; serialize shared CLI, schemas, registries,
  canonical contracts, and tests.
- Use independent review when risk requires it; reviewers do not implement fixes.

Stop conditions:
- Missing or contradictory approved contract
- Work outside scope or paths
- New public interface, dependency, migration, security, privacy, or product decision
- Changed/conflicting Git baseline
- Unresolved required check after bounded diagnosis
- Need for a forbidden provider or destructive action

Completion contract:
1. Required implementation complete
2. All named evidence passes against final head
3. No undeclared scope change
4. Durable Mission checkpoints/evidence updated by the orchestrator
5. Coherent commits created if authorized
6. Branch pushed and draft PR opened if authorized
7. No merge, deployment, remote deletion, or unauthorized cleanup

Return:
- status: complete | blocked | failed
- baseline and final head
- files changed
- tests/evidence with exact results
- decisions made: none unless explicitly owner-authorized
- assumptions and remaining limitations
- scope deviations
- PR/commit refs if created
- one next action
- universal return packet from ORCHESTRATION.md
```
