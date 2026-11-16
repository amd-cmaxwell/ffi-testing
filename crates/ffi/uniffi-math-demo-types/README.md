# uniffi-math-demo-types-ffi

This is the shared FFI types crate for the `uniffi-math-demo` family of
crates. It provides common UniFFI type definitions — errors and records —
so that every math-demo FFI crate uses a single, consistent set of types
across all generated language bindings.

See [`crates/ffi/README.md`](../README.md) for the general pattern.

## Exported types

| Type | Kind | Description |
|------|------|-------------|
| `MathError` | `uniffi::Error` | Shared error enum (`DivisionByZero`, `InvalidInput`, `NoRealSolution`, `Overflow`) |
| `Point2` | `uniffi::Record` | 2D point with `x` and `y` fields |

## Consuming crates

The following FFI crates depend on this crate:

- `uniffi-math-demo-complex`
- `uniffi-math-demo-linear-algebra`
- `uniffi-math-demo-polynomials`
- `uniffi-math-demo-sequences`
- `uniffi-math-demo-trig`

Each of these crates declares this crate as an external package in its
`uniffi.toml` so that UniFFI imports the shared types rather than
re-generating them:

```toml
[bindings.python.external_packages]
uniffi_math_demo_types_ffi = "uniffi_math_demo_types_ffi"
```

> **Important:** This crate does **not** have its own `uniffi.toml` listing
> itself as an external package — doing so would prevent UniFFI from
> generating the type definitions in the first place. Only consuming crates
> should declare it in `external_packages`.
