//go:build darwin && arm64

package uniffi_testing_linear_algebra_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -luniffi_testing_linear_algebra_ffi -lm -lpthread -framework Security
import "C"
