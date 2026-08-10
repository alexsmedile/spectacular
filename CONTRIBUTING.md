# Contributing

Use the root v2 module only. Before a release-facing change, run:

```sh
gofmt -w $(find cmd internal -name '*.go')
GOPROXY=off GOFLAGS=-mod=readonly go mod verify
GOPROXY=off GOFLAGS=-mod=readonly go test ./...
GOPROXY=off GOFLAGS=-mod=readonly go test -race ./...
bash install/test.sh
bash release/test.sh
```

`internal/command.Registry` is the single source for the generated mechanical
interface. Regenerate it when changing the registry; do not hand-maintain a
command inventory.

Use conventional commit messages. Do not reintroduce v1 compatibility,
migrations, generic record/search verbs, or a second package root.
