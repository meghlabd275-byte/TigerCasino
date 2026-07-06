//! Hacksaw Gaming Provider Integration
//! 
//! Integration with Hacksaw Gaming's game API.
//! Hacksaw Gaming offers slots, scratch cards, and instant win games.

use super::*;
use reqwest::Client;
use serde_json::json;

/// Hacksaw Gaming API client
pub struct HacksawGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl HacksawGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            client: Client::new(),
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    /// Get available games from Hacksaw Gaming
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // In production, this would call the actual Hacksaw Gaming API
        // For now, return simulated game list
        Ok(vec![
            GameInfo {
                id: "hacksaw_slots_001".to_string(),
                name: "Stick 'Em".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.3,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/stick-em/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_slots_002".to_string(),
                name: "The Great Stickings".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.2,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/great-stickings/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_slots_003".to_string(),
                name: "Cash Quest".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.1,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/cash-quest/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_slots_004".to_string(),
                name: "Alchymedes".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: false,
                thumbnail_url: "https://static.hacksawgaming.com/alchymedes/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_slots_005".to_string(),
                name: "Egyptian King".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.1,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/egyptian-king/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_scratch_001".to_string(),
                name: "Scratch! Https://static.hacksawgaming.com".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::ScratchCards,
                rtp: 96.5,
                volatility: Volatility::Low,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.hacksawgaming.com/scratch/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_scratch_002".to_string(),
                name: "Cash Compass".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::ScratchCards,
                rtp: 96.2,
                volatility: Volatility::Medium,
                min_bet: 0.50,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://static.hacksawgaming.com/cash-compass/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_instant_001".to_string(),
                name: "Chaos Crew".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/chaos-crew/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_instant_002".to_string(),
                name: "Dead Man's Treasure".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.1,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/dead-treasure/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "hacksaw_instant_003".to_string(),
                name: "Fortune NN".to_string(),
                provider: "Hacksaw Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.2,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.hacksawgaming.com/fortune-nn/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for HacksawGamingProvider {
    fn name(&self) -> &str {
        "Hacksaw Gaming"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_games()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        // In production, this would call the Hacksaw Gaming launch API
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
