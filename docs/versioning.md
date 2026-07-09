---
title: Versioning
description: How Spectacular versions itself — SemVer scheme, the breaking-change trigger, the single canonical version source, pre-release labels, and the optional marketing layer on the roadmap.
section: ""
status: stable
since: 1.9.0
updated: 2026-05-29
---

# Versioning

Spectacular is a **contract-driven dev tool** — it has a CLI surface (flags, verbs),
plugin manifests, and a `.spectacular/` file-format contract. So it versions with
**Semantic Versioning**, not dates or marketing numbers. This doc is the convention:
how the number is chosen, what triggers each bump, and where the single source of
truth lives.

> **TL;DR** — `MAJOR.MINOR.PATCH`. Default to the smallest bump the change warrants.
> Only MAJOR carries meaning beyond "it changed." There is no marketing layer in the
> version number itself — that lives on the [roadmap](#optional-marketing-layer-roadmap-only), opt-in.

---

## Scheme: Semantic Versioning

Format: `MAJOR.MINOR.PATCH` (e.g. `1.8.4`). Rules per [semver.org](https://semver.org/):

| Segment | Increment when… | Examples in Spectacular |
|---|---|---|
| **MAJOR** | A backward-**incompatible** change to a public contract | Renamed/removed a CLI verb or flag; changed `spectacular new` invocation syntax; changed the `.spectacular/` file-format contract in a way that breaks existing workspaces |
| **MINOR** | New, backward-**compatible** functionality | New CLI verb, new doc-type, new doctor area, new kit/pack — existing usage unaffected |
| **PATCH** | Backward-compatible **bug fix** only | Template-resolution fix, broken symlink, wrong path, ADR schema reconciliation |

Normative rules we follow (all from SemVer 2.0.0):

- **Numeric, not decimal.** `1.9.0 → 1.10.0`, never `2.0`. Segments are integers and can grow past 9.
- **No leading zeroes.** `1.8.4`, not `1.08.04`.
- **Reset on bump.** A MINOR bump zeroes PATCH (`1.8.4 → 1.9.0`); a MAJOR bump zeroes both (`1.9.0 → 2.0.0`).
- **Releases are immutable.** Never re-tag a shipped version. Ship a new one.
- **`0.y.z` was initial development.** Anything could break. Spectacular crossed into stable contract at `1.0.0`; we no longer use the `0.x` exemption.

---

## Choosing the next version

The default is **mechanical**: pick the smallest bump the change set warrants, by
highest-severity bucket.

| Highest-severity change in the set | Bump |
|---|---|
| Any breaking contract change | **MAJOR** |
| Any new backward-compatible capability (no breaks) | **MINOR** |
| Only fixes / internal changes | **PATCH** |

This is exactly what `/wrap-up` and `/update-docs` infer. In normal operation you
don't decide the version — the change set does.

### When the agent should ask first

Default to silent, correct increments. **Ask the user for an explicit target version only when:**

1. **A probable MAJOR is detected** — the change set includes a breaking contract
   change (removed/renamed verb or flag, altered invocation, file-contract break).
   A MAJOR is expensive and often a deliberate moment, so confirm the target
   (`2.0.0`?) rather than auto-bumping.
2. **The roadmap pins a milestone number** — if `ROADMAP.md` declares a marketing
   milestone (a deliberate `1.0`, `2.0`, etc. with stated goals), and the work in
   flight is what closes that milestone, confirm whether to land it *as* that
   milestone version vs. a routine increment. See below.

For MINOR and PATCH with no roadmap milestone in play: **don't ask, just increment.**

---

## The roadmap ledger — how builds map to versions

Before a change set *becomes* a release, it lives in the **roadmap ledger** — a table
at the top of `.spectacular/roadmaps/index.md` that is the **single source of truth for which
build ships in which version.** Understanding it is how you read "what's planned for
v1.x" without grepping.

**The model — a request never stores a version; it stores a build id.**

1. **`spectacular new <desc>`** scaffolds a request and stamps a **build id** on its
   `PLAN.md` frontmatter — `build: b17` — incrementing `last_build:` in `config.yaml`.
   The build id is permanent and never changes, even if the release target moves.
2. **You add a ledger row** when you slot the build into the roadmap (the CLI does *not*
   auto-insert rows — slotting is a human decision):

   ```
   | build | slug | title | tier | target-version | status |
   |-------|------|-------|------|----------------|--------|
   | b17   | roadmap-contract-docs | Document the ledger | themed | tbd     | planned |
   | b12   | lifecycle-undo        | Lifecycle undo      | full   | v1.22.0 | shipped |
   ```

3. **`target-version` is the only place a version number is written.** Not in PLAN
   frontmatter, not in prose. Reslotting a request to a different release is a **one-row
   edit** — change the cell, done.
4. **`tbd` means "slotted but not pinned yet."** When a build is real and prioritized
   but you haven't decided which release it lands in, its `target-version` is `tbd` —
   a committed sentinel, *not* a guessed number and *not* a blank. Prefer `tbd` over
   inventing a speculative version: false precision on unpinned work is exactly what the
   roadmap's precision gradient exists to avoid. Pin `tbd → vX.Y.Z` when the release is
   decided.
5. **Release-level `status` (`planned · active · shipped`) is distinct from request
   lifecycle** (`planned | active | review | verified` in PLAN.md). A request can be
   `verified` (the work is done and validated) while its ledger row is still `planned`
   (the release hasn't tagged yet). The row flips to `shipped` when the version tags.

**Why build ids instead of versions on the request?** Because versions move and builds
don't. If you store `v1.20.0` on a request and then reslot it to v1.22.0, you have to
find and edit every reference. With a build id, the request's own files never change —
only the one ledger cell does. The version is always a *read* from the ledger, never a
stored copy.

> **Shipped history lives in `CHANGELOG.md`, not the ledger.** The ledger tracks
> future + in-progress work; once a version tags, its detail belongs in the changelog.

Canonical schema (columns, tier legend, all the rules): `.spectacular/ARCHITECTURE.md`
§ Roadmap ledger. The two-layer roadmap model (ledger + per-version prose blocks):
`.spectacular/specs/roadmap.md`.

---

## The single canonical version source

Spectacular carries its version in several files. The convention: **one bump touches
all of them in lockstep — they are aliases, never independent.** This is the
"one canonical identifier" lesson (the alternative is the Windows
`20H2 = 2009 = 19042` trap, where aliases drift and confuse everyone).

Version-bearing locations, all bumped together at release:

| Location | Layer | Bumped by |
|---|---|---|
| `.claude-plugin/plugin.json` | project manifest | `bump-manifests.sh` (via `/wrap-up`) |
| `.claude-plugin/marketplace.json` (`.metadata` + `.plugins[]`) | project manifest | `bump-manifests.sh` |
| `.codex-plugin/plugin.json` | project manifest | `bump-manifests.sh` |
| `README.md` badge | project doc | `bump-manifests.sh` |
| `skills/spectacular/SKILL.md` (`version:` frontmatter) | component | **manual** — bumper does not touch component-level frontmatter |
| `cli/spectacular` (`SPECTACULAR_VERSION=`) | component | **manual** — bumper does not touch a shell constant |
| `CHANGELOG.md` (top entry) | record | written by `/update-docs` / `/wrap-up` |

> ⚠ **Known drift point.** The `bump-manifests.sh` automation only rewrites
> *project-level* manifests + the README badge. The two **component-level** sources —
> `SKILL.md` frontmatter and `SPECTACULAR_VERSION` in `cli/spectacular` — must be
> bumped by hand in the same release. `/wrap-up`'s Phase 6.5 post-write audit gate
> is what catches a missed manifest; it does **not** check the shell constant, so
> verify `cli/spectacular:5` matches the target before tagging.

**Rule:** a release is not done until every location above reads the target version.
No partial bumps.

---

## Pre-release labels

When shipping something unstable ahead of a stable release, use the standard SemVer
pre-release ladder. Lower precedence than the bare release; hyphen-delimited:

| Label | Meaning | Lockdown |
|---|---|---|
| `-alpha` | Cut from main, unstable, features may still be landing | features still in flux |
| `-beta.N` | Feature-complete, contract frozen, bugs expected | no new features / breaking changes |
| `-rc.N` | Release candidate — "this might be the one," QA rigor | bug fixes only |

```
2.0.0-alpha      → 2.0.0-beta.1 → 2.0.0-beta.2 → 2.0.0-rc.1 → 2.0.0
```

Build metadata (`+<meta>`) is allowed but ignored for precedence. Spectacular has not
needed pre-releases to date — most work ships straight to a stable PATCH/MINOR. Reach
for this ladder only when a MAJOR needs a soak period.

---

## Optional marketing layer (roadmap only)

The version number stays a pure technical/contract signal — **there is no marketing
layer inside `MAJOR.MINOR.PATCH`.** `2.0.0` means "a contract broke," not "a launch."

But the **roadmap** may carry a marketing layer, at the user's discretion. If the user
wants deliberate `1.0` / `2.0` milestones with stated goals — a "this is the launch"
narrative — that lives in `ROADMAP.md` as a named milestone with outcomes, *not* as a
reinterpretation of the version rules.

How the two layers interact:

- **User has roadmap milestones** (`ROADMAP.md` declares `2.0` with clear goals):
  when the work closing that milestone lands, the agent **confirms** whether to ship
  it as the milestone version (`2.0.0`) or as a routine increment. The milestone
  number and the SemVer number can deliberately coincide here — that's the user
  choosing to make `2.0.0` *both* the contract event and the marketing moment.
- **No roadmap milestones declared:** default behavior — normal mechanical
  increments, no asking, no marketing framing.

This keeps the default honest (smallest correct bump) while letting a user who *wants*
a launch narrative drive bigger version numbers intentionally — never the agent
inventing a marketing moment on its own.

---

## Optional: release-arc cadence (Apple-style)

A richer version of the marketing layer — also **opt-in, roadmap-only.** Where the
section above pins a single milestone *number*, this pins a milestone *shape*: a whole
major line as an arc, the way Apple runs an iOS year.

The pattern, observed across an iOS major line (e.g. iOS 18, Sep 2024 → Sep 2025):

| Phase | Example | Role |
|---|---|---|
| **`X.0`** | 18.0 | The launch — the marketing moment, the only one that gets a "keynote" |
| **`X.1`–`X.5`** | 18.1 Apple Intelligence, 18.2 image gen, 18.4 more AI… | **Staged feature delivery** — the rest of the promised set, drip-fed every ~1–3 months |
| **`X.6`–`X.7`** | 18.6 bug fixes | Convergence to **maintenance / security** only |
| **terminal `X.y`** | 18.7 (ships alongside 19.0) | **Terminal stable** — the version you rest on for the long tail before adopting `X+1.0` |

This is the `26.7-is-the-stable-before-27` behavior: the final point release of a line
becomes the durable stable, handed off the same week the next major launches.

**What's real vs. myth here:**

- **Real:** the *arc* — launch → staged features → maintenance → terminal stable. It's
  the emergent shape of a time-boxed major line.
- **Myth:** that a specific digit always means a specific thing (`.5` = always seasonal,
  `.7` = always terminal). It doesn't. The position in the arc carries meaning; the exact
  number is just cadence.
- **Note:** Apple's points are **cadence-driven, not SemVer-driven** — `18.1` shipping
  Apple Intelligence would be a *MINOR* under strict SemVer, yet Apple calls it a point.
  Apple flattens the engineering contract because the marketing layer drives their number.

**How Spectacular uses it (without breaking SemVer):**

The version number still obeys strict SemVer — the arc is a **roadmap narrative**, not a
re-interpretation of the bump rules. Concretely:

- The roadmap describes a major line as an arc: *"v2.0 launches the line; 2.x stages the
  remaining capabilities; the terminal 2.x is what we stabilize on before v3."*
- Each staged feature is still a real **MINOR** (`2.1.0`, `2.2.0`), each fix a real
  **PATCH** — the digits don't lie about the contract.
- The "terminal stable" is a roadmap designation, not a special version syntax: it's
  simply the last MINOR of the line, marked on the roadmap as the rest-here point.
- The agent never *invents* an arc. It applies one only when `ROADMAP.md` declares the
  major line that way — same opt-in rule as the milestone number above.

So: strict SemVer underneath (the contract), an optional Apple-style arc on top (the
narrative), and the roadmap is the only place the arc lives.

---

## Why SemVer and not CalVer / named releases

Considered and rejected for Spectacular:

- **CalVer** (`YYYY.MM`, like Ubuntu / JetBrains / Apple's `26`) — good for
  time-driven products with no meaningful API contract. Spectacular *is* a contract
  (CLI + file format), and "did this break my workspace?" is the question users
  actually have. CalVer can't answer it. Rejected.
- **Named / codename releases** (Apple big cats, Windows update names) — pure
  marketing layer, prone to alias sprawl. Rejected for the version; available as an
  optional roadmap-milestone narrative (above).
- **Apple-style point-arc** (`X.0` launch → staged `X.x` → terminal stable) —
  cadence-driven, not contract-driven; flattens MINOR/PATCH into one decimal. Not
  adopted as the *version* scheme, but its **arc shape** is available as an optional
  roadmap narrative on top of strict SemVer — see
  [release-arc cadence](#optional-release-arc-cadence-apple-style).

SemVer is the right scheme for a tool whose whole value is a stable, inspectable
contract.

---

## See also

- [commands.md](commands.md) — the CLI/skill surface that MAJOR protects
- [configuration.md](configuration.md) — `config.yaml` schema (a file-contract surface)
- `.spectacular/roadmaps/index.md` — where optional milestone narratives live
- `CHANGELOG.md` — the human-readable record of every version
