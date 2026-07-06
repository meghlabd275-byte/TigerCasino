//! Red Tiger Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct RedTigerProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl RedTigerProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "redtiger_001".to_string(), name: "Gonzo's Quest Megaways".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/gonzo-mega/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_002".to_string(), name: "Lightning Horseman".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/horseman/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_003".to_string(), name: "Treasure Mine".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/treasure-mine/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_004".to_string(), name: "Spartans".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/spartans/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_005".to_string(), name: "Snow Wild".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/snow-wild/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_006".to_string(), name: "Reel King".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/reel-king/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_007".to_string(), name: "Vikings".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/vikings/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_008".to_string(), name: "Primate King".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/primate-king/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_009".to_string(), name: "Thor".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/thor/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "redtiger_010".to_string(), name: "Dragon's Luck".to_string(), provider: "Red Tiger".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.redtiger.com/dragons-luck/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for RedTigerProvider {
    fn name(&self) -> &str { "Red Tiger" }
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> { self.fetch_games() }
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        Ok(LaunchGameResponse { game_url: format!("{}/game/{}", self.base_url, request.game_id), session_id: uuid::Uuid::new_v4().to_string(), token: "session_token".to_string(), expires_at: Utc::now().timestamp() + 3600 })
    }
    fn process_transaction(&self, _request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        Ok(TransactionResult { transaction_id: uuid::Uuid::new_v4().to_string(), status: TransactionStatus::Completed, amount: 0.0, balance_after: 0.0, game_round_id: "round_id".to_string(), timestamp: Utc::now().timestamp() })
    }
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.fetch_games()?;
        games.into_iter().find(|g| g.id == game_id).ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    fn is_available(&self) -> bool { self.config.enabled }
}
