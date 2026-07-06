use std::sync::Arc;
use sqlx::postgres::{PgPool, PgPoolOptions};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum DatabaseError {
    #[error("Connection error: {0}")]
    ConnectionError(String),
    #[error("Query error: {0}")]
    QueryError(String),
    #[error("Pool error: {0}")]
    PoolError(String),
}

pub struct Database {
    pool: Arc<PgPool>,
}

impl Database {
    pub async fn new(database_url: &str) -> Result<Self, DatabaseError> {
        let pool = PgPoolOptions::new()
            .max_connections(10)
            .min_connections(2)
            .acquire_timeout(std::time::Duration::from_secs(30))
            .connect(database_url)
            .await
            .map_err(|e| DatabaseError::ConnectionError(e.to_string()))?;

        Ok(Self {
            pool: Arc::new(pool),
        })
    }

    pub fn pool(&self) -> Arc<PgPool> {
        Arc::clone(&self.pool)
    }

    pub async fn close(&self) {
        self.pool.close().await;
    }
}

// User model for authentication
#[derive(sqlx::FromRow, Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct User {
    pub id: uuid::Uuid,
    pub email: String,
    pub username: String,
    pub password_hash: String,
    pub balance: f64,
    pub bonus_balance: f64,
    pub vip_level: i32,
    pub is_admin: bool,
    pub is_banned: bool,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

impl Database {
    // Get user by ID
    pub async fn get_user(&self, user_id: uuid::Uuid) -> Result<Option<User>, DatabaseError> {
        let user = sqlx::query_as!(
            User,
            "SELECT id, email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned, created_at FROM users WHERE id = $1",
            user_id
        )
        .fetch_optional(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(user)
    }

    // Get user by email
    pub async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, DatabaseError> {
        let user = sqlx::query_as!(
            User,
            "SELECT id, email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned, created_at FROM users WHERE email = $1",
            email
        )
        .fetch_optional(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(user)
    }

    // Get user by username
    pub async fn get_user_by_username(&self, username: &str) -> Result<Option<User>, DatabaseError> {
        let user = sqlx::query_as!(
            User,
            "SELECT id, email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned, created_at FROM users WHERE username = $1",
            username
        )
        .fetch_optional(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(user)
    }

    // Create user
    pub async fn create_user(
        &self,
        email: &str,
        username: &str,
        password_hash: &str,
    ) -> Result<User, DatabaseError> {
        let user = sqlx::query_as!(
            User,
            r#"INSERT INTO users (email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned)
               VALUES ($1, $2, $3, 0, 0, 0, false, false)
               RETURNING id, email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned, created_at"#,
            email,
            username,
            password_hash
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(user)
    }

    // Update user balance
    pub async fn update_balance(&self, user_id: uuid::Uuid, amount: f64) -> Result<(), DatabaseError> {
        sqlx::query!(
            "UPDATE users SET balance = balance + $1 WHERE id = $2",
            amount,
            user_id
        )
        .execute(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(())
    }

    // Check if user is admin
    pub async fn is_admin(&self, user_id: uuid::Uuid) -> Result<bool, DatabaseError> {
        let result = sqlx::query!(
            "SELECT is_admin FROM users WHERE id = $1",
            user_id
        )
        .fetch_optional(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(result.map(|r| r.is_admin).unwrap_or(false))
    }

    // Get all users (admin)
    pub async fn get_all_users(&self, limit: i32, offset: i32) -> Result<(Vec<User>, i64), DatabaseError> {
        let users = sqlx::query_as!(
            User,
            "SELECT id, email, username, password_hash, balance, bonus_balance, vip_level, is_admin, is_banned, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
            limit,
            offset
        )
        .fetch_all(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        let count: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM users")
            .fetch_one(&*self.pool)
            .await
            .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok((users, count.0))
    }

    // Ban/Unban user
    pub async fn set_user_banned(&self, user_id: uuid::Uuid, banned: bool) -> Result<(), DatabaseError> {
        sqlx::query!(
            "UPDATE users SET is_banned = $1 WHERE id = $2",
            banned,
            user_id
        )
        .execute(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(())
    }
}

// Transaction model
#[derive(sqlx::FromRow, Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Transaction {
    pub id: uuid::Uuid,
    pub user_id: uuid::Uuid,
    #[(rename = "type")]
    pub transaction_type: String,
    pub amount: f64,
    pub currency: String,
    pub status: String,
    pub tx_hash: Option<String>,
    pub address: Option<String>,
    pub fee: f64,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub processed_at: Option<chrono::DateTime<chrono::Utc>>,
}

impl Database {
    // Get user transactions
    pub async fn get_user_transactions(&self, user_id: uuid::Uuid, limit: i32) -> Result<Vec<Transaction>, DatabaseError> {
        let transactions = sqlx::query_as!(
            Transaction,
            "SELECT id, user_id, type, amount, currency, status, tx_hash, address, fee, created_at, processed_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2",
            user_id,
            limit
        )
        .fetch_all(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(transactions)
    }

    // Create transaction record
    pub async fn create_transaction(
        &self,
        user_id: uuid::Uuid,
        transaction_type: &str,
        amount: f64,
        currency: &str,
        status: &str,
        address: Option<&str>,
        fee: f64,
    ) -> Result<Transaction, DatabaseError> {
        let transaction = sqlx::query_as!(
            Transaction,
            r#"INSERT INTO transactions (user_id, type, amount, currency, status, address, fee)
               VALUES ($1, $2, $3, $4, $5, $6, $7)
               RETURNING id, user_id, type, amount, currency, status, tx_hash, address, fee, created_at, processed_at"#,
            user_id,
            transaction_type,
            amount,
            currency,
            status,
            address,
            fee
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(transaction)
    }
}

// Session model
#[derive(sqlx::FromRow, Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Session {
    pub id: uuid::Uuid,
    pub user_id: uuid::Uuid,
    pub token: String,
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
    pub expires_at: chrono::DateTime<chrono::Utc>,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

impl Database {
    // Create session
    pub async fn create_session(
        &self,
        user_id: uuid::Uuid,
        token: &str,
        expires_at: chrono::DateTime<chrono::Utc>,
        ip_address: Option<&str>,
        user_agent: Option<&str>,
    ) -> Result<Session, DatabaseError> {
        let session = sqlx::query_as!(
            Session,
            r#"INSERT INTO sessions (user_id, token, ip_address, user_agent, expires_at)
               VALUES ($1, $2, $3, $4, $5)
               RETURNING id, user_id, token, ip_address, user_agent, expires_at, created_at"#,
            user_id,
            token,
            ip_address,
            user_agent,
            expires_at
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(session)
    }

    // Get session by token
    pub async fn get_session(&self, token: &str) -> Result<Option<Session>, DatabaseError> {
        let session = sqlx::query_as!(
            Session,
            "SELECT id, user_id, token, ip_address, user_agent, expires_at, created_at FROM sessions WHERE token = $1 AND expires_at > NOW()",
            token
        )
        .fetch_optional(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(session)
    }

    // Delete session
    pub async fn delete_session(&self, token: &str) -> Result<(), DatabaseError> {
        sqlx::query!("DELETE FROM sessions WHERE token = $1", token)
            .execute(&*self.pool)
            .await
            .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(())
    }

    // Delete expired sessions
    pub async fn delete_expired_sessions(&self) -> Result<u64, DatabaseError> {
        let result = sqlx::query!("DELETE FROM sessions WHERE expires_at < NOW()")
            .execute(&*self.pool)
            .await
            .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(result.rows_affected())
    }
}

// Audit log model
#[derive(sqlx::FromRow, Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct AuditLog {
    pub id: uuid::Uuid,
    pub user_id: Option<uuid::Uuid>,
    pub action: String,
    pub details: Option<serde_json::Value>,
    pub ip_address: Option<String>,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

impl Database {
    // Create audit log
    pub async fn create_audit_log(
        &self,
        user_id: Option<uuid::Uuid>,
        action: &str,
        details: Option<serde_json::Value>,
        ip_address: Option<&str>,
    ) -> Result<AuditLog, DatabaseError> {
        let log = sqlx::query_as!(
            AuditLog,
            r#"INSERT INTO audit_logs (user_id, action, details, ip_address)
               VALUES ($1, $2, $3, $4)
               RETURNING id, user_id, action, details, ip_address, created_at"#,
            user_id,
            action,
            details,
            ip_address
        )
        .fetch_one(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(log)
    }

    // Get audit logs for user
    pub async fn get_user_audit_logs(&self, user_id: uuid::Uuid, limit: i32) -> Result<Vec<AuditLog>, DatabaseError> {
        let logs = sqlx::query_as!(
            AuditLog,
            "SELECT id, user_id, action, details, ip_address, created_at FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2",
            user_id,
            limit
        )
        .fetch_all(&*self.pool)
        .await
        .map_err(|e| DatabaseError::QueryError(e.to_string()))?;

        Ok(logs)
    }
}
