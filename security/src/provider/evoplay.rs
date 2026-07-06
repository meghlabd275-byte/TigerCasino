//! EvoPlay Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct EvoPlayProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EvoPlayProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "evoplay_001".to_string(), name: "The Great icescape".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/icecape/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_002".to_string(), name: "Chicago".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/chicago/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_003".to_string(), name: "Mine Field".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/mine-field/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_004".to_string(), name: "Indiana's Quest".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/indiana/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_005".to_string(), name: "Naughty Girls".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/naughty-girls/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_006".to_string(), name: "The Ming Dynasty".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/ming-dynasty/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_007".to_string(), name: "Star Cart".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/star-cart/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_008".to_string(), name: "Sea Harbor".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/sea-harbor/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_009".to_string(), name: "Irish Reels".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/irish-reels/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_010".to_string(), name: "Frogues".to_string(), provider: "EvoPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/frogues/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for EvoPlayProvider {
    fn name(&self) -> &str { "EvoPlay" }
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
