#[uniffi::export]
pub fn add(a: i32, b: i32) -> i32 {
    uniffi_testing_core::add(a, b)
}
