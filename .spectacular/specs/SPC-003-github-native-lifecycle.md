---
id: SPC-003
type: specification
status: draft
target_version: "tbd"
supersedes: ""
updated: 2026-08-03
summary: "GitHub-friendly collaboration projections for requests, bugs, ideas, pull requests, and verification evidence"
related:
  - ../decisions/DEC-018-keep-spectacular-as-the-primary-lifecycle.md
  - ../decisions/DEC-019-every-executed-spectacular-request-ends-through.md
  - ../decisions/DEC-020-keep-spectacular-questions-as-local-blockers.md
  - ../decisions/DEC-021-gate-parallel-request-execution-with-a.md
  - ../roadmaps/index.md
---

# SPC-003 — GitHub-friendly collaboration projections for requests, bugs, ideas, pull requests, and verification evidence

## Intent

Make Spectacular friendlier to collaborative projects by integrating with GitHub's native collaboration surfaces without turning GitHub into a competing source of lifecycle truth. Spectacular remains the primary driver; GitHub carries selected public/collaborative projections, review activity, and remote execution evidence.

Implementation remains unauthorized while this specification is grilled. This draft may carry a concrete staged plan, but it is not eligible to seed a request until review and explicit approval.

## Requirements

### Confirmed constraints

- Spectacular remains authoritative for intent, dependencies, approved scope, decisions, and verification interpretation. GitHub provides linked collaboration and evidence surfaces rather than a duplicate canonical database.
- Every executed Spectacular request must end through a pull request before integration. Merge remains human-gated.
- Spectacular opens a draft pull request after the first meaningful implementation commit is pushed. Passing request verification makes the pull request eligible to become ready for review; the user still owns approval and merge.
- Open `QUE` records remain local blockers by default; publishing internal ambiguity is never automatic.
- GitHub integration has three capability profiles: `observe` reads and links without remote mutation; `adapt` maps an existing repository's Issue types, labels, forms, and workflow vocabulary to normalized Spectacular meanings; `managed` may propose and, only after explicit authorization, maintain GitHub forms, labels, Actions, and rulesets.
- `managed` is the first-class happy path for repositories the user controls and uses Spectacular to build. `adapt` and `observe` remain required compatibility/safety profiles for mature, external, or read-only repositories. Managed setup detects existing conventions before proposing a baseline and never treats general setup authorization as permission to replace governance or disclose private material.
- Spectacular does not require a universal label taxonomy. It detects repository capabilities, asks the user to confirm ambiguous semantic mappings, and enforces approved-spec, request, verification, and HITL invariants independently of label names.
- GitHub remains authoritative for remote reports, conversations, reviews, checks, permissions, and merge state. Spectacular stores durable links and the team's accepted interpretation rather than bidirectionally copying complete remote records.
- Committed `.spectacular/` is the repository's shared operational knowledge: project-wide decisions/findings/questions/specs plus request folders are reviewed and merged through normal branches and pull requests. Gitignored `.spectacular.local/` is private machine/user working state for incomplete thoughts, sensitive notes, and undeclared material; code and executable tests become final implementation truth after integration.
- Request-owned questions and research use request-qualified identities such as `wayfinding-sequencer/Q1` and `wayfinding-sequencer/R1`, physically owned by `.spectacular/requests/wayfinding-sequencer/`. The durable request slug, not a temporary branch name, is the namespace; branch and commit remain provenance. A project-wide promotion may allocate `QUE-NNN` or distill accepted research into `FND-NNN` while preserving the request-local source reference.
- Sequential project-wide IDs are not safe to allocate independently on concurrent branches. The initial GitHub integration must detect duplicate canonical IDs across all collections and must not use branch-qualified identity such as `branch#Q1`; collision-safe allocation/promotion beyond the single-orchestrator request namespace remains a foundational compatibility requirement.
- Branch isolation is necessary but insufficient for parallel work. Before finalizing a request scaffold and again immediately before activation/branch creation, Spectacular runs a **traffic preflight** over live requests plus observable branches and pull requests. It looks for explicit dependency edges, overlapping files/modules/specs, incompatible assumptions, and shared migration or release boundaries.
- Traffic has four canonical launch states: `parallel` means independent enough to proceed; `conditional` means it may proceed only while named boundaries are respected; `serialized` means it must wait for another request or integration event; `unknown` means evidence is insufficient and the orchestrator must ask the user rather than assume safety.
- Actual relationships are durable request knowledge. Store confirmed `blocked_by`, `blocks`, and `conflicts_with` links plus any named conditional boundaries in the owning request and validate reciprocal consistency. Traffic state is a time-bound assessment, not permanent truth: retain its timestamp and assessed Git/GitHub baseline as provenance, then recalculate it whenever the live work graph may have changed.
- The metaphoric language layer distinguishes three related views: **fog** is what is not yet understood, **frontier** is what is dependency-ready, and **traffic** is whether ready requests can safely move together. Human phrasing such as “check traffic” invokes the preflight; exact CLI spelling remains subordinate to the activity and must not create a second scheduling model.
- The BUG schema is storage-neutral. In a collaborative GitHub repository it structures or interprets the native Issue, whose durable identity is `owner/repo#N`; Spectacular does not create a duplicate `BUG-NNN` record. In standalone/single-writer projects without a collaborative tracker, the reserved local BUG soft database may use the same schema, but its activation and collision-safe multi-writer identity rules require a separate specification.
- A defect found and resolved inside the current request remains a TASKS/SESSION checkpoint plus regression evidence; no Issue or BUG is manufactured. Unclear or multi-site symptoms use a request-linked debug trace/AUDIT until classification. A confirmed out-of-scope collaborative defect becomes a structured GitHub Issue. Security vulnerabilities use the protected security workflow and never enter an ordinary Issue or BUG collection.
- Defects are routed as `introduced-by-current-work`, `pre-existing-blocker`, or `independent-discovery`. The first two may be corrected inside the approved request only when intended behavior is already established and the change stays within its authorization; ambiguity, material scope expansion, or behavior/API/schema change creates the appropriate QUE/SPC/HITL gate.
- Multi-day unfinished work is protected by request/session checkpoints, debug traces when earned, checkpoint commits on the isolated branch, and a draft PR after the first meaningful commit. A new BUG file would not make uncommitted work durable.
- Debug traces are live raw working state only while investigation remains open. On resolution, distill any durable value into `AUDIT`, verified `FIX`, or owning-request evidence, then move the raw trace to `archive/debugs/` so normal agents do not load it. This supersedes the older keep-resolved-traces-in-place convention. Security-sensitive traces never enter committed `debugs/` or its archive.
- Unclassified GitHub reports stay authoritative on GitHub and appear through the normalized triage view; Spectacular creates no committed duplicate triage inbox.
- Remote-write authority is tiered: `observe` permits no GitHub mutation; `adapt` confirms each remote write; `managed` may use explicitly pre-authorized action classes such as labeling, commenting, linking, and opening draft PRs. Merge, judgmental Issue closure, publishing a private idea, security disclosure, ruleset/CODEOWNERS changes, and deletion always remain explicit HITL gates.
- GitHub integration state has three storage layers. Committed `.spectacular/config.yaml` holds repository identity, profile, semantic mappings, the maximum permitted remote-action classes, and permanent HITL gates. Gitignored `.spectacular.local/github.yaml` holds machine/account, fork/push-remote, and temporary capability-cache settings. The owning request or AFK run holds narrower goal-scoped grants and mutation evidence; those grants expire at goal/request completion.
- Authentication tokens are never stored by Spectacular; `gh` owns credentials. Committed project policy is a permission ceiling, not a standing operator grant, and local remote names may differ across collaborators while the canonical target repository remains shared.
- GitHub synchronization is demand-driven. Spectacular fetches current remote state only at relevant activities such as overview, triage, reconcile, request/PR handoff, verification, and archive; authorized writes happen at the owning action. There is no daemon, webhook listener, background polling, authoritative cache, or automatic bidirectional synchronization in the initial implementation.
- PR-to-Issue linkage expresses closure intent explicitly: use `Fixes #N` only when the accepted Issue is fully resolved by that request and its resolution point is `on_merge`; use ordinary `Refs #N` text for related or partial work, or when a user-facing report must remain open until release. `Refs` is a human-readable convention rather than a GitHub closing keyword. Selecting a closing link is part of the authorized PR handoff, not inferred from a shared label.
- Linked Issues declare `resolution: on_merge | on_release`. `on_merge` is the default for engineering and internal work. User-facing reports may use `on_release`; their PR references the Issue without a closing keyword, and Spectacular proposes the release response and closure only after release evidence exists. Judgmental Issue closure remains a HITL action in either mode.
- Merge and availability are separate events. The normal senior workflow closes the Issue when its fix integrates and lets the merged PR/GitHub Release notes communicate availability. `on_release` is exceptional and reserved for customer-facing reports whose acceptance condition is actual availability; it does not create a parallel Spectacular release-obligation database.
- GitHub Issue, pull-request, review, and Discussion comments may inform work and serve as attributable evidence, but comment text alone never grants authority, resolves a `QUE`, creates a `DEC`, approves a specification, expands request scope, or advances lifecycle state. Authorization remains with recognized human operators—such as the user, designated orchestrator, authorized collaborators, or repository administrators—acting within their declared project and GitHub permissions. Spectacular records the actor, authority source, scope, and resulting action rather than inferring authorization from conversational wording.
- Authorized collaborators and repository administrators may approve technical execution and review outcomes inside an already approved request when their declared role permits it. New product/business decisions, changes to product intent, and material scope expansion still require the user or a designated product orchestrator; repository access or technical review authority does not imply that authority.
- Each repository explicitly maps product approvers, technical approvers, repository-governance administrators, and security approvers. GitHub write/admin access may contribute repository capability evidence but never implies product or security authority. Changing a product/security authority mapping requires an already authorized owner for that domain.
- Every authorization event records actor identity, authority source, exact action/scope, timestamp, and expiry when temporary. Request/AFK grants expire with their owner; removing an actor from a configured role invalidates unused grants. Spectacular performs a fresh role, permission, repository, and target check before important remote actions rather than trusting a stale authorization cache.
- Once the user approves a closed plan and authorizes its request/AFK run, ordinary reversible steps inside that recorded scope proceed without repeated confirmation. Spectacular stops when traffic changes, a mismatch or blocker appears, scope expands, required access was undeclared, or a permanent HITL gate is reached; one run-level authorization never grants merge, destructive cleanup, security disclosure, or other excluded actions.
- `park this idea` defaults to an uncommitted draft under `.spectacular.local/ideas/`; explicit promotion creates a committed canonical `IDEA` under `.spectacular/ideas/`; explicit publication creates a linked GitHub Discussion. The IDEA records the Discussion reference and publication provenance. The Discussion is authoritative for its public conversation but remains an input bucket, not the implementation-intent or approval authority. A recognized human synthesizes the initial formal idea plus Discussion into the canonical Spectacular request.
- GitHub Actions and checks are first-class only as external verification evidence in the initial integration. Spectacular may optionally install and maintain one narrowly scoped `spectacular/contract` check in `managed` mode to validate Spectacular lifecycle invariants. The project's existing build, test, deployment, and release workflows remain project-owned; Spectacular is not a general-purpose Actions or CI/CD workflow manager.
- GitHub Releases are first-class release evidence. Spectacular may inspect published releases, associate verified work with a release, determine when an `on_release` Issue becomes eligible for its authorized response and closure, and prepare draft release notes. Creating or publishing a GitHub Release remains an explicit human-authorized action; release evidence never bypasses that gate.
- `CODEOWNERS` and repository branch/ruleset configuration are first-class governance constraints that Spectacular may inspect to explain required reviewers, checks, and merge restrictions. In `managed` mode Spectacular may propose a diff, but changing reviewer ownership or merge protections always requires explicit authorization from a repository administrator and must never be applied as an incidental setup action.
- Security alerts and private vulnerability reports are protected inputs, not ordinary GitHub collaboration records. Spectacular may surface only that an authorized security blocker exists unless a designated security authority grants narrower access. It must never copy sensitive details into ordinary Issues, Discussions, comments, PR bodies, release notes, committed `.spectacular/` artifacts, routine logs, or agent handoffs. Disclosure, declassification, and publication always require explicit content-, audience-, destination-, and time-bounded authorization from a designated security authority.
- Security-related code that contains no confidential vulnerability information follows the normal request/PR lifecycle. Confidential remediation follows the provider-native protected workflow—GitHub Security Advisory and temporary private security fork when available—with sensitive plans/traces kept uncommitted. The advisory's protected review and merge counts as the security equivalent of the normal PR gate; Spectacular never forces a public PR, ordinary Issue, or unavailable CI path merely to satisfy its generic lifecycle convention.
- GitHub Environments and deployments are read-only verification evidence when an approved specification requires deployment to staging, production, or another named environment. Spectacular may report deployment state and use it to satisfy an explicit verification check, but it does not create environments, change protection rules or required approvers, approve deployments, trigger deployments merely to advance lifecycle, or read/write environment secrets.
- GitHub failure handling is fail-safe and evidence-based. Offline or unavailable GitHub does not block authorized local-only work, but every remote-dependent step remains explicitly pending. Missing or mismatched authentication stops the affected operation and asks the user to authenticate; Spectacular never initiates login automatically. Insufficient permission produces the exact missing-capability path and degrades to read-only behavior where possible. Forks and multiple remotes resolve through the committed canonical repository plus each collaborator's gitignored push-remote/account setting; unresolved ambiguity stops rather than guesses. Spectacular never claims synchronization, push, PR creation, review state, merge, release, deployment, or Issue closure without fresh remote evidence.
- Review and conversation input routes by accepted meaning rather than GitHub surface or label name. Before mutation, triage shows the source, interpreted meaning, proposed local destination, whether a duplicate artifact is justified, required authority/gate, and routing rationale. Users can inspect the complete mapping and override a proposal when their authority covers the destination; Spectacular records the accepted route and provenance.
- A proposed triage route is ephemeral and may be recomputed from fresh GitHub state. Once accepted, its route and provenance live on the resulting request, `QUE`, `IDEA`, `AUDIT`, or other justified destination artifact. When triage creates no local artifact, GitHub remains the sole durable record and Spectacular retains no duplicate inbox or standalone routing record.
- External pull requests are evaluated from their actual current state without fabricating a retroactive Spectacular implementation request. Small or expected changes that fit established intent may proceed through direct review and verification. A change to product behavior, public contracts, API, schema, security posture, or material scope stops for the applicable product/specification authorization before merge. Any additional work starts as a forward-looking request from that point; Spectacular never rewrites history to imply the contributor followed a process they did not use.
- Direct inspection and testing of an external PR is advisory Spectacular assistance, not a Spectacular request lifecycle claim. GitHub owns the durable review/check result. Spectacular may call work `verified` only through a real request with its own checks and evidence log; when that certification is wanted for an external contribution, create a small forward-looking review/adoption request linked to the PR.
- Retrieval has a semantic Spectacular layer without replacing expert tools. `spectacular find` resolves canonical IDs/aliases, GitHub references/URLs, and phrases; ranks live frontmatter/index results before bodies; follows justified lifecycle links; excludes history unless `--history` is explicit; and returns a compact agent briefing with security redaction. `rg` remains the direct code/text search tool and `gh` remains the advanced GitHub interface.
- Debug retrieval stays lean: bare `spectacular debug` lists active jobs, `spectacular debug <slug>` returns one validated state/resume card, and `--history` admits archived traces. `spectacular status` mechanically surfaces open user-input questions rather than relying on the agent to remember a separate wayfinding call.
- Spectacular wraps GitHub only when it adds cross-system value: combines local lifecycle with fresh remote state, normalizes repository-specific conventions, filters/redacts output, enforces Spectacular authorization/safety gates, records canonical provenance, or reconciles discrepancies. If an operation needs only GitHub, use `gh` directly. Spectacular does not re-create generic Issue/PR browsing, authentication, Actions-log inspection, repository administration, or arbitrary API access.
- Agents use Spectacular for governed lifecycle mutations such as request PR handoff, accepted semantic links, mapped remote changes, Issue closure, idea publication, and protected actions. Humans remain free to use GitHub or `gh` directly; reconciliation observes those changes later rather than pretending Spectacular can enforce every external path.
- Every active request has exactly one authoritative Spectacular orchestrator. Specialists may investigate, implement, and verify delegated work, but only that orchestrator allocates request-local identities, interprets authority, chooses semantic routes, changes lifecycle state, accepts durable findings, or performs governed GitHub mutations. Specialists receive bounded job cards and return evidence: explorer maps, investigator diagnoses, researcher checks public/vendor evidence, builder/fixer applies a closed brief, reviewer finds problems, verifier tests, and auditor compares claims with evidence.
- Every specialist brief declares source references, accepted meaning, authorized scope, actor/authority provenance, sensitivity, allowed reads, allowed writes, allowed destinations, expected return artifact, and current commit/PR head when relevant. Specialists stop and return unexpected scope, authority, or security discoveries to the orchestrator; they do not independently publish, relabel, close, approve, or advance lifecycle.
- The fleet contract is platform-neutral and lives in shared skill references. Platform-specific agent definitions are adapters to that contract, so Claude, Codex, and other supported harnesses preserve the same roles and gates. Confidential context remains with the authorized orchestrator unless a named specialist and destination receive an explicit narrower grant.

### Confirmed activity model

The canonical activity backbone is `capture → triage → explore → resolve → scope → build → verify → close`, with `configure` and `reconcile` as supporting activities. Confirmed initial use cases are:

- GitHub bug report → explicit triage using the Spectacular bug schema while retaining `owner/repo#N` as identity → implementation request/PR when accepted → verified fix and closure.
- Private idea draft → explicit promotion to canonical local IDEA → optional publication to a linked Discussion → human synthesis into a Spectacular request; neither publication nor feedback automatically approves work.
- GitHub conversation informs a local `QUE`; only explicit user intent resolves it or creates a `DEC`.
- External pull request → inspect actual change → directly review/verify when it matches established intent, or stop for product/spec authorization when it changes behavior/contracts → create only necessary forward-looking follow-up work.
- Goal-authorized AFK resolution of an Issue grills unresolved product decisions, implements on an isolated branch, opens a draft PR, and stops at HITL gates or before merge.
- Parallel request launch checks **traffic** before branch creation, records real request relationships, proceeds as `parallel` or `conditional`, queues `serialized` work, and returns `unknown` conditions to the user.
- Merge closes fully resolved `on_merge` Issues. Release closes eligible `on_release` reports after release evidence and user authorization. The completed request then archives its temporary request/spec context while retaining durable decisions, findings, and verified fixes.

### Confirmed route-by-meaning map

| Accepted meaning | Spectacular route | Mutation boundary |
|---|---|---|
| Clear correction inside approved scope | Owning request task/checkpoint | Authorized technical collaborator may accept; add regression evidence when behavior changed |
| Confirmed collaborative defect requiring cross-session work | GitHub Issue interpreted through the BUG schema | GitHub identity/report/conversation remain canonical; accepted implementation context lives on the owning request |
| Unclear intent, ambiguity, or trade-off inside a request | `<request-slug>/Q<N>` | The one request orchestrator allocates it; user or designated product orchestrator resolves product/business meaning |
| Project-wide open question | `QUE` | Commit through a PR; user or designated product orchestrator resolves product/business meaning |
| Request-owned research | `<request-slug>/R<N>` | Research supplies evidence, never a decision by itself; promote only durable accepted results |
| Verified reusable observation | `FND` when that reserved workflow is activated; otherwise attributed request/audit evidence | Must be verified and worth durable reuse; no premature reserved-record allocation |
| Uncertain, unreproduced, or multi-site technical problem | Debug trace and/or `AUDIT` | Investigate before declaring defect, finding, or fix |
| Incomplete private thought | `.spectacular.local/ideas/` | Gitignored by default; no canonical ID until explicit promotion |
| Useful promoted improvement outside current scope | `IDEA` | Commit without expanding the active request; explicit publication links a Discussion without transferring implementation authority |
| Requested behavior, contract, or material scope change | Spec revision or new `SPC` | Requires user or designated product-orchestrator authority before execution |
| Informative feedback with no durable local work | Remain on GitHub with an attributable link if needed | Do not duplicate merely because the comment exists |
| Suspected security vulnerability | Protected security boundary | Never route through ordinary Issues/artifacts until sanitized or explicitly disclosed |

This table is an overview, not an automatic classifier. Labels and Issue types contribute evidence, while accepted meaning and authority determine the route. Triage shows `proposed`, `accepted`, `deferred`, or `blocked` state for each actionable input and gives a concise reason whenever it does not create a local artifact.

### Confirmed security containment

- **Protected source remains canonical.** Read GitHub Security Advisories, private vulnerability reporting, Dependabot, or code-scanning alerts on demand through the authenticated provider. Do not mirror their bodies into Spectacular.
- **Default storage is no local copy.** Keep only an opaque provider reference, minimum routing status, and authorization provenance under gitignored `.spectacular.local/security/` when continuity requires it; create files with restrictive permissions. Even the existence or title of a vulnerability may be sensitive, so no committed reference is assumed safe. The reserved committed `.spectacular/security/` collection is not activated as a confidential store; a future workflow may admit sanitized, explicitly declassified records only.
- **Sensitivity follows derived content.** Anything read from a protected source is marked security-sensitive in the live operation; summaries, excerpts, generated fixes, branch names, commit messages, and draft payloads inherit that classification until a designated security authority explicitly declassifies specific content.
- **Outbound publication is deny-by-default.** A shared remote-write helper rejects security-sensitive payloads bound for ordinary Issues, Discussions, comments, PR metadata, Releases, or logs. A disclosure grant must name the exact content, destination, audience, and expiry; broad `managed`, AFK, collaborator, or repository-admin authority cannot substitute for it.
- **AFK and delegation stop at the boundary.** AFK work may report a redacted blocker but cannot disclose, declassify, or route protected content. Protected content is not passed to a sub-agent, browser, external model/vendor, or undeclared account unless that destination was explicitly authorized for the security task.
- **Safe public work stays sanitized.** Security remediation may use neutral branch/request/commit wording and tests with synthetic fixtures. Public PRs describe the safe outcome rather than exploit detail; coordinated disclosure material stays in the protected provider surface.
- **Verification matches the protected environment.** Use authorized local/manual checks when provider integrations or normal CI cannot access a temporary private security fork. Record only redacted evidence until disclosure, then re-run appropriate public checks against the disclosed/integrated head. Never claim an unavailable check passed.
- **Prevention is layered.** The non-disclosure rule is a non-overridable integration invariant implemented in the GitHub adapter and shared outbound-write boundary. Phase policies reinforce it but are not the sole control because project policy overrides are configurable. `doctor github` scans committed provenance and generated payloads for forbidden protected-source references and likely leakage, while making clear that heuristic secret scanning cannot prove safety.
- **Security authorization is a distinct role.** Product authority, technical approval, repository administration, and broad GitHub write access do not automatically confer security-disclosure authority. The project explicitly designates who may authorize protected access and disclosure.

### Explicitly deferred beyond the initial implementation

- Broader security-alert remediation and advisory management remain later work beyond the confirmed protected-read subset.
- Collaborative roadmap projection, GitHub Projects, and GitHub Milestones remain deliberately deferred until their workflows, mappings, and duplication risks are understood through a separate grill. The initial integration must not project Spectacular roadmap or request state into either surface.

## Concrete implementation plan

### Delivery strategy

Ship the smallest complete slice as two ordered requests after this SPC is approved:

1. **GitHub observe/adapt foundation** — capability discovery, confirmed semantic mapping, durable links, triage, draft/ready PR handoff, and read-only reconciliation.
2. **Managed GitHub enforcement** — optional forms, labels, contract Action, CODEOWNERS/ruleset proposals, and explicitly authorized remote apply. Start only after the first request is dogfooded on personal and collaborative repositories.

GitHub Projects, GitHub Milestones, collaborative roadmap projection, environment administration, and bidirectional/event-driven synchronization are outside both initial requests.

### M1 — Repository capability and authority contract

- Add a `github:` block to `.spectacular/config.yaml` with `profile: observe|adapt|managed`, canonical repository identity, permission ceiling/permanent HITL gates, default Issue resolution point (`on_merge`), and user-confirmed semantic mappings for type, triage, missing-information, blocked, area, priority, and resolution signals.
- Define committed role mappings for product, technical, repository-governance, and security authority without storing credentials. Specify which existing domain owner may change each mapping and reject self-grant or privilege escalation through ordinary GitHub setup.
- Add optional gitignored `.spectacular.local/github.yaml` for preferred account/host, fork and push-remote selection, and an expendable capability cache; never store tokens.
- Store temporary, narrower allowed remote-action classes plus mutation evidence in the owning request or AFK run, expiring at goal/request completion. Use the authenticated `gh` CLI and fail with a precise login/permission recovery path.
- Before every important remote mutation, re-resolve authenticated account, canonical repository/target, current GitHub permission, configured domain role, goal/request grant, and expiry; a cached capability result is advisory only.
- Define flat, Bash-3.2-friendly provenance fields for linked Issue, Discussion, and PR references. Store remote URLs/identities and accepted local meaning, never copied comment bodies or secrets.
- Implement graceful degradation for missing `gh`, offline operation, insufficient permissions, forks, and multiple remotes before enabling any writer. Preserve explicit pending state for every skipped remote-dependent action and require a fresh fetch before reporting its completion.
- Add a distinct security-authority mapping and non-overridable protected-content contract. Do not treat the general GitHub profile, permission ceiling, or ordinary repository role as security-disclosure authorization.
- Define the committed-project, private-local, request-local, GitHub-collaboration, PR-delivery, and code/test truth boundaries. Make `.spectacular/` shared and reviewable, `.spectacular.local/` gitignored, and branch/commit identifiers provenance rather than canonical artifact identity.
- Define request-local Q/R identity as `<request-slug>/Q<N>` and `<request-slug>/R<N>`, allocated only by the request's one authoritative orchestrator. Specify promotion/back-reference behavior for project-wide QUE and durable FND outcomes.
- Define the durable request relationship schema for `blocked_by`, `blocks`, `conflicts_with`, and conditional boundaries, including reciprocal consistency checks and stable request-slug references. Define traffic-assessment provenance separately so a past `parallel` result can never masquerade as current truth.

**Primary files:** `cli/spectacular`, `skills/spectacular/references/github-integration.md`, `skills/spectacular/references/lifecycle-contract.md`, `docs/configuration.md`.

### M2 — Observe and adapt

- Add a lean `spectacular github` namespace: overview, `setup`, `triage`, `link`, `pr`, and `reconcile`; keep exact flags subordinate to these activities.
- Apply a wrapper-admission test to every proposed subcommand: it must combine/filter/normalize/gate/record/reconcile. Reject a wrapper that merely forwards a `gh` command or renames its flags.
- Add `spectacular find <query> [--history]` as a local-first semantic resolver. Contact GitHub only when the query is a GitHub identity/URL or an explicit GitHub option is supplied; never turn ordinary phrase lookup into an implicit network call.
- `spectacular github`/`setup` inspects repository owner/type, current labels and Issue types, templates/forms, Discussions, checks, CODEOWNERS, and rulesets, then proposes a profile and semantic mapping. `observe` writes nothing; `adapt` snapshots config and applies only the confirmed mapping.
- Capability discovery reports security integrations only as available/unavailable and authorized/unauthorized; it never prints alert titles, bodies, affected paths, or counts that the provider treats as sensitive.
- The agentic `github triage` flow reads the Issue plus mapped signals, proposes one route, and grills only unresolved choices. The deterministic CLI owns reads, links, and confirmed mutations.
- `github triage` renders the route-by-meaning overview on demand and a compact routing card per actionable input: source, interpreted meaning, proposed destination, duplication decision, required authority, rationale, and current routing state. `spectacular github` summarizes route counts and blockers without loading comment bodies.
- Treat `proposed` routing state as a live computed view. Display `accepted` from destination provenance; derive `deferred` or `blocked` without a local destination from fresh mapped GitHub state rather than persisting a parallel Spectacular record.
- Do not create a committed triage inbox: unclassified reports remain authoritative on GitHub and appear through the normalized live view.
- Add the traffic preflight as a combined local/remote interpretation: inspect live Spectacular requests first, then observable branches and pull requests when GitHub is available. Missing remote evidence yields `unknown` or an explicitly local-only assessment rather than a false green light.

**Primary files:** `cli/spectacular`, `skills/spectacular/SKILL.md`, new `skills/spectacular/references/github-integration.md`, `docs/commands.md`, `docs/integrations.md`.

### M3 — Platform-neutral agent fleet contract

- Define one shared specialist job-card envelope covering source refs, meaning, scope, authority, sensitivity, allowed reads/writes/destinations, expected return, and current head identity.
- Enforce exactly one authoritative orchestrator per active request. Keep request-local ID allocation, semantic routing, lifecycle writes, durable-finding acceptance, authorization interpretation, and governed GitHub mutations on that orchestrator. Specialists never call the mutation wrapper directly; they return structured evidence or bounce.
- Update explorer, investigator, researcher, builder/fixer, reviewer, verifier, and auditor roles to consume the shared envelope and preserve their existing narrow responsibilities.
- Add protected-return behavior: suspected security findings stop ordinary output/persistence, return a redacted blocker through the authorized channel, and cannot enter normal debug logs or PR comments.
- Put the normative contract in skill references and keep `agents/*.md` plus harness-specific dispatch instructions as thin adapters, preventing Claude-only behavior drift.

**Primary files:** new `skills/spectacular/references/agent-job-contract.md`, `github-integration.md`, `build-workflow.md`, `bug-workflow.md`, `review-sweep.md`, `agents/*.md`, `skills/spectacular/SKILL.md`, `tests/agents/`.

### M4 — Branch-aware defect routing and duplication control

- Classify implementation discoveries as `introduced-by-current-work`, `pre-existing-blocker`, or `independent-discovery` before deciding where they live.
- Keep a quick, in-scope defect in the owning request: add the corrective TASKS/SESSION checkpoint, fix it, add a regression test, and continue. Do not manufacture an Issue, AUDIT, or FIX entry for an ordinary intermediate mistake.
- Open a debug trace/AUDIT when the root cause is unclear, multi-site, unreproduced, or may not be a defect. Fold its disposition into the current request, a future request, a verified FIX, or a deliberate decline.
- On debug resolution, verify the distilled destinations and archive the raw trace through a deterministic mutator that preserves its slug, outcome, and back-links. Normal debug listing/resume ignores archived traces unless history is requested explicitly.
- Add bare `spectacular debug [<slug>] [--history]` for active-job overview and validated resume briefing; do not add separate list/show/resume subcommands until usage proves they are necessary.
- Reuse the Spectacular bug schema to create or normalize collaborative GitHub Issues, retaining `owner/repo#N` as the only durable bug identity and reading remote state on demand. Do not persist a mirrored `BUG-NNN` record in standard GitHub repositories.
- Keep the local BUG soft database reserved for standalone/single-writer projects without a collaborative tracker. Do not activate it here; define collision-safe multi-writer identity and its explicit lifecycle in a separate specification before supporting collaborative local BUG allocation.
- Project an independent defect to GitHub only when it needs durable cross-branch collaboration or visibility. Keep unreleased/private discoveries on the isolated branch until publication is authorized; link accepted Issue meaning from the owning request instead of copying Issue bodies/comments.
- Route any suspected vulnerability to the protected security boundary before ordinary defect publication. When classification is uncertain, fail closed and ask a designated security authority rather than opening an Issue.
- Record branch/request provenance and replace branch-only identity with an exact commit SHA once one exists. A GitHub Issue remains repository-level even when linked to that branch or PR.
- Add `doctor github` checks for one remote Issue linked inconsistently to multiple active local work owners, missing local targets, malformed remote identities, and forbidden copied sensitive content; judgment findings never auto-delete or auto-close. Expand canonical duplicate detection across every active collection, including advanced optional collections.

**Primary files:** `soft-db-index.md`, new `bugs-rules.md`, `bug-workflow.md`, `debug-trace.md`, `artifact-retention.md`, `audit-rules.md`, `fixes-rules.md`, `request-workflow.md`, `canonical-ids.md`, `doctor-areas.md`, `cli/spectacular`.

### M5 — Request and pull-request lifecycle

- Run traffic preflight while scaffolding a request so dependencies and conflict candidates are visible before the user commits to the plan. Persist confirmed relationships, named conditions, and the assessment baseline without treating provisional overlap guesses as decisions.
- Re-run traffic preflight immediately before request activation and branch creation. `parallel` proceeds; `conditional` records and enforces its named boundaries; `serialized` remains planned/held until its prerequisite changes; `unknown` requires user clarification. A changed result invalidates the earlier launch assessment without deleting its provenance.
- After closed-plan/run authorization, execute routine in-scope branch, implementation, verification, and draft-PR steps without stepwise confirmation. Stop only on traffic changes, mismatch, blocking work, undeclared access, scope expansion, or an existing HITL boundary.
- Extract general PR behavior from `cmd_afk_pr`; GitHub PR handling belongs to the integration layer while AFK supplies the authorization and branch-isolation context.
- `github pr open <request>` requires an active request, approved source spec, non-primary pushed branch, and a first meaningful implementation commit confirmed by the agent. It opens a draft PR containing stable request/spec/Issue references and the current head SHA.
- The draft PR uses `Fixes #N` only for a confirmed full resolution whose Issue declares `resolution: on_merge`. Partial work and `on_release` reports use `Refs #N`; the latter cannot receive a closing keyword until release evidence exists. Human-facing output says `Issue #N`, while stored provenance uses `owner/repo#N` or the full URL.
- For `on_release`, render a durable, machine-readable PR-body declaration equivalent to `Refs #N`, `Availability: vX.Y.Z|tbd`, and `Resolution: close on release`. The open Issue plus that PR declaration is the reminder; do not copy it into a separate local release queue.
- Before rendering or opening a PR, run the shared protected-content guard over title, body, branch, commit-derived text, linked references, and generated release-note fragments. Refuse publication on a positive or indeterminate classification until security authority supplies a sanitized payload or exact disclosure grant.
- Detect provider-native confidential remediation before enforcing the ordinary draft/ready PR sequence. Require the advisory's authorized protected review/merge evidence instead, while preserving the permanent human merge/disclosure gates.
- `github pr ready <request>` requires request verification against the current head, acceptable required checks, no blocking `QUE`, and explicit readiness confirmation; it never merges.
- Preserve `spectacular afk pr` as a compatibility path that delegates to the new PR lifecycle without weakening its AFK gates.
- Document the raw-tool escape hatch beside every wrapper: use `gh issue|pr|run|api` for GitHub-only inspection and administration; use the Spectacular wrapper when a governed local lifecycle mutation or combined interpretation is required.
- External PRs remain observable without forced backfill. In adapt mode material changes produce an approval warning; managed enforcement may make that warning a required check.
- When an external PR needs corrective or complementary work, create a new request linked to the PR and its accepted review findings; do not create a request whose timestamps or lifecycle claims pretend to precede the contribution.
- Keep direct-review output durable on GitHub. Do not create a local `VERIFY-LOG` or use the `verified` label unless a forward-looking review/adoption request owns that evidence.

**Primary files:** `cli/spectacular`, `afk-git-hygiene.md`, `request-workflow.md`, `verify.md`, PLAN/PR provenance handling, `docs/workflow.md`.

### M6 — Verification and reconciliation

- Import GitHub checks, review state, PR head SHA, merge state, and linked Issue state as evidence references; Spectacular still interprets whether they satisfy the spec.
- Import published GitHub Release identity and timestamps as evidence for release-gated work. Draft notes remain a proposal, and neither note generation nor release detection constitutes permission to publish a release or close an Issue.
- When a release publishes, inspect its included merged PRs for `close on release` declarations, match their declared version when present, and propose the response/closure for still-open referenced Issues. Leave `tbd` or mismatched declarations open and report them; never infer eligibility merely because some release exists.
- Import deployment/environment status only for verification checks that explicitly name it. Treat deployment approval and environment protection as external gates, never as mutations Spectacular may perform to make verification pass.
- Import comments as attributed evidence or suggestions only. Any action requiring authority must cite a separate valid authorization event and the actor's applicable local/GitHub role; quoted comment text cannot satisfy that gate.
- Validate the authorization event against both dimensions: the actor's declared role and the action's scope. A collaborator's valid technical approval cannot be reused as product authorization or permission to expand the approved request.
- `github reconcile` is read-only by default and reports: merged PR with live request, closed Issue with an unresolved linked AUDIT/request, merged `on_merge` Issue still open, unreleased `on_release` Issue closed early, released `on_release` Issue awaiting an authorized response/closure, verified FIX whose Issue lacks a resolution response, verification stamped against an older head, broken links, and changed/unmapped label semantics.
- Security reconciliation emits redacted state only and never imports protected bodies into the normal request, session, finding, memory, archive, or verification evidence paths.
- For protected remediation, distinguish redacted local/manual verification from later public-head verification and surface any required post-disclosure recheck; absence of inaccessible CI is not itself a failure, but an invented pass is forbidden.
- Reconciliation never chooses an authority silently. Remote corrections and local lifecycle changes remain separate, previewed actions.
- Session start and request/archive briefings surface only actionable linked discrepancies, not a full GitHub activity feed.
- Make raw `spectacular status` prepend open `requires_user_input` questions, matching the skill briefing and session-start contract; JSON output must expose the same blocker data without mixing prose into the request array.

**Primary files:** `cli/spectacular`, `verify.md`, `status.md`, `archive.md`, `doctor-areas.md`, `docs/integrations.md`.

### M7 — Managed enforcement, after dogfood

- Make `managed` the documented first-class setup for repositories the user controls while retaining `adapt`/`observe` for compatibility and read-only work. `github setup --profile managed` first renders a dry-run diff for proposed Issue Forms, native Issue-type/field mappings where available, fallback labels, PR template, optional Ideas/Q&A Discussion categories, contract Action, CODEOWNERS suggestions, and ruleset requirements.
- Prefer native Issue types and fields, then fall back to the minimal semantic labels `type:bug|feature|task|research|question`, `needs-triage`, `needs-info`, `blocked`, `breaking-change`, selected high-impact priority labels, and discovered repository-specific `area:*` labels. Do not require a Spectacular-branded label or duplicate native open/closed, assignment, linked-PR, merge, branch, or release state in labels.
- Treat GitHub repository topics as repository discovery metadata, not Issue classification. Use Discussion categories for public Ideas/Q&A and keep authorization questions local.
- Remote writes require authenticated capability checks plus explicit `--apply --yes`; destructive replacement of existing repository configuration is refused.
- Add a `spectacular/contract` check that validates committed local invariants: approved source, linked request/PR, current-head verification, declared breaking-change authorization, and no unresolved blocking question.
- Do not generate, rewrite, or assume ownership of unrelated build, test, deployment, or release workflows.
- Ruleset/CODEOWNERS application remains separately confirmed because it changes collaborator permissions and merge behavior.
- Inspect existing ownership and rules before proposing repository configuration; never treat `--apply` for general managed setup as sufficient authorization for a governance change.

**Primary files:** generated `.github/` templates/workflow assets, `cli/spectacular`, managed setup reference, docs, fixture-based tests.

### Verification and test plan

- Add `tests/cli/github-integration.test.sh` with a fake `gh` executable and fixtures for personal repos, organization Issue types, existing custom labels, forks, offline/auth failures, and insufficient permissions.
- Test wrong-account detection, ambiguous remotes, canonical-target plus local push-remote resolution, read-only degradation, pending remote actions, and refusal to report success from cached or assumed state.
- Test per-domain role separation, authorized mapping changes, rejected self-grants, complete approval provenance, request/AFK expiry, role-removal revocation, stale-cache refusal, and fresh pre-mutation role/permission checks.
- Extend `afk-git-hygiene.test.sh` for draft-open then verified-ready behavior and the compatibility path.
- Test current-work defect → request checkpoint only; unclear blocker → AUDIT; independent defect → optional linked Issue; verified reusable repair → FIX; and conflicting Issue → multiple active work-owner detection.
- Test that branch creation and immediate request-local fixes never allocate BUG; collaborative defects retain GitHub identity; the normalized view and owning request link rather than copy GitHub bodies/comments or mirror remote status; standalone BUG activation remains unavailable pending its separate contract.
- Test request-qualified Q/R allocation by the sole orchestrator, collision isolation across request slugs, stable identity after branch deletion, branch/SHA provenance, promotion back-links, and specialist allocation refusal.
- Test scaffold-time and activation-time traffic preflights; all four canonical states; durable relationship reciprocity; conditional-boundary enforcement; stale-baseline invalidation; local-only degradation; and refusal to infer `parallel` when remote or overlap evidence is insufficient.
- Test that one closed-plan/run authorization permits ordinary reversible in-scope steps while traffic changes, mismatch, blockers, undeclared access, scope expansion, merge, destructive cleanup, and protected actions still stop at their gates.
- Test duplicate canonical IDs across all active standard and optional collections, including same-ID/different-filename merges.
- Test every confirmed route-by-meaning row, full overview rendering, compact routing cards, authority-aware override, accepted-route provenance, no-artifact rationale, and summary counts that do not require copying comment bodies.
- Test semantic find by canonical alias, GitHub identity/URL, phrase, live-vs-history ranking, lifecycle link traversal, local-only default, concise output, and protected-content redaction; retain `rg`/`gh` escape-hatch documentation.
- Test wrapper admission with positive combined/gated cases and reject passthrough-only command proposals; verify that agents choose Spectacular for governed mutations while direct human/`gh` changes remain discoverable through reconciliation.
- Add cross-role contract tests for complete job cards, missing authority/sensitivity/head fields, specialist mutation refusal, out-of-scope bounce, orchestrator-only persistence, redacted security returns, and equivalent Claude/Codex routing behavior.
- Test open debug visibility, resolution distillation, deterministic `archive/debugs/` movement, explicit historical lookup, preserved back-links, normal-context exclusion, and refusal to commit/archive security-sensitive traces.
- Test lean debug overview/detail/resume cards without extra subcommands, plus identical open-question visibility across raw status, skill briefing, session start, and machine-readable status output.
- Test that proposed routes leave no durable record, accepted routes persist only on their destination, and no-destination outcomes remain GitHub-only across a fresh session.
- Test closing-keyword generation for full `on_merge`, partial `Refs`, and `on_release` references; verify that release-gated reports cannot close from PR merge and that release evidence does not bypass the Issue-closure HITL gate.
- Test durable `Availability`/`Resolution` PR-body declarations, release-to-merged-PR matching, `tbd` and wrong-version refusal, and absence of any duplicate local release-obligation ledger.
- Test that imperative wording in a comment cannot approve a spec, resolve a question, expand scope, advance lifecycle, or authorize a remote mutation without a separate valid authority event.
- Test collaborator/admin technical approval within approved scope and reject the same actor's attempted product decision or material scope expansion unless they are separately designated as a product orchestrator.
- Test direct review of a conforming external PR, product/spec gating for behavior/API/schema/security changes, forward-only follow-up request creation, and refusal to manufacture retroactive lifecycle provenance.
- Test that direct external-PR assistance cannot emit a Spectacular `verified` claim, while a linked forward-looking review/adoption request can complete the normal verification lifecycle.
- Test that Discussion publication requires explicit authorization and that Discussion feedback remains triaged input rather than directly mutating the canonical local idea or lifecycle.
- Test private idea parking under gitignored `.spectacular.local/ideas/`, explicit committed IDEA promotion, publication provenance/back-link, and human synthesis into a request without treating the Discussion as implementation authority.
- Test ingestion of existing check results separately from the optional `spectacular/contract` check, and ensure managed setup never rewrites unrelated Actions workflows.
- Test release lookup, verified-work association, draft-note preparation, and `on_release` eligibility without permitting automatic release publication or Issue closure.
- Test read-only governance discovery, required-reviewer/merge-restriction reporting, proposed diffs, administrator-role validation, and refusal to apply governance changes through general setup authorization.
- Test protected-source reads with synthetic payloads; taint propagation into summaries and derived text; denial across every ordinary outbound surface; restrictive local metadata storage; redacted errors/logs; AFK and delegation refusal; exact disclosure-grant scope and expiry; and false-negative-resistant fail-closed behavior for indeterminate content.
- Test normal non-confidential security code through the ordinary PR flow; confidential advisory/private-fork routing; no forced public PR/Issue; authorized local/manual verification when CI is unavailable; redacted evidence; protected merge equivalence; and post-disclosure public re-verification.
- Test named-environment deployment evidence and refuse environment creation, protection/approver changes, deployment approval, unsolicited deployment triggering, and secret access.
- Test that observe mode performs no writes; adapt mutates only local confirmed mapping; managed stays dry-run until both apply flags are present.
- Test managed capability fallback from native Issue types/fields to minimal semantic labels, mapping of established repository vocabulary by meaning, absence of mandatory Spectacular branding, and refusal to misuse repository topics as Issue classification.
- Run Bash 3.2 syntax, full CLI suite, lifecycle/links/github doctors, version guard, and two dogfood repositories before approval of the managed-enforcement request.

### Documentation and release impact

- Route the new GitHub activities from `SKILL.md` through one focused reference file rather than enlarging the orchestrator.
- Update `docs/integrations.md`, `docs/workflow.md`, `docs/commands.md`, `docs/configuration.md`, CLI help, README's short integration summary, changelog, capability index, and affected lifecycle/AFK references.
- Keep GitHub Projects, GitHub Milestones, and public roadmap usage explicitly deferred in the docs; do not imply adoption or automatic status projection.

## Evidence and decisions

- `DEC-018` — Spectacular is the primary lifecycle authority; GitHub is a collaboration/evidence projection.
- `DEC-019` — executed requests end through a pull request before integration.
- `DEC-020` — questions remain local blockers by default.
- `DEC-021` — parallel request execution is gated by a traffic preflight; real relationships persist while launch eligibility is recalculated from current evidence.
- GitHub Issue [#3](https://github.com/alexsmedile/spectacular/issues/3) — source report for dependency/conflict discovery, parallel-branch eligibility, and run-level authorization ergonomics.
- User-confirmed GitHub lifecycle design recorded through 2026-08-03, including committed shared `.spectacular/` knowledge, private `.spectacular.local/` drafts, request-qualified Q/R identities, one authoritative orchestrator per request, storage-neutral bug schema with GitHub-native collaborative identity, private→promoted→Discussion idea flow, and managed-mode-first semantic GitHub setup. Items explicitly deferred remain outside the initial implementation and require a future dedicated design pass.

## Confirmation

The design grill is complete. SPC-003 remains a draft, b40 remains a `candidate`, and its target version remains `tbd`; none of those states authorizes implementation. A separate review and explicit approval are required before any implementation request may be created.
