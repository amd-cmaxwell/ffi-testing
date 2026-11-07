#!/usr/bin/env bash
set -euo pipefail

if [ "${1:-}" = "--cached" ]; then
  changed_files=$(git diff --cached --name-only)
else
  base_ref="${1:-origin/main}"
  head_ref="${2:-HEAD}"
  changed_files=$(git diff --name-only "$base_ref"..."$head_ref")
fi

missing=()

for file in $changed_files; do
  case "$file" in
    crates/ffi/*) continue ;;
    crates/*/src/*)
      crate_name=$(echo "$file" | sed 's|crates/\([^/]*\)/.*|\1|')
      ffi_dir="crates/ffi/$crate_name"
      [ -d "$ffi_dir" ] || continue

      if ! echo "$changed_files" | grep -q "^crates/ffi/${crate_name}/"; then
        missing+=("$crate_name")
      fi
      ;;
  esac
done

# Deduplicate
IFS=$'\n' missing=($(printf '%s\n' "${missing[@]}" | sort -u))
unset IFS

if [ ${#missing[@]} -gt 0 ]; then
  echo "::warning::Library crates changed without corresponding FFI crate updates:"
  for name in "${missing[@]}"; do
    echo "  - crates/$name/src/ changed but crates/ffi/$name/ was not modified"
  done
  exit 1
fi
