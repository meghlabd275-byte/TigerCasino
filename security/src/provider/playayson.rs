//! Playson Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PlaysonProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PlaysonProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "playayson_001".to_string(), name: "Book of Gold".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/book-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_002".to_string(), name: "Legend of Cleopatra".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/cleopatra/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_003".to_string(), name: "Solar Queen".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/solar-queen/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_004".to_string(), name: "Imperial Fruits".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: false, thumbnail_url: "https://static.playayson.com/imperial-fruits/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_005".to_string(), name: "Crystal Land".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/crystal-land/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_006".to_string(), name: "Diamond Star".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: false, thumbnail_url: "https://static.playayson.com/diamond-star/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_007".to_string(), name: "Fruits and Jokers".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playayson.com/fruits-jokers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_008".to_string(), name: "Merge".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/merge/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_009".to_string(), name: "Royal Coins".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/royal-coins/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playayson_010".to_string(), name: "Eternal Cleopatra".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 95.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playayson.com/eternal-cleopatra/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PlaysonProvider {
    fn name(&self) -> &str { "Playson" }
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
