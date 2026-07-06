use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use std::sync::Arc;
use uuid::Uuid;
use chrono::{DateTime, Utc};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum CryptoError {
    #[error("Database error: {0}")]
    DatabaseError(String),
    #[error("Not found: {0}")]
    NotFound(String),
    #[error("Insufficient balance")]
    InsufficientBalance,
    #[error("Invalid operation: {0}")]
    InvalidOperation(String),
    #[error("Wallet not configured")]
    WalletNotConfigured,
}

impl From<sqlx::Error> for CryptoError {
    fn from(err: sqlx::Error) -> Self {
        CryptoError::DatabaseError(err.to_string())
    }
}

// ============== Data Models ==============

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoNetwork {
    pub id: Uuid,
    pub name: String,
    pub chain_id: Option<String>,
    pub symbol: String,
    pub explorer_url: Option<String>,
    pub rpc_url: Option<String>,
    pub is_withdrawal_enabled: bool,
    pub is_deposit_enabled: bool,
    pub is_active: bool,
    pub min_confirmation_blocks: i32,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoAsset {
    pub id: Uuid,
    pub name: String,
    pub symbol: String,
    pub decimals: i32,
    pub contract_address: Option<String>,
    pub is_active: bool,
    pub min_deposit_amount: f64,
    pub min_withdrawal_amount: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoAssetNetwork {
    pub id: Uuid,
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub deposit_enabled: bool,
    pub withdrawal_enabled: bool,
    pub contract_address: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    #[serde(default)]
    pub asset: Option<CryptoAsset>,
    #[serde(default)]
    pub network: Option<CryptoNetwork>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NetworkFee {
    pub id: Uuid,
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub deposit_fee: f64,
    pub withdrawal_fee: f64,
    pub deposit_fee_percent: f64,
    pub withdrawal_fee_percent: f64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BrandLevel {
    pub id: Uuid,
    pub name: String,
    pub level: i32,
    pub deposit_fee_discount_percent: f64,
    pub withdrawal_fee_discount_percent: f64,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoDeposit {
    pub id: Uuid,
    pub user_id: Uuid,
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub amount: f64,
    pub fee: f64,
    pub net_amount: f64,
    pub address: String,
    pub tx_hash: Option<String>,
    pub confirmations: i32,
    pub status: String,
    pub created_at: DateTime<Utc>,
    pub processed_at: Option<DateTime<Utc>>,
    #[serde(default)]
    pub asset: Option<CryptoAsset>,
    #[serde(default)]
    pub network: Option<CryptoNetwork>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoWithdrawal {
    pub id: Uuid,
    pub user_id: Uuid,
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub amount: f64,
    pub fee: f64,
    pub net_amount: f64,
    pub address: String,
    pub tx_hash: Option<String>,
    pub status: String,
    pub created_at: DateTime<Utc>,
    pub processed_at: Option<DateTime<Utc>>,
    #[serde(default)]
    pub asset: Option<CryptoAsset>,
    #[serde(default)]
    pub network: Option<CryptoNetwork>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AdminWalletAddress {
    pub id: Uuid,
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub address: String,
    pub private_key_encrypted: Option<String>,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoAssetWithNetwork {
    pub id: Uuid,
    pub name: String,
    pub symbol: String,
    pub decimals: i32,
    pub contract_address: Option<String>,
    pub is_active: bool,
    pub min_deposit_amount: f64,
    pub min_withdrawal_amount: f64,
    pub network_id: Uuid,
    pub deposit_enabled: bool,
    pub withdrawal_enabled: bool,
    pub deposit_fee: f64,
    pub withdrawal_fee: f64,
    pub deposit_fee_percent: f64,
    pub withdrawal_fee_percent: f64,
}

// ============== Request/Response Types ==============

#[derive(Debug, Deserialize)]
pub struct DepositRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
}

#[derive(Debug, Deserialize)]
pub struct WithdrawalRequest {
    pub asset_id: Uuid,
    pub network_id: Uuid,
    pub amount: f64,
    pub address: String,
}

#[derive(Debug, Serialize)]
pub struct DepositAddressResponse {
    pub address: String,
}

#[derive(Debug, Serialize)]
pub struct FeeCalculationResponse {
    pub deposit_fee: f64,
    pub withdrawal_fee: f64,
}

// ============== Crypto Service ==============

pub struct CryptoService {
    pool: Arc<PgPool>,
}

impl CryptoService {
    pub fn new(pool: Arc<PgPool>) -> Self {
        Self { pool }
    }

    // Get all active networks
    pub async fn get_all_networks(&self) -> Result<Vec<CryptoNetwork>, CryptoError> {
        let rows = sqlx::query_as!(
            CryptoNetwork,
            r#"SELECT id, name, chain_id, symbol, explorer_url, rpc_url,
               is_withdrawal_enabled, is_deposit_enabled, is_active,
               min_confirmation_blocks, created_at, updated_at
               FROM crypto_networks WHERE is_active = true"#
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Get all assets with network info
    pub async fn get_all_assets_with_networks(&self) -> Result<Vec<CryptoAssetWithNetwork>, CryptoError> {
        let rows = sqlx::query_as!(
            CryptoAssetWithNetwork,
            r#"SELECT ca.id, ca.name, ca.symbol, ca.decimals, ca.contract_address,
               ca.is_active, ca.min_deposit_amount, ca.min_withdrawal_amount,
               can.id as network_id,
               can.deposit_enabled, can.withdrawal_enabled, can.contract_address,
               COALESCE(nf.deposit_fee, 0) as deposit_fee,
               COALESCE(nf.withdrawal_fee, 0) as withdrawal_fee,
               COALESCE(nf.deposit_fee_percent, 0) as deposit_fee_percent,
               COALESCE(nf.withdrawal_fee_percent, 0) as withdrawal_fee_percent
               FROM crypto_assets ca
               JOIN crypto_asset_networks can ON ca.id = can.asset_id
               JOIN crypto_networks cn ON can.network_id = cn.id
               LEFT JOIN network_fees nf ON ca.id = nf.asset_id AND cn.id = nf.network_id
               WHERE ca.is_active = true AND cn.is_active = true
               ORDER BY ca.symbol, cn.name"#
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Get networks for a specific asset
    pub async fn get_networks_by_asset(&self, asset_id: Uuid) -> Result<Vec<CryptoAssetNetwork>, CryptoError> {
        let rows = sqlx::query_as!(
            CryptoAssetNetwork,
            r#"SELECT id, asset_id, network_id, deposit_enabled, withdrawal_enabled,
               contract_address, created_at, updated_at
               FROM crypto_asset_networks WHERE asset_id = $1"#,
            asset_id
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Get asset by ID
    pub async fn get_asset_by_id(&self, id: Uuid) -> Result<CryptoAsset, CryptoError> {
        let row = sqlx::query_as!(
            CryptoAsset,
            r#"SELECT id, name, symbol, decimals, contract_address, is_active,
               min_deposit_amount, min_withdrawal_amount, created_at, updated_at
               FROM crypto_assets WHERE id = $1"#,
            id
        )
        .fetch_one(&*self.pool)
        .await?;

        Ok(row)
    }

    // Calculate deposit fee
    pub async fn calculate_deposit_fee(
        &self,
        asset_id: Uuid,
        network_id: Uuid,
        amount: f64,
        user_level: i32,
    ) -> Result<f64, CryptoError> {
        let fee_row = sqlx::query!(
            r#"SELECT nf.deposit_fee, nf.deposit_fee_percent, bl.deposit_fee_discount_percent
               FROM network_fees nf
               LEFT JOIN brand_levels bl ON bl.level = $3 AND bl.is_active = true
               WHERE nf.asset_id = $1 AND nf.network_id = $2"#,
            asset_id,
            network_id,
            user_level
        )
        .fetch_optional(&*self.pool)
        .await?;

        let base_fee = fee_row.map(|r| r.deposit_fee.unwrap_or(0.0)).unwrap_or(0.0);
        let fee_percent = fee_row
            .and_then(|r| r.deposit_fee_percent)
            .unwrap_or(0.0);
        let discount = fee_row
            .and_then(|r| r.deposit_fee_discount_percent)
            .unwrap_or(0.0);

        let adjusted_percent = fee_percent * (1.0 - discount / 100.0);
        let total_fee = base_fee + (amount * adjusted_percent / 100.0);

        Ok((total_fee * 100_000_000.0).round() / 100_000_000.0)
    }

    // Calculate withdrawal fee
    pub async fn calculate_withdrawal_fee(
        &self,
        asset_id: Uuid,
        network_id: Uuid,
        amount: f64,
        user_level: i32,
    ) -> Result<f64, CryptoError> {
        let fee_row = sqlx::query!(
            r#"SELECT nf.withdrawal_fee, nf.withdrawal_fee_percent, bl.withdrawal_fee_discount_percent
               FROM network_fees nf
               LEFT JOIN brand_levels bl ON bl.level = $3 AND bl.is_active = true
               WHERE nf.asset_id = $1 AND nf.network_id = $2"#,
            asset_id,
            network_id,
            user_level
        )
        .fetch_optional(&*self.pool)
        .await?;

        let base_fee = fee_row.map(|r| r.withdrawal_fee.unwrap_or(0.0)).unwrap_or(0.0);
        let fee_percent = fee_row
            .and_then(|r| r.withdrawal_fee_percent)
            .unwrap_or(0.0);
        let discount = fee_row
            .and_then(|r| r.withdrawal_fee_discount_percent)
            .unwrap_or(0.0);

        let adjusted_percent = fee_percent * (1.0 - discount / 100.0);
        let total_fee = base_fee + (amount * adjusted_percent / 100.0);

        Ok((total_fee * 100_000_000.0).round() / 100_000_000.0)
    }

    // Get deposit address for user
    pub async fn get_deposit_address(
        &self,
        _user_id: Uuid,
        asset_id: Uuid,
        network_id: Uuid,
    ) -> Result<String, CryptoError> {
        // Check if deposit is enabled
        let asset_network = sqlx::query!(
            r#"SELECT deposit_enabled FROM crypto_asset_networks
               WHERE asset_id = $1 AND network_id = $2 AND deposit_enabled = true"#,
            asset_id,
            network_id
        )
        .fetch_optional(&*self.pool)
        .await?;

        if asset_network.is_none() {
            return Err(CryptoError::InvalidOperation(
                "Deposit not enabled for this asset on this network".to_string()
            ));
        }

        // Get admin wallet
        let wallet = sqlx::query!(
            r#"SELECT address FROM admin_wallet_addresses
               WHERE asset_id = $1 AND network_id = $2 AND is_active = true"#,
            asset_id,
            network_id
        )
        .fetch_optional(&*self.pool)
        .await?;

        match wallet {
            Some(w) => Ok(w.address),
            None => Err(CryptoError::WalletNotConfigured),
        }
    }

    // Create deposit record
    pub async fn create_deposit(
        &self,
        user_id: Uuid,
        asset_id: Uuid,
        network_id: Uuid,
        amount: f64,
        address: String,
        tx_hash: Option<String>,
    ) -> Result<CryptoDeposit, CryptoError> {
        // Get user's VIP level
        let user = sqlx::query!(
            "SELECT vip_level FROM users WHERE id = $1",
            user_id
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|_| CryptoError::NotFound("User not found".to_string()))?;

        let vip_level = user.vip_level;

        // Calculate fee
        let fee = self.calculate_deposit_fee(asset_id, network_id, amount, vip_level).await?;
        let net_amount = amount - fee;

        let id = Uuid::new_v4();
        let now = Utc::now();

        sqlx::query!(
            r#"INSERT INTO crypto_deposits
               (id, user_id, asset_id, network_id, amount, fee, net_amount, address, tx_hash, status, created_at)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending', $10)"#,
            id, user_id, asset_id, network_id, amount, fee, net_amount, address, tx_hash, now
        )
        .execute(&*self.pool)
        .await?;

        Ok(CryptoDeposit {
            id,
            user_id,
            asset_id,
            network_id,
            amount,
            fee,
            net_amount,
            address,
            tx_hash,
            confirmations: 0,
            status: "pending".to_string(),
            created_at: now,
            processed_at: None,
            asset: None,
            network: None,
        })
    }

    // Confirm deposit (admin)
    pub async fn confirm_deposit(&self, deposit_id: Uuid) -> Result<(), CryptoError> {
        let deposit = sqlx::query!(
            r#"SELECT id, user_id, net_amount, status FROM crypto_deposits WHERE id = $1"#,
            deposit_id
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|_| CryptoError::NotFound("Deposit not found".to_string()))?;

        if deposit.status != "pending" {
            return Err(CryptoError::InvalidOperation("Deposit already processed".to_string()));
        }

        let now = Utc::now();

        // Start transaction
        let mut tx = self.pool.begin().await?;

        // Update deposit status
        sqlx::query!(
            "UPDATE crypto_deposits SET status = 'completed', processed_at = $1 WHERE id = $2",
            now, deposit_id
        )
        .execute(&mut *tx)
        .await?;

        // Add funds to user balance
        sqlx::query!(
            "UPDATE users SET balance = balance + $1 WHERE id = $2",
            deposit.net_amount, deposit.user_id
        )
        .execute(&mut *tx)
        .await?;

        // Create transaction record
        let tx_id = Uuid::new_v4();
        sqlx::query!(
            r#"INSERT INTO transactions (id, user_id, type, amount, currency, status, tx_hash, fee, processed_at)
               VALUES ($1, $2, 'deposit', $3, $4, 'completed', NULL, $5, $6)"#,
            tx_id, deposit.user_id, deposit.net_amount, deposit_id.to_string(), deposit.net_amount, now
        )
        .execute(&mut *tx)
        .await?;

        tx.commit().await?;

        Ok(())
    }

    // Create withdrawal
    pub async fn create_withdrawal(
        &self,
        user_id: Uuid,
        asset_id: Uuid,
        network_id: Uuid,
        amount: f64,
        address: String,
    ) -> Result<CryptoWithdrawal, CryptoError> {
        // Check if withdrawal is enabled
        let asset_network = sqlx::query!(
            r#"SELECT withdrawal_enabled FROM crypto_asset_networks
               WHERE asset_id = $1 AND network_id = $2 AND withdrawal_enabled = true"#,
            asset_id,
            network_id
        )
        .fetch_optional(&*self.pool)
        .await?;

        if asset_network.is_none() {
            return Err(CryptoError::InvalidOperation(
                "Withdrawal not enabled for this asset on this network".to_string()
            ));
        }

        // Get asset for min withdrawal amount
        let asset = self.get_asset_by_id(asset_id).await?;

        if amount < asset.min_withdrawal_amount {
            return Err(CryptoError::InvalidOperation(
                format!("Amount below minimum withdrawal limit: {}", asset.min_withdrawal_amount)
            ));
        }

        // Get user's VIP level and balance
        let user = sqlx::query!(
            "SELECT vip_level, balance FROM users WHERE id = $1",
            user_id
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|_| CryptoError::NotFound("User not found".to_string()))?;

        // Calculate fee
        let fee = self.calculate_withdrawal_fee(asset_id, network_id, amount, user.vip_level).await?;
        let net_amount = amount - fee;

        if user.balance < amount {
            return Err(CryptoError::InsufficientBalance);
        }

        let id = Uuid::new_v4();
        let now = Utc::now();

        // Start transaction
        let mut tx = self.pool.begin().await?;

        // Deduct from balance
        sqlx::query!(
            "UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1",
            amount, user_id
        )
        .execute(&mut *tx)
        .await?;

        // Create withdrawal record
        sqlx::query!(
            r#"INSERT INTO crypto_withdrawals
               (id, user_id, asset_id, network_id, amount, fee, net_amount, address, status, created_at)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9)"#,
            id, user_id, asset_id, network_id, amount, fee, net_amount, address, now
        )
        .execute(&mut *tx)
        .await?;

        // Create transaction record
        let tx_id = Uuid::new_v4();
        sqlx::query!(
            r#"INSERT INTO transactions (id, user_id, type, amount, currency, status, address, fee)
               VALUES ($1, $2, 'withdrawal', $3, $4, 'pending', $5, $6)"#,
            tx_id, user_id, amount, asset_id.to_string(), address, fee
        )
        .execute(&mut *tx)
        .await?;

        tx.commit().await?;

        Ok(CryptoWithdrawal {
            id,
            user_id,
            asset_id,
            network_id,
            amount,
            fee,
            net_amount,
            address,
            tx_hash: None,
            status: "pending".to_string(),
            created_at: now,
            processed_at: None,
            asset: None,
            network: None,
        })
    }

    // Confirm withdrawal (admin)
    pub async fn confirm_withdrawal(&self, withdrawal_id: Uuid, tx_hash: String) -> Result<(), CryptoError> {
        let withdrawal = sqlx::query!(
            r#"SELECT id, user_id, status FROM crypto_withdrawals WHERE id = $1"#,
            withdrawal_id
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|_| CryptoError::NotFound("Withdrawal not found".to_string()))?;

        if withdrawal.status != "pending" {
            return Err(CryptoError::InvalidOperation("Withdrawal already processed".to_string()));
        }

        let now = Utc::now();

        // Update withdrawal status
        sqlx::query!(
            "UPDATE crypto_withdrawals SET status = 'completed', tx_hash = $1, processed_at = $2 WHERE id = $3",
            tx_hash, now, withdrawal_id
        )
        .execute(&*self.pool)
        .await?;

        // Update transaction record
        sqlx::query!(
            "UPDATE transactions SET status = 'completed', tx_hash = $1, processed_at = $2 WHERE user_id = $3 AND type = 'withdrawal' AND status = 'pending'",
            tx_hash, now, withdrawal.user_id
        )
        .execute(&*self.pool)
        .await?;

        Ok(())
    }

    // Cancel withdrawal (admin)
    pub async fn cancel_withdrawal(&self, withdrawal_id: Uuid) -> Result<(), CryptoError> {
        let withdrawal = sqlx::query!(
            r#"SELECT id, user_id, amount, status FROM crypto_withdrawals WHERE id = $1 AND status = 'pending'"#,
            withdrawal_id
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|_| CryptoError::NotFound("Withdrawal not found or not pending".to_string()))?;

        let now = Utc::now();

        // Start transaction
        let mut tx = self.pool.begin().await?;

        // Refund balance
        sqlx::query!(
            "UPDATE users SET balance = balance + $1 WHERE id = $2",
            withdrawal.amount, withdrawal.user_id
        )
        .execute(&mut *tx)
        .await?;

        // Update withdrawal status
        sqlx::query!(
            "UPDATE crypto_withdrawals SET status = 'cancelled', processed_at = $1 WHERE id = $2",
            now, withdrawal_id
        )
        .execute(&mut *tx)
        .await?;

        // Update transaction record
        sqlx::query!(
            "UPDATE transactions SET status = 'cancelled' WHERE user_id = $1 AND type = 'withdrawal' AND status = 'pending'",
            withdrawal.user_id
        )
        .execute(&mut *tx)
        .await?;

        tx.commit().await?;

        Ok(())
    }

    // Get user deposits
    pub async fn get_user_deposits(&self, user_id: Uuid) -> Result<Vec<CryptoDeposit>, CryptoError> {
        let rows = sqlx::query_as!(
            CryptoDeposit,
            r#"SELECT cd.id, cd.user_id, cd.asset_id, cd.network_id, cd.amount, cd.fee,
               cd.net_amount, cd.address, cd.tx_hash, cd.confirmations, cd.status,
               cd.created_at, cd.processed_at
               FROM crypto_deposits cd
               WHERE cd.user_id = $1
               ORDER BY cd.created_at DESC"#,
            user_id
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Get user withdrawals
    pub async fn get_user_withdrawals(&self, user_id: Uuid) -> Result<Vec<CryptoWithdrawal>, CryptoError> {
        let rows = sqlx::query_as!(
            CryptoWithdrawal,
            r#"SELECT cw.id, cw.user_id, cw.asset_id, cw.network_id, cw.amount, cw.fee,
               cw.net_amount, cw.address, cw.tx_hash, cw.status,
               cw.created_at, cw.processed_at
               FROM crypto_withdrawals cw
               WHERE cw.user_id = $1
               ORDER BY cw.created_at DESC"#,
            user_id
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Get all deposits (admin)
    pub async fn get_all_deposits(
        &self,
        status: Option<String>,
        limit: i32,
        offset: i32,
    ) -> Result<(Vec<CryptoDeposit>, i64), CryptoError> {
        let mut query = String::from(
            r#"SELECT cd.id, cd.user_id, cd.asset_id, cd.network_id, cd.amount, cd.fee,
               cd.net_amount, cd.address, cd.tx_hash, cd.confirmations, cd.status,
               cd.created_at, cd.processed_at
               FROM crypto_deposits cd"#
        );

        if status.is_some() {
            query.push_str(&format!(" WHERE cd.status = '{}'", status.as_ref().unwrap()));
        }
        query.push_str(" ORDER BY cd.created_at DESC LIMIT $1 OFFSET $2");

        let rows: Vec<CryptoDeposit> = if let Some(ref s) = status {
            sqlx::query_as(&query)
                .bind(limit)
                .bind(offset)
                .fetch_all(&*self.pool)
                .await?
        } else {
            sqlx::query_as(&query)
                .bind(limit)
                .bind(offset)
                .fetch_all(&*self.pool)
                .await?
        };

        let count_query = if status.is_some() {
            format!("SELECT COUNT(*) FROM crypto_deposits WHERE status = '{}'", status.as_ref().unwrap())
        } else {
            "SELECT COUNT(*) FROM crypto_deposits".to_string()
        };

        let count: (i64,) = sqlx::query_as(&count_query)
            .fetch_one(&*self.pool)
            .await?;

        Ok((rows, count.0))
    }

    // Get all withdrawals (admin)
    pub async fn get_all_withdrawals(
        &self,
        status: Option<String>,
        limit: i32,
        offset: i32,
    ) -> Result<(Vec<CryptoWithdrawal>, i64), CryptoError> {
        let mut query = String::from(
            r#"SELECT cw.id, cw.user_id, cw.asset_id, cw.network_id, cw.amount, cw.fee,
               cw.net_amount, cw.address, cw.tx_hash, cw.status,
               cw.created_at, cw.processed_at
               FROM crypto_withdrawals cw"#
        );

        if status.is_some() {
            query.push_str(&format!(" WHERE cw.status = '{}'", status.as_ref().unwrap()));
        }
        query.push_str(" ORDER BY cw.created_at DESC LIMIT $1 OFFSET $2");

        let rows: Vec<CryptoWithdrawal> = sqlx::query_as(&query)
            .bind(limit)
            .bind(offset)
            .fetch_all(&*self.pool)
            .await?;

        let count_query = if status.is_some() {
            format!("SELECT COUNT(*) FROM crypto_withdrawals WHERE status = '{}'", status.as_ref().unwrap())
        } else {
            "SELECT COUNT(*) FROM crypto_withdrawals".to_string()
        };

        let count: (i64,) = sqlx::query_as(&count_query)
            .fetch_one(&*self.pool)
            .await?;

        Ok((rows, count.0))
    }

    // Get all brand levels
    pub async fn get_all_brand_levels(&self) -> Result<Vec<BrandLevel>, CryptoError> {
        let rows = sqlx::query_as!(
            BrandLevel,
            r#"SELECT id, name, level, deposit_fee_discount_percent,
               withdrawal_fee_discount_percent, is_active, created_at, updated_at
               FROM brand_levels WHERE is_active = true ORDER BY level"#
        )
        .fetch_all(&*self.pool)
        .await?;

        Ok(rows)
    }

    // Admin: Create network
    pub async fn create_network(
        &self,
        name: String,
        chain_id: Option<String>,
        symbol: String,
        explorer_url: Option<String>,
        min_confirmation_blocks: i32,
    ) -> Result<CryptoNetwork, CryptoError> {
        let id = Uuid::new_v4();
        let now = Utc::now();

        sqlx::query!(
            r#"INSERT INTO crypto_networks
               (id, name, chain_id, symbol, explorer_url, is_withdrawal_enabled, is_deposit_enabled, is_active, min_confirmation_blocks, created_at, updated_at)
               VALUES ($1, $2, $3, $4, $5, true, true, true, $6, $7, $7)"#,
            id, name, chain_id, symbol, explorer_url, min_confirmation_blocks, now
        )
        .execute(&*self.pool)
        .await?;

        Ok(CryptoNetwork {
            id,
            name,
            chain_id,
            symbol,
            explorer_url,
            rpc_url: None,
            is_withdrawal_enabled: true,
            is_deposit_enabled: true,
            is_active: true,
            min_confirmation_blocks,
            created_at: now,
            updated_at: now,
        })
    }

    // Admin: Update network
    pub async fn update_network(
        &self,
        id: Uuid,
        name: String,
        chain_id: Option<String>,
        symbol: String,
        explorer_url: Option<String>,
        is_active: bool,
        is_deposit_enabled: bool,
        is_withdrawal_enabled: bool,
        min_confirmation_blocks: i32,
    ) -> Result<(), CryptoError> {
        let now = Utc::now();

        sqlx::query!(
            r#"UPDATE crypto_networks SET
               name = $1, chain_id = $2, symbol = $3, explorer_url = $4,
               is_active = $5, is_deposit_enabled = $6, is_withdrawal_enabled = $7,
               min_confirmation_blocks = $8, updated_at = $9
               WHERE id = $10"#,
            name, chain_id, symbol, explorer_url, is_active, is_deposit_enabled,
            is_withdrawal_enabled, min_confirmation_blocks, now, id
        )
        .execute(&*self.pool)
        .await?;

        Ok(())
    }

    // Admin: Create asset
    pub async fn create_asset(
        &self,
        name: String,
        symbol: String,
        decimals: i32,
        min_deposit_amount: f64,
        min_withdrawal_amount: f64,
    ) -> Result<CryptoAsset, CryptoError> {
        let id = Uuid::new_v4();
        let now = Utc::now();

        sqlx::query!(
            r#"INSERT INTO crypto_assets
               (id, name, symbol, decimals, is_active, min_deposit_amount, min_withdrawal_amount, created_at, updated_at)
               VALUES ($1, $2, $3, $4, true, $5, $6, $7, $7)"#,
            id, name, symbol, decimals, min_deposit_amount, min_withdrawal_amount, now
        )
        .execute(&*self.pool)
        .await?;

        Ok(CryptoAsset {
            id,
            name,
            symbol,
            decimals,
            contract_address: None,
            is_active: true,
            min_deposit_amount,
            min_withdrawal_amount,
            created_at: now,
            updated_at: now,
        })
    }

    // Admin: Update asset
    pub async fn update_asset(
        &self,
        id: Uuid,
        name: String,
        decimals: i32,
        min_deposit_amount: f64,
        min_withdrawal_amount: f64,
        is_active: bool,
    ) -> Result<(), CryptoError> {
        let now = Utc::now();

        sqlx::query!(
            r#"UPDATE crypto_assets SET
               name = $1, decimals = $2, min_deposit_amount = $3,
               min_withdrawal_amount = $4, is_active = $5, updated_at = $6
               WHERE id = $7"#,
            name, decimals, min_deposit_amount, min_withdrawal_amount, is_active, now, id
        )
        .execute(&*self.pool)
        .await?;

        Ok(())
    }

    // Admin: Link asset to network
    pub async fn link_asset_to_network(
        &self,
        asset_id: Uuid,
        network_id: Uuid,
        deposit_enabled: bool,
        withdrawal_enabled: bool,
        contract_address: Option<String>,
    ) -> Result<CryptoAssetNetwork, CryptoError> {
        let id = Uuid::new_v4();
        let now = Utc::now();

        sqlx::query!(
            r#"INSERT INTO crypto_asset_networks
               (id, asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address, created_at, updated_at)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $7)"#,
            id, asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address, now
        )
        .execute(&*self.pool)
        .await?;

        Ok(CryptoAssetNetwork {
            id,
            asset_id,
            network_id,
            deposit_enabled,
            withdrawal_enabled,
            contract_address,
            created_at: now,
            updated_at: now,
            asset: None,
            network: None,
        })
    }

    // Admin: Set network fee
    pub async fn set_network_fee(
        &self,
        asset_id: Uuid,
        network_id: Uuid,
        deposit_fee: f64,
        withdrawal_fee: f64,
        deposit_fee_percent: f64,
        withdrawal_fee_percent: f64,
    ) -> Result<NetworkFee, CryptoError> {
        let now = Utc::now();

        // Check if fee exists
        let existing = sqlx::query!(
            "SELECT id FROM network_fees WHERE asset_id = $1 AND network_id = $2",
            asset_id, network_id
        )
        .fetch_optional(&*self.pool)
        .await?;

        let id = if let Some(e) = existing {
            sqlx::query!(
                r#"UPDATE network_fees SET
                   deposit_fee = $1, withdrawal_fee = $2, deposit_fee_percent = $3,
                   withdrawal_fee_percent = $4, updated_at = $5
                   WHERE id = $6"#,
                deposit_fee, withdrawal_fee, deposit_fee_percent, withdrawal_fee_percent, now, e.id
            )
            .execute(&*self.pool)
            .await?;
            e.id
        } else {
            let new_id = Uuid::new_v4();
            sqlx::query!(
                r#"INSERT INTO network_fees
                   (id, asset_id, network_id, deposit_fee, withdrawal_fee, deposit_fee_percent, withdrawal_fee_percent, created_at, updated_at)
                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)"#,
                new_id, asset_id, network_id, deposit_fee, withdrawal_fee, deposit_fee_percent, withdrawal_fee_percent, now
            )
            .execute(&*self.pool)
            .await?;
            new_id
        };

        Ok(NetworkFee {
            id,
            asset_id,
            network_id,
            deposit_fee,
            withdrawal_fee,
            deposit_fee_percent,
            withdrawal_fee_percent,
            created_at: now,
            updated_at: now,
        })
    }

    // Admin: Update brand level
    pub async fn update_brand_level(
        &self,
        id: Uuid,
        name: String,
        level: i32,
        deposit_fee_discount_percent: f64,
        withdrawal_fee_discount_percent: f64,
        is_active: bool,
    ) -> Result<(), CryptoError> {
        let now = Utc::now();

        sqlx::query!(
            r#"UPDATE brand_levels SET
               name = $1, level = $2, deposit_fee_discount_percent = $3,
               withdrawal_fee_discount_percent = $4, is_active = $5, updated_at = $6
               WHERE id = $7"#,
            name, level, deposit_fee_discount_percent, withdrawal_fee_discount_percent, is_active, now, id
        )
        .execute(&*self.pool)
        .await?;

        Ok(())
    }

    // Admin: Set admin wallet
    pub async fn set_admin_wallet(
        &self,
        asset_id: Uuid,
        network_id: Uuid,
        address: String,
        private_key_encrypted: Option<String>,
        is_active: bool,
    ) -> Result<AdminWalletAddress, CryptoError> {
        let now = Utc::now();

        // Check if wallet exists
        let existing = sqlx::query!(
            "SELECT id FROM admin_wallet_addresses WHERE asset_id = $1 AND network_id = $2",
            asset_id, network_id
        )
        .fetch_optional(&*self.pool)
        .await?;

        let id = if let Some(e) = existing {
            sqlx::query!(
                r#"UPDATE admin_wallet_addresses SET
                   address = $1, private_key_encrypted = $2, is_active = $3, updated_at = $4
                   WHERE id = $5"#,
                address, private_key_encrypted, is_active, now, e.id
            )
            .execute(&*self.pool)
            .await?;
            e.id
        } else {
            let new_id = Uuid::new_v4();
            sqlx::query!(
                r#"INSERT INTO admin_wallet_addresses
                   (id, asset_id, network_id, address, private_key_encrypted, is_active, created_at, updated_at)
                   VALUES ($1, $2, $3, $4, $5, $6, $7, $7)"#,
                new_id, asset_id, network_id, address, private_key_encrypted, is_active, now
            )
            .execute(&*self.pool)
            .await?;
            new_id
        };

        Ok(AdminWalletAddress {
            id,
            asset_id,
            network_id,
            address,
            private_key_encrypted,
            is_active,
            created_at: now,
            updated_at: now,
        })
    }

    // Admin: Update user brand level
    pub async fn update_user_brand_level(&self, user_id: Uuid, level: i32) -> Result<(), CryptoError> {
        sqlx::query!(
            "UPDATE users SET vip_level = $1 WHERE id = $2",
            level, user_id
        )
        .execute(&*self.pool)
        .await?;

        Ok(())
    }
}
