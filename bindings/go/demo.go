package main

import (
	"fmt"

	ffi "uniffi-testing-go/uniffi_testing_ffi"
)

func main() {
	fmt.Println("=== Core ===")
	fmt.Printf("Add(2, 3) = %d\n", ffi.Add(2, 3))

	fmt.Println("\n=== Complex Numbers ===")
	z := ffi.ComplexNew(3.0, 4.0)
	fmt.Printf("z = %s\n", ffi.ComplexDisplay(z))
	fmt.Printf("  |z| = %f\n", ffi.ComplexMagnitude(z))
	fmt.Printf("  arg(z) = %.4f\n", ffi.ComplexArgument(z))
	conj := ffi.ComplexConjugate(z)
	fmt.Printf("  conj(z) = %s\n", ffi.ComplexDisplay(conj))

	w := ffi.ComplexNew(1.0, 2.0)
	fmt.Printf("w = %s\n", ffi.ComplexDisplay(w))
	fmt.Printf("  z + w = %s\n", ffi.ComplexDisplay(ffi.ComplexAdd(z, w)))
	fmt.Printf("  z * w = %s\n", ffi.ComplexDisplay(ffi.ComplexMul(z, w)))

	pair, err := ffi.QuadraticRootsComplex(1.0, 0.0, 1.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x^2 + 1 = 0: %s, %s\n", ffi.ComplexDisplay(pair.First), ffi.ComplexDisplay(pair.Second))

	fmt.Println("\n=== Polynomials ===")
	roots, err := ffi.QuadraticRoots(1.0, -5.0, 6.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x^2 - 5x + 6 = 0: x = %f, %f\n", roots.First, roots.Second)
	fmt.Printf("eval x^2 - 3x + 2 at x=1: %f\n", ffi.PolynomialEvaluate([]float64{1.0, -3.0, 2.0}, 1.0))

	lr, err := ffi.LinearRoot(2.0, 4.0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("linear_root(2x + 4 = 0): x = %f\n", lr)
	fmt.Printf("(x+1)(x-1) = %v\n", ffi.PolynomialMultiply([]float64{1.0, 1.0}, []float64{1.0, -1.0}))
	fmt.Printf("d/dx(3x^2 + 2x + 1) = %v\n", ffi.PolynomialDerivative([]float64{3.0, 2.0, 1.0}))

	fmt.Println("\n=== Linear Algebra ===")
	m := ffi.Matrix2x2{M00: 2, M01: 3, M10: 1, M11: 4}
	fmt.Printf("det([[2,3],[1,4]]) = %f\n", ffi.Determinant2x2(m))

	sol, err := ffi.Solve2x2(
		ffi.Matrix2x2{M00: 1, M01: 1, M10: 2, M11: -1},
		ffi.Point2{X: 3, Y: 0},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("solve x+y=3, 2x-y=0: x=%f, y=%f\n", sol.X, sol.Y)

	dp, err := ffi.DotProduct([]float64{1, 2, 3}, []float64{4, 5, 6})
	if err != nil {
		panic(err)
	}
	fmt.Printf("dot([1,2,3], [4,5,6]) = %f\n", dp)

	cp := ffi.CrossProduct3d(
		ffi.Point3{X: 1, Y: 0, Z: 0},
		ffi.Point3{X: 0, Y: 1, Z: 0},
	)
	fmt.Printf("cross(i, j) = (%f, %f, %f)\n", cp.X, cp.Y, cp.Z)

	fmt.Println("\n=== Trigonometry ===")
	fmt.Printf("90 deg = %.4f rad\n", ffi.DegreesToRadians(90.0))
	fmt.Printf("pi rad = %.1f deg\n", ffi.RadiansToDegrees(3.14159265358979))
	fmt.Printf("distance (0,0)-(3,4) = %f\n", ffi.Distance2d(0, 0, 3, 4))
	mid := ffi.Midpoint2d(0, 0, 4, 6)
	fmt.Printf("midpoint (0,0)-(4,6) = (%f, %f)\n", mid.X, mid.Y)
	area, err := ffi.TriangleAreaHeron(3, 4, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Heron area 3-4-5 = %f\n", area)
	fmt.Printf("circle area r=3 = %.4f\n", ffi.CircleArea(3.0))

	fmt.Println("\n=== Sequences & Number Theory ===")
	f, err := ffi.Factorial(10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("10! = %d\n", f)

	c, err := ffi.Combinations(10, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("C(10,3) = %d\n", c)

	p, err := ffi.Permutations(5, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("P(5,2) = %d\n", p)

	fib, err := ffi.Fibonacci(20)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fib(20) = %d\n", fib)

	fmt.Printf("gcd(48, 18) = %d\n", ffi.Gcd(48, 18))

	l, err := ffi.Lcm(4, 6)
	if err != nil {
		panic(err)
	}
	fmt.Printf("lcm(4, 6) = %d\n", l)

	fmt.Printf("is_prime(97) = %t\n", ffi.IsPrime(97))
	fmt.Printf("arithmetic sum 1..100 = %f\n", ffi.ArithmeticSum(1.0, 1.0, 100))

	gs, err := ffi.GeometricSumInfinite(1.0, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("geometric sum infinite a=1, r=0.5 = %f\n", gs)

	fmt.Println("\n=== Error Handling ===")
	_, err = ffi.QuadraticRoots(0.0, 1.0, 2.0)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	_, err = ffi.QuadraticRoots(1.0, 0.0, 1.0)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	_, err = ffi.Factorial(21)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	fmt.Println("\nAll demos passed!")
}
