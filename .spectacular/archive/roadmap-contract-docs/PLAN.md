---
status: archived
priority: high
owner: alex
updated: 2026-07-06
build: b17
summary: "Document the roadmap ledger as a real contract: spec the build-id→version model in specs/roadmap/SPEC.md + SPEC.md, define the `tbd` target-version sentinel (and stop roadmap-rules.md from rejecting it), and add user-facing docs + a short tutorial so build ids ↔ release versions are clear without reading internal ARCHITECTURE.md."
related:
  - PRD.md
  - ../../SPEC.md
  - ../../specs/roadmap/SPEC.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
archived: 2026-07-06
---

# Plan — roadmap-contract-docs

> **Origin (2026-06-28):** An audit found the roadmap ledger system (build ids `b1..bN`,
> `target-version` as single source of truth, release-level `planned/active/shipped`
> status, and the `tbd` "not-pinned-yet" sentinel) is **real and actively used** but
> **documented in exactly one internal place** — `ARCHITECTURE.md:208-269`. It is absent
> from both `SPEC.md` and the (stale, pre-ledger) `specs/roadmap/SPEC.md`, from all user
> docs (`docs/commands.md` has no `spectacular roadmap` section; `last_build:` config field
> undocumented), and from every onboarding/tutorial doc. The `tbd` convention is defined
> nowhere — worse, `roadmap-rules.md`'s universal placeholder check **rejects `<TBD>`**,
> contradicting the maintainer's own ledger usage. A user can only learn the build-id→version
> model by reading internal architecture reference.

## 1. Goal

Make the build-id→version model **discoverable and authoritative** without touching behavior:

1. **Spec the ledger** where specs live — `specs/roadmap/SPEC.md` (currently pre-ledger, stale) + the `SPEC.md` capability bullet — so the ledger schema, build ids, `target-version`-as-single-source, and ledger-status-vs-request-lifecycle are specified, not just architecture-noted.
2. **Define `tbd`** as the legitimate "slotted but not version-pinned yet" sentinel — in the ARCHITECTURE.md schema row + `roadmap-rules.md` — and stop the placeholder check from flagging it.
3. **Document for users** — a `spectacular roadmap` section + the build-id↔version model in `docs/`, plus a short tutorial/walkthrough so a new user understands the roadmap file without reading `.spectacular/ARCHITECTURE.md`.

## 2. Constraints

- **Docs/spec only — zero behavior change.** No CLI logic, no new verbs. This request closes a documentation gap; the ledger already works. (Pruning/scaling is the separate `roadmap-pruning` request — b18.)
- **Single source of truth stays single.** Documenting the schema in `specs/roadmap/SPEC.md` must *point to* ARCHITECTURE.md as the canonical schema home (or move it there cleanly) — not fork a second authoritative copy that can drift. Decide which file owns the schema vs which references it.
- **`tbd` fix must not weaken the real placeholder check.** `<TBD>` in a *prose slot* should still be rejected; `tbd` as a *ledger target-version value* is legitimate. The rule must distinguish the two contexts, not blanket-allow `tbd`.
- **Snapshot before overwriting `specs/roadmap/SPEC.md`** (it's a canonical doc — `spectacular snapshot` first per the versioning rule).

## Understanding

### How it works now (from the audit)

- **Schema home:** `ARCHITECTURE.md:208-269` is the *only* narrative spec of the ledger — columns (`build/slug/title/tier/target-version/status`, lines 212-231), tier legend (234-239), status values (241-247), the crucial "ledger status ≠ request lifecycle" note (249: a request can be `verified` while the row is still `planned`; flips to `shipped` at tag), and rules (251-257, incl. "shipped history lives in CHANGELOG, not the ledger").
- **`specs/roadmap/SPEC.md`** (51 lines, `updated: 2026-05-29` — predates the v1.17.0 ledger): specs only the *per-version prose blocks* (precision tier, 9-phase chain, outcome slot, 18-check review gate, Icebox). **Zero ledger awareness** — it doesn't know the ledger table sits above the blocks.
- **`SPEC.md:35`** structured-ROADMAP bullet: describes the prose blocks + review gate, says nothing about the ledger/build ids.
- **`docs/`:** `commands.md` has no `spectacular roadmap` section (only `init --with roadmap`, a grill override, and `roadmap` as a doc-id); `configuration.md` doesn't document `last_build:`; `versioning.md` covers SemVer but never links to "version is derived from the ledger."
- **`tbd`:** ad-hoc in ROADMAP.md (rows b3/b4/b10/b11/b16 + prose). Defined in no spec/rule/doc. `roadmap-rules.md` universal base check ("no `<TODO>`, `<TBD>`, `???` in any slot") would reject it.
- **Onboarding:** `onboarding.md` / `guided-first-run.md` / README / `.spectacular/AGENTS.md` mention ROADMAP only as "a doc to read" — none explains the model.

### What changes (docs + spec text only)

1. `specs/roadmap/SPEC.md` — add a "Ledger" section covering build ids, the table schema, `target-version` as the single written version, ledger-status-vs-lifecycle, and `tbd`. Refresh frontmatter (`status: published`, `updated: 2026-06-28`). Reference ARCHITECTURE.md as canonical schema home (constraint 2).
2. `SPEC.md:35` bullet — add a clause (or sibling bullet) naming the build-id ledger + pointing to `specs/roadmap/SPEC`.
3. `ARCHITECTURE.md:230` — `target-version` row notes `tbd` as the legitimate not-yet-pinned sentinel, distinct from a real `vX.Y.Z`.
4. `roadmap-rules.md` — add a ledger-section behavioral rule: "use `target-version: tbd` when a version isn't really planned yet"; scope the placeholder check so `tbd` in the *ledger column* is allowed while `<TBD>` in prose slots is still rejected.
5. `docs/commands.md` — a `spectacular roadmap` section + the build-id↔version model; `docs/configuration.md` — document `last_build:`.
6. A short tutorial/walkthrough (location TBD — `docs/` page or an onboarding addition) that walks: request gets a `build:` id → ledger row with `tbd` → row gets a `target-version` when slotted → flips to `shipped` at tag. Connect to `docs/versioning.md`.

### What stays the same

The ledger mechanism, the prose-block spec, the review gate, the CLI. No file moves, no behavior.

## 3. Milestones

### M1 — Spec the ledger
- `specs/roadmap/SPEC.md`: snapshot first, then add Ledger section (build ids, table schema, target-version single-source, status-vs-lifecycle, `tbd`); refresh frontmatter.
- `SPEC.md` bullet updated to name the ledger + link.
- Resolve constraint 2: which file is canonical schema home, which references.
- `doctor specs` green.

### M2 — Define `tbd` + fix the contradicting rule
- `ARCHITECTURE.md:230` target-version row documents `tbd`.
- `roadmap-rules.md`: ledger `tbd` rule + scope the placeholder check (prose `<TBD>` still rejected).
- Test/check: a ledger row with `tbd` no longer trips the placeholder gate; a prose slot with `<TBD>` still does.

### M3 — User docs + tutorial
- `docs/commands.md` `spectacular roadmap` section + build-id↔version model.
- `docs/configuration.md` documents `last_build:`.
- Short tutorial walkthrough; link from `docs/versioning.md`.
- VERIFY-LOG.

## Validation

Docs/spec-only request → verification is read-and-confirm, recorded in `VERIFY-LOG.md`. Success criteria:

- {assert} `specs/roadmap/SPEC.md` has a ledger section (build ids, target-version, tbd, status-vs-lifecycle); frontmatter `published`.
- {assert} `ARCHITECTURE.md` target-version row documents `tbd`; `roadmap-rules.md` placeholder check scoped to prose slots + has the ledger `tbd` rule.
- {assert} `docs/versioning.md` carries the ledger walkthrough; `configuration.md` documents `last_build:`; `commands.md` has the `spectacular roadmap` section.
- {judge} no second authoritative schema copy (ARCHITECTURE canonical, spec points to it); zero behavior change.
- `run: ./cli/spectacular doctor specs links docs` → 0 errors.

## M-questions (resolved at grill 2026-06-28)

1. **Schema canonical home:** ✅ ARCHITECTURE.md stays canonical; the spec adds a ledger section that *points to it* — no second drifting copy.
2. **Tutorial location:** ✅ folded into `docs/versioning.md` § "The roadmap ledger" (the natural home — how versions get assigned), linked from commands.md + configuration.md, rather than a standalone page.
