//! iSoftBet Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct ISoftBetProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl ISoftBetProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "isoft_001".to_string(), name: "Hot Shot".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/hot-shot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_002".to_string(), name: "Lucky Wizard".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/lucky-wizard/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_003".to_string(), name: "Rocks".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/rocks/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_004".to_string(), name: "Platoon".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/platoon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_005".to_string(), name: "The Champions".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/champions/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_006".to_string(), name: "Diamond Zone".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/diamond-zone/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_007".to_string(), name: "Great 88".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/great-88/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_008".to_string(), name: "Kung Fu".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/kung-fu/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_009".to_string(), name: "Robyn".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/robyn/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "isoft_010".to_string(), name: "Roxy's".to_string(), provider: "iSoftBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.isoftbet.com/roxys/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for ISoftBetProvider {
    fn name(&self) -> &str { "iSoftBet" }
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
