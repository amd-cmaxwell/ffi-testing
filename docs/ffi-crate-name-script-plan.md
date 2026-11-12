# Crate name resolution: from composed variables to Cargo.toml lookup

## Context

Package names were previously composed from `PKG_PREFIX`/`PKG_SUFFIX` Make
variables, which proved brittle — the composed names could diverge from the
actual Rust crate names in `Cargo.toml`. This was replaced with a shared
script that reads the canonical crate name from each crate's `Cargo.toml`.

## What changed

### 1. `scripts/ffi-crate-name.sh`

New script that takes a crate directory name (e.g. `core`, `linear-algebra`)
and outputs the crate name from `crates/ffi/<dir>/Cargo.toml`:

```
Usage: ffi-crate-name.sh [--us|--hy] <crate-dir>
  --hy (default): hyphenated form (normalized to hyphens)
  --us:           underscore form
```

### 2. Makefile

The repeated 3-line grep/sed/tr block in `generate-python`, `generate-go`,
and `generate-wasm-ts` was replaced with calls to the script:

```sh
crate_name_hy=$$($(CURDIR)/scripts/ffi-crate-name.sh $$crate_dir);
crate_name_us=$$($(CURDIR)/scripts/ffi-crate-name.sh --us $$crate_dir);
```

The following variables were removed entirely:

- `PKG_PREFIX`, `PKG_PREFIX_US`
- `PKG_SUFFIX`, `PKG_SUFFIX_US`
- `NPM_PKG_PREFIX`
- `REPO_NAME_US`
- Helper functions: `_strip_hyphens`, `_ensure_hyphen`, `_user_provided`

`NPM_ROOT_PKG_NAME` was kept (now defaults to `REPO_NAME`) since it's still
used for the pnpm workspace root package.json.

### 3. Workflows

- **Reusable build workflows** (`build-python-bindings.yml`,
  `build-go-bindings.yml`, `build-js-bindings.yml`): removed `pkg-prefix`,
  `pkg-suffix`, and `npm-pkg-prefix` inputs and env vars.
- **`ci.yml`**: removed those inputs from all `with:` blocks.
- **`release.yml`**: removed `PKG_PREFIX`, `PKG_SUFFIX`, `NPM_PKG_PREFIX` from
  the top-level `env:` block. Replaced composed-name logic in the "Tag
  bindings commit", "Create per-crate releases", and "Delete Go module tags"
  steps with calls to `scripts/ffi-crate-name.sh`. Simplified the "Derive FFI
  modules" step by removing the prefix/suffix normalization blocks.

### 4. Repository variables to remove

The following GitHub repository variables (`Settings > Variables`) are no
longer referenced and can be removed:

- `PKG_PREFIX`
- `PKG_SUFFIX`
- `NPM_PKG_PREFIX`
