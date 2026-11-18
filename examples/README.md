# Examples

Smoke-test demos that exercise the generated FFI bindings for each target language.

| Directory | Language | What it tests |
|-----------|----------|---------------|
| `py-demo/` | Python | Imports the maturin-built wheels and calls every exported function |
| `go-demo/` | Go | Imports the generated Go modules (with CGO static libs) and calls every exported function |
| `js-demo/` | Node.js | Installs the packed npm tarballs and runs `node --test` against the WASM bindings |
| `web-demo/` | Browser (TypeScript) | Builds with Vite to verify the WASM packages work in a frontend bundler context |

## Running

Each demo is run by its corresponding Makefile target:

```sh
make test-python   # py-demo
make test-go       # go-demo
make test-js       # js-demo
make test-web      # web-demo
make test-all      # all of the above
```

These targets build the bindings first if needed.

## Import paths

The demo source files contain import paths that reference the project's npm scope and Go module base. When the repository is forked or the owner changes, run:

```sh
make fix-example-imports
```

This rewrites the import paths in the demo files to match the current `GO_MODULE_BASE` and `NPM_SCOPE` derived from the git remote (or overridden via environment variables).
