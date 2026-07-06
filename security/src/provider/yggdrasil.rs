//! Yggdrasil Provider Integration
//! 
//! Integration with Yggdrasil's game API.

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct YggdrasilProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl YggdrasilProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "ygg_001".to_string(), name: "Valley of the Gods".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/valley-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_002".to_string(), name: "Ozwin's Jackpots".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/ozwin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_003".to_string(), name: "Temple of Tut".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/temple-tut/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_004".to_string(), name: "Jungle Books".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.10, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/jungle-books/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_005".to_string(), name: "Easter Island".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.70, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/easter-island/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_006".to_string(), name: "Penguin City".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.40, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/penguin-city/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_007".to_string(), name: "Hanzo's Dojo".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/hanzo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_008".to_string(), name: "Wolf Hunters".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.30, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/wolf-hunters/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_009".to_string(), name: "Jackpot Express".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/jackpot-express/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ygg_010".to_string(), name: "Hades".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/hades/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for YggdrasilProvider {
    fn name(&self) -> &str { "Yggdrasil" }
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
