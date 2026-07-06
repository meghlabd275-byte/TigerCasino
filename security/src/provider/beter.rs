//! Beter Provider Integration (Crash, Dice, Plinko)

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BeterProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BeterProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "beter_001".to_string(), name: "Aviator".to_string(), provider: "Beter".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/aviator/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_002".to_string(), name: "SpeedCrash".to_string(), provider: "Beter".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/speedcrash/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_003".to_string(), name: "JetX".to_string(), provider: "Beter".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/jetx/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_004".to_string(), name: "Plinko".to_string(), provider: "Beter".to_string(), category: GameCategory::Crash, rtp: 98.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/plinko/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_005".to_string(), name: "Mines".to_string(), provider: "Beter".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/mines/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_006".to_string(), name: "Dice".to_string(), provider: "Beter".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_007".to_string(), name: "Dice 2".to_string(), provider: "Beter".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/dice2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_008".to_string(), name: "Keno".to_string(), provider: "Beter".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/keno/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_009".to_string(), name: "Minesweeper".to_string(), provider: "Beter".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/minesweeper/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "beter_010".to_string(), name: "Tower".to_string(), provider: "Beter".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.beter.co/tower/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BeterProvider {
    fn name(&self) -> &str { "Beter" }
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
