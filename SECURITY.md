# Security policy

## Reporting a vulnerability

Report vulnerabilities privately through
[GitHub Security Advisories](https://github.com/alexsmedile/spectacular/security/advisories/new)
or the repository's private vulnerability-reporting channel. Do not disclose a
potential vulnerability publicly before maintainers can coordinate a response.

Please include the version you tested (`spectacular --version`), the platform,
and the smallest reproduction you have.

## Supported versions

The current 2.x release line is supported. See [`VERSION`](VERSION) for the
current version and [`CHANGELOG.md`](CHANGELOG.md) for what each release
changed. Earlier lines, including the v1 series and the `2.0.0-rc.*`
pre-releases, are unsupported and should not be installed.

## Verifying a release

Verify the SHA-256 checksum in `SHA256SUMS` before installing a release
archive. The installer refuses an archive whose selected checksum does not
match, so a failed verification is a refusal rather than a warning.

## Scope

Spectacular reads and writes Markdown records in a workspace you control and
shells out to `git`. It does not fetch remote code at runtime, and the
installer does not download binaries on your behalf. Reports that involve an
agent being granted authority it was never given by an owner gate are in
scope; reports that require an attacker to already have write access to your
workspace or your `git` history generally are not.
