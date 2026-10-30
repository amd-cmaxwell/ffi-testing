uniffi::setup_scaffolding!();

use uniffi_testing_complex as cx;
use uniffi_testing_types_ffi::MathError;

#[derive(Debug, Clone, uniffi::Record)]
pub struct Complex {
    pub re: f64,
    pub im: f64,
}

#[derive(Debug, Clone, uniffi::Record)]
pub struct ComplexPair {
    pub first: Complex,
    pub second: Complex,
}

impl From<cx::Complex> for Complex {
    fn from(c: cx::Complex) -> Self {
        Complex { re: c.re, im: c.im }
    }
}

impl From<&Complex> for cx::Complex {
    fn from(c: &Complex) -> Self {
        cx::Complex::new(c.re, c.im)
    }
}

#[uniffi::export]
pub fn complex_new(re: f64, im: f64) -> Complex {
    Complex { re, im }
}

#[uniffi::export]
pub fn complex_from_polar(r: f64, theta: f64) -> Complex {
    cx::Complex::from_polar(r, theta).into()
}

#[uniffi::export]
pub fn complex_magnitude(c: Complex) -> f64 {
    cx::Complex::from(&c).magnitude()
}

#[uniffi::export]
pub fn complex_argument(c: Complex) -> f64 {
    cx::Complex::from(&c).argument()
}

#[uniffi::export]
pub fn complex_conjugate(c: Complex) -> Complex {
    cx::Complex::from(&c).conjugate().into()
}

#[uniffi::export]
pub fn complex_add(a: Complex, b: Complex) -> Complex {
    (cx::Complex::from(&a) + cx::Complex::from(&b)).into()
}

#[uniffi::export]
pub fn complex_sub(a: Complex, b: Complex) -> Complex {
    (cx::Complex::from(&a) - cx::Complex::from(&b)).into()
}

#[uniffi::export]
pub fn complex_mul(a: Complex, b: Complex) -> Complex {
    (cx::Complex::from(&a) * cx::Complex::from(&b)).into()
}

#[uniffi::export]
pub fn complex_div(a: Complex, b: Complex) -> Complex {
    (cx::Complex::from(&a) / cx::Complex::from(&b)).into()
}

#[uniffi::export]
pub fn complex_display(c: Complex) -> String {
    cx::Complex::from(&c).to_string()
}

#[uniffi::export]
pub fn quadratic_roots_complex(a: f64, b: f64, c: f64) -> Result<ComplexPair, MathError> {
    let (r1, r2) = cx::quadratic_roots_complex(a, b, c).map_err(MathError::from)?;
    Ok(ComplexPair {
        first: r1.into(),
        second: r2.into(),
    })
}
