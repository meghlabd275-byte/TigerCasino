use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use rand::Rng;
use argon2::{
    password_hash::{PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};

#[no_mangle]
pub extern "C" fn rust_hash_password(password: *const c_char) -> *mut c_char {
    let c_str = unsafe { CStr::from_ptr(password) };
    let password_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return std::ptr::null_mut(),
    };

    let salt = SaltString::generate(&mut rand::thread_rng());
    let argon2 = Argon2::default();

    match argon2.hash_password(password_str.as_bytes(), &salt) {
        Ok(hash) => CString::new(hash.to_string()).unwrap().into_raw(),
        Err(_) => std::ptr::null_mut(),
    }
}

#[no_mangle]
pub extern "C" fn rust_verify_password(password: *const c_char, hash: *const c_char) -> bool {
    let c_password = unsafe { CStr::from_ptr(password) };
    let c_hash = unsafe { CStr::from_ptr(hash) };

    let password_str = match c_password.to_str() {
        Ok(s) => s,
        Err(_) => return false,
    };
    let hash_str = match c_hash.to_str() {
        Ok(s) => s,
        Err(_) => return false,
    };

    let parsed_hash = match PasswordHash::new(hash_str) {
        Ok(h) => h,
        Err(_) => return false,
    };

    Argon2::default()
        .verify_password(password_str.as_bytes(), &parsed_hash)
        .is_ok()
}

#[no_mangle]
pub extern "C" fn rust_generate_random(min: u64, max: u64) -> u64 {
    let mut rng = rand::thread_rng();
    let range = max - min + 1;
    min + (rng.gen::<u64>() % range)
}

#[no_mangle]
pub extern "C" fn rust_free_string(s: *mut c_char) {
    if s.is_null() {
        return;
    }
    unsafe {
        drop(CString::from_raw(s));
    }
}

#[no_mangle]
pub extern "C" fn rust_generate_outcome(server_seed: *const c_char, client_seed: *const c_char, nonce: i32) -> f64 {
    let c_server = unsafe { CStr::from_ptr(server_seed) };
    let c_client = unsafe { CStr::from_ptr(client_seed) };

    let combined = format!("{}:{}:{}", c_server.to_str().unwrap_or(""), c_client.to_str().unwrap_or(""), nonce);

    use sha2::{Sha256, Digest};
    let mut hasher = Sha256::new();
    hasher.update(combined.as_bytes());
    let result = hasher.finalize();

    let mut val: u64 = 0;
    for i in 0..8 {
        val = (val << 8) | result[i] as u64;
    }

    (val as f64) / (u64::MAX as f64)
}
