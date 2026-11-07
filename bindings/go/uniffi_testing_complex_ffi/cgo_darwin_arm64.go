//go:build darwin && arm64

package uniffi_testing_complex_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -luniffi_testing_complex_ffi -lm -lpthread -framework Security
import "C"
