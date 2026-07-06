//! Nolimit City Provider Integration
//! 
//! Integration with Nolimit City's game API.
//! Nolimit City offers high-volatility slots with innovative features.

use super::*;
use reqwest::Client;
use chrono::Utc;

/// Nolimit City API client
pub struct NolimitCityProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl NolimitCityProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            client: Client::new(),
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    /// Get available games from Nolimit City
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // Simulated game list based on Nolimit City's portfolio
        Ok(vec![
            GameInfo {
                id: "nolimit_slots_001".to_string(),
                name: "San Quentin".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 70.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/san-quentin/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_002".to_string(),
                name: "Book of Shadows".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/book-shadows/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_003".to_string(),
                name: "Mental".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 40.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/mental/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_004".to_string(),
                name: "xWays Hoarder".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.1,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/xways-hoarder/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_005".to_string(),
                name: "Punk Rocker".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.02,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/punk-rocker/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_006".to_string(),
                name: "Deadwood".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/deadwood/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_007".to_string(),
                name: "Derrick".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/derrick/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_008".to_string(),
                name: "Incinerator".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/incinerator/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_009".to_string(),
                name: "Golden Genie".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/golden-genie/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "nolimit_slots_010".to_string(),
                name: "The Descent".to_string(),
                provider: "Nolimit City".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.nolimitcity.com/descent/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for NolimitCityProvider {
    fn name(&self) -> &str {
        "Nolimit City"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_games()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        Ok(LaunchGameResponse {
            game_url: format!("{}/game/{}", self.base_url, request.game_id),
            session_id: uuid::Uuid::new_v4().to_string(),
            token: "session_token".to_string(),
            expires_at: Utc::now().timestamp() + 3600,
        })
    }
    
    fn process_transaction(&self, _request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        Ok(TransactionResult {
            transaction_id: uuid::Uuid::new_v4().to_string(),
            status: TransactionStatus::Completed,
            amount: 0.0,
            balance_after: 0.0,
            game_round_id: "round_id".to_string(),
            timestamp: Utc::now().timestamp(),
        })
    }
    
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.fetch_games()?;
        games.into_iter()
            .find(|g| g.id == game_id)
            .ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    
    fn is_available(&self) -> bool {
        self.config.enabled
    }
}
