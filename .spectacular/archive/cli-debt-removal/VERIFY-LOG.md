---
request: cli-debt-removal
build: b4
verified: 2026-06-30
verifier: alex (via spectacular skill)
artifact: PLAN § Validation (no VERIFY.md — 2-of-6 rule)
result: pass
---

# Verify log — cli-debt-removal

Validation walk run against `PLAN.md § 6. Validation`. All work was shipped to
`main` prior to this walk; verification confirmed each check against the live
CLI + repo, not against the PLAN's claims.

| # | Check (PLAN § Validation) | Kind | Evidence | Result |
|---|---|---|---|---|
| M1 | D6 ADR in DECISIONS.md (removal list + MINOR rationale) | assertable | `DECISIONS.md:17` — D6 present, lists verbs + `--global` + `docs-*` refs, justifies MINOR (banner-warned since v1.2.0, pageworks replacement) | ✅ |
| M2 | `spectacular docs <anything>` no longer runs old verb; pageworks hint surfaces | observable | `grep 'docs (init\|export\|new\|review\|status)\|cmd_docs_' cli/spectacular` → empty; pageworks handoff hint live at `cli/spectacular:4543-4548` | ✅ |
| M3 | Deprecated ref docs gone; `--global` removed | assertable | `docs-contract.md` / `docs-rules.md` / `docs-renderer-adapters.md` absent; `grep -- '--global' cli/spectacular` → empty (only `--skill-scope global`) | ✅ |
| M4 | `--help` no deprecated surface; tests green; `doctor docs` discovery passes | executable | `grep` of removed surface in `tests/` → empty; `./cli/spectacular doctor docs` → 0 errors, 0 warnings, discovery-only | ✅ |
| M5 | CHANGELOG Removed section; manifests aligned | assertable | `CHANGELOG.md:132` `### Removed` (verbs + `--global`); manifests aligned 1.23.1 × 6 | ✅ |

## Residual debt swept during verify (ponytail confrontation)

The verify walk ran a ponytail pass over the residual surface and found one
piece the original cleanup missed:

- **Deleted `skills/spectacular/templates/docs/`** (whole dir — `docs.yaml.tmpl`,
  `index.md.tmpl`, `page.md.tmpl`). The dir scaffolded the removed `docs export`
  machinery; `docs.yaml.tmpl:29` still pointed at the deleted
  `references/docs-renderer-adapters.md`, and `docs.yaml.tmpl:27` documented the
  removed `docs export <renderer>` verb. Nothing in `cli/`, `skills/`, or `docs/`
  loaded any of the three templates — `check_docs` now directs users to
  `pageworks init` to scaffold `docs.yaml`. Pure dead residue (~40 lines).

After the sweep: zero refs to `templates/docs` anywhere; `doctor docs skill`
both green.

## Verdict

All 5 validation checks pass. Residual debt swept. Request → verified, ready to archive.
