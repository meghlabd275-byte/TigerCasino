//! Additional Live Casino Games Provider

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct LiveCasinoExtraProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl LiveCasinoExtraProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // More Evolution Live Blackjack
            GameInfo { id: "live_extra_bj_101".to_string(), name: "Blackjack 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_102".to_string(), name: "Blackjack 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_103".to_string(), name: "Blackjack 3".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_104".to_string(), name: "Blackjack 4".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj4/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_105".to_string(), name: "Blackjack 5".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj5/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_106".to_string(), name: "Blackjack 6".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj6/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_107".to_string(), name: "Blackjack 7".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj7/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_108".to_string(), name: "Blackjack 8".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj8/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_109".to_string(), name: "Blackjack 9".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj9/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_110".to_string(), name: "Blackjack 10".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bj10/thumb.jpg".to_string(), game_url: "".to_string() },
            // More Roulette
            GameInfo { id: "live_extra_r_101".to_string(), name: "Roulette 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/r1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_102".to_string(), name: "Roulette 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/r2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_103".to_string(), name: "Roulette 3".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/r3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_104".to_string(), name: "Roulette 4".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/r4/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_105".to_string(), name: "Roulette 5".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/r5/thumb.jpg".to_string(), game_url: "".to_string() },
            // More Baccarat
            GameInfo { id: "live_extra_bac_101".to_string(), name: "Baccarat 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bac1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bac_102".to_string(), name: "Baccarat 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bac2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bac_103".to_string(), name: "Baccarat 3".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bac3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bac_104".to_string(), name: "Baccarat 4".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bac4/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bac_105".to_string(), name: "Baccarat 5".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/bac5/thumb.jpg".to_string(), game_url: "".to_string() },
            // Teen Patti & Andar Bahar
            GameInfo { id: "live_extra_tp_001".to_string(), name: "Teen Patti 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/tp1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_tp_002".to_string(), name: "Teen Patti 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/tp2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_ab_001".to_string(), name: "Andar Bahar 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/ab1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_ab_002".to_string(), name: "Andar Bahar 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/ab2/thumb.jpg".to_string(), game_url: "".to_string() },
            // Sic Bo
            GameInfo { id: "live_extra_sb_001".to_string(), name: "Sic Bo 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.22, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/sb1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_sb_002".to_string(), name: "Sic Bo 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.22, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/sb2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_sb_003".to_string(), name: "Sic Bo 3".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.22, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/sb3/thumb.jpg".to_string(), game_url: "".to_string() },
            // Dragon Tiger
            GameInfo { id: "live_extra_dt_001".to_string(), name: "Dragon Tiger 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/dt1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_dt_002".to_string(), name: "Dragon Tiger 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/dt2/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for LiveCasinoExtraProvider {
    fn name(&self) -> &str { "Live Casino Extra" }
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
