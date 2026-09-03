#!/bin/sh
set -eu

[ "$#" -eq 1 ] || { echo "usage: build-eval-cli.sh <absolute-output-path>" >&2; exit 2; }
output=$1
case "$output" in /*) ;; *) echo "output path must be absolute" >&2; exit 2 ;; esac

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
repo=$(CDPATH= cd "$script_dir/../../../.." && pwd)
case "$output" in "$repo"/*) echo "output must be outside the repository" >&2; exit 2 ;; esac

git -C "$repo" diff --quiet && git -C "$repo" diff --cached --quiet || {
  echo "repository must be clean so the binary is attributable to HEAD" >&2
  exit 3
}

version=$(sed -n '1p' "$repo/VERSION")
commit=$(git -C "$repo" rev-parse HEAD)
package=github.com/alexsmedile/spectacular/v2/internal/buildinfo

(cd "$repo" && go build -ldflags "-X $package.Version=$version -X $package.Commit=$commit" -o "$output" ./cmd/spectacular)
actual=$($output --version --json)
printf '%s\n' "$actual" | grep -q '"schema_version":"spectacular.build-info.v1"' || exit 3
printf '%s\n' "$actual" | grep -q '"version":"'"$version"'"' || exit 3
printf 'cli=%s\nversion=%s\ncommit=%s\n' "$output" "$version" "$commit"
