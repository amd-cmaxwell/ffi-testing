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

  IFS='.' read -r lib_major lib_minor lib_patch <<< "$lib_version"
  IFS='.' read -r ffi_major ffi_minor ffi_patch <<< "$ffi_version"

  # Only sync if the lib version is greater than or equal to the FFI version
  lib_gte=false
  if [ "$lib_major" -gt "$ffi_major" ] ||
     { [ "$lib_major" -eq "$ffi_major" ] && [ "$lib_minor" -gt "$ffi_minor" ]; } ||
     { [ "$lib_major" -eq "$ffi_major" ] && [ "$lib_minor" -eq "$ffi_minor" ] && [ "$lib_patch" -ge "$ffi_patch" ]; }; then
    lib_gte=true
  fi

  if ! $lib_gte; then
    echo >&2 "  warning: crates/ffi/$crate_dir v${ffi_version} is ahead of crates/$crate_dir v${lib_version} — keeping the newer version"
  elif [ "$lib_version" != "$ffi_version" ]; then
    # Both files must be updated:
    # - crates/ffi/*/Cargo.toml [package].version: the crate's own version (not inheritable from workspace)
    # - Cargo.toml [workspace.dependencies]: the version other crates pull via dep.workspace = true
    ffi_path="crates/ffi/${crate_dir}"
    sed "s/^version = \"${ffi_version}\"/version = \"${lib_version}\"/" "$ffi_toml" > "${ffi_toml}.tmp" && mv "${ffi_toml}.tmp" "$ffi_toml"
    sed "s|\(path = \"${ffi_path}\", version = \)\"${ffi_version}\"|\1\"${lib_version}\"|" "$REPO_ROOT/Cargo.toml" > "$REPO_ROOT/Cargo.toml.tmp" && mv "$REPO_ROOT/Cargo.toml.tmp" "$REPO_ROOT/Cargo.toml"
    crate_us="$(echo "$crate_dir" | tr '-' '_')"
    changed+=("$crate_us")
    echo >&2 "  synced $crate_dir: $ffi_version -> $lib_version"
  fi
done

if [ ${#changed[@]} -gt 0 ]; then
  echo "${changed[*]}"
fi
