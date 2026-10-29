package polynomials

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_ffi"
)

type F64Pair = ffi.F64Pair

type MathError = ffi.MathError
type MathErrorDivisionByZero = ffi.MathErrorDivisionByZero
type MathErrorInvalidInput = ffi.MathErrorInvalidInput
type MathErrorNoRealSolution = ffi.MathErrorNoRealSolution
type MathErrorOverflow = ffi.MathErrorOverflow

var ErrMathErrorDivisionByZero = ffi.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = ffi.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = ffi.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = ffi.ErrMathErrorOverflow

func Add(a []float64, b []float64) []float64 { return ffi.PolynomialAdd(a, b) }
func Derivative(coeffs []float64) []float64 { return ffi.PolynomialDerivative(coeffs) }
func Evaluate(coeffs []float64, x float64) float64 { return ffi.PolynomialEvaluate(coeffs, x) }
func LinearRoot(m float64, b float64) (float64, error) { return ffi.LinearRoot(m, b) }
func Multiply(a []float64, b []float64) []float64 { return ffi.PolynomialMultiply(a, b) }
func QuadraticRoots(a float64, b float64, c float64) (F64Pair, error) { return ffi.QuadraticRoots(a, b, c) }
