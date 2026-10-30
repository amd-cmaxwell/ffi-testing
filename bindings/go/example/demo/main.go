package main

import (
	"fmt"

	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_complex_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_core_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_linear_algebra_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_polynomials_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_sequences_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_trig_ffi"
	"github.com/amd-cmaxwell/ffi-testing/bindings/go/uniffi_testing_types_ffi"
)

func main() {
	fmt.Println("=== Core ===")
	fmt.Printf("Add(2, 3) = %d\n", uniffi_testing_core_ffi.Add(2, 3))

	fmt.Println("\n=== Complex Numbers ===")
	z := uniffi_testing_complex_ffi.ComplexNew(3.0, 4.0)
	fmt.Printf("z = %s\n", uniffi_testing_complex_ffi.ComplexDisplay(z))
	fmt.Printf("  |z| = %f\n", uniffi_testing_complex_ffi.ComplexMagnitude(z))
	fmt.Printf("  arg(z) = %.4f\n", uniffi_testing_complex_ffi.ComplexArgument(z))
	conj := uniffi_testing_complex_ffi.ComplexConjugate(z)
	fmt.Printf("  conj(z) = %s\n", uniffi_testing_complex_ffi.ComplexDisplay(conj))

	w := uniffi_testing_complex_ffi.ComplexNew(1.0, 2.0)
	fmt.Printf("w = %s\n", uniffi_testing_complex_ffi.ComplexDisplay(w))
	fmt.Printf("  z + w = %s\n", uniffi_testing_complex_ffi.ComplexDisplay(uniffi_testing_complex_ffi.ComplexAdd(z, w)))
	fmt.Printf("  z * w = %s\n", uniffi_testing_complex_ffi.ComplexDisplay(uniffi_testing_complex_ffi.ComplexMul(z, w)))

	pair, err := uniffi_testing_complex_ffi.QuadraticRootsComplex(1.0, 0.0, 1.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x^2 + 1 = 0: %s, %s\n", uniffi_testing_complex_ffi.ComplexDisplay(pair.First), uniffi_testing_complex_ffi.ComplexDisplay(pair.Second))

	fmt.Println("\n=== Polynomials ===")
	roots, err := uniffi_testing_polynomials_ffi.QuadraticRoots(1.0, -5.0, 6.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x^2 - 5x + 6 = 0: x = %f, %f\n", roots.First, roots.Second)
	fmt.Printf("eval x^2 - 3x + 2 at x=1: %f\n", uniffi_testing_polynomials_ffi.PolynomialEvaluate([]float64{1.0, -3.0, 2.0}, 1.0))

	lr, err := uniffi_testing_polynomials_ffi.LinearRoot(2.0, 4.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("linear_root(2x + 4 = 0): x = %f\n", lr)
	fmt.Printf("(x+1)(x-1) = %v\n", uniffi_testing_polynomials_ffi.PolynomialMultiply([]float64{1.0, 1.0}, []float64{1.0, -1.0}))
	fmt.Printf("d/dx(3x^2 + 2x + 1) = %v\n", uniffi_testing_polynomials_ffi.PolynomialDerivative([]float64{3.0, 2.0, 1.0}))

	fmt.Println("\n=== Linear Algebra ===")
	m := uniffi_testing_linear_algebra_ffi.Matrix2x2{M00: 2, M01: 3, M10: 1, M11: 4}
	fmt.Printf("det([[2,3],[1,4]]) = %f\n", uniffi_testing_linear_algebra_ffi.Determinant2x2(m))

	sol, err := uniffi_testing_linear_algebra_ffi.Solve2x2(
		uniffi_testing_linear_algebra_ffi.Matrix2x2{M00: 1, M01: 1, M10: 2, M11: -1},
		uniffi_testing_types_ffi.Point2{X: 3, Y: 0},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("solve x+y=3, 2x-y=0: x=%f, y=%f\n", sol.X, sol.Y)

	dp, err := uniffi_testing_linear_algebra_ffi.DotProduct([]float64{1, 2, 3}, []float64{4, 5, 6})
	if err != nil {
		panic(err)
	}
	fmt.Printf("dot([1,2,3], [4,5,6]) = %f\n", dp)

	cp := uniffi_testing_linear_algebra_ffi.CrossProduct3d(
		uniffi_testing_linear_algebra_ffi.Point3{X: 1, Y: 0, Z: 0},
		uniffi_testing_linear_algebra_ffi.Point3{X: 0, Y: 1, Z: 0},
	)
	fmt.Printf("cross(i, j) = (%f, %f, %f)\n", cp.X, cp.Y, cp.Z)

	fmt.Println("\n=== Trigonometry ===")
	fmt.Printf("90 deg = %.4f rad\n", uniffi_testing_trig_ffi.DegreesToRadians(90.0))
	fmt.Printf("pi rad = %.1f deg\n", uniffi_testing_trig_ffi.RadiansToDegrees(3.14159265358979))
	fmt.Printf("distance (0,0)-(3,4) = %f\n", uniffi_testing_trig_ffi.Distance2d(0, 0, 3, 4))
	mid := uniffi_testing_trig_ffi.Midpoint2d(0, 0, 4, 6)
	fmt.Printf("midpoint (0,0)-(4,6) = (%f, %f)\n", mid.X, mid.Y)
	area, err := uniffi_testing_trig_ffi.TriangleAreaHeron(3, 4, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Heron area 3-4-5 = %f\n", area)
	fmt.Printf("circle area r=3 = %.4f\n", uniffi_testing_trig_ffi.CircleArea(3.0))

	fmt.Println("\n=== Sequences & Number Theory ===")
	f, err := uniffi_testing_sequences_ffi.Factorial(10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("10! = %d\n", f)

	c, err := uniffi_testing_sequences_ffi.Combinations(10, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("C(10,3) = %d\n", c)

	p, err := uniffi_testing_sequences_ffi.Permutations(5, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("P(5,2) = %d\n", p)

	fib, err := uniffi_testing_sequences_ffi.Fibonacci(20)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fib(20) = %d\n", fib)

	fmt.Printf("gcd(48, 18) = %d\n", uniffi_testing_sequences_ffi.Gcd(48, 18))

	l, err := uniffi_testing_sequences_ffi.Lcm(4, 6)
	if err != nil {
		panic(err)
	}
	fmt.Printf("lcm(4, 6) = %d\n", l)

	fmt.Printf("is_prime(97) = %t\n", uniffi_testing_sequences_ffi.IsPrime(97))
	fmt.Printf("arithmetic sum 1..100 = %f\n", uniffi_testing_sequences_ffi.ArithmeticSum(1.0, 1.0, 100))

	gs, err := uniffi_testing_sequences_ffi.GeometricSumInfinite(1.0, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("geometric sum infinite a=1, r=0.5 = %f\n", gs)

	fmt.Println("\n=== Error Handling ===")
	_, err = uniffi_testing_polynomials_ffi.QuadraticRoots(0.0, 1.0, 2.0)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	_, err = uniffi_testing_polynomials_ffi.QuadraticRoots(1.0, 0.0, 1.0)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	_, err = uniffi_testing_sequences_ffi.Factorial(21)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	fmt.Println("\nAll demos passed!")
}
