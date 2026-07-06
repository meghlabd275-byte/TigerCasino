//! Vivo Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct VivoGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl VivoGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack
            GameInfo { id: "vivo_bj_001".to_string(), name: "Vivo Blackjack".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_bj_002".to_string(), name: "Vivo Blackjack 7".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/bj7/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_bj_003".to_string(), name: "Vivo Blackjack VIP".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 25.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/bjvip/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "vivo_r_001".to_string(), name: "Vivo Roulette".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_r_002".to_string(), name: "Vivo Auto Roulette".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/auto-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_r_003".to_string(), name: "Vivo Speed Roulette".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/speed-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "vivo_b_001".to_string(), name: "Vivo Baccarat".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_b_002".to_string(), name: "Vivo Dragon Tiger".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/dragon-tiger/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Poker
            GameInfo { id: "vivo_p_001".to_string(), name: "Vivo Casino Hold'em".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.80, volatility: Volatility::Medium, min_bet: 5.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/holdem/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_p_002".to_string(), name: "Vivo Caribbean Stud".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 96.30, volatility: Volatility::Medium, min_bet: 10.0, max_bet: 500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivogaming.com/caribbean/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for VivoGamingProvider {
    fn name(&self) -> &str { "Vivo Gaming" }
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
