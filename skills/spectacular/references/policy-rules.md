---
doc-id: policy
mode: structured
location: .spectacular/POLICY.md
scope: project-wide
template: templates/policy/base.md
snapshot-on-edit: true
summary: "Practice layer — policies filed under work-phase hooks"
status: active
always-set: true
---

# POLICY Rules

The practice layer. Always-set: scaffolded on every `spectacular init` with 8 prefilled policies (4 block / 4 warn). Structure-bound — CLI-managed *and* hand-editable. Full spec: [policies-contract.md](policies-contract.md). Runtime loop: [policy-injection.md](policy-injection.md).

**Verbs:**
- `grill` → polite no-op + hint: "POLICY ships prefilled. Add a policy by editing POLICY.md, or `--wide` to grill ad-hoc." (Policies are authored, not interviewed.)
- `refine` → tighten a policy's prose/check without changing its hook or id (snapshot first).
- `review` → **structure check** (not a placeholder check): every `### <id>` block has a `check`; `severity` is `block | warn` or absent; every `## @<hook>` is one of the locked 8; no orphan sections; `principle:` (if present) is a positive integer. Report violations; the mechanical subset also runs in `doctor policies`.

**Snapshot-on-edit: true** — POLICY is canonical project-wide; snapshot before any structural edit (e.g. `POLICY@v2.md`).

**Retrieval (not a doc-verb):** `spectacular policy [@hook | <id> | --principle N | --json]` is the runtime workhorse — the skill calls it on entering a phase, not the user. See [policy-injection.md](policy-injection.md).
