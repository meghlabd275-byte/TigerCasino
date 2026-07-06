//! Big Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BigGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BigGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "biggaming_001".to_string(), name: "Big Gaming Slots 1".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_002".to_string(), name: "Big Gaming Slots 2".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_003".to_string(), name: "Big Gaming Slots 3".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_004".to_string(), name: "Big Gaming Slots 4".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot4/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_005".to_string(), name: "Big Gaming Slots 5".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot5/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_006".to_string(), name: "Big Gaming Slots 6".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot6/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_007".to_string(), name: "Big Gaming Slots 7".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot7/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_008".to_string(), name: "Big Gaming Slots 8".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot8/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_009".to_string(), name: "Big Gaming Slots 9".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot9/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "biggaming_010".to_string(), name: "Big Gaming Slots 10".to_string(), provider: "Big Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.biggaming.com/slot10/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BigGamingProvider {
    fn name(&self) -> &str { "Big Gaming" }
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
