# Extract hardcoded versions into overridable variables

## Context

This repo is a template for bootstrapping Rust FFI devops tooling. Downstream consumers currently must hunt through the Makefile and CI workflows to update hardcoded version numbers for Node.js, Python, Go, nvm, pnpm, and the WASM target. Extracting these into overridable variables (Make `?=` and workflow `inputs` backed by GitHub repository variables) reduces the code changes downstream forks need to make.

The Makefile already does this well for project identity (REPO_OWNER, REPO_NAME, etc.) and output paths. The gap is **tool/runtime version numbers**.

## Changes

### 1. Makefile — new version variables

Added a `# ── Tool / runtime versions` section (after the logging block, before project identity) with:

```makefile
NODE_VERSION   ?= 20
NVM_VERSION    ?= v0.40.3
PNPM_VERSION   ?= latest
WASM_TARGET    ?= wasm32-unknown-unknown
```

All hardcoded references replaced with variable expansions. Variables appear in `make config` output.

### 2. CI workflows — GitHub repository variables with fallback defaults

Downstream repos override versions by setting GitHub repository variables
(Settings → Secrets and variables → Actions → Variables):

| Variable          | Default                  | Used by                            |
|-------------------|--------------------------|------------------------------------|
| `NODE_VERSION`    | `20`                     | ci.yml, build-js-bindings.yml      |
| `GO_VERSION`      | `stable` (ci) / `1.21` (build) | ci.yml, build-go-bindings.yml |
| `PYTHON_VERSION`  | `3.12`                   | build-python-bindings.yml, build-go-bindings.yml |
| `PNPM_VERSION`    | `9`                      | build-js-bindings.yml              |
| `WASM_TARGET`     | `wasm32-unknown-unknown` | build-js-bindings.yml              |

**Pattern used**: `${{ vars.X || 'default' }}` in callers (ci.yml, release.yml), passed as `workflow_call` inputs to reusable workflows. The reusable workflow input `default:` values serve as fallbacks when called without the input.

**Reusable workflows** (build-js-bindings, build-python-bindings, build-go-bindings):
- Accept version inputs with static defaults
- Use `${{ inputs.* }}` in steps

**ci.yml** — uses `env:` block with `${{ vars.X || 'default' }}` for test jobs (direct jobs, not reusable workflows). Passes `vars.*` through to reusable workflow calls.

**release.yml** — passes `vars.*` through to reusable workflow calls (same pattern as ci.yml).

### 3. Files NOT changed

- **Cargo.toml / crate versions**: Semantic, per-crate values managed by version-check scripts.
- **UniFFI ecosystem versions** (uniffi `"0.31"`, bindgen tags): Tightly coupled across Cargo.toml files; changing them requires Cargo resolution, not a simple variable swap.
- **GitHub Actions versions** (`actions/checkout@v4` etc.): Dependabot/Renovate handles these.
- **Rust edition, TypeScript target, engines fields**: Semantic project choices, not bootstrapping knobs.
- **Repo identity/URLs**: Already handled by existing REPO_OWNER/REPO_NAME/NPM_SCOPE vars.
