//go:build linux && amd64

package uniffi_testing_types_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -luniffi_testing_types_ffi -lm -ldl -lpthread
import "C"
