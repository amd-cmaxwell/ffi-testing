package polynomials

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_polynomials_ffi"
	types "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
)

type F64Pair = ffi.F64Pair

type MathError = types.MathError
type MathErrorDivisionByZero = types.MathErrorDivisionByZero
type MathErrorInvalidInput = types.MathErrorInvalidInput
type MathErrorNoRealSolution = types.MathErrorNoRealSolution
type MathErrorOverflow = types.MathErrorOverflow

var ErrMathErrorDivisionByZero = types.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = types.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = types.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = types.ErrMathErrorOverflow

func Add(a []float64, b []float64) []float64                                   { return ffi.PolynomialAdd(a, b) }
func Derivative(coeffs []float64) []float64                                     { return ffi.PolynomialDerivative(coeffs) }
func Evaluate(coeffs []float64, x float64) float64                              { return ffi.PolynomialEvaluate(coeffs, x) }
func LinearRoot(m float64, b float64) (float64, error)                          { return ffi.LinearRoot(m, b) }
func Multiply(a []float64, b []float64) []float64                               { return ffi.PolynomialMultiply(a, b) }
func QuadraticRoots(a float64, b float64, c float64) (F64Pair, error)           { return ffi.QuadraticRoots(a, b, c) }
