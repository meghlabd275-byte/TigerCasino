//! Wazdan Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct WazdanProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl WazdanProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "wazdan_001".to_string(), name: "9 Lions".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/9-lions/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_002".to_string(), name: "Magic Stars 3".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: false, thumbnail_url: "https://static.wazdan.com/magic-stars-3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_003".to_string(), name: "Sizzling 777".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.wazdan.com/sizzling-777/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_004".to_string(), name: "Hot Slot".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/hot-slot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_005".to_string(), name: "Larry the Leprechaun".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/larry-leprechaun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_006".to_string(), name: "Power of Gods".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/power-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_007".to_string(), name: "Jelly Boom".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/jelly-boom/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_008".to_string(), name: "Crazy Cars".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/crazy-cars/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_009".to_string(), name: "Turn Your Luck".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/turn-luck/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "wazdan_010".to_string(), name: "Super Hot".to_string(), provider: "Wazdan".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.wazdan.com/super-hot/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for WazdanProvider {
    fn name(&self) -> &str { "Wazdan" }
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
