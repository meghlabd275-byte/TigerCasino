//! Authentic Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct AuthenticGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl AuthenticGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Roulette from real casinos
            GameInfo { id: "auth_r_001".to_string(), name: "Authentic Roulette".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_r_002".to_string(), name: "Auto Roulette".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/auto-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_r_003".to_string(), name: "Roulette 1".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/roulette-1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_r_004".to_string(), name: "Roulette 2".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/roulette-2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_r_005".to_string(), name: "Arabic Roulette".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/arabic-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_r_006".to_string(), name: "Turkish Roulette".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/turkish-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Blackjack
            GameInfo { id: "auth_bj_001".to_string(), name: "Authentic Blackjack".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_bj_002".to_string(), name: "Blackjack 1".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/blackjack-1/thumb.jpg".to_string(), game_url: "".to_string() },
            // Baccarat
            GameInfo { id: "auth_b_001".to_string(), name: "Authentic Baccarat".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "auth_b_002".to_string(), name: "Baccarat 1".to_string(), provider: "Authentic Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.authenticgaming.com/baccarat-1/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for AuthenticGamingProvider {
    fn name(&self) -> &str { "Authentic Gaming" }
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
