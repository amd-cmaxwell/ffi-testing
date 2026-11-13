# ffi-testing

A toy project for evaluating Rust FFI approaches with the goal of building Python, Go, and JS/WASM bindings for the [virtee/sev](https://github.com/virtee/sev) crate.

## Motivation

The `virtee/sev` crate provides Rust APIs for AMD SEV (Secure Encrypted Virtualization). To make these capabilities accessible from Python, Go, and JavaScript, we need a reliable FFI bridge. This project uses a simple math library as a testbed to work through the challenges of cross-language bindings before applying the approach to the real crate.

## What's here

A Rust workspace with six domain crates and a per-crate FFI layer that exposes functions via [Mozilla UniFFI](https://github.com/mozilla/uniffi-rs). Each FFI crate acts as a boundary layer — it defines UniFFI-compatible wrapper types and thin functions that delegate to the corresponding lib crate, keeping domain code free of FFI concerns. Bindings are generated for Python (via maturin), Go (via uniffi-bindgen-go), and JS/WASM (via uniffi-bindgen-react-native's `wasm2` mode).

### Key challenges explored

- **Types UniFFI can't cross directly** — tuples, fixed-size arrays, and tuple enum variants all required wrapper Record types (see [docs/ffi-conversion-pitfalls.md](docs/ffi-conversion-pitfalls.md))
- **Per-crate FFI separation** — each lib crate has its own FFI crate under `crates/ffi/`, enabling independent versioning and releases (see [docs/adr-ffi-crate-separation.md](docs/adr-ffi-crate-separation.md))
- **Python packaging with maturin** — maturin builds Rust, runs uniffi-bindgen, and produces a wheel in one step
- **Go bindings via uniffi-bindgen-go** — community tool from NordSecurity, static linking with CGo
- **JS/WASM bindings** — WASM compilation with TypeScript wrappers and npm packaging, verified in both Node.js and browser contexts
- **Cascading version pipeline** — lib crate version bumps automatically propagate to FFI crates with semver-aware bumping
- **Per-crate release workflow** — independent GitHub releases with per-language artifacts and rollback support

## Project structure

```
crates/
  lib/
    core/            Shared error types and basic arithmetic
    complex/         Complex number operations
    polynomials/     Polynomial evaluation, roots, derivatives
    linear-algebra/  Matrix operations, solvers, vector math
    trig/            Trigonometry, geometry, distance functions
    sequences/       Sequences, combinatorics, number theory
  ffi/
    core/            FFI wrapper for core
    complex/         FFI wrapper for complex
    polynomials/     FFI wrapper for polynomials
    linear-algebra/  FFI wrapper for linear-algebra
    trig/            FFI wrapper for trig
    sequences/       FFI wrapper for sequences
    types/           Shared FFI types (no corresponding lib crate)
    tools/           Build tooling (excluded from workspace)

scripts/
  log.sh                      Shared logging helpers (emoji + GH Actions annotations)
  ffi-crate-name.sh           Resolve FFI crate name from Cargo.toml
  sync-ffi-versions.sh        Cascade lib version bumps to FFI crates
  check-ffi-coverage.sh       Verify FFI source changes for lib API changes
  check-ffi-versions.sh       Validate FFI versions >= lib versions
  sync-workspace-versions.sh  Sync root Cargo.toml workspace dep versions
  check-workspace-versions.sh Validate workspace dep versions match crates
  release-body.sh             Generate GitHub release notes with asset tables

.githooks/
  pre-commit               Version sync, FFI coverage, and bindings guard

bindings/                  Auto-generated language bindings (committed on release)

examples/
  py-demo/           Python smoke test
  go-demo/           Go smoke test
  js-demo/           JS smoke test (Node.js)
  web-demo/          Browser smoke test (Vite + vanilla TS)

docs/                Design decisions and pitfall documentation
```

## Setup

```bash
make setup    # configures git hooks and installs build tools
```

This sets up the `.githooks/` directory as the git hooks path and installs build tools.

## Git hooks

The `pre-commit` hook (see [`.githooks/README.md`](.githooks/README.md) for full details) runs automatically on every commit to keep versions and bindings consistent. All checks operate exclusively on staged files — unstaged working-tree changes are ignored. `git status` is printed automatically after any step stages new files or when the hook fails.

| Step | Script | Description |
|------|--------|-------------|
| 1 | — | Rejects staged edits to auto-generated `bindings/` files |
| 2 | `scripts/sync-ffi-versions.sh --cached` | Cascades staged lib version bumps to FFI crates (semver-aware) |
| 3 | `scripts/check-ffi-coverage.sh --cached` | Verifies staged lib source changes have corresponding FFI updates |
| 4 | `scripts/sync-workspace-versions.sh --cached` | Syncs workspace dep versions in root `Cargo.toml` for staged crates |
| 5 | `scripts/check-ffi-versions.sh --cached` | Validates staged FFI versions ≥ lib versions |
| 6 | `scripts/check-workspace-versions.sh --cached` | Validates staged workspace dep versions match crate versions |

### Hook environment variables

| Variable | Default | Effect |
|----------|---------|--------|
| `SKIP_BINDINGS_CHECK` | _(unset)_ | Set to `1` to skip the bindings edit guard (step 1) |
| `SKIP_FFI_COVERAGE` | _(unset)_ | Set to `1` to skip the FFI coverage check (step 3) |
| `GITHUB_ACTIONS` | _(unset)_ | When set (automatic in CI), log messages are prefixed with `::error::` / `::warning::` / `::notice::` / `::debug::` workflow commands so they appear as annotations in the GitHub Actions UI. When unset, messages use emoji-only prefixes (❌ ⚠️ ℹ️ 🔍 ✅). |

### Repository variables

When creating a new project from this template, configure these as GitHub repository variables (Settings > Secrets and variables > Actions > Variables):

| Variable | Required | Description |
|---|---|---|
| `REPO_OWNER` | yes | GitHub owner/org (used for package naming) |
| `REPO_NAME` | yes | Repository name |
| `NPM_SCOPE` | no | npm scope for JS packages (defaults to `@REPO_OWNER`) |
| `NODE_VERSION` | no | Node.js version (default: `20`) |
| `GO_VERSION` | no | Go version (default: `stable` for CI tests, `1.21` for builds) |
| `PYTHON_VERSION` | no | Python version (default: `3.12`) |
| `PNPM_VERSION` | no | pnpm version (default: `9`) |
| `WASM_TARGET` | no | Rust WASM compilation target (default: `wasm32-unknown-unknown`) |

The Makefile resolves `REPO_OWNER` and `REPO_NAME` dynamically from the git remote by default. See [`.github/workflows/README.md`](.github/workflows/README.md) for the full configurable variables reference.

## Building bindings

All required toolchains (Rust, Python/uv, Go, Node.js/pnpm) are automatically bootstrapped when you run a `make` target — no manual installation needed. Override versions via environment variables (run `make config` to see all resolved values):

```bash
make generate-python    # Build Python wheels for all FFI crates
make generate-go        # Generate Go modules with CGO static libs
make generate-js        # Build WASM + TS bindings and pack npm tarballs
make generate-all       # All of the above
make config             # Dump all resolved variables
```

### Makefile version variables

| Variable | Default | Description |
|---|---|---|
| `RUST_VERSION` | `stable` | Rust toolchain (channel, version, or nightly) |
| `GO_VERSION` | `1.21.0` | Go version |
| `NODE_VERSION` | `20` | Node.js version |
| `NVM_VERSION` | `v0.40.3` | nvm installer version |
| `PYTHON_VERSION` | `3.12` | Python version (installed via uv) |
| `PNPM_VERSION` | `latest` | pnpm version |
| `UV_VERSION` | _(latest)_ | uv version (empty = latest) |
| `MATURIN_VERSION` | _(latest)_ | maturin version (empty = latest) |
| `WASM_TARGET` | `wasm32-unknown-unknown` | Rust WASM compilation target |

### Makefile verbosity

`LOG_LEVEL` controls output verbosity for all `make` targets:

```bash
make ...                # default (LOG_LEVEL=1): progress messages
LOG_LEVEL=0 make ...    # quiet: errors only
LOG_LEVEL=2 make ...    # debug: resolved variables, paths, commands
```

## Testing

```bash
make test               # Run Rust tests
make test-python        # Run Python smoke test (examples/py-demo)
make test-go            # Run Go smoke test (examples/go-demo)
make test-js            # Run JS smoke test (examples/js-demo)
make test-web           # Build browser demo with Vite (examples/web-demo)
make test-all           # All binding tests
make clean-web          # Remove web-demo generated files
```

The `web-demo` is a browser-based interactive calculator that exercises every math module via WASM. It also runs the same smoke tests as `js-demo` on page load. After `make test-web`, you can open it locally:

```bash
cd examples/web-demo && npx vite
```

## CI / CD

For a comprehensive guide to the workflow architecture, configurable variables, artifact flow, release pipeline, and validation scripts, see [`.github/workflows/README.md`](.github/workflows/README.md).

### CI (`ci.yml`)

Runs on push to main, pull requests, and manual dispatch. Validates FFI coverage and version alignment, runs Rust tests on Linux/macOS/Windows, builds all three binding flavors, and runs smoke tests for each.

### Release (`release.yml`)

Triggered manually via `workflow_dispatch`. Supports optional inputs to scope the release to specific crates:
- **`crates`** — space-separated lib crate names (e.g. `core trig`)
- **`ffi-libs`** — space-separated FFI lib names

When inputs are empty, all changed crates are released (default).

The release flow:
1. Verifies CI passed on HEAD
2. Detects version changes by comparing against git tags
3. Builds Python, Go, and JS bindings for changed FFI crates
4. Commits bindings to the `bindings` branch and tags them
5. Creates per-crate GitHub releases with artifact tables
6. Rolls back tags, releases, and the bindings branch on failure

#### Versioning

- **Lib crates** are tagged as `crates/{name}/v{version}`
- **FFI crates** are tagged as `crates/ffi/{name}/v{version}` (Rust release) and `ffi/{name}/v{version}` (bindings)
- **Go modules** are tagged as `bindings/go/{pkg}/v{version}`

When a lib crate version changes, `sync-ffi-versions.sh` automatically bumps the corresponding FFI crate:
- If the FFI version is behind the lib version, it matches the lib version
- If the FFI version is already at or ahead, it gets a patch bump
- The sync is idempotent across pre-commit hook re-runs

FFI-only releases (no lib change) are distinguished with a `-ffi` title suffix and a release body notice.

## Dependencies

All dependencies are bootstrapped automatically by the Makefile when needed. If you prefer to install them manually:

- **Rust** — installed via [rustup](https://rustup.rs/) (`install-rust` target); WASM target added on demand (`ensure-wasm-target`)
- **Python** — installed via [uv](https://docs.astral.sh/uv/) (`install-python-toolchain` target)
- **maturin** — installed via `uv tool install` (`python-dependencies` target)
- **Go** — installed from [go.dev](https://go.dev/dl/) (`install-go-toolchain` target)
- **Node.js** — installed via [nvm](https://github.com/nvm-sh/nvm) (`install-node` target)
- **pnpm** — installed via corepack (`js-dependencies` target)

## Next steps

Apply this approach to [virtee/sev](https://github.com/virtee/sev) to provide Python, Go, and JS interfaces for AMD SEV attestation and key management.
