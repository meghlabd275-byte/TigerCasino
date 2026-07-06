//! Push Gaming Provider Integration
//! 
//! Integration with Push Gaming's game API.
//! Push Gaming offers high-quality slots with innovative mechanics.

use super::*;
use reqwest::Client;
use chrono::Utc;

/// Push Gaming API client
pub struct PushGamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PushGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            client: Client::new(),
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    /// Get available games from Push Gaming
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo {
                id: "push_slots_001".to_string(),
                name: "Jammin' Jars".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.83,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/jammin-jars/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_002".to_string(),
                name: "Jammin' Jars 2".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.60,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/jammin-jars-2/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_003".to_string(),
                name: "Dead Man's Trail".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.20,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/dead-mans-trail/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_004".to_string(),
                name: "Retro Tapes".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/retro-tapes/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_005".to_string(),
                name: "Big Bamboo".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/big-bamboo/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_006".to_string(),
                name: "Fire Hopper".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/fire-hopper/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_007".to_string(),
                name: "Mighty symbols: Diamonds".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/mighty-diamonds/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_008".to_string(),
                name: "Wild Swarm".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 97.00,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/wild-swarm/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_009".to_string(),
                name: "The Shadow".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/the-shadow/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "push_slots_010".to_string(),
                name: "Fat Rabbit".to_string(),
                provider: "Push Gaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.25,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://static.pushgaming.com/fat-rabbit/thumb.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for PushGamingProvider {
    fn name(&self) -> &str {
        "Push Gaming"
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
