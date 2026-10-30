package trig

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_trig_ffi"
	types "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
)

type Point2 = types.Point2

type MathError = types.MathError
type MathErrorDivisionByZero = types.MathErrorDivisionByZero
type MathErrorInvalidInput = types.MathErrorInvalidInput
type MathErrorNoRealSolution = types.MathErrorNoRealSolution
type MathErrorOverflow = types.MathErrorOverflow

var ErrMathErrorDivisionByZero = types.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = types.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = types.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = types.ErrMathErrorOverflow

func Area(radius float64) float64                                                                      { return ffi.CircleArea(radius) }
func AreaHeron(a float64, b float64, c float64) (float64, error)                                       { return ffi.TriangleAreaHeron(a, b, c) }
func Circumference(radius float64) float64                                                             { return ffi.CircleCircumference(radius) }
func CosinesAngle(a float64, b float64, c float64) (float64, error)                                    { return ffi.LawOfCosinesAngle(a, b, c) }
func CosinesSide(a float64, b float64, angleC float64) float64                                         { return ffi.LawOfCosinesSide(a, b, angleC) }
func Distance2d(x1 float64, y1 float64, x2 float64, y2 float64) float64                                { return ffi.Distance2d(x1, y1, x2, y2) }
func Distance3d(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) float64        { return ffi.Distance3d(x1, y1, z1, x2, y2, z2) }
func Midpoint2d(x1 float64, y1 float64, x2 float64, y2 float64) Point2                                 { return ffi.Midpoint2d(x1, y1, x2, y2) }
func SinesSide(a float64, angleA float64, angleB float64) float64                                       { return ffi.LawOfSinesSide(a, angleA, angleB) }
func ToDegrees(rad float64) float64                                                                     { return ffi.RadiansToDegrees(rad) }
func ToRadians(deg float64) float64                                                                     { return ffi.DegreesToRadians(deg) }
