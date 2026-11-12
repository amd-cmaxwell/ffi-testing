#!/usr/bin/env bash
set -euo pipefail

# Usage: check-ffi-versions.sh [--cached] [ffi_lib1 ffi_lib2 ...]
# If no arguments are given, all crates under crates/ffi/ are checked.
# FFI lib names use underscores (e.g. "core trig linear_algebra").
#
# --cached  Read versions from the git staging area instead of the working tree.
#           Only FFI crates with a staged Cargo.toml (or whose lib crate has one)
#           are checked.

. "$(dirname "$0")/log.sh"

use_cached=false
if [ "${1:-}" = "--cached" ]; then
  use_cached=true; shift
fi

staged_files=""
if $use_cached; then
  staged_files=$(git diff --cached --name-only)
fi

read_version() {
  if $use_cached; then
    git show :"$1" 2>/dev/null | grep '^version' | head -1 | sed 's/.*"\(.*\)"/\1/'
  else
    grep '^version' "$1" | head -1 | sed 's/.*"\(.*\)"/\1/'
  fi
}

log_group "FFI version alignment check"
failed=false
checked=0
skipped=0

if [ $# -gt 0 ]; then
  tomls=()
  for mod in "$@"; do
    crate_dir=$(echo "$mod" | tr '_' '-')
    tomls+=("crates/ffi/${crate_dir}/Cargo.toml")
  done
else
  tomls=(crates/ffi/*/Cargo.toml)
fi

for ffi_toml in "${tomls[@]}"; do
  [ -f "$ffi_toml" ] || continue
  crate_dir=$(basename "$(dirname "$ffi_toml")")
  lib_toml="crates/lib/${crate_dir}/Cargo.toml"
  if [ ! -f "$lib_toml" ]; then
    log_notice "Skipping Lib/FFI version parity check for crates/ffi/${crate_dir} — no corresponding lib crate"
    skipped=$((skipped + 1))
    continue
  fi

  if $use_cached; then
    ffi_staged=false; lib_staged=false
    echo "$staged_files" | grep -q "^${ffi_toml}$" && ffi_staged=true
    echo "$staged_files" | grep -q "^${lib_toml}$" && lib_staged=true
    $ffi_staged || $lib_staged || continue
  fi

  ffi_ver=$(read_version "$ffi_toml")
  lib_ver=$(read_version "$lib_toml")
  IFS='.' read -r ffi_major ffi_minor ffi_patch <<< "$ffi_ver"
  IFS='.' read -r lib_major lib_minor lib_patch <<< "$lib_ver"
  if [ "$ffi_major" -lt "$lib_major" ] ||
     { [ "$ffi_major" -eq "$lib_major" ] && [ "$ffi_minor" -lt "$lib_minor" ]; } ||
     { [ "$ffi_major" -eq "$lib_major" ] && [ "$ffi_minor" -eq "$lib_minor" ] && [ "$ffi_patch" -lt "$lib_patch" ]; }; then
    log_error --file "${ffi_toml}" "crates/ffi/${crate_dir} v${ffi_ver} is behind crates/lib/${crate_dir} v${lib_ver} — FFI version must be >= lib version"
    failed=true
  else
    log_success "crates/ffi/${crate_dir} v${ffi_ver} >= crates/lib/${crate_dir} v${lib_ver}"
  fi
  checked=$((checked + 1))
done
log_endgroup
[ $use_cached ] && maybe_staged="staged "
log_notice "Checked ${checked} ${maybe_staged}FFI/lib pairs, unconditionally skipped ${skipped} FFI-only crates"
if $failed; then
  log_error "One or more FFI crate versions are behind their lib crates"
  exit 1
fi
