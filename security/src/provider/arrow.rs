//! Arrow's Edge Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct ArrowsEdgeProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl ArrowsEdgeProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "ae_001".to_string(), name: "Slots Journey".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/journey/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_002".to_string(), name: "Vikings Plunder".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/vikings/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_003".to_string(), name: "Gold Rush".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/gold-rush/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_004".to_string(), name: "Treasure Nile".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/treasure-nile/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_005".to_string(), name: "Reels of Fortune".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/reels-fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_006".to_string(), name: "Mythic Wolf".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/mythic-wolf/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_007".to_string(), name: "Lucky Leprechaun".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/leprechaun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_008".to_string(), name: "Mighty Medusa".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/medusa/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_009".to_string(), name: "Dance of the Dead".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/dance-dead/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ae_010".to_string(), name: "Forgotten".to_string(), provider: "Arrow's Edge".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.arrowedge.com/forgotten/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for ArrowsEdgeProvider {
    fn name(&self) -> &str { "Arrow's Edge" }
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
