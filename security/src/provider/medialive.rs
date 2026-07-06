//! Medialive Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct MedialiveProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl MedialiveProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "media_bj_001".to_string(), name: "Medialive Blackjack".to_string(), provider: "Medialive".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.medialive.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "media_r_001".to_string(), name: "Medialive Roulette".to_string(), provider: "Medialive".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.medialive.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "media_b_001".to_string(), name: "Medialive Baccarat".to_string(), provider: "Medialive".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.medialive.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "media_tp_001".to_string(), name: "Medialive Teen Patti".to_string(), provider: "Medialive".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.medialive.com/teen-patti/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for MedialiveProvider {
    fn name(&self) -> &str { "Medialive" }
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
