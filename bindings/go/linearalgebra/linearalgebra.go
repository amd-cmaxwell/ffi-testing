package linearalgebra

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_ffi"
)

type Matrix2x2 = ffi.Matrix2x2
type Matrix3x3 = ffi.Matrix3x3
type Point2 = ffi.Point2
type Point3 = ffi.Point3

type MathError = ffi.MathError
type MathErrorDivisionByZero = ffi.MathErrorDivisionByZero
type MathErrorInvalidInput = ffi.MathErrorInvalidInput
type MathErrorNoRealSolution = ffi.MathErrorNoRealSolution
type MathErrorOverflow = ffi.MathErrorOverflow

var ErrMathErrorDivisionByZero = ffi.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = ffi.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = ffi.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = ffi.ErrMathErrorOverflow

func Determinant2x2(m Matrix2x2) float64 { return ffi.Determinant2x2(m) }
func Determinant3x3(m Matrix3x3) float64 { return ffi.Determinant3x3(m) }
func Inverse2x2(m Matrix2x2) (Matrix2x2, error) { return ffi.MatrixInverse2x2(m) }
func Multiply2x2(a Matrix2x2, b Matrix2x2) Matrix2x2 { return ffi.MatrixMultiply2x2(a, b) }
func Product(a []float64, b []float64) (float64, error) { return ffi.DotProduct(a, b) }
func Product3d(a Point3, b Point3) Point3 { return ffi.CrossProduct3d(a, b) }
func Solve2x2(a Matrix2x2, b Point2) (Point2, error) { return ffi.Solve2x2(a, b) }
func Solve3x3(a Matrix3x3, b Point3) (Point3, error) { return ffi.Solve3x3(a, b) }
func Transpose2x2(m Matrix2x2) Matrix2x2 { return ffi.Transpose2x2(m) }
