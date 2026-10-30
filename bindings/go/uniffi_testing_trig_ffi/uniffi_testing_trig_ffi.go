package uniffi_testing_trig_ffi

// #cgo LDFLAGS: -L${SRCDIR} -luniffi_testing_trig_ffi -lm -ldl -lpthread
// #include <uniffi_testing_trig_ffi.h>
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
		C.ffi_uniffi_testing_trig_ffi_rustbuffer_free(cb.inner, status)
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
		return C.ffi_uniffi_testing_trig_ffi_rustbuffer_from_bytes(foreign, status)
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
		return C.ffi_uniffi_testing_trig_ffi_uniffi_contract_version()
	})
	if bindingsContractVersion != int(scaffoldingContractVersion) {
		// If this happens try cleaning and rebuilding your project
		panic("uniffi_testing_trig_ffi: UniFFI contract version mismatch")
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_circle_area()
		})
		if checksum != 9437 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_circle_area: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_circle_circumference()
		})
		if checksum != 36709 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_circle_circumference: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_degrees_to_radians()
		})
		if checksum != 24173 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_degrees_to_radians: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_distance_2d()
		})
		if checksum != 47409 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_distance_2d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_distance_3d()
		})
		if checksum != 43261 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_distance_3d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_angle()
		})
		if checksum != 63631 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_angle: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_side()
		})
		if checksum != 44722 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_side: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_sines_side()
		})
		if checksum != 9177 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_sines_side: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_midpoint_2d()
		})
		if checksum != 63681 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_midpoint_2d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_radians_to_degrees()
		})
		if checksum != 53073 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_radians_to_degrees: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_trig_ffi_checksum_func_triangle_area_heron()
		})
		if checksum != 19364 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_trig_ffi: uniffi_uniffi_testing_trig_ffi_checksum_func_triangle_area_heron: UniFFI API checksum mismatch")
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

func CircleArea(radius float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_circle_area(FfiConverterFloat64INSTANCE.Lower(radius), _uniffiStatus)
	}))
}

func CircleCircumference(radius float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_circle_circumference(FfiConverterFloat64INSTANCE.Lower(radius), _uniffiStatus)
	}))
}

func DegreesToRadians(deg float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_degrees_to_radians(FfiConverterFloat64INSTANCE.Lower(deg), _uniffiStatus)
	}))
}

func Distance2d(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_distance_2d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), _uniffiStatus)
	}))
}

func Distance3d(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_distance_3d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(z1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), FfiConverterFloat64INSTANCE.Lower(z2), _uniffiStatus)
	}))
}

func LawOfCosinesAngle(a float64, b float64, c float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_angle(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus)
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
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_side(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(angleC), _uniffiStatus)
	}))
}

func LawOfSinesSide(a float64, angleA float64, angleB float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_law_of_sines_side(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(angleA), FfiConverterFloat64INSTANCE.Lower(angleB), _uniffiStatus)
	}))
}

func Midpoint2d(x1 float64, y1 float64, x2 float64, y2 float64) uniffi_testing_types_ffi.Point2 {
	return uniffi_testing_types_ffi.FfiConverterPoint2INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_trig_ffi_fn_func_midpoint_2d(FfiConverterFloat64INSTANCE.Lower(x1), FfiConverterFloat64INSTANCE.Lower(y1), FfiConverterFloat64INSTANCE.Lower(x2), FfiConverterFloat64INSTANCE.Lower(y2), _uniffiStatus),
		}
	}))
}

func RadiansToDegrees(rad float64) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_radians_to_degrees(FfiConverterFloat64INSTANCE.Lower(rad), _uniffiStatus)
	}))
}

func TriangleAreaHeron(a float64, b float64, c float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_trig_ffi_fn_func_triangle_area_heron(FfiConverterFloat64INSTANCE.Lower(a), FfiConverterFloat64INSTANCE.Lower(b), FfiConverterFloat64INSTANCE.Lower(c), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}
