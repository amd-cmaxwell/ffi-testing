package uniffi_testing_types_ffi

// #cgo LDFLAGS: -L${SRCDIR} -luniffi_testing_types_ffi -lm -ldl -lpthread
// #include <uniffi_testing_types_ffi.h>
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
		C.ffi_uniffi_testing_types_ffi_rustbuffer_free(cb.inner, status)
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
		return C.ffi_uniffi_testing_types_ffi_rustbuffer_from_bytes(foreign, status)
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
		return C.ffi_uniffi_testing_types_ffi_uniffi_contract_version()
	})
	if bindingsContractVersion != int(scaffoldingContractVersion) {
		// If this happens try cleaning and rebuilding your project
		panic("uniffi_testing_types_ffi: UniFFI contract version mismatch")
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
