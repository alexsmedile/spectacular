#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
mode="${1:-all}"
go_cache="${GOCACHE:-$(mktemp -d)}"

export GOCACHE="$go_cache"
export GOPROXY=off
export GOFLAGS=-mod=readonly

static_checks() {
  formatting="$(cd "$repo_root" && gofmt -l cmd internal acceptance)"
  if [[ -n "$formatting" ]]; then
    echo "gofmt required:" >&2
    echo "$formatting" >&2
    exit 1
  fi
  (cd "$repo_root" && go mod verify)
  (cd "$repo_root" && go vet ./...)
  (cd "$repo_root" && git diff --check)
}

quick_checks() {
  static_checks
  (cd "$repo_root" && go test ./cmd/... ./install/... ./internal/...)
}

acceptance_checks() {
  (cd "$repo_root" && go test -count=1 ./acceptance)
}

release_checks() {
  (cd "$repo_root" && bash install/test.sh)
  (cd "$repo_root" && bash release/test.sh)
}

case "$mode" in
  quick)
    quick_checks
    ;;
  acceptance)
    static_checks
    acceptance_checks
    ;;
  release)
    static_checks
    release_checks
    ;;
  all)
    quick_checks
    acceptance_checks
    (cd "$repo_root" && go test -race ./...)
    (cd "$repo_root" && go build ./...)
    release_checks
    ;;
  *)
    echo "usage: bash test/verify.sh [quick|acceptance|release|all]" >&2
    exit 2
    ;;
esac

echo "Spectacular $mode verification: PASS"
