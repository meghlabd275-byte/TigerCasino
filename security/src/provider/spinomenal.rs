//! Spinomenal Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SpinomenalProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SpinomenalProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "spinomenal_001".to_string(), name: "Book of Wolves".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/book-wolves/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_002".to_string(), name: "Wild Lanterns".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/wild-lanterns/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_003".to_string(), name: "Egyptian Gods".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/egyptian-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_004".to_string(), name: "Demons and Gold".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/demons-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_005".to_string(), name: "Story of Zeus".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/zeus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_006".to_string(), name: "Wolf Fang".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/wolf-fang/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_007".to_string(), name: "Lord of the Ocean".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/lord-ocean/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_008".to_string(), name: "Queen of the Sun".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/queen-sun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_009".to_string(), name: "Pirates Power".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/pirates-power/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinomenal_010".to_string(), name: "Juicy Pop".to_string(), provider: "Spinomenal".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinomenal.com/juicy-pop/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SpinomenalProvider {
    fn name(&self) -> &str { "Spinomenal" }
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
