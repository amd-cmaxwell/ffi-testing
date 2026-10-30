package sequences

import (
	ffi "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_sequences_ffi"
	types "github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
)

type MathError = types.MathError
type MathErrorDivisionByZero = types.MathErrorDivisionByZero
type MathErrorInvalidInput = types.MathErrorInvalidInput
type MathErrorNoRealSolution = types.MathErrorNoRealSolution
type MathErrorOverflow = types.MathErrorOverflow

var ErrMathErrorDivisionByZero = types.ErrMathErrorDivisionByZero
var ErrMathErrorInvalidInput = types.ErrMathErrorInvalidInput
var ErrMathErrorNoRealSolution = types.ErrMathErrorNoRealSolution
var ErrMathErrorOverflow = types.ErrMathErrorOverflow

func ArithmeticSum(a1 float64, d float64, n uint64) float64   { return ffi.ArithmeticSum(a1, d, n) }
func ArithmeticTerm(a1 float64, d float64, n uint64) float64  { return ffi.ArithmeticTerm(a1, d, n) }
func Combinations(n uint64, r uint64) (uint64, error)          { return ffi.Combinations(n, r) }
func Factorial(n uint64) (uint64, error)                       { return ffi.Factorial(n) }
func Fibonacci(n uint64) (uint64, error)                       { return ffi.Fibonacci(n) }
func Gcd(a uint64, b uint64) uint64                            { return ffi.Gcd(a, b) }
func GeometricSum(a1 float64, r float64, n uint64) float64    { return ffi.GeometricSum(a1, r, n) }
func GeometricTerm(a1 float64, r float64, n uint64) float64   { return ffi.GeometricTerm(a1, r, n) }
func IsPrime(n uint64) bool                                    { return ffi.IsPrime(n) }
func Lcm(a uint64, b uint64) (uint64, error)                   { return ffi.Lcm(a, b) }
func Permutations(n uint64, r uint64) (uint64, error)          { return ffi.Permutations(n, r) }
func SumInfinite(a1 float64, r float64) (float64, error)       { return ffi.GeometricSumInfinite(a1, r) }
