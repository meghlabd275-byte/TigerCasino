//! BetConstruct Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BetConstructProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BetConstructProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Casino
            GameInfo { id: "bc_bj_001".to_string(), name: "BetConstruct Blackjack".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_r_001".to_string(), name: "BetConstruct Roulette".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_b_001".to_string(), name: "BetConstruct Baccarat".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            // Virtual Sports
            GameInfo { id: "bc_vs_001".to_string(), name: "Virtual Football".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/v-football/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_vs_002".to_string(), name: "Virtual Basketball".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/v-basketball/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_vs_003".to_string(), name: "Virtual Tennis".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/v-tennis/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_vs_004".to_string(), name: "Virtual Horse Racing".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/v-horses/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_vs_005".to_string(), name: "Virtual Greyhound Racing".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betconstruct.com/v-greyhounds/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slots
            GameInfo { id: "bc_slots_001".to_string(), name: "Fruits and Jokers".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betconstruct.com/fruits-jokers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bc_slots_002".to_string(), name: "Hot 27".to_string(), provider: "BetConstruct".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betconstruct.com/hot-27/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BetConstructProvider {
    fn name(&self) -> &str { "BetConstruct" }
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
