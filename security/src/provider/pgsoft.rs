//! PGSoft Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PGSoftProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PGSoftProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "pgsoft_001".to_string(), name: "Mahjong Ways".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.92, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/mahjong-ways/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_002".to_string(), name: "Mahjong Ways 2".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.92, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/mahjong-ways-2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_003".to_string(), name: "Fortune Ox".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/fortune-ox/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_004".to_string(), name: "Fortune Mouse".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/fortune-mouse/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_005".to_string(), name: "Prosperity Fortune Tree".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/prosperity-tree/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_006".to_string(), name: "Dragon Hatch".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/dragon-hatch/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_007".to_string(), name: "Genie's Wish".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/genies-wish/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_008".to_string(), name: "Medusa II".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/medusa-2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_009".to_string(), name: "Jungle Delight".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/jungle-delight/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pgsoft_010".to_string(), name: "Treasures of Aztec".to_string(), provider: "PGSoft".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pgsoft.com/treasures-aztec/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PGSoftProvider {
    fn name(&self) -> &str { "PGSoft" }
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
