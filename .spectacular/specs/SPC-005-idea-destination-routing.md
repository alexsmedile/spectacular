---
id: SPC-005
type: specification
status: implemented
target_version: "tbd"
supersedes: ""
updated: 2026-08-04
summary: "Define destination-based idea routing and reduce local IDEA storage"
related:
  - "GitHub Issue alexsmedile/spectacular#6"
  - "skills/spectacular/references/idea-rules.md"
  - "skills/spectacular/references/github-work-bridge.md"
version: 1.1
approved_at: 2026-08-04
approved_by: alex
implemented_at: 2026-08-04
verified_against: uncommitted
---

# SPC-005 — Define destination-based idea routing and local IDEA deprecation

## Intent

Make the place where an idea will be developed and shared its authoritative
home. Reduce `.spectacular/ideas/` to a small, compatible local/private
capture surface without making quick capture ceremonial or assuming GitHub is
available, writable, or the only shared destination.

The resulting contract must distinguish discussion from prioritization and
execution: filing or linking a shared idea does not create a roadmap commitment
or a request.

## Requirements

### Destination model

| Destination | Authority and purpose | Audience and lifecycle |
|---|---|---|
| `issue` | Repository-specific collaborative report or proposal. The remote Issue owns its conversation. | Shared; may later route to direct work, a request, or spec-first. It is not an execution or roadmap commitment. |
| `shared` | A destination-neutral collaborative record, such as a Discussion, tracker card, or shared document. The referenced destination owns the conversation. | Shared; no GitHub-only assumption. |
| `roadmap` | An intentional project-priority record. Icebox is uncommitted; a version/tier carries progressively stronger commitment. | Project-wide; changes require a human choice. |
| local `IDEA` | Private/offline capture or locally useful refinement. It is not a default backlog or duplicate source of truth. | Local/private; capture → explore → hand off, archive, or delete. |
| `vision` | Request-scoped, human-reacted design evidence. | Optional after a request has been chosen; derives planning context but owns no product priority or execution lifecycle. |
| `request` | Durable execution coordination with PLAN/TASKS and validation. | Project/shared; created only after an explicit human decision. |

### Routing

Route by these facts, in order:

1. **Sensitivity:** protected/private material stays in a private or protected
   surface; it never enters ordinary shared or Issue publication.
2. **Visibility and scope:** cross-repo or collaborative discussion belongs in
   a shared surface; repo-specific actionable proposals normally belong in an
   Issue.
3. **Uncertainty:** an unformed or offline thought may remain local; a
   discussion surface is appropriate when collaborators must shape it.
4. **Commitment:** only an intentional priority decision enters the roadmap;
   only defined work needing durable coordination becomes a request.

No route is mandatory for fast capture. An IDEA is optional rather than a gate
between a thought and any other destination.

### CLI and skill contract

Keep the existing request promotion compatible, but make its destination
explicit in the replacement surface:

```bash
spectacular idea promote <slug> --to request
spectacular idea promote <slug> --to roadmap
spectacular idea promote <slug> --to shared --ref <stable-reference>
```

`request`, `roadmap`, and `shared` are the exact destination names. `issue` is
a documented subtype of `shared`, not a separate local authority or a generic
GitHub integration. A future convenience alias may accept `--to issue`, but it
must retain the same explicit-reference and no-remote-write guarantees.

`--to roadmap` must use the existing explicit Icebox/version/tier choice; it
must not add a version commitment implicitly. `--to shared` requires a stable
reference and only records a local handoff. It neither creates a remote record
nor copies a remote body. If no shared record exists, the skill presents a
copyable proposal envelope for the human's chosen filing tool.

The legacy omitted `--to` form remains a warning-producing alias for
`--to request` during the compatibility window.

### Deprecation and compatibility

1. Update guidance first: prefer the eventual development/shared destination;
   keep IDEAs for private/offline/local refinement only.
2. Keep `idea new|list|promote`, canonical IDEA IDs, existing archives, undo,
   and doctor checks working throughout the compatibility window.
3. Define private capture separately before adding it (for example,
   `.spectacular.local/ideas/`), including whether and how it is listed. Never
   silently migrate or commit private material.
4. Only after adoption evidence may a major release shrink IDEA counts,
   Wayfinding allocation, doctor behavior, and the legacy local directory.

### Boundaries

- No automatic mirroring or synchronized lifecycle across local and remote stores.
- No duplicate source of truth.
- No remote create, comment, label, close, or other side effect without a
  separate explicit human approval.
- No automatic roadmap commitment, request creation, or Vision creation.
- No expansion into the Issue #6 comment's broader bug/signal/security-agent
  capture proposal; that needs its own specification.

### Acceptance criteria

- Skill routing names the authority and route from uncertainty, visibility,
  cross-repo scope, commitment, and sensitivity.
- `idea promote --to request` retains the present scaffold/archive behavior;
  the legacy form works with a deprecation warning.
- `idea promote --to roadmap` refuses to mutate until the human chooses Icebox
  or a version/tier; it never infers a commitment.
- `idea promote --to shared` requires `--ref` and performs no network call.
- A shared Issue or Discussion remains authoritative; only links/provenance,
  never copied bodies, may be kept locally.
- Private/offline capture remains possible without making committed IDEA files
  the default store.
- Existing IDEA IDs, archive history, undo, and Wayfinding behavior stay
  compatible during the deprecation window.

### Proposed implementation and tests

Implementation is limited to idea-routing surfaces:

- `cli/spectacular`
- `skills/spectacular/SKILL.md`
- `skills/spectacular/references/idea-rules.md`
- `skills/spectacular/references/soft-db-index.md`
- `skills/spectacular/references/github-work-bridge.md`
- `skills/spectacular/references/roadmap-rules.md`
- `skills/spectacular/references/new-request.md`
- `skills/spectacular/templates/idea/base.md`
- `.spectacular/ARCHITECTURE.md`
- `tests/cli/idea.test.sh` (new)
- `tests/cli/undo.test.sh`
- `tests/cli/wayfinding-contract.test.sh`

The new CLI test file must cover `new`, `list`, every destination validation,
legacy compatibility, absent/invalid `--ref`, and the no-network guarantee.
Existing undo and Wayfinding tests must cover request compatibility and preserve
the no-active-scope-expansion invariant.

## Evidence and decisions

- [GitHub Issue #6](https://github.com/alexsmedile/spectacular/issues/6)
- Issue #6 maintainer direction: local IDEA storage is likely deprecated or
  reduced; routing must favour the eventual development and sharing surface.
- `idea-rules.md`: Issues and Discussions stay authoritative while discussed;
  do not mirror every capture.
- `github-work-bridge.md`: GitHub owns collaborative capture/discussion;
  Spectacular owns local reasoning and execution coordination only when earned.
- `roadmap-rules.md`: Icebox-to-version movement is intentionally manual
  because it is a commitment.

## Confirmation

draft — requires maintainer review and explicit approval before implementation.

**Approved 2026-08-04 by alex** — Maintainer approved implementation in this task
