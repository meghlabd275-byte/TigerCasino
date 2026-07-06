//! Amusnet Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct AmusnetProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl AmusnetProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "amusnet_001".to_string(), name: "Sugar Rush".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/sugar-rush/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_002".to_string(), name: "Fruits and Jokers".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/fruits-jokers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_003".to_string(), name: "Burning Hot".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/burning-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_004".to_string(), name: "40 Super Hot".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/40-super-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_005".to_string(), name: "Ultimate Hot".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/ultimate-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_006".to_string(), name: "Shining Crown".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/shining-crown/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_007".to_string(), name: "More Like a Cash".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/more-like-cash/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_008".to_string(), name: "Supreme Hot".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/supreme-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_009".to_string(), name: "Burning Hot Deluxe".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/burning-hot-deluxe/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "amusnet_010".to_string(), name: "Flaming Hot".to_string(), provider: "Amusnet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.amusnet.com/flaming-hot/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for AmusnetProvider {
    fn name(&self) -> &str { "Amusnet" }
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
