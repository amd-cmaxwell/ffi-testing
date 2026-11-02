.PHONY: build test generate-python generate-go generate-all install-python install-go test-python test-go clean

FFI_CRATE_DIRS  := $(wildcard crates/ffi/*/Cargo.toml)
FFI_CRATE_PATHS := $(patsubst %/Cargo.toml,%,$(FFI_CRATE_DIRS))
FFI_CRATE_NAMES := $(notdir $(FFI_CRATE_PATHS))
FFI_LIBS        ?= $(sort $(subst -,_,$(FFI_CRATE_NAMES)))

BINDINGS_DIR    ?= $(CURDIR)/bindings
PY_BINDINGS_DIR ?= $(BINDINGS_DIR)/python
GO_BINDINGS_DIR ?= $(BINDINGS_DIR)/go

CGO_OS   ?= $(shell go env GOOS 2>/dev/null || echo linux)
CGO_ARCH ?= $(shell go env GOARCH 2>/dev/null || echo amd64)

LIB_SUFFIX_linux  := .so
LIB_SUFFIX_darwin := .dylib
LIB_SUFFIX        := $(LIB_SUFFIX_$(CGO_OS))

CGO_LDFLAGS_linux  := -lm -ldl -lpthread
CGO_LDFLAGS_darwin := -lm -lpthread -framework Security
CGO_LDFLAGS_EXTRA  := $(CGO_LDFLAGS_$(CGO_OS))

build:
	cargo build --release

test:
	cargo test --workspace

generate-python:
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg_name=uniffi-testing-$${crate_dir}-ffi; \
		manifest_path=../../../crates/ffi/$$crate_dir/Cargo.toml; \
		proj_dir=$(PY_BINDINGS_DIR)/$$pkg_name; \
		mkdir -p $$proj_dir; \
		if [ ! -f $$proj_dir/pyproject.toml ]; then \
			echo "  creating python project for $$pkg_name"; \
			sed \
				-e "s|PKG_NAME|$$pkg_name|" \
				-e "s|MANIFEST_PATH|$$manifest_path|" \
				templates/pyproject.toml > $$proj_dir/pyproject.toml; \
		fi; \
		echo "  building python bindings for $$ffi_lib"; \
		cd $$proj_dir && maturin build --release --out $$proj_dir/dist; \
	done
	@echo "  python bindings written to $(PY_BINDINGS_DIR)/"

install-python: generate-python
	@for ffi_lib in $(FFI_LIBS); do \
		crate_dir=$$(echo $$ffi_lib | tr '_' '-'); \
		pkg_name=uniffi-testing-$${crate_dir}-ffi; \
		echo "  installing python bindings for $$ffi_lib"; \
		pip install --force-reinstall $(PY_BINDINGS_DIR)/$$pkg_name/dist/*.whl; \
	done
	@echo "  all python bindings installed"

generate-go: build
	@mkdir -p $(GO_BINDINGS_DIR)
	@for ffi_lib in $(FFI_LIBS); do \
		echo "  generating go bindings for $$ffi_lib ($(CGO_OS)/$(CGO_ARCH))"; \
		pkg=uniffi_testing_$${ffi_lib}_ffi; \
		lib=lib$${pkg}; \
		lib_dynamic=$${lib}$(LIB_SUFFIX); \
		lib_static=$${lib}.a; \
		cargo run --release -p uniffi-testing-types-ffi --features ffi --bin uniffi-bindgen-go -- \
			--library target/release/$${lib_dynamic} \
			--config uniffi.toml \
			--out-dir $(GO_BINDINGS_DIR); \
		mkdir -p $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH); \
		cp target/release/$${lib_static} $(GO_BINDINGS_DIR)/$$pkg/lib/$(CGO_OS)_$(CGO_ARCH)/; \
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

generate-all: generate-python generate-go

test-python: install-python
	python examples/demo.py

test-go: install-go
	cd $(GO_BINDINGS_DIR) && go build ./...

clean:
	cargo clean
