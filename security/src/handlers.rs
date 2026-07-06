use actix_web::{web, HttpResponse, Result};
use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::Utc;

use crate::database::Database;
use crate::crypto_service::CryptoService;
use crate::models::{User, UserResponse, AuthResponse, ErrorResponse};
use crate::config::Claims;

// ============== Auth Handlers ==============

#[derive(Debug, Deserialize)]
pub struct RegisterReq {
    pub email: String,
    pub username: String,
    pub password: String,
}

#[derive(Debug, Deserialize)]
pub struct LoginReq {
    pub email: String,
    pub password: String,
}

// Simple password hashing (in production use argon2)
fn hash_password(password: &str) -> String {
    use sha2::{Sha256, Digest};
    let mut hasher = Sha256::new();
    hasher.update(password.as_bytes());
    hex::encode(hasher.finalize())
}

fn verify_password(password: &str, hash: &str) -> bool {
    hash_password(password) == hash
}

pub async fn register(
    db: web::Data<Database>,
    req: web::Json<RegisterReq>,
) -> Result<HttpResponse> {
    // Check if email already exists
    if db.get_user_by_email(&req.email).await.ok().flatten().is_some() {
        return Ok(HttpResponse::BadRequest().json(ErrorResponse::new("error", "Email already registered")));
    }

    // Check if username already exists
    if db.get_user_by_username(&req.username).await.ok().flatten().is_some() {
        return Ok(HttpResponse::BadRequest().json(ErrorResponse::new("error", "Username already taken")));
    }

    let password_hash = hash_password(&req.password);

    match db.create_user(&req.email, &req.username, &password_hash).await {
        Ok(user) => {
            let user_response: UserResponse = user.into();
            Ok(HttpResponse::Created().json(user_response))
        }
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

pub async fn login(
    db: web::Data<Database>,
    req: web::Json<LoginReq>,
) -> Result<HttpResponse> {
    let user = match db.get_user_by_email(&req.email).await {
        Ok(Some(user)) => user,
        Ok(None) => return Ok(HttpResponse::Unauthorized().json(ErrorResponse::new("error", "Invalid credentials"))),
        Err(e) => return Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    };

    if !verify_password(&req.password, &user.password_hash) {
        return Ok(HttpResponse::Unauthorized().json(ErrorResponse::new("error", "Invalid credentials")));
    }

    if user.is_banned {
        return Ok(HttpResponse::Forbidden().json(ErrorResponse::new("error", "Account banned")));
    }

    // Generate JWT token (simplified - in production use proper JWT library)
    let token = format!("{}.{}", user.id, Utc::now().timestamp());
    
    // Create session
    let expires_at = Utc::now() + chrono::Duration::hours(24 * 7); // 7 days
    let _ = db.create_session(user.id, &token, expires_at, None, None).await;

    let user_response = UserResponse::from(user);

    Ok(HttpResponse::Ok().json(AuthResponse {
        token,
        user: user_response,
    }))
}

pub async fn logout(
    db: web::Data<Database>,
    token: web::ReqData<String>,
) -> Result<HttpResponse> {
    let _ = db.delete_session(&token).await;
    Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Logged out" })))
}

// ============== User Handlers ==============

pub async fn get_profile(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
) -> Result<HttpResponse> {
    match db.get_user(*user_id).await {
        Ok(Some(user)) => {
            let user_response: UserResponse = user.into();
            Ok(HttpResponse::Ok().json(user_response))
        }
        Ok(None) => Ok(HttpResponse::NotFound().json(ErrorResponse::new("error", "User not found"))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

#[derive(Debug, Deserialize)]
pub struct UpdateProfileReq {
    pub email: Option<String>,
    pub phone: Option<String>,
}

pub async fn update_profile(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
    req: web::Json<UpdateProfileReq>,
) -> Result<HttpResponse> {
    // In production, implement proper update logic
    Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Profile updated" })))
}

// ============== Wallet Handlers ==============

pub async fn get_balance(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
) -> Result<HttpResponse> {
    match db.get_user(*user_id).await {
        Ok(Some(user)) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "balance": user.balance,
            "bonus_balance": user.bonus_balance,
            "vip_level": user.vip_level,
        }))),
        Ok(None) => Ok(HttpResponse::NotFound().json(ErrorResponse::new("error", "User not found"))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

#[derive(Debug, Deserialize)]
pub struct DepositReq {
    pub amount: f64,
    pub currency: String,
}

pub async fn create_deposit(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
    req: web::Json<DepositReq>,
) -> Result<HttpResponse> {
    // Create transaction record
    match db.create_transaction(
        *user_id,
        "deposit",
        req.amount,
        &req.currency,
        "pending",
        None,
        0.0,
    ).await {
        Ok(tx) => Ok(HttpResponse::Created().json(tx)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

#[derive(Debug, Deserialize)]
pub struct WithdrawReq {
    pub amount: f64,
    pub currency: String,
    pub address: String,
}

pub async fn create_withdrawal(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
    req: web::Json<WithdrawReq>,
) -> Result<HttpResponse> {
    // Check balance
    let user = match db.get_user(*user_id).await {
        Ok(Some(u)) => u,
        Ok(None) => return Ok(HttpResponse::NotFound().json(ErrorResponse::new("error", "User not found"))),
        Err(e) => return Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    };

    if user.balance < req.amount {
        return Ok(HttpResponse::BadRequest().json(ErrorResponse::new("error", "Insufficient balance")));
    }

    // Deduct balance
    if let Err(e) = db.update_balance(*user_id, -req.amount).await {
        return Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string())));
    }

    // Create transaction record
    match db.create_transaction(
        *user_id,
        "withdrawal",
        req.amount,
        &req.currency,
        "pending",
        Some(&req.address),
        0.0,
    ).await {
        Ok(tx) => Ok(HttpResponse::Created().json(tx)),
        Err(e) => {
            // Refund balance on failure
            let _ = db.update_balance(*user_id, req.amount).await;
            Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string())))
        }
    }
}

pub async fn get_transactions(
    db: web::Data<Database>,
    user_id: web::ReqData<Uuid>,
) -> Result<HttpResponse> {
    match db.get_user_transactions(*user_id, 50).await {
        Ok(txs) => Ok(HttpResponse::Ok().json(txs)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

// ============== Admin Handlers ==============

pub async fn admin_get_users(
    db: web::Data<Database>,
    query: web::Query<std::collections::HashMap<String, String>>,
) -> Result<HttpResponse> {
    let limit: i32 = query.get("limit").and_then(|v| v.parse().ok()).unwrap_or(20);
    let offset: i32 = query.get("offset").and_then(|v| v.parse().ok()).unwrap_or(0);

    match db.get_all_users(limit, offset).await {
        Ok((users, total)) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "users": users,
            "total": total,
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

#[derive(Debug, Deserialize)]
pub struct BanUserReq {
    pub banned: bool,
}

pub async fn admin_ban_user(
    db: web::Data<Database>,
    path: web::Path<String>,
    req: web::Json<BanUserReq>,
) -> Result<HttpResponse> {
    let user_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ErrorResponse::new("error", "Invalid user ID"))),
    };

    match db.set_user_banned(user_id, req.banned).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "User banned status updated" }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}

pub async fn admin_get_audit_logs(
    db: web::Data<Database>,
    path: web::Path<String>,
) -> Result<HttpResponse> {
    let user_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ErrorResponse::new("error", "Invalid user ID"))),
    };

    match db.get_user_audit_logs(user_id, 100).await {
        Ok(logs) => Ok(HttpResponse::Ok().json(logs)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ErrorResponse::new("error", &e.to_string()))),
    }
}
