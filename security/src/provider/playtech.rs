//! Playtech Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PlaytechProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PlaytechProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "playtech_001".to_string(), name: "Age of Gods".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.02, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/age-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_002".to_string(), name: "Gladiator".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 94.46, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/gladiator/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_003".to_string(), name: "Rocky".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/rocky/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_004".to_string(), name: "Pink Panther".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 94.40, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/pink-panther/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_005".to_string(), name: "Great Blue".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.52, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/great-blue/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_006".to_string(), name: "King Kong".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/king-kong/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_007".to_string(), name: "Avengers".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/avengers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_008".to_string(), name: "Iron Man".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/iron-man/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_009".to_string(), name: "X-Men".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/x-men/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playtech_010".to_string(), name: "Batman".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/batman/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PlaytechProvider {
    fn name(&self) -> &str { "Playtech" }
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
