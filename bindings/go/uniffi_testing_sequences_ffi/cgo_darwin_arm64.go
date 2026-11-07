//go:build darwin && arm64

package uniffi_testing_sequences_ffi

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -luniffi_testing_sequences_ffi -lm -lpthread -framework Security
import "C"
