#!/usr/bin/env bash
# Shared logging helpers for scripts and hooks.
# Source this file: . "$(dirname "$0")/log.sh"  (or appropriate path)
#
# Emojis always prefix the message. GitHub Actions workflow commands
# (::error::, ::warning::, ::notice::, ::debug::) are prepended only
# when GITHUB_ACTIONS is set.

_log() {
  local level="$1" emoji="$2"
  shift 2

  local file=""
  if [ "${1:-}" = "--file" ]; then
    file="$2"; shift 2
  fi

  local msg="$*"

  if [ -n "${GITHUB_ACTIONS:-}" ]; then
    local attrs=""
    [ -n "$file" ] && attrs=" file=${file}"
    echo "::${level}${attrs}::${emoji}  ${msg}"
  else
    echo "${emoji}  ${msg}"
  fi
}

log_error()   { _log error   "❌" "$@"; }
log_warning() { _log warning "⚠️"  "$@"; }
log_notice()  { _log notice  "ℹ️"  "$@"; }
log_debug()   { _log debug   "🔍" "$@"; }
log_success() { echo "✅ $*"; }

log_group() {
  if [ -n "${GITHUB_ACTIONS:-}" ]; then
    echo "::group::$*"
  fi
}

log_endgroup() {
  if [ -n "${GITHUB_ACTIONS:-}" ]; then
    echo "::endgroup::"
  fi
}
