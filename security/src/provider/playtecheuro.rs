//! Playtech Live Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PlaytechLiveProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PlaytechLiveProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "ptlive_bj_001".to_string(), name: "Playtech Blackjack".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playtech.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ptlive_bj_002".to_string(), name: "Playtech VIP Blackjack".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 25.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playtech.com/bj-vip/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "ptlive_r_001".to_string(), name: "Playtech Roulette".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playtech.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ptlive_r_002".to_string(), name: "Playtech Speed Roulette".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playtech.com/speed-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "ptlive_b_001".to_string(), name: "Playtech Baccarat".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.playtech.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            // Game Shows
            GameInfo { id: "ptlive_gs_001".to_string(), name: "Playtech Crazy Time".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/crazy-time/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ptlive_gs_002".to_string(), name: "Playtech Adventures".to_string(), provider: "Playtech Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/adventures/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slots from Playtech
            GameInfo { id: "ptlive_slots_001".to_string(), name: "Age of Gods".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 95.02, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/age-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ptlive_slots_002".to_string(), name: "Jackpot Giant".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 94.22, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/jackpot-giant/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ptlive_slots_003".to_string(), name: "Gladiator".to_string(), provider: "Playtech".to_string(), category: GameCategory::Slots, rtp: 94.08, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playtech.com/gladiator/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PlaytechLiveProvider {
    fn name(&self) -> &str { "Playtech Live" }
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
