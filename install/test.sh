#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
v2_root="$(cd "$script_dir/.." && pwd)"
test_root="$(mktemp -d)"
binary_path="$test_root/spectacular"
go_cache="$test_root/go-cache"
release_version="$(sed -n '1p' "$v2_root/VERSION")"

mkdir -p "$go_cache"
(cd "$v2_root" && GOCACHE="$go_cache" GOFLAGS=-mod=readonly go build -trimpath -buildvcs=false \
  -ldflags "-s -w -X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Version=$release_version -X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Commit=local-stage" \
  -o "$binary_path" ./cmd/spectacular)

for runtime_name in codex claude; do
  destination="$test_root/$runtime_name"
  mkdir -p "$destination"
  "$script_dir/stage-local.sh" "$runtime_name" "$destination" "$binary_path" >/dev/null
  [[ -x "$destination/bin/spectacular" ]]
  [[ -f "$destination/plugins/spectacular/skills/spectacular/SKILL.md" ]]
  [[ -f "$destination/plugins/spectacular/.$runtime_name-plugin/plugin.json" ]]
  (cd "$destination" && "$destination/bin/spectacular" mission show M1 --json >/dev/null 2>&1) && exit 1 || code=$?
  [[ "$code" -eq 3 ]]
  (cd "$v2_root/testdata/scenario-a" && "$destination/bin/spectacular" mission show M1 --json) | grep -q 'spectacular.mission.show.v2'
done

home_probe="$test_root/home-probe"
mkdir -p "$home_probe"
HOME="$home_probe" "$script_dir/stage-local.sh" codex "$home_probe" "$binary_path" >/dev/null 2>&1 && exit 1 || code=$?
[[ "$code" -eq 3 ]]
HOME="$home_probe" "$script_dir/stage-local.sh" codex "$home_probe/." "$binary_path" >/dev/null 2>&1 && exit 1 || code=$?
[[ "$code" -eq 3 ]]
[[ -z "$(find "$home_probe" -mindepth 1 -maxdepth 1 -print -quit)" ]]

echo "Scenario S disposable Codex/Claude staging: PASS"
