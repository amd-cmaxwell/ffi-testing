.PHONY: build test generate-python generate-go generate-all test-python test-go clean

build:
	cargo build --release -p uniffi-testing-ffi

test:
	cargo test --workspace

generate-python-modules: build
	cargo run --release --bin uniffi-bindgen -- generate \
		--library target/release/libuniffi_testing_ffi.so \
		--language python --out-dir /tmp/uniffi-python-gen
	python3 bindings/python/scripts/generate_modules.py /tmp/uniffi-python-gen/uniffi_testing_ffi.py

generate-python: build generate-python-modules
	cd bindings/python && maturin develop --release

generate-go: build
	uniffi-bindgen-go --library target/release/libuniffi_testing_ffi.so --out-dir /tmp/uniffi-go-gen
	cp /tmp/uniffi-go-gen/uniffi_testing_ffi/uniffi_testing_ffi.go bindings/go/uniffi_testing_ffi/
	cp /tmp/uniffi-go-gen/uniffi_testing_ffi/uniffi_testing_ffi.h bindings/go/uniffi_testing_ffi/
	cp target/release/libuniffi_testing_ffi.a bindings/go/uniffi_testing_ffi/
	sed -i 's|// #include <uniffi_testing_ffi.h>|// #cgo LDFLAGS: -L$${SRCDIR} -luniffi_testing_ffi -lm -ldl -lpthread\n// #include <uniffi_testing_ffi.h>|' bindings/go/uniffi_testing_ffi/uniffi_testing_ffi.go

generate-all: generate-python generate-go

test-python: generate-python
	python examples/demo.py

test-go:
	cd bindings/go && go run ./example/demo/

clean:
	cargo clean
	rm -rf bindings/python/target
