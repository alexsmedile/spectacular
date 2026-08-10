# Contributor guide

This is Spectacular v2. The root Go module, `cmd/spectacular`, `skills/`,
`install/`, and `.spectacular/` are the only live product surface.

Run `gofmt -w`, `GOPROXY=off GOFLAGS=-mod=readonly go test ./...`, and
`bash release/test.sh` before release changes. Do not reintroduce v1 commands,
compatibility readers, migrations, generic record/search verbs, or a second
package root. Keep release version values aligned through `VERSION` and the
generated mechanical interface.
