# Makefile Debug Logging for External Commands

All external tool invocations in the Makefile now conditionally log the exact
command being run, visible only at `LOG_LEVEL=2` (debug). At the default level
(1) or quiet (0), these lines are suppressed.

## Usage

```sh
LOG_LEVEL=2 make generate-python   # shows every cargo, maturin, pip, etc. command
make generate-python               # default: progress messages only
LOG_LEVEL=0 make generate-python   # quiet: errors only
```

## Targets and commands covered

| Target | Commands logged |
|---|---|
| `build` | `cargo build --release` |
| `test` | `cargo test --workspace` |
| `clean` | `cargo clean` |
| `python-dependencies` | `pip install --user pipx`, `python3 -m pip install --user pipx`, `pipx install maturin` |
| `install-node` | `nvm install 20` |
| `js-dependencies` | `corepack enable && corepack prepare pnpm@latest --activate` |
| `build-tools-py` | `cargo install ... --bin uniffi-bindgen` |
| `build-tools-go` | `cargo install ... --bin uniffi-bindgen-go` |
| `build-tools-wasm` | `cargo install ... --bin uniffi-bindgen-react-native` |
| `generate-python` | `maturin build --release --out ...` |
| `install-python` | `pip install --force-reinstall <wheel>` |
| `generate-go` | `uniffi-bindgen-go ...`, `go mod init ...` |
| `install-go` | `go work init`, `go work use ...` |
| `generate-wasm-ts` | `cargo build --release --target wasm32-unknown-unknown ...`, `uniffi-bindgen-react-native ...` |
| `compile-wasm-ts` | `pnpm install`, `pnpm exec tsc` |
| `generate-js` | `npm pack ...` |
| `install-js` | `npm install <tarball>` (local and global) |
| `test-python` | `python examples/py-demo` |
| `test-go` | `go build ./...`, `go mod edit ...`, `go run main.go` |
| `test-js` | `node --test examples/js-demo/main.mjs` |
| `test-web` | `npm install`, `npx vite build` |

## Implementation

Two logging mechanisms are used depending on context:

- **Make-level** (`_log_debug`): for standalone recipe lines outside shell loops.
  Bare commands (without `@`) were silenced so output is controlled exclusively
  by `LOG_LEVEL`.
- **Shell-level** (`log_debug`): for commands inside `for` loops or conditional
  blocks that already use `$(_shell_log_helpers)`. The helpers were added to
  `python-dependencies`, `install-node`, and `js-dependencies` which previously
  lacked them.

Each `log_debug` call prints the command exactly as it will execute, including
any `cd` prefix, environment variables, and resolved arguments.
