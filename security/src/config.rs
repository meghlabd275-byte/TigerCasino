use serde::{Deserialize, Serialize};
use std::env;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigError {
    #[error("Environment variable not set: {0}")]
    MissingEnvVar(String),
    #[error("Parse error: {0}")]
    ParseError(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    // Server
    pub server_host: String,
    pub server_port: u16,
    
    // Database
    pub database_url: String,
    pub db_max_connections: u32,
    pub db_min_connections: u32,
    
    // JWT
    pub jwt_secret: String,
    pub jwt_expiration_hours: i64,
    
    // Security
    pub bcrypt_rounds: u32,
    pub session_expiration_hours: i64,
    
    // External Services
    pub smtp_host: Option<String>,
    pub smtp_port: Option<u16>,
    pub smtp_username: Option<String>,
    pub smtp_password: Option<String>,
    pub smtp_from: Option<String>,
    
    // Crypto (optional)
    pub crypto_rpc_urls: Option<String>,
}

impl Config {
    pub fn from_env() -> Result<Self, ConfigError> {
        Ok(Config {
            // Server
            server_host: env::var("SERVER_HOST").unwrap_or_else(|_| "0.0.0.0".to_string()),
            server_port: env::var("SERVER_PORT")
                .unwrap_or_else(|_| "8080".to_string())
                .parse()
                .map_err(|_| ConfigError::ParseError("SERVER_PORT".to_string()))?,
            
            // Database
            database_url: env::var("DATABASE_URL")
                .map_err(|_| ConfigError::MissingEnvVar("DATABASE_URL".to_string()))?,
            db_max_connections: env::var("DB_MAX_CONNECTIONS")
                .unwrap_or_else(|_| "10".to_string())
                .parse()
                .unwrap_or(10),
            db_min_connections: env::var("DB_MIN_CONNECTIONS")
                .unwrap_or_else(|_| "2".to_string())
                .parse()
                .unwrap_or(2),
            
            // JWT
            jwt_secret: env::var("JWT_SECRET")
                .map_err(|_| ConfigError::MissingEnvVar("JWT_SECRET".to_string()))?,
            jwt_expiration_hours: env::var("JWT_EXPIRATION_HOURS")
                .unwrap_or_else(|_| "24".to_string())
                .parse()
                .unwrap_or(24),
            
            // Security
            bcrypt_rounds: env::var("BCRYPT_ROUNDS")
                .unwrap_or_else(|_| "12".to_string())
                .parse()
                .unwrap_or(12),
            session_expiration_hours: env::var("SESSION_EXPIRATION_HOURS")
                .unwrap_or_else(|_| "168".to_string()) // 7 days
                .parse()
                .unwrap_or(168),
            
            // External Services
            smtp_host: env::var("SMTP_HOST").ok(),
            smtp_port: env::var("SMTP_PORT").ok().and_then(|p| p.parse().ok()),
            smtp_username: env::var("SMTP_USERNAME").ok(),
            smtp_password: env::var("SMTP_PASSWORD").ok(),
            smtp_from: env::var("SMTP_FROM").ok(),
            
            // Crypto
            crypto_rpc_urls: env::var("CRYPTO_RPC_URLS").ok(),
        })
    }
    
    pub fn database_url(&self) -> String {
        self.database_url.clone()
    }
    
    pub fn jwt_secret(&self) -> String {
        self.jwt_secret.clone()
    }
}

// JWT Claims
#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String,       // User ID
    pub email: String,
    pub exp: i64,
    pub iat: i64,
    pub is_admin: bool,
}

impl Claims {
    pub fn new(user_id: &str, email: &str, is_admin: bool, expiration: i64) -> Self {
        let now = chrono::Utc::now().timestamp();
        Self {
            sub: user_id.to_string(),
            email: email.to_string(),
            exp: now + (expiration * 3600),
            iat: now,
            is_admin,
        }
    }
    
    pub fn is_expired(&self) -> bool {
        chrono::Utc::now().timestamp() > self.exp
    }
}
