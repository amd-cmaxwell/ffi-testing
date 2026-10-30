package core

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_core_ffi"
)

func Add(a int32, b int32) int32 { return ffi.Add(a, b) }
