//! Habanero Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct HabaneroProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl HabaneroProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "habanero_001".to_string(), name: "Hot Hot Fruit".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/hot-fruit/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_002".to_string(), name: "Koi Gate".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/koi-gate/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_003".to_string(), name: "Fa Cai Shen".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/fa-cai-shen/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_004".to_string(), name: "4 Lucky Stars".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/4-lucky-stars/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_005".to_string(), name: "Jellyfish Flow".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/jellyfish-flow/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_006".to_string(), name: "Mighty Tips".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/mighty-tips/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_007".to_string(), name: "Nuwa".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/nuwa/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_008".to_string(), name: "Wealth Inn".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/wealth-inn/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_009".to_string(), name: "Christmas Gift".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/christmas-gift/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "habanero_010".to_string(), name: "Lucky Fortune".to_string(), provider: "Habanero".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.habanero.com/lucky-fortune/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for HabaneroProvider {
    fn name(&self) -> &str { "Habanero" }
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
