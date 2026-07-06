//! Apollo Games Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct ApolloGamesProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl ApolloGamesProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "apollo_001".to_string(), name: "Hot & Cold".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/hot-cold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_002".to_string(), name: "Neon Diamond".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/neon-diamond/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_003".to_string(), name: "Ruby Heart".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/ruby-heart/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_004".to_string(), name: "Thunder Zeus".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/thunder-zeus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_005".to_string(), name: "Super Rainbow".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/super-rainbow/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_006".to_string(), name: "Magic Cherry".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/magic-cherry/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_007".to_string(), name: "Lucky Dragon".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/lucky-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_008".to_string(), name: "Vegas Hot".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/vegas-hot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_009".to_string(), name: "Fruity 7".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/fruity-7/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "apollo_010".to_string(), name: "Speed Cash".to_string(), provider: "Apollo Games".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.apollo-games.com/speed-cash/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for ApolloGamesProvider {
    fn name(&self) -> &str { "Apollo Games" }
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
