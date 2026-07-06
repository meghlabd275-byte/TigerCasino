//! Pragmatic Play Provider Integration
//! 
//! Integration with Pragmatic Play's game API.
//! Pragmatic Play offers slots, live casino, bingo, and virtual sports.

use super::*;
use reqwest::Client;
use serde_json::json;

/// Pragmatic Play API client
pub struct PragmaticPlayProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PragmaticPlayProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            client: Client::new(),
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    /// Generate authentication signature
    fn generate_signature(&self, payload: &str, timestamp: i64) -> String {
        use hmac::{Hmac, Mac};
        use sha2::Sha256;
        
        type HmacSha256 = Hmac<Sha256>;
        
        let mut mac = HmacSha256::new_from_slice(self.config.secret_key.as_bytes())
            .expect("HMAC can take key of any size");
        mac.update(payload.as_bytes());
        mac.update(timestamp.to_string().as_bytes());
        
        hex::encode(mac.finalize().into_bytes())
    }
    
    /// Get available games from Pragmatic Play
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // In production, this would call the actual Pragmatic Play API
        // For now, return simulated game list based on their portfolio
        Ok(vec![
            GameInfo {
                id: "pragmatic_slots_001".to_string(),
                name: "Sweet Bonanza".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Slots,
                rtp: 96.48,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/sweet-bonanza/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_slots_002".to_string(),
                name: "Wolf Gold".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::Medium,
                min_bet: 0.25,
                max_bet: 125.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/wolf-gold/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_slots_003".to_string(),
                name: "The Dog House".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Slots,
                rtp: 96.51,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/dog-house/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_slots_004".to_string(),
                name: "Gates of Olympus".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/gates-olympus/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_slots_005".to_string(),
                name: "Starlight Princess".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/starlight-princess/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_live_001".to_string(),
                name: "Mega Wheel".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::GameShows,
                rtp: 96.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://static.pragmaticplay.net/mega-wheel/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_live_002".to_string(),
                name: "Speed Blackjack".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 99.0,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 10000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.pragmaticplay.net/speed-blackjack/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_crash_001".to_string(),
                name: "Spaceman".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Crash,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.pragmaticplay.net/spaceman/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_bingo_001".to_string(),
                name: "Bingo Blast".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::Bingo,
                rtp: 95.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.pragmaticplay.net/bingo-blast/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "pragmatic_virtual_001".to_string(),
                name: "Virtual Football".to_string(),
                provider: "Pragmatic Play".to_string(),
                category: GameCategory::VirtualSports,
                rtp: 95.0,
                volatility: Volatility::Medium,
                min_bet: 0.50,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.pragmaticplay.net/virtual-football/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for PragmaticPlayProvider {
    fn name(&self) -> &str {
        "Pragmatic Play"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_games()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        let timestamp = chrono::Utc::now().timestamp();
        let payload = json!({
            "user_id": request.user_id,
            "game_id": request.game_id,
            "currency": request.currency,
            "language": request.language,
            "stake": request.stake,
        }).to_string();
        
        let signature = self.generate_signature(&payload, timestamp);
        
        // In production, make API call to Pragmatic Play
        // For now, return mock response
        Ok(LaunchGameResponse {
            game_url: format!("{}/game/{}", self.base_url, request.game_id),
            session_id: uuid::Uuid::new_v4().to_string(),
            token: signature,
            expires_at: timestamp + 3600,
        })
    }
    
    fn process_transaction(&self, request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        // Process bet/win transaction with provider
        let timestamp = chrono::Utc::now().timestamp();
        let payload = json!({
            "session_id": request.session_id,
            "game_id": request.game_id,
            "user_id": request.user_id,
            "transaction_type": request.transaction_type,
            "amount": request.amount,
            "round_id": request.round_id,
            "timestamp": timestamp,
        }).to_string();
        
        let signature = self.generate_signature(&payload, timestamp);
        
        // In production, make API call
        Ok(TransactionResult {
            transaction_id: uuid::Uuid::new_v4().to_string(),
            status: TransactionStatus::Completed,
            amount: request.amount,
            balance_after: 0.0, // Would be returned from provider
            game_round_id: request.round_id,
            timestamp,
        })
    }
    
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.get_games()?;
        games.into_iter()
            .find(|g| g.id == game_id)
            .ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    
    fn is_available(&self) -> bool {
        self.config.enabled
    }
}
