//! High 5 Games Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct High5GamesProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl High5GamesProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "high5_001".to_string(), name: "Da Vinci Diamonds".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 94.90, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/davinci/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_002".to_string(), name: "Cleopatra".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/cleopatra/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_003".to_string(), name: "Siberian Storm".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/siberian/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_004".to_string(), name: "Cats".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.33, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/cats/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_005".to_string(), name: "Double Diamond".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.44, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/double-diamond/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_006".to_string(), name: "Wolf Run".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 94.98, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/wolf-run/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_007".to_string(), name: "Golden Goddess".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 94.99, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/golden-goddess/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_008".to_string(), name: "Red Phoenix".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/red-phoenix/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_009".to_string(), name: "Candy Bars".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.17, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/candy-bars/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "high5_010".to_string(), name: "Bird of Prey".to_string(), provider: "High 5 Games".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.high5games.com/bird-prey/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for High5GamesProvider {
    fn name(&self) -> &str { "High 5 Games" }
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
