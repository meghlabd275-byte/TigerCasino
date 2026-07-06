//! Iron Dog Studio Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct IronDogProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl IronDogProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "irondog_001".to_string(), name: "Vikings".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/vikings/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_002".to_string(), name: "Raiding Raiders".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/raiding-raiders/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_003".to_string(), name: "Foxy".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/foxy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_004".to_string(), name: "Gold".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_005".to_string(), name: "Creature".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/creature/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_006".to_string(), name: "Witches".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/witches/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_007".to_string(), name: "Jack".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/jack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_008".to_string(), name: "Dracula".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/dracula/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_009".to_string(), name: "Pirates".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/pirates/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "irondog_010".to_string(), name: "Wolf".to_string(), provider: "Iron Dog Studio".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.irondogstudio.com/wolf/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for IronDogProvider {
    fn name(&self) -> &str { "Iron Dog Studio" }
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
