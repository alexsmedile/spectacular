#!/usr/bin/env bash
set -euo pipefail

die() { echo "refused: $*" >&2; exit 3; }

usage() {
  die "usage: install.sh <install|update|select|rollback|uninstall|recover> --prefix <absolute-path> [--source <release-dir>] [--runtime <codex|claude>] [--version <semver>] [--os <darwin|linux>] [--arch <amd64|arm64>]"
}

command_name="${1:-}"
[[ -n "$command_name" ]] || usage
shift

prefix=""
source_dir=""
runtime_name=""
release_version=""
target_os=""
target_arch=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --prefix) [[ $# -ge 2 ]] || usage; prefix="$2"; shift 2 ;;
    --source) [[ $# -ge 2 ]] || usage; source_dir="$2"; shift 2 ;;
    --runtime) [[ $# -ge 2 ]] || usage; runtime_name="$2"; shift 2 ;;
    --version) [[ $# -ge 2 ]] || usage; release_version="$2"; shift 2 ;;
    --os) [[ $# -ge 2 ]] || usage; target_os="$2"; shift 2 ;;
    --arch) [[ $# -ge 2 ]] || usage; target_arch="$2"; shift 2 ;;
    *) usage ;;
  esac
done

[[ -n "$prefix" ]] || usage
case "$prefix" in /*) ;; *) die "prefix must be an absolute path" ;; esac
[[ "$prefix" != "/" ]] || die "prefix cannot be filesystem root"

if [[ -n "${HOME:-}" && -d "$HOME" ]]; then
  mkdir -p "$prefix"
  prefix="$(cd "$prefix" && pwd -P)"
  home_root="$(cd "$HOME" && pwd -P)"
  [[ "$prefix" != "$home_root" ]] || die "prefix cannot be the user home"
else
  mkdir -p "$prefix"
  prefix="$(cd "$prefix" && pwd -P)"
fi

for protected_path in \
  "$prefix/bin" \
  "$prefix/bin/spectacular" \
  "$prefix/plugins" \
  "$prefix/plugins/spectacular" \
  "$prefix/share" \
  "$prefix/share/spectacular" \
  "$prefix/share/spectacular/install.receipt" \
  "$prefix/share/spectacular/uninstall.receipt"; do
  [[ ! -L "$protected_path" ]] || die "installation path cannot be a symbolic link: $protected_path"
done

checksum_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    die "sha256sum or shasum is required"
  fi
}

detect_platform() {
  [[ -n "$target_os" ]] || target_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  [[ -n "$target_arch" ]] || target_arch="$(uname -m)"
  case "$target_os" in darwin|linux) ;; *) die "unsupported operating system: $target_os" ;; esac
  case "$target_arch" in
    x86_64) target_arch="amd64" ;;
    aarch64) target_arch="arm64" ;;
    amd64|arm64) ;;
    *) die "unsupported architecture: $target_arch" ;;
  esac
}

resolve_artifact() {
  [[ -n "$source_dir" ]] || usage
  [[ -d "$source_dir" ]] || die "release source is not a directory"
  source_dir="$(cd "$source_dir" && pwd -P)"
  [[ -n "$release_version" ]] || release_version="$(sed -n '1p' "$source_dir/VERSION" 2>/dev/null || true)"
  [[ "$release_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]] || die "release version is invalid"
  detect_platform
  artifact_name="spectacular-v${release_version}-${target_os}-${target_arch}.tar.gz"
  artifact_path="$source_dir/$artifact_name"
  checksums_path="$source_dir/SHA256SUMS"
  [[ -f "$artifact_path" ]] || die "release artifact is missing: $artifact_name"
  [[ -f "$checksums_path" ]] || die "SHA256SUMS is missing"
  expected_checksum=""
  checksum_matches=0
  while read -r digest listed_name extra; do
    if [[ "$listed_name" == "$artifact_name" ]]; then
      [[ -z "${extra:-}" ]] || die "malformed checksum entry for $artifact_name"
      expected_checksum="$digest"
      checksum_matches=$((checksum_matches + 1))
    fi
  done < "$checksums_path"
  [[ "$checksum_matches" -eq 1 ]] || die "checksum manifest must contain exactly one entry for $artifact_name"
  [[ "$expected_checksum" =~ ^[0-9a-f]{64}$ ]] || die "checksum entry is invalid"
  actual_checksum="$(checksum_file "$artifact_path")"
  [[ "$actual_checksum" == "$expected_checksum" ]] || die "checksum mismatch for $artifact_name"
}

verify_inventory() {
  inventory="$(tar -tzf "$artifact_path")" || die "release archive is unreadable"
  [[ -n "$inventory" ]] || die "release archive is empty"
  detailed_inventory="$(tar -tvzf "$artifact_path")" || die "release archive metadata is unreadable"
  # Installation copies and executes the payload.  Treat every archive entry
  # as hostile until tar has shown it to be an ordinary file.  In particular,
  # names alone cannot distinguish a regular binary from a symlink to an
  # attacker-controlled executable outside the extraction root.
  while IFS= read -r detail; do
    [[ "${detail:0:1}" == "-" ]] || die "release archive contains a non-regular entry"
  done <<< "$detailed_inventory"
  duplicates="$(printf '%s\n' "$inventory" | LC_ALL=C sort | uniq -d)"
  [[ -z "$duplicates" ]] || die "release archive contains duplicate entry names"
  while IFS= read -r entry; do
    [[ "$entry" != /* && "$entry" != *'//'* && "$entry" != *'/./'* && "$entry" != *$'\n'* && "$entry" != *$'\r'* ]] || die "unsafe release entry: $entry"
    case "$entry" in
      spectacular/RELEASE.json|spectacular/VERSION|spectacular/bin/spectacular) ;;
      spectacular/plugins/spectacular/.codex-plugin/plugin.json) ;;
      spectacular/plugins/spectacular/.claude-plugin/plugin.json) ;;
      spectacular/plugins/spectacular/skills/spectacular/*) ;;
      *) die "unexpected release entry: $entry" ;;
    esac
    case "$entry" in
      *../*|*/tests/*|*/testdata/*|*/cli/*|*/requests/*) die "v1 or unsafe release entry: $entry" ;;
    esac
  done <<< "$inventory"
  for required in \
    spectacular/RELEASE.json \
    spectacular/VERSION \
    spectacular/bin/spectacular \
    spectacular/plugins/spectacular/.codex-plugin/plugin.json \
    spectacular/plugins/spectacular/.claude-plugin/plugin.json \
    spectacular/plugins/spectacular/skills/spectacular/SKILL.md \
    spectacular/plugins/spectacular/skills/spectacular/generated/mechanical-interface.json; do
    grep -Fxq "$required" <<< "$inventory" || die "required release entry is missing: $required"
  done
}

receipt_value() {
  local receipt="$1" key="$2"
  sed -n "s/^${key}=//p" "$receipt" | sed -n '1p'
}

install_release() {
  case "$runtime_name" in codex|claude) ;; *) die "runtime must be codex or claude" ;; esac
  if [[ "$command_name" == install && -e "$prefix/share/spectacular/install.receipt" ]]; then
    die "v2 is already installed; use update"
  fi
  if [[ "$command_name" == update && ! -e "$prefix/share/spectacular/install.receipt" ]]; then
    die "no installed v2 release to update"
  fi
  resolve_artifact
  verify_inventory

  temporary="$(mktemp -d "${TMPDIR:-/tmp}/spectacular-install.XXXXXX")"
  trap 'rm -rf "$temporary"' EXIT
  tar -xzf "$artifact_path" -C "$temporary"
  payload="$temporary/spectacular"
  [[ "$($payload/bin/spectacular --version)" == "spectacular $release_version" ]] || die "binary version does not match selected release"
  grep -Fq '"release_version": "'"$release_version"'"' "$payload/plugins/spectacular/skills/spectacular/generated/mechanical-interface.json" || die "generated interface version mismatch"
  grep -Fq 'version: '"$release_version" "$payload/plugins/spectacular/skills/spectacular/SKILL.md" || die "Skill version mismatch"
  grep -Fq '"version": "'"$release_version"'"' "$payload/plugins/spectacular/.$runtime_name-plugin/plugin.json" || die "runtime manifest version mismatch"

  transaction="$prefix/.spectacular-install-$$"
  [[ ! -e "$transaction" ]] || die "installation transaction already exists"
  mkdir -p "$transaction/bin" "$transaction/plugins/spectacular" "$transaction/share/spectacular"
  cp "$payload/bin/spectacular" "$transaction/bin/spectacular"
  chmod 0755 "$transaction/bin/spectacular"
  cp -R "$payload/plugins/spectacular/skills" "$transaction/plugins/spectacular/skills"
  mkdir -p "$transaction/plugins/spectacular/.$runtime_name-plugin"
  cp "$payload/plugins/spectacular/.$runtime_name-plugin/plugin.json" "$transaction/plugins/spectacular/.$runtime_name-plugin/plugin.json"

  backup_relative="none"
  if [[ -e "$prefix/bin/spectacular" || -e "$prefix/plugins/spectacular" || -e "$prefix/share/spectacular/install.receipt" ]]; then
    backup_relative="share/spectacular/backups/$(date -u +%Y%m%dT%H%M%SZ)-$$"
    backup="$prefix/$backup_relative"
    mkdir -p "$backup/bin" "$backup/plugins" "$backup/share"
    [[ ! -e "$prefix/bin/spectacular" ]] || cp "$prefix/bin/spectacular" "$backup/bin/spectacular"
    [[ ! -e "$prefix/plugins/spectacular" ]] || cp -R "$prefix/plugins/spectacular" "$backup/plugins/spectacular"
    [[ ! -e "$prefix/share/spectacular/install.receipt" ]] || cp "$prefix/share/spectacular/install.receipt" "$backup/share/install.receipt"
  fi

  cat > "$transaction/share/spectacular/install.receipt" <<EOF
schema=spectacular.install.v1
version=$release_version
runtime=$runtime_name
platform=$target_os/$target_arch
artifact=$artifact_name
checksum=$actual_checksum
backup=$backup_relative
EOF

  apply_install() {
    mkdir -p "$prefix/bin" "$prefix/plugins" "$prefix/share/spectacular" || return 1
    rm -f "$prefix/bin/spectacular" || return 1
    rm -rf "$prefix/plugins/spectacular" || return 1
    mv "$transaction/bin/spectacular" "$prefix/bin/spectacular" || return 1
    mv "$transaction/plugins/spectacular" "$prefix/plugins/spectacular" || return 1
    mv "$transaction/share/spectacular/install.receipt" "$prefix/share/spectacular/install.receipt" || return 1
  }
  restore_previous() {
    rm -f "$prefix/bin/spectacular" || true
    rm -rf "$prefix/plugins/spectacular" || true
    rm -f "$prefix/share/spectacular/install.receipt" || true
    if [[ "$backup_relative" != none && -d "$prefix/$backup_relative" ]]; then
      [[ ! -f "$prefix/$backup_relative/bin/spectacular" ]] || cp "$prefix/$backup_relative/bin/spectacular" "$prefix/bin/spectacular"
      [[ ! -d "$prefix/$backup_relative/plugins/spectacular" ]] || cp -R "$prefix/$backup_relative/plugins/spectacular" "$prefix/plugins/spectacular"
      [[ ! -f "$prefix/$backup_relative/share/install.receipt" ]] || cp "$prefix/$backup_relative/share/install.receipt" "$prefix/share/spectacular/install.receipt"
    fi
  }
  if ! apply_install; then
    restore_previous
    die "installation transaction failed; previous v2 installation was restored"
  fi
  rm -rf "$transaction"
  trap - EXIT
  rm -rf "$temporary"
  printf 'version=%s\nruntime=%s\nplatform=%s/%s\nartifact=%s\nchecksum=%s\nresult=installed-v2-only\n' \
    "$release_version" "$runtime_name" "$target_os" "$target_arch" "$artifact_name" "$actual_checksum"
}

rollback_release() {
  receipt="$prefix/share/spectacular/install.receipt"
  [[ -f "$receipt" ]] || die "no installed v2 release receipt"
  backup_relative="$(receipt_value "$receipt" backup)"
  [[ "$backup_relative" == share/spectacular/backups/* ]] || die "no rollback point for current installation"
  [[ "$backup_relative" != *..* ]] || die "rollback pointer is invalid"
  backup="$prefix/$backup_relative"
  [[ -d "$backup" ]] || die "rollback point is missing"
  rm -f "$prefix/bin/spectacular"
  rm -rf "$prefix/plugins/spectacular"
  [[ ! -f "$backup/bin/spectacular" ]] || { mkdir -p "$prefix/bin"; mv "$backup/bin/spectacular" "$prefix/bin/spectacular"; }
  [[ ! -d "$backup/plugins/spectacular" ]] || { mkdir -p "$prefix/plugins"; mv "$backup/plugins/spectacular" "$prefix/plugins/spectacular"; }
  if [[ -f "$backup/share/install.receipt" ]]; then
    mv "$backup/share/install.receipt" "$receipt"
  else
    rm -f "$receipt"
  fi
  rmdir "$backup/bin" "$backup/plugins" "$backup/share" "$backup" 2>/dev/null || true
  printf 'restored=%s\nresult=rolled-back\n' "$backup_relative"
}

uninstall_release() {
  receipt="$prefix/share/spectacular/install.receipt"
  [[ -f "$receipt" ]] || die "no installed v2 release receipt"
  [[ ! -e "$prefix/share/spectacular/uninstall.receipt" ]] || die "an uninstall recovery point already exists"
  recovery_relative="share/spectacular/uninstalled/$(date -u +%Y%m%dT%H%M%SZ)-$$"
  recovery="$prefix/$recovery_relative"
  mkdir -p "$recovery/bin" "$recovery/plugins" "$recovery/share"
  [[ ! -f "$prefix/bin/spectacular" ]] || mv "$prefix/bin/spectacular" "$recovery/bin/spectacular"
  [[ ! -d "$prefix/plugins/spectacular" ]] || mv "$prefix/plugins/spectacular" "$recovery/plugins/spectacular"
  mv "$receipt" "$recovery/share/install.receipt"
  printf 'schema=spectacular.uninstall.v1\nrecovery=%s\n' "$recovery_relative" > "$prefix/share/spectacular/uninstall.receipt"
  printf 'recovery=%s\nresult=uninstalled-recoverable\n' "$recovery_relative"
}

recover_release() {
  uninstall_receipt="$prefix/share/spectacular/uninstall.receipt"
  [[ -f "$uninstall_receipt" ]] || die "no uninstall recovery point"
  recovery_relative="$(receipt_value "$uninstall_receipt" recovery)"
  [[ "$recovery_relative" == share/spectacular/uninstalled/* ]] || die "uninstall recovery pointer is invalid"
  [[ "$recovery_relative" != *..* ]] || die "uninstall recovery pointer is invalid"
  recovery="$prefix/$recovery_relative"
  [[ -d "$recovery" ]] || die "uninstall recovery point is missing"
  [[ ! -e "$prefix/bin/spectacular" && ! -e "$prefix/plugins/spectacular" && ! -e "$prefix/share/spectacular/install.receipt" ]] || die "current installation blocks recovery"
  [[ ! -f "$recovery/bin/spectacular" ]] || { mkdir -p "$prefix/bin"; mv "$recovery/bin/spectacular" "$prefix/bin/spectacular"; }
  [[ ! -d "$recovery/plugins/spectacular" ]] || { mkdir -p "$prefix/plugins"; mv "$recovery/plugins/spectacular" "$prefix/plugins/spectacular"; }
  mv "$recovery/share/install.receipt" "$prefix/share/spectacular/install.receipt"
  rm -f "$uninstall_receipt"
  rmdir "$recovery/bin" "$recovery/plugins" "$recovery/share" "$recovery" 2>/dev/null || true
  printf 'restored=%s\nresult=recovered\n' "$recovery_relative"
}

case "$command_name" in
  select)
    resolve_artifact
    verify_inventory
    printf 'version=%s\nplatform=%s/%s\nartifact=%s\nchecksum=%s\nresult=verified-local-artifact\n' \
      "$release_version" "$target_os" "$target_arch" "$artifact_path" "$actual_checksum"
    ;;
  install|update) install_release ;;
  rollback) rollback_release ;;
  uninstall) uninstall_release ;;
  recover) recover_release ;;
  *) usage ;;
esac
