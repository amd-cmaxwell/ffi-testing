#!/usr/bin/env bash
set -euo pipefail

# Usage: ffi-crate-name.sh [--us|--hy] <crate-dir>
#   --hy (default): hyphenated form (normalized to hyphens)
#   --us:           underscore form
#
# Reads the crate name from crates/ffi/<crate-dir>/Cargo.toml and prints it
# in the requested normalization.

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

format="hy"
while [[ "${1:-}" == --* ]]; do
  case "$1" in
    --us) format="us"; shift ;;
    --hy) format="hy"; shift ;;
    *) echo "error: unknown flag $1" >&2; exit 1 ;;
  esac
done

crate_dir="${1:?usage: ffi-crate-name.sh [--us|--hy] <crate-dir>}"
cargo_toml="${REPO_ROOT}/crates/ffi/${crate_dir}/Cargo.toml"

if [ ! -f "$cargo_toml" ]; then
  echo "error: $cargo_toml not found" >&2
  exit 1
fi

name=$(grep '^name' "$cargo_toml" | head -1 | sed 's/.*"\(.*\)"/\1/')

case "$format" in
  hy) echo "$name" | tr '_' '-' ;;
  us) echo "$name" | tr '-' '_' ;;
esac
