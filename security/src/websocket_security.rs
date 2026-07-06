//! WebSocket Security Module for TigerCasino
//!
//! Secure WebSocket handling for real-time gaming

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use chrono::{DateTime, Utc, Duration};
use sha2::{Sha256, Digest};

// ============== Message Types ==============

/// WebSocket message types
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "type", content = "data")]
pub enum WSMessageType {
    // Auth
    Auth { token: String },
    AuthResponse { success: bool, user_id: Option<String>, error: Option<String> },
    
    // Game messages
    GameAction { game_id: String, action: String, params: HashMap<String, serde_json::Value> },
    GameState { game_state: serde_json::Value },
    GameResult { result: serde_json::Value },
    
    // Chat
    Chat { message: String },
    ChatMessage { user_id: String, username: String, message: String, timestamp: i64 },
    
    // System
    Ping,
    Pong,
    Error { code: String, message: String },
    Heartbeat,
}

/// WebSocket message wrapper
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WSMessage {
    pub id: String,
    pub timestamp: i64,
    pub payload: WSMessageType,
}

impl WSMessage {
    pub fn new(payload: WSMessageType) -> Self {
        Self {
            id: uuid::Uuid::new_v4().to_string(),
            timestamp: Utc::now().timestamp(),
            payload,
        }
    }
}

// ============== Connection State ==============

/// Connection state
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ConnectionState {
    Connecting,
    Authenticating,
    Authenticated,
    InGame,
    Disconnecting,
    Disconnected,
}

/// WebSocket connection info
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WSConnection {
    pub connection_id: String,
    pub user_id: Option<String>,
    pub state: ConnectionState,
    pub ip_address: String,
    pub user_agent: String,
    pub connected_at: DateTime<Utc>,
    pub last_activity: DateTime<Utc>,
    pub current_game: Option<String>,
    pub messages_sent: u64,
    pub messages_received: u64,
    pub bytes_sent: u64,
    pub bytes_received: u64,
}

impl WSConnection {
    pub fn new(connection_id: String, ip_address: String, user_agent: String) -> Self {
        Self {
            connection_id,
            user_id: None,
            state: ConnectionState::Connecting,
            ip_address,
            user_agent,
            connected_at: Utc::now(),
            last_activity: Utc::now(),
            current_game: None,
            messages_sent: 0,
            messages_received: 0,
            bytes_sent: 0,
            bytes_received: 0,
        }
    }

    pub fn update_activity(&mut self) {
        self.last_activity = Utc::now();
    }

    pub fn is_expired(&self, timeout: Duration) -> bool {
        Utc::now() - self.last_activity > timeout
    }
}

// ============== Connection Manager ==============

/// WebSocket connection manager
pub struct ConnectionManager {
    connections: Arc<RwLock<HashMap<String, WSConnection>>>,
    user_connections: Arc<RwLock<HashMap<String, Vec<String>>>>,  // user_id -> connection_ids
    max_connections_per_user: usize,
    connection_timeout: Duration,
    message_rate_limit: u32,
}

impl ConnectionManager {
    pub fn new(max_connections_per_user: usize, connection_timeout_seconds: i64, message_rate_limit: u32) -> Self {
        Self {
            connections: Arc::new(RwLock::new(HashMap::new())),
            user_connections: Arc::new(RwLock::new(HashMap::new())),
            max_connections_per_user,
            connection_timeout: Duration::seconds(connection_timeout_seconds),
            message_rate_limit,
        }
    }

    /// Add a new connection
    pub async fn add_connection(&self, connection: WSConnection) -> Result<(), WSError> {
        // Check connection limit per user
        if let Some(ref user_id) = connection.user_id {
            let users = self.user_connections.read().await;
            if let Some(conns) = users.get(user_id) {
                if conns.len() >= self.max_connections_per_user {
                    return Err(WSError::TooManyConnections);
                }
            }
        }
        
        // Add connection
        let mut connections = self.connections.write().await;
        connections.insert(connection.connection_id.clone(), connection.clone());
        
        // Update user connections
        if let Some(user_id) = connection.user_id {
            let mut users = self.user_connections.write().await;
            users.entry(user_id).or_insert_with(Vec::new).push(connection.connection_id.clone());
        }
        
        Ok(())
    }

    /// Remove a connection
    pub async fn remove_connection(&self, connection_id: &str) {
        let mut connections = self.connections.write().await;
        
        if let Some(conn) = connections.remove(connection_id) {
            // Remove from user connections
            if let Some(user_id) = conn.user_id {
                let mut users = self.user_connections.write().await;
                if let Some(conns) = users.get_mut(&user_id) {
                    conns.retain(|id| id != connection_id);
                }
            }
        }
    }

    /// Get connection
    pub async fn get_connection(&self, connection_id: &str) -> Option<WSConnection> {
        let connections = self.connections.read().await;
        connections.get(connection_id).cloned()
    }

    /// Update connection state
    pub async fn update_state(&self, connection_id: &str, state: ConnectionState, user_id: Option<String>) {
        let mut connections = self.connections.write().await;
        if let Some(conn) = connections.get_mut(connection_id) {
            conn.state = state;
            if let Some(uid) = user_id {
                conn.user_id = Some(uid);
            }
            conn.update_activity();
        }
    }

    /// Clean up expired connections
    pub async fn cleanup_expired(&self) -> Vec<String> {
        let mut connections = self.connections.write().await;
        let mut expired = Vec::new();
        
        connections.retain(|id, conn| {
            if conn.is_expired(self.connection_timeout) {
                expired.push(id.clone());
                false
            } else {
                true
            }
        });
        
        expired
    }

    /// Get connection count
    pub async fn connection_count(&self) -> usize {
        let connections = self.connections.read().await;
        connections.len()
    }

    /// Get user connections
    pub async fn get_user_connections(&self, user_id: &str) -> Vec<WSConnection> {
        let connections = self.connections.read().await;
        let users = self.user_connections.read().await;
        
        if let Some(conn_ids) = users.get(user_id) {
            conn_ids.iter()
                .filter_map(|id| connections.get(id).cloned())
                .collect()
        } else {
            Vec::new()
        }
    }
}

// ============== Message Rate Limiter ==============

/// Per-connection message rate limiter
pub struct MessageRateLimiter {
    messages: Vec<DateTime<Utc>>,
    max_messages: u32,
    window_seconds: i64,
}

impl MessageRateLimiter {
    pub fn new(max_messages: u32, window_seconds: i64) -> Self {
        Self {
            messages: Vec::new(),
            max_messages,
            window_seconds,
        }
    }

    /// Check if message is allowed
    pub fn check(&mut self) -> bool {
        let now = Utc::now();
        let window_start = now - Duration::seconds(self.window_seconds);
        
        // Remove old messages
        self.messages.retain(|t| *t > window_start);
        
        // Check rate
        if self.messages.len() >= self.max_messages as usize {
            return false;
        }
        
        // Add current message
        self.messages.push(now);
        true
    }

    /// Reset the limiter
    pub fn reset(&mut self) {
        self.messages.clear();
    }
}

// ============== Security Validator ==============

/// WebSocket security validator
pub struct WSSecurityValidator {
    allowed_origins: Vec<String>,
    allowed_user_agents: Vec<String>,
    blocked_ips: Vec<String>,
    ip_whitelist: Vec<String>,
}

impl WSSecurityValidator {
    pub fn new() -> Self {
        Self {
            allowed_origins: Vec::new(),
            allowed_user_agents: Vec::new(),
            blocked_ips: Vec::new(),
            ip_whitelist: Vec::new(),
        }
    }

    pub fn with_allowed_origins(mut self, origins: Vec<String>) -> Self {
        self.allowed_origins = origins;
        self
    }

    pub fn with_blocked_ips(mut self, ips: Vec<String>) -> Self {
        self.blocked_ips = ips;
        self
    }

    pub fn with_ip_whitelist(mut self, ips: Vec<String>) -> Self {
        self.ip_whitelist = ips;
        self
    }

    /// Validate connection request
    pub fn validate_connection(&self, ip: &str, origin: &str, user_agent: &str) -> Result<(), WSError> {
        // Check blocked IPs
        if self.blocked_ips.contains(&ip.to_string()) {
            return Err(WSError::IPBlocked);
        }
        
        // Check whitelist (if set)
        if !self.ip_whitelist.is_empty() && !self.ip_whitelist.contains(&ip.to_string()) {
            return Err(WSError::IPNotWhitelisted);
        }
        
        // Check origin
        if !self.allowed_origins.is_empty() && !self.allowed_origins.contains(&origin.to_string()) {
            return Err(WSError::InvalidOrigin);
        }
        
        Ok(())
    }

    /// Validate message size
    pub fn validate_message_size(&self, size: usize, max_size: usize) -> Result<(), WSError> {
        if size > max_size {
            Err(WSError::MessageTooLarge)
        } else {
            Ok(())
        }
    }

    /// Generate connection fingerprint
    pub fn generate_fingerprint(ip: &str, user_agent: &str, timestamp: i64) -> String {
        let mut hasher = Sha256::new();
        hasher.update(ip.as_bytes());
        hasher.update(user_agent.as_bytes());
        hasher.update(timestamp.to_string().as_bytes());
        hex::encode(hasher.finalize())
    }
}

impl Default for WSSecurityValidator {
    fn default() -> Self {
        Self::new()
    }
}

// ============== Errors ==============

#[derive(Debug)]
pub enum WSError {
    TooManyConnections,
    IPBlocked,
    IPNotWhitelisted,
    InvalidOrigin,
    MessageTooLarge,
    AuthenticationFailed,
    RateLimited,
    InvalidMessage,
}

impl std::fmt::Display for WSError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            WSError::TooManyConnections => write!(f, "Too many connections"),
            WSError::IPBlocked => write!(f, "IP address blocked"),
            WSError::IPNotWhitelisted => write!(f, "IP not whitelisted"),
            WSError::InvalidOrigin => write!(f, "Invalid origin"),
            WSError::MessageTooLarge => write!(f, "Message too large"),
            WSError::AuthenticationFailed => write!(f, "Authentication failed"),
            WSError::RateLimited => write!(f, "Rate limited"),
            WSError::InvalidMessage => write!(f, "Invalid message"),
        }
    }
}

impl std::error::Error for WSError {}
