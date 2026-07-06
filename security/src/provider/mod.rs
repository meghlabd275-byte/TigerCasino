//! Game Provider API Integrations
//! 
//! Integrations with major game providers like Pragmatic Play, Play'n GO, NetEnt, 
//! Evolution Gaming, BGaming, Spribe, and more.

use serde::{Deserialize, Serialize};
use thiserror::Error;

pub mod pragmatic;
pub mod evolution;
pub mod playngo;
pub mod netent;
pub mod bgaming;
pub mod spribe;

/// Common errors for provider integrations
#[derive(Error, Debug)]
pub enum ProviderError {
    #[error("API request failed: {0}")]
    ApiError(String),
    
    #[error("Authentication failed: {0}")]
    AuthError(String),
    
    #[error("Game not found: {0}")]
    GameNotFound(String),
    
    #[error("Invalid response from provider")]
    InvalidResponse,
    
    #[error("Network error: {0}")]
    NetworkError(String),
    
    #[error("Rate limited: {0}")]
    RateLimited(String),
    
    #[error("Signature verification failed")]
    SignatureError,
}

/// Common game information structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GameInfo {
    pub id: String,
    pub name: String,
    pub provider: String,
    pub category: GameCategory,
    pub rtp: f64,
    pub volatility: Volatility,
    pub min_bet: f64,
    pub max_bet: f64,
    pub has_free_spins: bool,
    pub has_bonus_game: bool,
    pub thumbnail_url: String,
    pub game_url: String,
}

/// Game categories
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum GameCategory {
    Slots,
    LiveCasino,
    GameShows,
    TableGames,
    Crash,
    Arcade,
    VirtualSports,
    Lottery,
    Bingo,
    ScratchCards,
    Poker,
    Blackjack,
    Roulette,
    Baccarat,
}

/// Volatility levels
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum Volatility {
    Low,
    Medium,
    High,
}

/// Game launch request
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LaunchGameRequest {
    pub user_id: String,
    pub game_id: String,
    pub provider: String,
    pub currency: String,
    pub language: String,
    pub stake: f64,
    pub session_token: String,
}

/// Game launch response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LaunchGameResponse {
    pub game_url: String,
    pub session_id: String,
    pub token: String,
    pub expires_at: i64,
}

/// Provider configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProviderConfig {
    pub name: String,
    pub api_url: String,
    pub secret_key: String,
    pub public_key: String,
    pub enabled: bool,
}

/// Transaction result for game plays
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TransactionResult {
    pub transaction_id: String,
    pub status: TransactionStatus,
    pub amount: f64,
    pub balance_after: f64,
    pub game_round_id: String,
    pub timestamp: i64,
}

/// Transaction status
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum TransactionStatus {
    Pending,
    Completed,
    Failed,
    Cancelled,
}

/// Provider statistics
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProviderStats {
    pub provider: String,
    pub total_games: u32,
    pub active_games: u32,
    pub total_rtp: f64,
    pub avg_rtp: f64,
    pub total_bets: u64,
    pub total_wins: f64,
}

/// Trait for provider integrations
pub trait GameProvider: Send + Sync {
    /// Get provider name
    fn name(&self) -> &str;
    
    /// Get available games from provider
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError>;
    
    /// Launch a specific game
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError>;
    
    /// Process game transaction (bet/win)
    fn process_transaction(&self, request: TransactionRequest) -> Result<TransactionResult, ProviderError>;
    
    /// Get game info
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError>;
    
    /// Check if provider is available
    fn is_available(&self) -> bool;
}

/// Transaction request
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TransactionRequest {
    pub session_id: String,
    pub game_id: String,
    pub user_id: String,
    pub transaction_type: TransactionType,
    pub amount: f64,
    pub round_id: String,
    pub timestamp: i64,
}

/// Transaction types
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum TransactionType {
    Bet,
    Win,
    Refund,
    Bonus,
}

/// Provider manager for handling multiple providers
pub struct ProviderManager {
    providers: Vec<Box<dyn GameProvider>>,
}

impl ProviderManager {
    pub fn new() -> Self {
        Self {
            providers: Vec::new(),
        }
    }
    
    pub fn register_provider(&mut self, provider: Box<dyn GameProvider>) {
        self.providers.push(provider);
    }
    
    pub fn get_provider(&self, name: &str) -> Option<&Box<dyn GameProvider>> {
        self.providers.iter().find(|p| p.name() == name)
    }
    
    pub fn get_all_games(&self) -> Vec<GameInfo> {
        let mut games = Vec::new();
        for provider in &self.providers {
            if let Ok(provider_games) = provider.get_games() {
                games.extend(provider_games);
            }
        }
        games
    }
    
    pub fn get_stats(&self) -> Vec<ProviderStats> {
        self.providers.iter().map(|p| {
            let games = p.get_games().unwrap_or_default();
            ProviderStats {
                provider: p.name().to_string(),
                total_games: games.len() as u32,
                active_games: games.len() as u32,
                total_rtp: games.iter().map(|g| g.rtp).sum(),
                avg_rtp: if !games.is_empty() { 
                    games.iter().map(|g| g.rtp).sum::<f64>() / games.len() as f64 
                } else { 0.0 },
                total_bets: 0,
                total_wins: 0.0,
            }
        }).collect()
    }
}

impl Default for ProviderManager {
    fn default() -> Self {
        Self::new()
    }
}
