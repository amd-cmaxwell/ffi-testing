.PHONY: build test \
       python-dependencies js-dependencies install-node \
       build-tools-py build-tools-go build-tools-wasm \
       generate-python generate-go generate-wasm-ts generate-all \
       install-python install-go install-js compile-wasm-ts generate-js \
       test-python test-go \
       clean clean-tools clean-python clean-go clean-js clean-bindings clean-all purge

PKG_PREFIX := uniffi-testing

FFI_CRATE_DIRS  := $(wildcard crates/ffi/*/Cargo.toml)
FFI_CRATE_PATHS := $(patsubst %/Cargo.toml,%,$(FFI_CRATE_DIRS))
FFI_CRATE_NAMES := $(filter-out tools,$(notdir $(FFI_CRATE_PATHS)))
FFI_LIBS        ?= $(sort $(subst -,_,$(FFI_CRATE_NAMES)))

BINDINGS_DIR      ?= $(CURDIR)/bindings
PY_BINDINGS_DIR   ?= $(BINDINGS_DIR)/python
GO_BINDINGS_DIR   ?= $(BINDINGS_DIR)/go
JS_BINDINGS_DIR   ?= $(BINDINGS_DIR)/js

GO_MODULE_BASE ?= github.com/amd-cmaxwell/ffi-testing/bindings/go
CGO_OS         ?= $(shell go env GOOS 2>/dev/null || echo linux)
CGO_ARCH       ?= $(shell go env GOARCH 2>/dev/null || echo amd64)

LIB_SUFFIX_linux  := .so
LIB_SUFFIX_darwin := .dylib
LIB_SUFFIX        := $(LIB_SUFFIX_$(CGO_OS))

CGO_LDFLAGS_linux  := -lm -ldl -lpthread
CGO_LDFLAGS_darwin := -lm -lpthread -framework Security
CGO_LDFLAGS_EXTRA  := $(CGO_LDFLAGS_$(CGO_OS))

NPM_SCOPE       ?= @amd-cmaxwell
NPM_PKG_PREFIX  ?= $(PKG_PREFIX)
NPM_TARBALLS_DIR ?= $(CURDIR)/dist/npm

TOOLS_DIR         := $(CURDIR)/.tools
UNIFFI_BINDGEN    := $(TOOLS_DIR)/bin/uniffi-bindgen
UNIFFI_BINDGEN_GO := $(TOOLS_DIR)/bin/uniffi-bindgen-go
UNIFFI_BINDGEN_RN := $(TOOLS_DIR)/bin/uniffi-bindgen-react-native

build:
	cargo build --release

test:
	cargo test --workspace

python-dependencies:
	@command -v maturin >/dev/null 2>&1 || { \
		echo "  maturin not found, installing via pipx"; \
		command -v pipx >/dev/null 2>&1 || { \
			echo "  pipx not found, installing via pip"; \
			if command -v pip >/dev/null 2>&1; then \
				pip install --user pipx; \
			elif python3 -m pip --version >/dev/null 2>&1; then \
				python3 -m pip install --user pipx; \
			else \
				echo "  pip not found, installing"; \
				curl -sSL https://bootstrap.pypa.io/get-pip.py | python3; \
				python3 -m pip install --user pipx; \
			fi; \
		}; \
		pipx install maturin; \
	}

install-node:
	@command -v node >/dev/null 2>&1 || { \
		echo "  node not found, installing via nvm"; \
		export NVM_DIR="$${HOME}/.nvm"; \
		if [ ! -s "$$NVM_DIR/nvm.sh" ]; then \
			echo "  nvm not found, installing"; \
			curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash; \
		fi; \
		. "$$NVM_DIR/nvm.sh" && nvm install 20; \
	}

js-dependencies: install-node
	@command -v pnpm >/dev/null 2>&1 || { \
		echo "  pnpm not found, installing via corepack"; \
		corepack enable && corepack prepare pnpm@latest --activate; \
	}

build-tools-py:
	@test -x $(UNIFFI_BINDGEN) || { \
		echo "  installing uniffi-bindgen into $(TOOLS_DIR)/"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features py-ffi --bin uniffi-bindgen; \
	}

build-tools-go:
	@test -x $(UNIFFI_BINDGEN_GO) || { \
		echo "  installing uniffi-bindgen-go into $(TOOLS_DIR)/"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features go-ffi --bin uniffi-bindgen-go; \
	}

build-tools-wasm:
	@test -x $(UNIFFI_BINDGEN_RN) || { \
		echo "  installing uniffi-bindgen-react-native into $(TOOLS_DIR)/"; \
		cargo install --path crates/ffi/tools --root $(TOOLS_DIR) --features wasm-ffi --bin uniffi-bindgen-react-native; \
	}

generate-python: python-dependencies build-tools-py
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg_name=$(PKG_PREFIX)-$${crate_dir}-ffi; \
		manifest_path=../../../crates/ffi/$$crate_dir/Cargo.toml; \
		proj_dir=$(PY_BINDINGS_DIR)/$$pkg_name; \
		mkdir -p $$proj_dir; \
		if [ ! -f $$proj_dir/pyproject.toml ]; then \
			echo "  creating python project for $$pkg_name"; \
			sed \
				-e "s|PKG_NAME|$$pkg_name|" \
				-e "s|MANIFEST_PATH|$$manifest_path|" \
				$(CURDIR)/templates/pyproject.toml > $$proj_dir/pyproject.toml; \
		fi; \
		echo "  building python bindings for $$ffi_lib"; \
		( cd $$proj_dir && PATH=$(TOOLS_DIR)/bin:$$PATH maturin build --release --out $$proj_dir/dist ); \
	done
	@echo "  python bindings written to $(PY_BINDINGS_DIR)/"

install-python: generate-python
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg_name=$(PKG_PREFIX)-$${crate_dir}-ffi; \
		echo "  installing python bindings for $$ffi_lib"; \
		pip install --force-reinstall $(PY_BINDINGS_DIR)/$$pkg_name/dist/*.whl; \
	done
	@echo "  all python bindings installed"

generate-go: build build-tools-go
	@mkdir -p $(GO_BINDINGS_DIR)
	@for ffi_lib in $(FFI_LIBS); do \
		echo "  generating go bindings for $$ffi_lib ($(CGO_OS)/$(CGO_ARCH))"; \
		pkg=uniffi_testing_$${ffi_lib}_ffi; \
		lib=lib$${pkg}; \
		lib_dynamic=$${lib}$(LIB_SUFFIX); \
		lib_static=$${lib}.a; \
		$(UNIFFI_BINDGEN_GO) \
			--library target/release/$${lib_dynamic} \
			--config uniffi.toml \
			--out-dir $(GO_BINDINGS_DIR); \
		mkdir -p $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH); \
		cp target/release/$${lib_static} $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH)/; \
		if [ ! -f $(GO_BINDINGS_DIR)/$$pkg/go.mod ]; then \
			( cd $(GO_BINDINGS_DIR)/$$pkg && go mod init $(GO_MODULE_BASE)/$$pkg ); \
		fi; \
		sed \
			-e "s|CGO_OS_CGO_ARCH|$(CGO_OS)_$(CGO_ARCH)|g" \
			-e "s|CGO_OS|$(CGO_OS)|g" \
			-e "s|CGO_ARCH|$(CGO_ARCH)|g" \
			-e "s|PKG_NAME|$$pkg|g" \
			-e "s|CGO_LDFLAGS_EXTRA|$(CGO_LDFLAGS_EXTRA)|g" \
			templates/cgo_platform.go \
		> $(GO_BINDINGS_DIR)/$$pkg/cgo_$(CGO_OS)_$(CGO_ARCH).go; \
	done
	@echo "  go bindings written to $(GO_BINDINGS_DIR)/"

install-go: generate-go
	@echo "  go packages ready at $(GO_BINDINGS_DIR)/"

generate-wasm-ts: build-tools-wasm
	@mkdir -p $(JS_BINDINGS_DIR)
	@cp $(CURDIR)/templates/wasm/pnpm-workspace.yaml $(JS_BINDINGS_DIR)/pnpm-workspace.yaml
	@cp $(CURDIR)/templates/wasm/root-package.json $(JS_BINDINGS_DIR)/package.json
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg=uniffi_testing_$${ffi_lib}_ffi; \
		npm_pkg=$(NPM_SCOPE)/$(NPM_PKG_PREFIX)-$${crate_dir}-ffi; \
		pkg_dir=$(JS_BINDINGS_DIR)/$${npm_pkg}; \
		mkdir -p $$pkg_dir/src; \
		echo "  building wasm for $$pkg"; \
		cargo build --release --target wasm32-unknown-unknown \
			-p $(PKG_PREFIX)-$${crate_dir}-ffi --features wasm; \
		echo "  generating wasm/ts bindings for $$ffi_lib"; \
		$(UNIFFI_BINDGEN_RN) \
			generate wasm2 bindings \
			--library target/wasm32-unknown-unknown/release/$${pkg}.wasm \
			--no-format \
			--ts-dir $$pkg_dir/src; \
		cp target/wasm32-unknown-unknown/release/$${pkg}.wasm $$pkg_dir/src/; \
	done
	@echo "  wasm bindings written to $(JS_BINDINGS_DIR)/"


compile-wasm-ts: js-dependencies generate-wasm-ts
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		npm_pkg=$(NPM_SCOPE)/$(NPM_PKG_PREFIX)-$${crate_dir}-ffi; \
		pkg_dir=$(JS_BINDINGS_DIR)/$${npm_pkg}; \
		crate_version=$$(grep '^version' $(CURDIR)/crates/ffi/$$crate_dir/Cargo.toml | head -1 | sed 's/.*"\(.*\)"/\1/'); \
		echo "  stamping $$npm_pkg@$$crate_version"; \
		sed -e "s|PKG_NAME|$$npm_pkg|" \
			-e "s|PKG_VERSION|$$crate_version|" \
			$(CURDIR)/templates/wasm/package.json > $$pkg_dir/package.json; \
		cp $(CURDIR)/templates/wasm/tsconfig.json $$pkg_dir/tsconfig.json; \
	done
	@echo "  compiling typescript"
	@for pkg_dir in $(JS_BINDINGS_DIR)/$(NPM_SCOPE)/*/; do \
		( cd "$$pkg_dir" && pnpm exec tsc ); \
	done
	@echo "  typescript compiled at $(JS_BINDINGS_DIR)/"

generate-js: compile-wasm-ts
	@mkdir -p $(NPM_TARBALLS_DIR)
	@for pkg_dir in $(JS_BINDINGS_DIR)/$(NPM_SCOPE)/*/; do \
		echo "  packing $$(basename $$pkg_dir)"; \
		( cd "$$pkg_dir" && npm pack --pack-destination $(NPM_TARBALLS_DIR) ); \
	done
	@echo "  npm tarballs written to $(NPM_TARBALLS_DIR)/"

install-js: generate-js
	@test -n "$(JS_PROJECT_DIR)" || { echo "error: JS_PROJECT_DIR is required (e.g. make install-js JS_PROJECT_DIR=../my-app)"; exit 1; }
	@for tgz in $(NPM_TARBALLS_DIR)/*.tgz; do \
		echo "  installing $$(basename $$tgz) into $(JS_PROJECT_DIR)"; \
		( cd $(JS_PROJECT_DIR) && npm install $$tgz ); \
	done
	@echo "  all js packages installed into $(JS_PROJECT_DIR)/"

generate-all: generate-python generate-go generate-js

test-python: install-python
	python examples/demo.py

test-go: install-go
	cd $(GO_BINDINGS_DIR) && go build ./...

clean:
	cargo clean

clean-tools:
	rm -rf $(TOOLS_DIR)

clean-python:
	rm -rf $(PY_BINDINGS_DIR)

clean-go:
	rm -rf $(GO_BINDINGS_DIR)

clean-js:
	rm -rf $(JS_BINDINGS_DIR)

clean-bindings: clean-python clean-go clean-js

clean-all: clean clean-bindings

purge: clean-tools clean-all
