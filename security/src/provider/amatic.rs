//! Amatic Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct AmaticProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl AmaticProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "amatic_001".to_string(), name: "Admiral Nelson".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/nelson/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_002".to_string(), name: "Book of Fruits".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/book-fruits/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_003".to_string(), name: "Cool Wolf".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/cool-wolf/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_004".to_string(), name: "Golden Night".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/golden-night/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_005".to_string(), name: "Lucky Star".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/lucky-star/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_006".to_string(), name: "Mighty Dragon".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/mighty-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_007".to_string(), name: "Redrooster".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/redrooster/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_008".to_string(), name: "Royal Unicorn".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/royal-unicorn/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_009".to_string(), name: "Star Lanterns".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/star-lanterns/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amatic_010".to_string(), name: "Wolf Moon".to_string(), provider: "Amatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amatic.com/wolf-moon/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for AmaticProvider {
    fn name(&self) -> &str { "Amatic" }
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
