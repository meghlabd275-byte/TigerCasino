//! MrSlotty Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct MrSlottyProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl MrSlottyProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "mrslotty_001".to_string(), name: "Gods of Olympus".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/olympus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_002".to_string(), name: "The Money".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/the-money/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_003".to_string(), name: "Super Hot".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/super-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_004".to_string(), name: "Vegas".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/vegas/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_005".to_string(), name: "Miami".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/miami/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_006".to_string(), name: "Easter".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/easter/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_007".to_string(), name: "Alice".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/alice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_008".to_string(), name: "Monaco".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/monaco/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_009".to_string(), name: "Rockstar".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/rockstar/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "mrslotty_010".to_string(), name: "Neon".to_string(), provider: "MrSlotty".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.mrslotty.com/neon/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for MrSlottyProvider {
    fn name(&self) -> &str { "MrSlotty" }
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
