SHELL          := $(shell command -v bash)
.SHELLFLAGS    := -euo pipefail -c

# ── Logging ──────────────────────────────────────────────────────────────────
# LOG_LEVEL controls output verbosity:
#   0 = quiet   (errors only)
#   1 = info    (progress messages — the default)
#   2 = debug   (resolved variables, paths, file operations)
LOG_LEVEL      ?= 1
_LOG_LEVEL     := $(LOG_LEVEL)

# Make-level log helpers (use in simple recipe lines, NOT inside shell loops)
_log_info  = $(if $(filter 1 2,$(_LOG_LEVEL)),@echo "  $(1)",@:)
_log_debug = $(if $(filter 2,$(_LOG_LEVEL)),@echo "  [debug] $(1)",@:)

# Shell-level log helpers (paste into the preamble of for-loop recipe lines)
define _shell_log_helpers
log_info()  { [ $(_LOG_LEVEL) -ge 1 ] && echo "  $$*" || :; }; \
log_debug() { [ $(_LOG_LEVEL) -ge 2 ] && echo "  [debug] $$*" || :; };
endef

# ──────────────────────────────────────────────────────────────────────────────
# UniFFI Multi-Language Bindings Makefile
#
# Generates Python, Go, and WASM/TypeScript bindings from Rust FFI crates using
# UniFFI. Each crate under crates/ffi/ produces a language-specific package.
#
# Quick start:
#   make generate-python    # Build Python wheels for all FFI crates
#   make generate-go        # Generate Go modules with CGO static libs
#   make generate-js        # Build WASM + TS bindings and pack npm tarballs
#   make generate-all       # All of the above
#   make config             # Dump all resolved variables
#
# Verbosity:
#   make ...                # default (LOG_LEVEL=1): progress messages
#   LOG_LEVEL=0 make ...    # quiet: errors only
#   LOG_LEVEL=2 make ...    # debug: resolved variables, paths, commands
#
# Most variables are resolved dynamically from the git remote and directory
# name. Override REPO_SLUG, REPO_OWNER, REPO_NAME, PKG_SUFFIX, GIT_REPO,
# or NPM_SCOPE as needed.
# ──────────────────────────────────────────────────────────────────────────────

.PHONY: config setup build test \
       python-dependencies js-dependencies install-node \
       build-tools-py build-tools-go build-tools-wasm \
       generate-python generate-go generate-wasm-ts generate-all \
       install-python install-go install-js compile-wasm-ts generate-js \
       test-python test-go test-js test-web test-all \
       stamp-uniffi-config _check-repo-owner \
       clean clean-tools clean-python clean-go clean-js clean-bindings \
       clean-web clean-py-artifacts clean-js-artifacts clean-artifacts clean-all purge

# ── Project identity ─────────────────────────────────────────────────────────
# Resolution order: git remote -> REPO_SLUG -> REPO_OWNER + REPO_NAME.
# All three are overridable. REPO_NAME falls back to the directory name.
# REPO_OWNER is required for build targets (not clean). After resolution,
# REPO_SLUG is re-normalized to REPO_OWNER/REPO_NAME.
#
# REPO_SLUG    : owner/name pair; derived from the git remote URL
# REPO_OWNER   : GitHub username or org; derived from REPO_SLUG
# REPO_NAME    : kebab-case project name; derived from REPO_SLUG, then dir name
# PKG_PREFIX   : kebab-case prefix for generated packages; user-settable,
#                defaults to REPO_NAME when REPO_NAME is explicitly provided
# PKG_SUFFIX   : appended to every generated package name (e.g. "-ffi")
# *_US variants: underscore-normalized forms for contexts where kebab-case is
#                invalid (Rust crate names, Go package dirs, WASM filenames)
_git_remote_url  = $(shell git remote get-url $(1) 2>/dev/null)
_strip_git_url   = $(patsubst %.git,%,$(patsubst git@%,%,$(patsubst https://%,%,$(patsubst http://%,%,$(1)))))
_normalize_ssh   = $(subst :,/,$(_strip_git_url))
_GIT_REMOTE_RAW := $(or $(call _git_remote_url,upstream),$(call _git_remote_url,origin))
_GIT_REMOTE     := $(if $(_GIT_REMOTE_RAW),$(call _normalize_ssh,$(_GIT_REMOTE_RAW)),)
_REMOTE_PARTS   := $(subst /, ,$(_GIT_REMOTE))

REPO_SLUG      ?= $(if $(word 2,$(_REMOTE_PARTS)),$(word 2,$(_REMOTE_PARTS))/$(word 3,$(_REMOTE_PARTS)),)
REPO_OWNER     ?= $(firstword $(subst /, ,$(REPO_SLUG)))
REPO_NAME      ?= $(or $(word 2,$(subst /, ,$(REPO_SLUG))),$(notdir $(CURDIR)))
REPO_NAME_US   := $(subst -,_,$(REPO_NAME))
REPO_SLUG      := $(REPO_OWNER)/$(REPO_NAME)
PKG_SUFFIX     ?=
PKG_SUFFIX_US  := $(subst -,_,$(PKG_SUFFIX))
# Helper: normalize a value to end with exactly one hyphen (empty → empty)
_strip_hyphens = $(patsubst %-,%,$(patsubst %-,%,$(patsubst %-,%,$(1))))
_ensure_hyphen = $(if $(strip $(1)),$(call _strip_hyphens,$(strip $(1)))-,)

# True (non-empty) when the user explicitly set a variable via the command line
# or environment, as opposed to it being derived from the Makefile defaults.
_user_provided = $(filter command line environment,$(origin $(1)))

PKG_PREFIX     ?=
ifeq ($(strip $(PKG_PREFIX)),)
  ifneq ($(call _user_provided,REPO_NAME),)
    PKG_PREFIX := $(REPO_NAME)
  endif
endif
PKG_PREFIX     := $(call _ensure_hyphen,$(PKG_PREFIX))
PKG_PREFIX_US  := $(subst -,_,$(PKG_PREFIX))

_check-repo-owner:
	@test -n "$(REPO_OWNER)" || { echo "error: REPO_OWNER could not be resolved — configure a git remote, or set REPO_OWNER explicitly"; exit 1; }

# ── Crate discovery ──────────────────────────────────────────────────────────
# Discovers FFI crates under crates/ffi/ (excluding "tools"). FFI_LIBS is the
# underscore-normalized list used in loop bodies; defaults to "auto" which
# resolves to all discovered crates. Pass FFI_LIBS="foo bar" to build a subset.
FFI_CRATE_DIRS  := $(wildcard crates/ffi/*/Cargo.toml)
FFI_CRATE_PATHS := $(patsubst %/Cargo.toml,%,$(FFI_CRATE_DIRS))
FFI_CRATE_NAMES := $(filter-out tools,$(notdir $(FFI_CRATE_PATHS)))
AUTO_FFI_LIBS   := $(sort $(subst -,_,$(FFI_CRATE_NAMES)))
FFI_LIBS         ?= auto
override FFI_LIBS := $(if $(filter auto,$(FFI_LIBS)),$(AUTO_FFI_LIBS),$(FFI_LIBS))

# ── Output directories ───────────────────────────────────────────────────────
# Intermediate binding sources go under bindings/; distributable artifacts
# (wheels, tarballs) go under dist/.
BINDINGS_DIR      ?= $(CURDIR)/bindings
PY_BINDINGS_DIR   ?= $(BINDINGS_DIR)/python
GO_BINDINGS_DIR   ?= $(BINDINGS_DIR)/go
JS_BINDINGS_DIR   ?= $(BINDINGS_DIR)/js
PY_WHEELS_DIR     ?= $(CURDIR)/dist/pypi
NPM_TARBALLS_DIR  ?= $(CURDIR)/dist/npm

# ── Git remote / Go module path ──────────────────────────────────────────────
# GIT_REPO       : bare host/org/repo path (e.g. github.com/org/repo); resolved
#                  from the git remote, falling back to github.com/REPO_SLUG
# GO_MODULE_BASE : Go module import path root (GIT_REPO/bindings/go)
GIT_REPO        ?= $(or $(_GIT_REMOTE),github.com/$(REPO_SLUG))
GO_MODULE_BASE  ?= $(GIT_REPO)/bindings/go

# ── CGO cross-compilation ───────────────────────────────────────────────────
# CGO_OS / CGO_ARCH : target platform for cross-compiled static libs
CGO_OS         ?= $(shell go env GOOS 2>/dev/null || echo linux)
CGO_ARCH       ?= $(shell go env GOARCH 2>/dev/null || echo amd64)

# Platform-specific shared library extension (.so / .dylib)
LIB_SUFFIX_linux  := .so
LIB_SUFFIX_darwin := .dylib
LIB_SUFFIX        := $(LIB_SUFFIX_$(CGO_OS))

# Platform-specific linker flags for CGO
CGO_LDFLAGS_linux  := -lm -ldl -lpthread
CGO_LDFLAGS_darwin := -lm -lpthread -framework Security
CGO_LDFLAGS_EXTRA  := $(CGO_LDFLAGS_$(CGO_OS))

# ── npm / JS configuration ──────────────────────────────────────────────────
# NPM_SCOPE      : npm org scope (defaults to @REPO_OWNER; override for orgs)
# NPM_PKG_PREFIX : kebab-case prefix within the scope; user-settable,
#                  defaults to PKG_PREFIX, then REPO_NAME
NPM_SCOPE       ?= @$(REPO_OWNER)
NPM_PKG_PREFIX  ?=
ifeq ($(strip $(NPM_PKG_PREFIX)),)
  NPM_PKG_PREFIX := $(or $(strip $(PKG_PREFIX)),$(REPO_NAME))
endif
NPM_PKG_PREFIX  := $(call _ensure_hyphen,$(NPM_PKG_PREFIX))
NPM_ROOT_PKG_NAME := $(patsubst %-,%,$(NPM_PKG_PREFIX))

# ── Tooling paths ────────────────────────────────────────────────────────────
# Local installs of UniFFI bindgen tools, kept in .tools/ to avoid polluting
# the user's global cargo bin.
TOOLS_DIR         := $(CURDIR)/.tools
UNIFFI_BINDGEN    := $(TOOLS_DIR)/bin/uniffi-bindgen
UNIFFI_BINDGEN_GO := $(TOOLS_DIR)/bin/uniffi-bindgen-go
UNIFFI_BINDGEN_RN := $(TOOLS_DIR)/bin/uniffi-bindgen-react-native
UNIFFI_TOML_FILES := uniffi.toml $(wildcard crates/ffi/*/uniffi.toml)

# ── Configuration diagnostic ────────────────────────────────────────────────

config:
	@echo "resolved configuration:"
	@echo "  BINDINGS_DIR      = $(BINDINGS_DIR)"
	@echo "  CGO_ARCH          = $(CGO_ARCH)"
	@echo "  CGO_OS            = $(CGO_OS)"
	@echo "  FFI_LIBS          = $(FFI_LIBS)"
	@echo "  GIT_REPO          = $(GIT_REPO)"
	@echo "  GO_BINDINGS_DIR   = $(GO_BINDINGS_DIR)"
	@echo "  GO_MODULE_BASE    = $(GO_MODULE_BASE)"
	@echo "  JS_BINDINGS_DIR   = $(JS_BINDINGS_DIR)"
	@echo "  LIB_SUFFIX        = $(LIB_SUFFIX)"
	@echo "  LOG_LEVEL         = $(_LOG_LEVEL)"
	@echo "  NPM_PKG_PREFIX    = $(NPM_PKG_PREFIX)"
	@echo "  NPM_ROOT_PKG_NAME = $(NPM_ROOT_PKG_NAME)"
	@echo "  NPM_SCOPE         = $(NPM_SCOPE)"
	@echo "  NPM_TARBALLS_DIR  = $(NPM_TARBALLS_DIR)"
	@echo "  PKG_PREFIX        = $(PKG_PREFIX)"
	@echo "  PKG_PREFIX_US     = $(PKG_PREFIX_US)"
	@echo "  PKG_SUFFIX        = $(PKG_SUFFIX)"
	@echo "  PKG_SUFFIX_US     = $(PKG_SUFFIX_US)"
	@echo "  PY_BINDINGS_DIR   = $(PY_BINDINGS_DIR)"
	@echo "  PY_WHEELS_DIR     = $(PY_WHEELS_DIR)"
	@echo "  REPO_NAME         = $(REPO_NAME)"
	@echo "  REPO_NAME_US      = $(REPO_NAME_US)"
	@echo "  REPO_OWNER        = $(REPO_OWNER)"
	@echo "  REPO_SLUG         = $(REPO_SLUG)"
	@echo "  TOOLS_DIR         = $(TOOLS_DIR)"

# ── Setup ────────────────────────────────────────────────────────────────────

setup:
	git config core.hooksPath .githooks
	@echo "  git hooks configured to use .githooks/"

# ── Core build targets ───────────────────────────────────────────────────────

build:
	$(call _log_debug,cargo build --release)
	@cargo build --release

test:
	$(call _log_debug,cargo test --workspace)
	@cargo test --workspace

# ── Dependency bootstrapping ─────────────────────────────────────────────────
# These targets are no-ops when the required tools are already installed.

# Ensures maturin is available (installs via pipx → pip if needed)
python-dependencies:
	@$(_shell_log_helpers) \
	command -v maturin >/dev/null 2>&1 || { \
		echo "  maturin not found, installing via pipx"; \
		command -v pipx >/dev/null 2>&1 || { \
			echo "  pipx not found, installing via pip"; \
			if command -v pip >/dev/null 2>&1; then \
				log_debug "pip install --user pipx"; \
				pip install --user pipx; \
			elif python3 -m pip --version >/dev/null 2>&1; then \
				log_debug "python3 -m pip install --user pipx"; \
				python3 -m pip install --user pipx; \
			else \
				echo "  pip not found, installing"; \
				curl -sSL https://bootstrap.pypa.io/get-pip.py | python3; \
				log_debug "python3 -m pip install --user pipx"; \
				python3 -m pip install --user pipx; \
			fi; \
		}; \
		log_debug "pipx install maturin"; \
		pipx install maturin; \
	}

# Ensures Node.js is available (installs via nvm if needed)
install-node:
	@$(_shell_log_helpers) \
	command -v node >/dev/null 2>&1 || { \
		echo "  node not found, installing via nvm"; \
		export NVM_DIR="$${HOME}/.nvm"; \
		if [ ! -s "$$NVM_DIR/nvm.sh" ]; then \
			echo "  nvm not found, installing"; \
			curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash; \
		fi; \
		log_debug ". \$$NVM_DIR/nvm.sh && nvm install 20"; \
		. "$$NVM_DIR/nvm.sh" && nvm install 20; \
	}

# Ensures pnpm is available (installs via corepack if needed)
js-dependencies: install-node
	@$(_shell_log_helpers) \
	command -v pnpm >/dev/null 2>&1 || { \
		echo "  pnpm not found, installing via corepack"; \
		log_debug "corepack enable && corepack prepare pnpm@latest --activate"; \
		corepack enable && corepack prepare pnpm@latest --activate; \
	}

# ── UniFFI bindgen tool builds ──────────────────────────────────────────────
# Each target installs its bindgen binary into .tools/bin/ if not already present.

build-tools-py:
	@$(_shell_log_helpers) \
	if [ -x $(UNIFFI_BINDGEN) ]; then \
		log_debug "uniffi-bindgen already installed at $(UNIFFI_BINDGEN)"; \
	else \
		log_info "installing uniffi-bindgen into $(TOOLS_DIR)/"; \
		log_debug "cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features py-ffi --bin uniffi-bindgen"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features py-ffi --bin uniffi-bindgen; \
	fi

build-tools-go:
	@$(_shell_log_helpers) \
	if [ -x $(UNIFFI_BINDGEN_GO) ]; then \
		log_debug "uniffi-bindgen-go already installed at $(UNIFFI_BINDGEN_GO)"; \
	else \
		log_info "installing uniffi-bindgen-go into $(TOOLS_DIR)/"; \
		log_debug "cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features go-ffi --bin uniffi-bindgen-go"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features go-ffi --bin uniffi-bindgen-go; \
	fi

build-tools-wasm:
	@$(_shell_log_helpers) \
	if [ -x $(UNIFFI_BINDGEN_RN) ]; then \
		log_debug "uniffi-bindgen-react-native already installed at $(UNIFFI_BINDGEN_RN)"; \
	else \
		log_info "installing uniffi-bindgen-react-native into $(TOOLS_DIR)/"; \
		log_debug "cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features wasm-ffi --bin uniffi-bindgen-react-native"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features wasm-ffi --bin uniffi-bindgen-react-native; \
	fi

# ── Python bindings ─────────────────────────────────────────────────────────
# Generates a maturin-based Python project per FFI crate, then builds wheels
# into $(PY_WHEELS_DIR).  Package name pattern: $(PKG_PREFIX)<crate>$(PKG_SUFFIX)

generate-python: _check-repo-owner python-dependencies build-tools-py
	$(call _log_debug,PY_WHEELS_DIR = $(PY_WHEELS_DIR))
	$(call _log_debug,PY_BINDINGS_DIR = $(PY_BINDINGS_DIR))
	@mkdir -p $(PY_WHEELS_DIR)
	@$(_shell_log_helpers) \
	for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg_name=$(PKG_PREFIX)$${crate_dir}$(PKG_SUFFIX); \
		manifest_path=../../../crates/ffi/$$crate_dir/Cargo.toml; \
		proj_dir=$(PY_BINDINGS_DIR)/$$pkg_name; \
		log_debug "── $$ffi_lib ──"; \
		log_debug "crate dir:     $$crate_dir"; \
		log_debug "pkg name:      $$pkg_name"; \
		log_debug "manifest path: $$manifest_path"; \
		log_debug "proj dir:      $$proj_dir"; \
		mkdir -p $$proj_dir; \
		if [ ! -f $$proj_dir/pyproject.toml ]; then \
			log_info "creating python project for $$pkg_name"; \
			sed \
				-e "s|PKG_NAME|$$pkg_name|" \
				-e "s|MANIFEST_PATH|$$manifest_path|" \
				$(CURDIR)/templates/python/pyproject.toml > $$proj_dir/pyproject.toml; \
		fi; \
		mod_name=$$(echo $(PKG_PREFIX_US)$${ffi_lib}$(PKG_SUFFIX_US)); \
		log_debug "mod name:      $$mod_name"; \
		if grep -q "$$mod_name" $(CURDIR)/uniffi.toml 2>/dev/null && \
		   grep -A20 '\[bindings.python.external_packages\]' $(CURDIR)/uniffi.toml | grep -q "$$mod_name"; then \
			log_info "patching __init__.py for external package $$mod_name"; \
			mkdir -p $$proj_dir/$$mod_name; \
			sed "s|MOD_NAME|$$mod_name|g" $(CURDIR)/templates/python/init_external.py > $$proj_dir/$$mod_name/__init__.py; \
		fi; \
		log_info "building python bindings for $$ffi_lib"; \
		log_debug "cd $$proj_dir && PATH=$(TOOLS_DIR)/bin:\$$PATH maturin build --release --out $(PY_WHEELS_DIR)"; \
		( cd $$proj_dir && PATH=$(TOOLS_DIR)/bin:$$PATH maturin build --release --out $(PY_WHEELS_DIR) ); \
	done
	$(call _log_info,python wheels written to $(PY_WHEELS_DIR)/)

# Installs generated Python wheels into the current environment
install-python: generate-python
	@$(_shell_log_helpers) \
	for whl in $(PY_WHEELS_DIR)/*.whl; do \
		log_info "installing $$(basename $$whl)"; \
		log_debug "pip install --force-reinstall $$whl"; \
		pip install --force-reinstall $$whl; \
	done
	$(call _log_info,all python bindings installed)

# ── Go bindings ──────────────────────────────────────────────────────────────
# Generates Go modules with UniFFI-generated source and platform-specific static
# libs. Each module gets a go.mod and a CGO platform file from templates/.
# Package dir pattern: $(PKG_PREFIX_US)<lib>$(PKG_SUFFIX_US)

# Stamps the go_mod path into all uniffi.toml files so generated Go code imports
# from the correct module path.
stamp-uniffi-config:
	$(call _log_debug,GO_MODULE_BASE = $(GO_MODULE_BASE))
	$(call _log_debug,UNIFFI_TOML_FILES = $(UNIFFI_TOML_FILES))
	@$(_shell_log_helpers) \
	for f in $(UNIFFI_TOML_FILES); do \
		log_info "stamping go_mod in $$f"; \
		sed 's|go_mod = ".*"|go_mod = "$(GO_MODULE_BASE)"|' $$f > $$f.tmp && mv $$f.tmp $$f; \
	done

generate-go: _check-repo-owner build build-tools-go stamp-uniffi-config
	$(call _log_debug,GO_BINDINGS_DIR = $(GO_BINDINGS_DIR))
	$(call _log_debug,GO_MODULE_BASE = $(GO_MODULE_BASE))
	$(call _log_debug,CGO_OS = $(CGO_OS)  CGO_ARCH = $(CGO_ARCH))
	$(call _log_debug,LIB_SUFFIX = $(LIB_SUFFIX))
	$(call _log_debug,CGO_LDFLAGS_EXTRA = $(CGO_LDFLAGS_EXTRA))
	@mkdir -p $(GO_BINDINGS_DIR)
	@$(_shell_log_helpers) \
	for ffi_lib in $(FFI_LIBS); do \
		pkg=$(PKG_PREFIX_US)$${ffi_lib}$(PKG_SUFFIX_US); \
		lib=lib$${pkg}; \
		lib_dynamic=$${lib}$(LIB_SUFFIX); \
		lib_static=$${lib}.a; \
		log_info "generating go bindings for $$ffi_lib ($(CGO_OS)/$(CGO_ARCH))"; \
		log_debug "── $$ffi_lib ──"; \
		log_debug "pkg:         $$pkg"; \
		log_debug "lib:         $$lib"; \
		log_debug "lib_dynamic: $$lib_dynamic"; \
		log_debug "lib_static:  $$lib_static"; \
		log_debug "bindgen cmd: $(UNIFFI_BINDGEN_GO) --library target/release/$$lib_dynamic --config uniffi.toml --out-dir $(GO_BINDINGS_DIR)"; \
		$(UNIFFI_BINDGEN_GO) \
			--library target/release/$${lib_dynamic} \
			--config uniffi.toml \
			--out-dir $(GO_BINDINGS_DIR); \
		log_debug "copying $$lib_static → $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH)/"; \
		mkdir -p $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH); \
		cp target/release/$${lib_static} $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH)/; \
		if [ ! -f $(GO_BINDINGS_DIR)/$$pkg/go.mod ]; then \
			log_info "initializing go.mod for $$pkg ($(GO_MODULE_BASE)/$$pkg)"; \
			log_debug "cd $(GO_BINDINGS_DIR)/$$pkg && go mod init $(GO_MODULE_BASE)/$$pkg"; \
			( cd $(GO_BINDINGS_DIR)/$$pkg && go mod init $(GO_MODULE_BASE)/$$pkg ); \
		fi; \
		log_debug "writing cgo_$(CGO_OS)_$(CGO_ARCH).go from template"; \
		sed \
			-e "s|CGO_OS_CGO_ARCH|$(CGO_OS)_$(CGO_ARCH)|g" \
			-e "s|CGO_OS|$(CGO_OS)|g" \
			-e "s|CGO_ARCH|$(CGO_ARCH)|g" \
			-e "s|PKG_NAME|$$pkg|g" \
			-e "s|CGO_LDFLAGS_EXTRA|$(CGO_LDFLAGS_EXTRA)|g" \
			templates/go/cgo_platform.go \
		> $(GO_BINDINGS_DIR)/$$pkg/cgo_$(CGO_OS)_$(CGO_ARCH).go; \
	done
	$(call _log_info,go bindings written to $(GO_BINDINGS_DIR)/)

install-go: generate-go
	$(call _log_info,generating go.work for local module resolution)
	$(call _log_debug,GO_BINDINGS_DIR = $(GO_BINDINGS_DIR))
	@$(_shell_log_helpers) \
	cd $(GO_BINDINGS_DIR) && \
		log_debug "go work init"; \
		go work init && \
		for d in */go.mod; do \
			mod_dir=$$(dirname $$d); \
			log_debug "go work use $$mod_dir"; \
			go work use "$$mod_dir"; \
		done
	$(call _log_info,go packages ready at $(GO_BINDINGS_DIR)/)

# ── WASM / TypeScript bindings ───────────────────────────────────────────────
# Builds each FFI crate to wasm32-unknown-unknown, then runs uniffi-bindgen-
# react-native to produce TypeScript sources + .wasm blobs.
# npm package pattern: $(NPM_SCOPE)/$(NPM_PKG_PREFIX)<crate>$(PKG_SUFFIX)

generate-wasm-ts: _check-repo-owner build-tools-wasm
	$(call _log_debug,JS_BINDINGS_DIR = $(JS_BINDINGS_DIR))
	$(call _log_debug,NPM_SCOPE = $(NPM_SCOPE))
	$(call _log_debug,NPM_PKG_PREFIX = $(NPM_PKG_PREFIX))
	$(call _log_debug,NPM_ROOT_PKG_NAME = $(NPM_ROOT_PKG_NAME))
	@mkdir -p $(JS_BINDINGS_DIR)
	$(call _log_info,writing pnpm-workspace.yaml)
	@sed -e "s|NPM_SCOPE|$(NPM_SCOPE)|g" \
		$(CURDIR)/templates/js/pnpm-workspace.yaml > $(JS_BINDINGS_DIR)/pnpm-workspace.yaml
	$(call _log_info,writing root package.json ($(NPM_ROOT_PKG_NAME)))
	@sed -e "s|NPM_SCOPE|$(NPM_SCOPE)|g" \
		-e "s|ROOT_PKG_NAME|$(NPM_ROOT_PKG_NAME)|g" \
		$(CURDIR)/templates/js/root-package.json > $(JS_BINDINGS_DIR)/package.json
	@$(_shell_log_helpers) \
	for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg=$(PKG_PREFIX_US)$${ffi_lib}$(PKG_SUFFIX_US); \
		npm_pkg=$(NPM_SCOPE)/$(NPM_PKG_PREFIX)$${crate_dir}$(PKG_SUFFIX); \
		pkg_dir=$(JS_BINDINGS_DIR)/$${npm_pkg}; \
		log_info ""; \
		log_info "── $$ffi_lib ──"; \
		log_debug "crate dir:   $$crate_dir"; \
		log_debug "wasm pkg:    $$pkg"; \
		log_debug "npm pkg:     $$npm_pkg"; \
		log_debug "output dir:  $$pkg_dir"; \
		mkdir -p $$pkg_dir/src; \
		log_info "building wasm for $$pkg"; \
		log_debug "RUSTFLAGS=\"-C link-arg=--export-table -C link-arg=--growable-table\" cargo build --release --target wasm32-unknown-unknown -p $(PKG_PREFIX)$${crate_dir}$(PKG_SUFFIX) --features wasm"; \
		RUSTFLAGS="-C link-arg=--export-table -C link-arg=--growable-table" \
		cargo build --release --target wasm32-unknown-unknown \
			-p $(PKG_PREFIX)$${crate_dir}$(PKG_SUFFIX) --features wasm; \
		log_info "generating wasm/ts bindings for $$ffi_lib"; \
		log_debug "bindgen cmd: $(UNIFFI_BINDGEN_RN) generate wasm2 bindings --library target/wasm32-unknown-unknown/release/$$pkg.wasm --no-format --ts-dir $$pkg_dir/src"; \
		$(UNIFFI_BINDGEN_RN) \
			generate wasm2 bindings \
			--library target/wasm32-unknown-unknown/release/$${pkg}.wasm \
			--no-format \
			--ts-dir $$pkg_dir/src; \
		log_debug "copying $$pkg.wasm → $$pkg_dir/src/"; \
		cp target/wasm32-unknown-unknown/release/$${pkg}.wasm $$pkg_dir/src/; \
		crate_version=$$(grep '^version' $(CURDIR)/crates/ffi/$$crate_dir/Cargo.toml | head -1 | sed 's/.*"\(.*\)"/\1/'); \
		log_info "stamping $$npm_pkg@$$crate_version"; \
		sed -e "s|PKG_NAME|$$npm_pkg|" \
			-e "s|PKG_VERSION|$$crate_version|" \
			$(CURDIR)/templates/js/package.json > $$pkg_dir/package.json; \
		log_debug "copying tsconfig.json → $$pkg_dir/"; \
		cp $(CURDIR)/templates/js/tsconfig.json $$pkg_dir/tsconfig.json; \
	done
	$(call _log_info,fixing extensionless imports)
	@$(MAKE) fix-wasm-ts-imports
	$(call _log_info,wasm bindings written to $(JS_BINDINGS_DIR)/)

# The codegen emits extensionless relative imports. Add .ts extensions
# so tsc (with rewriteRelativeImportExtensions) rewrites them to .js
# in the compiled output. Skip specifiers that already have an extension.
fix-wasm-ts-imports:
	@$(_shell_log_helpers) \
	for f in $(JS_BINDINGS_DIR)/$(NPM_SCOPE)/*/src/*.ts; do \
		[ -f "$$f" ] || continue; \
		log_debug "fixing imports in $$f"; \
		sed -i -e "s/from '\.\/\([^'.]*\)'/from '.\/\1.ts'/g" \
		       -e 's/from "\.\/\([^".]*\)"/from ".\/\1.ts"/g' "$$f"; \
	done

# Compiles generated TypeScript sources
compile-wasm-ts: js-dependencies generate-wasm-ts
	$(call _log_info,installing JS workspace dependencies)
	$(call _log_debug,cd $(JS_BINDINGS_DIR) && pnpm install)
	@( cd $(JS_BINDINGS_DIR) && pnpm install )
	$(call _log_info,compiling typescript)
	@$(_shell_log_helpers) \
	for pkg_dir in $(JS_BINDINGS_DIR)/$(NPM_SCOPE)/*/; do \
		log_debug "cd $$pkg_dir && pnpm exec tsc"; \
		( cd "$$pkg_dir" && pnpm exec tsc ); \
	done
	$(call _log_info,typescript compiled at $(JS_BINDINGS_DIR)/)

# Packs each JS package into an npm tarball for local installation or publishing
generate-js: compile-wasm-ts
	$(call _log_debug,NPM_TARBALLS_DIR = $(NPM_TARBALLS_DIR))
	@mkdir -p $(NPM_TARBALLS_DIR)
	@$(_shell_log_helpers) \
	for pkg_dir in $(JS_BINDINGS_DIR)/$(NPM_SCOPE)/*/; do \
		log_info "packing $$(basename $$pkg_dir)"; \
		log_debug "cd $$pkg_dir && npm pack --pack-destination $(NPM_TARBALLS_DIR)"; \
		( cd "$$pkg_dir" && npm pack --pack-destination $(NPM_TARBALLS_DIR) ); \
	done
	$(call _log_info,npm tarballs written to $(NPM_TARBALLS_DIR)/)

# Installs packed tarballs into a target JS project, or globally if JS_PROJECT_DIR is not set
install-js: generate-js
	$(call _log_debug,JS_PROJECT_DIR = $(JS_PROJECT_DIR))
	@$(_shell_log_helpers) \
	if [ -n "$(JS_PROJECT_DIR)" ]; then \
		for tgz in $(NPM_TARBALLS_DIR)/*.tgz; do \
			log_info "installing $$(basename $$tgz) into $(JS_PROJECT_DIR)"; \
			log_debug "cd $(JS_PROJECT_DIR) && npm install $$tgz"; \
			( cd $(JS_PROJECT_DIR) && npm install $$tgz ); \
		done; \
		log_info "all js packages installed into $(JS_PROJECT_DIR)/"; \
	else \
		for tgz in $(NPM_TARBALLS_DIR)/*.tgz; do \
			log_info "installing $$(basename $$tgz) globally"; \
			log_debug "npm install -g $$tgz"; \
			npm install -g $$tgz; \
		done; \
		log_info "all js packages installed globally"; \
	fi

# ── Aggregate targets ────────────────────────────────────────────────────────

generate-all: generate-python generate-go generate-js

# ── Smoke tests ──────────────────────────────────────────────────────────────

test-python: install-python
	$(call _log_info,running python smoke test)
	$(call _log_debug,python examples/py-demo)
	@python examples/py-demo

test-go: install-go
	@$(_shell_log_helpers) \
	for pkg_dir in $(GO_BINDINGS_DIR)/*/; do \
		[ -f "$$pkg_dir/go.mod" ] || continue; \
		log_info "building $$(basename $$pkg_dir)"; \
		log_debug "cd $$pkg_dir && go build ./..."; \
		( cd "$$pkg_dir" && go build ./... ); \
	done
	$(call _log_info,running go smoke test)
	@$(_shell_log_helpers) \
	cd $(CURDIR)/examples/go-demo && \
		cp go.mod go.mod.bak && \
		for pkg_dir in $(GO_BINDINGS_DIR)/*/; do \
			[ -f "$$pkg_dir/go.mod" ] || continue; \
			mod_path=$$(head -1 "$$pkg_dir/go.mod" | awk '{print $$2}'); \
			log_debug "go mod edit -require=$$mod_path@v0.0.0 -replace=$$mod_path=$$(cd "$$pkg_dir" && pwd)"; \
			go mod edit -require="$$mod_path@v0.0.0" -replace="$$mod_path=$$(cd "$$pkg_dir" && pwd)"; \
		done && \
		log_debug "go run main.go"; \
		go run main.go; \
		status=$$?; \
		mv go.mod.bak go.mod; \
		rm -f go.sum; \
		exit $$status

test-js:
	$(call _log_debug,JS demo dir = $(CURDIR)/examples/js-demo)
	@cp $(CURDIR)/examples/js-demo/package.template.json $(CURDIR)/examples/js-demo/package.json
	$(MAKE) install-js JS_PROJECT_DIR=$(CURDIR)/examples/js-demo
	$(call _log_info,running js smoke test)
	$(call _log_debug,node --test examples/js-demo/main.mjs)
	@node --test examples/js-demo/main.mjs


# Builds the browser-based web-demo with Vite to verify wasm packages work in a frontend context
test-web:
	$(call _log_debug,web demo dir = $(CURDIR)/examples/web-demo)
	@cp $(CURDIR)/examples/web-demo/package.template.json $(CURDIR)/examples/web-demo/package.json
	$(MAKE) install-js JS_PROJECT_DIR=$(CURDIR)/examples/web-demo
	$(call _log_info,installing web-demo dev dependencies)
	$(call _log_debug,cd $(CURDIR)/examples/web-demo && npm install)
	@cd $(CURDIR)/examples/web-demo && npm install
	$(call _log_info,building web-demo (vite build))
	$(call _log_debug,cd $(CURDIR)/examples/web-demo && npx vite build)
	@cd $(CURDIR)/examples/web-demo && npx vite build
	$(call _log_info,web-demo build succeeded — run 'cd examples/web-demo && npx vite' to open in browser)

test-all: test-python test-go test-js test-web

# ── Cleanup ──────────────────────────────────────────────────────────────────
# clean            : cargo clean only (removes target/)
# clean-tools      : removes locally-installed bindgen tools
# clean-python     : removes generated Python binding sources
# clean-go         : removes generated Go bindings
# clean-js         : removes generated JS/WASM binding sources
# clean-bindings   : removes all generated binding sources (Python + Go + JS)
# clean-py-artifacts : removes built Python wheels (dist/pypi/)
# clean-js-artifacts : removes packed npm tarballs (dist/npm/)
# clean-artifacts  : removes all distributable artifacts (wheels + tarballs)
# clean-web        : removes web-demo generated files (node_modules, dist, package.json)
# clean-all        : cargo clean + bindings + artifacts
# purge            : everything including local tools

clean:
	$(call _log_debug,cargo clean)
	@cargo clean

clean-tools:
	rm -rf $(TOOLS_DIR)

clean-python:
	rm -rf $(PY_BINDINGS_DIR)

clean-go:
	rm -rf $(GO_BINDINGS_DIR)

clean-js:
	rm -rf $(JS_BINDINGS_DIR)

clean-bindings: clean-python clean-go clean-js

clean-py-artifacts:
	rm -rf $(PY_WHEELS_DIR)

clean-js-artifacts:
	rm -rf $(NPM_TARBALLS_DIR)

clean-artifacts: clean-py-artifacts clean-js-artifacts

clean-web:
	rm -rf $(CURDIR)/examples/web-demo/node_modules $(CURDIR)/examples/web-demo/dist
	rm -f $(CURDIR)/examples/web-demo/package.json $(CURDIR)/examples/web-demo/package-lock.json

clean-all: clean clean-bindings clean-artifacts clean-web

purge: clean-tools clean-all
