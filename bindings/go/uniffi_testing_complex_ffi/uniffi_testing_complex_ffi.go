package uniffi_testing_complex_ffi

// #cgo LDFLAGS: -L${SRCDIR} -luniffi_testing_complex_ffi -lm -ldl -lpthread
// #include <uniffi_testing_complex_ffi.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
	"unsafe"
)

// This is needed, because as of go 1.24
// type RustBuffer C.RustBuffer cannot have methods,
// RustBuffer is treated as non-local type
type GoRustBuffer struct {
	inner C.RustBuffer
}

type RustBufferI interface {
	AsReader() *bytes.Reader
	Free()
	ToGoBytes() []byte
	Data() unsafe.Pointer
	Len() uint64
	Capacity() uint64
}

// C.RustBuffer fields exposed as an interface so they can be accessed in different Go packages.
// See https://github.com/golang/go/issues/13467
type ExternalCRustBuffer interface {
	Data() unsafe.Pointer
	Len() uint64
	Capacity() uint64
}

func RustBufferFromC(b C.RustBuffer) ExternalCRustBuffer {
	return GoRustBuffer{
		inner: b,
	}
}

func CFromRustBuffer(b ExternalCRustBuffer) C.RustBuffer {
	return C.RustBuffer{
		capacity: C.uint64_t(b.Capacity()),
		len:      C.uint64_t(b.Len()),
		data:     (*C.uchar)(b.Data()),
	}
}

func RustBufferFromExternal(b ExternalCRustBuffer) GoRustBuffer {
	return GoRustBuffer{
		inner: C.RustBuffer{
			capacity: C.uint64_t(b.Capacity()),
			len:      C.uint64_t(b.Len()),
			data:     (*C.uchar)(b.Data()),
		},
	}
}

func (cb GoRustBuffer) Capacity() uint64 {
	return uint64(cb.inner.capacity)
}

func (cb GoRustBuffer) Len() uint64 {
	return uint64(cb.inner.len)
}

func (cb GoRustBuffer) Data() unsafe.Pointer {
	return unsafe.Pointer(cb.inner.data)
}

func (cb GoRustBuffer) AsReader() *bytes.Reader {
	b := unsafe.Slice((*byte)(cb.inner.data), C.uint64_t(cb.inner.len))
	return bytes.NewReader(b)
}

func (cb GoRustBuffer) Free() {
	rustCall(func(status *C.RustCallStatus) bool {
		C.ffi_uniffi_testing_complex_ffi_rustbuffer_free(cb.inner, status)
		return false
	})
}

func (cb GoRustBuffer) ToGoBytes() []byte {
	return C.GoBytes(unsafe.Pointer(cb.inner.data), C.int(cb.inner.len))
}

func stringToRustBuffer(str string) C.RustBuffer {
	return bytesToRustBuffer([]byte(str))
}

func bytesToRustBuffer(b []byte) C.RustBuffer {
	if len(b) == 0 {
		return C.RustBuffer{}
	}
	// We can pass the pointer along here, as it is pinned
	// for the duration of this call
	foreign := C.ForeignBytes{
		len:  C.int(len(b)),
		data: (*C.uchar)(unsafe.Pointer(&b[0])),
	}

	return rustCall(func(status *C.RustCallStatus) C.RustBuffer {
		return C.ffi_uniffi_testing_complex_ffi_rustbuffer_from_bytes(foreign, status)
	})
}

type BufLifter[GoType any] interface {
	Lift(value RustBufferI) GoType
}

type BufLowerer[GoType any] interface {
	Lower(value GoType) C.RustBuffer
}

type BufReader[GoType any] interface {
	Read(reader io.Reader) GoType
}

type BufWriter[GoType any] interface {
	Write(writer io.Writer, value GoType)
}

func LowerIntoRustBuffer[GoType any](bufWriter BufWriter[GoType], value GoType) C.RustBuffer {
	// This might be not the most efficient way but it does not require knowing allocation size
	// beforehand
	var buffer bytes.Buffer
	bufWriter.Write(&buffer, value)

	bytes, err := io.ReadAll(&buffer)
	if err != nil {
		panic(fmt.Errorf("reading written data: %w", err))
	}
	return bytesToRustBuffer(bytes)
}

func LiftFromRustBuffer[GoType any](bufReader BufReader[GoType], rbuf RustBufferI) GoType {
	defer rbuf.Free()
	reader := rbuf.AsReader()
	item := bufReader.Read(reader)
	if reader.Len() > 0 {
		// TODO: Remove this
		leftover, _ := io.ReadAll(reader)
		panic(fmt.Errorf("Junk remaining in buffer after lifting: %s", string(leftover)))
	}
	return item
}

func rustCallWithError[E any, U any](converter BufReader[*E], callback func(*C.RustCallStatus) U) (U, *E) {
	var status C.RustCallStatus
	returnValue := callback(&status)
	err := checkCallStatus(converter, status)
	return returnValue, err
}

func checkCallStatus[E any](converter BufReader[*E], status C.RustCallStatus) *E {
	switch status.code {
	case 0:
		return nil
	case 1:
		return LiftFromRustBuffer(converter, GoRustBuffer{inner: status.errorBuf})
	case 2:
		// when the rust code sees a panic, it tries to construct a rustBuffer
		// with the message.  but if that code panics, then it just sends back
		// an empty buffer.
		if status.errorBuf.len > 0 {
			panic(fmt.Errorf("%s", FfiConverterStringINSTANCE.Lift(GoRustBuffer{inner: status.errorBuf})))
		} else {
			panic(fmt.Errorf("Rust panicked while handling Rust panic"))
		}
	default:
		panic(fmt.Errorf("unknown status code: %d", status.code))
	}
}

func checkCallStatusUnknown(status C.RustCallStatus) error {
	switch status.code {
	case 0:
		return nil
	case 1:
		panic(fmt.Errorf("function not returning an error returned an error"))
	case 2:
		// when the rust code sees a panic, it tries to construct a C.RustBuffer
		// with the message.  but if that code panics, then it just sends back
		// an empty buffer.
		if status.errorBuf.len > 0 {
			panic(fmt.Errorf("%s", FfiConverterStringINSTANCE.Lift(GoRustBuffer{
				inner: status.errorBuf,
			})))
		} else {
			panic(fmt.Errorf("Rust panicked while handling Rust panic"))
		}
	default:
		return fmt.Errorf("unknown status code: %d", status.code)
	}
}

func rustCall[U any](callback func(*C.RustCallStatus) U) U {
	returnValue, err := rustCallWithError[error](nil, callback)
	if err != nil {
		panic(err)
	}
	return returnValue
}

type NativeError interface {
	AsError() error
}

func writeInt8(writer io.Writer, value int8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint8(writer io.Writer, value uint8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt16(writer io.Writer, value int16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint16(writer io.Writer, value uint16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt32(writer io.Writer, value int32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint32(writer io.Writer, value uint32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt64(writer io.Writer, value int64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint64(writer io.Writer, value uint64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat32(writer io.Writer, value float32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat64(writer io.Writer, value float64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func readInt8(reader io.Reader) int8 {
	var result int8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint8(reader io.Reader) uint8 {
	var result uint8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt16(reader io.Reader) int16 {
	var result int16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint16(reader io.Reader) uint16 {
	var result uint16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt32(reader io.Reader) int32 {
	var result int32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint32(reader io.Reader) uint32 {
	var result uint32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt64(reader io.Reader) int64 {
	var result int64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint64(reader io.Reader) uint64 {
	var result uint64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat32(reader io.Reader) float32 {
	var result float32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat64(reader io.Reader) float64 {
	var result float64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func init() {

	uniffiCheckChecksums()
}

func uniffiCheckChecksums() {
	// Get the bindings contract version from our ComponentInterface
	bindingsContractVersion := 29
	// Get the scaffolding contract version by calling the into the dylib
	scaffoldingContractVersion := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint32_t {
		return C.ffi_uniffi_testing_complex_ffi_uniffi_contract_version()
	})
	if bindingsContractVersion != int(scaffoldingContractVersion) {
		// If this happens try cleaning and rebuilding your project
		panic("uniffi_testing_complex_ffi: UniFFI contract version mismatch")
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_add()
		})
		if checksum != 8331 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_add: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_argument()
		})
		if checksum != 56874 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_argument: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_conjugate()
		})
		if checksum != 60397 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_conjugate: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_display()
		})
		if checksum != 42074 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_display: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_div()
		})
		if checksum != 47947 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_div: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_from_polar()
		})
		if checksum != 35957 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_from_polar: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_magnitude()
		})
		if checksum != 27223 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_magnitude: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_mul()
		})
		if checksum != 40097 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_mul: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_new()
		})
		if checksum != 57237 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_new: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_complex_sub()
		})
		if checksum != 47643 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_complex_sub: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_complex_ffi_checksum_func_quadratic_roots_complex()
		})
		if checksum != 19809 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_complex_ffi: uniffi_uniffi_testing_complex_ffi_checksum_func_quadratic_roots_complex: UniFFI API checksum mismatch")
		}
	}
}

type FfiConverterFloat64 struct{}

var FfiConverterFloat64INSTANCE = FfiConverterFloat64{}

func (FfiConverterFloat64) Lower(value float64) C.double {
	return C.double(value)
}

func (FfiConverterFloat64) Write(writer io.Writer, value float64) {
	writeFloat64(writer, value)
}

func (FfiConverterFloat64) Lift(value C.double) float64 {
	return float64(value)
}

func (FfiConverterFloat64) Read(reader io.Reader) float64 {
	return readFloat64(reader)
}

type FfiDestroyerFloat64 struct{}

func (FfiDestroyerFloat64) Destroy(_ float64) {}

type FfiConverterString struct{}

var FfiConverterStringINSTANCE = FfiConverterString{}

func (FfiConverterString) Lift(rb RustBufferI) string {
	defer rb.Free()
	reader := rb.AsReader()
	b, err := io.ReadAll(reader)
	if err != nil {
		panic(fmt.Errorf("reading reader: %w", err))
	}
	return string(b)
}

func (FfiConverterString) Read(reader io.Reader) string {
	length := readInt32(reader)
	buffer := make([]byte, length)
	read_length, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		panic(err)
	}
	if read_length != int(length) {
		panic(fmt.Errorf("bad read length when reading string, expected %d, read %d", length, read_length))
	}
	return string(buffer)
}

func (FfiConverterString) Lower(value string) C.RustBuffer {
	return stringToRustBuffer(value)
}

func (c FfiConverterString) LowerExternal(value string) ExternalCRustBuffer {
	return RustBufferFromC(stringToRustBuffer(value))
}

func (FfiConverterString) Write(writer io.Writer, value string) {
	if len(value) > math.MaxInt32 {
		panic("String is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	write_length, err := io.WriteString(writer, value)
	if err != nil {
		panic(err)
	}
	if write_length != len(value) {
		panic(fmt.Errorf("bad write length when writing string, expected %d, written %d", len(value), write_length))
	}
}

type FfiDestroyerString struct{}

func (FfiDestroyerString) Destroy(_ string) {}

type Complex struct {
	Re float64
	Im float64
}

func (r *Complex) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.Re)
	FfiDestroyerFloat64{}.Destroy(r.Im)
}

type FfiConverterComplex struct{}

var FfiConverterComplexINSTANCE = FfiConverterComplex{}

func (c FfiConverterComplex) Lift(rb RustBufferI) Complex {
	return LiftFromRustBuffer[Complex](c, rb)
}

func (c FfiConverterComplex) Read(reader io.Reader) Complex {
	return Complex{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterComplex) Lower(value Complex) C.RustBuffer {
	return LowerIntoRustBuffer[Complex](c, value)
}

func (c FfiConverterComplex) LowerExternal(value Complex) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[Complex](c, value))
}

func (c FfiConverterComplex) Write(writer io.Writer, value Complex) {
	FfiConverterFloat64INSTANCE.Write(writer, value.Re)
	FfiConverterFloat64INSTANCE.Write(writer, value.Im)
}

type FfiDestroyerComplex struct{}

func (_ FfiDestroyerComplex) Destroy(value Complex) {
	value.Destroy()
}

type ComplexPair struct {
	First  Complex
	Second Complex
}

func (r *ComplexPair) Destroy() {
	FfiDestroyerComplex{}.Destroy(r.First)
	FfiDestroyerComplex{}.Destroy(r.Second)
}

type FfiConverterComplexPair struct{}

var FfiConverterComplexPairINSTANCE = FfiConverterComplexPair{}

func (c FfiConverterComplexPair) Lift(rb RustBufferI) ComplexPair {
	return LiftFromRustBuffer[ComplexPair](c, rb)
}

func (c FfiConverterComplexPair) Read(reader io.Reader) ComplexPair {
	return ComplexPair{
		FfiConverterComplexINSTANCE.Read(reader),
		FfiConverterComplexINSTANCE.Read(reader),
	}
}

func (c FfiConverterComplexPair) Lower(value ComplexPair) C.RustBuffer {
	return LowerIntoRustBuffer[ComplexPair](c, value)
}

func (c FfiConverterComplexPair) LowerExternal(value ComplexPair) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[ComplexPair](c, value))
}

func (c FfiConverterComplexPair) Write(writer io.Writer, value ComplexPair) {
	FfiConverterComplexINSTANCE.Write(writer, value.First)
	FfiConverterComplexINSTANCE.Write(writer, value.Second)
}

type FfiDestroyerComplexPair struct{}

func (_ FfiDestroyerComplexPair) Destroy(value ComplexPair) {
	value.Destroy()
}

func ComplexAdd(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_add(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexArgument(c Complex) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_argument(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus)
	}))
}

func ComplexConjugate(c Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_conjugate(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus),
		}
	}))
}

func ComplexDisplay(c Complex) string {
	return FfiConverterStringINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_display(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus),
		}
	}))
}

func ComplexDiv(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_div(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexFromPolar(r float64, theta float64) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_from_polar(FfiConverterFloat64INSTANCE.Lower(r), FfiConverterFloat64INSTANCE.Lower(theta), _uniffiStatus),
		}
	}))
}

func ComplexMagnitude(c Complex) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_magnitude(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus)
	}))
}

func ComplexMul(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_mul(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexNew(re float64, im float64) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_new(FfiConverterFloat64INSTANCE.Lower(re), FfiConverterFloat64INSTANCE.Lower(im), _uniffiStatus),
		}
	}))
}

func ComplexSub(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_complex_sub(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func QuadraticRootsComplex(a float64, b float64, c float64) (ComplexPair, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_complex_ffi_fn_func_quadratic_roots_complex(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue ComplexPair
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterComplexPairINSTANCE.Lift(_uniffiRV), nil
	}
}
