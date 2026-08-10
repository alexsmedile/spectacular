#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
v2_root="$(cd "$script_dir/.." && pwd)"
test_root="$(mktemp -d)"
first="$test_root/first"
second="$test_root/second"
go_cache="$test_root/go-cache"
commit="scenario-r-reproducibility-proof"
clean_first="$test_root/clean-v2-first"
clean_second="$test_root/clean-v2-second"

mkdir -p "$go_cache"
cp -R "$v2_root" "$clean_first"
cp -R "$v2_root" "$clean_second"
(cd "$clean_first" && GOPROXY=off GOCACHE="$go_cache" GOFLAGS=-mod=readonly go run ./cmd/assemble-release --output "$first" --commit "$commit")
(cd "$clean_second" && GOPROXY=off GOCACHE="$go_cache" GOFLAGS=-mod=readonly go run ./cmd/assemble-release --output "$second" --commit "$commit")
release_root="$clean_first"

cmp "$first/VERSION" "$second/VERSION"
cmp "$first/SHA256SUMS" "$second/SHA256SUMS"
while read -r digest artifact; do
  [[ "$digest" =~ ^[0-9a-f]{64}$ ]]
  cmp "$first/$artifact" "$second/$artifact"
done < "$first/SHA256SUMS"

for target_os in darwin linux; do
  for target_arch in amd64 arm64; do
    prefix="$test_root/select-$target_os-$target_arch"
    "$release_root/install/install.sh" select --prefix "$prefix" --source "$first" --runtime codex --os "$target_os" --arch "$target_arch" | grep -q 'result=verified-local-artifact'
  done
done

case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
  darwin|linux) host_os="$(uname -s | tr '[:upper:]' '[:lower:]')" ;;
  *) echo "unsupported test host" >&2; exit 3 ;;
esac
case "$(uname -m)" in
  x86_64) host_arch="amd64" ;;
  aarch64|arm64) host_arch="arm64" ;;
  *) echo "unsupported test architecture" >&2; exit 3 ;;
esac

for runtime_name in codex claude; do
  prefix="$test_root/$runtime_name"
  "$release_root/install/install.sh" install --prefix "$prefix" --source "$first" --runtime "$runtime_name" | grep -q 'result=installed-v2-only'
  "$prefix/bin/spectacular" --version | grep -q "spectacular $(sed -n '1p' "$release_root/VERSION")"
  "$prefix/bin/spectacular" --version --json | grep -q '"schema_version":"spectacular.build-info.v1"'
  [[ -f "$prefix/plugins/spectacular/.$runtime_name-plugin/plugin.json" ]]
  if [[ "$runtime_name" == codex ]]; then
    [[ ! -e "$prefix/plugins/spectacular/.claude-plugin" ]]
  else
    [[ ! -e "$prefix/plugins/spectacular/.codex-plugin" ]]
  fi
  (cd "$release_root/testdata/scenario-a" && "$prefix/bin/spectacular" workspace context project --event @Orient --json) | grep -q 'spectacular.context.v1'
  "$release_root/install/install.sh" update --prefix "$prefix" --source "$first" --runtime "$runtime_name" | grep -q 'result=installed-v2-only'
  "$release_root/install/install.sh" rollback --prefix "$prefix" | grep -q 'result=rolled-back'
  "$prefix/bin/spectacular" --version >/dev/null
  "$release_root/install/install.sh" uninstall --prefix "$prefix" | grep -q 'result=uninstalled-recoverable'
  [[ ! -e "$prefix/bin/spectacular" ]]
  "$release_root/install/install.sh" recover --prefix "$prefix" | grep -q 'result=recovered'
  "$prefix/bin/spectacular" --version >/dev/null
  smoke_workspace="$test_root/smoke-$runtime_name"
  mkdir -p "$smoke_workspace"
  (cd "$release_root" && GOPROXY=off GOCACHE="$go_cache" GOFLAGS=-mod=readonly go run ./cmd/release-smoke \
    --binary "$prefix/bin/spectacular" --fixture testdata/scenario-b-c --workspace "$smoke_workspace") | \
    grep -q 'result=installed-binary-governed-closure-and-cold-resume-pass'
done

corrupt_source="$test_root/corrupt-source"
cp -R "$first" "$corrupt_source"
corrupt_artifact="spectacular-v$(sed -n '1p' "$release_root/VERSION")-$host_os-$host_arch.tar.gz"
printf 'corrupt' >> "$corrupt_source/$corrupt_artifact"
safe_prefix="$test_root/checksum-refusal"
mkdir -p "$safe_prefix"
printf 'preserve\n' > "$safe_prefix/sentinel"
before="$(shasum -a 256 "$safe_prefix/sentinel")"
if "$release_root/install/install.sh" install --prefix "$safe_prefix" --source "$corrupt_source" --runtime codex >/dev/null 2>&1; then
  echo "checksum mismatch unexpectedly installed" >&2
  exit 1
fi
[[ "$before" == "$(shasum -a 256 "$safe_prefix/sentinel")" ]]
[[ -z "$(find "$safe_prefix" -mindepth 1 ! -name sentinel -print -quit)" ]]

# A correctly checksummed archive is still untrusted.  This fixture keeps the
# allowlisted name inventory but substitutes the executable with an absolute
# symlink; installation must reject it before it writes anything below prefix.
hostile_source="$test_root/hostile-source"
cp -R "$first" "$hostile_source"
hostile_stage="$test_root/hostile-stage"
mkdir -p "$hostile_stage/spectacular/bin"
tar -xzf "$first/$corrupt_artifact" -C "$hostile_stage"
rm "$hostile_stage/spectacular/bin/spectacular"
ln -s "$test_root/outside-binary" "$hostile_stage/spectacular/bin/spectacular"
hostile_entries=()
while IFS= read -r entry; do
  hostile_entries+=("$entry")
done < <(tar -tzf "$first/$corrupt_artifact")
tar -czf "$hostile_source/$corrupt_artifact" -C "$hostile_stage" "${hostile_entries[@]}"
hostile_digest="$(shasum -a 256 "$hostile_source/$corrupt_artifact" | awk '{print $1}')"
awk -v name="$corrupt_artifact" -v digest="$hostile_digest" '$2 == name {$1 = digest} {print $1 "  " $2}' "$hostile_source/SHA256SUMS" > "$hostile_source/SHA256SUMS.new"
mv "$hostile_source/SHA256SUMS.new" "$hostile_source/SHA256SUMS"
hostile_prefix="$test_root/hostile-refusal"
mkdir -p "$hostile_prefix"
printf 'preserve\n' > "$hostile_prefix/sentinel"
if "$release_root/install/install.sh" install --prefix "$hostile_prefix" --source "$hostile_source" --runtime codex >/dev/null 2>&1; then
  echo "symlink archive unexpectedly installed" >&2
  exit 1
fi
[[ -f "$hostile_prefix/sentinel" ]]
[[ -z "$(find "$hostile_prefix" -mindepth 1 ! -name sentinel -print -quit)" ]]

unsupported_prefix="$test_root/unsupported-refusal"
mkdir -p "$unsupported_prefix"
printf 'preserve\n' > "$unsupported_prefix/sentinel"
if "$release_root/install/install.sh" select --prefix "$unsupported_prefix" --source "$first" --runtime codex --os windows --arch amd64 >/dev/null 2>&1; then
  echo "unsupported platform unexpectedly selected" >&2
  exit 1
fi
[[ -z "$(find "$unsupported_prefix" -mindepth 1 ! -name sentinel -print -quit)" ]]

outside="$test_root/outside"
symlink_prefix="$test_root/symlink-refusal"
mkdir -p "$outside" "$symlink_prefix"
ln -s "$outside" "$symlink_prefix/bin"
if "$release_root/install/install.sh" install --prefix "$symlink_prefix" --source "$first" --runtime codex >/dev/null 2>&1; then
  echo "symlinked install path unexpectedly accepted" >&2
  exit 1
fi
[[ -z "$(find "$outside" -mindepth 1 -print -quit)" ]]

dependencies="$(cd "$release_root" && GOPROXY=off GOCACHE="$go_cache" GOFLAGS=-mod=readonly go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ./cmd/spectacular)"
if grep -Fxq 'github.com/alexsmedile/spectacular' <<< "$dependencies"; then
  echo "v1 module dependency leaked into v2" >&2
  exit 1
fi

for archive in "$first"/*.tar.gz; do
  inventory="$(tar -tzf "$archive")"
  if grep -E '(^|/)(cli|tests|testdata|requests|archive)/|(^|/)skills/spectacular/versions/' <<< "$inventory"; then
    echo "v1 runtime/test path leaked into $archive" >&2
    exit 1
  fi
done

echo "Scenario R reproducible release/install/recovery proof: PASS"
