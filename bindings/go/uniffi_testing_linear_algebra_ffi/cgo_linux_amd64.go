//go:build linux && amd64

package uniffi_testing_linear_algebra_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -luniffi_testing_linear_algebra_ffi -lm -ldl -lpthread
import "C"
