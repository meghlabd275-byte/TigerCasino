//! Hollywood TV Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct HollywoodTVProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl HollywoodTVProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "htv_bj_001".to_string(), name: "Hollywood Blackjack".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_bj_002".to_string(), name: "Hollywood Speed Blackjack".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_r_001".to_string(), name: "Hollywood Roulette".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_r_002".to_string(), name: "Hollywood Auto Roulette".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/auto-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_b_001".to_string(), name: "Hollywood Baccarat".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_b_002".to_string(), name: "Hollywood Dragon Tiger".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 5.0, max_bet: 2500.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/dragon-tiger/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_p_001".to_string(), name: "Hollywood Casino Hold'em".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.80, volatility: Volatility::Medium, min_bet: 5.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/holdem/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_tp_001".to_string(), name: "Hollywood Teen Patti".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/teen-patti/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_ab_001".to_string(), name: "Hollywood Andar Bahar".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/andar-bahar/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "htv_war_001".to_string(), name: "Hollywood War".to_string(), provider: "Hollywood TV".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 1.0, max_bet: 5000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.hollywoodtv.com/war/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for HollywoodTVProvider {
    fn name(&self) -> &str { "Hollywood TV" }
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
