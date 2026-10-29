#!/usr/bin/env python3
"""Generate submodule re-export files from the uniffi-generated bindings.

Reads the __all__ list from the generated _bindings file, classifies each
name by prefix, and writes submodule .py files into the package directory.

Usage: python scripts/generate_modules.py <generated_bindings.py>
"""
import ast
import sys
from pathlib import Path

MODULES = {
    "core": {
        "prefixes": [],
        "names": {"add"},
        "types": set(),
    },
    "complex": {
        "prefixes": ["complex_"],
        "names": {"quadratic_roots_complex"},
        "types": {"Complex", "ComplexPair"},
    },
    "polynomials": {
        "prefixes": ["polynomial_"],
        "names": {"quadratic_roots", "linear_root"},
        "types": {"F64Pair"},
    },
    "linear_algebra": {
        "prefixes": ["determinant_", "solve_", "matrix_", "dot_", "cross_", "transpose_"],
        "names": set(),
        "types": {"Matrix2x2", "Matrix3x3", "Point2", "Point3"},
    },
    "trig": {
        "prefixes": ["degrees_", "radians_", "law_of_", "distance_", "midpoint_", "triangle_", "circle_"],
        "names": set(),
        "types": set(),
    },
    "sequences": {
        "prefixes": ["arithmetic_", "geometric_"],
        "names": {"factorial", "fibonacci", "permutations", "combinations", "gcd", "lcm", "is_prime"},
        "types": set(),
    },
}

SHARED = {"MathError", "InternalError"}


def extract_all_list(source_path: Path) -> list[str]:
    text = source_path.read_text()
    for node in ast.walk(ast.parse(text)):
        if isinstance(node, ast.Assign):
            for target in node.targets:
                if isinstance(target, ast.Name) and target.id == "__all__":
                    return ast.literal_eval(node.value)
    raise ValueError(f"__all__ not found in {source_path}")


def classify(name: str) -> str | None:
    for mod_name, spec in MODULES.items():
        if name in spec["types"] or name in spec["names"]:
            return mod_name
        for prefix in spec["prefixes"]:
            if name.startswith(prefix):
                return mod_name
    if name in SHARED:
        return None
    return None


def generate_module_file(mod_name: str, names: list[str]) -> str:
    imports = ",\n    ".join(sorted(names))
    return (
        f"from uniffi_testing_ffi._bindings import (  # noqa: F401\n"
        f"    {imports},\n"
        f")\n\n"
        f"__all__ = [\n"
        + "".join(f'    "{n}",\n' for n in sorted(names))
        + "]\n"
    )


def generate_init(module_names: list[str], shared_names: list[str]) -> str:
    lines = []
    for name in shared_names:
        lines.append(f"from uniffi_testing_ffi._bindings import {name}  # noqa: F401")
    lines.append("")
    for mod in sorted(module_names):
        lines.append(f"from uniffi_testing_ffi import {mod}  # noqa: F401")
    lines.append("")
    lines.append(f"__all__ = [")
    for name in sorted(shared_names):
        lines.append(f'    "{name}",')
    for mod in sorted(module_names):
        lines.append(f'    "{mod}",')
    lines.append("]")
    lines.append("")
    return "\n".join(lines)


def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <generated_bindings.py>", file=sys.stderr)
        sys.exit(1)

    source = Path(sys.argv[1])
    all_names = extract_all_list(source)
    pkg_dir = Path(__file__).resolve().parent.parent / "lib" / "uniffi_testing_ffi"

    classified: dict[str, list[str]] = {mod: [] for mod in MODULES}
    shared_names: list[str] = []

    for name in all_names:
        mod = classify(name)
        if mod:
            classified[mod].append(name)
        elif name in SHARED:
            shared_names.append(name)

    active_modules = []
    for mod_name, names in classified.items():
        if not names:
            continue
        active_modules.append(mod_name)
        out = pkg_dir / f"{mod_name}.py"
        out.write_text(generate_module_file(mod_name, names))
        print(f"  wrote {out.relative_to(pkg_dir.parent.parent)}")

    init = pkg_dir / "__init__.py"
    init.write_text(generate_init(active_modules, shared_names))
    print(f"  wrote {init.relative_to(pkg_dir.parent.parent)}")


if __name__ == "__main__":
    main()
