#!/usr/bin/env bash
set -euo pipefail

# Usage: sync-workspace-versions.sh [--cached] [ffi_lib1 ffi_lib2 ...]
# If arguments are given, syncs those FFI crates and their corresponding lib crates.
# FFI lib names use underscores (e.g. "core trig linear_algebra").
# If no arguments are given, all crates under crates/lib/ and crates/ffi/ are synced.
# Bumps workspace dependency versions in the root Cargo.toml to match each
# crate's own [package].version.
#
# --cached  Read crate versions from the git staging area instead of the working
#           tree. Only crates with a staged Cargo.toml are synced. The root
#           Cargo.toml is still written on disk (the caller is expected to stage it).

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

updated=false
for toml in "${tomls[@]}"; do
  [ -f "$toml" ] || continue
  crate_path=$(dirname "$toml")

  if $use_cached; then
    echo "$staged_files" | grep -q "^${toml}$" || continue
    crate_ver=$(git show :"${toml}" 2>/dev/null | grep '^version' | head -1 | sed 's/.*"\(.*\)"/\1/')
  else
    crate_ver=$(grep '^version' "$toml" | head -1 | sed 's/.*"\(.*\)"/\1/')
  fi

  ws_line=$(grep "path = \"${crate_path}\"" Cargo.toml) || continue
  ws_ver=$(echo "$ws_line" | sed 's/.*version = "\([^"]*\)".*/\1/')

  if [ "$crate_ver" != "$ws_ver" ]; then
    log_notice "${crate_path}: ${ws_ver} -> ${crate_ver}"
    sed "s|\(path = \"${crate_path}\", version = \)\"${ws_ver}\"|\1\"${crate_ver}\"|" Cargo.toml > Cargo.toml.tmp && mv Cargo.toml.tmp Cargo.toml
    updated=true
  fi
done

if $updated; then
  echo "updated"
else
  log_success "all workspace dependency versions already match"
fi
