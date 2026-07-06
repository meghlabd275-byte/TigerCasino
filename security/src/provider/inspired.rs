//! Inspired Provider Integration (Virtual Sports, Slots)

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct InspiredProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl InspiredProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Virtual Sports
            GameInfo { id: "inspired_vs_001".to_string(), name: "Virtual Football".to_string(), provider: "Inspired".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.inspired.com/v-football/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_vs_002".to_string(), name: "Virtual Basketball".to_string(), provider: "Inspired".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.inspired.com/v-basketball/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_vs_003".to_string(), name: "Virtual Tennis".to_string(), provider: "Inspired".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.inspired.com/v-tennis/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_vs_004".to_string(), name: "Virtual Horse Racing".to_string(), provider: "Inspired".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.inspired.com/v-horses/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_vs_005".to_string(), name: "Virtual Greyhound Racing".to_string(), provider: "Inspired".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.inspired.com/v-greyhounds/thumb.jpg".to_string(), game_url: "".to_string() },
            // Slots
            GameInfo { id: "inspired_slots_001".to_string(), name: "Rainbow Riches".to_string(), provider: "Inspired".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.inspired.com/rainbow-riches/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_slots_002".to_string(), name: "Monopoly".to_string(), provider: "Inspired".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.inspired.com/monopoly/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_slots_003".to_string(), name: "Deal or No Deal".to_string(), provider: "Inspired".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.inspired.com/deal-or-no-deal/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_slots_004".to_string(), name: "Who Wants to be a Millionaire".to_string(), provider: "Inspired".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.inspired.com/millionaire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "inspired_slots_005".to_string(), name: "Reel King".to_string(), provider: "Inspired".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.inspired.com/reel-king/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for InspiredProvider {
    fn name(&self) -> &str { "Inspired" }
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
