//! KA Gaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct KAGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl KAGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "ka_001".to_string(), name: "Dragon Ball".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/dragon-ball/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_002".to_string(), name: "Legend of Dragon".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/legend-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_003".to_string(), name: "Golden Basket".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/golden-basket/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_004".to_string(), name: "Money Tree".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/money-tree/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_005".to_string(), name: "Three Kingdoms".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/three-kingdoms/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_006".to_string(), name: "Jewels".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/jewels/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_007".to_string(), name: "Fortune Panda".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/fortune-panda/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_008".to_string(), name: "Golden Egg".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/golden-egg/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_009".to_string(), name: "Aladdin".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/aladdin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_010".to_string(), name: "Mulan".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/mulan/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for KAGamingProvider {
    fn name(&self) -> &str { "KA Gaming" }
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
