//! Booongo Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BooongoProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BooongoProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "booongo_001".to_string(), name: "Lord of the Sun".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/lord-sun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_002".to_string(), name: "Easter Fortune".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/easter-fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_003".to_string(), name: "Star Captain".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/star-captain/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_004".to_string(), name: "Gold of Sirens".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/gold-sirens/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_005".to_string(), name: "Fruits and SEAs".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/fruits-seas/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_006".to_string(), name: "Night Club".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/night-club/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_007".to_string(), name: "Temple of Treasures".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/temple-treasure/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_008".to_string(), name: "Scattered G".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/scattered-g/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_009".to_string(), name: "Freya".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/freya/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_010".to_string(), name: "Fairy Land".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/fairy-land/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BooongoProvider {
    fn name(&self) -> &str { "Booongo" }
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
