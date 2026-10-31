# Review-method research

Use this when: explaining provenance, refreshing external tips, or deciding whether
an upstream rule belongs in this skill. These sources are material to compare, not
runtime dependencies or instructions to install or execute. Consulted 2026-10-31;
links track upstream branches and may change. Guidance below is our synthesis.

| Source | Useful decision steering adopted |
|---|---|
| [Vibe Code Auditor](https://github.com/sickn33/agentic-awesome-skills/blob/main/skills/vibe-code-auditor/SKILL.md) | Check real integration, dependency/API validity and operational failure paths before cosmetic debt. |
| [BadTechBandit vibe-audit](https://github.com/BadTechBandit/skills/blob/master/vibe-audit/SKILL.md) | Identify stack first; explicitly distinguish applicable, skipped and unverified surfaces. |
| [noobnooc WTF](https://github.com/noobnooc/agent/blob/main/skills/wtf/SKILL.md) | Include manifests/runtime configuration and tracked-versus-local distinctions in reconnaissance. |
| [GabrielYMC security-audit](https://github.com/GabrielYMC/security-audit/blob/main/SKILL.md) | Security findings need attacker prerequisites, reachable source-to-sink evidence and actual exposure. |
| [Trail of Bits variant-analysis](https://github.com/trailofbits/skills/blob/main/plugins/variant-analysis/skills/variant-analysis/SKILL.md) | Calibrate a search against a confirmed seed, broaden incrementally, then validate each candidate. |
| [Trail of Bits differential-review](https://github.com/trailofbits/skills/blob/main/plugins/differential-review/skills/differential-review/SKILL.md) | Use caller reachability, risk and relevant history to separate regressions from pre-existing defects. |
| [Trail of Bits insecure-defaults audit](https://github.com/trailofbits/skills/blob/main/plugins/insecure-defaults/commands/audit.md) | Refute candidates by tracing defaults through effective guards before reporting. This upstream artifact is a command, not a standalone skill. |
| [Superpowers reviewer template](https://github.com/obra/superpowers/blob/main/skills/requesting-code-review/code-reviewer.md) | Review against requirements while explicitly retaining unresolved or deferred behaviors and their reasons. |

## Deliberately not imported

- Aggregate health/production-readiness scores, automatic deploy approval and severity
  assigned solely from a missing checklist item: they overstate covered evidence.
- Universal auth/database/payment/CSP requirements: apply only to actual stack/exposure.
- File-length, argument-count and duplication thresholds treated as defects by themselves.
- Automatic installs, cloud scanners, PR comments, extra report folders or forced
  multi-agent pipelines: these have independent scope, tools and authorization costs.
- “No callers found” as proof of dead code, or “no tests” as automatic severity inflation.
- Mandatory fixes during audit; this skill remains read-only unless fixes are requested.

Anthropic's code-review plugin appeared in search results, but the referenced source
returned 404 during verification; it is not relied upon for this update. Search results
and directories were discovery aids; the listed primary files were actually read.
The mutation/property guidance extends our existing test-credibility route; no claim
is made that these external methods have been empirically benchmarked here.
