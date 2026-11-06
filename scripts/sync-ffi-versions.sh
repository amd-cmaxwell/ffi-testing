#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CRATES_DIR="$REPO_ROOT/crates"
FFI_DIR="$CRATES_DIR/ffi"

changed=()

for lib_toml in "$CRATES_DIR"/*/Cargo.toml; do
  crate_dir="$(basename "$(dirname "$lib_toml")")"

  # Skip the ffi directory itself
  [ "$crate_dir" = "ffi" ] && continue

  ffi_toml="$FFI_DIR/$crate_dir/Cargo.toml"
  [ -f "$ffi_toml" ] || continue

  lib_version="$(grep '^version' "$lib_toml" | head -1 | sed 's/.*"\(.*\)"/\1/')"
  ffi_version="$(grep '^version' "$ffi_toml" | head -1 | sed 's/.*"\(.*\)"/\1/')"

  if [ "$lib_version" != "$ffi_version" ]; then
    sed -i "s/^version = \"$ffi_version\"/version = \"$lib_version\"/" "$ffi_toml"
    crate_us="$(echo "$crate_dir" | tr '-' '_')"
    changed+=("$crate_us")
    echo >&2 "  synced $crate_dir: $ffi_version -> $lib_version"
  fi
done

if [ ${#changed[@]} -gt 0 ]; then
  echo "${changed[*]}"
fi
