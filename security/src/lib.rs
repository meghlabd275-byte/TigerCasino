//! TigerCasino Security Module
//! 
//! Provides cryptographic operations for the casino platform including:
//! - Secure random number generation
//! - Password hashing with Argon2
//! - Encryption/Decryption with AES-GCM
//! - HMAC signature verification
//! - TOTP for two-factor authentication
//! - Fraud detection patterns
//! - Crypto deposit/withdrawal services
//! - Database operations
//! - HTTP handlers and middleware
//! - Game provider API integrations

// Re-export main modules
pub mod crypto_service;
pub mod crypto_handler;
pub mod database;
pub mod config;
pub mod models;
pub mod handlers;
pub mod middleware;
pub mod provider;

pub use crypto_service::CryptoService;
pub use crypto_handler::*;
pub use database::Database;
pub use config::{Config, Claims};
pub use models::*;
pub use handlers::*;
pub use middleware::*;
pub use provider::*;

use rand::Rng;
use sha2::{Sha256, Sha512, Digest};
use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm, Nonce,
};
use argon2::{
    password_hash::{PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};
use hmac::{Hmac, Mac};
use std::time::{SystemTime, UNIX_EPOCH};

/// Type alias for HMAC-SHA256
type HmacSha256 = Hmac<Sha256>;

/// Security module for TigerCasino
pub struct Security;

impl Security {
    /// Generate a cryptographically secure random number
    pub fn generate_random(min: u64, max: u64) -> u64 {
        let mut rng = rand::thread_rng();
        let range = max - min + 1;
        min + (rng.gen::<u64>() % range)
    }

    /// Generate a random string of specified length
    pub fn generate_random_string(length: usize) -> String {
        let chars: Vec<char> = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789".chars().collect();
        let mut rng = rand::thread_rng();
        (0..length).map(|_| chars[rng.gen::<usize>() % chars.len()]).collect()
    }

    /// Generate a random UUID
    pub fn generate_uuid() -> String {
        uuid::Uuid::new_v4().to_string()
    }

    /// Hash a password using Argon2
    pub fn hash_password(password: &str) -> Result<String, String> {
        let salt = SaltString::generate(&mut rand::thread_rng());
        let argon2 = Argon2::default();
        
        match argon2.hash_password(password.as_bytes(), &salt) {
            Ok(hash) => Ok(hash.to_string()),
            Err(e) => Err(format!("Failed to hash password: {}", e)),
        }
    }

    /// Verify a password against a hash
    pub fn verify_password(password: &str, hash: &str) -> Result<bool, String> {
        let parsed_hash = PasswordHash::new(hash)
            .map_err(|e| format!("Invalid hash: {}", e))?;
        
        Ok(Argon2::default()
            .verify_password(password.as_bytes(), &parsed_hash)
            .is_ok())
    }

    /// Encrypt data using AES-256-GCM
    pub fn encrypt(plaintext: &[u8], key: &[u8; 32]) -> Result<Vec<u8>, String> {
        let cipher = Aes256Gcm::new_from_slice(key)
            .map_err(|e| format!("Invalid key: {}", e))?;
        
        let mut rng = rand::thread_rng();
        let nonce_bytes: [u8; 12] = rng.gen();
        let nonce = Nonce::from_slice(&nonce_bytes);
        
        match cipher.encrypt(nonce, plaintext) {
            Ok(ciphertext) => {
                let mut result = nonce_bytes.to_vec();
                result.extend(ciphertext);
                Ok(result)
            }
            Err(e) => Err(format!("Encryption failed: {}", e)),
        }
    }

    /// Decrypt data using AES-256-GCM
    pub fn decrypt(ciphertext: &[u8], key: &[u8; 32]) -> Result<Vec<u8>, String> {
        if ciphertext.len() < 12 {
            return Err("Ciphertext too short".to_string());
        }

        let cipher = Aes256Gcm::new_from_slice(key)
            .map_err(|e| format!("Invalid key: {}", e))?;
        
        let nonce = Nonce::from_slice(&ciphertext[..12]);
        let encrypted = &ciphertext[12..];
        
        match cipher.decrypt(nonce, encrypted) {
            Ok(plaintext) => Ok(plaintext),
            Err(e) => Err(format!("Decryption failed: {}", e)),
        }
    }

    /// Create HMAC-SHA256 signature
    pub fn create_hmac(data: &[u8], key: &[u8]) -> String {
        let mut mac = HmacSha256::new_from_slice(key)
            .expect("HMAC can take key of any size");
        mac.update(data);
        hex::encode(mac.finalize().into_bytes())
    }

    /// Verify HMAC-SHA256 signature
    pub fn verify_hmac(data: &[u8], key: &[u8], signature: &str) -> bool {
        let expected = Self::create_hmac(data, key);
        expected == signature
    }

    /// Generate SHA-256 hash
    pub fn sha256(data: &[u8]) -> String {
        let mut hasher = Sha256::new();
        hasher.update(data);
        hex::encode(hasher.finalize())
    }

    /// Generate SHA-512 hash
    pub fn sha512(data: &[u8]) -> String {
        let mut hasher = Sha512::new();
        hasher.update(data);
        hex::encode(hasher.finalize())
    }

    /// Generate TOTP code (simplified implementation)
    pub fn generate_totp(secret: &str, timestamp: u64) -> String {
        // Simplified TOTP - in production use a proper TOTP library
        let time_step = timestamp / 30;
        let data = format!("{}:{}", secret, time_step);
        let hash = Self::sha256(data.as_bytes());
        let num = u64::from_str_radix(&hash[..8], 16).unwrap_or(0);
        format!("{:06}", num % 1_000_000)
    }

    /// Verify TOTP code
    pub fn verify_totp(secret: &str, code: &str) -> bool {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        // Check current and previous time steps
        let current = Self::generate_totp(secret, now);
        let previous = Self::generate_totp(secret, now - 30);
        
        code == current || code == previous
    }

    /// Encode to base64
    pub fn base64_encode(data: &[u8]) -> String {
        BASE64.encode(data)
    }

    /// Decode from base64
    pub fn base64_decode(data: &str) -> Result<Vec<u8>, String> {
        BASE64.decode(data).map_err(|e| format!("Base64 decode error: {}", e))
    }

    /// Detect suspicious patterns in betting
    pub fn detect_fraud_pattern(bet_amount: f64, balance: f64, win_rate: f64, bet_count: u64) -> FraudRisk {
        let mut score = 0.0;

        // High bet relative to balance
        if bet_amount / balance > 0.5 {
            score += 0.3;
        }

        // Suspiciously high win rate
        if win_rate > 0.95 && bet_count > 10 {
            score += 0.4;
        }

        // Very large bets
        if bet_amount > 10000.0 {
            score += 0.2;
        }

        // Multiple rapid bets
        if bet_count > 100 && win_rate > 0.8 {
            score += 0.3;
        }

        match score {
            s if s >= 0.7 => FraudRisk::High,
            s if s >= 0.4 => FraudRisk::Medium,
            _ => FraudRisk::Low,
        }
    }
}

/// Fraud risk level
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum FraudRisk {
    Low,
    Medium,
    High,
}

impl FraudRisk {
    pub fn as_str(&self) -> &'static str {
        match self {
            FraudRisk::Low => "low",
            FraudRisk::Medium => "medium",
            FraudRisk::High => "high",
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_random_number() {
        let num = Security::generate_random(1, 100);
        assert!(num >= 1 && num <= 100);
    }

    #[test]
    fn test_password_hash() {
        let hash = Security::hash_password("test_password123").unwrap();
        assert!(Security::verify_password("test_password123", &hash).unwrap());
    }

    #[test]
    fn test_encryption() {
        let key = [0u8; 32];
        let data = b"Hello, TigerCasino!";
        let encrypted = Security::encrypt(data, &key).unwrap();
        let decrypted = Security::decrypt(&encrypted, &key).unwrap();
        assert_eq!(data.to_vec(), decrypted);
    }

    #[test]
    fn test_hmac() {
        let data = b"test data";
        let key = b"secret key";
        let signature = Security::create_hmac(data, key);
        assert!(Security::verify_hmac(data, key, &signature));
    }
}
