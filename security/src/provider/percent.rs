//! Percent Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PercentProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PercentProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "percent_001".to_string(), name: "Starburst".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.09, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/starburst/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_002".to_string(), name: "Gonzo's Quest".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 95.97, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/gonzo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_003".to_string(), name: "Dead or Alive".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.80, volatility: Volatility::High, min_bet: 0.09, max_bet: 18.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/dead-alive/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_004".to_string(), name: "Twin Spin".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.60, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/twin-spin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_005".to_string(), name: "Jack and the Beanstalk".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.70, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/beanstalk/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_006".to_string(), name: "Blood Suckers".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 98.00, volatility: Volatility::Low, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/blood-suckers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_007".to_string(), name: "Flowers".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/flowers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_008".to_string(), name: "South Park".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.70, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/south-park/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_009".to_string(), name: "Blood Suckers 2".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Low, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/blood-suckers-2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "percent_010".to_string(), name: "Joker Pro".to_string(), provider: "Percent".to_string(), category: GameCategory::Slots, rtp: 96.80, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.percent.com/joker-pro/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PercentProvider {
    fn name(&self) -> &str { "Percent" }
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
