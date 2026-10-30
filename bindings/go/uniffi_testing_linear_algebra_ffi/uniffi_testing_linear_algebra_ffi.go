package uniffi_testing_linear_algebra_ffi

// #cgo LDFLAGS: -L${SRCDIR} -luniffi_testing_linear_algebra_ffi -lm -ldl -lpthread
// #include <uniffi_testing_linear_algebra_ffi.h>
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
		C.ffi_uniffi_testing_linear_algebra_ffi_rustbuffer_free(cb.inner, status)
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
		return C.ffi_uniffi_testing_linear_algebra_ffi_rustbuffer_from_bytes(foreign, status)
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
		return C.ffi_uniffi_testing_linear_algebra_ffi_uniffi_contract_version()
	})
	if bindingsContractVersion != int(scaffoldingContractVersion) {
		// If this happens try cleaning and rebuilding your project
		panic("uniffi_testing_linear_algebra_ffi: UniFFI contract version mismatch")
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_cross_product_3d()
		})
		if checksum != 51914 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_cross_product_3d: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_determinant_2x2()
		})
		if checksum != 40382 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_determinant_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_determinant_3x3()
		})
		if checksum != 48794 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_determinant_3x3: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_dot_product()
		})
		if checksum != 8730 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_dot_product: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_matrix_inverse_2x2()
		})
		if checksum != 48481 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_matrix_inverse_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_matrix_multiply_2x2()
		})
		if checksum != 26451 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_matrix_multiply_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_solve_2x2()
		})
		if checksum != 55904 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_solve_2x2: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_solve_3x3()
		})
		if checksum != 976 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_solve_3x3: UniFFI API checksum mismatch")
		}
	}
	{
		checksum := rustCall(func(_uniffiStatus *C.RustCallStatus) C.uint16_t {
			return C.uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_transpose_2x2()
		})
		if checksum != 28155 {
			// If this happens try cleaning and rebuilding your project
			panic("uniffi_testing_linear_algebra_ffi: uniffi_uniffi_testing_linear_algebra_ffi_checksum_func_transpose_2x2: UniFFI API checksum mismatch")
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

func CrossProduct3d(a Point3, b Point3) Point3 {
	return FfiConverterPoint3INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_cross_product_3d(FfiConverterPoint3INSTANCE.Lower(a), FfiConverterPoint3INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func Determinant2x2(m Matrix2x2) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_determinant_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus)
	}))
}

func Determinant3x3(m Matrix3x3) float64 {
	return FfiConverterFloat64INSTANCE.Lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_determinant_3x3(FfiConverterMatrix3x3INSTANCE.Lower(m), _uniffiStatus)
	}))
}

func DotProduct(a []float64, b []float64) (float64, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) C.double {
		return C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_dot_product(FfiConverterSequenceFloat64INSTANCE.Lower(a), FfiConverterSequenceFloat64INSTANCE.Lower(b), _uniffiStatus)
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue float64
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return FfiConverterFloat64INSTANCE.Lift(_uniffiRV), nil
	}
}

func MatrixInverse2x2(m Matrix2x2) (Matrix2x2, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_matrix_inverse_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus),
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
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_matrix_multiply_2x2(FfiConverterMatrix2x2INSTANCE.Lower(a), FfiConverterMatrix2x2INSTANCE.Lower(b), _uniffiStatus),
		}
	}))
}

func Solve2x2(a Matrix2x2, b uniffi_testing_types_ffi.Point2) (uniffi_testing_types_ffi.Point2, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_solve_2x2(FfiConverterMatrix2x2INSTANCE.Lower(a),
				CFromRustBuffer(uniffi_testing_types_ffi.FfiConverterPoint2INSTANCE.LowerExternal(b)), _uniffiStatus),
		}
	})
	if _uniffiErr != nil {
		var _uniffiDefaultValue uniffi_testing_types_ffi.Point2
		return _uniffiDefaultValue, _uniffiErr
	} else {
		return uniffi_testing_types_ffi.FfiConverterPoint2INSTANCE.Lift(_uniffiRV), nil
	}
}

func Solve3x3(a Matrix3x3, b Point3) (Point3, error) {
	_uniffiRV, _uniffiErr := rustCallWithError[MathError](uniffi_testing_types_ffi.FfiConverterMathError{}, func(_uniffiStatus *C.RustCallStatus) RustBufferI {
		return GoRustBuffer{
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_solve_3x3(FfiConverterMatrix3x3INSTANCE.Lower(a), FfiConverterPoint3INSTANCE.Lower(b), _uniffiStatus),
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
			inner: C.uniffi_uniffi_testing_linear_algebra_ffi_fn_func_transpose_2x2(FfiConverterMatrix2x2INSTANCE.Lower(m), _uniffiStatus),
		}
	}))
}
