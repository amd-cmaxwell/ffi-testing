# CI & Release Workflows

This directory contains the GitHub Actions workflows that power continuous integration, multi-language binding generation, and automated releases for this UniFFI Rust FFI template repository.

## Architecture overview

The CI system is built around **reusable workflows**: the three `build-*-bindings.yml` files are standalone `workflow_call` workflows that are invoked by both `ci.yml` (on every push/PR) and `release.yml` (on manual dispatch). This avoids duplicating build logic across the two pipelines.

```
ci.yml (push / PR / dispatch)              release.yml (manual dispatch)
├── check-ffi-coverage (PR only)           ├── ci-gate
├── check-ffi-versions ─────────────┐      ├── detect-changes
├── check-workspace-versions ───────┤      ├── check-ffi-versions ───────────┐
├── rust-test (matrix: 3 OS) ───────┤      ├── check-workspace-versions ─────┤
│                                   │      │                                 │
├── build-python-bindings.yml ◄─────┘      ├── build-python-bindings.yml ◄───┘
├── build-go-bindings.yml ◄─────────┘      ├── build-go-bindings.yml ◄───────┘
├── build-js-bindings.yml ◄─────────┘      ├── build-js-bindings.yml ◄───────┘
│                                          │
├── test-python-bindings                   ├── publish-bindings
├── test-go-bindings                       ├── release-rust-crates
└── test-js-bindings                       └── rollback (on failure)
```

## Workflow files

### `ci.yml` — Continuous Integration

| | |
|---|---|
| **Triggers** | Push to `main`, pull request to `main`, manual dispatch |
| **Purpose** | Run tests, build all bindings, and smoke-test the generated artifacts |

**Jobs:**

| Job | Description |
|---|---|
| `check-ffi-coverage` | (PR only) Ensures changes to `crates/lib/*/src/` have corresponding `crates/ffi/` updates |
| `check-ffi-versions` | Verifies FFI crate versions are ≥ their corresponding lib crate versions |
| `check-workspace-versions` | Verifies workspace dependency versions in `Cargo.toml` match each crate's own version |
| `rust-test` | Runs `cargo test --workspace` on Linux, macOS, and Windows |
| `build-python-bindings` | Calls `build-python-bindings.yml` to generate Python wheels (3 OS matrix) |
| `build-go-bindings` | Calls `build-go-bindings.yml` to generate Go bindings (Linux + macOS) |
| `build-js-bindings` | Calls `build-js-bindings.yml` to generate WASM/TypeScript bindings (Linux) |
| `test-python-bindings` | Installs wheels into a venv and runs `examples/py-demo` |
| `test-go-bindings` | Downloads Go bindings artifact and runs `examples/go-demo` |
| `test-js-bindings` | Downloads npm tarballs and runs `examples/js-demo/main.mjs` |

---

### `release.yml` — Release

| | |
|---|---|
| **Triggers** | Manual dispatch only (`workflow_dispatch`) |
| **Inputs** | `crates` — space-separated lib crates to release; `ffi-libs` — space-separated FFI libs to release. Leave either empty for auto-detection. |
| **Permissions** | `contents: write` (creates tags and releases) |

**Jobs:**

| Job | Description |
|---|---|
| `ci-gate` | Verifies CI has passed on HEAD before proceeding |
| `detect-changes` | Compares current crate versions against the latest release tag to find changed crates |
| `check-ffi-versions` | Same version-parity check as CI, scoped to changed FFI crates |
| `check-workspace-versions` | Same workspace version check as CI, scoped to changed crates |
| `build-python-bindings` | Builds wheels for changed FFI crates |
| `build-go-bindings` | Builds Go bindings for changed FFI crates |
| `build-js-bindings` | Builds WASM/TS bindings for changed FFI crates |
| `publish-bindings` | Commits generated bindings to the `bindings` branch and creates tags |
| `release-rust-crates` | Creates GitHub Releases with attached wheel and tarball assets |
| `rollback` | On failure: deletes tags, releases, and resets the `bindings` branch |

See [Release flow](#release-flow) below for the detailed sequence.

---

### `build-python-bindings.yml` — Reusable: Python

| | |
|---|---|
| **Trigger** | `workflow_call` only |
| **Matrix** | `ubuntu-latest`, `macos-latest`, `windows-latest` |
| **Runs** | `make generate-python` via maturin + uniffi-bindgen |

**Inputs:**

| Input | Required | Default | Description |
|---|---|---|---|
| `repo-owner` | yes | — | GitHub org or user |
| `repo-name` | yes | — | Repository name |
| `upload-artifacts` | no | `false` | Upload wheels and binding sources as artifacts |
| `ffi-libs` | no | `"auto"` | Space-separated FFI crates to build |
| `python-version` | no | `"3.12"` | Python version for the build |

**Artifacts uploaded** (when `upload-artifacts: true`):

| Artifact | Contents |
|---|---|
| `python-bindings-{os}-{arch}` | Generated Python source under `bindings/python/` |
| `wheel-{runner-os}` | Compiled `.whl` files under `dist/pypi/` |

---

### `build-go-bindings.yml` — Reusable: Go

| | |
|---|---|
| **Trigger** | `workflow_call` only |
| **Matrix** | `ubuntu-latest` (linux/amd64), `macos-latest` (darwin/arm64) |
| **Runs** | `make generate-go` via uniffi-bindgen-go + CGO static libs |

**Inputs:**

| Input | Required | Default | Description |
|---|---|---|---|
| `repo-owner` | yes | — | GitHub org or user |
| `repo-name` | yes | — | Repository name |
| `upload-artifacts` | no | `false` | Upload generated Go modules as artifacts |
| `ffi-libs` | no | `"auto"` | Space-separated FFI crates to build |
| `go-version` | no | `"1.21"` | Go version for the build |
| `python-version` | no | `"3.12"` | Python version (used by Go binding generation tooling) |

**Artifacts uploaded:**

| Artifact | Contents |
|---|---|
| `go-bindings-{os}-{arch}` | Generated Go source + static libs under `bindings/go/` |

---

### `build-js-bindings.yml` — Reusable: JS/TS/WASM

| | |
|---|---|
| **Trigger** | `workflow_call` only |
| **Runs on** | `ubuntu-latest` only |
| **Runs** | `make generate-js` (default) via uniffi-bindgen-react-native + tsc + npm pack |

**Inputs:**

| Input | Required | Default | Description |
|---|---|---|---|
| `repo-owner` | yes | — | GitHub org or user |
| `repo-name` | yes | — | Repository name |
| `upload-artifacts` | no | `false` | Upload bindings and tarballs as artifacts |
| `make-target` | no | `generate-js` | Make target to invoke |
| `ffi-libs` | no | `"auto"` | Space-separated FFI crates to build |
| `node-version` | no | `"20"` | Node.js version |
| `pnpm-version` | no | `"9"` | pnpm version |
| `wasm-target` | no | `"wasm32-unknown-unknown"` | Rust WASM compilation target |

**Artifacts uploaded:**

| Artifact | Contents |
|---|---|
| `js-bindings` | Generated TS source + `.wasm` blobs under `bindings/js/` |
| `npm-tarballs` | Packed `.tgz` files under `dist/npm/` |

---

### `check-ffi-versions.yml` — Reusable: FFI version parity

| | |
|---|---|
| **Trigger** | `workflow_call` only |
| **Runs** | `scripts/check-ffi-versions.sh` |

Asserts that every FFI crate version is ≥ its corresponding lib crate version. Accepts an optional `ffi-libs` input to scope the check.

---

### `check-workspace-versions.yml` — Reusable: workspace dep versions

| | |
|---|---|
| **Trigger** | `workflow_call` only |
| **Runs** | `scripts/check-workspace-versions.sh` |

Verifies that workspace dependency versions in the root `Cargo.toml` match each crate's own `[package].version`. Accepts an optional `ffi-libs` input to scope the check.

## Configurable variables

All tool and runtime versions are overridable via **GitHub repository variables** (Settings → Secrets and variables → Actions → Variables). Callers pass them to reusable workflows using the `${{ vars.X || 'default' }}` pattern.

### Required variables

These must be set for binding generation to work:

| Variable | Example | Description |
|---|---|---|
| `REPO_OWNER` | `my-org` | GitHub org or username |
| `REPO_NAME` | `my-project` | Repository name (kebab-case) |

### Optional variables

Set these to override tool version defaults. If unset, the fallback defaults are used.

| Variable | Default | Used by | Description |
|---|---|---|---|
| `NPM_SCOPE` | `@{REPO_OWNER}` | release.yml | npm scope for JS packages |
| `NODE_VERSION` | `20` | ci.yml, build-js-bindings.yml | Node.js version |
| `GO_VERSION` | `stable` (ci) / `1.21` (build) | ci.yml, build-go-bindings.yml | Go version |
| `PYTHON_VERSION` | `3.12` | build-python-bindings.yml, build-go-bindings.yml | Python version |
| `PNPM_VERSION` | `9` | build-js-bindings.yml | pnpm version |
| `WASM_TARGET` | `wasm32-unknown-unknown` | build-js-bindings.yml | Rust WASM compilation target |

> **Note:** The Makefile has its own overridable variables for local builds, including all of the above plus `RUST_VERSION`, `NVM_VERSION`, `UV_VERSION`, and `MATURIN_VERSION`. The Makefile also auto-bootstraps all toolchains (Rust, Python, Go, Node.js) when missing. Run `make config` to see all resolved values.

## Release flow

The release workflow is triggered manually and executes a multi-stage pipeline with automatic rollback on failure.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        release.yml (dispatch)                          │
│                                                                        │
│  1. ci-gate                                                            │
│     └─ Verify CI passed on HEAD commit                                 │
│                                                                        │
│  2. detect-changes                                                     │
│     ├─ Snapshot all pre-existing tags (for safe rollback)               │
│     ├─ Compare crate versions against latest release tags              │
│     └─ Output: ffi-libs, crate-versions, has-ffi-changes              │
│                                                                        │
│  3. check-ffi-versions + check-workspace-versions                      │
│     └─ Guard: version parity and workspace consistency                 │
│                                                                        │
│  4. build-python / build-go / build-js (parallel)                      │
│     └─ Build bindings only for changed FFI crates                      │
│                                                                        │
│  5. publish-bindings                                                   │
│     ├─ Download all binding artifacts                                  │
│     ├─ Merge main into `bindings` branch                               │
│     ├─ Commit generated bindings (force-add past .gitignore)           │
│     ├─ Tag: ffi/{crate}/v{version} (pip/pnpm consumers)               │
│     └─ Tag: bindings/go/{pkg}/v{version} (Go module consumers)        │
│                                                                        │
│  6. release-rust-crates                                                │
│     ├─ Download wheels + npm tarballs                                  │
│     ├─ Tag: crates/{path}/v{version}                                   │
│     ├─ Create GitHub Release with assets                               │
│     └─ Generate release body via scripts/release-body.sh               │
│                                                                        │
│  7. rollback (on failure of any prior step)                            │
│     ├─ Delete GitHub releases + crate tags (skip pre-existing)         │
│     ├─ Delete Go module tags (skip pre-existing)                       │
│     ├─ Delete FFI bindings tags (skip pre-existing)                    │
│     └─ Reset bindings branch (force-with-lease, only if HEAD matches)  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Tag conventions

Each release creates up to three tags per changed crate:

| Tag pattern | Example | Purpose |
|---|---|---|
| `crates/{name}/v{version}` | `crates/core/v0.2.0` | Lib crate release (GitHub Release) |
| `crates/ffi/{name}/v{version}` | `crates/ffi/core/v0.2.0` | FFI crate release (GitHub Release) |
| `ffi/{name}/v{version}` | `ffi/core/v0.2.0` | FFI bindings snapshot on `bindings` branch |
| `bindings/go/{pkg}/v{version}` | `bindings/go/uniffi_math_demo_core/v0.2.0` | Go module version tag |

### Change detection logic

The `detect-changes` job determines what to release:

1. **Library crates** (`crates/lib/*/Cargo.toml`): Compares the current version against the version at the most recent `crates/{name}/v*` tag. A version change triggers a release.
2. **FFI crates** (`crates/ffi/*/Cargo.toml`): Same comparison against `crates/ffi/{name}/v*` tags. If the corresponding lib crate also changed, the FFI release is folded into the lib release. If only the FFI crate changed, it gets its own release.
3. **First-time crates** (no previous tag): Only released if the corresponding lib crate also changed, to avoid releasing unfinished FFI wrappers.

## Validation scripts

The CI invokes several scripts under `scripts/` as guardrails. Each supports a `--cached` mode for use in git pre-commit hooks.

### `scripts/check-ffi-coverage.sh`

Runs on PRs to verify that changes to library crate source files (`crates/lib/*/src/`) are accompanied by changes in the corresponding FFI crate (`crates/ffi/*/`). For minor or major version bumps, it additionally requires `.rs` source changes in the FFI crate's `src/` directory (not just `Cargo.toml` changes).

### `scripts/check-ffi-versions.sh`

Compares each FFI crate's `[package].version` against its corresponding lib crate. Fails if any FFI version is strictly less than the lib version (semver component comparison). This ensures downstream consumers always get bindings that match or exceed the underlying library.

### `scripts/check-workspace-versions.sh`

Reads each crate's version from its own `Cargo.toml` and compares it to the version declared in the root workspace `Cargo.toml` dependency table. Fails on any mismatch. The companion `scripts/sync-workspace-versions.sh` can auto-fix mismatches.

### `scripts/ffi-crate-name.sh`

Utility used by the Makefile and release scripts. Reads the `name` field from a crate's `Cargo.toml` and normalizes it to either hyphenated (`--hy`, default) or underscored (`--us`) form.

### `scripts/release-body.sh`

Generates the markdown body for GitHub Releases. Produces a platform table mapping assets (wheels, tarballs) to their language, platform, architecture, and source links. Handles both lib crate and FFI-only releases with appropriate context.

## Artifact flow

Artifacts are the primary mechanism for passing build outputs between jobs.

```
build-python-bindings ──┬── wheel-ubuntu-latest ──────────► test-python-bindings
(matrix: 3 OS)          ├── wheel-macos-latest
                        ├── wheel-windows-latest
                        ├── python-bindings-linux-amd64
                        ├── python-bindings-darwin-arm64 ─► publish-bindings
                        └── python-bindings-windows-amd64

build-go-bindings ──────┬── go-bindings-linux-amd64 ──────► test-go-bindings
(matrix: 2 OS)          └── go-bindings-darwin-arm64 ─────► publish-bindings

build-js-bindings ──────┬── js-bindings ──────────────────► publish-bindings
(single OS)             └── npm-tarballs ─────────────────► test-js-bindings

                                            release-rust-crates ◄── wheel-*
                                                                ◄── npm-tarballs
```

In the CI pipeline, test jobs download the ubuntu/linux artifacts to run smoke tests. In the release pipeline, `publish-bindings` downloads all binding artifacts to commit them to the `bindings` branch, and `release-rust-crates` downloads wheels and tarballs to attach to GitHub Releases.

## Caching

Each build workflow caches its UniFFI bindgen tool binary using `actions/cache@v4`, keyed on `runner.os` + the hash of `crates/ffi/tools/Cargo.toml`. This avoids rebuilding the bindgen tools on every run (they take several minutes to compile from source).

| Cache key pattern | Tool cached |
|---|---|
| `uniffi-bindgen-py-{os}-{hash}` | `uniffi-bindgen` (Python bindings) |
| `uniffi-bindgen-go-{os}-{hash}` | `uniffi-bindgen-go` (Go bindings) |
| `uniffi-bindgen-rn-{os}-{hash}` | `uniffi-bindgen-react-native` (JS/WASM bindings) |

Rust compilation artifacts are cached via `Swatinem/rust-cache@v2`.

## Debug logging

All build workflows support GitHub Actions debug logging. When a workflow is re-run with debug logging enabled (`runner.debug == '1'`), the Makefile's `LOG_LEVEL` is automatically set to `2`, which prints resolved variables, paths, and full commands.
