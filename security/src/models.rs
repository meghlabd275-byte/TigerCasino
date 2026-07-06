use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

// ============== User Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: uuid::Uuid,
    pub email: String,
    pub username: String,
    #[serde(skip_serializing)]
    pub password_hash: String,
    pub wallet_address: Option<String>,
    pub wallet_type: Option<String>,
    pub balance: f64,
    pub bonus_balance: f64,
    pub vip_level: i32,
    pub kyc_status: String,
    pub is_verified: bool,
    pub is_admin: bool,
    pub is_banned: bool,
    pub ban_reason: Option<String>,
    #[serde(skip_serializing)]
    pub two_fa_secret: Option<String>,
    pub is_2fa_enabled: bool,
    pub email_verified: bool,
    pub phone_verified: bool,
    pub phone: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_login: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserResponse {
    pub id: uuid::Uuid,
    pub email: String,
    pub username: String,
    pub balance: f64,
    pub bonus_balance: f64,
    pub vip_level: i32,
    pub is_verified: bool,
    pub is_admin: bool,
}

impl From<User> for UserResponse {
    fn from(user: User) -> Self {
        Self {
            id: user.id,
            email: user.email,
            username: user.username,
            balance: user.balance,
            bonus_balance: user.bonus_balance,
            vip_level: user.vip_level,
            is_verified: user.is_verified,
            is_admin: user.is_admin,
        }
    }
}

// ============== Transaction Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Transaction {
    pub id: uuid::Uuid,
    pub user_id: uuid::Uuid,
    #[serde(rename = "type")]
    pub transaction_type: String,
    pub amount: f64,
    pub currency: String,
    pub status: String,
    pub tx_hash: Option<String>,
    pub address: Option<String>,
    pub fee: f64,
    pub chain: Option<String>,
    pub platform_fee: Option<f64>,
    pub brand_revenue: Option<f64>,
    pub created_at: DateTime<Utc>,
    pub processed_at: Option<DateTime<Utc>>,
}

// ============== Game Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Game {
    pub id: uuid::Uuid,
    pub name: String,
    pub game_type: String,
    pub provider: Option<String>,
    pub rtp: Option<f64>,
    pub min_bet: Option<f64>,
    pub max_bet: Option<f64>,
    pub is_active: bool,
    pub thumbnail_url: Option<String>,
    pub game_data: Option<serde_json::Value>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Bet {
    pub id: uuid::Uuid,
    pub user_id: uuid::Uuid,
    pub game_id: uuid::Uuid,
    pub bet_amount: f64,
    pub win_amount: f64,
    pub multiplier: f64,
    pub game_data: Option<serde_json::Value>,
    pub status: String,
    pub settled_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
}

// ============== Session Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Session {
    pub id: uuid::Uuid,
    pub user_id: uuid::Uuid,
    pub token: String,
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
    pub expires_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
}

// ============== White Label Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WhiteLabel {
    pub id: uuid::Uuid,
    pub name: String,
    pub domain: String,
    pub brand_color: Option<String>,
    pub logo_url: Option<String>,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
}

// ============== Request/Response Types ==============

#[derive(Debug, Deserialize)]
pub struct RegisterRequest {
    pub email: String,
    pub username: String,
    pub password: String,
}

#[derive(Debug, Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(Debug, Serialize)]
pub struct AuthResponse {
    pub token: String,
    pub user: UserResponse,
}

#[derive(Debug, Deserialize)]
pub struct UpdateProfileRequest {
    pub email: Option<String>,
    pub phone: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct ChangePasswordRequest {
    pub current_password: String,
    pub new_password: String,
}

// ============== Error Types ==============

#[derive(Debug, Serialize)]
pub struct ErrorResponse {
    pub error: String,
    pub message: String,
}

impl ErrorResponse {
    pub fn new(error: &str, message: &str) -> Self {
        Self {
            error: error.to_string(),
            message: message.to_string(),
        }
    }
}
