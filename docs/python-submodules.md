# Python Submodule Structure

The Python bindings are organized into domain submodules instead of a single flat namespace.

## Usage

```python
from uniffi_testing_ffi import MathError
from uniffi_testing_ffi.core import add
from uniffi_testing_ffi.complex import Complex, complex_magnitude
from uniffi_testing_ffi.polynomials import quadratic_roots
from uniffi_testing_ffi.linear_algebra import Matrix2x2, determinant_2x2
from uniffi_testing_ffi.trig import distance_2d
from uniffi_testing_ffi.sequences import factorial, is_prime
```

## Modules

| Module | Contents |
|--------|----------|
| `uniffi_testing_ffi` | `MathError`, `InternalError` (shared error types) |
| `uniffi_testing_ffi.core` | `add` |
| `uniffi_testing_ffi.complex` | `Complex`, `ComplexPair`, all `complex_*` functions, `quadratic_roots_complex` |
| `uniffi_testing_ffi.polynomials` | `F64Pair`, all `polynomial_*` functions, `quadratic_roots`, `linear_root` |
| `uniffi_testing_ffi.linear_algebra` | `Matrix2x2`, `Matrix3x3`, `Point2`, `Point3`, `determinant_*`, `solve_*`, `matrix_*`, `transpose_*`, `dot_product`, `cross_product_3d` |
| `uniffi_testing_ffi.trig` | `degrees_*`, `radians_*`, `law_of_*`, `distance_*`, `midpoint_*`, `triangle_*`, `circle_*` |
| `uniffi_testing_ffi.sequences` | `arithmetic_*`, `geometric_*`, `factorial`, `fibonacci`, `permutations`, `combinations`, `gcd`, `lcm`, `is_prime` |

## How it works

UniFFI always generates a single flat Python module. We layer submodules on top:

1. Maturin builds the Rust crate and generates the flat bindings as `uniffi_testing_ffi._bindings` (configured via `module-name` in `pyproject.toml`)
2. A generator script (`bindings/python/scripts/generate_modules.py`) reads the `__all__` list from the generated file, classifies each name by prefix, and writes submodule `.py` files with explicit re-exports
3. The submodule files are checked into `bindings/python/lib/uniffi_testing_ffi/` for full IDE/autocomplete support

## Keeping submodules in sync

When you add new Rust functions to the FFI crate:

```bash
make generate-python-modules
```

This regenerates the submodule files. New functions matching existing prefixes (e.g., `complex_power`) are picked up automatically. Only a new prefix pattern or a new module requires editing the `MODULES` dict in `generate_modules.py`.

## Prefix classification rules

The generator classifies names using this mapping (defined in `generate_modules.py`):

- **core**: explicit names (`add`)
- **complex**: prefix `complex_`, types `Complex`/`ComplexPair`, explicit `quadratic_roots_complex`
- **polynomials**: prefix `polynomial_`, types `F64Pair`, explicit `quadratic_roots`/`linear_root`
- **linear_algebra**: prefixes `determinant_`/`solve_`/`matrix_`/`dot_`/`cross_`/`transpose_`, types `Matrix2x2`/`Matrix3x3`/`Point2`/`Point3`
- **trig**: prefixes `degrees_`/`radians_`/`law_of_`/`distance_`/`midpoint_`/`triangle_`/`circle_`
- **sequences**: prefixes `arithmetic_`/`geometric_`, explicit `factorial`/`fibonacci`/`permutations`/`combinations`/`gcd`/`lcm`/`is_prime`
- **Shared** (`MathError`, `InternalError`): exported from the top-level `__init__.py`
