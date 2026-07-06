//! EBet Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct EBetProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EBetProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "ebet_bj_001".to_string(), name: "EBet Blackjack".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ebet_bj_002".to_string(), name: "EBet Speed Blackjack".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "ebet_r_001".to_string(), name: "EBet Roulette".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ebet_r_002".to_string(), name: "EBet Auto Roulette".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/auto-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "ebet_b_001".to_string(), name: "EBet Baccarat".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ebet_b_002".to_string(), name: "EBet Dragon Tiger".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/dragon-tiger/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Sic Bo
            GameInfo { id: "ebet_sb_001".to_string(), name: "EBet Sic Bo".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.22, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/sicbo/thumb.jpg".to_string(), game_url: "".to_string() },
            // Teen Patti
            GameInfo { id: "ebet_tp_001".to_string(), name: "EBet Teen Patti".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/teen-patti/thumb.jpg".to_string(), game_url: "".to_string() },
            // Andar Bahar
            GameInfo { id: "ebet_ab_001".to_string(), name: "EBet Andar Bahar".to_string(), provider: "EBet".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.ebet.com/andar-bahar/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slot Games
            GameInfo { id: "ebet_slots_001".to_string(), name: "EBet Fortune".to_string(), provider: "EBet".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ebet.com/fortune/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for EBetProvider {
    fn name(&self) -> &str { "EBet" }
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
