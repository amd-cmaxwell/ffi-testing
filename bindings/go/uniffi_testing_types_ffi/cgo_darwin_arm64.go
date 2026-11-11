//go:build darwin && arm64

package uniffi_testing_types_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -luniffi_testing_types_ffi -lm -lpthread -framework Security
import "C"
