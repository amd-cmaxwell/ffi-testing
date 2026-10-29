use crate::types::*;
use uniffi_testing_polynomials as poly;

#[uniffi::export]
pub fn polynomial_evaluate(coeffs: Vec<f64>, x: f64) -> f64 {
    poly::evaluate(&coeffs, x)
}

#[uniffi::export]
pub fn quadratic_roots(a: f64, b: f64, c: f64) -> Result<F64Pair, MathError> {
    let (r1, r2) = poly::quadratic_roots(a, b, c).map_err(MathError::from)?;
    Ok(F64Pair {
        first: r1,
        second: r2,
    })
}

#[uniffi::export]
pub fn linear_root(m: f64, b: f64) -> Result<f64, MathError> {
    poly::linear_root(m, b).map_err(MathError::from)
}

#[uniffi::export]
pub fn polynomial_add(a: Vec<f64>, b: Vec<f64>) -> Vec<f64> {
    poly::add(&a, &b)
}

#[uniffi::export]
pub fn polynomial_multiply(a: Vec<f64>, b: Vec<f64>) -> Vec<f64> {
    poly::multiply(&a, &b)
}

#[uniffi::export]
pub fn polynomial_derivative(coeffs: Vec<f64>) -> Vec<f64> {
    poly::derivative(&coeffs)
}
