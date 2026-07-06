//! Dragon Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct DragonGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl DragonGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "dragon_001".to_string(), name: "Dragon's Fortune".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_002".to_string(), name: "Golden Dragon".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/golden-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_003".to_string(), name: "Dragon Phoenix".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/phoenix/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_004".to_string(), name: "Dragon's Realm".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/realm/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_005".to_string(), name: "Dragon's Gold".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_006".to_string(), name: "Dragon's Treasure".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/treasure/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_007".to_string(), name: "Dragon's Lair".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/lair/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_008".to_string(), name: "Dragon's Eye".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/eye/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_009".to_string(), name: "Dragon's Fire".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/fire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "dragon_010".to_string(), name: "Dragon's Champion".to_string(), provider: "Dragon Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.dragongaming.com/champion/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for DragonGamingProvider {
    fn name(&self) -> &str { "Dragon Gaming" }
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
