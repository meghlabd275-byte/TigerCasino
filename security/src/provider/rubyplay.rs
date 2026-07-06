//! RubyPlay Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct RubyPlayProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl RubyPlayProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "ruby_001".to_string(), name: "Ruby Champions".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/champions/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_002".to_string(), name: "Ruby Hit".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/ruby-hit/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_003".to_string(), name: "Ruby Fortune".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_004".to_string(), name: "Ruby Magic".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/magic/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_005".to_string(), name: "Ruby Win".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/win/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_006".to_string(), name: "Ruby Riches".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/riches/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_007".to_string(), name: "Ruby Star".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/star/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_008".to_string(), name: "Ruby Galaxy".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/galaxy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_009".to_string(), name: "Ruby Crown".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/crown/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ruby_010".to_string(), name: "Ruby Quest".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/quest/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for RubyPlayProvider {
    fn name(&self) -> &str { "RubyPlay" }
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
