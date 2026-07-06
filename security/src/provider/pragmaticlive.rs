//! Pragmatic Play Live Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PragmaticLiveProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PragmaticLiveProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "pplive_bj_001".to_string(), name: "Pragmatic Blackjack".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.pragmaticplay.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_bj_002".to_string(), name: "Pragmatic Speed Blackjack".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.pragmaticplay.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "pplive_r_001".to_string(), name: "Pragmatic Roulette".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.pragmaticplay.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "pplive_b_001".to_string(), name: "Pragmatic Baccarat".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.pragmaticplay.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            // Game Shows - Pragmatic's strength
            GameInfo { id: "pplive_gs_001".to_string(), name: "Sweet Bonanza CandyLand".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/sweet-bonanza-candyland/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_gs_002".to_string(), name: "PowerUP Roulette".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/powerup-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_gs_003".to_string(), name: "Mega Wheel".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/mega-wheel/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_gs_004".to_string(), name: "Money Drop".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/money-drop/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_gs_005".to_string(), name: "Boom City".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/boom-city/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "pplive_gs_006".to_string(), name: "Gonzo's Treasure Hunt".to_string(), provider: "Pragmatic Play Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.com/gonzo-treasure-hunt/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PragmaticLiveProvider {
    fn name(&self) -> &str { "Pragmatic Play Live" }
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
