.PHONY: build test generate-python generate-go generate-all install-python install-go test-python test-go clean

FFI_CRATE_DIRS  := $(wildcard crates/ffi/*/Cargo.toml)
FFI_CRATE_PATHS := $(patsubst %/Cargo.toml,%,$(FFI_CRATE_DIRS))
FFI_CRATE_NAMES := $(notdir $(FFI_CRATE_PATHS))
FFI_LIBS        := $(sort $(subst -,_,$(FFI_CRATE_NAMES)))

build:
	cargo build --release

test:
	cargo test --workspace

generate-python: build
	@rm -rf /tmp/uniffi-python-gen
	@mkdir -p /tmp/uniffi-python-gen
	@for lib in $(FFI_LIBS); do \
		echo "  generating python bindings for $$lib"; \
		cargo run --release -p uniffi-testing-types-ffi --features ffi --bin uniffi-bindgen -- generate \
			--library target/release/libuniffi_testing_$${lib}_ffi.so \
			--language python --out-dir /tmp/uniffi-python-gen; \
	done
	@echo "  python bindings written to /tmp/uniffi-python-gen/"

install-python: generate-python
	@echo "  copying python bindings and shared libraries"
	@cp /tmp/uniffi-python-gen/*.py bindings/python/lib/uniffi_testing_ffi/_bindings/
	@for lib in $(FFI_LIBS); do \
		cp target/release/libuniffi_testing_$${lib}_ffi.so \
			bindings/python/lib/uniffi_testing_ffi/_bindings/; \
	done
	@pip install -e bindings/python
	@echo "  python package installed"

generate-go: build
	@rm -rf /tmp/uniffi-go-gen
	@mkdir -p /tmp/uniffi-go-gen
	@for lib in $(FFI_LIBS); do \
		echo "  generating go bindings for $$lib"; \
		cargo run --release -p uniffi-testing-types-ffi --features ffi --bin uniffi-bindgen-go -- \
			--library target/release/libuniffi_testing_$${lib}_ffi.so \
			--out-dir /tmp/uniffi-go-gen; \
	done
	@echo "  go bindings written to /tmp/uniffi-go-gen/"

GO_MODULE = github.com/amd-cmaxwell/ffi-testing/bindings/go

install-go: generate-go
	@echo "  copying go bindings and static libraries"
	@for lib in $(FFI_LIBS); do \
		pkg=uniffi_testing_$${lib}_ffi; \
		mkdir -p bindings/go/$$pkg; \
		cp /tmp/uniffi-go-gen/$$pkg/$$pkg.go bindings/go/$$pkg/; \
		cp /tmp/uniffi-go-gen/$$pkg/$$pkg.h bindings/go/$$pkg/; \
		cp target/release/libuniffi_testing_$${lib}_ffi.a bindings/go/$$pkg/; \
		sed -i 's|// #include|// #cgo LDFLAGS: -L$${SRCDIR} -luniffi_testing_'"$${lib}"'_ffi -lm -ldl -lpthread\n// #include|' \
			bindings/go/$$pkg/$$pkg.go; \
		sed -i 's|"uniffi_testing_types_ffi"|"$(GO_MODULE)/uniffi_testing_types_ffi"|g' \
			bindings/go/$$pkg/$$pkg.go; \
	done
	@echo "  go packages installed"

generate-all: generate-python generate-go

test-python: install-python
	python examples/demo.py

test-go: install-go
	cd bindings/go && go build ./...

clean:
	cargo clean
