//! NetEnt Live Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct NetEntLiveProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl NetEntLiveProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "netentlive_bj_001".to_string(), name: "NetEnt Blackjack".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.netent.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netentlive_bj_002".to_string(), name: "NetEnt Speed Blackjack".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.netent.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "netentlive_r_001".to_string(), name: "NetEnt Roulette".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.netent.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netentlive_r_002".to_string(), name: "NetEnt Auto Roulette".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.netent.com/auto-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "netentlive_b_001".to_string(), name: "NetEnt Baccarat".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.netent.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            // Game Shows
            GameInfo { id: "netentlive_gs_001".to_string(), name: "NetEnt Live Dream Catcher".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.netent.com/dream-catcher/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netentlive_gs_002".to_string(), name: "NetEnt Live Monopoly".to_string(), provider: "NetEnt Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.netent.com/monopoly/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slots
            GameInfo { id: "netentlive_slots_001".to_string(), name: "Starburst".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.09, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/starburst/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netentlive_slots_002".to_string(), name: "Gonzo's Quest".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 95.97, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/gonzo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netentlive_slots_003".to_string(), name: "Divine Fortune".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/divine-fortune/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for NetEntLiveProvider {
    fn name(&self) -> &str { "NetEnt Live" }
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
