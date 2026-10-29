# ffi-testing

A toy project for evaluating Rust FFI approaches with the goal of building Python and Go bindings for the [virtee/sev](https://github.com/virtee/sev) crate.

## Motivation

The `virtee/sev` crate provides Rust APIs for AMD SEV (Secure Encrypted Virtualization). To make these capabilities accessible from Python and Go, we need a reliable FFI bridge. This project uses a simple math library as a testbed to work through the challenges of cross-language bindings before applying the approach to the real crate.

## What's here

A Rust workspace with six domain crates (core, complex, polynomials, linear-algebra, trig, sequences) and an FFI crate that exposes 52 functions via [Mozilla UniFFI](https://github.com/mozilla/uniffi-rs). The FFI crate acts as a boundary layer — it defines UniFFI-compatible wrapper types and thin functions that delegate to the internal crates, keeping domain code free of FFI concerns.

### Key challenges explored

- **Types UniFFI can't cross directly** — tuples, fixed-size arrays, and tuple enum variants all required wrapper Record types (see [docs/ffi-conversion-pitfalls.md](docs/ffi-conversion-pitfalls.md))
- **Python packaging with maturin** — maturin builds Rust, runs uniffi-bindgen, and produces a wheel in one step
- **Python submodules over a flat namespace** — UniFFI generates one file; a [generator script](bindings/python/scripts/generate_modules.py) splits exports into domain submodules
- **Go bindings via uniffi-bindgen-go** — community tool from NordSecurity, static linking with CGo
- **Monorepo Go module** — subdirectory module path (`github.com/amd-cmaxwell/ffi-testing/bindings/go`)

## Project structure

```
crates/
  core/              Shared error types and basic arithmetic
  complex/           Complex number operations
  polynomials/       Polynomial evaluation, roots, derivatives
  linear-algebra/    Matrix operations, solvers, vector math
  trig/              Trigonometry, geometry, distance functions
  sequences/         Sequences, combinatorics, number theory
  ffi/               UniFFI wrapper layer (cdylib + staticlib)

bindings/
  python/            Maturin-packaged Python bindings
  go/                Go module with static library

examples/
  demo.py            Python smoke test
  demo.go            Go smoke test (run from bindings/go/)

docs/                Design decisions and pitfall documentation
```

## Quick start

### Python

```bash
pipx install git+https://github.com/amd-cmaxwell/ffi-testing.git@math-monorepo-uniffi#subdirectory=bindings/python
uniffi-testing-ffi-python
```

Or from a local branch:

```bash
pip install bindings/python/
uniffi-testing-ffi-python
```

### Go

```bash
go install github.com/amd-cmaxwell/ffi-testing/bindings/go/cmd/uniffi-testing-ffi-go@math-monorepo-uniffi
uniffi-testing-ffi-go
```

Or from a local branch:

```bash
cd bindings/go
go run cmd/uniffi-testing-ffi-go/main.go
```


## Regenerating bindings

```bash
make generate-python    # maturin develop --release
make generate-go        # uniffi-bindgen-go + copy + inject CGo LDFLAGS
make test-python        # run Python smoke test
make test-go            # run Go smoke test
```

## Dependencies

- Rust toolchain (for building the FFI crate)
- [maturin](https://www.maturin.rs/) (for Python packaging)
- [uniffi-bindgen-go](https://github.com/NordSecurity/uniffi-bindgen-go) v0.5.0+v0.29.5 (for Go bindings)
- Python >= 3.8, Go >= 1.21

## Next steps

Apply this approach to [virtee/sev](https://github.com/virtee/sev) to provide Python and Go interfaces for AMD SEV attestation and key management.
