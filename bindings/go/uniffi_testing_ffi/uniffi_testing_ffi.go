package uniffi_testing_ffi

// #cgo LDFLAGS: -L${SRCDIR} -luniffi_testing_ffi -lm -ldl -lpthread
// #include <uniffi_testing_ffi.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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
		C.ffi_uniffi_testing_ffi_rustbuffer_free(cb.inner, status)
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
		return C.ffi_uniffi_testing_ffi_rustbuffer_from_bytes(foreign, status)
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
		return C.ffi_uniffi_testing_ffi_uniffi_contract_version()
	})
	if bindingsContractVersion != int(scaffoldingContractVersion) {
		// If this happens try cleaning and rebuilding your project
		panic("uniffi_testing_ffi: UniFFI contract version mismatch")
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_add()
		})
		if checksum != 19761 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_add: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_arithmetic_sum()
		})
		if checksum != 15501 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_arithmetic_sum: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_arithmetic_term()
		})
		if checksum != 61609 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_arithmetic_term: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_circle_area()
		})
		if checksum != 31582 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_circle_area: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_circle_circumference()
		})
		if checksum != 17655 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_circle_circumference: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_combinations()
		})
		if checksum != 17739 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_combinations: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_add()
		})
		if checksum != 20051 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_add: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_argument()
		})
		if checksum != 64529 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_argument: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_conjugate()
		})
		if checksum != 35988 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_conjugate: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_display()
		})
		if checksum != 27305 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_display: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_div()
		})
		if checksum != 45406 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_div: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_from_polar()
		})
		if checksum != 13900 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_from_polar: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_magnitude()
		})
		if checksum != 35072 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_magnitude: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_mul()
		})
		if checksum != 13239 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_mul: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_new()
		})
		if checksum != 15280 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_new: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_complex_sub()
		})
		if checksum != 57267 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_complex_sub: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_cross_product_3d()
		})
		if checksum != 18556 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_cross_product_3d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_degrees_to_radians()
		})
		if checksum != 23307 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_degrees_to_radians: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_determinant_2x2()
		})
		if checksum != 35624 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_determinant_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_determinant_3x3()
		})
		if checksum != 544 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_determinant_3x3: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_distance_2d()
		})
		if checksum != 34209 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_distance_2d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_distance_3d()
		})
		if checksum != 35912 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_distance_3d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_dot_product()
		})
		if checksum != 29989 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_dot_product: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_factorial()
		})
		if checksum != 22251 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_factorial: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_fibonacci()
		})
		if checksum != 55122 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_fibonacci: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_gcd()
		})
		if checksum != 20820 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_gcd: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_geometric_sum()
		})
		if checksum != 35321 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_geometric_sum: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_geometric_sum_infinite()
		})
		if checksum != 64413 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_geometric_sum_infinite: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_geometric_term()
		})
		if checksum != 13330 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_geometric_term: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_is_prime()
		})
		if checksum != 26969 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_is_prime: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_law_of_cosines_angle()
		})
		if checksum != 62990 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_law_of_cosines_angle: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_law_of_cosines_side()
		})
		if checksum != 60299 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_law_of_cosines_side: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_law_of_sines_side()
		})
		if checksum != 55508 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_law_of_sines_side: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_lcm()
		})
		if checksum != 48366 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_lcm: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_linear_root()
		})
		if checksum != 53951 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_linear_root: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_matrix_inverse_2x2()
		})
		if checksum != 53505 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_matrix_inverse_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_matrix_multiply_2x2()
		})
		if checksum != 8287 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_matrix_multiply_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_midpoint_2d()
		})
		if checksum != 25997 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_midpoint_2d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_permutations()
		})
		if checksum != 11958 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_permutations: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_polynomial_add()
		})
		if checksum != 24660 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_polynomial_add: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_polynomial_derivative()
		})
		if checksum != 43011 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_polynomial_derivative: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_polynomial_evaluate()
		})
		if checksum != 52683 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_polynomial_evaluate: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_polynomial_multiply()
		})
		if checksum != 34187 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_polynomial_multiply: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_quadratic_roots()
		})
		if checksum != 45306 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_quadratic_roots: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_quadratic_roots_complex()
		})
		if checksum != 57483 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_quadratic_roots_complex: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_radians_to_degrees()
		})
		if checksum != 64310 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_radians_to_degrees: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_solve_2x2()
		})
		if checksum != 15807 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_solve_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_solve_3x3()
		})
		if checksum != 4838 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_solve_3x3: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_transpose_2x2()
		})
		if checksum != 51104 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_transpose_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_ffi_checksum_func_triangle_area_heron()
		})
		if checksum != 52634 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_ffi: uniffi_uniffi_testing_ffi_checksum_func_triangle_area_heron: UniFFI API checksum mismatch")
		}
	}
}

type FfiConverterInt32 struct{}

var FfiConverterInt32INSTANCE = FfiConverterInt32{}

func (FfiConverterInt32) Lower(value int32) C.int32_t {
	return C.int32_t(value)
}

func (FfiConverterInt32) Write(writer io.Writer, value int32) {
	writeInt32(writer, value)
}

func (FfiConverterInt32) Lift(value C.int32_t) int32 {
	return int32(value)
}

func (FfiConverterInt32) Read(reader io.Reader) int32 {
	return readInt32(reader)
}

type FfiDestroyerInt32 struct{}

func (FfiDestroyerInt32) Destroy(_ int32) {}

type FfiConverterUint64 struct{}

var FfiConverterUint64INSTANCE = FfiConverterUint64{}

func (FfiConverterUint64) Lower(value uint64) C.uint64_t {
	return C.uint64_t(value)
}

func (FfiConverterUint64) Write(writer io.Writer, value uint64) {
	writeUint64(writer, value)
}

func (FfiConverterUint64) Lift(value C.uint64_t) uint64 {
	return uint64(value)
}

func (FfiConverterUint64) Read(reader io.Reader) uint64 {
	return readUint64(reader)
}

type FfiDestroyerUint64 struct{}

func (FfiDestroyerUint64) Destroy(_ uint64) {}

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

type FfiConverterBool struct{}

var FfiConverterBoolINSTANCE = FfiConverterBool{}

func (FfiConverterBool) Lower(value bool) C.int8_t {
	if value {
		return C.int8_t(1)
	}
	return C.int8_t(0)
}

func (FfiConverterBool) Write(writer io.Writer, value bool) {
	if value {
		writeInt8(writer, 1)
	} else {
		writeInt8(writer, 0)
	}
}

func (FfiConverterBool) Lift(value C.int8_t) bool {
	return value != 0
}

func (FfiConverterBool) Read(reader io.Reader) bool {
	return readInt8(reader) != 0
}

type FfiDestroyerBool struct{}

func (FfiDestroyerBool) Destroy(_ bool) {}

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

type F64Pair struct {
	First  float64
	Second float64
}

func (r *F64Pair) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.First)
	FfiDestroyerFloat64{}.Destroy(r.Second)
}

type FfiConverterF64Pair struct{}

var FfiConverterF64PairINSTANCE = FfiConverterF64Pair{}

func (c FfiConverterF64Pair) Lift(rb RustBufferI) F64Pair {
	return LiftFromRustBuffer[F64Pair](c, rb)
}

func (c FfiConverterF64Pair) Read(reader io.Reader) F64Pair {
	return F64Pair{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterF64Pair) Lower(value F64Pair) C.RustBuffer {
	return LowerIntoRustBuffer[F64Pair](c, value)
}

func (c FfiConverterF64Pair) LowerExternal(value F64Pair) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[F64Pair](c, value))
}

func (c FfiConverterF64Pair) Write(writer io.Writer, value F64Pair) {
	FfiConverterFloat64INSTANCE.Write(writer, value.First)
	FfiConverterFloat64INSTANCE.Write(writer, value.Second)
}

type FfiDestroyerF64Pair struct{}

func (_ FfiDestroyerF64Pair) Destroy(value F64Pair) {
	value.Destroy()
}

type Matrix2x2 struct {
	M00 float64
	M01 float64
	M10 float64
	M11 float64
}

func (r *Matrix2x2) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.M00)
	FfiDestroyerFloat64{}.Destroy(r.M01)
	FfiDestroyerFloat64{}.Destroy(r.M10)
	FfiDestroyerFloat64{}.Destroy(r.M11)
}

type FfiConverterMatrix2x2 struct{}

var FfiConverterMatrix2x2INSTANCE = FfiConverterMatrix2x2{}

func (c FfiConverterMatrix2x2) Lift(rb RustBufferI) Matrix2x2 {
	return LiftFromRustBuffer[Matrix2x2](c, rb)
}

func (c FfiConverterMatrix2x2) Read(reader io.Reader) Matrix2x2 {
	return Matrix2x2{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterMatrix2x2) Lower(value Matrix2x2) C.RustBuffer {
	return LowerIntoRustBuffer[Matrix2x2](c, value)
}

func (c FfiConverterMatrix2x2) LowerExternal(value Matrix2x2) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[Matrix2x2](c, value))
}

func (c FfiConverterMatrix2x2) Write(writer io.Writer, value Matrix2x2) {
	FfiConverterFloat64INSTANCE.Write(writer, value.M00)
	FfiConverterFloat64INSTANCE.Write(writer, value.M01)
	FfiConverterFloat64INSTANCE.Write(writer, value.M10)
	FfiConverterFloat64INSTANCE.Write(writer, value.M11)
}

type FfiDestroyerMatrix2x2 struct{}

func (_ FfiDestroyerMatrix2x2) Destroy(value Matrix2x2) {
	value.Destroy()
}

type Matrix3x3 struct {
	M00 float64
	M01 float64
	M02 float64
	M10 float64
	M11 float64
	M12 float64
	M20 float64
	M21 float64
	M22 float64
}

func (r *Matrix3x3) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.M00)
	FfiDestroyerFloat64{}.Destroy(r.M01)
	FfiDestroyerFloat64{}.Destroy(r.M02)
	FfiDestroyerFloat64{}.Destroy(r.M10)
	FfiDestroyerFloat64{}.Destroy(r.M11)
	FfiDestroyerFloat64{}.Destroy(r.M12)
	FfiDestroyerFloat64{}.Destroy(r.M20)
	FfiDestroyerFloat64{}.Destroy(r.M21)
	FfiDestroyerFloat64{}.Destroy(r.M22)
}

type FfiConverterMatrix3x3 struct{}

var FfiConverterMatrix3x3INSTANCE = FfiConverterMatrix3x3{}

func (c FfiConverterMatrix3x3) Lift(rb RustBufferI) Matrix3x3 {
	return LiftFromRustBuffer[Matrix3x3](c, rb)
}

func (c FfiConverterMatrix3x3) Read(reader io.Reader) Matrix3x3 {
	return Matrix3x3{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterMatrix3x3) Lower(value Matrix3x3) C.RustBuffer {
	return LowerIntoRustBuffer[Matrix3x3](c, value)
}

func (c FfiConverterMatrix3x3) LowerExternal(value Matrix3x3) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[Matrix3x3](c, value))
}

func (c FfiConverterMatrix3x3) Write(writer io.Writer, value Matrix3x3) {
	FfiConverterFloat64INSTANCE.Write(writer, value.M00)
	FfiConverterFloat64INSTANCE.Write(writer, value.M01)
	FfiConverterFloat64INSTANCE.Write(writer, value.M02)
	FfiConverterFloat64INSTANCE.Write(writer, value.M10)
	FfiConverterFloat64INSTANCE.Write(writer, value.M11)
	FfiConverterFloat64INSTANCE.Write(writer, value.M12)
	FfiConverterFloat64INSTANCE.Write(writer, value.M20)
	FfiConverterFloat64INSTANCE.Write(writer, value.M21)
	FfiConverterFloat64INSTANCE.Write(writer, value.M22)
}

type FfiDestroyerMatrix3x3 struct{}

func (_ FfiDestroyerMatrix3x3) Destroy(value Matrix3x3) {
	value.Destroy()
}

type Point2 struct {
	X float64
	Y float64
}

func (r *Point2) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.X)
	FfiDestroyerFloat64{}.Destroy(r.Y)
}

type FfiConverterPoint2 struct{}

var FfiConverterPoint2INSTANCE = FfiConverterPoint2{}

func (c FfiConverterPoint2) Lift(rb RustBufferI) Point2 {
	return LiftFromRustBuffer[Point2](c, rb)
}

func (c FfiConverterPoint2) Read(reader io.Reader) Point2 {
	return Point2{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterPoint2) Lower(value Point2) C.RustBuffer {
	return LowerIntoRustBuffer[Point2](c, value)
}

func (c FfiConverterPoint2) LowerExternal(value Point2) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[Point2](c, value))
}

func (c FfiConverterPoint2) Write(writer io.Writer, value Point2) {
	FfiConverterFloat64INSTANCE.Write(writer, value.X)
	FfiConverterFloat64INSTANCE.Write(writer, value.Y)
}

type FfiDestroyerPoint2 struct{}

func (_ FfiDestroyerPoint2) Destroy(value Point2) {
	value.Destroy()
}

type Point3 struct {
	X float64
	Y float64
	Z float64
}

func (r *Point3) Destroy() {
	FfiDestroyerFloat64{}.Destroy(r.X)
	FfiDestroyerFloat64{}.Destroy(r.Y)
	FfiDestroyerFloat64{}.Destroy(r.Z)
}

type FfiConverterPoint3 struct{}

var FfiConverterPoint3INSTANCE = FfiConverterPoint3{}

func (c FfiConverterPoint3) Lift(rb RustBufferI) Point3 {
	return LiftFromRustBuffer[Point3](c, rb)
}

func (c FfiConverterPoint3) Read(reader io.Reader) Point3 {
	return Point3{
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
		FfiConverterFloat64INSTANCE.Read(reader),
	}
}

func (c FfiConverterPoint3) Lower(value Point3) C.RustBuffer {
	return LowerIntoRustBuffer[Point3](c, value)
}

func (c FfiConverterPoint3) LowerExternal(value Point3) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[Point3](c, value))
}

func (c FfiConverterPoint3) Write(writer io.Writer, value Point3) {
	FfiConverterFloat64INSTANCE.Write(writer, value.X)
	FfiConverterFloat64INSTANCE.Write(writer, value.Y)
	FfiConverterFloat64INSTANCE.Write(writer, value.Z)
}

type FfiDestroyerPoint3 struct{}

func (_ FfiDestroyerPoint3) Destroy(value Point3) {
	value.Destroy()
}

type MathError struct {
	err error
}

// Convience method to turn *MathError into error
// Avoiding treating nil pointer as non nil error interface
func (err *MathError) AsError() error {
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (err MathError) Error() string {
	return fmt.Sprintf("MathError: %s", err.err.Error())
}

func (err MathError) Unwrap() error {
	return err.err
}

// Err* are used for checking error type with `errors.Is`
var ErrMathErrorDivisionByZero = fmt.Errorf("MathErrorDivisionByZero")
var ErrMathErrorInvalidInput = fmt.Errorf("MathErrorInvalidInput")
var ErrMathErrorNoRealSolution = fmt.Errorf("MathErrorNoRealSolution")
var ErrMathErrorOverflow = fmt.Errorf("MathErrorOverflow")

// Variant structs
type MathErrorDivisionByZero struct {
}

func NewMathErrorDivisionByZero() *MathError {
	return &MathError{err: &MathErrorDivisionByZero{}}
}

func (e MathErrorDivisionByZero) destroy() {
}

func (err MathErrorDivisionByZero) Error() string {
	return fmt.Sprint("DivisionByZero")
}

func (self MathErrorDivisionByZero) Is(target error) bool {
	return target == ErrMathErrorDivisionByZero
}

type MathErrorInvalidInput struct {
	Message string
}

func NewMathErrorInvalidInput(
	message string,
) *MathError {
	return &MathError{err: &MathErrorInvalidInput{
		Message: message}}
}

func (e MathErrorInvalidInput) destroy() {
	FfiDestroyerString{}.Destroy(e.Message)
}

func (err MathErrorInvalidInput) Error() string {
	return fmt.Sprint("InvalidInput",
		": ",

		"Message=",
		err.Message,
	)
}

func (self MathErrorInvalidInput) Is(target error) bool {
	return target == ErrMathErrorInvalidInput
}

type MathErrorNoRealSolution struct {
}

func NewMathErrorNoRealSolution() *MathError {
	return &MathError{err: &MathErrorNoRealSolution{}}
}

func (e MathErrorNoRealSolution) destroy() {
}

func (err MathErrorNoRealSolution) Error() string {
	return fmt.Sprint("NoRealSolution")
}

func (self MathErrorNoRealSolution) Is(target error) bool {
	return target == ErrMathErrorNoRealSolution
}

type MathErrorOverflow struct {
}

func NewMathErrorOverflow() *MathError {
	return &MathError{err: &MathErrorOverflow{}}
}

func (e MathErrorOverflow) destroy() {
}

func (err MathErrorOverflow) Error() string {
	return fmt.Sprint("Overflow")
}

func (self MathErrorOverflow) Is(target error) bool {
	return target == ErrMathErrorOverflow
}

type FfiConverterMathError struct{}

var FfiConverterMathErrorINSTANCE = FfiConverterMathError{}

func (c FfiConverterMathError) Lift(eb RustBufferI) *MathError {
	return LiftFromRustBuffer[*MathError](c, eb)
}

func (c FfiConverterMathError) Lower(value *MathError) C.RustBuffer {
	return LowerIntoRustBuffer[*MathError](c, value)
}

func (c FfiConverterMathError) LowerExternal(value *MathError) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[*MathError](c, value))
}

func (c FfiConverterMathError) Read(reader io.Reader) *MathError {
	errorID := readUint32(reader)

	switch errorID {
	case 1:
		return &MathError{&MathErrorDivisionByZero{}}
	case 2:
		return &MathError{&MathErrorInvalidInput{
			Message: FfiConverterStringINSTANCE.Read(reader),
		}}
	case 3:
		return &MathError{&MathErrorNoRealSolution{}}
	case 4:
		return &MathError{&MathErrorOverflow{}}
	default:
		panic(fmt.Sprintf("Unknown error code %d in FfiConverterMathError.Read()", errorID))
	}
}

func (c FfiConverterMathError) Write(writer io.Writer, value *MathError) {
	switch variantValue := value.err.(type) {
	case *MathErrorDivisionByZero:
		writeInt32(writer, 1)
	case *MathErrorInvalidInput:
		writeInt32(writer, 2)
		FfiConverterStringINSTANCE.Write(writer, variantValue.Message)
	case *MathErrorNoRealSolution:
		writeInt32(writer, 3)
	case *MathErrorOverflow:
		writeInt32(writer, 4)
	default:
		_ = variantValue
		panic(fmt.Sprintf("invalid error value `%v` in FfiConverterMathError.Write", value))
	}
}

type FfiDestroyerMathError struct{}

func (_ FfiDestroyerMathError) Destroy(value *MathError) {
	switch variantValue := value.err.(type) {
	case MathErrorDivisionByZero:
		variantValue.destroy()
	case MathErrorInvalidInput:
		variantValue.destroy()
	case MathErrorNoRealSolution:
		variantValue.destroy()
	case MathErrorOverflow:
		variantValue.destroy()
	default:
		_ = variantValue
		panic(fmt.Sprintf("invalid error value `%v` in FfiDestroyerMathError.Destroy", value))
	}
}

type FfiConverterSequenceFloat64 struct{}

var FfiConverterSequenceFloat64INSTANCE = FfiConverterSequenceFloat64{}

func (c FfiConverterSequenceFloat64) Lift(rb RustBufferI) []float64 {
	return LiftFromRustBuffer[[]float64](c, rb)
}

func (c FfiConverterSequenceFloat64) Read(reader io.Reader) []float64 {
	length := readInt32(reader)
	if length == 0 {
		return nil
	}
	result := make([]float64, 0, length)
	for i := int32(0); i < length; i++ {
		result = append(result, FfiConverterFloat64INSTANCE.Read(reader))
	}
	return result
}

func (c FfiConverterSequenceFloat64) Lower(value []float64) C.RustBuffer {
	return LowerIntoRustBuffer[[]float64](c, value)
}

func (c FfiConverterSequenceFloat64) LowerExternal(value []float64) ExternalCRustBuffer {
	return RustBufferFromC(LowerIntoRustBuffer[[]float64](c, value))
}

func (c FfiConverterSequenceFloat64) Write(writer io.Writer, value []float64) {
	if len(value) > math.MaxInt32 {
		panic("[]float64 is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	for _, item := range value {
		FfiConverterFloat64INSTANCE.Write(writer, item)
	}
}

type FfiDestroyerSequenceFloat64 struct{}

func (FfiDestroyerSequenceFloat64) Destroy(sequence []float64) {
	for _, value := range sequence {
		FfiDestroyerFloat64{}.Destroy(value)
	}
}

func Add(a int32, b int32) int32 {
	return FfiConverterInt32INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.int32_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_add(FfiConverterInt32INSTANCE.Lower(a), FfiConverterInt32INSTANCE.Lower(b), _uniffiStatus)
	}))
}

func ArithmeticSum(a1 float64, d float64, n uint64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_arithmetic_sum(FfiConverterFloat64INSTANCE.Lower(a1), FfiConverterFloat64INSTANCE.Lower(d), FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	}))
}

func ArithmeticTerm(a1 float64, d float64, n uint64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_arithmetic_term(FfiConverterFloat64INSTANCE.Lower(a1), FfiConverterFloat64INSTANCE.Lower(d), FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	}))
}

func CircleArea(radius float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_circle_area(FfiConverterFloat64INSTANCE.Lower(radius), _uniffiStatus)
	}))
}

func CircleCircumference(radius float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_circle_circumference(FfiConverterFloat64INSTANCE.Lower(radius), _uniffiStatus)
	}))
}

func Combinations(n uint64, r uint64) (uint64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_combinations(FfiConverterUint64INSTANCE.Lower(n), FfiConverterUint64INSTANCE.Lower(r), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uint64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterUint64INSTANCE.Lift(_uniffiRV), nil
	}
}

func ComplexAdd(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_add(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexArgument(c Complex) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_complex_argument(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus)
	}))
}

func ComplexConjugate(c Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_conjugate(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus),
		}
	}))
}

func ComplexDisplay(c Complex) string {
	return FfiConverterStringINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_display(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus),
		}
	}))
}

func ComplexDiv(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_div(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexFromPolar(r float64, theta float64) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_from_polar(FfiConverterFloat64INSTANCE.Lower(r), FfiConverterFloat64INSTANCE.Lower(theta), _uniffiStatus),
		}
	}))
}

func ComplexMagnitude(c Complex) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_complex_magnitude(FfiConverterComplexINSTANCE.Lower(c), _uniffiStatus)
	}))
}

func ComplexMul(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_mul(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func ComplexNew(re float64, im float64) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_new(FfiConverterFloat64INSTANCE.Lower(re), FfiConverterFloat64INSTANCE.Lower(im), _uniffiStatus),
		}
	}))
}

func ComplexSub(a Complex, b Complex) Complex {
	return FfiConverterComplexINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_complex_sub(FfiConverterComplexINSTANCE.Lower(a), FfiConverterComplexINSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func CrossProduct3d(a Point3, b Point3) Point3 {
	return FfiConverterPoint3INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_cross_product_3d(FfiConverterPoint3INSTANCE.Lower(a), FfiConverterPoint3INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func DegreesToRadians(deg float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_degrees_to_radians(FfiConverterFloat64INSTANCE.Lower(deg), _uniffiStatus)
	}))
}

func Determinant2x2(m Matrix2x2) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_determinant_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus)
	}))
}

func Determinant3x3(m Matrix3x3) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_determinant_3x3(FfiConverterMatrix3x3INSTANCE.Lower(m), _uniffiStatus)
	}))
}

func Distance2d(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_distance_2d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), _uniffiStatus)
	}))
}

func Distance3d(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_distance_3d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(z1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), FfiConverterFloat64INSTANCE.Lower(z2), _uniffiStatus)
	}))
}

func DotProduct(a []float64, b []float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_dot_product(FfiConverterSequenceFloat64INSTANCE.Lower(a), FfiConverterSequenceFloat64INSTANCE.Lower(b), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}

func Factorial(n uint64) (uint64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_factorial(FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uint64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterUint64INSTANCE.Lift(_uniffiRV), nil
	}
}

func Fibonacci(n uint64) (uint64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_fibonacci(FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uint64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterUint64INSTANCE.Lift(_uniffiRV), nil
	}
}

func Gcd(a uint64, b uint64) uint64 {
	return FfiConverterUint64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_gcd(FfiConverterUint64INSTANCE.Lower(a), FfiConverterUint64INSTANCE.Lower(b), _uniffiStatus)
	}))
}

func GeometricSum(a1 float64, r float64, n uint64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_geometric_sum(FfiConverterFloat64INSTANCE.Lower(a1), FfiConverterFloat64INSTANCE.Lower(r), FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	}))
}

func GeometricSumInfinite(a1 float64, r float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_geometric_sum_infinite(FfiConverterFloat64INSTANCE.Lower(a1), FfiConverterFloat64INSTANCE.Lower(r), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}

func GeometricTerm(a1 float64, r float64, n uint64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_geometric_term(FfiConverterFloat64INSTANCE.Lower(a1), FfiConverterFloat64INSTANCE.Lower(r), FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	}))
}

func IsPrime(n uint64) bool {
	return FfiConverterBoolINSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.int8_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_is_prime(FfiConverterUint64INSTANCE.Lower(n), _uniffiStatus)
	}))
}

func LawOfCosinesAngle(a float64, b float64, c float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_law_of_cosines_angle(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}

func LawOfCosinesSide(a float64, b float64, angleC float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_law_of_cosines_side(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(angleC), _uniffiStatus)
	}))
}

func LawOfSinesSide(a float64, angleA float64, angleB float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_law_of_sines_side(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(angleA), FfiConverterFloat64INSTANCE.Lower(angleB), _uniffiStatus)
	}))
}

func Lcm(a uint64, b uint64) (uint64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_lcm(FfiConverterUint64INSTANCE.Lower(a), FfiConverterUint64INSTANCE.Lower(b), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uint64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterUint64INSTANCE.Lift(_uniffiRV), nil
	}
}

func LinearRoot(m float64, b float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_linear_root(FfiConverterFloat64INSTANCE.Lower(m), FfiConverterFloat64INSTANCE.Lower(b), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}

func MatrixInverse2x2(m Matrix2x2) (Matrix2x2, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_matrix_inverse_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue Matrix2x2
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterMatrix2x2INSTANCE.Lift(_uniffiRV), nil
	}
}

func MatrixMultiply2x2(a Matrix2x2, b Matrix2x2) Matrix2x2 {
	return FfiConverterMatrix2x2INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_matrix_multiply_2x2(FfiConverterMatrix2x2INSTANCE.Lower(a), FfiConverterMatrix2x2INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func Midpoint2d(x1 float64, y1 float64, x2 float64, y2 float64) Point2 {
	return FfiConverterPoint2INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_midpoint_2d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), _uniffiStatus),
		}
	}))
}

func Permutations(n uint64, r uint64) (uint64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.uint64_t {
		return C.uniffi_uniffi_testing_ffi_fn_func_permutations(FfiConverterUint64INSTANCE.Lower(n), FfiConverterUint64INSTANCE.Lower(r), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uint64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterUint64INSTANCE.Lift(_uniffiRV), nil
	}
}

func PolynomialAdd(a []float64, b []float64) []float64 {
	return FfiConverterSequenceFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_polynomial_add(FfiConverterSequenceFloat64INSTANCE.Lower(a), FfiConverterSequenceFloat64INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func PolynomialDerivative(coeffs []float64) []float64 {
	return FfiConverterSequenceFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_polynomial_derivative(FfiConverterSequenceFloat64INSTANCE.Lower(coeffs), _uniffiStatus),
		}
	}))
}

func PolynomialEvaluate(coeffs []float64, x float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_polynomial_evaluate(FfiConverterSequenceFloat64INSTANCE.Lower(coeffs), FfiConverterFloat64INSTANCE.Lower(x), _uniffiStatus)
	}))
}

func PolynomialMultiply(a []float64, b []float64) []float64 {
	return FfiConverterSequenceFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_polynomial_multiply(FfiConverterSequenceFloat64INSTANCE.Lower(a), FfiConverterSequenceFloat64INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func QuadraticRoots(a float64, b float64, c float64) (F64Pair, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_quadratic_roots(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue F64Pair
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterF64PairINSTANCE.Lift(_uniffiRV), nil
	}
}

func QuadraticRootsComplex(a float64, b float64, c float64) (ComplexPair, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_quadratic_roots_complex(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue ComplexPair
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterComplexPairINSTANCE.Lift(_uniffiRV), nil
	}
}

func RadiansToDegrees(rad float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_radians_to_degrees(FfiConverterFloat64INSTANCE.Lower(rad), _uniffiStatus)
	}))
}

func Solve2x2(a Matrix2x2, b Point2) (Point2, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_solve_2x2(FfiConverterMatrix2x2INSTANCE.Lower(a), FfiConverterPoint2INSTANCE.Lower(b), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue Point2
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterPoint2INSTANCE.Lift(_uniffiRV), nil
	}
}

func Solve3x3(a Matrix3x3, b Point3) (Point3, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_solve_3x3(FfiConverterMatrix3x3INSTANCE.Lower(a), FfiConverterPoint3INSTANCE.Lower(b), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue Point3
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterPoint3INSTANCE.Lift(_uniffiRV), nil
	}
}

func Transpose2x2(m Matrix2x2) Matrix2x2 {
	return FfiConverterMatrix2x2INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_ffi_fn_func_transpose_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus),
		}
	}))
}

func TriangleAreaHeron(a float64, b float64, c float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_ffi_fn_func_triangle_area_heron(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}
