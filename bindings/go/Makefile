REPO_ROOT := $(shell git rev-parse --show-toplevel)

.PHONY: build clean demo

build:
	cd $(REPO_ROOT) && cargo build --release -p uniffi-testing-ffi
	cp $(REPO_ROOT)/target/release/libuniffi_testing_ffi.a uniffi_testing_ffi/

clean:
	rm -f uniffi_testing_ffi/libuniffi_testing_ffi.a

demo: build
	go run ./example/demo/
