package complexnums

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_complex_ffi"
	types "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
)

type Complex = ffi.Complex
type ComplexPair = ffi.ComplexPair

type MathError = types.MathError
type MathErrorDivisionByZero = types.MathErrorDivisionByZero
type MathErrorInvalidInput = types.MathErrorInvalidInput
type MathErrorNoRealSolution = types.MathErrorNoRealSolution
type MathErrorOverflow = types.MathErrorOverflow

var ErrMathErrorDivisionByZero = types.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = types.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = types.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = types.ErrMathErrorOverflow

func Add(a Complex, b Complex) Complex                                     { return ffi.ComplexAdd(a, b) }
func Argument(c Complex) float64                                           { return ffi.ComplexArgument(c) }
func Conjugate(c Complex) Complex                                          { return ffi.ComplexConjugate(c) }
func Display(c Complex) string                                             { return ffi.ComplexDisplay(c) }
func Div(a Complex, b Complex) Complex                                     { return ffi.ComplexDiv(a, b) }
func FromPolar(r float64, theta float64) Complex                           { return ffi.ComplexFromPolar(r, theta) }
func Magnitude(c Complex) float64                                          { return ffi.ComplexMagnitude(c) }
func Mul(a Complex, b Complex) Complex                                     { return ffi.ComplexMul(a, b) }
func New(re float64, im float64) Complex                                   { return ffi.ComplexNew(re, im) }
func QuadraticRootsComplex(a float64, b float64, c float64) (ComplexPair, error) { return ffi.QuadraticRootsComplex(a, b, c) }
func Sub(a Complex, b Complex) Complex                                     { return ffi.ComplexSub(a, b) }
