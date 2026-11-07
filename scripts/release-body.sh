#!/usr/bin/env bash
set -euo pipefail

# Usage: release-body.sh <tag> <asset1> [asset2 ...]
# Generates a markdown release body with a platform table for the given assets.
# Output is written to stdout.

tag="$1"
shift
assets=("$@")

repo_url="${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}"

if [ ${#assets[@]} -eq 0 ]; then
  exit 0
fi

cat <<'HEADER'
## 📦 Assets by Platform

| Language | Platform | Architecture | Asset |
| :--- | :--- | :--- | :--- |
HEADER

for asset in "${assets[@]}"; do
  fname=$(basename "$asset")
  link="[${fname}](${repo_url}/releases/download/${tag}/${fname})"

  if echo "$fname" | grep -qi 'manylinux.*x86_64'; then
    echo "| **Python** | Linux | x86_64 | ${link} |"
  elif echo "$fname" | grep -qi 'manylinux.*aarch64'; then
    echo "| **Python** | Linux | aarch64 | ${link} |"
  elif echo "$fname" | grep -qi 'macosx.*arm64'; then
    echo "| **Python** | macOS | Apple Silicon | ${link} |"
  elif echo "$fname" | grep -qi 'macosx.*x86_64'; then
    echo "| **Python** | macOS | Intel | ${link} |"
  elif echo "$fname" | grep -qi 'win.*amd64\|windows.*x86_64'; then
    echo "| **Python** | Windows | x64 | ${link} |"
  elif echo "$fname" | grep -qi '\.tgz$'; then
    echo "| **JS/TS** | All | All | ${link} |"
  else
    echo "| — | — | — | ${link} |"
  fi
done
