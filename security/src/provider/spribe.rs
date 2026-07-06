//! Spribe Provider Integration
//! 
//! Integration with Spribe's game API.
//! Spribe is known for innovative crash games and lottery-style games.

use super::*;
use serde_json::json;

/// Spribe API client
pub struct SpribeProvider {
    config: ProviderConfig,
    base_url: String,
}

impl SpribeProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // Spribe portfolio - innovative crash and lottery games
        Ok(vec![
            // Aviator - their flagship game
            GameInfo {
                id: "spribe_aviator".to_string(),
                name: "Aviator".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Crash,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/aviator.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Mine games
            GameInfo {
                id: "spribe_mine".to_string(),
                name: "Mines".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/mines.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Hi-Lo
            GameInfo {
                id: "spribe_hilo".to_string(),
                name: "Hi-Lo".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 99.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/hilo.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Dice
            GameInfo {
                id: "spribe_dice".to_string(),
                name: "Dice".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 99.0,
                volatility: Volatility::Low,
                min_bet: 0.10,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/dice.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Keno
            GameInfo {
                id: "spribe_keno".to_string(),
                name: "Keno".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Lottery,
                rtp: 97.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/keno.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Plinko
            GameInfo {
                id: "spribe_plinko".to_string(),
                name: "Plinko".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 98.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/plinko.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Penalty
            GameInfo {
                id: "spribe_penalty".to_string(),
                name: "Penalty".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::Medium,
                min_bet: 0.50,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/penalty.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Cave
            GameInfo {
                id: "spribe_cave".to_string(),
                name: "Cave".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/cave.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Dino
            GameInfo {
                id: "spribe_dino".to_string(),
                name: "Dino".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/dino.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Wheel
            GameInfo {
                id: "spribe_wheel".to_string(),
                name: "Wheel".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::GameShows,
                rtp: 97.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/wheel.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Christmas
            GameInfo {
                id: "spribe_christmas".to_string(),
                name: "Christmas".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/christmas.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Heads & Tails
            GameInfo {
                id: "spribe_heads".to_string(),
                name: "Heads & Tails".to_string(),
                provider: "Spribe".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.spribe.co/heads-tails.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for SpribeProvider {
    fn name(&self) -> &str {
        "Spribe"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_games()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        Ok(LaunchGameResponse {
            game_url: format!("{}/game/{}", self.base_url, request.game_id),
            session_id: uuid::Uuid::new_v4().to_string(),
            token: "session_token".to_string(),
            expires_at: chrono::Utc::now().timestamp() + 3600,
        })
    }
    
    fn process_transaction(&self, request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        Ok(TransactionResult {
            transaction_id: uuid::Uuid::new_v4().to_string(),
            status: TransactionStatus::Completed,
            amount: request.amount,
            balance_after: 0.0,
            game_round_id: request.round_id,
            timestamp: chrono::Utc::now().timestamp(),
        })
    }
    
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.get_games()?;
        games.into_iter()
            .find(|g| g.id == game_id)
            .ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    
    fn is_available(&self) -> bool {
        self.config.enabled
    }
}
