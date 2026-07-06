use actix_web::{web, HttpResponse, Result};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::crypto_service::{
    BrandLevel, CryptoAssetWithNetwork, CryptoDeposit, CryptoNetwork, CryptoService, CryptoWithdrawal,
    CryptoError, DepositRequest, FeeCalculationResponse, WithdrawalRequest,
};

#[derive(Debug, Serialize)]
pub struct ApiResponse<T> {
    pub success: bool,
    pub data: Option<T>,
    pub error: Option<String>,
}

impl<T> ApiResponse<T> {
    pub fn success(data: T) -> Self {
        Self {
            success: true,
            data: Some(data),
            error: None,
        }
    }

    pub fn error(message: String) -> Self {
        Self {
            success: false,
            data: None,
            error: Some(message),
        }
    }
}

// ============== User Endpoints ==============

// Get all networks
pub async fn get_networks(
    service: web::Data<CryptoService>,
) -> Result<HttpResponse> {
    match service.get_all_networks().await {
        Ok(networks) => Ok(HttpResponse::Ok().json(ApiResponse::success(networks))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Get all assets with network info
pub async fn get_assets(
    service: web::Data<CryptoService>,
) -> Result<HttpResponse> {
    match service.get_all_assets_with_networks().await {
        Ok(assets) => Ok(HttpResponse::Ok().json(ApiResponse::success(assets))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Get networks for a specific asset
pub async fn get_asset_networks(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
) -> Result<HttpResponse> {
    let asset_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid asset ID".to_string()))),
    };

    match service.get_networks_by_asset(asset_id).await {
        Ok(networks) => Ok(HttpResponse::Ok().json(ApiResponse::success(networks))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Get deposit address
#[derive(Debug, Deserialize)]
pub struct DepositAddressRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
}

pub async fn get_deposit_address(
    service: web::Data<CryptoService>,
    user_id: Uuid,
    req: web::Json<DepositAddressRequest>,
) -> Result<HttpResponse> {
    match service.get_deposit_address(user_id, req.asset_id, req.network_id).await {
        Ok(address) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "address": address })))),
        Err(e) => Ok(HttpResponse::BadRequest().json(ApiResponse::error(e.to_string()))),
    }
}

// Get user deposits
pub async fn get_user_deposits(
    service: web::Data<CryptoService>,
    user_id: Uuid,
) -> Result<HttpResponse> {
    match service.get_user_deposits(user_id).await {
        Ok(deposits) => Ok(HttpResponse::Ok().json(ApiResponse::success(deposits))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Create withdrawal
pub async fn create_withdrawal(
    service: web::Data<CryptoService>,
    user_id: Uuid,
    req: web::Json<WithdrawalRequest>,
) -> Result<HttpResponse> {
    match service.create_withdrawal(user_id, req.asset_id, req.network_id, req.amount, req.address.clone()).await {
        Ok(withdrawal) => Ok(HttpResponse::Created().json(ApiResponse::success(withdrawal))),
        Err(e) => Ok(HttpResponse::BadRequest().json(ApiResponse::error(e.to_string()))),
    }
}

// Get user withdrawals
pub async fn get_user_withdrawals(
    service: web::Data<CryptoService>,
    user_id: Uuid,
) -> Result<HttpResponse> {
    match service.get_user_withdrawals(user_id).await {
        Ok(withdrawals) => Ok(HttpResponse::Ok().json(ApiResponse::success(withdrawals))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Calculate fees
#[derive(Debug, Deserialize)]
pub struct FeeCalculationRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub amount: f64,
}

pub async fn calculate_fees(
    service: web::Data<CryptoService>,
    user_id: Uuid,
    req: web::Json<FeeCalculationRequest>,
) -> Result<HttpResponse> {
    // Get user level (simplified - in production get from database)
    let user_level = 0;

    let deposit_fee = service.calculate_deposit_fee(req.asset_id, req.network_id, req.amount, user_level).await.unwrap_or(0.0);
    let withdrawal_fee = service.calculate_withdrawal_fee(req.asset_id, req.network_id, req.amount, user_level).await.unwrap_or(0.0);

    Ok(HttpResponse::Ok().json(ApiResponse::success(FeeCalculationResponse {
        deposit_fee,
        withdrawal_fee,
    })))
}

// ============== Admin Endpoints ==============

// Admin: Get all networks
pub async fn admin_get_networks(
    service: web::Data<CryptoService>,
) -> Result<HttpResponse> {
    match service.get_all_networks().await {
        Ok(networks) => Ok(HttpResponse::Ok().json(ApiResponse::success(networks))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Create network
#[derive(Debug, Deserialize)]
pub struct CreateNetworkRequest {
    pub name: String,
    pub chain_id: Option<String>,
    pub symbol: String,
    pub explorer_url: Option<String>,
    pub min_confirmation_blocks: Option<i32>,
}

pub async fn admin_create_network(
    service: web::Data<CryptoService>,
    req: web::Json<CreateNetworkRequest>,
) -> Result<HttpResponse> {
    match service.create_network(
        req.name.clone(),
        req.chain_id.clone(),
        req.symbol.clone(),
        req.explorer_url.clone(),
        req.min_confirmation_blocks.unwrap_or(12),
    ).await {
        Ok(network) => Ok(HttpResponse::Created().json(ApiResponse::success(network))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Update network
#[derive(Debug, Deserialize)]
pub struct UpdateNetworkRequest {
    pub name: String,
    pub chain_id: Option<String>,
    pub symbol: String,
    pub explorer_url: Option<String>,
    pub is_active: bool,
    pub is_deposit_enabled: bool,
    pub is_withdrawal_enabled: bool,
    pub min_confirmation_blocks: i32,
}

pub async fn admin_update_network(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
    req: web::Json<UpdateNetworkRequest>,
) -> Result<HttpResponse> {
    let network_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid network ID".to_string()))),
    };

    match service.update_network(
        network_id,
        req.name.clone(),
        req.chain_id.clone(),
        req.symbol.clone(),
        req.explorer_url.clone(),
        req.is_active,
        req.is_deposit_enabled,
        req.is_withdrawal_enabled,
        req.min_confirmation_blocks,
    ).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Network updated" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Get all assets
pub async fn admin_get_assets(
    service: web::Data<CryptoService>,
) -> Result<HttpResponse> {
    match service.get_all_assets_with_networks().await {
        Ok(assets) => Ok(HttpResponse::Ok().json(ApiResponse::success(assets))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Create asset
#[derive(Debug, Deserialize)]
pub struct CreateAssetRequest {
    pub name: String,
    pub symbol: String,
    pub decimals: Option<i32>,
    pub min_deposit_amount: Option<f64>,
    pub min_withdrawal_amount: Option<f64>,
}

pub async fn admin_create_asset(
    service: web::Data<CryptoService>,
    req: web::Json<CreateAssetRequest>,
) -> Result<HttpResponse> {
    match service.create_asset(
        req.name.clone(),
        req.symbol.clone(),
        req.decimals.unwrap_or(18),
        req.min_deposit_amount.unwrap_or(0.0),
        req.min_withdrawal_amount.unwrap_or(0.0),
    ).await {
        Ok(asset) => Ok(HttpResponse::Created().json(ApiResponse::success(asset))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Update asset
#[derive(Debug, Deserialize)]
pub struct UpdateAssetRequest {
    pub name: String,
    pub decimals: i32,
    pub min_deposit_amount: f64,
    pub min_withdrawal_amount: f64,
    pub is_active: bool,
}

pub async fn admin_update_asset(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
    req: web::Json<UpdateAssetRequest>,
) -> Result<HttpResponse> {
    let asset_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid asset ID".to_string()))),
    };

    match service.update_asset(
        asset_id,
        req.name.clone(),
        req.decimals,
        req.min_deposit_amount,
        req.min_withdrawal_amount,
        req.is_active,
    ).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Asset updated" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Link asset to network
#[derive(Debug, Deserialize)]
pub struct LinkAssetNetworkRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub deposit_enabled: Option<bool>,
    pub withdrawal_enabled: Option<bool>,
    pub contract_address: Option<String>,
}

pub async fn admin_link_asset_to_network(
    service: web::Data<CryptoService>,
    req: web::Json<LinkAssetNetworkRequest>,
) -> Result<HttpResponse> {
    match service.link_asset_to_network(
        req.asset_id,
        req.network_id,
        req.deposit_enabled.unwrap_or(true),
        req.withdrawal_enabled.unwrap_or(true),
        req.contract_address.clone(),
    ).await {
        Ok(link) => Ok(HttpResponse::Created().json(ApiResponse::success(link))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Set network fee
#[derive(Debug, Deserialize)]
pub struct SetFeeRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub deposit_fee: Option<f64>,
    pub withdrawal_fee: Option<f64>,
    pub deposit_fee_percent: Option<f64>,
    pub withdrawal_fee_percent: Option<f64>,
}

pub async fn admin_set_network_fee(
    service: web::Data<CryptoService>,
    req: web::Json<SetFeeRequest>,
) -> Result<HttpResponse> {
    match service.set_network_fee(
        req.asset_id,
        req.network_id,
        req.deposit_fee.unwrap_or(0.0),
        req.withdrawal_fee.unwrap_or(0.0),
        req.deposit_fee_percent.unwrap_or(0.0),
        req.withdrawal_fee_percent.unwrap_or(0.0),
    ).await {
        Ok(fee) => Ok(HttpResponse::Ok().json(ApiResponse::success(fee))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Get brand levels
pub async fn admin_get_brand_levels(
    service: web::Data<CryptoService>,
) -> Result<HttpResponse> {
    match service.get_all_brand_levels().await {
        Ok(levels) => Ok(HttpResponse::Ok().json(ApiResponse::success(levels))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Update brand level
#[derive(Debug, Deserialize)]
pub struct UpdateBrandLevelRequest {
    pub name: String,
    pub level: i32,
    pub deposit_fee_discount_percent: f64,
    pub withdrawal_fee_discount_percent: f64,
    pub is_active: bool,
}

pub async fn admin_update_brand_level(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
    req: web::Json<UpdateBrandLevelRequest>,
) -> Result<HttpResponse> {
    let level_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid brand level ID".to_string()))),
    };

    match service.update_brand_level(
        level_id,
        req.name.clone(),
        req.level,
        req.deposit_fee_discount_percent,
        req.withdrawal_fee_discount_percent,
        req.is_active,
    ).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Brand level updated" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Set admin wallet
#[derive(Debug, Deserialize)]
pub struct SetWalletRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub address: String,
    pub private_key_encrypted: Option<String>,
    pub is_active: Option<bool>,
}

pub async fn admin_set_admin_wallet(
    service: web::Data<CryptoService>,
    req: web::Json<SetWalletRequest>,
) -> Result<HttpResponse> {
    match service.set_admin_wallet(
        req.asset_id,
        req.network_id,
        req.address.clone(),
        req.private_key_encrypted.clone(),
        req.is_active.unwrap_or(true),
    ).await {
        Ok(wallet) => Ok(HttpResponse::Ok().json(ApiResponse::success(wallet))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Get all deposits
#[derive(Debug, Deserialize)]
pub struct ListQuery {
    pub status: Option<String>,
    pub page: Option<i32>,
    pub limit: Option<i32>,
}

pub async fn admin_get_deposits(
    service: web::Data<CryptoService>,
    query: web::Query<ListQuery>,
) -> Result<HttpResponse> {
    let limit = query.limit.unwrap_or(20);
    let offset = ((query.page.unwrap_or(1)) - 1) * limit;

    match service.get_all_deposits(query.status.clone(), limit, offset).await {
        Ok((deposits, total)) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({
            "deposits": deposits,
            "total": total,
            "page": query.page.unwrap_or(1),
            "limit": limit,
        })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Get all withdrawals
pub async fn admin_get_withdrawals(
    service: web::Data<CryptoService>,
    query: web::Query<ListQuery>,
) -> Result<HttpResponse> {
    let limit = query.limit.unwrap_or(20);
    let offset = ((query.page.unwrap_or(1)) - 1) * limit;

    match service.get_all_withdrawals(query.status.clone(), limit, offset).await {
        Ok((withdrawals, total)) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({
            "withdrawals": withdrawals,
            "total": total,
            "page": query.page.unwrap_or(1),
            "limit": limit,
        })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Confirm deposit
pub async fn admin_confirm_deposit(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
) -> Result<HttpResponse> {
    let deposit_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid deposit ID".to_string()))),
    };

    match service.confirm_deposit(deposit_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Deposit confirmed" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Confirm withdrawal
#[derive(Debug, Deserialize)]
pub struct ConfirmWithdrawalRequest {
    pub tx_hash: String,
}

pub async fn admin_confirm_withdrawal(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
    req: web::Json<ConfirmWithdrawalRequest>,
) -> Result<HttpResponse> {
    let withdrawal_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid withdrawal ID".to_string()))),
    };

    match service.confirm_withdrawal(withdrawal_id, req.tx_hash.clone()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Withdrawal confirmed" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Cancel withdrawal
pub async fn admin_cancel_withdrawal(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
) -> Result<HttpResponse> {
    let withdrawal_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid withdrawal ID".to_string()))),
    };

    match service.cancel_withdrawal(withdrawal_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "Withdrawal cancelled" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}

// Admin: Update user brand level
#[derive(Debug, Deserialize)]
pub struct UpdateUserBrandLevelRequest {
    pub level: i32,
}

pub async fn admin_update_user_brand_level(
    service: web::Data<CryptoService>,
    path: web::Path<String>,
    req: web::Json<UpdateUserBrandLevelRequest>,
) -> Result<HttpResponse> {
    let user_id = match Uuid::parse_str(&path) {
        Ok(id) => id,
        Err(_) => return Ok(HttpResponse::BadRequest().json(ApiResponse::error("Invalid user ID".to_string()))),
    };

    match service.update_user_brand_level(user_id, req.level).await {
        Ok(_) => Ok(HttpResponse::Ok().json(ApiResponse::success(serde_json::json!({ "message": "User brand level updated" })))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(ApiResponse::error(e.to_string()))),
    }
}
