//! Spinmatic Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SpinmaticProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SpinmaticProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "spinmatic_001".to_string(), name: "Book of Magic".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/book-magic/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_002".to_string(), name: "Fruit Fantasy".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/fruit-fantasy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_003".to_string(), name: "Joker's Luck".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/jokers-luck/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_004".to_string(), name: "Lucky Dragon".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/lucky-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_005".to_string(), name: "Egyptian Dreams".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/egyptian-dreams/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_006".to_string(), name: "Gangsters".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/gangsters/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_007".to_string(), name: "La Granaventura".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/granaventura/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_008".to_string(), name: "Bamboo Tower".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/bamboo-tower/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_009".to_string(), name: "888 Dragons".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/888-dragons/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spinmatic_010".to_string(), name: "Soccer Mania".to_string(), provider: "Spinmatic".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.spinmatic.com/soccer-mania/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SpinmaticProvider {
    fn name(&self) -> &str { "Spinmatic" }
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
