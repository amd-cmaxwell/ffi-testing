#!/usr/bin/env python3
"""Generate Go sub-packages that re-export uniffi-generated bindings by module.

Reads the generated uniffi_testing_ffi.go file, classifies each exported
function and type by prefix, and writes Go sub-package directories with
thin wrapper functions and type aliases.

Usage: python scripts/generate_modules.py <generated_bindings.go>
"""
import re
import sys
from pathlib import Path

MODULES = {
    "core": {
        "prefixes": [],
        "names": {"Add"},
        "types": set(),
        "go_pkg": "core",
    },
    "complexnums": {
        "prefixes": ["Complex"],
        "names": {"QuadraticRootsComplex"},
        "types": {"Complex", "ComplexPair"},
        "go_pkg": "complexnums",
    },
    "polynomials": {
        "prefixes": ["Polynomial"],
        "names": {"QuadraticRoots", "LinearRoot"},
        "types": {"F64Pair"},
        "go_pkg": "polynomials",
    },
    "linearalgebra": {
        "prefixes": ["Determinant", "Solve", "Matrix", "Dot", "Cross", "Transpose"],
        "names": set(),
        "types": {"Matrix2x2", "Matrix3x3", "Point2", "Point3"},
        "go_pkg": "linearalgebra",
    },
    "trig": {
        "prefixes": ["Degrees", "Radians", "LawOf", "Distance", "Midpoint", "Triangle", "Circle"],
        "names": set(),
        "types": set(),
        "go_pkg": "trig",
    },
    "sequences": {
        "prefixes": ["Arithmetic", "Geometric"],
        "names": {"Factorial", "Fibonacci", "Permutations", "Combinations", "Gcd", "Lcm", "IsPrime"},
        "types": set(),
        "go_pkg": "sequences",
    },
}

SHARED_TYPES = {"MathError"}

SHARED_ERROR_VARS = {
    "ErrMathErrorDivisionByZero",
    "ErrMathErrorInvalidInput",
    "ErrMathErrorNoRealSolution",
    "ErrMathErrorOverflow",
}

SHARED_ERROR_TYPES = {
    "MathErrorDivisionByZero",
    "MathErrorInvalidInput",
    "MathErrorNoRealSolution",
    "MathErrorOverflow",
}

SKIP_PREFIXES = {"FfiConverter", "FfiDestroyer", "GoRustBuffer", "RustBuffer",
                 "NewMathError", "CFromRust", "LowerInto", "LiftFrom"}

FUNC_RE = re.compile(
    r'^func ([A-Z]\w+)\(([^)]*)\)\s*(.*?)\s*\{',
    re.MULTILINE,
)

TYPE_RE = re.compile(
    r'^type ([A-Z]\w+) struct',
    re.MULTILINE,
)


def should_skip(name: str) -> bool:
    return any(name.startswith(p) for p in SKIP_PREFIXES)


def extract_functions(source: str) -> list[dict]:
    funcs = []
    for m in FUNC_RE.finditer(source):
        name = m.group(1)
        if should_skip(name):
            continue
        params_raw = m.group(2).strip()
        returns_raw = m.group(3).strip()
        funcs.append({
            "name": name,
            "params": params_raw,
            "returns": returns_raw,
        })
    return funcs


def extract_types(source: str) -> list[str]:
    domain_types = set()
    for mod_spec in MODULES.values():
        domain_types |= mod_spec["types"]
    domain_types |= SHARED_TYPES | SHARED_ERROR_TYPES

    found = []
    for m in TYPE_RE.finditer(source):
        name = m.group(1)
        if name in domain_types:
            found.append(name)
    return found


def classify_func(name: str) -> str | None:
    for mod_name, spec in MODULES.items():
        if name in spec["names"]:
            return mod_name
        for prefix in spec["prefixes"]:
            if name.startswith(prefix):
                return mod_name
    return None


def classify_type(name: str) -> str | None:
    for mod_name, spec in MODULES.items():
        if name in spec["types"]:
            return mod_name
    return None


def strip_prefix(func_name: str, mod_name: str) -> str:
    spec = MODULES[mod_name]
    if func_name in spec["names"]:
        return func_name
    for prefix in sorted(spec["prefixes"], key=len, reverse=True):
        if func_name.startswith(prefix):
            stripped = func_name[len(prefix):]
            if stripped and stripped[0].isalpha():
                return stripped
    return func_name


def resolve_collisions(funcs: list[dict]) -> None:
    """If multiple functions got the same stripped name, revert them to originals."""
    from collections import Counter
    counts = Counter(f["stripped_name"] for f in funcs)
    for f in funcs:
        if counts[f["stripped_name"]] > 1:
            f["stripped_name"] = f["name"]


def param_names(params_raw: str) -> str:
    if not params_raw:
        return ""
    parts = []
    for chunk in params_raw.split(","):
        chunk = chunk.strip()
        if not chunk:
            continue
        tokens = chunk.split()
        parts.append(tokens[0])
    return ", ".join(parts)


def uses_error(returns_raw: str) -> bool:
    return "error" in returns_raw


def uses_type(params_raw: str, returns_raw: str, type_name: str) -> bool:
    combined = params_raw + " " + returns_raw
    return re.search(rf'\b{re.escape(type_name)}\b', combined) is not None


def read_module_path(go_dir: Path) -> str:
    go_mod = go_dir / "go.mod"
    for line in go_mod.read_text().splitlines():
        if line.startswith("module "):
            return line.split()[1]
    raise ValueError(f"Could not find module path in {go_mod}")


ALL_DOMAIN_TYPES: set[str] = set()
for _spec in MODULES.values():
    ALL_DOMAIN_TYPES |= _spec["types"]
ALL_DOMAIN_TYPES |= SHARED_TYPES


def find_referenced_types(funcs: list[dict]) -> set[str]:
    """Find all domain types referenced in function signatures."""
    referenced = set()
    for f in funcs:
        combined = f["params"] + " " + f["returns"]
        for t in ALL_DOMAIN_TYPES:
            if re.search(rf'\b{re.escape(t)}\b', combined):
                referenced.add(t)
    return referenced


def generate_package_file(
    mod_name: str,
    go_pkg: str,
    module_path: str,
    funcs: list[dict],
    types: list[str],
    needs_errors: bool,
) -> str:
    referenced = find_referenced_types(funcs)
    all_types = sorted(set(types) | referenced)

    lines = []
    lines.append(f"package {go_pkg}")
    lines.append("")
    lines.append("import (")
    lines.append(f'\tffi "{module_path}/uniffi_testing_ffi"')
    lines.append(")")
    lines.append("")

    for t in all_types:
        lines.append(f"type {t} = ffi.{t}")
    if all_types:
        lines.append("")

    if needs_errors:
        for t in sorted(SHARED_TYPES):
            if t not in all_types:
                lines.append(f"type {t} = ffi.{t}")
        for t in sorted(SHARED_ERROR_TYPES):
            lines.append(f"type {t} = ffi.{t}")
        lines.append("")
        for v in sorted(SHARED_ERROR_VARS):
            lines.append(f"var {v} = ffi.{v}")
        lines.append("")

    for f in sorted(funcs, key=lambda x: x["stripped_name"]):
        stripped = f["stripped_name"]
        params = f["params"]
        returns = f["returns"]
        args = param_names(params)
        orig = f["name"]

        if returns:
            lines.append(f"func {stripped}({params}) {returns} {{ return ffi.{orig}({args}) }}")
        else:
            lines.append(f"func {stripped}({params}) {{ ffi.{orig}({args}) }}")

    lines.append("")
    return "\n".join(lines)


def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <generated_bindings.go>", file=sys.stderr)
        sys.exit(1)

    source_path = Path(sys.argv[1])
    source = source_path.read_text()
    go_dir = Path(__file__).resolve().parent.parent

    module_path = read_module_path(go_dir)

    all_funcs = extract_functions(source)
    all_types = extract_types(source)

    classified_funcs: dict[str, list[dict]] = {mod: [] for mod in MODULES}
    classified_types: dict[str, list[str]] = {mod: [] for mod in MODULES}

    for func in all_funcs:
        mod = classify_func(func["name"])
        if mod:
            func["stripped_name"] = strip_prefix(func["name"], mod)
            classified_funcs[mod].append(func)

    for type_name in all_types:
        mod = classify_type(type_name)
        if mod:
            classified_types[mod].append(type_name)

    for mod_funcs in classified_funcs.values():
        resolve_collisions(mod_funcs)

    for mod_name, spec in MODULES.items():
        funcs = classified_funcs[mod_name]
        types = classified_types[mod_name]
        if not funcs and not types:
            continue

        go_pkg = spec["go_pkg"]
        pkg_dir = go_dir / go_pkg
        pkg_dir.mkdir(exist_ok=True)

        needs_errors = any(uses_error(f["returns"]) for f in funcs)

        content = generate_package_file(
            mod_name, go_pkg, module_path, funcs, types, needs_errors,
        )

        out_file = pkg_dir / f"{go_pkg}.go"
        out_file.write_text(content)
        print(f"  wrote {out_file.relative_to(go_dir)}")


if __name__ == "__main__":
    main()
