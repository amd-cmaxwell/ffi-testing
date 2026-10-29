use std::fmt;

#[derive(Debug, Clone, PartialEq, uniffi::Error)]
pub enum MathError {
    DivisionByZero,
    InvalidInput { message: String },
    NoRealSolution,
    Overflow,
}

impl fmt::Display for MathError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            MathError::DivisionByZero => write!(f, "division by zero"),
            MathError::InvalidInput { message } => write!(f, "invalid input: {message}"),
            MathError::NoRealSolution => write!(f, "no real solution"),
            MathError::Overflow => write!(f, "arithmetic overflow"),
        }
    }
}

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

#[derive(Debug, Clone, uniffi::Record)]
pub struct F64Pair {
    pub first: f64,
    pub second: f64,
}

#[derive(Debug, Clone, uniffi::Record)]
pub struct Point2 {
    pub x: f64,
    pub y: f64,
}

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
