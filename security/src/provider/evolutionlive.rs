//! Evolution Live Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct EvolutionLiveProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EvolutionLiveProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Blackjack - Evolution is the best
            GameInfo { id: "evolive_bj_001".to_string(), name: "Evolution Blackjack".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolutiongaming.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_bj_002".to_string(), name: "Evolution Speed Blackjack".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolutiongaming.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_bj_003".to_string(), name: "Evolution Infinite Blackjack".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 1.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolutiongaming.com/infinite-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Roulette
            GameInfo { id: "evolive_r_001".to_string(), name: "Evolution Roulette".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolutiongaming.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            // Live Baccarat
            GameInfo { id: "evolive_b_001".to_string(), name: "Evolution Baccarat".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 10000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolutiongaming.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            // Game Shows - Evolution's specialty
            GameInfo { id: "evolive_gs_001".to_string(), name: "Evolution Crazy Time".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolutiongaming.com/crazy-time/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_gs_002".to_string(), name: "Evolution Monopoly Live".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolutiongaming.com/monopoly-live/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_gs_003".to_string(), name: "Evolution Lightning Roulette".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolutiongaming.com/lightning-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_gs_004".to_string(), name: "Evolution Dream Catcher".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolutiongaming.com/dream-catcher/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evolive_gs_005".to_string(), name: "Evolution Deal or No Deal".to_string(), provider: "Evolution Live".to_string(), category: GameCategory::GameShows, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolutiongaming.com/deal-or-no-deal/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for EvolutionLiveProvider {
    fn name(&self) -> &str { "Evolution Live" }
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
