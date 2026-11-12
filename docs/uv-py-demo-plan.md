# Plan: Convert py-demo to a uv-installable package

## Context

On Ubuntu 23.04+ and other distros using PEP 668, `pip install` into the system Python fails with "externally-managed-environment". Instead of a generic `.venv`, make `examples/py-demo` a proper Python package installable via `uv tool install`. The `test-python` target handles installing the demo app via `uv` with binding wheels bundled as `--with` dependencies into an isolated tool venv. The `install-python` target stays generic — it installs wheels into the current environment (useful when the user is already in their own venv).

## Changes

### 1. Create `examples/py-demo/pyproject.toml`

Uses hatchling as the build backend:

```toml
[project]
name = "uniffi-math-demo"
version = "0.1.0"
requires-python = ">=3.9"

[project.scripts]
uniffi-math-demo = "py_demo:main"

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[tool.hatch.build.targets.wheel]
packages = ["py_demo.py"]
```

### 2. Restructure `examples/py-demo/`

```
examples/py-demo/
├── pyproject.toml
├── py_demo.py          # main() + all demo code (from __main__.py)
└── __main__.py          # thin: from py_demo import main; main()
```

### 3. Update `python-dependencies`

Ensure `uv` is available (installs via official installer if missing). Use `uv tool install` for maturin instead of pipx.

### 4. Update Makefile targets

- **`install-python`**: Unchanged — `pip install --force-reinstall <whl>` for each wheel. Stays generic for users in their own venv.
- **`test-python`**: Depends on `generate-python` (not `install-python`). Runs:
  1. `uv tool install --force examples/py-demo/ --with <whl>...` (single command installs demo + all bindings)
  2. `uniffi-math-demo` (the console script)
- **`clean-py-demo`**: `uv tool uninstall uniffi-math-demo`
- **`clean-demos`**: Depends on `clean-py-demo clean-web`
- **`clean-all`**: Includes `clean-demos`

## Verification

1. `make test-python` → installs demo via uv, bundles bindings, runs CLI
2. `make install-python` → still works for manual `pip install` in user's venv
3. `make clean-py-demo` → uninstalls cleanly
4. `make clean-demos` → cleans both py-demo and web-demo
