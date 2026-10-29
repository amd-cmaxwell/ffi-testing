use crate::types::*;
use uniffi_testing_complex as complex_crate;
use uniffi_testing_core as core;

impl From<core::MathError> for MathError {
    fn from(e: core::MathError) -> Self {
        match e {
            core::MathError::DivisionByZero => MathError::DivisionByZero,
            core::MathError::InvalidInput(msg) => MathError::InvalidInput { message: msg },
            core::MathError::NoRealSolution => MathError::NoRealSolution,
            core::MathError::Overflow => MathError::Overflow,
        }
    }
}

impl From<complex_crate::Complex> for Complex {
    fn from(c: complex_crate::Complex) -> Self {
        Complex { re: c.re, im: c.im }
    }
}

impl From<&Complex> for complex_crate::Complex {
    fn from(c: &Complex) -> Self {
        complex_crate::Complex::new(c.re, c.im)
    }
}

impl From<[[f64; 2]; 2]> for Matrix2x2 {
    fn from(m: [[f64; 2]; 2]) -> Self {
        Matrix2x2 {
            m00: m[0][0],
            m01: m[0][1],
            m10: m[1][0],
            m11: m[1][1],
        }
    }
}

impl From<&Matrix2x2> for [[f64; 2]; 2] {
    fn from(m: &Matrix2x2) -> Self {
        [[m.m00, m.m01], [m.m10, m.m11]]
    }
}

impl From<[[f64; 3]; 3]> for Matrix3x3 {
    fn from(m: [[f64; 3]; 3]) -> Self {
        Matrix3x3 {
            m00: m[0][0],
            m01: m[0][1],
            m02: m[0][2],
            m10: m[1][0],
            m11: m[1][1],
            m12: m[1][2],
            m20: m[2][0],
            m21: m[2][1],
            m22: m[2][2],
        }
    }
}

impl From<&Matrix3x3> for [[f64; 3]; 3] {
    fn from(m: &Matrix3x3) -> Self {
        [
            [m.m00, m.m01, m.m02],
            [m.m10, m.m11, m.m12],
            [m.m20, m.m21, m.m22],
        ]
    }
}

impl From<[f64; 2]> for Point2 {
    fn from(v: [f64; 2]) -> Self {
        Point2 { x: v[0], y: v[1] }
    }
}

impl From<&Point2> for [f64; 2] {
    fn from(p: &Point2) -> Self {
        [p.x, p.y]
    }
}

impl From<[f64; 3]> for Point3 {
    fn from(v: [f64; 3]) -> Self {
        Point3 {
            x: v[0],
            y: v[1],
            z: v[2],
        }
    }
}

impl From<&Point3> for [f64; 3] {
    fn from(p: &Point3) -> Self {
        [p.x, p.y, p.z]
    }
}
