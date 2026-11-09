#!/usr/bin/env bash
set -euo pipefail

# Usage: release-body.sh <tag> <asset1> [asset2 ...]
# Generates a markdown release body with a platform table for the given assets.
# Output is written to stdout.

tag="$1"
shift
assets=("$@")

repo_url="${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}"

if [ ${#assets[@]} -eq 0 ] && [ -z "${BINDINGS_COMMIT:-}" ]; then
  exit 0
fi

if [ -n "${BINDINGS_COMMIT:-}" ]; then
  echo "**Bindings commit:** [\`${BINDINGS_COMMIT:0:7}\`](${repo_url}/commit/${BINDINGS_COMMIT}) (on \`bindings\` branch)"
  echo
fi

if [ ${#assets[@]} -eq 0 ]; then
  exit 0
fi

cat <<'HEADER'
## 📦 Assets by Platform

| Language | Platform | Architecture | Asset | Bindings |
| :--- | :--- | :--- | :--- | :--- |
HEADER

# Extract crate name and version from the tag (crates/{name}/vX.Y.Z or crates/ffi/{name}/vX.Y.Z)
version=$(echo "$tag" | grep -o 'v[0-9].*')
crate_dir=$(echo "$tag" | sed "s|/${version}$||" | xargs basename)

# Build bindings source links using the ffi/ and bindings/go/ tags
repo_us=$(echo "${REPO_NAME:-}" | tr '-' '_')
suffix_us=$(echo "${PKG_SUFFIX:-}" | tr '-' '_')
crate_us=$(echo "$crate_dir" | tr '-' '_')
ffi_tag="ffi/${crate_dir}/${version}"
go_pkg="${repo_us}_${crate_us}${suffix_us}"
go_tag="bindings/go/${go_pkg}/${version}"

py_link="[source](${repo_url}/tree/${ffi_tag}/bindings/python)"
go_link="[source](${repo_url}/tree/${go_tag}/bindings/go/${go_pkg})"
js_link="[source](${repo_url}/tree/${ffi_tag}/bindings/js)"

for asset in "${assets[@]}"; do
  fname=$(basename "$asset")
  link="[${fname}](${repo_url}/releases/download/${tag}/${fname})"

  if echo "$fname" | grep -qi 'manylinux.*x86_64'; then
    echo "| **Python** | Linux | x86_64 | ${link} | ${py_link} |"
  elif echo "$fname" | grep -qi 'manylinux.*aarch64'; then
    echo "| **Python** | Linux | aarch64 | ${link} | ${py_link} |"
  elif echo "$fname" | grep -qi 'macosx.*arm64'; then
    echo "| **Python** | macOS | Apple Silicon | ${link} | ${py_link} |"
  elif echo "$fname" | grep -qi 'macosx.*x86_64'; then
    echo "| **Python** | macOS | Intel | ${link} | ${py_link} |"
  elif echo "$fname" | grep -qi 'win.*amd64\|windows.*x86_64'; then
    echo "| **Python** | Windows | x64 | ${link} | ${py_link} |"
  elif echo "$fname" | grep -qi '\.tgz$'; then
    echo "| **JS/TS** | All | All | ${link} | ${js_link} |"
  else
    echo "| — | — | — | ${link} | — |"
  fi
done

go_get="\`go get ${repo_url#https://}/bindings/go/${go_pkg}@${version}\`"
echo "| **Go** | All | All | ${go_get} | ${go_link} |"
