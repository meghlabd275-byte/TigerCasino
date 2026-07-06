//! Tom Horn Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct TomHornProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl TomHornProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "tomhorn_001".to_string(), name: "Wild Fire".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/wild-fire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_002".to_string(), name: "Hot Slot".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/hot-slot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_003".to_string(), name: "Wolf Sierra".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/wolf-sierra/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_004".to_string(), name: "Crystal Garden".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/crystal-garden/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_005".to_string(), name: "Dynasty".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/dynasty/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_006".to_string(), name: " Shaolin".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/shaolin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_007".to_string(), name: "Almighty".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/almighty/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_008".to_string(), name: "Vegas".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/vegas/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_009".to_string(), name: "Enchanted".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/enchanted/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "tomhorn_010".to_string(), name: "The Master".to_string(), provider: "Tom Horn".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.01, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.tomhorn.com/the-master/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for TomHornProvider {
    fn name(&self) -> &str { "Tom Horn" }
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
