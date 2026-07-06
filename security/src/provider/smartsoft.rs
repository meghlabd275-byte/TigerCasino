//! SmartSoft Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SmartSoftProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SmartSoftProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "smartsoft_001".to_string(), name: "JetX".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/jetx/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_002".to_string(), name: "Balloon".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/balloon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_003".to_string(), name: "Keno".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/keno/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_004".to_string(), name: "Plinko".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Crash, rtp: 98.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/plinko/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_005".to_string(), name: "Mines".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/mines/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_006".to_string(), name: "Dice".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_007".to_string(), name: "Roulette".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Roulette, rtp: 97.30, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_008".to_string(), name: "Baccarat".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Baccarat, rtp: 98.94, volatility: Volatility::Low, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_009".to_string(), name: "Blackjack".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Blackjack, rtp: 99.50, volatility: Volatility::Low, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "smartsoft_010".to_string(), name: "Poker".to_string(), provider: "SmartSoft Gaming".to_string(), category: GameCategory::Poker, rtp: 97.80, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.smartsoftgaming.com/poker/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SmartSoftProvider {
    fn name(&self) -> &str { "SmartSoft Gaming" }
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
