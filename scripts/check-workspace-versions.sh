#!/usr/bin/env bash
set -euo pipefail

# Usage: check-workspace-versions.sh [--cached] [ffi_lib1 ffi_lib2 ...]
# If arguments are given, checks those FFI crates and their corresponding lib crates.
# FFI lib names use underscores (e.g. "core trig linear_algebra").
# If no arguments are given, all crates under crates/lib/ and crates/ffi/ are checked.
# Exits non-zero if any workspace dependency version in the root Cargo.toml
# does not match the crate's own [package].version.
#
# --cached  Read versions from the git staging area instead of the working tree.
#           Only crates with a staged Cargo.toml are checked.

. "$(dirname "$0")/log.sh"

use_cached=false
if [ "${1:-}" = "--cached" ]; then
  use_cached=true; shift
fi

if [ $# -gt 0 ]; then
  tomls=()
  for mod in "$@"; do
    crate_dir=$(echo "$mod" | tr '_' '-')
    tomls+=("crates/lib/${crate_dir}/Cargo.toml")
    tomls+=("crates/ffi/${crate_dir}/Cargo.toml")
  done
else
  tomls=(crates/lib/*/Cargo.toml crates/ffi/*/Cargo.toml)
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

failed=false

log_group "Check workspace dependency versions"
for toml in "${tomls[@]}"; do
  [ -f "$toml" ] || continue
  crate_path=$(dirname "$toml")

  if $use_cached; then
    echo "$staged_files" | grep -q "^${toml}$" || continue
  fi

  crate_ver=$(read_version "$toml")
  if $use_cached; then
    ws_line=$(git show :Cargo.toml 2>/dev/null | grep "path = \"${crate_path}\"") || continue
  else
    ws_line=$(grep "path = \"${crate_path}\"" Cargo.toml) || continue
  fi
  ws_ver=$(echo "$ws_line" | sed 's/.*version = "\([^"]*\)".*/\1/')

  if [ "$crate_ver" != "$ws_ver" ]; then
    log_error --file Cargo.toml "Workspace dep version ${ws_ver} for ${crate_path} does not match crate version ${crate_ver}"
    failed=true
  else
    log_success "${crate_path} v${crate_ver} matches workspace dep"
  fi
done
log_endgroup

if $failed; then
  log_error "Workspace dependency versions must match crate versions — run ./scripts/sync-workspace-versions.sh to fix"
  exit 1
fi
