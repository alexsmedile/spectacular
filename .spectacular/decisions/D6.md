# D6 — Remove docs * verbs + deprecation_notice() + docs-* refs as MINOR (v1.17.0) — banner-warned since v1.2.0, pageworks is t

**Context:**
docs * surface (init/export/new/review/status) was extracted to pageworks in v1.2.0 with in-product deprecation banners pointing at v2.0.0 removal. Five releases later it is still in the CLI. Strict SemVer would call verb removal MAJOR; the justification for MINOR is that the removal was announced in-product continuously, pageworks is the documented replacement, and no current undeprecated surface changes behavior.

**Decision:**
Remove docs * verbs + deprecation_notice() + docs-* refs as MINOR (v1.17.0) — banner-warned since v1.2.0, pageworks is the documented replacement, no current documented behavior changes

**Consequences:**
v2.0.0 major is left with a single breaking concern (the file-contract change). The CLI shrinks by ~700 lines. Any user still calling docs * gets a clear error pointing at pageworks instead of a deprecation banner.
