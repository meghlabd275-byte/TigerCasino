//! BBIN Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BBINProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BBINProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "bbin_bj_001".to_string(), name: "BBIN Blackjack".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bbin_bj_002".to_string(), name: "BBIN VIP Blackjack".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 25.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/bj-vip/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "bbin_r_001".to_string(), name: "BBIN Roulette".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bbin_r_002".to_string(), name: "BBIN Speed Roulette".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/speed-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "bbin_b_001".to_string(), name: "BBIN Baccarat".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bbin_b_002".to_string(), name: "BBIN Dragon Tiger".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/dragon-tiger/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Sic Bo
            GameInfo { id: "bbin_sb_001".to_string(), name: "BBIN Sic Bo".to_string(), provider: "BBIN".to_string(), category: GameCategory::LiveCasino, rtp: 97.22, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/sicbo/thumb.jpg".to_string(), game_url: "".to_string() },
            // Lottery
            GameInfo { id: "bbin_lt_001".to_string(), name: "BBIN Lottery".to_string(), provider: "BBIN".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bbin.com/lottery/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slot Games
            GameInfo { id: "bbin_slots_001".to_string(), name: "BBIN Fortune".to_string(), provider: "BBIN".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bbin.com/fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bbin_slots_002".to_string(), name: "BBIN Dragon".to_string(), provider: "BBIN".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bbin.com/dragon/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BBINProvider {
    fn name(&self) -> &str { "BBIN" }
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
