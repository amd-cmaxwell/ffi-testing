use crate::types::*;
use uniffi_testing_linear_algebra as la;

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
