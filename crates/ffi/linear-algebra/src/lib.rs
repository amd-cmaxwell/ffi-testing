uniffi::setup_scaffolding!();

use uniffi_testing_linear_algebra as la;
use uniffi_testing_types_ffi::{MathError, Point2};

#[derive(Debug, Clone, uniffi::Record)]
pub struct Point3 {
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

#[derive(Debug, Clone, uniffi::Record)]
pub struct Matrix2x2 {
    pub m00: f64,
    pub m01: f64,
    pub m10: f64,
    pub m11: f64,
}

#[derive(Debug, Clone, uniffi::Record)]
pub struct Matrix3x3 {
    pub m00: f64,
    pub m01: f64,
    pub m02: f64,
    pub m10: f64,
    pub m11: f64,
    pub m12: f64,
    pub m20: f64,
    pub m21: f64,
    pub m22: f64,
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

#[uniffi::export]
pub fn determinant_2x2(m: Matrix2x2) -> f64 {
    la::determinant_2x2((&m).into())
}

#[uniffi::export]
pub fn determinant_3x3(m: Matrix3x3) -> f64 {
    la::determinant_3x3((&m).into())
}

#[uniffi::export]
pub fn solve_2x2(a: Matrix2x2, b: Point2) -> Result<Point2, MathError> {
    let result = la::solve_2x2((&a).into(), (&b).into()).map_err(MathError::from)?;
    Ok(result.into())
}

#[uniffi::export]
pub fn solve_3x3(a: Matrix3x3, b: Point3) -> Result<Point3, MathError> {
    let result = la::solve_3x3((&a).into(), (&b).into()).map_err(MathError::from)?;
    Ok(result.into())
}

#[uniffi::export]
pub fn matrix_multiply_2x2(a: Matrix2x2, b: Matrix2x2) -> Matrix2x2 {
    la::matrix_multiply_2x2((&a).into(), (&b).into()).into()
}

#[uniffi::export]
pub fn matrix_inverse_2x2(m: Matrix2x2) -> Result<Matrix2x2, MathError> {
    let result = la::matrix_inverse_2x2((&m).into()).map_err(MathError::from)?;
    Ok(result.into())
}

#[uniffi::export]
pub fn transpose_2x2(m: Matrix2x2) -> Matrix2x2 {
    la::transpose_2x2((&m).into()).into()
}

#[uniffi::export]
pub fn dot_product(a: Vec<f64>, b: Vec<f64>) -> Result<f64, MathError> {
    la::dot_product(&a, &b).map_err(MathError::from)
}

#[uniffi::export]
pub fn cross_product_3d(a: Point3, b: Point3) -> Point3 {
    la::cross_product_3d((&a).into(), (&b).into()).into()
}
